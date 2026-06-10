import { parseUnits } from "viem";
import { apiClient } from "./client";
import { signedPost } from "./signedPost";
import type { DepositRecord } from "../../types";

// 对账重构（design §3.3 / handoff §1.2/1.3）：
// - 充值/提现端点**仅写 pending 意向**，不据 200 / tx_hash 置终态。
// - 终态唯一由后端 event_sync 监听链上事件（DepositMade / DepositWithdrawn）回填。
// - 余额以链上 getDepositAmount 为准（source of truth, design §9）；历史/终态轮询后端 status。
// - 精度：USDT 6 位最小单位（旧 parseEther 18 位差 10^12 倍，资损红线）。
//
// usdtDecimals 全链路恒 6（MockUSDT / Arbitrum USDT）。api 层无 chainId 上下文，
// 这里以 6 解析金额；与 contracts.getUsdtDecimals 一致（deployments.usdtDecimals=6）。
const USDT_DECIMALS = 6;

export const depositApi = {
  async getHistory(wallet: string): Promise<DepositRecord[]> {
    const data = await apiClient.get<any, any[]>(
      `/api/deposit/${wallet}/history`,
    );
    if (!Array.isArray(data)) return [];
    return data.map((r) => ({
      id: String(r.id),
      type: r.type === "withdraw" ? "withdraw" : "deposit",
      amount: (() => {
        try {
          return BigInt(r.amount ?? "0");
        } catch {
          return 0n;
        }
      })(),
      currency: "USDT",
      // 后端 status：pending / confirmed（event_sync 回填）；缺省按 pending（不暗示已到账）。
      status: r.status === "confirmed" ? "confirmed" : "pending",
      timestamp: r.created_at,
      txHash: r.tx_hash ?? "",
    }));
  },

  /**
   * 充值 pending 意向（design §3.3）。仅上报「我已发起充值」，**不据 200 置终态**。
   * 余额以链上 getDepositAmount 确认；不带 tx_hash 作为记账依据。
   * T5：写端点经 signedPost 带 WalletAuth 身份签名（action="deposit"）；拒签抛 WalletAuthRejectedError。
   */
  async postDepositIntent(wallet: string, amount: string): Promise<void> {
    const amountMinUnit = parseUnits(amount, USDT_DECIMALS).toString();
    await signedPost(
      "/api/deposit",
      { wallet, amount: amountMinUnit },
      { wallet, action: "deposit" },
    );
  },

  async getDepositAmount(wallet: string): Promise<string> {
    const data = await apiClient.get<any, { amount: string }>(
      `/api/deposit/${wallet}`,
    );
    return data.amount;
  },

  /**
   * 提现 pending 意向（design §3.3 / handoff §1.2）。
   * **废弃凭 tx_hash 记账**：仅上报 wallet；记账唯一由后端监听 DepositWithdrawn 事件回填。
   * T5：经 signedPost 带 WalletAuth 身份签名（action="withdraw"）。
   */
  async postWithdrawIntent(wallet: string): Promise<void> {
    await signedPost("/api/withdraw", { wallet }, { wallet, action: "withdraw" });
  },
};
