# Stage: plan — 后端落地计划任务拆解（子项目 backend 2/3）

> **状态**: v1 | **日期**: 2026-06-09 | **Gate**: 2（计划） | **子项目**: backend(2/3) | **角色**: 工程师/计划 | **分支**: backend/align-arbitrum-usdt
> **输入**：design.md v2（§8 落地计划 T1–T8 + §7.0 资损预防控制 + §11 闭合表 + §12 handoff-web）+ arch-review.md（PASS，§七 2 条 implement 验收红线 N1/L2 + 占位常量/operatorId 遗留）+ handoff-backend.md（冻结 ABI/selector + 6 位精度 + 分批 N≤25/≤50）+ 2 份后端基线（alignment-surface / blockchain-integration）+ 真实源码核对。
> **后端根目录**：`packages/backend/`（module `linkworld-backend`，go-ethereum v1.13.5）。本文所有 Go 路径以此为根。
> **注**：本 stage 文件按子项目串行复用，合约 1/3 的 plan 前版已入库（见 git history）。
>
> ⚠️ 本文档只产计划，不写 .go 代码。implement 阶段按本表逐任务串行执行，每任务完成审查后才派下一个（CLAUDE.md 主控-分派 + implement 始终串行规则）。
>
> **铁律前置**：T1 前置 = **合约 PR#1 merge（ABI 冻结）**，非真机上链；真机 421614 联调待上链，不阻塞 implement/test（design §7.4）。

---

## 1. 任务拆解（串行 T1→T8）

> 串行依据：design §3.1 依赖图 + §8 顺序 `T1→T2→(T3→T4)→T5→T6→T7→T8`。design 原 T3/T4 合并为本表 **T3**（同文件 client.go，避免拆同文件并行冲突）；design 原 T9 handoff-web 已在 design 阶段产出（§12），本轮不重复，故本表收敛为 8 个串行任务 T1–T8。**implement 阶段始终串行**：一个任务 gate 绿且审查通过，才派下一个。

