# Stage test — checkpoint（后端子项目 2/3）

> 状态 **DONE（全绿，可进 review）** | 日期 2026-06-09 | 分支 backend/align-arbitrum-usdt
> 角色：test runner（只跑测试/验证 + 产文档，未改任何业务 .go）

## Phase 1 — 静态门（packages/backend）
- `go build ./...` → **0 error**（BUILD_EXIT=0）
- `go vet ./...` → **0 warning**（VET_EXIT=0）
- go 工具链：go1.26.1 darwin/arm64；模块 `linkworld-backend`；`github.com/ethereum/go-ethereum v1.13.5`

## Phase 2 — go test ./... -cover
- **75 个用例全 PASS，0 FAIL**（`-count=1 -v` 统计：`--- PASS` 75 / `--- FAIL` 0）
- `-race -count=1` 全包通过（blockchain/config/handlers/middleware/repository/services/sync 均 ok）

各包覆盖率：

| 包 | 覆盖率 |
|---|---|
| internal/blockchain | 72.9% |
| internal/middleware | 73.7% |
| internal/sync | 65.6% |
| internal/config | 48.0% |
| internal/services | 39.8% |
| internal/handlers | 18.3% |
| internal/repository | 6.2% |
| cmd / bindings / abis / models / genbindings | 无测试（装配/生成产物/纯模型） |

## Phase 3 — 覆盖率门（资损口径）
资损红线路径逐条核对，均有红线用例且绿：

| 资损红线 | 红线用例 | 绿 |
|---|---|---|
| 三层金额闸 L1/L2/L3 不漂移（单一常量源） | GateInterlock_LINK01..04、Pricing_PRICE03（L1）、Settlement_BATCH02/CB01（L2）、L3GateRejects*（L3） | ✅ |
| L2 冷启动回退（无样本不除零，回退绝对闸）【arch-review N2】 | Settlement_CB02_ColdStartFallback（已核源码 settlement.go:172/216 avg==nil 跳月均闸） | ✅ |
| usage 上界（天文账单入口拦截） | Pricing_PRICE02_UsageUpperBound、SubmitUsage_OverMax_400 | ✅ |
| 纯整数计价无浮点 | Pricing_PRICE04_PureIntegerNoFloat、PRICE01_FixedAmount | ✅ |
| reorg 回退 + 两阶段确认 + 去重 | SyncOnce_Reorg_RollbackUnconfirmedSeen、DepositMade_TwoPhaseConfirmation、Dedup_NoDoubleLedger | ✅ |
| 对账单一路径：IsPaid 唯一事件回填 | SyncOnce_BillPaid_OnlyEventSetsIsPaid、PayBill_WritesIntent_NotPaid、Reconcile_AmountMismatch_Alerts | ✅ |
| WalletAuth nonce 重放拒绝（消费式台账，非纯 timestamp）【arch-review N1】 | WalletAuth_ReplayNonce_Rejected/UnknownNonce/WalletMismatch/WrongChainID、E2E_WalletAuth_ReplayNonce_FullChain（已核源码 repository.go:411 Consume 原子条件 UPDATE） | ✅ |
| operatorId 固定映射 + sanity check（不靠 name 比对） | OperatorID_OPID01_FixedMapping/OPID02_SanityCheck/OPID03_ZeroPaymentAddress（已核 cmd/main.go:38/128 seed ID=链上 operatorId） | ✅ |
| 组批熔断 + 分批≤25 + 幂等键 month+batchIndex + 失败续跑 | Settlement_BATCH01_Slicing/CB01_CircuitBreak/IDEM01_ConfirmedNotResent/FAIL01_FailedBatchRetriable | ✅ |
| withdraw 降级 pending 不计入余额 | Withdraw_PendingIntent_NotCountedInBalance、Repo_GetTotalByUserID_ConfirmedDepositOnly | ✅ |
| 6 位精度 abigen Filterer 事件解码 | ParseBillCreated_TotalAmountIncludesFee、ParseUsageDataSubmitted_OnlyUserIndexed、ParseDepositWithdrawn_PrincipalPlusInterest | ✅ |
| 占位零地址不订阅事件 | SyncOnce_SkipsPlaceholderContracts、Deployments_PlaceholderContracts | ✅ |
| chainID 一致校验 | ChainIDMismatchRejected、ValidateChainID、WalletAuth_WrongChainID_Rejected | ✅ |
| AdminAuth 常量时间 + 缺 key fail | AdminAuth_CorrectKey/WrongKey_401/NoHeader_401/MissingKey_StartupFail | ✅ |
| abiHash 校验 | ABIHashMatchesDeployments、VerifyABIHashDetectsMismatch、SignatureTopicsMatchBindings | ✅ |

