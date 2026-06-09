# Task 01 — T1 修编译（P1 compile gate）

> 子项目 contracts(1/3) | 状态 DONE_WITH_CONCERNS | 2026-06-08

### 产出文件
- packages/contracts/contracts/Oracle.sol（声明 UsageDataSubmitted 事件 + deposit/payment 改 IDeposit/IPayment 接口类型 + 新增 setPayment + 收敛内联接口改 import interfaces/）
- packages/contracts/contracts/interfaces/IDeposit.sol（补 issueMonthlyTrafficCards 声明）
- packages/contracts/contracts/interfaces/IPayment.sol（补 applyTrafficCardToBill 声明 + TrafficCardApplied 事件）
- packages/contracts/contracts/TrafficCardNFT.sol（删 DeductionCredit dead code + 删与接口重复的 CardMinted/CardDestroyed 声明）
- packages/contracts/contracts/Payment.sol（applyTrafficCardToBill 受限桩：onlyOracle + require bill 存在 + emit，不转资金）

### git commit
e4bd81d fix: 合约 T1 修编译（Oracle/接口/TrafficCardNFT/Payment 桩）

### TDD
T1 为结构性修编译（非新行为），不走红→绿→重构。验证靠 compile 绿 + 现有测试不回归；applyTrafficCardToBill 桩的单测（ATC-01/02）按 plan 留 T6。

### 测试结果
hardhat compile：Compiled 13 Solidity files successfully (evm: cancun)，0 error（改前 Oracle.sol:71 DeclarationError 中止）。hardhat test：30 passing / 0 failing（改前因编译失败 0 可运行）。主 Agent 已独立复跑确认（Nothing to compile + 30 passing）。

### code-simplifier
改动以删除/类型替换/补声明为主，净新增 < 50 行；遵循现有代码风格，无需额外简化。

### spec review
严格按 design §5.6 编译错误清单 ①②④⑤ + §4.2/v2-B 受限桩 + §4.4/§4.5 执行。createBill onlyOracle（B2）与 monthlySettlement amounts[] 签名（B1/v2-A）按任务边界刻意未做，留 T3/T4——对编译无影响，未越界。

### 设计还原
合约任务无 UI 还原。以"design §5.6 编译清单逐条落地"替代：①UsageDataSubmitted 事件已声明；②Oracle deposit/payment 已改接口类型；④TrafficCardNFT DeductionCredit dead code 已删；⑤内联接口已收敛到 interfaces/ 并补全 2 个缺失声明；Payment 受限桩已实现。

### 复用检查
复用现有 interfaces/IDeposit.sol、IPayment.sol（补声明而非新建）；复用 TrafficCardNFT 继承自 ITrafficCardNFT 的 CardMinted/CardDestroyed 事件（删重复声明后 emit 解析到接口继承）；无新建合约。

### 设计稿对照
合约编译/测试数值对照：编译文件数 13 个 vs 预期全量 ✅；compile error 0 个 vs gate 要求 0 ✅；测试 30 passing / 0 failing vs 基线 30 ✅（无回归）；改动文件 5 个 vs T1 范围 5 个 ✅；Oracle 改动 4 处（事件1+state2类型+setPayment1+import收敛）vs design §4.4 列出 4 项 ✅。

### 新增组件
无新增合约（T1 仅修编译）。新增接口声明 2 个：IDeposit.issueMonthlyTrafficCards、IPayment.applyTrafficCardToBill；新增事件 2 个：Oracle.UsageDataSubmitted、IPayment.TrafficCardApplied。

### 新增色值
无（合约任务，无色值）。
