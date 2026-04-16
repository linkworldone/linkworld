# 设计文档 — LinkWorld 全栈联调 (Round 2)

> 2026-04-14 | 架构师角色

## 一、架构总览

### 调用路径分类

| 类型 | 路径 | 适用场景 |
|------|------|----------|
| 合约直调 | 前端 → wagmi → 合约 | 资金操作（充值/提现/支付） |
| 后端 API | 前端 → axios → Go 后端 → PostgreSQL | 查询/索引（运营商/账单/用量） |
| 双写 | 前端 → 合约 → 等 tx → 后端 API | 注册、激活服务 |

### 前端新增分层

```
src/
├── config/
│   ├── abis/              # 新增：合约 ABI
│   │   ├── UserRegistry.ts
│   │   ├── Deposit.ts
│   │   ├── ServiceManager.ts
│   │   ├── Payment.ts
│   │   └── FeeManager.ts
│   ├── contracts.ts       # 新增：合约地址（按 chainId 区分）
│   └── api.ts             # 新增：API_BASE_URL
├── services/
│   ├── api/               # 新增：后端 API 调用层
│   │   ├── client.ts      # axios 实例 + 拦截器
│   │   ├── userApi.ts
│   │   ├── operatorApi.ts
│   │   ├── depositApi.ts
│   │   ├── billingApi.ts
│   │   └── usageApi.ts
│   ├── mock/              # 保留：通知仍用 mock
│   └── index.ts           # 改：统一导出 api + mock(notification)
├── hooks/
│   ├── contracts/         # 新增：合约交互 hooks
│   │   ├── useUserRegistry.ts
│   │   ├── useDepositContract.ts
│   │   ├── useServiceManager.ts
│   │   └── usePaymentContract.ts
│   ├── useUser.ts         # 改：内部调用 api + contract
│   ├── useOperator.ts     # 改：内部调用 api
│   ├── useDeposit.ts      # 改：内部调用 contract + api
│   ├── useBilling.ts      # 改：内部调用 api + contract(payBill)
│   └── useNotification.ts # 不变：继续 mock
```

## 二、合约交互设计

### 2.1 ABI 来源

从 `packages/contracts/artifacts/contracts/*.sol/*.json` 提取 abi 字段，存为 TypeScript 常量。

### 2.2 合约地址配置

```typescript
// src/config/contracts.ts
export const CONTRACTS = {
  // Hardhat localhost
  31337: {
    UserRegistry: "0x5FbDB2315678afecb367f032d93F642f64180aa3",
    FeeManager: "0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512",
    Deposit: "0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0",
    ServiceManager: "0xCf7Ed3AccA5a467e9e704C703E8D87F634fB0Fc9",
    Payment: "0xDc64a140Aa3E981100a9becA4E685f962f0cF6C9",
    Oracle: "", // deploy 后填入
  },
  // 0G Testnet
  16601: {
    UserRegistry: "",
    FeeManager: "",
    Deposit: "",
    ServiceManager: "",
    Payment: "",
    Oracle: "",
  },
} as const;
```

### 2.3 合约 Hook 设计

#### useUserRegistry

```typescript
// 读：检查是否已注册
useReadContract({
  address: CONTRACTS[chainId].UserRegistry,
  abi: UserRegistryABI,
  functionName: "isRegistered",
  args: [address],
})

// 写：注册（铸 NFT）
useWriteContract() → writeContract({
  address: CONTRACTS[chainId].UserRegistry,
  abi: UserRegistryABI,
  functionName: "register",
  args: [email],
})
```

#### useDepositContract

```typescript
// 读：查余额
useReadContract({
  functionName: "getDepositAmount",
  args: [address],
}) → returns uint256 (wei)

// 写：充值
writeContract({
  functionName: "deposit",
  value: parseEther(amount), // ETH 原生币
})

// 写：提现
writeContract({
  functionName: "withdraw",
})
```

#### useServiceManager

```typescript
// 读：查用户服务
useReadContract({
  functionName: "getUserService",
  args: [address],
}) → returns { user, operatorId, virtualNumber, password, activatedAt, isActive }

// 读：查运营商（备选，优先走后端）
useReadContract({
  functionName: "getActiveOperators",
})

// 写：激活服务
writeContract({
  functionName: "activateService",
  args: [operatorId, virtualNumber, password],
})
```

#### usePaymentContract

```typescript
// 读：查账单（备选，优先走后端）
useReadContract({
  functionName: "getUserBills",
  args: [address],
})

// 写：支付账单
writeContract({
  functionName: "payBill",
  args: [billId],
  value: totalAmount, // operatorFee + platformFee
})
```

## 三、API Client 设计