**低覆盖区性质判定**：services(39.8)/handlers(18.3)/repository(6.2) 数值偏低，逐项核对低覆盖区均为**非资损路径**——展示/查询 CRUD、SMTP stub（NotificationService 本轮非阻塞）、装配 wiring、虚拟号生成展示等。资损红线路径覆盖密集且红线用例齐全。

**口径结论**：以资损红线清单覆盖为准判 **PASS**。低数字覆盖率不构成资损缺口；未发现资损路径无测的缺口。

## Phase 4 — grep 自检
| 项 | 结果 |
|---|---|
| 无硬编码私钥（ORACLE_OWNER/DEPLOYER） | ✅ 仅 `os.Getenv` 读 + 日志 WARN；无 64-hex 私钥字面量（命中的 64-hex 均为 event topic hash） |
| 无硬编码精度绕过 | ✅ 精度从 `deployments.usdtDecimals`(=6) 读；命中的 `1_000_000` 是 MAX_DATA_MB usage 上界常量，非精度除数 |
| deployments.json 无 0G 残留 | ✅ chainId=421614、rpcUrl=sepolia-rollup.arbitrum.io/rpc；无 0g.ai/16602/0g_testnet（地址为占位零，等合约 1/3 上链回填，已 _note 标注） |
| CORS 无 *+credentials | ✅ AllowOrigins 固定白名单（localhost:5173/127.0.0.1:5173/localhost:3000）+ AllowCredentials:true；无 `"*"`（唯一 `"*"` 出现在注释里） |
| IsPaid 无 HTTP 直置 | ✅ MarkAsPaid 仅 services.go 内部（event_sync 调用），handlers.go 无 IsPaid/MarkAsPaid 写；pay 端点仅写 PayIntentTxHash |

## Phase 5 — 验收映射
见 test.md §验收映射（requirement §五② + design §8 + arch-review §七红线 → 实际用例）。子系统② 后端验收项（#5 deployments.json schema、#6 event_sync 真实监听落库）+ design §8 T1–T8 全项 + arch-review N1/N2/operatorId 红线均有用例对应且绿。

## 已知限制（如实记录，不算失败）
- **业务合约真机端到端**（金额数组语义/事件回填上链）未跑：go-ethereum v1.13.5 simulated.Backend 是 London 上限（ethash faker），跑不了项目 Cancun 字节码（transient storage/PUSH0）。本轮链交互用接口 mock（SettlementClient/SettlementBatchStore）+ 手工构造 types.Log（abigen Filterer.Parse* 解码）覆盖；client 机制层（L3/nonce/降级/chainID）用 London 兼容桩字节码 + 真实绑定验证。
- 真机端到端待**合约 1/3 真·上链**后用 hardhat 31337 或升级 geth v1.16+（PoS 后端支持 Cancun）。design §7.4 已声明此依赖不阻塞 implement/test。

## 三自检
1. **承诺**：只跑测试/验证 + 产文档，不修业务代码——✅ 未改任何 .go / pipeline.json / web / contracts；仅写本 checkpoint + test.md。
2. **交付物存在**：docs/pipeline/checkpoints/stage-test-done.md + docs/pipeline/stages/test.md——✅ 均已产出。
3. **正确连接**：Phase 1-5 结果与实际命令输出一致；红线用例名经 `go test -v` 实跑列出核对；N1/N2/operatorId 经 codegraph 源码核对。

## 结论
**DONE — 全绿，可进 review。** go build/vet 0 问题；75 用例全绿（含 -race）；资损红线逐条覆盖；4 项 grep 自检全过；唯一遗留（业务合约真机端到端）为环境依赖，design §7.4 已声明不阻塞。
