// SPDX-License-Identifier: MIT
pragma solidity 0.8.29;

/**
 * @title IInterestRateStrategy
 * @author DeFi Lending Platform
 * @notice Interface for the InterestRateStrategy
 */
interface IInterestRateStrategy {
    /**
     * @dev Returns the base variable borrow rate
     * @return The base variable borrow rate
     */
    function getBaseVariableBorrowRate() external view returns (uint256);

    /**
     * @dev Returns the maximum variable borrow rate
     * @return The maximum variable borrow rate
     */
    function getMaxVariableBorrowRate() external view returns (uint256);

    /**
     * @dev Calculates the interest rates based on utilization
     * @param reserve The address of the reserve
     * @param availableLiquidity The available liquidity
     * @param totalBorrows The total borrows
     * @param reserveFactor The reserve factor
     * @return liquidityRate The liquidity rate
     * @return stableBorrowRate The stable borrow rate
     * @return variableBorrowRate The variable borrow rate
     */
    function calculateInterestRates(
        address reserve,
        uint256 availableLiquidity,
        uint256 totalBorrows,
        uint256 reserveFactor
    ) external view returns (uint256 liquidityRate, uint256 stableBorrowRate, uint256 variableBorrowRate);
}
