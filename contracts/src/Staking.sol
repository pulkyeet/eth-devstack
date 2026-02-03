// SPDX-License-Identifier: MIT

pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";

contract Staking is Ownable, ReentrancyGuard {
    using SafeERC20 for IERC20;

    // immutable so its set once in the constructor. 
    // token that the users stake
    IERC20 public immutable stakingToken;

    // token that the users get for staking
    IERC20 public immutable rewardToken;

    // 1 token per second
    uint256 public rewardRate = 1 * 10 ** 18;

    // last time rewards were calculated
    uint256 public lastUpdateTime;

    uint256 public rewardPerTokenStored;

    // last snapshot of user's reward pertokenstores
    mapping(address=>uint256) public userRewardPerTokenPaid;

    // unclaimed reward for each user
    mapping(address=>uint256) public rewards;

    // amount staked by user
    mapping(address=>uint256) public stakedBalance;

    // timestamp of staking done by the user
    mapping(address=>uint256) public stakeTimestamp;

    // total staked coins
    uint256 public totalStaked;

    uint256 public constant MIN_STAKE_DURATION = 1 days;

    uint256 public earlyWithdrawalPenalty = 1000;

    uint256 public constant PENALTY_DENOMINATOR = 10000;

    event Staked(address indexed user, uint256 amount);
    event Withdrawn(address indexed user, uint256 amount, uint256 penalty);
    event RewardsClaimed(address indexed user, uint256 amount);

    // staking reward rate is updated after each new staking/withdrawal event
    event RewardRateUpdated(uint256 newRate);
    event EmergencyWithdraw(address indexed user, uint256 amount);

    constructor(address _stakingToken, address _rewardToken) Ownable(msg.sender) {
        stakingToken = IERC20(_stakingToken);
        rewardToken = IERC20(_rewardToken);
        lastUpdateTime = block.timestamp;
    }

    // calculates current reward per token (time * rewardrate) per totalstakedtokens
    function rewardPerToken() public view returns (uint256) {
        if (totalStaked == 0) {
            return rewardPerTokenStored;
        }

        uint256 timeDelta = block.timestamp - lastUpdateTime;
        uint256 newRewardPerToken = (timeDelta * rewardRate * 1e18) / totalStaked;

        return rewardPerTokenStored + newRewardPerToken;
    }

    function earned(address account) public view returns (uint256) {
        uint256 stakedAmount = stakedBalance[account];
        uint256 rewardPerTokenDelta = rewardPerToken() - userRewardPerTokenPaid[account];
        uint256 newlyEarned = (stakedAmount * rewardPerTokenDelta) / 1e18;

        return rewards[account] + newlyEarned;
    }   

    // updates each and every time a function that uses it is called.
    modifier updateReward(address account) {
        rewardPerTokenStored = rewardPerToken();
        lastUpdateTime = block.timestamp;

        if (account!=address(0)) {
            rewards[account] = earned(account);
            userRewardPerTokenPaid[account] = rewardPerTokenStored;
        }

        _;
    }

    function stake(uint256 amount) external nonReentrant updateReward(msg.sender) {
        require(amount > 0, "Cannot stake 0.");

        totalStaked += amount;
        stakedBalance[msg.sender] += amount;
        stakeTimestamp[msg.sender] = block.timestamp;

        stakingToken.safeTransferFrom(msg.sender, address(this), amount);
        emit Staked(msg.sender, amount);
    }

    function withdraw(uint256 amount) external nonReentrant updateReward(msg.sender) {
        require(amount > 0, "Cannot withdraw 0.");
        require(stakedBalance[msg.sender] >= amount, "Insufficient balance.");

        uint256 penalty = 0;

        if (block.timestamp < stakeTimestamp[msg.sender] + MIN_STAKE_DURATION) {
            penalty = (amount*earlyWithdrawalPenalty)/PENALTY_DENOMINATOR;
        }

        stakedBalance[msg.sender] -= amount;

        totalStaked -= amount;
        uint256 amountAfterPenalty = amount - penalty;

        stakingToken.safeTransfer(msg.sender, amountAfterPenalty);

        if (penalty > 0) {
            stakingToken.safeTransfer(owner(), penalty);
        }

        emit Withdrawn(msg.sender, amountAfterPenalty, penalty);
    }

    function claimRewards() external nonReentrant updateReward(msg.sender) {
        uint256 reward = rewards[msg.sender];
        require(reward > 0, "No rewards.");

        rewards[msg.sender] = 0;

        rewardToken.safeTransfer(msg.sender, reward);

        emit RewardsClaimed(msg.sender, reward);
    }

    function emergencyWithdraw() external nonReentrant {
        uint256 amount = stakedBalance[msg.sender];
        require(amount >0, "Nothing staked.");

        totalStaked -= amount;
        stakedBalance[msg.sender] = 0;
        rewards[msg.sender] = 0;

        stakingToken.safeTransfer(msg.sender, amount);

        emit EmergencyWithdraw(msg.sender, amount);
    }

    function setRewardRate(uint256 _rewardRate) external onlyOwner updateReward(address(0)) {
        rewardRate = _rewardRate;
        emit RewardRateUpdated(_rewardRate);
    }

    function setEarlyWithdrawalPenalty(uint256 _penalty) external onlyOwner {
        require(_penalty <= 5000, "Penalty too high.");
        earlyWithdrawalPenalty = _penalty;
    }

    function recoverToken(address token, uint256 amount) external onlyOwner {
        require(token!=address(stakingToken), "Cannot recover staking token.");
        IERC20(token).safeTransfer(owner(), amount);
    }

    function getStakeInfo(address account) external view returns (uint256 staked, uint256 earned_, bool canWithdrawWithoutPenalty) {
        staked = stakedBalance[account];
        earned_ = earned(account);
        canWithdrawWithoutPenalty = block.timestamp >= stakeTimestamp[account] + MIN_STAKE_DURATION;
    }
}