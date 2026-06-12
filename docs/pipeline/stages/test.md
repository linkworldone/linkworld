# test 阶段报告 — LinkWorld 后端子项目（2/3）

> 状态：**DONE（全绿，可进 review）** | 日期：2026-06-09 | 子项目：packages/backend | Round：1
> 角色：test runner（只跑验证 + 产出文档，**不改任何业务 .go / web / contracts / pipeline.json**）
> 资损敏感后端（owner=平台 root 权限 + 月度结算代发），标准从严。
> 注：合约子项目（1/3）的 test 报告见 git 历史（本文件前一版本），1/3 已全绿进 review。

---

## 〇、结论

**packages/backend 静态门 + 单测全绿，资损红线逐条覆盖，质量门全过，可进 review 阶段。**

- `go build ./...` / `go vet ./...`：0 error / 0 warning
- `go test ./...`：**75 用例全 PASS，0 FAIL**；`-race` 全包通过
- 资损红线（三层金额闸不漂移 + 冷启动回退 / usage 上界 / reorg 两阶段 / 对账单一路径 IsPaid 唯一事件回填 / WalletAuth nonce 重放拒绝 / operatorId 固定映射 / 组批熔断）逐条有红线用例且绿
- grep 自检 4 项全过（无硬编码私钥 / 无精度绕过 / 无 0G 残留 / CORS 无 *+credentials / IsPaid 无 HTTP 直置）
- 唯一遗留：业务合约真机端到端（go-ethereum v1.13.5 simulated.Backend 跑不了 Cancun 字节码）——design §7.4 已声明不阻塞 implement/test

---

## 一、Phase 1 — 静态门

| 命令 | 结果 |
|---|---|
| `go build ./...` | 0 error（BUILD_EXIT=0） |
| `go vet ./...` | 0 warning（VET_EXIT=0） |

环境：go1.26.1 darwin/arm64；模块 `linkworld-backend`；`github.com/ethereum/go-ethereum v1.13.5`。

---

## 二、Phase 2 — go test ./... -cover

- `go test ./... -count=1 -v`：`--- PASS` **75** / `--- FAIL` **0**
- `go test ./... -race -count=1`：blockchain / config / handlers / middleware / repository / services / sync 全 ok

各包覆盖率：

| 包 | 覆盖率 | 备注 |
|---|---|---|
| internal/blockchain | 72.9% | 含 client L3/nonce/降级/chainID + abigen 解码 |
| internal/middleware | 73.7% | AdminAuth + WalletAuth nonce 台账 |
| internal/sync | 65.6% | reorg/两阶段/去重/对账/占位跳过 |
| internal/config | 48.0% | deployments 解析 + RPC 优先级 |
| internal/services | 39.8% | 计价 L1/组批 L2 红线密；展示/装配未单测 |
| internal/handlers | 18.3% | bills/pay/withdraw/usage 鉴权降级红线已覆盖；展示查询未单测 |
| internal/repository | 6.2% | confirmed-only 余额 / SetPayIntent 红线已覆盖；CRUD 展示未单测 |
| cmd / bindings / abis / models / genbindings | 无测试 | 装配 / 生成产物 / 纯模型 |

---

## 三、Phase 3 — 覆盖率门（资损口径）

**口径**：资损红线路径必须逐条有红线用例且绿；低数字覆盖区若全为非资损（展示/装配/SMTP stub/CRUD）则以红线清单覆盖为准判 PASS。

低覆盖区核对结论：services/handlers/repository 数值偏低主因为**非资损展示/查询/装配/SMTP stub** 未单测（本轮范围 out）。资损红线路径覆盖密集，逐条核对如下（用例名经 `go test -v` 实跑列出）：

