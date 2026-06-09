# Stage: design — 后端技术设计（子项目 backend 2/3）

> **状态**: v1 | **日期**: 2026-06-09 | **Gate**: 1（设计） | **子项目**: backend(2/3) | **角色**: 架构师 | **分支**: backend/align-arbitrum-usdt
> **输入**：合约 handoff-backend.md（冻结 ABI/selector + USDT 6 位 + 计价归后端 + 批量上限 N）+ scan.md（后端 8 大发现）+ 4 份基线（project-scan/blockchain-integration/services-api/alignment-surface）+ requirement.md §三②/§五B + 真实源码逐文件核对。
> **用户已拍板锁定 3 决策**（见 §1.3），本设计据此展开，不再二次澄清。
> **注**：合约子项目(1/3) 的 design 见 git commit `2b8d4a0`（本 stage 文件按子项目串行复用，前版已入库）。
>
> ⚠️ 本文档只产设计，不写 .go 代码。implement 阶段按 §8 任务划分执行。

---

## 1. 背景与目标

### 1.1 为什么（问题陈述）

合约子项目(1/3)已完成 ERC20(USDT) 迁移 + selector 变更 + 计价归后端 + onlyOracle 权限收口，并产出冻结 ABI（handoff-backend.md）。但后端(2/3)**链集成几乎全为 stub**（scan 8 大发现，逐文件核对确认）：

- `blockchain/client.go`：业务方法全部返回零值/`not implemented`，无任何合约 ABI 绑定、无链上写交易。
- `sync/event_sync.go`：主循环只 `time.Sleep(30s)` 空转，`FilterLogs`/`SubscribeFilterLogs` 从未调用，`process*` 永不触发。
- `services/oracle.go`：「计价」是 `OperatorAPISimulator.GetBill` 返回 `rand.Intn(5000)+500` 随机数；`FetchAndCreateBills` 只写 DB Bill，**完全不碰 `monthlySettlement`/`createBill`**；`SignData` 是 SHA256 拼字符串（非 ECDSA）。
- `config/config.go`：struct 读 `proxies` 键，但 `deployments.json` 写的是 `contracts` 键 → `Proxies` 永远反序列化为空 map（**真实 bug**，event sync 静默失效）。
- `configs/deployments.json`：`chainId=16602`(0G) + `rpcUrl=evm-testnet.0g.ai` + 7 个 0G 旧地址；缺 `usdt`/`usdtDecimals`/`abiHash`。
- `signatures.go`：`BillCreated` 是旧 5 参签名（冻结 ABI 是 4 参，topic hash 错匹配不到）；缺 `TrafficCardApplied`/`UsageDataSubmitted`。

→ 本子项目的「对齐」**不是改几处调用签名，而是从零接入链上写 + 真实事件同步 + 真实计价 + owner 签名 + 端点鉴权**。规模评估：链侧接近重写，业务 DB 层（services.go/repository/models）基本复用。

### 1.2 本轮范围边界

| 做（in） | 不做（out / 留后续） |
|----------|----------------------|
| abigen 生成 8 份合约绑定（7 业务 + MockUSDT） | ❌ 不改主流程业务路由结构 |
| `client.go` 真实写调用（monthlySettlement/issueMonthly…）+ 读调用 | ❌ 不动 web、不动 packages/contracts |
| `event_sync.go` 真实 FilterLogs 多事件落库（6 位精度） | ❌ 不做流量卡资金抵扣对账（合约侧仅桩，handoff §6） |
| 真实计价：费率表 × usage → amounts[] → monthlySettlement 分批 | ❌ 不补 NotificationService 真实 SMTP（保持 stub，本轮非阻塞） |
| owner key 从 env 读 + `bind.NewKeyedTransactorWithChainID` 真实签名发交易 | ❌ 不接 phone 字段 / 本地号码注册（requirement R10） |
| 敏感端点（oracle/monthly-bill、usage/submit、withdraw 回调）加管理员鉴权 | ❌ 不引入 KMS（本轮 env 私钥 + 部署侧最小权限，KMS 留后续 Round） |
| config 键名 bug 修复 + deployments.json 421614 schema + 本地 31337 联调机制 | ❌ 端到端真机联调**待合约 1/3 真·上链**（依赖项，标注不阻塞 implement/test） |
| go test 覆盖（计价/分批/精度/签名构造/事件解码，含 mock RPC） | |

### 1.3 用户已拍板锁定决策（必须落进实现）

1. **计价 = 按用量 × 费率表**：后端维护费率表（每 MB 流量单价 + 每分钟通话单价，可按运营商/地区），`usage × 单价` 算出 USDT 金额（6 位最小单位），传 `monthlySettlement` 的 `amounts[]`。**替换现状 `rand.Intn` 随机数。**
2. **范围 = 建真实链集成 + 单测**：abigen 8 份绑定 → client.go 真实写（后端作 owner 角色调 `Oracle.monthlySettlement`）→ event_sync 真实监听落库；链配置 421614；go test 覆盖。端到端联调待合约真·上链（依赖项标注）。
3. **敏感端点加管理员鉴权**：`oracle/monthly-bill`、`usage/submit` 等加 API key 或签名校验。

