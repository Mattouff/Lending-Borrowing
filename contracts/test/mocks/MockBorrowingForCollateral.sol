// SPDX-License-Identifier: MIT
pragma solidity 0.8.29;

/// @title MockBorrowingForCollateral
/// @notice A minimal mock for Borrowing used for testing Collateral; it allows setting and retrieving a user's borrowed balance.
contract MockBorrowingForCollateral {
    mapping(address => uint256) private _borrowedBalance;

    /// @notice Sets the borrowed balance for a user.
    /// @param user The address of the user.
    /// @param amount The borrowed balance to set.
    function setBorrowedBalance(address user, uint256 amount) external {
        _borrowedBalance[user] = amount;
    }
    
    /// @notice Returns the borrowed balance for a user.
    /// @param user The address of the user.
    /// @return The borrowed balance.
    function getBorrowedBalance(address user) external view returns (uint256) {
        return _borrowedBalance[user];
    }
    
    /// @notice Mimics the public getter in the original Borrowing contract.
    /// @param user The address of the user.
    /// @return The borrowed balance.
    function borrowedBalance(address user) external view returns (uint256) {
        return _borrowedBalance[user];
    }
}
