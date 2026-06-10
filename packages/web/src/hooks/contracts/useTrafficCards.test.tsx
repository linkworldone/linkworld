import { describe, it, expect, vi, beforeEach } from "vitest";

// LOG-01：getLogs 失败 → error 态（非静默置空）。区分「加载失败」vs「真无卡」。
// 现状 bug：catch → setTokenIds([]) 静默吞错，会把 RPC 失败误显示为「没有卡」。

let getLogsImpl: () => Promise<unknown[]> = async () => [];

// 稳定引用（真实 wagmi 用 useMemo 稳定返回；mock 也须稳定，否则 useEffect 依赖每次变 → 死循环）。
const stablePublicClient = {
  getLogs: () => getLogsImpl(),
  getBlockNumber: async () => 1000n,
};
const stableReadContracts = { data: [], isLoading: false, isError: false, error: null, refetch: () => {} };

vi.mock("wagmi", () => ({
  useChainId: () => 31337,
  usePublicClient: () => stablePublicClient,
  useReadContracts: () => stableReadContracts,
  useReadContract: () => ({ data: undefined, isLoading: false, refetch: () => {} }),
  useWriteContract: () => ({ writeContract: vi.fn(), data: undefined, isPending: false, error: null, reset: vi.fn() }),
  useWaitForTransactionReceipt: () => ({ isPending: false, isSuccess: false }),
}));

vi.mock("../../config/contracts", () => ({
  getContractAddress: () => "0xNFT" as `0x${string}`,
}));

import { renderHook, waitFor } from "@testing-library/react";
import { useTrafficCards } from "./useTrafficCard";

const USER = "0xabc" as `0x${string}`;

describe("useTrafficCards getLogs error 态（LOG-01）", () => {
  beforeEach(() => {
    getLogsImpl = async () => [];
  });

  it("getLogs 失败 → isError=true（不静默置空，区分加载失败 vs 真无卡）", async () => {
    getLogsImpl = async () => {
      throw new Error("RPC rate limited");
    };
    const { result } = renderHook(() => useTrafficCards(USER));
    await waitFor(() => expect(result.current.isError).toBe(true));
    // 加载失败时不可冒充「真无卡」：cards 为空但 isError 标记区分语义。
    expect(result.current.isError).toBe(true);
  });

  it("getLogs 成功且无日志 → isError=false（真无卡，非错误）", async () => {
    getLogsImpl = async () => [];
    const { result } = renderHook(() => useTrafficCards(USER));
    await waitFor(() => expect(result.current.isLoading).toBe(false));
    expect(result.current.isError).toBe(false);
    expect(result.current.cards).toEqual([]);
  });
});