---

## 2. 已核对真实代码清单（文件:行 / 结论）

| 文件:行 | 核对结论（权威，本设计据此） |
|---------|------------------------------|
| `internal/blockchain/client.go:22-33` | `NewClient` 真 `ethclient.Dial`（唯一生效部分）；`:40-71` 业务方法全 stub 返零值/`not implemented`，**无 ABI 绑定** |
| `internal/blockchain/signatures.go:20` | `BillCreated` 写死 **5 参** `BillCreated(uint256,address,uint256,uint256,uint256)` → 与冻结 ABI 4 参不符，topic 错 |
| `internal/blockchain/signatures.go:13-26` | 只定义 5 个 topic；缺 `TrafficCardApplied(uint256)`、`UsageDataSubmitted`、`DepositWithdrawn` |
| `internal/config/config.go:13-17` | struct tag `json:"proxies"`，但 JSON 键是 `contracts` → **Proxies 永远空 map**（bug 根因）；无 Usdt/UsdtDecimals/AbiHash 字段 |
| `internal/sync/event_sync.go:33-45` | `syncUserRegistered` 主循环只 `time.Sleep(30s)`，**从未调 FilterLogs**；`processDepositMade:61` 空返回 nil |
| `internal/sync/event_sync.go:47-59` | `processUserRegistered` 落库 `RegisteredAt=time.Unix(0,0)`，不解析 email/tokenId |
| `internal/services/oracle.go:130-134` | `GetBill` 返回 `rand.Intn(5000)+500` + 2.5% fee，**随机数，无精度语义** |
| `internal/services/oracle.go:160-194` | `FetchAndCreateBills` 只写 DB Bill，**不调任何合约**，无分批，无 amounts[] |
| `internal/services/oracle.go:215-219` | `SignData` 是 `sha256.Sum256(拼字符串)`，**非 ECDSA/secp256k1**，链上不可验证 |
| `internal/models/models.go:29-50` | `Deposit.Amount`/`Bill.Amount`/`Bill.PlatformFee`/`Bill.TrafficCardDeduction` 全 `string`；`UsageData.DataUsage/CallUsage` 是 `uint64` |
| `cmd/main.go:40-50` | 11 运营商 seed，`RequiredDeposit="0.01"` 等小数字符串，单位语义不明 |
| `cmd/main.go:78-90` | 仅 `RPC_URL` 非空才起 event sync；用 `deployments.Proxies`（当前空 map）；**无 owner key 装配** |
| `cmd/main.go:103-120` | 18 业务端点，**无任何鉴权中间件**（仅 CORS） |
| `configs/deployments.json:1-13` | `chainId:16602` + `rpcUrl:evm-testnet.0g.ai` + 键名 `contracts` + 7 个 0G 旧地址；缺 usdt/usdtDecimals/abiHash |
| `packages/contracts/deployments/hardhat.json` | **后端应对齐的范式**：键名 `proxies`，含 `usdt`(0x5FbD…aa3) + `usdtDecimals:6` + `abiHash`(7 合约) |
| `packages/contracts/contracts/interfaces/IPayment.sol:15-17` | `BillCreated(uint256,address,uint256,uint256)`（4 参，权威）、`BillPaid(uint256,address,uint256,uint256)`、`TrafficCardApplied(uint256)` |
| `packages/contracts/contracts/interfaces/IDeposit.sol:5-16` | `DepositMade(address,uint256)`、`DepositWithdrawn(address,uint256,uint256)`、`TrafficCardMinted(address,uint256,uint256)`；`deposit(uint256)`/`withdraw()`/`issueMonthlyTrafficCards(address[])` |
| `packages/contracts/contracts/Oracle.sol:54-86` | `monthlySettlement(address[],uint256[],uint256[]) onlyOwner`，per-bill emit **`UsageDataSubmitted(user,operatorId,amount)`**，内部 `payment.createBill` + `deposit.issueMonthlyTrafficCards` + `applyTrafficCardToBill` 桩 |
| `packages/contracts/contracts/Oracle.sol:89` | `submitUsage(address,uint256,uint256,uint256)` 预留接口，v2-A 不参与计价 |

---

## 3. 整体方案

### 3.1 链集成补全顺序（依赖驱动，串行）

```
T1 abigen 8 份绑定 ──┐
T2 config 修复 + schema ──┴─> T3 client.go 真实读/写 ──> T4 owner 签名(transactor)
                                      │
T5 event_sync 真实监听落库 <──────────┘（依赖 T1 绑定 + T2 地址表）
T6 计价费率表 + amounts[] + 分批 ──> 依赖 T3(写) + T4(签名)
T7 敏感端点鉴权中间件（独立，可与 T6 并行）
T8 go test 覆盖（依赖 T1–T7）
```

