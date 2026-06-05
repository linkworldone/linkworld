# Link World Web — 项目基线扫描 (project-scan)

> 本文档反映 **当前 HEAD** 的 `packages/web` 真实代码状态（shadcn + oklch 暗色主题）。
> scan 阶段仅做现状基线，不含任何改动或设计建议。

## 1. 技术栈版本表

版本号取自 `packages/web/package.json`（均为 semver caret 区间）。

### 运行时依赖

| 包 | 版本 | 用途 |
|----|------|------|
| react | ^19.0.0 | UI 框架 |
| react-dom | ^19.0.0 | DOM 渲染 |
| react-router-dom | ^7.1.0 | 路由（BrowserRouter + lazy Routes） |
| wagmi | ^2.14.0 | 以太坊 React hooks（账户/合约读写） |
| viem | ^2.21.0 | EVM 底层客户端（wagmi 依赖） |
| @rainbow-me/rainbowkit | ^2.2.0 | 钱包连接 UI（darkTheme） |
| @tanstack/react-query | ^5.62.0 | 服务端状态/缓存（所有数据 hooks 基座） |
| axios | ^1.15.0 | HTTP 客户端（services/api） |
| @base-ui/react | ^1.3.0 | 无样式组件原语（button/badge/tabs 基于此） |
| shadcn | ^4.2.0 | 设计系统 CSS（`shadcn/tailwind.css` 被 index.css import） |
| class-variance-authority | ^0.7.1 | 变体类生成（cva） |
| clsx | ^2.1.1 | className 拼接 |
| tailwind-merge | ^3.5.0 | Tailwind class 去冲突（cn 工具） |
| tw-animate-css | ^1.4.0 | 动画类（index.css 顶部 import） |
| vaul | ^1.1.2 | 底部抽屉 Drawer（BottomSheet / RegisterSheet） |
| sonner | ^2.0.7 | Toast 通知（main.tsx Toaster theme=dark） |
| lucide-react | ^1.8.0 | 图标库 |
| @fontsource-variable/geist | ^5.2.8 | Geist Variable 字体（index.css import + --font-sans） |

### 开发依赖

| 包 | 版本 | 用途 |
|----|------|------|
| typescript | ^5.6.0 | 类型系统 |
| vite | ^6.0.0 | 构建/dev server |
| @vitejs/plugin-react | ^4.3.0 | React 插件 |
| tailwindcss | ^3.4.17 | 原子化 CSS（注意：**v3**，非 v4） |
| postcss | ^8.4.49 | CSS 处理 |
| autoprefixer | ^10.4.20 | 浏览器前缀 |
| @types/react / @types/react-dom | ^19.0.0 | 类型定义 |

> 关键提示：`tailwindcss` 是 **3.4**，但 `index.css` 通过 `@import "shadcn/tailwind.css"` 引入了 shadcn 的 token 层，并用 `oklch` 定义 `:root` 变量。tailwind.config.ts 又用 **hex** 定义了一套并行的语义色（surface/brand/status/text）。两套色彩体系并存（详见 color-mapping.md）。

## 2. 目录结构概览

`packages/web/src/` 共 69 个索引文件，分层如下：

```
src/
├── App.tsx                 # 路由表（lazy + Suspense）
├── main.tsx                # 应用入口 + Provider 栈
├── index.css               # 全局样式 + oklch :root 变量
├── vite-env.d.ts
├── components/
│   ├── ui/                 # shadcn 原语：badge / button / tabs
│   ├── shared/             # 业务通用：AmountDisplay / BottomSheet / EmptyState / GuardCard / StatusBadge
│   ├── layout/             # AppLayout / Header / TabBar
│   └── wallet/             # ConnectButton / RegisterSheet
├── pages/                  # 9 个页面组件（懒加载）
├── hooks/                  # 业务数据 hooks（react-query）
│   └── contracts/          # 链上合约 hooks（wagmi）
├── services/
│   ├── api/                # axios 真实后端 API
│   └── mock/               # 内存 mock 数据/服务
├── config/                 # 链/合约/wagmi/api 配置 + abis/
├── types/index.ts          # 全局领域类型
├── lib/utils.ts            # cn()
└── utils/                  # format.ts / pendingSync.ts
```

