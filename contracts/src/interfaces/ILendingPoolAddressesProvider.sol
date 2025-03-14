// SPDX-License-Identifier: MIT
pragma solidity 0.8.29;

/**
 * @title ILendingPoolAddressesProvider
 * @author DeFi Lending Platform
 * @notice Defines the basic interface for a LendingPool addresses provider
 */
interface ILendingPoolAddressesProvider {
    /**
     * @dev Emitted when the lending pool address is updated
     * @param oldAddress The old address
     * @param newAddress The new address
     */
    event LendingPoolUpdated(address indexed oldAddress, address indexed newAddress);

    /**
     * @dev Emitted when the lending pool core address is updated
     * @param oldAddress The old address
     * @param newAddress The new address
     */
    event LendingPoolCoreUpdated(address indexed oldAddress, address indexed newAddress);

    /**
     * @dev Emitted when the price oracle address is updated
     * @param oldAddress The old address
     * @param newAddress The new address
     */
    event PriceOracleUpdated(address indexed oldAddress, address indexed newAddress);

    /**
     * @dev Emitted when the lending pool configurator address is updated
     * @param oldAddress The old address
     * @param newAddress The new address
     */
    event LendingPoolConfiguratorUpdated(address indexed oldAddress, address indexed newAddress);

    /**
     * @dev Emitted when the lending pool data provider address is updated
     * @param oldAddress The old address
     * @param newAddress The new address
     */
    event LendingPoolDataProviderUpdated(address indexed oldAddress, address indexed newAddress);

    /**
     * @dev Emitted when the lending pool parameters provider address is updated
     * @param oldAddress The old address
     * @param newAddress The new address
     */
    event LendingPoolParametersProviderUpdated(address indexed oldAddress, address indexed newAddress);

    /**
     * @dev Emitted when the collateral manager address is updated
     * @param oldAddress The old address
     * @param newAddress The new address
     */
    event CollateralManagerUpdated(address indexed oldAddress, address indexed newAddress);

    /**
     * @dev Emitted when the liquidation manager address is updated
     * @param oldAddress The old address
     * @param newAddress The new address
     */
    event LiquidationManagerUpdated(address indexed oldAddress, address indexed newAddress);

    /**
     * @dev Returns the address of the LendingPool proxy
     * @return The LendingPool proxy address
     */
    function getLendingPool() external view returns (address);

    /**
     * @dev Updates the address of the LendingPool
     * @param _pool The new address of the LendingPool
     */
    function setLendingPoolImpl(address _pool) external;

    /**
     * @dev Returns the address of the LendingPoolCore proxy
     * @return The LendingPoolCore proxy address
     */
    function getLendingPoolCore() external view returns (address);

    /**
     * @dev Updates the address of the LendingPoolCore
     * @param _lendingPoolCore The new address of the LendingPoolCore
     */
    function setLendingPoolCoreImpl(address _lendingPoolCore) external;

    /**
     * @dev Returns the address of the price oracle
     * @return The address of the PriceOracle
     */
    function getPriceOracle() external view returns (address);

    /**
     * @dev Updates the address of the price oracle
     * @param _priceOracle The new address of the PriceOracle
     */
    function setPriceOracle(address _priceOracle) external;

    /**
     * @dev Returns the address of the LendingPoolConfigurator proxy
     * @return The LendingPoolConfigurator proxy address
     */
    function getLendingPoolConfigurator() external view returns (address);

    /**
     * @dev Updates the address of the LendingPoolConfigurator
     * @param _configurator The new address of the LendingPoolConfigurator
     */
    function setLendingPoolConfiguratorImpl(address _configurator) external;

    /**
     * @dev Returns the address of the LendingPoolDataProvider proxy
     * @return The LendingPoolDataProvider proxy address
     */
    function getLendingPoolDataProvider() external view returns (address);

    /**
     * @dev Updates the address of the LendingPoolDataProvider
     * @param _provider The new address of the LendingPoolDataProvider
     */
    function setLendingPoolDataProviderImpl(address _provider) external;

    /**
     * @dev Returns the address of the LendingPoolParametersProvider proxy
     * @return The LendingPoolParametersProvider proxy address
     */
    function getLendingPoolParametersProvider() external view returns (address);

    /**
     * @dev Updates the address of the LendingPoolParametersProvider
     * @param _parametersProvider The new address of the LendingPoolParametersProvider
     */
    function setLendingPoolParametersProviderImpl(address _parametersProvider) external;

    /**
     * @dev Returns the address of the CollateralManager
     * @return The CollateralManager address
     */
    function getCollateralManager() external view returns (address);

    /**
     * @dev Updates the address of the CollateralManager
     * @param _collateralManager The new address of the CollateralManager
     */
    function setCollateralManagerImpl(address _collateralManager) external;

    /**
     * @dev Returns the address of the LiquidationManager
     * @return The LiquidationManager address
     */
    function getLiquidationManager() external view returns (address);

    /**
     * @dev Updates the address of the LiquidationManager
     * @param _liquidationManager The new address of the LiquidationManager
     */
    function setLiquidationManagerImpl(address _liquidationManager) external;
}
