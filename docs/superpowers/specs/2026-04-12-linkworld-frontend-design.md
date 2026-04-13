# LinkWorld 前端实现设计规格

> **日期**: 2026-04-12 | **状态**: Approved | **Round**: 1

## 1. 项目概述

LinkWorld 是基于 0G Chain 的去中心化全球通信服务平台。用户通过钱包连接，存入押金后即可获取全球虚拟号码，享受跨境流量与通话服务。无需 KYC，先用后付。

本文档定义前端 Web 应用的完整实现设计规格。当前阶段合约和后端均使用 Mock 数据，前端按真实接口签名开发，后续替换 Mock 实现即可接入真实服务。

## 2. 技术栈

| 层级 | 选型 | 状态 |
|------|------|------|
| 框架 | React 19 + Vite 6 | 已配置 |
| 路由 | React Router 7 | 已配置 |
| 样式 | TailwindCSS 3.4 + shadcn/ui（按需引入） | shadcn/ui 待安装 |
| Web3 | wagmi 2.14 + viem 2.21 + RainbowKit 2.2 | 已安装未集成 |
| 数据层 | TanStack React Query 5.62 | 已安装 |
| 状态管理 | React Query + wagmi hooks（无额外 Context） | 决策确认 |
| Mock 策略 | Service 层 Mock + 真实 RainbowKit 钱包连接 | 决策确认 |
| 链 | 0G Chain Testnet (16601) | 已配置 |

### 2.1 shadcn/ui 集成策略

**按需引入**：只使用 Sheet、Dialog、Button、Badge、Tabs 等基础交互组件。卡片和布局手写 Tailwind，保持品牌风格一致性。

### 2.2 Mock 策略

- **钱包连接**：真实 RainbowKit，连接真实测试网钱包（MetaMask 等）
- **业务逻辑**：注册、押金、账单、运营商、通知全部走 Mock Service
- **接口签名**：Mock Service 的方法签名与未来真实 API/合约调用完全一致
- **切换方式**：`services/index.ts` 统一导出，替换 mock 实现即可接入真实服务

## 3. 架构设计

### 3.1 三层数据流

```
Pages / Components
       ↓
React Query Hooks (useUser, useDeposit, useBilling, useOperator, useNotification)
       ↓
Mock Service 层 (当前) → 未来替换为真实 API / 合约调用
       ↓
Mock Data  |  [未来] 0G Chain  |  [未来] Backend API
```

- **Pages/Components**：纯 UI 层，只消费 hooks 提供的数据和 mutation
- **Hooks**：React Query 封装，负责缓存、loading/error 状态、乐观更新
- **Services**：数据接口层，当前返回 mock 数据，接口签名与合约/API 一致

### 3.2 状态管理方案

**React Query + wagmi hooks，不使用额外 Context**：

- 服务端状态（用户 profile、押金、账单等）→ React Query 自定义 hooks
- 钱包状态（连接状态、地址、余额）→ wagmi hooks（useAccount, useBalance 等）
- 少量 UI 状态（注册弹窗开关等）→ 组件 local state 或 URL params

### 3.3 Provider 嵌套顺序

```tsx
<WagmiProvider config={wagmiConfig}>
  <QueryClientProvider client={queryClient}>
    <RainbowKitProvider>
      <BrowserRouter>
        <App />
      </BrowserRouter>
    </RainbowKitProvider>
  </QueryClientProvider>
</WagmiProvider>
```

### 3.4 目录结构

```
packages/web/src/
├── pages/              # 页面组件
│   ├── Landing.tsx
│   ├── Dashboard.tsx
│   ├── Deposit.tsx
│   ├── Services.tsx
│   ├── RegionDetail.tsx
│   ├── Billing.tsx
│   ├── BillDetail.tsx
│   └── Notifications.tsx
├── components/         # 通用组件
│   ├── ui/             # shadcn/ui 组件
│   ├── layout/
│   │   ├── AppLayout.tsx       # 登录后布局（Header + TabBar + Content + 守卫）
│   │   ├── Header.tsx          # 顶部栏
│   │   └── TabBar.tsx          # 底部 Tab 导航
│   ├── wallet/
│   │   ├── ConnectButton.tsx   # RainbowKit 封装
│   │   └── RegisterSheet.tsx   # 注册底部弹窗（邮箱验证）
│   └── shared/
│       ├── StatusBadge.tsx     # 账户状态标签
│       ├── AmountDisplay.tsx   # 金额显示（支持多币种）
│       ├── GuardCard.tsx       # 状态引导卡片
│       └── EmptyState.tsx      # 空状态占位
├── hooks/              # React Query hooks
│   ├── useUser.ts
│   ├── useDeposit.ts
│   ├── useBilling.ts
│   ├── useOperator.ts
│   └── useNotification.ts
├── services/           # 数据接口层 (mock)
│   ├── mock/
│   │   ├── data.ts             # Mock 数据定义
│   │   ├── userService.ts
│   │   ├── depositService.ts
│   │   ├── billingService.ts
│   │   ├── operatorService.ts
│   │   └── notificationService.ts
│   └── index.ts               # 统一导出，未来切换入口
├── config/
│   ├── chains.ts               # 0G Chain 配置 (已有)
│   ├── wagmi.ts                # wagmi + RainbowKit 配置
│   └── constants.ts            # 业务常量
├── types/
│   └── index.ts                # 类型定义
├── utils/
│   └── format.ts               # 格式化工具
├── App.tsx
└── main.tsx
```

