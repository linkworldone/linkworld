import { apiClient } from "./client";
import type { Bill } from "../../types";

interface ApiBill {
  id: number;
  user_id: number;
  operator_id: number;
  amount: string;
  platform_fee: string;
  is_paid: boolean;
  created_at: string;
  paid_at?: string;
  tx_hash?: string;
}

function toBill(api: ApiBill): Bill {
  const createdAt = new Date(api.created_at);
  const dueDate = new Date(createdAt);
  dueDate.setDate(dueDate.getDate() + 14);

  let status: Bill["status"] = "unpaid";
  if (api.is_paid) {
    status = "paid";
  } else if (new Date() > dueDate) {
    status = "overdue";
  }

  const operatorFee = api.amount;
  const platformFee = api.platform_fee;
  const total = (parseFloat(operatorFee) + parseFloat(platformFee)).toFixed(2);

  return {
    id: String(api.id),
    month: `${createdAt.getFullYear()}-${String(createdAt.getMonth() + 1).padStart(2, "0")}`,
    status,
    operatorFee,
    platformFee,
    totalAmount: total,
    dueDate: dueDate.toISOString(),
    paidAt: api.paid_at || undefined,
    usage: { dataGB: 0, callMinutes: 0 }, // 需额外调 usage API
  };
}

export const billingApi = {
  async getBills(
    wallet: string,
    filter?: Bill["status"],
  ): Promise<Bill[]> {
    const data = await apiClient.get<any, ApiBill[]>(
      `/api/bills/${wallet}`,
    );
    const bills = data.map(toBill);
    if (filter) {
      return bills.filter((b) => b.status === filter);
    }
    return bills;
  },

  async recordPayment(
    wallet: string,
    billId: string,
    txHash?: string,
  ): Promise<void> {
    await apiClient.post("/api/bills/pay", {
      wallet,
      bill_id: Number(billId),
      tx_hash: txHash,
    });
  },
};
