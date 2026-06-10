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
  currency: "USDT" | "ETH";
}

export interface DepositRecord {
  id: string;
  type: "deposit" | "withdraw" | "deduction";
  amount: bigint;
  currency: string;
  // 对账三态（design §3.1）：pending（仅意向，不计入余额）/ confirmed（event_sync 回填终态）。
  status: "pending" | "confirmed";
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
  countryCode?: string;
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
  // 对账三态扩展（design §3.3）：paying = 已发起支付意向、等后端 BillPaid 事件回填（不据 200 置 paid）。
  status: "unpaid" | "paying" | "paid" | "overdue";
  operatorFee: string;
  platformFee: string;
  trafficCardDeduction?: string;
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

// --- Backend API response types (snake_case) ---
export interface ApiUser {
  id: number;
  wallet_addr: string;
  email: string;
  token_id: number;
  is_active: boolean;
  registered_at: string;
  created_at: string;
  updated_at: string;
}

export interface ApiOperator {
  id: number;
  name: string;
  region: string;
  country_code: string;
  required_deposit: string;
  is_active: boolean;
}

export interface ApiBill {
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

export interface ApiUsage {
  data_usage: number;
  call_usage: number;
  signature: string;
}
