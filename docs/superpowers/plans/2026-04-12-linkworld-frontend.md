# LinkWorld 前端实现计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 构建 LinkWorld 移动端优先的 Web3 DApp 前端，包含钱包连接、押金管理、虚拟号码申请、账单管理、通知中心 6 大模块，合约和后端使用 Mock 数据。

**Architecture:** 三层架构 — Pages/Components 消费 React Query Hooks，Hooks 调用 Service 层，Service 层当前返回 Mock 数据。移动端优先布局，底部 Tab Bar 导航，暗色主题。

**Tech Stack:** React 19, Vite 6, TailwindCSS 3.4, shadcn/ui, wagmi 2.14, RainbowKit 2.2, TanStack React Query 5.62, React Router 7

**Design Spec:** `docs/superpowers/specs/2026-04-12-linkworld-frontend-design.md`

---

## File Structure

```
packages/web/src/
├── types/index.ts                    # 所有 TypeScript 类型定义
├── config/
│   ├── chains.ts                     # 已有 — 0G Chain 配置
│   ├── wagmi.ts                      # wagmi + RainbowKit provider 配置
│   └── constants.ts                  # 业务常量（fee rate, min deposit 等）
├── services/
│   ├── mock/
│   │   ├── data.ts                   # Mock 原始数据
│   │   ├── userService.ts            # 用户注册/查询
│   │   ├── depositService.ts         # 押金存取/历史
│   │   ├── billingService.ts         # 账单查询/支付
│   │   ├── operatorService.ts        # 地区/运营商/号码
│   │   └── notificationService.ts    # 通知列表/已读
│   └── index.ts                      # 统一导出
├── hooks/
│   ├── useUser.ts                    # 用户 profile + 注册
│   ├── useDeposit.ts                 # 押金余额 + 存取
│   ├── useBilling.ts                 # 账单列表 + 支付
│   ├── useOperator.ts               # 地区/运营商/号码申请
│   └── useNotification.ts           # 通知 + 未读数
├── contexts/
│   └── UserContext.tsx               # 全局用户状态
├── components/
│   ├── ui/                           # shadcn/ui 组件（自动生成）
│   ├── layout/
│   │   ├── AppLayout.tsx             # Header + Content + TabBar
│   │   ├── Header.tsx                # 顶部栏
│   │   ├── TabBar.tsx                # 底部 Tab 导航
│   │   └── BottomSheet.tsx           # 底部弹出 Sheet
│   ├── wallet/
│   │   ├── ConnectButton.tsx         # RainbowKit 封装
│   │   └── RegisterSheet.tsx         # 注册流程（邮箱验证）
│   └── guards/
│       └── AuthGuard.tsx             # 路由守卫
├── pages/
│   ├── Landing.tsx                   # 营销首页
│   ├── Dashboard.tsx                 # 用户主页
│   ├── Deposit.tsx                   # 押金管理
│   ├── Services.tsx                  # 地区列表
│   ├── RegionDetail.tsx              # 运营商列表 + 号码申请
│   ├── Billing.tsx                   # 账单列表
│   ├── BillDetail.tsx                # 账单详情 + 支付
│   └── Notifications.tsx             # 通知中心
├── utils/
│   └── format.ts                     # 格式化工具
├── App.tsx                           # 路由 + providers
├── main.tsx                          # 入口
└── index.css                         # 全局样式 + 暗色主题
```

---

## Task 1: 基础设施 — shadcn/ui + TailwindCSS 暗色主题

**Files:**
- Modify: `packages/web/tailwind.config.ts`
- Modify: `packages/web/src/index.css`
- Modify: `packages/web/tsconfig.json`
- Create: `packages/web/components.json` (shadcn/ui config)

- [ ] **Step 1: 安装 shadcn/ui 及依赖**

```bash
cd packages/web
pnpm add tailwindcss-animate class-variance-authority clsx tailwind-merge lucide-react
```

- [ ] **Step 2: 创建 shadcn/ui 配置文件**

创建 `packages/web/components.json`:

```json
{
  "$schema": "https://ui.shadcn.com/schema.json",
  "style": "new-york",
  "rsc": false,
  "tsx": true,
  "tailwind": {
    "config": "tailwind.config.ts",
    "css": "src/index.css",
    "baseColor": "zinc",
    "cssVariables": true
  },
  "aliases": {
    "components": "@/components",
    "utils": "@/lib/utils",
    "ui": "@/components/ui",
    "hooks": "@/hooks",
    "lib": "@/lib"
  }
}
```

- [ ] **Step 3: 创建 shadcn/ui 工具函数**

创建 `packages/web/src/lib/utils.ts`:

```typescript
import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}
```

- [ ] **Step 4: 更新 tailwind.config.ts — 暗色主题 + shadcn/ui**

替换 `packages/web/tailwind.config.ts` 完整内容:

```typescript
import type { Config } from "tailwindcss";
import tailwindcssAnimate from "tailwindcss-animate";

const config: Config = {
  darkMode: "class",
  content: ["./index.html", "./src/**/*.{js,ts,jsx,tsx}"],
  theme: {
    extend: {
      colors: {
        background: "#0a0a14",
        foreground: "#e0e0e0",
        card: "#0f0f1a",
        "card-foreground": "#e0e0e0",
        muted: "#1a1a2e",
        "muted-foreground": "#888888",
        primary: "#3b82f6",
        "primary-foreground": "#ffffff",
        success: "#22c55e",
        warning: "#f59e0b",
        danger: "#ef4444",
        info: "#06b6d4",
        accent: "#8b5cf6",
        border: "#1a1a2e",
        ring: "#3b82f6",
      },
      fontFamily: {
        sans: ["Inter", "system-ui", "sans-serif"],
      },
      borderRadius: {
        lg: "12px",
        md: "8px",
        sm: "6px",
      },
    },
  },
  plugins: [tailwindcssAnimate],
};

export default config;
```

- [ ] **Step 5: 更新 index.css — 暗色基础样式**

替换 `packages/web/src/index.css` 完整内容:

```css
@tailwind base;
@tailwind components;
@tailwind utilities;

@layer base {
  * {
    border-color: theme("colors.border");
  }

  body {
    background-color: theme("colors.background");
    color: theme("colors.foreground");
    font-family: "Inter", system-ui, sans-serif;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
  }
}

@layer utilities {
  .safe-bottom {
    padding-bottom: env(safe-area-inset-bottom, 0px);
  }
}
```

- [ ] **Step 6: 安装 shadcn/ui 基础组件**

```bash
cd packages/web
npx shadcn@latest add button card badge input dialog sheet toast sonner separator tabs scroll-area
```

如果 `npx shadcn@latest` 交互式提示确认配置，选 yes。组件会安装到 `src/components/ui/`。

- [ ] **Step 7: 验证开发服务器启动**

```bash
cd packages/web && pnpm dev
```

打开 http://localhost:3000 确认页面正常渲染（暗色背景）。

- [ ] **Step 8: 提交**

```bash
git add packages/web/
git commit -m "chore: 配置 shadcn/ui + TailwindCSS 暗色主题

- 安装 shadcn/ui 及基础组件（button, card, badge, input, dialog, sheet, toast 等）
- 配置暗色主题色板（background, card, muted, primary 等）
- 添加 safe-area 工具类"
```

---

## Task 2: 基础设施 — wagmi + RainbowKit + Providers

**Files:**
- Create: `packages/web/src/config/wagmi.ts`
- Modify: `packages/web/src/main.tsx`

- [ ] **Step 1: 创建 wagmi 配置**

创建 `packages/web/src/config/wagmi.ts`:

```typescript
import { getDefaultConfig } from "@rainbow-me/rainbowkit";
import { zgMainnet, zgTestnet } from "./chains";

export const wagmiConfig = getDefaultConfig({
  appName: "LinkWorld",
  projectId: "linkworld-dev", // WalletConnect projectId, 开发阶段用占位符
  chains: [zgMainnet, zgTestnet],
});
```

- [ ] **Step 2: 更新 main.tsx — 集成所有 Providers**

替换 `packages/web/src/main.tsx` 完整内容:

```typescript
import React from "react";
import ReactDOM from "react-dom/client";
import { BrowserRouter } from "react-router-dom";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { WagmiProvider } from "wagmi";
import { RainbowKitProvider, darkTheme } from "@rainbow-me/rainbowkit";
import "@rainbow-me/rainbowkit/styles.css";
import App from "./App";
import "./index.css";

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 30_000,
      retry: 1,
    },
  },
});

// wagmi config 延迟导入避免循环依赖
import { wagmiConfig } from "./config/wagmi";

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <WagmiProvider config={wagmiConfig}>
      <QueryClientProvider client={queryClient}>
        <RainbowKitProvider
          theme={darkTheme({
            accentColor: "#3b82f6",
            borderRadius: "medium",
          })}
        >
          <BrowserRouter>
            <App />
          </BrowserRouter>
        </RainbowKitProvider>
      </QueryClientProvider>
    </WagmiProvider>
  </React.StrictMode>
);
```

- [ ] **Step 3: 验证启动无报错**

```bash
cd packages/web && pnpm dev
```

确认控制台无报错，RainbowKit provider 正常加载。

- [ ] **Step 4: 提交**

```bash
git add packages/web/src/config/wagmi.ts packages/web/src/main.tsx
git commit -m "feat: 集成 wagmi + RainbowKit provider

- 配置 0G Chain mainnet/testnet
- 集成 QueryClient + WagmiProvider + RainbowKitProvider
- 暗色主题配色"
```

---

## Task 3: 类型定义 + 常量 + 工具函数

**Files:**
- Create: `packages/web/src/types/index.ts`
- Create: `packages/web/src/config/constants.ts`
- Create: `packages/web/src/utils/format.ts`

- [ ] **Step 1: 创建类型定义**

创建 `packages/web/src/types/index.ts`:

```typescript
// === User ===
export type UserStatus = "inactive" | "active" | "suspended";

export interface User {
  address: string;
  email: string;
  status: UserStatus;
  registeredAt: string;
  nftTokenId: number;
}

// === Deposit ===
export interface DepositInfo {
  balance: string; // 格式化后的数字字符串，如 "50.00"
  minimumRequired: string;
  currency: "USDT" | "ETH";
}

export interface DepositRecord {
  id: string;
  type: "deposit" | "withdraw" | "deduction";
  amount: string;
  currency: string;
  timestamp: string;
  txHash: string;
}

// === Operator & Number ===
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
  requiredDeposit: string;
  dataRate: number;
  callRate: number;
  isActive: boolean;
}

export interface VirtualNumber {
  id: string;
  number: string;
  region: string;
  regionName: string;
  regionFlag: string;
  operator: string;
  operatorName: string;
  status: "active" | "inactive";
  activatedAt: string;
}

// === Billing ===
export type BillStatus = "unpaid" | "paid" | "overdue";

export interface Bill {
  id: string;
  month: string;
  status: BillStatus;
  operatorFee: number;
  platformFee: number;
  totalAmount: number;
  dueDate: string;
  paidAt?: string;
  usage: {
    dataGB: number;
    callMinutes: number;
  };
}

export interface MonthEstimate {
  dataGB: number;
  callMinutes: number;
  estimatedCost: number;
}

// === Notification ===
export type NotificationType =
  | "bill_due"
  | "payment_confirmed"
  | "deposit_confirmed"
  | "service_suspended"
  | "system";

export interface Notification {
  id: string;
  type: NotificationType;
  title: string;
  message: string;
  read: boolean;
  createdAt: string;
}
```

