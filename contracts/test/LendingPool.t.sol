// SPDX-License-Identifier: MIT
pragma solidity 0.8.29;

import "forge-std/Test.sol";
import "../src/Token.sol";
import "../src/LendingPool.sol";
import "./mocks/MockFailingToken.sol";
import "./mocks/MockReturnFalseOnTransferToken.sol";
import "@openzeppelin/contracts/token/ERC20/ERC20.sol";

/// @notice Token simulating a failure on transferFrom (for testing deposit)
contract FailingTransferFromToken is ERC20 {
    constructor(uint256 initialSupply) ERC20("FailingTransferFromToken", "FTFT") {
        _mint(msg.sender, initialSupply);
    }
    /// @dev This function simulates a systematic failure of transferFrom.

    function transferFrom(address, address, uint256) public pure override returns (bool) {
        return false;
    }
}

/// @notice Token simulating a failure on transfer (for testing withdraw)
contract FailingTransferToken is Token {
    constructor(uint256 initialSupply) Token(initialSupply) {}
    /// @notice Overrides transfer so that transfers fail unless initiated by the deployer.
    /// @dev This allows initial allocation to the user but forces a revert when the LendingPool attempts to transfer tokens.

    function transfer(address recipient, uint256 amount) public override returns (bool) {
        if (msg.sender == address(this)) {
            return super.transfer(recipient, amount);
        }
        revert("Token transfer failed");
    }
}

contract LendingPoolTest is Test {
    Token token;
    LendingPool lendingPool;
    address user = address(0x1);
    uint256 initialSupply = 1000 * 10 ** 18;

    function setUp() public {
        token = new Token(initialSupply);
        lendingPool = new LendingPool(address(token));
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
        assertEq(userUnderlyingBalance, (100 - 50 + 20) * 10 ** 18, "User underlying balance after withdraw incorrect");
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
        FailingTransferFromToken failingToken = new FailingTransferFromToken(initialSupply);
        failingToken.transfer(user, 100 * 10 ** 18);
        LendingPool poolWithFailingToken = new LendingPool(address(failingToken));
        uint256 depositAmount = 50 * 10 ** 18;
        vm.prank(user);
        failingToken.approve(address(poolWithFailingToken), depositAmount);
        vm.prank(user);
        vm.expectRevert("Token transfer failed");
        poolWithFailingToken.deposit(depositAmount);
    }

    /// @notice Test simulating a failure in withdraw via a failing transfer.
    function testWithdrawFailsWhenTransferFails() public {
        MockFailingToken failingToken = new MockFailingToken(initialSupply);
        failingToken.transfer(user, 100 * 10 ** 18);
        LendingPool poolWithFailingToken = new LendingPool(address(failingToken));
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
        LendingPool poolWithFailingToken = new LendingPool(address(failingToken));
        uint256 depositAmount = 50 * 10 ** 18;
        vm.prank(user);
        failingToken.approve(address(poolWithFailingToken), depositAmount);
        vm.prank(user);
        poolWithFailingToken.deposit(depositAmount);
        vm.prank(user);
        vm.expectRevert("Token transfer failed");
        poolWithFailingToken.withdraw(20 * 10 ** 18);
    }

    /// @notice Tests getLendingToken returns the correct amount for a user.
    function testGetLendingToken() public {
        uint256 depositAmount = 50 * 10 ** 18;
        vm.prank(user);
        token.approve(address(lendingPool), depositAmount);
        vm.prank(user);
        lendingPool.deposit(depositAmount);
        uint256 lending = lendingPool.getLendingToken(user);
        assertEq(lending, depositAmount, "getLendingToken should return deposit amount");
    }

    /// @notice Tests getAllLendingToken returns the total lending across all users.
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
}