## 4. 路由与状态守卫

### 4.1 路由表

| 路径 | 页面 | Layout | 守卫 |
|------|------|--------|------|
| `/` | Landing | 无（全屏） | 已登录 → 重定向 `/dashboard` |
| `/dashboard` | Dashboard | AppLayout | 需登录 |
| `/deposit` | Deposit | AppLayout | 需登录；SUSPENDED → 引导去 Billing |
| `/services` | Services | AppLayout | 需登录；INACTIVE → 引导去 Deposit；SUSPENDED → 引导去 Billing |
| `/services/:regionCode` | RegionDetail | AppLayout | 同 Services |
| `/billing` | Billing | AppLayout | 需登录 |
| `/billing/:billId` | BillDetail | AppLayout | 需登录 |
| `/notifications` | Notifications | AppLayout | 需登录 |

### 4.2 用户状态机

四种账户状态：

| 状态 | 说明 | 可访问页面 |
|------|------|-----------|
| UNREGISTERED | 未连接钱包或未注册 | Landing Page only |
| INACTIVE | 已注册，无押金或押金不足 | Dashboard, Deposit |
| ACTIVE | 押金充足，正常使用中 | 全部页面 |
| SUSPENDED | 存在未结账单，服务暂停 | Dashboard, Billing, Notifications |

状态流转：

```
UNREGISTERED --[连接钱包+注册]--> INACTIVE
INACTIVE     --[存入押金≥最低要求]--> ACTIVE
ACTIVE       --[月初有未结账单]--> SUSPENDED
SUSPENDED    --[结清账单]--> ACTIVE
ACTIVE       --[提取押金(需结清账单)]--> INACTIVE
SUSPENDED    --[逾期2周]--> 扣除押金 --> INACTIVE
```

### 4.3 守卫实现（AppLayout 层）

AppLayout 内读取用户状态，根据 status + 当前路由决定渲染页面内容还是引导卡片：

```
用户访问任意 AppLayout 路由
  ↓
useAccount() → 未连接钱包？→ redirect /
  ↓
useUser() → 获取 profile
  ↓
status === INACTIVE?
  → 访问 Services/Billing/Notifications → GuardCard："请先存入押金"+ 跳转 Deposit
  → 访问 Dashboard/Deposit → 正常渲染
  ↓
status === SUSPENDED?
  → 访问 Services/Deposit → GuardCard："请先结清账单"+ 跳转 Billing
  → 访问 Dashboard/Billing/Notifications → 正常渲染
  ↓
status === ACTIVE → 全部正常渲染
```

### 4.4 注册流程

```
点击 Connect Wallet → RainbowKit 弹窗（真实连接）
  ↓ 连接成功
useUser(address) 查询 → 已注册？→ 跳转 /dashboard
  ↓ 未注册
弹出 RegisterSheet（底部弹出）
  → 输入邮箱 → sendVerificationCode() → 输入验证码 → verifyEmail()
  → 注册成功 → 跳转 /dashboard（status: INACTIVE）
```

## 5. 页面设计

### 5.1 Landing Page (`/`)

未连接钱包时展示的营销首页，全屏竖向单页。

- 顶部导航：Logo（渐变蓝紫） + Connect Wallet 按钮
- Hero 区域：地球 icon + 主标语 "Global Seamless Communication" + 副标题 + Get Started（渐变按钮）/ Learn More（描边按钮）双 CTA
- 底部统计：50+ Countries / 2.5% Fee / 0 KYC
- 品牌标识：Powered by 0G Chain & Chainlink

交互：
- 点击 Get Started / Connect Wallet → 触发 RainbowKit 连接弹窗
- 首次连接的钱包 → 弹出 RegisterSheet
- 已注册钱包 → 直接跳转 Dashboard

### 5.2 Dashboard (`/dashboard`)

登录后主页，一屏展示所有关键信息。

