import { apiClient } from "./client";
import { signedPost } from "./signedPost";

// SIM 一步式领取（销毁流量卡 → 领取 SIM）。
// - claim：写端点，经 signedPost 带 WalletAuth 身份签名（action="sim/claim"）。
//   body 对齐后端契约：{ destination, recipient, addressLine, tokenIds, txHash }。
// - getMySims：公开读端点，返回该钱包已领 SIM 列表。

export interface SimClaimPayload {
  destination: string;
  recipient: string;
  addressLine: string;
  /** 本次销毁的流量卡 tokenId（每张 = 1 天，SIM 天数 = 卡数）。 */
  tokenIds: number[];
  /** redeemForSim 链上交易哈希。 */
  txHash: string;
}

export type SimStatus = "pending" | "confirmed";

export interface SimRecord {
  id: string;
  /** SIM 天数余额（= 销毁卡数）。 */
  days: number;
  destination: string;
  recipient: string;
  addressLine: string;
  status: SimStatus;
  createdAt: string;
}

export const simApi = {
  /**
   * 领取 SIM（链上 redeemForSim 成功后调用）。
   * 经 signedPost 带 WalletAuth 身份签名（action="sim/claim"）；拒签抛 WalletAuthRejectedError。
   */
  async claim(wallet: string, payload: SimClaimPayload): Promise<void> {
    await signedPost("/api/sim/claim", payload, {
      wallet,
      action: "sim/claim",
    });
  },

  /** 我的 SIM 列表（公开读）。 */
  async getMySims(wallet: string): Promise<SimRecord[]> {
    const data = await apiClient.get<unknown, unknown[]>(`/api/sim/${wallet}`);
    if (!Array.isArray(data)) return [];
    return data.map((raw) => {
      const r = raw as Record<string, unknown>;
      return {
        id: String(r.id ?? ""),
        days: Number(r.days ?? 0),
        destination: String(r.destination ?? ""),
        recipient: String(r.recipient ?? ""),
        addressLine: String(r.addressLine ?? ""),
        status: r.status === "confirmed" ? "confirmed" : "pending",
        createdAt: String(r.createdAt ?? ""),
      };
    });
  },
};
