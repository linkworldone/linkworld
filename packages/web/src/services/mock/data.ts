import type { User, DepositInfo, DepositRecord, Region, Operator, VirtualNumber, Bill, Notification } from "@/types";

let mockUser: User | null = null;

let mockDeposit: DepositInfo = {
  balance: 150000000000000000000n,
  minimumRequired: 100000000000000000000n,
  currency: "USDT",
};

const mockDepositHistory: DepositRecord[] = [
  { id: "dr-1", type: "deposit", amount: 50000000000000000000n, currency: "USDT", timestamp: "2026-04-10T10:00:00Z", txHash: "0xabc123def456" },
  { id: "dr-2", type: "deposit", amount: 100000000000000000000n, currency: "USDT", timestamp: "2026-04-01T08:00:00Z", txHash: "0x789abcdef123" },
  { id: "dr-3", type: "deduction", amount: 15300000000000000000n, currency: "USDT", timestamp: "2026-03-31T12:00:00Z", txHash: "0xdef789abc123" },
];

const mockRegions: Region[] = [
  { code: "JP", name: "Japan", flag: "🇯🇵", operatorCount: 2, startingPrice: 10 },
  { code: "US", name: "United States", flag: "🇺🇸", operatorCount: 3, startingPrice: 8 },
  { code: "GB", name: "United Kingdom", flag: "🇬🇧", operatorCount: 2, startingPrice: 10 },
  { code: "SG", name: "Singapore", flag: "🇸🇬", operatorCount: 2, startingPrice: 12 },
  { code: "KR", name: "South Korea", flag: "🇰🇷", operatorCount: 2, startingPrice: 9 },
  { code: "DE", name: "Germany", flag: "🇩🇪", operatorCount: 2, startingPrice: 11 },
  { code: "AU", name: "Australia", flag: "🇦🇺", operatorCount: 1, startingPrice: 14 },
  { code: "TH", name: "Thailand", flag: "🇹🇭", operatorCount: 2, startingPrice: 7 },
];

const mockOperators: Record<string, Operator[]> = {
  JP: [
    { id: "op-jp-1", name: "NTT Docomo", region: "JP", requiredDeposit: 100000000000000000000n, dataRate: 3.5, callRate: 0.15, isActive: true },
    { id: "op-jp-2", name: "SoftBank", region: "JP", requiredDeposit: 80000000000000000000n, dataRate: 4.0, callRate: 0.12, isActive: true },
  ],
  US: [
    { id: "op-us-1", name: "T-Mobile", region: "US", requiredDeposit: 80000000000000000000n, dataRate: 2.5, callRate: 0.10, isActive: true },
    { id: "op-us-2", name: "AT&T", region: "US", requiredDeposit: 100000000000000000000n, dataRate: 3.0, callRate: 0.12, isActive: true },
    { id: "op-us-3", name: "Verizon", region: "US", requiredDeposit: 120000000000000000000n, dataRate: 2.8, callRate: 0.11, isActive: true },
  ],
  GB: [
    { id: "op-gb-1", name: "Vodafone UK", region: "GB", requiredDeposit: 100000000000000000000n, dataRate: 3.5, callRate: 0.18, isActive: true },
    { id: "op-gb-2", name: "EE", region: "GB", requiredDeposit: 90000000000000000000n, dataRate: 3.8, callRate: 0.15, isActive: true },
  ],
  SG: [
    { id: "op-sg-1", name: "Singtel", region: "SG", requiredDeposit: 120000000000000000000n, dataRate: 4.0, callRate: 0.20, isActive: true },
    { id: "op-sg-2", name: "StarHub", region: "SG", requiredDeposit: 100000000000000000000n, dataRate: 4.5, callRate: 0.18, isActive: true },
  ],
  KR: [
    { id: "op-kr-1", name: "SK Telecom", region: "KR", requiredDeposit: 90000000000000000000n, dataRate: 3.0, callRate: 0.12, isActive: true },
    { id: "op-kr-2", name: "KT", region: "KR", requiredDeposit: 85000000000000000000n, dataRate: 3.2, callRate: 0.14, isActive: true },
  ],
};

let mockNumbers: VirtualNumber[] = [
  {
    id: "vn-1", number: "+81 90-1234-5678", region: "JP", operator: "NTT Docomo", status: "active",
    activatedAt: "2026-03-15T10:00:00Z",
    credentials: { type: "esim", config: "LPA:1$rsp.example.com$ABCDEF1234567890" },
  },
];

const mockBills: Bill[] = [
  { id: "bill-2026-04", month: "2026-04", status: "unpaid", operatorFee: "12.00", platformFee: "0.30", totalAmount: "12.30", dueDate: "2026-04-30T12:00:00Z", usage: { dataGB: 2.4, callMinutes: 47 } },
  { id: "bill-2026-03", month: "2026-03", status: "paid", operatorFee: "14.93", platformFee: "0.37", totalAmount: "15.30", dueDate: "2026-03-31T12:00:00Z", paidAt: "2026-03-28T09:00:00Z", usage: { dataGB: 3.1, callMinutes: 62 } },
];

const mockNotifications: Notification[] = [
  { id: "notif-1", type: "bill_due", title: "Bill Payment Due", message: "Your April 2026 bill of $12.30 is due by Apr 30.", read: false, createdAt: new Date(Date.now() - 2 * 3600 * 1000).toISOString() },
  { id: "notif-2", type: "service_suspended", title: "Service Warning", message: "Your deposit balance is approaching the minimum threshold.", read: false, createdAt: new Date(Date.now() - 5 * 3600 * 1000).toISOString() },
  { id: "notif-3", type: "deposit_confirmed", title: "Deposit Confirmed", message: "50.00 USDT has been deposited to your account.", read: true, createdAt: "2026-04-10T10:00:00Z" },
  { id: "notif-4", type: "payment_confirmed", title: "Payment Confirmed", message: "March 2026 bill of $15.30 has been paid.", read: true, createdAt: "2026-03-28T09:00:00Z" },
  { id: "notif-5", type: "system", title: "System Update", message: "New regions added: South Korea, Singapore.", read: true, createdAt: "2026-03-25T08:00:00Z" },
];

export { mockDepositHistory, mockRegions, mockOperators, mockNumbers, mockBills, mockNotifications };

export function getMockUser(): User | null { return mockUser; }
export function setMockUser(user: User | null) { mockUser = user; }
export function getMockDeposit(): DepositInfo { return mockDeposit; }
export function setMockDeposit(deposit: DepositInfo) { mockDeposit = deposit; }
export function addMockDepositRecord(record: DepositRecord) { mockDepositHistory.unshift(record); }
export function addMockNumber(num: VirtualNumber) { mockNumbers.push(num); }
export function updateMockBill(billId: string, updates: Partial<Bill>) {
  const idx = mockBills.findIndex((b) => b.id === billId);
  if (idx !== -1) Object.assign(mockBills[idx], updates);
}
export function addMockNotification(notif: Notification) { mockNotifications.unshift(notif); }
