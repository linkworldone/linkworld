import { formatAmount } from "@/utils/format";

/*
 * 金色用色铁律（DESIGN.md §B2，WCAG）：金额默认色按所在底色分流——
 * 卡内（暖米白底）金额用 navy（text-on-light-primary，≈12:1）；
 * 深底（navy 画布上）金额用 gold（text-on-dark-gold，≈6.5:1）。
 * 金在暖米白上 ≈2:1 不达标，禁用。默认 navy；深底显式传 tone="gold-on-dark"。
 * 显式 colorClass 优先级最高（覆盖 tone）。金额走等宽数字（font-data + tabular-nums）。
 */
type Tone = "auto" | "gold-on-dark";

interface AmountDisplayProps {
  amount: bigint | string;
  currency?: string;
  size?: "sm" | "md" | "lg";
  /** 底色分流：auto=卡内 navy（默认）；gold-on-dark=深底金 */
  tone?: Tone;
  /** 显式覆盖（优先于 tone） */
  colorClass?: string;
}

const toneClass: Record<Tone, { amount: string; currency: string }> = {
  auto: { amount: "text-text-on-light-primary", currency: "text-text-on-light-secondary" },
  "gold-on-dark": { amount: "text-text-on-dark-gold", currency: "text-text-on-dark-secondary" },
};

export function AmountDisplay({
  amount,
  currency,
  size = "md",
  tone = "auto",
  colorClass,
}: AmountDisplayProps) {
  const display = typeof amount === "bigint" ? formatAmount(amount) : amount;
  const sizeClasses = {
    sm: "text-sm font-semibold",
    md: "text-lg font-bold",
    lg: "text-4xl font-extrabold",
  };
  const tones = toneClass[tone];
  const amountColor = colorClass ?? tones.amount;

  return (
    <span className={`font-data tabular-nums ${sizeClasses[size]} ${amountColor}`}>
      {display}
      {currency && <span className={`text-xs ml-1 ${tones.currency}`}>{currency}</span>}
    </span>
  );
}
