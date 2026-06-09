# Stage: design — 后端技术设计（子项目 backend 2/3）

> **状态**: v2（按 arch-review 返工，闭合 B1–B5 + ⚠️ 收紧项） | **日期**: 2026-06-09 | **Gate**: 1（设计） | **子项目**: backend(2/3) | **角色**: 架构师 | **分支**: backend/align-arbitrum-usdt
> **输入**：合约 handoff-backend.md（冻结 ABI/selector + USDT 6 位 + 计价归后端 + 批量上限 N）+ scan.md（后端 8 大发现）+ 4 份基线（project-scan/blockchain-integration/services-api/alignment-surface）+ requirement.md §三②/§五B + 真实源码逐文件核对 + **arch-review.md（BLOCKED：5 ❌ + 6 ⚠️）**。
> **用户已拍板锁定 3 决策**（见 §1.3），本设计据此展开，不再二次澄清。
> **注**：合约子项目(1/3) 的 design 见 git commit `2b8d4a0`（本 stage 文件按子项目串行复用，前版已入库）。
>
> ⚠️ 本文档只产设计，不写 .go 代码。implement 阶段按 §8 任务划分执行。

---

## 0. v2 返工说明（相对 v1 的变更）

v1 在 arch-review 被 **BLOCK**：安全审计 5 ❌ + CEO 1 ❌（与 B2 重叠）+ Eng 6 ⚠️。核心批评：**资损「预防控制」系统性缺失**——v1 全篇依赖事后对账告警（检测）而非「发交易前硬上限 / reorg 确认 / 单一对账路径」（预防）。

v2 逐条闭合，关键变更：

1. **新增 §7.0「资损预防控制」专节**：把「金额硬闸 / usage 上界 / reorg 确认 / 单一对账路径」集中表述，正面回应「预防 vs 检测」批评。
2. **新增 §11「arch-review 阻塞闭合对照表」**：B1–B5 各自怎么解，逐条可核。
3. **§4.3 对账模型重写**：删除「保留但标注不可信」的二选一，**IsPaid / withdraw 记账唯一由链上事件回填**；HTTP 写端点降级为 pending 意向，绝不直接置终态。
4. **§6.6 鉴权清单重排**：`/api/bills/pay`、`/api/withdraw` 改钱包签名鉴权（用户操作）而非 AdminAuth；补全裸奔端点。
5. **§7.0 + §11 锁定金额硬闸常量**：`MAX_BILL_PER_USER` / 单批总额上限 / 异常熔断 / usage 上界。
6. **§6.3 event_sync 加 reorg + 两阶段确认**（pending→confirmed，等 K 块）。
7. **§4.5 operatorId 固定 ID 映射**：seed 显式写 `ID = 链上 operatorId`，启动 sanity check，不靠 name 比对。
8. **新增 `handoff-web.md`**（docs/design/linkworld-backend/handoff-web.md，**移交 web 3/3**）：exact-amount approve 禁 infinite、对账契约 breaking change、deposit/withdraw 状态机。
9. owner key 不落 .env 明文长存（启动注入内存）；§7.1 措辞升级 owner = 平台 root 权限。
10. ⚠️ 收紧项全部落地（abigen Filterer 解码、nonce 计数器 + WaitMined、simulated.Backend、T1 前置拆分、验收语言澄清、幂等键、crypto/rand、CORS 收口等，见 §11 下半。

### 0.1 v2 追加锁定决策（在 §1.3 三决策之上）

| # | 锁定项 | 取值 / 口径 |
|---|--------|------------|
| L4 | 单 user 单 bill 金额硬上限 | `MAX_BILL_PER_USER`（config 可审计，占位 1000 USDT = 1_000_000_000 最小单位），超界**拒绝该 bill + 告警**，不发交易 |
| L5 | 单批总额硬上限 + 异常熔断 | 单批 `amounts[]` 求和 ≤ `MAX_BATCH_TOTAL`（占位 5000 USDT）；且单批总额 > 历史月度均值 × `N`（占位 3 倍）→ **熔断暂停，人工放行** |
| L6 | usage 上界 | `MAX_DATA_MB` / `MAX_CALL_MIN`（config，占位 1_000_000 MB / 100_000 min），PricingService 超界拒绝 + 告警；SubmitUsage gin binding `max=` 同步校验 |
| L7 | 对账单一路径 | `IsPaid`/`PaidAt`、withdraw 记账 **唯一**由 event_sync 监听链上事件回填；HTTP 端点最多写 pending 意向，**绝不置终态** |
| L8 | 资金事件两阶段 + reorg | 资金事件 pending（监听到）→ confirmed（等 `CONFIRMATIONS` K 块，占位 5）；reorg 检测（last-block 记 blockHash，父哈希不连续→回退重扫） |
| L9 | owner key 不长存明文 | 启动时从注入源读入内存（env 仅本地/CI；生产经 secret manager 注入进程环境，进程退出即失），日志只打 owner address |
| L10 | operatorId 固定映射 | seed 显式写 `Operator.ID = 链上 operatorId(1..11)`，启动读链 sanity check 不一致 fail/warn，**不靠 name 比对** |

