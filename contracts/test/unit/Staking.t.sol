// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../../src/Token.sol";
import "../../src/Staking.sol";

/**
 * @title StakingTest
 * @notice Unit tests for Staking contract
 * 
 * Test Coverage:
 * - Deployment and initialization
 * - Staking tokens
 * - Reward calculation over time
 * - Withdrawing with/without penalty
 * - Claiming rewards
 * - Emergency withdrawal
 * - Owner functions
 */
contract StakingTest is Test {
    Token2178 public token;
    Staking public staking;
    
    address public owner = address(1);
    address public alice = address(2);
    address public bob = address(3);
    
    // Initial token distribution
    uint256 constant INITIAL_SUPPLY = 100_000 * 10**18;
    uint256 constant USER_BALANCE = 10_000 * 10**18;
    uint256 constant REWARD_POOL = 100_000 * 10**18;  // 100K tokens for rewards
    
    function setUp() public {
        // Deploy token as owner
        vm.prank(owner);
        token = new Token2178();

        vm.prank(owner);
        token.mint(owner, REWARD_POOL);
        
        // Deploy staking as owner (same token for staking and rewards)
        vm.prank(owner);
        staking = new Staking(address(token), address(token));
        
        // Give users tokens
        vm.startPrank(owner);
        token.transfer(alice, USER_BALANCE);
        token.transfer(bob, USER_BALANCE);
        
        // Fund staking contract with reward tokens
        token.transfer(address(staking), REWARD_POOL);
        vm.stopPrank();
    }
    
    // ========== DEPLOYMENT TESTS ==========
    
    function testDeployment() public {
        assertEq(address(staking.stakingToken()), address(token));
        assertEq(address(staking.rewardToken()), address(token));
        assertEq(staking.rewardRate(), 1 * 10**18);
        assertEq(staking.totalStaked(), 0);
        assertEq(staking.owner(), owner);
    }
    
    // ========== STAKING TESTS ==========
    
    function testStake() public {
        uint256 stakeAmount = 1000 * 10**18;
        
        vm.startPrank(alice);
        token.approve(address(staking), stakeAmount);
        staking.stake(stakeAmount);
        vm.stopPrank();
        
        assertEq(staking.stakedBalance(alice), stakeAmount);
        assertEq(staking.totalStaked(), stakeAmount);
        assertEq(token.balanceOf(alice), USER_BALANCE - stakeAmount);
        assertEq(token.balanceOf(address(staking)), REWARD_POOL + stakeAmount);
    }
    
    function testStakeZeroAmount() public {
        vm.startPrank(alice);
        token.approve(address(staking), 1000);
        
        vm.expectRevert("Cannot stake 0.");
        staking.stake(0);
        vm.stopPrank();
    }
    
    function testStakeWithoutApproval() public {
        vm.prank(alice);
        vm.expectRevert();  // Will revert with ERC20 transfer error
        staking.stake(1000);
    }
    
    function testMultipleStakes() public {
        vm.startPrank(alice);
        token.approve(address(staking), 2000 * 10**18);
        
        staking.stake(1000 * 10**18);
        assertEq(staking.stakedBalance(alice), 1000 * 10**18);
        
        staking.stake(500 * 10**18);
        assertEq(staking.stakedBalance(alice), 1500 * 10**18);
        
        vm.stopPrank();
    }
    
    // ========== REWARD CALCULATION TESTS ==========
    
    function testRewardsSingleStaker() public {
        uint256 stakeAmount = 1000 * 10**18;
        
        // Alice stakes
        vm.startPrank(alice);
        token.approve(address(staking), stakeAmount);
        staking.stake(stakeAmount);
        vm.stopPrank();
        
        // Fast forward 1 hour
        vm.warp(block.timestamp + 1 hours);
        
        // Check earned rewards
        // 1 hour = 3600 seconds
        // 3600 seconds * 100 tokens/sec = 360,000 tokens
        uint256 expectedReward = 3600 * 1 * 10**18;
        uint256 earned = staking.earned(alice);
        
        assertEq(earned, expectedReward);
    }
    
    function testRewardsMultipleStakers() public {
        uint256 stakeAmount = 1000 * 10**18;
        
        // Alice stakes at T=0
        vm.startPrank(alice);
        token.approve(address(staking), stakeAmount);
        staking.stake(stakeAmount);
        vm.stopPrank();
        
        // Fast forward 1 hour
        vm.warp(block.timestamp + 1 hours);
        
        // Bob stakes at T=1hr
        vm.startPrank(bob);
        token.approve(address(staking), stakeAmount);
        staking.stake(stakeAmount);
        vm.stopPrank();
        
        // Fast forward another hour (T=2hr)
        vm.warp(block.timestamp + 1 hours);
        
        // Alice earned:
        // - 1 hour alone: 3600 * 100 = 360,000 tokens
        // - 1 hour with Bob (50%): 3600 * 100 * 0.5 = 180,000 tokens
        // Total: 540,000 tokens
        uint256 aliceExpected = (3600 * 1 * 10**18) + (3600 * 1 * 10**18 / 2);
        uint256 aliceEarned = staking.earned(alice);
        assertApproxEqRel(aliceEarned, aliceExpected, 0.01e18);  // 1% tolerance
        
        // Bob earned:
        // - 1 hour with Alice (50%): 3600 * 100 * 0.5 = 180,000 tokens
        uint256 bobExpected = 3600 * 1 * 10**18 / 2;
        uint256 bobEarned = staking.earned(bob);
        assertApproxEqRel(bobEarned, bobExpected, 0.01e18);
    }
    
    function testRewardsAfterClaim() public {
        uint256 stakeAmount = 1000 * 10**18;
        
        // Alice stakes
        vm.startPrank(alice);
        token.approve(address(staking), stakeAmount);
        staking.stake(stakeAmount);
        vm.stopPrank();
        
        // Fast forward 1 hour
        vm.warp(block.timestamp + 1 hours);
        
        // Alice claims
        vm.prank(alice);
        staking.claimRewards();
        
        // Check earned is now 0
        assertEq(staking.earned(alice), 0);
        
        // Fast forward another hour
        vm.warp(block.timestamp + 1 hours);
        
        // Should earn for the second hour
        uint256 expectedReward = 3600 * 1 * 10**18;
        assertApproxEqRel(staking.earned(alice), expectedReward, 0.01e18);
    }
    
    // ========== WITHDRAWAL TESTS ==========
    
    function testWithdrawWithoutPenalty() public {
        uint256 stakeAmount = 1000 * 10**18;
        
        // Alice stakes
        vm.startPrank(alice);
        token.approve(address(staking), stakeAmount);
        staking.stake(stakeAmount);
        vm.stopPrank();
        
        // Fast forward 1 day + 1 second (past MIN_STAKE_DURATION)
        vm.warp(block.timestamp + 1 days + 1);
        
        // Alice withdraws
        uint256 balanceBefore = token.balanceOf(alice);
        
        vm.prank(alice);
        staking.withdraw(stakeAmount);
        
        uint256 balanceAfter = token.balanceOf(alice);
        
        // Should receive full amount (no penalty)
        assertEq(balanceAfter - balanceBefore, stakeAmount);
        assertEq(staking.stakedBalance(alice), 0);
        assertEq(staking.totalStaked(), 0);
    }
    
    function testWithdrawWithPenalty() public {
        uint256 stakeAmount = 1000 * 10**18;
        
        // Alice stakes
        vm.startPrank(alice);
        token.approve(address(staking), stakeAmount);
        staking.stake(stakeAmount);
        vm.stopPrank();
        
        // Fast forward only 12 hours (less than 1 day)
        vm.warp(block.timestamp + 12 hours);
        
        // Alice withdraws early
        uint256 aliceBalanceBefore = token.balanceOf(alice);
        uint256 ownerBalanceBefore = token.balanceOf(owner);
        
        vm.prank(alice);
        staking.withdraw(stakeAmount);
        
        uint256 aliceBalanceAfter = token.balanceOf(alice);
        uint256 ownerBalanceAfter = token.balanceOf(owner);
        
        // Calculate expected penalty (10% = 1000/10000)
        uint256 expectedPenalty = (stakeAmount * 1000) / 10000;
        uint256 expectedReceived = stakeAmount - expectedPenalty;
        
        // Check Alice received amount minus penalty
        assertEq(aliceBalanceAfter - aliceBalanceBefore, expectedReceived);
        
        // Check owner received penalty
        assertEq(ownerBalanceAfter - ownerBalanceBefore, expectedPenalty);
        
        assertEq(staking.stakedBalance(alice), 0);
    }
    
    function testWithdrawZeroAmount() public {
        vm.prank(alice);
        vm.expectRevert("Cannot withdraw 0.");
        staking.withdraw(0);
    }
    
    function testWithdrawMoreThanStaked() public {
        uint256 stakeAmount = 1000 * 10**18;
        
        vm.startPrank(alice);
        token.approve(address(staking), stakeAmount);
        staking.stake(stakeAmount);
        
        vm.expectRevert("Insufficient balance.");
        staking.withdraw(stakeAmount + 1);
        vm.stopPrank();
    }
    
    function testPartialWithdraw() public {
        uint256 stakeAmount = 1000 * 10**18;
        
        vm.startPrank(alice);
        token.approve(address(staking), stakeAmount);
        staking.stake(stakeAmount);
        vm.stopPrank();
        
        // Wait past penalty period
        vm.warp(block.timestamp + 1 days + 1);
        
        // Withdraw half
        vm.prank(alice);
        staking.withdraw(500 * 10**18);
        
        assertEq(staking.stakedBalance(alice), 500 * 10**18);
        assertEq(staking.totalStaked(), 500 * 10**18);
    }
    
    // ========== CLAIM REWARDS TESTS ==========
    
    function testClaimRewards() public {
        uint256 stakeAmount = 1000 * 10**18;
        
        // Alice stakes
        vm.startPrank(alice);
        token.approve(address(staking), stakeAmount);
        staking.stake(stakeAmount);
        vm.stopPrank();
        
        // Fast forward 1 hour
        vm.warp(block.timestamp + 1 hours);
        
        uint256 expectedReward = 3600 * 1 * 10**18;
        uint256 balanceBefore = token.balanceOf(alice);
        
        // Alice claims
        vm.prank(alice);
        staking.claimRewards();
        
        uint256 balanceAfter = token.balanceOf(alice);
        
        assertApproxEqRel(balanceAfter - balanceBefore, expectedReward, 0.01e18);
        assertEq(staking.rewards(alice), 0);
    }
    
    function testClaimRewardsWithNoRewards() public {
        vm.prank(alice);
        vm.expectRevert("No rewards.");
        staking.claimRewards();
    }
    
    function testRewardsAccumulateAfterWithdraw() public {
        uint256 stakeAmount = 1000 * 10**18;
        
        // Alice stakes
        vm.startPrank(alice);
        token.approve(address(staking), stakeAmount);
        staking.stake(stakeAmount);
        vm.stopPrank();
        
        // Fast forward 1 hour
        vm.warp(block.timestamp + 1 hours);
        
        // Alice withdraws (but doesn't claim rewards)
        vm.warp(block.timestamp + 1 days);  // Wait for no penalty
        vm.prank(alice);
        staking.withdraw(stakeAmount);
        
        // Rewards should still be claimable
        uint256 earned = staking.rewards(alice);
        assertTrue(earned > 0);
        
        // Can still claim even after withdrawing all stake
        vm.prank(alice);
        staking.claimRewards();
        
        assertEq(staking.rewards(alice), 0);
    }
    
    // ========== EMERGENCY WITHDRAW TESTS ==========
    
    function testEmergencyWithdraw() public {
        uint256 stakeAmount = 1000 * 10**18;
        
        // Alice stakes
        vm.startPrank(alice);
        token.approve(address(staking), stakeAmount);
        staking.stake(stakeAmount);
        vm.stopPrank();
        
        // Fast forward to earn some rewards
        vm.warp(block.timestamp + 1 hours);
        
        // Check Alice has earned rewards
        uint256 earned = staking.earned(alice);
        assertTrue(earned > 0);
        
        uint256 balanceBefore = token.balanceOf(alice);
        
        // Emergency withdraw
        vm.prank(alice);
        staking.emergencyWithdraw();
        
        uint256 balanceAfter = token.balanceOf(alice);
        
        // Should get principal back
        assertEq(balanceAfter - balanceBefore, stakeAmount);
        
        // Should forfeit rewards
        assertEq(staking.rewards(alice), 0);
        assertEq(staking.stakedBalance(alice), 0);
    }
    
    function testEmergencyWithdrawWithNoStake() public {
        vm.prank(alice);
        vm.expectRevert("Nothing staked.");
        staking.emergencyWithdraw();
    }
    
    // ========== OWNER FUNCTIONS TESTS ==========
    
    function testSetRewardRate() public {
        uint256 newRate = 2 * 10**18;
        
        vm.prank(owner);
        staking.setRewardRate(newRate);
        
        assertEq(staking.rewardRate(), newRate);
    }
    
    function testSetRewardRateNotOwner() public {
        vm.prank(alice);
        vm.expectRevert();
        staking.setRewardRate(2 * 10**18);
    }
    
    function testSetRewardRateUpdatesRewards() public {
        uint256 stakeAmount = 1000 * 10**18;
        
        // Alice stakes
        vm.startPrank(alice);
        token.approve(address(staking), stakeAmount);
        staking.stake(stakeAmount);
        vm.stopPrank();
        
        // Fast forward 1 hour at 100 tokens/sec
        vm.warp(block.timestamp + 1 hours);
        uint256 earnedAt100 = staking.earned(alice);
        
        // Owner changes rate to 2 tokens/sec
        vm.prank(owner);
        staking.setRewardRate(2 * 10**18);
        
        // Fast forward another hour at 200 tokens/sec
        vm.warp(block.timestamp + 1 hours);
        uint256 earnedAt200 = staking.earned(alice);
        
        // Second hour should have earned double
        uint256 firstHour = 3600 * 1 * 10**18;
        uint256 secondHour = 3600 * 2 * 10**18;
        uint256 expected = firstHour + secondHour;
        
        assertApproxEqRel(earnedAt200, expected, 0.01e18);
    }
    
    function testSetEarlyWithdrawalPenalty() public {
        vm.prank(owner);
        staking.setEarlyWithdrawalPenalty(2000);  // 20%
        
        assertEq(staking.earlyWithdrawalPenalty(), 2000);
    }
    
    function testSetEarlyWithdrawalPenaltyTooHigh() public {
        vm.prank(owner);
        vm.expectRevert("Penalty too high.");
        staking.setEarlyWithdrawalPenalty(6000);  // 60% > max 50%
    }
    
    function testRecoverToken() public {
        // Deploy different token
        vm.prank(owner);
        Token2178 otherToken = new Token2178();
        
        // Send some to staking contract by mistake
        vm.prank(owner);
        otherToken.transfer(address(staking), 1000 * 10**18);
        
        // Owner recovers
        uint256 balanceBefore = otherToken.balanceOf(owner);
        
        vm.prank(owner);
        staking.recoverToken(address(otherToken), 1000 * 10**18);
        
        uint256 balanceAfter = otherToken.balanceOf(owner);
        assertEq(balanceAfter - balanceBefore, 1000 * 10**18);
    }
    
    function testRecoverTokenCannotRecoverStakingToken() public {
        vm.prank(owner);
        vm.expectRevert("Cannot recover staking token.");
        staking.recoverToken(address(token), 1000);
    }
    
    // ========== HELPER FUNCTION TESTS ==========
    
    function testGetStakeInfo() public {
        uint256 stakeAmount = 1000 * 10**18;
        
        // Alice stakes
        vm.startPrank(alice);
        token.approve(address(staking), stakeAmount);
        staking.stake(stakeAmount);
        vm.stopPrank();
        
        // Fast forward 12 hours (less than MIN_STAKE_DURATION)
        vm.warp(block.timestamp + 12 hours);
        
        (uint256 staked, uint256 earned, bool canWithdrawWithoutPenalty) = 
            staking.getStakeInfo(alice);
        
        assertEq(staked, stakeAmount);
        assertTrue(earned > 0);
        assertFalse(canWithdrawWithoutPenalty);
        
        // Fast forward past MIN_STAKE_DURATION
        vm.warp(block.timestamp + 13 hours);
        
        (,, canWithdrawWithoutPenalty) = staking.getStakeInfo(alice);
        assertTrue(canWithdrawWithoutPenalty);
    }
}