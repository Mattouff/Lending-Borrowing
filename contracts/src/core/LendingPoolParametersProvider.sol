// SPDX-License-Identifier: MIT
pragma solidity 0.8.29;

import "@openzeppelin/contracts/access/Ownable.sol";

/**
 * @title LendingPoolParametersProvider
 * @author DeFi Lending Platform
 * @notice Provides the parameters for the lending pool
 * @dev Contains the default parameters for the lending pool
 */
contract LendingPoolParametersProvider is Ownable {
    // Default parameters
    uint256 private _minHealthFactor = 1.05 * 1e18; // 105%
    uint256 private _maxNumberOfReserves = 128;
    uint256 private _defaultInterestRateSlope1 = 0.1 * 1e18; // 10%
    uint256 private _defaultInterestRateSlope2 = 0.4 * 1e18; // 40%
    uint256 private _defaultOptimalUtilizationRate = 0.8 * 1e18; // 80%
    uint256 private _defaultBaseBorrowRate = 0.01 * 1e18; // 1%
    uint256 private _defaultMaxBorrowRate = 1.0 * 1e18; // 100%
    uint256 private _defaultElasticityFactor = 2 * 1e18; // 2.0 (quadratic)
    uint256 private _defaultReserveFactor = 0.1 * 1e18; // 10%
    uint256 private _defaultLiquidationBonus = 1.1 * 1e18; // 110%
    uint256 private _defaultCloseFactor = 0.5 * 1e18; // 50%

    // Events
    event MinHealthFactorChanged(uint256 oldValue, uint256 newValue);
    event MaxNumberOfReservesChanged(uint256 oldValue, uint256 newValue);
    event InterestRateSlopeChanged(uint256 oldSlope1, uint256 newSlope1, uint256 oldSlope2, uint256 newSlope2);
    event OptimalUtilizationRateChanged(uint256 oldValue, uint256 newValue);
    event BaseBorrowRateChanged(uint256 oldValue, uint256 newValue);
    event MaxBorrowRateChanged(uint256 oldValue, uint256 newValue);
    event ElasticityFactorChanged(uint256 oldValue, uint256 newValue);
    event ReserveFactorChanged(uint256 oldValue, uint256 newValue);
    event LiquidationBonusChanged(uint256 oldValue, uint256 newValue);
    event CloseFactorChanged(uint256 oldValue, uint256 newValue);

    constructor() Ownable(msg.sender) { }

    /**
     * @dev Gets the minimum health factor
     * @return The minimum health factor
     */
    function getMinHealthFactor() external view returns (uint256) {
        return _minHealthFactor;
    }

    /**
     * @dev Sets the minimum health factor
     * @param newMinHealthFactor The new minimum health factor
     */
    function setMinHealthFactor(uint256 newMinHealthFactor) external onlyOwner {
        require(newMinHealthFactor > 1e18, "LendingPoolParametersProvider: Health factor must be greater than 1");

        uint256 oldMinHealthFactor = _minHealthFactor;
        _minHealthFactor = newMinHealthFactor;

        emit MinHealthFactorChanged(oldMinHealthFactor, newMinHealthFactor);
    }

    /**
     * @dev Gets the maximum number of reserves
     * @return The maximum number of reserves
     */
    function getMaxNumberOfReserves() external view returns (uint256) {
        return _maxNumberOfReserves;
    }

    /**
     * @dev Sets the maximum number of reserves
     * @param newMaxNumberOfReserves The new maximum number of reserves
     */
    function setMaxNumberOfReserves(uint256 newMaxNumberOfReserves) external onlyOwner {
        require(
            newMaxNumberOfReserves > 0, "LendingPoolParametersProvider: Max number of reserves must be greater than 0"
        );

        uint256 oldMaxNumberOfReserves = _maxNumberOfReserves;
        _maxNumberOfReserves = newMaxNumberOfReserves;

        emit MaxNumberOfReservesChanged(oldMaxNumberOfReserves, newMaxNumberOfReserves);
    }

    /**
     * @dev Gets the default interest rate parameters
     * @return slope1 The default interest rate slope1
     * @return slope2 The default interest rate slope2
     * @return optimalUtilizationRate The default optimal utilization rate
     * @return baseBorrowRate The default base borrow rate
     * @return maxBorrowRate The default max borrow rate
     * @return elasticityFactor The default elasticity factor
     */
    function getDefaultInterestRateParameters()
        external
        view
        returns (
            uint256 slope1,
            uint256 slope2,
            uint256 optimalUtilizationRate,
            uint256 baseBorrowRate,
            uint256 maxBorrowRate,
            uint256 elasticityFactor
        )
    {
        return (
            _defaultInterestRateSlope1,
            _defaultInterestRateSlope2,
            _defaultOptimalUtilizationRate,
            _defaultBaseBorrowRate,
            _defaultMaxBorrowRate,
            _defaultElasticityFactor
        );
    }

    /**
     * @dev Sets the default interest rate slopes
     * @param newSlope1 The new slope1
     * @param newSlope2 The new slope2
     */
    function setDefaultInterestRateSlopes(uint256 newSlope1, uint256 newSlope2) external onlyOwner {
        uint256 oldSlope1 = _defaultInterestRateSlope1;
        uint256 oldSlope2 = _defaultInterestRateSlope2;

        _defaultInterestRateSlope1 = newSlope1;
        _defaultInterestRateSlope2 = newSlope2;

        emit InterestRateSlopeChanged(oldSlope1, newSlope1, oldSlope2, newSlope2);
    }

    /**
     * @dev Sets the default optimal utilization rate
     * @param newOptimalUtilizationRate The new optimal utilization rate
     */
    function setDefaultOptimalUtilizationRate(uint256 newOptimalUtilizationRate) external onlyOwner {
        require(
            newOptimalUtilizationRate <= 1e18,
            "LendingPoolParametersProvider: Optimal utilization rate must be less than or equal to 100%"
        );

        uint256 oldOptimalUtilizationRate = _defaultOptimalUtilizationRate;
        _defaultOptimalUtilizationRate = newOptimalUtilizationRate;

        emit OptimalUtilizationRateChanged(oldOptimalUtilizationRate, newOptimalUtilizationRate);
    }

    /**
     * @dev Sets the default base borrow rate
     * @param newBaseBorrowRate The new base borrow rate
     */
    function setDefaultBaseBorrowRate(uint256 newBaseBorrowRate) external onlyOwner {
        uint256 oldBaseBorrowRate = _defaultBaseBorrowRate;
        _defaultBaseBorrowRate = newBaseBorrowRate;

        emit BaseBorrowRateChanged(oldBaseBorrowRate, newBaseBorrowRate);
    }

    /**
     * @dev Sets the default max borrow rate
     * @param newMaxBorrowRate The new max borrow rate
     */
    function setDefaultMaxBorrowRate(uint256 newMaxBorrowRate) external onlyOwner {
        uint256 oldMaxBorrowRate = _defaultMaxBorrowRate;
        _defaultMaxBorrowRate = newMaxBorrowRate;

        emit MaxBorrowRateChanged(oldMaxBorrowRate, newMaxBorrowRate);
    }

    /**
     * @dev Sets the default elasticity factor
     * @param newElasticityFactor The new elasticity factor
     */
    function setDefaultElasticityFactor(uint256 newElasticityFactor) external onlyOwner {
        uint256 oldElasticityFactor = _defaultElasticityFactor;
        _defaultElasticityFactor = newElasticityFactor;

        emit ElasticityFactorChanged(oldElasticityFactor, newElasticityFactor);
    }

    /**
     * @dev Gets the default reserve factor
     * @return The default reserve factor
     */
    function getDefaultReserveFactor() external view returns (uint256) {
        return _defaultReserveFactor;
    }

    /**
     * @dev Sets the default reserve factor
     * @param newReserveFactor The new reserve factor
     */
    function setDefaultReserveFactor(uint256 newReserveFactor) external onlyOwner {
        require(
            newReserveFactor <= 1e18, "LendingPoolParametersProvider: Reserve factor must be less than or equal to 100%"
        );

        uint256 oldReserveFactor = _defaultReserveFactor;
        _defaultReserveFactor = newReserveFactor;

        emit ReserveFactorChanged(oldReserveFactor, newReserveFactor);
    }

    /**
     * @dev Gets the default liquidation bonus
     * @return The default liquidation bonus
     */
    function getDefaultLiquidationBonus() external view returns (uint256) {
        return _defaultLiquidationBonus;
    }

    /**
     * @dev Sets the default liquidation bonus
     * @param newLiquidationBonus The new liquidation bonus
     */
    function setDefaultLiquidationBonus(uint256 newLiquidationBonus) external onlyOwner {
        require(
            newLiquidationBonus >= 1e18,
            "LendingPoolParametersProvider: Liquidation bonus must be greater than or equal to 100%"
        );

        uint256 oldLiquidationBonus = _defaultLiquidationBonus;
        _defaultLiquidationBonus = newLiquidationBonus;

        emit LiquidationBonusChanged(oldLiquidationBonus, newLiquidationBonus);
    }

    /**
     * @dev Gets the default close factor
     * @return The default close factor
     */
    function getDefaultCloseFactor() external view returns (uint256) {
        return _defaultCloseFactor;
    }

    /**
     * @dev Sets the default close factor
     * @param newCloseFactor The new close factor
     */
    function setDefaultCloseFactor(uint256 newCloseFactor) external onlyOwner {
        require(
            newCloseFactor <= 1e18, "LendingPoolParametersProvider: Close factor must be less than or equal to 100%"
        );

        uint256 oldCloseFactor = _defaultCloseFactor;
        _defaultCloseFactor = newCloseFactor;

        emit CloseFactorChanged(oldCloseFactor, newCloseFactor);
    }
}