- Header：Welcome back + 钱包地址缩写 + 通知铃铛（带未读角标）+ 头像（渐变圆）
- 状态大卡片（渐变深色背景，2x2 网格）：Account Status（指示灯）/ Deposit Balance / Virtual Number / Region（国旗）
- 本月用量卡片：Data Used / Calls / Est. Bill 三栏分隔
- Quick Actions 2x2 网格：Top Up / Switch Region / Bills / Pay Now
- 用量详情区域：当前计费周期的实时流量/通话使用记录（Data Usage / Call Log），Mock 用定时递增数据模拟实时感
- 底部 TabBar

数据依赖：useUser, useDeposit, useBilling (current month estimate), useOperator (getMyNumbers), useNotification (unread count)

### 5.3 Deposit (`/deposit`)

押金管理页面。

- 余额大卡片（渐变背景，居中）：余额数字（36px bold）+ 币种 + 最低要求进度条 + 百分比
- 双按钮：Deposit（绿色实心）/ Withdraw（描边）
- 历史记录列表：时间 + 类型 + 金额（绿色+/红色-）

交互（底部 Sheet）：
- Deposit → 选择币种（USDT/ETH）→ 输入金额 → 确认 → Mock 合约交互（loading → success toast）
- Withdraw → 显示可提取金额 → 条件校验（需结清账单 + 无新费用）→ 确认提取

数据依赖：useDeposit

### 5.4 Services (`/services`)

浏览地区和运营商，申请虚拟号码。

- 搜索栏：Search country or region...
- My Numbers Tab：展示所有已申请号码列表（active/inactive），支持查看号码详情和接入凭证（eSIM 配置/VoIP 账号 mock）
- 当前号码卡片（如已有）：国旗 + 号码 + 运营商 + Manage 按钮
- 国家列表：国旗 + 国名 + 运营商数量 + 起步价 + 箭头

子页 RegionDetail (`/services/:regionCode`)：
- 返回箭头 + 国家名
- 运营商列表卡片：名称 + 资费详情 + 所需最低押金
- 选择运营商 → 号码申请 Sheet → 确认运营商和地区 → 查看所需押金 → 确认申请 → Mock 分配虚拟号码

数据依赖：useOperator

### 5.5 Billing (`/billing`)

账单列表和支付。

- 顶部分段切换：Unpaid / History
- 未付账单告警条（红色背景，显示数量和截止日期）
- 账单卡片列表：月份 + 状态标签 + 费用明细（Operator Fee + Platform Fee 2.5%）+ 总额 + Pay Now

子页 BillDetail (`/billing/:billId`)：
- 完整账单详情 + 费用分项 + 用量明细（Data / Calls）+ 支付按钮

交互：Pay Now → 支付确认 Sheet → 金额 + 支付来源 → 确认 → Mock 合约交互

数据依赖：useBilling

### 5.6 Notifications (`/notifications`)

通知中心。

- 顶部：标题 + Mark all read 链接
- 分组：New / Earlier
- 未读通知：蓝色背景 + 左侧彩色边条 + 蓝色未读圆点
- 已读通知：暗色背景，文字降低对比度

通知类型及边条颜色：
- bill_due：蓝色 #3b82f6
- deposit_confirmed / payment_confirmed：绿色 #22c55e
- service_suspended：红色 #ef4444
- system：灰色 #555

数据依赖：useNotification

## 6. 通用组件

### 6.1 AppLayout

登录后所有页面的外壳：固定 Header + 可滚动内容区 + 固定 TabBar + 状态守卫逻辑。

### 6.2 TabBar

5 个 Tab，固定底部：

| Tab | Icon | 路由 | 角标 |
|-----|------|------|------|
| Home | 🏠 | /dashboard | - |
| Services | 📱 | /services | - |
| Deposit | 💰 | /deposit | - |
| Bills | 📄 | /billing | 未付数量 |
| Alerts | 🔔 | /notifications | 未读数量 |

当前 Tab 高亮为主题蓝色 #3b82f6。

### 6.3 底部 Sheet

所有操作确认统一使用底部弹出 Sheet（shadcn/ui Sheet 组件）：
- 底部弹出，带遮罩 + 顶部拖拽条
- 内容区 + 底部确认按钮（宽度 100%）
- 操作中显示 loading spinner

### 6.4 GuardCard

状态引导卡片，用于 INACTIVE/SUSPENDED 用户访问受限页面时展示：
- 图标 + 提示文案 + 跳转按钮

## 7. Mock 数据类型

### 7.1 User
```typescript
interface User {
  address: string
  email: string
  status: 'inactive' | 'active' | 'suspended'
  registeredAt: string
  nftTokenId: number
}
```

### 7.2 DepositInfo & DepositRecord
```typescript
interface DepositInfo {
  balance: bigint
  minimumRequired: bigint
  currency: 'USDT' | 'ETH'
}

interface DepositRecord {
  id: string
  type: 'deposit' | 'withdraw' | 'deduction'
  amount: bigint
  currency: string
  timestamp: string
  txHash: string
}
```

