import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, fireEvent } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import type { Operator, Region } from "@/types";

// ── 可变 mock 状态 ───────────────────────────────────────────────
let operators: Operator[] = [];
let regions: Region[] = [];
const applyNumberFn = vi.fn();
const invalidateFn = vi.fn();
let txStatus = "idle";
let txError: string | undefined;

function makeOperator(over: Partial<Operator> = {}): Operator {
  return {
    id: "1",
    name: "Acme Telecom",
    region: "United States",
    // 6 位最小单位 bigint（T2）：50 USDT = 50_000_000；不应再 ×10^6 当成 18 位/美元。
    requiredDeposit: 50_000_000n,
    dataRate: 0.05,
    callRate: 0.02,
    isActive: true,
    ...over,
  };
}

vi.mock("wagmi", () => ({
  useAccount: () => ({ address: "0xabc" as `0x${string}`, chainId: 31337 }),
  http: () => ({}),
}));

vi.mock("@/hooks/useOperator", () => ({
  useRegions: () => ({ data: regions }),
  useOperatorsByRegion: () => ({ data: operators }),
  useApplyNumber: () => ({
    applyNumber: applyNumberFn,
    txState: { status: txStatus, error: txError },
    invalidate: invalidateFn,
  }),
}));

vi.mock("@/services/api/client", () => ({
  apiClient: { post: vi.fn().mockResolvedValue({ virtual_number: "+1555", password: "pw" }) },
}));

// FeeBreakdown 读链；页面测隔离成静态行（fee 读链测在 FeeBreakdown.test）。
vi.mock("@/components/shared/FeeBreakdown", () => ({
  FeeBreakdown: ({ amount }: { amount: bigint }) => (
    <div data-testid="fee-breakdown" data-amount={String(amount)}>
      平台手续费 (1.5%)
    </div>
  ),
}));

import RegionDetail from "./RegionDetail";

function renderPage() {
  return render(
    <MemoryRouter initialEntries={["/region/US"]}>
      <RegionDetail />
    </MemoryRouter>,
  );
}

beforeEach(() => {
  operators = [];
  regions = [{ code: "US", name: "United States", flag: "🇺🇸", operatorCount: 1, startingPrice: 50 }];
  txStatus = "idle";
  txError = undefined;
  applyNumberFn.mockReset();
  invalidateFn.mockReset();
});

describe("RegionDetail 页（T10）", () => {
  it("REG-01: requiredDeposit 50_000_000(6 位最小单位) → 展示 50.00 USDT，不二次缩放/不当 18 位美元", () => {
    operators = [makeOperator()];
    const { container } = renderPage();
    // 正确：50.00（formatAmount 6 位）。
    expect(screen.getAllByText("50.00").length).toBeGreaterThanOrEqual(1);
    // 反向：绝不出现二次缩放（×10^6）或当 18 位（极小数）的文本。
    expect(container.textContent).not.toContain("50000000");
    expect(container.textContent).not.toContain("0.00000000005");
  });

  it("REG-02: 申请弹层展示费用明细（FeeBreakdown 读链 calculateFee(押金本金) + 押金本金）", () => {
    operators = [makeOperator()];
    renderPage();
    // 打开申请弹层
    fireEvent.click(screen.getByText("Apply for Number"));
    // FeeBreakdown 渲染，amount = 押金本金（读链 calculateFee 的入参，非自算）。
    const fee = screen.getByTestId("fee-breakdown");
    expect(fee).toBeInTheDocument();
    expect(fee.getAttribute("data-amount")).toBe("50000000");
    // 弹层内押金本金仍展示 50.00 USDT（AmountDisplay）。
    expect(screen.getAllByText("50.00").length).toBeGreaterThanOrEqual(1);
    // 身份签名提示（§3.7，走 signedPost 意向，无臆造链上付费步骤）。
    expect(screen.getByText(/不消耗 gas/)).toBeInTheDocument();
  });

  it("REG-03: 深蓝金换肤——卡用 ui/card(暖米白+金线)，无旧 token，无装饰 emoji 图标", () => {
    operators = [makeOperator()];
    const { container } = renderPage();
    // ui/card 原子（暖米白 + 金线描边）取代手写 bg-surface-card 卡。
    expect(container.querySelector('[data-slot="card"]')).not.toBeNull();
    // 国旗 emoji 是数据（国家标识），保留；但不得出现旧暗色 token / 旧手写卡类名。
    const html = container.innerHTML;
    expect(html).not.toContain("text-text-muted");
    expect(html).not.toContain("text-text-secondary");
    expect(html).not.toContain("bg-surface-secondary");
    // lucide 图标（svg）已注入（Data/Call/Deposit 行图标）。
    expect(container.querySelectorAll("svg").length).toBeGreaterThan(0);
  });

  it("REG-04: 申请走身份签名意向，拒签提示对齐(不进 pending)", () => {
    operators = [makeOperator()];
    txStatus = "idle";
    const { rerender } = renderPage();
    fireEvent.click(screen.getByText("Apply for Number"));
    // 拒签态：txState.status="error" + error 文案 → 弹层内展示拒签提示，不出现 pending/成功文案。
    txStatus = "error";
    txError = "身份签名被取消，操作未提交";
    rerender(
      <MemoryRouter initialEntries={["/region/US"]}>
        <RegionDetail />
      </MemoryRouter>,
    );
    // 拒签提示渲染（对齐 §3.7）。
    expect(screen.getByText("身份签名被取消，操作未提交")).toBeInTheDocument();
    // BottomSheet 经 Portal 渲染到 document.body，按 document 查询 reject-msg slot。
    expect(document.querySelector('[data-slot="reject-msg"]')).not.toBeNull();
    // error 状态下按钮不应处于 Applying...（不进 pending）。
    expect(screen.queryByText("Applying...")).toBeNull();
  });
});
