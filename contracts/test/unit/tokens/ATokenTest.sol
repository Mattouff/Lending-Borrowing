// SPDX-License-Identifier: MIT
pragma solidity 0.8.29;

import "forge-std/Test.sol";
import "../../../src/tokens/AToken.sol";
import "../../helpers/TestingHelper.sol";
import "../../mocks/MockERC20.sol";

/**
 * @title ATokenTest
 * @author DeFi Lending Platform
 * @notice Unit tests for the AToken contract
 */
contract ATokenTest is Test {
    // Contract under test
    AToken private aToken;

    // Helper contract
    TestingHelper private helper;

    // Test accounts
    address private admin;
    address private lendingPool;
    address private user1;
    address private user2;

    // Mock contracts
    MockERC20 private dai;

    // Constants
    uint256 private constant INITIAL_SUPPLY = 10000 * 1e18;
    uint256 private constant MINT_AMOUNT = 1000 * 1e18;
    uint256 private constant INITIAL_INDEX = 1e18; // 1.0 in ray format

    /**
     * @dev Set up the test fixture
     */
    function setUp() public {
        helper = new TestingHelper();

        admin = makeAddr("admin");
        lendingPool = makeAddr("lendingPool");
        user1 = makeAddr("user1");
        user2 = makeAddr("user2");

        // Create DAI token
        dai = MockERC20(helper.createDAI(INITIAL_SUPPLY, admin));

        // Create aToken
        vm.startPrank(admin);
        aToken = new AToken(address(dai), lendingPool, "aDAI Token", "aDAI");
        vm.stopPrank();

        // Transfer DAI to users
        vm.startPrank(admin);
        dai.transfer(user1, 1000 * 1e18);
        dai.transfer(user2, 1000 * 1e18);
        vm.stopPrank();
    }

    /**
     * @dev Test initialization parameters
     */
    function testInitialization() public view {
        assertEq(aToken.name(), "aDAI Token", "Name should be correctly set");
        assertEq(aToken.symbol(), "aDAI", "Symbol should be correctly set");
        assertEq(aToken.getUnderlyingAssetAddress(), address(dai), "Underlying asset should be DAI");
        assertEq(aToken.owner(), lendingPool, "Owner should be lending pool");
    }

    /**
     * @dev Test mint functionality
     */
    function testMint() public {
        // Only lending pool should be able to mint
        vm.startPrank(admin);
        vm.expectRevert("AToken: Caller must be lending pool");
        aToken.mint(user1, MINT_AMOUNT, INITIAL_INDEX);
        vm.stopPrank();

        // Lending pool mints tokens to user1
        vm.startPrank(lendingPool);
        aToken.mint(user1, MINT_AMOUNT, INITIAL_INDEX);
        vm.stopPrank();

        // Check results
        assertEq(aToken.balanceOf(user1), MINT_AMOUNT, "User should have received aTokens");
        assertEq(aToken.scaledBalanceOf(user1), MINT_AMOUNT, "Scaled balance should be correct with initial index");
        assertEq(aToken.scaledTotalSupply(), MINT_AMOUNT, "Total scaled supply should be updated");
    }

    /**
     * @dev Test burn functionality
     */
    function testBurn() public {
        // Setup: mint tokens first
        vm.startPrank(lendingPool);
        aToken.mint(user1, MINT_AMOUNT, INITIAL_INDEX);
        vm.stopPrank();

        // Only lending pool should be able to burn
        vm.startPrank(admin);
        vm.expectRevert("AToken: Caller must be lending pool");
        aToken.burn(user1, MINT_AMOUNT / 2, INITIAL_INDEX);
        vm.stopPrank();

        // Lending pool burns half the tokens
        vm.startPrank(lendingPool);
        uint256 burnAmount = MINT_AMOUNT / 2;
        uint256 actualBurned = aToken.burn(user1, burnAmount, INITIAL_INDEX);
        vm.stopPrank();

        // Check results
        assertEq(actualBurned, burnAmount, "Burn should return the amount burned");
        assertEq(aToken.balanceOf(user1), MINT_AMOUNT - burnAmount, "User balance should be reduced");
        assertEq(aToken.scaledBalanceOf(user1), MINT_AMOUNT - burnAmount, "Scaled balance should be reduced");
        assertEq(aToken.scaledTotalSupply(), MINT_AMOUNT - burnAmount, "Total scaled supply should be reduced");
    }

    /**
     * @dev Test burn with insufficient balance
     */
    function testBurnInsufficientBalance() public {
        // Setup: mint tokens first
        vm.startPrank(lendingPool);
        aToken.mint(user1, MINT_AMOUNT, INITIAL_INDEX);
        vm.stopPrank();

        // Try to burn more than the user has
        vm.startPrank(lendingPool);
        vm.expectRevert("AToken: Insufficient balance");
        aToken.burn(user1, MINT_AMOUNT * 2, INITIAL_INDEX);
        vm.stopPrank();
    }

    /**
     * @dev Test transferUnderlyingTo function
     */
    function testTransferUnderlyingTo() public {
        // Send DAI to aToken contract
        vm.startPrank(admin);
        dai.transfer(address(aToken), MINT_AMOUNT);
        vm.stopPrank();

        // Only lending pool should be able to transfer underlying
        vm.startPrank(admin);
        vm.expectRevert("AToken: Caller must be lending pool");
        aToken.transferUnderlyingTo(user1, MINT_AMOUNT);
        vm.stopPrank();

        // Check DAI balance before
        uint256 user1DaiBefore = dai.balanceOf(user1);

        // Lending pool transfers underlying to user1
        vm.startPrank(lendingPool);
        aToken.transferUnderlyingTo(user1, MINT_AMOUNT);
        vm.stopPrank();

        // Check results
        assertEq(dai.balanceOf(user1), user1DaiBefore + MINT_AMOUNT, "User should have received DAI");
    }

    /**
     * @dev Test scaling with different indices
     */
    function testScalingWithDifferentIndices() public {
        // Initial mint at index 1.0
        vm.startPrank(lendingPool);
        aToken.mint(user1, MINT_AMOUNT, INITIAL_INDEX);
        vm.stopPrank();

        // Second mint at index 1.1 (10% interest accrued)
        uint256 higherIndex = INITIAL_INDEX * 110 / 100; // 1.1 in ray
        uint256 secondMintAmount = MINT_AMOUNT;

        vm.startPrank(lendingPool);
        aToken.mint(user1, secondMintAmount, higherIndex);
        vm.stopPrank();

        // Check results
        uint256 expectedScaledAmount = MINT_AMOUNT + (secondMintAmount * INITIAL_INDEX / higherIndex);
        assertApproxEqRel(
            aToken.scaledBalanceOf(user1),
            expectedScaledAmount,
            0.0001e18,
            "Scaled balance should account for different indices"
        );

        // The actual balance should be the sum of the minted amounts
        assertEq(aToken.balanceOf(user1), MINT_AMOUNT + secondMintAmount, "Balance should be the sum of minted amounts");
    }

    /**
     * @dev Test transfer between users
     */
    function testTransfer() public {
        // Setup: mint tokens first
        vm.startPrank(lendingPool);
        aToken.mint(user1, MINT_AMOUNT, INITIAL_INDEX);
        vm.stopPrank();

        // User1 transfers to User2
        vm.startPrank(user1);
        bool success = aToken.transfer(user2, MINT_AMOUNT / 2);
        vm.stopPrank();

        // Check results
        assertTrue(success, "Transfer should be successful");
        assertEq(aToken.balanceOf(user1), MINT_AMOUNT / 2, "User1 balance should be reduced");
        assertEq(aToken.balanceOf(user2), MINT_AMOUNT / 2, "User2 balance should be increased");

        // Check scaled balances
        assertEq(aToken.scaledBalanceOf(user1), MINT_AMOUNT / 2, "User1 scaled balance should be reduced");
        assertEq(aToken.scaledBalanceOf(user2), MINT_AMOUNT / 2, "User2 scaled balance should be increased");
    }

    /**
     * @dev Test transferFrom functionality
     */
    function testTransferFrom() public {
        // Setup: mint tokens first
        vm.startPrank(lendingPool);
        aToken.mint(user1, MINT_AMOUNT, INITIAL_INDEX);
        vm.stopPrank();

        // User1 approves User2 to spend tokens
        vm.startPrank(user1);
        aToken.approve(user2, MINT_AMOUNT / 2);
        vm.stopPrank();

        // User2 transfers tokens from User1 to themselves
        vm.startPrank(user2);
        bool success = aToken.transferFrom(user1, user2, MINT_AMOUNT / 2);
        vm.stopPrank();

        // Check results
        assertTrue(success, "TransferFrom should be successful");
        assertEq(aToken.balanceOf(user1), MINT_AMOUNT / 2, "User1 balance should be reduced");
        assertEq(aToken.balanceOf(user2), MINT_AMOUNT / 2, "User2 balance should be increased");

        // Check allowance
        assertEq(aToken.allowance(user1, user2), 0, "Allowance should be spent");
    }
}
