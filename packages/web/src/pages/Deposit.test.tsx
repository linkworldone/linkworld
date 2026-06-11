import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, fireEvent } from "@testing-library/react";

// ── 可变 mock 状态 ───────────────────────────────────────────────
type TrancheMock = { amount: bigint; unlockAt: bigint; withdrawn: boolean };
let depositBalance: bigint = 5_000_000n; // 5 USDT 链上本金
let usdtWalletBalance: bigint = 100_000_000n; // 100 USDT 钱包
let lockExpiry: bigint = 0n;
let tranchesData: TrancheMock[] = [];
const depositFn = vi.fn();
const withdrawFn = vi.fn();
const refetchTranchesFn = vi.fn();

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
  useTranches: () => ({ tranches: tranchesData, refetch: refetchTranchesFn }),
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
  tranchesData = [];
  depositFn.mockReset();
  withdrawFn.mockReset();
  refetchTranchesFn.mockReset();
}

beforeEach(resetMocks);

describe("Deposit 页（T7）", () => {
  it("渲染余额卡 + 充值/提现按钮 + 逐笔锁仓说明", () => {
    render(<Deposit />);
    expect(screen.getByText(/保证金余额/)).toBeInTheDocument();
    expect(screen.getByText("充值")).toBeInTheDocument();
    expect(screen.getByText("提现")).toBeInTheDocument();
    expect(screen.getByText(/每笔充值独立锁仓 30 天/)).toBeInTheDocument();
  });

  it("DEP-01: 充值抽屉只含 4 个档位按钮（10/20/50/100），各标注得卡张数", () => {
    render(<Deposit />);
    fireEvent.click(screen.getByText("充值"));
    // BottomSheet 内容经 Portal 渲染到 document.body，故查 document 而非 container。
    const options = document.querySelectorAll('[data-slot="tier-option"]');
    expect(options).toHaveLength(4);
    expect(screen.getByText("10 USDT")).toBeInTheDocument();
    expect(screen.getByText("20 USDT")).toBeInTheDocument();
    expect(screen.getByText("50 USDT")).toBeInTheDocument();
    expect(screen.getByText("100 USDT")).toBeInTheDocument();
    expect(screen.getByText(/→ 得 1 张无限流量卡/)).toBeInTheDocument();
    expect(screen.getByText(/→ 得 10 张无限流量卡/)).toBeInTheDocument();
  });

  it("DEP-02: 未选档位 → 存入禁用；选中合法档位 → 存入可用", () => {
    render(<Deposit />);
    fireEvent.click(screen.getByText("充值"));
    expect(screen.getByTestId("two-step")).toBeDisabled();
    fireEvent.click(screen.getByText("20 USDT"));
    expect(screen.getByTestId("two-step")).not.toBeDisabled();
  });

  it("DEP-03: 选中档位 > 钱包 USDT 余额 → 拦截（错误文案 + 存入禁用）", () => {
    usdtWalletBalance = 30_000_000n; // 30 USDT
    render(<Deposit />);
    fireEvent.click(screen.getByText("充值"));
    fireEvent.click(screen.getByText("50 USDT")); // 50 > 30
    expect(screen.getByText(/超出钱包 USDT 余额/)).toBeInTheDocument();
    expect(screen.getByTestId("two-step")).toBeDisabled();
  });

  it("WDR-01: 无未取回笔 → 提现按钮禁用", () => {
    tranchesData = [];
    render(<Deposit />);
    const withdrawBtn = screen.getByText("提现").closest("button")!;
    expect(withdrawBtn).toBeDisabled();
  });

  it("WDR-02: 有未取回笔 → 提现按钮可点，弹层按笔列出", () => {
    const now = Math.floor(Date.now() / 1000);
    tranchesData = [
      { amount: 10_000_000n, unlockAt: BigInt(now - 10), withdrawn: false }, // 可取回
      { amount: 20_000_000n, unlockAt: BigInt(now + 86_400 * 5), withdrawn: false }, // 锁定中
    ];
    render(<Deposit />);
    const withdrawBtn = screen.getByText("提现").closest("button")!;
    expect(withdrawBtn).not.toBeDisabled();
    fireEvent.click(withdrawBtn);
    const rows = document.querySelectorAll('[data-slot="tranche-row"]');
    expect(rows).toHaveLength(2);
    // 一笔可取回（status-success）、一笔锁定中
    expect(screen.getByText(/可取回/)).toBeInTheDocument();
    expect(screen.getByText(/锁定中/)).toBeInTheDocument();
  });

  it("WDR-03: 到期笔的「取回」按钮可点 → 调 withdraw(index)；锁定笔禁用", () => {
    const now = Math.floor(Date.now() / 1000);
    tranchesData = [
      { amount: 10_000_000n, unlockAt: BigInt(now - 10), withdrawn: false }, // index 0 可取回
      { amount: 20_000_000n, unlockAt: BigInt(now + 86_400 * 5), withdrawn: false }, // index 1 锁定中
    ];
    render(<Deposit />);
    fireEvent.click(screen.getByText("提现").closest("button")!);
    const withdrawButtons = screen
      .getAllByText("取回")
      .map((el) => el.closest("button")!);
    expect(withdrawButtons).toHaveLength(2);
    expect(withdrawButtons[0]).not.toBeDisabled(); // 到期可取
    expect(withdrawButtons[1]).toBeDisabled(); // 锁定中
    fireEvent.click(withdrawButtons[0]);
    expect(withdrawFn).toHaveBeenCalledWith(0);
  });

  it("WDR-04: 已取回笔显示「已取回」且无取回按钮", () => {
    const now = Math.floor(Date.now() / 1000);
    tranchesData = [
      { amount: 10_000_000n, unlockAt: BigInt(now - 10), withdrawn: true },
    ];
    render(<Deposit />);
    fireEvent.click(screen.getByText("提现").closest("button")!);
    expect(screen.getByText(/已取回/)).toBeInTheDocument();
    expect(screen.queryByText("取回")).not.toBeInTheDocument();
  });

  it("换肤校验：页面无 ETH 币种、无旧 token 文本", () => {
    const { container } = render(<Deposit />);
    expect(container.textContent).not.toContain("ETH");
  });
});
