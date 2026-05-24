import { parseEther } from "viem";
import { apiClient } from "./client";
import type { DepositRecord } from "../../types";

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
      currency: "ETH",
      timestamp: r.created_at,
      txHash: r.tx_hash ?? "",
    }));
  },
  async recordDeposit(
    wallet: string,
    amount: string,
    txHash?: string,
  ): Promise<void> {
    const amountWei = parseEther(amount).toString();
    await apiClient.post("/api/deposit", {
      wallet,
      amount: amountWei,
      tx_hash: txHash,
    });
  },

  async getDepositAmount(wallet: string): Promise<string> {
    const data = await apiClient.get<any, { amount: string }>(
      `/api/deposit/${wallet}`,
    );
    return data.amount;
  },

  async recordWithdraw(
    wallet: string,
    txHash?: string,
  ): Promise<void> {
    await apiClient.post("/api/withdraw", { wallet, tx_hash: txHash });
  },
};