### 3.1 axios 实例

```typescript
// src/services/api/client.ts
import axios from "axios";

export const apiClient = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || "http://localhost:8080",
  timeout: 10000,
  headers: { "Content-Type": "application/json" },
});

// 响应拦截器：统一错误处理
apiClient.interceptors.response.use(
  (res) => res.data,
  (err) => Promise.reject(err.response?.data?.error || err.message)
);
```

### 3.2 API Service 接口

#### userApi

```typescript
register(wallet: string, email: string, tokenId?: number)
  → POST /api/register { wallet, email }
  → 返回 { message }

getUser(wallet: string)
  → GET /api/user/{wallet}
  → 返回后端 User JSON → 转换为前端 User interface
```

**转换逻辑（后端 → 前端）：**

| 后端字段 | 前端字段 | 转换 |
|----------|----------|------|
| wallet_addr | address | 直接映射 |
| email | email | 直接映射 |
| is_active | status | true → "active", false 需看 has_unpaid_bills → "suspended" 或 "inactive" |
| registered_at | registeredAt | ISO string |
| token_id | nftTokenId | number |

> 注意：后端 User 没有 status enum（只有 is_active bool）。需要在前端转换层判断：
> - 刚注册没充值 → "inactive"
> - is_active + 有充值 → "active"
> - 有逾期账单 → "suspended"
> 本轮简化：is_active=true → "active"，is_active=false → "inactive"

#### operatorApi

```typescript
getOperators()
  → GET /api/operators
  → 返回 Operator[] → 需聚合为 Region[]

getCountries()
  → GET /api/countries
  → 返回 [{ code, name, prefix }]
```

**转换逻辑：**

后端 Operator 扁平结构 → 前端需要 Region 层级：
1. 从后端拿全量 Operator[]
2. 按 region 字段 groupBy → 生成 Region[]
3. 每个 Region: { code（从 countries API 取）, name, flag（前端 hardcode emoji map）, operatorCount, startingPrice（min(requiredDeposit)）}

#### depositApi

```typescript
recordDeposit(wallet: string, amount: string, txHash?: string)
  → POST /api/deposit { wallet, amount }

getDepositAmount(wallet: string)
  → GET /api/deposit/{wallet}
  → 返回 { amount: string }
```

> 注意：真实余额从合约读（Deposit.getDepositAmount），后端只做记录。

#### billingApi

```typescript
getBills(wallet: string)
  → GET /api/bills/{wallet}
  → 返回 Bill[] → 转换为前端 Bill interface

recordPayment(wallet: string, billId: string, txHash?: string)
  → POST /api/bills/pay { wallet, bill_id, tx_hash }
```

**转换逻辑（后端 Bill → 前端 Bill）：**

| 后端字段 | 前端字段 | 转换 |
|----------|----------|------|
| id | id | string(id) |
| amount | operatorFee | 直接映射 |
| platform_fee | platformFee | 直接映射 |
| is_paid | status | true → "paid", false + 过期 → "overdue", false → "unpaid" |
| created_at | month | 提取 YYYY-MM |
| created_at + 14 days | dueDate | 计算 |
| paid_at | paidAt | ISO string or undefined |
| — | totalAmount | amount + platform_fee |
| — | usage | 需额外调 usage API 或固定值 |

#### usageApi

```typescript
getUsage(wallet: string)
  → GET /api/usage/{wallet}
  → 返回 { data_usage: number, call_usage: number, signature: string }
  → 转换为 MonthEstimate
```

## 四、业务流程详细序列

### 4.1 注册（双写）

```
用户点击 Register
  → RegisterSheet Step 1: 邮箱验证（保持 mock）
  → RegisterSheet Step 2:
    1. useWriteContract → UserRegistry.register(email)
    2. useWaitForTransactionReceipt → 等 tx 确认
    3. 从 tx receipt 解析 tokenId（UserRegistered event）
    4. apiClient.post("/api/register", { wallet, email })
    5. invalidateQueries(["user"])
    6. navigate → /dashboard
```

### 4.2 充值（合约 + 后端记录）

```
用户输入金额，点击 Deposit
  1. useWriteContract → Deposit.deposit() { value: parseEther(amount) }
  2. useWaitForTransactionReceipt → 等 tx 确认
  3. apiClient.post("/api/deposit", { wallet, amount })  // 后端记录
  4. invalidateQueries(["deposit", "depositHistory"])
  // 余额始终从合约读：useReadContract → Deposit.getDepositAmount(address)
```

### 4.3 提现（合约 + 后端记录）

