import { apiClient } from "./client";
import { signedPost } from "./signedPost";

// SIM 一步式领取（销毁流量卡 → 领取 SIM）。
// - claim：写端点，经 signedPost 带 WalletAuth 身份签名（action="sim/claim"）。
//   body 对齐后端契约：{ destination, recipient, addressLine, tokenIds, txHash }。
// - getMySims：公开读端点，返回该钱包已领 SIM 列表。

/** 交付方式：eSIM（扫码激活，无需邮寄）或 physical（实体卡邮寄）。 */
export type SimDeliveryType = "esim" | "physical";

export interface SimClaimPayload {
  destination: string;
  /** 交付方式：esim 无需地址；physical 需 recipient + addressLine。 */
  deliveryType: SimDeliveryType;
  /** physical 必填；esim 传空串。 */
  recipient?: string;
  /** physical 必填；esim 传空串。 */
  addressLine?: string;
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
  /** 交付方式：esim | physical（后端非 "esim" 一律归一为 "physical"）。 */
  deliveryType: SimDeliveryType;
  recipient: string;
  addressLine: string;
  /** eSIM 激活链接（esim 有值，physical 为 ""）。 */
  activationUrl: string;
  status: SimStatus;
  createdAt: string;
}

// 后端原始记录 → SimRecord 归一化（claim 响应里的 sim 与 getMySims 共用）。
function normalizeSim(raw: unknown): SimRecord {
  const r = raw as Record<string, unknown>;
  return {
    id: String(r.id ?? ""),
    days: Number(r.days ?? 0),
    destination: String(r.destination ?? ""),
    deliveryType: r.deliveryType === "esim" ? "esim" : "physical",
    recipient: String(r.recipient ?? ""),
    addressLine: String(r.addressLine ?? ""),
    activationUrl: String(r.activationUrl ?? ""),
    status: r.status === "confirmed" ? "confirmed" : "pending",
    createdAt: String(r.createdAt ?? ""),
  };
}

export const simApi = {
  /**
   * 领取 SIM（链上 redeemForSim 成功后调用），返回归一化后的 SimRecord。
   * 经 signedPost 带 WalletAuth 身份签名（action="sim/claim"）；拒签抛 WalletAuthRejectedError。
   * 解析响应里的 sim 字段，使兑换成功后能立刻拿到 activationUrl 渲染二维码。
   */
  async claim(wallet: string, payload: SimClaimPayload): Promise<SimRecord> {
    const res = await signedPost<{ sim?: unknown }>("/api/sim/claim", payload, {
      wallet,
      action: "sim/claim",
    });
    return normalizeSim(res?.sim);
  },

  /** 我的 SIM 列表（公开读）。 */
  async getMySims(wallet: string): Promise<SimRecord[]> {
    const data = await apiClient.get<unknown, unknown[]>(`/api/sim/${wallet}`);
    if (!Array.isArray(data)) return [];
    return data.map(normalizeSim);
  },
};
