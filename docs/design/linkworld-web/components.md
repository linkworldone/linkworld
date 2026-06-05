# Link World Web — 组件清单 (components)

> 反映 **当前 HEAD** 代码，仅描述现状。`src/components/` 下共 **13 个 .tsx 组件**，分 4 组：ui(3) / shared(5) / layout(3) / wallet(2)。
> 「耦合度」一列描述各组件当前对路由/链上状态/硬编码样式的依赖程度（高=纯展示低耦合；中=有样式或少量逻辑耦合；高耦合=强依赖 location/account 或样式硬编码到具体 class），仅作现状标注，不含改动判断。

## ui — shadcn 原语（基于 @base-ui/react + cva）

| 文件 | 组件 | 用途 | 关键 props | 基座 | 耦合度 |
|------|------|------|-----------|------|--------|
| `components/ui/button.tsx` | `Button` / `buttonVariants` | 通用按钮 | `variant`(default/outline/secondary/ghost/destructive/link), `size`(default/xs/sm/lg/icon/icon-xs/icon-sm/icon-lg), `className` | @base-ui/react + cva | 高 |
| `components/ui/badge.tsx` | `Badge` / `badgeVariants` | 标签徽标 | `variant`(default/secondary/destructive/outline/ghost/link), `render`, `className` | @base-ui/react + cva | 高 |
| `components/ui/tabs.tsx` | `Tabs` / `TabsList` / `TabsTrigger` / `TabsContent` / `tabsListVariants` | 标签页切换 | `orientation`, `variant`(default/line), `className` | @base-ui/react + cva | 高 |

ui 组颜色现状：
- **全部走 shadcn 语义 token**（`bg-primary`/`text-primary-foreground`/`bg-muted`/`border-border`/`text-destructive` 等），**不含硬编码 hex**，色值来源是 index.css 的 oklch token。
- 含大量 `dark:` 变体类（依赖 tailwind `darkMode: "class"`，但当前 `:root` 已是暗色，未见显式给 root 加 `.dark`/`.theme` class）。
- 不耦合路由/location，纯展示组件。

## shared — 业务通用组件

| 文件 | 组件 | 用途 | 关键 props | 依赖基础组件 | 耦合度 |
|------|------|------|-----------|------------|--------|
| `components/shared/AmountDisplay.tsx` | `AmountDisplay` | 金额展示（bigint→格式化，调 `formatAmount`） | `amount`(bigint\|string), `currency?`, `size?`(sm/md/lg), `colorClass?` | 无 | 中 |
| `components/shared/BottomSheet.tsx` | `BottomSheet` | 底部抽屉容器（vaul Drawer） | `open`, `onOpenChange`, `children` | vaul | 中 |
| `components/shared/EmptyState.tsx` | `EmptyState` | 空状态占位 | `icon`(emoji string), `message` | 无 | 高 |
| `components/shared/GuardCard.tsx` | `GuardCard` | 守卫拦截卡片（充值/停用提示） | `icon`, `title`, `message`, `actionLabel`, `actionPath` | ui/Button + react-router `useNavigate` | 中 |
| `components/shared/StatusBadge.tsx` | `StatusBadge` | 用户状态徽标（active/inactive/suspended） | `status`(UserStatus) | 无 | 中 |

shared 组色值/耦合现状：
- **AmountDisplay**：`colorClass` 默认值 `"text-status-warning"`（tailwind.config 自定义色，非 shadcn token）；`size` 对应的文本字号类内联在 `sizeClasses` map。
- **BottomSheet**：用 tailwind.config 自定义色 `bg-surface-card` / `bg-surface-secondary`、`bg-black/60` 遮罩、`max-w-mobile`。无路由耦合。
- **EmptyState**：纯展示，仅用 `text-text-secondary`。
- **GuardCard**：耦合路由（内部 `useNavigate()`，靠 `actionPath` 跳转）；渲染 ui/Button；文字用 `text-text-primary` / `text-text-secondary`。
- **StatusBadge**：内置 `statusConfig` 把 UserStatus 映射到 `{label, color: text-status-*, dot: bg-status-*}`（tailwind.config status 色）。颜色语义集中在此 map。
- 共性：icon 均为 emoji（string），非图标组件。

## layout — 布局/外壳

| 文件 | 组件 | 用途 | 关键 props | 依赖基础组件 | 耦合度 |
|------|------|------|-----------|------------|--------|
| `components/layout/AppLayout.tsx` | `AppLayout` | 路由外壳 + 守卫 + Header/TabBar | 无（渲染 `<Outlet/>`） | Header / TabBar / GuardCard + react-router + wagmi | 高耦合 |
| `components/layout/Header.tsx` | `Header` | 顶栏（dashboard 欢迎语 / 子页标题+返回） | 无 | react-router + wagmi + useUnreadCount | 高耦合 |
| `components/layout/TabBar.tsx` | `TabBar` | 底部固定导航（5 tab + 未付账单 badge） | 无 | react-router + wagmi + useBills | 高耦合 |

