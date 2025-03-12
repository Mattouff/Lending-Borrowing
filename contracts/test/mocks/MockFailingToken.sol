// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import "../../src/Token.sol";

/// @title MockFailingToken
/// @notice A mock token that behaves like Token, except that calls to transfer revert if not initiated by the owner.
///         This allows simulating a failure during withdraw operations in tests.
contract MockFailingToken is Token {
    address public owner;

    /// @notice Constructor that mints the initial supply and stores the deployer as owner.
    /// @param initialSupply The initial token supply.
    constructor(uint256 initialSupply) Token(initialSupply) {
        owner = msg.sender;
    }

    /// @notice Overrides transfer so that transfers fail unless initiated by the owner.
    /// @param recipient The address receiving the tokens.
    /// @param amount The amount of tokens to transfer.
    /// @return A boolean indicating whether the transfer was successful.
    function transfer(address recipient, uint256 amount) public override returns (bool) {
        if (msg.sender == owner) {
            return super.transfer(recipient, amount);
        }
        revert("Token transfer failed");
    }
}
