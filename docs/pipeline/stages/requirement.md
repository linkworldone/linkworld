# 需求文档 — LinkWorld 全栈联调 (Round 2)

> 2026-04-14 | 产品经理角色

## 一、背景

Round 1 完成了前端 8 页面 + Mock Service 的完整 UI 实现。Round 2 目标是将前端、后端（Go/Gin）、合约（Solidity/Hardhat）三端串通，使用户操作真正触达链上和数据库。

## 二、核心决策

| 决策项 | 选择 | 理由 |
|--------|------|------|
| 资金路径 | **混合模式** | 充值/提现走合约（用户签名），账单/用量走后端（Oracle 代理） |
| 注册流程 | **合约 + 后端双写** | 先调合约铸 NFT → 成功后调后端存 email/profile |
| API Client | **axios + React Query** | axios 拦截器方便后续加认证 |
| 后端认证 | **暂不加** | 先串通功能，认证后续补 |

## 三、业务流程串通定义

### 3.1 注册流程
```
前端 RegisterSheet
  ├── Step 1: 邮箱验证（保持现有 mock，或后端发真邮件 — 本轮用 mock）
  └── Step 2: 注册
       ├── ① wagmi writeContract → UserRegistry.register(email) → 链上铸 NFT
       ├── ② 等待 tx confirmation
       └── ③ axios POST /api/register { wallet, email, tokenId } → 后端存档
```

### 3.2 充值流程（合约优先）
```
前端 Deposit 页面
  ├── 用户输入金额
  ├── wagmi writeContract → Deposit.deposit() { value: amount }
  ├── 等待 tx confirmation
  ├── axios POST /api/deposit { wallet, amount, txHash } → 后端记录
  └── 刷新余额: wagmi readContract → Deposit.getDepositAmount(address)
```

### 3.3 提现流程（合约优先）
```
前端 Deposit 页面
  ├── wagmi writeContract → Deposit.withdraw()
  ├── 等待 tx confirmation
  ├── axios POST /api/withdraw { wallet, txHash } → 后端记录
  └── 刷新余额
```

### 3.4 运营商列表（后端优先）
```
前端 Services 页面
  ├── axios GET /api/operators → 后端返回运营商列表
  └── 或 wagmi readContract → ServiceManager.getActiveOperators()（备选）

选择后端优先原因：后端可做分页、搜索、缓存，合约读全量数组 gas 高
```

### 3.5 激活服务（混合）
```
前端 RegionDetail 页面
  ├── axios POST /api/virtual-number/generate { country_code } → 后端生成虚拟号码
  ├── wagmi writeContract → ServiceManager.activateService(operatorId, virtualNumber, password)
  ├── 等待 tx confirmation
  └── axios POST /api/service/activate { wallet, operator_id, virtual_number, password } → 后端存档
```

### 3.6 账单列表（后端优先）
```
前端 Billing 页面
  ├── axios GET /api/bills/:wallet → 后端返回账单列表
  └── 合约 Payment.getUserBills() 作为可选校验源
```

### 3.7 支付账单（合约优先）
```
前端 BillDetail 页面
  ├── wagmi writeContract → Payment.payBill(billId) { value: totalAmount }
  ├── 等待 tx confirmation
  ├── axios POST /api/bills/pay { wallet, billId, txHash } → 后端更新状态
  └── 刷新账单列表
```

### 3.8 用量查询（后端）
```
前端 Dashboard
  └── axios GET /api/usage/:wallet → 后端返回用量数据
```

### 3.9 通知（后端）
```
前端 Notifications 页面
  └── 需后端新增通知 API（当前缺失）
  └── 本轮保持 mock，后续补
```

## 四、前端改造范围

### 4.1 新增层
| 层 | 文件 | 职责 |
|----|------|------|
| API Client | src/services/api/client.ts | axios 实例，baseURL，拦截器 |
| API Services | src/services/api/*.ts | 每个业务域一个文件，调后端 API |
| Contract Hooks | src/hooks/contracts/*.ts | wagmi useReadContract / useWriteContract 封装 |
| ABI | src/config/abis/*.ts | 合约 ABI（从 artifacts 复制） |
| Contract Addresses | src/config/contracts.ts | 各合约地址（按网络区分） |

### 4.2 改造层
| 层 | 变更 |
|----|------|
| src/services/index.ts | 从导出 mock → 导出 api services + contract hooks |
| src/hooks/*.ts | React Query hooks 内部调用从 mock → api/contract |
| src/types/index.ts | 对齐后端 response 结构 |
| src/config/constants.ts | 新增 API_BASE_URL |

### 4.3 不变
- 页面组件（8 个 pages）— 不改，只换底层数据源
- UI 组件 — 不改
- 布局组件 — 不改
- Provider 树 — 不改（已有 wagmi + RainbowKit）

## 五、后端改造范围

### 5.1 补全
| 项 | 说明 |
|----|------|
| CORS 中间件 | 允许 localhost:5173 跨域 |
| PayBill handler | 从 stub → 实际更新 DB |
| Withdraw handler | 从 stub → 实际记录 |
| Operator 种子数据 | 确保 DB 有 11 个运营商 |

### 5.2 新增
| 项 | 说明 |
|----|------|
| GET /api/notifications/:wallet | 通知列表（可选，本轮可保持 mock） |

### 5.3 不变
- 其余 13 个已实现的 API

## 六、合约改造范围

### 6.1 部署补全
| 项 | 说明 |
|----|------|
| 部署 Oracle 合约 | 目前 deploy.ts 未包含 |
| link Deposit → ServiceManager | deposit.setServiceManager() |
| link Oracle → Payment | oracle.setPayment() |
| set Payment oracle | payment.setOracle(oracleAddress) |

### 6.2 不变
- 6 个合约代码本身不改

## 七、数据模型对齐

### 前端 types 需适配后端 response

| 前端 Interface | 后端 JSON | 差异 |
|---------------|-----------|------|
| User | { wallet_addr, email, token_id, is_active } | 字段名 snake_case vs camelCase |
| Operator | { id, name, region, required_deposit, is_active } | 前端多 Region 层级 |
| DepositInfo | { amount: string } | 前端用 bigint，后端返回 string |
| Bill | { amount, platform_fee, is_paid } | 前端用 string 金额 + status enum |

### 策略
- 前端 API service 层做 response → 前端 type 的转换
- 不改后端 response 格式

## 八、本轮不做

| 项 | 理由 |
|----|------|
| 后端认证 | 先串通，后续补 wallet 签名认证 |
| 真邮件验证 | 保持 mock 验证码（always pass） |
| 通知 API | 后端未实现，前端继续用 mock |
| 合约测试补全 | 聚焦串通，测试后续补 |
| 生产环境部署 | 本轮限本地环境 |

## 九、验收标准

1. 前端注册 → 链上铸 NFT + 后端存用户 ✓
2. 前端充值 → 链上 Deposit 合约余额增加 ✓
3. 前端提现 → 链上 Deposit 合约余额减少 ✓
4. 前端查运营商 → 后端返回真实数据 ✓
5. 前端激活服务 → 链上 ServiceManager 记录 + 后端存档 ✓
6. 前端查账单 → 后端返回真实账单 ✓
7. 前端支付账单 → 链上 Payment 完成支付 ✓
8. 前端查余额 → 读链上合约真实值 ✓
9. 三端本地环境可一键启动调试 ✓
