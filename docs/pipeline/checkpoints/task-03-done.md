# Task 03 — T2 ERC20 迁移（Deposit 合约）(P2)

> 子项目 contracts(1/3) | 状态 DONE | 2026-06-08

### 产出文件
- packages/contracts/contracts/Deposit.sol（引入 SafeERC20/IERC20 + `IERC20 public usdt` state + initialize(userRegistry,usdt) 注入 + deposit() payable→deposit(uint256) 含 MIN_DEPOSIT require + safeTransferFrom + withdraw 改 CEI+safeTransfer）
- packages/contracts/contracts/interfaces/IDeposit.sol（deposit() payable→deposit(uint256 amount)）
- packages/contracts/contracts/mocks/MockUSDT.sol（新建：OZ ERC20，decimals=6，public mint，测试/测试网用）
- packages/contracts/test/erc20.ts（新建：T2 TDD 用例）
- packages/contracts/test/linkworld.ts（deployAndWire 补 MockUSDT+2参 initialize；DP-01/02 由 {value} 改写为 approve+deposit）

### git commit
5fcccff feat: 合约 T2 Deposit ERC20 迁移（USDT approve/transferFrom + 最小额）

### TDD
先红后绿：先写 MockUSDT 桩 + erc20.ts → 跑 RED（initialize 2 参不存在 fragment 报错）→ 实现 Deposit ERC20 化 → GREEN 8 passing → 修 deployAndWire/DP 旧用例 → 全量 38 passing。

### 测试结果
hardhat compile 0 error。hardhat test：38 passing / 0 failing（原 30 + 新 8：DEC-01/MIN-01/MIN-02/ERC-01/ERC-02/ERC-02b/ERC-03/REG-02）。主 Agent 已独立复跑确认 38 passing。

### code-simplifier
ERC20 改写以替换资金通道为主，using SafeERC20 简化转账调用；改动聚焦 Deposit，无冗余。

### spec review
严格按 design §4.1/§5.1 执行：去 payable、initialize 注入 usdt（非 setUSDT）、MIN_DEPOSIT=10*10**usdt.decimals() 不硬编码、CEI safeTransfer、SafeERC20 全覆盖、DepositMade emit 不变。唯一合理偏差：withdraw 的 nonReentrant 按 plan-review 列入 T4（A1 发卡路径统一加），T2 用纯 CEI（状态先清零再 safeTransfer 到 msg.sender），未越界引入 mixin。

### 设计还原
合约无 UI。design §4.1 接口前后对比表逐项落地：deposit 签名/withdraw 资金通道/MIN_DEPOSIT/usdt state/IDeposit 同步全部按表实现并测试覆盖。

### 复用检查
复用 OZ SafeERC20/IERC20（无新依赖，OZ 5.6.1 自带）；复用现有 userRegistry.isRegistered 校验、锁仓续期算法（未改）；新建 MockUSDT 仅测试用。

### 设计稿对照
数值对照：MIN_DEPOSIT=10*10**6=10_000_000（6位精度）vs 验收 A.2「10 USDT 起」✅；MIN-01 存 9_999_999 拒绝 / MIN-02 存 10_000_000 通过 ✅；测试 38 passing（原30+新8）vs 预期无回归 ✅；compile error 0 vs gate 0 ✅；改动文件 5 个 vs T2 范围 ✅；mock USDT decimals=6 vs 真实 USDT 6 位 ✅。

### 新增组件
新增合约：MockUSDT.sol（测试/测试网用 ERC20，decimals=6）。无新增业务合约。

### 新增色值
无（合约任务）。
