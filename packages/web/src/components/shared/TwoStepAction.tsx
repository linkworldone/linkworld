import { useEffect, useState } from "react";
import { ShieldCheck, Loader2, Check } from "lucide-react";
import { Button } from "@/components/ui/button";
import { useApprove, useAllowance } from "@/hooks/contracts/useUsdtContract";
import { type TxState, useTxState } from "@/hooks/useTransactionFlow";

/**
 * TwoStepAction —— USDT approve → action 两笔串行的两步态编排组件（design §8 / B5）。
 *
 * 状态机：
 *   idle (校验在调用方做)
 *     ├─[allowance ≥ 需求额]──────────────▶ 跳过 Approve，直达 ② Action
 *     └─[allowance < 需求额] ① Approve:
 *           approve-sign → (拒签→idle)
 *                        → approving → confirming-approval → ② Action
 *      ② Action:
 *           action-sign → (拒签→★approved-idle)
 *                       → acting → (失败→★approved-idle)
 *                                → 成功 → onSuccess()（调用方写 pending 意向）
 *
 *   ★ 关键回退：Approve 成功后 Action 拒签/失败 → 回 `approved-idle`
 *     （已授权可重试，**绝不 re-approve**；重试时 allowance≥需求额直跳 ②）。
 *
 * 充值/付账复用同一组件，仅 spender / amount(需求额) / actionLabel / onAction 不同。
 *
 * 内部编排：useAllowance（读授权）+ useApprove（第一笔）+ 调用方传入的 action tx 状态（第二笔）。
 * approve 走组件内部；action 的 writeContract + tx 状态由调用方（T7 充值 / T9 付账）注入，
 * 这样 action 合约调用（deposit/payBill）保持各页自有，组件只负责两步编排与 UI。
 */

export type TwoStepPhase =
  | "idle"
  | "approve-sign"
  | "approving"
  | "confirming-approval"
  | "approved-idle"
  | "action-sign"
  | "acting"
  | "confirming"
  | "done";

export interface ActionTx {
  /** 触发第二笔交易（deposit/payBill）。 */
  write: () => void;
  /** 第二笔交易的统一 TxState（调用方用 useTxState 包 useContractDeposit/PayBill 产出）。 */
  state: TxState;
}

export interface TwoStepActionProps {
  /** 当前用户地址（读 allowance 的 owner）。 */
  owner: `0x${string}` | undefined;
  /** 授权花费方：充值=Deposit 地址；付账=Payment 地址。 */
  spender: `0x${string}` | undefined;
  /** 需求额（最小单位 bigint）：充值=amount；付账=amount+calculateFee(amount)。 */
  amount: bigint;
  /** 第二步动作（充值/付账）的 writeContract 触发 + tx 状态。 */
  action: ActionTx;
  /** 第二步按钮文案，如「存入」「支付」。 */
  actionLabel: string;
  /** 第二笔交易确认成功回调（调用方写 pending 意向 / 刷新）。 */
  onSuccess?: () => void;
  /** 整体禁用（如金额未校验通过）。 */
  disabled?: boolean;
  className?: string;
}