- [ ] **Step 2: 创建业务常量**

创建 `packages/web/src/config/constants.ts`:

```typescript
export const PLATFORM_FEE_RATE = 0.025; // 2.5%

export const MIN_DEPOSIT_USDT = "20.00";

export const SUPPORTED_CURRENCIES = ["USDT", "ETH"] as const;

export const BILL_DUE_DAY = 15; // 每月 15 号为账单截止日

export const OVERDUE_GRACE_DAYS = 14; // 逾期 2 周扣押金

export const NOTIFICATION_TYPE_COLORS: Record<string, string> = {
  bill_due: "#3b82f6",
  payment_confirmed: "#22c55e",
  deposit_confirmed: "#22c55e",
  service_suspended: "#ef4444",
  system: "#888888",
};
```

- [ ] **Step 3: 创建格式化工具**

创建 `packages/web/src/utils/format.ts`:

```typescript
/**
 * 缩写钱包地址: 0x1234...5678
 */
export function shortenAddress(address: string, chars = 4): string {
  return `${address.slice(0, chars + 2)}...${address.slice(-chars)}`;
}

/**
 * 格式化美元金额: $12.40
 */
export function formatUSD(amount: number): string {
  return `$${amount.toFixed(2)}`;
}

/**
 * 格式化代币金额: 50.00 USDT
 */
export function formatToken(amount: string, currency: string): string {
  return `${amount} ${currency}`;
}

/**
 * 格式化相对时间: 2h ago, 3d ago, Mar 15
 */
export function formatRelativeTime(isoDate: string): string {
  const now = Date.now();
  const date = new Date(isoDate).getTime();
  const diffMs = now - date;
  const diffMin = Math.floor(diffMs / 60000);
  const diffHour = Math.floor(diffMs / 3600000);
  const diffDay = Math.floor(diffMs / 86400000);

  if (diffMin < 60) return `${diffMin}m ago`;
  if (diffHour < 24) return `${diffHour}h ago`;
  if (diffDay < 7) return `${diffDay}d ago`;
  return new Date(isoDate).toLocaleDateString("en-US", {
    month: "short",
    day: "numeric",
  });
}

/**
 * 格式化月份: April 2026
 */
export function formatMonth(monthStr: string): string {
  const [year, month] = monthStr.split("-");
  const date = new Date(Number(year), Number(month) - 1);
  return date.toLocaleDateString("en-US", { month: "long", year: "numeric" });
}
```

- [ ] **Step 4: 提交**

```bash
git add packages/web/src/types/ packages/web/src/config/constants.ts packages/web/src/utils/
git commit -m "feat: 类型定义、业务常量和格式化工具

- 定义 User, Deposit, Operator, Bill, Notification 等全部业务类型
- 业务常量：费率、最低押金、通知类型配色
- 工具函数：地址缩写、金额格式化、相对时间"
```

---

## Task 4: Mock 数据 + Service 层

**Files:**
- Create: `packages/web/src/services/mock/data.ts`
- Create: `packages/web/src/services/mock/userService.ts`
- Create: `packages/web/src/services/mock/depositService.ts`
- Create: `packages/web/src/services/mock/billingService.ts`
- Create: `packages/web/src/services/mock/operatorService.ts`
- Create: `packages/web/src/services/mock/notificationService.ts`
- Create: `packages/web/src/services/index.ts`

- [ ] **Step 1: 创建 Mock 数据**

创建 `packages/web/src/services/mock/data.ts`:

```typescript
import type {
  User,
  DepositInfo,
  DepositRecord,
  Region,
  Operator,
  VirtualNumber,
  Bill,
  Notification,
} from "@/types";

// 模拟延迟
export const delay = (ms = 300) =>
  new Promise((resolve) => setTimeout(resolve, ms));

// === User ===
export const mockUsers: Record<string, User> = {};

// === Deposit ===
export const mockDeposit: DepositInfo = {
  balance: "50.00",
  minimumRequired: "20.00",
  currency: "USDT",
};

export const mockDepositHistory: DepositRecord[] = [
  {
    id: "d1",
    type: "deposit",
    amount: "20.00",
    currency: "USDT",
    timestamp: new Date(Date.now() - 2 * 3600000).toISOString(),
    txHash: "0xabc123...def456",
  },
  {
    id: "d2",
    type: "deduction",
    amount: "15.30",
    currency: "USDT",
    timestamp: new Date(Date.now() - 3 * 86400000).toISOString(),
    txHash: "0x789abc...123def",
  },
  {
    id: "d3",
    type: "deposit",
    amount: "50.00",
    currency: "USDT",
    timestamp: new Date(Date.now() - 28 * 86400000).toISOString(),
    txHash: "0xfed987...654cba",
  },
];

// === Regions & Operators ===
export const mockRegions: Region[] = [
  { code: "JP", name: "Japan", flag: "🇯🇵", operatorCount: 3, startingPrice: 8 },
  { code: "US", name: "United States", flag: "🇺🇸", operatorCount: 5, startingPrice: 10 },
  { code: "GB", name: "United Kingdom", flag: "🇬🇧", operatorCount: 4, startingPrice: 7 },
  { code: "SG", name: "Singapore", flag: "🇸🇬", operatorCount: 2, startingPrice: 6 },
  { code: "KR", name: "South Korea", flag: "🇰🇷", operatorCount: 3, startingPrice: 9 },
  { code: "DE", name: "Germany", flag: "🇩🇪", operatorCount: 3, startingPrice: 8 },
  { code: "TH", name: "Thailand", flag: "🇹🇭", operatorCount: 2, startingPrice: 5 },
  { code: "AU", name: "Australia", flag: "🇦🇺", operatorCount: 3, startingPrice: 11 },
];

export const mockOperators: Record<string, Operator[]> = {
  JP: [
    { id: "jp-sb", name: "SoftBank", region: "JP", requiredDeposit: "20.00", dataRate: 3.5, callRate: 0.15, isActive: true },
    { id: "jp-dc", name: "NTT Docomo", region: "JP", requiredDeposit: "25.00", dataRate: 4.0, callRate: 0.12, isActive: true },
    { id: "jp-au", name: "AU by KDDI", region: "JP", requiredDeposit: "20.00", dataRate: 3.8, callRate: 0.14, isActive: true },
  ],
  US: [
    { id: "us-tm", name: "T-Mobile", region: "US", requiredDeposit: "30.00", dataRate: 5.0, callRate: 0.10, isActive: true },
    { id: "us-vz", name: "Verizon", region: "US", requiredDeposit: "35.00", dataRate: 5.5, callRate: 0.08, isActive: true },
    { id: "us-at", name: "AT&T", region: "US", requiredDeposit: "30.00", dataRate: 4.8, callRate: 0.10, isActive: true },
  ],
  GB: [
    { id: "gb-vo", name: "Vodafone UK", region: "GB", requiredDeposit: "20.00", dataRate: 3.0, callRate: 0.12, isActive: true },
    { id: "gb-ee", name: "EE", region: "GB", requiredDeposit: "22.00", dataRate: 3.2, callRate: 0.11, isActive: true },
  ],
  SG: [
    { id: "sg-st", name: "Singtel", region: "SG", requiredDeposit: "18.00", dataRate: 2.8, callRate: 0.10, isActive: true },
    { id: "sg-m1", name: "M1", region: "SG", requiredDeposit: "18.00", dataRate: 2.5, callRate: 0.12, isActive: true },
  ],
  KR: [
    { id: "kr-sk", name: "SK Telecom", region: "KR", requiredDeposit: "25.00", dataRate: 4.0, callRate: 0.10, isActive: true },
    { id: "kr-kt", name: "KT", region: "KR", requiredDeposit: "22.00", dataRate: 3.8, callRate: 0.11, isActive: true },
  ],
};

// === Virtual Numbers ===
export const mockNumbers: VirtualNumber[] = [
  {
    id: "vn1",
    number: "+81 90-1234-5678",
    region: "JP",
    regionName: "Japan",
    regionFlag: "🇯🇵",
    operator: "jp-sb",
    operatorName: "SoftBank",
    status: "active",
    activatedAt: new Date(Date.now() - 28 * 86400000).toISOString(),
  },
];

// === Bills ===
export const mockBills: Bill[] = [
  {
    id: "bill-2026-04",
    month: "2026-04",
    status: "unpaid",
    operatorFee: 12.1,
    platformFee: 0.3,
    totalAmount: 12.4,
    dueDate: "2026-04-15",
    usage: { dataGB: 2.4, callMinutes: 45 },
  },
  {
    id: "bill-2026-03",
    month: "2026-03",
    status: "paid",
    operatorFee: 14.9,
    platformFee: 0.37,
    totalAmount: 15.27,
    dueDate: "2026-03-15",
    paidAt: "2026-03-31",
    usage: { dataGB: 3.1, callMinutes: 62 },
  },
  {
    id: "bill-2026-02",
    month: "2026-02",
    status: "paid",
    operatorFee: 11.51,
    platformFee: 0.29,
    totalAmount: 11.8,
    dueDate: "2026-02-15",
    paidAt: "2026-02-28",
    usage: { dataGB: 1.8, callMinutes: 38 },
  },
];

// === Notifications ===
export const mockNotifications: Notification[] = [
  {
    id: "n1",
    type: "bill_due",
    title: "Bill Due Reminder",
    message: "Your April bill of $12.40 is due in 3 days. Pay now to avoid service suspension.",
    read: false,
    createdAt: new Date(Date.now() - 2 * 3600000).toISOString(),
  },
  {
    id: "n2",
    type: "deposit_confirmed",
    title: "Deposit Confirmed",
    message: "Your deposit of 20.00 USDT has been confirmed on 0G Chain.",
    read: false,
    createdAt: new Date(Date.now() - 86400000).toISOString(),
  },
  {
    id: "n3",
    type: "payment_confirmed",
    title: "March Bill Paid",
    message: "Payment of $15.27 received. Your account remains active.",
    read: true,
    createdAt: "2026-03-31T12:00:00Z",
  },
  {
    id: "n4",
    type: "system",
    title: "Number Activated",
    message: "Your virtual number +81 90-1234-5678 is now active in Japan.",
    read: true,
    createdAt: "2026-03-15T10:00:00Z",
  },
  {
    id: "n5",
    type: "system",
    title: "Welcome to LinkWorld!",
    message: "Your account has been created. Deposit funds to get started.",
    read: true,
    createdAt: "2026-03-15T09:00:00Z",
  },
];
```

- [ ] **Step 2: 创建 userService**

创建 `packages/web/src/services/mock/userService.ts`:

