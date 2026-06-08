# project-scan.md — Link World Web 基线·项目结构

> 扫描对象：`packages/web`（React 19 + Vite 6 + TS + Tailwind 3 + base-ui/shadcn + wagmi/viem/RainbowKit + react-router 7 + TanStack Query）
> 用途：重构前摸底「已有什么」。server/backend/contracts 不在范围内。
> 扫描时间：2026-06-06

## 1. 目录树（packages/web/src）

```
src/
├── App.tsx                 # 路由表（react-router Routes/Route，lazy 加载页面）
├── main.tsx                # 应用入口，Provider 树（Wagmi/Query/RainbowKit/Router/Toaster）
├── index.css               # Tailwind 入口 + shadcn CSS 变量（oklch 深色主题）
├── vite-env.d.ts
├── components/
│   ├── layout/             # 布局骨架：AppLayout(路由守卫) / Header / TabBar
│   ├── shared/             # 通用展示组件：AmountDisplay / BottomSheet / EmptyState / GuardCard / StatusBadge
│   ├── ui/                 # shadcn/base-ui 基础组件：button / badge / tabs
│   └── wallet/             # 钱包相关：ConnectButton(RainbowKit 封装) / RegisterSheet(注册流程)
├── config/
│   ├── abis/               # 7 个合约 ABI（UserRegistry/Deposit/ServiceManager/Payment/FeeManager/Oracle/TrafficCardNFT）+ index 桶导出
│   ├── api.ts              # API_BASE_URL（env VITE_API_BASE_URL）
│   ├── chains.ts           # viem defineChain：0G 主网/测试网 + Hardhat 本地
│   ├── constants.ts        # 业务常量（费率/最低押金/逾期天数/mock 延迟）
│   ├── contracts.ts        # 各链合约地址表 + getContractAddress()
│   └── wagmi.ts            # RainbowKit getDefaultConfig（按 env 选链）
├── hooks/
│   ├── contracts/          # 链上读写 hook（按合约分文件，基于 wagmi useReadContract/useWriteContract）
│   ├── useUser / useDeposit / useBilling / useOperator / useNotification  # 业务数据 hook（TanStack Query）
│   └── useTransactionFlow  # 交易状态/错误解析工具 hook
├── lib/
│   └── utils.ts            # cn() — clsx + tailwind-merge
├── pages/                  # 9 个页面（Landing/Dashboard/Deposit/Services/RegionDetail/Billing/BillDetail/Notifications/Cards）
├── services/
│   ├── api/                # 真实后端 API 封装（axios，client + 5 个 *Api 对象）
│   └── mock/               # mock 数据与服务（notification 仍走 mock）
├── types/
│   └── index.ts            # 领域类型 + 后端 snake_case API 响应类型
└── utils/
    ├── format.ts           # 格式化函数（地址/金额/USD/日期/相对时间）
    └── pendingSync.ts      # localStorage 待同步队列 + 指数退避重试
```

## 2. 各目录职责

| 目录 | 职责 | 重构关注点 |
|------|------|-----------|
| `components/ui` | shadcn/base-ui 原子组件，用 cva 管理变体，全部走 CSS 变量（`bg-primary`/`text-foreground`） | 主题色改 CSS 变量即可全局生效，重构最友好 |
| `components/shared` | 业务无关展示组件 | 大量硬编码 `text-status-*`/`text-text-*` Tailwind token，未走 CSS 变量 |
| `components/layout` | 路由守卫 + 顶/底导航 | Header/TabBar 用 emoji 当图标、含 `text-white`/`bg-gradient` 硬编码 |
| `components/wallet` | 钱包连接 + 注册流程 | ConnectButton 直接 `bg-brand-blue`，未复用 Button 组件 |
| `config` | 链/合约/wagmi/常量配置 | 与 UI 无关，重构基本不动 |
| `hooks` | 数据层（Query + 链上交互） | 与 UI 解耦良好，可整体复用 |
| `services` | API + mock 数据访问 | 与 UI 解耦良好，可整体复用 |
| `pages` | 页面级组合 | 重构主战场，含较多内联/硬编码样式 |
| `lib`/`utils` | 纯工具函数 | 全部可复用 |

## 3. 路由清单（来自 `App.tsx`）

react-router 7（`react-router-dom`），`BrowserRouter` 在 `main.tsx`，页面全部 `React.lazy` + `Suspense`（fallback 为简单 Loading）。

