import type { Bill, MonthEstimate } from "@/types";
import { delay } from "./delay";
import { mockBills, updateMockBill } from "./data";

export const billingService = {
  async getBills(_address: string, filter?: "unpaid" | "paid"): Promise<Bill[]> {
    await delay();
    if (filter === "unpaid") return mockBills.filter((b) => b.status === "unpaid" || b.status === "overdue");
    if (filter === "paid") return mockBills.filter((b) => b.status === "paid");
    return [...mockBills];
  },
  async getBillDetail(billId: string): Promise<Bill> {
    await delay();
    const bill = mockBills.find((b) => b.id === billId);
    if (!bill) throw new Error(`Bill ${billId} not found`);
    return { ...bill };
  },
  async payBill(billId: string): Promise<Bill> {
    await delay(1200);
    updateMockBill(billId, { status: "paid", paidAt: new Date().toISOString() });
    return { ...mockBills.find((b) => b.id === billId)! };
  },
  async getCurrentMonthEstimate(_address: string): Promise<MonthEstimate> {
    await delay();
    return { dataGB: 2.4, callMinutes: 47, estimatedCost: "12.30" };
  },
};
