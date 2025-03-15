// SPDX-License-Identifier: MIT
pragma solidity 0.8.29;

import "@openzeppelin/contracts/access/Ownable.sol";
import "../../src/interfaces/ILendingPoolAddressesProvider.sol";

/**
 * @title MockLendingPoolAddressesProvider
 * @author DeFi Lending Platform
 * @notice Mock Lending Pool Addresses Provider for testing purposes
 * @dev Implements the ILendingPoolAddressesProvider interface
 */
contract MockLendingPoolAddressesProvider is ILendingPoolAddressesProvider, Ownable {
    // Core protocol addresses
    address private _lendingPool;
    address private _lendingPoolCore;
    address private _lendingPoolConfigurator;
    address private _lendingPoolDataProvider;
    address private _lendingPoolParametersProvider;
    address private _priceOracle;
    address private _collateralManager;
    address private _liquidationManager;
    
    // Flag to control whether operations revert
    bool private _shouldRevert;

    /**
     * @dev Constructor to create a new MockLendingPoolAddressesProvider
     * @param owner The owner of the provider
     */
    constructor(address owner) Ownable(msg.sender) {
        // Transfer ownership
        if (owner != msg.sender) {
            transferOwnership(owner);
        }
    }

    /**
     * @inheritdoc ILendingPoolAddressesProvider
     */
    function getLendingPool() external view override returns (address) {
        if (_shouldRevert) {
            revert("MockLendingPoolAddressesProvider: Forced failure");
        }
        
        return _lendingPool;
    }

    /**
     * @inheritdoc ILendingPoolAddressesProvider
     */
    function setLendingPoolImpl(address pool) external override onlyOwner {
        address oldAddress = _lendingPool;
        _lendingPool = pool;
        emit LendingPoolUpdated(oldAddress, pool);
    }

    /**
     * @inheritdoc ILendingPoolAddressesProvider
     */
    function getLendingPoolCore() external view override returns (address) {
        if (_shouldRevert) {
            revert("MockLendingPoolAddressesProvider: Forced failure");
        }
        
        return _lendingPoolCore;
    }

    /**
     * @inheritdoc ILendingPoolAddressesProvider
     */
    function setLendingPoolCoreImpl(address lendingPoolCore) external override onlyOwner {
        address oldAddress = _lendingPoolCore;
        _lendingPoolCore = lendingPoolCore;
        emit LendingPoolCoreUpdated(oldAddress, lendingPoolCore);
    }

    /**
     * @inheritdoc ILendingPoolAddressesProvider
     */
    function getPriceOracle() external view override returns (address) {
        if (_shouldRevert) {
            revert("MockLendingPoolAddressesProvider: Forced failure");
        }
        
        return _priceOracle;
    }

    /**
     * @inheritdoc ILendingPoolAddressesProvider
     */
    function setPriceOracle(address priceOracle) external override onlyOwner {
        address oldAddress = _priceOracle;
        _priceOracle = priceOracle;
        emit PriceOracleUpdated(oldAddress, priceOracle);
    }

    /**
     * @inheritdoc ILendingPoolAddressesProvider
     */
    function getLendingPoolConfigurator() external view override returns (address) {
        if (_shouldRevert) {
            revert("MockLendingPoolAddressesProvider: Forced failure");
        }
        
        return _lendingPoolConfigurator;
    }

    /**
     * @inheritdoc ILendingPoolAddressesProvider
     */
    function setLendingPoolConfiguratorImpl(address configurator) external override onlyOwner {
        address oldAddress = _lendingPoolConfigurator;
        _lendingPoolConfigurator = configurator;
        emit LendingPoolConfiguratorUpdated(oldAddress, configurator);
    }

    /**
     * @inheritdoc ILendingPoolAddressesProvider
     */
    function getLendingPoolDataProvider() external view override returns (address) {
        if (_shouldRevert) {
            revert("MockLendingPoolAddressesProvider: Forced failure");
        }
        
        return _lendingPoolDataProvider;
    }

    /**
     * @inheritdoc ILendingPoolAddressesProvider
     */
    function setLendingPoolDataProviderImpl(address provider) external override onlyOwner {
        address oldAddress = _lendingPoolDataProvider;
        _lendingPoolDataProvider = provider;
        emit LendingPoolDataProviderUpdated(oldAddress, provider);
    }

    /**
     * @inheritdoc ILendingPoolAddressesProvider
     */
    function getLendingPoolParametersProvider() external view override returns (address) {
        if (_shouldRevert) {
            revert("MockLendingPoolAddressesProvider: Forced failure");
        }
        
        return _lendingPoolParametersProvider;
    }

    /**
     * @inheritdoc ILendingPoolAddressesProvider
     */
    function setLendingPoolParametersProviderImpl(address parametersProvider) external override onlyOwner {
        address oldAddress = _lendingPoolParametersProvider;
        _lendingPoolParametersProvider = parametersProvider;
        emit LendingPoolParametersProviderUpdated(oldAddress, parametersProvider);
    }

    /**
     * @inheritdoc ILendingPoolAddressesProvider
     */
    function getCollateralManager() external view override returns (address) {
        if (_shouldRevert) {
            revert("MockLendingPoolAddressesProvider: Forced failure");
        }
        
        return _collateralManager;
    }

    /**
     * @inheritdoc ILendingPoolAddressesProvider
     */
    function setCollateralManagerImpl(address collateralManager) external override onlyOwner {
        address oldAddress = _collateralManager;
        _collateralManager = collateralManager;
        emit CollateralManagerUpdated(oldAddress, collateralManager);
    }

    /**
     * @inheritdoc ILendingPoolAddressesProvider
     */
    function getLiquidationManager() external view override returns (address) {
        if (_shouldRevert) {
            revert("MockLendingPoolAddressesProvider: Forced failure");
        }
        
        return _liquidationManager;
    }

    /**
     * @inheritdoc ILendingPoolAddressesProvider
     */
    function setLiquidationManagerImpl(address liquidationManager) external override onlyOwner {
        address oldAddress = _liquidationManager;
        _liquidationManager = liquidationManager;
        emit LiquidationManagerUpdated(oldAddress, liquidationManager);
    }

    /**
     * @dev Sets whether operations should revert
     * @param shouldRevert Whether operations should revert
     */
    function setShouldRevert(bool shouldRevert) external onlyOwner {
        _shouldRevert = shouldRevert;
    }
}