顺序理由：绑定（T1）是一切链调用的前提；config（T2）提供地址表/精度/RPC；client 写（T3）依赖绑定；owner 签名（T4）是写交易能落链的前提；event_sync（T5）依赖绑定+地址表；计价（T6）依赖写+签名；鉴权（T7）独立；测试（T8）收口。

### 3.2 关键选型

| 选型点 | 决策 | 理由 |
|--------|------|------|
| ABI 绑定方式 | **abigen 生成 Go 绑定**（`go-ethereum/cmd/abigen`），8 份 | vs 手写 `bind.NewBoundContract`：abigen 自动产出类型安全的 `*Caller`/`*Transactor`/`*Filterer`，编码/解码事件无需手写 topic 匹配；与冻结 ABI 一一对应，`abiHash` 可比对。手写易错且无类型保护。ABI 来源 = `packages/contracts/artifacts/contracts/<Name>.sol/<Name>.json` 的 `.abi`（handoff §1） |
| 绑定生成时机 | **构建期生成 + 提交到仓库**（`internal/blockchain/bindings/`），非运行时 | 可审查、可 diff、`go test` 离线可跑；提供 `make abigen` 脚本由 implement 阶段定义，从 contracts artifacts 取 ABI |
| owner key 管理 | **从 env `ORACLE_OWNER_PRIVATE_KEY` 读**，绝不硬编码、绝不入库、绝不进日志；启动时若缺失则**链上写功能降级关闭**（只读 + 事件同步仍可跑） | handoff §7：后端对外只需一个能调 `Oracle.monthlySettlement` 的 owner(=deployer) 私钥，链上权限传导由合约 setter 拓扑保证。本轮 env 注入，KMS 留后续 |
| 签名/发交易 | **`bind.NewKeyedTransactorWithChainID(privKey, chainID)`** + nonce 管理 + gas 估算 | 替换 `SignData` 的 SHA256；这是 go-ethereum 标准发交易路径 |
| 金额内部表示 | 链上/计价用 **`*big.Int`（6 位最小单位）**；DB 仍存 string（最小单位字符串）；展示层除 10^usdtDecimals | 避免浮点；与合约 `uint256` 对齐；精度从 `deployments.json.usdtDecimals` 读不硬编码 |
| 计价模拟器去留 | **保留 `OperatorAPISimulator.GetUsage`（产 usage），新增 `PricingService` 做 usage×费率**；废弃 `GetBill` 随机金额 | usage 仍可模拟（真实运营商 API 留后续），但金额必须确定可复算 |

---

## 4. 数据流与领域模型

### 4.1 月度结算主数据流（计价 → 上链 → 落库对账）

```
[触发] POST /api/oracle/monthly-bill (管理员鉴权 T7)
   │
   ▼
OracleServiceV2.FetchAndCreateBills() 重写：
   1. userRepo.FindAll() → 过滤激活服务用户
   2. 对每个 (user, operatorId)：
        usage = OperatorAPISimulator.GetUsage(userID, opID)   // (dataMB, callMin)
        amount6 = PricingService.Price(operatorId, usage)     // *big.Int, 6 位最小单位
        若 amount6 == 0 → 跳过（合约侧 amount>0 才 createBill）
   3. 构造三等长数组 (users[], operatorIds[], amounts[])
   4. 分批切片：每批 users ≤ 25（handoff §5.1 monthlySettlement 全链路最紧约束）
   5. 逐批：client.OracleMonthlySettlement(owner transactor, batch) → 发交易 → 等回执
        失败批：记录 + 幂等重试（合约逐用户 continue，issueMonthly 幂等）
   6. 本地 DB Bill 落库策略：见 §4.3「DB 与链的关系」
   │
   ▼（链上）Oracle.monthlySettlement → payment.createBill(emit BillCreated)
            + issueMonthlyTrafficCards(emit TrafficCardMinted)
            + applyTrafficCardToBill 桩(emit TrafficCardApplied)
            + per-bill emit UsageDataSubmitted
   │
   ▼（异步）event_sync FilterLogs 监听上述事件 → 6 位精度解码 → 回填/对账 DB Bill（IsPaid/TxHash/链上 billId）
```

### 4.2 计价费率表（数据结构，费率值占位，结构锁定）

费率表是后端新增的领域模型。结构如下（implement 阶段落为 Go struct + seed，可选后续转 DB 表）：

```go
// 计价单位约定：data 以 MB 计，call 以 minute 计；单价以 USDT 6 位最小单位/单位用量计
type OperatorRate struct {
    OperatorID    uint     // 对应 models.Operator.ID（与链上 operatorId 对齐，见 §4.5 风险）
    Region        string   // 冗余便于按地区批量配置
    DataUnitPrice *big.Int // 每 1 MB 流量单价（6 位最小单位 USDT）。占位示例：10000 = 0.01 USDT/MB
    CallUnitPrice *big.Int // 每 1 分钟通话单价（6 位最小单位 USDT）。占位示例：5000 = 0.005 USDT/min
    MinBillAmount *big.Int // 最小出账金额（低于此不出账，避免 dust）。占位示例：0
}

// PricingService.Price(operatorId, dataMB, callMin) → *big.Int(6位)
//   amount6 = dataMB * rate.DataUnitPrice + callMin * rate.CallUnitPrice
//   if amount6 < rate.MinBillAmount → return 0（不出账）
```

