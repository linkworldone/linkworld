# utils.md — Link World Web 基线·公共方法/工具清单

> 覆盖 `lib/` `utils/` `hooks/` `services/` `config/`
> 用途：重构时优先复用，避免重复造轮子。这些层与 UI 解耦良好，基本可整体保留。
> 扫描时间：2026-06-06

## 1. `lib/` —— 类名工具

| 文件 | 函数 | 签名 | 用途 |
|------|------|------|------|
| `lib/utils.ts` | `cn` | `cn(...inputs: ClassValue[]): string` | clsx + tailwind-merge 合并/去重 Tailwind 类名（shadcn 标准） |

## 2. `utils/` —— 格式化 & 离线同步

### `utils/format.ts`
| 函数 | 签名 | 用途 |
|------|------|------|
| `shortenAddress` | `(address: string, chars = 4) => string` | 钱包地址缩写 `0x12...abcd` |
| `parseUnits` | `(value: string, decimals = 18) => bigint` | 字符串金额 → wei（手写，非 viem） |
| `formatAmount` | `(wei: bigint, decimals = 18, displayDecimals = 2) => string` | wei → 显示金额字符串 |
| `formatUSD` | `(amount: string) => string` | 数字字符串 → `$x.xx` |
| `formatDate` | `(isoDate: string) => string` | ISO → `Jun 6, 2026`（en-US） |
| `timeAgo` | `(isoDate: string) => string` | 相对时间（just now / x minutes ago / …，超 1 天回退 formatDate） |

### `utils/pendingSync.ts`（localStorage 待同步队列，前缀 `linkworld_pending_`）
| 函数 | 签名 | 用途 |
|------|------|------|
| `savePendingSync` | `<T>(key: string, data: T) => void` | 暂存待回传后端的数据（含 createdAt/retryCount） |
| `getPendingSync` | `<T>(key: string) => PendingSyncItem<T> \| null` | 读取暂存项 |
| `clearPendingSync` | `(key: string) => void` | 删除暂存项 |
| `incrementRetryCount` | `(key: string) => number` | 重试次数 +1 并返回 |
| `retryWithBackoff` | `(fn: () => Promise<void>, maxRetries = 3, baseDelay = 2000) => Promise<boolean>` | 指数退避重试，全部失败返回 false |
| `PendingSyncItem<T>` | `interface { data: T; createdAt: number; retryCount: number }` | 暂存项类型 |

## 3. `config/` —— 配置常量与函数

| 文件 | 导出 | 签名/值 | 用途 |
|------|------|---------|------|
| `config/api.ts` | `API_BASE_URL` | `string`（env `VITE_API_BASE_URL` ?? `http://localhost:8080`） | 后端基址 |
| `config/constants.ts` | `PLATFORM_FEE_RATE` | `0.025` | 平台费率 |
| | `MIN_DEPOSIT_USDT` | `100n * 10n**18n` | 最低押金 |
| | `OVERDUE_DAYS` | `14` | 账单逾期天数 |
| | `MOCK_DELAY_MS` | `600` | mock 延迟 |
| | `SUPPORTED_CURRENCIES` / `SupportedCurrency` | `["USDT","ETH"]` const + 派生 type | 支持币种 |
| `config/chains.ts` | `zgMainnet` / `zgTestnet` / `hardhatLocal` | viem `defineChain` 结果 | 链定义（16600/16601/31337） |
| `config/contracts.ts` | `CONTRACTS` | `Record<number, ContractAddresses>` | 各链 7 合约地址表 |
| | `getContractAddress` | `(chainId: number, name: keyof ContractAddresses) => \`0x${string}\`` | 取地址，零地址/缺失抛错 |
| `config/wagmi.ts` | `wagmiConfig` | RainbowKit `getDefaultConfig` 结果 | wagmi 全局配置 |
| `config/abis/index.ts` | `UserRegistryABI` `DepositABI` `ServiceManagerABI` `PaymentABI` `FeeManagerABI` `TrafficCardNFTABI` | ABI 数组 | 合约 ABI 桶导出（注意 OracleABI 文件存在但未在此 index 导出） |

## 4. `hooks/` —— 业务数据 hook（TanStack Query 封装）

