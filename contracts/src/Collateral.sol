// SPDX-License-Identifier: MIT
pragma solidity 0.8.29;

import "./Token.sol";
import "./Borrowing.sol";

/// @title Collateral - A contract to manage collateral deposits for a decentralized lending platform.
/// @notice Users can deposit and withdraw collateral. This contract interacts with the Borrowing contract to define borrowing limits based on collateral.
contract Collateral {
    Token public immutable token;
    Borrowing public immutable borrowing;

    /// @notice The minimum collateral ratio required, expressed as a percentage (e.g., 150 means 150%).
    uint256 public constant MIN_COLLATERAL_RATIO = 150;

    /// @notice Mapping of user addresses to their deposited collateral amounts.
    mapping(address => uint256) public collateralBalance;

    /// @notice Constructor that sets the underlying token and the Borrowing contract.
    /// @param _token The address of the token used as collateral.
    /// @param _borrowing The address of the Borrowing contract defining collateral conditions.
    constructor(address _token, address _borrowing) {
        token = Token(_token);
        borrowing = Borrowing(_borrowing);
    }

    /// @notice Deposits collateral into the contract.
    /// @param amount The amount of tokens to deposit as collateral.
    /// @dev The user must first approve the token transfer.
    function depositCollateral(uint256 amount) external {
        require(amount > 0, "Amount must be greater than zero");
        require(token.transferFrom(msg.sender, address(this), amount), "Collateral transfer failed");
        collateralBalance[msg.sender] += amount;
    }

    /// @notice Withdraws collateral from the contract.
    /// @param amount The amount of collateral to withdraw.
    /// @dev Withdrawal is allowed only if the remaining collateral maintains the minimum collateral ratio based on the user's borrowed amount.
    function withdrawCollateral(uint256 amount) external {
        require(amount > 0, "Amount must be greater than zero");
        require(collateralBalance[msg.sender] >= amount, "Withdrawal exceeds collateral balance");

        uint256 newCollateral = collateralBalance[msg.sender] - amount;
        uint256 borrowed = borrowing.borrowedBalance(msg.sender);
        if (borrowed > 0) {
            require(newCollateral * 100 >= borrowed * MIN_COLLATERAL_RATIO, "Collateral ratio too low after withdrawal");
        }
        collateralBalance[msg.sender] = newCollateral;
        
        // Appel bas niveau via .call pour intercepter un revert ou un retour false
        (bool success, bytes memory returndata) = address(token).call(
            abi.encodeWithSelector(token.transfer.selector, msg.sender, amount)
        );
        if (!success) {
            revert("Collateral withdrawal transfer failed");
        }
        if (returndata.length > 0) {
            // Décode le retour pour vérifier qu'il renvoie true.
            bool transferSuccess = abi.decode(returndata, (bool));
            require(transferSuccess, "Collateral withdrawal transfer failed");
        }
    }

    /// @notice Gets the collateral ratio for a specific user.
    /// @param user The address of the user.
    /// @return The collateral ratio expressed as a percentage.
    function getCollateralRatio(address user) external view returns (uint256) {
        uint256 borrowed = borrowing.borrowedBalance(user);
        if (borrowed == 0) {
            return type(uint256).max;
        }
        return (collateralBalance[user] * 100) / borrowed;
    }

    /// @notice Checks if a user can borrow an additional amount based on their collateral.
    /// @param user The address of the user.
    /// @param borrowAmount The additional amount the user intends to borrow.
    /// @return True if the user can borrow the additional amount while maintaining the minimum collateral ratio.
    function canBorrow(address user, uint256 borrowAmount) external view returns (bool) {
        uint256 totalBorrowed = borrowing.borrowedBalance(user) + borrowAmount;
        if (totalBorrowed == 0) {
            return true;
        }
        return (collateralBalance[user] * 100) >= (totalBorrowed * MIN_COLLATERAL_RATIO);
    }
}
