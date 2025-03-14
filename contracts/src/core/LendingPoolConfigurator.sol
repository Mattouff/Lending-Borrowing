// SPDX-License-Identifier: MIT
pragma solidity 0.8.29;

import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "../interfaces/ILendingPoolAddressesProvider.sol";
import "../interfaces/ILendingPoolCore.sol";
import "../interfaces/ICollateralManager.sol";
import "../tokens/AToken.sol";
import "../tokens/DebtToken.sol";

/**
 * @title LendingPoolConfigurator
 * @author DeFi Lending Platform
 * @notice Handles the configuration of the lending pool
 * @dev Contains administrative functions to manage reserves
 */
contract LendingPoolConfigurator is Ownable {
    // Address provider
    ILendingPoolAddressesProvider private immutable _addressesProvider;

    // Events
    event ReserveInitialized(
        address indexed asset, address indexed aToken, address indexed debtToken, address interestRateStrategy
    );

    event BorrowingEnabledOnReserve(address indexed asset, bool enabled);

    event CollateralConfigurationChanged(
        address indexed asset, uint256 ltv, uint256 liquidationThreshold, uint256 liquidationBonus
    );

    /**
     * @dev Constructor
     * @param addressesProvider The address of the LendingPoolAddressesProvider
     */
    constructor(address addressesProvider) Ownable(msg.sender) {
        require(addressesProvider != address(0), "LendingPoolConfigurator: Invalid addresses provider");
        _addressesProvider = ILendingPoolAddressesProvider(addressesProvider);
    }

    /**
     * @dev Initializes a reserve
     * @param asset The address of the asset
     * @param aTokenName The name of the aToken
     * @param aTokenSymbol The symbol of the aToken
     * @param debtTokenName The name of the debt token
     * @param debtTokenSymbol The symbol of the debt token
     * @param interestRateStrategy The address of the interest rate strategy
     */
    function initReserve(
        address asset,
        string memory aTokenName,
        string memory aTokenSymbol,
        string memory debtTokenName,
        string memory debtTokenSymbol,
        address interestRateStrategy
    ) external onlyOwner {
        require(asset != address(0), "LendingPoolConfigurator: Invalid asset");
        require(interestRateStrategy != address(0), "LendingPoolConfigurator: Invalid interest rate strategy");

        // Get the lending pool and lending pool core addresses
        address lendingPool = _addressesProvider.getLendingPool();
        address lendingPoolCore = _addressesProvider.getLendingPoolCore();

        // Deploy aToken
        AToken aToken = new AToken(asset, lendingPool, aTokenName, aTokenSymbol);

        // Deploy debt token
        DebtToken debtToken = new DebtToken(asset, lendingPool, debtTokenName, debtTokenSymbol);

        // Initialize the reserve in the lending pool core
        // Call the initReserve function on the lending pool core
        ILendingPoolCore(lendingPoolCore).initReserve(asset, address(aToken), address(debtToken), interestRateStrategy);

        emit ReserveInitialized(asset, address(aToken), address(debtToken), interestRateStrategy);
    }

    /**
     * @dev Enables or disables borrowing on a reserve
     * @param asset The address of the asset
     * @param enabled Whether borrowing is enabled
     */
    function enableBorrowingOnReserve(address asset, bool enabled) external onlyOwner {
        require(asset != address(0), "LendingPoolConfigurator: Invalid asset");

        address lendingPoolCore = _addressesProvider.getLendingPoolCore();
        ILendingPoolCore(lendingPoolCore).enableBorrowingOnReserve(asset, enabled);

        emit BorrowingEnabledOnReserve(asset, enabled);
    }

    /**
     * @dev Configures the reserve as collateral
     * @param asset The address of the asset
     * @param ltv The loan to value ratio
     * @param liquidationThreshold The liquidation threshold
     * @param liquidationBonus The liquidation bonus
     */
    function configureReserveAsCollateral(
        address asset,
        uint256 ltv,
        uint256 liquidationThreshold,
        uint256 liquidationBonus
    ) external onlyOwner {
        require(asset != address(0), "LendingPoolConfigurator: Invalid asset");
        require(ltv <= 10000, "LendingPoolConfigurator: Invalid LTV");
        require(liquidationThreshold <= 10000, "LendingPoolConfigurator: Invalid liquidation threshold");
        require(liquidationBonus >= 10000, "LendingPoolConfigurator: Invalid liquidation bonus");
        require(liquidationThreshold >= ltv, "LendingPoolConfigurator: Liquidation threshold must be >= LTV");

        address collateralManager = _getCollateralManager();
        ICollateralManager(collateralManager).configureAsCollateral(
            asset, true, ltv, liquidationThreshold, liquidationBonus
        );

        emit CollateralConfigurationChanged(asset, ltv, liquidationThreshold, liquidationBonus);
    }

    /**
     * @dev Updates the interest rate strategy of a reserve
     * @param asset The address of the asset
     * @param interestRateStrategy The address of the interest rate strategy
     */
    function updateInterestRateStrategy(address asset, address interestRateStrategy) external onlyOwner {
        require(asset != address(0), "LendingPoolConfigurator: Invalid asset");
        require(interestRateStrategy != address(0), "LendingPoolConfigurator: Invalid interest rate strategy");

        address lendingPoolCore = _addressesProvider.getLendingPoolCore();
        ILendingPoolCore(lendingPoolCore).updateInterestRateStrategy(asset, interestRateStrategy);
    }

    /**
     * @dev Get the address of the collateral manager
     * @return The address of the collateral manager
     */
    function _getCollateralManager() internal view returns (address) {
        return _addressesProvider.getCollateralManager();
    }
}
