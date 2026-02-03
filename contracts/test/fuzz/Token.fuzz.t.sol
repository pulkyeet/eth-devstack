// SPDX-License-Identifier: MIT

pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../../src/Token.sol";

contract TokenFuzzTest is Test {
    Token2178 token;
    address owner = address(this);

    function setUp() public {
        token = new Token2178();
    }

    function testFuzzMint(address to, uint256 amount) public {
        // dont mint to 0x0
        vm.assume(to != address(0));
        vm.assume(to!=owner);
        vm.assume(amount > 0);
        vm.assume(amount <= token.MAX_SUPPLY() - token.totalSupply());

        uint256 supplyBefore = token.totalSupply();

        token.mint(to, amount);

        assertEq(token.balanceOf(to), amount);
        assertEq(token.totalSupply(), supplyBefore + amount);
        assertTrue(token.totalSupply() <= token.MAX_SUPPLY());
    }

    function testFuzzTransfer(address to, uint256 amount) public {
        vm.assume(to!=address(0));
        //cant send to self
        vm.assume(to!=owner);
        vm.assume(amount>0);
        vm.assume(amount <= token.balanceOf(owner));

        uint256 ownerBalanceBefore = token.balanceOf(owner);
        uint256 toBalanceBefore = token.balanceOf(to);
        uint256 totalSupplyBefore = token.totalSupply();

        token.transfer(to, amount); 

        assertEq(token.balanceOf(owner), ownerBalanceBefore - amount);
        assertEq(token.balanceOf(to), toBalanceBefore + amount);
        assertEq(token.totalSupply(), totalSupplyBefore);
    }

    function testFuzzBurn(uint256 amount) public {
        vm.assume(amount > 0);
        vm.assume(amount <= token.balanceOf(owner));
        
        uint256 balanceBefore = token.balanceOf(owner);
        uint256 supplyBefore = token.totalSupply();
        
        token.burn(amount);
        
        assertEq(token.balanceOf(owner), balanceBefore - amount);
        assertEq(token.totalSupply(), supplyBefore - amount);
        assertTrue(token.totalSupply() <= token.MAX_SUPPLY());
    }

    function testFuzzApproveTranferFrom(address spender, address to, uint256 approvalAmount, uint256 transferAmount) public {
        vm.assume(spender!=address(0));
        vm.assume(to!=address(0));
        vm.assume(to!=owner);
        vm.assume(spender!=owner);
        vm.assume(approvalAmount>0 && approvalAmount<type(uint256).max);
        vm.assume(transferAmount>0 && approvalAmount<type(uint256).max);
        vm.assume(transferAmount <=approvalAmount);
        vm.assume(transferAmount<=token.balanceOf(owner));

        token.approve(spender, approvalAmount);
        assertEq(token.allowance(owner, spender), approvalAmount);

        vm.prank(spender);
        token.transferFrom(owner, to, transferAmount);

        assertEq(token.allowance(owner, spender), approvalAmount - transferAmount);
        assertEq(token.balanceOf(to), transferAmount);
    }

    function testFuzzLockTokens(address user, uint256 lockDuration) public {
        vm.assume(user!=address(0));
        vm.assume(lockDuration>0);
        vm.assume(lockDuration<365 days * 100);

        token.transfer(user, 1000);

        uint256 unlockTime = block.timestamp + lockDuration;
        token.lockTokens(user, unlockTime);

        vm.prank(user);
        vm.expectRevert("Tokens are locked.");
        token.transfer(owner, 100);

        vm.warp(unlockTime+1);

        vm.prank(user);
        token.transfer(owner, 100);
        assertEq(token.balanceOf(owner), token.balanceOf(owner));
    }
}