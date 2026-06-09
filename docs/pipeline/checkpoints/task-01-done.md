# Task 01 — T1 abigen 合约绑定 + abiHash 校验（后端 2/3）

> 子项目 backend(2/3) | 状态 DONE | 2026-06-09

### 产出文件
- packages/backend/internal/blockchain/bindings/（8 份 abigen 绑定：feemanager/userregistry/servicemanager/trafficcardnft/payment/deposit/oracle/mockusdt，package bindings）
- packages/backend/internal/blockchain/abihash.go（纯 Go 复刻 ethers formatJson() 规范化 + ComputeABIHash/VerifyABIHash）+ abihash_test.go
- packages/backend/internal/blockchain/abis/（8 份原始 ABI JSON 嵌入 + abis.FS/Names；同步 Deposit/UserRegistry 陈旧文件到 artifacts）
- packages/backend/scripts/genbindings/main.go + scripts/gen-bindings.sh + Makefile(gen-bindings target)

### git commit
a6b26a3 feat: 后端 T1 abigen 8 合约绑定 + abiHash 校验

### TDD
T1 为绑定 codegen（非业务行为），不走红绿。验证靠 go build 0 error + abiHash 单测（TestABIHashMatchesDeployments）；abihash.go 用单测锁定 7 合约哈希与 deployments 一致。

### 测试结果
go build ./... → 0 error。go test ./internal/blockchain/... → PASS（abiHash 7 业务合约与 hardhat.json abiHash 全一致）。主 Agent 已独立复跑确认 build exit 0 + blockchain test ok。

### code-simplifier
绑定为生成产物；abihash.go 复刻 formatJson 规范化逻辑紧凑无冗余。

### spec review
按 design v2 §5/§6 + handoff abiHash 算法执行。关键发现：abiHash = keccak256(formatJson()) 非原始 ABI 哈希，已纯 Go 复刻（离线无 JS 依赖）。abigen CLI 在 Go 1.26 下 memsize 不兼容 → 改用同源 accounts/abi/bind.Bind 库（产物一致），合理偏差。未越界（未碰 client/config/event_sync/services/handlers 逻辑）。

### 设计还原
后端无 UI。以 design §5/§6 落地替代：8 合约绑定生成入库 + abiHash 启动自检工具就位，供 T3 client init 守卫。

### 复用检查
复用 go-ethereum v1.13.5 accounts/abi/bind（abigen 同源）；复用 packages/contracts/artifacts ABI 来源；绑定入库供 go test 离线跑；abis.FS 嵌入供 T2/T3 复用。

### 设计稿对照
数值对照：生成绑定 8 份 vs 合约总数 8 ✅；abiHash 一致 7/7 业务合约（MockUSDT 部署侧无 abiHash）vs handoff 校验要求 ✅；go build error 0 vs gate 0 ✅；新增 .go 文件数（8 绑定+abihash+test+genbindings）vs T1 范围 ✅。

### 新增组件
新增 Go 包：internal/blockchain/bindings（8 合约绑定）；新增 internal/blockchain/abihash.go（ComputeABIHash/VerifyABIHash）；新增 scripts/genbindings + gen-bindings.sh + Makefile target。无新增业务合约。

### 新增色值
无（后端任务）。