费率表数据来源（本轮）：**代码内 seed map（按 operatorId）**，与 `cmd/main.go` 的 11 运营商 seed 并列；费率值由产品/运营后续填真实数，本轮用占位值并在测试中固定断言可复算。

> 设计要点：金额必须**纯整数运算**（MB/min 为整数，单价为整数最小单位），杜绝浮点。usage 的 dataMB/callMin 来自 `OperatorAPISimulator.GetUsage`（本轮模拟，单位需在 implement 明确为 MB/min）。

### 4.3 DB 与链的关系（对账模型）

后端是 oracle/owner 角色，月度结算的**金额事实源在后端计价**，**资金事实源在链上**。三者关系：

- **Bill 金额**：后端计价算出，传入 `createBill`；event_sync 监听 `BillCreated(billId,user,totalAmount,platformFee)` 回填**链上 billId**（DB 当前无此字段，需补 `OnChainBillID` + `TxHash`）。
- **Bill 已付状态**：链上 `BillPaid` 事件为准；event_sync 监听后置 `IsPaid=true`+`PaidAt`+`TxHash`。废弃现状「前端传 txHash 即标记已付」的无确认路径（`BillingService.MarkAsPaid` 改为只受 event_sync 驱动，或保留但标注不可信）。
- **TrafficCardDeduction**：合约侧 `applyTrafficCardToBill` 为受限桩（handoff §6 不转资金）→ **本字段本轮恒为 "0"，不参与抵扣对账**。event_sync 监听 `TrafficCardApplied(billId)` 仅记录事件发生，不改金额。
- **Deposit/Withdraw**：event_sync 监听 `DepositMade(user,amount)`/`DepositWithdrawn(user,principal,interest)` 落库（6 位精度）；`/api/deposit` `/api/withdraw` 端点改为**记账由事件回填驱动**，前端调用仅做意向记录（标注 pending，事件到达后转 confirmed）。

### 4.4 模型字段调整（最小改动）

| 模型 | 调整 | 理由 |
|------|------|------|
| `models.Bill` | 新增 `OnChainBillID uint64`；`TxHash`/`PaidAt` 由事件回填；`TrafficCardDeduction` 恒 "0" | 对账需链上 billId 关联 |
| `models.Deposit` | `Amount` 语义统一为「6 位最小单位字符串」；可加 `Status`(pending/confirmed) | 精度一致 + 事件回填 |
| `models.Operator` | `RequiredDeposit` 语义重定义为 6 位最小单位（当前 "0.01" 含义不明）；与费率表 OperatorID 对应 | 精度一致 |
| `models.UsageData` | `DataUsage`/`CallUsage` 保持 uint64，**明确单位 = MB / minute** | 计价输入 |

> 字段调整属 implement 阶段执行，design 仅锁定语义。AutoMigrate 兼容（GORM 加列）。

### 4.5 operatorId 映射风险（需 arch-review 确认）

后端 `models.Operator.ID`（DB 自增）与合约 `ServiceManager` 的 `operatorId` **是否一致未经证实**。`monthlySettlement` 传的 `operatorIds[]` 必须是**链上 operatorId**，否则 `createBill` 关联错运营商/分账地址错。implement 阶段必须建立**后端 operatorID ↔ 链上 operatorId 映射**（读 `ServiceManager.getActiveOperators` 比对 name/region，或部署时固定映射表）。**列为 arch-review 重点。**

---

## 5. ABI / 接口对齐（对照 handoff 真实值）

### 5.1 函数 selector 对照（写调用，handoff §2 权威）

| 函数 | selector | 后端调用方 | 调用前置 | 现状 |
|------|----------|-----------|----------|------|
| `Oracle.monthlySettlement(address[],uint256[],uint256[])` | `0x01eb00ca` | 月度结算主入口（owner 签名） | owner=deployer 账户；三数组等长；每批 ≤25 | 完全未实现 |
| `Payment.createBill(address,uint256,uint256)` | `0xceb323e8` | **不直接调**（由 Oracle 内部 onlyOracle 调） | — | — |
| `Deposit.issueMonthlyTrafficCards(address[])` | `0x0948eaad` | 可选单独发卡（owner→Oracle，或经 monthlySettlement 内部） | 每批 ≤50 | 完全未实现 |
| `Deposit.deposit(uint256)` | `0xb6b55f25` | **后端不代发**（用户钱包侧）；后端只读 `getDepositAmount` | 用户先 `usdt.approve` | client stub |
| `Payment.payBill(uint256)` | `0xf0975190` | **后端不代发**（用户钱包侧） | 用户先 `usdt.approve(amount+fee)` | — |
| `ServiceManager.setOperatorPaymentAddress(uint256,address)` | `0xadb76801` | 部署侧一次性（onlyOwner） | 非零地址 | — |

