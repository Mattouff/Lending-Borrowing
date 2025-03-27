// SPDX-License-Identifier: MIT
pragma solidity 0.8.29;

import "forge-std/Test.sol";
import "../../src/interfaces/ILendingPool.sol";
import "../../src/interfaces/ILendingPoolCore.sol";
import "../../src/interfaces/IAToken.sol";
import "../../src/interfaces/IDebtToken.sol";
import "../../src/interfaces/IPriceOracle.sol";
import "../mocks/MockERC20.sol";
import "../mocks/MockAToken.sol";
import "../mocks/MockDebtToken.sol";

/**
 * @title TestingHelper
 * @author DeFi Lending Platform
 * @notice Helper contract for testing the lending platform
 * @dev Provides utility functions for common testing operations
 */
contract TestingHelper is Test {
    // Constants for testing
    uint256 public constant WAD = 1e18;
    uint256 public constant RAY = 1e27;
    uint256 public constant DEFAULT_COLLATERAL_RATIO = 75 * 1e16; // 75%
    uint256 public constant DEFAULT_LIQUIDATION_THRESHOLD = 80 * 1e16; // 80%
    uint256 public constant DEFAULT_LIQUIDATION_BONUS = 110 * 1e16; // 110%
    uint256 public constant DEFAULT_INTEREST_RATE = 5 * 1e16; // 5%

    /**
     * @dev Creates a new ERC20 token for testing
     * @param name The name of the token
     * @param symbol The symbol of the token
     * @param initialSupply The initial supply to mint
     * @param recipient The address to receive the initial supply
     * @param decimals The number of decimals for the token
     * @return The address of the created token
     */
    function createToken(
        string memory name,
        string memory symbol,
        uint256 initialSupply,
        address recipient,
        uint8 decimals
    ) public returns (address) {
        MockERC20 token = new MockERC20(name, symbol, initialSupply, decimals, recipient);
        return address(token);
    }

    /**
     * @dev Creates a DAI token with 18 decimals
     * @param initialSupply The initial supply to mint
     * @param recipient The address to receive the initial supply
     * @return The address of the DAI token
     */
    function createDAI(uint256 initialSupply, address recipient) public returns (address) {
        return createToken("DAI Stablecoin", "DAI", initialSupply, recipient, 18);
    }

    /**
     * @dev Creates a USDC token with 6 decimals
     * @param initialSupply The initial supply to mint
     * @param recipient The address to receive the initial supply
     * @return The address of the USDC token
     */
    function createUSDC(uint256 initialSupply, address recipient) public returns (address) {
        return createToken("USD Coin", "USDC", initialSupply, recipient, 6);
    }

    /**
     * @dev Creates a WETH token with 18 decimals
     * @param initialSupply The initial supply to mint
     * @param recipient The address to receive the initial supply
     * @return The address of the WETH token
     */
    function createWETH(uint256 initialSupply, address recipient) public returns (address) {
        return createToken("Wrapped Ether", "WETH", initialSupply, recipient, 18);
    }

    /**
     * @dev Creates an aToken for testing
     * @param name The name of the token
     * @param symbol The symbol of the token
     * @param underlyingAsset The address of the underlying asset
     * @param owner The owner of the token
     * @return The address of the created aToken
     */
    function createAToken(string memory name, string memory symbol, address underlyingAsset, address owner)
        public
        returns (address)
    {
        MockAToken aToken = new MockAToken(name, symbol, underlyingAsset, owner);
        return address(aToken);
    }

    /**
     * @dev Creates a debt token for testing
     * @param name The name of the token
     * @param symbol The symbol of the token
     * @param underlyingAsset The address of the underlying asset
     * @param lendingPool The address of the lending pool
     * @param owner The owner of the token
     * @return The address of the created debt token
     */
    function createDebtToken(
        string memory name,
        string memory symbol,
        address underlyingAsset,
        address lendingPool,
        address owner
    ) public returns (address) {
        MockDebtToken debtToken = new MockDebtToken(name, symbol, underlyingAsset, lendingPool, owner);
        return address(debtToken);
    }

    /**
     * @dev Approve tokens for transfer to another address
     * @param token The token to approve
     * @param from The address to approve from
     * @param to The address to approve to
     * @param amount The amount to approve
     */
    function approveTokens(address token, address from, address to, uint256 amount) public {
        vm.prank(from);
        MockERC20(token).approve(to, amount);
    }

    /**
     * @dev Mint tokens to an address
     * @param token The token to mint
     * @param to The address to mint to
     * @param amount The amount to mint
     */
    function mintTokens(address token, address to, uint256 amount) public {
        address owner = MockERC20(token).owner();
        vm.prank(owner);
        MockERC20(token).mint(to, amount);
    }

    /**
     * @dev Helper to perform a deposit into the lending pool
     * @param lendingPool The lending pool address
     * @param user The user address
     * @param asset The asset address
     * @param amount The amount to deposit
     */
    function deposit(address lendingPool, address user, address asset, uint256 amount) public {
        // Approve tokens first
        approveTokens(asset, user, lendingPool, amount);

        // Perform deposit
        vm.prank(user);
        ILendingPool(lendingPool).deposit(asset, amount, 0);
    }

    /**
     * @dev Helper to perform a borrow from the lending pool
     * @param lendingPool The lending pool address
     * @param user The user address
     * @param asset The asset address
     * @param amount The amount to borrow
     */
    function borrow(address lendingPool, address user, address asset, uint256 amount) public {
        vm.prank(user);
        ILendingPool(lendingPool).borrow(asset, amount, 1, 0); // Variable rate
    }

    /**
     * @dev Helper to perform a repayment to the lending pool
     * @param lendingPool The lending pool address
     * @param user The user address
     * @param asset The asset address
     * @param amount The amount to repay
     */
    function repay(address lendingPool, address user, address asset, uint256 amount) public {
        // Approve tokens first
        approveTokens(asset, user, lendingPool, amount);

        // Perform repayment
        vm.prank(user);
        ILendingPool(lendingPool).repay(asset, amount, 1); // Variable rate
    }

    /**
     * @dev Helper to perform a liquidation
     * @param lendingPool The lending pool address
     * @param liquidator The liquidator address
     * @param borrower The borrower address
     * @param collateralAsset The collateral asset address
     * @param debtAsset The debt asset address
     * @param debtToCover The amount of debt to cover
     */
    function liquidate(
        address lendingPool,
        address liquidator,
        address borrower,
        address collateralAsset,
        address debtAsset,
        uint256 debtToCover
    ) public {
        // Approve tokens first
        approveTokens(debtAsset, liquidator, lendingPool, debtToCover);

        // Perform liquidation
        vm.prank(liquidator);
        ILendingPool(lendingPool).liquidationCall(collateralAsset, debtAsset, borrower, debtToCover);
    }

    /**
     * @dev Helper to check a user's account data
     * @param lendingPool The lending pool address
     * @param user The user address
     * @return collateral The total collateral in ETH
     * @return debt The total debt in ETH
     * @return availableBorrows The available borrows in ETH
     * @return liquidationThreshold The liquidation threshold
     * @return ltv The loan to value ratio
     * @return healthFactor The health factor
     */
    function getUserData(address lendingPool, address user)
        public
        view
        returns (
            uint256 collateral,
            uint256 debt,
            uint256 availableBorrows,
            uint256 liquidationThreshold,
            uint256 ltv,
            uint256 healthFactor
        )
    {
        return ILendingPool(lendingPool).getUserAccountData(user);
    }

    /**
     * @dev Converts a RAY value to a human-readable percentage
     * @param ray The RAY value to convert
     * @return The percentage value (e.g., 5.0% for 0.05 * 1e27)
     */
    function rayToPercent(uint256 ray) public pure returns (uint256) {
        return (ray * 100) / RAY;
    }

    /**
     * @dev Converts a WAD value to a human-readable percentage
     * @param wad The WAD value to convert
     * @return The percentage value (e.g., 5.0% for 0.05 * 1e18)
     */
    function wadToPercent(uint256 wad) public pure returns (uint256) {
        return (wad * 100) / WAD;
    }

    /**
     * @dev Helper function to advance time and trigger interest accrual
     * @param secondsToAdvance The number of seconds to advance
     */
    function advanceTimeAndUpdateReserve(address lendingPool, address asset, uint256 secondsToAdvance) public {
        // Advance time
        vm.warp(block.timestamp + secondsToAdvance);

        // Trigger reserve update
        vm.prank(address(this));
        ILendingPool(lendingPool).deposit(asset, 0, 0); // Zero deposit to trigger update
    }
}
