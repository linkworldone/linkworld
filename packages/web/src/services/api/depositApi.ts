import { apiClient } from "./client";

export const depositApi = {
  async recordDeposit(
    wallet: string,
    amount: string,
    txHash?: string,
  ): Promise<void> {
    await apiClient.post("/api/deposit", {
      wallet,
      amount,
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
