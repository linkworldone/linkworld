# LinkWorld Frontend Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Build the complete LinkWorld frontend web app — 6 pages, 5 mock services, wallet integration, dark theme mobile-first UI.

**Architecture:** Three-layer data flow (Pages → React Query Hooks → Mock Services). Real RainbowKit wallet connection + mock business logic. AppLayout with route guards based on user status state machine (UNREGISTERED/INACTIVE/ACTIVE/SUSPENDED).

**Tech Stack:** React 19, Vite 6, TailwindCSS 3.4, shadcn/ui (selective), wagmi 2.14, viem 2.21, RainbowKit 2.2, TanStack React Query 5.62, vaul (drawer)

---

## Arch-Review Risk Items (addressed in tasks)

| Risk | Addressed In |
|------|-------------|
| Tailwind 配色需 semantic tokens | Task 0 Step 2 |
| BrowserRouter 位置需从 main.tsx 移除 | Task 0 Step 8 |
| shadcn/ui 初始化需前置 | Task 0 Step 3 |
| 桌面端需 max-w 容器约束 | Task 3 (AppLayout) |
| Sheet 拖拽手势需 vaul 库 | Task 0 Step 1 |
| Mock Service 需加延迟模拟 | Task 1 (delay helper) |
| Bill 金额精度问题 | Task 1 (use string for amounts) |
| 多虚拟号码场景 | Task 7 (Services My Numbers tab) |
| 提前结账入口 | Task 8 (Billing) |
| 逾期扣款倒计时警告 | Task 5 (Dashboard warning) |

---

### Task 0: Infrastructure Setup

**Files:**
- Modify: `packages/web/package.json`
- Modify: `packages/web/tailwind.config.ts`
- Modify: `packages/web/src/index.css`
- Modify: `packages/web/src/main.tsx`
- Create: `packages/web/src/lib/utils.ts`
- Create: `packages/web/src/config/wagmi.ts`
- Create: `packages/web/src/config/constants.ts`
- Create: `packages/web/components.json`

- [ ] **Step 1: Install vaul + clsx + tailwind-merge**
- [ ] **Step 2: Update tailwind.config.ts with semantic color tokens** (surface, brand, status, text, border, maxWidth.mobile=430px)
- [ ] **Step 3: Initialize shadcn/ui + install Button, Badge, Tabs**
- [ ] **Step 4: Verify vaul installed**
- [ ] **Step 5: Update index.css with dark theme CSS variables for shadcn**
- [ ] **Step 6: Create wagmi config** (getDefaultConfig, zgTestnet, darkTheme)
- [ ] **Step 7: Create business constants** (PLATFORM_FEE_RATE, MIN_DEPOSIT_USDT, MOCK_DELAY_MS, SUPPORTED_CURRENCIES)
- [ ] **Step 8: Restructure main.tsx Provider tree** (WagmiProvider > QueryClientProvider > RainbowKitProvider > BrowserRouter > App)
- [ ] **Step 9: Verify build**
- [ ] **Step 10: Commit**

---

### Task 1: Types, Mock Data & Services

**Files:**
- Create: `src/types/index.ts` — User, DepositInfo, DepositRecord, Region, Operator, VirtualNumber(with credentials), Bill(string amounts), MonthEstimate, Notification
- Create: `src/utils/format.ts` — shortenAddress, formatAmount(bigint), formatUSD(string), formatDate, timeAgo
- Create: `src/services/mock/delay.ts` — delay(ms) helper
- Create: `src/services/mock/data.ts` — all mock data + mutation helpers
- Create: `src/services/mock/userService.ts` — getUserProfile(→User|null), register, sendVerificationCode, verifyEmail
- Create: `src/services/mock/depositService.ts` — getDeposit, getDepositHistory, deposit, withdraw
- Create: `src/services/mock/billingService.ts` — getBills(filter), getBillDetail, payBill, getCurrentMonthEstimate
- Create: `src/services/mock/operatorService.ts` — getRegions, getOperatorsByRegion, getMyNumbers, applyNumber(with eSIM credentials)
- Create: `src/services/mock/notificationService.ts` — getNotifications, getUnreadCount, markAsRead, markAllAsRead
- Create: `src/services/index.ts` — unified re-export

- [ ] **Step 1-10: Create all files above**
- [ ] **Step 11: Verify tsc --noEmit**
- [ ] **Step 12: Commit**

---

### Task 2: React Query Hooks

**Files:**
- Create: `src/hooks/useUser.ts` — useUser(address), useRegister, useSendVerificationCode, useVerifyEmail
- Create: `src/hooks/useDeposit.ts` — useDeposit, useDepositHistory, useDepositMutation, useWithdrawMutation
- Create: `src/hooks/useBilling.ts` — useBills(address,filter), useBillDetail, usePayBill, useMonthEstimate
- Create: `src/hooks/useOperator.ts` — useRegions, useOperatorsByRegion, useMyNumbers, useApplyNumber
- Create: `src/hooks/useNotification.ts` — useNotifications, useUnreadCount(refetchInterval:30s), useMarkAsRead, useMarkAllAsRead