> 关键判断：**后端只代发 `monthlySettlement`（owner 角色）**。`deposit`/`payBill` 是用户钱包侧操作（web 调），后端不持用户私钥、不代发。后端对这两者只做**只读 + 事件回填**。

### 5.2 事件 topic 对照（event_sync 监听，interfaces 权威）

| 事件签名（冻结 ABI） | 后端用途 | signatures.go 现状 | 改动 |
|----------------------|----------|---------------------|------|
| `UserRegistered(address,string,uint256)` | 注册落库（解析 email/tokenId） | ✅ 正确，但 process 不解析 data | 修 process 解码 email/tokenId/时间 |
| `DepositMade(address,uint256)` | 押金落库（6 位） | ✅ topic 对，但无落库 | 补 process |
| `DepositWithdrawn(address,uint256,uint256)` | 提现落库（principal+interest） | 🔴 缺 | 新增 topic + process |
| `BillCreated(uint256,address,uint256,uint256)` | 账单落库 + 回填 OnChainBillID | 🔴 **现状 5 参错** | 改 4 参 `BillCreated(uint256,address,uint256,uint256)` |
| `BillPaid(uint256,address,uint256,uint256)` | 标记已付（事实源） | ✅ 正确，无 process | 补 process（置 IsPaid/PaidAt/TxHash） |
| `TrafficCardMinted(address,uint256,uint256)` | NFT 发卡落库 | ✅ topic 对，无 process | 补 process |
| `TrafficCardApplied(uint256)` | 桩事件，仅记录（不改金额） | 🔴 缺 | 新增 topic + process（只记录） |
| `UsageDataSubmitted(user,operatorId,amount)` | monthlySettlement per-bill emit，对账金额 | 🔴 缺 | 新增 topic + process（对账校验后端计价 = 链上入账） |

> abigen 生成的 `*Filterer` 自带类型化事件解析，可替代手写 topic 常量；但 `signatures.go` 仍保留供轻量过滤/日志用，需同步修正避免误导。

### 5.3 abigen 绑定清单（8 份，handoff §1 来源）

从 `packages/contracts/artifacts/contracts/<Name>.sol/<Name>.json` 的 `.abi` 取全量：

| # | 合约 | 绑定用途 | abiHash 比对源 |
|---|------|----------|----------------|
| 1 | `Oracle` | monthlySettlement（写，owner）+ verifyServiceActive（读） | hardhat.json.abiHash.Oracle |
| 2 | `Payment` | createBill/payBill/getUnpaidBills（读）+ 事件 | .Payment |
| 3 | `Deposit` | getDepositAmount/getLockExpiry（读）+ issueMonthly（写）+ 事件 | .Deposit |
| 4 | `FeeManager` | getFeeRate/calculateFee（读，展示用） | .FeeManager |
| 5 | `ServiceManager` | getActiveOperators/getOperator（读，operatorId 映射） | .ServiceManager |
| 6 | `TrafficCardNFT` | getCardInfo/getUserCardCount（读）+ 事件 | .TrafficCardNFT |
| 7 | `UserRegistry` | isRegistered/getUserInfo（读）+ UserRegistered 事件 | .UserRegistry |
| 8 | `MockUSDT` | decimals/balanceOf/allowance（读，校验精度/余额） | （从 artifacts/mocks 取，正式 USDT 同 6 位） |

> implement 阶段提供 `make abigen` 或脚本：读 artifacts → abigen 生成到 `internal/blockchain/bindings/<name>.go`，并校验 `abiHash` 与 `deployments.json` 一致（不一致 fail，提示重新生成）。

---

## 6. 关键技术设计（按模块）

### 6.1 `blockchain/client.go` — abigen 绑定 + 真实读/写

- **结构扩展**：`Client` 增加 `bindings`（8 份 abigen 实例）、`transactor *bind.TransactOpts`（owner，可为 nil = 只读降级）、`usdtDecimals uint8`。
- **读调用（替换 stub 返零值）**：
  - `GetDepositAmount(addr)` → `bindings.Deposit.GetDepositAmount(callOpts, addr)`（返回 6 位最小单位 *big.Int）
  - `GetLockExpiry(addr)` → `bindings.Deposit.GetLockExpiry(...)`
  - `VerifyServiceActive(addr)` → `bindings.Oracle.VerifyServiceActive(...)`
  - 新增 `GetFeeRate()` / `CalculateFee(amount)`（FeeManager，展示用）
- **写调用（owner 签名，新增）**：
  - `MonthlySettlement(users []common.Address, opIds []*big.Int, amounts []*big.Int) (txHash, error)`：`bindings.Oracle.MonthlySettlement(transactor, ...)`；调用前断言 `transactor != nil`（否则返回 ErrNoOwnerKey），断言三数组等长、`len ≤ 25`。
  - `IssueMonthlyTrafficCards(users)`（可选单独发卡，≤50）。
  - 写调用统一：估 gas → 发交易 → 返回 txHash（回执确认交由调用方/event_sync 异步）。
