// SPDX-License-Identifier: MIT
pragma solidity 0.8.29;

import "../../src/Token.sol";

/// @title MockFailTransferFromToken
/// @notice A mock token that simulates a failure on transferFrom by always returning false.
contract MockFailTransferFromToken is Token {
    /// @notice Constructor that mints the initial supply.
    /// @param initialSupply The amount of tokens to mint.
    constructor(uint256 initialSupply) Token(initialSupply) {}

    /// @notice Overrides transferFrom to always return false.
    /// @param sender The sender of the tokens.
    /// @param recipient The recipient of the tokens.
    /// @param amount The amount of tokens to transfer.
    /// @return Always false.
    function transferFrom(address sender, address recipient, uint256 amount) public pure override returns (bool) {
        return false;
    }
}
