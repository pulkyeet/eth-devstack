// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../../src/Token.sol";

contract TokenTest is Test {
    Token2178 token;

    address owner = address(this);
    address user1 = address(0x1);
    address user2 = address(0x2);

    function setUp() public {
        token = new Token2178();
    }

    function testInitialSupply() public {
        assertEq(token.totalSupply(), 100_000 * 10 ** 18);
        assertEq(token.balanceOf(owner), 100_000 * 10 ** 18);
    }

    function testMetadata() public {
        assertEq(token.name(), "TK2178");
        assertEq(token.symbol(), "T21");
        assertEq(token.decimals(), 18);
    }

    function testMint() public {
        token.mint(user1, 1000 * 10 ** 18);

        assertEq(token.balanceOf(user1), 1000 * 10 ** 18);
        assertEq(token.totalSupply(), 101_000 * 10 ** 18);
    }

    function testMintRevertsIfExceedsMax() public {
        uint256 tooMuch = 1_000_000 * 10 ** 18;

        vm.expectRevert("Exceeds max supply");
        token.mint(user1, tooMuch);
    }

    function testMintRevertsIfNotOwner() public {
        // should fail
        vm.prank(user1);
        vm.expectRevert();
        token.mint(user1, 1000);
    }

    function testPause() public {
        //transfer tokens
        token.transfer(user1, 1000);
        assertEq(token.balanceOf(user1), 1000);

        // Token transfer paused
        token.pause();

        // should revert to og state
        vm.expectRevert();
        token.transfer(user2, 100);
    }

    function testUnpause() public {
        token.pause();
        token.unpause();

        token.transfer(user1, 1000);
        assertEq(token.balanceOf(user1), 1000);
    }

    function testLockTokens() public {
        token.transfer(user1, 1000 * 10 ** 18);

        uint256 unlockTime = block.timestamp + 1 days;

        token.lockTokens(user1, unlockTime);

        // user1 tries to send to user2 and fails
        vm.prank(user1);
        vm.expectRevert("Tokens are locked.");
        token.transfer(user2, 100);

        // FF time
        vm.warp(unlockTime + 1);
        vm.prank(user1);
        token.transfer(user2, 100);
        assertEq(token.balanceOf(user2), 100);
    }

    function testLockRevertsIfPastTimestamp() public {
        vm.expectRevert("Must be future timestamp");
        token.lockTokens(user1, block.timestamp - 1);
    }

    // test to check if owner burns tokens
    function testBurn() public {
        token.burn(1000 * 10 ** 18);

        assertEq(token.balanceOf(owner), 99_000 * 10 ** 18);
        assertEq(token.totalSupply(), 99_000 * 10 ** 18);
    }

    function testBurnFrom() public {
        // owner approves user1 to burn 1000 tokens
        token.approve(user1, 1000 * 10 ** 18);

        //user1 tries to burn owner's tokens
        vm.prank(user1);
        token.burnFrom(owner, 500 * 10 ** 18);

        assertEq(token.totalSupply(), 99_500 * 10 ** 18);
        assertEq(token.allowance(owner, user1), 500 * 10 ** 18);
    }

    function testTransfer() public {
        token.transfer(user1, 1000);

        assertEq(token.balanceOf(owner), 100_000 * 10 ** 18 - 1000);
        assertEq(token.balanceOf(user1), 1000);
    }

    function testTransferFrom() public {
        token.approve(user1, 1000);

        vm.prank(user1);
        token.transferFrom(owner, user2, 500);

        assertEq(token.balanceOf(user2), 500);
        assertEq(token.allowance(owner, user1), 500);
    }

    function testPauseRevertsIfNotOwner() public {
        vm.prank(user1);
        vm.expectRevert();
        token.pause();
    }

    function testUnpauseRevertsIfNotOwner() public {
        token.pause();

        vm.prank(user1);
        vm.expectRevert();
        token.unpause();
    }

    function testLockTokensRevertsIfNotOwner() public {
        vm.prank(user1);
        vm.expectRevert();
        token.lockTokens(user2, block.timestamp + 1 days);
    }

    function testBurnFromRevertsWithoutApproval() public {
        vm.prank(user1);
        vm.expectRevert();
        token.burnFrom(owner, 100);
    }

    function testTransferToZeroAddressReverts() public {
        vm.expectRevert();
        token.transfer(address(0), 100);
    }

    function testApproveZero() public {
        token.approve(user1, 1000);
        assertEq(token.allowance(owner, user1), 1000);

        // Approve 0 to revoke
        token.approve(user1, 0);
        assertEq(token.allowance(owner, user1), 0);
    }

    function testApproveMax() public {
        token.approve(user1, type(uint256).max);
        assertEq(token.allowance(owner, user1), type(uint256).max);
    }

    function testTransferFromExceedingAllowance() public {
        token.approve(user1, 100);

        vm.prank(user1);
        vm.expectRevert();
        token.transferFrom(owner, user2, 101);
    }

    function testTransferFromZeroAddress() public {
        // Can't prank address(0) but can test it reverts
        vm.expectRevert();
        token.transferFrom(address(0), user1, 100);
    }

    function testLockAtExactTimestamp() public {
        token.transfer(user1, 1000);

        uint256 unlockTime = block.timestamp + 1 days;
        token.lockTokens(user1, unlockTime);

        // At exact unlock time, should still be locked
        vm.warp(unlockTime);
        vm.prank(user1);
        token.transfer(user2, 100); // Should work at exact timestamp
        assertEq(token.balanceOf(user2), 100);
    }

    function testBurnExceedingBalance() public {
        uint256 tooMuch = token.balanceOf(owner) + 1;

        vm.expectRevert(
            abi.encodeWithSignature(
                "ERC20InsufficientBalance(address,uint256,uint256)",
                owner,
                token.balanceOf(owner),
                tooMuch
            )
        );
        token.burn(tooMuch);
    }

    function testTransferExceedingBalance() public {
        uint256 tooMuch = token.balanceOf(owner) + 1;

        vm.expectRevert(
            abi.encodeWithSignature(
                "ERC20InsufficientBalance(address,uint256,uint256)",
                owner,
                token.balanceOf(owner),
                tooMuch
            )
        );
        token.transfer(user1, tooMuch);
    }

    function testMintToZeroAddress() public {
        vm.expectRevert();
        token.mint(address(0), 1000);
    }

    function testPauseWhenAlreadyPaused() public {
        token.pause();

        // Pausing again should revert
        vm.expectRevert();
        token.pause();
    }

    function testUnpauseWhenNotPaused() public {
        // Unpausing when not paused should revert
        vm.expectRevert();
        token.unpause();
    }
}
