// SPDX-License-Identifier: MIT
pragma solidity 0.8.29;

import "forge-std/Test.sol";
import "../src/Token.sol";
import "../src/LendingPool.sol";
import "./mocks/MockFailingToken.sol";
import "./mocks/MockReturnFalseOnTransferToken.sol";
import "./mocks/MockFailTransferFromToken.sol";
import "./mocks/MockFailTransferToken.sol";
import "@openzeppelin/contracts/token/ERC20/ERC20.sol";

contract LendingPoolTest is Test {
    Token token;
    LendingPool lendingPool;
    address user = address(0x1);
    uint256 initialSupply = 1000 * 10 ** 18;

    // Pour les tests d'intérêt, nous utilisons un taux annuel de 5% (5e16 avec 18 décimales)
    function setUp() public {
        token = new Token(initialSupply);
        lendingPool = new LendingPool(address(token), 5e16);
        // Transférer 100 tokens à l'utilisateur pour qu'il puisse déposer
        token.transfer(user, 100 * 10 ** 18);
    }

    function testDeposit() public {
        uint256 depositAmount = 50 * 10 ** 18;
        vm.prank(user);
        token.approve(address(lendingPool), depositAmount);
        vm.prank(user);
        lendingPool.deposit(depositAmount);
        uint256 poolBalance = token.balanceOf(address(lendingPool));
        assertEq(poolBalance, depositAmount, "Pool balance incorrect");
        uint256 receiptBalance = lendingPool.balanceOf(user);
        assertEq(receiptBalance, depositAmount, "Receipt token balance incorrect");
    }

    function testWithdraw() public {
        uint256 depositAmount = 50 * 10 ** 18;
        vm.prank(user);
        token.approve(address(lendingPool), depositAmount);
        vm.prank(user);
        lendingPool.deposit(depositAmount);
        uint256 withdrawAmount = 20 * 10 ** 18;
        vm.prank(user);
        lendingPool.withdraw(withdrawAmount);
        uint256 poolBalance = token.balanceOf(address(lendingPool));
        assertEq(poolBalance, depositAmount - withdrawAmount, "Pool balance after withdraw incorrect");
        uint256 receiptBalance = lendingPool.balanceOf(user);
        assertEq(receiptBalance, depositAmount - withdrawAmount, "Receipt token balance after withdraw incorrect");
        uint256 userUnderlyingBalance = token.balanceOf(user);
        // User initial 100, après dépôt = 50, après retrait = 50+20 = 70.
        assertEq(userUnderlyingBalance, 70 * 10 ** 18, "User underlying balance after withdraw incorrect");
    }

    function testDepositZeroShouldFail() public {
        vm.prank(user);
        token.approve(address(lendingPool), 0);
        vm.prank(user);
        vm.expectRevert("Amount must be greater than zero");
        lendingPool.deposit(0);
    }

    function testWithdrawZeroShouldFail() public {
        uint256 depositAmount = 10 * 10 ** 18;
        vm.prank(user);
        token.approve(address(lendingPool), depositAmount);
        vm.prank(user);
        lendingPool.deposit(depositAmount);
        vm.prank(user);
        vm.expectRevert("Amount must be greater than zero");
        lendingPool.withdraw(0);
    }

    function testWithdrawMoreThanBalanceShouldFail() public {
        uint256 depositAmount = 10 * 10 ** 18;
        vm.prank(user);
        token.approve(address(lendingPool), depositAmount);
        vm.prank(user);
        lendingPool.deposit(depositAmount);
        vm.prank(user);
        vm.expectRevert("Insufficient deposit balance");
        lendingPool.withdraw(depositAmount + 1);
    }

    /// @notice Test simulating a failure in deposit via a failing transferFrom.
    function testDepositFailsWhenTransferFromFails() public {
        MockFailTransferFromToken failingToken = new MockFailTransferFromToken(initialSupply);
        failingToken.transfer(user, 100 * 10 ** 18);
        LendingPool poolWithFailingToken = new LendingPool(address(failingToken), 5e16);
        uint256 depositAmount = 50 * 10 ** 18;
        vm.prank(user);
        failingToken.approve(address(poolWithFailingToken), depositAmount);
        vm.prank(user);
        vm.expectRevert("Token transfer failed");
        poolWithFailingToken.deposit(depositAmount);
    }

    /// @notice Test simulating a failure in withdraw via a failing transfer.
    function testWithdrawFailsWhenTransferFails() public {
        // Utilise le mock qui reverte sur transfer.
        MockFailingToken failingToken = new MockFailingToken(initialSupply);
        failingToken.transfer(user, 100 * 10 ** 18);
        LendingPool poolWithFailingToken = new LendingPool(address(failingToken), 5e16);
        uint256 depositAmount = 50 * 10 ** 18;
        vm.prank(user);
        failingToken.approve(address(poolWithFailingToken), depositAmount);
        vm.prank(user);
        poolWithFailingToken.deposit(depositAmount);
        vm.prank(user);
        vm.expectRevert(bytes("Token transfer failed"));
        poolWithFailingToken.withdraw(20 * 10 ** 18);
    }

    /// @notice Test simulating a failure in withdraw when underlying.transfer returns false.
    function testWithdrawFailsWhenTransferReturnsFalse() public {
        MockReturnFalseOnTransferToken failingToken = new MockReturnFalseOnTransferToken(initialSupply);
        failingToken.transfer(user, 100 * 10 ** 18);
        LendingPool poolWithFailingToken = new LendingPool(address(failingToken), 5e16);
        uint256 depositAmount = 50 * 10 ** 18;
        vm.prank(user);
        failingToken.approve(address(poolWithFailingToken), depositAmount);
        vm.prank(user);
        poolWithFailingToken.deposit(depositAmount);
        vm.prank(user);
        vm.expectRevert("Token transfer failed");
        poolWithFailingToken.withdraw(20 * 10 ** 18);
    }

    /// @notice Tests getLendingToken returns the correct deposit amount for a user.
    function testGetLendingToken() public {
        uint256 depositAmount = 50 * 10 ** 18;
        vm.prank(user);
        token.approve(address(lendingPool), depositAmount);
        vm.prank(user);
        lendingPool.deposit(depositAmount);
        uint256 lending = lendingPool.getLendingToken(user);
        assertEq(lending, depositAmount, "getLendingToken should return deposit amount");
    }

    /// @notice Tests getAllLendingToken returns the total deposits across all users.
    function testGetAllLendingToken() public {
        uint256 deposit1 = 30 * 10 ** 18;
        uint256 deposit2 = 20 * 10 ** 18;
        vm.prank(user);
        token.approve(address(lendingPool), deposit1);
        vm.prank(user);
        lendingPool.deposit(deposit1);

        address user2 = address(0x2);
        token.transfer(user2, 50 * 10 ** 18);
        vm.prank(user2);
        token.approve(address(lendingPool), deposit2);
        vm.prank(user2);
        lendingPool.deposit(deposit2);

        uint256 total = lendingPool.getAllLendingToken();
        assertEq(total, deposit1 + deposit2, "getAllLendingToken should return the sum of deposits");
    }

    /// @notice Tests that interest accrues over time on deposits.
    function testInterestAccrualOnDeposit() public {
        uint256 depositAmount = 100 * 10 ** 18;
        vm.prank(user);
        token.approve(address(lendingPool), depositAmount);
        vm.prank(user);
        lendingPool.deposit(depositAmount);

        // Transférer des tokens supplémentaires à l'utilisateur pour permettre un dépôt additionnel.
        vm.prank(address(this));
        token.transfer(user, 50 * 10 ** 18);

        // Simuler le passage d'un an (365 jours)
        uint256 oneYear = 365 days;
        uint256 initialLending = lendingPool.getLendingToken(user);

        vm.warp(block.timestamp + oneYear);

        // Effectuer un dépôt additionnel pour déclencher l'actualisation des intérêts
        uint256 additionalDeposit = 50 * 10 ** 18;
        vm.prank(user);
        token.approve(address(lendingPool), additionalDeposit);
        vm.prank(user);
        lendingPool.deposit(additionalDeposit);

        uint256 updatedLending = lendingPool.getLendingToken(user);

        // Calcul théorique de l'intérêt :
        // intérêt = depositAmount * annualInterestRate * oneYear / (365 days * 1e18)
        uint256 expectedInterest = depositAmount * lendingPool.annualInterestRate() * oneYear / (365 days * 1e18);

        // Le solde mis à jour en dTokens doit être au moins égal à (initialLending + additionalDeposit + intérêt)
        assertGe(
            updatedLending,
            initialLending + additionalDeposit + expectedInterest,
            "Interest not accrued correctly on deposit"
        );
    }

    /// @notice Tests that interest accrues over time on deposits and is correctly updated on withdrawal.
    function testInterestAccrualOnWithdraw() public {
        uint256 depositAmount = 100 * 10 ** 18;
        vm.prank(user);
        token.approve(address(lendingPool), depositAmount);
        vm.prank(user);
        lendingPool.deposit(depositAmount);

        // Avancer le temps d'un an (365 jours) pour générer un intérêt plus important
        uint256 oneYear = 365 days;
        vm.warp(block.timestamp + oneYear);

        // Forcer la mise à jour des intérêts
        lendingPool.updateUserInterest(user);

        uint256 lendingBefore = lendingPool.getLendingToken(user);
        // On s'attend à ce que l'intérêt ait été accumulé, donc lendingBefore > depositAmount
        assertGt(lendingBefore, depositAmount, "Interest should accrue on deposit");

        // Retirer une partie du dépôt
        uint256 withdrawAmount = 50 * 10 ** 18;
        vm.prank(user);
        lendingPool.withdraw(withdrawAmount);

        uint256 lendingAfter = lendingPool.getLendingToken(user);
        // La balance en dTokens doit être diminuée exactement du montant retiré
        assertEq(
            lendingAfter, lendingBefore - withdrawAmount, "Lending token balance not updated correctly after withdrawal"
        );
    }
}
