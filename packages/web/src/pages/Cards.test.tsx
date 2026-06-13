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
    expect(screen.getByText(/Get Traffic Cards Instantly/)).toBeInTheDocument();
    expect(screen.getByText(/non-transferable/)).toBeInTheDocument();
    // 旧单张「销毁激活」按钮已移除
    expect(screen.queryByText(/Burn to activate/)).not.toBeInTheDocument();
  });

  it("CARD-02: isError → 显「加载失败·重试」非空态；refetch 可触发", () => {
    cardsError = true;
    renderCards();
    expect(screen.getByText(/Failed to load traffic cards/)).toBeInTheDocument();
    expect(screen.queryByText(/No traffic cards yet/)).not.toBeInTheDocument();
    fireEvent.click(screen.getByText("Retry"));
    expect(refetchFn).toHaveBeenCalled();
  });

  it("CARD-02: 无卡且非 error → 空态「暂无流量卡」", () => {
    renderCards();
    expect(screen.getByText(/No traffic cards yet/)).toBeInTheDocument();
    expect(screen.queryByText(/Failed to load traffic cards/)).not.toBeInTheDocument();
  });

  it("CARD-02: 有卡 → 渲染读链数据（tokenId / 不可转卖）", () => {
    cardsData = [
      { tokenId: 7n, dataAmount: 1024n * 1024n * 500n, createdAt: 1_700_000_000n, isDestroyed: false },
    ];
    renderCards();
    expect(screen.getByText(/Traffic Card #7/)).toBeInTheDocument();
    expect(screen.getByText(/Issued on .* · non-transferable/)).toBeInTheDocument();
  });

  it("CARD-05: dataAmount = type(uint256).max → 显「无限流量 · 1 天」（非天文数字）", () => {
    const UINT256_MAX = 2n ** 256n - 1n;
    cardsData = [
      { tokenId: 9n, dataAmount: UINT256_MAX, createdAt: 1_700_000_000n, isDestroyed: false },
    ];
    const { container } = renderCards();
    expect(screen.getByText("Unlimited Data")).toBeInTheDocument();
    expect(screen.getByText("1 day")).toBeInTheDocument();
    expect(container.textContent).not.toContain(UINT256_MAX.toString());
  });

  it("CARD-06: < 3 张时禁用并提示门槛，选满 3 张后可点", () => {
    cardsData = [
      { tokenId: 1n, dataAmount: 100n, createdAt: 1n, isDestroyed: false },
      { tokenId: 2n, dataAmount: 100n, createdAt: 2n, isDestroyed: false },
      { tokenId: 3n, dataAmount: 100n, createdAt: 3n, isDestroyed: false },
    ];
    renderCards();
    // 初始无操作条
    expect(screen.queryByText(/and redeem SIM/)).not.toBeInTheDocument();
    // 选 2 张 → 出现操作条但禁用 + 门槛提示，不出现可点的「Burn N and redeem SIM」
    fireEvent.click(screen.getByText(/Traffic Card #1/));
    fireEvent.click(screen.getByText(/Traffic Card #3/));
    expect(screen.queryByText(/Burn \d+ and redeem SIM/)).not.toBeInTheDocument();
    const hintBtn = screen.getByText(/at least 3/i).closest("button")!;
    expect(hintBtn).toBeDisabled();
    // 选满 3 张 → 可点，文案切回 Burn 3 and redeem SIM (3 days)
    fireEvent.click(screen.getByText(/Traffic Card #2/));
    const burnBtn = screen.getByText(/Burn 3 and redeem SIM \(3 days\)/).closest("button")!;
    expect(burnBtn).toBeEnabled();
  });

  it("CARD-07: 选满 3 张 → 弹窗提交调 redeemForSim(选中 tokenIds)", async () => {
    cardsData = [
      { tokenId: 5n, dataAmount: 100n, createdAt: 1n, isDestroyed: false },
      { tokenId: 8n, dataAmount: 100n, createdAt: 2n, isDestroyed: false },
      { tokenId: 11n, dataAmount: 100n, createdAt: 3n, isDestroyed: false },
    ];
    renderCards();
    fireEvent.click(screen.getByText(/Traffic Card #5/));
    fireEvent.click(screen.getByText(/Traffic Card #8/));
    fireEvent.click(screen.getByText(/Traffic Card #11/));
    // 打开弹窗
    fireEvent.click(screen.getByText(/Burn 3 and redeem SIM/));
    fireEvent.change(screen.getByPlaceholderText("Recipient name"), { target: { value: "Alice" } });
    fireEvent.change(screen.getByPlaceholderText("Full shipping address"), {
      target: { value: "123 Main St" },
    });
    fireEvent.click(screen.getByText(/Burn and redeem SIM/));

    await waitFor(() => expect(redeemFn).toHaveBeenCalledTimes(1));
    const arg = redeemFn.mock.calls[0][0];
    expect(arg.map((b: bigint) => Number(b)).sort((a: number, b: number) => a - b)).toEqual([5, 8, 11]);
  });

  it("CARD-08: 链上销毁成功 → 调 simApi.claim（含 tokenIds/txHash/收件信息）", async () => {
    cardsData = [
      { tokenId: 5n, dataAmount: 100n, createdAt: 1n, isDestroyed: false },
      { tokenId: 8n, dataAmount: 100n, createdAt: 2n, isDestroyed: false },
      { tokenId: 11n, dataAmount: 100n, createdAt: 3n, isDestroyed: false },
    ];
    const qc = new QueryClient({ defaultOptions: { queries: { retry: false } } });
    const { rerender } = render(
      <QueryClientProvider client={qc}>
        <Cards />
      </QueryClientProvider>,
    );
    // 真实流程：先选 3 张卡 + 填表 + 提交（此时链上尚未成功）
    fireEvent.click(screen.getByText(/Traffic Card #5/));
    fireEvent.click(screen.getByText(/Traffic Card #8/));
    fireEvent.click(screen.getByText(/Traffic Card #11/));
    fireEvent.click(screen.getByText(/Burn 3 and redeem SIM/));
    fireEvent.change(screen.getByPlaceholderText("Recipient name"), { target: { value: "Bob" } });
    fireEvent.change(screen.getByPlaceholderText("Full shipping address"), {
      target: { value: "9 Park Ave" },
    });
    fireEvent.click(screen.getByText(/Burn and redeem SIM/));
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
    expect([...payload.tokenIds].sort((a: number, b: number) => a - b)).toEqual([5, 8, 11]);
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
    fireEvent.click(screen.getByText("My SIM"));
    expect(screen.getByText("3 days")).toBeInTheDocument();
    expect(screen.getByText(/Japan/)).toBeInTheDocument();
    expect(screen.getByText("Processing")).toBeInTheDocument();
    expect(screen.getByText(/Alice/)).toBeInTheDocument();
  });

  it("CARD-09: 我的 SIM 空态 → 「销毁流量卡即可领取 SIM」", () => {
    renderCards();
    fireEvent.click(screen.getByText("My SIM"));
    expect(screen.getByText(/Burn traffic cards to redeem a SIM/)).toBeInTheDocument();
  });
});
