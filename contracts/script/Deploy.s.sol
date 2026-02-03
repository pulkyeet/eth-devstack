// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Script.sol";
import "../src/Token.sol";
import "../src/Staking.sol";
import "../src/DEX.sol";

contract DeployScript is Script {
    function run() external {
        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");
        
        vm.startBroadcast(deployerPrivateKey);

        console.log("Deploying contracts...");
        console.log("Deployer:", vm.addr(deployerPrivateKey));

        // Deploy Token
        Token2178 token = new Token2178();
        console.log("Token:", address(token));

        // Deploy Staking
        Staking staking = new Staking(
            address(token),
            address(token)
        );
        console.log("Staking:", address(staking));

        // Deploy second token for DEX
        Token2178 token1 = new Token2178();
        console.log("Token1:", address(token1));

        // Deploy DEX
        DEX2178 dex = new DEX2178(
            address(token),
            address(token1)
        );
        console.log("DEX:", address(dex));

        // Fund staking with rewards (10M tokens)
        token.transfer(address(staking), 10_000_000 * 10**18);
        console.log("\nStaking funded with 10M tokens");

        vm.stopBroadcast();

        console.log("\n=== DEPLOYMENT COMPLETE ===");
        console.log("Save these addresses:");
        console.log("Token:", address(token));
        console.log("Token1:", address(token1));
        console.log("Staking:", address(staking));
        console.log("DEX:", address(dex));
    }
}