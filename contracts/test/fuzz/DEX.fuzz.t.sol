// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../../src/DEX.sol";
import "../../src/Token.sol";

contract DEXFuzzTest is Test {
    Token2178 public t0;
    Token2178 public t1;
    DEX2178 public dex;

    address public alice = address(0xA11CE);
    address public bob = address(0xB0B);

    function setUp() public {
        t0 = new Token2178();
        t1 = new Token2178();
        dex = new DEX2178(address(t0), address(t1));

        t0.transfer(alice, 1_000_000 * 10 ** 18);
        t0.transfer(bob, 1_000_000 * 10 ** 18);
        t1.transfer(alice, 1_000_000 * 10 ** 18);
        t1.transfer(bob, 1_000_000 * 10 ** 18);
    }

    // a0 and a1 are added amounts
    function testFuzz_AddInitialLiquidity(uint256 a0, uint256 a1) public {
        a0 = bound(a0, 1001, 100_00 * 10 ** 18);
        a1 = bound(a1, 1001, 100_00 * 10 ** 18);

        vm.startPrank(alice);
        t0.approve(address(dex), a0);
        t1.approve(address(dex), a1);

        uint256 lpAmount = dex.AddLiquidity(a0, a1);
        vm.stopPrank();

        assertEq(dex.reserve0(), a0);
        assertEq(dex.reserve1(), a1);

        assertGt(lpAmount, 0);

        uint256 expectedLP = sqrt(a0 * a1) - 1000;
        assertEq(lpAmount, expectedLP);

        assertEq(dex.lpToken().balanceOf(address(1)), 1000);
    }

    // i0 and i1 are Initial amounts.
    function testFuzz_AddLiquidityAfterInitial(
        uint256 i0,
        uint256 i1,
        uint256 a0,
        uint256 a1
    ) public {
        i0 = bound(i0, 10_000 * 10 ** 18, 50_000 * 10 ** 18);
        i1 = bound(i1, 10_000 * 10 ** 18, 50_000 * 10 ** 18);
        a0 = bound(a0, 1000 * 10 ** 18, 50_000 * 10 ** 18);
        a1 = bound(a1, 1000 * 10 ** 18, 50_000 * 10 ** 18);

        vm.startPrank(alice);
        t0.approve(address(dex), i0);
        t1.approve(address(dex), i1);
        dex.AddLiquidity(i0, i1);
        vm.stopPrank();

        uint256 r0b = dex.reserve0();
        uint256 r1b = dex.reserve1();

        vm.startPrank(bob);
        t0.approve(address(dex), a0);
        t1.approve(address(dex), a1);
        uint256 bobLP = dex.AddLiquidity(a0, a1);
        vm.stopPrank();

        assertGt(dex.reserve0(), r0b);
        assertGt(dex.reserve1(), r1b);

        assertGt(bobLP, 0);
        assertEq(dex.lpToken().balanceOf(bob), bobLP);
    }

    function testFuzz_RemoveLiquidity(
        uint256 a0,
        uint256 a1,
        uint256 removePercentage
    ) public {
        a0 = bound(a0, 1 * 10 ** 18, 100_000 * 10 ** 18);
        a1 = bound(a1, 1 * 10 ** 18, 100_000 * 10 ** 18);
        removePercentage = bound(removePercentage, 1, 100);

        vm.startPrank(alice);
        t0.approve(address(dex), a0);
        t1.approve(address(dex), a1);
        uint256 lpAmount = dex.AddLiquidity(a0, a1);

        uint256 removeAmount = (removePercentage * lpAmount) / 100;
        if (removeAmount == 0) {
            removeAmount = 1;
        }

        (uint256 rem0, uint256 rem1) = dex.RemoveLiquidity(removeAmount);
        vm.stopPrank();

        assertGt(rem0, 0);
        assertGt(rem1, 0);

        assertEq(dex.lpToken().balanceOf(alice), lpAmount - removeAmount);
    }

    function testFuzz_removeAllLiquidity(uint256 a0, uint256 a1) public {
        a0 = bound(a0, 10_000 * 10 ** 18, 100_000 * 10 ** 18);
        a1 = bound(a1, 10_000 * 10 ** 18, 100_000 * 10 ** 18);

        vm.startPrank(bob);
        t0.approve(address(dex), a0);
        t1.approve(address(dex), a1);
        uint256 lpAmount = dex.AddLiquidity(a0, a1);

        uint256 balanceBefore0 = t0.balanceOf(bob);
        uint256 balanceBefore1 = t1.balanceOf(bob);

        (uint256 rem0, uint256 rem1) = dex.RemoveLiquidity(lpAmount);
        vm.stopPrank();

        assertGt(rem0, (a0 * 95) / 100);
        assertGt(rem1, (a1 * 95) / 100);

        assertEq(t0.balanceOf(bob), balanceBefore0 + rem0);
        assertEq(t1.balanceOf(bob), balanceBefore1 + rem1);
    }

    function testFuzz_Swap(
        uint256 liquidity0,
        uint256 liquidity1,
        uint256 swapAmount,
        bool swapToken0
    ) public {
        liquidity0 = bound(liquidity0, 10_000 * 10 ** 18, 1_000_000 * 10 ** 18);
        liquidity1 = bound(liquidity1, 10_000 * 10 ** 18, 1_000_000 * 10 ** 18);

        vm.startPrank(alice);
        t0.approve(address(dex), liquidity0);
        t1.approve(address(dex), liquidity1);
        dex.AddLiquidity(liquidity0, liquidity1);
        vm.stopPrank();
        // only swap 10% of the reserves to avoid excessive slippage
        uint256 max = swapToken0 ? liquidity0 / 10 : liquidity1 / 10;
        uint256 maxSwap = swapToken0
            ? (type(uint256).max / dex.reserve1()) - dex.reserve0()
            : (type(uint256).max / dex.reserve0()) - dex.reserve1();
        swapAmount = bound(swapAmount, 1e18, min(max, maxSwap));

        address tokenIn = swapToken0 ? address(t0) : address(t1);

        IERC20 tIn = swapToken0 ? t0 : t1;
        IERC20 tOut = swapToken0 ? t1 : t0;

        uint256 balBefore = tOut.balanceOf(bob);
        uint256 kBefore = dex.reserve0() * dex.reserve1();

        vm.startPrank(bob);
        tIn.approve(address(dex), swapAmount);
        uint256 amountOut = dex.swap(tokenIn, swapAmount, 0);
        vm.stopPrank();

        uint256 balAfter = tOut.balanceOf(bob);
        uint256 kAfter = dex.reserve0() * dex.reserve1();

        assertGt(amountOut, 0);
        assertEq(balAfter - balBefore, amountOut);

        // kAfter should be higher because of removal
        assertGe(kAfter, kBefore);
    }

    function testFuzz_MultipleSwaps(
        uint256 liq0,
        uint256 liq1,
        uint8 numSwaps
    ) public {
        liq0 = bound(liq0, 10_000 * 10 ** 18, 100_000 * 10 ** 18);
        liq1 = bound(liq1, 10_000 * 10 ** 18, 100_000 * 10 ** 18);
        numSwaps = uint8(bound(numSwaps, 1, 10));

        vm.startPrank(alice);
        t0.approve(address(dex), liq0);
        t1.approve(address(dex), liq1);
        dex.AddLiquidity(liq0, liq1);
        vm.stopPrank();

        uint256 kInitial = dex.reserve0() * dex.reserve1();

        for (uint256 i = 0; i < numSwaps; i++) {
            bool swapToken0 = i % 2 == 0;
            uint256 maxSwap = swapToken0
                ? dex.reserve0() / 20
                : dex.reserve1() / 20;
            uint256 swapAmount = (maxSwap * (i + 1)) / (numSwaps + 1);

            vm.startPrank(bob);
            if (swapToken0) {
                t0.approve(address(dex), swapAmount);
                dex.swap(address(t0), swapAmount, 0);
            } else {
                t1.approve(address(dex), swapAmount);
                dex.swap(address(t1), swapAmount, 0);
            }
            vm.stopPrank();
        }
        uint256 kFinal = dex.reserve0() * dex.reserve1();

        assertGt(kFinal, kInitial);
    }

    function testFuzz_PriceConsistency(uint256 a0, uint256 a1) public {
        a0 = bound(a0, 1 * 10 ** 18, 1_000_000 * 10 ** 18);
        a1 = bound(a1, 1 * 10 ** 18, 1_000_000 * 10 ** 18);

        vm.startPrank(alice);
        t0.approve(address(dex), a0);
        t1.approve(address(dex), a1);
        dex.AddLiquidity(a0, a1);
        vm.stopPrank();

        uint256 price0 = dex.getPrice0();
        uint256 price1 = dex.getPrice1();

        uint256 product = (price0 * price1) / 1e18;
        assertApproxEqRel(product, 1e18, 0.01e18);
    }

    function testFuzz_GetAmountOut(
        uint256 liq0,
        uint256 liq1,
        uint256 swapAmount,
        bool swapToken0
    ) public {
        liq0 = bound(liq0, 10_000 * 10 ** 18, 100_000 * 10 ** 18);
        liq1 = bound(liq1, 10_000 * 10 ** 18, 100_000 * 10 ** 18);

        vm.startPrank(alice);
        t0.approve(address(dex), liq0);
        t1.approve(address(dex), liq1);
        dex.AddLiquidity(liq0, liq1);
        vm.stopPrank();

        uint256 max = swapToken0 ? liq0 / 10 : liq1 / 10;
        swapAmount = bound(swapAmount, 1 * 10 ** 18, max);

        address tIn = swapToken0 ? address(t0) : address(t1);

        uint256 quote = dex.getAmountOut(tIn, swapAmount);

        vm.startPrank(bob);
        IERC20(tIn).approve(address(dex), swapAmount);
        uint256 actualOut = dex.swap(tIn, swapAmount, 0);
        vm.stopPrank();

        assertEq(quote, actualOut);
    }

    function testFuzz_ConstantProductNeverDecreases(
        uint256 liquidity0,
        uint256 liquidity1,
        uint256 swap1,
        uint256 swap2
    ) public {
        liquidity0 = bound(liquidity0, 10_000 * 10 ** 18, 100_000 * 10 ** 18);
        liquidity1 = bound(liquidity1, 10_000 * 10 ** 18, 100_000 * 10 ** 18);
        swap1 = bound(swap1, 1 * 10 ** 18, liquidity0 / 20);
        swap2 = bound(swap2, 1 * 10 ** 18, liquidity1 / 20);

        vm.startPrank(alice);
        t0.approve(address(dex), liquidity0);
        t1.approve(address(dex), liquidity1);
        dex.AddLiquidity(liquidity0, liquidity1);
        vm.stopPrank();

        uint256 k0 = dex.reserve0() * dex.reserve1();

        vm.startPrank(bob);
        t0.approve(address(dex), swap1);
        dex.swap(address(t0), swap1, 0);
        vm.stopPrank();

        uint256 k1 = dex.reserve0() * dex.reserve1();
        assertGe(k1, k0);

        vm.startPrank(bob);
        t1.approve(address(dex), swap2);
        dex.swap(address(t1), swap2, 0);
        vm.stopPrank();

        uint256 k2 = dex.reserve0() * dex.reserve1();
        assertGe(k2, k1);
    }

    function sqrt(uint256 y) internal pure returns (uint256 z) {
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
