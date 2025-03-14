// SPDX-License-Identifier: MIT
pragma solidity 0.8.29;

import "forge-std/Script.sol";
import "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";
import "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";
import "../src/core/LendingPool.sol";

/**
 * @title UpgradeLendingPool
 * @author DeFi Lending Platform
 * @notice Script to upgrade the LendingPool implementation
 */
contract UpgradeLendingPool is Script {
    // Replace these addresses with your deployed contract addresses
    address constant PROXY_ADMIN_ADDRESS = address(0); // Replace with your ProxyAdmin address
    address constant LENDING_POOL_PROXY_ADDRESS = address(0); // Replace with your LendingPool proxy address

    function run() external {
        vm.startBroadcast();

        // Deploy new implementation
        LendingPool newImplementation = new LendingPool();

        // Get the ProxyAdmin instance
        ProxyAdmin proxyAdmin = ProxyAdmin(PROXY_ADMIN_ADDRESS);

        // Upgrade the proxy to point to the new implementation
        proxyAdmin.upgradeAndCall(
            ITransparentUpgradeableProxy(payable(LENDING_POOL_PROXY_ADDRESS)),
            address(newImplementation),
            new bytes(0) // Empty call data since we don't need to call any function
        );

        vm.stopBroadcast();

        console.log("LendingPool upgraded to new implementation: ", address(newImplementation));
    }
}
