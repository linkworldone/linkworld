# Link World Web — 公共方法 / Hooks / Services 清单 (utils)

> 反映 **当前 HEAD** 代码。覆盖：`lib/utils.ts`、`utils/format.ts`、`utils/pendingSync.ts`、业务 hooks、合约 hooks、services/api 与 services/mock。

## 1. lib/utils.ts

| 导出 | 签名 | 用途 |
|------|------|------|
| `cn` | `cn(...inputs: ClassValue[]): string` | `twMerge(clsx(inputs))` —— 合并/去冲突 Tailwind className。ui 原语普遍使用。 |

## 2. utils/format.ts

| 导出 | 签名 | 用途 |
|------|------|------|
| `shortenAddress` | `(address: string, chars = 4) => string` | 地址缩写，如 `0x1234...abcd`（取前 `chars+2` 与后 `chars`） |
| `parseUnits` | `(value: string, decimals = 18) => bigint` | 小数字符串 → 最小单位 bigint |
| `formatAmount` | `(wei: bigint, decimals = 18, displayDecimals = 2) => string` | 最小单位 bigint → 展示用小数字符串 |
| `formatUSD` | `(amount: string) => string` | `parseFloat` 后 → `$x.xx` |
| `formatDate` | `(isoDate: string) => string` | ISO → `en-US` 短日期（Mon DD, YYYY） |
| `timeAgo` | `(isoDate: string) => string` | 相对时间（just now / N minutes ago / N hours ago / 超 1 天回退 formatDate） |

## 3. utils/pendingSync.ts

localStorage 持久化的「待同步」重试机制，key 前缀 `linkworld_pending_`。

| 导出 | 签名 | 用途 |
|------|------|------|
| `PendingSyncItem<T>` | `interface { data: T; createdAt: number; retryCount: number }` | 待同步项结构 |
| `savePendingSync` | `<T>(key: string, data: T) => void` | 写入待同步项（try/catch 静默失败） |
| `getPendingSync` | `<T>(key: string) => PendingSyncItem<T> \| null` | 读取（解析失败返回 null） |
| `clearPendingSync` | `(key: string) => void` | 删除 |
| `incrementRetryCount` | `(key: string) => number` | retryCount +1 并回写，返回新值 |
| `retryWithBackoff` | `(fn: () => Promise<void>, maxRetries = 3, baseDelay = 2000) => Promise<boolean>` | 指数退避重试，全部失败返回 false |

## 4. 业务数据 hooks（src/hooks/*.ts，react-query 封装）

### useUser.ts
| 导出 | 签名 | 用途 |
|------|------|------|
| `useUser` | `(address?: string)` | 查询用户（query） |
| `useRegister` | `()` | 合约注册 + 成功后自动 `backendSync`（返回 `register/backendSync/isContractPending/isSuccess`） |
| `useSendVerificationCode` | `()` | 发送邮箱验证码（mutation） |
| `useVerifyEmail` | `()` | 校验验证码（mutation，返回是否通过） |

### useNotification.ts
| 导出 | 签名 | 用途 |
|------|------|------|
| `useNotifications` | `(address: string \| undefined)` | 通知列表 |
| `useUnreadCount` | `(address: string \| undefined)` | 未读数（Header 铃铛用） |
| `useMarkAsRead` | `()` | 标记单条已读（mutation） |
| `useMarkAllAsRead` | `()` | 全部已读（mutation） |

### useBilling.ts
| 导出 | 签名 | 用途 |
|------|------|------|
| `useBills` | `(address?: string, filter?: "unpaid" \| "paid")` | 账单列表（TabBar 用 unpaid 算 badge） |
| `useBillDetail` | `(billId?: string)` | 账单详情 |
| `useMonthEstimate` | `(address?: string)` | 当月预估用量/费用 |
| `usePayBill` | `()` | 支付账单（mutation） |

### useDeposit.ts
| 导出 | 签名 | 用途 |
|------|------|------|
| `useDeposit` | `(address?: string)` | 押金/余额信息 |
| `useDepositHistory` | `(address?: string)` | 充值/提现历史 |
| `useDepositMutation` | `()` | 充值（mutation） |
| `useWithdrawMutation` | `()` | 提现（mutation） |

### useOperator.ts
| 导出 | 签名 | 用途 |
|------|------|------|
| `useRegions` | `()` | 地区列表 |
| `useOperatorsByRegion` | `(regionCode?: string)` | 某地区运营商 |
| `useMyNumbers` | `(address?: string)` | 我的虚拟号码 |
| `useApplyNumber` | `()` | 申请号码（mutation） |

### useTransactionFlow.ts
| 导出 | 签名 | 用途 |
|------|------|------|
| `TxStatus` | `type = "idle" \| "pending-signature" \| "pending-confirmation" \| "success" \| "error"` | 交易状态枚举 |
| `TxState` | `interface` | 交易状态对象结构（status 等字段） |
| `parseContractError` | `(error: unknown) => string` | 合约错误归一为可读文案 |
| `useTxState` | `(params: {...})` | 交易状态机封装（pending/success/error 流转，配合 toast） |

