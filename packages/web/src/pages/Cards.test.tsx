import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";

// ── 可变 mock 状态 ───────────────────────────────────────────────
let cardsData: Array<{ tokenId: bigint; dataAmount: bigint; createdAt: bigint; isDestroyed: boolean }> = [];
let cardsLoading = false;
let cardsError = false;
const refetchFn = vi.fn();

// redeemForSim mock：调用即捕获 tokenIds，并切到 success（同步）以驱动后续 claim。
const redeemFn = vi.fn();
let redeemSuccess = false;
const redeemReset = vi.fn();
const redeemHash = "0xtxhash";

// simApi.claim mock（经 hook）
const claimMutate = vi.fn();
let claimPending = false;

// 我的 SIM 列表
let simsData: Array<{
  id: string;
  days: number;
  destination: string;
  recipient: string;
  addressLine: string;
  status: "pending" | "confirmed";
  createdAt: string;
}> = [];
let simsLoading = false;
let simsError = false;

vi.mock("wagmi", () => ({
  useAccount: () => ({ address: "0xabc" as `0x${string}`, chainId: 31337 }),
  http: () => ({}),
}));

vi.mock("@/hooks/contracts", () => ({
  useTrafficCards: () => ({
    cards: cardsData,
    isLoading: cardsLoading,
    isError: cardsError,
    error: cardsError ? new Error("RPC rate limited") : null,
    refetch: refetchFn,
  }),
  useRedeemForSim: () => ({
    redeem: redeemFn,
    hash: redeemSuccess ? redeemHash : undefined,
    isPending: false,
    isConfirming: false,
    isSuccess: redeemSuccess,
    error: null,
    reset: redeemReset,
  }),
}));

vi.mock("@/hooks/useSim", () => ({
  useMySims: () => ({
    data: simsData,
    isLoading: simsLoading,
    isError: simsError,
    refetch: vi.fn(),
  }),
  useClaimSim: () => ({
    mutate: claimMutate,
    isPending: claimPending,
  }),
}));

import Cards from "./Cards";

function renderCards() {
  const qc = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return render(
    <QueryClientProvider client={qc}>
      <Cards />
    </QueryClientProvider>,
  );
}

function resetState() {
  cardsData = [];
  cardsLoading = false;
  cardsError = false;
  refetchFn.mockReset();
  redeemFn.mockReset();
  redeemReset.mockReset();
  redeemSuccess = false;
  claimMutate.mockReset();
  claimPending = false;
  simsData = [];
  simsLoading = false;
  simsError = false;
}

beforeEach(resetState);

