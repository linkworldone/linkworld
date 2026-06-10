# Task 02 — T1 链配置 421614 + ABI 重生成 + 单一出口（web 3/3）

> 子项目 web(3/3) | 状态 DONE | 2026-06-10

### 产出文件
- packages/web/src/config/chains.ts（新增 arbitrumSepolia 421614，删 zgMainnet/zgTestnet，留 hardhatLocal）
- packages/web/src/config/wagmi.ts（非本地链 → arbitrumSepolia，VITE_CHAIN_ID 切换沿用）
- packages/web/src/config/contracts.ts（单一出口重写：从 deployments json 派生地址修 31337 漂移，ContractAddresses 补 usdt/usdtDecimals，新增 getUsdt/getUsdtDecimals）
- packages/web/src/config/deployments/{hardhat,arbitrum_sepolia}.json（新增：31337 真源 + 421614 零地址占位待回填）
- packages/web/src/config/abis/*.ts（8 个从 artifacts 重生成：deposit(uint256)/payBill 去 payable/getLockExpiry/getFeeRate/calculateFee/Oracle/MockUSDT approve·allowance·decimals）+ index.ts 补 export OracleABI/MockUSDTABI
- packages/web/tsconfig.json（开 resolveJsonModule）
- packages/web/src/hooks/contracts/*（最小适配保编译：去 value、as never 占位已删方法，标 T3/T4 完善）

### git commit
559fc78 feat: web T1 链配置 421614 + ABI 重生成(补 Oracle) + contracts 单一出口

### TDD
T1 为接链基建（配置/ABI/类型），非新业务行为。验证靠 tsc 0 + build 绿 + T0 smoke 不回归；业务单测 T2-T12 各自写。

### 测试结果
npx tsc --noEmit 0 error。npm run build ✓ built。npm test 3 passed（T0 smoke 不回归）。主 Agent 已独立复跑确认 tsc exit 0 + build ✓ + 3 测过。

### code-simplifier
contracts 单一出口消除手抄地址漂移；ABI 从 artifacts 重生成保证与合约一致；最小 hook 适配保编译不超范围。

### spec review
按 design v2 §1.2/§5 + arch-review + handoff-backend 执行：chains 421614、contracts deployments 单一出口、abis 重生成补 OracleABI、wagmi 链选择。最小 hook 适配（去 value/as never）仅为保 tsc 绿，approve 两步态/精度/对账逻辑严格留 T2/T3/T4，注释标注。未越界换肤/WalletAuth。

### 设计还原
T1 为接链基建无 UI 还原。design §1.2 deployments 单一出口 + §5 契约映射（deposit(uint256)/payBill 去 payable/getLockExpiry/getFeeRate）逐项落地。

### 复用检查
复用 viem defineChain、artifacts ABI 来源、现有 getContractAddress 零地址保护；新增 getUsdt/getUsdtDecimals。

### 设计稿对照
数值对照：新增链 421614 vs R2 ✅；contracts 31337 地址=deployments.json proxies（修漂移）✅；abis 导出 6→8（补 Oracle/MockUSDT）✅；usdtDecimals 6 ✅；tsc error 0 ✅；build 绿 ✅；测试 3 不回归 ✅。

### 新增组件
新增配置：arbitrumSepolia 链、deployments json（拷入 src）、MockUSDT ABI、getUsdt/getUsdtDecimals。无新增业务组件。

### 新增色值
无（T1 接链基建，换肤留 T6+）。

### ⚠️ 遗留（带入 T2/T3/T4）
- T2：format.ts/constants 全链路 6 位 + currency ETH→USDT 未动；deposit 内已用 usdtDecimals 解析可复用。
- T3：approve allowance/exact/TwoStepAction 未实现；payBill 的 _value 占位参数待接 approve 后移除；MockUSDT ABI 就绪。
- T4：useServiceManager（运营商新模型 getOperator/addOperator）+ useTrafficCard（getCardInfo 逐卡，无 destroyedAt）需按新 ABI 重写去掉 as never。
- 421614 真上链阻塞：arbitrum_sepolia.json 零地址占位，端到端验收待合约上链；本地 31337 可验。
- contracts 单一出口当前方案=deployments json 拷入 src，合约重部署需同步该 json，后续可脚本固化。
