// SPDX-License-Identifier: MIT
pragma solidity 0.8.29;

import "forge-std/Test.sol";
import "../helpers/TestingHelper.sol";
import "../mocks/MockERC20.sol";
import "../mocks/MockChainlinkAggregator.sol";
import "../mocks/MockPriceOracle.sol";
import "../mocks/MockInterestRateStrategy.sol";
import "../mocks/MockLendingPoolAddressesProvider.sol";
import "../mocks/MockLendingPoolCore.sol";
import "../mocks/MockAToken.sol";
import "../mocks/MockDebtToken.sol";
import "../mocks/MockCollateralManager.sol";

/**
 * @title BaseTest
 * @author DeFi Lending Platform
 * @notice Base test contract that sets up the testing environment
 * @dev Other test contracts should inherit from this
 */
contract BaseTest is Test {
    // Helper contract
    TestingHelper internal helper;

    // Test accounts
    address internal admin;
    address internal alice;
    address internal bob;
    address internal carol;
    address internal liquidator;

    // Mock contracts
    MockLendingPoolAddressesProvider internal addressesProvider;
    MockPriceOracle internal priceOracle;
    MockLendingPoolCore internal mockLendingPoolCore;
    MockCollateralManager internal mockCollateralManager;

    // Mock tokens
    MockERC20 internal dai;
    MockERC20 internal usdc;
    MockERC20 internal weth;

    // Mock price feeds
    MockChainlinkAggregator internal ethUsdFeed;
    MockChainlinkAggregator internal daiUsdFeed;
    MockChainlinkAggregator internal usdcUsdFeed;

    // Mock interest rate strategy
    MockInterestRateStrategy internal interestRateStrategy;

    // Constants
    uint256 internal constant WAD = 1e18;
    uint256 internal constant RAY = 1e27;
    uint256 internal constant INITIAL_ETH_PRICE = 2000 * 1e8; // $2000 with 8 decimals
    uint256 internal constant INITIAL_TOKEN_AMOUNT = 10000 * WAD; // 10,000 tokens

    /**
     * @dev Sets up the test environment before each test
     */
    function setUp() public virtual {
        // Create helper
        helper = new TestingHelper();

        // Set up test accounts
        admin = makeAddr("admin");
        alice = makeAddr("alice");
        bob = makeAddr("bob");
        carol = makeAddr("carol");
        liquidator = makeAddr("liquidator");

        // Deal ETH to accounts
        vm.deal(admin, 100 ether);
        vm.deal(alice, 10 ether);
        vm.deal(bob, 10 ether);
        vm.deal(carol, 10 ether);
        vm.deal(liquidator, 10 ether);

        // Set up mock price feeds
        ethUsdFeed = new MockChainlinkAggregator(int256(INITIAL_ETH_PRICE), 8, "ETH / USD");
        daiUsdFeed = new MockChainlinkAggregator(int256(1 * 1e8), 8, "DAI / USD"); // $1
        usdcUsdFeed = new MockChainlinkAggregator(int256(1 * 1e8), 8, "USDC / USD"); // $1

        // Set up price oracle
        priceOracle = new MockPriceOracle(admin, INITIAL_ETH_PRICE);

        // Set up addresses provider
        addressesProvider = new MockLendingPoolAddressesProvider(admin);

        // Set up lending pool core
        mockLendingPoolCore = new MockLendingPoolCore(admin);

        // Set up collateral manager
        mockCollateralManager = new MockCollateralManager(admin);

        // Set up mock tokens
        dai = MockERC20(helper.createDAI(INITIAL_TOKEN_AMOUNT * 10, admin));
        usdc = MockERC20(helper.createUSDC(INITIAL_TOKEN_AMOUNT * 10, admin));
        weth = MockERC20(helper.createWETH(INITIAL_TOKEN_AMOUNT, admin));

        // Set up interest rate strategy
        interestRateStrategy = new MockInterestRateStrategy(
            WAD / 100, // 1% base borrow rate
            WAD / 10, // 10% max borrow rate
            WAD / 50, // 2% liquidity rate
            0, // 0% stable borrow rate (not used)
            WAD / 20 // 5% variable borrow rate
        );

        // Register addresses in provider
        vm.startPrank(admin);
        addressesProvider.setLendingPoolCoreImpl(address(mockLendingPoolCore));
        addressesProvider.setPriceOracle(address(priceOracle));
        addressesProvider.setCollateralManagerImpl(address(mockCollateralManager));
        vm.stopPrank();

        // Set up price oracle
        vm.startPrank(admin);
        priceOracle.setAssetPrice(address(dai), WAD); // 1 DAI = 1 ETH * $1 / $2000 = 0.0005 ETH
        priceOracle.setAssetPrice(address(usdc), WAD); // 1 USDC = 1 ETH * $1 / $2000 = 0.0005 ETH
        priceOracle.setAssetPrice(address(weth), 2000 * WAD); // 1 WETH = 1 ETH = $2000
        vm.stopPrank();

        // Distribute tokens to test accounts
        distributeTokens();
    }

    /**
     * @dev Distributes tokens to test accounts
     */
    function distributeTokens() internal {
        address[] memory users = new address[](3);
        users[0] = alice;
        users[1] = bob;
        users[2] = carol;

        // Distribute DAI
        vm.startPrank(admin);
        for (uint256 i = 0; i < users.length; i++) {
            dai.transfer(users[i], INITIAL_TOKEN_AMOUNT);
            usdc.transfer(users[i], INITIAL_TOKEN_AMOUNT / 1e12); // Adjust for decimals
            weth.transfer(users[i], INITIAL_TOKEN_AMOUNT / 10); // Less WETH
        }

        // Give liquidator some stablecoins for liquidations
        dai.transfer(liquidator, INITIAL_TOKEN_AMOUNT);
        usdc.transfer(liquidator, INITIAL_TOKEN_AMOUNT / 1e12);
        vm.stopPrank();
    }

    /**
     * @dev Helper method to create and initialize a reserve
     * @param asset The asset address
     * @param aTokenName The name of the aToken
     * @param aTokenSymbol The symbol of the aToken
     * @param debtTokenName The name of the debt token
     * @param debtTokenSymbol The symbol of the debt token
     * @return aToken The aToken address
     * @return debtToken The debt token address
     */
    function setupReserve(
        address asset,
        string memory aTokenName,
        string memory aTokenSymbol,
        string memory debtTokenName,
        string memory debtTokenSymbol
    ) internal returns (address aToken, address debtToken) {
        aToken = helper.createAToken(aTokenName, aTokenSymbol, asset, address(mockLendingPoolCore));

        debtToken = helper.createDebtToken(
            debtTokenName, debtTokenSymbol, asset, address(mockLendingPoolCore), address(mockLendingPoolCore)
        );

        // Initialize the reserve
        vm.prank(admin);
        mockLendingPoolCore.initReserve(asset, aToken, debtToken, address(interestRateStrategy));

        // Configure collateral settings
        vm.prank(admin);
        mockCollateralManager.configureAsCollateral(
            asset,
            true,
            75 * 1e16, // 75% LTV
            80 * 1e16, // 80% liquidation threshold
            110 * 1e16 // 110% liquidation bonus
        );

        return (aToken, debtToken);
    }
}
