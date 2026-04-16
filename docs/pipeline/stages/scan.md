# Scan 阶段产出 — LinkWorld 全栈基线

> Round 2 | 2026-04-14

## 一、项目结构

```
linkworld/ (pnpm monorepo)
├── packages/web/          # React 19 + Vite 6 前端
├── packages/backend/      # Go + Gin + GORM 后端
├── packages/contracts/    # Solidity 0.8.27 + Hardhat 合约
└── packages/server/       # NestJS（本轮暂不涉及）
```

## 二、前端基线 (packages/web)

### 页面 (8)
| 页面 | 路由 | 功能 |
|------|------|------|
| Landing | / | 公开首页，ConnectButton + RegisterSheet |
| Dashboard | /dashboard | 主控台：余额、号码、用量、快捷操作 |
| Deposit | /deposit | 充值/提现 + 历史记录 |
| Services | /services | 地区浏览 + 我的号码 |
| RegionDetail | /services/:regionCode | 运营商列表 + 申请号码 |
| Billing | /billing | 账单列表（未付/历史） |
| BillDetail | /billing/:billId | 单笔账单详情 + 支付 |
| Notifications | /notifications | 通知中心 |

### 组件 (13)
- Layout: AppLayout, Header, TabBar
- Shared: GuardCard, BottomSheet, StatusBadge, AmountDisplay, EmptyState
- Wallet: ConnectButton, RegisterSheet
- UI: Button, Badge, Tabs (shadcn/base-ui)

### Hooks (5) — 全部包装 Mock Service
| Hook | 服务 | 查询/变更 |
|------|------|-----------|
| useUser | userService | useUser, useRegister, useSendVerificationCode, useVerifyEmail |
| useOperator | operatorService | useRegions, useOperatorsByRegion, useMyNumbers, useApplyNumber |
| useDeposit | depositService | useDeposit, useDepositHistory, useDepositMutation, useWithdrawMutation |
| useBilling | billingService | useBills, useBillDetail, useMonthEstimate, usePayBill |
| useNotification | notificationService | useNotifications, useUnreadCount, useMarkAsRead, useMarkAllAsRead |

### Services — 全部 Mock
5 个 mock service + delay helper，存于 `src/services/mock/`，通过 `src/services/index.ts` 统一导出。

### 类型 (9 interfaces)
User, DepositInfo, DepositRecord, Region, Operator, VirtualNumber, Bill, MonthEstimate, Notification

### Provider 树
WagmiProvider → QueryClientProvider → RainbowKitProvider → BrowserRouter → App

## 三、后端基线 (packages/backend)

### API 端点 (15)
| Method | Path | Handler | 状态 |
|--------|------|---------|------|
| POST | /api/register | Register | ✅ 实现 |
| GET | /api/user/:wallet | GetUser | ✅ 实现 |
| GET | /api/operators | GetOperators | ✅ 实现 |
| POST | /api/service/activate | ActivateService | ✅ 实现 |
| POST | /api/service/deactivate | DeactivateService | ✅ 实现 |
| GET | /api/service/:wallet | GetUserService | ✅ 实现 |
| GET | /api/bills/:wallet | GetBills | ✅ 实现 |
| POST | /api/bills/pay | PayBill | ⚠️ Stub |
| POST | /api/deposit | Deposit | ✅ 实现 |
| GET | /api/deposit/:wallet | GetDeposit | ✅ 实现 |
| POST | /api/withdraw | Withdraw | ⚠️ Stub |
| POST | /api/virtual-number/generate | GenerateVirtualNumber | ✅ 实现 |
| GET | /api/countries | GetCountryList | ✅ 实现 |
| GET | /api/usage/:wallet | GetUsage | ✅ 实现 |
| POST | /api/oracle/monthly-bill | TriggerMonthlyBill | ✅ 实现 |

### 架构分层
Handler → Service → Repository → GORM/PostgreSQL

### Models (6)
User, Operator, Deposit, Bill, UserService, UsageData

