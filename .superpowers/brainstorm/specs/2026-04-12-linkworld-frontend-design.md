# LinkWorld 前端设计规格

## 1. 项目概述

LinkWorld 是基于 0G Chain 的去中心化全球通信服务平台。用户通过钱包连接，存入押金后即可获取全球虚拟号码，享受跨境流量与通话服务。无需 KYC，先用后付。

本文档定义前端 Web 应用的完整设计规格。当前阶段合约和后端均使用 Mock 数据，前端按真实接口签名开发，后续替换 Mock 实现即可接入真实服务。

## 2. 技术栈

| 层级 | 选型 | 说明 |
|------|------|------|
| 框架 | React 19 + Vite 6 | 已配置 |
| 路由 | React Router 7 | 已配置 |
| 样式 | TailwindCSS 3.4 + shadcn/ui | shadcn/ui 待安装 |
| Web3 | wagmi 2.14 + viem 2.21 + RainbowKit 2.2 | 已安装未集成 |
| 数据层 | TanStack React Query 5.62 | 已安装 |
| 状态管理 | React Context（用户信息、钱包状态） | 轻量全局状态 |
| Mock 策略 | Service 层 Mock | 接口签名与真实 API/合约一致，切换时只改 service 内部实现 |
| 链 | 0G Chain Mainnet (16600) / Testnet (16601) | 已配置 |

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

### 3.2 目录结构

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
│   │   ├── AppLayout.tsx       # 登录后布局（Header + TabBar + Content）
│   │   ├── Header.tsx          # 顶部栏
│   │   └── TabBar.tsx          # 底部 Tab 导航
│   ├── wallet/
│   │   ├── ConnectButton.tsx   # RainbowKit 封装
│   │   └── RegisterModal.tsx   # 注册弹窗（邮箱验证）
│   └── shared/
│       ├── StatusBadge.tsx     # 账户状态标签
│       ├── AmountDisplay.tsx   # 金额显示（支持多币种）
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
├── contexts/
│   └── UserContext.tsx          # 用户全局状态 (profile, status)
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

## 4. 用户状态机

### 4.1 四种账户状态

| 状态 | 说明 | 可访问页面 |
|------|------|-----------|
| UNREGISTERED | 未连接钱包或未注册 | Landing Page only |
| INACTIVE | 已注册，无押金或押金不足 | Dashboard, Deposit |
| ACTIVE | 押金充足，正常使用中 | 全部页面 |
| SUSPENDED | 存在未结账单，服务暂停 | Dashboard, Billing |

### 4.2 状态流转

```
UNREGISTERED --[连接钱包+注册]--> INACTIVE
INACTIVE     --[存入押金≥最低要求]--> ACTIVE
ACTIVE       --[月初有未结账单]--> SUSPENDED
SUSPENDED    --[结清账单]--> ACTIVE
ACTIVE       --[提取押金(需结清账单)]--> INACTIVE
SUSPENDED    --[逾期2周]--> 扣除押金 --> INACTIVE
```

### 4.3 路由守卫

- 未连接钱包：所有路由重定向到 `/`（Landing）
- INACTIVE：访问 Services/Billing/Notifications 时，显示提示卡片引导去 Deposit 页面
- SUSPENDED：访问 Services/Deposit 时，显示提示卡片引导去 Billing 页面结清账单

## 5. 页面设计

### 5.1 Landing Page (`/`)

未连接钱包时展示的营销首页。

**布局**：全屏竖向单页
- 顶部导航：Logo + Connect Wallet 按钮
- Hero 区域：标语 + 副标题 + Get Started / Learn More 双 CTA
- 底部统计：50+ Countries / 2.5% Fee / 0 KYC
- 品牌标识：Powered by 0G Chain & Chainlink

**交互**：
- 点击 Get Started / Connect Wallet → 触发 RainbowKit 连接弹窗
- 首次连接的钱包 → 弹出注册 Sheet（输入邮箱 + 验证码确认）
- 已注册钱包 → 直接跳转 Dashboard

### 5.2 Dashboard (`/dashboard`)

