// SPDX-License-Identifier: MIT
pragma solidity 0.8.29;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "../../src/interfaces/IAToken.sol";

/**
 * @title MockAToken
 * @author DeFi Lending Platform
 * @notice Mock AToken for testing purposes
 * @dev Implements the IAToken interface
 */
contract MockAToken is IAToken, ERC20, Ownable {
    // Underlying asset address
    address private _underlyingAsset;
    
    // Scaled balances for users
    mapping(address => uint256) private _scaledBalances;
    
    // Total scaled supply
    uint256 private _totalScaledSupply;
    
    // Flag to control whether operations revert
    bool private _shouldRevert;
    
    // Track calls for testing verification
    uint256 public mintCalls;
    uint256 public burnCalls;
    uint256 public transferUnderlyingToCalls;

    /**
     * @dev Constructor to create a new MockAToken
     * @param name The name of the token
     * @param symbol The symbol of the token
     * @param underlyingAsset The address of the underlying asset
     * @param owner The owner of the token
     */
    constructor(
        string memory name,
        string memory symbol,
        address underlyingAsset,
        address owner
    ) ERC20(name, symbol) Ownable(msg.sender) {
        require(underlyingAsset != address(0), "MockAToken: Zero address");
        
        _underlyingAsset = underlyingAsset;
        
        // Transfer ownership
        if (owner != msg.sender) {
            transferOwnership(owner);
        }
    }

    /**
     * @inheritdoc IAToken
     */
    function mint(address user, uint256 amount, uint256 index) external override {
        if (_shouldRevert) {
            revert("MockAToken: Forced failure");
        }
        
        mintCalls++;
        
        uint256 scaledAmount = amount / index;
        _scaledBalances[user] += scaledAmount;
        _totalScaledSupply += scaledAmount;
        
        _mint(user, amount);
        
        emit Mint(user, amount, index);
    }

    /**
     * @inheritdoc IAToken
     */
    function burn(address user, uint256 amount, uint256 index) external override returns (uint256) {
        if (_shouldRevert) {
            revert("MockAToken: Forced failure");
        }
        
        burnCalls++;
        
        require(balanceOf(user) >= amount, "MockAToken: Insufficient balance");
        
        uint256 scaledAmount = amount / index;
        _scaledBalances[user] -= scaledAmount;
        _totalScaledSupply -= scaledAmount;
        
        _burn(user, amount);
        
        emit Burn(user, amount, index);
        
        return amount;
    }

    /**
     * @inheritdoc IAToken
     */
    function transferUnderlyingTo(address user, uint256 amount) external override returns (uint256) {
        if (_shouldRevert) {
            revert("MockAToken: Forced failure");
        }
        
        transferUnderlyingToCalls++;
        
        // In a real implementation, this would transfer the underlying asset
        // For our mock, we just track the call
        
        return amount;
    }

    /**
     * @inheritdoc IAToken
     */
    function getUnderlyingAssetAddress() external view override returns (address) {
        return _underlyingAsset;
    }

    /**
     * @inheritdoc IAToken
     */
    function scaledBalanceOf(address user) external view override returns (uint256) {
        return _scaledBalances[user];
    }

    /**
     * @inheritdoc IAToken
     */
    function scaledTotalSupply() external view override returns (uint256) {
        return _totalScaledSupply;
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
        transferUnderlyingToCalls = 0;
    }
}