| 任务ID | 范围 | 受影响 Go 文件 | 依赖 | gate（go build/test） |
|--------|------|----------------|------|------------------------|
| **T1** abigen 绑定 | make 脚本 + 生成 8 份 abigen 绑定（7 业务 + MockUSDT）+ abiHash 校验工具 | `Makefile`（新增 `abigen` target）、`internal/blockchain/bindings/<name>.go`（生成入库）、`internal/blockchain/abihash.go`（指纹比对） | **合约 PR#1 merge（ABI 冻结）** | `go build ./...` 绿；生成绑定可编译；`abiHash` 比对工具单测过 |
| **T2** config + 链配置 | config.go 键名 bug 修（`contracts`→`proxies`）+ Deployments struct 扩展（Usdt/UsdtDecimals/AbiHash）+ deployments.json 改 421614 schema（7 占位地址 + usdt + usdtDecimals:6 + abiHash）+ RPC 优先级统一（json 为准 env 覆盖）+ chainID 一致校验 + `.env.example` 同步 + 同步脚本 | `internal/config/config.go`、`configs/deployments.json`、`.env.example`、`scripts/sync-deployments.sh`（新增） | T1 | `go build ./...` 绿；config 加载单测过；占位零地址不进 event 过滤断言过 |
| **T3** client 读/写 + owner 签名 | 读调用真实化（GetDepositAmount/GetLockExpiry/VerifyServiceActive/GetFeeRate）+ 写调用 MonthlySettlement/IssueMonthlyTrafficCards + 本地 nonce 计数器（mutex）+ WaitMined 等回执 + 金额闸 L3（发交易前断言 amount∈(0,MAX_BILL_PER_USER]、sum≤MAX_BATCH_TOTAL）+ owner key 注入内存（不长存明文）+ transactor + chainID 校验 + 缺 key 降级关闭 + 日志脱敏 | `internal/blockchain/client.go`、`cmd/main.go`（装配 transactor） | T1, T2 | `go build ./...` 绿；simulated.Backend 写交易 + L3 拒发单测过 |
| **T4** event_sync | FilterLogs 轮询（区块分页 + last-block 持久化）+ 8 事件 abigen Filterer 解码 + 6 位精度（SetString 校验 ok fail-fast）+ (txHash,logIndex) 去重 + reorg 检测（last-block 记 blockHash 父哈希断→回退重扫）+ pending→seen→confirmed 两阶段（资金事件等 K 块）+ signatures.go 修正（BillCreated 4 参 + 新增 3 topic）+ models 加列（Status/BlockHash/PayIntentTxHash/OnChainBillID）| `internal/sync/event_sync.go`、`internal/blockchain/signatures.go`、`internal/models/models.go` | T1, T2 | `go build ./...` 绿；simulated.Backend 事件解码 + reorg 回退 + 两阶段确认单测过 |
| **T5** 计价 + operatorId | PricingService 费率表（OperatorRate big.Int 6 位 PLACEHOLDER + 启动 warn）+ usage 上界 L1（Price 入口 MAX_DATA_MB/MAX_CALL_MIN 拒绝）+ amount6 单 bill 上限（L1）+ 单位锁 MB/min + operatorId 固定映射（seed ID=链上 1..11）+ 启动 sanity check（读链校验 fail/warn）| `internal/services/oracle.go`（+`pricing.go`）、`cmd/main.go`（seed ID） | T3 | `go build ./...` 绿；计价复算 + 上界拒绝 + 精度 + operatorId 映射断言单测过 |
| **T6** 组批/熔断 + 幂等 | FetchAndCreateBills 重写：组批切片每批 ≤25 + 金额闸 L2（单批 sum≤MAX_BATCH_TOTAL + 超月均×N 熔断 + 冷启动回退绝对闸）+ 逐批发交易（依赖 T3 client）+ 结算批次幂等键（month+batchIndex 落 DB 已确认不重发）+ TriggerMonthlyBill handler 返回分批摘要 | `internal/services/oracle.go`、`internal/models/models.go`（批次幂等表）、`internal/handlers/handlers.go`（摘要返回） | T5 | `go build ./...` 绿；分批切片 + 熔断（含冷启动回退）+ 幂等不重发单测过 |
| **T7** 鉴权 | AdminAuth 中间件（subtle.ConstantTimeCompare，缺 key 启动 fail）+ WalletAuth 中间件（EIP-712 ecrecover 绑 wallet + 服务端一次性 nonce 台账，禁纯 timestamp）+ 端点按 §6.6 挂载（bills/pay 降级 pending+PayIntentTxHash 解耦 IsPaid；withdraw 事件回填+WalletAuth；usage/submit gin `max=`）+ MarkAsPaid 不导出 HTTP + CORS 禁 `*`+credentials（allowlist）+ password crypto/rand | `internal/handlers/handlers.go`（+`middleware.go`）、`internal/services/services.go`（MarkAsPaid 收口 + RecordWithdraw 废弃）、`cmd/main.go`（挂中间件 + CORS） | T4（依赖 PayIntentTxHash/Status 列） | `go build ./...` 绿；AdminAuth/WalletAuth/nonce 防重放/降级 pending 单测过 |
| **T8** 测试 | simulated.Backend 部署真实绑定测分批结算/事件解码/operatorId 映射断言 + 计价复算/上界/精度单测 + 鉴权中间件测 + nonce 防重放测 + reorg 回退测 + 幂等不重发测；`go test ./...` 覆盖收口 | `internal/services/*_test.go`、`internal/blockchain/*_test.go`、`internal/sync/*_test.go`、`internal/handlers/*_test.go` | T1–T7 | `go test ./...` 全绿；simulated.Backend 端到端断言过 |

---

## 2. 受影响模块与改动

