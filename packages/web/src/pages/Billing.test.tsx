import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, fireEvent } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import type { Bill } from "@/types";

// ── 可变 mock 状态 ───────────────────────────────────────────────
let bills: Bill[] = [];
const payBillFn = vi.fn();
const recordIntentFn = vi.fn();

function makeBill(over: Partial<Bill> = {}): Bill {
  return {
    id: "1",
    month: "2026-05",
    status: "unpaid",
    // 6 位最小单位字符串（T2）：5 USDT = 5_000_000；不应再 ×10^6。
    operatorFee: "4925000",
    platformFee: "75000",
    trafficCardDeduction: "0",
    totalAmount: "5000000",
    dueDate: "2026-06-15T00:00:00.000Z",
    usage: { dataGB: 0, callMinutes: 0 },
    ...over,
  };
}

vi.mock("wagmi", () => ({
  useAccount: () => ({ address: "0xabc" as `0x${string}`, chainId: 31337 }),
  // config/wagmi 在模块加载期调 http()（经 signedPost 链）→ 须补 http 防 mock 报错。
  http: () => ({}),
}));

vi.mock("@/hooks/useBilling", () => ({
  useBills: () => ({ data: bills }),
  usePayBill: () => ({
    payBill: payBillFn,
    txState: { status: "idle" },
    recordIntent: recordIntentFn,
  }),
}));

vi.mock("@/hooks/contracts", () => ({
  useCalculateFee: () => ({ fee: 75_000n, isLoading: false, isError: false }),
  useFeeRate: () => ({ label: "1.5%", isLoading: false, isError: false }),
}));

vi.mock("@/config/contracts", () => ({
  getContractAddress: () => "0xPAYMENT" as `0x${string}`,
}));

// FeeBreakdown 读链；页面测隔离成静态行（fee 测在 FeeBreakdown.test）。
vi.mock("@/components/shared/FeeBreakdown", () => ({
  FeeBreakdown: () => <div data-testid="fee-breakdown">平台手续费 (1.5%)</div>,
}));

// TwoStepAction 依赖 wagmi allowance/approve；页面测隔离成简单按钮。
vi.mock("@/components/shared/TwoStepAction", () => ({
  TwoStepAction: ({ amount, actionLabel }: { amount: bigint; actionLabel: string }) => (
    <button data-testid="two-step" data-amount={String(amount)}>
      {actionLabel}
    </button>
  ),
}));

import Billing from "./Billing";

function renderPage() {
  return render(
    <MemoryRouter>
      <Billing />
    </MemoryRouter>,
  );
}

beforeEach(() => {
  bills = [];
  payBillFn.mockReset();
  recordIntentFn.mockReset();
});

describe("Billing 页（T9）", () => {
  it("BILL-02: totalAmount 5_000_000(6 位最小单位) → 展示 5.00 USDT，不二次缩放 ×10^6", () => {
    bills = [makeBill()];
    const { container } = renderPage();
    // 正确：5.00（formatAmount 6 位）；AmountDisplay 把数值与 USDT 拆成相邻 span，按值节点断言。
    expect(screen.getAllByText("5.00").length).toBeGreaterThanOrEqual(1);
    // 反向：绝不出现二次缩放后的 5000000(×10^6) 文本。
    expect(container.textContent).not.toContain("5000000");
    expect(container.textContent).not.toContain("5,000,000");
  });

  it("BILL-03: paying 态用 TxStatusBadge「处理中」(不染绿)，不据成功置「已支付」", () => {
    bills = [makeBill({ status: "paying" })];
    const { container } = renderPage();
    // paying → 处理中徽章（pending，禁绿）。
    const badge = container.querySelector('[data-slot="tx-status-badge"][data-status="pending"]');
    expect(badge).not.toBeNull();
    expect(screen.getByText(/处理中/)).toBeInTheDocument();
    // 不显示「已支付」终态文案，不出现 confirmed 绿态徽章。
    expect(screen.queryByText("已支付")).toBeNull();
    expect(container.querySelector('[data-status="confirmed"]')).toBeNull();
    // paying 不提供「立即支付」按钮（已发起，等回填）。
    expect(screen.queryByText("立即支付")).toBeNull();
  });

  it("付账授权额 = 本金 + calculateFee（读链 75_000n），注入 TwoStepAction amount", () => {
    bills = [makeBill()];
    renderPage();
    fireEvent.click(screen.getByText("立即支付"));
    const twoStep = screen.getByTestId("two-step");
    // 5_000_000 + 75_000 = 5_075_000（不自算，calculateFee 读链值相加）。
    expect(twoStep.getAttribute("data-amount")).toBe("5075000");
  });

  it("paid 账单进「历史」Tab 展示「已支付」，无支付按钮", () => {
    bills = [makeBill({ status: "paid", paidAt: "2026-05-20T00:00:00.000Z" })];
    renderPage();
    fireEvent.click(screen.getByText("历史"));
    expect(screen.getByText("已支付")).toBeInTheDocument();
    expect(screen.queryByText("立即支付")).toBeNull();
  });
});