```
用户点击 Withdraw
  1. useWriteContract → Deposit.withdraw()
     - 合约自动检查：无活跃服务 + 无未付账单
  2. useWaitForTransactionReceipt
  3. apiClient.post("/api/withdraw", { wallet })  // 后端记录
  4. invalidateQueries(["deposit", "depositHistory"])
```

### 4.4 激活服务（三步）

```
用户选运营商，点击 Apply
  1. apiClient.post("/api/virtual-number/generate", { country_code })
     → 返回 { virtual_number, password }
  2. useWriteContract → ServiceManager.activateService(operatorId, virtualNumber, password)
  3. useWaitForTransactionReceipt
  4. apiClient.post("/api/service/activate", { wallet, operator_id, virtual_number, password })
  5. invalidateQueries(["userService", "myNumbers"])
```

### 4.5 支付账单（合约 + 后端更新）

```
用户点击 Pay Bill
  1. useWriteContract → Payment.payBill(billId) { value: totalAmount }
  2. useWaitForTransactionReceipt
  3. apiClient.post("/api/bills/pay", { wallet, bill_id })  // 后端标记已付
  4. invalidateQueries(["bills", "billDetail"])
```

## 五、后端改造设计

### 5.1 CORS 中间件

```go
// cmd/main.go 添加
r.Use(cors.New(cors.Config{
  AllowOrigins: []string{"http://localhost:5173"},
  AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
  AllowHeaders: []string{"Content-Type", "Authorization"},
}))
```

依赖：`github.com/gin-contrib/cors`

### 5.2 Register handler 改造

接收 tokenId 参数（合约铸造后传入）：

```go
type RegisterRequest struct {
  Wallet  string `json:"wallet" binding:"required"`
  Email   string `json:"email" binding:"required,email"`
  TokenID uint   `json:"token_id"`
}
```

### 5.3 PayBill handler 改造（从 stub → 实际）

```go
func (h *Handler) PayBill(c *gin.Context) {
  var req struct {
    Wallet string `json:"wallet" binding:"required"`
    BillID uint   `json:"bill_id" binding:"required"`
    TxHash string `json:"tx_hash"`
  }
  // 查 bill → 标记 is_paid=true, paid_at=now
}
```

### 5.4 Withdraw handler 改造（从 stub → 记录）

```go
func (h *Handler) Withdraw(c *gin.Context) {
  var req struct {
    Wallet string `json:"wallet" binding:"required"`
    TxHash string `json:"tx_hash"`
  }
  // 记录提现事件（金额从合约事件中获取）
}
```

### 5.5 种子数据确认

后端 AutoMigrate 建表后 operators 表为空。需要：
- 方案 A：启动时自动 seed（检查空表则插入）
- 方案 B：运行 migrations/001_init.sql 手动 seed

**选择方案 A**：在 main.go 中 AutoMigrate 后加 seed 逻辑。

## 六、合约部署补全

### deploy.ts 需补充

```typescript
// 6. Oracle
const Oracle = await ethers.getContractFactory("Oracle");
const oracle = await Oracle.deploy();
await oracle.waitForDeployment();

// 7. 关联
await deposit.setServiceManager(await serviceManager.getAddress());
await oracle.setPayment(await payment.getAddress());
await payment.setOracle(await oracle.getAddress());

// 8. 输出所有地址到 JSON（供前端 config 读取）
```

## 七、wagmi 配置改造

当前 wagmi config 连的是 0G Testnet (16601)。本地开发需加 Hardhat 网络：

```typescript
// src/config/chains.ts 新增
export const hardhatLocal = {
  id: 31337,
  name: "Hardhat Local",
  nativeCurrency: { name: "Ether", symbol: "ETH", decimals: 18 },
  rpcUrls: { default: { http: ["http://127.0.0.1:8545"] } },
};

// src/config/wagmi.ts
// chains 数组根据环境切换：dev 用 hardhatLocal，prod 用 zgTestnet
```

## 八、数据模型对齐汇总

### 前端 types/index.ts 需调整

| Interface | 变更 |
|-----------|------|
| User | 新增 `id?: number`（后端返回），status 从后端 is_active 转换 |
| DepositInfo | balance 改为从合约读（uint256 → bigint），currency 固定 "ETH"（本地链用 ETH） |
| DepositRecord | 保留结构，新增 txHash 从链上获取 |
| Region | 保持不变，前端从 Operator[] 聚合生成 |
| Operator | id 改为 number（后端返回 uint），新增 countryCode |
| Bill | id 改为 number，新增映射逻辑 |

### 新增类型

