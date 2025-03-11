// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "forge-std/Script.sol";
import "../src/Token.sol";

contract DeployToken is Script {
    function run() external {
        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");
        
        vm.startBroadcast(deployerPrivateKey);
        
        uint256 initialSupply = 1000 * 10 ** 18;
        
        Token token = new Token(initialSupply);
        
        vm.stopBroadcast();
        
        console.log("Token deployed :", address(token));
    }
}
