import { useEffect, useMemo, useRef, useState } from "react";
import { Lock, LockOpen } from "lucide-react";
import { useTranslation } from "react-i18next";
import { useLockExpiry } from "@/hooks/contracts";

/**
 * LockCountdown —— 保证金锁仓倒计时 / 解锁态（design §3.4 / DESIGN.md §3）。
 *
 * 合约依据：`getLockExpiry(addr)` 返回到期 unix 秒；提现 `require(block.timestamp >= _lockExpiry)`。
 * 边界 **`>=`**：`now >= expiry` 即解锁可提（用 `remainingSec <= 0` 判定，禁用 `<`——差 1 秒会卡死到期用户）。
 * 累加语义：每次充值 `lockExpiry += 30 days` 顺延，提现成功归 0。
 *
 * 三态：
 *   expiry == 0   → 无锁仓（无押金 / 已全提）
 *   now <  expiry → 锁仓中 · 剩余 Dd Hh Mm + Lock，Withdraw 由父禁用
 *   now >= expiry → 锁仓已满，可提取本金 + LockOpen，Withdraw 可点
 *
 * **去「利息」**：合约 `DepositWithdrawn(user, principal, 0)` interest 恒 0，文案只说「本金」。
 *
 * 解锁态经 `onStateChange` 暴露给父（Deposit.tsx）以联动 Withdraw 按钮禁用。
 */

export type LockStatus = "none" | "locked" | "unlocked";

export interface LockStateResult {
  status: LockStatus;
  /** 距解锁剩余秒数（locked 时 >0；unlocked/none 为 0）。 */
  remainingSec: number;
}

/**
 * 纯函数：由 expiry(秒) 与 now(秒) 推导锁仓态（边界 `>=` 对齐合约，便于单测）。
 * @param expirySec 链上 getLockExpiry（0=无锁仓）
 * @param nowSec    当前 unix 秒
 */
export function lockState(expirySec: bigint, nowSec: number): LockStateResult {
  if (expirySec === 0n) return { status: "none", remainingSec: 0 };
  const remaining = Number(expirySec) - nowSec;
  // 边界：now >= expiry（remaining <= 0）即解锁。禁用 `<`。
  if (remaining <= 0) return { status: "unlocked", remainingSec: 0 };
  return { status: "locked", remainingSec: remaining };
}

/** 剩余秒数 → 「Dd Hh Mm」。 */
export function formatRemaining(remainingSec: number): string {
  const total = Math.max(0, remainingSec);
  const d = Math.floor(total / 86_400);
  const h = Math.floor((total % 86_400) / 3_600);
  const m = Math.floor((total % 3_600) / 60);
  return `${d}d ${h}h ${m}m`;
}

export interface LockCountdownProps {
  address: `0x${string}` | undefined;
  /** 锁仓态变化回调（父据此联动 Withdraw 禁用 + 是否可提）。 */
  onStateChange?: (state: LockStateResult) => void;
  className?: string;
}

export function LockCountdown({ address, onStateChange, className }: LockCountdownProps) {
  const { t } = useTranslation();
  const { data } = useLockExpiry(address);
  const expiry = (data as bigint | undefined) ?? 0n;

  // 每分钟刷新；剩余 ≤1 天时每秒（接近解锁时刻精确）。
  const [nowSec, setNowSec] = useState(() => Math.floor(Date.now() / 1000));
  const state = useMemo(() => lockState(expiry, nowSec), [expiry, nowSec]);

  useEffect(() => {
    if (state.status !== "locked") return; // none/unlocked 无需 tick
    const fast = state.remainingSec <= 86_400;
    const id = setInterval(
      () => setNowSec(Math.floor(Date.now() / 1000)),
      fast ? 1_000 : 60_000,
    );
    return () => clearInterval(id);
  }, [state.status, state.remainingSec <= 86_400]);

  // 暴露态给父（避免渲染期 setState：用 effect）。
  const lastNotified = useRef<string>("");
  useEffect(() => {
    const key = `${state.status}:${state.status === "unlocked" ? 1 : 0}`;
    if (key !== lastNotified.current) {
      lastNotified.current = key;
      onStateChange?.(state);
    }
  }, [state, onStateChange]);

  if (state.status === "none") {
    return (
      <div
        className={className}
        data-slot="lock-countdown"
        data-status="none"
      >
        <p className="text-xs text-text-on-light-secondary">{t("lockCountdown.none")}</p>
      </div>
    );
  }

  if (state.status === "locked") {
    return (
      <div
        className={`flex items-center gap-2 rounded-md bg-surface-input px-3 py-2 ${className ?? ""}`}
        data-slot="lock-countdown"
        data-status="locked"
      >
        <Lock className="size-4 text-brand-royal" />
        <span className="text-xs font-medium text-text-on-light-primary">
          {t("lockCountdown.remaining", { time: formatRemaining(state.remainingSec) })}
        </span>
      </div>
    );
  }

  // unlocked
  return (
    <div
      className={`flex items-center gap-2 rounded-md bg-status-success/10 px-3 py-2 ${className ?? ""}`}
      data-slot="lock-countdown"
      data-status="unlocked"
    >
      <LockOpen className="size-4 text-status-success" />
      <span className="text-xs font-medium text-status-success">
        {t("lockCountdown.unlocked")}
      </span>
    </div>
  );
}
