import { apiClient } from "./client";
import type { Bill } from "../../types";

interface ApiBill {
  id: number;
  user_id: number;
  operator_id: number;
  amount: string;
  platform_fee: string;
  traffic_card_deduction?: string;
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

  // 金额字段均为 USDT 6 位最小单位字符串：全程 BigInt 加减，绝不用 parseFloat
  //（最小单位当元做浮点会单位语义错 + 大额超 MAX_SAFE_INTEGER 丢精度，资损红线）。
  // totalAmount 保持 6 位最小单位字符串，不落 number；展示侧经 formatAmount(total, usdtDecimals)。
  const operatorFee = api.amount;
  const platformFee = api.platform_fee;
  const trafficCardDeduction = api.traffic_card_deduction || "0";
  const total = (
    BigInt(operatorFee) + BigInt(platformFee) - BigInt(trafficCardDeduction)
  ).toString();

  return {
    id: String(api.id),
    month: `${createdAt.getFullYear()}-${String(createdAt.getMonth() + 1).padStart(2, "0")}`,
    status,
    operatorFee,
    platformFee,
    trafficCardDeduction,
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
