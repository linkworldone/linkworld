import { describe, it, expect, vi, beforeEach } from "vitest";

// mock wagmi useReadContract —— 受控返回 getFeeRate / calculateFee 读链结果。
let readResult: { data?: bigint; isLoading: boolean; isError: boolean } = {
  data: undefined,
  isLoading: false,
  isError: false,
};
const readSpy = vi.fn();
vi.mock("wagmi", () => ({
  useChainId: () => 31337,
  useReadContract: (cfg: unknown) => {
    readSpy(cfg);
    return readResult;
  },
}));

vi.mock("../../config/contracts", () => ({
  getContractAddress: () => "0xFEE" as `0x${string}`,
}));

import { renderHook } from "@testing-library/react";
import { useFeeRate, useCalculateFee } from "./useFeeManager";

beforeEach(() => {
  readResult = { data: undefined, isLoading: false, isError: false };
  readSpy.mockReset();
});

describe("useFeeRate —— 读链费率基点（FEE-01）", () => {
  it("getFeeRate 返回基点 150n → percent 1.5、label '1.5%'（/10000）", () => {
    readResult = { data: 150n, isLoading: false, isError: false };
    const { result } = renderHook(() => useFeeRate());
    expect(result.current.bps).toBe(150n);
    expect(result.current.percent).toBe(1.5);
    expect(result.current.label).toBe("1.5%");
    expect(result.current.isError).toBe(false);

    // 调的是 getFeeRate（读链，不自算）。
    const cfg = readSpy.mock.calls[0][0] as { functionName: string };
    expect(cfg.functionName).toBe("getFeeRate");
  });

  it("250n 基点 → 2.5%（整数）", () => {
    readResult = { data: 250n, isLoading: false, isError: false };
    const { result } = renderHook(() => useFeeRate());
    expect(result.current.label).toBe("2.5%");
  });

  it("FEE-02: 读链失败 → label undefined（调用方兜底 '--'，不写死）", () => {
    readResult = { data: undefined, isLoading: false, isError: true };
    const { result } = renderHook(() => useFeeRate());
    expect(result.current.label).toBeUndefined();
    expect(result.current.percent).toBeUndefined();
    expect(result.current.isError).toBe(true);
  });

  it("读链 loading → isLoading=true、label undefined（skeleton 兜底）", () => {
    readResult = { data: undefined, isLoading: true, isError: false };
    const { result } = renderHook(() => useFeeRate());
    expect(result.current.isLoading).toBe(true);
    expect(result.current.label).toBeUndefined();
  });
});

describe("useCalculateFee —— 读链精确费额（FEE-01 不自算）", () => {
  it("calculateFee(amount) 读链返回费额，args=[amount]，不前端算 amount*rate", () => {
    readResult = { data: 1_500_000n, isLoading: false, isError: false };
    const amount = 100_000_000n; // 100 USDT(6 位)
    const { result } = renderHook(() => useCalculateFee(amount));
    expect(result.current.fee).toBe(1_500_000n);

    const cfg = readSpy.mock.calls[0][0] as { functionName: string; args?: unknown[] };
    expect(cfg.functionName).toBe("calculateFee");
    expect(cfg.args).toEqual([amount]);
  });

  it("amount<=0 → 不读链（enabled=false，fee undefined）", () => {
    readResult = { data: undefined, isLoading: false, isError: false };
    const { result } = renderHook(() => useCalculateFee(0n));
    expect(result.current.fee).toBeUndefined();
    const cfg = readSpy.mock.calls[0][0] as { args?: unknown[] };
    expect(cfg.args).toBeUndefined();
  });
});
