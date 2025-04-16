// SPDX-License-Identifier: MIT
pragma solidity 0.8.29;

import "../../src/Token.sol";

/// @title MockRevertTransferTokenForCollateral
/// @notice A mock token that simulates a revert on transfer.
contract MockRevertTransferTokenForCollateral is Token {
    /// @notice Constructor that mints the initial supply.
    /// @param initialSupply The amount of tokens to mint.
    constructor(uint256 initialSupply) Token(initialSupply) {}

    /// @notice Overrides transfer to always revert with the expected message.
    /// @param recipient The recipient of the tokens.
    /// @param amount The amount of tokens to transfer.
    function transfer(address recipient, uint256 amount) public pure override returns (bool) {
        revert("Collateral withdrawal transfer failed");
    }
}
