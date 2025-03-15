// SPDX-License-Identifier: MIT
pragma solidity 0.8.29;

import "forge-std/Test.sol";
import "../../../src/tokens/DebtToken.sol";
import "../../helpers/TestingHelper.sol";
import "../../mocks/MockERC20.sol";

/**
 * @title DebtTokenTest
 * @author DeFi Lending Platform
 * @notice Unit tests for the DebtToken contract
 */
contract DebtTokenTest is Test {
    // Contract under test
    DebtToken private debtToken;
    
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
        
        // Create debt token
        vm.startPrank(admin);
        debtToken = new DebtToken(address(dai), lendingPool, "debtDAI Token", "debtDAI");
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
        assertEq(debtToken.name(), "debtDAI Token", "Name should be correctly set");
        assertEq(debtToken.symbol(), "debtDAI", "Symbol should be correctly set");
        assertEq(debtToken.getUnderlyingAssetAddress(), address(dai), "Underlying asset should be DAI");
        assertEq(debtToken.owner(), lendingPool, "Owner should be lending pool");
    }
    
    /**
     * @dev Test mint functionality
     */
    function testMint() public {
        // Only lending pool should be able to mint
        vm.startPrank(admin);
        vm.expectRevert("DebtToken: Caller must be lending pool");
        debtToken.mint(user1, MINT_AMOUNT, INITIAL_INDEX);
        vm.stopPrank();
        
        // Lending pool mints tokens to user1
        vm.startPrank(lendingPool);
        debtToken.mint(user1, MINT_AMOUNT, INITIAL_INDEX);
        vm.stopPrank();
        
        // Check results
        assertEq(debtToken.balanceOf(user1), MINT_AMOUNT, "User should have received debt tokens");
        assertEq(debtToken.scaledBalanceOf(user1), MINT_AMOUNT, "Scaled balance should be correct with initial index");
        assertEq(debtToken.scaledTotalSupply(), MINT_AMOUNT, "Total scaled supply should be updated");
    }
    
    /**
     * @dev Test burn functionality
     */
    function testBurn() public {
        // Setup: mint tokens first
        vm.startPrank(lendingPool);
        debtToken.mint(user1, MINT_AMOUNT, INITIAL_INDEX);
        vm.stopPrank();
        
        // Only lending pool should be able to burn
        vm.startPrank(admin);
        vm.expectRevert("DebtToken: Caller must be lending pool");
        debtToken.burn(user1, MINT_AMOUNT / 2, INITIAL_INDEX);
        vm.stopPrank();
        
        // Lending pool burns half the tokens
        vm.startPrank(lendingPool);
        uint256 burnAmount = MINT_AMOUNT / 2;
        uint256 actualBurned = debtToken.burn(user1, burnAmount, INITIAL_INDEX);
        vm.stopPrank();
        
        // Check results
        assertEq(actualBurned, burnAmount, "Burn should return the amount burned");
        assertEq(debtToken.balanceOf(user1), MINT_AMOUNT - burnAmount, "User balance should be reduced");
        assertEq(debtToken.scaledBalanceOf(user1), MINT_AMOUNT - burnAmount, "Scaled balance should be reduced");
        assertEq(debtToken.scaledTotalSupply(), MINT_AMOUNT - burnAmount, "Total scaled supply should be reduced");
    }
    
    /**
     * @dev Test burn with insufficient balance
     */
    function testBurnInsufficientBalance() public {
        // Setup: mint tokens first
        vm.startPrank(lendingPool);
        debtToken.mint(user1, MINT_AMOUNT, INITIAL_INDEX);
        vm.stopPrank();
        
        // Try to burn more than the user has
        vm.startPrank(lendingPool);
        vm.expectRevert("DebtToken: Insufficient balance");
        debtToken.burn(user1, MINT_AMOUNT * 2, INITIAL_INDEX);
        vm.stopPrank();
    }
    
    /**
     * @dev Test burn with max uint value to burn entire balance
     */
    function testBurnMax() public {
        // Setup: mint tokens first
        vm.startPrank(lendingPool);
        debtToken.mint(user1, MINT_AMOUNT, INITIAL_INDEX);
        vm.stopPrank();
        
        // Lending pool burns max amount (entire balance)
        vm.startPrank(lendingPool);
        uint256 actualBurned = debtToken.burn(user1, type(uint256).max, INITIAL_INDEX);
        vm.stopPrank();
        
        // Check results
        assertEq(actualBurned, MINT_AMOUNT, "Burn should return the entire amount");
        assertEq(debtToken.balanceOf(user1), 0, "User balance should be zero");
        assertEq(debtToken.scaledBalanceOf(user1), 0, "Scaled balance should be zero");
        assertEq(debtToken.scaledTotalSupply(), 0, "Total scaled supply should be zero");
    }
    
    /**
     * @dev Test scaling with different indices
     */
    function testScalingWithDifferentIndices() public {
        // Initial mint at index 1.0
        vm.startPrank(lendingPool);
        debtToken.mint(user1, MINT_AMOUNT, INITIAL_INDEX);
        vm.stopPrank();
        
        // Second mint at index 1.1 (10% interest accrued)
        uint256 higherIndex = INITIAL_INDEX * 110 / 100; // 1.1 in ray
        uint256 secondMintAmount = MINT_AMOUNT;
        
        vm.startPrank(lendingPool);
        debtToken.mint(user1, secondMintAmount, higherIndex);
        vm.stopPrank();
        
        // Check results
        uint256 expectedScaledAmount = MINT_AMOUNT + (secondMintAmount * INITIAL_INDEX / higherIndex);
        assertApproxEqRel(debtToken.scaledBalanceOf(user1), expectedScaledAmount, 0.0001e18, "Scaled balance should account for different indices");
        
        // The actual balance should be the sum of the minted amounts
        assertEq(debtToken.balanceOf(user1), MINT_AMOUNT + secondMintAmount, "Balance should be the sum of minted amounts");
    }
    
    /**
     * @dev Test transfers are disabled
     */
    function testTransfersDisabled() public {
        // Setup: mint tokens first
        vm.startPrank(lendingPool);
        debtToken.mint(user1, MINT_AMOUNT, INITIAL_INDEX);
        vm.stopPrank();
        
        // Try to transfer
        vm.startPrank(user1);
        vm.expectRevert("DebtToken: Transfer not allowed");
        debtToken.transfer(user2, MINT_AMOUNT / 2);
        vm.stopPrank();
        
        // Try transferFrom
        vm.prank(user1);
        debtToken.approve(user2, MINT_AMOUNT);
        
        vm.startPrank(user2);
        vm.expectRevert("DebtToken: Transfer not allowed");
        debtToken.transferFrom(user1, user2, MINT_AMOUNT / 2);
        vm.stopPrank();
    }
    
    /**
     * @dev Test approvals are disabled
     */
    function testApprovalsDisabled() public {
        vm.startPrank(user1);
        vm.expectRevert("DebtToken: Approve not allowed");
        debtToken.approve(user2, MINT_AMOUNT);
        vm.stopPrank();
    }

    /**
     * @dev Test getLendingPool function
     */
    function testGetLendingPool() public view {
        assertEq(debtToken.getLendingPool(), lendingPool, "getLendingPool should return the correct address");
    }
}
