# Task 03 — T3 client 真实读写 + owner 发交易（后端 2/3）

> 子项目 backend(2/3) | 状态 DONE_WITH_CONCERNS | 2026-06-09

### 产出文件
- packages/backend/internal/blockchain/client.go（重写：读真实化 via bindings Caller；MonthlySettlement/IssueMonthlyTrafficCards owner 发交易 via Transactor；本地 nonce 计数器 mutex+markNonceDirty resync；WaitMined 等回执；L3 金额闸 assertAmountGateL3；EnableOwnerWrites/CanWrite；NewClientWithBackend 注入测试后端；chainID 校验；删死 stub ListenEvents）
- packages/backend/internal/blockchain/client_test.go（新建 simulated.Backend TDD 7 用例）
- packages/backend/internal/blockchain/testdata_bytecode_test.go（新建 London 兼容桩字节码）
- packages/backend/cmd/main.go（装配 EnableOwnerWrites，env 读 owner key，降级日志不打私钥）
- go.mod/go.sum（simulated.Backend 测试依赖）

### git commit
b3dd5dd feat: 后端 T3 client 真实读写 + owner 发交易 + nonce + WaitMined + 金额闸 L3

### TDD
先红后绿：先写 client_test.go → 红（NewClientWithBackend/MaxBillPerUser undefined）→ 实现 client.go → 绿；过程真捕获并修 2 bug（WaitMined/Commit 死锁→autoCommit；committer goroutine 与 Close 竞态→stop 同步等待）。

### 测试结果
go build ./... 0 error。go test ./internal/blockchain/... ok 5.2s（含 -race ok 6.1s）：3 abihash + 7 T3（MonthlySettlementSendsTxSuccessfully/NonceNoCollisionAcrossBatches/L3GateRejects Amount·BatchTotal·Zero/OwnerKeyMissingDegradesWrites/ChainIDMismatchRejected/ReadMethodsReturnOnChainValues）。主 Agent 已独立复跑确认 build exit 0 + test ok。

### code-simplifier
client.go 重写聚焦；nonce/L3/降级各自单一职责；删死 stub。

### spec review
按 design v2 §6.1/§7.0 + arch-review B1/nonce/WaitMined/owner=root/chainID 执行。L3 闸（发交易前断言 amount∈(0,MaxBillPerUser]+sum≤MaxBatchTotal）= B1 最后防线；owner key 仅内存不落盘不日志、缺失降级不 fatal。未越界（event_sync/signatures 留 T4，services 计价留 T5，组批 L2/熔断留 T6，鉴权留 T7）。

### 设计还原
后端无 UI。design §6.1 client 写路径逐项落地：nonce 计数器/WaitMined/L3/降级/chainID 校验，simulated.Backend 验证编码→签名→发交易→解码机制对真实链成立。

### 复用检查
复用 T1 bindings（Caller/Transactor/Deploy）、abihash；复用 T2 config（ResolveRPCURL/ValidateChainID）；NewKeyedTransactorWithChainID（go-ethereum v1.13.5）。

### 设计稿对照
数值对照：L3 闸常量 MaxBillPerUser=1e9/MaxBatchTotal=1e10（6 位，PLACEHOLDER+启动 warn）vs B1 要求 ✅；nonce 连发 3 笔不复用（NonceNoCollision 测）vs 设计 ✅；测试 10 用例全绿（3+7）vs 无回归 ✅；owner key 缺失降级返 error 不 panic vs 设计 ✅；go build error 0 ✅。

### 新增组件
新增 client 方法：MonthlySettlement/IssueMonthlyTrafficCards/EnableOwnerWrites/CanWrite/OwnerNonce/assertAmountGateL3/NewClientWithBackend；无新增业务合约。

### 新增色值
无（后端任务）。

### ⚠️ 遗留（带入 T8/test 阶段）
- **go-ethereum v1.13.5 simulated.Backend 是 London 上限（ethash faker），跑不了项目 Cancun 字节码（transient storage/PUSH0）**。T3 测试用 London 兼容桩字节码+真实绑定，验证了 client 机制（L3/降级/chainID 纯逻辑全真实），但**业务合约端到端结算（金额数组语义/事件回填）需 hardhat 31337 跑，或升级 go-ethereum v1.16+（ethclient/simulated PoS 后端支持 Cancun）**。T8 测试规划须据此调整，不要假定 v1.13.5 simulated.Backend 能部署业务合约。
- T6 失败批续跑用 receipt.Status；发送失败已 markNonceDirty。
- 占位常量 MaxBillPerUser/MaxBatchTotal 待产品拍真值。
