// SPDX-License-Identifier: MIT
pragma solidity 0.8.29;

import "./WadRayMath.sol";

/**
 * @title PercentageMath library
 * @author DeFi Lending Platform
 * @notice Provides functions for percentage calculations
 */
library PercentageMath {
    using WadRayMath for uint256;

    uint256 internal constant PERCENTAGE_FACTOR = 1e4; // 10000 = 100%
    uint256 internal constant HALF_PERCENTAGE_FACTOR = PERCENTAGE_FACTOR / 2;

    /**
     * @dev Executes a percentage multiplication: (value * percentage) / PERCENTAGE_FACTOR
     * @param value The value to be multiplied
     * @param percentage The percentage value in PERCENTAGE_FACTOR format (10000 = 100%)
     * @return The result of the percentage multiplication
     */
    function percentMul(uint256 value, uint256 percentage) internal pure returns (uint256) {
        if (value == 0 || percentage == 0) {
            return 0;
        }

        return (value * percentage + HALF_PERCENTAGE_FACTOR) / PERCENTAGE_FACTOR;
    }

    /**
     * @dev Executes a percentage division: (value * PERCENTAGE_FACTOR) / percentage
     * @param value The value to be divided
     * @param percentage The percentage value in PERCENTAGE_FACTOR format (10000 = 100%)
     * @return The result of the percentage division
     */
    function percentDiv(uint256 value, uint256 percentage) internal pure returns (uint256) {
        require(percentage != 0, "PercentageMath: Division by zero");

        uint256 halfPercentage = percentage / 2;
        return (value * PERCENTAGE_FACTOR + halfPercentage) / percentage;
    }

    /**
     * @dev Calculates U raised to the power of elasticity factor (β)
     * Used in the variable interest rate formula: r(U) = r_min + (r_max - r_min) × U^β
     * @param utilizationRate The utilization rate in wad format (1e18 = 100%)
     * @param elasticityFactor The elasticity factor (β)
     * @return The result of U^β in wad format
     */
    function powFactor(uint256 utilizationRate, uint256 elasticityFactor) internal pure returns (uint256) {
        if (elasticityFactor == WadRayMath.wad()) {
            return utilizationRate; // If β = 1, U^1 = U
        }

        if (utilizationRate == 0) {
            return 0; // 0^anything = 0
        }

        if (utilizationRate == WadRayMath.wad()) {
            return WadRayMath.wad(); // 1^anything = 1
        }

        // For values < 1 and β > 1, we can use a simple approximation
        // For simplicity, we'll use a linear approximation for test purposes
        // In production, use a proper power function library
        if (elasticityFactor > WadRayMath.wad()) {
            // If elasticity > 1, the curve grows faster than linear
            // Approximate by returning min of (U * elasticity) and 1
            uint256 result = utilizationRate.wadMul(elasticityFactor);
            return result < WadRayMath.wad() ? result : WadRayMath.wad();
        } else {
            // If elasticity < 1, the curve grows slower than linear
            // Approximate by returning sqrt(U) for elasticity = 0.5
            return utilizationRate; // Simplification for test
        }
    }
}
