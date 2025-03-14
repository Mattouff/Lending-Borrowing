// SPDX-License-Identifier: MIT
pragma solidity 0.8.29;

/**
 * @title ILendingPool
 * @author DeFi Lending Platform
 * @notice Interface for the LendingPool contract
 */
interface ILendingPool {
    /**
     * @dev Emitted when a deposit is made
     * @param user The address of the user
     * @param asset The address of the asset
     * @param amount The amount deposited
     * @param referral The referral code
     */
    event Deposit(address indexed user, address indexed asset, uint256 amount, uint16 indexed referral);

    /**
     * @dev Emitted when a withdrawal is made
     * @param user The address of the user
     * @param asset The address of the asset
     * @param amount The amount withdrawn
     */
    event Withdraw(address indexed user, address indexed asset, uint256 amount);

    /**
     * @dev Emitted when a borrow is made
     * @param user The address of the user
     * @param asset The address of the asset
     * @param amount The amount borrowed
     * @param interestRateMode The interest rate mode (variable/stable)
     * @param borrowRate The borrow rate
     * @param referral The referral code
     */
    event Borrow(
        address indexed user,
        address indexed asset,
        uint256 amount,
        uint256 interestRateMode,
        uint256 borrowRate,
        uint16 indexed referral
    );

    /**
     * @dev Emitted when a repayment is made
     * @param user The address of the user
     * @param asset The address of the asset
     * @param amount The amount repaid
     * @param interestPaid The interest paid
     */
    event Repay(address indexed user, address indexed asset, uint256 amount, uint256 interestPaid);

    /**
     * @dev Emitted when a liquidation is performed
     * @param liquidator The address of the liquidator
     * @param user The address of the user being liquidated
     * @param asset The address of the asset being liquidated
     * @param amount The amount liquidated
     * @param liquidationBonus The liquidation bonus
     */
    event Liquidation(
        address indexed liquidator,
        address indexed user,
        address indexed asset,
        uint256 amount,
        uint256 liquidationBonus
    );

    /**
     * @dev Deposits an `amount` of assets into the lending pool
     * @param asset The address of the asset to deposit
     * @param amount The amount to deposit
     * @param referralCode The referral code
     */
    function deposit(address asset, uint256 amount, uint16 referralCode) external;

    /**
     * @dev Withdraws an `amount` of assets from the lending pool
     * @param asset The address of the asset to withdraw
     * @param amount The amount to withdraw
     * @return The actual amount withdrawn
     */
    function withdraw(address asset, uint256 amount) external returns (uint256);

    /**
     * @dev Borrows an `amount` of assets from the lending pool
     * @param asset The address of the asset to borrow
     * @param amount The amount to borrow
     * @param interestRateMode The interest rate mode (1 for variable, 2 for stable)
     * @param referralCode The referral code
     */
    function borrow(address asset, uint256 amount, uint256 interestRateMode, uint16 referralCode) external;

    /**
     * @dev Repays an `amount` of borrowed assets
     * @param asset The address of the asset to repay
     * @param amount The amount to repay
     * @param interestRateMode The interest rate mode (1 for variable, 2 for stable)
     * @return The actual amount repaid
     */
    function repay(address asset, uint256 amount, uint256 interestRateMode) external returns (uint256);

    /**
     * @dev Liquidates an undercollateralized position
     * @param collateralAsset The address of the collateral asset
     * @param debtAsset The address of the debt asset
     * @param user The address of the user
     * @param debtToCover The debt amount to cover
     * @return The amount of collateral liquidated
     */
    function liquidationCall(address collateralAsset, address debtAsset, address user, uint256 debtToCover)
        external
        returns (uint256);

    /**
     * @dev Returns the user account data across all reserves
     * @param user The address of the user
     * @return totalCollateralETH The total collateral in ETH
     * @return totalDebtETH The total debt in ETH
     * @return availableBorrowsETH The available borrows in ETH
     * @return currentLiquidationThreshold The current liquidation threshold
     * @return ltv The loan to value ratio
     * @return healthFactor The health factor
     */
    function getUserAccountData(address user)
        external
        view
        returns (
            uint256 totalCollateralETH,
            uint256 totalDebtETH,
            uint256 availableBorrowsETH,
            uint256 currentLiquidationThreshold,
            uint256 ltv,
            uint256 healthFactor
        );
}