| 模块（文件） | 现状（源码核对） | 改动 | 来源（design §） | 风险 |
|--------------|------------------|------|------------------|------|
| `internal/blockchain/bindings/`（新建） | 不存在；abis/ 仅 2 份手写裁剪 ABI | abigen 生成 8 份类型安全绑定（Caller/Transactor/Filterer）入库 + abiHash 比对 | §5.3 / §3.2 | ABI 与链上不符 → revert；须 PR#1 merge 后生成防返工 |
| `internal/blockchain/client.go` | 业务方法全 stub 返零值/not implemented，无 ABI 绑定、无写交易、无 transactor | 读调用真实化 + MonthlySettlement/IssueMonthly 写 + nonce 计数器 + WaitMined + L3 金额断言 + transactor 装配 + chainID 校验 + 降级 | §6.1 / §6.2 / §7.0① | selector/编码错；nonce 复用；私钥泄露 |
| `internal/blockchain/signatures.go` | BillCreated 写死 5 参（错）；缺 TrafficCardApplied/UsageDataSubmitted/DepositWithdrawn | BillCreated 改 4 参 + 新增 3 topic；仅日志/轻量过滤用，字段解码走 abigen Filterer | §5.2 / §6.3 | topic 错匹配不到事件 |
| `internal/sync/event_sync.go` | 主循环只 sleep 30s 空转，从未调 FilterLogs，process* 永不触发；processUserRegistered 落 Unix(0,0) | FilterLogs 轮询 + 8 事件 Filterer 解码 + 6 位精度 + 去重 + reorg + 两阶段确认 + last-block(blockHash) 持久化 + process 修正 | §6.3 / §4.3 / §7.0③ | 漏/重复/回滚落库 → 已付/押金虚高资损 |
| `internal/services/oracle.go` | GetBill 返回 rand.Intn 随机数；FetchAndCreateBills 只写 DB 不上链无分批；SignData 是 SHA256 非 ECDSA | 新增 PricingService 费率表 + 上界 L1 + FetchAndCreateBills 重写（组批 ≤25 + L2 熔断 + 逐批上链 + 幂等键）+ operatorId 固定映射；废弃 GetBill 随机金额、SignData SHA256 | §4.1 / §4.2 / §4.5 / §6.4 / §7.0① | 计价错算/映射错 = 分账打错地址（最隐蔽资损） |
| `internal/config/config.go` | struct tag `proxies` 但 JSON 键 `contracts` → Proxies 永远空 map（真 bug）；无 Usdt/UsdtDecimals/AbiHash | 修键名 + struct 扩展 + abiHash 校验 + RPC 优先级统一（json 为准 env 覆盖）+ chainID 一致校验 | §6.5 / §7.3 | 键名不修 sync 静默失效；连 A 链发 B 链 |
| `configs/deployments.json` | chainId 16602(0G) + 0G 旧 RPC + 键名 contracts + 7 个 0G 旧地址；缺 usdt/usdtDecimals/abiHash | 改 421614 schema（proxies 7 占位 + usdt 占位 + usdtDecimals:6 + abiHash 占位）；占位零地址不订阅事件；同步脚本从 contracts/deployments/<net>.json 回填 | §6.7 / §7.3 | 占位零地址误匹配事件；连错网 |
| `internal/handlers/handlers.go`（+ middleware.go 新建） | 18 端点全裸奔（仅 CORS）；CORS 可能 *+credentials | AdminAuth + WalletAuth 双中间件 + 端点按 §6.6 挂载 + bills/pay 降级 pending + usage/submit `max=` + CORS 收口 | §6.6 / §7.1 / §7.0④ | 端点裸奔/白嫖/伪造记账 |
| `internal/services/services.go` | MarkAsPaid 可由 HTTP 置 IsPaid（白嫖面）；RecordWithdraw 凭前端 txHash 记账 | MarkAsPaid 不导出给 HTTP（仅 event_sync 内部调）；RecordWithdraw 废弃（withdraw 唯一 DepositWithdrawn 事件回填）；余额计算只算 confirmed | §4.3 / §6.6 / §7.0④ | 旁路写终态 = 资损 |
| `internal/models/models.go` | Bill/Deposit 金额全 string；无 Status/BlockHash/PayIntentTxHash/OnChainBillID | Bill 加 OnChainBillID/PayIntentTxHash（与 IsPaid 解耦）；Deposit 加 Status/BlockHash；UsageData 加 max 校验；批次幂等表 | §4.4 | AutoMigrate 加列（GORM 兼容） |
| `cmd/main.go` | 11 运营商 seed RequiredDeposit="0.01" 语义不明；仅 RPC_URL 非空起 sync 用空 map；无 owner key 装配；无鉴权中间件 | seed 显式 ID=1..11 + RequiredDeposit 改 6 位整数 + 存量迁移；装配 transactor + sanity check；挂鉴权中间件 + CORS 收口 | §4.5 / §6.2 / §6.6 | operatorId 映射错；私钥硬编码 |

