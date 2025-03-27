// SPDX-License-Identifier: MIT
pragma solidity 0.8.29;

import "forge-std/Test.sol";
import "../../../src/oracles/PriceOracle.sol";
import "../../helpers/TestingHelper.sol";
import "../../mocks/MockERC20.sol";
import "../../mocks/MockChainlinkAggregator.sol";

/**
 * @title PriceOracleTest
 * @author DeFi Lending Platform
 * @notice Unit tests for the PriceOracle contract
 */
contract PriceOracleTest is Test {
    // Contract under test
    PriceOracle private priceOracle;

    // Helper contract
    TestingHelper private helper;

    // Test accounts
    address private admin;
    address private user;

    // Mock contracts
    MockERC20 private dai;
    MockERC20 private usdc;
    MockERC20 private weth;
    MockChainlinkAggregator private ethUsdFeed;
    MockChainlinkAggregator private daiUsdFeed;
    MockChainlinkAggregator private usdcUsdFeed;

    // Constants
    uint256 private constant WAD = 1e18;
    uint256 private constant INITIAL_SUPPLY = 10000 * WAD;
    uint256 private constant ETH_USD_PRICE = 2000 * 1e8; // $2000 with 8 decimals
    uint256 private constant DAI_USD_PRICE = 1 * 1e8; // $1 with 8 decimals
    uint256 private constant USDC_USD_PRICE = 1 * 1e8; // $1 with 8 decimals

    /**
     * @dev Set up the test fixture
     */
    function setUp() public {
        helper = new TestingHelper();

        admin = makeAddr("admin");
        user = makeAddr("user");

        // Create mock tokens
        dai = MockERC20(helper.createDAI(INITIAL_SUPPLY, admin));
        usdc = MockERC20(helper.createUSDC(INITIAL_SUPPLY / 1e12, admin)); // USDC has 6 decimals
        weth = MockERC20(helper.createWETH(INITIAL_SUPPLY, admin));

        // Create mock price feeds
        ethUsdFeed = new MockChainlinkAggregator(int256(ETH_USD_PRICE), 8, "ETH / USD");
        daiUsdFeed = new MockChainlinkAggregator(int256(DAI_USD_PRICE), 8, "DAI / USD");
        usdcUsdFeed = new MockChainlinkAggregator(int256(USDC_USD_PRICE), 8, "USDC / USD");

        // Create price oracle
        vm.startPrank(admin);
        priceOracle = new PriceOracle(admin, address(ethUsdFeed));
        vm.stopPrank();
    }

    /**
     * @dev Test initialization parameters
     */
    function testInitialization() public view {
        assertEq(priceOracle.owner(), admin, "Owner should be set correctly");
        assertEq(priceOracle.getEthUsdPrice(), ETH_USD_PRICE, "ETH/USD price should be set correctly");
    }

    /**
     * @dev Test setting asset sources
     */
    function testSetAssetSources() public {
        vm.startPrank(admin);
        priceOracle.setAssetSource(address(dai), address(daiUsdFeed));
        priceOracle.setAssetSource(address(usdc), address(usdcUsdFeed));
        vm.stopPrank();

        // Non-admin should not be able to set sources
        vm.startPrank(user);
        vm.expectRevert();
        priceOracle.setAssetSource(address(weth), address(ethUsdFeed));
        vm.stopPrank();

        // Test batch setting
        address[] memory assets = new address[](1);
        assets[0] = address(weth);

        address[] memory sources = new address[](1);
        sources[0] = address(ethUsdFeed);

        vm.startPrank(admin);
        priceOracle.setAssetsSources(assets, sources);
        vm.stopPrank();
    }

    /**
     * @dev Test getting asset prices
     */
    function testGetAssetPrice() public {
        // Setup price feeds
        vm.startPrank(admin);
        priceOracle.setAssetSource(address(dai), address(daiUsdFeed));
        priceOracle.setAssetSource(address(usdc), address(usdcUsdFeed));
        priceOracle.setAssetSource(address(weth), address(ethUsdFeed));
        vm.stopPrank();

        // Calculate expected ETH prices
        // 1 DAI = $1, 1 ETH = $2000, so 1 DAI = 0.0005 ETH
        uint256 expectedDaiEthPrice = WAD / 2000;

        // 1 USDC = $1, 1 ETH = $2000, so 1 USDC = 0.0005 ETH
        uint256 expectedUsdcEthPrice = WAD / 2000;

        // 1 WETH = 1 ETH
        uint256 expectedWethEthPrice = WAD;

        // Get actual prices
        uint256 daiEthPrice = priceOracle.getAssetPrice(address(dai));
        uint256 usdcEthPrice = priceOracle.getAssetPrice(address(usdc));
        uint256 wethEthPrice = priceOracle.getAssetPrice(address(weth));

        // Check results
        assertApproxEqRel(daiEthPrice, expectedDaiEthPrice, 0.001e18, "DAI/ETH price should be correct");
        assertApproxEqRel(usdcEthPrice, expectedUsdcEthPrice, 0.001e18, "USDC/ETH price should be correct");
        assertApproxEqRel(wethEthPrice, expectedWethEthPrice, 0.001e18, "WETH/ETH price should be correct");
    }

    /**
     * @dev Test getting multiple asset prices
     */
    function testGetAssetsPrices() public {
        // Setup price feeds
        vm.startPrank(admin);
        priceOracle.setAssetSource(address(dai), address(daiUsdFeed));
        priceOracle.setAssetSource(address(usdc), address(usdcUsdFeed));
        priceOracle.setAssetSource(address(weth), address(ethUsdFeed));
        vm.stopPrank();

        // Create arrays for batch query
        address[] memory assets = new address[](3);
        assets[0] = address(dai);
        assets[1] = address(usdc);
        assets[2] = address(weth);

        // Get prices in batch
        uint256[] memory prices = priceOracle.getAssetsPrices(assets);

        // Check array length
        assertEq(prices.length, assets.length, "Should return the same number of prices as assets");

        // Calculate expected ETH prices
        uint256 expectedDaiEthPrice = WAD / 2000;
        uint256 expectedUsdcEthPrice = WAD / 2000;
        uint256 expectedWethEthPrice = WAD;

        // Check individual prices
        assertApproxEqRel(prices[0], expectedDaiEthPrice, 0.001e18, "DAI/ETH price should be correct");
        assertApproxEqRel(prices[1], expectedUsdcEthPrice, 0.001e18, "USDC/ETH price should be correct");
        assertApproxEqRel(prices[2], expectedWethEthPrice, 0.001e18, "WETH/ETH price should be correct");
    }

    /**
     * @dev Test price updates
     */
    function testPriceUpdates() public {
        // Setup price feeds
        vm.startPrank(admin);
        priceOracle.setAssetSource(address(dai), address(daiUsdFeed));
        vm.stopPrank();

        // Initial price check
        uint256 initialDaiEthPrice = priceOracle.getAssetPrice(address(dai));

        // Update DAI price (from $1 to $1.1)
        vm.startPrank(ethUsdFeed.owner());
        ethUsdFeed.setPrice(int256(ETH_USD_PRICE * 110 / 100)); // ETH price increased by 10%
        vm.stopPrank();

        // Update ETH/USD price in the oracle
        vm.startPrank(admin);
        priceOracle.updateEthUsdPrice();
        vm.stopPrank();

        // Get new price
        uint256 newDaiEthPrice = priceOracle.getAssetPrice(address(dai));

        // DAI/ETH price should be less when ETH is worth more
        assertTrue(newDaiEthPrice < initialDaiEthPrice, "DAI/ETH price should decrease when ETH/USD price increases");
        assertApproxEqRel(
            newDaiEthPrice, initialDaiEthPrice * 100 / 110, 0.001e18, "Price change should be proportional"
        );
    }

    /**
     * @dev Test handling of stale prices
     */
    function testStalePrices() public {
        // Setup price feeds
        vm.startPrank(admin);
        priceOracle.setAssetSource(address(dai), address(daiUsdFeed));
        vm.stopPrank();

        // Ensure we have a non-zero block.timestamp to work with
        vm.warp(10 hours);

        // Make the price feed stale by setting timestamp in the past
        // The PRICE_EXPIRATION_TIME is 1 hour
        vm.startPrank(daiUsdFeed.owner());
        daiUsdFeed.setTimestampAgo(2 hours); // 2 hours ago (stale)

        // Verify the timestamp is actually stale
        (,,, uint256 updatedAt,) = daiUsdFeed.latestRoundData();
        console.log("Block timestamp:", block.timestamp);
        console.log("Updated timestamp:", updatedAt);
        console.log("Difference:", block.timestamp - updatedAt);
        assertTrue(block.timestamp - updatedAt > 1 hours, "Price should be stale");

        vm.stopPrank();

        // Try to get the stale price
        vm.expectRevert("PriceOracle: Stale price");
        priceOracle.getAssetPrice(address(dai));

        // Fix the price feed timestamp
        vm.startPrank(daiUsdFeed.owner());
        daiUsdFeed.setTimestampAgo(30 minutes); // 30 minutes ago (fresh)
        vm.stopPrank();

        // Now should be able to get the price
        uint256 price = priceOracle.getAssetPrice(address(dai));
        assertTrue(price > 0, "Should get a valid price after fixing the timestamp");
    }

    /**
     * @dev Test ETH/USD price management
     */
    function testEthUsdPriceManagement() public {
        // Initial price check
        assertEq(priceOracle.getEthUsdPrice(), ETH_USD_PRICE, "Initial ETH/USD price should be set from constructor");

        // Update price via feed
        vm.startPrank(ethUsdFeed.owner());
        ethUsdFeed.setPrice(int256(ETH_USD_PRICE * 110 / 100)); // 10% increase
        vm.stopPrank();

        vm.startPrank(admin);
        priceOracle.updateEthUsdPrice();
        vm.stopPrank();

        assertEq(priceOracle.getEthUsdPrice(), ETH_USD_PRICE * 110 / 100, "ETH/USD price should be updated from feed");

        // Manual override
        uint256 manualPrice = 3000 * 1e8; // $3000

        vm.startPrank(admin);
        priceOracle.setEthUsdPrice(manualPrice);
        vm.stopPrank();

        assertEq(priceOracle.getEthUsdPrice(), manualPrice, "ETH/USD price should be set manually");
    }

    /**
     * @dev Test handling of negative prices
     */
    function testNegativePrices() public {
        // Setup price feed with negative price
        vm.startPrank(daiUsdFeed.owner());
        daiUsdFeed.setPrice(-1 * int256(DAI_USD_PRICE));
        vm.stopPrank();

        // Setup price feeds in oracle
        vm.startPrank(admin);
        priceOracle.setAssetSource(address(dai), address(daiUsdFeed));
        vm.stopPrank();

        // Try to get the negative price
        vm.expectRevert("PriceOracle: Negative price");
        priceOracle.getAssetPrice(address(dai));
    }
}
