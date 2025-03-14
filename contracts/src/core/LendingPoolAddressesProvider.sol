// SPDX-License-Identifier: MIT
pragma solidity 0.8.29;

import "@openzeppelin-upgradeable/contracts/access/OwnableUpgradeable.sol";
import "@openzeppelin-upgradeable/contracts/proxy/utils/Initializable.sol";
import "@openzeppelin-upgradeable/contracts/proxy/utils/UUPSUpgradeable.sol";
import "../interfaces/ILendingPoolAddressesProvider.sol";

/**
 * @title LendingPoolAddressesProvider
 * @author DeFi Lending Platform
 * @notice Registry for all the addresses used in the protocol
 * @dev Acts as a service registry, allowing for upgradability of the core components
 * This contract is upgradeable using the UUPS pattern
 */
contract LendingPoolAddressesProvider is
    Initializable,
    ILendingPoolAddressesProvider,
    OwnableUpgradeable,
    UUPSUpgradeable
{
    // Main addresses
    address private _lendingPool;
    address private _lendingPoolCore;
    address private _lendingPoolConfigurator;
    address private _lendingPoolDataProvider;
    address private _lendingPoolParametersProvider;
    address private _priceOracle;
    address private _collateralManager;
    address private _liquidationManager;

    // Storage gap for future upgrades
    uint256[48] private __gap;

    /**
     * @custom:oz-upgrades-unsafe-allow constructor
     */
    constructor() {
        _disableInitializers();
    }

    /**
     * @dev Initializer function
     * @param owner The owner address
     */
    function initialize(address owner) external initializer {
        __Ownable_init(owner);
        __UUPSUpgradeable_init();
    }

    /**
     * @inheritdoc ILendingPoolAddressesProvider
     */
    function getLendingPool() external view override returns (address) {
        return _lendingPool;
    }

    /**
     * @inheritdoc ILendingPoolAddressesProvider
     */
    function setLendingPoolImpl(address _pool) external override onlyOwner {
        address oldAddress = _lendingPool;
        _lendingPool = _pool;
        emit LendingPoolUpdated(oldAddress, _pool);
    }

    /**
     * @inheritdoc ILendingPoolAddressesProvider
     */
    function getLendingPoolCore() external view override returns (address) {
        return _lendingPoolCore;
    }

    /**
     * @inheritdoc ILendingPoolAddressesProvider
     */
    function setLendingPoolCoreImpl(address _lendingPoolCore_) external override onlyOwner {
        address oldAddress = _lendingPoolCore;
        _lendingPoolCore = _lendingPoolCore_;
        emit LendingPoolCoreUpdated(oldAddress, _lendingPoolCore_);
    }

    /**
     * @inheritdoc ILendingPoolAddressesProvider
     */
    function getLendingPoolConfigurator() external view override returns (address) {
        return _lendingPoolConfigurator;
    }

    /**
     * @inheritdoc ILendingPoolAddressesProvider
     */
    function setLendingPoolConfiguratorImpl(address _configurator) external override onlyOwner {
        address oldAddress = _lendingPoolConfigurator;
        _lendingPoolConfigurator = _configurator;
        emit LendingPoolConfiguratorUpdated(oldAddress, _configurator);
    }

    /**
     * @inheritdoc ILendingPoolAddressesProvider
     */
    function getLendingPoolDataProvider() external view override returns (address) {
        return _lendingPoolDataProvider;
    }

    /**
     * @inheritdoc ILendingPoolAddressesProvider
     */
    function setLendingPoolDataProviderImpl(address _provider) external override onlyOwner {
        address oldAddress = _lendingPoolDataProvider;
        _lendingPoolDataProvider = _provider;
        emit LendingPoolDataProviderUpdated(oldAddress, _provider);
    }

    /**
     * @inheritdoc ILendingPoolAddressesProvider
     */
    function getLendingPoolParametersProvider() external view override returns (address) {
        return _lendingPoolParametersProvider;
    }

    /**
     * @inheritdoc ILendingPoolAddressesProvider
     */
    function setLendingPoolParametersProviderImpl(address _parametersProvider) external override onlyOwner {
        address oldAddress = _lendingPoolParametersProvider;
        _lendingPoolParametersProvider = _parametersProvider;
        emit LendingPoolParametersProviderUpdated(oldAddress, _parametersProvider);
    }

    /**
     * @inheritdoc ILendingPoolAddressesProvider
     */
    function getPriceOracle() external view override returns (address) {
        return _priceOracle;
    }

    /**
     * @inheritdoc ILendingPoolAddressesProvider
     */
    function setPriceOracle(address _priceOracle_) external override onlyOwner {
        address oldAddress = _priceOracle;
        _priceOracle = _priceOracle_;
        emit PriceOracleUpdated(oldAddress, _priceOracle_);
    }

    /**
     * @inheritdoc ILendingPoolAddressesProvider
     */
    function getCollateralManager() external view override returns (address) {
        return _collateralManager;
    }

    /**
     * @inheritdoc ILendingPoolAddressesProvider
     */
    function setCollateralManagerImpl(address _collateralManager_) external override onlyOwner {
        address oldAddress = _collateralManager;
        _collateralManager = _collateralManager_;
        emit CollateralManagerUpdated(oldAddress, _collateralManager_);
    }

    /**
     * @inheritdoc ILendingPoolAddressesProvider
     */
    function getLiquidationManager() external view override returns (address) {
        return _liquidationManager;
    }

    /**
     * @inheritdoc ILendingPoolAddressesProvider
     */
    function setLiquidationManagerImpl(address _liquidationManager_) external override onlyOwner {
        address oldAddress = _liquidationManager;
        _liquidationManager = _liquidationManager_;
        emit LiquidationManagerUpdated(oldAddress, _liquidationManager_);
    }

    /**
     * @dev Function that should revert when msg.sender is not authorized to upgrade the contract
     * @param newImplementation The address of the new implementation
     */
    function _authorizeUpgrade(address newImplementation) internal override onlyOwner { }
}
