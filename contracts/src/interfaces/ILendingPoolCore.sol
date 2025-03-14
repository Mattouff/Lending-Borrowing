// SPDX-License-Identifier: MIT
pragma solidity 0.8.29;

import "../libraries/ReserveLogic.sol";

/**
 * @title ILendingPoolCore
 * @author DeFi Lending Platform
 * @notice Interface for the LendingPoolCore contract
 */
interface ILendingPoolCore {
    /**
     * @dev Initializes a reserve
     * @param asset The address of the asset
     * @param aToken The address of the aToken
     * @param debtToken The address of the debt token
     * @param interestRateStrategy The address of the interest rate strategy
     */
    function initReserve(address asset, address aToken, address debtToken, address interestRateStrategy) external;

    /**
     * @dev Updates the interest rate strategy of a reserve
     * @param asset The address of the asset
     * @param interestRateStrategy The address of the interest rate strategy
     */
    function updateInterestRateStrategy(address asset, address interestRateStrategy) external;

    /**
     * @dev Updates the reserve state with accumulated interest
     * @param asset The address of the asset
     */
    function updateReserveState(address asset) external;

    /**
     * @dev Transfers an amount of asset to the user
     * @param asset The address of the asset
     * @param to The address of the recipient
     * @param amount The amount to transfer
     */
    function transferToUser(address asset, address to, uint256 amount) external;

    /**
     * @dev Transfers an amount of asset from the user to the reserve
     * @param asset The address of the asset
     * @param from The address of the sender
     * @param amount The amount to transfer
     */
    function transferToReserve(address asset, address from, uint256 amount) external;

    /**
     * @dev Gets the reserve data for an asset
     * @param asset The address of the asset
     * @return The reserve data
     */
    function getReserveData(address asset) external view returns (ReserveLogic.ReserveData memory);

    /**
     * @dev Gets the list of all reserves
     * @return The list of all reserves
     */
    function getReservesList() external view returns (address[] memory);

    /**
     * @dev Enables or disables borrowing on a reserve
     * @param asset The address of the asset
     * @param enabled Whether borrowing is enabled
     */
    function enableBorrowingOnReserve(address asset, bool enabled) external;
}
