# Task 05 — T5 计价 PricingService + operatorId 固定映射（后端 2/3）

> 子项目 backend(2/3) | 状态 DONE | 2026-06-09

### 产出文件
- packages/backend/internal/services/pricing.go（新：PricingService + OperatorRate 费率表 big.Int 6 位 + Price() + 占位费率/上界常量）
- packages/backend/internal/services/operatorid.go（新：SeedOperators 单一事实源 + ChainOperatorID + SanityCheckOperators + OperatorChainReader 接口 + ServiceManagerCaller 适配器）
- packages/backend/internal/services/{pricing_test.go,operatorid_test.go}（新：PRICE-01..05 + OPID-01..03）
- packages/backend/internal/services/oracle.go（删 GetBill rand.Intn 随机金额；OperatorAPI 收敛为只 GetUsage；OracleServiceV2 注入 PricingService；FetchAndCreateBills 改 GetUsage→Price 确定计价）
- packages/backend/cmd/main.go（seed 改 SeedOperators 显式 ID=1..11 + 构造 PricingService + 启动 operatorId sanity check 占位降级）

### git commit
cb0576e feat: 后端 T5 计价 PricingService(费率表+usage上界L1) + operatorId 固定映射

### TDD
先红后绿：先写 pricing_test.go+operatorid_test.go → 红(NewPricingService/SeedOperators/ChainOperatorID/MaxDataMB undefined build failed) → 实现 → 绿(ok 0.68s)。

### 测试结果
go build ./... 0 error。go test ./... 全绿无回归。T5 用例：PRICE-01 固定金额/PRICE-02 usage 上界/PRICE-03 单 bill 上限/PRICE-04 纯整数无浮点(5 子)/PRICE-05 未知运营商/OPID-01 固定映射/OPID-02 sanity check/OPID-03 零地址/SeedShape。主 Agent 已独立复跑确认 build exit 0 + services test ok。

### code-simplifier
SeedOperators 单一事实源消除 seed 重复；MaxBillPerUser 复用 blockchain 常量（L1/L3 单一事实源不漂移）；OperatorAPI 接口收敛只剩 GetUsage。

### spec review
按 design v2 §4.2/§4.5/§7.0 + arch-review B4 + CEO(占位刺眼化/usage 仍模拟) 执行。usage 上界 L1+单 bill 上限+纯整数计价、operatorId 固定映射不靠 name 比对+sanity check。未越界（组批/熔断 L2 留 T6，handler/鉴权留 T7）。

### 设计还原
后端无 UI。design §4.2 计价 + §4.5 operatorId 逐项落地：amount6=dataMB×单价+callMin×单价 纯 big.Int，seed ID=链上 operatorId 强契约。

### 复用检查
复用 blockchain.MaxBillPerUser（L1/L3 同源）；复用现有 OracleServiceV2 注入扩展；SeedOperators 与 ServiceManager.initialize(id=1..11,countryCode US..PH) 逐一核对对齐。

### 设计稿对照
数值对照：PRICE-01 精确断言 100MB×10000+10min×5000=1_050_000 vs 计价公式 ✅；usage 上界 MaxDataMB=1e6/MaxCallMin=1e5 vs B4 ✅；单 bill 上限==blockchain.MaxBillPerUser（PRICE-03 断言相等，L1/L3 不漂移）✅；operatorId seed ID 1..11 vs 链上 ✅；占位费率刺眼 PLACEHOLDER+启动 warn vs CEO ✅；9 测全绿 ✅；go build 0 ✅。

### 新增组件
新增 PricingService/OperatorRate/SeedOperators/ChainOperatorID/SanityCheckOperators/OperatorChainReader；删 GetBill。无新增业务合约。

### 新增色值
无（后端任务）。

### ⚠️ 遗留（带入 T6/T7）
- T6：FetchAndCreateBills 当前是逐 user 写 DB Bill 骨架（已用 PricingService 出确定金额），T6 重写为组批 ≤25 + L2 熔断 + 逐批 client.MonthlySettlement + 月度幂等键；组批传 operatorIds[] 用 services.ChainOperatorID(op.ID)。Bill.PlatformFee 当前 "0"，待用链上 FeeManager 回填。
- T7：TriggerMonthlyBill handler 待加 AdminAuth + 分批结果摘要。
- usage 仍 rand 模拟（计价引擎正确≠计费正确，真实计量后续 Round）；占位费率/上界待产品拍真值。
