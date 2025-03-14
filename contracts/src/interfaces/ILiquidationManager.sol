// SPDX-License-Identifier: MIT
pragma solidity 0.8.29;

/**
 * @title ILiquidationManager
 * @author DeFi Lending Platform
 * @notice Interface for the LiquidationManager contract
 */
interface ILiquidationManager {
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
     * @dev Calculates the liquidation amount
     * @param debtToCover The amount of debt to cover
     * @param collateralAsset The address of the collateral asset
     * @param debtAsset The address of the debt asset
     * @return The amount of collateral to be liquidated
     */
    function calculateLiquidationAmount(uint256 debtToCover, address collateralAsset, address debtAsset)
        external
        view
        returns (uint256);
}