layout 组色值/耦合现状：
- **AppLayout**：耦合 `useLocation` + `useAccount` + `useUser`，内置守卫分支（未连钱包/无 user→`Navigate("/")`；inactive 拦 `/services|/billing|/notifications`；suspended 拦 `/services|/deposit`，受限路径白/黑名单写死在局部数组里）。容器用 `bg-surface`、`max-w-mobile`、`text-text-secondary`。守卫逻辑与布局同文件。
- **Header**：耦合 `useLocation`（`isDashboard`/`isSubPage` 判定）+ `useNavigate`；`pageTitles` map 写死路径→标题；含 `bg-gradient-to-br from-brand-blue to-brand-purple`（头像圆点）+ `bg-status-danger`（未读红点）；emoji 铃铛 `\u{1F514}`。
- **TabBar**：`tabs` 常量数组写死（label/emoji icon/path/badgeKey）；`isActive` 靠 `pathname.startsWith`；active 色 `text-brand-blue`、非 active `text-text-muted`、badge `bg-status-danger`；`fixed bottom-0 max-w-mobile`，含 `pb-[env(safe-area-inset-bottom)]`。emoji 图标 + 品牌色直接写到 class。

## wallet — 钱包相关

| 文件 | 组件 | 用途 | 关键 props | 依赖基础组件 | 耦合度 |
|------|------|------|-----------|------------|--------|
| `components/wallet/ConnectButton.tsx` | `ConnectButton` | 连接钱包按钮（RainbowKit `ConnectButton.Custom` 包装） | `label?`（默认 "Connect Wallet"） | RainbowKit | 中 |
| `components/wallet/RegisterSheet.tsx` | `RegisterSheet` | 注册抽屉（邮箱→验证码→合约注册两步流） | `address`, `open`, `onClose`, `onSuccess` | vaul + useRegister/useSendVerificationCode/useVerifyEmail | 高耦合 |

wallet 组色值/耦合现状：
- **ConnectButton**：包 RainbowKit `ConnectButton.Custom`，按钮用原生 `<button>`，样式 `bg-brand-blue text-white`（未走 ui/Button 原语）。
- **RegisterSheet**：内置两步状态机（email/verify）+ 三个 mutation hook（useRegister/useSendVerificationCode/useVerifyEmail）+ 邮箱正则校验；UI 用 vaul Drawer，色值 `bg-surface-card`/`bg-surface-secondary`、`focus:border-brand-blue`、`bg-black/60`、`text-status-danger`；文案中英混排（如 "请输入有效的邮箱地址"）。RainbowKit `darkTheme({ accentColor })` 在 main.tsx 写死 `#3b82f6`。

## 业务页面（pages/）

9 个页面，全部 `export default`，由 App.tsx `lazy` 加载。各页消费的 hooks / 复用组件 / 工具如下（均现状）：

| 页面 | 路由 | 消费 hooks | 复用组件 | 工具 / 其他 |
|------|------|-----------|---------|------------|
| `Landing.tsx` | `/` | `useUser` + `useAccount` | ConnectButton / RegisterSheet | `useNavigate`；含 `font-orbitron` logo + 2 处 brand 渐变 |
| `Dashboard.tsx` | `/dashboard` | `useUser` / `useDeposit` / `useMonthEstimate` / `useMyNumbers` + `useAccount` | StatusBadge / AmountDisplay | `useNavigate`；余额卡 surface 渐变 |
| `Deposit.tsx` | `/deposit` | `useDeposit` / `useDepositHistory` / `useDepositMutation` / `useWithdrawMutation` + `useAccount` | ui/Button / BottomSheet / AmountDisplay | `formatAmount` / `formatDate` / `toast`；余额卡 surface 渐变 |
| `Services.tsx` | `/services` | `useRegions` / `useMyNumbers` + `useAccount` | EmptyState | `useNavigate` |
| `RegionDetail.tsx` | `/services/:regionCode` | `useOperatorsByRegion` / `useApplyNumber` / `useRegions` + `useAccount` | ui/Button / BottomSheet | `useParams` / `formatAmount`；直接调 `apiClient` |
| `Billing.tsx` | `/billing` | `useBills` / `usePayBill` + `useAccount` | ui/Button / ui/Badge / BottomSheet | `useNavigate` / `formatDate` / `formatUSD` / `parseUnits` |
| `BillDetail.tsx` | `/billing/:billId` | `useBillDetail` / `usePayBill` | ui/Button / ui/Badge / BottomSheet | `useParams` / `formatUSD` / `formatDate` / `parseUnits` |
| `Notifications.tsx` | `/notifications` | `useNotifications` / `useMarkAsRead` / `useMarkAllAsRead` + `useAccount` | EmptyState | `timeAgo`；用 `types.Notification` |
| `Cards.tsx` | `/cards` | TrafficCard 合约 hooks + `useAccount` | ui/Button | 卡面 surface 渐变 |

> 备注：`RegionDetail` 是唯一一个绕过 hooks 直接 `import { apiClient }` 调后端的页面（其余页面统一走 hooks）。

## 颜色硬编码总览（组件层，现状）

- **组件/页面源码内无任何裸 hex**——所有颜色都走 class 名（要么 shadcn 语义 token，要么 tailwind.config 的 `surface/brand/status/text` hex 别名）。
- gradient 工具类共 6 处（Header 1 + Landing 2 + Dashboard/Deposit/Cards 各 1）：brand-blue→brand-purple 3 处、surface-gradient-from→to 3 处。详见 color-mapping.md 第 3 节。
- 半透明遮罩 `bg-black/60` 2 处（BottomSheet、RegisterSheet）。
- 色值来源两套：① tailwind.config.ts 的 hex（surface/brand/status/text/border）；② index.css 的 oklch token（shadcn）；另有 main.tsx 的 RainbowKit accentColor 与 ConnectButton 的 `bg-brand-blue` 两处链上 UI 硬编码。
</content>