### 种子数据
migrations/001_init.sql 预置 11 个运营商（T-Mobile US, Vodafone UK, Orange FR, MTS RU, SoftBank JP, Viettel VN, Unitel LA, Smart KH, AIS TH, Maxis MY, Globe PH）

## 四、合约基线 (packages/contracts)

### 合约 (6)
| 合约 | 继承 | 功能 |
|------|------|------|
| UserRegistry | ERC721, Ownable | 用户注册 + NFT 铸造 |
| FeeManager | Ownable | 平台费率管理（默认 2.5%） |
| Deposit | Ownable, ReentrancyGuard | 用户押金管理 |
| ServiceManager | Ownable | 运营商注册 + 用户服务激活 |
| Payment | Ownable, ReentrancyGuard | 账单创建 + 支付 |
| Oracle | Ownable | 用量上报 + 自动出账 |

### 合约调用链
```
User → UserRegistry.register()           → NFT 铸造
User → Deposit.deposit()                 → 充值（需已注册）
User → ServiceManager.activateService()  → 选运营商开通
Oracle → Oracle.submitUsage()            → 上报用量 → Payment.createBill()
User → Payment.payBill()                 → 支付账单（fee → platformWallet）
User → Deposit.withdraw()               → 提现（需无活跃服务 + 无未付账单）
```

### 部署脚本已部署
UserRegistry → FeeManager → Deposit → ServiceManager → Payment → link(Deposit→Payment)

### 缺失
- Oracle 未在 deploy.ts 中部署
- Deposit 未 link ServiceManager
- Payment 未设 oracle 地址
- 测试仅覆盖 UserRegistry（2 cases）

## 五、三端 API 映射对比

| 业务 | 前端 Mock | 后端 API | 合约 |
|------|-----------|----------|------|
| 注册 | userService.register | POST /api/register | UserRegistry.register |
| 查用户 | userService.getUserProfile | GET /api/user/:wallet | UserRegistry.getUserInfo |
| 运营商列表 | operatorService.getRegions | GET /api/operators | ServiceManager.getActiveOperators |
| 充值 | depositService.deposit | POST /api/deposit | Deposit.deposit |
| 查余额 | depositService.getDeposit | GET /api/deposit/:wallet | Deposit.getDepositAmount |
| 提现 | depositService.withdraw | POST /api/withdraw (stub) | Deposit.withdraw |
| 激活服务 | operatorService.applyNumber | POST /api/service/activate | ServiceManager.activateService |
| 查服务 | — | GET /api/service/:wallet | ServiceManager.getUserService |
| 账单列表 | billingService.getBills | GET /api/bills/:wallet | Payment.getUserBills |
| 支付账单 | billingService.payBill | POST /api/bills/pay (stub) | Payment.payBill |
| 用量查询 | billingService.getCurrentMonthEstimate | GET /api/usage/:wallet | Oracle.getLatestUsage |
| 通知 | notificationService.* | — | — |

## 六、串通差距分析

### 前端 → 后端
1. **接口签名不匹配**：前端 mock 的 request/response 格式与后端不完全一致
2. **数据模型差异**：前端 Region/Operator 结构 vs 后端 Operator 扁平结构
3. **缺失 API**：后端无通知相关 API
4. **认证**：后端无 wallet 签名认证中间件

### 前端 → 合约
1. **前端无合约交互**：当前全部走 mock，未集成 wagmi contract hooks
2. **充值/提现**：应直接调合约而非后端 API
3. **注册**：应调合约 + 后端双写

### 后端 → 合约
1. **后端未读合约状态**：如余额、账单应以合约为准
2. **Oracle 流程**：后端 OracleServiceV2 未连接合约 Oracle
3. **支付**：后端 PayBill 是 stub，实际应由合约处理

### 关键串通任务
1. 前端新增 API client 层（替换 mock）
2. 前端合约交互层（wagmi useContractRead/Write）
3. 后端补齐认证中间件
4. 合约部署脚本补全（Oracle + linking）
5. 数据模型对齐（前端 types ↔ 后端 models ↔ 合约 structs）
6. 混合调用策略：读合约状态 + 写合约交易 + 后端做索引/通知
