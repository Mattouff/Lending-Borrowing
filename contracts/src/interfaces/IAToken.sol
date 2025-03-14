// SPDX-License-Identifier: MIT
pragma solidity 0.8.29;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";

/**
 * @title IAToken
 * @author DeFi Lending Platform
 * @notice Interface for the AToken contract
 */
interface IAToken is IERC20 {
    /**
     * @dev Emitted when tokens are minted
     * @param user The address of the user
     * @param amount The amount minted
     * @param index The index of the user
     */
    event Mint(address indexed user, uint256 amount, uint256 index);

    /**
     * @dev Emitted when tokens are burned
     * @param user The address of the user
     * @param amount The amount burned
     * @param index The index of the user
     */
    event Burn(address indexed user, uint256 amount, uint256 index);

    /**
     * @dev Mints aTokens to the user
     * @param user The address of the user
     * @param amount The amount to mint
     * @param index The index of the user
     */
    function mint(address user, uint256 amount, uint256 index) external;

    /**
     * @dev Burns aTokens from the user
     * @param user The address of the user
     * @param amount The amount to burn
     * @param index The index of the user
     * @return The amount burned
     */
    function burn(address user, uint256 amount, uint256 index) external returns (uint256);

    /**
     * @dev Transfers the underlying asset to the user
     * @param user The address of the user
     * @param amount The amount to transfer
     * @return The amount transferred
     */
    function transferUnderlyingTo(address user, uint256 amount) external returns (uint256);

    /**
     * @dev Returns the underlying asset address
     * @return The underlying asset address
     */
    function getUnderlyingAssetAddress() external view returns (address);

    /**
     * @dev Returns the scaled balance of the user
     * @param user The address of the user
     * @return The scaled balance
     */
    function scaledBalanceOf(address user) external view returns (uint256);

    /**
     * @dev Returns the scaled total supply
     * @return The scaled total supply
     */
    function scaledTotalSupply() external view returns (uint256);
}
