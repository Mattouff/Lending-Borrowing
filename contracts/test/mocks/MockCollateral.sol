// SPDX-License-Identifier: MIT
pragma solidity 0.8.29;

/// @title MockCollateral
/// @notice A mock collateral contract for testing purposes that always allows borrowing.
contract MockCollateral {
    /// @notice Always returns true for any borrowing request.
    /// @param user The address of the borrower.
    /// @param borrowAmount The additional amount the borrower intends to borrow.
    /// @return True, indicating that the borrower can always borrow.
    function canBorrow(address user, uint256 borrowAmount) external pure returns (bool) {
        return true;
    }
}
