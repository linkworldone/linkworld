# Stage: plan — 任务拆解（子项目 backend 2/3）

> **状态**: completed（用户已审批） | **日期**: 2026-06-09 | **Gate**: 2 | **子项目**: backend(2/3)
> 产出：plan-review.md（人读，4 节含表格）+ plan-review.json（hook 读，8 pages 各含 components+colors，已过 advance 校验）
> 领域适配：plan 模板为 web/UI 形状，后端为 Go——任务=page、components=Go 文件、colors=[]。

## 8 个串行任务（implement 铁律：串行，一个完成审查再下一个）
| 任务 | 范围 | gate |
|------|------|------|
| T1 abigen | 8 份合约绑定生成入库 + abiHash 校验工具（前置=合约 PR#1 merge ABI 冻结；本分支栈式 artifacts 可用） | go build 绿 |
| T2 config+链配置 | config.go 键名 bug(contracts→proxies)修 + deployments.json 421614 schema + RPC/chainID 统一 + .env.example | 配置加载正确 |
| T3 client 读/写 | 读真实化 + MonthlySettlement/IssueMonthly owner 发交易 + 本地 nonce 计数器 + WaitMined + 金额闸 L3 + key 内存注入 + chainID 校验 + 降级 | simulated.Backend 测 |
| T4 event_sync | FilterLogs 轮询 + 8 事件 abigen Filterer 解码 + 6 位精度 + (txHash,logIndex)去重 + reorg+K 块确认 + pending→seen→confirmed 两阶段 + signatures.go 修正 + models 加列 | 事件解码测 |
| T5 计价+operatorId | PricingService 费率表(big.Int 6 位)+usage 上界 L1+单 bill 上限+单位锁 MB/min；operatorId 固定映射(seed ID=链上 1..11)+启动 sanity check | 计价复算测 |
| T6 组批/熔断+幂等 | 组批 ≤25 + 金额闸 L2 + 熔断(冷启动回退 MAX_BATCH_TOTAL 绝对闸) + month+batchIndex 幂等键 | 组批测 |
| T7 鉴权 | AdminAuth(ConstantTimeCompare,缺 key fail) + WalletAuth(EIP-712 ecrecover + 服务端一次性 nonce 台账,禁纯 timestamp) + 端点挂载 + bills/pay 降级 pending(PayIntentTxHash 解耦)+MarkAsPaid 收口 + withdraw 事件回填 + CORS 收口 | 鉴权/防重放测 |
| T8 测试 | simulated.Backend 真实绑定端到端(分批结算/事件解码/operatorId 断言) + 计价/上界/精度/鉴权/nonce/reorg/幂等单测 | go test ./... 覆盖 |

## arch-review 红线映射
N1 WalletAuth 一次性 nonce 台账(禁纯 timestamp)→T7；N2/L2 熔断冷启动回退绝对闸→T6；operatorId 固定映射→T5；金额三层闸 L1→T5/L2→T6/L3→T3；对账单一路径(IsPaid 唯一事件回填)→T4+T7。

## 待拍真值/遗留（不阻塞 plan）
- 占位常量 MAX_BILL_PER_USER/MAX_BATCH_TOTAL/熔断 N/确认数 K 待产品/运营/安全拍真值（design 已刺眼 PLACEHOLDER+启动 warn+测试固定断言）。
- T1 前置=合约 PR#1 merge（ABI 真冻结）；本分支栈式于合约分支，artifacts 可用可先行，但 PR#1 review 若改 ABI 需 regen 绑定。
- operatorId seed 顺序待与合约 ServiceManager.initialize 最终核对。
- 真机 421614 联调待合约真·上链（handoff §10）；implement/test 用本地 hardhat.json + simulated.Backend。

## 移交 implement
从 T1 严格串行 + TDD（Go 用 simulated.Backend）。每任务完成→主 Agent 审查+写 checkpoint→再派下一个。
