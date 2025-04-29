// SPDX-License-Identifier: MIT
pragma solidity 0.8.29;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "./Token.sol";

/// @title LendingPool - A contract to manage deposits for a decentralized lending platform.
/// @notice Users deposit ERC20 tokens (defined in Token.sol) and receive deposit tokens (dToken) representing their share in the pool.
/// They earn interest on their deposits based on the deposit amount and the elapsed time.
contract LendingPool is ERC20 {
    Token public immutable underlying;

    /// @notice Mapping of user addresses to their deposited amounts (including accrued interest).
    mapping(address => uint256) public lendingBalance;

    /// @notice Total deposited tokens (including accrued interest) across all users.
    uint256 public totalLending;

    /// @notice Annual interest rate, expressed with 18 decimals (e.g. 5% = 5e16).
    uint256 public annualInterestRate;

    /// @notice Mapping of user addresses to the last timestamp when interest was accrued.
    mapping(address => uint256) public lastUpdate;

    /// @notice Constructor for the LendingPool contract.
    /// @param _underlying The address of the underlying token (Token.sol).
    /// @param _annualInterestRate The annual interest rate (with 18 decimals).
    constructor(address _underlying, uint256 _annualInterestRate) ERC20("Deposit Token", "dTOKEN") {
        underlying = Token(_underlying);
        annualInterestRate = _annualInterestRate;
    }

    /// @notice Internal function to update accrued interest for a user.
    /// @param user The address of the user.
    function updateInterest(address user) internal {
        uint256 lastTime = lastUpdate[user];
        lastUpdate[user] = block.timestamp;

        // If it's the first update or no time has passed, exit early
        if (lastTime == 0) {
            return;
        }

        // Check if time has elapsed and user has a balance
        uint256 timeElapsed = block.timestamp - lastTime;
        uint256 balance = lendingBalance[user];

        if (timeElapsed == 0 || balance == 0) {
            return;
        }

        // Calculate compound interest: A = P × (1 + r/n)^(n×t)
        // For simplicity and as agreed, we'll use daily compounding (n = 365)
        // timeElapsed is in seconds, so we convert to days for the formula
        uint256 daysElapsed = timeElapsed / (1 days);

        if (daysElapsed > 0) {
            uint256 dailyRate = annualInterestRate / 365; // Daily rate

            // Calculate compound factor: (1 + r/n)^(days)
            uint256 compoundFactor = 1e18; // Start with 1 in fixed-point
            for (uint256 i = 0; i < daysElapsed; i++) {
                compoundFactor = (compoundFactor * (1e18 + dailyRate)) / 1e18;
            }

            // Calculate new balance with compound interest
            uint256 newBalance = (balance * compoundFactor) / 1e18;
            uint256 interest = newBalance - balance;

            if (interest > 0) {
                lendingBalance[user] += interest;
                totalLending += interest;
                _mint(user, interest);
            }
        }
    }

    /// @notice External function to update accrued interest for a user.
    /// @param user The address of the user.
    function updateUserInterest(address user) external {
        updateInterest(user);
    }

    /// @notice Deposit tokens into the lending pool.
    /// @param amount The amount of tokens to deposit.
    /// @dev The user must first approve the token transfer to this contract.
    function deposit(uint256 amount) external {
        require(amount > 0, "Amount must be greater than zero");
        updateInterest(msg.sender);
        lendingBalance[msg.sender] += amount;
        totalLending += amount;

        bool success = underlying.transferFrom(msg.sender, address(this), amount);
        require(success, "Token transfer failed");

        _mint(msg.sender, amount);
    }

    /// @notice Withdraw tokens from the pool by burning the corresponding deposit tokens.
    /// @param amount The amount of tokens to withdraw.
    function withdraw(uint256 amount) external {
        require(amount > 0, "Amount must be greater than zero");
        updateInterest(msg.sender);
        require(balanceOf(msg.sender) >= amount, "Insufficient deposit balance");
        lendingBalance[msg.sender] -= amount;
        totalLending -= amount;

        _burn(msg.sender, amount);

        bool success = underlying.transfer(msg.sender, amount);
        require(success, "Token transfer failed");
    }

    /// @notice Returns the lending token balance (deposit amount including accrued interest) for a specific user.
    /// @param user The address of the user.
    /// @return The amount of tokens deposited by the user.
    function getLendingToken(address user) external view returns (uint256) {
        return lendingBalance[user];
    }

    /// @notice Returns the total lending tokens (including accrued interest) across all users.
    /// @return The total amount deposited in the pool.
    function getAllLendingToken() external view returns (uint256) {
        return totalLending;
    }
}