- **nonce 管理**：单 owner 账户串行发交易（分批结算本就串行），用 `PendingNonceAt` 取 nonce；并发保护用 mutex（本轮单 worker，简单 mutex 足够）。

### 6.2 `blockchain/signatures.go` + owner key — 真实签名

- **`SignData`（services/oracle.go）废弃 SHA256**：本轮 `submitUsage` 是合约预留接口、v2-A 不参与计价（Oracle.sol:89），后端**不需要为计价做链上用量签名**。`UsageData.Signature` 字段可保留为本地审计哈希（标注非链上签名），或直接置空。**结论：本轮不实现链上用量 ECDSA 签名**，避免无用复杂度。
- **owner 交易签名**：在 `client.go` 用 `bind.NewKeyedTransactorWithChainID(privKey, big.NewInt(chainID))` 构造 transactor。私钥从 `os.Getenv("ORACLE_OWNER_PRIVATE_KEY")` 读（hex），`crypto.HexToECDSA` 解析。**禁止硬编码、禁止入库、禁止进日志**（日志只打 owner address 不打私钥）。缺失 → transactor=nil，链上写降级关闭，启动 log warn。
- `signatures.go` 修正：`BillCreated` 改 4 参；新增 `TrafficCardApplied`/`UsageDataSubmitted`/`DepositWithdrawn` topic（与 §5.2 一致）。

### 6.3 `sync/event_sync.go` — 真实订阅/轮询 + 6 位精度 + 多事件

- **真实拉取**：替换 `time.Sleep(30s)` 空转为**轮询 `FilterLogs`**（区块范围分页，从 last-synced block 推进；持久化 last block 到 DB 或文件，断点续传）。Arbitrum Sepolia 出块快，优先 `FilterLogs` 轮询（比 `SubscribeFilterLogs` 对公共 RPC 更稳）。
- **地址表来源**：`s.contracts`（修 config bug 后非空，§6.5）；按合约地址 + topic 过滤。
- **多事件分发**：`UserRegistered`/`DepositMade`/`DepositWithdrawn`/`BillCreated`/`BillPaid`/`TrafficCardMinted`/`TrafficCardApplied`/`UsageDataSubmitted` 各自 process，用 abigen `*Filterer.Parse*` 解码（替代手写 topic 索引解析）。
- **6 位精度落库**：所有金额事件解码出的 `*big.Int` 为 6 位最小单位，DB 存最小单位字符串；不做除法（展示层处理）。
- **process 修正**：`processUserRegistered` 解析 email/tokenId/真实时间（不再 `Unix(0,0)`）；`processBillCreated` 回填 `OnChainBillID`+`TxHash`；`processBillPaid` 置 `IsPaid`/`PaidAt`；`processUsageDataSubmitted` 对账（后端计价金额 == 链上入账金额，不一致告警）。
- **幂等**：按 (txHash, logIndex) 或 OnChainBillID 去重，避免重复落库（轮询重叠区块时）。

### 6.4 services 计价 — 费率表 + amounts[] + 分批

- 新增 `PricingService`（§4.2 费率表 + `Price()`）。
- `OracleServiceV2.FetchAndCreateBills` 重写（§4.1 流程）：
  - 用 `OperatorAPISimulator.GetUsage` 取 (dataMB, callMin)（保留模拟，单位明确 MB/min）。
  - `PricingService.Price(operatorId, usage)` → amount6（确定可复算，**废弃 `GetBill` 随机金额**）。
  - operatorId 映射（§4.5）：后端 opID → 链上 operatorId。
  - 构造三等长数组，**分批切片每批 ≤25**（handoff §5.1），逐批调 `client.MonthlySettlement`，失败批记录 + 幂等续跑。
- **handler 入口 `TriggerMonthlyBill`** 保持，但加鉴权（§6.6）+ 返回分批结果摘要（成功批数/失败批/txHashes）。

### 6.5 `config.go` — 键名 bug 修复 + schema 扩展

- **修 bug（二选一，建议对齐合约产物）**：struct tag `json:"proxies"` 保持，**deployments.json 键名从 `contracts` 改为 `proxies`**（与 `hardhat.json` 范式一致），消除双重不一致。
- **Deployments struct 扩展**：新增 `Usdt common.Address json:"usdt"`、`UsdtDecimals uint8 json:"usdtDecimals"`、`AbiHash map[string]string json:"abiHash"`。
- **abiHash 校验**：加载后比对 abigen 绑定的 ABI 指纹与 `deployments.json.abiHash`，不一致 → 启动 fail（提示重新生成绑定），防止 ABI 与链上实现脱节。
- **RPC 来源统一**：`main.go` 当前用 `os.Getenv("RPC_URL")` 起 sync、地址走 json。统一为 **json.rpcUrl 为准，env `RPC_URL` 可覆盖**（明确优先级，避免连 A 链发 B 链）。

### 6.6 handlers — 敏感端点鉴权中间件

