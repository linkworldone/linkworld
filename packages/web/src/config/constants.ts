export const PLATFORM_FEE_RATE = 0.025;
export const MIN_DEPOSIT_USDT = 100n * 10n ** 18n;
export const OVERDUE_DAYS = 14;
export const MOCK_DELAY_MS = 600;

export const SUPPORTED_CURRENCIES = ["USDT", "ETH"] as const;
export type SupportedCurrency = (typeof SUPPORTED_CURRENCIES)[number];
