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
    /// @notice The liquidation threshold ratio, expressed as a percentage (e.g., 125 means liquidation is allowed if collateral ratio < 125%).
    uint256 public constant LIQUIDATION_THRESHOLD = 125;
    /// @notice The liquidation bonus (pénalité) expressed as a percentage bonus for the liquidator (e.g., 10 means +10% bonus collateral).
    uint256 public constant LIQUIDATION_BONUS = 10;

    /// @notice Mapping of user addresses to their deposited collateral amounts.
    mapping(address => uint256) public collateralBalance;

    /// @notice Constructor that sets the underlying token and the Borrowing contract.
    /// @param _token The address of the token used as collateral.
    /// @param _collateral The address of the Borrowing contract defining collateral conditions.
    constructor(address _token, address _collateral) {
        token = Token(_token);
        borrowing = Borrowing(_collateral);
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
        uint256 borrowed = borrowing.borrowedPrincipal(msg.sender);
        if (borrowed > 0) {
            require(newCollateral * 100 >= borrowed * MIN_COLLATERAL_RATIO, "Collateral ratio too low after withdrawal");
        }
        collateralBalance[msg.sender] = newCollateral;

        // Appel bas niveau pour transférer le collatéral à l'utilisateur
        (bool success, bytes memory returndata) =
            address(token).call(abi.encodeWithSelector(token.transfer.selector, msg.sender, amount));
        if (!success) {
            revert("Collateral withdrawal transfer failed");
        }
        if (returndata.length > 0) {
            bool transferSuccess = abi.decode(returndata, (bool));
            require(transferSuccess, "Collateral withdrawal transfer failed");
        }
    }

    /// @notice Gets the collateral ratio for a specific user.
    /// @param user The address of the user.
    /// @return The collateral ratio expressed as a percentage.
    function getCollateralRatio(address user) external view returns (uint256) {
        uint256 borrowed = borrowing.borrowedPrincipal(user);
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
        uint256 totalBorrowed = borrowing.borrowedPrincipal(user) + borrowAmount;
        if (totalBorrowed == 0) {
            return true;
        }
        return (collateralBalance[user] * 100) >= (totalBorrowed * MIN_COLLATERAL_RATIO);
    }

    /// @notice Liquidates an undercollateralized borrower's position if the collateral crash.
    /// @param borrower The address of the borrower to be liquidated.
    /// @param repayAmount The amount of debt the liquidator is willing to repay.
    /// @dev The liquidator must transfer repayAmount tokens to the Borrowing contract.
    ///      In exchange, the liquidator se voit attribuer une portion du collatéral avec un bonus.
    function liquidate(address borrower, uint256 repayAmount) external {
        require(repayAmount > 0, "Repay amount must be greater than zero");
        uint256 borrowed = borrowing.borrowedPrincipal(borrower);
        require(borrowed > 0, "Borrower has no debt");

        uint256 userCollateral = collateralBalance[borrower];
        uint256 currentRatio = (userCollateral * 100) / borrowed;
        require(currentRatio < LIQUIDATION_THRESHOLD, "Collateral ratio is sufficient for liquidation");

        require(token.transferFrom(msg.sender, address(borrowing), repayAmount), "Repayment transfer failed");

        uint256 collateralToSeize = (repayAmount * (100 + LIQUIDATION_BONUS)) / 100;
        if (collateralToSeize > userCollateral) {
            collateralToSeize = userCollateral;
        }

        borrowing.reduceDebt(borrower, repayAmount);
        collateralBalance[borrower] -= collateralToSeize;

        require(token.transfer(msg.sender, collateralToSeize), "Collateral transfer to liquidator failed");
    }
}