```typescript
import type { User } from "@/types";
import { delay, mockUsers } from "./data";

export const userService = {
  async getUserProfile(address: string): Promise<User | null> {
    await delay();
    return mockUsers[address.toLowerCase()] ?? null;
  },

  async register(address: string, email: string): Promise<User> {
    await delay(500);
    const user: User = {
      address: address.toLowerCase(),
      email,
      status: "inactive",
      registeredAt: new Date().toISOString(),
      nftTokenId: Math.floor(Math.random() * 10000),
    };
    mockUsers[address.toLowerCase()] = user;
    return user;
  },

  async verifyEmail(_address: string, _code: string): Promise<boolean> {
    await delay(300);
    // Mock: 任何验证码都通过
    return true;
  },
};
```

- [ ] **Step 3: 创建 depositService**

创建 `packages/web/src/services/mock/depositService.ts`:

```typescript
import type { DepositInfo, DepositRecord } from "@/types";
import { delay, mockDeposit, mockDepositHistory } from "./data";

export const depositService = {
  async getDeposit(_address: string): Promise<DepositInfo> {
    await delay();
    return { ...mockDeposit };
  },

  async getDepositHistory(_address: string): Promise<DepositRecord[]> {
    await delay();
    return [...mockDepositHistory];
  },

  async deposit(
    _address: string,
    amount: string,
    currency: string
  ): Promise<DepositRecord> {
    await delay(800);
    const record: DepositRecord = {
      id: `d${Date.now()}`,
      type: "deposit",
      amount,
      currency,
      timestamp: new Date().toISOString(),
      txHash: `0x${Math.random().toString(16).slice(2, 14)}...`,
    };
    mockDeposit.balance = (
      parseFloat(mockDeposit.balance) + parseFloat(amount)
    ).toFixed(2);
    mockDepositHistory.unshift(record);
    return record;
  },

  async withdraw(_address: string, amount: string): Promise<DepositRecord> {
    await delay(800);
    const record: DepositRecord = {
      id: `d${Date.now()}`,
      type: "withdraw",
      amount,
      currency: mockDeposit.currency,
      timestamp: new Date().toISOString(),
      txHash: `0x${Math.random().toString(16).slice(2, 14)}...`,
    };
    mockDeposit.balance = (
      parseFloat(mockDeposit.balance) - parseFloat(amount)
    ).toFixed(2);
    mockDepositHistory.unshift(record);
    return record;
  },
};
```

- [ ] **Step 4: 创建 billingService**

创建 `packages/web/src/services/mock/billingService.ts`:

```typescript
import type { Bill, MonthEstimate } from "@/types";
import { delay, mockBills } from "./data";

export const billingService = {
  async getBills(
    _address: string,
    filter?: "unpaid" | "paid"
  ): Promise<Bill[]> {
    await delay();
    if (filter) {
      return mockBills.filter((b) => b.status === filter);
    }
    return [...mockBills];
  },

  async getBillDetail(billId: string): Promise<Bill> {
    await delay();
    const bill = mockBills.find((b) => b.id === billId);
    if (!bill) throw new Error(`Bill ${billId} not found`);
    return { ...bill };
  },

  async payBill(billId: string): Promise<Bill> {
    await delay(800);
    const bill = mockBills.find((b) => b.id === billId);
    if (!bill) throw new Error(`Bill ${billId} not found`);
    bill.status = "paid";
    bill.paidAt = new Date().toISOString();
    return { ...bill };
  },

  async getCurrentMonthEstimate(_address: string): Promise<MonthEstimate> {
    await delay();
    return {
      dataGB: 2.4,
      callMinutes: 45,
      estimatedCost: 12.4,
    };
  },
};
```

- [ ] **Step 5: 创建 operatorService**

创建 `packages/web/src/services/mock/operatorService.ts`:

```typescript
import type { Region, Operator, VirtualNumber } from "@/types";
import {
  delay,
  mockRegions,
  mockOperators,
  mockNumbers,
} from "./data";

export const operatorService = {
  async getRegions(): Promise<Region[]> {
    await delay();
    return [...mockRegions];
  },

  async getOperatorsByRegion(regionCode: string): Promise<Operator[]> {
    await delay();
    return mockOperators[regionCode] ?? [];
  },

  async getMyNumbers(_address: string): Promise<VirtualNumber[]> {
    await delay();
    return [...mockNumbers];
  },

  async applyNumber(
    _address: string,
    operatorId: string
  ): Promise<VirtualNumber> {
    await delay(1000);
    const allOps = Object.values(mockOperators).flat();
    const op = allOps.find((o) => o.id === operatorId);
    const region = mockRegions.find((r) => r.code === op?.region);
    const number: VirtualNumber = {
      id: `vn${Date.now()}`,
      number: `+${Math.floor(Math.random() * 90 + 10)} ${Math.floor(
        Math.random() * 9000 + 1000
      )}-${Math.floor(Math.random() * 9000 + 1000)}`,
      region: op?.region ?? "",
      regionName: region?.name ?? "",
      regionFlag: region?.flag ?? "",
      operator: operatorId,
      operatorName: op?.name ?? "",
      status: "active",
      activatedAt: new Date().toISOString(),
    };
    mockNumbers.push(number);
    return number;
  },
};
```

- [ ] **Step 6: 创建 notificationService**

创建 `packages/web/src/services/mock/notificationService.ts`:

```typescript
import type { Notification } from "@/types";
import { delay, mockNotifications } from "./data";

export const notificationService = {
  async getNotifications(_address: string): Promise<Notification[]> {
    await delay();
    return [...mockNotifications];
  },

  async getUnreadCount(_address: string): Promise<number> {
    await delay(100);
    return mockNotifications.filter((n) => !n.read).length;
  },

  async markAsRead(notificationId: string): Promise<void> {
    await delay(200);
    const n = mockNotifications.find((n) => n.id === notificationId);
    if (n) n.read = true;
  },

  async markAllAsRead(_address: string): Promise<void> {
    await delay(200);
    mockNotifications.forEach((n) => (n.read = true));
  },
};
```

- [ ] **Step 7: 创建统一导出**

创建 `packages/web/src/services/index.ts`:

```typescript
export { userService } from "./mock/userService";
export { depositService } from "./mock/depositService";
export { billingService } from "./mock/billingService";
export { operatorService } from "./mock/operatorService";
export { notificationService } from "./mock/notificationService";
```

- [ ] **Step 8: 提交**

```bash
git add packages/web/src/services/
git commit -m "feat: Mock 数据和 Service 层

- 5 个 service 模块：user, deposit, billing, operator, notification
- 完整 mock 数据集：8 个地区、多个运营商、账单历史、通知列表
- 统一导出入口，未来可替换为真实 API"
```

---

## Task 5: React Query Hooks

**Files:**
- Create: `packages/web/src/hooks/useUser.ts`
- Create: `packages/web/src/hooks/useDeposit.ts`
- Create: `packages/web/src/hooks/useBilling.ts`
- Create: `packages/web/src/hooks/useOperator.ts`
- Create: `packages/web/src/hooks/useNotification.ts`

- [ ] **Step 1: 创建 useUser hook**

创建 `packages/web/src/hooks/useUser.ts`:

```typescript
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { userService } from "@/services";

export function useUserProfile(address: string | undefined) {
  return useQuery({
    queryKey: ["user", address],
    queryFn: () => userService.getUserProfile(address!),
    enabled: !!address,
  });
}

export function useRegister() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ address, email }: { address: string; email: string }) =>
      userService.register(address, email),
    onSuccess: (user) => {
      queryClient.setQueryData(["user", user.address], user);
    },
  });
}

export function useVerifyEmail() {
  return useMutation({
    mutationFn: ({ address, code }: { address: string; code: string }) =>
      userService.verifyEmail(address, code),
  });
}
```

- [ ] **Step 2: 创建 useDeposit hook**

创建 `packages/web/src/hooks/useDeposit.ts`:

```typescript
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { depositService } from "@/services";

export function useDepositInfo(address: string | undefined) {
  return useQuery({
    queryKey: ["deposit", address],
    queryFn: () => depositService.getDeposit(address!),
    enabled: !!address,
  });
}

export function useDepositHistory(address: string | undefined) {
  return useQuery({
    queryKey: ["deposit-history", address],
    queryFn: () => depositService.getDepositHistory(address!),
    enabled: !!address,
  });
}

export function useDoDeposit() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({
      address,
      amount,
      currency,
    }: {
      address: string;
      amount: string;
      currency: string;
    }) => depositService.deposit(address, amount, currency),
    onSuccess: (_, { address }) => {
      queryClient.invalidateQueries({ queryKey: ["deposit", address] });
      queryClient.invalidateQueries({ queryKey: ["deposit-history", address] });
    },
  });
}

export function useWithdraw() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({
      address,
      amount,
    }: {
      address: string;
      amount: string;
    }) => depositService.withdraw(address, amount),
    onSuccess: (_, { address }) => {
      queryClient.invalidateQueries({ queryKey: ["deposit", address] });
      queryClient.invalidateQueries({ queryKey: ["deposit-history", address] });
    },
  });
}
```

- [ ] **Step 3: 创建 useBilling hook**

创建 `packages/web/src/hooks/useBilling.ts`:

```typescript
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { billingService } from "@/services";

export function useBills(address: string | undefined, filter?: "unpaid" | "paid") {
  return useQuery({
    queryKey: ["bills", address, filter],
    queryFn: () => billingService.getBills(address!, filter),
    enabled: !!address,
  });
}

export function useBillDetail(billId: string | undefined) {
  return useQuery({
    queryKey: ["bill", billId],
    queryFn: () => billingService.getBillDetail(billId!),
    enabled: !!billId,
  });
}

export function usePayBill() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (billId: string) => billingService.payBill(billId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["bills"] });
      queryClient.invalidateQueries({ queryKey: ["bill"] });
    },
  });
}

export function useMonthEstimate(address: string | undefined) {
  return useQuery({
    queryKey: ["month-estimate", address],
    queryFn: () => billingService.getCurrentMonthEstimate(address!),
    enabled: !!address,
  });
}
```

- [ ] **Step 4: 创建 useOperator hook**

创建 `packages/web/src/hooks/useOperator.ts`:

```typescript
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { operatorService } from "@/services";

export function useRegions() {
  return useQuery({
    queryKey: ["regions"],
    queryFn: () => operatorService.getRegions(),
  });
}

export function useOperatorsByRegion(regionCode: string | undefined) {
  return useQuery({
    queryKey: ["operators", regionCode],
    queryFn: () => operatorService.getOperatorsByRegion(regionCode!),
    enabled: !!regionCode,
  });
}

export function useMyNumbers(address: string | undefined) {
  return useQuery({
    queryKey: ["my-numbers", address],
    queryFn: () => operatorService.getMyNumbers(address!),
    enabled: !!address,
  });
}

export function useApplyNumber() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({
      address,
      operatorId,
    }: {
      address: string;
      operatorId: string;
    }) => operatorService.applyNumber(address, operatorId),
    onSuccess: (_, { address }) => {
      queryClient.invalidateQueries({ queryKey: ["my-numbers", address] });
    },
  });
}
```

- [ ] **Step 5: 创建 useNotification hook**