describe("Cards 页 — 流量卡多选销毁领 SIM", () => {
  it("CARD-01: 流量卡 tab 有充值即发说明 + 无单张销毁按钮", () => {
    renderCards();
    expect(screen.getByText(/充值即得流量卡/)).toBeInTheDocument();
    expect(screen.getByText(/流量卡不可转卖/)).toBeInTheDocument();
    // 旧单张「销毁激活」按钮已移除
    expect(screen.queryByText(/销毁激活/)).not.toBeInTheDocument();
  });

  it("CARD-02: isError → 显「加载失败·重试」非空态；refetch 可触发", () => {
    cardsError = true;
    renderCards();
    expect(screen.getByText(/流量卡加载失败/)).toBeInTheDocument();
    expect(screen.queryByText(/暂无流量卡/)).not.toBeInTheDocument();
    fireEvent.click(screen.getByText("重试"));
    expect(refetchFn).toHaveBeenCalled();
  });

  it("CARD-02: 无卡且非 error → 空态「暂无流量卡」", () => {
    renderCards();
    expect(screen.getByText(/暂无流量卡/)).toBeInTheDocument();
    expect(screen.queryByText(/流量卡加载失败/)).not.toBeInTheDocument();
  });

  it("CARD-02: 有卡 → 渲染读链数据（tokenId / 不可转卖）", () => {
    cardsData = [
      { tokenId: 7n, dataAmount: 1024n * 1024n * 500n, createdAt: 1_700_000_000n, isDestroyed: false },
    ];
    renderCards();
    expect(screen.getByText(/流量卡 #7/)).toBeInTheDocument();
    expect(screen.getByText(/发放于 .* · 不可转卖/)).toBeInTheDocument();
  });

  it("CARD-05: dataAmount = type(uint256).max → 显「无限流量 · 1 天」（非天文数字）", () => {
    const UINT256_MAX = 2n ** 256n - 1n;
    cardsData = [
      { tokenId: 9n, dataAmount: UINT256_MAX, createdAt: 1_700_000_000n, isDestroyed: false },
    ];
    const { container } = renderCards();
    expect(screen.getByText("无限流量")).toBeInTheDocument();
    expect(screen.getByText("1 天")).toBeInTheDocument();
    expect(container.textContent).not.toContain(UINT256_MAX.toString());
  });

  it("CARD-06: 多选 → 底部操作条显示正确张数/天数", () => {
    cardsData = [
      { tokenId: 1n, dataAmount: 100n, createdAt: 1n, isDestroyed: false },
      { tokenId: 2n, dataAmount: 100n, createdAt: 2n, isDestroyed: false },
      { tokenId: 3n, dataAmount: 100n, createdAt: 3n, isDestroyed: false },
    ];
    renderCards();
    // 初始无操作条
    expect(screen.queryByText(/并领取 SIM/)).not.toBeInTheDocument();
    // 选 2 张
    fireEvent.click(screen.getByText(/流量卡 #1/));
    fireEvent.click(screen.getByText(/流量卡 #3/));
    expect(screen.getByText(/销毁 2 张并领取 SIM（2 天）/)).toBeInTheDocument();
    // 取消一张 → 1 张
    fireEvent.click(screen.getByText(/流量卡 #1/));
    expect(screen.getByText(/销毁 1 张并领取 SIM（1 天）/)).toBeInTheDocument();
  });

  it("CARD-07: 弹窗提交 → 调 redeemForSim(选中 tokenIds)", async () => {
    cardsData = [
      { tokenId: 5n, dataAmount: 100n, createdAt: 1n, isDestroyed: false },
      { tokenId: 8n, dataAmount: 100n, createdAt: 2n, isDestroyed: false },
    ];
    renderCards();
    fireEvent.click(screen.getByText(/流量卡 #5/));
    fireEvent.click(screen.getByText(/流量卡 #8/));
    // 打开弹窗
    fireEvent.click(screen.getByText(/销毁 2 张并领取 SIM/));
    fireEvent.change(screen.getByPlaceholderText("收件人姓名"), { target: { value: "Alice" } });
    fireEvent.change(screen.getByPlaceholderText("详细收件地址"), {
      target: { value: "123 Main St" },
    });
    fireEvent.click(screen.getByText(/销毁并领取 SIM/));

    await waitFor(() => expect(redeemFn).toHaveBeenCalledTimes(1));
    const arg = redeemFn.mock.calls[0][0];
    expect(arg.map((b: bigint) => Number(b)).sort()).toEqual([5, 8]);
  });

  it("CARD-08: 链上销毁成功 → 调 simApi.claim（含 tokenIds/txHash/收件信息）", async () => {
    cardsData = [{ tokenId: 5n, dataAmount: 100n, createdAt: 1n, isDestroyed: false }];
    const qc = new QueryClient({ defaultOptions: { queries: { retry: false } } });
    const { rerender } = render(
      <QueryClientProvider client={qc}>
        <Cards />
      </QueryClientProvider>,
    );
    // 真实流程：先选卡 + 填表 + 提交（此时链上尚未成功）
    fireEvent.click(screen.getByText(/流量卡 #5/));
    fireEvent.click(screen.getByText(/销毁 1 张并领取 SIM/));
    fireEvent.change(screen.getByPlaceholderText("收件人姓名"), { target: { value: "Bob" } });
    fireEvent.change(screen.getByPlaceholderText("详细收件地址"), {
      target: { value: "9 Park Ave" },
    });
    fireEvent.click(screen.getByText(/销毁并领取 SIM/));
    expect(redeemFn).toHaveBeenCalled();

    // 模拟链上交易确认成功 → 重渲染触发 claim effect
    redeemSuccess = true;
    rerender(
      <QueryClientProvider client={qc}>
        <Cards />
      </QueryClientProvider>,
    );

    await waitFor(() => expect(claimMutate).toHaveBeenCalled());
    const payload = claimMutate.mock.calls[0][0];
    expect(payload.tokenIds).toEqual([5]);
    expect(payload.txHash).toBe("0xtxhash");
    expect(payload.recipient).toBe("Bob");
    expect(payload.addressLine).toBe("9 Park Ave");
  });

  it("CARD-09: 我的 SIM tab — 列表展示天数/目的地/状态", () => {
    simsData = [
      {
        id: "s1",
        days: 3,
        destination: "JP",
        recipient: "Alice",
        addressLine: "123 Main St",
        status: "pending",
        createdAt: "2026-01-01",
      },
    ];
    renderCards();
    fireEvent.click(screen.getByText("我的 SIM"));
    expect(screen.getByText("3 天")).toBeInTheDocument();
    expect(screen.getByText(/日本/)).toBeInTheDocument();
    expect(screen.getByText("处理中")).toBeInTheDocument();
    expect(screen.getByText(/Alice/)).toBeInTheDocument();
  });

  it("CARD-09: 我的 SIM 空态 → 「销毁流量卡即可领取 SIM」", () => {
    renderCards();
    fireEvent.click(screen.getByText("我的 SIM"));
    expect(screen.getByText(/销毁流量卡即可领取 SIM/)).toBeInTheDocument();
  });
});
