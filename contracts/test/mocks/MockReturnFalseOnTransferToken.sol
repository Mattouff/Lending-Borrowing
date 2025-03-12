// SPDX-License-Identifier: MIT
pragma solidity 0.8.29;

import "../../src/Token.sol";

/// @title MockReturnFalseOnTransferToken
/// @notice A mock token that behaves like Token, but its transfer function returns false (instead of reverting)
///         when called by an address other than the owner. This allows simulating a failure during withdraw operations.
contract MockReturnFalseOnTransferToken is Token {
    address public owner;

    /// @notice Constructor that mints the initial supply and stores the deployer as owner.
    /// @param initialSupply The initial token supply.
    constructor(uint256 initialSupply) Token(initialSupply) {
        owner = msg.sender;
    }

    /// @notice Overrides transfer so that transfers return false unless initiated by the owner.
    /// @param recipient The address receiving the tokens.
    /// @param amount The amount of tokens to transfer.
    /// @return A boolean indicating whether the transfer was successful.
    function transfer(address recipient, uint256 amount) public override returns (bool) {
        // Allow transfers if initiated by the owner (e.g. for initial allocation)
        if (msg.sender == owner) {
            return super.transfer(recipient, amount);
        }
        // Otherwise, simulate a failure by returning false (instead of reverting)
        return false;
    }
}