创建 `packages/web/src/hooks/useNotification.ts`:

```typescript
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { notificationService } from "@/services";

export function useNotifications(address: string | undefined) {
  return useQuery({
    queryKey: ["notifications", address],
    queryFn: () => notificationService.getNotifications(address!),
    enabled: !!address,
  });
}

export function useUnreadCount(address: string | undefined) {
  return useQuery({
    queryKey: ["unread-count", address],
    queryFn: () => notificationService.getUnreadCount(address!),
    enabled: !!address,
    refetchInterval: 30_000, // 30 秒轮询
  });
}

export function useMarkAsRead() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (notificationId: string) =>
      notificationService.markAsRead(notificationId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["notifications"] });
      queryClient.invalidateQueries({ queryKey: ["unread-count"] });
    },
  });
}

export function useMarkAllAsRead() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (address: string) =>
      notificationService.markAllAsRead(address),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["notifications"] });
      queryClient.invalidateQueries({ queryKey: ["unread-count"] });
    },
  });
}
```

- [ ] **Step 6: 提交**

```bash
git add packages/web/src/hooks/
git commit -m "feat: React Query hooks

- useUser: profile 查询 + 注册 mutation
- useDeposit: 余额/历史查询 + 存取 mutation
- useBilling: 账单查询 + 支付 mutation + 月估算
- useOperator: 地区/运营商查询 + 号码申请 mutation
- useNotification: 通知列表 + 未读数 + 标记已读"
```

---

## Task 6: UserContext + AuthGuard + 路由

**Files:**
- Create: `packages/web/src/contexts/UserContext.tsx`
- Create: `packages/web/src/components/guards/AuthGuard.tsx`
- Modify: `packages/web/src/App.tsx`

- [ ] **Step 1: 创建 UserContext**

创建 `packages/web/src/contexts/UserContext.tsx`:

```typescript
import { createContext, useContext, type ReactNode } from "react";
import { useAccount } from "wagmi";
import { useUserProfile } from "@/hooks/useUser";
import type { User, UserStatus } from "@/types";

interface UserContextValue {
  address: string | undefined;
  isConnected: boolean;
  user: User | null | undefined;
  isLoading: boolean;
  status: UserStatus | "unregistered";
}

const UserContext = createContext<UserContextValue>({
  address: undefined,
  isConnected: false,
  user: undefined,
  isLoading: false,
  status: "unregistered",
});

export function UserProvider({ children }: { children: ReactNode }) {
  const { address, isConnected } = useAccount();
  const { data: user, isLoading } = useUserProfile(address);

  const status: UserStatus | "unregistered" =
    !isConnected || !user ? "unregistered" : user.status;

  return (
    <UserContext.Provider
      value={{ address, isConnected, user, isLoading, status }}
    >
      {children}
    </UserContext.Provider>
  );
}

export function useUserContext() {
  return useContext(UserContext);
}
```

- [ ] **Step 2: 创建 AuthGuard**

创建 `packages/web/src/components/guards/AuthGuard.tsx`:

```typescript
import { Navigate } from "react-router-dom";
import { useUserContext } from "@/contexts/UserContext";

interface AuthGuardProps {
  children: React.ReactNode;
  /** 哪些状态可以访问此页面 */
  allowedStatuses?: Array<"inactive" | "active" | "suspended">;
  /** 被拒绝时跳转到哪 */
  fallbackPath?: string;
}

export function AuthGuard({
  children,
  allowedStatuses = ["inactive", "active", "suspended"],
  fallbackPath = "/",
}: AuthGuardProps) {
  const { isConnected, status, isLoading } = useUserContext();

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen bg-background">
        <div className="animate-spin w-8 h-8 border-2 border-primary border-t-transparent rounded-full" />
      </div>
    );
  }

  if (!isConnected || status === "unregistered") {
    return <Navigate to={fallbackPath} replace />;
  }

  if (!allowedStatuses.includes(status as "inactive" | "active" | "suspended")) {
    // INACTIVE 想访问受限页面 → 去 deposit
    if (status === "inactive") return <Navigate to="/deposit" replace />;
    // SUSPENDED 想访问受限页面 → 去 billing
    if (status === "suspended") return <Navigate to="/billing" replace />;
    return <Navigate to="/dashboard" replace />;
  }

  return <>{children}</>;
}
```

- [ ] **Step 3: 更新 App.tsx — 完整路由**

替换 `packages/web/src/App.tsx` 完整内容:

```typescript
import { Routes, Route, Navigate } from "react-router-dom";
import { UserProvider, useUserContext } from "@/contexts/UserContext";
import { AuthGuard } from "@/components/guards/AuthGuard";
import Landing from "@/pages/Landing";
import Dashboard from "@/pages/Dashboard";
import Deposit from "@/pages/Deposit";
import Services from "@/pages/Services";
import RegionDetail from "@/pages/RegionDetail";
import Billing from "@/pages/Billing";
import BillDetail from "@/pages/BillDetail";
import Notifications from "@/pages/Notifications";

function AppRoutes() {
  const { isConnected, status } = useUserContext();

  return (
    <Routes>
      <Route
        path="/"
        element={
          isConnected && status !== "unregistered" ? (
            <Navigate to="/dashboard" replace />
          ) : (
            <Landing />
          )
        }
      />
      <Route
        path="/dashboard"
        element={
          <AuthGuard>
            <Dashboard />
          </AuthGuard>
        }
      />
      <Route
        path="/deposit"
        element={
          <AuthGuard>
            <Deposit />
          </AuthGuard>
        }
      />
      <Route
        path="/services"
        element={
          <AuthGuard allowedStatuses={["active"]}>
            <Services />
          </AuthGuard>
        }
      />
      <Route
        path="/services/:regionCode"
        element={
          <AuthGuard allowedStatuses={["active"]}>
            <RegionDetail />
          </AuthGuard>
        }
      />
      <Route
        path="/billing"
        element={
          <AuthGuard allowedStatuses={["active", "suspended"]}>
            <Billing />
          </AuthGuard>
        }
      />
      <Route
        path="/billing/:billId"
        element={
          <AuthGuard allowedStatuses={["active", "suspended"]}>
            <BillDetail />
          </AuthGuard>
        }
      />
      <Route
        path="/notifications"
        element={
          <AuthGuard allowedStatuses={["active"]}>
            <Notifications />
          </AuthGuard>
        }
      />
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  );
}

export default function App() {
  return (
    <UserProvider>
      <AppRoutes />
    </UserProvider>
  );
}
```

- [ ] **Step 4: 创建占位页面**

为了让路由不报错，先创建所有页面的最小占位版本。每个文件内容格式相同：

`packages/web/src/pages/Landing.tsx`:
```typescript
export default function Landing() {
  return <div className="min-h-screen bg-background text-foreground p-4">Landing</div>;
}
```

对以下页面重复此模式（替换组件名和文字）：
- `packages/web/src/pages/Dashboard.tsx` → "Dashboard"
- `packages/web/src/pages/Deposit.tsx` → "Deposit"
- `packages/web/src/pages/Services.tsx` → "Services"
- `packages/web/src/pages/RegionDetail.tsx` → "RegionDetail"
- `packages/web/src/pages/Billing.tsx` → "Billing"
- `packages/web/src/pages/BillDetail.tsx` → "BillDetail"
- `packages/web/src/pages/Notifications.tsx` → "Notifications"

删除旧的 `packages/web/src/pages/Home.tsx`。

- [ ] **Step 5: 验证路由正常**

```bash
cd packages/web && pnpm dev
```

访问 http://localhost:3000 应看到 Landing 占位页面。

- [ ] **Step 6: 提交**

```bash
git add packages/web/src/
git rm packages/web/src/pages/Home.tsx
git commit -m "feat: UserContext + AuthGuard + 完整路由

- UserContext 管理全局用户状态（钱包地址、注册状态）
- AuthGuard 根据用户状态控制页面访问权限
- 8 个页面路由 + 占位组件
- 删除旧 Home.tsx"
```

---

## Task 7: Layout 组件 — AppLayout + TabBar + Header + BottomSheet

**Files:**
- Create: `packages/web/src/components/layout/TabBar.tsx`
- Create: `packages/web/src/components/layout/Header.tsx`
- Create: `packages/web/src/components/layout/AppLayout.tsx`
- Create: `packages/web/src/components/layout/BottomSheet.tsx`

- [ ] **Step 1: 创建 TabBar**

创建 `packages/web/src/components/layout/TabBar.tsx`:

```typescript
import { useLocation, useNavigate } from "react-router-dom";
import { Home, Smartphone, Wallet, FileText, Bell } from "lucide-react";
import { useUnreadCount } from "@/hooks/useNotification";
import { useBills } from "@/hooks/useBilling";
import { useUserContext } from "@/contexts/UserContext";

const tabs = [
  { path: "/dashboard", label: "Home", icon: Home },
  { path: "/services", label: "Services", icon: Smartphone },
  { path: "/deposit", label: "Deposit", icon: Wallet },
  { path: "/billing", label: "Bills", icon: FileText },
  { path: "/notifications", label: "Alerts", icon: Bell },
] as const;

export function TabBar() {
  const location = useLocation();
  const navigate = useNavigate();
  const { address } = useUserContext();
  const { data: unreadCount } = useUnreadCount(address);
  const { data: unpaidBills } = useBills(address, "unpaid");

  const getBadge = (path: string): number | undefined => {
    if (path === "/notifications" && unreadCount && unreadCount > 0)
      return unreadCount;
    if (path === "/billing" && unpaidBills && unpaidBills.length > 0)
      return unpaidBills.length;
    return undefined;
  };

  return (
    <nav className="fixed bottom-0 left-0 right-0 h-[60px] bg-card border-t border-border flex justify-around items-center safe-bottom z-50">
      {tabs.map(({ path, label, icon: Icon }) => {
        const isActive =
          location.pathname === path ||
          location.pathname.startsWith(path + "/");
        const badge = getBadge(path);

        return (
          <button
            key={path}
            onClick={() => navigate(path)}
            className="flex flex-col items-center gap-0.5 min-w-[48px] min-h-[44px] justify-center relative"
          >
            <div className="relative">
              <Icon
                size={20}
                className={isActive ? "text-primary" : "text-muted-foreground"}
              />
              {badge !== undefined && (
                <span className="absolute -top-1 -right-2 bg-danger text-white text-[8px] rounded-full min-w-[14px] h-[14px] flex items-center justify-center px-0.5">
                  {badge}
                </span>
              )}
            </div>
            <span
              className={`text-[10px] ${
                isActive ? "text-primary" : "text-muted-foreground"
              }`}
            >
              {label}
            </span>
          </button>
        );
      })}
    </nav>
  );
}
```

- [ ] **Step 2: 创建 Header**

创建 `packages/web/src/components/layout/Header.tsx`:

