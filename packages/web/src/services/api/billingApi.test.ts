import { describe, it, expect, vi, beforeEach } from "vitest";

// billingApi.toBill 经由 apiClient.get 取数后映射，这里 mock client 验证 toBill 的金额计算。
// totalAmount 应为 6 位最小单位字符串，全程 BigInt 加减，绝不用 parseFloat（资损红线）。
vi.mock("./client", () => ({
  apiClient: {
    get: vi.fn(),
    post: vi.fn(),
  },
}));

import { apiClient } from "./client";
import { billingApi } from "./billingApi";

const mockGet = apiClient.get as unknown as ReturnType<typeof vi.fn>;

function makeApiBill(over: Record<string, unknown> = {}) {
  return {
    id: 1,
    user_id: 1,
    operator_id: 1,
    amount: "0",
    platform_fee: "0",
    traffic_card_deduction: "0",
    is_paid: false,
    created_at: "2026-06-01T00:00:00.000Z",
    ...over,
  };
}

describe("billingApi.toBill bigint 金额计算（BILL-01）", () => {
  beforeEach(() => {
    mockGet.mockReset();
  });

  it("operatorFee + platformFee - deduction 用 BigInt 加减（6 位最小单位字符串）", async () => {
    // 100 USDT 运营费 + 2.5 USDT 平台费 - 1 USDT 抵扣 = 101.5 USDT = 101_500_000
    mockGet.mockResolvedValue([
      makeApiBill({
        amount: "100000000",
        platform_fee: "2500000",
        traffic_card_deduction: "1000000",
      }),
    ]);

    const [bill] = await billingApi.getBills("0xabc");
    expect(bill.totalAmount).toBe("101500000");
  });

  it("大额不丢精度（超 Number.MAX_SAFE_INTEGER）", async () => {
    // 9_007_199_254_740_993 = MAX_SAFE_INTEGER(9_007_199_254_740_991) + 2，浮点会丢精度
    mockGet.mockResolvedValue([
      makeApiBill({
        amount: "9007199254740993",
        platform_fee: "0",
        traffic_card_deduction: "0",
      }),
    ]);

    const [bill] = await billingApi.getBills("0xabc");
    expect(bill.totalAmount).toBe("9007199254740993");
  });

  it("缺省抵扣字段按 0 处理", async () => {
    mockGet.mockResolvedValue([
      makeApiBill({ amount: "5000000", platform_fee: "125000", traffic_card_deduction: undefined }),
    ]);

    const [bill] = await billingApi.getBills("0xabc");
    expect(bill.totalAmount).toBe("5125000");
  });
});