export function TwoStepAction({
  owner,
  spender,
  amount,
  action,
  actionLabel,
  onSuccess,
  disabled = false,
  className,
}: TwoStepActionProps) {
  const allowanceQ = useAllowance(owner, spender);
  const allowance = (allowanceQ.data as bigint | undefined) ?? 0n;
  const allowanceLoaded = allowanceQ.data !== undefined;
  // allowance ≥ 需求额 → 视为已授权（跳步 / 不 re-approve）。amount=0 视为无需授权。
  const hasAllowance = amount > 0n ? allowance >= amount : true;

  const approveTx = useApprove();
  const approveState = useTxState({
    hash: approveTx.hash,
    isPending: approveTx.isPending,
    isConfirming: approveTx.isConfirming,
    isSuccess: approveTx.isSuccess,
    error: approveTx.error,
  });

  const actionState = action.state;

  // ★ 一旦本会话内 approve 成功过，即使 allowance 重读尚未刷新，也锁定「已授权」语义，
  //   保证 action 失败回退到 approved-idle（不 re-approve）。
  const [approvedOnce, setApprovedOnce] = useState(false);
  useEffect(() => {
    if (approveState.status === "success") setApprovedOnce(true);
  }, [approveState.status]);

  // action 成功回调（只触发一次）。
  const [notifiedSuccess, setNotifiedSuccess] = useState(false);
  useEffect(() => {
    if (actionState.status === "success" && !notifiedSuccess) {
      setNotifiedSuccess(true);
      onSuccess?.();
    }
  }, [actionState.status, notifiedSuccess, onSuccess]);

  const phase = derivePhase({
    hasAllowance,
    approvedOnce,
    approveStatus: approveState.status,
    actionStatus: actionState.status,
  });

  const atAction = phase === "approved-idle" || hasAllowance || approvedOnce;
  const approveDone = atAction;

  const busy =
    phase === "approve-sign" ||
    phase === "approving" ||
    phase === "confirming-approval" ||
    phase === "action-sign" ||
    phase === "acting" ||
    phase === "confirming";

  const handleApprove = () => approveTx.approve(spender!, amount);
  const handleAction = () => action.write();

  // 主按钮：未授权时先 Approve；已授权/回退后直接 action（含「重试」语义，不再 approve）。
  const showApproveButton = !approveDone && !hasAllowance;

  const stepLabel = (p: TwoStepPhase): string => {
    switch (p) {
      case "approve-sign":
        return "请在钱包中确认授权…";
      case "approving":
      case "confirming-approval":
        return "授权确认中…";
      case "action-sign":
        return "请在钱包中确认…";
      case "acting":
      case "confirming":
        return "交易确认中…";
      default:
        return "";
    }
  };

  return (
    <div className={className} data-slot="two-step-action" data-phase={phase}>
      {/* Stepper */}
      <ol className="flex items-center gap-3 mb-4" aria-label="交易步骤">
        <li
          className="flex items-center gap-2 text-sm"
          data-step="approve"
          data-state={approveDone ? "done" : showApproveButton ? "current" : "todo"}
        >
          <span
            className={
              "flex h-6 w-6 items-center justify-center rounded-full border text-xs font-semibold " +
              (approveDone
                ? "border-status-success text-status-success"
                : "border-surface-card-line text-text-on-light-secondary")
            }
          >
            {approveDone ? <Check className="size-3.5" /> : "1"}
          </span>
          <span className={approveDone ? "text-status-success" : "text-text-on-light-primary"}>
            {approveDone ? "已授权" : "授权 USDT"}
          </span>
        </li>
        <li
          className="flex items-center gap-2 text-sm"
          data-step="action"
          data-state={
            phase === "done" ? "done" : atAction ? "current" : "todo"
          }
        >
          <span
            className={
              "flex h-6 w-6 items-center justify-center rounded-full border text-xs font-semibold " +
              (phase === "done"
                ? "border-status-success text-status-success"
                : atAction
                  ? "border-brand-royal text-brand-royal"
                  : "border-surface-card-line text-text-on-light-secondary")
            }
          >
            {phase === "done" ? <Check className="size-3.5" /> : "2"}
          </span>
          <span className="text-text-on-light-primary">{actionLabel}</span>
        </li>
      </ol>

      {/* 主按钮 */}
      {showApproveButton ? (
        <Button
          className="w-full"
          disabled={disabled || busy || !allowanceLoaded || !spender || amount <= 0n}
          onClick={handleApprove}
          data-action="approve"
        >
          {busy ? (
            <Loader2 className="size-4 animate-spin" />
          ) : (
            <ShieldCheck className="size-4" />
          )}
          授权 USDT
        </Button>
      ) : (
        <Button
          className="w-full"
          disabled={disabled || busy || !spender}
          onClick={handleAction}
          data-action="action"
        >
          {busy && <Loader2 className="size-4 animate-spin" />}
          {phase === "approved-idle" ? `重试${actionLabel}` : actionLabel}
        </Button>
      )}

      {/* 中间态文案 */}
      {busy && (
        <p className="mt-2 text-xs text-text-on-light-secondary" data-slot="step-status">
          {stepLabel(phase)}
        </p>
      )}

      {/* 错误文案（拒签/失败回退已在 phase 反映，这里给原因） */}
      {approveState.status === "error" && !approveDone && (
        <p className="mt-2 text-xs text-status-danger" data-slot="approve-error">
          {approveState.error}
        </p>
      )}
      {actionState.status === "error" && (
        <p className="mt-2 text-xs text-status-danger" data-slot="action-error">
          {actionState.error}
        </p>
      )}
    </div>
  );
}

/** 纯函数：由 allowance/approve/action 状态推导当前 phase（便于推理与单测）。 */
export function derivePhase(params: {
  hasAllowance: boolean;
  approvedOnce: boolean;
  approveStatus: TxState["status"];
  actionStatus: TxState["status"];
}): TwoStepPhase {
  const { hasAllowance, approvedOnce, approveStatus, actionStatus } = params;

  const approved = hasAllowance || approvedOnce;

  // 第二步进行中 / 完成优先反映。
  if (actionStatus === "success") return "done";
  if (actionStatus === "pending-signature") return "action-sign";
  if (actionStatus === "pending-confirmation") return "acting";

  if (approved) {
    // Approve 已完成（跳步或本会话签过）。action 拒签/失败 → approved-idle（可重试，不 re-approve）。
    if (actionStatus === "error") return "approved-idle";
    return "approved-idle";
  }

  // 尚未授权 → 第一步态。
  if (approveStatus === "pending-signature") return "approve-sign";
  if (approveStatus === "pending-confirmation") return "confirming-approval";
  // approve 拒签/失败 → 回 idle（可重试 Approve）。
  return "idle";
}
