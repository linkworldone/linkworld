import { apiClient } from "./client";
import { signedPost } from "./signedPost";
import type { Bill } from "../../types";

interface ApiBill {
  id: number;
  user_id: number;
  operator_id: number;
  amount: string;
  platform_fee: string;
  traffic_card_deduction?: string;
  is_paid: boolean;
  // 后端 pending 支付意向标记（handoff §1.1 PayIntentTxHash）：
  // is_paid=false 但已发起支付 → status=paying（等 BillPaid 事件回填 is_paid）。
  pay_intent_tx_hash?: string;
  created_at: string;
  paid_at?: string;
  tx_hash?: string;
}

function toBill(api: ApiBill): Bill {
  const createdAt = new Date(api.created_at);
  const dueDate = new Date(createdAt);
  dueDate.setDate(dueDate.getDate() + 14);

  // 对账三态（design §3.3）：is_paid 唯一由后端 BillPaid 事件回填——**不据 HTTP 200 置 paid**。
  // 已发起支付意向（pay_intent_tx_hash）但事件未回填 → paying（确认中，UI info 蓝 + Loader2，禁绿）。
  let status: Bill["status"] = "unpaid";
  if (api.is_paid) {
    status = "paid";
  } else if (api.pay_intent_tx_hash) {
    status = "paying";
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

  /**
   * 付账 pending 意向（design §3.3 / handoff §1.1）。
   * 仅上报「我已发起支付」，**不据 200 置 is_paid**；is_paid 唯一由后端 BillPaid 事件回填。
   * 不带 tx_hash 作为终态依据（后端 event_sync 监听链上事件）。
   * T5：经 signedPost 带 WalletAuth 身份签名（action="bills/pay"）。
   */
  async payIntent(wallet: string, billId: string): Promise<void> {
    await signedPost(
      "/api/bills/pay",
      { wallet, bill_id: Number(billId) },
      { wallet, action: "bills/pay" },
    );
  },
};
