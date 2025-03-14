// SPDX-License-Identifier: MIT
pragma solidity 0.8.29;

import "@openzeppelin/contracts/access/Ownable.sol";
import "@chainlink/v0.8/shared/interfaces/AggregatorV3Interface.sol";
import "../interfaces/IPriceOracle.sol";

/**
 * @title PriceOracle
 * @author DeFi Lending Platform
 * @notice Provides prices for assets using Chainlink price feeds
 */
contract PriceOracle is IPriceOracle, Ownable {
    // Mapping from asset address to price feed address
    mapping(address => address) private assetToPriceFeed;

    // ETH/USD price feed address
    address private ethUsdPriceFeed;

    // ETH/USD price
    uint256 private ethUsdPrice;

    // Time window for price validity
    uint256 public constant PRICE_EXPIRATION_TIME = 1 hours;

    /**
     * @dev Constructor
     * @param _owner The owner of the oracle
     * @param _ethUsdPriceFeed The ETH/USD price feed address
     */
    constructor(address _owner, address _ethUsdPriceFeed) Ownable(msg.sender) {
        transferOwnership(_owner);
        ethUsdPriceFeed = _ethUsdPriceFeed;
        updateEthUsdPrice();
    }

    /**
     * @inheritdoc IPriceOracle
     */
    function getAssetPrice(address asset) external view override returns (uint256) {
        address priceFeed = assetToPriceFeed[asset];
        require(priceFeed != address(0), "PriceOracle: Price feed not found");

        return _getAssetPrice(priceFeed);
    }

    /**
     * @inheritdoc IPriceOracle
     */
    function getAssetsPrices(address[] calldata assets) external view override returns (uint256[] memory) {
        uint256[] memory prices = new uint256[](assets.length);

        for (uint256 i = 0; i < assets.length; i++) {
            address priceFeed = assetToPriceFeed[assets[i]];
            require(priceFeed != address(0), "PriceOracle: Price feed not found");

            prices[i] = _getAssetPrice(priceFeed);
        }

        return prices;
    }

    /**
     * @inheritdoc IPriceOracle
     */
    function setAssetSource(address asset, address source) external override onlyOwner {
        assetToPriceFeed[asset] = source;
        emit AssetSourceUpdated(asset, source);
    }

    /**
     * @inheritdoc IPriceOracle
     */
    function setAssetsSources(address[] calldata assets, address[] calldata sources) external override onlyOwner {
        require(assets.length == sources.length, "PriceOracle: Arrays length mismatch");

        for (uint256 i = 0; i < assets.length; i++) {
            assetToPriceFeed[assets[i]] = sources[i];
            emit AssetSourceUpdated(assets[i], sources[i]);
        }
    }

    /**
     * @inheritdoc IPriceOracle
     */
    function setEthUsdPrice(uint256 price) external override onlyOwner {
        ethUsdPrice = price;
        emit EthUsdPriceUpdated(price);
    }

    /**
     * @inheritdoc IPriceOracle
     */
    function getEthUsdPrice() external view override returns (uint256) {
        return ethUsdPrice;
    }

    /**
     * @dev Updates the ETH/USD price from the price feed
     */
    function updateEthUsdPrice() public {
        if (ethUsdPriceFeed == address(0)) {
            return;
        }

        (, int256 price,, uint256 updatedAt,) = AggregatorV3Interface(ethUsdPriceFeed).latestRoundData();

        // Check for stale price
        require(block.timestamp - updatedAt <= PRICE_EXPIRATION_TIME, "PriceOracle: Stale price");

        // Price must be positive
        require(price > 0, "PriceOracle: Negative price");

        ethUsdPrice = uint256(price);
        emit EthUsdPriceUpdated(ethUsdPrice);
    }

    /**
     * @dev Gets the asset price from the price feed
     * @param priceFeed The price feed address
     * @return The asset price in ETH
     */
    function _getAssetPrice(address priceFeed) internal view returns (uint256) {
        (, int256 price,, uint256 updatedAt,) = AggregatorV3Interface(priceFeed).latestRoundData();

        // Check for stale price
        require(block.timestamp - updatedAt <= PRICE_EXPIRATION_TIME, "PriceOracle: Stale price");

        // Price must be positive
        require(price > 0, "PriceOracle: Negative price");

        // Convert to 18 decimals
        uint8 decimals = AggregatorV3Interface(priceFeed).decimals();
        uint256 priceUint = uint256(price);

        if (decimals < 18) {
            return priceUint * 10 ** (18 - decimals);
        } else {
            return priceUint / 10 ** (decimals - 18);
        }
    }
}
