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

        // For simplicity, we're using a basic approximation method
        // In a production environment, consider using a more sophisticated method or a library

        // Convert to RAY for higher precision in intermediate calculations
        uint256 rayUtilization = WadRayMath.wadToRay(utilizationRate);

        // Using natural logarithm approximation
        // ln(x) ≈ 2 * ((x - 1) / (x + 1)) for x close to 1
        // Then U^β = e^(β * ln(U))

        // Simple approximation for demo purposes
        // In practice, use a more accurate method or a library
        return WadRayMath.rayToWad(WadRayMath.rayPow(rayUtilization, elasticityFactor));
    }
}
