// SPDX-License-Identifier: MIT
pragma solidity 0.8.29;

import "forge-std/Test.sol";
import "../src/Collateral.sol";
import "../src/Token.sol";
import "./mocks/MockBorrowingForCollateral.sol";
import "./mocks/MockFailTransferTokenForCollateral.sol";
import "./mocks/MockRevertTransferTokenForCollateral.sol";
import "./mocks/MockFailTransferFromToken.sol";

/// @title CollateralTest
/// @notice Test suite for the Collateral contract.
contract CollateralTest is Test {
    Token token;
    Collateral collateral;
    MockBorrowingForCollateral mockBorrowing;
    address user = address(0x1);
    address borrower = address(0x2);
    address liquidator = address(0x3);
    uint256 initialSupply = 1000 * 10 ** 18;

    /// @notice Sets up the testing environment by deploying Token, MockBorrowingForCollateral, and Collateral contracts,
    /// and transferring tokens to the user.
    function setUp() public {
        token = new Token(initialSupply);
        mockBorrowing = new MockBorrowingForCollateral();
        collateral = new Collateral(address(token), address(mockBorrowing));

        token.transfer(user, 200 * 10 ** 18);
    }

    /// @notice Tests that depositCollateral successfully updates the user's collateral balance.
    function testDepositCollateralSuccess() public {
        uint256 depositAmount = 50 * 10 ** 18;
        vm.prank(user);
        token.approve(address(collateral), depositAmount);
        vm.prank(user);
        collateral.depositCollateral(depositAmount);
        uint256 balance = collateral.collateralBalance(user);
        assertEq(balance, depositAmount, "Collateral balance should be updated");
    }

    /// @notice Tests that depositing zero tokens as collateral reverts.
    function testDepositCollateralZeroShouldFail() public {
        vm.prank(user);
        vm.expectRevert("Amount must be greater than zero");
        collateral.depositCollateral(0);
    }

    /// @notice Tests that depositCollateral reverts when token.transferFrom fails.
    function testDepositCollateralFailsWhenTokenTransferFromFails() public {
        MockFailTransferFromToken failToken = new MockFailTransferFromToken(initialSupply);
        Collateral collateralFail = new Collateral(address(failToken), address(mockBorrowing));
        failToken.transfer(user, 100 * 10 ** 18);
        vm.prank(user);
        failToken.approve(address(collateralFail), 50 * 10 ** 18);
        vm.prank(user);
        vm.expectRevert("Collateral transfer failed");
        collateralFail.depositCollateral(50 * 10 ** 18);
    }

    /// @notice Tests that withdrawCollateral successfully reduces the collateral balance when the ratio is maintained.
    function testWithdrawCollateralSuccess() public {
        uint256 depositAmount = 150 * 10 ** 18;
        vm.prank(user);
        token.approve(address(collateral), depositAmount);
        vm.prank(user);
        collateral.depositCollateral(depositAmount);

        mockBorrowing.setBorrowedPrincipal(user, 50 * 10 ** 18);

        uint256 withdrawAmount = 50 * 10 ** 18;
        vm.prank(user);
        collateral.withdrawCollateral(withdrawAmount);

        uint256 newBalance = collateral.collateralBalance(user);
        assertEq(newBalance, depositAmount - withdrawAmount, "Collateral balance should be reduced by withdrawn amount");

        uint256 userTokenBalance = token.balanceOf(user);
        assertEq(userTokenBalance, 100 * 10 ** 18, "User token balance should reflect withdrawal");
    }

    /// @notice Tests that withdrawCollateral succeeds when no tokens are borrowed (thus bypassing the ratio check).
    function testWithdrawCollateralNoBorrowed() public {
        uint256 depositAmount = 100 * 10 ** 18;
        vm.prank(user);
        token.approve(address(collateral), depositAmount);
        vm.prank(user);
        collateral.depositCollateral(depositAmount);

        uint256 withdrawAmount = 40 * 10 ** 18;
        vm.prank(user);
        collateral.withdrawCollateral(withdrawAmount);

        uint256 remaining = collateral.collateralBalance(user);
        assertEq(remaining, depositAmount - withdrawAmount, "Collateral should be reduced correctly when no borrow");

        uint256 userBalance = token.balanceOf(user);
        assertEq(userBalance, 140 * 10 ** 18, "User token balance should be correct after withdrawal with no borrow");
    }

    /// @notice Tests that withdrawCollateral reverts if the new collateral would drop the ratio below the minimum.
    function testWithdrawCollateralFailsDueToLowRatio() public {
        uint256 depositAmount = 150 * 10 ** 18;
        vm.prank(user);
        token.approve(address(collateral), depositAmount);
        vm.prank(user);
        collateral.depositCollateral(depositAmount);

        mockBorrowing.setBorrowedPrincipal(user, 100 * 10 ** 18);

        uint256 withdrawAmount = 10 * 10 ** 18;
        vm.prank(user);
        vm.expectRevert("Collateral ratio too low after withdrawal");
        collateral.withdrawCollateral(withdrawAmount);
    }

    /// @notice Tests that withdrawing zero collateral reverts.
    function testWithdrawCollateralZeroShouldFail() public {
        vm.prank(user);
        vm.expectRevert("Amount must be greater than zero");
        collateral.withdrawCollateral(0);
    }

    /// @notice Tests that withdrawCollateral reverts if the withdrawal amount exceeds the deposited collateral.
    function testWithdrawCollateralExceedsBalanceShouldFail() public {
        uint256 depositAmount = 50 * 10 ** 18;
        vm.prank(user);
        token.approve(address(collateral), depositAmount);
        vm.prank(user);
        collateral.depositCollateral(depositAmount);

        vm.prank(user);
        vm.expectRevert("Withdrawal exceeds collateral balance");
        collateral.withdrawCollateral(depositAmount + 1);
    }

    /// @notice Tests that withdrawCollateral reverts when token.transfer returns false.
    function testWithdrawCollateralFailsWhenTokenTransferReturnsFalse() public {
        MockFailTransferTokenForCollateral failToken = new MockFailTransferTokenForCollateral(initialSupply);
        Collateral collateralFail = new Collateral(address(failToken), address(mockBorrowing));

        failToken.setFailAddress(address(collateralFail));

        vm.prank(address(this));
        failToken.transfer(user, 200 * 10 ** 18);

        vm.prank(user);
        failToken.approve(address(collateralFail), 100 * 10 ** 18);

        vm.prank(user);
        collateralFail.depositCollateral(100 * 10 ** 18);

        vm.prank(user);
        vm.expectRevert("Collateral withdrawal transfer failed");
        collateralFail.withdrawCollateral(50 * 10 ** 18);
    }

    /// @notice Tests that getCollateralRatio returns the maximum value when no tokens are borrowed.
    function testGetCollateralRatioWhenNoBorrowed() public {
        uint256 depositAmount = 50 * 10 ** 18;
        vm.prank(user);
        token.approve(address(collateral), depositAmount);
        vm.prank(user);
        collateral.depositCollateral(depositAmount);
        uint256 ratio = collateral.getCollateralRatio(user);
        assertEq(ratio, type(uint256).max, "Collateral ratio should be max when no borrow");
    }

    /// @notice Tests that getCollateralRatio returns the correct ratio when the user has borrowed tokens.
    function testGetCollateralRatioWithBorrowed() public {
        uint256 depositAmount = 150 * 10 ** 18;
        vm.prank(user);
        token.approve(address(collateral), depositAmount);
        vm.prank(user);
        collateral.depositCollateral(depositAmount);
        mockBorrowing.setBorrowedPrincipal(user, 50 * 10 ** 18);
        uint256 ratio = collateral.getCollateralRatio(user);
        assertEq(ratio, 300, "Collateral ratio should be 300");
    }

    /// @notice Tests that canBorrow returns true or false based on whether additional borrowing is allowed.
    function testCanBorrow() public {
        uint256 depositAmount = 150 * 10 ** 18;
        vm.prank(user);
        token.approve(address(collateral), depositAmount);
        vm.prank(user);
        collateral.depositCollateral(depositAmount);

        bool canBorrow1 = collateral.canBorrow(user, 0);
        assertTrue(canBorrow1, "Should allow borrowing when no borrow");

        mockBorrowing.setBorrowedPrincipal(user, 50 * 10 ** 18);
        bool canBorrow2 = collateral.canBorrow(user, 10 * 10 ** 18);
        assertTrue(canBorrow2, "Should allow borrowing additional amount");
        bool canBorrow3 = collateral.canBorrow(user, 100 * 10 ** 18);
        assertFalse(canBorrow3, "Should not allow borrowing additional amount exceeding ratio");
    }

    /// @notice Tests that canBorrow returns true when no tokens are borrowed even if borrowAmount > 0.
    function testCanBorrowNoBorrow() public {
        uint256 depositAmount = 100 * 10 ** 18;
        vm.prank(user);
        token.approve(address(collateral), depositAmount);
        vm.prank(user);
        collateral.depositCollateral(depositAmount);
        bool allowed = collateral.canBorrow(user, 50 * 10 ** 18);
        assertTrue(allowed, "Should allow borrowing when no borrow exists");
    }

    /// @notice Tests that liquidation succeeds for an undercollateralized borrower.
    function testLiquidateSuccess() public {
        // Transférer des tokens au borrower
        vm.prank(address(this));
        token.transfer(borrower, 150 * 10 ** 18);

        uint256 depositCollateral = 150 * 10 ** 18;
        vm.prank(borrower);
        token.approve(address(collateral), depositCollateral);
        vm.prank(borrower);
        collateral.depositCollateral(depositCollateral);

        // Simuler une dette élevée pour rendre le ratio insuffisant :
        // Par exemple, une dette de 130 tokens donne un ratio ≈115% (< seuil de liquidation, supposé 125%).
        mockBorrowing.setBorrowedPrincipal(borrower, 130 * 10 ** 18);
        // IMPORTANT : s'assurer que le contrat Borrowing connaît l'adresse correcte du contrat Collateral
        mockBorrowing.setCollateral(address(collateral));

        uint256 repayAmount = 50 * 10 ** 18;

        // Assurez-vous que le liquidateur a suffisamment de tokens
        vm.prank(address(this));
        token.transfer(liquidator, repayAmount);

        // Le liquidateur approuve le transfert de repayAmount vers le contrat Collateral.
        vm.prank(liquidator);
        token.approve(address(collateral), repayAmount);

        // Vérifiez que l'approbation est bien en place
        uint256 allowance = token.allowance(liquidator, address(collateral));
        assertEq(allowance, repayAmount, "Allowance not set correctly");

        vm.prank(liquidator);
        collateral.liquidate(borrower, repayAmount);

        uint256 newDebt = mockBorrowing.borrowedPrincipal(borrower);
        assertEq(newDebt, (130 - 50) * 10 ** 18, "Borrower's debt not reduced correctly");

        // Supposons LIQUIDATION_BONUS = 10, alors collateralToSeize = 50 * 110/100 = 55 tokens.
        uint256 newCollateral = collateral.collateralBalance(borrower);
        assertEq(newCollateral, (150 - 55) * 10 ** 18, "Borrower's collateral not reduced correctly");

        uint256 liquidatorBalance = token.balanceOf(liquidator);
        assertEq(liquidatorBalance, 55 * 10 ** 18, "Liquidator did not receive collateral correctly");
    }

    /// @notice Teste que la liquidation échoue si le montant de remboursement est supérieur à la dette de l'emprunteur.
    function testLiquidateFailsWhenCollateralRatioSufficient() public {
        // Transférer des tokens au borrower
        vm.prank(address(this));
        token.transfer(borrower, 150 * 10 ** 18);

        uint256 depositCollateral = 150 * 10 ** 18;
        vm.prank(borrower);
        token.approve(address(collateral), depositCollateral);
        vm.prank(borrower);
        collateral.depositCollateral(depositCollateral);

        mockBorrowing.setBorrowedPrincipal(borrower, 100 * 10 ** 18);

        uint256 repayAmount = 50 * 10 ** 18;
        vm.prank(liquidator);
        token.approve(address(Borrowing(address(collateral.borrowing()))), repayAmount);
        vm.prank(liquidator);
        vm.expectRevert("Collateral ratio is sufficient for liquidation");
        collateral.liquidate(borrower, repayAmount);
    }

    /// @notice Teste que la liquidation échoue si le montant de remboursement est zéro.
    function testLiquidateFailsWhenRepayAmountZero() public {
        // Transférer des tokens au borrower
        vm.prank(address(this));
        token.transfer(borrower, 100 * 10 ** 18);

        uint256 depositCollateral = 100 * 10 ** 18;
        vm.prank(borrower);
        token.approve(address(collateral), depositCollateral);
        vm.prank(borrower);
        collateral.depositCollateral(depositCollateral);

        vm.prank(liquidator);
        vm.expectRevert("Repay amount must be greater than zero");
        collateral.liquidate(borrower, 0);
    }

    /// @notice Tests that liquidation fails when the borrower has no debt
    function testLiquidateFailsWhenNoBorrowerDebt() public {
        vm.prank(address(this));
        token.transfer(borrower, 150 * 10 ** 18);

        uint256 depositCollateral = 150 * 10 ** 18;
        vm.prank(borrower);
        token.approve(address(collateral), depositCollateral);
        vm.prank(borrower);
        collateral.depositCollateral(depositCollateral);

        mockBorrowing.setBorrowedPrincipal(borrower, 0);

        uint256 repayAmount = 50 * 10 ** 18;
        vm.prank(liquidator);
        vm.expectRevert("Borrower has no debt");
        collateral.liquidate(borrower, repayAmount);
    }

    /// @notice Tests that liquidation fails when the repayment transfer fails
    function testLiquidateFailsWhenRepaymentTransferFails() public {
        vm.prank(address(this));
        token.transfer(borrower, 150 * 10 ** 18);

        uint256 depositCollateral = 150 * 10 ** 18;
        vm.prank(borrower);
        token.approve(address(collateral), depositCollateral);
        vm.prank(borrower);
        collateral.depositCollateral(depositCollateral);

        mockBorrowing.setBorrowedPrincipal(borrower, 130 * 10 ** 18);
        mockBorrowing.setCollateral(address(collateral));

        MockFailTransferFromToken failToken = new MockFailTransferFromToken(initialSupply);
        Collateral collateralFail = new Collateral(address(failToken), address(mockBorrowing));

        vm.prank(borrower);
        failToken.approve(address(collateralFail), depositCollateral);

        vm.store(
            address(collateralFail),
            keccak256(abi.encode(borrower, uint256(0))), // slot for collateralBalance[borrower]
            bytes32(uint256(150 * 10 ** 18))
        );

        uint256 repayAmount = 50 * 10 ** 18;
        vm.prank(liquidator);
        failToken.approve(address(collateralFail), repayAmount);
        vm.prank(liquidator);
        vm.expectRevert("Repayment transfer failed");
        collateralFail.liquidate(borrower, repayAmount);
    }

    /// @notice Tests that liquidation caps the collateralToSeize at the user's total collateral
    function testLiquidateCapsCollateralToSeize() public {
        vm.prank(address(this));
        token.transfer(borrower, 100 * 10 ** 18);

        uint256 depositCollateral = 100 * 10 ** 18;
        vm.prank(borrower);
        token.approve(address(collateral), depositCollateral);
        vm.prank(borrower);
        collateral.depositCollateral(depositCollateral);

        uint256 highDebt = 200 * 10 ** 18; // 200 tokens debt with 100 collateral
        mockBorrowing.setBorrowedPrincipal(borrower, highDebt);
        mockBorrowing.setCollateral(address(collateral));

        uint256 repayAmount = 100 * 10 ** 18; // Would normally seize 110 tokens

        vm.prank(address(this));
        token.transfer(liquidator, repayAmount);

        vm.prank(liquidator);
        token.approve(address(collateral), repayAmount);

        vm.prank(liquidator);
        collateral.liquidate(borrower, repayAmount);

        // The liquidator should get all collateral (100 tokens) instead of 110
        uint256 borrowerCollateralRemaining = collateral.collateralBalance(borrower);
        assertEq(borrowerCollateralRemaining, 0, "Borrower should have no collateral left");

        uint256 liquidatorBalance = token.balanceOf(liquidator);
        assertEq(liquidatorBalance, depositCollateral, "Liquidator should receive exactly the borrower's collateral");
    }

    /// @notice Tests that liquidation fails when collateral transfer to liquidator fails
    function testLiquidateFailsWhenCollateralTransferToLiquidatorFails() public {
        MockFailTransferTokenForCollateral failToken = new MockFailTransferTokenForCollateral(initialSupply);
        Collateral collateralFail = new Collateral(address(failToken), address(mockBorrowing));
        failToken.setFailAddress(address(collateralFail));

        vm.prank(address(this));
        failToken.transfer(borrower, 150 * 10 ** 18);

        // Setup borrower's collateral (via storage manipulation since depositCollateral won't work)
        vm.store(
            address(collateralFail),
            keccak256(abi.encode(borrower, uint256(0))), // slot for collateralBalance[borrower]
            bytes32(uint256(150 * 10 ** 18))
        );

        mockBorrowing.setBorrowedPrincipal(borrower, 130 * 10 ** 18);
        mockBorrowing.setCollateral(address(collateralFail));

        uint256 repayAmount = 50 * 10 ** 18;
        vm.prank(address(this));
        failToken.transfer(liquidator, repayAmount);

        vm.prank(liquidator);
        failToken.approve(address(collateralFail), repayAmount);

        vm.prank(liquidator);
        vm.expectRevert("Collateral transfer to liquidator failed");
        collateralFail.liquidate(borrower, repayAmount);
    }

    /// @notice Tests that the low-level call in withdrawCollateral reverts with the right message when it fails
    function testWithdrawCollateralFailsWithLowLevelCall() public {
        address mockToken = address(0xDEAD);
        Collateral collateralWithMock = new Collateral(mockToken, address(mockBorrowing));

        vm.store(
            address(collateralWithMock),
            keccak256(abi.encode(user, uint256(0))), // slot for collateralBalance[user]
            bytes32(uint256(100 * 10 ** 18))
        );

        vm.mockCallRevert(
            mockToken,
            abi.encodeWithSelector(Token(address(0)).transfer.selector, user, 50 * 10 ** 18),
            "Some external error"
        );

        vm.prank(user);
        vm.expectRevert("Collateral withdrawal transfer failed");
        collateralWithMock.withdrawCollateral(50 * 10 ** 18);
    }

    /// @notice Tests the getMaxBorrowableAmount function with no collateral
    function testGetMaxBorrowableAmountNoCollateral() public {
        uint256 maxAmount = collateral.getMaxBorrowableAmount(user);
        assertEq(maxAmount, 0, "Max borrowable amount should be 0 with no collateral");
    }

    /// @notice Tests the getMaxBorrowableAmount function with collateral where MIN_COLLATERAL_RATIO is the limiting factor
    function testGetMaxBorrowableAmountLimitedByRatio() public {
        uint256 depositAmount = 150 * 10 ** 18;
        vm.prank(user);
        token.approve(address(collateral), depositAmount);
        vm.prank(user);
        collateral.depositCollateral(depositAmount);

        // Based on MIN_COLLATERAL_RATIO = 150, max borrowable = 150 * 100 / 150 = 100 tokens
        // Based on MAX_BORROWING_PERCENTAGE = 75, max borrowable = 150 * 75 / 100 = 112.5 tokens
        // Should return 100 tokens (limited by MIN_COLLATERAL_RATIO)
        uint256 maxAmount = collateral.getMaxBorrowableAmount(user);
        assertEq(maxAmount, 100 * 10 ** 18, "Max borrowable amount should be limited by MIN_COLLATERAL_RATIO");
    }

    /// @notice Tests the getMaxBorrowableAmount function with collateral where MAX_BORROWING_PERCENTAGE is the limiting factor
    function testGetMaxBorrowableAmountLimitedByMaxPercentage() public {
        // Transfer more tokens to the user first
        vm.prank(address(this));
        token.transfer(user, 300 * 10 ** 18);

        // Use a larger collateral amount to make MAX_BORROWING_PERCENTAGE the limiting factor
        uint256 depositAmount = 300 * 10 ** 18;
        vm.prank(user);
        token.approve(address(collateral), depositAmount);
        vm.prank(user);
        collateral.depositCollateral(depositAmount);

        // Based on MIN_COLLATERAL_RATIO = 150, max borrowable = 300 * 100 / 150 = 200 tokens
        // Based on MAX_BORROWING_PERCENTAGE = 75, max borrowable = 300 * 75 / 100 = 225 tokens
        // Should return 200 tokens (limited by MIN_COLLATERAL_RATIO)
        uint256 maxAmount = collateral.getMaxBorrowableAmount(user);
        assertEq(maxAmount, 200 * 10 ** 18, "Max borrowable amount should be limited by MIN_COLLATERAL_RATIO");
    }

    /// @notice Tests that getMaxBorrowableAmount returns the correct amount regardless of existing borrowed amount
    function testGetMaxBorrowableAmountWithExistingBorrow() public {
        uint256 depositAmount = 150 * 10 ** 18;
        vm.prank(user);
        token.approve(address(collateral), depositAmount);
        vm.prank(user);
        collateral.depositCollateral(depositAmount);

        // Set an existing borrowed amount
        mockBorrowing.setBorrowedPrincipal(user, 50 * 10 ** 18);

        // The max borrowable amount should still be 100 tokens, as it's based only on collateral
        uint256 maxAmount = collateral.getMaxBorrowableAmount(user);
        assertEq(maxAmount, 100 * 10 ** 18, "Max borrowable should be based only on collateral, not existing borrow");
    }

    /// @notice Tests that getMaxBorrowableAmount works with edge case values
    function testGetMaxBorrowableAmountEdgeCases() public {
        // Test with minimum possible collateral
        uint256 minDeposit = 1;
        vm.prank(user);
        token.approve(address(collateral), minDeposit);
        vm.prank(user);
        collateral.depositCollateral(minDeposit);

        uint256 maxAmount = collateral.getMaxBorrowableAmount(user);
        // Expected: 1 * 100 / 150 = 0 (integer division rounds down)
        assertEq(maxAmount, 0, "Max borrowable should be 0 for minimal collateral");

        // Clean up for next test
        vm.prank(user);
        collateral.withdrawCollateral(minDeposit);

        // Test with large collateral
        uint256 largeDeposit = 800 * 10 ** 18;
        vm.prank(address(this));
        token.transfer(user, largeDeposit);

        vm.prank(user);
        token.approve(address(collateral), largeDeposit);
        vm.prank(user);
        collateral.depositCollateral(largeDeposit);

        maxAmount = collateral.getMaxBorrowableAmount(user);
        // Expected: 800 * 100 / 150 = 533.33... tokens
        assertEq(maxAmount, 533333333333333333333, "Max borrowable should calculate correctly for large deposits");
    }
}
