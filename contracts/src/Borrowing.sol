// SPDX-License-Identifier: MIT
pragma solidity 0.8.29;

import "./Token.sol";
import "./Collateral.sol";

/// @title Borrowing - A contract to manage borrowing operations for a decentralized lending platform.
/// @notice Users can borrow tokens and repay their loans. Borrowing limits are defined by the collateral managed in Collateral.sol.
///         A dynamic interest rate r(U) = rMin + (rMax - rMin) * U^beta est appliqué sur le montant emprunté.
contract Borrowing {
    Token public immutable token;
    Collateral public immutable collateral;

    /// @notice Mapping of user addresses to their principal borrowed token amounts (without accrued interest).
    mapping(address => uint256) public borrowedPrincipal;

    /// @notice Mapping of user addresses to the last time we updated their borrowed balance.
    mapping(address => uint256) public lastUpdateTime;

    /// @notice Total borrowed tokens across all users (somme de borrowedPrincipal).
    uint256 public totalBorrowed;

    /// @notice Paramètres pour le calcul du taux d'intérêt dynamique.
    uint256 public rMin; // Taux minimum, en 1e18 (ex: 5e16 = 5%).
    uint256 public rMax; // Taux maximum, en 1e18 (ex: 20e16 = 20%).
    uint256 public beta; // Facteur d’élasticité, en 1e18 (ex: 1e18 = exponent 1, 2e18 = exponent 2, etc.).

    /// @notice Constructor for the Borrowing contract.
    /// @param _token The address of the token used for borrowing.
    /// @param _collateral The address of the Collateral contract defining collateral conditions.
    /// @param _rMin Taux d’intérêt minimum (ex: 5e16 = 5%).
    /// @param _rMax Taux d’intérêt maximum (ex: 20e16 = 20%).
    /// @param _beta Facteur d’élasticité (ex: 1e18 = exponent 1, 2e18 = exponent 2).
    constructor(address _token, address _collateral, uint256 _rMin, uint256 _rMax, uint256 _beta) {
        token = Token(_token);
        collateral = Collateral(_collateral);
        rMin = _rMin;
        rMax = _rMax;
        beta = _beta;
    }

    /// @notice Gets the current interest rate for borrowed tokens.
    /// @return The current interest rate in 1e18.
    /// @dev The interest rate is calculated based on the total borrowed amount and the available balance in the pool.
    ///         The interest rate is determined by the formula: r(U) = rMin + (rMax - rMin) * U^beta.
    ///         Where U is the utilization rate of the pool: U = totalBorrowed / (balanceDisponible + totalBorrowed).
    ///         The utilization rate is calculated as a fraction of the total borrowed amount over the total pool capacity.
    function getCurrentRate() public view returns (uint256) {
        // Déterminer la capacité du pool : ce qui reste dans le contrat + totalBorrowed
        uint256 balanceDisponible = token.balanceOf(address(this));
        uint256 capacity = balanceDisponible + totalBorrowed;
        if (capacity == 0) {
            // Si rien n'est dans le pool, on renvoie le taux minimum
            return rMin;
        }

        uint256 U = (totalBorrowed * 1e18) / capacity;
        // Calculate U^beta properly
        uint256 UtoBeta = power(U, beta);
        uint256 diff = rMax - rMin;
        uint256 variablePart = (diff * UtoBeta) / 1e18;
        return rMin + variablePart;
    }

    /// @notice Updates the borrowed balance for a user by calculating the interest accrued.
    /// @param user The address of the user.
    /// @dev The interest is calculated based on the time elapsed since the last update.
    function updateBorrowedPrincipal(address user) internal {
        uint256 previousTime = lastUpdateTime[user];
        lastUpdateTime[user] = block.timestamp;

        // If the first time, do nothing
        if (previousTime == 0) {
            return;
        }

        uint256 timeElapsed = block.timestamp - previousTime;
        if (timeElapsed == 0) {
            return;
        }

        uint256 principal = borrowedPrincipal[user];
        if (principal == 0) {
            return;
        }

        uint256 currentRate = getCurrentRate();

        // Calculate compound interest: A = P × (1 + r/n)^(n×t)
        // For simplicity and as agreed, we'll use daily compounding (n = 365)
        // timeElapsed is in seconds, so we convert to days for the formula
        uint256 daysElapsed = timeElapsed / (1 days);
        if (daysElapsed > 0) {
            uint256 dailyRate = currentRate / 365; // Daily rate

            // Calculate compound factor: (1 + r/n)^(days)
            uint256 compoundFactor = 1e18; // Start with 1 in fixed-point
            for (uint256 i = 0; i < daysElapsed; i++) {
                compoundFactor = (compoundFactor * (1e18 + dailyRate)) / 1e18;
            }

            // Calculate new balance with compound interest
            uint256 newBalance = (principal * compoundFactor) / 1e18;
            uint256 interest = newBalance - principal;

            if (interest > 0) {
                borrowedPrincipal[user] += interest;
                totalBorrowed += interest;
            }
        }
    }

    /// @notice Allows a user to borrow tokens.
    /// @param amount The amount of tokens to borrow.
    /// @dev Collateral conditions are enforced via Collateral.canBorrow.
    function borrow(uint256 amount) external {
        require(amount > 0, "Amount must be greater than zero");
        require(collateral.canBorrow(msg.sender, amount), "Insufficient collateral");
        updateBorrowedPrincipal(msg.sender);

        borrowedPrincipal[msg.sender] += amount;
        totalBorrowed += amount;
        require(token.transfer(msg.sender, amount), "Borrow transfer failed");
    }

    /// @notice Repays a portion or all of the borrowed tokens.
    /// @param amount The amount of tokens to repay.
    /// @dev The caller must have approved the contract to transfer tokens on their behalf.
    function repay(uint256 amount) external {
        require(amount > 0, "Amount must be greater than zero");
        updateBorrowedPrincipal(msg.sender);
        require(borrowedPrincipal[msg.sender] >= amount, "Repay amount exceeds borrowed balance");
        require(token.transferFrom(msg.sender, address(this), amount), "Repay transfer failed");
        borrowedPrincipal[msg.sender] -= amount;
        totalBorrowed -= amount;
    }

    /// @notice Reduces the debt of a borrower by the specified amount.
    /// @dev Can only be called by the associated Collateral contract if we need to liquidate the user.
    /// @param borrower The address of the borrower.
    /// @param amount The amount by which to reduce the debt.
    function reduceDebt(address borrower, uint256 amount) external {
        require(msg.sender == address(collateral), "Not authorized");
        require(borrowedPrincipal[borrower] >= amount, "Insufficient debt");
        borrowedPrincipal[borrower] -= amount;
        totalBorrowed -= amount;
    }

    /// @notice Returns the current borrowed balance for a user, including accrued interest.
    /// @param user The address of the user.
    /// @return The current borrowed balance including interest.
    function getBorrowToken(address user) external view returns (uint256) {
        if (lastUpdateTime[user] == 0 || borrowedPrincipal[user] == 0) {
            return borrowedPrincipal[user];
        }

        uint256 timeElapsed = block.timestamp - lastUpdateTime[user];
        if (timeElapsed == 0) {
            return borrowedPrincipal[user];
        }

        uint256 principal = borrowedPrincipal[user];
        uint256 currentRate = getCurrentRate();
        uint256 interest = (principal * currentRate * timeElapsed) / (365 days * 1e18);

        return principal + interest;
    }

    /// @notice Gets the total borrowed tokens across all users.
    /// @return The total amount of tokens borrowed.
    function getAllBorrowToken() external view returns (uint256) {
        return totalBorrowed;
    }

    /// @notice Calculates the power of a base raised to an exponent.
    /// @param base The base number.
    /// @param exponent The exponent to raise the base to.
    /// @return The result of base^exponent in fixed-point representation.
    function power(uint256 base, uint256 exponent) internal pure returns (uint256) {
        uint256 result = 1e18; // Start with 1 in fixed-point representation

        for (uint256 i = 0; i < exponent / 1e18; i++) {
            result = (result * base) / 1e18;
        }

        return result;
    }
}
