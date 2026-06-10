/**
 * @deprecated 写死费率废弃。手续费改读链：`useFeeRate` 读 `getFeeRate()`（基点 150=1.5%，/10000），
 * 精确费额走 `calculateFee(amount)`（T9 实现）。请勿在新代码引用此常量。
 */
export const PLATFORM_FEE_RATE = 0.025;

// USDT 6 位最小单位；值 10 对齐链上 require(amount >= 10 USDT) 下限。
// （此前 100n*10n**18n 含两处 bug：精度 18→6 差 10^12 倍，值 100→10。）
export const MIN_DEPOSIT_USDT = 10n * 10n ** 6n;
export const OVERDUE_DAYS = 14;
export const MOCK_DELAY_MS = 600;

export const SUPPORTED_CURRENCIES = ["USDT"] as const;
export type SupportedCurrency = (typeof SUPPORTED_CURRENCIES)[number];