## 3. 路由清单

入口在 `main.tsx`（`BrowserRouter` 包裹 `App`），路由表在 `App.tsx`（全部页面 `lazy` + `Suspense fallback`）。

| 路径 | 页面组件 | 布局 | 说明 |
|------|----------|------|------|
| `/` | `Landing` | 无（裸页） | 落地页 / 连接钱包入口 |
| `/dashboard` | `Dashboard` | AppLayout | 首页（Header 显示欢迎语 + 地址 + 通知铃铛） |
| `/deposit` | `Deposit` | AppLayout | 充值/余额 |
| `/services` | `Services` | AppLayout | 服务/地区列表 |
| `/services/:regionCode` | `RegionDetail` | AppLayout | 地区下运营商详情（子页，返回按钮） |
| `/billing` | `Billing` | AppLayout | 账单列表 |
| `/billing/:billId` | `BillDetail` | AppLayout | 账单详情（子页，返回按钮） |
| `/notifications` | `Notifications` | AppLayout | 通知列表 |
| `/cards` | `Cards` | AppLayout | 流量卡 NFT |

路由特征：
- 共 **9 条 Route**（1 个裸页 + 8 个 AppLayout 嵌套子路由）。
- 除 `/` 外，所有页面包在 `<Route element={<AppLayout />}>` 嵌套路由下。
- `AppLayout` 内做 **守卫**：未连接钱包 → `Navigate("/")`；无 user → `Navigate("/")`；`status==="inactive"` 时 `/services|/billing|/notifications` 被 GuardCard 拦截；`status==="suspended"` 时 `/deposit|/services` 被拦截。
- 子页判定靠 `location.pathname.includes("/services/" | "/billing/")`（Header）。
- TabBar 5 个 tab：Home `/dashboard`、Services `/services`、Deposit `/deposit`、Bills `/billing`(带未付账单 badge)、Cards `/cards`。

## 4. 数据流概览

三层分离：**hooks（react-query 封装） → services（数据源） → config/abis（链上）**。

### Provider 栈（main.tsx，由外到内）
`WagmiProvider(wagmiConfig)` → `QueryClientProvider(queryClient)` → `RainbowKitProvider(darkTheme accentColor #3b82f6)` → `BrowserRouter` → `App` + `Toaster(theme=dark, position=top-center, richColors)`。

`QueryClient` 默认：`staleTime: 30_000`，`retry: 1`。

### Web3 链 / wagmi 配置（config/chains.ts + config/wagmi.ts）
`config/chains.ts` 用 viem `defineChain` 定义 3 条链：

| 导出 | chainId | 名称 | nativeCurrency | RPC | 区块浏览器 |
|------|---------|------|----------------|-----|-----------|
| `zgMainnet` | 16600 | 0G Newton Mainnet | A0GI (18) | `https://evmrpc.0g.ai` | `https://chainscan.0g.ai` |
| `zgTestnet` | 16601 | 0G Galileo Testnet | A0GI (18) | `https://evmrpc-testnet.0g.ai` | `https://chainscan-galileo.0g.ai` |
| `hardhatLocal` | 31337 | Hardhat Local | ETH (18) | `http://127.0.0.1:8545` | 无 |

`config/wagmi.ts`：用 RainbowKit `getDefaultConfig` 构建 `wagmiConfig`，`appName: "LinkWorld"`，`projectId: "21fef48091f12692cad574a6f7753643"`（WalletConnect，硬编码）。**单链运行**——按 `VITE_CHAIN_ID === "31337"` 选 `hardhatLocal`，否则 `zgTestnet`（`zgMainnet` 定义了但当前未接入 chains 数组）。transports 对选中链用 `http()` 默认 RPC。