---

## 3. 验收标准映射

> 验收点来源：requirement.md §五B（计价引擎正确性 — 注意 design §10 防伪胜利：usage 仍模拟，本轮验收 = 计价引擎正确性 + 链集成 + 资损预防，**非「计费痛点已解决」**）+ arch-review §七 2 条 implement 验收红线（N1 WalletAuth 防重放 / L2 冷启动熔断回退）+ design §7.0 资损预防控制四闸。

| 验收点 | 对应任务 | 验收方式 |
|--------|----------|----------|
| abigen 8 份绑定可编译 + abiHash 与 deployments.json 一致（不一致 fail） | T1 | `go build ./...` + abiHash 比对单测 |
| config 键名 bug 修复，Proxies 非空 map；占位零地址不进 event 过滤 | T2 | config 加载单测 + 零地址跳过断言 |
| RPC/chainID 单一优先级（json 为准 env 覆盖）+ transactor chainID == 链上 chainID（不一致 fail） | T2, T3 | chainID 一致校验单测 |
| owner key 注入内存不长存明文；缺 key 写降级关闭只读仍可跑；日志不打私钥 | T3 | 降级路径单测 + 日志脱敏人工核对 + .env.example 留空 |
| 金额硬闸 L3：发交易前 amount∈(0,MAX_BILL_PER_USER]、sum≤MAX_BATCH_TOTAL，越界拒发 | T3 | simulated.Backend L3 拒发单测（固定占位常量断言） |
| MonthlySettlement 本地 nonce 计数器（发一笔+1）+ WaitMined 等回执判成败 | T3 | simulated.Backend 连发多批 nonce 不复用 + 回执判定单测 |
| 8 事件 abigen Filterer 解码；UsageDataSubmitted 只 user indexed；BillCreated 第三参含费总额 | T4 | simulated.Backend 事件解码字段断言单测 |
| reorg 检测（父哈希断回退重扫未确认）+ 资金事件等 K 块 confirmed 才置终态 | T4 | simulated.Backend reorg 回退 + 两阶段确认单测 |
| (txHash,logIndex) 去重；last-block(blockHash) 持久化断点续传；string→big.Int 校验 ok fail-fast | T4 | 重叠区块去重单测 + 非法金额 fail-fast 单测 |
| 计价纯整数 big.Int 6 位可复算（费率 × usage）；废弃 rand 随机金额 | T5 | 计价复算固定断言单测 |
| usage 上界 L1：Price 入口 dataMB>MAX_DATA_MB / callMin>MAX_CALL_MIN 拒绝 + 告警 | T5 | 上界拒绝单测 |
| operatorId 固定映射（seed ID=链上 1..11）+ 启动 sanity check 不一致 fail/warn | T5 | operatorId 映射断言 + sanity check 单测 |
| 组批切片每批 users≤25（handoff §5.1 最紧约束） | T6 | 分批切片单测（边界 25/26） |
| 金额闸 L2：单批 sum≤MAX_BATCH_TOTAL + 超月均×N 熔断；**冷启动/无样本回退绝对闸**（arch-review N2/L2 红线） | T6 | 熔断单测 + 月均=0 冷启动回退 MAX_BATCH_TOTAL 单测 |
| 结算批次幂等键（month+batchIndex 落 DB 已确认不重发） | T6 | 重复触发不重发单测 |
| AdminAuth 常量时间比较（subtle.ConstantTimeCompare）+ 缺 key 启动 fail | T7 | 中间件单测 + 缺 key fail 路径 |
| **WalletAuth 防重放（arch-review N1 最关键红线）**：EIP-712 ecrecover 绑 wallet + 服务端一次性 nonce 台账，**禁纯 timestamp** | T7 | 同一签名第二次提交被拒（nonce 已消费）单测 + ecrecover 地址绑定单测 |
| bills/pay 降级 pending（写 PayIntentTxHash 不动 IsPaid）；withdraw 不收 txHash 记账；MarkAsPaid 不导出 HTTP | T7 | 端点不置终态单测 + MarkAsPaid 仅内部可调 |
| usage/submit gin binding `max=` 上界校验；CORS 禁 *+credentials（allowlist）；password crypto/rand | T7 | binding 拒绝超界单测 + CORS 配置核对 |
| simulated.Backend 部署真实绑定跑分批结算 + 事件解码端到端 | T8 | `go test ./...` 全绿 |

