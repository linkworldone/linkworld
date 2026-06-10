import { describe, it, expect, vi, beforeEach } from "vitest";

// 付账对账重构（design §3.3 / handoff §1.1）：POST /api/bills/pay 仅写 pending 意向，
// 不据 200 置 is_paid；is_paid 唯一由后端 BillPaid 事件回填。
// Bill status 扩展 unpaid/paying/paid/overdue：is_paid=false 且有支付意向 → paying。
vi.mock("./client", () => ({
  apiClient: { get: vi.fn(), post: vi.fn() },
}));

import { apiClient } from "./client";
import { billingApi } from "./billingApi";

const mockGet = apiClient.get as unknown as ReturnType<typeof vi.fn>;
const mockPost = apiClient.post as unknown as ReturnType<typeof vi.fn>;

function makeApiBill(over: Record<string, unknown> = {}) {
  return {
    id: 1,
    user_id: 1,
    operator_id: 1,
    amount: "100000000",
    platform_fee: "1500000",
    traffic_card_deduction: "0",
    is_paid: false,
    created_at: "2026-06-01T00:00:00.000Z",
    ...over,
  };
}

describe("billingApi 付账 pending 意向（REC-03）", () => {
  beforeEach(() => {
    mockGet.mockReset();
    mockPost.mockReset();
    mockPost.mockResolvedValue(undefined);
  });

  it("REC-03a: payIntent POST /api/bills/pay 仅 pending 意向，不带 tx_hash 终态", async () => {
    await billingApi.payIntent("0xabc", "1");
    expect(mockPost).toHaveBeenCalledTimes(1);
    const [path, body] = mockPost.mock.calls[0];
    expect(path).toBe("/api/bills/pay");
    expect(body.wallet).toBe("0xabc");
    expect(body.bill_id).toBe(1);
    expect(body.tx_hash).toBeUndefined();
  });

  it("REC-03b: is_paid=false 且后端标记支付意向(pay_intent_tx_hash) → status=paying（不据 200 置 paid）", async () => {
    mockGet.mockResolvedValue([
      makeApiBill({ is_paid: false, pay_intent_tx_hash: "0xdeadbeef" }),
    ]);
    const [bill] = await billingApi.getBills("0xabc");
    expect(bill.status).toBe("paying");
  });

  it("REC-03c: is_paid=true → status=paid（后端事件回填后才置已付）", async () => {
    mockGet.mockResolvedValue([
      makeApiBill({ is_paid: true, paid_at: "2026-06-02T00:00:00.000Z" }),
    ]);
    const [bill] = await billingApi.getBills("0xabc");
    expect(bill.status).toBe("paid");
  });

  it("REC-03d: is_paid=false 且无意向 → unpaid", async () => {
    mockGet.mockResolvedValue([makeApiBill({ is_paid: false })]);
    const [bill] = await billingApi.getBills("0xabc");
    expect(bill.status).toBe("unpaid");
  });
});
