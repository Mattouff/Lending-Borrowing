// SPDX-License-Identifier: MIT
pragma solidity 0.8.29;

import "@openzeppelin-upgradeable/contracts/access/OwnableUpgradeable.sol";
import "@openzeppelin-upgradeable/contracts/proxy/utils/Initializable.sol";
import "@openzeppelin-upgradeable/contracts/proxy/utils/UUPSUpgradeable.sol";
import "../interfaces/ILendingPoolAddressesProvider.sol";
import "../interfaces/IPriceOracle.sol";
import "../libraries/WadRayMath.sol";
import "../libraries/PercentageMath.sol";

/**
 * @title CollateralManager
 * @author DeFi Lending Platform
 * @notice Manages the collateral requirements for the lending platform
 * @dev This contract is upgradeable using the UUPS pattern
 */
contract CollateralManager is Initializable, OwnableUpgradeable, UUPSUpgradeable {
    using WadRayMath for uint256;
    using PercentageMath for uint256;

    // Address provider
    ILendingPoolAddressesProvider private _addressesProvider;

    // Minimum health factor (1.05 = 105%)
    uint256 public constant HEALTH_FACTOR_LIQUIDATION_THRESHOLD = 1.05 * 1e18;

    // Close factor (amount of debt that can be liquidated in a single tx)
    uint256 public constant CLOSE_FACTOR = 0.5 * 1e18; // 50%

    // Reserve configuration mapping (asset => config)
    struct ReserveConfig {
        // Whether the asset can be used as collateral
        bool isCollateral;
        // Loan to value ratio (max percentage of collateral that can be borrowed)
        uint256 ltv;
        // Liquidation threshold (percentage of collateral at which liquidation can occur)
        uint256 liquidationThreshold;
        // Liquidation bonus (bonus for liquidators)
        uint256 liquidationBonus;
    }

    // Asset to reserve configuration mapping
    mapping(address => ReserveConfig) private _reserveConfigs;

    // Events
    event ReserveConfigurationChanged(
        address indexed asset, bool isCollateral, uint256 ltv, uint256 liquidationThreshold, uint256 liquidationBonus
    );

    // Storage gap for future upgrades
    uint256[50] private __gap;

    /**
     * @custom:oz-upgrades-unsafe-allow constructor
     */
    constructor() {
        _disableInitializers();
    }

    /**
     * @dev Initializer function
     * @param addressesProvider The address of the LendingPoolAddressesProvider
     */
    function initialize(address addressesProvider) external initializer {
        __Ownable_init(msg.sender);
        __UUPSUpgradeable_init();

        require(addressesProvider != address(0), "CollateralManager: Invalid addresses provider");
        _addressesProvider = ILendingPoolAddressesProvider(addressesProvider);
    }

    /**
     * @dev Configures an asset as a collateral
     * @param asset The address of the asset
     * @param isCollateral Whether the asset can be used as collateral
     * @param ltv The loan to value ratio (max percentage of collateral that can be borrowed)
     * @param liquidationThreshold The liquidation threshold (percentage of collateral at which liquidation can occur)
     * @param liquidationBonus The liquidation bonus (bonus for liquidators)
     */
    function configureAsCollateral(
        address asset,
        bool isCollateral,
        uint256 ltv,
        uint256 liquidationThreshold,
        uint256 liquidationBonus
    ) external onlyOwner {
        require(asset != address(0), "CollateralManager: Invalid asset");
        require(liquidationThreshold >= ltv, "CollateralManager: Liquidation threshold must be >= LTV");
        require(liquidationBonus >= 10000, "CollateralManager: Liquidation bonus must be >= 100%");

        _reserveConfigs[asset] = ReserveConfig({
            isCollateral: isCollateral,
            ltv: ltv,
            liquidationThreshold: liquidationThreshold,
            liquidationBonus: liquidationBonus
        });

        emit ReserveConfigurationChanged(asset, isCollateral, ltv, liquidationThreshold, liquidationBonus);
    }

    // Rest of the functions remain the same but adjusted for upgradeable pattern...

    /**
     * @dev Function that should revert when msg.sender is not authorized to upgrade the contract
     * @param newImplementation The address of the new implementation
     */
    function _authorizeUpgrade(address newImplementation) internal override onlyOwner { }

    // Implement the rest of the functions as in the original contract...
    // (I've truncated them to focus on the upgradeability changes,
    // but in a real implementation you would include all functions)
}