```typescript
import { useNavigate, useLocation } from "react-router-dom";
import { ArrowLeft, Bell } from "lucide-react";
import { useUserContext } from "@/contexts/UserContext";
import { useUnreadCount } from "@/hooks/useNotification";
import { shortenAddress } from "@/utils/format";

interface HeaderProps {
  title?: string;
  showBack?: boolean;
}

export function Header({ title, showBack }: HeaderProps) {
  const navigate = useNavigate();
  const location = useLocation();
  const { address } = useUserContext();
  const { data: unreadCount } = useUnreadCount(address);

  const isDashboard = location.pathname === "/dashboard";

  return (
    <header className="sticky top-0 z-40 bg-background/80 backdrop-blur-sm px-5 pt-[env(safe-area-inset-top)] pb-3">
      <div className="flex items-center justify-between h-11">
        <div className="flex items-center gap-3">
          {showBack && (
            <button
              onClick={() => navigate(-1)}
              className="min-w-[44px] min-h-[44px] flex items-center justify-center -ml-3"
            >
              <ArrowLeft size={20} className="text-foreground" />
            </button>
          )}
          {isDashboard ? (
            <div>
              <div className="text-[13px] text-muted-foreground">
                Welcome back
              </div>
              <div className="text-[11px] text-muted-foreground/60">
                {address ? shortenAddress(address) : ""}
              </div>
            </div>
          ) : (
            <h1 className="text-[17px] font-bold text-foreground">{title}</h1>
          )}
        </div>
        <div className="flex items-center gap-3">
          <button
            onClick={() => navigate("/notifications")}
            className="relative min-w-[44px] min-h-[44px] flex items-center justify-center"
          >
            <Bell size={20} className="text-muted-foreground" />
            {unreadCount !== undefined && unreadCount > 0 && (
              <span className="absolute top-2 right-2 bg-danger w-2 h-2 rounded-full" />
            )}
          </button>
          <div className="w-8 h-8 rounded-full bg-gradient-to-br from-primary to-accent flex items-center justify-center">
            <span className="text-[11px] text-white font-semibold">LW</span>
          </div>
        </div>
      </div>
    </header>
  );
}
```

- [ ] **Step 3: 创建 BottomSheet**

创建 `packages/web/src/components/layout/BottomSheet.tsx`:

```typescript
import { type ReactNode } from "react";
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
} from "@/components/ui/sheet";

interface BottomSheetProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  title: string;
  children: ReactNode;
}

export function BottomSheet({
  open,
  onOpenChange,
  title,
  children,
}: BottomSheetProps) {
  return (
    <Sheet open={open} onOpenChange={onOpenChange}>
      <SheetContent
        side="bottom"
        className="bg-card border-t border-border rounded-t-2xl px-5 pb-8 safe-bottom max-h-[85vh] overflow-y-auto"
      >
        <div className="w-10 h-1 bg-muted rounded-full mx-auto mb-4 mt-2" />
        <SheetHeader className="mb-4">
          <SheetTitle className="text-foreground">{title}</SheetTitle>
        </SheetHeader>
        {children}
      </SheetContent>
    </Sheet>
  );
}
```

- [ ] **Step 4: 创建 AppLayout**

创建 `packages/web/src/components/layout/AppLayout.tsx`:

```typescript
import { type ReactNode } from "react";
import { Header } from "./Header";
import { TabBar } from "./TabBar";

interface AppLayoutProps {
  children: ReactNode;
  title?: string;
  showBack?: boolean;
}

export function AppLayout({ children, title, showBack }: AppLayoutProps) {
  return (
    <div className="min-h-screen bg-background flex flex-col">
      <Header title={title} showBack={showBack} />
      <main className="flex-1 overflow-y-auto pb-[76px]">{children}</main>
      <TabBar />
    </div>
  );
}
```

- [ ] **Step 5: 提交**

```bash
git add packages/web/src/components/layout/
git commit -m "feat: Layout 组件 — AppLayout + TabBar + Header + BottomSheet

- TabBar: 5 个 Tab + 未读/未付角标
- Header: Dashboard 欢迎语 / 其他页面标题 + 返回按钮
- BottomSheet: shadcn/ui Sheet 封装，用于操作确认
- AppLayout: 组合 Header + Content + TabBar"
```

---

## Task 8: Landing Page

**Files:**
- Modify: `packages/web/src/pages/Landing.tsx`
- Create: `packages/web/src/components/wallet/ConnectButton.tsx`
- Create: `packages/web/src/components/wallet/RegisterSheet.tsx`

- [ ] **Step 1: 创建 ConnectButton**

创建 `packages/web/src/components/wallet/ConnectButton.tsx`:

```typescript
import { ConnectButton as RainbowConnect } from "@rainbow-me/rainbowkit";

export function ConnectButton({ className }: { className?: string }) {
  return (
    <RainbowConnect.Custom>
      {({ openConnectModal, account, mounted }) => {
        const connected = mounted && account;
        return (
          <button
            onClick={openConnectModal}
            className={`bg-primary text-primary-foreground font-semibold rounded-xl ${
              connected ? "hidden" : ""
            } ${className ?? ""}`}
          >
            Connect Wallet
          </button>
        );
      }}
    </RainbowConnect.Custom>
  );
}
```

- [ ] **Step 2: 创建 RegisterSheet**

创建 `packages/web/src/components/wallet/RegisterSheet.tsx`:

```typescript
import { useState } from "react";
import { BottomSheet } from "@/components/layout/BottomSheet";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useRegister, useVerifyEmail } from "@/hooks/useUser";
import { Loader2 } from "lucide-react";

interface RegisterSheetProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  address: string;
  onSuccess: () => void;
}

export function RegisterSheet({
  open,
  onOpenChange,
  address,
  onSuccess,
}: RegisterSheetProps) {
  const [step, setStep] = useState<"email" | "verify">("email");
  const [email, setEmail] = useState("");
  const [code, setCode] = useState("");

  const register = useRegister();
  const verify = useVerifyEmail();

  const handleSendCode = () => {
    if (!email) return;
    // Mock: 直接进入验证步骤
    setStep("verify");
  };

  const handleVerify = async () => {
    const ok = await verify.mutateAsync({ address, code });
    if (ok) {
      await register.mutateAsync({ address, email });
      onSuccess();
      onOpenChange(false);
    }
  };

  return (
    <BottomSheet open={open} onOpenChange={onOpenChange} title="Create Account">
      {step === "email" ? (
        <div className="flex flex-col gap-4">
          <p className="text-sm text-muted-foreground">
            Enter your email to receive billing notifications.
          </p>
          <Input
            type="email"
            placeholder="you@example.com"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            className="bg-muted border-border h-12"
          />
          <Button
            onClick={handleSendCode}
            disabled={!email}
            className="h-12 w-full bg-primary text-primary-foreground rounded-xl"
          >
            Send Verification Code
          </Button>
        </div>
      ) : (
        <div className="flex flex-col gap-4">
          <p className="text-sm text-muted-foreground">
            Enter the verification code sent to {email}
          </p>
          <Input
            type="text"
            placeholder="Enter code"
            value={code}
            onChange={(e) => setCode(e.target.value)}
            className="bg-muted border-border h-12 text-center text-lg tracking-widest"
          />
          <Button
            onClick={handleVerify}
            disabled={!code || register.isPending}
            className="h-12 w-full bg-primary text-primary-foreground rounded-xl"
          >
            {register.isPending ? (
              <Loader2 className="animate-spin" size={18} />
            ) : (
              "Verify & Register"
            )}
          </Button>
        </div>
      )}
    </BottomSheet>
  );
}
```

- [ ] **Step 3: 实现 Landing 页面**

替换 `packages/web/src/pages/Landing.tsx` 完整内容:

```typescript
import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { useAccount } from "wagmi";
import { ConnectButton } from "@/components/wallet/ConnectButton";
import { RegisterSheet } from "@/components/wallet/RegisterSheet";
import { useUserProfile } from "@/hooks/useUser";

export default function Landing() {
  const navigate = useNavigate();
  const { address, isConnected } = useAccount();
  const { data: user } = useUserProfile(address);
  const [showRegister, setShowRegister] = useState(false);

  useEffect(() => {
    if (isConnected && user) {
      navigate("/dashboard", { replace: true });
    } else if (isConnected && user === null) {
      setShowRegister(true);
    }
  }, [isConnected, user, navigate]);

  return (
    <div className="min-h-screen bg-gradient-to-b from-background to-[#0f1a2e] flex flex-col">
      {/* Nav */}
      <div className="flex justify-between items-center px-5 pt-[env(safe-area-inset-top)] pb-2">
        <div className="pt-3">
          <span className="text-lg font-bold text-white">Link</span>
          <span className="text-lg font-bold text-primary">World</span>
        </div>
        <div className="pt-3">
          <ConnectButton className="px-4 py-2 text-sm" />
        </div>
      </div>

      {/* Hero */}
      <div className="flex-1 flex flex-col items-center justify-center text-center px-6">
        <p className="text-[10px] text-primary tracking-[2px] uppercase mb-3">
          Powered by 0G Chain & Chainlink
        </p>
        <h1 className="text-[28px] font-bold text-white leading-tight">
          Global Seamless
          <br />
          Communication
        </h1>
        <p className="text-[14px] text-muted-foreground mt-4 leading-relaxed">
          No KYC. No bank account.
          <br />
          Just your wallet.
        </p>
        <div className="w-full mt-7 flex flex-col gap-3">
          <ConnectButton className="w-full py-3.5 text-[15px]" />
          <button className="w-full py-3.5 text-[15px] border border-border text-foreground/80 rounded-xl">
            Learn More
          </button>
        </div>
      </div>

      {/* Stats */}
      <div className="flex justify-around py-5 border-t border-white/5">
        <div className="text-center">
          <div className="text-xl font-bold text-primary">50+</div>
          <div className="text-[10px] text-muted-foreground">Countries</div>
        </div>
        <div className="text-center">
          <div className="text-xl font-bold text-primary">2.5%</div>
          <div className="text-[10px] text-muted-foreground">Fee</div>
        </div>
        <div className="text-center">
          <div className="text-xl font-bold text-primary">0</div>
          <div className="text-[10px] text-muted-foreground">KYC</div>
        </div>
      </div>

      {/* Register Sheet */}
      {address && (
        <RegisterSheet
          open={showRegister}
          onOpenChange={setShowRegister}
          address={address}
          onSuccess={() => navigate("/dashboard", { replace: true })}
        />
      )}
    </div>
  );
}
```

- [ ] **Step 4: 验证 Landing 页面**

```bash
cd packages/web && pnpm dev
```

访问 http://localhost:3000 — 应看到暗色背景的 Landing Page，含 Connect Wallet 按钮。

- [ ] **Step 5: 提交**

```bash
git add packages/web/src/pages/Landing.tsx packages/web/src/components/wallet/
git commit -m "feat: Landing Page + 钱包连接 + 注册流程

- Landing 页：Hero + CTA + 统计数据
- ConnectButton: RainbowKit Custom 封装
- RegisterSheet: 邮箱 + 验证码两步注册"
```

---

## Task 9: Dashboard 页面

**Files:**
- Modify: `packages/web/src/pages/Dashboard.tsx`

- [ ] **Step 1: 实现 Dashboard**

替换 `packages/web/src/pages/Dashboard.tsx` 完整内容:

```typescript
import { useNavigate } from "react-router-dom";
import { AppLayout } from "@/components/layout/AppLayout";
import { useUserContext } from "@/contexts/UserContext";
import { useDepositInfo } from "@/hooks/useDeposit";
import { useMonthEstimate } from "@/hooks/useBilling";
import { useMyNumbers } from "@/hooks/useOperator";
import { formatToken, formatUSD } from "@/utils/format";
import { Wallet, Globe, FileText, Zap } from "lucide-react";

export default function Dashboard() {
  const navigate = useNavigate();
  const { address, status } = useUserContext();
  const { data: deposit } = useDepositInfo(address);
  const { data: estimate } = useMonthEstimate(address);
  const { data: numbers } = useMyNumbers(address);

  const activeNumber = numbers?.find((n) => n.status === "active");

  return (
    <AppLayout>
      <div className="px-4 flex flex-col gap-3">
        {/* Status Card */}
        <div className="bg-gradient-to-br from-[#1a2a4a] to-[#1a1a3a] rounded-2xl p-5">
          <div className="flex justify-between items-start mb-4">
            <div>
              <div className="text-[11px] text-muted-foreground">
                Account Status
              </div>
              <div className="flex items-center gap-1.5 mt-1">
                <div
                  className={`w-2 h-2 rounded-full ${
                    status === "active"
                      ? "bg-success"
                      : status === "suspended"
                      ? "bg-danger"
                      : "bg-warning"
                  }`}
                />
                <span
                  className={`text-[15px] font-semibold capitalize ${
                    status === "active"
                      ? "text-success"
                      : status === "suspended"
                      ? "text-danger"
                      : "text-warning"
                  }`}
                >
                  {status}
                </span>
              </div>
            </div>
            <div className="text-right">
              <div className="text-[11px] text-muted-foreground">Deposit</div>
              <div className="text-[15px] font-semibold text-warning mt-1">
                {deposit ? formatToken(deposit.balance, deposit.currency) : "—"}
              </div>
            </div>
          </div>
          {activeNumber && (
            <div className="border-t border-white/[0.08] pt-3 flex justify-between">
              <div>
                <div className="text-[10px] text-muted-foreground">
                  Virtual Number
                </div>
                <div className="text-[13px] text-primary mt-0.5">
                  {activeNumber.number}
                </div>
              </div>
              <div className="text-right">
                <div className="text-[10px] text-muted-foreground">Region</div>
                <div className="text-[13px] text-foreground mt-0.5">
                  {activeNumber.regionFlag} {activeNumber.regionName}
                </div>
              </div>
            </div>
          )}
        </div>

        {/* This Month */}
        <div className="bg-card border border-border rounded-xl p-4">
          <div className="flex justify-between items-center mb-3">
            <span className="text-[13px] font-semibold text-foreground">
              This Month
            </span>
            <span className="text-[11px] text-muted-foreground">
              {new Date().toLocaleDateString("en-US", { month: "long" })}
            </span>
          </div>
          <div className="flex gap-2.5">
            <div className="flex-1 bg-muted rounded-lg p-2.5 text-center">
              <div className="text-base font-semibold text-primary">
                {estimate?.dataGB ?? "—"} GB
              </div>
              <div className="text-[10px] text-muted-foreground mt-0.5">
                Data Used
              </div>
            </div>
            <div className="flex-1 bg-muted rounded-lg p-2.5 text-center">
              <div className="text-base font-semibold text-accent">
                {estimate?.callMinutes ?? "—"} min
              </div>
              <div className="text-[10px] text-muted-foreground mt-0.5">
                Calls
              </div>
            </div>
            <div className="flex-1 bg-muted rounded-lg p-2.5 text-center">
              <div className="text-base font-semibold text-foreground">
                {estimate ? formatUSD(estimate.estimatedCost) : "—"}
              </div>
              <div className="text-[10px] text-muted-foreground mt-0.5">
                Est. Bill
              </div>
            </div>
          </div>
        </div>

        {/* Quick Actions */}
        <div className="grid grid-cols-2 gap-2">
          {[
            { icon: Wallet, label: "Top Up", path: "/deposit" },
            { icon: Globe, label: "Switch Region", path: "/services" },
            { icon: FileText, label: "Bills", path: "/billing" },
            { icon: Zap, label: "Pay Now", path: "/billing" },
          ].map(({ icon: Icon, label, path }) => (
            <button
              key={label}
              onClick={() => navigate(path)}
              className="bg-card border border-border rounded-xl p-4 flex flex-col items-center gap-1.5"
            >
              <Icon size={22} className="text-primary" />
              <span className="text-[11px] text-primary">{label}</span>
            </button>
          ))}
        </div>
      </div>
    </AppLayout>
  );
}
```

- [ ] **Step 2: 提交**

```bash
git add packages/web/src/pages/Dashboard.tsx
git commit -m "feat: Dashboard 页面

- 状态大卡片：账户状态 + 押金余额 + 虚拟号码 + 地区
- 本月用量：Data / Calls / Est. Bill 三栏
- Quick Actions 2x2 网格"
```

---

## Task 10: Deposit 页面

**Files:**
- Modify: `packages/web/src/pages/Deposit.tsx`

- [ ] **Step 1: 实现 Deposit 页面**

替换 `packages/web/src/pages/Deposit.tsx` 完整内容:

```typescript
import { useState } from "react";
import { AppLayout } from "@/components/layout/AppLayout";
import { BottomSheet } from "@/components/layout/BottomSheet";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useUserContext } from "@/contexts/UserContext";
import { useDepositInfo, useDepositHistory, useDoDeposit, useWithdraw } from "@/hooks/useDeposit";
import { formatRelativeTime } from "@/utils/format";
import { Loader2 } from "lucide-react";

export default function Deposit() {
  const { address } = useUserContext();
  const { data: deposit } = useDepositInfo(address);
  const { data: history } = useDepositHistory(address);
  const doDeposit = useDoDeposit();
  const doWithdraw = useWithdraw();

  const [sheetType, setSheetType] = useState<"deposit" | "withdraw" | null>(null);
  const [amount, setAmount] = useState("");
  const [currency, setCurrency] = useState("USDT");

  const balanceNum = parseFloat(deposit?.balance ?? "0");
  const minNum = parseFloat(deposit?.minimumRequired ?? "20");
  const pct = Math.min(Math.round((balanceNum / minNum) * 100), 999);

  const handleConfirm = async () => {
    if (!address || !amount) return;
    if (sheetType === "deposit") {
      await doDeposit.mutateAsync({ address, amount, currency });
    } else {
      await doWithdraw.mutateAsync({ address, amount });
    }
    setSheetType(null);
    setAmount("");
  };

  const isPending = doDeposit.isPending || doWithdraw.isPending;

  return (
    <AppLayout title="Deposit">
      <div className="px-4 flex flex-col gap-4">
        {/* Balance Card */}
        <div className="bg-gradient-to-br from-[#2d1f0a] to-[#1a1a0a] border border-warning/20 rounded-2xl p-5 text-center">
          <div className="text-[11px] text-muted-foreground">
            Available Balance
          </div>
          <div className="text-[32px] font-bold text-warning mt-2">
            {deposit?.balance ?? "—"}
          </div>
          <div className="text-[12px] text-muted-foreground">
            {deposit?.currency ?? "USDT"}
          </div>
          <div className="text-[10px] text-muted-foreground mt-2">
            Minimum required: {deposit?.minimumRequired ?? "20.00"} USDT
          </div>
          {/* Progress bar */}
          <div className="bg-muted rounded h-1.5 mt-3 overflow-hidden">
            <div
              className="h-full rounded bg-gradient-to-r from-warning to-success transition-all"
              style={{ width: `${Math.min(pct, 100)}%` }}
            />
          </div>
          <div className="text-[9px] text-success mt-1">{pct}% of minimum</div>
        </div>

        {/* Action Buttons */}
        <div className="flex gap-3">
          <Button
            onClick={() => { setSheetType("deposit"); setAmount(""); }}
            className="flex-1 h-12 bg-success hover:bg-success/90 text-white rounded-xl font-semibold"
          >
            Deposit
          </Button>
          <Button
            onClick={() => { setSheetType("withdraw"); setAmount(""); }}
            variant="outline"
            className="flex-1 h-12 border-border rounded-xl font-semibold"
          >
            Withdraw
          </Button>
        </div>

        {/* History */}
        <div>
          <h3 className="text-[11px] text-muted-foreground uppercase tracking-wider mb-2">
            History
          </h3>
          <div className="bg-card rounded-xl overflow-hidden">
            {history?.map((record, i) => (
              <div
                key={record.id}
                className={`flex justify-between items-center px-4 py-3 ${
                  i < (history.length - 1) ? "border-b border-border" : ""
                }`}
              >
                <div>
                  <div className="text-[12px] text-foreground capitalize">
                    {record.type === "deduction" ? "Bill Deduction" : record.type}
                  </div>
                  <div className="text-[10px] text-muted-foreground">
                    {formatRelativeTime(record.timestamp)}
                  </div>
                </div>
                <div
                  className={`font-semibold ${
                    record.type === "deposit" ? "text-success" : "text-danger"
                  }`}
                >
                  {record.type === "deposit" ? "+" : "-"}
                  {record.amount} {record.currency}
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>

      {/* Sheet */}
      <BottomSheet
        open={sheetType !== null}
        onOpenChange={(open) => !open && setSheetType(null)}
        title={sheetType === "deposit" ? "Deposit Funds" : "Withdraw Funds"}
      >
        <div className="flex flex-col gap-4">
          {sheetType === "deposit" && (
            <div className="flex gap-2">
              {["USDT", "ETH"].map((c) => (
                <button
                  key={c}
                  onClick={() => setCurrency(c)}
                  className={`flex-1 py-2 rounded-lg text-sm font-medium ${
                    currency === c
                      ? "bg-primary text-white"
                      : "bg-muted text-muted-foreground"
                  }`}
                >
                  {c}
                </button>
              ))}
            </div>
          )}
          <Input
            type="number"
            placeholder="Enter amount"
            value={amount}
            onChange={(e) => setAmount(e.target.value)}
            className="bg-muted border-border h-12 text-lg"
          />
          {sheetType === "withdraw" && (
            <p className="text-xs text-muted-foreground">
              Max withdrawal: {deposit?.balance} {deposit?.currency}. All bills
              must be settled first.
            </p>
          )}
          <Button
            onClick={handleConfirm}
            disabled={!amount || isPending}
            className="h-12 w-full bg-primary text-white rounded-xl font-semibold"
          >
            {isPending ? (
              <Loader2 className="animate-spin" size={18} />
            ) : (
              `Confirm ${sheetType === "deposit" ? "Deposit" : "Withdrawal"}`
            )}
          </Button>
        </div>
      </BottomSheet>
    </AppLayout>
  );
}
```

- [ ] **Step 2: 提交**

```bash
git add packages/web/src/pages/Deposit.tsx
git commit -m "feat: Deposit 页面

- 余额大卡片 + 进度条
- 存入/提取双按钮 + BottomSheet 操作
- 历史记录列表"
```