### 7.3 Region, Operator & VirtualNumber
```typescript
interface Region {
  code: string
  name: string
  flag: string
  operatorCount: number
  startingPrice: number
}

interface Operator {
  id: string
  name: string
  region: string
  requiredDeposit: bigint
  dataRate: number
  callRate: number
  isActive: boolean
}

interface VirtualNumber {
  id: string
  number: string
  region: string
  operator: string
  status: 'active' | 'inactive'
  activatedAt: string
}
```

### 7.4 Bill
```typescript
interface Bill {
  id: string
  month: string
  status: 'unpaid' | 'paid' | 'overdue'
  operatorFee: number
  platformFee: number
  totalAmount: number
  dueDate: string
  paidAt?: string
  usage: { dataGB: number; callMinutes: number }
}
```

### 7.5 Notification
```typescript
interface Notification {
  id: string
  type: 'bill_due' | 'payment_confirmed' | 'deposit_confirmed' | 'service_suspended' | 'system'
  title: string
  message: string
  read: boolean
  createdAt: string
}
```

## 8. Service 接口定义

所有 service 方法签名与未来真实 API 一致，当前实现返回 mock 数据。

### 8.1 userService
```typescript
getUserProfile(address: string): Promise<User | null>
register(address: string, email: string): Promise<User>
sendVerificationCode(address: string, email: string): Promise<{ success: boolean }>
verifyEmail(address: string, code: string): Promise<boolean>
```

### 8.2 depositService
```typescript
getDeposit(address: string): Promise<DepositInfo>
getDepositHistory(address: string): Promise<DepositRecord[]>
deposit(address: string, amount: bigint, currency: string): Promise<DepositRecord>
withdraw(address: string, amount: bigint): Promise<DepositRecord>
```

### 8.3 billingService
```typescript
getBills(address: string, filter?: 'unpaid' | 'paid'): Promise<Bill[]>
getBillDetail(billId: string): Promise<Bill>
payBill(billId: string): Promise<Bill>
getCurrentMonthEstimate(address: string): Promise<{ dataGB: number; callMinutes: number; estimatedCost: number }>
```

### 8.4 operatorService
```typescript
getRegions(): Promise<Region[]>
getOperatorsByRegion(regionCode: string): Promise<Operator[]>
getMyNumbers(address: string): Promise<VirtualNumber[]>
applyNumber(address: string, operatorId: string): Promise<VirtualNumber>
```

### 8.5 notificationService
```typescript
getNotifications(address: string): Promise<Notification[]>
getUnreadCount(address: string): Promise<number>
markAsRead(notificationId: string): Promise<void>
markAllAsRead(address: string): Promise<void>
```

## 9. 设计规范

### 9.1 配色

| 用途 | 色值 |
|------|------|
| 背景主色 | `#0a0a14`（最深）|
| 卡片背景 | `#0f0f1a` |
| 次级元素 | `#1a1a2e` |
| 渐变卡片 | `linear-gradient(135deg, #1a1a3e, #0f1a2e)` |
| 主题蓝 | `#3b82f6`（操作/高亮/Tab 选中）|
| 品牌渐变 | `linear-gradient(135deg, #3b82f6, #8b5cf6)` |
| 成功绿 | `#22c55e`（Active、存入、已付）|
| 警告黄 | `#f59e0b`（Deposit 余额、Inactive）|
| 危险红 | `#ef4444`（Unpaid、Suspended、扣款）|
| 信息青 | `#06b6d4`（通知）|
| 紫色 | `#8b5cf6`（通话用量）|
| 主文字 | `#e0e0e0` |
| 次级文字 | `#888888` |
| 第三级文字 | `#555555` |

### 9.2 间距与圆角

- 页面内边距：16px
- 卡片间距：8-12px
- 卡片圆角：12-16px
- 按钮圆角：10-12px
- 最小点击区域：44x44px

### 9.3 字体大小

- 页面标题：17px bold
- 大数字展示：36px bold（余额）/ 22px bold（用量）
- 卡片标题：13-14px semibold
- 正文：12-13px
- 辅助文字：10-11px

### 9.4 移动端适配

- 设计基准：375px 宽度（iPhone SE/8）
- 底部 TabBar 高度：60px + safe area padding
- Header 高度：44px + status bar
- 所有操作使用底部弹出 Sheet，不跳转新页面
- 列表项高度至少 48px（触控友好）

## 10. 交付范围

Round 1 完整交付 6 个页面：
1. Landing Page — 营销首页 + 钱包连接 + 注册流程
2. Dashboard — 核心仪表盘
3. Deposit — 押金管理
4. Services — 地区浏览 + 虚拟号码申请
5. Billing — 账单管理 + 支付
6. Notifications — 通知中心

所有业务逻辑走 Mock Service，钱包连接使用真实 RainbowKit。
