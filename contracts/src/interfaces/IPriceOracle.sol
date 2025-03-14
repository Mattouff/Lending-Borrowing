// SPDX-License-Identifier: MIT
pragma solidity 0.8.29;

/**
 * @title IPriceOracle
 * @author DeFi Lending Platform
 * @notice Interface for the Price Oracle
 */
interface IPriceOracle {
    /**
     * @dev Emitted when a price source is updated
     * @param asset The address of the asset
     * @param source The price source
     */
    event AssetSourceUpdated(address indexed asset, address indexed source);

    /**
     * @dev Emitted when the ETH/USD price is updated
     * @param price The new ETH/USD price
     */
    event EthUsdPriceUpdated(uint256 price);

    /**
     * @dev Returns the asset price in ETH
     * @param asset The address of the asset
     * @return The asset price in ETH
     */
    function getAssetPrice(address asset) external view returns (uint256);

    /**
     * @dev Returns the prices of multiple assets in ETH
     * @param assets The addresses of the assets
     * @return The asset prices in ETH
     */
    function getAssetsPrices(address[] calldata assets) external view returns (uint256[] memory);

    /**
     * @dev Sets the price source for an asset
     * @param asset The address of the asset
     * @param source The price source
     */
    function setAssetSource(address asset, address source) external;

    /**
     * @dev Sets the price sources for multiple assets
     * @param assets The addresses of the assets
     * @param sources The price sources
     */
    function setAssetsSources(address[] calldata assets, address[] calldata sources) external;

    /**
     * @dev Sets the ETH/USD price
     * @param price The ETH/USD price
     */
    function setEthUsdPrice(uint256 price) external;

    /**
     * @dev Returns the ETH/USD price
     * @return The ETH/USD price
     */
    function getEthUsdPrice() external view returns (uint256);
}