---

## Task 11: Services + RegionDetail 页面

**Files:**
- Modify: `packages/web/src/pages/Services.tsx`
- Modify: `packages/web/src/pages/RegionDetail.tsx`

- [ ] **Step 1: 实现 Services 页面**

替换 `packages/web/src/pages/Services.tsx` 完整内容:

```typescript
import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { AppLayout } from "@/components/layout/AppLayout";
import { useRegions } from "@/hooks/useOperator";
import { useMyNumbers } from "@/hooks/useOperator";
import { useUserContext } from "@/contexts/UserContext";
import { Search, ChevronRight } from "lucide-react";

export default function Services() {
  const navigate = useNavigate();
  const { address } = useUserContext();
  const { data: regions } = useRegions();
  const { data: numbers } = useMyNumbers(address);
  const [search, setSearch] = useState("");

  const activeNumber = numbers?.find((n) => n.status === "active");
  const filtered = regions?.filter(
    (r) =>
      r.name.toLowerCase().includes(search.toLowerCase()) ||
      r.code.toLowerCase().includes(search.toLowerCase())
  );

  return (
    <AppLayout title="Services">
      <div className="px-4 flex flex-col gap-4">
        {/* Search */}
        <div className="bg-card border border-border rounded-xl px-3.5 py-2.5 flex items-center gap-2">
          <Search size={16} className="text-muted-foreground" />
          <input
            type="text"
            placeholder="Search country or region..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="flex-1 bg-transparent text-sm text-foreground placeholder:text-muted-foreground outline-none"
          />
        </div>

        {/* Current Number */}
        {activeNumber && (
          <div className="bg-gradient-to-br from-[#1a2a4a] to-[#1a1a3a] rounded-xl p-4">
            <div className="text-[10px] text-muted-foreground mb-1.5">
              Current Number
            </div>
            <div className="flex justify-between items-center">
              <div>
                <div className="text-[15px] font-semibold text-primary">
                  {activeNumber.number}
                </div>
                <div className="text-[10px] text-muted-foreground mt-0.5">
                  {activeNumber.regionFlag} {activeNumber.regionName} ·{" "}
                  {activeNumber.operatorName} · Active
                </div>
              </div>
              <button className="bg-muted rounded-lg px-3 py-1.5 text-[10px] text-primary">
                Manage
              </button>
            </div>
          </div>
        )}

        {/* Region List */}
        <div>
          <h3 className="text-[11px] text-muted-foreground uppercase tracking-wider mb-2">
            {search ? "Search Results" : "Popular Regions"}
          </h3>
          <div className="flex flex-col gap-2">
            {filtered?.map((region) => (
              <button
                key={region.code}
                onClick={() => navigate(`/services/${region.code}`)}
                className="bg-card border border-border rounded-xl p-4 flex justify-between items-center text-left"
              >
                <div className="flex items-center gap-3">
                  <span className="text-[22px]">{region.flag}</span>
                  <div>
                    <div className="text-[13px] text-foreground font-medium">
                      {region.name}
                    </div>
                    <div className="text-[10px] text-muted-foreground">
                      {region.operatorCount} operators · from $
                      {region.startingPrice}/mo
                    </div>
                  </div>
                </div>
                <ChevronRight size={16} className="text-muted-foreground" />
              </button>
            ))}
          </div>
        </div>
      </div>
    </AppLayout>
  );
}
```

- [ ] **Step 2: 实现 RegionDetail 页面**

替换 `packages/web/src/pages/RegionDetail.tsx` 完整内容:

```typescript
import { useState } from "react";
import { useParams } from "react-router-dom";
import { AppLayout } from "@/components/layout/AppLayout";
import { BottomSheet } from "@/components/layout/BottomSheet";
import { Button } from "@/components/ui/button";
import { useOperatorsByRegion, useApplyNumber } from "@/hooks/useOperator";
import { useRegions } from "@/hooks/useOperator";
import { useUserContext } from "@/contexts/UserContext";
import { formatUSD } from "@/utils/format";
import { Loader2 } from "lucide-react";
import type { Operator } from "@/types";

export default function RegionDetail() {
  const { regionCode } = useParams<{ regionCode: string }>();
  const { address } = useUserContext();
  const { data: regions } = useRegions();
  const { data: operators } = useOperatorsByRegion(regionCode);
  const applyNumber = useApplyNumber();

  const [selected, setSelected] = useState<Operator | null>(null);

  const region = regions?.find((r) => r.code === regionCode);

  const handleApply = async () => {
    if (!address || !selected) return;
    await applyNumber.mutateAsync({ address, operatorId: selected.id });
    setSelected(null);
  };

  return (
    <AppLayout title={region ? `${region.flag} ${region.name}` : "Region"} showBack>
      <div className="px-4 flex flex-col gap-3">
        {operators?.map((op) => (
          <div
            key={op.id}
            className="bg-card border border-border rounded-xl p-4"
          >
            <div className="flex justify-between items-start mb-3">
              <div className="text-[14px] font-semibold text-foreground">
                {op.name}
              </div>
              {op.isActive && (
                <span className="text-[10px] bg-success/10 text-success px-2 py-0.5 rounded">
                  Available
                </span>
              )}
            </div>
            <div className="flex gap-4 text-[11px] text-muted-foreground mb-3">
              <span>Data: {formatUSD(op.dataRate)}/GB</span>
              <span>Calls: {formatUSD(op.callRate)}/min</span>
            </div>
            <div className="flex justify-between items-center pt-3 border-t border-border">
              <div className="text-[11px] text-muted-foreground">
                Min deposit: {op.requiredDeposit} USDT
              </div>
              <Button
                onClick={() => setSelected(op)}
                size="sm"
                className="bg-primary text-white rounded-lg text-xs h-8 px-4"
              >
                Get Number
              </Button>
            </div>
          </div>
        ))}
      </div>

      {/* Apply Sheet */}
      <BottomSheet
        open={selected !== null}
        onOpenChange={(open) => !open && setSelected(null)}
        title="Apply for Virtual Number"
      >
        {selected && (
          <div className="flex flex-col gap-4">
            <div className="bg-muted rounded-xl p-4 space-y-2">
              <div className="flex justify-between text-sm">
                <span className="text-muted-foreground">Region</span>
                <span className="text-foreground">
                  {region?.flag} {region?.name}
                </span>
              </div>
              <div className="flex justify-between text-sm">
                <span className="text-muted-foreground">Operator</span>
                <span className="text-foreground">{selected.name}</span>
              </div>
              <div className="flex justify-between text-sm">
                <span className="text-muted-foreground">Required Deposit</span>
                <span className="text-warning">
                  {selected.requiredDeposit} USDT
                </span>
              </div>
            </div>
            <Button
              onClick={handleApply}
              disabled={applyNumber.isPending}
              className="h-12 w-full bg-primary text-white rounded-xl font-semibold"
            >
              {applyNumber.isPending ? (
                <Loader2 className="animate-spin" size={18} />
              ) : (
                "Confirm & Get Number"
              )}
            </Button>
          </div>
        )}
      </BottomSheet>
    </AppLayout>
  );
}
```

- [ ] **Step 3: 提交**

```bash
git add packages/web/src/pages/Services.tsx packages/web/src/pages/RegionDetail.tsx
git commit -m "feat: Services + RegionDetail 页面

- Services: 搜索栏 + 当前号码卡片 + 国家列表
- RegionDetail: 运营商列表 + 号码申请 Sheet"
```

---

## Task 12: Billing + BillDetail 页面

**Files:**
- Modify: `packages/web/src/pages/Billing.tsx`
- Modify: `packages/web/src/pages/BillDetail.tsx`

- [ ] **Step 1: 实现 Billing 页面**

替换 `packages/web/src/pages/Billing.tsx` 完整内容:

```typescript
import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { AppLayout } from "@/components/layout/AppLayout";
import { useUserContext } from "@/contexts/UserContext";
import { useBills } from "@/hooks/useBilling";
import { formatUSD, formatMonth } from "@/utils/format";
import { AlertTriangle } from "lucide-react";

export default function Billing() {
  const navigate = useNavigate();
  const { address } = useUserContext();
  const [tab, setTab] = useState<"unpaid" | "paid">("unpaid");
  const { data: bills } = useBills(address, tab === "unpaid" ? "unpaid" : "paid");

  const unpaidCount = bills?.filter((b) => b.status === "unpaid").length ?? 0;

  return (
    <AppLayout title="Bills">
      <div className="px-4 flex flex-col gap-4">
        {/* Tab Switch */}
        <div className="flex bg-card rounded-xl p-1">
          {(["unpaid", "paid"] as const).map((t) => (
            <button
              key={t}
              onClick={() => setTab(t)}
              className={`flex-1 py-2 rounded-lg text-[12px] font-medium capitalize ${
                tab === t
                  ? "bg-muted text-foreground"
                  : "text-muted-foreground"
              }`}
            >
              {t === "unpaid" ? "Unpaid" : "History"}
            </button>
          ))}
        </div>

        {/* Unpaid Alert */}
        {tab === "unpaid" && unpaidCount > 0 && (
          <div className="bg-danger/10 border border-danger/20 rounded-xl p-3.5 flex items-center gap-3">
            <AlertTriangle size={20} className="text-danger" />
            <div>
              <div className="text-[12px] text-danger font-medium">
                {unpaidCount} unpaid bill{unpaidCount > 1 ? "s" : ""}
              </div>
              <div className="text-[10px] text-muted-foreground">
                Pay before due date to avoid service suspension
              </div>
            </div>
          </div>
        )}

        {/* Bill List */}
        <div className="flex flex-col gap-3">
          {bills?.map((bill) => (
            <button
              key={bill.id}
              onClick={() => navigate(`/billing/${bill.id}`)}
              className="bg-card border border-border rounded-xl p-4 text-left"
            >
              <div className="flex justify-between items-center mb-2.5">
                <div className="text-[13px] font-semibold text-foreground">
                  {formatMonth(bill.month)}
                </div>
                <span
                  className={`text-[10px] font-medium px-2.5 py-0.5 rounded ${
                    bill.status === "unpaid"
                      ? "bg-danger/10 text-danger"
                      : "bg-success/10 text-success"
                  }`}
                >
                  {bill.status === "unpaid" ? "Unpaid" : "Paid"}
                </span>
              </div>
              <div className="flex justify-between text-[11px] text-muted-foreground mb-1">
                <span>Operator Fee</span>
                <span className="text-foreground">
                  {formatUSD(bill.operatorFee)}
                </span>
              </div>
              <div className="flex justify-between text-[11px] text-muted-foreground mb-2">
                <span>Platform Fee (2.5%)</span>
                <span className="text-foreground">
                  {formatUSD(bill.platformFee)}
                </span>
              </div>
              <div className="border-t border-border pt-2 flex justify-between items-center">
                <div className="text-[15px] font-bold text-foreground">
                  {formatUSD(bill.totalAmount)}
                </div>
                {bill.status === "paid" && bill.paidAt && (
                  <div className="text-[10px] text-muted-foreground">
                    Paid{" "}
                    {new Date(bill.paidAt).toLocaleDateString("en-US", {
                      month: "short",
                      day: "numeric",
                    })}
                  </div>
                )}
              </div>
            </button>
          ))}
        </div>
      </div>
    </AppLayout>
  );
}
```

