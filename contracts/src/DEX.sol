// SPDX-License-Identifier: MIT

pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";

// LPToken is an ERC20 token representing liquidity provider shares

contract LPToken is ERC20 {
    address public immutable dex;

    constructor(string memory name, string memory symbol) ERC20(name, symbol) {
        dex = msg.sender;
    }

    function mint(address to, uint256 amount) external {
        require(msg.sender == dex, "Only DEX can mint.");
        _mint(to, amount);
    }

    function burn(address from, uint256 amount) external {
        require(msg.sender == dex, "Only DEX can burn.");
        _burn(from, amount);
    }
}

// Basic uniswap v2 AMM (x*y=constant)
contract DEX2178 is ReentrancyGuard {
    using SafeERC20 for IERC20;

    IERC20 public immutable token0;
    IERC20 public immutable token1;
    LPToken public immutable lpToken;

    uint256 public reserve0;
    uint256 public reserve1;

    uint256 private constant FEE = 3;
    uint256 private constant FEE_DENOMINATOR = 1000;
    uint256 private constant MINIMUM_LIQUIDITY = 1000;

    event Swap(
        address indexed user,
        address tokenIn,
        uint256 amountIn,
        uint256 amountOut
    );

    event LiquidityAdded(
        address indexed provider,
        uint256 amount0,
        uint256 amount1,
        uint256 lpAmount
    );

    event LiquidityRemoved(
        address indexed provider,
        uint256 amount0,
        uint256 amount1,
        uint256 lpAmount
    );

    constructor(address _token0, address _token1) {
        require(_token0!=_token1, "Identical tokens");
        require(_token0!=address(0) && _token1!=address(0), "Zero address");

        token0 = IERC20(_token0);
        token1 = IERC20(_token1);

        lpToken = new LPToken("DEX LP Token", "DLP");
    }

    // Adding liquidity to the pool (token0 and token1 added and returns lpToken to the user for providing liquidity)
    function AddLiquidity(uint256 amount0, uint256 amount1) external nonReentrant returns (uint256 lpAmount) {
        require(amount0>0 && amount1>0, "Amounts must be > 0.");

        if (reserve0==0 && reserve1==0) {
            // if no reserves, lp amount should be mroe than min liq
            lpAmount = sqrt(amount0*amount1);
            require(lpAmount>MINIMUM_LIQUIDITY, "Insufficient liquidity");

            lpToken.mint(address(1), MINIMUM_LIQUIDITY);
            lpAmount -= MINIMUM_LIQUIDITY;
        } else {
            uint256 amount1Optimal = (amount0*reserve1)/reserve0;
            if (amount1Optimal <= amount1) {
                amount1 = amount1Optimal;
            } else {
                amount0 = (amount1*reserve0) / reserve1;
            }

            lpAmount = min((amount0*lpToken.totalSupply())/reserve0, (amount1*lpToken.totalSupply())/reserve1);
        }

        require(lpAmount>0, "LP amount is too small.");
        // send tokens from user to the pool
        token0.safeTransferFrom(msg.sender, address(this), amount0);
        token1.safeTransferFrom(msg.sender, address(this), amount1);

        // add amounts to pool reserves
        reserve0 += amount0;
        reserve1 += amount1;

        // give user the lpToken for proof of providing liquidity
        lpToken.mint(msg.sender, lpAmount);

        emit LiquidityAdded(msg.sender, amount0, amount1, lpAmount);
    }

    function RemoveLiquidity(uint256 lpAmount) external nonReentrant returns (uint256 amount0, uint256 amount1) {
        require(lpAmount>0, "Amount must be more than 0.");

        uint256 totalSupply = lpToken.totalSupply();

        // calculating the tokens to return to user
        amount0 = (lpAmount*reserve0)/totalSupply;
        amount1 = (lpAmount*reserve1)/totalSupply;

        require(amount0>0&&amount1>0, "Insufficient liquidity burned");

        lpToken.burn(msg.sender, lpAmount);

        reserve0 -= amount0;
        reserve1 -= amount1;

        token0.safeTransfer(msg.sender, amount0);
        token1.safeTransfer(msg.sender, amount1);

        emit LiquidityRemoved(msg.sender, amount0, amount1, lpAmount);
    }

    // function to swap tokens: amount of token out is set to require a minimum amount(not exact)
    function swap(address tokenIn, uint256 amountIn, uint256 minAmountOut) external nonReentrant returns (uint256 amountOut) {
        require(amountIn>0, "Amount must be more than 0.");
        require(tokenIn==address(token0)||tokenIn==address(token1), "Invalid token");

        bool isToken0 = tokenIn == address(token0);

        (IERC20 tIn, IERC20 tOut, uint256 resIn, uint256 resOut) = isToken0 ? (token0, token1, reserve0, reserve1) : (token1, token0, reserve1, reserve0);

        // transfer in
        tIn.safeTransferFrom(msg.sender, address(this), amountIn);

        // calculate output with fees taken
        uint256 amountInWithFee = amountIn * (FEE_DENOMINATOR-FEE);
        amountOut = (amountInWithFee * resOut) / (resIn * FEE_DENOMINATOR + amountInWithFee);

        require(amountOut >= minAmountOut, "Slippage exceeded");
        require(amountOut < resOut, "Insufficient liquidity");

        if (isToken0) {
            reserve0 += amountIn;
            reserve1 -= amountOut;
        } else {
            reserve1 += amountIn;
            reserve0 -= amountOut;
        }

        tOut.safeTransfer(msg.sender, amountOut);

        emit Swap(msg.sender, tokenIn, amountIn, amountOut);
    }

    function getAmountOut(address tokenIn, uint256 amountIn) public view returns (uint256 amountOut) {
        require(amountIn >0, "Amount must be more than 0.");
        require(
            tokenIn == address(token0) || tokenIn == address(token1),
            "Invalid token"
        );

        bool isToken0 = tokenIn == address(token0);

        (uint256 resIn, uint256 resOut) = isToken0 ? (reserve0, reserve1) : (reserve1, reserve0);

        uint256 amountInWithFee = amountIn * (FEE_DENOMINATOR - FEE);
        amountOut = (amountInWithFee * resOut) / (resIn * FEE_DENOMINATOR + amountInWithFee);
    }

    function getPrice0() external view returns (uint256 price) {
        require(reserve0>0, "No liquidity");
        price = (reserve1*1e18)/reserve0;
    }

    function getPrice1() external view returns (uint256 price) {
        require(reserve1>0, "No liquidity");
        price = (reserve0*1e18)/reserve1;
    }

    function sqrt(uint256 y) private pure returns (uint256 z) {
        if (y > 3) {
            z = y;
            uint256 x = y / 2 + 1;
            while (x < z) {
                z = x;
                x = (y / x + x) / 2;
            }
        } else if (y != 0) {
            z = 1;
        }
    }

    function min(uint256 a, uint256 b) private pure returns (uint256) {
        return a < b ? a : b;
    }
}