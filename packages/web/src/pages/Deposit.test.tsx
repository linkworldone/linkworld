import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, fireEvent } from "@testing-library/react";

// ── 可变 mock 状态 ───────────────────────────────────────────────
let depositBalance: bigint = 5_000_000n; // 5 USDT 链上本金
let usdtWalletBalance: bigint = 100_000_000n; // 100 USDT 钱包
let lockExpiry: bigint = 0n;
const depositFn = vi.fn();
const withdrawFn = vi.fn();

vi.mock("wagmi", () => ({
  useAccount: () => ({ address: "0xabc" as `0x${string}`, chainId: 31337 }),
  // config/wagmi 在模块加载期调 http()（经 signedPost 链）→ 须补 http 防 mock 报错。
  http: () => ({}),
}));

vi.mock("@/hooks/useDeposit", () => ({
  useDeposit: () => ({
    data: { balance: depositBalance, currency: "USDT" },
    refetch: vi.fn(),
  }),
  useDepositHistory: () => ({ data: [] }),
  useDepositMutation: () => ({
    deposit: depositFn,
    txState: { status: "idle" },
    recordIntent: vi.fn(),
    isContractPending: false,
    isConfirming: false,
    isSuccess: false,
  }),
  useWithdrawMutation: () => ({
    withdraw: withdrawFn,
    txState: { status: "idle" },
    recordIntent: vi.fn(),
    isContractPending: false,
    isConfirming: false,
    isSuccess: false,
  }),
}));

vi.mock("@/hooks/contracts", () => ({
  useUsdtBalance: () => ({ data: usdtWalletBalance }),
  // LockCountdown 内部读 useLockExpiry
  useLockExpiry: () => ({ data: lockExpiry }),
}));

vi.mock("@/config/contracts", () => ({
  getContractAddress: () => "0xDEPOSIT" as `0x${string}`,
  getUsdt: () => "0xUSDT" as `0x${string}`,
  getUsdtDecimals: () => 6,
}));

// TwoStepAction 依赖 useAllowance/useApprove（wagmi）；为页面测隔离，替身成简单渲染。
vi.mock("@/components/shared/TwoStepAction", () => ({
  TwoStepAction: ({ disabled, actionLabel }: { disabled?: boolean; actionLabel: string }) => (
    <button data-testid="two-step" disabled={disabled}>
      {actionLabel}
    </button>
  ),
}));

import Deposit from "./Deposit";

function resetMocks() {
  depositBalance = 5_000_000n;
  usdtWalletBalance = 100_000_000n;
  lockExpiry = 0n;
  depositFn.mockReset();
  withdrawFn.mockReset();
}

beforeEach(resetMocks);

describe("Deposit 页（T7）", () => {
  it("渲染余额卡 + 充值/提现按钮 + 顺延提示", () => {
    render(<Deposit />);
    expect(screen.getByText(/保证金余额/)).toBeInTheDocument();
    expect(screen.getByText("充值")).toBeInTheDocument();
    expect(screen.getByText("提现")).toBeInTheDocument();
    expect(screen.getByText(/再次充值将把锁仓期顺延 30 天/)).toBeInTheDocument();
  });

  it("DEP-03: 充值金额 < 10 USDT → 校验拦截（TwoStepAction disabled + 错误文案）", () => {
    render(<Deposit />);
    fireEvent.click(screen.getByText("充值"));
    const input = screen.getByPlaceholderText(/输入充值金额/);
    fireEvent.change(input, { target: { value: "5" } }); // 5 < 10
    expect(screen.getByText(/单次充值不少于/)).toBeInTheDocument();
    expect(screen.getByTestId("two-step")).toBeDisabled();
  });

  it("DEP-03: 充值金额 ≥ 10 且 ≤ 钱包余额 → 校验通过（TwoStepAction 可用）", () => {
    render(<Deposit />);
    fireEvent.click(screen.getByText("充值"));
    const input = screen.getByPlaceholderText(/输入充值金额/);
    fireEvent.change(input, { target: { value: "20" } });
    expect(screen.queryByText(/单次充值不少于/)).not.toBeInTheDocument();
    expect(screen.getByTestId("two-step")).not.toBeDisabled();
  });

  it("充值金额 > 钱包 USDT 余额 → 拦截", () => {
    usdtWalletBalance = 10_000_000n; // 10 USDT
    render(<Deposit />);
    fireEvent.click(screen.getByText("充值"));
    fireEvent.change(screen.getByPlaceholderText(/输入充值金额/), {
      target: { value: "50" },
    });
    expect(screen.getByText(/超出钱包 USDT 余额/)).toBeInTheDocument();
    expect(screen.getByTestId("two-step")).toBeDisabled();
  });

  it("锁仓未满 → 提现按钮禁用", () => {
    const now = Math.floor(Date.now() / 1000);
    lockExpiry = BigInt(now + 86_400 * 5); // 锁仓中
    render(<Deposit />);
    const withdrawBtn = screen.getByText("提现").closest("button")!;
    expect(withdrawBtn).toBeDisabled();
  });

  it("锁仓已满 + 有余额 → 提现按钮可点", () => {
    const now = Math.floor(Date.now() / 1000);
    lockExpiry = BigInt(now - 10); // 已解锁
    render(<Deposit />);
    const withdrawBtn = screen.getByText("提现").closest("button")!;
    expect(withdrawBtn).not.toBeDisabled();
  });

  it("换肤校验：页面无 ETH 币种、无旧 token 文本", () => {
    const { container } = render(<Deposit />);
    expect(container.textContent).not.toContain("ETH");
  });
});