| 文件 | hook | 签名概要 | 用途 |
|------|------|----------|------|
| `useUser.ts` | `useUser` | `(address?) => Query<User\|null>` | 拉取用户 |
| | `useRegister` | `() => { register, backendSync, isContractPending, isSuccess }` | 合约注册→后端同步（含 useEffect 自动回传） |
| | `useSendVerificationCode` | `() => Mutation<{address,email}>` | 发邮箱验证码 |
| | `useVerifyEmail` | `() => Mutation<{address,code}> => boolean` | 校验验证码 |
| `useDeposit.ts` | `useDeposit` | `(address?) => Query<DepositInfo>` | 押金余额 |
| | `useDepositHistory` | `(address?) => Query<DepositRecord[]>` | 充值/提现历史 |
| | `useDepositMutation` | `() => Mutation` | 充值（链上+后端记录） |
| | `useWithdrawMutation` | `() => Mutation` | 提现 |
| `useBilling.ts` | `useBills` | `(address?, filter?: "unpaid"\|"paid") => Query<Bill[]>` | 账单列表 |
| | `useBillDetail` | `(billId?) => Query<Bill>` | 账单详情 |
| | `useMonthEstimate` | `(address?) => Query<MonthEstimate>` | 当月预估 |
| | `usePayBill` | `() => Mutation` | 支付账单 |
| `useOperator.ts` | `useRegions` | `() => Query<Region[]>` | 地区列表 |
| | `useOperatorsByRegion` | `(regionCode?) => Query<Operator[]>` | 地区运营商 |
| | `useMyNumbers` | `(address?) => Query<VirtualNumber[]>` | 我的虚拟号码 |
| | `useApplyNumber` | `() => Mutation` | 申请号码 |
| `useNotification.ts` | `useNotifications` | `(address?) => Query<Notification[]>` | 通知列表（mock） |
| | `useUnreadCount` | `(address?) => Query<number>` | 未读数 |
| | `useMarkAsRead` / `useMarkAllAsRead` | `() => Mutation` | 标记已读 |
| `useTransactionFlow.ts` | `parseContractError` | `(error: unknown) => string` | 合约错误 → 友好文案 |
| | `useTxState` | `(params) => 交易状态` | 统一交易状态管理 |

### 链上 hook（`hooks/contracts/`，wagmi useReadContract/useWriteContract）
| 文件 | hook | 用途 |
|------|------|------|
| `useUserRegistry.ts` | `useIsRegistered(address)` / `useContractRegister()` | 注册状态/链上注册 |
| `useDepositContract.ts` | `useDepositBalance(address)` / `useContractDeposit()` / `useContractWithdraw()` | 押金读写 |
| `usePaymentContract.ts` | `useContractPayBill()` | 链上付账单 |
| `useServiceManager.ts` | `useContractUserService(address)` / `useContractActivateService()` | 服务激活 |
| `useTrafficCard.ts` | `useTrafficCardCredit(address)` / `useTrafficCards(address)` / `useBurnCard()` / `useIssueMonthlyCards()` | 流量卡 NFT（最大/最新模块，21 符号） |
| `contracts/index.ts` | 桶导出 | — |

## 5. `services/` —— API & mock 数据层

### `services/api/`（真实后端，axios）
| 文件 | 导出 | 主要方法 | 用途 |
|------|------|----------|------|
| `client.ts` | `apiClient` | axios 实例（baseURL=API_BASE_URL，timeout 10s，响应拦截器直返 data、统一错误文案） | HTTP 客户端 |
| `userApi.ts` | `userApi` | `getUser(wallet)` / `register(wallet,email,...)`（含 ApiUser→User snake→camel 映射） | 用户 |
| `operatorApi.ts` | `operatorApi` | `getRegions()` / `getOperatorsByRegion(regionCode)`（含 FLAG_MAP 国旗 emoji、按国家分组） | 运营商/地区 |
| `depositApi.ts` | `depositApi` | `getHistory(wallet)` / `recordDeposit(wallet,amount)` / `getDepositAmount(wallet)` / `recordWithdraw(wallet,...)` | 押金记录 |
| `billingApi.ts` | `billingApi` | `getBills(wallet,filter?)` / `recordPayment(wallet,billId,...)`（含 ApiBill→Bill 映射、费用合计） | 账单 |
| `usageApi.ts` | `usageApi` | `getUsage(wallet) => MonthEstimate` | 用量 |
| `index.ts` | 桶导出 5 个 *Api | — | — |

### `services/mock/`（mock 数据，仍在用的只剩 notification）
| 文件 | 导出 | 用途 |
|------|------|------|
| `notificationService.ts` | `notificationService` | 通知 mock 服务（后端暂无通知 API，`services/index.ts` 仍导出它） |
| `delay.ts` | `delay(ms = MOCK_DELAY_MS)` | 模拟网络延迟 |
| `data.ts` | mock 数据集（18 符号） | 旧 mock 数据源 |
| `userService` / `depositService` / `billingService` / `operatorService` | 各 mock 服务 | **已被 api/ 取代，仅 notification 仍引用**，重构可清理 |
| `services/index.ts` | 统一出口：5 个真实 *Api + `notificationService`(mock) | 服务层入口 |

## 6. `types/index.ts` —— 领域 & API 类型

- 领域（camelCase）：`UserStatus` `User` `DepositInfo` `DepositRecord` `Region` `Operator` `VirtualNumber` `Bill` `MonthEstimate` `Notification`
- 后端响应（snake_case）：`ApiUser` `ApiOperator` `ApiBill` `ApiUsage` —— 在各 *Api.ts 内做 snake→camel 映射。

## 7. 重构提示（工具层）

- 工具/hook/service/config 层**与 UI 解耦良好，基本可整体复用**，重构聚焦 components/pages 即可。
- `services/mock/` 下 user/deposit/billing/operator 四个 service 已被真实 API 取代但文件仍在，可在重构时清理（仅 notification 仍 mock）。
- `config/abis` 桶导出漏了 OracleABI（文件存在），若有链上 Oracle 调用需求需补；同时 `useTrafficCard` 是最新最复杂的链上模块，重构 Cards 页时重点回归。
