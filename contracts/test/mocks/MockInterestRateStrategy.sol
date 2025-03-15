// SPDX-License-Identifier: MIT
pragma solidity 0.8.29;

import "@openzeppelin/contracts/access/Ownable.sol";
import "../../src/interfaces/IInterestRateStrategy.sol";

/**
 * @title MockInterestRateStrategy
 * @author DeFi Lending Platform
 * @notice Mock Interest Rate Strategy for testing purposes
 * @dev Implements the IInterestRateStrategy interface with configurable rates
 */
contract MockInterestRateStrategy is IInterestRateStrategy, Ownable {
    // Base variable borrow rate
    uint256 private _baseVariableBorrowRate;
    
    // Max variable borrow rate
    uint256 private _maxVariableBorrowRate;
    
    // Configurable rates to return
    uint256 private _liquidityRate;
    uint256 private _stableBorrowRate;
    uint256 private _variableBorrowRate;
    
    // Flag to control whether operations revert
    bool private _shouldRevert;
    
    // Track calls for testing verification
    uint256 public calculateInterestRatesCalls;

    /**
     * @dev Constructor to create a new MockInterestRateStrategy
     * @param baseVariableBorrowRate The base variable borrow rate
     * @param maxVariableBorrowRate The max variable borrow rate
     * @param initialLiquidityRate The initial liquidity rate to return
     * @param initialStableBorrowRate The initial stable borrow rate to return
     * @param initialVariableBorrowRate The initial variable borrow rate to return
     */
    constructor(
        uint256 baseVariableBorrowRate,
        uint256 maxVariableBorrowRate,
        uint256 initialLiquidityRate,
        uint256 initialStableBorrowRate,
        uint256 initialVariableBorrowRate
    ) Ownable(msg.sender) {
        _baseVariableBorrowRate = baseVariableBorrowRate;
        _maxVariableBorrowRate = maxVariableBorrowRate;
        _liquidityRate = initialLiquidityRate;
        _stableBorrowRate = initialStableBorrowRate;
        _variableBorrowRate = initialVariableBorrowRate;
    }

    /**
     * @inheritdoc IInterestRateStrategy
     */
    function getBaseVariableBorrowRate() external view override returns (uint256) {
        if (_shouldRevert) {
            revert("MockInterestRateStrategy: Forced failure");
        }
        
        return _baseVariableBorrowRate;
    }

    /**
     * @inheritdoc IInterestRateStrategy
     */
    function getMaxVariableBorrowRate() external view override returns (uint256) {
        if (_shouldRevert) {
            revert("MockInterestRateStrategy: Forced failure");
        }
        
        return _maxVariableBorrowRate;
    }

    /**
     * @inheritdoc IInterestRateStrategy
     */
    function calculateInterestRates(
        address reserve,
        uint256 availableLiquidity,
        uint256 totalBorrows,
        uint256 reserveFactor
    ) external view override returns (uint256 liquidityRate, uint256 stableBorrowRate, uint256 variableBorrowRate) {
        if (_shouldRevert) {
            revert("MockInterestRateStrategy: Forced failure");
        }
        
        // Can use these parameters to calculate dynamic rates for more complex tests
        // For now, just return the configured rates
        return (_liquidityRate, _stableBorrowRate, _variableBorrowRate);
    }

    /**
     * @dev External function to track calls for testing
     */
    function trackCalculateInterestRatesCall() external {
        calculateInterestRatesCalls++;
    }

    /**
     * @dev Sets the base variable borrow rate
     * @param baseVariableBorrowRate The new base variable borrow rate
     */
    function setBaseVariableBorrowRate(uint256 baseVariableBorrowRate) external onlyOwner {
        _baseVariableBorrowRate = baseVariableBorrowRate;
    }

    /**
     * @dev Sets the max variable borrow rate
     * @param maxVariableBorrowRate The new max variable borrow rate
     */
    function setMaxVariableBorrowRate(uint256 maxVariableBorrowRate) external onlyOwner {
        _maxVariableBorrowRate = maxVariableBorrowRate;
    }

    /**
     * @dev Sets the liquidity rate to return
     * @param liquidityRate The new liquidity rate
     */
    function setLiquidityRate(uint256 liquidityRate) external onlyOwner {
        _liquidityRate = liquidityRate;
    }

    /**
     * @dev Sets the stable borrow rate to return
     * @param stableBorrowRate The new stable borrow rate
     */
    function setStableBorrowRate(uint256 stableBorrowRate) external onlyOwner {
        _stableBorrowRate = stableBorrowRate;
    }

    /**
     * @dev Sets the variable borrow rate to return
     * @param variableBorrowRate The new variable borrow rate
     */
    function setVariableBorrowRate(uint256 variableBorrowRate) external onlyOwner {
        _variableBorrowRate = variableBorrowRate;
    }

    /**
     * @dev Sets all rates at once
     * @param liquidityRate The new liquidity rate
     * @param stableBorrowRate The new stable borrow rate
     * @param variableBorrowRate The new variable borrow rate
     */
    function setAllRates(
        uint256 liquidityRate,
        uint256 stableBorrowRate,
        uint256 variableBorrowRate
    ) external onlyOwner {
        _liquidityRate = liquidityRate;
        _stableBorrowRate = stableBorrowRate;
        _variableBorrowRate = variableBorrowRate;
    }

    /**
     * @dev Sets whether operations should revert
     * @param shouldRevert Whether operations should revert
     */
    function setShouldRevert(bool shouldRevert) external onlyOwner {
        _shouldRevert = shouldRevert;
    }

    /**
     * @dev Resets the call counter for testing
     */
    function resetCallCounter() external onlyOwner {
        calculateInterestRatesCalls = 0;
    }
}