```typescript
// 后端原始 response 类型（snake_case）
interface ApiUser {
  id: number;
  wallet_addr: string;
  email: string;
  token_id: number;
  is_active: boolean;
  registered_at: string;
}

interface ApiOperator {
  id: number;
  name: string;
  region: string;
  required_deposit: string;
  is_active: boolean;
}

interface ApiBill {
  id: number;
  user_id: number;
  operator_id: number;
  amount: string;
  platform_fee: string;
  is_paid: boolean;
  created_at: string;
  paid_at?: string;
}
```

转换函数放在各 api service 文件中。

## 九、环境配置

### .env.local（前端）

```
VITE_API_BASE_URL=http://localhost:8080
VITE_CHAIN_ID=31337
```

### 启动顺序

1. `npx hardhat node` → 本地链 (8545)
2. `npx hardhat run scripts/deploy.ts --network localhost` → 部署合约
3. `go run cmd/main.go` → 后端 (8080)
4. `pnpm dev` → 前端 (5173)

## 八-A、双写失败补偿设计

### 策略：链上状态兜底 + 前端重试

**原则**：合约是 source of truth，后端是索引。合约成功但后端失败不应阻塞用户操作。

### 每个双写流程的失败处理

#### 注册（合约成功 + 后端失败）
- 前端 catch 后端 API 错误 → 将 { wallet, email, txHash } 存入 localStorage
- useUser hook 启动时检测：`isRegistered(address)` 返回 true 但 `GET /api/user` 404 → 自动重试 `POST /api/register`
- 用户可见：toast 提示"注册成功，正在同步数据..."

#### 充值（合约成功 + 后端失败）
- 余额始终从合约读，后端仅记录历史
- 后端 POST 失败 → toast 提示但不阻塞，余额已正确更新
- 充值历史暂时缺失该笔记录，不影响核心功能

#### 激活服务（合约成功 + 后端失败）
- 前端 catch → 将 { wallet, operatorId, virtualNumber, password, txHash } 存 localStorage
- useOperator hook 检测：合约 getUserService 有记录但后端 404 → 自动重试 activate
- 用户可见：toast "服务已激活，正在同步..."

#### 支付账单（合约成功 + 后端失败）
- 前端 catch → 存 { wallet, billId, txHash } 到 localStorage
- useBilling hook 检测：合约 bill.isPaid=true 但后端 is_paid=false → 自动重试 payBill API
- 用户可见：toast "支付成功，正在同步..."

### 实现要点
- 统一工具函数：`savePendingSync(key, data)` / `getPendingSync(key)` / `clearPendingSync(key)`
- 每个业务 hook 在初始化时检查 localStorage 中的 pending sync
- 重试最多 3 次，间隔 2s，全部失败后 toast 提示用户手动刷新
- 后续 Round 3 用事件索引器替代前端补偿

## 八-B、链上交易 UX 规范

### Transaction Pending 状态

所有合约写操作统一使用 `useTransactionFlow` hook：

```typescript
// 三阶段状态机
type TxState = 
  | { status: "idle" }
  | { status: "pending-signature" }    // 钱包弹窗等签名
  | { status: "pending-confirmation" } // 等区块确认
  | { status: "success"; txHash: string }
  | { status: "error"; message: string }
```

### UI 行为
- **pending-signature**: BottomSheet 显示"请在钱包中确认交易..."，按钮 disabled
- **pending-confirmation**: BottomSheet 显示 spinner + "交易确认中..."，预估 5-15 秒
- **success**: 绿色 ✓ + "交易成功"，1.5 秒后自动关闭并刷新数据
- **error**: 红色提示 + 错误原因（中文映射），显示"重试"按钮

### 合约 Revert Reason 映射表

| Revert Message | 中文提示 |
|---------------|----------|
| Not registered | 请先注册账户 |
| Already registered | 该钱包已注册 |
| Zero deposit | 充值金额不能为零 |
| No deposit | 没有可提取的余额 |
| Service still active | 请先停用服务后再提现 |
| Has unpaid bills | 请先支付所有账单后再提现 |
| Not your bill | 该账单不属于您 |
| Already paid | 该账单已支付 |
| Insufficient payment | 支付金额不足 |
| Only oracle | 仅限预言机操作 |

### 用户切换页面处理
- 交易提交后（有 txHash），即使用户切换页面，useWaitForTransactionReceipt 仍在后台等待
- 交易完成后通过 invalidateQueries 更新全局数据
- 不需要额外处理，React Query + wagmi 已支持此场景

## 十、本轮不涉及

| 项 | 理由 |
|----|------|
| 后端认证中间件 | 决策暂不加 |
| 通知 API | 后端未实现，前端继续 mock |
| 合约事件监听 | 复杂度高，后续补 |
| 合约测试 | 聚焦串通 |
| 错误重试/回退 | 先跑通 happy path |
