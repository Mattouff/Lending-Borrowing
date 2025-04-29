// SPDX-License-Identifier: MIT
pragma solidity 0.8.29;

/// @title CompoundInterest - A library for calculating compound interest
/// @notice Provides functions to calculate compound interest with daily compounding
library CompoundInterest {
    /// @notice Calculates compound interest with daily compounding
    /// @param principal The initial amount
    /// @param annualRate The annual interest rate with 18 decimals (e.g., 5% = 5e16)
    /// @param timeElapsed The time elapsed in seconds
    /// @return newBalance The new balance after applying compound interest
    /// @return interest The interest earned (newBalance - principal)
    function calculateCompoundInterest(uint256 principal, uint256 annualRate, uint256 timeElapsed)
        internal
        pure
        returns (uint256 newBalance, uint256 interest)
    {
        if (principal == 0 || timeElapsed == 0) {
            return (principal, 0);
        }

        uint256 daysElapsed = timeElapsed / (1 days);
        if (daysElapsed == 0) {
            return (principal, 0);
        }

        uint256 dailyRate = annualRate / 365; // Daily rate

        // Calculate compound factor: (1 + r/n)^(days)
        uint256 compoundFactor = 1e18; // Start with 1 in fixed-point
        for (uint256 i = 0; i < daysElapsed; i++) {
            compoundFactor = (compoundFactor * (1e18 + dailyRate)) / 1e18;
        }

        // Calculate new balance with compound interest
        newBalance = (principal * compoundFactor) / 1e18;
        interest = newBalance - principal;

        return (newBalance, interest);
    }
}
