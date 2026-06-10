import { describe, it, expect, vi, beforeEach } from "vitest";

// 对账重构（design §3.3 / handoff §1）：充值/提现端点仅写 pending 意向，
// 不据 200 / tx_hash 置终态。余额以链上 getDepositAmount 为准，历史以后端事件回填为准。
vi.mock("./client", () => ({
  apiClient: { get: vi.fn(), post: vi.fn() },
}));

import { apiClient } from "./client";
import { depositApi } from "./depositApi";

const mockPost = apiClient.post as unknown as ReturnType<typeof vi.fn>;

describe("depositApi pending 意向（REC-01 / REC-02）", () => {
  beforeEach(() => {
    mockPost.mockReset();
    mockPost.mockResolvedValue(undefined);
  });

  it("REC-01: 充值意向 POST /api/deposit 带 6 位精度 amount，不把 tx_hash 当终态依据", async () => {
    // 1.5 USDT → 6 位最小单位 = 1_500_000（绝非 parseEther 的 18 位）
    await depositApi.postDepositIntent("0xabc", "1.5");
    expect(mockPost).toHaveBeenCalledTimes(1);
    const [path, body] = mockPost.mock.calls[0];
    expect(path).toBe("/api/deposit");
    expect(body.wallet).toBe("0xabc");
    expect(body.amount).toBe("1500000");
    // pending 意向：不依赖 tx_hash 置终态（后端 event_sync 唯一回填）。
    expect(body.tx_hash).toBeUndefined();
  });

  it("REC-02: 提现意向 POST /api/withdraw 仅 wallet，废弃凭 tx_hash 记账", async () => {
    await depositApi.postWithdrawIntent("0xabc");
    expect(mockPost).toHaveBeenCalledTimes(1);
    const [path, body] = mockPost.mock.calls[0];
    expect(path).toBe("/api/withdraw");
    expect(body.wallet).toBe("0xabc");
    expect(body.tx_hash).toBeUndefined();
  });
});