- [ ] **Step 2: 实现 BillDetail 页面**

替换 `packages/web/src/pages/BillDetail.tsx` 完整内容:

```typescript
import { useParams } from "react-router-dom";
import { AppLayout } from "@/components/layout/AppLayout";
import { BottomSheet } from "@/components/layout/BottomSheet";
import { Button } from "@/components/ui/button";
import { useBillDetail, usePayBill } from "@/hooks/useBilling";
import { formatUSD, formatMonth } from "@/utils/format";
import { Loader2 } from "lucide-react";
import { useState } from "react";

export default function BillDetail() {
  const { billId } = useParams<{ billId: string }>();
  const { data: bill } = useBillDetail(billId);
  const payBill = usePayBill();
  const [showPay, setShowPay] = useState(false);

  const handlePay = async () => {
    if (!billId) return;
    await payBill.mutateAsync(billId);
    setShowPay(false);
  };

  if (!bill) {
    return (
      <AppLayout title="Bill" showBack>
        <div className="flex items-center justify-center h-60">
          <div className="animate-spin w-6 h-6 border-2 border-primary border-t-transparent rounded-full" />
        </div>
      </AppLayout>
    );
  }

  return (
    <AppLayout title={formatMonth(bill.month)} showBack>
      <div className="px-4 flex flex-col gap-4">
        {/* Status */}
        <div className="bg-card border border-border rounded-xl p-5 text-center">
          <span
            className={`text-[12px] font-medium px-3 py-1 rounded ${
              bill.status === "unpaid"
                ? "bg-danger/10 text-danger"
                : "bg-success/10 text-success"
            }`}
          >
            {bill.status === "unpaid" ? "Unpaid" : "Paid"}
          </span>
          <div className="text-[32px] font-bold text-foreground mt-3">
            {formatUSD(bill.totalAmount)}
          </div>
          {bill.status === "unpaid" && (
            <div className="text-[11px] text-muted-foreground mt-1">
              Due by{" "}
              {new Date(bill.dueDate).toLocaleDateString("en-US", {
                month: "long",
                day: "numeric",
                year: "numeric",
              })}
            </div>
          )}
        </div>

        {/* Breakdown */}
        <div className="bg-card border border-border rounded-xl p-4">
          <h3 className="text-[13px] font-semibold text-foreground mb-3">
            Fee Breakdown
          </h3>
          <div className="space-y-2">
            <div className="flex justify-between text-[12px]">
              <span className="text-muted-foreground">Operator Fee</span>
              <span className="text-foreground">{formatUSD(bill.operatorFee)}</span>
            </div>
            <div className="flex justify-between text-[12px]">
              <span className="text-muted-foreground">Platform Fee (2.5%)</span>
              <span className="text-foreground">{formatUSD(bill.platformFee)}</span>
            </div>
            <div className="border-t border-border pt-2 flex justify-between text-[13px] font-semibold">
              <span className="text-foreground">Total</span>
              <span className="text-foreground">{formatUSD(bill.totalAmount)}</span>
            </div>
          </div>
        </div>

        {/* Usage */}
        <div className="bg-card border border-border rounded-xl p-4">
          <h3 className="text-[13px] font-semibold text-foreground mb-3">
            Usage Summary
          </h3>
          <div className="flex gap-3">
            <div className="flex-1 bg-muted rounded-lg p-3 text-center">
              <div className="text-base font-semibold text-primary">
                {bill.usage.dataGB} GB
              </div>
              <div className="text-[10px] text-muted-foreground mt-0.5">
                Data
              </div>
            </div>
            <div className="flex-1 bg-muted rounded-lg p-3 text-center">
              <div className="text-base font-semibold text-accent">
                {bill.usage.callMinutes} min
              </div>
              <div className="text-[10px] text-muted-foreground mt-0.5">
                Calls
              </div>
            </div>
          </div>
        </div>

        {/* Pay Button */}
        {bill.status === "unpaid" && (
          <Button
            onClick={() => setShowPay(true)}
            className="h-12 w-full bg-primary text-white rounded-xl font-semibold"
          >
            Pay {formatUSD(bill.totalAmount)}
          </Button>
        )}
      </div>

      {/* Pay Sheet */}
      <BottomSheet
        open={showPay}
        onOpenChange={setShowPay}
        title="Confirm Payment"
      >
        <div className="flex flex-col gap-4">
          <div className="bg-muted rounded-xl p-4 space-y-2">
            <div className="flex justify-between text-sm">
              <span className="text-muted-foreground">Amount</span>
              <span className="text-foreground font-semibold">
                {formatUSD(bill.totalAmount)}
              </span>
            </div>
            <div className="flex justify-between text-sm">
              <span className="text-muted-foreground">Payment Source</span>
              <span className="text-foreground">Deposit Balance</span>
            </div>
          </div>
          <Button
            onClick={handlePay}
            disabled={payBill.isPending}
            className="h-12 w-full bg-primary text-white rounded-xl font-semibold"
          >
            {payBill.isPending ? (
              <Loader2 className="animate-spin" size={18} />
            ) : (
              "Confirm Payment"
            )}
          </Button>
        </div>
      </BottomSheet>
    </AppLayout>
  );
}
```

- [ ] **Step 3: 提交**

```bash
git add packages/web/src/pages/Billing.tsx packages/web/src/pages/BillDetail.tsx
git commit -m "feat: Billing + BillDetail 页面

- Billing: Unpaid/History 分段 + 未付告警 + 账单卡片列表
- BillDetail: 费用明细 + 用量统计 + 支付确认 Sheet"
```

---

## Task 13: Notifications 页面

**Files:**
- Modify: `packages/web/src/pages/Notifications.tsx`

- [ ] **Step 1: 实现 Notifications 页面**

替换 `packages/web/src/pages/Notifications.tsx` 完整内容:

```typescript
import { AppLayout } from "@/components/layout/AppLayout";
import { useUserContext } from "@/contexts/UserContext";
import {
  useNotifications,
  useMarkAsRead,
  useMarkAllAsRead,
} from "@/hooks/useNotification";
import { formatRelativeTime } from "@/utils/format";
import { NOTIFICATION_TYPE_COLORS } from "@/config/constants";

export default function Notifications() {
  const { address } = useUserContext();
  const { data: notifications } = useNotifications(address);
  const markAsRead = useMarkAsRead();
  const markAllAsRead = useMarkAllAsRead();

  const unread = notifications?.filter((n) => !n.read) ?? [];
  const read = notifications?.filter((n) => n.read) ?? [];

  const handleMarkAllRead = () => {
    if (address) markAllAsRead.mutate(address);
  };

  const handleTap = (id: string, isRead: boolean) => {
    if (!isRead) markAsRead.mutate(id);
  };

  return (
    <AppLayout title="Notifications">
      <div className="px-4 flex flex-col gap-4">
        {/* Header action */}
        {unread.length > 0 && (
          <div className="flex justify-end">
            <button
              onClick={handleMarkAllRead}
              className="text-[11px] text-primary"
            >
              Mark all read
            </button>
          </div>
        )}

        {/* Unread */}
        {unread.length > 0 && (
          <div>
            <h3 className="text-[10px] text-muted-foreground uppercase tracking-wider mb-2">
              New
            </h3>
            <div className="flex flex-col gap-2">
              {unread.map((n) => (
                <button
                  key={n.id}
                  onClick={() => handleTap(n.id, n.read)}
                  className="bg-[#0f1a2e] border border-primary/20 rounded-xl p-3.5 text-left"
                  style={{
                    borderLeftWidth: 3,
                    borderLeftColor:
                      NOTIFICATION_TYPE_COLORS[n.type] ?? "#888",
                  }}
                >
                  <div className="flex justify-between items-start">
                    <div className="flex-1">
                      <div className="text-[12px] text-foreground font-medium">
                        {n.title}
                      </div>
                      <div className="text-[11px] text-muted-foreground mt-1 leading-relaxed">
                        {n.message}
                      </div>
                    </div>
                    <div className="w-2 h-2 rounded-full bg-primary mt-1 ml-2 flex-shrink-0" />
                  </div>
                  <div className="text-[10px] text-muted-foreground/60 mt-2">
                    {formatRelativeTime(n.createdAt)}
                  </div>
                </button>
              ))}
            </div>
          </div>
        )}

        {/* Read */}
        {read.length > 0 && (
          <div>
            <h3 className="text-[10px] text-muted-foreground uppercase tracking-wider mb-2">
              Earlier
            </h3>
            <div className="flex flex-col gap-2">
              {read.map((n) => (
                <div
                  key={n.id}
                  className="bg-card border border-border rounded-xl p-3.5"
                >
                  <div className="text-[12px] text-muted-foreground font-medium">
                    {n.title}
                  </div>
                  <div className="text-[11px] text-muted-foreground/60 mt-1 leading-relaxed">
                    {n.message}
                  </div>
                  <div className="text-[10px] text-muted-foreground/40 mt-2">
                    {formatRelativeTime(n.createdAt)}
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}

        {/* Empty state */}
        {notifications?.length === 0 && (
          <div className="flex flex-col items-center justify-center py-20 text-muted-foreground">
            <div className="text-4xl mb-3">🔔</div>
            <div className="text-sm">No notifications yet</div>
          </div>
        )}
      </div>
    </AppLayout>
  );
}
```

- [ ] **Step 2: 提交**

```bash
git add packages/web/src/pages/Notifications.tsx
git commit -m "feat: Notifications 页面

- New/Earlier 分组显示
- 未读蓝点 + 左侧彩色边条区分通知类型
- Mark all read + 点击标记已读"
```

---

## Task 14: 最终集成验证

- [ ] **Step 1: 确保所有文件存在且无 TypeScript 错误**

```bash
cd packages/web && pnpm build
```

修复所有 TS 编译错误。

- [ ] **Step 2: 启动开发服务器，逐页验证**

```bash
cd packages/web && pnpm dev
```

验证清单:
1. http://localhost:3000 → Landing Page 正常渲染
2. Connect Wallet → RainbowKit 弹窗弹出
3. 注册后跳转 Dashboard → 状态卡片 + 用量 + Quick Actions 正常
4. Tab Bar → 5 个 Tab 切换正常
5. Deposit → 余额卡片 + 存入/提取 Sheet 正常
6. Services → 地区列表 + 搜索 + 点击进入 RegionDetail
7. RegionDetail → 运营商列表 + 号码申请 Sheet
8. Billing → Unpaid/History 切换 + 账单卡片
9. BillDetail → 费用明细 + 支付 Sheet
10. Notifications → 未读/已读分组 + Mark all read

- [ ] **Step 3: 提交集成修复（如有）**

```bash
git add -A packages/web/
git commit -m "fix: 集成验证修复"
```

- [ ] **Step 4: 更新 pipeline 状态**

将 `.claude/pipeline.json` 中 `plan` 阶段标记为 `completed`。
