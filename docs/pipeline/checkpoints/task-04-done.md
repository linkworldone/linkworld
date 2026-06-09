# Task 04 — T4 event_sync 真实事件同步（后端 2/3）

> 子项目 backend(2/3) | 状态 DONE | 2026-06-09

### 产出文件
- packages/backend/internal/sync/event_sync.go（重写：FilterLogs 轮询分页 + reorg + K 块两阶段 + 8 事件 abigen Filterer 解码 + 去重 + 对账回填）
- packages/backend/internal/sync/event_sync_test.go（新建 9 用例）
- packages/backend/internal/blockchain/signatures.go（BillCreated 改 4 参 + 补 DepositWithdrawn/TrafficCardApplied/UsageDataSubmitted/CardMinted/ServiceActivated topic）
- packages/backend/internal/models/models.go（Bill 加 OnChainBillID/PayIntentTxHash；Deposit 加 Status/BlockHash；新增 SyncState/ChainEvent + 状态常量）
- packages/backend/internal/repository/repository.go（SetOnChainID/MarkPaidByOnChainID/CreateConfirmed/SyncStateRepository/ChainEventRepository；GetTotalByUserID 仅计 confirmed）
- go.mod/go.sum（测试用纯 Go sqlite glebarez/sqlite；gorm 1.25.7）

### git commit
fe7b98f feat: 后端 T4 event_sync 真实同步（Filterer 解码+reorg+K块确认+两阶段+对账回填）

### TDD
先红后绿（真捕获 3 bug）：①内存库共享致唯一约束冲突→每测独立命名库；②SetOnChainID operator_id=0 死过滤匹配不到→按 user 关联；③reorg 回退只删 ChainEvent 未删关联 pending→rollbackFrom 补删。

### 测试结果
go build ./... 0 error。go test ./internal/sync/... ./internal/blockchain/... 全绿；go test ./... 整模块无回归。9 新用例 PASS：ParseBillCreated_TotalAmountIncludesFee/ParseUsageDataSubmitted_OnlyUserIndexed/ParseDepositWithdrawn_PrincipalPlusInterest/SignatureTopicsMatchBindings/SyncOnce_DepositMade_TwoPhaseConfirmation/Dedup_NoDoubleLedger/BillPaid_OnlyEventSetsIsPaid/Reorg_RollbackUnconfirmedSeen/SkipsPlaceholderContracts。主 Agent 已独立复跑确认 build exit 0 + sync/blockchain test ok。

### code-simplifier
dispatch 按合约地址+topic0 分发；signatures 仅日志/粗过滤，解码统一走 Filterer，消除手写 topic 重复。

### spec review
按 design v2 §6.3/§4.3/§4.4 + arch-review B5/B2/B3 + abigen 解码执行。reorg 两阶段(seen→K块→confirmed)、IsPaid 唯一 BillPaid 回填、withdraw 唯一 DepositWithdrawn、占位零地址跳过。未越界（计价/组批留 T5/T6，handler 侧 bills-pay 降级/withdraw 鉴权留 T7）。

### 设计还原
后端无 UI。design §6.3 event_sync + §4.3 状态机逐项落地：pending→seen→confirmed、reorg 回退重扫、(txHash,logIndex)去重幂等。

### 复用检查
复用 T1 bindings Filterer/Parse*；复用 T2 config.IsPlaceholder 跳过占位；复用现有 repository 扩展而非重写；SignatureTopicsMatchBindings 锁 topic0 与 bindings 一致。

### 设计稿对照
数值对照：8 事件全 Filterer 解码 vs 设计 ✅；BillCreated 第三参=amount+platformFee 含费总额（专测断言）vs handoff ✅；UsageDataSubmitted 只 user indexed（专测）vs 合约 ✅；K 确认数=5(PLACEHOLDER)vs B5 ✅；去重唯一索引(txHash,logIndex) vs 防重复记账 ✅；9 测全绿 vs 无回归 ✅；go build 0 ✅。

### 新增组件
新增模型 SyncState/ChainEvent + 状态常量；新增 repository 方法（SetOnChainID/MarkPaidByOnChainID/CreateConfirmed/SyncState·ChainEventRepository）；event_sync 从 stub→真实。无新增业务合约。

### 新增色值
无（后端任务）。

### ⚠️ 遗留（带入 T6/T7/T8）
- T7 handler 侧：bills/pay 写 PayIntentTxHash 降级 pending（IsPaid 已收口只由 event 置）；withdraw/deposit 写 Status=pending（GetTotalByUserID 已只计 confirmed）；废弃 RecordWithdraw 凭前端 txHash 记账。
- UsageDataSubmitted 金额不一致告警已实现，依赖 T6 计价把 Bill.Amount/OperatorID 写对才能比对。
- CONFIRMATIONS=5 占位待安全/运营按 Arbitrum Sepolia 最终性确认。
- 业务合约端到端仍受 T3 的 Cancun/simulated.Backend 限制（事件解码用手工构造 log 测，已绕开）；T8 端到端需 hardhat 31337。