## 5. 合约 hooks（src/hooks/contracts/*.ts，wagmi）

统一通过 `hooks/contracts/index.ts` re-export。地址读 `config/contracts.ts`，ABI 读 `config/abis/`。

| 文件 | 导出 | 签名/形态 | 用途 |
|------|------|-----------|------|
| useUserRegistry.ts | `useIsRegistered` | `(address: \`0x${string}\` \| undefined)` | 链上是否已注册（read） |
| useUserRegistry.ts | `useContractRegister` | `()` | 链上注册（write） |
| useDepositContract.ts | `useDepositBalance` | `(address: \`0x${string}\` \| undefined)` | 链上押金余额（read） |
| useDepositContract.ts | `useContractDeposit` | `()` | 链上充值（write） |
| useDepositContract.ts | `useContractWithdraw` | `()` | 链上提现（write） |
| usePaymentContract.ts | `useContractPayBill` | `()` | 链上支付账单（write） |
| useServiceManager.ts | `useContractUserService` | `(address: \`0x${string}\` \| undefined)` | 用户服务状态（read） |
| useServiceManager.ts | `useContractActivateService` | `()` | 激活服务（write） |
| useTrafficCard.ts | `useTrafficCardCredit` | `(address: \`0x${string}\` \| undefined)` | 流量卡额度（read） |
| useTrafficCard.ts | `useTrafficCards` | `(address: \`0x${string}\` \| undefined)` | 持有的流量卡 NFT 列表（read） |
| useTrafficCard.ts | `useBurnCard` | `()` | 销卡（write） |
| useTrafficCard.ts | `useIssueMonthlyCards` | `()` | 发月卡（write） |
| useTrafficCard.ts | `TrafficCardItem` | `type` | 流量卡条目类型（导出供页面用） |

## 6. services/api（真实后端，axios）

`client.ts`：`apiClient`（baseURL=`API_BASE_URL`，timeout 10s，响应拦截器返回 `res.data`、错误归一 `Error`）。各模块含 DTO→领域类型转换。

| 模块 | 方法 | 签名 |
|------|------|------|
| `userApi` | getUser | `(wallet: string) => Promise<User \| null>`（内部 `toUser` 转换 snake_case DTO） |
| `userApi` | register | `(wallet: string, email: string, tokenId?: number) => Promise<void>` |
| `operatorApi` | getRegions | `() => Promise<Region[]>` |
| `operatorApi` | getOperatorsByRegion | `(regionCode: string) => Promise<Operator[]>` |
| `depositApi` | getHistory | `(wallet: string) => Promise<DepositRecord[]>` |
| `depositApi` | recordDeposit | `(...)` 记录充值 |
| `depositApi` | getDepositAmount | `(wallet: string) => Promise<string>` |
| `depositApi` | recordWithdraw | `(...)` 记录提现 |
| `billingApi` | getBills | `(...)` 账单列表 |
| `billingApi` | recordPayment | `(...)` 记录支付 |
| `usageApi` | getUsage | `(wallet: string) => Promise<MonthEstimate>` |

`services/api/index.ts` re-export：`userApi / operatorApi / depositApi / billingApi / usageApi`。

> 注意：号码申请/我的号码等并不在 operatorApi 中（仅 getRegions / getOperatorsByRegion），其余靠 mock 层补齐。

## 7. services/mock（内存 mock）

| 文件 | 导出形态 | 用途 |
|------|----------|------|
| `data.ts` | 可变数据 + getter/setter：`mockDepositHistory / mockRegions / mockOperators / mockNumbers / mockBills / mockNotifications`；`getMockUser/setMockUser`、`getMockDeposit/setMockDeposit`、`addMockDepositRecord`、`addMockNumber`、`updateMockBill`、`addMockNotification` | 内存数据源（可变状态模拟后端） |
| `delay.ts` | `delay(ms = MOCK_DELAY_MS): Promise<void>` + 常量 | 模拟网络延迟 |
| `userService.ts` | `userService` 对象 | 用户相关 mock 服务 |
| `depositService.ts` | `depositService` 对象 | 押金相关 mock 服务 |
| `billingService.ts` | `billingService` 对象 | 账单相关 mock 服务 |
| `operatorService.ts` | `operatorService` 对象 | 运营商/号码相关 mock 服务 |
| `notificationService.ts` | `notificationService` 对象 | 通知相关 mock 服务 |

`services/index.ts` 当前混合 re-export：`userApi/operatorApi/depositApi/billingApi/usageApi`（来自 api）+ `notificationService`（来自 mock）—— **api 与 mock 混用**，是数据层一个待统一的点。
</content>