| 资损红线 | 红线用例 | 状态 |
|---|---|---|
| 三层金额闸 L1/L2/L3 不漂移（单一常量源 blockchain.MaxBatchTotal） | GateInterlock_LINK01..04 / Pricing_PRICE03 / Settlement_BATCH02·CB01 / L3GateRejects* | ✅ |
| **L2 冷启动熔断回退**（无样本不除零，回退绝对闸）【arch-review N2】 | Settlement_CB02_ColdStartFallback（源码核：settlement.go historicalAverage 无样本返回 nil，settleOneBatch avg==nil 跳月均闸保留绝对闸） | ✅ |
| usage 上界（天文账单入口拦截 + Price 入口拒绝） | Pricing_PRICE02_UsageUpperBound / SubmitUsage_OverMax_400 / SubmitUsage_WithinMax_OK | ✅ |
| 纯整数计价无浮点 | Pricing_PRICE04_PureIntegerNoFloat / PRICE01_FixedAmount / PRICE05_UnknownOperator | ✅ |
| reorg 回退 + 两阶段确认（pending→seen→confirmed 等 K 块）+ 去重 | SyncOnce_Reorg_RollbackUnconfirmedSeen / DepositMade_TwoPhaseConfirmation / Dedup_NoDoubleLedger | ✅ |
| **对账单一路径：IsPaid 唯一 BillPaid 事件回填**，HTTP 不置终态 | SyncOnce_BillPaid_OnlyEventSetsIsPaid / PayBill_WritesIntent_NotPaid / Reconcile_AmountMismatch_Alerts / Reconcile_UsageDataSubmitted_ImmediateConfirmed | ✅ |
| **WalletAuth nonce 重放拒绝**（消费式台账，绑 chainId，非纯 timestamp）【arch-review N1】 | WalletAuth_ReplayNonce/UnknownNonce/WalletMismatch/WrongChainID_Rejected / WalletAuth_ValidSignature_Passes / E2E_WalletAuth_ReplayNonce_FullChain（源码核：repository.go Consume 原子条件 UPDATE + RowsAffected==1） | ✅ |
| **operatorId 固定映射 + sanity check**（seed ID=链上 operatorId，不靠 name 比对） | OperatorID_OPID01_FixedMapping / OPID02_SanityCheck / OPID03_ZeroPaymentAddress / OperatorID_SeedShape（源码核：cmd/main.go seed ID 显式写 + 启动读链校验） | ✅ |
| 组批熔断 + 分批 ≤25 + 幂等键 month+batchIndex + 失败续跑 | Settlement_BATCH01_Slicing / CB01_CircuitBreak / IDEM01_ConfirmedNotResent / FAIL01_FailedBatchRetriable | ✅ |
| withdraw 降级 pending 不计入余额 | Withdraw_PendingIntent_NotCountedInBalance / Repo_GetTotalByUserID_ConfirmedDepositOnly / Repo_SetPayIntent_IntentOnly | ✅ |
| 6 位精度 abigen Filterer 事件解码（BillCreated 含费 / Usage 仅 user indexed / Withdrawn 本金+利息） | ParseBillCreated_TotalAmountIncludesFee / ParseUsageDataSubmitted_OnlyUserIndexed / ParseDepositWithdrawn_PrincipalPlusInterest | ✅ |
| 占位零地址不订阅事件 | SyncOnce_SkipsPlaceholderContracts / Deployments_PlaceholderContracts / TestIsPlaceholder | ✅ |
| chainID 一致校验（transactor==链上） | ChainIDMismatchRejected / ValidateChainID | ✅ |
| AdminAuth 常量时间 + 缺 key 启动 fail | AdminAuth_CorrectKey/WrongKey_401/NoHeader_401/MissingKey_StartupFail | ✅ |
| abiHash 校验 + topic 一致 | ABIHashMatchesDeployments / VerifyABIHashDetectsMismatch / SignatureTopicsMatchBindings / MockUSDTHashStable | ✅ |
| owner key 缺失降级关闭写 | OwnerKeyMissingDegradesWrites / MonthlySettlementSendsTxSuccessfully / NonceNoCollisionAcrossBatches | ✅ |

**结论：资损红线全覆盖且绿，未发现资损路径无测的缺口 → 覆盖率门 PASS（资损口径）。**

---

## 四、Phase 4 — grep 自检

| 项 | 结果 | 证据 |
|---|---|---|
| 无硬编码私钥 | ✅ | ORACLE_OWNER 仅 `os.Getenv` 读 + 缺失 WARN；无 64-hex 私钥字面量（命中的 64-hex 全为 event topic hash / abigen 注释） |
| 无硬编码精度绕过 | ✅ | 精度从 `deployments.usdtDecimals`(=6) 读；命中的 `1_000_000` 是 services.MaxDataMB usage 上界常量，非精度除数 |
| deployments.json 无 0G 残留 | ✅ | chainId=421614、rpcUrl=`https://sepolia-rollup.arbitrum.io/rpc`；grep 0g.ai/16602/0g_testnet 无命中。合约地址为占位零（等合约 1/3 上链回填，_note 已标注 sync-deployments.sh 流程） |
| CORS 无 *+credentials | ✅ | AllowOrigins 固定白名单 + AllowCredentials:true；代码无 `"*"` origin（唯一 `"*"` 在注释中说明禁用） |
| IsPaid 无 HTTP 直置 | ✅ | MarkAsPaid 仅 services.go 内部（注释标注唯一由 event_sync 调用）；handlers.go 无 IsPaid/MarkAsPaid 写，pay 端点仅写 PayIntentTxHash |