- [ ] **Steps 1-5: Create all hook files**
- [ ] **Step 6: Verify tsc --noEmit**
- [ ] **Step 7: Commit**

---

### Task 3: Shared Components + AppLayout

**Files:**
- Create: `src/components/shared/StatusBadge.tsx` — active(green)/inactive(yellow)/suspended(red)
- Create: `src/components/shared/AmountDisplay.tsx` — bigint or string, sm/md/lg sizes
- Create: `src/components/shared/GuardCard.tsx` — icon + title + message + action button
- Create: `src/components/shared/EmptyState.tsx` — icon + message centered
- Create: `src/components/layout/TabBar.tsx` — 5 tabs, active highlight, badges (unpaid/unread)
- Create: `src/components/layout/Header.tsx` — Dashboard variant + generic variant
- Create: `src/components/layout/AppLayout.tsx` — max-w-mobile mx-auto + Header + Outlet + TabBar + route guards

- [ ] **Steps 1-7: Create all component files**
- [ ] **Step 8: Verify tsc --noEmit**
- [ ] **Step 9: Commit**

---

### Task 4: Wallet Components + Landing Page + Router

**Files:**
- Create: `src/components/wallet/ConnectButton.tsx` — RainbowKit.Custom wrapper
- Create: `src/components/wallet/RegisterSheet.tsx` — vaul Drawer, email→verify two-step
- Create: `src/pages/Landing.tsx` — Logo + Hero + Stats + RegisterSheet auto-trigger
- Modify: `src/App.tsx` — lazy routes: / Landing, AppLayout wrapping all inner pages
- Delete: `src/pages/Home.tsx`

- [ ] **Steps 1-5: Create files, update router, delete Home.tsx**
- [ ] **Step 6: Create placeholder pages** (Dashboard, Deposit, Services, RegionDetail, Billing, BillDetail, Notifications)
- [ ] **Step 7: Verify pnpm build + dev**
- [ ] **Step 8: Commit**

---

### Task 5: Dashboard Page

**Files:** Modify `src/pages/Dashboard.tsx`

Full implementation: Status card (2x2 grid), usage card (3-col), quick actions 2x2, SUSPENDED overdue warning.
Data: useUser + useDeposit + useMonthEstimate + useMyNumbers + useUnreadCount

- [ ] **Step 1: Implement Dashboard**
- [ ] **Step 2: Verify build**
- [ ] **Step 3: Commit**

---

### Task 6: Deposit Page

**Files:** Modify `src/pages/Deposit.tsx`

Full implementation: Balance card (progress bar), Deposit(green)/Withdraw(outline) buttons, history list, vaul Drawer sheets.
Data: useDeposit + useDepositHistory + useDepositMutation + useWithdrawMutation

- [ ] **Step 1: Implement Deposit**
- [ ] **Step 2: Verify build**
- [ ] **Step 3: Commit**

---

### Task 7: Services + RegionDetail Pages

**Files:** Modify `src/pages/Services.tsx` + `src/pages/RegionDetail.tsx`

Services: Regions/My Numbers dual tabs, search, country list, number cards with eSIM credentials.
RegionDetail: Operator cards + Apply Number drawer.

- [ ] **Step 1: Implement Services**
- [ ] **Step 2: Implement RegionDetail**
- [ ] **Step 3: Verify build**
- [ ] **Step 4: Commit**

---

### Task 8: Billing + BillDetail Pages

**Files:** Modify `src/pages/Billing.tsx` + `src/pages/BillDetail.tsx`

Billing: Unpaid/History tabs, alert banner, bill cards with fee breakdown, pay drawer.
BillDetail: Full detail + usage stats + pay.

- [ ] **Step 1: Implement Billing**
- [ ] **Step 2: Implement BillDetail**
- [ ] **Step 3: Verify build**
- [ ] **Step 4: Commit**

---

### Task 9: Notifications Page

**Files:** Modify `src/pages/Notifications.tsx`

New/Earlier groups, colored left borders, unread blue bg + dot, mark read/all.

- [ ] **Step 1: Implement Notifications**
- [ ] **Step 2: Verify build**
- [ ] **Step 3: Commit**

---

### Task 10: Final Integration & Smoke Test

- [ ] **Step 1: tsc --noEmit** → 0 errors
- [ ] **Step 2: pnpm build** → success
- [ ] **Step 3: Visual smoke test** — all 6 pages, wallet, tabs, sheets, guards
- [ ] **Step 4: Final commit**
