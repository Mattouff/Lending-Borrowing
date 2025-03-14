// SPDX-License-Identifier: MIT
pragma solidity 0.8.29;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/token/ERC20/extensions/IERC20Metadata.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import "../interfaces/IAToken.sol";
import "../libraries/WadRayMath.sol";

/**
 * @title AToken
 * @author DeFi Lending Platform
 * @notice Implementation of the interest bearing token for the DeFi lending platform
 * @dev The token is pegged 1:1 to the value of the underlying asset and represents a deposit in the lending pool
 */
contract AToken is ERC20, Ownable, ReentrancyGuard, IAToken {
    using SafeERC20 for IERC20;
    using WadRayMath for uint256;

    // Underlying asset
    IERC20 private _underlyingAsset;

    // Address of the lending pool
    address private immutable _lendingPool;

    // Maps from user address to their scaled balance
    mapping(address => uint256) private _userScaledBalances;

    // Total scaled supply
    uint256 private _totalScaledSupply;

    /**
     * @dev Constructor
     * @param underlying The address of the underlying asset
     * @param lendingPool The address of the lending pool
     * @param name The name of the token
     * @param symbol The symbol of the token
     */
    constructor(address underlying, address lendingPool, string memory name, string memory symbol)
        ERC20(name, symbol)
        Ownable(msg.sender)
    {
        require(underlying != address(0), "AToken: Underlying asset cannot be zero address");
        require(lendingPool != address(0), "AToken: Lending pool cannot be zero address");

        _underlyingAsset = IERC20(underlying);
        _lendingPool = lendingPool;

        // Derive the number of decimals from the underlying asset
        uint8 underlyingDecimals = 18;
        if (address(_underlyingAsset) != address(0)) {
            try IERC20Metadata(address(_underlyingAsset)).decimals() returns (uint8 decimals) {
                underlyingDecimals = decimals;
            } catch { }
        }

        transferOwnership(lendingPool);
    }

    /**
     * @dev Ensures the caller is the lending pool
     */
    modifier onlyLendingPool() {
        require(msg.sender == _lendingPool, "AToken: Caller must be lending pool");
        _;
    }

    /**
     * @inheritdoc IAToken
     */
    function mint(address user, uint256 amount, uint256 index) external override onlyLendingPool {
        require(user != address(0), "AToken: Cannot mint to zero address");

        // Calculate scaled amount
        uint256 scaledAmount = amount.wadDiv(index);

        // Update user's scaled balance
        _userScaledBalances[user] = _userScaledBalances[user] + scaledAmount;

        // Update total scaled supply
        _totalScaledSupply = _totalScaledSupply + scaledAmount;

        // Mint tokens
        _mint(user, amount);

        emit Mint(user, amount, index);
    }

    /**
     * @inheritdoc IAToken
     */
    function burn(address user, uint256 amount, uint256 index) external override onlyLendingPool returns (uint256) {
        require(user != address(0), "AToken: Cannot burn from zero address");

        uint256 userBalance = balanceOf(user);
        require(userBalance >= amount, "AToken: Insufficient balance");

        // Calculate scaled amount
        uint256 scaledAmount = amount.wadDiv(index);

        // Update user's scaled balance
        _userScaledBalances[user] = _userScaledBalances[user] - scaledAmount;

        // Update total scaled supply
        _totalScaledSupply = _totalScaledSupply - scaledAmount;

        // Burn tokens
        _burn(user, amount);

        emit Burn(user, amount, index);

        return amount;
    }

    /**
     * @inheritdoc IAToken
     */
    function transferUnderlyingTo(address user, uint256 amount) external override onlyLendingPool returns (uint256) {
        require(user != address(0), "AToken: Cannot transfer to zero address");

        // Transfer underlying token to user
        _underlyingAsset.safeTransfer(user, amount);

        return amount;
    }

    /**
     * @inheritdoc IAToken
     */
    function getUnderlyingAssetAddress() external view override returns (address) {
        return address(_underlyingAsset);
    }

    /**
     * @inheritdoc IAToken
     */
    function scaledBalanceOf(address user) external view override returns (uint256) {
        return _userScaledBalances[user];
    }

    /**
     * @inheritdoc IAToken
     */
    function scaledTotalSupply() external view override returns (uint256) {
        return _totalScaledSupply;
    }

    /**
     * @dev Transfer tokens to another address
     * @param to Address to transfer tokens to
     * @param amount Amount of tokens to transfer
     * @return Whether transfer was successful
     */
    function transfer(address to, uint256 amount) public override(ERC20, IERC20) returns (bool) {
        // Track the transfer in terms of scaled balances as well
        uint256 scaledAmount = amount.wadDiv(WadRayMath.wad()); // Simplified, should use current index

        _userScaledBalances[msg.sender] = _userScaledBalances[msg.sender] - scaledAmount;
        _userScaledBalances[to] = _userScaledBalances[to] + scaledAmount;

        return super.transfer(to, amount);
    }

    /**
     * @dev Transfer tokens from one address to another
     * @param from Address to transfer tokens from
     * @param to Address to transfer tokens to
     * @param amount Amount of tokens to transfer
     * @return Whether transfer was successful
     */
    function transferFrom(address from, address to, uint256 amount) public override(ERC20, IERC20) returns (bool) {
        // Track the transfer in terms of scaled balances as well
        uint256 scaledAmount = amount.wadDiv(WadRayMath.wad()); // Simplified, should use current index

        _userScaledBalances[from] = _userScaledBalances[from] - scaledAmount;
        _userScaledBalances[to] = _userScaledBalances[to] + scaledAmount;

        return super.transferFrom(from, to, amount);
    }
}
