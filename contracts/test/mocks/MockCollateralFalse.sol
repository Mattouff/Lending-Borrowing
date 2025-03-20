// SPDX-License-Identifier: MIT
pragma solidity 0.8.29;

/// @title MockCollateralFalse
/// @notice A mock collateral contract for testing purposes that always disallows borrowing.
contract MockCollateralFalse {
    /// @notice Always returns false for any borrowing request.
    /// @param user The address of the borrower.
    /// @param borrowAmount The additional amount the borrower intends to borrow.
    /// @return False, indicating that the borrower cannot borrow.
    function canBorrow(address user, uint256 borrowAmount) external pure returns (bool) {
        return false;
    }
}
