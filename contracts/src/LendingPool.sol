// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "./Token.sol";

/// @title LendingPool - A contract to manage deposits for a decentralized lending platform
/// @notice Users deposit ERC20 tokens (defined in Token.sol) and receive deposit tokens (dToken) representing their share in the pool.
contract LendingPool is ERC20 {
    Token public immutable underlying;

    /// @notice Constructor for the LendingPool contract.
    /// @param _underlying The address of the underlying token (Token.sol)
    constructor(address _underlying) ERC20("Deposit Token", "dTOKEN") {
        underlying = Token(_underlying);
    }

    /// @notice Deposit tokens into the lending pool.
    /// @param amount The amount of tokens to deposit.
    /// @dev The user must first approve the transfer of tokens to this contract.
    function deposit(uint256 amount) external {
        require(amount > 0, "Amount must be greater than zero");

        bool success = underlying.transferFrom(msg.sender, address(this), amount);
        require(success, "Token transfer failed");

        _mint(msg.sender, amount);
    }

    /// @notice Withdraw tokens from the pool by burning the corresponding deposit tokens.
    /// @param amount The amount of tokens to withdraw.
    function withdraw(uint256 amount) external {
        require(amount > 0, "Amount must be greater than zero");
        require(balanceOf(msg.sender) >= amount, "Insufficient deposit balance");

        _burn(msg.sender, amount);

        bool success = underlying.transfer(msg.sender, amount);
        require(success, "Token transfer failed");
    }
}
