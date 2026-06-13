import { useTranslation } from "react-i18next";
import { useFeeRate, useCalculateFee } from "@/hooks/contracts";
import { formatAmount } from "@/utils/format";

/**
 * FeeBreakdown —— 平台手续费明细（读链费率，design §3.6 / D12）。
 *
 * 铁律：
 * - 费率读 `FeeManager.getFeeRate()`（基点 150=1.5%，/10000）。
 * - 费额来源二选一：
 *   1) 付账场景传 `fee`（= 后端 Bill.platformFee = 链上 calculateFee(amount)）→ **直接展示，不读链重算**（避免在已含费的合计上二次算费）。
 *   2) 未传 `fee` 时用 `amount` 读链 `calculateFee(amount)`（**不自算**）。
 * - 删写死 `PLATFORM_FEE_RATE` 与写死「2.5%」文案。
 * - loading→skeleton（占位条）、失败→「--」，**绝不写死兜底**。
 *
 * 展示：「平台手续费 (1.5%)：N USDT」一行。
 */
export interface FeeBreakdownProps {
  /** 付账本金（USDT 6 位最小单位 bigint）；仅在未传 `fee` 时用于读链 calculateFee。 */
  amount?: bigint;
  /** 直接展示的费额（USDT 6 位最小单位 bigint）。付账场景传后端 Bill.platformFee，**不再读链重算**。 */
  fee?: bigint;
  className?: string;
  /** 文字尺寸跟随所在卡片，默认 xs。 */
  size?: "xs" | "sm";
}

function SkeletonBar({ w }: { w: string }) {
  return (
    <span
      className={`inline-block h-3 rounded bg-surface-input animate-pulse align-middle ${w}`}
      data-slot="fee-skeleton"
      aria-hidden
    />
  );
}

export function FeeBreakdown({ amount, fee, className, size = "xs" }: FeeBreakdownProps) {
  const { t } = useTranslation();
  const rate = useFeeRate();
  // 传入 fee（付账场景）→ 直接展示，不读链；否则用 amount 读链 calculateFee。
  const directFee = fee !== undefined;
  const feeQ = useCalculateFee(directFee ? undefined : amount);

  const textSize = size === "sm" ? "text-sm" : "text-xs";

  // 费率标签：loading→skeleton，失败/未回→「--」。
  const rateLabel = rate.isLoading ? (
    <SkeletonBar w="w-8" />
  ) : rate.label ? (
    rate.label
  ) : (
    "--"
  );

  // 费额：传入 fee → 直接展示；否则走读链（loading→skeleton，失败/未回→「--」）。
  const feeText = directFee ? (
    `${formatAmount(fee)} USDT`
  ) : feeQ.isLoading ? (
    <SkeletonBar w="w-12" />
  ) : feeQ.fee !== undefined ? (
    `${formatAmount(feeQ.fee)} USDT`
  ) : (
    "--"
  );

  return (
    <div
      className={`flex justify-between ${textSize} ${className ?? ""}`}
      data-slot="fee-breakdown"
    >
      <span className="text-text-on-light-secondary">
        {t("feeBreakdown.platformFee")} ({rateLabel})
      </span>
      <span className="font-data tabular-nums text-text-on-light-primary">{feeText}</span>
    </div>
  );
}
