// SPDX-License-Identifier: MIT
pragma solidity 0.8.29;

import "../../src/Token.sol";

/// @title MockFailTransferToken
/// @notice A mock token that simulates a failure on transfer by always returning false.
contract MockFailTransferToken is Token {
    /// @notice Constructor that mints the initial supply.
    /// @param initialSupply The amount of tokens to mint.
    constructor(uint256 initialSupply) Token(initialSupply) {}

    /// @notice Overrides transfer to always return false.
    /// @param recipient The recipient of the tokens.
    /// @param amount The amount of tokens to transfer.
    /// @return Always false.
    function transfer(address recipient, uint256 amount) public pure override returns (bool) {
        return false;
    }
}