| 路径 | 页面组件 | 布局 | 守卫 | 说明 |
|------|----------|------|------|------|
| `/` | `Landing` | 无（裸路由） | 无 | 落地页/连接钱包入口 |
| `/dashboard` | `Dashboard` | `AppLayout` | 需连接+已注册 | 首页总览 |
| `/deposit` | `Deposit` | `AppLayout` | 同上；suspended 拦截 | 充值/提现 |
| `/services` | `Services` | `AppLayout` | 同上；inactive/suspended 拦截 | 地区/运营商列表 |
| `/services/:regionCode` | `RegionDetail` | `AppLayout` | 同上 | 地区运营商详情/申请号码 |
| `/billing` | `Billing` | `AppLayout` | 同上；inactive 拦截 | 账单列表 |
| `/billing/:billId` | `BillDetail` | `AppLayout` | 同上 | 账单详情/支付 |
| `/notifications` | `Notifications` | `AppLayout` | 同上；inactive 拦截 | 通知中心 |
| `/cards` | `Cards` | `AppLayout` | 需连接+已注册 | 流量卡 NFT |

**守卫逻辑**（`AppLayout.tsx`）：未连接钱包 → 跳 `/`；查不到 user → 跳 `/`；`inactive` 访问 services/billing/notifications → GuardCard(去充值)；`suspended` 访问 services/deposit → GuardCard(去还款)。所有受守卫页面套在 `max-w-mobile`(430px) 移动端容器内，含 Header + 底部 TabBar。

## 4. 关键依赖说明（Web3 栈接入）

### Provider 树（`main.tsx`）
```
React.StrictMode
└── WagmiProvider (config=wagmiConfig)
    └── QueryClientProvider (staleTime 30s, retry 1)
        └── RainbowKitProvider (darkTheme, accentColor="#3b82f6")  ← 硬编码强调色
            └── BrowserRouter
                ├── App (路由)
                └── Toaster (sonner, top-center, dark, richColors)
```

### wagmi / RainbowKit 配置（`config/wagmi.ts`）
- 用 `@rainbow-me/rainbowkit` 的 `getDefaultConfig`，`appName: "LinkWorld"`，含 WalletConnect `projectId`。
- 链选择：`import.meta.env.VITE_CHAIN_ID === "31337"` → `hardhatLocal`，否则 `zgTestnet`。
- transport：`http()` 默认 RPC。

### 链定义（`config/chains.ts`）
- `zgMainnet`（16600，A0GI）/ `zgTestnet`（16601，A0GI）/ `hardhatLocal`（31337，ETH，http://127.0.0.1:8545）。

### 合约地址（`config/contracts.ts`）
- `CONTRACTS: Record<chainId, ContractAddresses>`，含 7 个合约。31337 已填地址，16601 全为零地址（TODO 待测试网部署）。
- `getContractAddress(chainId, name)`：链不存在或地址为零会抛错。

### 链上 hook（`hooks/contracts/`）
- 基于 wagmi `useReadContract` / `useWriteContract` / `useWaitForTransactionReceipt` / `useChainId`，配合 `config/abis`（桶导出 6 个 ABI，注意 Oracle ABI 文件存在但未在 index 桶导出）与 `getContractAddress`。

## 5. 状态管理方式

- **服务端状态**：TanStack Query v5。全局 `QueryClient`（staleTime 30s，retry 1）。业务 hook 统一用 `useQuery`/`useMutation` 封装，按 query key 缓存。
- **链上状态**：wagmi hooks（`useAccount`、读写合约），与 Query 并存。
- **本地 UI 状态**：组件内 `useState`（如 RegisterSheet 的 step/email/code）。
- **离线/待同步**：`utils/pendingSync.ts` 用 localStorage 暂存待回传后端的数据，配合指数退避重试。
- **无 Redux / Zustand / 全局 Context**：跨组件共享一律走 Query 缓存或 wagmi。

## 6. 关键依赖版本（`package.json`）

| 类别 | 包 | 版本 |
|------|----|------|
| 框架 | react / react-dom | ^19.0.0 |
| 构建 | vite ^6.0.0；@vitejs/plugin-react ^4.3 | — |
| 语言 | typescript ^5.6 | — |
| 样式 | tailwindcss ^3.4.17；tw-animate-css ^1.4；class-variance-authority ^0.7；clsx ^2.1；tailwind-merge ^3.5 | — |
| UI | @base-ui/react ^1.3；shadcn ^4.2；lucide-react ^1.8（已装但页面多用 emoji 图标）；vaul ^1.1（抽屉）；sonner ^2.0（toast） | — |
| 字体 | @fontsource-variable/geist ^5.2 | — |
| Web3 | wagmi ^2.14；viem ^2.21；@rainbow-me/rainbowkit ^2.2 | — |
| 数据 | @tanstack/react-query ^5.62；axios ^1.15 | — |
| 路由 | react-router-dom ^7.1 | — |

### 别名（`components.json` / vite）
`@/components` `@/lib` `@/lib/utils` `@/components/ui` `@/hooks`。shadcn style=`base-nova`，baseColor=`neutral`，cssVariables=true，iconLibrary=lucide。