### 全局配置文件清单

| 文件 | 关键内容 |
|------|----------|
| `tailwind.config.ts` | `darkMode: "class"`；`content` 扫 index.html + src；`theme.extend` 含 hex colors（surface/brand/status/text/border）、`fontFamily.orbitron`、`maxWidth.mobile: 430px`；`plugins: []` |
| `src/index.css` | import `tw-animate-css` / `shadcn/tailwind.css` / `@fontsource-variable/geist` + `@tailwind base/components/utilities`；`.theme` 字体 token；`:root` oklch 暗色 token；`@layer base` 全局规则（body `bg-surface text-text-primary`、font-family 写死 Inter） |
| `src/config/api.ts` | `API_BASE_URL = import.meta.env.VITE_API_BASE_URL \|\| "http://localhost:8080"` |
| `src/config/contracts.ts` | `CONTRACTS: Record<chainId, {7 个合约地址}>` + `getContractAddress(chainId, name)`（地址为零或缺失时抛错）。已填 31337(Hardhat) 真实地址；16601(0G Testnet) 全为 `0x0`（标注 TODO 待部署）。被 contracts hooks 读取 |
| `src/config/abis/` | 7 个 ABI 文件：UserRegistry / Deposit / Payment / ServiceManager / TrafficCardNFT / FeeManager / Oracle。但 `index.ts` 只 re-export 6 个（**Oracle ABI 未在 index 导出**） |

> 环境变量：`VITE_CHAIN_ID`（选链）、`VITE_API_BASE_URL`（后端地址）。tsconfig 路径别名 `@/` → `src/`（App.tsx/hooks 普遍使用）。

### services 层
- **services/api/**：真实后端，基于 `axios`。`client.ts` 创建 `apiClient`（baseURL=`API_BASE_URL`，timeout 10s，响应拦截器直接返回 `res.data`、错误归一为 `Error(message)`）。`API_BASE_URL` 取 `import.meta.env.VITE_API_BASE_URL || "http://localhost:8080"`（config/api.ts）。模块：`userApi / operatorApi / depositApi / billingApi / usageApi`，各自含 `toXxx` DTO→领域类型转换（如 `userApi.toUser` 把 `wallet_addr/is_active` 映射为 `address/status`）。
- **services/mock/**：内存 mock。`data.ts` 持有可变数组与 getter/setter（`getMockUser/setMockUser/addMockNumber/updateMockBill...`），`delay.ts` 模拟网络延迟。服务模块：`billingService / depositService / notificationService / operatorService / userService`。
- `services/index.ts` 当前 **混合导出**：api 系列（userApi/operatorApi/depositApi/billingApi/usageApi）+ mock 的 `notificationService`。

### hooks 层
- **业务 hooks**（`src/hooks/*.ts`，react-query）：`useUser`(+useRegister/useSendVerificationCode/useVerifyEmail)、`useNotification`(useNotifications/useUnreadCount/useMarkAsRead/useMarkAllAsRead)、`useBilling`(useBills/useBillDetail/useMonthEstimate/usePayBill)、`useDeposit`(useDeposit/useDepositHistory/useDepositMutation/useWithdrawMutation)、`useOperator`(useRegions/useOperatorsByRegion/useMyNumbers/useApplyNumber)、`useTransactionFlow`(parseContractError/useTxState)。
- **合约 hooks**（`src/hooks/contracts/*.ts`，wagmi）：`useUserRegistry`、`useDepositContract`、`usePaymentContract`、`useServiceManager`、`useTrafficCard`。读 `config/contracts.ts` 地址 + `config/abis/` ABI。

### 组件消费
组件直接调 hooks 取数（如 `AppLayout→useUser`、`Header→useUnreadCount`、`TabBar→useBills`、`RegisterSheet→useRegister/useSendVerificationCode/useVerifyEmail`），不直接碰 services。
</content>
