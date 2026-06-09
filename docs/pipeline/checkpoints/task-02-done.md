# Task 02 — T2 config 修复 + 链配置 421614（后端 2/3）

> 子项目 backend(2/3) | 状态 DONE | 2026-06-09

### 产出文件
- packages/backend/internal/config/config.go（Deployments struct 扩展 Usdt/UsdtDecimals/AbiHash + 新增 IsPlaceholder/HasPlaceholders/PlaceholderContracts/ValidateChainID/ResolveRPCURL + nil map 兜底）
- packages/backend/internal/config/config_test.go（新建 table-test 8 函数）
- packages/backend/configs/deployments.json（重写 421614 schema：chainId 421614 + Arbitrum Sepolia rpcUrl + 7 占位零地址 + usdt 占位 + usdtDecimals:6 + abiHash + _note；清除 0G 残留）
- packages/backend/.env.example（RPC/CHAIN_ID=421614 + ORACLE_OWNER_PRIVATE_KEY/DEPLOYER_PRIVATE_KEY/ADMIN_API_KEY 留空占位）
- packages/backend/cmd/main.go（blockchain-init 改 ResolveRPCURL + ValidateChainID + 占位 warn；仅此部分，路由/鉴权留 T7）
- packages/backend/scripts/sync-deployments.sh（上链后回填占位地址）

### git commit
260cc11 feat: 后端 T2 config 键名修复 + deployments.json 421614 schema + RPC/chainID 统一

### TDD
先红后绿：先写 config_test.go → 红（d.Usdt/IsPlaceholder/HasPlaceholders undefined build failed）→ 实现 struct 扩展+5 方法 → 绿（8 函数全 PASS）。

### 测试结果
go build ./... 0 error。go test ./internal/config/... ok（8 函数全 PASS）。go test ./internal/blockchain/... ok（T1 未破坏）。go test ./... 全绿。主 Agent 已独立复跑确认 build exit 0 + config test ok。

### code-simplifier
config 方法聚焦、复用 struct；deployments.json 加 _note 自解释；无冗余。

### spec review
按 design v2 §6.5/§6.7/§7.0 + arch-review 执行：键名 contracts→proxies 对齐、421614 schema、ResolveRPCURL 单一优先级、ValidateChainID、IsPlaceholder 占位保护。未越界（client 发交易/event_sync 订阅/鉴权 均留后续；main.go 仅动 blockchain-init RPC/chainID 部分）。

### 设计还原
后端无 UI。design §6.5/§6.7 逐项落地：config bug 修复（TestLoadDeployments_ParsesProxiesKey 断言 7 项非空）、421614 schema（TestRealDeploymentsJSON 断言无 0G+全占位）、RPC/chainID 统一。

### 复用检查
复用 T1 的 abis/abiHash（abiHash 拷自 hardhat.json）；复用现有 config.Load 结构扩展而非重写；sync-deployments.sh 复用 contracts/deployments 产出。

### 设计稿对照
数值对照：deployments chainId 421614 vs R2 目标 ✅；usdtDecimals 6 vs 合约 ✅；proxies 解析 7/7 非空 vs bug 修复 ✅；0G 残留 0 处 vs 清除要求 ✅；config 测试 8 函数 PASS vs 无回归 ✅；go build error 0 ✅。

### 新增组件
新增 config 方法：IsPlaceholder/HasPlaceholders/PlaceholderContracts/ValidateChainID/ResolveRPCURL；新增 scripts/sync-deployments.sh。无新增业务合约。

### 新增色值
无（后端任务）。
