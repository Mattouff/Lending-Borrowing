// SPDX-License-Identifier: MIT
pragma solidity 0.8.29;

import "../../src/Token.sol";

/// @title MockFailTransferTokenForCollateral
/// @notice A mock token that simulates a failure on transfer when called by a designated address.
contract MockFailTransferTokenForCollateral is Token {
    /// @notice The address for which transfers should fail.
    address public failAddress;

    /// @notice Constructor that mints the initial supply.
    /// @param initialSupply The amount of tokens to mint.
    constructor(uint256 initialSupply) Token(initialSupply) {}

    /// @notice Sets the address for which transfers should fail.
    /// @param _failAddress The address that, when msg.sender equals it, transfer returns false.
    function setFailAddress(address _failAddress) external {
        failAddress = _failAddress;
    }

    /// @notice Overrides transfer so that if msg.sender equals failAddress, it returns false.
    /// @param recipient The recipient of the tokens.
    /// @param amount The amount of tokens to transfer.
    /// @return Returns false if msg.sender is failAddress; otherwise, behaves normally.
    function transfer(address recipient, uint256 amount) public override returns (bool) {
        if (msg.sender == failAddress) {
            return false;
        }
        return super.transfer(recipient, amount);
    }
}
