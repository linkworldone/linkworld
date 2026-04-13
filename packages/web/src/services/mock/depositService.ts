import type { DepositInfo, DepositRecord } from "@/types";
import { delay } from "./delay";
import { mockDeposit, mockDepositHistory, setMockDeposit, addMockDepositRecord } from "./data";

export const depositService = {
  async getDeposit(_address: string): Promise<DepositInfo> {
    await delay();
    return { ...mockDeposit };
  },
  async getDepositHistory(_address: string): Promise<DepositRecord[]> {
    await delay();
    return [...mockDepositHistory];
  },
  async deposit(_address: string, amount: bigint, currency: string): Promise<DepositRecord> {
    await delay(1200);
    setMockDeposit({ ...mockDeposit, balance: mockDeposit.balance + amount });
    const record: DepositRecord = { id: `dr-${Date.now()}`, type: "deposit", amount, currency, timestamp: new Date().toISOString(), txHash: `0x${Math.random().toString(16).slice(2, 14)}` };
    addMockDepositRecord(record);
    return record;
  },
  async withdraw(_address: string, amount: bigint): Promise<DepositRecord> {
    await delay(1200);
    const newBalance = mockDeposit.balance - amount;
    setMockDeposit({ ...mockDeposit, balance: newBalance > 0n ? newBalance : 0n });
    const record: DepositRecord = { id: `dr-${Date.now()}`, type: "withdraw", amount, currency: mockDeposit.currency, timestamp: new Date().toISOString(), txHash: `0x${Math.random().toString(16).slice(2, 14)}` };
    addMockDepositRecord(record);
    return record;
  },
};
