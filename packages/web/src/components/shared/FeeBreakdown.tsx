import { useFeeRate, useCalculateFee } from "@/hooks/contracts";
import { formatAmount } from "@/utils/format";

/**
 * FeeBreakdown —— 平台手续费明细（读链费率，design §3.6 / D12）。
 *
 * 铁律：
 * - 费率读 `FeeManager.getFeeRate()`（基点 150=1.5%，/10000）；费额读 `calculateFee(amount)`（**不自算**）。
 * - 删写死 `PLATFORM_FEE_RATE` 与写死「2.5%」文案。
 * - loading→skeleton（占位条）、失败→「--」，**绝不写死兜底**。
 *
 * 展示：「平台手续费 (1.5%)：N USDT」一行。amount=付账本金（最小单位 bigint）。
 */
export interface FeeBreakdownProps {
  /** 付账本金（USDT 6 位最小单位 bigint），用于读链 calculateFee。 */
  amount: bigint;
  className?: string;
  /** 文字尺寸跟随所在卡片，默认 xs。 */
  size?: "xs" | "sm";
}

function SkeletonBar({ w }: { w: string }) {
  return (
    <span
      className={`inline-block h-3 rounded bg-surface-secondary animate-pulse align-middle ${w}`}
      data-slot="fee-skeleton"
      aria-hidden
    />
  );
}

export function FeeBreakdown({ amount, className, size = "xs" }: FeeBreakdownProps) {
  const rate = useFeeRate();
  const feeQ = useCalculateFee(amount);

  const textSize = size === "sm" ? "text-sm" : "text-xs";

  // 费率标签：loading→skeleton，失败/未回→「--」。
  const rateLabel = rate.isLoading ? (
    <SkeletonBar w="w-8" />
  ) : rate.label ? (
    rate.label
  ) : (
    "--"
  );

  // 费额：loading→skeleton，失败/未回→「--」。
  const feeText = feeQ.isLoading ? (
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
      <span className="text-text-secondary">平台手续费 ({rateLabel})</span>
      <span className="font-data tabular-nums">{feeText}</span>
    </div>
  );
}
