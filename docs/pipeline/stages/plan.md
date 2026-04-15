# 实现计划 — LinkWorld 全栈联调 (Round 2)

> 2026-04-14 | 11 个 Task，预估串行 3-4 小时

## 任务依赖图

```
Task 0 (合约部署补全)
  ↓
Task 1 (前端合约 ABI + 地址配置) ← Task 0
  ↓
Task 2 (前端 axios client) — 独立
  ↓
Task 3 (后端改造) — 独立
  ↓
Task 4 (前端 API service 层) ← Task 2 + Task 3
  ↓
Task 5 (前端合约 hooks) ← Task 1
  ↓
Task 6 (前端 useTransactionFlow) ← Task 5
  ↓
Task 7 (前端双写补偿工具) — 独立
  ↓
Task 8 (业务 hooks 改造) ← Task 4 + Task 5 + Task 6 + Task 7
  ↓
Task 9 (前端 types 对齐 + services/index.ts) ← Task 4
  ↓
Task 10 (wagmi 配置 + 环境变量 + 冒烟测试)
```

## 并行策略

| 批次 | Tasks | 理由 |
|------|-------|------|
| 批次 1 | Task 0 + Task 2 + Task 3 | 无文件交集，完全独立 |
| 批次 2 | Task 1（依赖 Task 0） | |
| 批次 3 | Task 4 + Task 5 + Task 7 | 无文件交集 |
| 批次 4 | Task 6（依赖 Task 5） | |
| 批次 5 | Task 8 + Task 9 | Task 8 改 hooks，Task 9 改 types + index |
| 批次 6 | Task 10 | 最终验证 |

---

## Task 0: 合约部署脚本补全

**改动文件**: `packages/contracts/scripts/deploy.ts`

**内容**:
1. 部署 Oracle 合约
2. 关联 Deposit → ServiceManager: `deposit.setServiceManager(serviceManager)`
3. 关联 Oracle → Payment: `oracle.setPayment(payment)`
4. 关联 Payment → Oracle: `payment.setOracle(oracle)`
5. 脚本末尾输出所有合约地址到 `packages/contracts/deployments/localhost.json`

**验收**: `npx hardhat run scripts/deploy.ts --network localhost` 成功，输出 JSON 含 6 个合约地址

---

## Task 1: 前端合约 ABI + 地址配置

**改动文件**:
- 新增 `packages/web/src/config/abis/UserRegistry.ts`
- 新增 `packages/web/src/config/abis/Deposit.ts`
- 新增 `packages/web/src/config/abis/ServiceManager.ts`
- 新增 `packages/web/src/config/abis/Payment.ts`
- 新增 `packages/web/src/config/abis/FeeManager.ts`
- 新增 `packages/web/src/config/contracts.ts`

**内容**:
1. 从 `packages/contracts/artifacts/` 提取每个合约的 abi 字段，导出为 `const XxxABI = [...] as const`
2. `contracts.ts` 按 chainId 映射合约地址，从 Task 0 输出的 JSON 读取 localhost 地址
3. 导出 `getContractAddress(chainId, contractName)` 工具函数，空地址时 throw

**验收**: TypeScript 编译无错误

---

## Task 2: 前端 axios client

**改动文件**:
- 新增 `packages/web/src/services/api/client.ts`
- 新增 `packages/web/src/config/api.ts`（API_BASE_URL 常量）

**内容**:
1. axios 实例：baseURL 从 `import.meta.env.VITE_API_BASE_URL` 读取，默认 `http://localhost:8080`
2. 响应拦截器：`res => res.data`，错误拦截器提取 error message
3. `api.ts` 导出 API_BASE_URL

**验收**: 实例可正常创建，TypeScript 无错误

---

## Task 3: 后端改造

**改动文件**:
- `packages/backend/cmd/main.go` — 加 CORS + seed
- `packages/backend/internal/models/models.go` — Operator 加 CountryCode, Bill 加 TxHash
- `packages/backend/internal/handlers/handlers.go` — PayBill 补全, Withdraw 补全, Register 加 token_id
- `packages/backend/internal/repository/repository.go` — BillRepo 补 MarkAsPaid(id, txHash)
- `packages/backend/go.mod` — 加 gin-contrib/cors

**内容**:
1. CORS 允许 localhost:5173
2. Operator model 加 `CountryCode string \`json:"country_code"\``
3. Bill model 加 `TxHash string \`json:"tx_hash"\``
4. Register handler 接收 token_id 参数
5. PayBill handler: 接收 { wallet, bill_id, tx_hash } → 查 bill → 标记已付
6. Withdraw handler: 接收 { wallet, tx_hash } → 记录提现
7. main.go AutoMigrate 后 seed 11 个运营商（检查空表再插入），含 country_code
8. 安装 gin-contrib/cors

**验收**: `go build ./...` 成功，启动后 GET /api/operators 返回 11 个含 country_code 的运营商

---

## Task 4: 前端 API service 层

**改动文件**:
- 新增 `packages/web/src/services/api/userApi.ts`
- 新增 `packages/web/src/services/api/operatorApi.ts`
- 新增 `packages/web/src/services/api/depositApi.ts`
- 新增 `packages/web/src/services/api/billingApi.ts`
- 新增 `packages/web/src/services/api/usageApi.ts`