---

## 4. 实现红线与铁律

> 以下为 implement 阶段不可逾越的铁律，违反任一即返工。来源：design §7.0/§11 + arch-review §七 + CLAUDE.md 主控-分派规则。

1. **串行执行**：implement 阶段始终串行，一个任务 gate（go build/test）绿且主 Agent 审查通过，才派下一个 subagent（CLAUDE.md）。两任务编辑同一文件必串行（T3/T5/T6 同碰 oracle.go/client.go，严格按 T1→…→T8 顺序）。
2. **T1 前置 = 合约 PR#1 merge（ABI 冻结）**，非真机上链。基于未 merge ABI 生成绑定会返工（design §7.4）。真机 421614 联调待上链，**不阻塞 implement/test**（全程本地 hardhat 31337 + hardhat.json + simulated.Backend）。
3. **WalletAuth nonce 台账（arch-review N1，最关键）**：必须服务端一次性 nonce 台账（per-wallet 单调递增或消费式）+ 绑 chainId/domain（EIP-712），**禁止退化为纯 timestamp 时间窗**（窗口内签名可重放）。design §6.6 措辞收紧为「nonce 强制 + 服务端状态」。
4. **熔断冷启动回退（arch-review N2/L2）**：历史月均熔断在首月/无样本（均值=0）时须回退到 `MAX_BATCH_TOTAL` 绝对闸，**不得因均值=0 失效**。L1/L3 绝对上限始终在。
5. **operatorId 固定映射**：seed 显式写 `Operator.ID = 链上 operatorId(1..11)`，**绝不靠 name 比对**；启动读链 sanity check 不一致 fail/warn。seed 顺序待 PR#1 merge 后与 `ServiceManager.initialize` 核对（design §4.5）。
6. **对账单一路径**：资金终态字段（IsPaid/PaidAt、withdraw 记账）**唯一**由 event_sync 链上事件回填；HTTP 写端点最多写 pending 意向，**绝不置终态**。MarkAsPaid 不导出给 HTTP，RecordWithdraw（凭前端 txHash）废弃（design §4.3/§7.0④）。
7. **金额三层硬闸**：L1 计价层（≤MAX_BILL_PER_USER）/ L2 组批层（单批 sum≤MAX_BATCH_TOTAL + 超月均×N 熔断）/ L3 client 发交易前断言。三闸独立，即使上游漏校验 L3 也挡（design §7.0①）。
8. **owner = 平台 root 权限**：私钥启动注入进程内存（本地/CI env、生产 secret manager），不落 .env 明文长存、不硬编码、不入库、不进日志（仅打 owner address）；缺失则写降级关闭；chainID 一致校验防签错链（design §6.2/§7.1）。
9. **占位常量刺眼化**：MAX_BILL_PER_USER/MAX_BATCH_TOTAL/N/K + 费率表全部 PLACEHOLDER 注释 + 启动 `log.Warn`，测试用固定值断言；真值待产品/运营/安全在 implement 前确认（design §0.1 注 / §10）。
10. **精度 6 位全链路**：链/计价 *big.Int 6 位最小单位，DB string，展示除 10^usdtDecimals（从 deployments 读不硬编码 6）；string→big.Int `SetString(s,10)` 必须校验 ok，失败 fail-fast 不静默（design §6.3/§7.2）。
