# Task 06 — T5 部署脚本 + 421614 网络 + wiring + handoff (P5)

> 子项目 contracts(1/3) | 状态 DONE | 2026-06-08

### 产出文件
- packages/contracts/hardhat.config.ts（新增 arbitrum_sepolia 421614 网络，RPC env fallback，无 DEPLOYER_PRIVATE_KEY 时空数组不报错；保留 31337/16600/16602 + cancun/viaIR）
- packages/contracts/scripts/deploy.ts（步骤0 部署 MockUSDT；Deposit/Payment initialize 参数数组同步新签名；部署顺序调整 MockUSDT+ServiceManager+FeeManager 先于 Deposit/Payment；补 payment.setOracle/oracle.setPayment/11 运营商 setOperatorPaymentAddress/NFT transferOwnership；部署后 wiring 断言；deployments.json 加 usdt/usdtDecimals/abiHash/storageLayout）
- packages/contracts/.env.example（新建：DEPLOYER_PRIVATE_KEY/ARBITRUM_SEPOLIA_RPC 样例）
- docs/design/linkworld-contracts/handoff-backend.md（新建：后端子项 2/3 handoff——ABI/selector 变更 + 金额精度语义 + usdt 来源）

### git commit
ebc4b41 feat: 合约 T5 部署脚本+421614 网络+MockUSDT+wiring+后端 handoff

### TDD
T5 为部署/配置基础设施，非新合约行为。验证靠本地部署冒烟（in-process hardhat，部署+wiring+断言全过）+ test 66 passing 不回归。

### 测试结果
hardhat compile 0 error。本地部署冒烟（--network hardhat）：7 proxy + MockUSDT 部署成功，11 运营商 paymentAddress 注入，NFT ownership 转 Deposit，All wiring assertions passed（payment.oracle/oracle.payment/oracle.deposit/deposit.oracle/usdt/serviceManager/operator≠0 全断言通过）。hardhat test：66 passing。主 Agent 已独立复跑部署冒烟确认断言全过。

### code-simplifier
deploy.ts 增量复用现有 deployProxy 骨架，wiring 集中断言；无冗余。

### spec review
严格按 design §7.0/§7.1/§7.2/§7.3 + arch-review deploy initializer 同步执行：MockUSDT 步骤0、initialize 参数同步、★ B2 必补 wiring（setOracle/setPayment）、11 运营商零地址注入、部署后断言、deployments.json + handoff 交付物。无偏差。

### 设计还原
合约无 UI。design §7.0.3 三条权限链（createBill/发卡/桩）前置 wiring 由部署后断言兜底验证；真实 monthlySettlement 一笔交易由 T6 MS-03 集成测试覆盖。

### 复用检查
复用现有 deploy.ts UUPS deployProxy 骨架、contracts/mocks/MockUSDT.sol（T2 建）、NFT transferOwnership(deposit)（现状 L88）；新增 arbitrum_sepolia 网络配置 + handoff 文档。

### 设计稿对照
数值对照：新增网络 chainId 421614 vs R2 ✅；MockUSDT decimals=6 vs v2-C ✅；wiring 补全项 5+ 处（setOracle/setPayment/11×setOperatorPaymentAddress/transferOwnership）vs §7.2 ✅；部署后断言项 9 项 vs §7.0.3 ✅；deployments.json 新增字段 usdt/usdtDecimals:6/abiHash/storageLayout vs §7.3 ✅；test 66 passing 无回归 ✅。

### 新增组件
无新增合约。新增配置：arbitrum_sepolia 网络、.env.example；新增文档：handoff-backend.md。

### 新增色值
无（合约任务）。

### 遗留（带入 T6 / 上链）
- 421614 真·上链未执行（无 DEPLOYER_PRIVATE_KEY/RPC）→ handoff §10：配 key 后 `hardhat run scripts/deploy.ts --network arbitrum_sepolia` + §6.4 TSTORE 链上实测（实测 payBill 确认 transient guard，否则降级）。
- T6：GAS-01 压测量 monthlySettlement/issueMonthlyTrafficCards 单批安全上限 N 回填 handoff；MS-03 端到端跑真实 monthlySettlement 一笔。
- deployments/hardhat.json 是 in-process 冒烟产物（确定性地址非真部署），未纳入 commit。
