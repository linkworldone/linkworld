export function shortenAddress(address: string, chars = 4): string {
  return `${address.slice(0, chars + 2)}...${address.slice(-chars)}`;
}

// USDT 6 位精度：金额一律按 6 位最小单位解析/展示（usdtDecimals 从 deployments 读，见 config/contracts）。
// 默认 6（非 18）——18 当 6 差 10^12 倍，资损红线。调用方可显式传精度。
const USDT_DECIMALS = 6;

export function parseUnits(value: string, decimals = USDT_DECIMALS): bigint {
  const [whole = "0", fraction = ""] = value.split(".");
  const paddedFraction = fraction.slice(0, decimals).padEnd(decimals, "0");
  return BigInt(whole + paddedFraction);
}

export function formatAmount(
  wei: bigint,
  decimals = USDT_DECIMALS,
  displayDecimals = 2,
): string {
  const divisor = 10n ** BigInt(decimals);
  const whole = wei / divisor;
  const fraction = wei % divisor;
  if (displayDecimals <= 0) return `${whole}`;
  const fractionStr = fraction.toString().padStart(decimals, "0").slice(0, displayDecimals);
  return `${whole}.${fractionStr}`;
}

export function formatUSD(amount: string): string {
  const num = parseFloat(amount);
  return `$${num.toFixed(2)}`;
}

export function formatDate(isoDate: string): string {
  return new Date(isoDate).toLocaleDateString("en-US", {
    month: "short",
    day: "numeric",
    year: "numeric",
  });
}

export function timeAgo(isoDate: string): string {
  const seconds = Math.floor((Date.now() - new Date(isoDate).getTime()) / 1000);
  if (seconds < 60) return "just now";
  if (seconds < 3600) return `${Math.floor(seconds / 60)} minutes ago`;
  if (seconds < 86400) return `${Math.floor(seconds / 3600)} hours ago`;
  return formatDate(isoDate);
}
