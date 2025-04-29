// SPDX-License-Identifier: MIT
pragma solidity 0.8.29;

import "forge-std/Test.sol";
import "../src/libraries/CompoundInterest.sol";

contract CompoundInterestTest is Test {
    // Helper function to call the library function directly
    function calculateCompoundInterest(uint256 principal, uint256 annualRate, uint256 timeElapsed)
        public
        pure
        returns (uint256 newBalance, uint256 interest)
    {
        return CompoundInterest.calculateCompoundInterest(principal, annualRate, timeElapsed);
    }

    /// @notice Tests the early return condition when principal is zero
    function testCalculateCompoundInterestZeroPrincipal() pure public {
        uint256 annualRate = 5e16; // 5%
        uint256 timeElapsed = 365 days;

        (uint256 newBalance, uint256 interest) = calculateCompoundInterest(0, annualRate, timeElapsed);

        assertEq(newBalance, 0, "New balance should be 0 when principal is 0");
        assertEq(interest, 0, "Interest should be 0 when principal is 0");
    }

    /// @notice Tests the early return condition when timeElapsed is zero
    function testCalculateCompoundInterestZeroTimeElapsed() pure public {
        uint256 principal = 100e18;
        uint256 annualRate = 5e16; // 5%

        (uint256 newBalance, uint256 interest) = calculateCompoundInterest(principal, annualRate, 0);

        assertEq(newBalance, principal, "New balance should equal principal when timeElapsed is 0");
        assertEq(interest, 0, "Interest should be 0 when timeElapsed is 0");
    }

    /// @notice Tests the early return condition when days elapsed is less than 1
    function testCalculateCompoundInterestLessThanOneDay() pure public {
        uint256 principal = 100e18;
        uint256 annualRate = 5e16; // 5%
        uint256 timeElapsed = 23 hours; // Less than 1 day

        (uint256 newBalance, uint256 interest) = calculateCompoundInterest(principal, annualRate, timeElapsed);

        assertEq(newBalance, principal, "New balance should equal principal when less than 1 day has passed");
        assertEq(interest, 0, "Interest should be 0 when less than 1 day has passed");
    }

    /// @notice Tests exactly one day of interest
    function testCalculateCompoundInterestExactlyOneDay() pure public {
        uint256 principal = 100e18;
        uint256 annualRate = 5e16; // 5%
        uint256 timeElapsed = 1 days; // Exactly 1 day

        (uint256 newBalance, uint256 interest) = calculateCompoundInterest(principal, annualRate, timeElapsed);

        uint256 dailyRate = annualRate / 365;
        uint256 expectedNewBalance = (principal * (1e18 + dailyRate)) / 1e18;
        uint256 expectedInterest = expectedNewBalance - principal;

        assertEq(newBalance, expectedNewBalance, "New balance calculation incorrect for 1 day");
        assertEq(interest, expectedInterest, "Interest calculation incorrect for 1 day");
    }

    /// @notice Tests normal computation with multiple days
    function testCalculateCompoundInterestNormalCase() pure public {
        uint256 principal = 1000e18;
        uint256 annualRate = 5e16; // 5%
        uint256 timeElapsed = 30 days; // 30 days

        (uint256 newBalance, uint256 interest) = calculateCompoundInterest(principal, annualRate, timeElapsed);

        // Manual calculation of expected result
        uint256 dailyRate = annualRate / 365;
        uint256 compoundFactor = 1e18;
        for (uint256 i = 0; i < 30; i++) {
            compoundFactor = (compoundFactor * (1e18 + dailyRate)) / 1e18;
        }
        uint256 expectedNewBalance = (principal * compoundFactor) / 1e18;
        uint256 expectedInterest = expectedNewBalance - principal;

        assertEq(newBalance, expectedNewBalance, "New balance calculation incorrect for 30 days");
        assertEq(interest, expectedInterest, "Interest calculation incorrect for 30 days");
    }

    /// @notice Tests that interest compounds correctly for a full year
    function testCalculateCompoundInterestFullYear() pure public {
        uint256 principal = 1000e18;
        uint256 annualRate = 5e16; // 5%
        uint256 timeElapsed = 365 days; // Full year

        (uint256 newBalance, uint256 interest) = calculateCompoundInterest(principal, annualRate, timeElapsed);

        // With daily compounding over a year, the result should be approximately
        // principal * (1 + rate)^1, but slightly higher due to compounding
        uint256 simpleInterest = (principal * annualRate) / 1e18;
        uint256 simpleTotal = principal + simpleInterest;

        // The compound interest should be greater than simple interest
        assertGt(newBalance, simpleTotal, "Compound interest should exceed simple interest");
        assertGt(interest, simpleInterest, "Compound interest difference should exceed simple interest");

        // The result should be approximately principal * (1 + rate/365)^365
        uint256 dailyRate = annualRate / 365;
        uint256 compoundFactor = 1e18;
        for (uint256 i = 0; i < 365; i++) {
            compoundFactor = (compoundFactor * (1e18 + dailyRate)) / 1e18;
        }
        uint256 expectedNewBalance = (principal * compoundFactor) / 1e18;

        // Allow for small precision errors with a relative tolerance
        assertApproxEqRel(newBalance, expectedNewBalance, 0.0001e18, "Full year compound calculation incorrect");
    }
}
