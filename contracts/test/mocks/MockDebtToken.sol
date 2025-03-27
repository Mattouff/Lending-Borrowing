// SPDX-License-Identifier: MIT
pragma solidity 0.8.29;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "../../src/interfaces/IDebtToken.sol";

/**
 * @title MockDebtToken
 * @author DeFi Lending Platform
 * @notice Mock DebtToken for testing purposes
 * @dev Implements the IDebtToken interface
 */
contract MockDebtToken is IDebtToken, ERC20, Ownable {
    // Underlying asset address
    address private _underlyingAsset;

    // Lending pool address
    address private _lendingPool;

    // Scaled balances for users
    mapping(address => uint256) private _scaledBalances;

    // Total scaled supply
    uint256 private _totalScaledSupply;

    // Flag to control whether operations revert
    bool private _shouldRevert;

    // Track calls for testing verification
    uint256 public mintCalls;
    uint256 public burnCalls;

    /**
     * @dev Constructor to create a new MockDebtToken
     * @param name The name of the token
     * @param symbol The symbol of the token
     * @param underlyingAsset The address of the underlying asset
     * @param lendingPool The address of the lending pool
     * @param owner The owner of the token
     */
    constructor(string memory name, string memory symbol, address underlyingAsset, address lendingPool, address owner)
        ERC20(name, symbol)
        Ownable(msg.sender)
    {
        require(underlyingAsset != address(0), "MockDebtToken: Zero underlying asset");
        require(lendingPool != address(0), "MockDebtToken: Zero lending pool");

        _underlyingAsset = underlyingAsset;
        _lendingPool = lendingPool;

        // Transfer ownership
        if (owner != msg.sender) {
            transferOwnership(owner);
        }
    }

    /**
     * @inheritdoc IDebtToken
     */
    function mint(address user, uint256 amount, uint256 index) external override {
        if (_shouldRevert) {
            revert("MockDebtToken: Forced failure");
        }

        require(msg.sender == _lendingPool, "MockDebtToken: Only lending pool");
        mintCalls++;

        uint256 scaledAmount = amount / index;
        _scaledBalances[user] += scaledAmount;
        _totalScaledSupply += scaledAmount;

        _mint(user, amount);

        emit Mint(user, amount, index);
    }

    /**
     * @inheritdoc IDebtToken
     */
    function burn(address user, uint256 amount, uint256 index) external override returns (uint256) {
        if (_shouldRevert) {
            revert("MockDebtToken: Forced failure");
        }

        require(msg.sender == _lendingPool, "MockDebtToken: Only lending pool");
        burnCalls++;

        uint256 currentBalance = balanceOf(user);
        uint256 amountToBurn = amount == type(uint256).max ? currentBalance : amount;

        require(currentBalance >= amountToBurn, "MockDebtToken: Insufficient balance");

        uint256 scaledAmount = amountToBurn / index;
        _scaledBalances[user] -= scaledAmount;
        _totalScaledSupply -= scaledAmount;

        _burn(user, amountToBurn);

        emit Burn(user, amountToBurn, index);

        return amountToBurn;
    }

    /**
     * @inheritdoc IDebtToken
     */
    function getUnderlyingAssetAddress() external view override returns (address) {
        return _underlyingAsset;
    }

    /**
     * @inheritdoc IDebtToken
     */
    function scaledBalanceOf(address user) external view override returns (uint256) {
        return _scaledBalances[user];
    }

    /**
     * @inheritdoc IDebtToken
     */
    function scaledTotalSupply() external view override returns (uint256) {
        return _totalScaledSupply;
    }

    /**
     * @dev Returns the lending pool address
     * @return The address of the lending pool
     */
    function getLendingPool() external view returns (address) {
        return _lendingPool;
    }

    /**
     * @dev Sets the lending pool address
     * @param lendingPool The new lending pool address
     */
    function setLendingPool(address lendingPool) external onlyOwner {
        require(lendingPool != address(0), "MockDebtToken: Zero lending pool");
        _lendingPool = lendingPool;
    }

    /**
     * @dev Sets whether operations should revert
     * @param shouldRevert Whether operations should revert
     */
    function setShouldRevert(bool shouldRevert) external onlyOwner {
        _shouldRevert = shouldRevert;
    }

    /**
     * @dev Sets the scaled balance for a user
     * @param user The user address
     * @param scaledBalance The scaled balance
     */
    function setScaledBalance(address user, uint256 scaledBalance) external onlyOwner {
        _scaledBalances[user] = scaledBalance;
    }

    /**
     * @dev Sets the total scaled supply
     * @param totalScaledSupply The total scaled supply
     */
    function setTotalScaledSupply(uint256 totalScaledSupply) external onlyOwner {
        _totalScaledSupply = totalScaledSupply;
    }

    /**
     * @dev Resets the call counters for testing
     */
    function resetCallCounters() external onlyOwner {
        mintCalls = 0;
        burnCalls = 0;
    }

    /**
     * @dev Disables transfers for debt tokens
     */
    function transfer(address, uint256) public pure override(ERC20, IERC20) returns (bool) {
        revert("MockDebtToken: Transfers not allowed");
    }

    /**
     * @dev Disables transfer from for debt tokens
     */
    function transferFrom(address, address, uint256) public pure override(ERC20, IERC20) returns (bool) {
        revert("MockDebtToken: Transfers not allowed");
    }
}