**内容**:
每个文件导出异步函数，调 apiClient，返回前端 interface 类型（内部做 snake_case → camelCase 转换）。

具体函数签名见设计文档"三、API Client 设计"章节。

关键转换：
- `operatorApi.getRegions()`: 调 GET /api/operators + GET /api/countries → groupBy region → 返回 Region[]
- `billingApi.getBills()`: 后端 Bill → 前端 Bill（计算 totalAmount, dueDate, status）
- `userApi.getUser()`: 后端 User → 前端 User（is_active → status 枚举）

**验收**: 所有 API service 函数 TypeScript 类型正确

---

## Task 5: 前端合约 hooks

**改动文件**:
- 新增 `packages/web/src/hooks/contracts/useUserRegistry.ts`
- 新增 `packages/web/src/hooks/contracts/useDepositContract.ts`
- 新增 `packages/web/src/hooks/contracts/useServiceManager.ts`
- 新增 `packages/web/src/hooks/contracts/usePaymentContract.ts`

**内容**:
每个 hook 封装 wagmi `useReadContract` / `useWriteContract` / `useWaitForTransactionReceipt`。

- useUserRegistry: isRegistered(address), register(email)
- useDepositContract: getDepositAmount(address), deposit(amount), withdraw()
- useServiceManager: getUserService(address), activateService(operatorId, vn, pwd)
- usePaymentContract: payBill(billId, value)

**验收**: TypeScript 编译无错误，hook 返回类型正确

---

## Task 6: useTransactionFlow hook

**改动文件**:
- 新增 `packages/web/src/hooks/useTransactionFlow.ts`

**内容**:
统一的链上交易状态管理 hook，封装五阶段状态机（idle → pending-signature → pending-confirmation → success/error）。

包含 revert reason 中文映射表（10 条）。

供所有合约写操作复用。

**验收**: TypeScript 编译无错误

---

## Task 7: 双写补偿工具

**改动文件**:
- 新增 `packages/web/src/utils/pendingSync.ts`

**内容**:
- `savePendingSync(key, data)` — 存 localStorage
- `getPendingSync(key)` — 读 localStorage
- `clearPendingSync(key)` — 清 localStorage
- `retryWithBackoff(fn, maxRetries=3, baseDelay=2000)` — 重试工具

**验收**: 函数签名正确，TypeScript 无错误

---

## Task 8: 业务 hooks 改造

**改动文件**:
- `packages/web/src/hooks/useUser.ts`
- `packages/web/src/hooks/useOperator.ts`
- `packages/web/src/hooks/useDeposit.ts`
- `packages/web/src/hooks/useBilling.ts`
- （useNotification.ts 不改，继续 mock）

**内容**:
将每个 hook 内部的 mock service 调用替换为：
- 查询类 → api service（axios）
- 合约交互类 → contract hooks（wagmi）
- 写操作 → useTransactionFlow 包装
- 双写流程 → 合约写 + API 写 + pendingSync 补偿

具体：
- **useUser**: useRegister 改为先调合约 register → 再调 userApi.register → pendingSync 兜底
- **useOperator**: useRegions 改调 operatorApi.getRegions，useApplyNumber 改三步流程
- **useDeposit**: useDeposit 余额从合约读，history 从后端读，充值/提现走合约 + 后端双写
- **useBilling**: useBills 从后端读，usePayBill 走合约 + 后端双写

**验收**: 所有 hooks TypeScript 编译无错误，函数签名与页面调用兼容

---

## Task 9: 前端 types 对齐 + services/index.ts

**改动文件**:
- `packages/web/src/types/index.ts`
- `packages/web/src/services/index.ts`

**内容**:
1. types/index.ts:
   - Operator.id 改为 number
   - 新增 Operator.countryCode: string
   - DepositInfo.currency 改为 string（支持 ETH）
   - 新增 ApiUser, ApiOperator, ApiBill 中间类型
2. services/index.ts:
   - 从导出 mock → 导出 api services
   - notificationService 继续导出 mock

**验收**: TypeScript 编译无错误

---

## Task 10: wagmi 配置 + 环境变量 + 冒烟测试

**改动文件**:
- `packages/web/src/config/chains.ts` — 新增 hardhatLocal 链定义
- `packages/web/src/config/wagmi.ts` — dev 环境用 hardhatLocal
- 新增 `packages/web/.env.local` — VITE_API_BASE_URL + VITE_CHAIN_ID
- `packages/web/vite.config.ts` — 确认 env 变量可用

**内容**:
1. chains.ts 新增 hardhatLocal（chainId: 31337, rpc: http://127.0.0.1:8545）
2. wagmi.ts: 开发环境 chains=[hardhatLocal]，生产环境 chains=[zgTestnet]
3. .env.local: VITE_API_BASE_URL=http://localhost:8080, VITE_CHAIN_ID=31337

**冒烟测试**:
- 启动 Hardhat 节点 + 部署合约
- 启动后端
- 启动前端 dev server
- 打开浏览器连接钱包（Hardhat 测试账户）
- 走通注册流程（合约铸 NFT + 后端存用户）
- 走通充值流程（合约 deposit + 后端记录）
- 检查余额（从合约读取）

**验收**: 前端 dev server 启动无错误，钱包可连接，至少注册 + 充值两个流程端到端通过
