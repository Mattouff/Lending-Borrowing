// SPDX-License-Identifier: MIT
pragma solidity 0.8.29;

import "forge-std/Test.sol";
import "../src/Borrowing.sol";
import "../src/Token.sol";
import "./mocks/MockCollateral.sol";
import "./mocks/MockCollateralFalse.sol";
import "./mocks/MockFailTransferToken.sol";
import "./mocks/MockFailTransferFromToken.sol";

/// @notice Test suite for the Borrowing contract.
contract BorrowingTest is Test {
    Token token;
    Borrowing borrowing;
    MockCollateral mockCollateral;
    address user = address(0x1);
    address liquidator = address(0x2);
    uint256 initialSupply = 1000 * 10 ** 18;

    /// @notice Sets up the testing environment by deploying Token, MockCollateral, and Borrowing contracts,
    /// and funding the Borrowing contract with tokens for lending.
    function setUp() public {
        token = new Token(initialSupply);
        mockCollateral = new MockCollateral();
        borrowing = new Borrowing(address(token), address(mockCollateral), 1e16, 2e16, 1e18);

        token.transfer(address(borrowing), 500 * 10 ** 18);
        token.transfer(user, 100 * 10 ** 18);
        token.transfer(liquidator, 100 * 10 ** 18);
    }

    /// @notice Tests that borrowing a valid amount succeeds and updates balances correctly.
    function testBorrowSuccess() public {
        vm.warp(1000); // Set a fixed timestamp to start

        vm.prank(user);
        borrowing.borrow(50 * 10 ** 18);

        assertEq(borrowing.borrowedPrincipal(user), 50 * 10 ** 18, "Borrowed balance should be 50 tokens");
        assertEq(borrowing.getAllBorrowToken(), 50 * 10 ** 18, "Total borrowed should be 50 tokens");
        assertEq(token.balanceOf(user), 150 * 10 ** 18, "User token balance should be 150 tokens");
    }

    /// @notice Tests that borrowing zero tokens reverts.
    function testBorrowZeroShouldFail() public {
        vm.prank(user);
        vm.expectRevert("Amount must be greater than zero");
        borrowing.borrow(0);
    }

    /// @notice Tests that borrowing fails when collateral conditions are not met.
    function testBorrowFailsDueToCollateral() public {
        Token localToken = new Token(initialSupply);
        MockCollateralFalse falseCollateral = new MockCollateralFalse();
        Borrowing localBorrowing = new Borrowing(address(localToken), address(falseCollateral), 1e16, 2e16, 1e18);
        localToken.transfer(address(localBorrowing), 500 * 10 ** 18);
        localToken.transfer(user, 100 * 10 ** 18);

        vm.prank(user);
        vm.expectRevert("Insufficient collateral");
        localBorrowing.borrow(50 * 10 ** 18);
    }

    /// @notice Tests that borrowing fails when token.transfer returns false.
    function testBorrowFailsWhenTokenTransferFails() public {
        MockFailTransferToken failToken = new MockFailTransferToken(initialSupply);
        MockCollateral collateral = new MockCollateral();
        Borrowing localBorrowing = new Borrowing(address(failToken), address(collateral), 1e16, 2e16, 1e18);
        failToken.transfer(address(localBorrowing), 500 * 10 ** 18);
        failToken.transfer(user, 100 * 10 ** 18);

        vm.prank(user);
        vm.expectRevert("Borrow transfer failed");
        localBorrowing.borrow(50 * 10 ** 18);
    }

    /// @notice Tests that a successful partial repayment reduces the borrowed balance accordingly.
    function testRepaySuccess() public {
        vm.warp(1000);

        vm.prank(user);
        borrowing.borrow(50 * 10 ** 18);

        // Stay at the same timestamp to avoid interest accrual

        vm.prank(user);
        token.approve(address(borrowing), 50 * 10 ** 18);

        vm.prank(user);
        borrowing.repay(20 * 10 ** 18);

        assertEq(
            borrowing.borrowedPrincipal(user), 30 * 10 ** 18, "Borrowed balance should be 30 tokens after repayment"
        );
        assertEq(borrowing.getAllBorrowToken(), 30 * 10 ** 18, "Total borrowed should be 30 tokens after repayment");
        assertEq(token.balanceOf(user), 130 * 10 ** 18, "User token balance should be 130 tokens after repayment");
    }

    /// @notice Tests that repaying zero tokens reverts.
    function testRepayZeroShouldFail() public {
        vm.prank(user);
        borrowing.borrow(50 * 10 ** 18);

        vm.prank(user);
        token.approve(address(borrowing), 50 * 10 ** 18);

        vm.prank(user);
        vm.expectRevert("Amount must be greater than zero");
        borrowing.repay(0);
    }

    /// @notice Tests that repaying more than the borrowed amount reverts.
    function testRepayMoreThanBorrowedShouldFail() public {
        vm.prank(user);
        borrowing.borrow(50 * 10 ** 18);

        vm.prank(user);
        token.approve(address(borrowing), 100 * 10 ** 18);

        vm.prank(user);
        vm.expectRevert("Repay amount exceeds borrowed balance");
        borrowing.repay(60 * 10 ** 18);
    }

    /// @notice Tests that repaying fails when token.transferFrom returns false.
    function testRepayFailsWhenTokenTransferFromFails() public {
        MockFailTransferFromToken failFromToken = new MockFailTransferFromToken(initialSupply);
        MockCollateral collateral = new MockCollateral();
        Borrowing localBorrowing = new Borrowing(address(failFromToken), address(collateral), 1e16, 2e16, 1e18);
        failFromToken.transfer(address(localBorrowing), 500 * 10 ** 18);
        failFromToken.transfer(user, 100 * 10 ** 18);

        vm.prank(user);
        localBorrowing.borrow(50 * 10 ** 18);

        vm.prank(user);
        failFromToken.approve(address(localBorrowing), 50 * 10 ** 18);

        vm.prank(user);
        vm.expectRevert("Repay transfer failed");
        localBorrowing.repay(20 * 10 ** 18);
    }

    /// @notice Tests the getBorrowToken function to verify it returns the correct borrowed amount for a user.
    function testGetBorrowToken() public {
        vm.prank(user);
        borrowing.borrow(50 * 10 ** 18);
        uint256 userBorrowed = borrowing.borrowedPrincipal(user);
        assertEq(userBorrowed, 50 * 10 ** 18, "getBorrowToken should return 50 tokens");
    }

    /// @notice Tests the getAllBorrowToken function to verify it returns the correct total borrowed amount.
    function testGetAllBorrowToken() public {
        vm.prank(user);
        borrowing.borrow(50 * 10 ** 18);
        uint256 totalBorrowed = borrowing.getAllBorrowToken();
        assertEq(totalBorrowed, 50 * 10 ** 18, "getAllBorrowToken should return 50 tokens");
    }

    /// @notice Tests the reduceDebt function called by the collateral contract.
    function testReduceDebt() public {
        vm.warp(1000);

        vm.prank(user);
        borrowing.borrow(50 * 10 ** 18);

        vm.prank(address(mockCollateral));
        borrowing.reduceDebt(user, 20 * 10 ** 18);

        assertEq(borrowing.borrowedPrincipal(user), 30 * 10 ** 18, "Debt should be reduced to 30 tokens");
        assertEq(borrowing.getAllBorrowToken(), 30 * 10 ** 18, "Total borrowed should be reduced to 30 tokens");
    }

    /// @notice Tests that reduceDebt fails when called by an unauthorized address.
    function testReduceDebtUnauthorized() public {
        vm.prank(user);
        borrowing.borrow(50 * 10 ** 18);

        vm.prank(address(0x3));
        vm.expectRevert("Not authorized");
        borrowing.reduceDebt(user, 20 * 10 ** 18);
    }

    /// @notice Tests that reduceDebt fails when trying to reduce more than the borrowed amount.
    function testReduceDebtInsufficientDebt() public {
        vm.prank(user);
        borrowing.borrow(50 * 10 ** 18);

        vm.prank(address(mockCollateral));
        vm.expectRevert("Insufficient debt");
        borrowing.reduceDebt(user, 60 * 10 ** 18);
    }

    /// @notice Tests the getCurrentRate function to verify it returns the correct interest rate.
    function testGetCurrentRate() public {
        // Case 1: Empty pool (no borrowing)
        uint256 rate = borrowing.getCurrentRate();
        assertEq(rate, 1e16, "Initial rate should be the minimum rate (1%)");

        // Case 2: Some borrowing - should increase the rate
        vm.prank(user);
        borrowing.borrow(250 * 10 ** 18); // Borrow 50% of the pool

        rate = borrowing.getCurrentRate();
        // With 50% utilization, the rate should be between rMin and rMax
        assertTrue(rate > 1e16 && rate <= 2e16, "Rate should increase with utilization");

        // Case 3: High utilization - should approach maximum rate
        vm.prank(user);
        borrowing.borrow(200 * 10 ** 18); // Borrow more, pushing utilization higher

        rate = borrowing.getCurrentRate();
        assertTrue(rate <= 2e16 && rate > 1e16, "Rate should be close to maximum with high utilization");
    }

    /// @notice Tests the getCurrentRate function with an empty pool (capacity = 0).
    function testGetCurrentRateEmptyPool() public {
        Token emptyToken = new Token(100 * 10 ** 18);
        Borrowing emptyBorrowing = new Borrowing(address(emptyToken), address(mockCollateral), 1e16, 2e16, 1e18);

        uint256 rate = emptyBorrowing.getCurrentRate();
        assertEq(rate, 1e16, "Empty pool should return minimum rate");
    }

    /// @notice Tests the internal power function by checking the rate calculation with different beta values.
    function testPowerFunctionViaRateCalculation() public {
        Borrowing borrowingBeta1 = new Borrowing(address(token), address(mockCollateral), 1e16, 2e16, 1e18); // Beta = 1
        Borrowing borrowingBeta2 = new Borrowing(address(token), address(mockCollateral), 1e16, 2e16, 2e18); // Beta = 2

        token.transfer(address(borrowingBeta1), 100 * 10 ** 18);
        token.transfer(address(borrowingBeta2), 100 * 10 ** 18);

        vm.startPrank(user);
        borrowingBeta1.borrow(50 * 10 ** 18);
        borrowingBeta2.borrow(50 * 10 ** 18);
        vm.stopPrank();

        uint256 rateBeta1 = borrowingBeta1.getCurrentRate();
        uint256 rateBeta2 = borrowingBeta2.getCurrentRate();

        // With beta = 1, the rate should be linear: rMin + (rMax - rMin) * 0.5 = 1.5%
        // With beta = 2, the rate should be quadratic: rMin + (rMax - rMin) * 0.5^2 = 1.25%
        assertTrue(rateBeta1 > rateBeta2, "Rate with beta=1 should be higher than with beta=2 at 50% utilization");
    }

    /// @notice Tests the updateBorrowedPrincipal function by checking interest accrual over time.
    function testUpdateBorrowedPrincipal() public {
        vm.prank(user);
        borrowing.borrow(100 * 10 ** 18);

        uint256 initialBalance = borrowing.borrowedPrincipal(user);

        vm.warp(block.timestamp + 365 days);

        vm.prank(user);
        borrowing.borrow(1); // Small amount to trigger updateBorrowedPrincipal

        uint256 newBalance = borrowing.borrowedPrincipal(user);
        assertTrue(newBalance > initialBalance, "Interest should be added after time passes");

        uint256 expectedInterest = (initialBalance * 1e16 * 365 days) / (365 days * 1e18);
        uint256 expectedBalance = initialBalance + expectedInterest + 1; // +1 for the extra borrow

        assertApproxEqRel(newBalance, expectedBalance, 0.01e18, "Interest calculation should match expected amount");
    }

    /// @notice Tests the updateBorrowedPrincipal function with no time elapsed.
    function testUpdateBorrowedPrincipalNoTimeElapsed() public {
        vm.prank(user);
        borrowing.borrow(100 * 10 ** 18);

        uint256 initialBalance = borrowing.borrowedPrincipal(user);

        vm.prank(user);
        borrowing.borrow(1);

        uint256 newBalance = borrowing.borrowedPrincipal(user);
        assertEq(newBalance, initialBalance + 1, "No interest should be added when no time has elapsed");
    }

    /// @notice Tests the updateBorrowedPrincipal function with zero borrowed balance.
    function testUpdateBorrowedPrincipalZeroBalance() public {
        vm.prank(user);
        borrowing.borrow(1);

        vm.prank(user);
        token.approve(address(borrowing), 1);

        vm.prank(user);
        borrowing.repay(1);

        vm.warp(block.timestamp + 365 days);

        vm.prank(user);
        borrowing.borrow(50 * 10 ** 18);

        assertEq(borrowing.borrowedPrincipal(user), 50 * 10 ** 18, "No interest should be added to zero balance");
    }

    /// @notice Tests that updateBorrowedPrincipal handles the first-time update properly.
    function testUpdateBorrowedPrincipalFirstTimeUpdate() public {
        vm.prank(user);
        borrowing.borrow(100 * 10 ** 18);

        assertEq(borrowing.borrowedPrincipal(user), 100 * 10 ** 18, "First update should not add interest");
    }

    /// @notice Tests the complete borrowing and repayment flow with interest accrual.
    function testCompleteFlowWithInterest() public {
        vm.warp(1000);

        vm.prank(user);
        borrowing.borrow(100 * 10 ** 18);

        uint256 initialBorrowedAmount = borrowing.getBorrowToken(user);

        vm.warp(1000 + 365 days);

        uint256 newBorrowedAmount = borrowing.getBorrowToken(user);
        assertGt(newBorrowedAmount, initialBorrowedAmount, "Interest should accrue after time passes");

        // Calculate expected compound interest
        uint256 currentRate = borrowing.getCurrentRate();
        uint256 dailyRate = currentRate / 365;

        // Calculate the compound factor: (1 + dailyRate)^365
        uint256 compoundFactor = 1e18;
        for (uint256 i = 0; i < 365; i++) {
            compoundFactor = (compoundFactor * (1e18 + dailyRate)) / 1e18;
        }

        uint256 expectedTotal = (initialBorrowedAmount * compoundFactor) / 1e18;

        // Allow for small precision errors with a delta
        assertApproxEqAbs(newBorrowedAmount, expectedTotal, 100, "Interest calculation should be approximately correct");

        // Test repayment after interest accrual
        vm.prank(user);
        token.approve(address(borrowing), 200 * 10 ** 18); // Approve more than needed

        vm.prank(user);
        borrowing.repay(50 * 10 ** 18);

        // Check that the state is updated correctly after repayment
        uint256 afterRepayAmount = borrowing.getBorrowToken(user);
        assertEq(afterRepayAmount, newBorrowedAmount - 50 * 10 ** 18, "Balance should decrease by repay amount");
    }

    /// @notice Tests the early return conditions in getBorrowToken function.
    function testGetBorrowTokenEarlyReturns() public {
        // Case 1: lastUpdateTime is 0 (user has never borrowed)
        assertEq(borrowing.lastUpdateTime(user), 0, "Initial lastUpdateTime should be 0");
        assertEq(borrowing.getBorrowToken(user), 0, "getBorrowToken should return 0 for user who never borrowed");

        // Case 2: User borrows and then repays everything
        vm.prank(user);
        borrowing.borrow(50 * 10 ** 18);

        vm.prank(user);
        token.approve(address(borrowing), 50 * 10 ** 18);

        vm.prank(user);
        borrowing.repay(50 * 10 ** 18);

        assertEq(borrowing.borrowedPrincipal(user), 0, "borrowedPrincipal should be 0 after full repayment");
        assertNotEq(borrowing.lastUpdateTime(user), 0, "lastUpdateTime should not be 0 after activity");

        vm.warp(block.timestamp + 365 days);

        assertEq(borrowing.getBorrowToken(user), 0, "getBorrowToken should return 0 when borrowedPrincipal is 0");
    }
}
