// SPDX-License-Identifier: MIT
pragma solidity 0.8.29;

/// @title MockBorrowingForCollateral
/// @notice A minimal mock for Borrowing used for testing Collateral; it allows setting and retrieving a user's borrowed balance.
contract MockBorrowingForCollateral {
    mapping(address => uint256) private _borrowedPrincipal;
    address private _collateral;

    /// @notice Sets the borrowed balance for a user.
    /// @param user The address of the user.
    /// @param amount The borrowed balance to set.
    function setBorrowedPrincipal(address user, uint256 amount) external {
        _borrowedPrincipal[user] = amount;
    }

    /// @notice Returns the borrowed balance for a user.
    /// @param user The address of the user.
    /// @return The borrowed balance.
    function getBorrowedPrincipal(address user) external view returns (uint256) {
        return _borrowedPrincipal[user];
    }

    /// @notice Mimics the public getter in the original Borrowing contract.
    /// @param user The address of the user.
    /// @return The borrowed balance.
    function borrowedPrincipal(address user) external view returns (uint256) {
        return _borrowedPrincipal[user];
    }

    function getBorrowToken(address user) external view returns (uint256) {
        return _borrowedPrincipal[user];
    }

    /// @notice Sets the collateral contract address.
    /// @param collateralAddress The address of the Collateral contract.
    function setCollateral(address collateralAddress) external {
        _collateral = collateralAddress;
    }

    /// @notice Returns the collateral contract address.
    /// @return The address of the Collateral contract.
    function getCollateral() external view returns (address) {
        return _collateral;
    }

    /// @notice Reduces the debt of a borrower by the specified amount.
    /// @dev Only callable by the set collateral contract.
    /// @param borrower The address of the borrower.
    /// @param amount The amount by which to reduce the debt.
    function reduceDebt(address borrower, uint256 amount) external {
        require(msg.sender == _collateral, "Not authorized");
        require(_borrowedPrincipal[borrower] >= amount, "Insufficient debt");
        _borrowedPrincipal[borrower] -= amount;
    }
}
