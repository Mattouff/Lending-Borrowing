// SPDX-License-Identifier: MIT
pragma solidity 0.8.29;

/**
 * @title ICollateralManager
 * @author DeFi Lending Platform
 * @notice Interface for the CollateralManager contract
 */
interface ICollateralManager {
    /**
     * @dev Configures an asset as a collateral
     * @param asset The address of the asset
     * @param isCollateral Whether the asset can be used as collateral
     * @param ltv The loan to value ratio
     * @param liquidationThreshold The liquidation threshold
     * @param liquidationBonus The liquidation bonus
     */
    function configureAsCollateral(
        address asset,
        bool isCollateral,
        uint256 ltv,
        uint256 liquidationThreshold,
        uint256 liquidationBonus
    ) external;

    /**
     * @dev Gets the reserve configuration
     * @param asset The address of the asset
     * @return isCollateral Whether the asset can be used as collateral
     * @return ltv The loan to value ratio
     * @return liquidationThreshold The liquidation threshold
     * @return liquidationBonus The liquidation bonus
     */
    function getReserveConfig(address asset)
        external
        view
        returns (bool isCollateral, uint256 ltv, uint256 liquidationThreshold, uint256 liquidationBonus);

    /**
     * @dev Validates if a borrow is allowed
     * @param user The address of the user
     * @param asset The address of the asset
     * @param amount The amount to borrow
     * @return Whether the borrow is valid
     */
    function validateBorrow(address user, address asset, uint256 amount) external view returns (bool);

    /**
     * @dev Checks if a liquidation is valid
     * @param user The address of the user
     * @return Whether the liquidation is valid and the health factor
     */
    function isLiquidationValid(address user) external view returns (bool, uint256);
}
