// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "forge-std/Test.sol";
import "../src/Token.sol";

contract TokenTest is Test {
    Token public token;
    address public owner;
    uint256 public initialSupply = 1000 * 10**18; 

    function setUp() public {
        owner = address(this); 
        token = new Token(initialSupply);
    }

    function testInitialBalance() public view{
        uint256 balance = token.balanceOf(owner);
        assertEq(balance, initialSupply, "Insufficient initial balance");
    }

    function testNameAndSymbol() public view {
        assertEq(token.name(), "LendingBorrowing", "Incorrect name");
        assertEq(token.symbol(), "LBT", "Incorrect symbol");
    }

    function testTransfer() public {
        address receiver = address(0xBEEF);
        uint256 amount = 100 * 10**18; 

        token.transfer(receiver, amount);

        assertEq(token.balanceOf(receiver), amount, "Recipient balance incorrect");
        assertEq(token.balanceOf(owner), initialSupply - amount, "Owner balance incorrect");
    }

    function testApprovalAndTransferFrom() public {
        address spender = address(0xCAFE);
        uint256 amount = 50 * 10**18;

        token.approve(spender, amount);
        assertEq(token.allowance(owner, spender), amount, "Allocation incorrect");

        address receiver = address(0xDAD);

        vm.prank(spender);
        token.transferFrom(owner, receiver, amount);

        assertEq(token.balanceOf(receiver), amount, "Recipient balance incorrect");
        assertEq(token.allowance(owner, spender), 0, "Allocation incorrect");
    }
}
