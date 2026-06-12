import { describe, it, expect } from "vitest";
import { MIN_DEPOSIT_USDT, SUPPORTED_CURRENCIES } from "./constants";

// T2 常量对齐：精度 6 位 + 值对齐链上 require ≥10 USDT + 币种仅 USDT。
describe("constants 对齐（PREC-02 / PREC-03）", () => {
  it("MIN_DEPOSIT_USDT === 10n * 10n ** 6n（10 USDT，6 位最小单位）", () => {
    expect(MIN_DEPOSIT_USDT).toBe(10n * 10n ** 6n);
  });

  it("MIN_DEPOSIT_USDT 数值为 10_000_000n", () => {
    expect(MIN_DEPOSIT_USDT).toBe(10_000_000n);
  });

  it("SUPPORTED_CURRENCIES === ['USDT']（去 ETH）", () => {
    expect(SUPPORTED_CURRENCIES).toEqual(["USDT"]);
  });
});