> 占位常量值（1000 / 5000 / 3× / K=5 等）由产品/运营/安全在 implement 前确认；本轮代码内带刺眼 `PLACEHOLDER` 注释 + 启动 warn，测试用固定值断言。

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
| owner key 从注入源读入内存（不长存明文）+ `bind.NewKeyedTransactorWithChainID` 真实签名发交易 | ❌ 不接 phone 字段 / 本地号码注册（requirement R10） |
| 敏感端点鉴权：写端点分两类——管理触发端点（oracle/monthly-bill、notification/send）走 AdminAuth；用户资金端点（bills/pay、withdraw）走钱包签名鉴权 | ❌ 不引入 KMS（本轮注入内存私钥 + 部署侧最小权限，KMS 留后续 Round） |
| **资损预防硬闸（v2 新增）**：金额上限 `MAX_BILL_PER_USER` + 单批总额 + 异常熔断 + usage 上界，发交易前校验（§7.0） | ❌ 不做实时风控引擎（本轮静态阈值 + 熔断，动态风控留后续） |
| **对账单一路径（v2 新增）**：IsPaid/withdraw 唯一由 event_sync 链上事件回填；HTTP 写端点降级 pending 意向（§4.3/§6.6） | |
| **reorg 防护（v2 新增）**：资金事件等 K 块确认 + 父哈希连续性检测（§6.3） | |
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
        ── B4 上界校验：usage 超 MAX_DATA_MB/MAX_CALL_MIN → 拒绝该条 + 告警 ──
        amount6 = PricingService.Price(operatorId, usage)     // *big.Int, 6 位最小单位
        ── B1 金额硬闸：amount6 > MAX_BILL_PER_USER → 拒绝该 bill + 告警，不入数组 ──
        若 amount6 == 0 → 跳过（合约侧 amount>0 才 createBill）
   3. 构造三等长数组 (users[], operatorIds[], amounts[])（operatorId = §4.5 固定映射）
   4. 分批切片：每批 users ≤ 25（handoff §5.1 monthlySettlement 全链路最紧约束）
   5. ── B1 单批熔断：每批 sum(amounts) > MAX_BATCH_TOTAL 或 > 历史月均 × N → 熔断暂停，人工放行 ──
   6. 逐批（幂等键 = month+batchIndex，落 DB 已确认批不重发）：
        client.MonthlySettlement(owner transactor, batch)
        → 本地 nonce 计数器取 nonce（发一笔 +1）→ 发交易 → **WaitMined 等回执** → 据回执判成败
        失败批：记录 + 幂等重试（合约逐用户 continue，issueMonthly 幂等）
   7. 本地 DB Bill 落库策略：见 §4.3「DB 与链的关系」（终态字段不在此置，待 event_sync 回填）
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
// 计价单位约定（锁定）：data 以 MB 计，call 以 minute 计；单价以 USDT 6 位最小单位/单位用量计
// !!! PLACEHOLDER 费率，非真实定价，上线前必须由运营填实数并挂真机验收 gate；启动时 log.Warn 提示占位 !!!
type OperatorRate struct {
    OperatorID    uint     // == models.Operator.ID == 链上 operatorId（§4.5 固定映射，非 name 比对）
    Region        string   // 冗余便于按地区批量配置
    DataUnitPrice *big.Int // 每 1 MB 流量单价（6 位最小单位 USDT）。PLACEHOLDER：10000 = 0.01 USDT/MB
    CallUnitPrice *big.Int // 每 1 分钟通话单价（6 位最小单位 USDT）。PLACEHOLDER：5000 = 0.005 USDT/min
    MinBillAmount *big.Int // 最小出账金额（低于此不出账，避免 dust）。PLACEHOLDER：0
}

// PricingService.Price(operatorId, dataMB, callMin) → (*big.Int(6位), error)
//   ── B4 上界（fail-fast，非静默）：
//      if dataMB > MAX_DATA_MB || callMin > MAX_CALL_MIN → return err（拒绝 + 告警）
//   amount6 = dataMB * rate.DataUnitPrice + callMin * rate.CallUnitPrice   // 纯 big.Int，无浮点
//   ── B1 单 bill 硬上限：
//      if amount6 > MAX_BILL_PER_USER → return err（拒绝 + 告警，不出账）
//   if amount6 < rate.MinBillAmount → return 0（不出账）
```

费率表数据来源（本轮）：**代码内 seed map（按 operatorId）**，与 `cmd/main.go` 的 11 运营商 seed 并列；费率值由产品/运营后续填真实数，本轮用 **PLACEHOLDER 占位值**（注释刺眼 + 启动 `log.Warn("PRICING TABLE IS PLACEHOLDER...")`），并在测试中固定断言可复算。

> 设计要点：① 金额必须**纯整数运算**（MB/min 为整数，单价为整数最小单位），杜绝浮点；② usage 的 dataMB/callMin 来自 `OperatorAPISimulator.GetUsage`（本轮模拟 rand，**单位锁定 MB/min，测试固定断言**）；③ 计价上界是 B4 的第一道闸（Price 内拒绝），金额上限是 B1 的第二道闸（amount6 > MAX_BILL_PER_USER 拒绝），两闸独立。

### 4.3 DB 与链的关系（对账模型，v2 重写：单一路径 + 两阶段状态机）

后端是 oracle/owner 角色，月度结算的**金额事实源在后端计价**，**资金事实源唯一在链上**。

> **v2 铁律（闭合 B2/B3）：所有「资金终态」字段（Bill.IsPaid / PaidAt、Deposit/Withdraw 记账）唯一由 event_sync 监听链上事件回填。HTTP 写端点不接受前端 txHash 作为记账依据，最多写 pending 意向，绝不直接置终态。** 删除 v1「保留但标注不可信」的二选一——不保留任何旁路写终态的代码路径。

**资金记录两阶段状态机（统一 §6.3 reorg）**：

```
                 HTTP 意向(可选)            event_sync 监听到          等 K 块确认
   [无记录] ───────────────────────▶ [pending] ──────────▶ [seen] ──────────▶ [confirmed]
              (绝不置终态字段)                                  │ reorg(父哈希断)
                                                              ▼
                                                          [回退重扫] → 删 seen 记录，重新 FilterLogs