- **新增 Gin 中间件 `AdminAuth`**：校验请求头 `X-Admin-Key` == `os.Getenv("ADMIN_API_KEY")`（常量时间比较 `subtle.ConstantTimeCompare`），缺失/不匹配 → 401。
- **加固端点**：`POST /api/oracle/monthly-bill`、`POST /api/usage/submit`、`POST /api/withdraw`（若仍保留写路径）、`POST /api/notification/send`。
- **不加固**：用户读端点（user/bills/deposit 查询）、register、operators 等公开端点保持。
- 中间件方案优先 API key（实现简单、本轮够用）；HMAC 签名校验留后续 Round 升级（标注）。
- `ADMIN_API_KEY` 从 env 读，缺失 → 启动 fail（不允许敏感端点裸奔）。

### 6.7 `configs/deployments.json` — 421614 schema + 本地 31337 机制

- **本地联调（先行）**：新增/对齐读 `packages/contracts/deployments/hardhat.json`（chainId 31337）的 `proxies`+`usdt`+`usdtDecimals`+`abiHash`，作为 implement/test 阶段离线联调来源。
- **Arbitrum Sepolia 421614**：`configs/deployments.json` schema 改为：
  ```json
  {
    "chainId": 421614,
    "rpcUrl": "<Arbitrum Sepolia RPC>",
    "proxies": { "FeeManager": "0x…", "...7 个待上链填": "" },
    "usdt": "0x… (MockUSDT，待上链填)",
    "usdtDecimals": 6,
    "abiHash": { "...7 合约指纹（取自 arbitrum_sepolia.json，待生成）": "" }
  }
  ```
- **占位 + 读 arbitrum_sepolia.json 机制**：合约 1/3 真·上链前，`deployments/arbitrum_sepolia.json` 不存在（handoff §10）。本轮 421614 地址先占位（零地址或注释），**implement 提供「从 contracts/deployments/<net>.json 同步到 backend/configs/deployments.json」的脚本**，上链后一键回填。test 阶段优先用 31337 hardhat.json 跑通。
- 同步修正 `.env.example`：`CHAIN_ID=421614`、`RPC_URL=<Arbitrum Sepolia>`，新增 `ORACLE_OWNER_PRIVATE_KEY=`（空，注释说明）、`ADMIN_API_KEY=`。

---

## 7. 非功能性设计

### 7.1 安全（资损敏感，arch-review/security-review 重点）

- **owner key 管理**：`ORACLE_OWNER_PRIVATE_KEY` 仅 env 注入；不硬编码、不入库、不进日志（日志仅 owner address）；缺失则写功能降级。私钥泄露 = 可任意发起月度结算（高危资损面）。**部署侧建议最小权限账户 + 后续 KMS**（本轮标注遗留）。
- **端点鉴权**：敏感端点 API key + 常量时间比较；`ADMIN_API_KEY` 缺失启动 fail。
- **资损 = 计价错算**：费率表纯整数运算（无浮点）；amounts[] 与链上 `UsageDataSubmitted` 事件双向对账（不一致告警）；`amount > 0` 才出账（合约侧也 fail-fast）。`MIN_DEPOSIT`/最小出账金额校验。
- **fee-on-transfer 禁用**（handoff §3）：后端记账信任 amount=实收，正式 USDT 无此问题；mock 阶段确认 MockUSDT 标准 ERC20。
- **applyTrafficCardToBill 桩**（handoff §6）：后端不据此做资金抵扣，`TrafficCardDeduction` 恒 0。

### 7.2 USDT 6 位精度一致性

- 全链路 6 位最小单位（链/计价 *big.Int，DB string，展示除 10^decimals）。
- `usdtDecimals` 从 `deployments.json` 读，**不硬编码 6**（正式 USDT 同 6 位，但读字段防漂移）。
- seed `RequiredDeposit` + 费率表单价 + `MIN_DEPOSIT(10×10^6)` 全部 6 位语义统一。

### 7.3 RPC 统一 Arbitrum

- 三端（hardhat / 后端 / 前端）统一 Arbitrum Sepolia 421614；后端 RPC = json.rpcUrl 为准 + env 可覆盖，单一优先级。
- 杜绝 0G 残留（旧 `evm-testnet.0g.ai` / `evmrpc-testnet.0g.ai` 全清）。

### 7.4 联调依赖（前置项，明确标注）

- ⏳ **端到端真机联调依赖合约 1/3 真·上链**（handoff §10：缺 DEPLOYER_PRIVATE_KEY/RPC，`arbitrum_sepolia.json` 未生成、Arbitrum MockUSDT 地址不存在）。
- **不阻塞 implement/test**：implement 阶段全部用本地 hardhat(31337) + hardhat.json 联调；abigen 绑定/计价/分批/事件解码/鉴权均可本地 + mock RPC 单测验证。真机 421614 联调待上链后补跑（test 阶段标注）。
- 上链后需复测：在 421614 校准分批上限 N（L2 calldata 计价差异，handoff §5.1）。

