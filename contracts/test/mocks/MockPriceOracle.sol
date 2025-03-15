// SPDX-License-Identifier: MIT
pragma solidity 0.8.29;

import "@openzeppelin/contracts/access/Ownable.sol";
import "../../src/interfaces/IPriceOracle.sol";

/**
 * @title MockPriceOracle
 * @author DeFi Lending Platform
 * @notice Mock Price Oracle for testing purposes
 * @dev Implements the IPriceOracle interface with configurable prices
 */
contract MockPriceOracle is IPriceOracle, Ownable {
    // Mapping from asset address to price in ETH
    mapping(address => uint256) private _assetPrices;
    
    // ETH/USD price
    uint256 private _ethUsdPrice;
    
    // Flag to control whether operations revert
    bool private _shouldRevert;

    /**
     * @dev Constructor to create a new MockPriceOracle
     * @param owner The owner of the oracle
     * @param ethUsdPrice The initial ETH/USD price
     */
    constructor(
        address owner,
        uint256 ethUsdPrice
    ) Ownable(msg.sender) {
        _ethUsdPrice = ethUsdPrice;
        
        // Transfer ownership
        if (owner != msg.sender) {
            transferOwnership(owner);
        }
    }

    /**
     * @inheritdoc IPriceOracle
     */
    function getAssetPrice(address asset) external view override returns (uint256) {
        if (_shouldRevert) {
            revert("MockPriceOracle: Forced failure");
        }
        
        require(_assetPrices[asset] > 0, "MockPriceOracle: Price not set for asset");
        return _assetPrices[asset];
    }

    /**
     * @inheritdoc IPriceOracle
     */
    function getAssetsPrices(address[] calldata assets) external view override returns (uint256[] memory) {
        if (_shouldRevert) {
            revert("MockPriceOracle: Forced failure");
        }
        
        uint256[] memory prices = new uint256[](assets.length);
        
        for (uint256 i = 0; i < assets.length; i++) {
            require(_assetPrices[assets[i]] > 0, "MockPriceOracle: Price not set for asset");
            prices[i] = _assetPrices[assets[i]];
        }
        
        return prices;
    }

    /**
     * @inheritdoc IPriceOracle
     */
    function setAssetSource(address asset, address source) external override onlyOwner {
        // No-op in the mock, but emit the event for testing event listeners
        emit AssetSourceUpdated(asset, source);
    }

    /**
     * @inheritdoc IPriceOracle
     */
    function setAssetsSources(address[] calldata assets, address[] calldata sources) external override onlyOwner {
        require(assets.length == sources.length, "MockPriceOracle: Arrays length mismatch");
        
        for (uint256 i = 0; i < assets.length; i++) {
            emit AssetSourceUpdated(assets[i], sources[i]);
        }
    }

    /**
     * @inheritdoc IPriceOracle
     */
    function setEthUsdPrice(uint256 price) external override onlyOwner {
        _ethUsdPrice = price;
        emit EthUsdPriceUpdated(price);
    }

    /**
     * @inheritdoc IPriceOracle
     */
    function getEthUsdPrice() external view override returns (uint256) {
        if (_shouldRevert) {
            revert("MockPriceOracle: Forced failure");
        }
        
        return _ethUsdPrice;
    }

    /**
     * @dev Sets the price for an asset
     * @param asset The asset address
     * @param price The price in ETH
     */
    function setAssetPrice(address asset, uint256 price) external onlyOwner {
        _assetPrices[asset] = price;
    }

    /**
     * @dev Sets the prices for multiple assets
     * @param assets The asset addresses
     * @param prices The prices in ETH
     */
    function setAssetsPricesDirectly(address[] calldata assets, uint256[] calldata prices) external onlyOwner {
        require(assets.length == prices.length, "MockPriceOracle: Arrays length mismatch");
        
        for (uint256 i = 0; i < assets.length; i++) {
            _assetPrices[assets[i]] = prices[i];
        }
    }

    /**
     * @dev Sets whether operations should revert
     * @param shouldRevert Whether operations should revert
     */
    function setShouldRevert(bool shouldRevert) external onlyOwner {
        _shouldRevert = shouldRevert;
    }
}
