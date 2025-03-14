// SPDX-License-Identifier: MIT
pragma solidity 0.8.29;

/**
 * @title WadRayMath library
 * @author DeFi Lending Platform
 * @notice Provides functions for fixed point math calculations with Wad (18 decimals) and Ray (27 decimals)
 */
library WadRayMath {
    // WAD and RAY precision units
    uint256 internal constant WAD = 1e18;
    uint256 internal constant RAY = 1e27;

    uint256 internal constant HALF_WAD = WAD / 2;
    uint256 internal constant HALF_RAY = RAY / 2;

    /**
     * @dev Multiplies two wad numbers and returns a wad (decimal 18) fixed point number
     * @param a First wad number
     * @param b Second wad number
     * @return Result of a*b in wad units
     */
    function wadMul(uint256 a, uint256 b) internal pure returns (uint256) {
        if (a == 0 || b == 0) {
            return 0;
        }

        return (a * b + HALF_WAD) / WAD;
    }

    /**
     * @dev Divides two wad numbers and returns a wad (decimal 18) fixed point number
     * @param a Wad number (dividend)
     * @param b Wad number (divisor)
     * @return Result of a/b in wad units
     */
    function wadDiv(uint256 a, uint256 b) internal pure returns (uint256) {
        require(b != 0, "WadRayMath: Division by zero");

        uint256 halfB = b / 2;
        return (a * WAD + halfB) / b;
    }

    /**
     * @dev Multiplies two ray numbers and returns a ray (decimal 27) fixed point number
     * @param a First ray number
     * @param b Second ray number
     * @return Result of a*b in ray units
     */
    function rayMul(uint256 a, uint256 b) internal pure returns (uint256) {
        if (a == 0 || b == 0) {
            return 0;
        }

        return (a * b + HALF_RAY) / RAY;
    }

    /**
     * @dev Divides two ray numbers and returns a ray (decimal 27) fixed point number
     * @param a Ray number (dividend)
     * @param b Ray number (divisor)
     * @return Result of a/b in ray units
     */
    function rayDiv(uint256 a, uint256 b) internal pure returns (uint256) {
        require(b != 0, "WadRayMath: Division by zero");

        uint256 halfB = b / 2;
        return (a * RAY + halfB) / b;
    }

    /**
     * @dev Converts a wad (decimal 18) to ray (decimal 27) fixed point number
     * @param a Wad number
     * @return Ray fixed point number
     */
    function wadToRay(uint256 a) internal pure returns (uint256) {
        return a * 10 ** 9;
    }

    /**
     * @dev Converts a ray (decimal 27) to wad (decimal 18) fixed point number
     * @param a Ray number
     * @return Wad fixed point number
     */
    function rayToWad(uint256 a) internal pure returns (uint256) {
        uint256 result = a / 10 ** 9;
        return result;
    }

    /**
     * @dev Returns 1 in ray units
     */
    function ray() internal pure returns (uint256) {
        return RAY;
    }

    /**
     * @dev Returns 1 in wad units
     */
    function wad() internal pure returns (uint256) {
        return WAD;
    }

    /**
     * @dev Calculates pow(x, n) using ray math
     * @param x Base in ray
     * @param n Exponent
     * @return x^n in ray
     */
    function rayPow(uint256 x, uint256 n) internal pure returns (uint256) {
        if (n == 0) {
            return ray();
        }

        uint256 result = ray();

        for (uint256 i = 0; i < n; i++) {
            result = rayMul(result, x);
        }

        return result;
    }
}
