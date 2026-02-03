// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../../src/Staking.sol";
import "../../src/Token.sol";

contract StakingFuzzTest is Test {
    Token2178 public token;
    Token2178 public rewardToken;
    Staking public staking;

    address public owner = address(1);
    address public alice = address(2);
    address public bob = address(3);

    function setUp() public {
        vm.startPrank(owner);
        token = new Token2178();
        rewardToken = new Token2178();
        staking = new Staking(address(token), address(rewardToken));

        deal(address(rewardToken), owner, 10_000_000 * 10 ** 18);

        rewardToken.transfer(address(staking), 1_000_000 * 10 ** 18);

        token.transfer(alice, 10_000 * 10 ** 18);
        token.transfer(bob, 10_000 * 10 ** 18);

        vm.stopPrank();
    }

    function testFuzz_Stake(uint256 amount) public {
        amount = bound(amount, 1, 10_000 * 10 ** 18);

        vm.startPrank(alice);
        token.approve(address(staking), amount);
        staking.stake(amount);
        vm.stopPrank();

        assertEq(staking.stakedBalance(alice), amount);
        assertEq(staking.totalStaked(), amount);
    }

    function testFuzz_Withdraw(
        uint256 stakeAmount,
        uint256 withdrawAmount
    ) public {
        // limits var between amounts
        stakeAmount = bound(stakeAmount, 1000, 10_000 * 10 ** 18);
        withdrawAmount = bound(withdrawAmount, 1, stakeAmount);

        vm.startPrank(alice);
        token.approve(address(staking), stakeAmount);
        staking.stake(stakeAmount);

        // no penalty
        vm.warp(block.timestamp + 1 days + 1);

        uint256 balanceBefore = token.balanceOf(alice);
        staking.withdraw(withdrawAmount);
        uint256 balanceAfter = token.balanceOf(alice);
        vm.stopPrank();

        assertEq(balanceAfter - balanceBefore, withdrawAmount);
        assertEq(staking.stakedBalance(alice), stakeAmount - withdrawAmount);
    }

    function testFuzz_EarlyWithdrawalPenalty(
        uint256 amount,
        uint256 timeElapsed
    ) public {
        amount = bound(amount, 1000, 10_000 * 10 ** 18);
        timeElapsed = bound(timeElapsed, 1, 1 days - 1);

        vm.startPrank(bob);
        token.approve(address(staking), amount);
        staking.stake(amount);

        vm.warp(block.timestamp + timeElapsed);

        uint256 balanceBefore = token.balanceOf(bob);
        staking.withdraw(amount);
        uint256 balanceAfter = token.balanceOf(bob);
        vm.stopPrank();

        uint256 expectedPenalty = (amount * 1000) / 10000;
        uint256 expectedReceived = amount - expectedPenalty;

        assertEq(balanceAfter - balanceBefore, expectedReceived);
    }

    function testFuzz_RewardsOverTime(
        uint256 stakeAmount,
        uint256 timeElapsed
    ) public {
        stakeAmount = bound(stakeAmount, 1000 * 10 ** 18, 10_000 * 10 ** 18);
        timeElapsed = bound(timeElapsed, 1 hours, 30 days);

        vm.startPrank(alice);
        token.approve(address(staking), stakeAmount);
        staking.stake(stakeAmount);

        vm.warp(block.timestamp + timeElapsed);

        uint256 earned = staking.earned(alice);

        uint256 rewardRate = staking.rewardRate();
        uint256 expectedReward = timeElapsed * rewardRate;

        // allow 1% margin
        assertApproxEqRel(earned, expectedReward, 0.01e18);
    }

    function testFuzz_MultipleStakers(
        uint256 aliceStake,
        uint256 bobStake,
        uint256 aliceTime,
        uint256 bobTime
    ) public {
        aliceStake = bound(aliceStake, 10 * 10 ** 18, 100 * 10 ** 18);
        bobStake = bound(bobStake, 10 * 10 ** 18, 100 * 10 ** 18);
        aliceTime = bound(aliceTime, 1 hours, 30 days);
        bobTime = bound(bobTime, 1 hours, 30 days);

        vm.startPrank(alice);
        token.approve(address(staking), aliceStake);
        staking.stake(aliceStake);
        vm.stopPrank();

        vm.warp(block.timestamp + aliceTime);

        vm.startPrank(bob);
        token.approve(address(staking), bobStake);
        staking.stake(bobStake);
        vm.stopPrank();

        vm.warp(block.timestamp + bobTime);

        uint256 aliceEarned = staking.earned(alice);
        uint256 bobEarned = staking.earned(bob);

        assertTrue(aliceEarned > 0);

        assertTrue(bobEarned > 0);

        uint256 totalTime = aliceTime + bobTime;
        uint256 maxRewards = totalTime * staking.rewardRate();

        assertTrue(aliceEarned + bobEarned <= maxRewards * 2);
    }

    function testFuzz_partialWithdrawals(
        uint256 stakeAmount,
        uint256 firstWithdrawal,
        uint256 secondWithdrawal
    ) public {
        stakeAmount = bound(stakeAmount, 1000 * 10 ** 18, 10_000 * 10 ** 18);
        firstWithdrawal = bound(
            firstWithdrawal,
            100 * 10 ** 18,
            stakeAmount / 2
        );

        uint256 remaining = stakeAmount - firstWithdrawal;
        secondWithdrawal = bound(secondWithdrawal, 1, remaining);

        vm.startPrank(alice);
        token.approve(address(staking), stakeAmount);
        staking.stake(stakeAmount);

        vm.warp(block.timestamp + 1 days + 12);

        staking.withdraw(firstWithdrawal);
        assertEq(staking.stakedBalance(alice), stakeAmount - firstWithdrawal);

        staking.withdraw(secondWithdrawal);
        assertEq(
            staking.stakedBalance(alice),
            stakeAmount - firstWithdrawal - secondWithdrawal
        );

        vm.stopPrank();
    }

    function testFuzz_RewardRateChange(
        uint256 newRate,
        uint256 stakeTime
    ) public {
        newRate = bound(newRate, 0.1 * 10 ** 18, 2.5 * 10 ** 18);
        stakeTime = bound(stakeTime, 1 hours, 1 days);

        uint256 stakeAmount = 10_000 * 10 ** 18;

        vm.startPrank(alice);
        token.approve(address(staking), stakeAmount);
        staking.stake(stakeAmount);
        vm.stopPrank();

        vm.warp(block.timestamp + stakeTime / 2);
        uint256 earnedFirst = staking.earned(alice);

        vm.prank(owner);
        staking.setRewardRate(newRate);

        vm.warp(block.timestamp + stakeTime / 2);

        uint256 earnedTotal = staking.earned(alice);

        assertTrue(earnedTotal > 0);
        assertTrue(earnedTotal - earnedFirst > 0);
    }

    function testFuzz_ClaimAndRestake(
        uint256 cycles,
        uint256 stakeAmount
    ) public {
        cycles = bound(cycles, 1, 5);
        stakeAmount = bound(stakeAmount, 1 * 10 ** 18, 10 * 10 ** 18);

        console.log("Bounded cycles:", cycles);
        console.log("Bounded stakeAmount:", stakeAmount);
        console.log("Total needed:", cycles * stakeAmount);
        console.log("Alice balance:", token.balanceOf(alice));

        vm.startPrank(alice);

        // Approve entire amount upfront
        token.approve(address(staking), cycles * stakeAmount);

        for (uint256 i = 0; i < cycles; i++) {
            staking.stake(stakeAmount);
            vm.warp(block.timestamp + 1 days);

            // Only claim if rewards > 0
            if (staking.rewards(alice) > 0) {
                staking.claimRewards();
            }
        }

        vm.stopPrank();

        assertEq(staking.stakedBalance(alice), stakeAmount * cycles);
    }
}