```

各字段关系：

- **Bill 金额**：后端计价算出 → 传入 `monthlySettlement.amounts[]`（发交易前先过 §7.0 金额硬闸）；event_sync 监听 `BillCreated(billId,user,totalAmount,platformFee)` 回填**链上 billId**（DB 补 `OnChainBillID`+`TxHash`）+ 校验 `totalAmount` == 后端计价（不一致告警）。
- **Bill 已付状态（B2 核心）**：唯一以链上 `BillPaid` 事件为准。event_sync 监听 → 资金事件须等 `CONFIRMATIONS`(K) 块（§6.3）→ 置 `IsPaid=true`+`PaidAt`+`TxHash`。**`BillingService.MarkAsPaid` 仅供 event_sync 内部调用**（不导出给 HTTP handler）。`POST /api/bills/pay` 改为：(a) 删写能力——直接 410/404 弃用；或 (b) 降级为 pending 意向记录（写 `Bill.PayIntentTxHash` + 不动 `IsPaid`，等 `BillPaid` 事件确认）。**本轮取 (b)**，保留前端「我已发起支付」UX，但 IsPaid 绝不由此置。
- **TrafficCardDeduction**：合约侧 `applyTrafficCardToBill` 为受限桩（handoff §6 不转资金）→ **本字段本轮恒为 "0"，不参与抵扣对账**。event_sync 监听 `TrafficCardApplied(billId)` 仅记录事件发生，不改金额。
- **Deposit/Withdraw（B3 核心）**：event_sync 监听 `DepositMade(user,amount)`/`DepositWithdrawn(user,principal,interest)` 落库（6 位精度，等 K 块确认）。`/api/deposit` `/api/withdraw` 端点**不接受前端 txHash 记账**：
  - `withdraw` 记账唯一由 `DepositWithdrawn` 事件回填；`RecordWithdraw(wallet, txHash)` 现状（凭前端 txHash 直接写记账）**废弃**。
  - HTTP 端点最多写 pending 意向（标注 `Status=pending`，事件到达 confirmed 后回填真实 principal/interest），且该意向**不计入余额计算**（`GetTotalByUserID` 只算 confirmed）。
- **breaking change 声明**：`/api/bills/pay`、`/api/withdraw` 的语义从「写终态」变「写意向」，是**对 web 3/3 的接口契约 breaking change**，详见 `docs/design/linkworld-backend/handoff-web.md`。

### 4.4 模型字段调整（最小改动）

| 模型 | 调整 | 理由 |
|------|------|------|
| `models.Bill` | 新增 `OnChainBillID uint64`；新增 `PayIntentTxHash string`(pending 意向，与 IsPaid 解耦)；`TxHash`/`PaidAt`/`IsPaid` **仅事件回填**；`TrafficCardDeduction` 恒 "0" | 对账需链上 billId 关联 + 意向与终态解耦（B2） |
| `models.Deposit` | `Amount` 语义统一为「6 位最小单位字符串」；**新增 `Status`(pending/confirmed)**；新增 `BlockHash`(reorg 检测) | 精度一致 + 两阶段事件回填（B3/B5） |
| `models.Operator` | `RequiredDeposit` 语义重定义为 6 位最小单位整数字符串（当前 "0.01" 含义不明，需存量迁移脚本）；`ID` 固定为链上 operatorId（见 §4.5） | 精度一致 + operatorId 映射（⚠️） |
| `models.UsageData` | 新增校验：`DataUsage ≤ MAX_DATA_MB`、`CallUsage ≤ MAX_CALL_MIN`（B4 上界） | usage 上界防天文账单 |
| `models.UsageData` | `DataUsage`/`CallUsage` 保持 uint64，**明确单位 = MB / minute** | 计价输入 |

> 字段调整属 implement 阶段执行，design 仅锁定语义。AutoMigrate 兼容（GORM 加列）。

### 4.5 operatorId 固定 ID 映射（v2 锁定，闭合 ⚠️ 最高优先项）

**决策（不靠 name 比对）**：后端 seed **显式写 `Operator.ID = 链上 operatorId`（1..11，合约 `initialize` 冻结的固定序号）**，即后端 DB 主键直接等于链上 operatorId。`monthlySettlement` 传 `operatorIds[]` 时直接用 `models.Operator.ID`，无需任何转换/比对。

为什么不靠 name 比对：后端 seed 的 `"T-Mobile"` 与链上 `"T-Mobile US"` 字符串不相等会静默 miss → 分账打到错误运营商地址（**最隐蔽资损**）。name 比对脆弱，固定 ID 是单一事实源。

**启动期 sanity check（防漂移）**：启动时读链 `ServiceManager.getActiveOperators` / `getOperator(id)`，逐个校验「后端 seed 的 operatorId 在链上存在且 paymentAddress 非零」：
- 不一致 → 默认 **fail-fast**（拒启动，避免带病结算）；可配置降级为 warn（仅非资金环境）。
- 链上 paymentAddress == 零地址 → 该 operator 标记不可结算，跳过其 bill（合约侧 `createBill` 也会 fail-fast，但后端提前拦截避免整批浪费 gas）。

> seed 落地：`cmd/main.go` 的 11 运营商 seed 显式指定 `ID: 1..11`（GORM 允许显式主键），与链上 `ServiceManager.initialize` 的注入顺序一一对应。implement 阶段需与合约 1/3 的 initialize 顺序核对（取自 `packages/contracts/scripts/deploy.ts` 或 ServiceManager 初始化清单）。

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
  - `MonthlySettlement(users []common.Address, opIds []*big.Int, amounts []*big.Int) (receipt, error)`：`bindings.Oracle.MonthlySettlement(transactor, ...)`；调用前断言链：
    - `transactor != nil`（否则 ErrNoOwnerKey）；
    - 三数组等长、`len ≤ 25`；
    - **每个 amount > 0 且 ≤ MAX_BILL_PER_USER**（B1 client 侧最后一道防线，即使上游漏校验也挡）；
    - **sum(amounts) ≤ MAX_BATCH_TOTAL**（B1 单批闸）。
  - `IssueMonthlyTrafficCards(users)`（可选单独发卡，≤50）。
  - 写调用统一：估 gas → 发交易 → **`bind.WaitMined` 等回执** → 据 `receipt.Status` 判成败后返回（v1「回执交由异步」与 §4.1 表述不一致，v2 统一：发交易方必须 WaitMined 确认上链成功，event_sync 仅做事件落库对账，不替代回执判定）。
- **nonce 管理（统一 §4.1，闭合 ⚠️）**：**本地 nonce 计数器**——初始化/恢复时用 `PendingNonceAt` 取一次，之后每发一笔本地 +1（不每次查 RPC，避免分批连发时 pending 未刷新导致 nonce 复用）。单 owner 账户串行发交易，mutex 保护计数器。每批 WaitMined 后再发下一批。

### 6.2 `blockchain/signatures.go` + owner key — 真实签名

- **`SignData`（services/oracle.go）废弃 SHA256**：本轮 `submitUsage` 是合约预留接口、v2-A 不参与计价（Oracle.sol:89），后端**不需要为计价做链上用量签名**。`UsageData.Signature` 字段可保留为本地审计哈希（标注非链上签名），或直接置空。**结论：本轮不实现链上用量 ECDSA 签名**，避免无用复杂度。
- **owner 交易签名 + key 不长存明文（闭合 B1 一环）**：在 `client.go` 用 `bind.NewKeyedTransactorWithChainID(privKey, big.NewInt(chainID))` 构造 transactor。
  - 私钥**启动时注入进程内存**：本地/CI 经 `os.Getenv("ORACLE_OWNER_PRIVATE_KEY")`（hex）读入，生产经 secret manager 注入进程环境变量（进程退出即失，不写磁盘）。**`.env` 仅本地用且 gitignore，`.env.example` 留空**——禁止明文私钥长存于版本库/配置文件。
  - `crypto.HexToECDSA` 解析后立即清空原始字符串引用（尽力而为）。
  - **禁止硬编码、禁止入库、禁止进日志**（日志只打 owner address 不打私钥）。
  - **chainID 一致校验**：`NewKeyedTransactorWithChainID` 用的 chainID 必须 == 链上 `ethclient.ChainID()`，不一致 fail-fast（防签错链重放）。
  - 缺失 → transactor=nil，链上写降级关闭，启动 log warn。
- `signatures.go` 修正：`BillCreated` 改 4 参；新增 `TrafficCardApplied`/`UsageDataSubmitted`/`DepositWithdrawn` topic（与 §5.2 一致）。

### 6.3 `sync/event_sync.go` — 真实订阅/轮询 + 6 位精度 + 多事件

- **真实拉取**：替换 `time.Sleep(30s)` 空转为**轮询 `FilterLogs`**（区块范围分页，从 last-synced block 推进；持久化 last block 到 DB，断点续传）。Arbitrum Sepolia 出块快，优先 `FilterLogs` 轮询（比 `SubscribeFilterLogs` 对公共 RPC 更稳）。
- **地址表来源**：`s.contracts`（修 config bug 后非空，§6.5）；按合约地址 + topic 过滤。**占位零地址不做 event 过滤**（零地址合约视为未上链，拒绝订阅其事件，避免误匹配；deployments.json 关键地址为零 → 该合约事件同步跳过 + warn）。
- **多事件分发**：`UserRegistered`/`DepositMade`/`DepositWithdrawn`/`BillCreated`/`BillPaid`/`TrafficCardMinted`/`TrafficCardApplied`/`UsageDataSubmitted` 各自 process，**一律走 abigen `*Filterer.Parse*` 解码**（类型安全，替代手写 topic 索引解析；`signatures.go` 仅日志/轻量过滤用，不参与字段解码）。
  - **解码歧义澄清（闭合 ⚠️）**：`UsageDataSubmitted(user,operatorId,amount)` **只有 `user` indexed**（operatorId/amount 在 data 区，不能当 topic 取）；`BillCreated(billId,user,totalAmount,platformFee)` 的**第三参 totalAmount = amount + platformFee 含费总额**，对账时勿当裸 amount。
- **6 位精度落库**：所有金额事件解码出的 `*big.Int` 为 6 位最小单位，DB 存最小单位字符串；不做除法（展示层处理）。string→big.Int 转换 `SetString(s,10)` 必须校验 `ok`，失败 fail-fast 不静默。

- **reorg 防护 + 两阶段确认（闭合 B5，统一 §4.3 状态机）**：
  1. **资金事件须等确认**：`DepositMade`/`DepositWithdrawn`/`BillPaid`/`BillCreated` 解码后**不立即置终态**——先落 `pending/seen`，待该 log 所在块 ≥ `latestBlock - CONFIRMATIONS(K=5)` 才置 `confirmed`（IsPaid/记账生效）。非资金事件（UserRegistered/TrafficCardMinted/TrafficCardApplied）可即时落库（无资损面）。
  2. **reorg 检测**：last-synced 记录不仅存 blockNumber，还存该块 `blockHash`。每轮拉取前校验「上次记录的块 hash 在当前链上仍存在 / 新块父哈希连续」；若父哈希断裂（reorg）→ **回退到分叉点前重扫**，删除受影响区间的 `seen`（未 confirmed）记录后重新 FilterLogs。已 confirmed（深度 ≥ K）的记录视为最终，不回退。
- **process 修正**：`processUserRegistered` 解析 email/tokenId/真实时间（不再 `Unix(0,0)`）；`processBillCreated` 回填 `OnChainBillID`+`TxHash`+校验金额；`processBillPaid` 置 `IsPaid`/`PaidAt`（**唯一置 IsPaid 的路径**，B2）；`processDepositWithdrawn` 回填 withdraw 记账（**唯一 withdraw 记账路径**，B3）；`processUsageDataSubmitted` 对账（后端计价金额 == 链上入账金额，不一致告警）。
- **幂等**：按 (txHash, logIndex) 去重，避免重复落库（轮询重叠区块时）；结算批次另有 (month, batchIndex) 幂等键（§4.1）。

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

### 6.6 handlers — 敏感端点鉴权中间件（v2 重排：两类鉴权 + 鉴权清单）

v1 错把 `withdraw` 归 AdminAuth（鉴权模型错配：withdraw 是用户操作，B3）。v2 按「操作主体」分两类鉴权：

| 端点 | 鉴权类型 | 理由 | 现状 |
|------|----------|------|------|
| `POST /api/oracle/monthly-bill` (TriggerMonthlyBill) | **AdminAuth** | 平台触发结算，管理操作 | 裸奔 |
| `POST /api/notification/send` | **AdminAuth** | 平台触发通知 | 裸奔 |
| `POST /api/usage/submit` (SubmitUsage) | **AdminAuth** | 平台/oracle 写 usage（非终端用户自助）+ **gin binding `max=` 上界校验**（B4） | 裸奔、无范围校验 |
| `POST /api/bills/pay` (PayBill) | **钱包签名（WalletAuth）** | 用户操作；且**仅写 pending 意向不置 IsPaid**（B2，§4.3） | 裸奔、直接置 IsPaid（白嫖面） |
| `POST /api/withdraw` (Withdraw) | **钱包签名（WalletAuth）** | 用户操作；且**不接受 txHash 记账，仅 pending**（B3，§4.3） | 裸奔、凭前端 txHash 记账 |
| `POST /api/deposit` (Deposit) | **钱包签名（WalletAuth）** | 用户操作（意向记录） | 裸奔 |
| `POST /api/service/activate` `deactivate` 等用户写端点 | WalletAuth | 用户操作 | 裸奔 |
| 读端点（GET user/bills/deposit/usage 查询）、register、operators | 公开 | 无写/无资损 | 公开（保持） |

- **`AdminAuth` 中间件**：校验请求头 `X-Admin-Key` == `ADMIN_API_KEY`（**常量时间比较 `subtle.ConstantTimeCompare`**），缺失/不匹配 → 401。`ADMIN_API_KEY` 注入内存，缺失 → 启动 fail（不允许管理端点裸奔）。
- **`WalletAuth` 中间件**：用户用钱包私钥对「请求体 + nonce/timestamp」签名，后端 `ecrecover` 还原地址 == 请求 `wallet` 字段（绑定 msg.sender 语义），防止他人冒充任意 wallet 提交意向。本轮实现钱包签名校验（非 AdminAuth）；防重放用 nonce/时间窗。
- HMAC/更强签名方案留后续 Round 升级（标注）。

> 这是对 web 3/3 的 breaking change（pay/withdraw 需带钱包签名头），详见 handoff-web.md。

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

### 7.0 资损预防控制（v2 新增专节，正面回应安全审计「预防 vs 检测」核心批评）

> arch-review 核心批评：v1 全篇依赖**事后对账告警（检测）**，缺**发交易前的硬控制（预防）**。本节集中表述四类预防控制，均在「资金动作发生前 / 终态落库前」生效，而非事后发现。

**① 金额硬闸（B1，预防凭空造大单）**——三层防线，发交易前层层拦：

| 层 | 位置 | 校验 | 越界动作 |
|----|------|------|----------|
| L1 计价层 | `PricingService.Price` | `amount6 ≤ MAX_BILL_PER_USER` | 拒绝该 bill + 告警 |
| L2 组批层 | `FetchAndCreateBills` 组批后 | 单批 `sum(amounts) ≤ MAX_BATCH_TOTAL`；且 ≤ 历史月均 × N（异常熔断） | 熔断暂停，人工放行 |
| L3 client 层 | `client.MonthlySettlement` 发交易前 | 每 amount ∈ (0, MAX_BILL_PER_USER]；sum ≤ MAX_BATCH_TOTAL | 拒发交易（最后防线，即使上游漏校验也挡） |

常量（config 可审计，PLACEHOLDER）：`MAX_BILL_PER_USER`=1000 USDT、`MAX_BATCH_TOTAL`=5000 USDT、熔断倍数 N=3。

**② usage 上界（B4，预防天文账单输入）**：`MAX_DATA_MB`/`MAX_CALL_MIN`（config），三处校验——`SubmitUsage` gin binding `max=`、`PricingService.Price` 入口、模拟器产值不超界。超界拒绝 + 告警。

**③ reorg 确认（B5，预防虚高已付/押金）**：资金事件落 `confirmed` 前等 `CONFIRMATIONS`(K=5) 块；reorg（父哈希断）回退重扫未确认记录（§6.3）。监听到 ≠ 终态，深度够才生效。

**④ 单一对账路径（B2/B3，预防伪造记账白嫖）**：资金终态字段（IsPaid/withdraw 记账）**唯一**由 event_sync 链上事件回填；HTTP 写端点只写 pending 意向不置终态（§4.3）。无任何旁路写终态的代码路径。

> 预防（本节）+ 检测（§7.5 对账告警）双层：硬闸/确认/单一路径在事前挡住资损，对账告警在事后兜底发现漏网。

### 7.1 安全（资损敏感，arch-review/security-review 重点）

- **owner = 平台 root 权限（措辞升级，B1）**：owner=deployer 不只是「结算权限」，而是**平台 root 权限**——能调 `setOperatorPaymentAddress`(改分账地址)、`Oracle.setPayment`/`Payment.setOracle`(改授权拓扑)、`monthlySettlement`(凭空造单 amounts[] 扣款)。私钥泄露 = 平台被接管（可改分账地址把全平台资金导走 + 任意造单），是**最高危资损面**。故 §7.0 金额硬闸是「即使 owner key 被盗用 / 误操作也限制单笔/单批爆炸半径」的纵深防御。
- **owner key 管理（L9）**：启动注入进程内存（本地/CI env、生产 secret manager），**不落 .env 明文长存、不硬编码、不入库、不进日志**（日志仅 owner address）；缺失则写功能降级。chainID 一致校验防签错链。部署侧最小权限账户 + 后续 KMS（标注遗留）。
- **端点鉴权（§6.6）**：AdminAuth（管理端点，常量时间比较，缺 key 启动 fail）+ WalletAuth（用户端点，钱包签名 ecrecover 绑 wallet）。
- **资损 = 计价错算**：费率表纯整数运算（无浮点）；amounts[] 与链上 `UsageDataSubmitted` 事件双向对账（不一致告警）；`amount > 0 且 ≤ MAX_BILL_PER_USER` 才出账（合约侧也 fail-fast）。`MIN_DEPOSIT`/最小出账金额校验。
- **fee-on-transfer 禁用**（handoff §3）：后端记账信任 amount=实收，正式 USDT 无此问题；mock 阶段确认 MockUSDT 标准 ERC20。
- **applyTrafficCardToBill 桩**（handoff §6）：后端不据此做资金抵扣，`TrafficCardDeduction` 恒 0。
- **其它收紧（⚠️）**：虚拟号 password 用 `crypto/rand`（非弱随机 `math/rand`）；**CORS 生产禁 `*` + credentials 并存**（显式 allowlist origin）；`.env` 确认 gitignore + `.env.example` 留空（无真值）。

### 7.2 USDT 6 位精度一致性

- 全链路 6 位最小单位（链/计价 *big.Int，DB string，展示除 10^decimals）。
- `usdtDecimals` 从 `deployments.json` 读，**不硬编码 6**（正式 USDT 同 6 位，但读字段防漂移）。
- seed `RequiredDeposit` + 费率表单价 + `MIN_DEPOSIT(10×10^6)` 全部 6 位语义统一。

### 7.3 RPC 统一 Arbitrum

- 三端（hardhat / 后端 / 前端）统一 Arbitrum Sepolia 421614；后端 RPC = json.rpcUrl 为准 + env 可覆盖，单一优先级。
- 杜绝 0G 残留（旧 `evm-testnet.0g.ai` / `evmrpc-testnet.0g.ai` 全清）。

### 7.4 联调依赖（前置项，明确标注；T1 前置拆分，闭合 CEO ⚠️）

**T1（abigen 绑定）的真前置 = 合约 PR#1 merge（ABI 冻结），不是「真机上链」**。两者必须拆开，避免基于未 merge 的 ABI 生成绑定导致返工：

| 前置事件 | 阻塞什么 | 当前状态 |
|----------|----------|----------|
| 合约 PR#1 merge（ABI 冻结） | **阻塞 T1 abigen 绑定**（ABI 必须是最终版才能生成稳定绑定） | handoff-backend.md 已声明 selector 冻结，但需确认 PR#1 已 merge 再开 T1 |
| 真机上链（421614 部署，arbitrum_sepolia.json 生成） | **仅阻塞端到端真机联调**，不阻塞 implement/test | ⏳ 待 DEPLOYER_PRIVATE_KEY/RPC（handoff §10） |

- **不阻塞 implement/test**：全程本地 hardhat(31337) + hardhat.json；abigen 绑定/计价/分批/事件解码/reorg/鉴权/金额硬闸均可本地 + **simulated.Backend** 单测验证。真机 421614 联调待上链后补跑。
- 上链后需复测：在 421614 校准分批上限 N（L2 calldata 计价差异，handoff §5.1）。

### 7.5 异常与监控/对账（检测层，与 §7.0 预防层互补）

> 本节是**事后检测兜底**，不是主防线；主防线见 §7.0 预防控制。

- 分批失败：记录失败批 users + 错误，幂等续跑（合约逐用户 continue + issueMonthly 幂等；批次幂等键 month+batchIndex）。
- 交易回执确认：发交易后 `WaitMined` 同步等回执判成败（§6.1）；event_sync 为事件落库对账事实源，不替代回执判定。
- 对账告警（检测）：`UsageDataSubmitted`/`BillCreated` 链上金额 ≠ 后端计价金额 → 告警日志（资损信号）。
- 熔断告警（检测→预防）：单批总额异常（> 月均 × N）→ 暂停 + 告警 + 人工放行（§7.0 ②）。
- event_sync last-block(含 blockHash) 持久化 + 断点续传 + reorg 回退 + 幂等去重（防漏/重复/回滚）。

---

## 8. 落地计划（给 plan 阶段输入）

| 任务 | 内容 | 文件 | 依赖 | 风险 |
|------|------|------|------|------|
| **T1** abigen 绑定 | make 脚本 + 生成 8 份绑定 + abiHash 校验 | `internal/blockchain/bindings/*`, Makefile | 合约 artifacts 存在 | ABI 与链上不符 → revert |
| **T2** config 修复 | 键名 bug + struct 扩展(Usdt/UsdtDecimals/AbiHash) + RPC 优先级 + .env.example | `config/config.go`, `configs/deployments.json`, `.env.example` | — | 键名不修 sync 静默失效 |
| **T3** client 读/写 | 读调用真实化 + MonthlySettlement/IssueMonthly 写 + **本地 nonce 计数器** + WaitMined + **金额硬闸 L3 断言**(B1) | `blockchain/client.go` | T1,T2 | selector/编码错 |
| **T4** owner 签名 | 注入内存私钥(不长存明文,B1/L9) + transactor + chainID 校验 + 降级关闭 + 日志脱敏 | `blockchain/client.go`, `cmd/main.go` | T3 | 私钥泄露/硬编码 |
| **T5** event_sync | FilterLogs 轮询 + 8 事件 abigen Filterer 解码 + 6 位精度 + 幂等 + **reorg 检测 + 两阶段确认 K 块**(B5) + last-block(含 blockHash) 持久化 + signatures.go 修正 | `sync/event_sync.go`, `signatures.go`, `models.go`(加列 Status/BlockHash/PayIntentTxHash) | T1,T2 | 漏/重复/回滚落库 |
| **T6** 计价+护栏 | PricingService 费率表(PLACEHOLDER) + **usage 上界**(B4) + **金额 L1/L2 硬闸 + 熔断**(B1) + FetchAndCreateBills 重写 + **operatorId 固定映射 + sanity check** + amounts[] + 分批≤25 + 批次幂等键 | `services/oracle.go`(+pricing), `cmd/main.go`(seed ID) | T3,T4 | 计价错算/映射错 |
| **T7** 鉴权+对账路径 | **AdminAuth + WalletAuth 双中间件**(B2/B3) + 端点按 §6.6 挂载 + **bills/pay/withdraw 降级 pending 不置终态**(B2/B3) + SubmitUsage max= 校验(B4) | `handlers/handlers.go`(+middleware), `services/services.go`(MarkAsPaid 不导出 HTTP), `cmd/main.go` | — | 端点裸奔/白嫖/伪造记账 |
| **T8** 测试 | 计价复算/上界拒绝/金额硬闸/分批切片+熔断/6位精度/abigen 事件解码/reorg 回退/鉴权中间件/对账单一路径 单测 + **simulated.Backend** 真实绑定跑分批结算+事件解码 | `services/*_test.go`, `blockchain/*_test.go`, `sync/*_test.go` | T1–T7 | 覆盖不足 |
| **T9** handoff-web | 产出 `docs/design/linkworld-backend/handoff-web.md`：exact approve / 对账 breaking change / deposit-withdraw 状态机（移交 web 3/3） | `docs/design/linkworld-backend/handoff-web.md` | T3–T7 | web 不知契约变更 |

**顺序**：T1→T2→(T3→T4)→T5→T6→T7→T8→T9。T5 与 T6 串行（都依赖 T1/T3）。T7 可与 T6 并行（不同文件）。**implement 阶段始终串行**（CLAUDE.md 规则）。

**风险总览**：① operatorId 固定映射 + sanity check（§4.5）；② owner=root key 安全 + 金额硬闸纵深防御（§6.2/7.0/7.1）；③ 计价正确性 + usage 上界（§4.2/7.0）；④ reorg/两阶段确认（§6.3/7.0）；⑤ 对账单一路径不被旁路（§4.3/6.6）；⑥ 联调依赖上链（§7.4）。

---

## 9. arch-review / security-review 重审清单（v2 已闭合，供复审核对）

**arch-review 复审（v2 闭合状态）**：
1. **operatorId 固定映射**（§4.5）：seed `ID=链上 operatorId` + 启动 sanity check，不靠 name 比对 → ✅ 已锁定。
2. **后端只代发 monthlySettlement**（§5.1）：deposit/payBill 属用户钱包侧，后端只读 + 事件回填 → ✅ 三方认可。
3. **对账单一路径**（§4.3）：删「保留不可信」二选一，IsPaid/withdraw 唯一事件回填，HTTP 仅 pending → ✅ 闭合 B2/B3，breaking change 移交 web（§12/handoff-web）。
4. **计价 + 上界 + 金额硬闸**（§4.2/§7.0）：纯整数 6 位 + usage 上界 + 三层金额硬闸 + 熔断 → ✅ 闭合 B4/B1。
5. **reorg + 两阶段确认**（§6.3/§7.0）：等 K 块 + 父哈希检测回退 → ✅ 闭合 B5。
6. **submitUsage 不实现链上签名**（§6.2）：v2-A 不参与计价 → 取舍保留。

**security-review 复审**：
1. **owner=root key 全链路**（注入内存不长存→transactor→chainID 校验→日志脱敏→降级）无泄露面；金额硬闸限制爆炸半径（§7.0/7.1）。
2. **双鉴权中间件**：AdminAuth 常量时间 + 缺 key fail；WalletAuth ecrecover 绑 wallet；覆盖 §6.6 全清单无遗漏。
3. **资损预防（§7.0 专节）**：金额硬闸 / usage 上界 / reorg 确认 / 单一对账路径——预防而非仅检测。
4. **分批幂等**：失败续跑不重复结算（合约幂等 + 后端去重 + month/batchIndex 幂等键）。
5. **RPC/链一致**：无 0G 残留、env/json 优先级单一、chainID 校验、占位零地址不订阅事件。
6. **其它**：crypto/rand password、CORS 收口、.env gitignore + example 留空。

---

## 10. 遗留 / 依赖

- ⏳ 端到端真机 421614 联调 + 分批上限校准 → **依赖合约 1/3 真·上链**（handoff §10），不阻塞 implement/test（本地 31337）。
- ⏳ KMS 私钥托管、HMAC 端点签名、真实运营商 usage API、NotificationService SMTP、动态风控引擎 → 留后续 Round。
- ⏳ 费率表真实费率值 + 金额硬闸/usage 上界/熔断倍数/K 块确认数等 PLACEHOLDER 常量 → 产品/运营/安全填，本轮占位 + 启动 warn + 测试固定断言。
- ⚠️ **验收语言澄清（防伪胜利）**：本轮验收 = **计价引擎正确性 + 链集成 + 资损预防控制**；**usage 仍是 `OperatorAPISimulator` 模拟（rand）**，真实计量是独立后续 Round。「计价引擎正确 ≠ 计费正确」——§5B 验收**不得读成「计费痛点已解决」**。
- ⏳ operatorId seed 顺序需与合约 1/3 `ServiceManager.initialize` 注入顺序最终核对（PR#1 merge 后）。

---

## 11. arch-review 阻塞闭合对照表（B1–B5）

| # | 原阻塞（arch-review §二） | v2 闭合方案 | 落点 |
|---|--------------------------|-------------|------|
| **B1** 🔴 | owner 单 EOA「凭空造单 + 扣款」复合权限无金额护栏 | ① **三层金额硬闸**：L1 计价层（≤MAX_BILL_PER_USER）/L2 组批层（单批总额上限 + 超月均×N 熔断人工放行）/L3 client 发交易前断言；② owner key **注入内存不落 .env 明文长存**（L9）；③ §7.1 措辞升级 **owner=平台 root 权限**（改分账地址/授权拓扑）；④ web 侧 **exact-amount approve 禁 infinite** 写入 handoff-web | §0.1(L4/L5/L9)、§4.1、§6.1、§7.0①、§7.1、§12/handoff-web |
| **B2** 🔴 | `/api/bills/pay`→直接置 IsPaid 不验链上 + §4.3「保留不可信」二选一 = 白嫖 | ① **删除「保留」选项**；② **IsPaid 唯一由 event_sync 监听 BillPaid 回填**；③ `/api/bills/pay` 降级为 **pending 意向**（写 PayIntentTxHash 不动 IsPaid），`MarkAsPaid` 不导出给 HTTP；④ §6.6 鉴权清单补 /api/bills/pay（WalletAuth）；⑤ breaking change 声明移交 web | §4.3、§4.4、§6.6、§7.0④、§12/handoff-web |
| **B3** | `/api/withdraw`→凭前端 txHash 写记账 + 错归 AdminAuth | ① withdraw 记账**唯一由 DepositWithdrawn 事件回填**；② HTTP **不接受前端 txHash 记账**（最多 pending，不计入余额）；③ 鉴权改 **WalletAuth（钱包签名 ecrecover）非 AdminAuth**；④ §6.6 把 withdraw 移出 AdminAuth | §4.3、§6.3、§6.6、§7.0④ |
| **B4** | usage 无上界 → 天文账单；SubmitUsage 无范围校验 | ① **PricingService 对 (dataMB,callMin) 设上界**超界拒绝+告警；② **amount6 单 bill 硬上限**（同 B1 L1）；③ **SubmitUsage gin binding `max=`** 校验 | §4.1、§4.2、§4.4、§6.6、§7.0② |
| **B5** | event_sync 无 reorg/确认数 → 已付/押金虚高 | ① 资金事件落 **confirmed 前等 K 块确认**；② **reorg 检测**（last-block 记 blockHash，父哈希断→回退重扫未确认）；③ 资金事件 **pending→seen→confirmed 两阶段**，与 §4.3 状态机统一 | §4.3、§4.4、§6.3、§7.0③ |

**⚠️ 收紧项闭合**（arch-review §三）：operatorId 固定映射+sanity check（§4.5）/ abigen Filterer 解码 + UsageDataSubmitted 只 user indexed + BillCreated 含费总额（§6.3）/ nonce 本地计数器+WaitMined 统一（§6.1）/ simulated.Backend 测试（§7.4/T8）/ T1 前置拆「ABI 冻结 vs 真机上链」（§7.4）/ 验收语言澄清（§10）/ 费率表 PLACEHOLDER 刺眼+owner+真机 gate（§4.2）/ usage 单位锁 MB/min（§4.2）/ seed RequiredDeposit 整数最小单位+迁移（§4.4）/ string→big.Int 校验 ok fail-fast（§6.3）/ chainID 一致校验（§6.2）/ 占位零地址不订阅事件（§6.3）/ 批次幂等键 month+batchIndex（§4.1/7.5）/ password crypto/rand + CORS 禁 *+credentials + .env gitignore+example 留空（§7.1）。

---

## 12. 跨子项 handoff-web 要点（移交 web 3/3）

> 详见独立文件 `docs/design/linkworld-backend/handoff-web.md`。本节为摘要，后端 design 锁定、web 子项目据此对齐。

1. **exact-amount approve 禁 infinite（B1 跨子项）**：web 在 `deposit`/`payBill` 前的 `usdt.approve` **必须按本次精确金额授权**（`approve(spender, amount/amount+fee)`），**禁止 infinite/max approve**。infinite approve 放大 owner=root 被盗用时的资金抽取面。
2. **对账契约 breaking change（B2/B3）**：
   - `POST /api/bills/pay`：语义从「标记已付」变「记录支付意向」。**后端不再据此置 IsPaid**；web 必须等链上 `BillPaid` 被后端 event_sync 确认后（轮询 GET bills 看 is_paid，或后端推送）才显示「已付」。乐观 UI 可标「支付确认中」。
   - `POST /api/withdraw`：不再接受 txHash 做记账依据；后端凭 `DepositWithdrawn` 事件回填。web 流程不变（仍先链上 withdraw），但提现历史以后端确认为准。
   - 两端点新增**钱包签名鉴权头**（WalletAuth）：web 需对请求体 + nonce 用钱包签名。
3. **deposit/withdraw 状态机（pending→confirmed）**：web 显示余额/历史时区分 pending（意向，未上链确认）与 confirmed（事件确认）；余额计算只认 confirmed。
4. **金额精度**：全链路 USDT 6 位最小单位（big.Int / 字符串），web 展示除 10^usdtDecimals（从 deployments 读，勿硬编码 6）。
