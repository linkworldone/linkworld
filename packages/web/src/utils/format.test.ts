import { describe, it, expect } from "vitest";
import { parseUnits, formatAmount } from "./format";

// T2 精度 6 位：USDT 金额按 6 位最小单位解析/展示，默认不再 18 位。
describe("parseUnits / formatAmount 6 位精度（PREC-01）", () => {
  it("parseUnits 默认按 6 位解析：'10' === 10_000_000n", () => {
    expect(parseUnits("10")).toBe(10_000_000n);
  });

  it("parseUnits 显式 6 位：'10' === 10_000_000n", () => {
    expect(parseUnits("10", 6)).toBe(10_000_000n);
  });

  it("parseUnits 带小数 6 位：'1.5' === 1_500_000n", () => {
    expect(parseUnits("1.5", 6)).toBe(1_500_000n);
  });

  it("formatAmount 默认按 6 位展示：10_000_000n → '10.00'（displayDecimals=2）", () => {
    expect(formatAmount(10_000_000n)).toBe("10.00");
  });

  it("formatAmount 显式 6 位 + displayDecimals=0：10_000_000n → '10'", () => {
    expect(formatAmount(10_000_000n, 6, 0)).toBe("10");
  });

  it("formatAmount 显式 6 位 2 位展示：1_500_000n → '1.50'", () => {
    expect(formatAmount(1_500_000n, 6, 2)).toBe("1.50");
  });

  it("边界 0：parseUnits('0')===0n，formatAmount(0n)==='0.00'", () => {
    expect(parseUnits("0")).toBe(0n);
    expect(formatAmount(0n)).toBe("0.00");
  });

  it("大额不丢精度：parseUnits('1000000.5') === 1_000_000_500_000n", () => {
    expect(parseUnits("1000000.5")).toBe(1_000_000_500_000n);
  });
});
