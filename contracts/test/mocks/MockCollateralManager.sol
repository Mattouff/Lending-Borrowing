// SPDX-License-Identifier: MIT
pragma solidity 0.8.29;

import "@openzeppelin/contracts/access/Ownable.sol";
import "../../src/interfaces/ICollateralManager.sol";

/**
 * @title MockCollateralManager
 * @author DeFi Lending Platform
 * @notice Mock Collateral Manager for testing purposes
 * @dev Implements the ICollateralManager interface
 */
contract MockCollateralManager is ICollateralManager, Ownable {
    // Reserve configuration mapping (asset => config)
    struct ReserveConfig {
        bool isCollateral;
        uint256 ltv;
        uint256 liquidationThreshold;
        uint256 liquidationBonus;
    }

    // Asset to reserve configuration mapping
    mapping(address => ReserveConfig) private _reserveConfigs;

    // Flag to control whether operations revert
    bool private _shouldRevert;

    // Flags to control validation results
    bool private _validateBorrowShouldPass;
    bool private _isLiquidationValidShouldPass;
    uint256 private _mockHealthFactor;

    // Track calls for testing verification
    uint256 public configureAsCollateralCalls;
    uint256 public validateBorrowCalls;
    uint256 public isLiquidationValidCalls;

    // Events
    event ReserveConfigurationChanged(
        address indexed asset, bool isCollateral, uint256 ltv, uint256 liquidationThreshold, uint256 liquidationBonus
    );

    /**
     * @dev Constructor to create a new MockCollateralManager
     * @param owner The owner of the contract
     */
    constructor(address owner) Ownable(msg.sender) {
        // Initialize defaults
        _validateBorrowShouldPass = true;
        _isLiquidationValidShouldPass = false;
        _mockHealthFactor = 2 * 1e18; // 2.0, healthy

        // Transfer ownership
        if (owner != msg.sender) {
            transferOwnership(owner);
        }
    }

    /**
     * @inheritdoc ICollateralManager
     */
    function configureAsCollateral(
        address asset,
        bool isCollateral,
        uint256 ltv,
        uint256 liquidationThreshold,
        uint256 liquidationBonus
    ) external override onlyOwner {
        if (_shouldRevert) {
            revert("MockCollateralManager: Forced failure");
        }

        configureAsCollateralCalls++;

        require(asset != address(0), "MockCollateralManager: Zero asset");

        _reserveConfigs[asset] = ReserveConfig({
            isCollateral: isCollateral,
            ltv: ltv,
            liquidationThreshold: liquidationThreshold,
            liquidationBonus: liquidationBonus
        });

        emit ReserveConfigurationChanged(asset, isCollateral, ltv, liquidationThreshold, liquidationBonus);
    }

    /**
     * @inheritdoc ICollateralManager
     */
    function getReserveConfig(address asset)
        external
        view
        override
        returns (bool isCollateral, uint256 ltv, uint256 liquidationThreshold, uint256 liquidationBonus)
    {
        if (_shouldRevert) {
            revert("MockCollateralManager: Forced failure");
        }

        ReserveConfig memory config = _reserveConfigs[asset];
        return (config.isCollateral, config.ltv, config.liquidationThreshold, config.liquidationBonus);
    }

    /**
     * @inheritdoc ICollateralManager
     */
    function validateBorrow(address user, address asset, uint256 amount) external view override returns (bool) {
        if (_shouldRevert) {
            revert("MockCollateralManager: Forced failure");
        }

        // Increment would make this function non-view, so we don't track this call count

        // Return the configured validation result
        return _validateBorrowShouldPass;
    }

    /**
     * @inheritdoc ICollateralManager
     */
    function isLiquidationValid(address user) external view override returns (bool, uint256) {
        if (_shouldRevert) {
            revert("MockCollateralManager: Forced failure");
        }

        // Increment would make this function non-view, so we don't track this call count

        // Return the configured validation result
        return (_isLiquidationValidShouldPass, _mockHealthFactor);
    }

    /**
     * @dev Sets the validation result for validateBorrow
     * @param shouldPass Whether validateBorrow should pass
     */
    function setValidateBorrowResult(bool shouldPass) external onlyOwner {
        _validateBorrowShouldPass = shouldPass;
    }

    /**
     * @dev Sets the validation result for isLiquidationValid
     * @param shouldPass Whether isLiquidationValid should pass
     * @param healthFactor The health factor to return
     */
    function setIsLiquidationValidResult(bool shouldPass, uint256 healthFactor) external onlyOwner {
        _isLiquidationValidShouldPass = shouldPass;
        _mockHealthFactor = healthFactor;
    }

    /**
     * @dev Sets a reserve's configuration directly
     * @param asset The asset address
     * @param isCollateral Whether the asset can be used as collateral
     * @param ltv The loan to value ratio
     * @param liquidationThreshold The liquidation threshold
     * @param liquidationBonus The liquidation bonus
     */
    function setReserveConfigDirect(
        address asset,
        bool isCollateral,
        uint256 ltv,
        uint256 liquidationThreshold,
        uint256 liquidationBonus
    ) external onlyOwner {
        _reserveConfigs[asset] = ReserveConfig({
            isCollateral: isCollateral,
            ltv: ltv,
            liquidationThreshold: liquidationThreshold,
            liquidationBonus: liquidationBonus
        });
    }

    /**
     * @dev Sets whether operations should revert
     * @param shouldRevert Whether operations should revert
     */
    function setShouldRevert(bool shouldRevert) external onlyOwner {
        _shouldRevert = shouldRevert;
    }

    /**
     * @dev Resets the call counters for testing
     */
    function resetCallCounters() external onlyOwner {
        configureAsCollateralCalls = 0;
        validateBorrowCalls = 0;
        isLiquidationValidCalls = 0;
    }

    /**
     * @dev Tracks validateBorrow calls (for use in test functions)
     */
    function trackValidateBorrowCall() external onlyOwner {
        validateBorrowCalls++;
    }

    /**
     * @dev Tracks isLiquidationValid calls (for use in test functions)
     */
    function trackIsLiquidationValidCall() external onlyOwner {
        isLiquidationValidCalls++;
    }
}
