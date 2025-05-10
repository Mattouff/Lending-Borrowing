// SPDX-License-Identifier: MIT
pragma solidity 0.8.29;

import "forge-std/Script.sol";
import "../src/Token.sol";
import "../src/LendingPool.sol";
import "../src/Borrowing.sol";
import "../src/Collateral.sol";

contract Deploy is Script {
    // Configuration
    uint256 public constant INITIAL_TOKEN_SUPPLY = 1_000_000 * 1e18; // 1 million tokens
    uint256 public constant LENDING_INTEREST_RATE = 5 * 1e16; // 5% annual interest rate
    uint256 public constant BORROW_MIN_RATE = 3 * 1e16; // 3% minimum borrow rate
    uint256 public constant BORROW_MAX_RATE = 20 * 1e16; // 20% maximum borrow rate
    uint256 public constant BETA_FACTOR = 1 * 1e18; // Elasticity factor = 1

    // Deployed contract addresses
    Token public token;
    LendingPool public lendingPool;
    Borrowing public borrowing;
    Collateral public collateral;

    function run() external {
        // Get the private key from the environment variable
        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");

        // Start broadcasting transactions
        vm.startBroadcast(deployerPrivateKey);

        // Deploy Token contract
        token = new Token(INITIAL_TOKEN_SUPPLY);
        console.log("Token deployed at:", address(token));

        // Deploy LendingPool contract
        lendingPool = new LendingPool(address(token), LENDING_INTEREST_RATE);
        console.log("LendingPool deployed at:", address(lendingPool));

        // Set placeholder for Collateral - will be replaced later with actual address
        address collateralAddress = address(1);

        // Deploy Borrowing contract with placeholder
        borrowing = new Borrowing(address(token), collateralAddress, BORROW_MIN_RATE, BORROW_MAX_RATE, BETA_FACTOR);
        console.log("Borrowing deployed at:", address(borrowing));

        // Now deploy Collateral with the actual borrowing contract address
        collateral = new Collateral(address(token), address(borrowing));
        console.log("Collateral deployed at:", address(collateral));

        // Now we need to redeploy Borrowing with the correct Collateral address
        borrowing = new Borrowing(address(token), address(collateral), BORROW_MIN_RATE, BORROW_MAX_RATE, BETA_FACTOR);
        console.log("Borrowing redeployed at:", address(borrowing));

        // Create a small buffer of tokens in the LendingPool and Borrowing contracts
        token.transfer(address(lendingPool), 100_000 * 1e18);
        token.transfer(address(borrowing), 100_000 * 1e18);

        // Allocate tokens to Anvil test accounts for easier testing
        address[] memory testAccounts = new address[](3);
        testAccounts[0] = 0x70997970C51812dc3A010C7d01b50e0d17dc79C8; // Second Anvil account
        testAccounts[1] = 0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC; // Third Anvil account
        testAccounts[2] = 0x90F79bf6EB2c4f870365E785982E1f101E93b906; // Fourth Anvil account

        for (uint256 i = 0; i < testAccounts.length; i++) {
            token.transfer(testAccounts[i], 100_000 * 1e18);
            console.log("Transferred 100,000 tokens to", testAccounts[i]);
        }

        // Stop broadcasting transactions
        vm.stopBroadcast();

        // Output contract addresses in format for .env file
        string memory contractAddresses = string(
            abi.encodePacked(
                "LendingPool=",
                vm.toString(address(lendingPool)),
                ",Token=",
                vm.toString(address(token)),
                ",Borrowing=",
                vm.toString(address(borrowing)),
                ",Collateral=",
                vm.toString(address(collateral))
            )
        );
        console.log("\n--- Contract Addresses for .env ---");
        console.log("CONTRACT_ADDRESSES=", contractAddresses);
    }
}
