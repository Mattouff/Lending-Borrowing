// SPDX-License-Identifier: MIT
pragma solidity 0.8.29;

import "@openzeppelin-upgradeable/contracts/access/OwnableUpgradeable.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin-upgradeable/contracts/utils/ReentrancyGuardUpgradeable.sol";
import "@openzeppelin-upgradeable/contracts/proxy/utils/Initializable.sol";
import "@openzeppelin-upgradeable/contracts/proxy/utils/UUPSUpgradeable.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "../interfaces/ILendingPoolAddressesProvider.sol";
import "../interfaces/IAToken.sol";
import "../libraries/WadRayMath.sol";
import "../libraries/ReserveLogic.sol";

/**
 * @title LendingPoolCore
 * @author DeFi Lending Platform
 * @notice Manages the core logic of the lending pool
 * @dev Holds the reserves of assets and handles the transfers
 * This contract is upgradeable using the UUPS pattern
 */
contract LendingPoolCore is Initializable, OwnableUpgradeable, ReentrancyGuardUpgradeable, UUPSUpgradeable {
    using SafeERC20 for IERC20;
    using WadRayMath for uint256;
    using ReserveLogic for ReserveLogic.ReserveData;

    // Address provider
    ILendingPoolAddressesProvider private _addressesProvider;

    // Reserves mapping (asset => reserve data)
    mapping(address => ReserveLogic.ReserveData) private _reserves;

    // List of all reserves
    address[] private _reservesList;

    // Reserve factor (percentage of interest that goes to the protocol)
    uint256 public constant RESERVE_FACTOR = 0.1 * 1e18; // 10%

    // Events
    event ReserveInitialized(
        address indexed asset, address indexed aToken, address indexed debtToken, address interestRateStrategy
    );

    event ReserveUpdated(
        address indexed asset, uint256 liquidityRate, uint256 borrowRate, uint256 liquidityIndex, uint256 borrowIndex
    );

    // Storage gap for future upgrades
    uint256[50] private __gap;

    /**
     * @custom:oz-upgrades-unsafe-allow constructor
     */
    constructor() {
        _disableInitializers();
    }

    /**
     * @dev Initializer function
     * @param addressesProvider The address of the LendingPoolAddressesProvider
     */
    function initialize(address addressesProvider) external initializer {
        __Ownable_init(msg.sender);
        __ReentrancyGuard_init();
        __UUPSUpgradeable_init();

        require(addressesProvider != address(0), "LendingPoolCore: Invalid addresses provider");
        _addressesProvider = ILendingPoolAddressesProvider(addressesProvider);
    }

    /**
     * @dev Initializes a reserve
     * @param asset The address of the asset
     * @param aToken The address of the aToken
     * @param debtToken The address of the debt token
     * @param interestRateStrategy The address of the interest rate strategy
     */
    function initReserve(address asset, address aToken, address debtToken, address interestRateStrategy)
        external
        onlyOwner
    {
        require(asset != address(0), "LendingPoolCore: Invalid asset");
        require(aToken != address(0), "LendingPoolCore: Invalid aToken");
        require(debtToken != address(0), "LendingPoolCore: Invalid debtToken");
        require(interestRateStrategy != address(0), "LendingPoolCore: Invalid interest rate strategy");

        // Ensure the reserve is not already initialized
        require(_reserves[asset].aTokenAddress == address(0), "LendingPoolCore: Reserve already initialized");

        // Initialize the reserve
        _reserves[asset].aTokenAddress = aToken;
        _reserves[asset].debtTokenAddress = debtToken;
        _reserves[asset].interestRateStrategyAddress = interestRateStrategy;
        _reserves[asset].lastUpdateTimestamp = block.timestamp;

        // Add the asset to the list of reserves
        _reservesList.push(asset);

        emit ReserveInitialized(asset, aToken, debtToken, interestRateStrategy);
    }

    /**
     * @dev Updates the interest rate strategy of a reserve
     * @param asset The address of the asset
     * @param interestRateStrategy The address of the interest rate strategy
     */
    function updateInterestRateStrategy(address asset, address interestRateStrategy) external onlyOwner {
        require(asset != address(0), "LendingPoolCore: Invalid asset");
        require(interestRateStrategy != address(0), "LendingPoolCore: Invalid interest rate strategy");

        // Ensure the reserve is initialized
        require(_reserves[asset].aTokenAddress != address(0), "LendingPoolCore: Reserve not initialized");

        // Update the interest rate strategy
        _reserves[asset].interestRateStrategyAddress = interestRateStrategy;
    }

    /**
     * @dev Updates the reserve state with accumulated interest
     * @param asset The address of the asset
     */
    function updateReserveState(address asset) external {
        require(asset != address(0), "LendingPoolCore: Invalid asset");

        // Ensure the reserve is initialized
        require(_reserves[asset].aTokenAddress != address(0), "LendingPoolCore: Reserve not initialized");

        // Update the reserve state
        _reserves[asset].updateReserveState();
    }

    /**
     * @dev Transfers an amount of asset to the user
     * @param asset The address of the asset
     * @param to The address of the recipient
     * @param amount The amount to transfer
     */
    function transferToUser(address asset, address to, uint256 amount) external onlyLendingPool {
        require(asset != address(0), "LendingPoolCore: Invalid asset");
        require(to != address(0), "LendingPoolCore: Invalid recipient");
        require(amount > 0, "LendingPoolCore: Invalid amount");

        // Transfer the asset to the user
        IERC20(asset).safeTransfer(to, amount);
    }

    /**
     * @dev Transfers an amount of asset from the user to the reserve
     * @param asset The address of the asset
     * @param from The address of the sender
     * @param amount The amount to transfer
     */
    function transferToReserve(address asset, address from, uint256 amount) external onlyLendingPool {
        require(asset != address(0), "LendingPoolCore: Invalid asset");
        require(from != address(0), "LendingPoolCore: Invalid sender");
        require(amount > 0, "LendingPoolCore: Invalid amount");

        // Transfer the asset from the user to the reserve
        IERC20(asset).safeTransferFrom(from, address(this), amount);
    }

    /**
     * @dev Gets the reserve data for an asset
     * @param asset The address of the asset
     * @return The reserve data
     */
    function getReserveData(address asset) external view returns (ReserveLogic.ReserveData memory) {
        return _reserves[asset];
    }

    /**
     * @dev Gets the list of all reserves
     * @return The list of all reserves
     */
    function getReservesList() external view returns (address[] memory) {
        return _reservesList;
    }

    /**
     * @dev Ensures the caller is the lending pool
     */
    modifier onlyLendingPool() {
        require(msg.sender == _addressesProvider.getLendingPool(), "LendingPoolCore: Caller must be lending pool");
        _;
    }

    /**
     * @dev Function that should revert when msg.sender is not authorized to upgrade the contract
     * @param newImplementation The address of the new implementation
     */
    function _authorizeUpgrade(address newImplementation) internal override onlyOwner { }
}