登录后主页，一屏展示所有关键信息。

**布局**：
- 顶部 Header：Welcome back + 钱包地址 + 通知铃铛（带角标）+ 头像
- 状态大卡片（渐变背景）：
  - 左上：Account Status（Active/Inactive/Suspended + 指示灯）
  - 右上：Deposit 余额 + 币种
  - 左下：当前虚拟号码
  - 右下：所在地区 + 国旗
- 本月用量卡片：Data Used / Calls / Est. Bill 三栏
- Quick Actions 2x2 网格：Top Up / Switch Region / Bills / Pay Now
- 底部 Tab Bar

**数据依赖**：useUser (profile + status), useDeposit (balance), useBilling (current month estimate)

### 5.3 Deposit (`/deposit`)

押金管理，存入和提取操作。

**布局**：
- 余额大卡片（渐变背景）：
  - 余额数字（大字号居中）
  - 币种标识
  - 最低要求进度条 + 百分比
- 双按钮：Deposit（绿色实心）/ Withdraw（描边）
- 历史记录列表：时间 + 类型 + 金额（绿色+/红色-）

**交互**：
- Deposit 按钮 → 底部弹出 Sheet：
  - 选择币种（USDT / ETH）
  - 输入金额
  - 确认按钮 → Mock 模拟合约交互（loading → success toast）
- Withdraw 按钮 → 底部弹出 Sheet：
  - 显示可提取金额
  - 条件校验提示（需结清账单 + 无新费用）
  - 确认提取

**数据依赖**：useDeposit (balance, history, deposit mutation, withdraw mutation)

### 5.4 Services (`/services`)

浏览地区和运营商，申请虚拟号码。

**布局**：
- 搜索栏：Search country or region...
- 当前号码卡片（如已有）：号码 + 地区 + 运营商 + Manage 按钮
- 国家列表：国旗 + 国名 + 运营商数量 + 起步价 + 箭头

**子页 — ServiceDetail (`/services/:regionCode`)**：
- 返回箭头 + 国家名
- 运营商列表卡片：名称 + 资费详情 + 所需最低押金
- 选择运营商 → 号码申请 Sheet：
  - 确认运营商和地区
  - 查看所需押金
  - 确认申请 → Mock 分配虚拟号码

**数据依赖**：useOperator (regions, operators, applyNumber mutation)

### 5.5 Billing (`/billing`)

账单列表和支付。

**布局**：
- 顶部分段切换：Unpaid / History
- 未付账单告警条（如有）：红色背景，显示数量和截止日期
- 账单卡片列表：
  - 月份标题 + 状态标签（Unpaid 红色 / Paid 绿色）
  - 费用明细：Operator Fee + Platform Fee (2.5%)
  - 总额 + Pay Now 按钮（未付时显示）

**子页 — BillDetail (`/billing/:billId`)**：
- 完整账单详情
- 费用分项（运营商费 + 平台手续费）
- 用量明细（Data / Calls / Duration）
- 支付按钮

**交互**：
- Pay Now → 底部弹出支付确认 Sheet：
  - 显示金额 + 支付来源（押金扣除或钱包直付）
  - 确认支付 → Mock 合约交互

**数据依赖**：useBilling (unpaid bills, history, pay mutation)

### 5.6 Notifications (`/notifications`)

通知中心。

**布局**：
- 顶部：标题 + Mark all read 链接
- 分组：New / Earlier
- 通知卡片：
  - 未读：蓝色背景 + 左侧彩色边条 + 蓝色未读圆点
  - 已读：暗色背景，文字降低对比度
  - 内容：标题 + 描述 + 时间

**通知类型及边条颜色**：
- 账单到期提醒：蓝色
- 存款确认：绿色
- 支付确认：绿色
- 服务暂停警告：红色
- 系统通知：灰色

**数据依赖**：useNotification (list, unread count, markRead mutation, markAllRead mutation)

## 6. 通用组件设计

### 6.1 AppLayout

登录后所有页面的外壳，包含 Header + 内容区 + TabBar。

