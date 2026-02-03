// SPDX-License-Identifier: MIT

pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../../src/DEX.sol";
import "../../src/Token.sol";

contract DEXTest is Test {
    Token2178 public token0;
    Token2178 public token1;
    DEX2178 public dex;

    address public owner = address(this);
    address public alice = address(0x21);
    address public bob = address(0x78);

    function setUp() public {
        token0 = new Token2178();
        token1 = new Token2178();

        dex = new DEX2178(address(token0), address(token1));

        token0.transfer(alice, 10_000 * 10 ** 18);
        token1.transfer(alice, 10_000 * 10 ** 18);

        token0.transfer(bob, 10_000 * 10 ** 18);
        token1.transfer(bob, 10_000 * 10 ** 18);
    }

    function testDeployment() public {
        assertEq(address(dex.token0()), address(token0));
        assertEq(address(dex.token1()), address(token1));
        assertEq(dex.reserve0(), 0);
        assertEq(dex.reserve1(), 0);

        assertTrue(address(dex.lpToken()) != address(0));
    }

    function testCannotDeployWithIdenticalTokens() public {
        vm.expectRevert("Identical tokens");
        new DEX2178(address(token0), address(token0));
    }

    function testCannotDeployWithZeroAddress() public {
        vm.expectRevert("Zero address");
        new DEX2178(address(0), address(token1));

        vm.expectRevert("Zero address");
        new DEX2178(address(token0), address(0));
    }

    function testAddInitialLiquidity() public {
        uint256 amount0 = 1000 * 10 ** 18;
        uint256 amount1 = 1000 * 10 ** 18;

        vm.startPrank(alice);
        token0.approve(address(dex), amount0);
        token1.approve(address(dex), amount1);

        uint256 lpAmount = dex.AddLiquidity(amount0, amount1);
        vm.stopPrank();

        uint256 expectedLP = 1000 * 10 ** 18 - 1000;

        assertEq(lpAmount, expectedLP);
        assertEq(dex.reserve0(), amount0);
        assertEq(dex.reserve1(), amount1);
        assertEq(dex.lpToken().balanceOf(alice), expectedLP);

        // min liqiuidity check
        assertEq(dex.lpToken().balanceOf(address(1)), 1000);
    }

    function testAddLiquidityAfterInitial() public {
        vm.startPrank(alice);
        token0.approve(address(dex), 1000 * 10 ** 18);
        token1.approve(address(dex), 1000 * 10 ** 18);
        dex.AddLiquidity(1000 * 10 ** 18, 1000 * 10 ** 18);
        vm.stopPrank();

        vm.startPrank(bob);
        token0.approve(address(dex), 500 * 10 ** 18);
        token1.approve(address(dex), 500 * 10 ** 18);

        uint256 lpAmount = dex.AddLiquidity(500 * 10 ** 18, 500 * 10 ** 18);
        vm.stopPrank();

        assertEq(lpAmount, 500 * 10 ** 18);
        assertEq(dex.reserve0(), 1500 * 10 ** 18);
        assertEq(dex.reserve1(), 1500 * 10 ** 18);
    }

    function testAddLiquidityAdjustsRatio() public {
        vm.startPrank(alice);
        token0.approve(address(dex), 1000 * 10 ** 18);
        token1.approve(address(dex), 1000 * 10 ** 18);
        dex.AddLiquidity(1000 * 10 ** 18, 1000 * 10 ** 18);
        vm.stopPrank();

        vm.startPrank(bob);
        token0.approve(address(dex), 500 * 10 ** 18);
        token1.approve(address(dex), 1000 * 10 ** 18);
        uint256 lpAmount = dex.AddLiquidity(500 * 10 ** 18, 1000 * 10 ** 18);
        vm.stopPrank();

        assertEq(dex.reserve0(), 1500 * 10 ** 18);
        assertEq(dex.reserve1(), 1500 * 10 ** 18);
        assertEq(lpAmount, 500 * 10 ** 18);
    }

    function testAddLiquidityRevertsWithZero() public {
        vm.startPrank(alice);
        token0.approve(address(dex), 1000 * 10 ** 18);
        token1.approve(address(dex), 1000 * 10 ** 18);

        vm.expectRevert("Amounts must be > 0.");
        dex.AddLiquidity(0, 1000 * 10 ** 18);

        vm.expectRevert("Amounts must be > 0.");
        dex.AddLiquidity(1000 * 10 ** 18, 0);

        vm.stopPrank();
    }

    function testAddLiquidityRevertsInsufficientInitial() public {
        vm.startPrank(alice);
        token0.approve(address(dex), 100);
        token1.approve(address(dex), 100);

        vm.expectRevert("Insufficient liquidity");
        dex.AddLiquidity(100, 100);
        vm.stopPrank();
    }

    function testRemoveLiquidity() public {
        vm.startPrank(alice);
        token0.approve(address(dex), 1000 * 10 ** 18);
        token1.approve(address(dex), 1000 * 10 ** 18);
        uint256 lpAmount = dex.AddLiquidity(1000 * 10 ** 18, 1000 * 10 ** 18);

        uint256 balanceBefore0 = token0.balanceOf(alice);
        uint256 balanceBefore1 = token1.balanceOf(alice);

        // remove half liquidiity
        (uint256 amount0, uint256 amount1) = dex.RemoveLiquidity(lpAmount / 2);
        vm.stopPrank();

        assertApproxEqRel(amount0, 500 * 10 ** 18, 0.01e18);
        assertApproxEqRel(amount1, 500 * 10 ** 18, 0.01e18);
    }

    function testRemoveAllLiquidity() public {
        vm.startPrank(alice);
        token0.approve(address(dex), 1000 * 10 ** 18);
        token1.approve(address(dex), 1000 * 10 ** 18);
        uint256 lpAmount = dex.AddLiquidity(1000 * 10 ** 18, 1000 * 10 ** 18);

        (uint256 amount0, uint256 amount1) = dex.RemoveLiquidity(lpAmount);
        vm.stopPrank();

        assertTrue(amount0 > 999 * 10 ** 18);
        assertTrue(amount1 > 999 * 10 ** 18);
    }

    function testRemoveLiquidityRevertsZero() public {
        vm.expectRevert("Amount must be more than 0.");
        dex.RemoveLiquidity(0);
    }

    function testRemoveLiquidityRevertsInsufficientBalance() public {
        vm.startPrank(alice);

        vm.expectRevert();
        dex.RemoveLiquidity(100 * 10 ** 18);

        vm.stopPrank();
    }

    function testSwap() public {
        vm.startPrank(alice);
        token0.approve(address(dex), 1000 * 10 ** 18);
        token1.approve(address(dex), 1000 * 10 ** 18);
        dex.AddLiquidity(1000 * 10 ** 18, 1000 * 10 ** 18);
        vm.stopPrank();

        vm.startPrank(bob);
        token0.approve(address(dex), 100 * 10 ** 18);
        uint256 balanceBefore = token1.balanceOf(bob);
        uint256 amountOut = dex.swap(
            address(token0),
            100 * 10 ** 18,
            90 * 10 ** 18
        );
        uint256 balanceAfter = token1.balanceOf(bob);

        vm.stopPrank();

        assertTrue(amountOut > 0);
        assertTrue(amountOut < 100 * 10 ** 18);
        assertEq(balanceAfter - balanceBefore, amountOut);
    }

    function testSwapRevertsSlippageExceeded() public {
        vm.startPrank(alice);
        token0.approve(address(dex), 1000 * 10 ** 18);
        token1.approve(address(dex), 1000 * 10 ** 18);
        dex.AddLiquidity(1000 * 10 ** 18, 1000 * 10 ** 18);
        vm.stopPrank();

        vm.startPrank(bob);
        token0.approve(address(dex), 100 * 10 ** 18);

        vm.expectRevert("Slippage exceeded");
        dex.swap(address(token0), 100 * 10 ** 18, 100 * 10 ** 18);

        vm.stopPrank();
    }

    function testSwapRevertsInvalidToken() public {
        vm.startPrank(bob);

        vm.expectRevert("Invalid token");
        dex.swap(address(0x123), 100, 0);

        vm.stopPrank();
    }

    function testSwapRevertsZeroAmount() public {
        vm.startPrank(bob);

        vm.expectRevert("Amount must be more than 0.");
        dex.swap(address(token0), 0, 0);

        vm.stopPrank();
    }

    function testSwapBothDirections() public {
        vm.startPrank(alice);
        token0.approve(address(dex), 1000 * 10 ** 18);
        token1.approve(address(dex), 1000 * 10 ** 18);
        dex.AddLiquidity(1000 * 10 ** 18, 1000 * 10 ** 18);
        vm.stopPrank();

        vm.startPrank(bob);
        token0.approve(address(dex), 100 * 10 ** 18);
        uint256 out1 = dex.swap(address(token0), 100 * 10 ** 18, 0);

        token1.approve(address(dex), out1);
        uint256 out0 = dex.swap(address(token1), out1, 0);
        vm.stopPrank();

        // should get back less than original due to fees
        assertTrue(out0 < 100 * 10 ** 18);
    }

    function testgetAmountOut() public {
        vm.startPrank(alice);
        token0.approve(address(dex), 1000 * 10 ** 18);
        token1.approve(address(dex), 1000 * 10 ** 18);
        dex.AddLiquidity(1000 * 10 ** 18, 1000 * 10 ** 18);
        vm.stopPrank();

        uint256 amountOut = dex.getAmountOut(address(token0), 100 * 10 ** 18);

        assertApproxEqRel(amountOut, 99.69 * 10 ** 18, 0.02e18);
    }

    function testgetPrice() public {
        vm.startPrank(alice);
        token0.approve(address(dex), 1000 * 10 ** 18);
        token1.approve(address(dex), 2000 * 10 ** 18);
        dex.AddLiquidity(1000 * 10 ** 18, 2000 * 10 ** 18);
        vm.stopPrank();

        uint256 price0 = dex.getPrice0();
        assertEq(price0, 2 * 10 ** 18);

        uint256 price1 = dex.getPrice1();
        assertEq(price1, 0.5 * 10 ** 18);
    }

    function testGetPriceRevertsNoLiquidity() public {
        vm.expectRevert("No liquidity");
        dex.getPrice0();

        vm.expectRevert("No liquidity");
        dex.getPrice1();
    }

    function testConstantProductMaintained() public {
        vm.startPrank(alice);
        token0.approve(address(dex), 1000 * 10 ** 18);
        token1.approve(address(dex), 1000 * 10 ** 18);
        dex.AddLiquidity(1000 * 10 ** 18, 1000 * 10 ** 18);
        vm.stopPrank();

        uint256 kBefore = dex.reserve0() * dex.reserve1();

        vm.startPrank(bob);
        token0.approve(address(dex), 100 * 10 ** 18);
        dex.swap(address(token0), 100 * 10 ** 18, 0);
        vm.stopPrank();

        uint256 kAfter = dex.reserve0() * dex.reserve1();

        assertTrue(kAfter >= kBefore);
    }
}
