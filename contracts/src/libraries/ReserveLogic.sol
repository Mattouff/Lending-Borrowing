// SPDX-License-Identifier: MIT
pragma solidity 0.8.29;

import "./WadRayMath.sol";
import "./PercentageMath.sol";

/**
 * @title ReserveLogic library
 * @author DeFi Lending Platform
 * @notice Implements the logic to update reserves' state
 */
library ReserveLogic {
    using WadRayMath for uint256;
    using PercentageMath for uint256;

    // Constants for time calculations
    uint256 internal constant SECONDS_PER_YEAR = 365 days;

    /**
     * @dev Reserve state structure
     */
    struct ReserveData {
        // Principal state
        uint256 totalLiquidity; // Total liquidity in the reserve
        uint256 totalBorrows; // Total borrowed from the reserve
        // Interest rate state
        uint256 currentLiquidityRate; // Current interest rate for lenders
        uint256 currentBorrowRate; // Current interest rate for borrowers
        uint256 lastUpdateTimestamp; // Last time the reserve was updated
        // Reserve configuration
        address interestRateStrategyAddress; // The address of interest rate strategy
        uint256 baseLTVasCollateral; // The LTV used for borrowing
        uint256 liquidationThreshold; // Threshold for liquidation
        uint256 liquidationBonus; // Bonus for liquidators
        uint256 borrowingEnabled; // Flag for enabling/disabling borrowing
        // Addresses
        address aTokenAddress; // Address of the aToken representing deposits
        address debtTokenAddress; // Address of the debt token representing borrows
    }

    /**
     * @dev Updates the reserve state with accumulated interest
     * @param reserve The reserve to be updated
     */
    function updateReserveState(ReserveData storage reserve) internal {
        if (reserve.lastUpdateTimestamp == block.timestamp) {
            return; // Already updated in the same block
        }

        uint256 timeDelta = block.timestamp - reserve.lastUpdateTimestamp;

        if (timeDelta == 0) {
            return; // No time has passed
        }

        // Calculate accumulated interest
        uint256 borrowInterest = calculateCompoundedInterest(reserve.currentBorrowRate, timeDelta);

        // Update total borrows with accumulated interest
        uint256 previousBorrows = reserve.totalBorrows;

        if (previousBorrows > 0) {
            reserve.totalBorrows = previousBorrows.wadMul(borrowInterest);
        }

        // Update liquidity with accumulated interest
        reserve.totalLiquidity = reserve.totalLiquidity + (reserve.totalBorrows - previousBorrows);

        // Update interest rates
        updateInterestRates(reserve);

        // Update timestamp
        reserve.lastUpdateTimestamp = block.timestamp;
    }

    /**
     * @dev Updates the interest rates based on current utilization
     * @param reserve The reserve to update
     */
    function updateInterestRates(ReserveData storage reserve) internal {
        // Call external interest rate strategy contract
        // (Simplified for demonstration - in practice, use interface call)

        // Simple interest rate calculation based on utilization
        uint256 utilizationRate = calculateUtilizationRate(reserve.totalBorrows, reserve.totalLiquidity);

        (uint256 newLiquidityRate, uint256 newBorrowRate) =
            calculateInterestRates(utilizationRate, reserve.interestRateStrategyAddress);

        reserve.currentLiquidityRate = newLiquidityRate;
        reserve.currentBorrowRate = newBorrowRate;
    }

    /**
     * @dev Calculates the utilization rate of the reserve
     * @param totalBorrows Total borrowed from the reserve
     * @param totalLiquidity Total liquidity in the reserve
     * @return The utilization rate as a wad percentage
     */
    function calculateUtilizationRate(uint256 totalBorrows, uint256 totalLiquidity) internal pure returns (uint256) {
        if (totalLiquidity == 0) {
            return 0;
        }

        return totalBorrows.wadDiv(totalLiquidity);
    }

    /**
     * @dev Placeholder for interest rate calculation
     * In a real implementation, this would call the interest rate strategy contract
     */
    function calculateInterestRates(uint256 utilizationRate, address /*strategyAddress*/ )
        internal
        pure
        returns (uint256 liquidityRate, uint256 borrowRate)
    {
        // Simplified implementation for demonstration
        // In production, use a proper interest rate strategy contract

        uint256 baseRate = 0.01 * 1e18; // 1% base rate
        uint256 rSlope1 = 0.1 * 1e18; // 10% slope below optimal utilization
        uint256 rSlope2 = 0.4 * 1e18; // 40% slope above optimal utilization
        uint256 optimalUtilization = 0.8 * 1e18; // 80% optimal utilization

        if (utilizationRate <= optimalUtilization) {
            borrowRate = baseRate + (utilizationRate.wadDiv(optimalUtilization).wadMul(rSlope1));
        } else {
            borrowRate = baseRate + rSlope1
                + ((utilizationRate - optimalUtilization).wadDiv(WadRayMath.wad() - optimalUtilization).wadMul(rSlope2));
        }

        // Liquidity rate is the borrow rate multiplied by utilization rate
        liquidityRate = borrowRate.wadMul(utilizationRate);

        return (liquidityRate, borrowRate);
    }

    /**
     * @dev Calculates the compounded interest over a time period
     * @param rate The annual interest rate in ray
     * @param timeDelta The time period in seconds
     * @return The compounded interest factor
     */
    function calculateCompoundedInterest(uint256 rate, uint256 timeDelta) internal pure returns (uint256) {
        // Convert annual rate to per-second rate
        uint256 ratePerSecond = rate / SECONDS_PER_YEAR;

        // A = (1 + r)^t
        // Using approximation for small time periods
        return WadRayMath.wad() + (ratePerSecond * timeDelta);
    }

    /**
     * @dev Calculates the APY for lenders
     * @param rate The interest rate per second
     * @return The annual percentage yield
     */
    function calculateAPY(uint256 rate) internal pure returns (uint256) {
        // APY = (1 + r)^(365*24*60*60) - 1
        return calculateCompoundedInterest(rate, SECONDS_PER_YEAR) - WadRayMath.wad();
    }

    /**
     * @dev Calculates the maximum borrowable amount based on collateral
     * @param collateralValue The value of the collateral in ETH
     * @param minCollateralRatio The minimum collateralization ratio
     * @return The maximum borrowable amount
     */
    function calculateMaxBorrowAmount(uint256 collateralValue, uint256 minCollateralRatio)
        internal
        pure
        returns (uint256)
    {
        return collateralValue.wadDiv(minCollateralRatio);
    }

    /**
     * @dev Calculates the health factor for a user
     * @param totalCollateralETH The total collateral in ETH
     * @param totalDebtETH The total debt in ETH
     * @param liquidationThreshold The liquidation threshold
     * @return The health factor (>1 is healthy, <1 is liquidatable)
     */
    function calculateHealthFactor(uint256 totalCollateralETH, uint256 totalDebtETH, uint256 liquidationThreshold)
        internal
        pure
        returns (uint256)
    {
        if (totalDebtETH == 0) {
            return type(uint256).max;
        }

        return totalCollateralETH.percentMul(liquidationThreshold).wadDiv(totalDebtETH);
    }
}
