// SPDX-License-Identifier: MIT
pragma solidity 0.8.29;

import "forge-std/Script.sol";
import "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";
import "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";

import "../src/core/LendingPoolAddressesProvider.sol";
import "../src/core/LendingPool.sol";
import "../src/core/LendingPoolCore.sol";
import "../src/core/LendingPoolConfigurator.sol";
import "../src/core/LendingPoolDataProvider.sol";
import "../src/core/LendingPoolParametersProvider.sol";
import "../src/oracles/PriceOracle.sol";
import "../src/risk/CollateralManager.sol";
import "../src/risk/LiquidationManager.sol";
import "../src/core/InterestRateStrategy.sol";
import "../src/tokens/PlatformToken.sol";

/**
 * @title Deploy
 * @author DeFi Lending Platform
 * @notice Deployment script for the lending platform using upgradeable contracts
 */
contract Deploy is Script {
    // ProxyAdmin to manage proxies
    ProxyAdmin public proxyAdmin;

    // Implementation contracts
    LendingPoolAddressesProvider public addressesProviderImpl;
    LendingPool public lendingPoolImpl;
    LendingPoolCore public lendingPoolCoreImpl;
    CollateralManager public collateralManagerImpl;
    LiquidationManager public liquidationManagerImpl;

    // Proxy contracts
    TransparentUpgradeableProxy public addressesProviderProxy;
    TransparentUpgradeableProxy public lendingPoolProxy;
    TransparentUpgradeableProxy public lendingPoolCoreProxy;
    TransparentUpgradeableProxy public collateralManagerProxy;
    TransparentUpgradeableProxy public liquidationManagerProxy;

    // Non-upgradeable contracts
    PriceOracle public priceOracle;
    LendingPoolConfigurator public lendingPoolConfigurator;
    LendingPoolDataProvider public lendingPoolDataProvider;
    LendingPoolParametersProvider public lendingPoolParametersProvider;
    InterestRateStrategy public interestRateStrategy;
    PlatformToken public platformToken;

    function run() external {
        // Start broadcast to record and send transactions
        vm.startBroadcast();

        address admin = msg.sender;
        address ethUsdPriceFeed = 0x694AA1769357215DE4FAC081bf1f309aDC325306; // Sepolia ETH/USD Feed

        // Deploy ProxyAdmin
        proxyAdmin = new ProxyAdmin(admin);

        // Deploy implementations
        addressesProviderImpl = new LendingPoolAddressesProvider();
        lendingPoolImpl = new LendingPool();
        lendingPoolCoreImpl = new LendingPoolCore();
        collateralManagerImpl = new CollateralManager();
        liquidationManagerImpl = new LiquidationManager();

        // Deploy proxies

        // 1. AddressesProvider Proxy
        bytes memory addressesProviderData = abi.encodeWithSignature("initialize(address)", admin);
        addressesProviderProxy =
            new TransparentUpgradeableProxy(address(addressesProviderImpl), address(proxyAdmin), addressesProviderData);

        // Get the proxied AddressesProvider
        LendingPoolAddressesProvider addressesProvider = LendingPoolAddressesProvider(address(addressesProviderProxy));

        // Deploy PriceOracle
        priceOracle = new PriceOracle(admin, ethUsdPriceFeed);

        // Register PriceOracle
        addressesProvider.setPriceOracle(address(priceOracle));

        // 2. LendingPoolCore Proxy
        bytes memory lendingPoolCoreData =
            abi.encodeWithSignature("initialize(address)", address(addressesProviderProxy));
        lendingPoolCoreProxy =
            new TransparentUpgradeableProxy(address(lendingPoolCoreImpl), address(proxyAdmin), lendingPoolCoreData);

        // 3. LendingPool Proxy
        bytes memory lendingPoolData = abi.encodeWithSignature("initialize(address)", address(addressesProviderProxy));
        lendingPoolProxy =
            new TransparentUpgradeableProxy(address(lendingPoolImpl), address(proxyAdmin), lendingPoolData);

        // 4. CollateralManager Proxy
        bytes memory collateralManagerData =
            abi.encodeWithSignature("initialize(address)", address(addressesProviderProxy));
        collateralManagerProxy =
            new TransparentUpgradeableProxy(address(collateralManagerImpl), address(proxyAdmin), collateralManagerData);

        // 5. LiquidationManager Proxy
        bytes memory liquidationManagerData =
            abi.encodeWithSignature("initialize(address)", address(addressesProviderProxy));
        liquidationManagerProxy = new TransparentUpgradeableProxy(
            address(liquidationManagerImpl), address(proxyAdmin), liquidationManagerData
        );

        // Deploy other contracts
        lendingPoolParametersProvider = new LendingPoolParametersProvider();
        lendingPoolParametersProvider.transferOwnership(admin);

        lendingPoolConfigurator = new LendingPoolConfigurator(address(addressesProviderProxy));
        lendingPoolConfigurator.transferOwnership(admin);

        lendingPoolDataProvider = new LendingPoolDataProvider(address(addressesProviderProxy));

        // Get the default parameters
        (
            uint256 slope1,
            uint256 slope2,
            uint256 optimalUtilizationRate,
            uint256 baseBorrowRate,
            ,
            uint256 elasticityFactor
        ) = lendingPoolParametersProvider.getDefaultInterestRateParameters();

        // Deploy InterestRateStrategy
        interestRateStrategy =
            new InterestRateStrategy(baseBorrowRate, slope1, slope2, optimalUtilizationRate, elasticityFactor);
        interestRateStrategy.transferOwnership(admin);

        // Deploy PlatformToken
        platformToken = new PlatformToken(
            "DeFi Lending Platform Token",
            "DLPT",
            10_000_000 * 1e18, // 10 million initial supply
            100_000_000 * 1e18, // 100 million max cap
            admin
        );

        // Register addresses in AddressesProvider
        addressesProvider.setLendingPoolImpl(address(lendingPoolProxy));
        addressesProvider.setLendingPoolCoreImpl(address(lendingPoolCoreProxy));
        addressesProvider.setLendingPoolConfiguratorImpl(address(lendingPoolConfigurator));
        addressesProvider.setLendingPoolDataProviderImpl(address(lendingPoolDataProvider));
        addressesProvider.setLendingPoolParametersProviderImpl(address(lendingPoolParametersProvider));
        addressesProvider.setCollateralManagerImpl(address(collateralManagerProxy));
        addressesProvider.setLiquidationManagerImpl(address(liquidationManagerProxy));

        // Transfer ownership of proxied contracts
        LendingPoolCore(address(lendingPoolCoreProxy)).transferOwnership(admin);
        CollateralManager(address(collateralManagerProxy)).transferOwnership(admin);
        LiquidationManager(address(liquidationManagerProxy)).transferOwnership(admin);

        vm.stopBroadcast();
    }
}
