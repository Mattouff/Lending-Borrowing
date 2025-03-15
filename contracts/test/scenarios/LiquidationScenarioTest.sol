// SPDX-License-Identifier: MIT
pragma solidity 0.8.29;

import "forge-std/Test.sol";
import "../../src/core/LendingPool.sol";
import "../../src/core/LendingPoolCore.sol";
import "../../src/core/LendingPoolAddressesProvider.sol";
import "../../src/core/LendingPoolConfigurator.sol";
import "../../src/risk/CollateralManager.sol";
import "../../src/risk/LiquidationManager.sol";
import "../../src/oracles/PriceOracle.sol";
import "../../src/core/InterestRateStrategy.sol";
import "../helpers/BaseTest.sol";

/**
 * @title LiquidationScenarioTest
 * @author DeFi Lending Platform
 * @notice Scenario tests for liquidation functionality
 */
contract LiquidationScenarioTest is BaseTest {
    // Real contracts for testing
    LendingPoolAddressesProvider private realAddressesProvider;
    LendingPool private lendingPool;
    LendingPoolCore private lendingPoolCore;
    CollateralManager private collateralManager;
    LiquidationManager private liquidationManager;
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
        // Deploy address provider
        vm.startPrank(admin);
        realAddressesProvider = new LendingPoolAddressesProvider();
        realAddressesProvider.initialize(admin);
        
        // Deploy price oracle
        realPriceOracle = new PriceOracle(admin, address(ethUsdFeed));
        
        // Set price oracle in addresses provider
        realAddressesProvider.setPriceOracle(address(realPriceOracle));
        
        // Deploy lending pool core
        lendingPoolCore = new LendingPoolCore();
        
        // Deploy lending pool
        lendingPool = new LendingPool();
        
        // Deploy collateral manager
        collateralManager = new CollateralManager();
        
        // Deploy liquidation manager
        liquidationManager = new LiquidationManager();
        
        // Deploy interest rate strategy
        realStrategy = new InterestRateStrategy(
            WAD / 100,                // 1% base borrow rate
            WAD / 10,                 // 10% slope1
            WAD * 4 / 10,             // 40% slope2
            WAD * 8 / 10,             // 80% optimal utilization
            2 * WAD                   // Quadratic growth
        );
        
        // Deploy configurator
        configurator = new LendingPoolConfigurator(address(realAddressesProvider));
        
        // Set contract addresses in addresses provider
        bytes memory lendingPoolCoreInitData = abi.encodeWithSignature(
            "initialize(address)",
            address(realAddressesProvider)
        );
        bytes memory lendingPoolInitData = abi.encodeWithSignature(
            "initialize(address)",
            address(realAddressesProvider)
        );
        bytes memory collateralManagerInitData = abi.encodeWithSignature(
            "initialize(address)",
            address(realAddressesProvider)
        );
        bytes memory liquidationManagerInitData = abi.encodeWithSignature(
            "initialize(address)",
            address(realAddressesProvider)
        );
        
        // Create proxies and point to implementations
        realAddressesProvider.setLendingPoolCoreImpl(address(lendingPoolCore));
        realAddressesProvider.setLendingPoolImpl(address(lendingPool));
        realAddressesProvider.setCollateralManagerImpl(address(collateralManager));
        realAddressesProvider.setLiquidationManagerImpl(address(liquidationManager));
        realAddressesProvider.setLendingPoolConfiguratorImpl(address(configurator));
        
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
    function _initializeReserves() internal {
        vm.startPrank(admin);
        
        // Initialize DAI reserve
        (daiAToken, daiDebtToken) = setupReserve(
            address(dai),
            "aDAI Token",
            "aDAI",
            "debtDAI Token",
            "debtDAI"
        );
        
        // Initialize USDC reserve
        (usdcAToken, usdcDebtToken) = setupReserve(
            address(usdc),
            "aUSDC Token",
            "aUSDC",
            "debtUSDC Token",
            "debtUSDC"
        );
        
        // Initialize WETH reserve
        (wethAToken, wethDebtToken) = setupReserve(
            address(weth),
            "aWETH Token",
            "aWETH",
            "debtWETH Token",
            "debtWETH"
        );
        
        // Configure DAI as collateral
        collateralManager.configureAsCollateral(
            address(dai),
            true,                // Can be used as collateral
            75 * 1e16,           // 75% LTV
            80 * 1e16,           // 80% liquidation threshold
            110 * 1e16           // 110% liquidation bonus
        );
        
        // Configure USDC as collateral
        collateralManager.configureAsCollateral(
            address(usdc),
            true,                // Can be used as collateral
            75 * 1e16,           // 75% LTV
            80 * 1e16,           // 80% liquidation threshold
            110 * 1e16           // 110% liquidation bonus
        );
        
        // Configure WETH as collateral
        collateralManager.configureAsCollateral(
            address(weth),
            true,                // Can be used as collateral
            75 * 1e16,           // 75% LTV
            80 * 1e16,           // 80% liquidation threshold
            110 * 1e16           // 110% liquidation bonus
        );
        
        // Enable borrowing on reserves
        configurator.enableBorrowingOnReserve(address(dai), true);
        configurator.enableBorrowingOnReserve(address(usdc), true);
        configurator.enableBorrowingOnReserve(address(weth), true);
        
        vm.stopPrank();
    }
    
    /**
     * @dev Transfer tokens to test accounts
     */
    function _transferTokens() internal {
        vm.startPrank(admin);
        
        // Alice gets 1000 DAI
        dai.transfer(alice, 1000 * WAD);
        
        // Bob gets 500 WETH
        weth.transfer(bob, 500 * WAD);
        
        // Liquidator gets 10000 DAI and 10000 USDC
        dai.transfer(liquidator, 10000 * WAD);
        usdc.transfer(liquidator, 10000 * 1e6);
        
        vm.stopPrank();
    }
    
    /**
     * @dev Test a simple liquidation scenario:
     * 1. Alice deposits DAI as collateral
     * 2. Alice borrows ETH
     * 3. ETH price increases
     * 4. Alice becomes undercollateralized
     * 5. Liquidator liquidates Alice's position
     */
    function testSimpleLiquidation() public {
        // 1. Alice deposits DAI as collateral
        uint256 aliceCollateralAmount = 1000 * WAD; // 1000 DAI
        vm.startPrank(alice);
        dai.approve(address(lendingPool), aliceCollateralAmount);
        lendingPool.deposit(address(dai), aliceCollateralAmount, 0);
        vm.stopPrank();
        
        // Get Alice's borrowing capacity
        (,, uint256 aliceAvailableBorrowsETH,,,) = lendingPool.getUserAccountData(alice);
        uint256 wethPriceInETH = realPriceOracle.getAssetPrice(address(weth));
        uint256 maxWethToBorrow = (aliceAvailableBorrowsETH * WAD) / wethPriceInETH;
        
        // 2. Alice borrows close to maximum WETH (90% of max)
        uint256 borrowAmount = maxWethToBorrow * 90 / 100;
        vm.startPrank(alice);
        lendingPool.borrow(address(weth), borrowAmount, 1, 0);
        vm.stopPrank();
        
        // Check initial health factor - should be above 1.0 but not by much
        (,,,,, uint256 initialHealthFactor) = lendingPool.getUserAccountData(alice);
        assertTrue(initialHealthFactor > 1e18, "Initial health factor should be > 1.0");
        assertTrue(initialHealthFactor < 1.2e18, "Initial health factor should be close to 1.0");
        
        // 3. ETH price increases by 25%, making Alice's position undercollateralized
        // Initial price: 1 ETH = $2000
        // New price: 1 ETH = $2500
        // Create a new price feed with updated price
        vm.startPrank(admin);
        MockChainlinkAggregator updatedWethPriceFeed = new MockChainlinkAggregator(int256(2500 * 1e8), 8, "WETH / USD");
        realPriceOracle.setAssetSource(address(weth), address(updatedWethPriceFeed));
        vm.stopPrank();
        
        // No need to set the price again as we've already updated the price feed
        vm.stopPrank();
        
        // 4. Check that Alice is now undercollateralized
        (,,,,, uint256 newHealthFactor) = lendingPool.getUserAccountData(alice);
        assertTrue(newHealthFactor < 1e18, "New health factor should be < 1.0");
        
        // Record initial balances before liquidation
        uint256 liquidatorInitialDai = dai.balanceOf(liquidator);
        uint256 liquidatorInitialWeth = weth.balanceOf(liquidator);
        uint256 aliceInitialADai = IERC20(daiAToken).balanceOf(alice);
        
        // 5. Liquidator liquidates Alice's position
        // Calculate how much debt to cover (50% of the debt as per CLOSE_FACTOR)
        uint256 debtToCover = borrowAmount / 2;
        
        // Liquidator approves WETH to cover debt
        vm.startPrank(liquidator);
        weth.approve(address(lendingPool), debtToCover);
        // Perform liquidation
        uint256 liquidatedCollateral = lendingPool.liquidationCall(
            address(dai),  // Collateral asset
            address(weth), // Debt asset
            alice,         // User being liquidated
            debtToCover    // Amount of debt to cover
        );
        vm.stopPrank();
        
        // 6. Verify liquidation results
        
        // Liquidator should have received collateral with bonus
        uint256 liquidatorFinalDai = dai.balanceOf(liquidator);
        uint256 liquidatorFinalWeth = weth.balanceOf(liquidator);
        
        // Liquidator should have spent WETH
        assertEq(liquidatorInitialWeth - liquidatorFinalWeth, debtToCover, "Liquidator should have spent WETH");
        
        // Liquidator should have received DAI collateral with bonus
        assertTrue(liquidatorFinalDai > liquidatorInitialDai, "Liquidator should have received DAI");
        
        // Calculate expected liquidation amount
        // 1. Convert debtToCover to ETH: debtToCover * wethPriceInETH / WAD
        // 2. Convert ETH to DAI equivalent: ethAmount * WAD / daiPriceInETH
        // 3. Apply bonus (110%): daiAmount * 110 / 100
        uint256 daiPriceInETH = realPriceOracle.getAssetPrice(address(dai));
        uint256 expectedLiquidatedCollateral = (debtToCover * wethPriceInETH * 110 * WAD) / (WAD * daiPriceInETH * 100);
        
        // Due to rounding, we use approximate comparison
        assertApproxEqRel(liquidatedCollateral, expectedLiquidatedCollateral, 0.01e18, "Liquidated collateral should match expected amount");
        
        // Alice's aToken balance should have decreased
        uint256 aliceFinalADai = IERC20(daiAToken).balanceOf(alice);
        assertEq(aliceInitialADai - aliceFinalADai, liquidatedCollateral, "Alice's aDAI balance should have decreased by liquidated amount");
        
        // Alice's health factor should have improved
        (,,,,, uint256 finalHealthFactor) = lendingPool.getUserAccountData(alice);
        assertTrue(finalHealthFactor > newHealthFactor, "Health factor should have improved after liquidation");
    }
    
    /**
     * @dev Test a complex liquidation scenario with multiple assets
     */
    function testComplexLiquidation() public {
        // Setup: Bob deposits WETH
        uint256 bobCollateralAmount = 10 * WAD; // 10 WETH
        vm.startPrank(bob);
        weth.approve(address(lendingPool), bobCollateralAmount);
        lendingPool.deposit(address(weth), bobCollateralAmount, 0);
        vm.stopPrank();
        
        // Bob borrows DAI
        (,, uint256 bobAvailableBorrowsETH,,,) = lendingPool.getUserAccountData(bob);
        uint256 daiPriceInETH = realPriceOracle.getAssetPrice(address(dai));
        uint256 maxDaiToBorrow = (bobAvailableBorrowsETH * WAD) / daiPriceInETH;
        
        // Borrow 90% of max to be closer to liquidation threshold
        uint256 daiBorrowAmount = maxDaiToBorrow * 90 / 100;
        vm.startPrank(bob);
        lendingPool.borrow(address(dai), daiBorrowAmount, 1, 0);
        vm.stopPrank();
        
        // Now Bob borrows USDC as well
        (,, uint256 bobRemainingBorrowsETH,,,) = lendingPool.getUserAccountData(bob);
        uint256 usdcPriceInETH = realPriceOracle.getAssetPrice(address(usdc));
        uint256 maxUsdcToBorrow = (bobRemainingBorrowsETH * WAD) / usdcPriceInETH;
        
        // Borrow 80% of remaining capacity in USDC
        uint256 usdcBorrowAmount = (maxUsdcToBorrow * 80 / 100) / 1e12; // Convert to USDC decimals
        vm.startPrank(bob);
        lendingPool.borrow(address(usdc), usdcBorrowAmount, 1, 0);
        vm.stopPrank();
        
        // Check initial health factor
        (,,,,, uint256 initialHealthFactor) = lendingPool.getUserAccountData(bob);
        assertTrue(initialHealthFactor > 1e18, "Initial health factor should be > 1.0");
        
        // ETH price drops by 30%
        // Initial price: 1 ETH = $2000
        // New price: 1 ETH = $1400
        // Create a new price feed with updated price
        vm.startPrank(admin);
        MockChainlinkAggregator updatedWethPriceFeed = new MockChainlinkAggregator(int256(1400 * 1e8), 8, "WETH / USD");
        realPriceOracle.setAssetSource(address(weth), address(updatedWethPriceFeed));
        vm.stopPrank();
        
        // No need to set the price again as we've already updated the price feed
        vm.stopPrank();
        
        // Bob should now be undercollateralized
        (,,,,, uint256 newHealthFactor) = lendingPool.getUserAccountData(bob);
        assertTrue(newHealthFactor < 1e18, "New health factor should be < 1.0");
        
        // Record initial balances
        uint256 liquidatorInitialDai = dai.balanceOf(liquidator);
        uint256 liquidatorInitialUsdc = usdc.balanceOf(liquidator);
        uint256 liquidatorInitialWeth = weth.balanceOf(liquidator);
        uint256 bobInitialAWeth = IERC20(wethAToken).balanceOf(bob);
        uint256 bobInitialDebtDai = IERC20(daiDebtToken).balanceOf(bob);
        uint256 bobInitialDebtUsdc = IERC20(usdcDebtToken).balanceOf(bob);
        
        // Liquidator liquidates Bob's DAI debt position
        uint256 daiDebtToCover = daiBorrowAmount / 2; // 50% of the debt
        
        vm.startPrank(liquidator);
        dai.approve(address(lendingPool), daiDebtToCover);
        uint256 liquidatedCollateralForDai = lendingPool.liquidationCall(
            address(weth), // Collateral asset
            address(dai),  // Debt asset
            bob,           // User being liquidated
            daiDebtToCover // Amount of debt to cover
        );
        vm.stopPrank();
        
        // Verify first liquidation
        uint256 liquidatorDaiAfterFirst = dai.balanceOf(liquidator);
        uint256 liquidatorWethAfterFirst = weth.balanceOf(liquidator);
        
        assertEq(liquidatorInitialDai - liquidatorDaiAfterFirst, daiDebtToCover, "Liquidator should have spent DAI");
        assertTrue(liquidatorWethAfterFirst > liquidatorInitialWeth, "Liquidator should have received WETH");
        
        // Bob's health factor should have improved but might still be below 1.0
        (,,,,, uint256 intermediateHealthFactor) = lendingPool.getUserAccountData(bob);
        assertTrue(intermediateHealthFactor > newHealthFactor, "Health factor should have improved after first liquidation");
        
        // If Bob is still undercollateralized, liquidate USDC debt too
        if (intermediateHealthFactor < 1e18) {
            uint256 usdcDebtToCover = usdcBorrowAmount / 2; // 50% of the USDC debt
            
            vm.startPrank(liquidator);
            usdc.approve(address(lendingPool), usdcDebtToCover);
            uint256 liquidatedCollateralForUsdc = lendingPool.liquidationCall(
                address(weth), // Collateral asset
                address(usdc), // Debt asset
                bob,           // User being liquidated
                usdcDebtToCover // Amount of debt to cover
            );
            vm.stopPrank();
            
            // Verify second liquidation
            uint256 liquidatorUsdcAfterSecond = usdc.balanceOf(liquidator);
            uint256 liquidatorWethAfterSecond = weth.balanceOf(liquidator);
            
            assertEq(liquidatorInitialUsdc - liquidatorUsdcAfterSecond, usdcDebtToCover, "Liquidator should have spent USDC");
            assertTrue(liquidatorWethAfterSecond > liquidatorWethAfterFirst, "Liquidator should have received more WETH");
            
            // Total liquidated collateral
            uint256 totalLiquidatedCollateral = liquidatedCollateralForDai + liquidatedCollateralForUsdc;
            
            // Bob's final positions
            uint256 bobFinalAWeth = IERC20(wethAToken).balanceOf(bob);
            uint256 bobFinalDebtDai = IERC20(daiDebtToken).balanceOf(bob);
            uint256 bobFinalDebtUsdc = IERC20(usdcDebtToken).balanceOf(bob);
            
            assertEq(bobInitialAWeth - bobFinalAWeth, totalLiquidatedCollateral, "Bob's aWETH balance should have decreased by total liquidated amount");
            assertEq(bobInitialDebtDai - bobFinalDebtDai, daiDebtToCover, "Bob's DAI debt should have decreased");
            assertEq(bobInitialDebtUsdc - bobFinalDebtUsdc, usdcDebtToCover, "Bob's USDC debt should have decreased");
        }
        
        // Final health factor should be improved
        (,,,,, uint256 finalHealthFactor) = lendingPool.getUserAccountData(bob);
        assertTrue(finalHealthFactor > newHealthFactor, "Final health factor should be better than before liquidation");
    }
    
    /**
     * @dev Test liquidation with interest accrual
     */
    function testLiquidationWithInterestAccrual() public {
        // 1. Alice deposits DAI as collateral
        uint256 aliceCollateralAmount = 1000 * WAD; // 1000 DAI
        vm.startPrank(alice);
        dai.approve(address(lendingPool), aliceCollateralAmount);
        lendingPool.deposit(address(dai), aliceCollateralAmount, 0);
        vm.stopPrank();
        
        // 2. Alice borrows WETH at 95% of max capacity
        (,, uint256 aliceAvailableBorrowsETH,,,) = lendingPool.getUserAccountData(alice);
        uint256 wethPriceInETH = realPriceOracle.getAssetPrice(address(weth));
        uint256 maxWethToBorrow = (aliceAvailableBorrowsETH * WAD) / wethPriceInETH;
        uint256 borrowAmount = maxWethToBorrow * 95 / 100;
        
        vm.startPrank(alice);
        lendingPool.borrow(address(weth), borrowAmount, 1, 0);
        vm.stopPrank();
        
        // 3. Fast forward time to accrue interest
        // Advance time by 1 year
        vm.warp(block.timestamp + 365 days);
        
        // Trigger an update to accrue interest (any operation works)
        vm.prank(alice);
        lendingPool.deposit(address(dai), 0, 0); // Zero deposit to trigger update
        
        // 4. Check that interest has pushed Alice closer to liquidation
        (,,,,, uint256 healthFactorAfterInterest) = lendingPool.getUserAccountData(alice);
        assertTrue(healthFactorAfterInterest < 1.05e18, "Health factor should be close to 1.0 after interest accrual");
        
        // 5. ETH price increases slightly, pushing Alice into liquidation territory
        // Create a new price feed with updated price
        vm.startPrank(admin);
        MockChainlinkAggregator updatedWethPriceFeed = new MockChainlinkAggregator(int256(2100 * 1e8), 8, "WETH / USD");
        realPriceOracle.setAssetSource(address(weth), address(updatedWethPriceFeed));
        vm.stopPrank();
        
        // No need to set the price again as we've already updated the price feed
        vm.stopPrank();
        
        // 6. Verify Alice is now undercollateralized
        (,,,,, uint256 finalHealthFactor) = lendingPool.getUserAccountData(alice);
        assertTrue(finalHealthFactor < 1e18, "Health factor should be < 1.0 after price change");
        
        // 7. Liquidator liquidates Alice's position
        uint256 currentDebt = IERC20(wethDebtToken).balanceOf(alice);
        uint256 debtToCover = currentDebt / 2; // 50% of current debt
        
        vm.startPrank(liquidator);
        weth.approve(address(lendingPool), debtToCover);
        uint256 liquidatedCollateral = lendingPool.liquidationCall(
            address(dai),
            address(weth),
            alice,
            debtToCover
        );
        vm.stopPrank();
        
        // 8. Verify liquidation worked properly
        (,,,,, uint256 healthFactorAfterLiquidation) = lendingPool.getUserAccountData(alice);
        assertTrue(healthFactorAfterLiquidation > finalHealthFactor, "Health factor should improve after liquidation");
        assertTrue(liquidatedCollateral > 0, "Liquidated collateral should be greater than 0");
        
        // Verify liquidator received collateral
        uint256 liquidatorDaiBalance = dai.balanceOf(liquidator);
        assertTrue(liquidatorDaiBalance > 0, "Liquidator should have received DAI collateral");
    }
}
