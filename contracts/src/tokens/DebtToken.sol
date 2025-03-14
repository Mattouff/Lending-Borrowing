// SPDX-License-Identifier: MIT
pragma solidity 0.8.29;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "../libraries/WadRayMath.sol";

/**
 * @title DebtToken
 * @author DeFi Lending Platform
 * @notice Implementation of the debt token for the DeFi lending platform
 * @dev The token represents the debt of users in the lending pool
 */
contract DebtToken is ERC20, Ownable {
    using WadRayMath for uint256;

    // Address of the underlying asset
    address private immutable _underlyingAsset;

    // Address of the lending pool
    address private immutable _lendingPool;

    // Maps from user address to their scaled debt balance
    mapping(address => uint256) private _userScaledBalances;

    // Total scaled debt
    uint256 private _totalScaledSupply;

    // Events
    event Mint(address indexed user, uint256 amount, uint256 index);
    event Burn(address indexed user, uint256 amount, uint256 index);

    /**
     * @dev Constructor
     * @param underlyingAsset The address of the underlying asset
     * @param lendingPool The address of the lending pool
     * @param name The name of the token
     * @param symbol The symbol of the token
     */
    constructor(address underlyingAsset, address lendingPool, string memory name, string memory symbol)
        ERC20(name, symbol)
        Ownable(msg.sender)
    {
        require(underlyingAsset != address(0), "DebtToken: Underlying asset cannot be zero address");
        require(lendingPool != address(0), "DebtToken: Lending pool cannot be zero address");

        _underlyingAsset = underlyingAsset;
        _lendingPool = lendingPool;

        transferOwnership(lendingPool);
    }

    /**
     * @dev Ensures the caller is the lending pool
     */
    modifier onlyLendingPool() {
        require(msg.sender == _lendingPool, "DebtToken: Caller must be lending pool");
        _;
    }

    /**
     * @dev Mints debt tokens to a user
     * @param user The address of the user
     * @param amount The amount to mint
     * @param index The current debt index
     */
    function mint(address user, uint256 amount, uint256 index) external onlyLendingPool {
        require(user != address(0), "DebtToken: Cannot mint to zero address");

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
     * @dev Burns debt tokens from a user
     * @param user The address of the user
     * @param amount The amount to burn
     * @param index The current debt index
     * @return The actual amount burned
     */
    function burn(address user, uint256 amount, uint256 index) external onlyLendingPool returns (uint256) {
        require(user != address(0), "DebtToken: Cannot burn from zero address");

        uint256 userBalance = balanceOf(user);

        // If amount is uint256 max, burn the entire balance
        uint256 amountToBurn = amount == type(uint256).max ? userBalance : amount;
        require(userBalance >= amountToBurn, "DebtToken: Insufficient balance");

        // Calculate scaled amount
        uint256 scaledAmount = amountToBurn.wadDiv(index);

        // Update user's scaled balance
        _userScaledBalances[user] = _userScaledBalances[user] - scaledAmount;

        // Update total scaled supply
        _totalScaledSupply = _totalScaledSupply - scaledAmount;

        // Burn tokens
        _burn(user, amountToBurn);

        emit Burn(user, amountToBurn, index);

        return amountToBurn;
    }

    /**
     * @dev Returns the address of the underlying asset
     * @return The address of the underlying asset
     */
    function getUnderlyingAssetAddress() external view returns (address) {
        return _underlyingAsset;
    }

    /**
     * @dev Returns the address of the lending pool
     * @return The address of the lending pool
     */
    function getLendingPool() external view returns (address) {
        return _lendingPool;
    }

    /**
     * @dev Returns the scaled balance of a user
     * @param user The address of the user
     * @return The scaled balance of the user
     */
    function scaledBalanceOf(address user) external view returns (uint256) {
        return _userScaledBalances[user];
    }

    /**
     * @dev Returns the total scaled supply
     * @return The total scaled supply
     */
    function scaledTotalSupply() external view returns (uint256) {
        return _totalScaledSupply;
    }

    /**
     * @dev Disable transfers for debt tokens
     */
    function transfer(address, uint256) public pure override returns (bool) {
        revert("DebtToken: Transfer not allowed");
    }

    /**
     * @dev Disable transfer from for debt tokens
     */
    function transferFrom(address, address, uint256) public pure override returns (bool) {
        revert("DebtToken: Transfer not allowed");
    }

    /**
     * @dev Disable approvals for debt tokens
     */
    function approve(address, uint256) public pure override returns (bool) {
        revert("DebtToken: Approve not allowed");
    }
}
