// SPDX-License-Identifier: MIT
pragma solidity 0.8.29;

import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "../../src/interfaces/ILendingPoolCore.sol";
import "../../src/libraries/ReserveLogic.sol";

/**
 * @title MockLendingPoolCore
 * @author DeFi Lending Platform
 * @notice Mock Lending Pool Core for testing purposes
 * @dev Implements the ILendingPoolCore interface
 */
contract MockLendingPoolCore is ILendingPoolCore, Ownable {
    using SafeERC20 for IERC20;

    // Reserves mapping (asset => reserve data)
    mapping(address => ReserveLogic.ReserveData) private _reserves;

    // List of reserves
    address[] private _reservesList;

    // Lending pool address (only this address can call certain functions)
    address private _lendingPool;

    // Flag to control whether operations revert
    bool private _shouldRevert;

    // Track calls for testing verification
    uint256 public initReserveCalls;
    uint256 public updateInterestRateStrategyCalls;
    uint256 public updateReserveStateCalls;
    uint256 public transferToUserCalls;
    uint256 public transferToReserveCalls;

    // Events
    event ReserveInitialized(
        address indexed asset, address indexed aToken, address indexed debtToken, address interestRateStrategy
    );

    /**
     * @dev Constructor to create a new MockLendingPoolCore
     * @param owner The owner of the contract
     */
    constructor(address owner) Ownable(msg.sender) {
        // Transfer ownership
        if (owner != msg.sender) {
            transferOwnership(owner);
        }
    }

    /**
     * @dev Only the lending pool can call this
     */
    modifier onlyLendingPool() {
        require(
            _lendingPool != address(0) && msg.sender == _lendingPool, "MockLendingPoolCore: Caller is not lending pool"
        );
        _;
    }

    /**
     * @inheritdoc ILendingPoolCore
     */
    function initReserve(address asset, address aToken, address debtToken, address interestRateStrategy)
        external
        override
        onlyOwner
    {
        if (_shouldRevert) {
            revert("MockLendingPoolCore: Forced failure");
        }

        initReserveCalls++;

        require(asset != address(0), "MockLendingPoolCore: Zero asset");
        require(aToken != address(0), "MockLendingPoolCore: Zero aToken");
        require(debtToken != address(0), "MockLendingPoolCore: Zero debtToken");
        require(interestRateStrategy != address(0), "MockLendingPoolCore: Zero interest rate strategy");

        _reserves[asset].aTokenAddress = aToken;
        _reserves[asset].debtTokenAddress = debtToken;
        _reserves[asset].interestRateStrategyAddress = interestRateStrategy;
        _reserves[asset].lastUpdateTimestamp = block.timestamp;

        _reservesList.push(asset);

        emit ReserveInitialized(asset, aToken, debtToken, interestRateStrategy);
    }

    /**
     * @inheritdoc ILendingPoolCore
     */
    function updateInterestRateStrategy(address asset, address interestRateStrategy) external override onlyOwner {
        if (_shouldRevert) {
            revert("MockLendingPoolCore: Forced failure");
        }

        updateInterestRateStrategyCalls++;

        require(asset != address(0), "MockLendingPoolCore: Zero asset");
        require(interestRateStrategy != address(0), "MockLendingPoolCore: Zero interest rate strategy");

        _reserves[asset].interestRateStrategyAddress = interestRateStrategy;
    }

    /**
     * @inheritdoc ILendingPoolCore
     */
    function updateReserveState(address asset) external override {
        if (_shouldRevert) {
            revert("MockLendingPoolCore: Forced failure");
        }

        updateReserveStateCalls++;

        require(asset != address(0), "MockLendingPoolCore: Zero asset");

        // In a real implementation, this would update the interest rates
        // For our mock, we just update the timestamp
        _reserves[asset].lastUpdateTimestamp = block.timestamp;
    }

    /**
     * @inheritdoc ILendingPoolCore
     */
    function transferToUser(address asset, address to, uint256 amount) external override onlyLendingPool {
        if (_shouldRevert) {
            revert("MockLendingPoolCore: Forced failure");
        }

        transferToUserCalls++;

        require(asset != address(0), "MockLendingPoolCore: Zero asset");
        require(to != address(0), "MockLendingPoolCore: Zero recipient");

        // Transfer the tokens to the user
        IERC20(asset).safeTransfer(to, amount);
    }

    /**
     * @inheritdoc ILendingPoolCore
     */
    function transferToReserve(address asset, address from, uint256 amount) external override onlyLendingPool {
        if (_shouldRevert) {
            revert("MockLendingPoolCore: Forced failure");
        }

        transferToReserveCalls++;

        require(asset != address(0), "MockLendingPoolCore: Zero asset");
        require(from != address(0), "MockLendingPoolCore: Zero sender");

        // Transfer the tokens from the user to the reserve
        IERC20(asset).safeTransferFrom(from, address(this), amount);
    }

    /**
     * @inheritdoc ILendingPoolCore
     */
    function getReserveData(address asset) external view override returns (ReserveLogic.ReserveData memory) {
        if (_shouldRevert) {
            revert("MockLendingPoolCore: Forced failure");
        }

        return _reserves[asset];
    }

    /**
     * @inheritdoc ILendingPoolCore
     */
    function getReservesList() external view override returns (address[] memory) {
        if (_shouldRevert) {
            revert("MockLendingPoolCore: Forced failure");
        }

        return _reservesList;
    }

    /**
     * @inheritdoc ILendingPoolCore
     */
    function enableBorrowingOnReserve(address asset, bool enabled) external override onlyOwner {
        if (_shouldRevert) {
            revert("MockLendingPoolCore: Forced failure");
        }

        require(asset != address(0), "MockLendingPoolCore: Zero asset");

        _reserves[asset].borrowingEnabled = enabled ? 1 : 0;
    }

    /**
     * @dev Sets the lending pool address
     * @param lendingPool The lending pool address
     */
    function setLendingPool(address lendingPool) external onlyOwner {
        require(lendingPool != address(0), "MockLendingPoolCore: Zero lending pool");
        _lendingPool = lendingPool;
    }

    /**
     * @dev Sets a reserve's liquidity and borrow rates
     * @param asset The asset address
     * @param liquidityRate The liquidity rate
     * @param borrowRate The borrow rate
     */
    function setReserveRates(address asset, uint256 liquidityRate, uint256 borrowRate) external onlyOwner {
        require(asset != address(0), "MockLendingPoolCore: Zero asset");

        _reserves[asset].currentLiquidityRate = liquidityRate;
        _reserves[asset].currentBorrowRate = borrowRate;
    }

    /**
     * @dev Sets a reserve's total liquidity and borrows
     * @param asset The asset address
     * @param totalLiquidity The total liquidity
     * @param totalBorrows The total borrows
     */
    function setReserveTotals(address asset, uint256 totalLiquidity, uint256 totalBorrows) external onlyOwner {
        require(asset != address(0), "MockLendingPoolCore: Zero asset");

        _reserves[asset].totalLiquidity = totalLiquidity;
        _reserves[asset].totalBorrows = totalBorrows;
    }

    /**
     * @dev Sets whether operations should revert
     * @param shouldRevert Whether operations should revert
     */
    function setShouldRevert(bool shouldRevert) external onlyOwner {
        _shouldRevert = shouldRevert;
    }

    /**
     * @dev Resets the call counters for testing
     */
    function resetCallCounters() external onlyOwner {
        initReserveCalls = 0;
        updateInterestRateStrategyCalls = 0;
        updateReserveStateCalls = 0;
        transferToUserCalls = 0;
        transferToReserveCalls = 0;
    }

    /**
     * @dev Manually set a complete reserve data structure
     * @param asset The asset address
     * @param reserveData The reserve data
     */
    function setReserveData(address asset, ReserveLogic.ReserveData calldata reserveData) external onlyOwner {
        require(asset != address(0), "MockLendingPoolCore: Zero asset");

        _reserves[asset] = reserveData;

        // Add to reserves list if not already there
        bool found = false;
        for (uint256 i = 0; i < _reservesList.length; i++) {
            if (_reservesList[i] == asset) {
                found = true;
                break;
            }
        }

        if (!found) {
            _reservesList.push(asset);
        }
    }
}
