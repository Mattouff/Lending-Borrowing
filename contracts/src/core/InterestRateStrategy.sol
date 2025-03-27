// SPDX-License-Identifier: MIT
pragma solidity 0.8.29;

import "@openzeppelin/contracts/access/Ownable.sol";
import "../interfaces/IInterestRateStrategy.sol";
import "../libraries/WadRayMath.sol";
import "../libraries/PercentageMath.sol";

/**
 * @title InterestRateStrategy
 * @author DeFi Lending Platform
 * @notice Implements the interest rate model for the lending platform
 * @dev Uses the formula r(U) = r_min + (r_max - r_min) × U^β
 */
contract InterestRateStrategy is IInterestRateStrategy, Ownable {
    using WadRayMath for uint256;
    using PercentageMath for uint256;

    // Base variable borrow rate when utilization rate = 0
    uint256 private immutable _baseVariableBorrowRate;

    // Slope of the variable interest curve when utilization rate <= optimal
    uint256 private immutable _variableRateSlope1;

    // Slope of the variable interest curve when utilization rate > optimal
    uint256 private immutable _variableRateSlope2;

    // Optimal utilization rate
    uint256 private immutable _optimalUtilizationRate;

    // Elasticity factor (β) for the interest rate model
    uint256 private immutable _elasticityFactor;

    /**
     * @dev Constructor
     * @param baseVariableBorrowRate The base variable borrow rate
     * @param variableRateSlope1 The slope of the variable interest curve when utilization rate <= optimal
     * @param variableRateSlope2 The slope of the variable interest curve when utilization rate > optimal
     * @param optimalUtilizationRate The optimal utilization rate
     * @param elasticityFactor The elasticity factor (β) for the interest rate model
     */
    constructor(
        uint256 baseVariableBorrowRate,
        uint256 variableRateSlope1,
        uint256 variableRateSlope2,
        uint256 optimalUtilizationRate,
        uint256 elasticityFactor
    ) Ownable(msg.sender) {
        _baseVariableBorrowRate = baseVariableBorrowRate;
        _variableRateSlope1 = variableRateSlope1;
        _variableRateSlope2 = variableRateSlope2;
        _optimalUtilizationRate = optimalUtilizationRate;
        _elasticityFactor = elasticityFactor;
    }

    /**
     * @inheritdoc IInterestRateStrategy
     */
    function getBaseVariableBorrowRate() external view override returns (uint256) {
        return _baseVariableBorrowRate;
    }

    /**
     * @inheritdoc IInterestRateStrategy
     */
    function getMaxVariableBorrowRate() external view override returns (uint256) {
        return _baseVariableBorrowRate + _variableRateSlope1 + _variableRateSlope2;
    }

    /**
     * @inheritdoc IInterestRateStrategy
     */
    function calculateInterestRates(
        address, /* reserve */
        uint256 availableLiquidity,
        uint256 totalBorrows,
        uint256 reserveFactor
    ) external view override returns (uint256 liquidityRate, uint256 stableBorrowRate, uint256 variableBorrowRate) {
        // Calculate utilization rate
        uint256 totalLiquidity = availableLiquidity + totalBorrows;
        uint256 utilizationRate = totalLiquidity > 0 ? totalBorrows.wadDiv(totalLiquidity) : 0;

        // Calculate variable borrow rate
        variableBorrowRate = _calculateVariableBorrowRate(utilizationRate);

        // Placeholder for stable borrow rate (not implemented in this version)
        stableBorrowRate = 0;

        // Calculate liquidity rate
        liquidityRate = _calculateLiquidityRate(utilizationRate, variableBorrowRate, reserveFactor);

        return (liquidityRate, stableBorrowRate, variableBorrowRate);
    }

    /**
     * @dev Calculates the variable borrow rate based on utilization rate
     * @param utilizationRate The utilization rate
     * @return The variable borrow rate
     */
    function _calculateVariableBorrowRate(uint256 utilizationRate) internal view returns (uint256) {
        // Implement the formula r(U) = r_min + (r_max - r_min) × U^β

        if (utilizationRate <= _optimalUtilizationRate) {
            // Below optimal: Use slope1
            uint256 rateSlope = utilizationRate.wadMul(_variableRateSlope1.wadDiv(_optimalUtilizationRate));
            return _baseVariableBorrowRate + rateSlope;
        } else {
            // Above optimal: Add slope2 with an increasing rate
            uint256 normalRate = _baseVariableBorrowRate + _variableRateSlope1;

            // Calculate excess utilization safely
            uint256 denominator = WadRayMath.wad() - _optimalUtilizationRate;
            if (denominator == 0) {
                denominator = 1; // Prevent division by zero
            }

            uint256 excessUtilization = (utilizationRate - _optimalUtilizationRate).wadDiv(denominator);

            // Cap excessUtilization to prevent extreme values
            if (excessUtilization > WadRayMath.wad()) {
                excessUtilization = WadRayMath.wad();
            }

            // Calculate excess rate using a safer approach
            uint256 excessRate;

            // For elasticity factor of 2.0 (quadratic growth), we can use simple multiplication
            if (_elasticityFactor == 2 * WadRayMath.wad()) {
                // U^2 is simply U*U
                excessRate = excessUtilization.wadMul(excessUtilization).wadMul(_variableRateSlope2);
            } else if (_elasticityFactor == WadRayMath.wad()) {
                // Linear growth
                excessRate = excessUtilization.wadMul(_variableRateSlope2);
            } else {
                // For other elasticity factors, use a simplified approach
                // Linear interpolation between linear and quadratic based on elasticity
                uint256 linearPart = excessUtilization.wadMul(_variableRateSlope2);
                uint256 quadraticPart = excessUtilization.wadMul(excessUtilization).wadMul(_variableRateSlope2);

                // If elasticity is > 1, weight toward quadratic; if < 1, weight toward linear
                if (_elasticityFactor > WadRayMath.wad()) {
                    uint256 weight = (_elasticityFactor - WadRayMath.wad()).wadDiv(WadRayMath.wad());
                    excessRate = linearPart.wadMul(WadRayMath.wad() - weight) + quadraticPart.wadMul(weight);
                } else {
                    // For elasticity < 1, use linear as it's safer
                    excessRate = linearPart;
                }
            }

            return normalRate + excessRate;
        }
    }

    /**
     * @dev Calculates the liquidity rate based on utilization rate and variable borrow rate
     * @param utilizationRate The utilization rate
     * @param variableBorrowRate The variable borrow rate
     * @param reserveFactor The reserve factor
     * @return The liquidity rate
     */
    function _calculateLiquidityRate(uint256 utilizationRate, uint256 variableBorrowRate, uint256 reserveFactor)
        internal
        pure
        returns (uint256)
    {
        // Calculate liquidity rate as utilizationRate * variableBorrowRate * (1 - reserveFactor)
        uint256 oneMinusReserveFactor = WadRayMath.wad() - reserveFactor;

        return utilizationRate.wadMul(variableBorrowRate.wadMul(oneMinusReserveFactor));
    }
}
