import { formatAmount } from "@/utils/format";

interface AmountDisplayProps {
  amount: bigint | string;
  currency?: string;
  size?: "sm" | "md" | "lg";
  colorClass?: string;
}

export function AmountDisplay({ amount, currency, size = "md", colorClass = "text-status-warning" }: AmountDisplayProps) {
  const display = typeof amount === "bigint" ? formatAmount(amount) : amount;
  const sizeClasses = {
    sm: "text-sm font-semibold",
    md: "text-lg font-bold",
    lg: "text-4xl font-extrabold",
  };

  return (
    <span className={`${sizeClasses[size]} ${colorClass}`}>
      {display}
      {currency && <span className="text-xs text-text-secondary ml-1">{currency}</span>}
    </span>
  );
}
