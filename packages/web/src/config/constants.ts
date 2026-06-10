// 手续费率不写死：改读链 `useFeeRate`（getFeeRate 基点 150=1.5%，/10000）+ `useCalculateFee`
// （calculateFee 精确费额），见 hooks/contracts/useFeeManager.ts（T9）。原 PLATFORM_FEE_RATE 已删。

// USDT 6 位最小单位；值 10 对齐链上 require(amount >= 10 USDT) 下限。
// （此前 100n*10n**18n 含两处 bug：精度 18→6 差 10^12 倍，值 100→10。）
export const MIN_DEPOSIT_USDT = 10n * 10n ** 6n;
export const OVERDUE_DAYS = 14;
export const MOCK_DELAY_MS = 600;

export const SUPPORTED_CURRENCIES = ["USDT"] as const;
export type SupportedCurrency = (typeof SUPPORTED_CURRENCIES)[number];