```
┌──────────────────────┐
│       Header         │  固定顶部
├──────────────────────┤
│                      │
│    Page Content      │  可滚动区域
│    (slot)            │
│                      │
├──────────────────────┤
│       TabBar         │  固定底部 (safe area)
└──────────────────────┘
```

### 6.2 TabBar

5 个 Tab，固定底部：

| Tab | Icon | 路由 | 角标 |
|-----|------|------|------|
| Home | 🏠 | /dashboard | - |
| Services | 📱 | /services | - |
| Deposit | 💰 | /deposit | - |
| Bills | 📄 | /billing | 未付数量 |
| Alerts | 🔔 | /notifications | 未读数量 |

当前 Tab 图标和文字高亮为主题蓝色。

### 6.3 Header

- Dashboard 页面：Welcome back + 钱包地址缩写 + 通知铃铛 + 头像
- 其他页面：页面标题 + 返回箭头（二级页面时）

### 6.4 底部 Sheet

用于所有操作确认：存入押金、提取押金、支付账单、申请号码。

- 底部弹出，带遮罩
- 顶部拖拽条
- 内容区 + 底部确认按钮（宽度 100%）
- 操作中显示 loading spinner

## 7. Mock 数据设计

### 7.1 Mock 用户

```typescript
interface User {
  address: string           // 钱包地址
  email: string             // 注册邮箱
  status: 'inactive' | 'active' | 'suspended'
  registeredAt: string      // ISO date
  nftTokenId: number        // 身份 NFT ID
}
```

### 7.2 Mock 押金

```typescript
interface DepositInfo {
  balance: bigint           // 当前余额 (wei)
  minimumRequired: bigint   // 最低要求 (wei)
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

### 7.3 Mock 运营商 & 号码

```typescript
interface Region {
  code: string              // 'JP', 'US', 'GB'...
  name: string
  flag: string              // emoji
  operatorCount: number
  startingPrice: number     // USD/month
}

interface Operator {
  id: string
  name: string
  region: string
  requiredDeposit: bigint
  dataRate: number          // USD per GB
  callRate: number          // USD per minute
  isActive: boolean
}

interface VirtualNumber {
  id: string
  number: string            // e.g. '+81 90-1234-5678'
  region: string
  operator: string
  status: 'active' | 'inactive'
  activatedAt: string
}
```

### 7.4 Mock 账单

```typescript
interface Bill {
  id: string
  month: string             // '2026-04'
  status: 'unpaid' | 'paid' | 'overdue'
  operatorFee: number       // USD
  platformFee: number       // USD (2.5% of operatorFee)
  totalAmount: number       // USD
  dueDate: string
  paidAt?: string
  usage: {
    dataGB: number
    callMinutes: number
  }
}
```

### 7.5 Mock 通知

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
getUserProfile(address: string): Promise<User>
register(address: string, email: string): Promise<User>
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

- 背景主色：`#0a0a14`（最深）、`#0f0f1a`（卡片）、`#1a1a2e`（次级元素）
- 主题蓝：`#3b82f6`（操作/高亮）
- 成功绿：`#22c55e`（Active、存入、已付）
- 警告黄：`#f59e0b`（Deposit 余额、Inactive）
- 危险红：`#ef4444`（Unpaid、Suspended、扣款）
- 信息青：`#06b6d4`（通知）
- 紫色：`#8b5cf6`（通话用量）
- 文字白：`#e0e0e0`（主文字）、`#888`（次级）、`#555`（第三级）

### 9.2 间距与圆角

- 页面内边距：16-20px
- 卡片间距：12px
- 卡片圆角：12-16px
- 按钮圆角：8-12px
- 最小点击区域：44x44px

### 9.3 字体大小

- 页面标题：17px bold
- 卡片标题：13px semibold
- 正文：12-13px
- 辅助文字：10-11px
- 大数字展示：28-32px bold

### 9.4 移动端适配

- 设计基准：375px 宽度（iPhone SE/8）
- 底部 TabBar 高度：60px + safe area padding
- Header 高度：44px + status bar
- 所有操作使用底部弹出 Sheet，不跳转新页面
- 列表项高度至少 48px（触控友好）