### 7.5 异常与监控/对账

- 分批失败：记录失败批 users + 错误，幂等续跑（合约逐用户 continue + issueMonthly 幂等）。
- 交易回执确认：发交易后异步等回执；event_sync 为最终对账事实源。
- 对账告警：`UsageDataSubmitted` 链上金额 ≠ 后端计价金额 → 告警日志（资损信号）。
- event_sync last-block 持久化 + 断点续传 + 幂等去重（防漏/重复）。

---

## 8. 落地计划（给 plan 阶段输入）

| 任务 | 内容 | 文件 | 依赖 | 风险 |
|------|------|------|------|------|
| **T1** abigen 绑定 | make 脚本 + 生成 8 份绑定 + abiHash 校验 | `internal/blockchain/bindings/*`, Makefile | 合约 artifacts 存在 | ABI 与链上不符 → revert |
| **T2** config 修复 | 键名 bug + struct 扩展(Usdt/UsdtDecimals/AbiHash) + RPC 优先级 + .env.example | `config/config.go`, `configs/deployments.json`, `.env.example` | — | 键名不修 sync 静默失效 |
| **T3** client 读/写 | 读调用真实化 + MonthlySettlement/IssueMonthly 写 + nonce mutex | `blockchain/client.go` | T1,T2 | selector/编码错 |
| **T4** owner 签名 | env 私钥 + transactor + 降级关闭 + 日志脱敏 | `blockchain/client.go`, `cmd/main.go` | T3 | 私钥泄露/硬编码 |
| **T5** event_sync | FilterLogs 轮询 + 8 事件 process + 6 位精度 + 幂等 + last-block 持久化 + signatures.go 修正 | `sync/event_sync.go`, `signatures.go`, `models.go`(加列) | T1,T2 | 漏/重复落库 |
| **T6** 计价 | PricingService 费率表 + FetchAndCreateBills 重写 + operatorId 映射 + amounts[] + 分批≤25 | `services/oracle.go`(+pricing) | T3,T4 | 计价错算/映射错 |
| **T7** 鉴权 | AdminAuth 中间件 + 敏感端点挂载 + ADMIN_API_KEY env | `handlers/handlers.go`(+middleware), `cmd/main.go` | — | 端点裸奔 |
| **T8** 测试 | 计价复算/分批切片/6位精度/事件解码/鉴权中间件 单测 + mock RPC | `services/*_test.go`, `blockchain/*_test.go` | T1–T7 | 覆盖不足 |

**顺序**：T1→T2→(T3→T4)→T5→T6→T7→T8。T5 与 T6 串行（都依赖 T1/T3）。T7 可与 T6 并行（不同文件）。**implement 阶段始终串行**（CLAUDE.md 规则）。

**风险总览**：① operatorId 映射（§4.5，arch-review 必审）；② owner key 安全（§6.2/7.1）；③ 计价正确性（§4.2/7.1）；④ 联调依赖上链（§7.4）；⑤ 事件幂等/精度（§6.3/7.2）。

---

## 9. arch-review / security-review 重点清单

**arch-review 必须确认**：
1. **operatorId 映射**（§4.5）：后端 DB ID 与链上 operatorId 是否一致？映射建立方式（读 ServiceManager 比对 vs 固定表）是否可靠？传错 = createBill 关联错运营商/分账。
2. **后端只代发 monthlySettlement**（§5.1）的判断是否正确？deposit/payBill 确属用户钱包侧、后端不代发？
3. **DB 与链对账模型**（§4.3）：废弃「前端传 txHash 即标记已付」改事件驱动，是否影响现有 web 流程（需与 web 3/3 对齐）？
4. **计价费率表结构**（§4.2）：纯整数 6 位最小单位运算是否覆盖所有计价场景？MinBillAmount/dust 处理？
5. **submitUsage 不实现链上签名**（§6.2）的取舍是否合理（v2-A 不参与计价）？

**security-review 必须审**：
1. **owner key 全链路**（env→transactor→日志脱敏→降级）无泄露面；硬编码/入库/日志泄露 grep 核查。
2. **AdminAuth 中间件**：常量时间比较、缺 key 启动 fail、覆盖全部敏感端点无遗漏。
3. **资损路径**：计价错算、amounts[] 与链上事件对账、amount>0 校验、6 位精度一致性、fee-on-transfer 信任假设。
4. **分批幂等**：失败续跑不重复结算（合约幂等 + 后端去重）。
5. **RPC 统一**：无 0G 残留、env/json 优先级单一、不连错链。

---

## 10. 遗留 / 依赖

- ⏳ 端到端真机 421614 联调 + 分批上限校准 → **依赖合约 1/3 真·上链**（handoff §10），不阻塞 implement/test（本地 31337）。
- ⏳ KMS 私钥托管、HMAC 端点签名、真实运营商 usage API、NotificationService SMTP → 留后续 Round。
- ⏳ 费率表真实费率值 → 产品/运营填，本轮占位 + 测试固定断言。
