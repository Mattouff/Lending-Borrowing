// SPDX-License-Identifier: MIT
pragma solidity 0.8.29;

import "./Token.sol";
import "./Collateral.sol";

/// @title Borrowing - A contract to manage borrowing operations for a decentralized lending platform.
/// @notice Users can borrow tokens and repay their loans. Borrowing limits are defined by the collateral managed in Collateral.sol.
contract Borrowing {
    Token public immutable token;
    Collateral public immutable collateral;

    /// @notice Mapping of user addresses to their borrowed token amounts.
    mapping(address => uint256) public borrowedBalance;
    /// @notice Total borrowed tokens across all users.
    uint256 public totalBorrowed;

    /// @notice Constructor for the Borrowing contract.
    /// @param _token The address of the token used for borrowing.
    /// @param _collateral The address of the Collateral contract defining collateral conditions.
    constructor(address _token, address _collateral) {
        token = Token(_token);
        collateral = Collateral(_collateral);
    }

    /// @notice Allows a user to borrow tokens.
    /// @param amount The amount of tokens to borrow.
    /// @dev Collateral conditions are enforced via Collateral.canBorrow.
    function borrow(uint256 amount) external {
        require(amount > 0, "Amount must be greater than zero");
        require(collateral.canBorrow(msg.sender, amount), "Insufficient collateral");
        borrowedBalance[msg.sender] += amount;
        totalBorrowed += amount;
        require(token.transfer(msg.sender, amount), "Borrow transfer failed");
    }

    /// @notice Repays a portion or all of the borrowed tokens.
    /// @param amount The amount of tokens to repay.
    /// @dev The caller must have approved the contract to transfer tokens on their behalf.
    function repay(uint256 amount) external {
        require(amount > 0, "Amount must be greater than zero");
        require(borrowedBalance[msg.sender] >= amount, "Repay amount exceeds borrowed balance");
        require(token.transferFrom(msg.sender, address(this), amount), "Repay transfer failed");
        borrowedBalance[msg.sender] -= amount;
        totalBorrowed -= amount;
    }

    /// @notice Gets the borrowed token amount for a specific user.
    /// @param user The address of the user.
    /// @return The amount of tokens borrowed by the user.
    function getBorrowToken(address user) external view returns (uint256) {
        return borrowedBalance[user];
    }

    /// @notice Gets the total borrowed tokens across all users.
    /// @return The total amount of tokens borrowed.
    function getAllBorrowToken() external view returns (uint256) {
        return totalBorrowed;
    }
}
