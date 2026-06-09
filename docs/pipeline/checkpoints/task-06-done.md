# Task 06 — T6 组批结算 + L2 熔断 + 幂等键（后端 2/3）

> 子项目 backend(2/3) | 状态 DONE | 2026-06-09

### 产出文件
- packages/backend/internal/services/settlement.go（新：SettlementOrchestrator.SettleMonth 组批→L2闸→发交易→据 receipt 落幂等状态；SettlementClient/SettlementBatchStore 接口 mock 友好；NewSettlementBatchRepoStore 含 string→big.Int fail-fast）
- packages/backend/internal/services/settlement_test.go（新：6 用例 + fake client/store）
- packages/backend/internal/services/oracle.go（抽 collectSettlementItems 复用 L1 计价；新增 SettleMonthlyOnChain + SetSettlementOrchestrator）
- packages/backend/internal/models/models.go（新增 SettlementBatch + 唯一键 idx_month_batch + 4 状态常量）
- packages/backend/internal/repository/repository.go（新增 SettlementBatchRepository：Get/Save upsert 按幂等键 + HistoricalConfirmedTotals 供月均 + Migrate）

### git commit
5720277 feat: 后端 T6 组批结算(≤25)+L2熔断(冷启动回退)+month/batchIndex 幂等键

### TDD
先红后绿：先写 settlement_test.go(6 用例) → 红(SettlementItem/NewSettlementOrchestrator/MaxSettlementBatch undefined build failed) → 补 model+repo+settlement.go → 绿(6 PASS，ALERT/WARN 日志可见熔断与失败批续跑路径被走到)。

### 测试结果
go build ./... 0 error。go vet 0 warning。go test ./... 全绿不回归。T6 用例：BATCH01_Slicing/BATCH02_OverMaxBatchTotal/CB01_CircuitBreak/CB02_ColdStartFallback/IDEM01_ConfirmedNotResent/FAIL01_FailedBatchRetriable。主 Agent 已独立复跑确认 build exit 0 + services test ok。

### code-simplifier
collectSettlementItems 抽取复用 L1 计价；MaxBatchTotal 复用 blockchain 常量（L2/L3 同源不漂移）；接口 mock 隔离链交互。

### spec review
按 design v2 §7.0/§6.1 + arch-review B1-L2/冷启动回退红线/幂等键 执行。组批≤25、L2 绝对闸+均值熔断、冷启动回退仅查绝对闸（绝不除零）、month+batchIndex 幂等、失败批续跑。未越界（handler/鉴权留 T7，main.go 装配留 T7，T8 端到端留 T8）。

### 设计还原
后端无 UI。design §7.0 金额闸 L2 + 组批 + 幂等逐项落地：30 user→2 批(25+5)、超绝对闸阻断、超月均×N 熔断、冷启动跳过均值闸、confirmed 批不重发、failed 批续跑。

### 复用检查
复用 blockchain.MaxBatchTotal（L2/L3 同源）+ T5 PricingService(L1 计价)+ChainOperatorID + T3 client.MonthlySettlement+receipt + 现有 repository 扩展；编译期断言 *blockchain.Client 满足 SettlementClient。

### 设计稿对照
数值对照：MaxSettlementBatch=25 vs handoff §5.1 N≤25 ✅；30 user→2 批(25+5)(BATCH-01) ✅；单批超 MaxBatchTotal 阻断(BATCH-02) vs L2 ✅；冷启动无均值仅查绝对闸不除零(CB-02 红线) ✅；幂等键 month+batchIndex confirmed 不重发(IDEM-01) ✅；6 测全绿 ✅；go build 0 ✅。

### 新增组件
新增 SettlementOrchestrator/SettlementClient/SettlementBatchStore/SettlementBatch 模型/SettlementBatchRepository/SettleMonthlyOnChain；无新增业务合约。

### 新增色值
无（后端任务）。

### ⚠️ 遗留（带入 T7/T8）
- T7：SettleMonthlyOnChain(ctx,month)→SettlementSummary 已就绪，TriggerMonthlyBill handler 直接调用返回分批摘要 + 加 AdminAuth；main.go 装配 NewSettlementBatchRepoStore→NewSettlementOrchestrator→oracle.SetSettlementOrchestrator（未装配时返明确 error 降级安全）。
- T8：补 simulated.Backend 端到端真实绑定分批结算（本轮用接口 mock 隔离）。
- 熔断倍数 N=3 占位待产品/运营/安全拍真值。
- 绝对闸超限批归 BlockedBatches 计数但 DB 状态用 pending_review（靠 Note 区分原因），如需单列 blocked 终态可 T7/T8 调整。
