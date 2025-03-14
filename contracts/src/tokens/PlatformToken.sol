// SPDX-License-Identifier: MIT
pragma solidity 0.8.29;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/token/ERC20/extensions/ERC20Burnable.sol";
import "@openzeppelin/contracts/access/AccessControl.sol";
import "@openzeppelin/contracts/utils/Pausable.sol";

/**
 * @title PlatformToken
 * @author DeFi Lending Platform
 * @notice The governance token for the lending platform
 * @dev ERC20 token with additional governance functions
 * Note: ERC20Snapshot was removed as it's not available in this version of OpenZeppelin
 */
contract PlatformToken is ERC20, ERC20Burnable, AccessControl, Pausable {
    // Roles
    bytes32 public constant SNAPSHOT_ROLE = keccak256("SNAPSHOT_ROLE");
    bytes32 public constant PAUSER_ROLE = keccak256("PAUSER_ROLE");
    bytes32 public constant MINTER_ROLE = keccak256("MINTER_ROLE");

    // Cap
    uint256 private immutable _cap;

    // Simple snapshot implementation
    uint256 private _currentSnapshotId;
    mapping(address => mapping(uint256 => uint256)) private _accountBalanceSnapshots;
    mapping(uint256 => uint256) private _totalSupplySnapshots;

    // Events
    event CapUpdated(uint256 oldCap, uint256 newCap);
    event Snapshot(uint256 id);

    /**
     * @dev Constructor
     * @param name The name of the token
     * @param symbol The symbol of the token
     * @param initialSupply The initial supply of the token
     * @param capAmount The maximum cap of the token
     * @param admin The address of the admin
     */
    constructor(string memory name, string memory symbol, uint256 initialSupply, uint256 capAmount, address admin)
        ERC20(name, symbol)
    {
        require(capAmount > 0, "PlatformToken: Cap is 0");
        require(initialSupply <= capAmount, "PlatformToken: Initial supply exceeds cap");
        require(admin != address(0), "PlatformToken: Admin cannot be zero address");

        _cap = capAmount;

        // Setup roles
        _grantRole(DEFAULT_ADMIN_ROLE, admin);
        _grantRole(SNAPSHOT_ROLE, admin);
        _grantRole(PAUSER_ROLE, admin);
        _grantRole(MINTER_ROLE, admin);

        // Mint initial supply
        _mint(admin, initialSupply);
    }

    /**
     * @dev Creates a new snapshot
     */
    function snapshot() external onlyRole(SNAPSHOT_ROLE) returns (uint256) {
        _currentSnapshotId += 1;
        uint256 currentId = _currentSnapshotId;

        _totalSupplySnapshots[currentId] = totalSupply();

        emit Snapshot(currentId);
        return currentId;
    }

    /**
     * @dev Get the balance of an account at a specific snapshot id
     */
    function balanceOfAt(address account, uint256 snapshotId) public view returns (uint256) {
        require(snapshotId <= _currentSnapshotId, "PlatformToken: Invalid snapshot id");

        // If we have a specific snapshot balance, return it
        if (_accountBalanceSnapshots[account][snapshotId] > 0) {
            return _accountBalanceSnapshots[account][snapshotId];
        }

        // Otherwise return current balance
        return balanceOf(account);
    }

    /**
     * @dev Get the total supply at a specific snapshot id
     */
    function totalSupplyAt(uint256 snapshotId) public view returns (uint256) {
        require(snapshotId <= _currentSnapshotId, "PlatformToken: Invalid snapshot id");

        // If we have a specific snapshot of total supply, return it
        if (_totalSupplySnapshots[snapshotId] > 0) {
            return _totalSupplySnapshots[snapshotId];
        }

        // Otherwise return current supply
        return totalSupply();
    }

    /**
     * @dev Pauses all token transfers
     */
    function pause() external onlyRole(PAUSER_ROLE) {
        _pause();
    }

    /**
     * @dev Unpauses all token transfers
     */
    function unpause() external onlyRole(PAUSER_ROLE) {
        _unpause();
    }

    /**
     * @dev Mints new tokens
     * @param to The address that will receive the minted tokens
     * @param amount The amount of tokens to mint
     */
    function mint(address to, uint256 amount) external onlyRole(MINTER_ROLE) {
        require(ERC20.totalSupply() + amount <= cap(), "PlatformToken: Cap exceeded");
        _mint(to, amount);
    }

    /**
     * @dev Returns the cap on the token's total supply
     * @return The cap
     */
    function cap() public view returns (uint256) {
        return _cap;
    }

    /**
     * @dev Updates the account balances and snapshots during transfers
     */
    function _update(address from, address to, uint256 amount) internal virtual override {
        if (!paused()) {
            super._update(from, to, amount);

            // Snapshot logic here
            if (_currentSnapshotId > 0) {
                if (from != address(0)) {
                    // not minting
                    _accountBalanceSnapshots[from][_currentSnapshotId] = balanceOf(from);
                }
                if (to != address(0)) {
                    // not burning
                    _accountBalanceSnapshots[to][_currentSnapshotId] = balanceOf(to);
                }
            }
        } else {
            revert("PlatformToken: token transfer while paused");
        }
    }
}
