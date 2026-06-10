import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, fireEvent } from "@testing-library/react";

// ── 可变 wagmi 状态（测试逐项设置）──────────────────────────────
let allowanceData: bigint | undefined = undefined;
let approveIsPending = false;
let approveIsConfirming = false;
let approveIsSuccess = false;
let approveHash: `0x${string}` | undefined = undefined;
const approveWrite = vi.fn();

function resetWagmi() {
  allowanceData = undefined;
  approveIsPending = false;
  approveIsConfirming = false;
  approveIsSuccess = false;
  approveHash = undefined;
  approveWrite.mockReset();
}

vi.mock("wagmi", () => ({
  useChainId: () => 31337,
  // useAllowance → useReadContract
  useReadContract: () => ({ data: allowanceData }),
  // useApprove → useWriteContract
  useWriteContract: () => ({
    writeContract: (cfg: unknown) => {
      approveWrite(cfg);
    },
    data: approveHash,
    isPending: approveIsPending,
    error: null,
  }),
  // useApprove → useWaitForTransactionReceipt
  useWaitForTransactionReceipt: () => ({
    isPending: approveIsConfirming,
    isSuccess: approveIsSuccess,
  }),
}));

vi.mock("../../config/contracts", () => ({
  getUsdt: () => "0xUSDT" as `0x${string}`,
}));

import { TwoStepAction, derivePhase, type ActionTx } from "./TwoStepAction";
import type { TxState } from "@/hooks/useTransactionFlow";

const OWNER = "0xOWNER" as `0x${string}`;
const SPENDER = "0xSPENDER" as `0x${string}`;
const AMOUNT = 10_000_000n; // 10 USDT

function makeAction(state: TxState, write = vi.fn()): ActionTx {
  return { write, state };
}

function renderTSA(over: Partial<Parameters<typeof TwoStepAction>[0]> = {}) {
  const action = over.action ?? makeAction({ status: "idle" });
  const utils = render(
    <TwoStepAction
      owner={OWNER}
      spender={SPENDER}
      amount={AMOUNT}
      action={action}
      actionLabel="存入"
      {...over}
    />
  );
  return { ...utils, action };
}

beforeEach(resetWagmi);

describe("TwoStepAction 状态机（design §8 / B5）", () => {
  it("TSA-01: allowance < amount → 显示 Approve 步，点击走 approve(exact)", () => {
    allowanceData = 0n; // < 10 USDT
    renderTSA();

    const approveBtn = screen.getByRole("button");
    expect(approveBtn).toHaveAttribute("data-action", "approve");
    expect(approveBtn).toHaveTextContent("授权 USDT");

    fireEvent.click(approveBtn);
    expect(approveWrite).toHaveBeenCalledTimes(1);
    const cfg = approveWrite.mock.calls[0][0];
    expect(cfg.functionName).toBe("approve");
    expect(cfg.args).toEqual([SPENDER, AMOUNT]); // exact，禁 infinite
  });

  it("TSA-02: allowance ≥ amount → 跳过 Approve 直达 Action（主按钮=action，点击不触发 approve）", () => {
    allowanceData = 20_000_000n; // ≥ 10 USDT
    const action = makeAction({ status: "idle" });
    renderTSA({ action });

    const btn = screen.getByRole("button");
    expect(btn).toHaveAttribute("data-action", "action");
    expect(btn).toHaveTextContent("存入");
    // Step1 显示「已授权」
    expect(screen.getByText("已授权")).toBeInTheDocument();

    fireEvent.click(btn);
    expect(action.write).toHaveBeenCalledTimes(1);
    expect(approveWrite).not.toHaveBeenCalled();
  });

  it("TSA-03: approve 成功 + action 失败 → 回 approved-idle，重试不再 approve（直接 action）", () => {
    // 初始 allowance < amount，approve 成功（isSuccess），action 处于 error。
    allowanceData = 0n; // 链上读尚未刷新到新授权
    approveIsSuccess = true;
    approveHash = "0xAPPROVEHASH";
    const action = makeAction({ status: "error", error: "交易失败，请重试" });
    const { container } = renderTSA({ action });

    const root = container.querySelector("[data-slot='two-step-action']");
    expect(root).toHaveAttribute("data-phase", "approved-idle");

    // 主按钮应为「重试存入」且 data-action=action（不是 approve）。
    const btn = screen.getByRole("button");
    expect(btn).toHaveAttribute("data-action", "action");
    expect(btn).toHaveTextContent("重试存入");

    fireEvent.click(btn);
    expect(action.write).toHaveBeenCalledTimes(1);
    expect(approveWrite).not.toHaveBeenCalled(); // ★ 绝不 re-approve
  });

  it("TSA-04: approve 调用 args 为 exact amount（不是 MaxUint256）", () => {
    allowanceData = 0n;
    renderTSA();
    fireEvent.click(screen.getByRole("button"));
    const cfg = approveWrite.mock.calls[0][0];
    expect(cfg.args[1]).toBe(AMOUNT);
    expect(cfg.args[1]).not.toBe(2n ** 256n - 1n);
  });

  it("action 成功 → phase=done + 触发 onSuccess", () => {
    allowanceData = 20_000_000n;
    const onSuccess = vi.fn();
    const action = makeAction({ status: "success", txHash: "0xDONE" });
    const { container } = renderTSA({ action, onSuccess });
    const root = container.querySelector("[data-slot='two-step-action']");
    expect(root).toHaveAttribute("data-phase", "done");
    expect(onSuccess).toHaveBeenCalledTimes(1);
  });
});

describe("derivePhase 纯函数", () => {
  it("未授权 + 全 idle → idle", () => {
    expect(
      derivePhase({
        hasAllowance: false,
        approvedOnce: false,
        approveStatus: "idle",
        actionStatus: "idle",
      })
    ).toBe("idle");
  });

  it("approve 签名中 → approve-sign", () => {
    expect(
      derivePhase({
        hasAllowance: false,
        approvedOnce: false,
        approveStatus: "pending-signature",
        actionStatus: "idle",
      })
    ).toBe("approve-sign");
  });

  it("approve 确认中 → confirming-approval", () => {
    expect(
      derivePhase({
        hasAllowance: false,
        approvedOnce: false,
        approveStatus: "pending-confirmation",
        actionStatus: "idle",
      })
    ).toBe("confirming-approval");
  });

  it("allowance 充足 + action idle → approved-idle（跳步）", () => {
    expect(
      derivePhase({
        hasAllowance: true,
        approvedOnce: false,
        approveStatus: "idle",
        actionStatus: "idle",
      })
    ).toBe("approved-idle");
  });

  it("approvedOnce + action error → approved-idle（回退不 re-approve）", () => {
    expect(
      derivePhase({
        hasAllowance: false,
        approvedOnce: true,
        approveStatus: "success",
        actionStatus: "error",
      })
    ).toBe("approved-idle");
  });

  it("action 签名中 → action-sign", () => {
    expect(
      derivePhase({
        hasAllowance: true,
        approvedOnce: false,
        approveStatus: "idle",
        actionStatus: "pending-signature",
      })
    ).toBe("action-sign");
  });

  it("action 成功 → done", () => {
    expect(
      derivePhase({
        hasAllowance: true,
        approvedOnce: true,
        approveStatus: "success",
        actionStatus: "success",
      })
    ).toBe("done");
  });
});
