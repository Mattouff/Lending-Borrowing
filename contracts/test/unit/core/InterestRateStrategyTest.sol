// SPDX-License-Identifier: MIT
pragma solidity 0.8.29;

import "forge-std/Test.sol";
import "../../../src/core/InterestRateStrategy.sol";
import "../../helpers/TestingHelper.sol";

/**
 * @title InterestRateStrategyTest
 * @author DeFi Lending Platform
 * @notice Unit tests for the InterestRateStrategy contract
 */
contract InterestRateStrategyTest is Test {
    // Contract under test
    InterestRateStrategy private strategy;

    // Helper contract
    TestingHelper private helper;

    // Test accounts
    address private admin;
    address private user;

    // Constants
    uint256 private constant WAD = 1e18;
    uint256 private constant BASE_BORROW_RATE = WAD / 100; // 1%
    uint256 private constant SLOPE1 = WAD / 10; // 10%
    uint256 private constant SLOPE2 = WAD * 4 / 10; // 40%
    uint256 private constant OPTIMAL_UTILIZATION_RATE = WAD * 8 / 10; // 80%
    uint256 private constant ELASTICITY_FACTOR = 2 * WAD; // 2.0 (quadratic)
    uint256 private constant RESERVE_FACTOR = WAD / 10; // 10%

    /**
     * @dev Set up the test fixture
     */
    function setUp() public {
        helper = new TestingHelper();
        
        admin = makeAddr("admin");
        user = makeAddr("user");
        
        vm.startPrank(admin);
        strategy = new InterestRateStrategy(
            BASE_BORROW_RATE,
            SLOPE1,
            SLOPE2,
            OPTIMAL_UTILIZATION_RATE,
            ELASTICITY_FACTOR
        );
        vm.stopPrank();
    }

    /**
     * @dev Test the initialization of the interest rate strategy
     */
    function testInitialization() public view {
        assertEq(strategy.getBaseVariableBorrowRate(), BASE_BORROW_RATE, "Base borrow rate should be set correctly");
        
        uint256 maxRate = BASE_BORROW_RATE + SLOPE1 + SLOPE2;
        assertEq(strategy.getMaxVariableBorrowRate(), maxRate, "Max borrow rate should be calculated correctly");
    }

    /**
     * @dev Test interest rate calculation when utilization is 0
     */
    function testZeroUtilization() public view {
        address asset = address(0x1); // Dummy asset address
        uint256 availableLiquidity = 1000 * WAD;
        uint256 totalBorrows = 0;
        
        (
            uint256 liquidityRate,
            uint256 stableBorrowRate,
            uint256 variableBorrowRate
        ) = strategy.calculateInterestRates(asset, availableLiquidity, totalBorrows, RESERVE_FACTOR);
        
        // With zero utilization, borrow rate should equal base rate
        assertEq(variableBorrowRate, BASE_BORROW_RATE, "Variable borrow rate should equal base rate at 0 utilization");
        
        // Liquidity rate should be 0 at 0% utilization (no borrowers)
        assertEq(liquidityRate, 0, "Liquidity rate should be 0 at 0% utilization");
        
        // We're not implementing stable borrow rate in this version
        assertEq(stableBorrowRate, 0, "Stable borrow rate should be 0");
    }

    /**
     * @dev Test interest rate calculation when utilization is at optimal level
     */
    function testOptimalUtilization() public view {
        address asset = address(0x1); // Dummy asset address
        uint256 availableLiquidity = 200 * WAD;
        uint256 totalBorrows = 800 * WAD;
        
        // This is 80% utilization (optimal)
        
        (
            uint256 liquidityRate,
            uint256 stableBorrowRate,
            uint256 variableBorrowRate
        ) = strategy.calculateInterestRates(asset, availableLiquidity, totalBorrows, RESERVE_FACTOR);
        
        // At optimal utilization, borrow rate should be base rate + slope1
        uint256 expectedBorrowRate = BASE_BORROW_RATE + SLOPE1;
        assertApproxEqRel(variableBorrowRate, expectedBorrowRate, 0.0001e18, "Variable borrow rate should be base + slope1 at optimal utilization");
        
        // Liquidity rate = borrow rate * utilization * (1 - reserve factor)
        uint256 utilizationRate = (totalBorrows * WAD) / (availableLiquidity + totalBorrows);
        uint256 expectedLiquidityRate = (variableBorrowRate * utilizationRate * (WAD - RESERVE_FACTOR)) / (WAD * WAD);
        assertApproxEqRel(liquidityRate, expectedLiquidityRate, 0.0001e18, "Liquidity rate should be calculated correctly");
    }

    /**
     * @dev Test interest rate calculation when utilization is above optimal level
     */
    function testHighUtilization() public view {
        address asset = address(0x1); // Dummy asset address
        uint256 availableLiquidity = 100 * WAD;
        uint256 totalBorrows = 900 * WAD;
        
        // This is 90% utilization (above optimal)
        
        (
            uint256 liquidityRate,
            uint256 stableBorrowRate,
            uint256 variableBorrowRate
        ) = strategy.calculateInterestRates(asset, availableLiquidity, totalBorrows, RESERVE_FACTOR);
        
        // Calculate the excess utilization (normalized to 0-1 scale)
        uint256 excessUtilization = (90 * WAD / 100 - OPTIMAL_UTILIZATION_RATE) * WAD / (WAD - OPTIMAL_UTILIZATION_RATE);
        
        // Calculate expected variable rate using the model in InterestRateStrategy.sol
        uint256 normalRate = BASE_BORROW_RATE + SLOPE1;
        
        // At high utilization, rate increases non-linearly based on elasticity factor
        uint256 expectedBorrowRate = normalRate + SLOPE2 / 4; // Simplified for test, should grow faster
        
        // Check that rate is higher than at optimal utilization
        uint256 optimalBorrowRate = BASE_BORROW_RATE + SLOPE1;
        assertTrue(variableBorrowRate > optimalBorrowRate, "Variable borrow rate should be higher than at optimal utilization");
        
        // Liquidity rate should follow the same pattern
        uint256 utilizationRate = (totalBorrows * WAD) / (availableLiquidity + totalBorrows);
        uint256 expectedLiquidityRate = (variableBorrowRate * utilizationRate * (WAD - RESERVE_FACTOR)) / (WAD * WAD);
        assertApproxEqRel(liquidityRate, expectedLiquidityRate, 0.0001e18, "Liquidity rate should be calculated correctly");
    }
    
    /**
     * @dev Test that interest rates properly respond to utilization changes
     */
    function testInterestRateResponseToUtilization() public view {
        address asset = address(0x1); // Dummy asset address
        
        // Array of utilization rates to test
        uint256[] memory utilRates = new uint256[](5);
        utilRates[0] = 0; // 0%
        utilRates[1] = 40 * WAD / 100; // 40%
        utilRates[2] = 80 * WAD / 100; // 80% (optimal)
        utilRates[3] = 90 * WAD / 100; // 90%
        utilRates[4] = 95 * WAD / 100; // 95%
        
        uint256[] memory borrowRates = new uint256[](5);
        
        // Calculate rates at each utilization
        for (uint256 i = 0; i < utilRates.length; i++) {
            uint256 totalLiquidity = 1000 * WAD;
            uint256 totalBorrows = (totalLiquidity * utilRates[i]) / WAD;
            uint256 availableLiquidity = totalLiquidity - totalBorrows;
            
            (, , uint256 variableBorrowRate) = strategy.calculateInterestRates(
                asset, 
                availableLiquidity, 
                totalBorrows, 
                RESERVE_FACTOR
            );
            
            borrowRates[i] = variableBorrowRate;
        }
        
        // Rates should increase monotonically with utilization
        for (uint256 i = 1; i < borrowRates.length; i++) {
            assertTrue(borrowRates[i] > borrowRates[i-1], "Interest rates should increase with utilization");
        }
        
        // Rate increase should accelerate after optimal utilization
        uint256 diffAt40To80 = borrowRates[2] - borrowRates[1];
        uint256 diffAt80To90 = borrowRates[3] - borrowRates[2];
        uint256 diffAt90To95 = borrowRates[4] - borrowRates[3];
        
        // Comparing per-percent increase in rate to verify non-linear growth
        uint256 ratePerPercentBeforeOptimal = diffAt40To80 / 40; // 40% difference in utilization
        uint256 ratePerPercentAfterOptimal1 = diffAt80To90 / 10;  // 10% difference
        uint256 ratePerPercentAfterOptimal2 = diffAt90To95 / 5;   // 5% difference
        
        assertTrue(ratePerPercentAfterOptimal1 > ratePerPercentBeforeOptimal, 
            "Rate increase should accelerate after optimal utilization");
        assertTrue(ratePerPercentAfterOptimal2 > ratePerPercentAfterOptimal1, 
            "Rate increase should continue to accelerate as utilization approaches 100%");
    }
    
    /**
     * @dev Test ownership management
     */
    function testOwnership() public {
        // Only owner should be able to transfer ownership
        vm.startPrank(user);
        vm.expectRevert();
        strategy.transferOwnership(user);
        vm.stopPrank();
        
        // Owner should be able to transfer ownership
        vm.startPrank(admin);
        strategy.transferOwnership(user);
        vm.stopPrank();
        
        assertEq(strategy.owner(), user, "Owner should be updated immediately");
    }
}
