import { describe, it, expect, vi, beforeEach } from "vitest";

// mock wagmi —— useApprove 内部走 useWriteContract / useWaitForTransactionReceipt / useChainId。
// 验证 approve 调用 args = exact amount（禁 infinite/MaxUint256）。
const writeContract = vi.fn();
vi.mock("wagmi", () => ({
  useChainId: () => 31337,
  useWriteContract: () => ({
    writeContract,
    data: undefined,
    isPending: false,
    error: null,
  }),
  useWaitForTransactionReceipt: () => ({ isPending: false, isSuccess: false }),
  useReadContract: vi.fn(),
}));

// mock contracts 地址出口（避免依赖 deployments 真地址）。
vi.mock("../../config/contracts", () => ({
  getUsdt: () => "0xUSDT" as `0x${string}`,
}));

import { renderHook } from "@testing-library/react";
import { useApprove } from "./useUsdtContract";

describe("useApprove —— exact approve（TSA-04 单元）", () => {
  beforeEach(() => writeContract.mockReset());

  it("approve(spender, amount) 透传 exact amount 作为第二参，绝不传 MaxUint256", () => {
    const { result } = renderHook(() => useApprove());
    const spender = "0xSPENDER" as `0x${string}`;
    const amount = 1_500_000n; // 1.5 USDT(6 位)

    result.current.approve(spender, amount);

    expect(writeContract).toHaveBeenCalledTimes(1);
    const callArg = writeContract.mock.calls[0][0];
    expect(callArg.functionName).toBe("approve");
    expect(callArg.args).toEqual([spender, amount]);

    // 反向断言：不是 infinite（MaxUint256）。
    const MAX_UINT256 = 2n ** 256n - 1n;
    expect(callArg.args[1]).not.toBe(MAX_UINT256);
    expect(callArg.args[1]).toBe(amount);
  });
});
