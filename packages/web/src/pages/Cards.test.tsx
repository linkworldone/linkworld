import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, fireEvent } from "@testing-library/react";

// ── 可变 mock 状态 ───────────────────────────────────────────────
let cardsData: Array<{ tokenId: bigint; dataAmount: bigint; createdAt: bigint; isDestroyed: boolean }> = [];
let cardsLoading = false;
let cardsError = false;
const refetchFn = vi.fn();
const burnFn = vi.fn();

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
  useBurnCard: () => ({
    burnCard: burnFn,
    isPending: false,
    isConfirming: false,
    isSuccess: false,
    error: null,
    reset: vi.fn(),
  }),
}));

// pendingSync 真实实现走 localStorage（jsdom 提供），CARD-03 直接断言写入。
import Cards from "./Cards";
import { getPendingSync, clearPendingSync } from "@/utils/pendingSync";

function resetState() {
  cardsData = [];
  cardsLoading = false;
  cardsError = false;
  refetchFn.mockReset();
  burnFn.mockReset();
  clearPendingSync("sim_claim");
}

beforeEach(resetState);

describe("Cards 页 双 Tab（T8）", () => {
  it("CARD-01: Tab1 无 Admin 发卡按钮 + 有自动发放说明", () => {
    render(<Cards />);
    // 移除 onlyOracle 的「Issue Monthly Card」发卡按钮
    expect(screen.queryByText(/Issue Monthly Card/i)).not.toBeInTheDocument();
    expect(screen.queryByText(/发卡/)).not.toBeInTheDocument();
    // 自动发放说明卡
    expect(screen.getByText(/流量卡自动发放/)).toBeInTheDocument();
    expect(screen.getByText(/锁仓满 1 个月后，系统将按你的保证金自动发放/)).toBeInTheDocument();
    expect(screen.getByText(/流量卡不可转卖/)).toBeInTheDocument();
  });

  it("CARD-02: isError → 显「加载失败·重试」非空态；refetch 可触发", () => {
    cardsError = true;
    render(<Cards />);
    expect(screen.getByText(/流量卡加载失败/)).toBeInTheDocument();
    // 加载失败不得冒充空态
    expect(screen.queryByText(/暂无流量卡/)).not.toBeInTheDocument();
    fireEvent.click(screen.getByText("重试"));
    expect(refetchFn).toHaveBeenCalled();
  });

  it("CARD-02: 无卡且非 error → 空态「暂无流量卡」", () => {
    cardsData = [];
    cardsError = false;
    render(<Cards />);
    expect(screen.getByText(/暂无流量卡/)).toBeInTheDocument();
    expect(screen.queryByText(/流量卡加载失败/)).not.toBeInTheDocument();
  });

  it("CARD-02: 有卡 → 渲染读链数据（tokenId / 不可转卖）", () => {
    cardsData = [
      { tokenId: 7n, dataAmount: 1024n * 1024n * 500n, createdAt: 1_700_000_000n, isDestroyed: false },
    ];
    render(<Cards />);
    expect(screen.getByText(/流量卡 #7/)).toBeInTheDocument();
    expect(screen.getByText(/发放于 .* · 不可转卖/)).toBeInTheDocument();
    // 卡列表区脚注「销毁后 30 天」
    expect(screen.getAllByText(/销毁后剩余流量额度有效期为 30 天/).length).toBeGreaterThan(0);
  });

  it("CARD-03: SIM 表单提交 → 写 pendingSync + 成功提示", () => {
    render(<Cards />);
    fireEvent.click(screen.getByText("SIM 领取"));
    fireEvent.change(screen.getByPlaceholderText("收件人姓名"), { target: { value: "Alice" } });
    fireEvent.change(screen.getByPlaceholderText("详细收件地址"), {
      target: { value: "123 Main St" },
    });
    fireEvent.click(screen.getByText("提交领取申请"));

    // 写入 pendingSync
    const pending = getPendingSync<{ recipient: string; addressLine: string }>("sim_claim");
    expect(pending).not.toBeNull();
    expect(pending!.data.recipient).toBe("Alice");
    expect(pending!.data.addressLine).toBe("123 Main St");
    // 内联成功提示
    expect(screen.getByText(/领取申请已提交/)).toBeInTheDocument();
  });

  it("CARD-03: 收件信息为空 → 提交按钮禁用", () => {
    render(<Cards />);
    fireEvent.click(screen.getByText("SIM 领取"));
    expect(screen.getByText("提交领取申请").closest("button")).toBeDisabled();
  });

  it("CARD-04: 双 Tab 切换渲染", () => {
    render(<Cards />);
    // 默认 Tab1 流量卡
    expect(screen.getByText(/流量卡自动发放/)).toBeInTheDocument();
    // 切到 Tab2 SIM
    fireEvent.click(screen.getByText("SIM 领取"));
    expect(screen.getByText(/领取实体 SIM/)).toBeInTheDocument();
    expect(screen.getByText(/全球通 eSIM 即将推出/)).toBeInTheDocument();
  });
});
