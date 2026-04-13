export type UserStatus = "inactive" | "active" | "suspended";

export interface User {
  address: string;
  email: string;
  status: UserStatus;
  registeredAt: string;
  nftTokenId: number;
}

export interface DepositInfo {
  balance: bigint;
  minimumRequired: bigint;
  currency: "USDT" | "ETH";
}

export interface DepositRecord {
  id: string;
  type: "deposit" | "withdraw" | "deduction";
  amount: bigint;
  currency: string;
  timestamp: string;
  txHash: string;
}

export interface Region {
  code: string;
  name: string;
  flag: string;
  operatorCount: number;
  startingPrice: number;
}

export interface Operator {
  id: string;
  name: string;
  region: string;
  requiredDeposit: bigint;
  dataRate: number;
  callRate: number;
  isActive: boolean;
}

export interface VirtualNumber {
  id: string;
  number: string;
  region: string;
  operator: string;
  status: "active" | "inactive";
  activatedAt: string;
  credentials?: {
    type: "esim" | "voip";
    config: string;
  };
}

export interface Bill {
  id: string;
  month: string;
  status: "unpaid" | "paid" | "overdue";
  operatorFee: string;
  platformFee: string;
  totalAmount: string;
  dueDate: string;
  paidAt?: string;
  usage: { dataGB: number; callMinutes: number };
}

export interface MonthEstimate {
  dataGB: number;
  callMinutes: number;
  estimatedCost: string;
}

export interface Notification {
  id: string;
  type: "bill_due" | "payment_confirmed" | "deposit_confirmed" | "service_suspended" | "system";
  title: string;
  message: string;
  read: boolean;
  createdAt: string;
}
