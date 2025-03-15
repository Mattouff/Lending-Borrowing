// SPDX-License-Identifier: MIT
pragma solidity 0.8.29;

import "forge-std/Test.sol";
import "../../src/core/LendingPool.sol";
import "../../src/core/LendingPoolCore.sol";
import "../../src/core/LendingPoolAddressesProvider.sol";
import "../../src/core/LendingPoolConfigurator.sol";
import "../../src/risk/CollateralManager.sol";
import "../../src/oracles/PriceOracle.sol";
import "../../src/core/InterestRateStrategy.sol";
import "../../src/tokens/AToken.sol";
import "../../src/tokens/DebtToken.sol";
import "../helpers/BaseTest.sol";
import "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";
import "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";

/**
 * @title DepositBorrowTest
 * @author DeFi Lending Platform
 * @notice Integration tests for deposit and borrow interactions
 */
contract DepositBorrowTest is BaseTest {
    // Real contracts for testing
    LendingPoolAddressesProvider private realAddressesProvider;
    LendingPool private lendingPool;
    LendingPoolCore private lendingPoolCore;
    CollateralManager private collateralManager;
    PriceOracle private realPriceOracle;
    LendingPoolConfigurator private configurator;
    InterestRateStrategy private realStrategy;

    // Token contract addresses
    address private daiAToken;
    address private daiDebtToken;
    address private usdcAToken;
    address private usdcDebtToken;
    address private wethAToken;
    address private wethDebtToken;

    /**
     * @dev Set up the test environment extending the BaseTest setup
     */
    function setUp() public override {
        // Call parent setup to initialize mocks
        super.setUp();

        // Deploy actual contracts
        _deployContracts();

        // Initialize reserves
        _initializeReserves();

        // Transfer tokens to users
        _transferTokens();
    }

    /**
     * @dev Deploy the actual contracts
     */
    function _deployContracts() internal {
        vm.startPrank(admin);
        console.log("Deploying with admin:", admin);

        // Deploy ProxyAdmin
        ProxyAdmin proxyAdmin = new ProxyAdmin(admin);
        console.log("ProxyAdmin deployed at:", address(proxyAdmin));

        // Deploy implementations
        LendingPoolAddressesProvider addressesProviderImpl = new LendingPoolAddressesProvider();
        LendingPool lendingPoolImpl = new LendingPool();
        LendingPoolCore lendingPoolCoreImpl = new LendingPoolCore();
        CollateralManager collateralManagerImpl = new CollateralManager();

        console.log("Implementations deployed:");
        console.log("- AddressesProvider:", address(addressesProviderImpl));
        console.log("- LendingPool:", address(lendingPoolImpl));
        console.log("- LendingPoolCore:", address(lendingPoolCoreImpl));
        console.log("- CollateralManager:", address(collateralManagerImpl));

        // Create initialization data for AddressesProvider proxy
        bytes memory addressesProviderData = abi.encodeWithSignature("initialize(address)", admin);

        // Deploy AddressesProvider Proxy
        TransparentUpgradeableProxy addressesProviderProxy =
            new TransparentUpgradeableProxy(address(addressesProviderImpl), address(proxyAdmin), addressesProviderData);

        // Get the proxied AddressesProvider
        realAddressesProvider = LendingPoolAddressesProvider(address(addressesProviderProxy));
        console.log("AddressesProvider proxy deployed at:", address(realAddressesProvider));

        // Deploy price oracle
        realPriceOracle = new PriceOracle(admin, address(ethUsdFeed));
        console.log("PriceOracle deployed at:", address(realPriceOracle));

        // Set price oracle in addresses provider
        realAddressesProvider.setPriceOracle(address(realPriceOracle));

        // Create initialization data for other proxies
        bytes memory lendingPoolCoreData =
            abi.encodeWithSignature("initialize(address)", address(realAddressesProvider));

        bytes memory lendingPoolData = abi.encodeWithSignature("initialize(address)", address(realAddressesProvider));

        bytes memory collateralManagerData =
            abi.encodeWithSignature("initialize(address)", address(realAddressesProvider));

        // Deploy other proxies
        TransparentUpgradeableProxy lendingPoolCoreProxy =
            new TransparentUpgradeableProxy(address(lendingPoolCoreImpl), address(proxyAdmin), lendingPoolCoreData);

        TransparentUpgradeableProxy lendingPoolProxy =
            new TransparentUpgradeableProxy(address(lendingPoolImpl), address(proxyAdmin), lendingPoolData);

        TransparentUpgradeableProxy collateralManagerProxy =
            new TransparentUpgradeableProxy(address(collateralManagerImpl), address(proxyAdmin), collateralManagerData);

        console.log("Proxies deployed:");
        console.log("- LendingPoolCore:", address(lendingPoolCoreProxy));
        console.log("- LendingPool:", address(lendingPoolProxy));
        console.log("- CollateralManager:", address(collateralManagerProxy));

        // Set the variables to the proxied contracts
        lendingPoolCore = LendingPoolCore(address(lendingPoolCoreProxy));
        lendingPool = LendingPool(address(lendingPoolProxy));
        collateralManager = CollateralManager(address(collateralManagerProxy));

        // Deploy interest rate strategy
        realStrategy = new InterestRateStrategy(
            WAD / 100, // 1% base borrow rate
            WAD / 10, // 10% slope1
            WAD * 4 / 10, // 40% slope2
            WAD * 8 / 10, // 80% optimal utilization
            2 * WAD // Quadratic growth
        );
        console.log("InterestRateStrategy deployed at:", address(realStrategy));

        // Deploy configurator
        configurator = new LendingPoolConfigurator(address(realAddressesProvider));
        console.log("LendingPoolConfigurator deployed at:", address(configurator));

        // Check ownership before transfer
        address currentOwner = lendingPoolCore.owner();
        console.log("LendingPoolCore owner before transfer:", currentOwner);

        // Transfer ownership of LendingPoolCore to the Configurator
        lendingPoolCore.transferOwnership(address(configurator));

        // Check ownership after transfer
        currentOwner = lendingPoolCore.owner();
        console.log("LendingPoolCore owner after transfer:", currentOwner);

        // Register addresses in AddressesProvider
        console.log("Registering addresses in provider...");
        realAddressesProvider.setLendingPoolCoreImpl(address(lendingPoolCore));
        realAddressesProvider.setLendingPoolImpl(address(lendingPool));
        realAddressesProvider.setCollateralManagerImpl(address(collateralManager));
        realAddressesProvider.setLendingPoolConfiguratorImpl(address(configurator));

        // Log the addresses stored in the provider
        console.log("Stored addresses:");
        console.log("- LendingPool:", realAddressesProvider.getLendingPool());
        console.log("- LendingPoolCore:", realAddressesProvider.getLendingPoolCore());
        console.log("- Configurator:", realAddressesProvider.getLendingPoolConfigurator());

        // Set up price oracle price feeds
        MockChainlinkAggregator daiPriceFeed = new MockChainlinkAggregator(int256(1 * 1e8), 8, "DAI / USD");
        MockChainlinkAggregator usdcPriceFeed = new MockChainlinkAggregator(int256(1 * 1e8), 8, "USDC / USD");
        MockChainlinkAggregator wethPriceFeed = new MockChainlinkAggregator(int256(2000 * 1e8), 8, "WETH / USD");

        realPriceOracle.setAssetSource(address(dai), address(daiPriceFeed));
        realPriceOracle.setAssetSource(address(usdc), address(usdcPriceFeed));
        realPriceOracle.setAssetSource(address(weth), address(wethPriceFeed));

        vm.stopPrank();
    }

    /**
     * @dev Initialize reserves with tokens
     */
    function setupTestReserve(
        address asset,
        string memory aTokenName,
        string memory aTokenSymbol,
        string memory debtTokenName,
        string memory debtTokenSymbol
    ) internal returns (address aToken, address debtToken) {
        vm.startPrank(admin);

        console.log("Initializing reserve for asset:", asset);
        console.log("Using strategy:", address(realStrategy));

        // Initialize the reserve through the configurator
        configurator.initReserve(asset, aTokenName, aTokenSymbol, debtTokenName, debtTokenSymbol, address(realStrategy));

        // Get the deployed token addresses from the core
        ReserveLogic.ReserveData memory reserveData = lendingPoolCore.getReserveData(asset);

        aToken = reserveData.aTokenAddress;
        debtToken = reserveData.debtTokenAddress;

        console.log("Reserve initialized:");
        console.log("- aToken:", aToken);
        console.log("- debtToken:", debtToken);
        console.log("- borrowingEnabled:", reserveData.borrowingEnabled);

        vm.stopPrank();

        return (aToken, debtToken);
    }

    function _initializeReserves() internal {
        // Initialize reserves
        (daiAToken, daiDebtToken) = setupTestReserve(address(dai), "aDAI Token", "aDAI", "debtDAI Token", "debtDAI");

        (usdcAToken, usdcDebtToken) =
            setupTestReserve(address(usdc), "aUSDC Token", "aUSDC", "debtUSDC Token", "debtUSDC");

        (wethAToken, wethDebtToken) =
            setupTestReserve(address(weth), "aWETH Token", "aWETH", "debtWETH Token", "debtWETH");

        // Configure collateral and enable borrowing
        vm.startPrank(admin);

        // Configure DAI as collateral
        console.log("Configuring DAI as collateral");
        collateralManager.configureAsCollateral(
            address(dai),
            true, // Can be used as collateral
            75 * 1e16, // 75% LTV
            80 * 1e16, // 80% liquidation threshold
            110 * 1e16 // 110% liquidation bonus
        );

        // Configure USDC as collateral
        console.log("Configuring USDC as collateral");
        collateralManager.configureAsCollateral(
            address(usdc),
            true, // Can be used as collateral
            75 * 1e16, // 75% LTV
            80 * 1e16, // 80% liquidation threshold
            110 * 1e16 // 110% liquidation bonus
        );

        // Configure WETH as collateral
        console.log("Configuring WETH as collateral");
        collateralManager.configureAsCollateral(
            address(weth),
            true, // Can be used as collateral
            75 * 1e16, // 75% LTV
            80 * 1e16, // 80% liquidation threshold
            110 * 1e16 // 110% liquidation bonus
        );

        // Enable borrowing
        console.log("Enabling borrowing on DAI");
        console.log("Using configurator:", address(configurator));
        console.log("DAI address:", address(dai));

        // Check reserve data before enabling borrowing
        ReserveLogic.ReserveData memory reserveData = lendingPoolCore.getReserveData(address(dai));
        console.log("Reserve data before enabling borrowing:");
        console.log("- aToken:", reserveData.aTokenAddress);
        console.log("- debtToken:", reserveData.debtTokenAddress);
        console.log("- borrowingEnabled:", reserveData.borrowingEnabled);

        try configurator.enableBorrowingOnReserve(address(dai), true) {
            console.log("Borrowing successfully enabled");
        } catch Error(string memory reason) {
            console.log("Error enabling borrowing:", reason);
        } catch {
            console.log("Unknown error enabling borrowing");
        }

        console.log("Enabling borrowing on USDC");
        console.log("USDC address:", address(usdc));
        try configurator.enableBorrowingOnReserve(address(usdc), true) {
            console.log("Borrowing successfully enabled");
        } catch Error(string memory reason) {
            console.log("Error enabling borrowing:", reason);
        } catch {
            console.log("Unknown error enabling borrowing");
        }

        console.log("Enabling borrowing on WETH");
        console.log("WETH address:", address(weth));
        try configurator.enableBorrowingOnReserve(address(weth), true) {
            console.log("Borrowing successfully enabled");
        } catch Error(string memory reason) {
            console.log("Error enabling borrowing:", reason);
        } catch {
            console.log("Unknown error enabling borrowing");
        }

        vm.stopPrank();
    }

    /**
     * @dev Transfer tokens to test accounts
     */
    function _transferTokens() internal {
        // Alice gets 1000 DAI
        vm.startPrank(admin);
        dai.transfer(alice, 1000 * WAD);

        // Bob gets 1000 USDC
        usdc.transfer(bob, 1000 * 1e6);

        // Carol gets 10 WETH
        weth.transfer(carol, 10 * WAD);

        vm.stopPrank();
    }

    /**
     * @dev Test simple deposit
     */
    function testDeposit() public {
        uint256 depositAmount = 100 * WAD;

        // Alice approves and deposits DAI
        vm.startPrank(alice);
        dai.approve(address(lendingPool), depositAmount);
        lendingPool.deposit(address(dai), depositAmount, 0);
        vm.stopPrank();

        // Check aToken balance
        assertEq(IERC20(daiAToken).balanceOf(alice), depositAmount, "Alice should receive aTokens");

        // Check user account data
        (uint256 totalCollateralETH, uint256 totalDebtETH, uint256 newAvailableBorrowsETH,,, uint256 healthFactor) =
            lendingPool.getUserAccountData(alice);

        // Calculate expected values
        uint256 daiPriceInETH = realPriceOracle.getAssetPrice(address(dai));
        uint256 expectedCollateralETH = (depositAmount * daiPriceInETH) / WAD;
        uint256 expectedAvailableBorrows = (expectedCollateralETH * 75) / 100; // 75% LTV

        assertEq(totalCollateralETH, expectedCollateralETH, "Total collateral should match");
        assertEq(totalDebtETH, 0, "Total debt should be zero");
        assertEq(newAvailableBorrowsETH, expectedAvailableBorrows, "Available borrows should match LTV");
        assertEq(healthFactor, type(uint256).max, "Health factor should be infinite with no debt");
    }

    /**
     * @dev Test deposit and borrow
     */
    function testDepositAndBorrow() public {
        // Alice deposits DAI as collateral
        uint256 depositAmount = 1000 * WAD;
        vm.startPrank(alice);
        dai.approve(address(lendingPool), depositAmount);
        lendingPool.deposit(address(dai), depositAmount, 0);
        vm.stopPrank();

        // Calculate how much WETH Alice can borrow
        (,, uint256 availableBorrowsETH,,,) = lendingPool.getUserAccountData(alice);
        uint256 wethPriceInETH = realPriceOracle.getAssetPrice(address(weth));
        uint256 maxWethToBorrow = (availableBorrowsETH * WAD) / wethPriceInETH;

        // Borrow 50% of max available WETH
        uint256 borrowAmount = maxWethToBorrow / 2;
        vm.startPrank(alice);
        lendingPool.borrow(address(weth), borrowAmount, 1, 0); // Variable rate
        vm.stopPrank();

        // Check debt token balance
        assertEq(IERC20(wethDebtToken).balanceOf(alice), borrowAmount, "Alice should have debt tokens");

        // Check WETH balance
        assertEq(weth.balanceOf(alice), borrowAmount, "Alice should receive borrowed WETH");

        // Check updated account data
        (uint256 totalCollateralETH, uint256 totalDebtETH, /*uint256 newAvailableBorrowsETH*/,,, uint256 healthFactor) =
            lendingPool.getUserAccountData(alice);

        // Expected values
        uint256 expectedDebtETH = (borrowAmount * wethPriceInETH) / WAD;
        uint256 daiPriceInETH = realPriceOracle.getAssetPrice(address(dai));
        uint256 expectedCollateralETH = (depositAmount * daiPriceInETH) / WAD;

        assertEq(totalCollateralETH, expectedCollateralETH, "Total collateral should be unchanged");
        assertEq(totalDebtETH, expectedDebtETH, "Total debt should match borrowed amount");
        assertTrue(availableBorrowsETH < (expectedCollateralETH * 75) / 100, "Available borrows should be reduced");
        assertTrue(healthFactor >= 1e18, "Health factor should be >= 1");
        assertTrue(healthFactor < type(uint256).max, "Health factor should be finite with debt");
    }

    /**
     * @dev Test multiple users depositing and borrowing
     */
    function testMultipleUsersDepositAndBorrow() public {
        // Alice deposits DAI
        uint256 aliceDepositAmount = 500 * WAD;
        vm.startPrank(alice);
        dai.approve(address(lendingPool), aliceDepositAmount);
        lendingPool.deposit(address(dai), aliceDepositAmount, 0);
        vm.stopPrank();

        // Bob deposits USDC
        uint256 bobDepositAmount = 500 * 1e6; // USDC has 6 decimals
        vm.startPrank(bob);
        usdc.approve(address(lendingPool), bobDepositAmount);
        lendingPool.deposit(address(usdc), bobDepositAmount, 0);
        vm.stopPrank();

        // Carol deposits WETH
        uint256 carolDepositAmount = 5 * WAD;
        vm.startPrank(carol);
        weth.approve(address(lendingPool), carolDepositAmount);
        lendingPool.deposit(address(weth), carolDepositAmount, 0);
        vm.stopPrank();

        // Alice borrows USDC
        (,, uint256 aliceAvailableBorrowsETH,,,) = lendingPool.getUserAccountData(alice);
        uint256 usdcPriceInETH = realPriceOracle.getAssetPrice(address(usdc));
        uint256 maxUsdcToBorrow = (aliceAvailableBorrowsETH * WAD) / usdcPriceInETH;
        uint256 aliceBorrowAmount = (maxUsdcToBorrow / 2) / 1e12; // Convert to USDC decimals

        vm.startPrank(alice);
        lendingPool.borrow(address(usdc), aliceBorrowAmount, 1, 0);
        vm.stopPrank();

        // Bob borrows WETH
        (,, uint256 bobAvailableBorrowsETH,,,) = lendingPool.getUserAccountData(bob);
        uint256 wethPriceInETH = realPriceOracle.getAssetPrice(address(weth));
        uint256 maxWethToBorrow = (bobAvailableBorrowsETH * WAD) / wethPriceInETH;
        uint256 bobBorrowAmount = maxWethToBorrow / 2;

        vm.startPrank(bob);
        lendingPool.borrow(address(weth), bobBorrowAmount, 1, 0);
        vm.stopPrank();

        // Carol borrows DAI
        (,, uint256 carolAvailableBorrowsETH,,,) = lendingPool.getUserAccountData(carol);
        uint256 daiPriceInETH = realPriceOracle.getAssetPrice(address(dai));
        uint256 maxDaiToBorrow = (carolAvailableBorrowsETH * WAD) / daiPriceInETH;
        uint256 carolBorrowAmount = maxDaiToBorrow / 2;

        vm.startPrank(carol);
        lendingPool.borrow(address(dai), carolBorrowAmount, 1, 0);
        vm.stopPrank();

        // Check all users have debt and collateral
        (uint256 aliceCollateralETH, uint256 aliceDebtETH,,,, uint256 aliceHealthFactor) =
            lendingPool.getUserAccountData(alice);
        (uint256 bobCollateralETH, uint256 bobDebtETH,,,, uint256 bobHealthFactor) = lendingPool.getUserAccountData(bob);
        (uint256 carolCollateralETH, uint256 carolDebtETH,,,, uint256 carolHealthFactor) =
            lendingPool.getUserAccountData(carol);

        assertTrue(aliceCollateralETH > 0, "Alice should have collateral");
        assertTrue(aliceDebtETH > 0, "Alice should have debt");
        assertTrue(aliceHealthFactor >= 1e18, "Alice's health factor should be healthy");

        assertTrue(bobCollateralETH > 0, "Bob should have collateral");
        assertTrue(bobDebtETH > 0, "Bob should have debt");
        assertTrue(bobHealthFactor >= 1e18, "Bob's health factor should be healthy");

        assertTrue(carolCollateralETH > 0, "Carol should have collateral");
        assertTrue(carolDebtETH > 0, "Carol should have debt");
        assertTrue(carolHealthFactor >= 1e18, "Carol's health factor should be healthy");
    }

    /**
     * @dev Test deposit and withdraw
     */
    function testDepositAndWithdraw() public {
        // Alice deposits DAI
        uint256 depositAmount = 100 * WAD;
        vm.startPrank(alice);
        dai.approve(address(lendingPool), depositAmount);
        lendingPool.deposit(address(dai), depositAmount, 0);
        vm.stopPrank();

        // Check aToken balance
        assertEq(IERC20(daiAToken).balanceOf(alice), depositAmount, "Alice should receive aTokens");

        // Initial DAI balance
        uint256 initialDaiBalance = dai.balanceOf(alice);

        // Alice withdraws half of the deposit
        uint256 withdrawAmount = depositAmount / 2;
        vm.startPrank(alice);
        uint256 actualWithdrawn = lendingPool.withdraw(address(dai), withdrawAmount);
        vm.stopPrank();

        // Check results
        assertEq(actualWithdrawn, withdrawAmount, "Withdraw should return the requested amount");
        assertEq(dai.balanceOf(alice), initialDaiBalance + withdrawAmount, "Alice should receive DAI");
        assertEq(IERC20(daiAToken).balanceOf(alice), depositAmount - withdrawAmount, "aToken balance should be reduced");

        // Withdraw the rest
        vm.startPrank(alice);
        actualWithdrawn = lendingPool.withdraw(address(dai), type(uint256).max); // Withdraw all
        vm.stopPrank();

        assertEq(actualWithdrawn, depositAmount - withdrawAmount, "Should withdraw remaining amount");
        assertEq(dai.balanceOf(alice), initialDaiBalance + depositAmount, "Alice should receive all DAI back");
        assertEq(IERC20(daiAToken).balanceOf(alice), 0, "aToken balance should be zero");
    }

    /**
     * @dev Test deposit, borrow and repay
     */
    function testDepositBorrowAndRepay() public {
        // Alice deposits DAI as collateral
        uint256 depositAmount = 1000 * WAD;
        vm.startPrank(alice);
        dai.approve(address(lendingPool), depositAmount);
        lendingPool.deposit(address(dai), depositAmount, 0);
        vm.stopPrank();

        // Alice borrows WETH
        (,, uint256 availableBorrowsETH,,,) = lendingPool.getUserAccountData(alice);
        uint256 wethPriceInETH = realPriceOracle.getAssetPrice(address(weth));
        uint256 maxWethToBorrow = (availableBorrowsETH * WAD) / wethPriceInETH;
        uint256 borrowAmount = maxWethToBorrow / 2;

        vm.startPrank(alice);
        lendingPool.borrow(address(weth), borrowAmount, 1, 0);
        vm.stopPrank();

        // Check debt token balance
        assertEq(IERC20(wethDebtToken).balanceOf(alice), borrowAmount, "Alice should have debt tokens");

        // Alice approves and repays half the loan
        uint256 repayAmount = borrowAmount / 2;
        vm.startPrank(alice);
        weth.approve(address(lendingPool), repayAmount);
        uint256 actualRepaid = lendingPool.repay(address(weth), repayAmount, 1);
        vm.stopPrank();

        // Check results
        assertEq(actualRepaid, repayAmount, "Repay should return the amount repaid");
        assertEq(IERC20(wethDebtToken).balanceOf(alice), borrowAmount - repayAmount, "Debt should be reduced");

        // Alice repays the rest
        vm.startPrank(alice);
        weth.approve(address(lendingPool), borrowAmount); // Approve more than needed
        actualRepaid = lendingPool.repay(address(weth), type(uint256).max, 1); // Repay all
        vm.stopPrank();

        assertEq(actualRepaid, borrowAmount - repayAmount, "Should repay remaining amount");
        assertEq(IERC20(wethDebtToken).balanceOf(alice), 0, "Debt should be zero");

        // Check user account data - should have no debt
        ( /*uint256 totalCollateralETH*/ , uint256 totalDebtETH, uint256 newAvailableBorrowsETH,,, uint256 healthFactor)
        = lendingPool.getUserAccountData(alice);

        assertEq(totalDebtETH, 0, "Total debt should be zero");
        assertTrue(newAvailableBorrowsETH > 0, "Available borrows should be restored");
        assertEq(healthFactor, type(uint256).max, "Health factor should be infinity with no debt");
    }
}
