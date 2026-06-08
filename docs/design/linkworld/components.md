# components.md — Link World Web 基线·组件清单

> 扫描对象：`packages/web/src/components` + `pages`
> 用途：重构时区分「可复用 UI 原子」与「需重做的业务组件」。
> 扫描时间：2026-06-06

## A. UI 基础组件（shadcn / base-ui 封装，`components/ui/`）

特点：基于 `@base-ui/react` 原语 + `class-variance-authority`(cva) 管理变体，样式全部走 **CSS 变量**（`bg-primary` / `text-foreground` / `border-border`），换主题改 `index.css` 变量即可全局生效。**重构友好，建议保留。**

| 文件路径 | 组件名 | 职责 | props / 变体概要 |
|----------|--------|------|------------------|
| `components/ui/button.tsx` | `Button` / `buttonVariants` | 基础按钮，包裹 base-ui Button | variant: default/outline/secondary/ghost/destructive/link；size: default/xs/sm/lg/icon/icon-xs/icon-sm/icon-lg；透传 `ButtonPrimitive.Props` |
| `components/ui/badge.tsx` | `Badge` / `badgeVariants` | 标签徽章，基于 `useRender` 可多态渲染 | variant: default/secondary/destructive/outline/ghost/link；`useRender.ComponentProps<"span">` + render |
| `components/ui/tabs.tsx` | `Tabs` / `TabsList` / `TabsTrigger` / `TabsContent` / `tabsListVariants` | 选项卡，包裹 base-ui Tabs | Tabs: orientation horizontal/vertical；TabsList variant: default/line；透传各 `TabsPrimitive.*.Props` |

> 注意：实际只有 3 个 UI 原子组件。常见的 Card / Input / Dialog / Sheet 等**尚未沉淀为 ui 组件** —— 页面里用裸 `<div>`/`<input>` + Tailwind 手写，或直接用 `vaul` 的 Drawer。重构时可考虑补齐。

## B. 业务组件

### B1. 布局组件（`components/layout/`）

| 文件路径 | 组件名 | 职责 | props 概要 |
|----------|--------|------|-----------|
| `components/layout/AppLayout.tsx` | `AppLayout` | 路由守卫 + 移动端外壳（max-w-mobile + Header + Outlet + TabBar）；按 user.status 拦截受限路由并渲染 GuardCard | 无 props，内部读 `useAccount`/`useUser`/`useLocation` |
| `components/layout/Header.tsx` | `Header` | 顶部栏；Dashboard 显示「欢迎+地址+通知铃铛+头像渐变球」，其他页显示标题+返回箭头 | 无 props，内部读路由+`useUnreadCount` |
| `components/layout/TabBar.tsx` | `TabBar` | 固定底部 5 项导航（Home/Services/Deposit/Bills/Cards），emoji 图标，Bills 带未付账单角标 | 无 props，内部读 `useBills(unpaid)` |

### B2. 通用展示组件（`components/shared/`）

| 文件路径 | 组件名 | 职责 | props 概要 |
|----------|--------|------|-----------|
| `components/shared/AmountDisplay.tsx` | `AmountDisplay` | 金额展示（bigint 自动 formatAmount，可带币种后缀） | `{ amount: bigint\|string; currency?: string; size?: "sm"\|"md"\|"lg"; colorClass?: string }`（colorClass 默认 `text-status-warning`） |
| `components/shared/BottomSheet.tsx` | `BottomSheet` | 通用底部抽屉（基于 vaul Drawer） | `{ open: boolean; onOpenChange: (open)=>void; children: ReactNode }` |
| `components/shared/EmptyState.tsx` | `EmptyState` | 空状态占位（emoji + 文案） | `{ icon: string; message: string }` |
| `components/shared/GuardCard.tsx` | `GuardCard` | 守卫拦截卡片（图标+标题+说明+跳转按钮，复用 Button） | `{ icon: string; title: string; message: string; actionLabel: string; actionPath: string }` |
| `components/shared/StatusBadge.tsx` | `StatusBadge` | 用户状态徽章（active/inactive/suspended，圆点+文案，颜色 = status token） | `{ status: UserStatus }` |

### B3. 钱包组件（`components/wallet/`）

| 文件路径 | 组件名 | 职责 | props 概要 |
|----------|--------|------|-----------|
| `components/wallet/ConnectButton.tsx` | `ConnectButton` | 包裹 RainbowKit `ConnectButton.Custom`，未连接时显示自定义按钮 | `{ label?: string }`（默认 "Connect Wallet"）；样式 `bg-brand-blue text-white`（硬编码，未复用 Button） |
| `components/wallet/RegisterSheet.tsx` | `RegisterSheet` | 注册流程抽屉（两步：邮箱→验证码→合约注册+后端同步），含邮箱校验 | `{ address: string; open: boolean; onClose: ()=>void; onSuccess: ()=>void }`；内部用 `useSendVerificationCode`/`useVerifyEmail`/`useRegister` |

### B4. 页面组件（`pages/`，均为 `export default function`）

| 文件路径 | 组件名 | 职责 | 主要依赖 hook |
|----------|--------|------|---------------|
| `pages/Landing.tsx` | `Landing` | 落地页/连接钱包入口，已注册自动跳 dashboard | `useUser`、`useNavigate` |
| `pages/Dashboard.tsx` | `Dashboard` | 首页总览（押金/账单预估/号码概览快捷入口） | `useUser`/`useDeposit`/`useMonthEstimate`/`useMyNumbers` |
| `pages/Deposit.tsx` | `Deposit` | 充值/提现 + 历史记录 | `useDeposit`/`useDepositHistory`/`useDepositMutation`/`useWithdrawMutation` |
| `pages/Services.tsx` | `Services` | 地区列表 + 我的号码 | `useRegions`/`useMyNumbers` |
| `pages/RegionDetail.tsx` | `RegionDetail` | 地区下运营商列表 + 申请号码 | `useOperatorsByRegion`/`useApplyNumber`/`useRegions`（path 参 regionCode） |
| `pages/Billing.tsx` | `Billing` | 账单列表（含未付筛选）+ 支付 | `useBills`/`usePayBill` |
| `pages/BillDetail.tsx` | `BillDetail` | 单账单详情 + 支付 | `useBillDetail`/`usePayBill`（path 参 billId） |
| `pages/Notifications.tsx` | `Notifications` | 通知中心 + 已读/全部已读 | `useNotifications`/`useMarkAsRead`/`useMarkAllAsRead` |
| `pages/Cards.tsx` | `Cards` | 流量卡 NFT 列表/管理 | `hooks/contracts/useTrafficCard`（链上） |

## C. 重构提示（组件层）

- **可直接复用**：`components/ui/*`（cva + CSS 变量，最规范）、`shared/EmptyState`、`shared/GuardCard`、`shared/BottomSheet`。
- **硬编码颜色需清洗**：`Header`（`text-white`、`bg-gradient-to-br from-brand-blue to-brand-purple`）、`TabBar`（`text-brand-blue`/`text-white` 角标）、`ConnectButton`（`bg-brand-blue text-white`，建议改用 `<Button>`）、`AmountDisplay`/`StatusBadge`（直接吃 `text-status-*` token）。
- **emoji 当图标**：AppLayout/Header/TabBar/Dashboard 大量用 emoji（🏠📱💰📄🎟️🔔），已装 lucide-react 却基本没用 —— 新主题（深蓝+金）下建议统一换 lucide 图标。
- **缺失原子组件**：无 Card/Input/Dialog 封装，页面手写较多，重构可借机沉淀。
