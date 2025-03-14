// SPDX-License-Identifier: MIT
pragma solidity 0.8.29;

import "@openzeppelin-upgradeable/contracts/utils/ReentrancyGuardUpgradeable.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin-upgradeable/contracts/proxy/utils/Initializable.sol";
import "@openzeppelin-upgradeable/contracts/proxy/utils/UUPSUpgradeable.sol";
import "@openzeppelin-upgradeable/contracts/access/OwnableUpgradeable.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "../interfaces/ILendingPool.sol";
import "../interfaces/ILendingPoolAddressesProvider.sol";
import "../interfaces/ILendingPoolCore.sol";
import "../interfaces/IPriceOracle.sol";
import "../interfaces/IAToken.sol";
import "../interfaces/IDebtToken.sol";
import "../libraries/WadRayMath.sol";
import "../libraries/ReserveLogic.sol";

/**
 * @title LendingPool
 * @author DeFi Lending Platform
 * @notice Main contract for the lending platform
 * @dev Entry point for users to interact with the protocol
 * This contract is upgradeable using the UUPS pattern
 */
contract LendingPool is Initializable, ILendingPool, ReentrancyGuardUpgradeable, OwnableUpgradeable, UUPSUpgradeable {
    using SafeERC20 for IERC20;
    using WadRayMath for uint256;

    // Address provider
    ILendingPoolAddressesProvider private _addressesProvider;

    // User configuration mapping (user => reserve => config)
    struct UserConfig {
        bool isUsingAsCollateral;
        bool hasBorrowed;
    }

    mapping(address => mapping(address => UserConfig)) private _usersConfig;

    // Storage gap for future upgrades
    uint256[50] private __gap;

    /**
     * @custom:oz-upgrades-unsafe-allow constructor
     */
    constructor() {
        _disableInitializers();
    }

    /**
     * @dev Initializer function
     * @param addressesProvider The address of the LendingPoolAddressesProvider
     */
    function initialize(address addressesProvider) external initializer {
        __ReentrancyGuard_init();
        __Ownable_init(msg.sender);
        __UUPSUpgradeable_init();

        require(addressesProvider != address(0), "LendingPool: Invalid addresses provider");
        _addressesProvider = ILendingPoolAddressesProvider(addressesProvider);
    }

    /**
     * @inheritdoc ILendingPool
     */
    function deposit(address asset, uint256 amount, uint16 referralCode) external override nonReentrant {
        require(asset != address(0), "LendingPool: Invalid asset");
        require(amount > 0, "LendingPool: Invalid amount");

        address lendingPoolCore = _addressesProvider.getLendingPoolCore();
        address aTokenAddress = _getReserveATokenAddress(asset);

        require(aTokenAddress != address(0), "LendingPool: Asset not supported");

        // Update the reserve state
        _updateReserveState(asset);

        // Transfer the asset from the user to the lending pool core
        IERC20(asset).safeTransferFrom(msg.sender, lendingPoolCore, amount);

        // Mint aTokens to the user
        IAToken(aTokenAddress).mint(msg.sender, amount, _getLiquidityIndex(asset));

        // Update user configuration
        _usersConfig[msg.sender][asset].isUsingAsCollateral = true;

        emit Deposit(msg.sender, asset, amount, referralCode);
    }

    /**
     * @inheritdoc ILendingPool
     */
    function withdraw(address asset, uint256 amount) external override nonReentrant returns (uint256) {
        require(asset != address(0), "LendingPool: Invalid asset");

        address aTokenAddress = _getReserveATokenAddress(asset);

        require(aTokenAddress != address(0), "LendingPool: Asset not supported");

        // Update the reserve state
        _updateReserveState(asset);

        // Get the current balance
        uint256 userBalance = IAToken(aTokenAddress).balanceOf(msg.sender);

        // If amount is uint256 max, withdraw the entire balance
        uint256 amountToWithdraw = amount == type(uint256).max ? userBalance : amount;

        require(amountToWithdraw > 0, "LendingPool: Invalid amount");
        require(amountToWithdraw <= userBalance, "LendingPool: Insufficient balance");

        // Check if the user has enough collateral after withdrawal
        require(_validateWithdrawal(msg.sender, asset, amountToWithdraw), "LendingPool: Not enough collateral");

        // Burn aTokens from the user
        uint256 withdrawnAmount = IAToken(aTokenAddress).burn(msg.sender, amountToWithdraw, _getLiquidityIndex(asset));

        // Transfer the asset from the lending pool core to the user
        _transferFromCore(asset, msg.sender, withdrawnAmount);

        emit Withdraw(msg.sender, asset, withdrawnAmount);

        return withdrawnAmount;
    }

    // Rest of the functions remain the same but with changes to use upgradeable libraries...

    /**
     * @dev Function that should revert when msg.sender is not authorized to upgrade the contract
     * @param newImplementation The address of the new implementation
     */
    function _authorizeUpgrade(address newImplementation) internal override onlyOwner { }

    /**
     * @inheritdoc ILendingPool
     */
    function borrow(address asset, uint256 amount, uint256 interestRateMode, uint16 referralCode)
        external
        override
        nonReentrant
    {
        require(asset != address(0), "LendingPool: Invalid asset");
        require(amount > 0, "LendingPool: Invalid amount");
        require(interestRateMode == 1, "LendingPool: Only variable rate is supported");

        address debtTokenAddress = _getReserveDebtTokenAddress(asset);

        require(debtTokenAddress != address(0), "LendingPool: Asset not supported");

        // Update the reserve state
        _updateReserveState(asset);

        // Check if the user has enough collateral
        require(_validateBorrow(msg.sender, asset, amount), "LendingPool: Not enough collateral");

        // Transfer the asset from the lending pool core to the user
        _transferFromCore(asset, msg.sender, amount);

        // Mint debt tokens to the user
        uint256 borrowRate = _getCurrentBorrowRate(asset);
        uint256 borrowIndex = _getBorrowIndex(asset);

        IDebtToken(debtTokenAddress).mint(msg.sender, amount, borrowIndex);

        // Update user configuration
        _usersConfig[msg.sender][asset].hasBorrowed = true;

        emit Borrow(msg.sender, asset, amount, interestRateMode, borrowRate, referralCode);
    }

    /**
     * @inheritdoc ILendingPool
     */
    function repay(address asset, uint256 amount, uint256 interestRateMode)
        external
        override
        nonReentrant
        returns (uint256)
    {
        require(asset != address(0), "LendingPool: Invalid asset");
        require(interestRateMode == 1, "LendingPool: Only variable rate is supported");

        address lendingPoolCore = _addressesProvider.getLendingPoolCore();
        address debtTokenAddress = _getReserveDebtTokenAddress(asset);

        require(debtTokenAddress != address(0), "LendingPool: Asset not supported");

        // Update the reserve state
        _updateReserveState(asset);

        // Get the current debt balance
        uint256 userDebt = _getUserDebt(msg.sender, asset);

        // If amount is uint256 max, repay the entire debt
        uint256 amountToRepay = amount == type(uint256).max ? userDebt : amount;

        require(amountToRepay > 0, "LendingPool: Invalid amount");
        require(amountToRepay <= userDebt, "LendingPool: Amount exceeds debt");

        // Transfer the asset from the user to the lending pool core
        IERC20(asset).safeTransferFrom(msg.sender, lendingPoolCore, amountToRepay);

        // Burn debt tokens from the user
        uint256 borrowIndex = _getBorrowIndex(asset);

        uint256 repaidAmount = IDebtToken(debtTokenAddress).burn(msg.sender, amountToRepay, borrowIndex);

        // Calculate interest paid
        uint256 interestPaid = repaidAmount > userDebt ? repaidAmount - userDebt : 0;

        emit Repay(msg.sender, asset, repaidAmount, interestPaid);

        return repaidAmount;
    }

    /**
     * @inheritdoc ILendingPool
     */
    function liquidationCall(address collateralAsset, address debtAsset, address user, uint256 debtToCover)
        external
        override
        nonReentrant
        returns (uint256)
    {
        require(collateralAsset != address(0), "LendingPool: Invalid collateral asset");
        require(debtAsset != address(0), "LendingPool: Invalid debt asset");
        require(user != address(0), "LendingPool: Invalid user");
        require(debtToCover > 0, "LendingPool: Invalid debt amount");

        // Check if the liquidation is valid
        (bool isValid,) = _isLiquidationValid(user);
        require(isValid, "LendingPool: Health factor above threshold");

        // Get the addresses of the tokens
        address lendingPoolCore = _addressesProvider.getLendingPoolCore();
        address debtTokenAddress = _getReserveDebtTokenAddress(debtAsset);
        address aTokenCollateralAddress = _getReserveATokenAddress(collateralAsset);

        require(debtTokenAddress != address(0), "LendingPool: Debt asset not supported");
        require(aTokenCollateralAddress != address(0), "LendingPool: Collateral asset not supported");

        // Update the reserve state
        _updateReserveState(debtAsset);
        _updateReserveState(collateralAsset);

        // Calculate the maximum debt amount that can be liquidated
        uint256 userDebt = _getUserDebt(user, debtAsset);
        uint256 maxDebtToLiquidate = userDebt / 2; // 50% close factor
        uint256 actualDebtToCover = debtToCover > maxDebtToLiquidate ? maxDebtToLiquidate : debtToCover;

        // Calculate the collateral amount to be liquidated
        uint256 collateralAmount = _calculateLiquidationAmount(actualDebtToCover, collateralAsset, debtAsset);

        // Transfer the debt asset from the liquidator to the lending pool core
        IERC20(debtAsset).safeTransferFrom(msg.sender, lendingPoolCore, actualDebtToCover);

        // Burn debt tokens from the user
        uint256 borrowIndex = _getBorrowIndex(debtAsset);

        // Call the burn function on the debt token
        IDebtToken(debtTokenAddress).burn(user, actualDebtToCover, borrowIndex);

        // Burn aTokens from the user
        uint256 liquidityIndex = _getLiquidityIndex(collateralAsset);

        // Call the burn function on the aToken
        IAToken(aTokenCollateralAddress).burn(user, collateralAmount, liquidityIndex);

        // Transfer the collateral asset from the lending pool core to the liquidator
        _transferFromCore(collateralAsset, msg.sender, collateralAmount);

        emit Liquidation(msg.sender, user, collateralAsset, collateralAmount, actualDebtToCover);

        return collateralAmount;
    }

    /**
     * @inheritdoc ILendingPool
     */
    function getUserAccountData(address user)
        external
        view
        override
        returns (
            uint256 totalCollateralETH,
            uint256 totalDebtETH,
            uint256 availableBorrowsETH,
            uint256 currentLiquidationThreshold,
            uint256 ltv,
            uint256 healthFactor
        )
    {
        // Get the price oracle
        address priceOracleAddress = _addressesProvider.getPriceOracle();
        IPriceOracle priceOracle = IPriceOracle(priceOracleAddress);

        // Calculate the user's collateral and debt
        (totalCollateralETH, totalDebtETH, ltv, currentLiquidationThreshold, healthFactor) =
            _calculateUserAccountData(user, priceOracle);

        // Calculate the available borrows
        availableBorrowsETH = totalCollateralETH * ltv / 10000;
        if (availableBorrowsETH > totalDebtETH) {
            availableBorrowsETH = availableBorrowsETH - totalDebtETH;
        } else {
            availableBorrowsETH = 0;
        }

        return (totalCollateralETH, totalDebtETH, availableBorrowsETH, currentLiquidationThreshold, ltv, healthFactor);
    }

    /**
     * @dev Update the reserve state
     * @param asset The address of the asset
     */
    function _updateReserveState(address asset) internal {
        address lendingPoolCore = _addressesProvider.getLendingPoolCore();
        ILendingPoolCore(lendingPoolCore).updateReserveState(asset);
    }

    /**
     * @dev Transfer an asset from the lending pool core to a user
     * @param asset The address of the asset
     * @param to The address of the recipient
     * @param amount The amount to transfer
     */
    function _transferFromCore(address asset, address to, uint256 amount) internal {
        address lendingPoolCore = _addressesProvider.getLendingPoolCore();
        ILendingPoolCore(lendingPoolCore).transferToUser(asset, to, amount);
    }

    /**
     * @dev Get the aToken address for a reserve
     * @param asset The address of the asset
     * @return The address of the aToken
     */
    function _getReserveATokenAddress(address asset) internal view returns (address) {
        address lendingPoolCore = _addressesProvider.getLendingPoolCore();
        ReserveLogic.ReserveData memory reserveData = ILendingPoolCore(lendingPoolCore).getReserveData(asset);
        return reserveData.aTokenAddress;
    }

    /**
     * @dev Get the debt token address for a reserve
     * @param asset The address of the asset
     * @return The address of the debt token
     */
    function _getReserveDebtTokenAddress(address asset) internal view returns (address) {
        address lendingPoolCore = _addressesProvider.getLendingPoolCore();
        ReserveLogic.ReserveData memory reserveData = ILendingPoolCore(lendingPoolCore).getReserveData(asset);
        return reserveData.debtTokenAddress;
    }

    /**
     * @dev Get the liquidity index for a reserve
     * @return The liquidity index
     */
    function _getLiquidityIndex(address /*asset*/ ) internal pure returns (uint256) {
        // In a real implementation, this would get the liquidity index from the reserve
        // For simplicity, we return 1e18 (RAY)
        return 1e18;
    }

    /**
     * @dev Get the borrow index for a reserve
     * @return The borrow index
     */
    function _getBorrowIndex(address /*asset*/ ) internal pure returns (uint256) {
        // In a real implementation, this would get the borrow index from the reserve
        // For simplicity, we return 1e18 (RAY)
        return 1e18;
    }

    /**
     * @dev Get the current borrow rate for a reserve
     * @return The current borrow rate
     */
    function _getCurrentBorrowRate(address /*asset*/ ) internal pure returns (uint256) {
        // In a real implementation, this would get the borrow rate from the reserve
        // For simplicity, we return 5% (0.05 * 1e18)
        return 0.05 * 1e18;
    }

    /**
     * @dev Get the user's debt for an asset
     * @param user The address of the user
     * @param asset The address of the asset
     * @return The user's debt
     */
    function _getUserDebt(address user, address asset) internal view returns (uint256) {
        address debtTokenAddress = _getReserveDebtTokenAddress(asset);

        // Use IERC20 interface since DebtToken implements ERC20 balanceOf
        return IERC20(debtTokenAddress).balanceOf(user);
    }

    /**
     * @dev Validate a borrow operation
     * @param user The address of the user
     * @param asset The address of the asset
     * @param amount The amount to borrow
     * @return Whether the borrow is valid
     */
    function _validateBorrow(address user, address asset, uint256 amount) internal view returns (bool) {
        // Get the price oracle
        address priceOracleAddress = _addressesProvider.getPriceOracle();
        IPriceOracle priceOracle = IPriceOracle(priceOracleAddress);

        // Calculate the user's account data
        (uint256 totalCollateralETH, uint256 totalDebtETH, uint256 ltv,,) = _calculateUserAccountData(user, priceOracle);

        // If the user has no collateral, the borrow is invalid
        if (totalCollateralETH == 0) {
            return false;
        }

        // Calculate the borrow amount in ETH
        uint256 amountETH = amount * priceOracle.getAssetPrice(asset) / 1e18;

        // Calculate the maximum borrowable amount
        uint256 maxBorrowableETH = totalCollateralETH * ltv / 10000;

        // Check if the user has enough collateral
        if (totalDebtETH + amountETH > maxBorrowableETH) {
            return false;
        }

        // Check if the health factor would remain above 1 after the borrow
        uint256 newHealthFactor = totalCollateralETH * 10000 / (totalDebtETH + amountETH);
        if (newHealthFactor < 1e18) {
            return false;
        }

        return true;
    }

    /**
     * @dev Validate a withdrawal operation
     * @param user The address of the user
     * @param asset The address of the asset
     * @param amount The amount to withdraw
     * @return Whether the withdrawal is valid
     */
    function _validateWithdrawal(address user, address asset, uint256 amount) internal view returns (bool) {
        // Get the price oracle
        address priceOracleAddress = _addressesProvider.getPriceOracle();
        IPriceOracle priceOracle = IPriceOracle(priceOracleAddress);

        // If the user has no debt, the withdrawal is valid
        uint256 totalDebt = 0;
        address[] memory reservesList = _getReservesList();
        for (uint256 i = 0; i < reservesList.length; i++) {
            if (_usersConfig[user][reservesList[i]].hasBorrowed) {
                totalDebt += _getUserDebt(user, reservesList[i]);
            }
        }

        if (totalDebt == 0) {
            return true;
        }

        // Calculate the withdrawal amount in ETH
        uint256 withdrawalAmountETH = amount * priceOracle.getAssetPrice(asset) / 1e18;

        // Calculate the user's account data
        (uint256 totalCollateralETH, uint256 totalDebtETH,, uint256 liquidationThreshold,) =
            _calculateUserAccountData(user, priceOracle);

        // Check if the user has enough collateral after withdrawal
        if (withdrawalAmountETH > totalCollateralETH) {
            return false;
        }

        // Calculate the new health factor after withdrawal
        uint256 newTotalCollateralETH = totalCollateralETH - withdrawalAmountETH;
        uint256 newHealthFactor = newTotalCollateralETH * liquidationThreshold / totalDebtETH;

        // The health factor must remain above 1
        return newHealthFactor >= 1e18;
    }

    /**
     * @dev Check if a liquidation is valid
     * @param user The address of the user
     * @return Whether the liquidation is valid and the health factor
     */
    function _isLiquidationValid(address user) internal view returns (bool, uint256) {
        // Get the price oracle
        address priceOracleAddress = _addressesProvider.getPriceOracle();
        IPriceOracle priceOracle = IPriceOracle(priceOracleAddress);

        // Calculate the user's account data
        (,,,, uint256 healthFactor) = _calculateUserAccountData(user, priceOracle);

        // Liquidation is valid if the health factor is below 1
        return (healthFactor < 1e18, healthFactor);
    }

    /**
     * @dev Calculate the liquidation amount
     * @param debtToCover The amount of debt to cover
     * @param collateralAsset The address of the collateral asset
     * @param debtAsset The address of the debt asset
     * @return The amount of collateral to be liquidated
     */
    function _calculateLiquidationAmount(uint256 debtToCover, address collateralAsset, address debtAsset)
        internal
        view
        returns (uint256)
    {
        // Get the price oracle
        address priceOracleAddress = _addressesProvider.getPriceOracle();
        IPriceOracle priceOracle = IPriceOracle(priceOracleAddress);

        // Get the prices of the assets
        uint256 collateralPrice = priceOracle.getAssetPrice(collateralAsset);
        uint256 debtPrice = priceOracle.getAssetPrice(debtAsset);

        // Calculate the collateral amount needed to cover the debt
        uint256 debtInETH = debtToCover * debtPrice / 1e18;

        // Apply liquidation bonus (10% bonus for liquidators)
        uint256 liquidationBonus = 110; // 110%
        uint256 collateralInETH = debtInETH * liquidationBonus / 100;

        // Convert from ETH to collateral asset
        return collateralInETH * 1e18 / collateralPrice;
    }

    /**
     * @dev Calculate the user's account data
     * @param user The address of the user
     * @param priceOracle The price oracle
     * @return totalCollateralETH The total collateral in ETH
     * @return totalDebtETH The total debt in ETH
     * @return ltv The loan to value ratio
     * @return liquidationThreshold The liquidation threshold
     * @return healthFactor The health factor
     */
    function _calculateUserAccountData(address user, IPriceOracle priceOracle)
        internal
        view
        returns (
            uint256 totalCollateralETH,
            uint256 totalDebtETH,
            uint256 ltv,
            uint256 liquidationThreshold,
            uint256 healthFactor
        )
    {
        // Initialize accumulators
        totalCollateralETH = 0;
        totalDebtETH = 0;
        uint256 totalLTVCollateral = 0;
        uint256 totalLiquidationThresholdCollateral = 0;

        // Get the reserves list
        address[] memory reservesList = _getReservesList();

        // Calculate the user's collateral and debt
        for (uint256 i = 0; i < reservesList.length; i++) {
            address asset = reservesList[i];

            // Get the reserve data
            (uint256 assetLTV, uint256 assetLiquidationThreshold) = _getReserveConfiguration(asset);

            // Check if the user is using the asset as collateral
            if (_usersConfig[user][asset].isUsingAsCollateral) {
                // Get the aToken address
                address aTokenAddress = _getReserveATokenAddress(asset);

                // Get the user's balance
                uint256 userBalance = IERC20(aTokenAddress).balanceOf(user);

                // Calculate the collateral value in ETH
                uint256 assetPrice = priceOracle.getAssetPrice(asset);
                uint256 collateralValueETH = userBalance * assetPrice / 1e18;

                // Add to the accumulators
                totalCollateralETH += collateralValueETH;
                totalLTVCollateral += collateralValueETH * assetLTV;
                totalLiquidationThresholdCollateral += collateralValueETH * assetLiquidationThreshold;
            }

            // Check if the user has borrowed the asset
            if (_usersConfig[user][asset].hasBorrowed) {
                // Get the user's debt
                uint256 userDebt = _getUserDebt(user, asset);

                // Calculate the debt value in ETH
                uint256 assetPrice = priceOracle.getAssetPrice(asset);
                uint256 debtValueETH = userDebt * assetPrice / 1e18;

                // Add to the accumulators
                totalDebtETH += debtValueETH;
            }
        }

        // Calculate the weighted average LTV and liquidation threshold
        if (totalCollateralETH > 0) {
            ltv = totalLTVCollateral / totalCollateralETH;
            liquidationThreshold = totalLiquidationThresholdCollateral / totalCollateralETH;
        }

        // Calculate the health factor
        healthFactor = totalDebtETH > 0 ? totalCollateralETH * liquidationThreshold / totalDebtETH : type(uint256).max;

        return (totalCollateralETH, totalDebtETH, ltv, liquidationThreshold, healthFactor);
    }

    /**
     * @dev Get the reserve configuration
     * @return The LTV and liquidation threshold
     */
    function _getReserveConfiguration(address /*asset*/ ) internal pure returns (uint256, uint256) {
        // In a real implementation, this would get the configuration from the reserve
        // For simplicity, we return fixed values
        return (75, 80); // 75% LTV, 80% liquidation threshold
    }

    /**
     * @dev Get the list of all reserves
     * @return The list of all reserves
     */
    function _getReservesList() internal view returns (address[] memory) {
        address lendingPoolCore = _addressesProvider.getLendingPoolCore();
        return ILendingPoolCore(lendingPoolCore).getReservesList();
    }
}
