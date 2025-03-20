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
    uint256 initialSupply = 1000 * 10 ** 18;

    /// @notice Sets up the testing environment by deploying Token, MockCollateral, and Borrowing contracts,
    /// and funding the Borrowing contract with tokens for lending.
    function setUp() public {
        token = new Token(initialSupply);
        mockCollateral = new MockCollateral();
        borrowing = new Borrowing(address(token), address(mockCollateral));

        token.transfer(address(borrowing), 500 * 10 ** 18);
        token.transfer(user, 100 * 10 ** 18);
    }

    /// @notice Tests that borrowing a valid amount succeeds and updates balances correctly.
    function testBorrowSuccess() public {
        vm.prank(user);
        borrowing.borrow(50 * 10 ** 18);

        assertEq(borrowing.borrowedBalance(user), 50 * 10 ** 18, "Borrowed balance should be 50 tokens");
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
        Borrowing localBorrowing = new Borrowing(address(localToken), address(falseCollateral));
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
        Borrowing localBorrowing = new Borrowing(address(failToken), address(collateral));
        failToken.transfer(address(localBorrowing), 500 * 10 ** 18);
        failToken.transfer(user, 100 * 10 ** 18);

        vm.prank(user);
        vm.expectRevert("Borrow transfer failed");
        localBorrowing.borrow(50 * 10 ** 18);
    }

    /// @notice Tests that a successful partial repayment reduces the borrowed balance accordingly.
    function testRepaySuccess() public {
        vm.prank(user);
        borrowing.borrow(50 * 10 ** 18);

        vm.prank(user);
        token.approve(address(borrowing), 50 * 10 ** 18);

        vm.prank(user);
        borrowing.repay(20 * 10 ** 18);

        assertEq(borrowing.borrowedBalance(user), 30 * 10 ** 18, "Borrowed balance should be 30 tokens after repayment");
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
        Borrowing localBorrowing = new Borrowing(address(failFromToken), address(collateral));
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
        uint256 userBorrowed = borrowing.getBorrowToken(user);
        assertEq(userBorrowed, 50 * 10 ** 18, "getBorrowToken should return 50 tokens");
    }

    /// @notice Tests the getAllBorrowToken function to verify it returns the correct total borrowed amount.
    function testGetAllBorrowToken() public {
        vm.prank(user);
        borrowing.borrow(50 * 10 ** 18);
        uint256 totalBorrowed = borrowing.getAllBorrowToken();
        assertEq(totalBorrowed, 50 * 10 ** 18, "getAllBorrowToken should return 50 tokens");
    }
}