---

## 五、Phase 5 — 验收映射

### requirement §五 子系统②（后端）
| # | 验收项 | 用例 / 证据 | 状态 |
|---|---|---|---|
| 5 | deployments.json chainId=421614 + Arbitrum RPC + schema（proxies/usdt/usdtDecimals/abiHash） | RealDeploymentsJSON / LoadDeployments_ParsesProxiesKey / ParsesExtendedFields；Phase 4 grep 核 | ✅（地址占位零待上链，schema/链配置已对齐） |
| 6 | event_sync 能监听 ERC20 改写后真实事件、落库、占位零地址不报错跳过 | SyncOnce_*（DepositMade/BillPaid/Reorg/Dedup/SkipsPlaceholder）+ Parse*（6 位精度解码） | ✅（机制层；真机落库待上链） |

> §五⑰ 主流程端到端真机走查（31337/421614）属 web 3/3 端到端阶段 + 合约上链后联调，非后端单测范围（design §7.4 不阻塞）。

### design §8 落地计划（T1–T8）
| 任务 | 红线用例 | 状态 |
|---|---|---|
| T1 abigen 绑定 + abiHash | ABIHashMatchesDeployments / SignatureTopicsMatchBindings / MockUSDTHashStable | ✅ |
| T2 config 修复 + schema | LoadDeployments_ParsesProxiesKey/ParsesExtendedFields / ResolveRPCURL | ✅ |
| T3 client 读/写 + nonce + L3 闸 | ReadMethodsReturnOnChainValues / MonthlySettlement* / Nonce* / L3GateRejects* | ✅ |
| T4 owner 签名 + chainID + 降级 | OwnerKeyMissingDegradesWrites / ChainIDMismatchRejected / ValidateChainID | ✅ |
| T5 event_sync + reorg + 两阶段 | SyncOnce_* / Parse* / Reconcile_* | ✅ |
| T6 计价 + 护栏 + operatorId + 分批 | Pricing_PRICE* / Settlement_* / OperatorID_* | ✅ |
| T7 鉴权 + 对账降级路径 | AdminAuth_* / WalletAuth_* / E2E_WalletAuth_* / PayBill_WritesIntent / Withdraw_PendingIntent / SubmitUsage_* | ✅ |
| T8 测试收口 + simulated.Backend | GateInterlock_* / Reconcile_* / Trigger* / Repo_* + 真实绑定（London 兼容桩） | ✅ |

### arch-review §七 验收红线
| 红线 | 状态 |
|---|---|
| N1 WalletAuth nonce 台账（服务端一次性消费 + 绑 chainId，禁纯 timestamp） | ✅ 源码 + 用例核 |
| N2 L2 冷启动熔断回退（首月/无样本回退绝对闸不失效，不除零） | ✅ 源码 + 用例核 |
| operatorId 固定映射 + sanity check（不靠 name 比对） | ✅ 源码 + 用例核 |

---

## 六、已知限制（如实记录，不算失败）

- **业务合约真机端到端**（monthlySettlement 金额数组语义 + 事件回填上链）未跑：go-ethereum v1.13.5 `simulated.Backend` 为 London 上限（ethash faker），跑不了项目 Cancun 字节码（transient storage / PUSH0）。
- 本轮策略：① 链交互逻辑用接口 mock（SettlementClient / SettlementBatchStore）+ 手工构造 `types.Log`（abigen Filterer.Parse* 解码）覆盖；② client 机制层（L3 / nonce / 降级 / chainID）用 London 兼容桩字节码 + 真实绑定验证。
- 真机端到端待**合约 1/3 真·上链**后用 hardhat 31337 或升级 geth v1.16+（PoS 后端支持 Cancun）。design §7.4 已声明此依赖不阻塞 implement/test。

---

## 七、可否进 review

**可进 review。** 静态门 + 75 单测全绿（含 -race）；资损红线（三层金额闸不漂移 / 冷启动回退 / usage 上界 / reorg 两阶段 / 对账单一路径 IsPaid 唯一事件回填 / WalletAuth nonce 重放拒绝 / operatorId 固定映射 / 组批熔断）逐条覆盖且绿；4 项 grep 自检全过；唯一遗留为环境依赖（业务合约真机端到端），design §7.4 已明确不阻塞。
