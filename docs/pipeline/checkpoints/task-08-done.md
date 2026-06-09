# Task 08 — T8 测试扫尾（端到端鉴权 + 三层闸联动 + 对账 + 覆盖率）（后端 2/3）

> 子项目 backend(2/3) | 状态 DONE | 2026-06-09 | 最后一棒

### 产出文件（仅新增测试，未改任何业务 .go）
- packages/backend/internal/handlers/e2e_auth_test.go（新：端到端鉴权流——WalletAuth 中间件 + 受保护端点 HTTP 全链路 4 用例）
- packages/backend/internal/services/gate_interlock_test.go（新：L1/L2/L3 三层金额闸联动 + 不漂移 4 用例）
- packages/backend/internal/sync/reconcile_test.go（新：对账金额不一致告警路径 3 用例）
- packages/backend/internal/handlers/trigger_bill_test.go（新：SettlementSummary=null 降级序列化 2 用例）
- packages/backend/internal/repository/repository_crud_test.go（新：余额 confirmed-only + SetPayIntent 意向 3 用例）

### git commit
（见 commit hash）test: 后端 T8 测试扫尾（端到端鉴权+三层闸联动+对账+覆盖率）

### 覆盖审计表（design §8 逐项）
| design §8 / arch-review 项 | 状态 | 落点 |
|---|---|---|
| 计价复算/usage 上界/L1 金额闸/纯整数 | 已覆盖 | pricing_test (PRICE-01..05) |
| 金额硬闸 L3 (client 发交易前) | 已覆盖 | client_test (L3Gate*) |
| 金额硬闸 L2 + 熔断 + 冷启动回退 + 分批≤25 + 幂等 + 失败续跑 | 已覆盖 | settlement_test (BATCH/CB/IDEM/FAIL) |
| 6 位精度 abigen 事件解码 (BillCreated 含费/Usage 仅user/Withdrawn) | 已覆盖 | event_sync_test (Parse*) |
| reorg 回退 / 两阶段确认 / 去重 / 占位跳过 | 已覆盖 | event_sync_test (SyncOnce*) |
| 对账单一路径 (IsPaid 唯一事件回填) | 已覆盖 | event_sync_test (BillPaid_OnlyEvent) |
| operatorId 固定映射 + sanity check | 已覆盖 | operatorid_test (OPID-01..03) |
| AdminAuth / WalletAuth nonce 重放 / chainId / wallet 不符 | 已覆盖 | middleware_test (ADMIN/WALLET) |
| bills/pay 降级 / withdraw pending / usage max | 已覆盖（handler 直挂） | handlers_t7_test |
| abiHash 校验 | 已覆盖 | abihash_test |
| **端到端鉴权流（中间件+受保护端点 HTTP 全链路）** | **新补** | e2e_auth_test (E2E-AUTH-01..04) |
| **三层闸 L1/L2/L3 联动（金额闸不漂移）** | **新补** | gate_interlock_test (LINK-01..04) |
| **对账金额不一致告警路径 (processUsageDataSubmitted)** | **新补** | reconcile_test (RECONCILE-01..03) |
| **SettlementSummary=null 降级序列化** | **新补** | trigger_bill_test (TRIGGER-01/02) |
| **repository CRUD 关键路径 (余额 confirmed-only / SetPayIntent)** | **新补** | repository_crud_test (REPO-01..02) |
| 业务合约端到端真机结算（金额数组语义/事件回填上链） | **N-A（遗留）** | 见下「simulated.Backend 决策」 |

### go test -cover 最终统计
75 个 PASS 用例，0 失败。go build ./... 0 error，go vet ./... 0 warning。
新增测试涉及包 -race 全绿（handlers/services/sync/repository）。

各包覆盖率：
- internal/blockchain        72.9%
- internal/middleware        73.7%
- internal/sync              65.6%
- internal/config            48.0%
- internal/services          39.8%
- internal/handlers          18.3%
- internal/repository         6.2%
- bindings/abis/models/cmd/genbindings  无测试（生成产物/装配/纯模型）

> 覆盖率口径说明：services/handlers/repository 数值偏低主因——大量「展示查询/装配/SMTP stub」非资损路径未单测（本轮范围 out）；
> 资损敏感路径（计价 L1/组批 L2/client L3/event_sync reorg+两阶段/鉴权 nonce 重放/对账告警）覆盖率高且红线用例齐全。
> blockchain 72.9% 已含 simulated.Backend 能跑的 client 机制（nonce/L3/降级/chainID）；剩余未覆盖为业务合约端到端（受 Cancun 限制，遗留）。

### simulated.Backend / Cancun 决策结论（T3 遗留闭合）
- **结论（本轮策略 ①②，不强行升级/不强接 hardhat）**：
  ① 链交互逻辑用接口 mock（SettlementClient/SettlementBatchStore）+ 手工构造 types.Log（abigen Filterer.Parse* 解码）覆盖——T3/T4 已做，T8 复用并补联动/对账缺口。
  ② go-ethereum v1.13.5 simulated.Backend 是 London 上限（ethash faker），跑不了项目 Cancun 字节码（transient storage/PUSH0）。client 机制层（L3/nonce/降级/chainID）已用 London 兼容桩字节码 + 真实绑定验证。
- **明确记录为 test 阶段/后续遗留**：**业务合约端到端真机联调 = 待合约 1/3 真·上链后，用 hardhat 31337 或升级 go-ethereum v1.16+（ethclient/simulated PoS 后端支持 Cancun）跑分批结算金额数组语义 + 事件回填上链**。不阻塞本轮 implement/test（design §7.4 已声明此依赖不阻塞）。

### 发现的 bug
无。审计 + 新增 16 用例全部一次性或经合法测试调整后通过，未触及任何业务逻辑 bug（严格遵守「只补测试不改业务」）。

### 偏差记录
无。仅新增 5 个 *_test.go 文件，未碰任何业务 .go / web / contracts / pipeline.json / cmd 装配。

### 遗留问题（带入 test 阶段 / 后续 Round）
- 业务合约端到端真机结算（simulated.Backend/Cancun 限制）→ 待合约上链 + hardhat 31337 / geth v1.16+。
- WalletNonce 过期清理 job（不影响正确性，台账会增长）→ 后续可加。
- services/handlers/repository 非资损展示路径覆盖率可后续 Round 补（本轮聚焦资损红线）。

### 可否进 test 阶段的结论
**可进 test 阶段。** 资损预防控制（三层金额闸联动不漂移 / usage 上界 / reorg 两阶段确认 / 对账单一路径 + 不一致告警）+ 端到端鉴权流（WalletAuth 中间件→受保护端点→降级落库）+ repository 记账契约（余额 confirmed-only / 意向不置终态）均有红线用例覆盖且全绿。唯一遗留（业务合约真机端到端）为环境/联调依赖，design §7.4 已明确不阻塞 implement/test。
