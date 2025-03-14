// SPDX-License-Identifier: MIT
pragma solidity 0.8.29;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "../interfaces/ILendingPool.sol";
import "../interfaces/ILendingPoolAddressesProvider.sol";
import "../interfaces/ILendingPoolCore.sol";
import "../interfaces/ICollateralManager.sol";
import "../interfaces/IPriceOracle.sol";
import "../interfaces/IAToken.sol";
import "../libraries/WadRayMath.sol";
import "../libraries/ReserveLogic.sol";

/**
 * @title LendingPoolDataProvider
 * @author DeFi Lending Platform
 * @notice Provides data about the lending pool state
 * @dev Contains view functions for frontends
 */
contract LendingPoolDataProvider {
    using WadRayMath for uint256;

    // Address provider
    ILendingPoolAddressesProvider private immutable _addressesProvider;

    // Constants
    uint256 public constant SECONDS_PER_YEAR = 365 days;

    /**
     * @dev Constructor
     * @param addressesProvider The address of the LendingPoolAddressesProvider
     */
    constructor(address addressesProvider) {
        require(addressesProvider != address(0), "LendingPoolDataProvider: Invalid addresses provider");
        _addressesProvider = ILendingPoolAddressesProvider(addressesProvider);
    }

    /**
     * @dev Gets the reserve data
     * @param asset The address of the asset
     * @return totalLiquidity The total liquidity
     * @return availableLiquidity The available liquidity
     * @return totalBorrows The total borrows
     * @return liquidityRate The liquidity rate
     * @return variableBorrowRate The variable borrow rate
     * @return utilizationRate The utilization rate
     * @return liquidityIndex The liquidity index
     * @return variableBorrowIndex The variable borrow index
     * @return ltv The loan to value ratio
     * @return liquidationThreshold The liquidation threshold
     * @return liquidationBonus The liquidation bonus
     * @return aTokenAddress The address of the aToken
     * @return lastUpdateTimestamp The last update timestamp
     */
    function getReserveData(address asset)
        external
        view
        returns (
            uint256 totalLiquidity,
            uint256 availableLiquidity,
            uint256 totalBorrows,
            uint256 liquidityRate,
            uint256 variableBorrowRate,
            uint256 utilizationRate,
            uint256 liquidityIndex,
            uint256 variableBorrowIndex,
            uint256 ltv,
            uint256 liquidationThreshold,
            uint256 liquidationBonus,
            address aTokenAddress,
            uint256 lastUpdateTimestamp
        )
    {
        address lendingPoolCore = _addressesProvider.getLendingPoolCore();

        // Get the reserve data from the lending pool core
        ReserveLogic.ReserveData memory reserveData = ILendingPoolCore(lendingPoolCore).getReserveData(asset);

        // Get the collateral configuration
        (ltv, liquidationThreshold, liquidationBonus) = _getReserveConfiguration(asset);

        // Get the available liquidity
        availableLiquidity = IERC20(asset).balanceOf(lendingPoolCore);

        return (
            reserveData.totalLiquidity,
            availableLiquidity,
            reserveData.totalBorrows,
            reserveData.currentLiquidityRate,
            reserveData.currentBorrowRate,
            _calculateUtilizationRate(reserveData.totalBorrows, reserveData.totalLiquidity),
            liquidityIndex,
            variableBorrowIndex,
            ltv,
            liquidationThreshold,
            liquidationBonus,
            reserveData.aTokenAddress,
            reserveData.lastUpdateTimestamp
        );
    }

    /**
     * @dev Gets the user reserve data
     * @param asset The address of the asset
     * @param user The address of the user
     * @return currentATokenBalance The current aToken balance
     * @return currentVariableDebt The current variable debt
     * @return liquidityRate The liquidity rate
     * @return variableBorrowRate The variable borrow rate
     * @return originationFee The origination fee
     * @return borrowCap The borrow cap
     * @return isUserUseReserveAsCollateral Whether the user is using the reserve as collateral
     */
    function getUserReserveData(address asset, address user)
        external
        view
        returns (
            uint256 currentATokenBalance,
            uint256 currentVariableDebt,
            uint256 liquidityRate,
            uint256 variableBorrowRate,
            uint256 originationFee,
            uint256 borrowCap,
            bool isUserUseReserveAsCollateral
        )
    {
        address lendingPoolCore = _addressesProvider.getLendingPoolCore();

        // Get the reserve data from the lending pool core
        ReserveLogic.ReserveData memory reserveData = ILendingPoolCore(lendingPoolCore).getReserveData(asset);

        // Get the aToken and debt token addresses
        address aTokenAddress = reserveData.aTokenAddress;
        address debtTokenAddress = reserveData.debtTokenAddress;

        // Get the user's aToken balance
        currentATokenBalance = IERC20(aTokenAddress).balanceOf(user);

        // Get the user's debt
        currentVariableDebt = IERC20(debtTokenAddress).balanceOf(user);

        // Get the interest rates
        liquidityRate = reserveData.currentLiquidityRate;
        variableBorrowRate = reserveData.currentBorrowRate;

        // Origination fee and borrow cap are hardcoded for simplicity
        originationFee = 0;
        borrowCap = type(uint256).max;

        // Whether the user is using the reserve as collateral
        isUserUseReserveAsCollateral = _isUserUsingAsCollateral(user, asset);

        return (
            currentATokenBalance,
            currentVariableDebt,
            liquidityRate,
            variableBorrowRate,
            originationFee,
            borrowCap,
            isUserUseReserveAsCollateral
        );
    }

    /**
     * @dev Gets the APY for a reserve
     * @param asset The address of the asset
     * @return The deposit APY and borrow APY
     */
    function getReserveAPY(address asset) external view returns (uint256, uint256) {
        address lendingPoolCore = _addressesProvider.getLendingPoolCore();

        // Get the reserve data from the lending pool core
        ReserveLogic.ReserveData memory reserveData = ILendingPoolCore(lendingPoolCore).getReserveData(asset);

        // Calculate the APY
        uint256 depositAPY = _calculateAPY(reserveData.currentLiquidityRate);
        uint256 borrowAPY = _calculateAPY(reserveData.currentBorrowRate);

        return (depositAPY, borrowAPY);
    }

    /**
     * @dev Gets the health factor for a user
     * @param user The address of the user
     * @return The health factor
     */
    function getUserHealthFactor(address user) external view returns (uint256) {
        address lendingPool = _addressesProvider.getLendingPool();

        (,,,,, uint256 healthFactor) = ILendingPool(lendingPool).getUserAccountData(user);

        return healthFactor;
    }

    /**
     * @dev Calculates the APY for a rate
     * @param rate The rate
     * @return The APY
     */
    function _calculateAPY(uint256 rate) internal pure returns (uint256) {
        // APY = (1 + r/n)^n - 1
        // Where r is the annual rate and n is the number of compounding periods per year
        // For simplicity, we use r as the annual rate and n=1

        return rate; // Simplified calculation
    }

    /**
     * @dev Calculates the utilization rate
     * @param totalBorrows The total borrows
     * @param totalLiquidity The total liquidity
     * @return The utilization rate
     */
    function _calculateUtilizationRate(uint256 totalBorrows, uint256 totalLiquidity) internal pure returns (uint256) {
        if (totalLiquidity == 0) {
            return 0;
        }

        return totalBorrows.wadDiv(totalLiquidity);
    }

    /**
     * @dev Gets the reserve configuration
     * @param asset The address of the asset
     * @return The LTV, liquidation threshold, and liquidation bonus
     */
    function _getReserveConfiguration(address asset) internal view returns (uint256, uint256, uint256) {
        address collateralManager = _getCollateralManager();

        try ICollateralManager(collateralManager).getReserveConfig(asset) returns (
            bool, /*isCollateral*/ uint256 ltv, uint256 liquidationThreshold, uint256 liquidationBonus
        ) {
            return (ltv, liquidationThreshold, liquidationBonus);
        } catch {
            // Return default values if the call fails
            return (75, 80, 110); // 75% LTV, 80% liquidation threshold, 110% liquidation bonus
        }
    }

    /**
     * @dev Checks if a user is using an asset as collateral
     * @param user The address of the user
     * @param asset The address of the asset
     * @return Whether the user is using the asset as collateral
     */
    function _isUserUsingAsCollateral(address user, address asset) internal view returns (bool) {
        // In a real implementation, this would check the user's configuration
        // For simplicity, we return true if the user has any aTokens
        address aTokenAddress = _getReserveATokenAddress(asset);

        return IERC20(aTokenAddress).balanceOf(user) > 0;
    }

    /**
     * @dev Get the aToken address for a reserve
     * @param asset The address of the asset
     * @return The address of the aToken
     */
    function _getReserveATokenAddress(address asset) internal view returns (address) {
        address lendingPoolCore = _addressesProvider.getLendingPoolCore();

        try ILendingPoolCore(lendingPoolCore).getReserveData(asset) returns (
            ReserveLogic.ReserveData memory reserveData
        ) {
            return reserveData.aTokenAddress;
        } catch {
            return address(0);
        }
    }

    /**
     * @dev Get the address of the collateral manager
     * @return The address of the collateral manager
     */
    function _getCollateralManager() internal view returns (address) {
        return _addressesProvider.getCollateralManager();
    }
}
