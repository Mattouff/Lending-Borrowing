// SPDX-License-Identifier: MIT
pragma solidity 0.8.29;

import "@openzeppelin/contracts/access/Ownable.sol";

/**
 * @title MockChainlinkAggregator
 * @author DeFi Lending Platform
 * @notice Mock Chainlink price feed aggregator for testing purposes
 * @dev Implements the AggregatorV3Interface with configurable responses
 */
contract MockChainlinkAggregator is Ownable {
    int256 private _price;
    uint8 private _decimals;
    string private _description;
    uint256 private _version;
    uint80 private _roundId;
    uint256 private _updatedAt;
    
    /**
     * @dev Constructor to create a new MockChainlinkAggregator
     * @param price The initial price to return
     * @param decimalsValue The decimal places for the price feed
     * @param descriptionText A description of the price feed
     */
    constructor(
        int256 price,
        uint8 decimalsValue,
        string memory descriptionText
    ) Ownable(msg.sender) {
        _price = price;
        _decimals = decimalsValue;
        _description = descriptionText;
        _version = 1;
        _roundId = 1;
        _updatedAt = block.timestamp;
    }
    
    /**
     * @dev Returns the decimal places of the price feed
     * @return The number of decimal places
     */
    function decimals() external view returns (uint8) {
        return _decimals;
    }
    
    /**
     * @dev Returns the description of the price feed
     * @return A description string
     */
    function description() external view returns (string memory) {
        return _description;
    }
    
    /**
     * @dev Returns the version of the price feed
     * @return The version number
     */
    function version() external view returns (uint256) {
        return _version;
    }
    
    /**
     * @dev Returns the most recent round data
     * @return roundId The round ID
     * @return answer The price
     * @return startedAt When the round started
     * @return updatedAt When the round was updated
     * @return answeredInRound The round in which the answer was computed
     */
    function latestRoundData() external view returns (
        uint80 roundId,
        int256 answer,
        uint256 startedAt,
        uint256 updatedAt,
        uint80 answeredInRound
    ) {
        return (
            _roundId,
            _price,
            _updatedAt,
            _updatedAt,
            _roundId
        );
    }
    
    /**
     * @dev Returns the latest price
     * @return The latest price
     */
    function latestAnswer() external view returns (int256) {
        return _price;
    }
    
    /**
     * @dev Returns the latest timestamp
     * @return The latest timestamp
     */
    function latestTimestamp() external view returns (uint256) {
        return _updatedAt;
    }
    
    /**
     * @dev Returns the latest round
     * @return The latest round ID
     */
    function latestRound() external view returns (uint256) {
        return _roundId;
    }
    
    /**
     * @dev Sets the price to be returned
     * @param price The new price
     */
    function setPrice(int256 price) external onlyOwner {
        _price = price;
        _roundId++;
        _updatedAt = block.timestamp;
    }
    
    /**
     * @dev Sets the timestamp to be returned
     * @param timestamp The new timestamp
     */
    function setUpdatedAt(uint256 timestamp) external onlyOwner {
        _updatedAt = timestamp;
    }
    
    /**
     * @dev Sets the decimals for the price feed
     * @param decimalsValue The new decimal places
     */
    function setDecimals(uint8 decimalsValue) external onlyOwner {
        _decimals = decimalsValue;
    }
    
    /**
     * @dev Sets the description for the price feed
     * @param descriptionText The new description
     */
    function setDescription(string memory descriptionText) external onlyOwner {
        _description = descriptionText;
    }
    
    /**
     * @dev Increments the round ID
     */
    function nextRound() external onlyOwner {
        _roundId++;
        _updatedAt = block.timestamp;
    }

    /**
     * @dev Creates a timestamp in the past to test stale prices
     * @param secondsAgo How many seconds ago to set the timestamp
     */
    function setTimestampAgo(uint256 secondsAgo) external onlyOwner {
        if (secondsAgo > block.timestamp) {
            _updatedAt = 0;
        } else {
            _updatedAt = block.timestamp - secondsAgo;
        }
    }
}
