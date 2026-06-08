# Stage: design — 合约技术设计（子项目 contracts 1/3）

> **状态**: completed | **日期**: 2026-06-08 | **Gate**: 1（设计） | **子项目**: contracts(1/3) | **角色**: 架构师
> 输入：requirement.md（PRD §三①/§五A·D/§六）+ scan.md（编译不过 + 7 发现）+ 4 份基线（project-scan / contracts-inventory / erc20-migration-surface / test-deploy-baseline）。
> 范围：**只产设计**，不改任何 .sol / 脚本 / 测试 / 配置（实现留 implement 阶段）。
> 注：本文件取代 2026-06-07 旧 design.md（web 视觉版，属 web 子项 3/3，已因全栈反转重新分项）。
> 锁定决策（用户已拍板，不再翻）：① 全补齐不留 stub；② 运营商分账=链上直分（amount→operator.paymentAddress，fee→platform）；③ 自动发卡入口由后端/Oracle 定时调用；④ mock USDT 6 位精度，合约一律读 `usdt.decimals()` 不硬编码；⑤ 部署 Arbitrum Sepolia(421614) + 本地 hardhat(31337)，不上主网。

---

## 〇、已核对真实代码清单（含行号引用）

本设计基于逐行核对，引用以下源码（均已 Read）：

| 文件 | 关键行号 / 核对结论 |
|------|---------------------|
| `contracts/Deposit.sol` | `deposit() payable` L41-53；`withdraw() .call{value}` L56-68；`mintTrafficCard onlyOwner` L71-81，`dataAmount=_deposits/100000` L76；`issueMonthlyTrafficCards` L104-107 **空实现**；state L13-21（无 usdt） |
| `contracts/Payment.sol` | `payBill() payable` L76-97；fee→platformWallet L86-89，**amount 无出口**；`createBill onlyOwner` L55-73；`ReentrancyGuardTransient` L11 + 自定义 `_reentrancyGuardInit` assembly L33-38；**无 `applyTrafficCardToBill`**；**无 `setUSDT`** |
| `contracts/FeeManager.sol` | `FEE_DENOMINATOR=10000` L10、`MAX_FEE_RATE=1000` L11、`calculateFee=amount*rate/10000` L34-36；满足 R5，本轮不改逻辑 |
| `contracts/TrafficCardNFT.sol` | `mint onlyOwner` L45-62；`burn` L83-94 emit `ServiceActivated(now+30days)`；**L18 `DeductionCredit` 类型未定义（编译错误②）**；`onlyDeposit`/`_deductionCredits` dead code L18/L27-30 |
| `contracts/Oracle.sol` | `deposit`/`payment` 声明为 `address` L10-11 却调 `.createBill`/`.issueMonthlyTrafficCards`/`.applyTrafficCardToBill`（**编译错误①根因**）；`UsageDataSubmitted` L71 **未声明（实测唯一报错点）**；**无 `setPayment`**；文件尾内联 `IPayment`/`IDeposit` L105-122 |
| `contracts/ServiceManager.sol` | 11 内置运营商 L21-153，`requiredDeposit` 用 `0.0xx ether`（18 位）；**全部 `paymentAddress=address(0)`**；`addOperator` 已带 paymentAddress 入参 L156-177，`updateOperator` L179-192 **无 paymentAddress 入参** |
| `contracts/UserRegistry.sol` | 邮箱注册 L22-37，无 phone（对齐 R10）；本轮不改 |
| `contracts/interfaces/*.sol` | `IDeposit.deposit() payable` L9；`IPayment.payBill() payable` L19，`IPayment` **无 applyTrafficCardToBill**；`IServiceManager.Operator.paymentAddress` 已有 L12 |
| `scripts/deploy.ts` | UUPS 全量 L9-89；wiring **缺 `payment.setOracle`、`oracle.setPayment`、operator.paymentAddress、mock USDT**；`transferOwnership(deposit)` L88 |
| `hardhat.config.ts` | solidity 0.8.27 / `evmVersion:"cancun"` L15 / `viaIR:true` L16；networks 仅 31337/16600/16602 L19-35，**无 421614** |
| `test/linkworld.ts` | 直接 `Factory.deploy()`+手动 `initialize`（不走 proxy）；全用 `parseEther`（18 位）；**无 withdraw/payBill/mintTrafficCard/最小额/升级 测试** |

**实跑 `npx hardhat compile` 结论**：编译在 `Oracle.sol:71 UsageDataSubmitted` 第一个 DeclarationError 处中止（编译器只报首错）。修掉它后会暴露后续错误（Oracle address-vs-contract、TrafficCardNFT `DeductionCredit`），见 §六编译错误清单。

---

## 一、背景与目标

### 1.1 为什么做这次合约改造
合并自 origin 的合约（commit `7ef9677`）是**编译不过的半成品**，且 PRD Round1 全栈反转后确定把资金通道从原生币切到 ERC20 USDT、目标链切到 Arbitrum Sepolia。本轮合约层（子项 1/3）是后端（2/3）/web（3/3）的地基：合约不编译、ABI 不稳、地址不产出，后两子项无法对齐。

### 1.2 本轮范围边界
**做**：① 修编译（3 类错误）；② Deposit/Payment ERC20 USDT 改写（approve+transferFrom、SafeERC20、最小额 10 USDT）；③ Payment 链上直分（amount→operator.paymentAddress，fee→platform）；④ TrafficCardNFT mint/burn/有效期 wiring；⑤ Oracle wiring（setPayment）+ `issueMonthlyTrafficCards` 真实自动发卡逻辑（不留 stub）；⑥ mock USDT(6 decimals) 自部署；⑦ hardhat.config 加 421614 + deploy.ts 补 wiring；⑧ 补测（验收 A.2 等）。

**不做（红线）**：❌ 不上 Arbitrum One 主网；❌ 不加 phone 注册（R10）；❌ 不改业务路由/前端/后端；❌ 不引入新业务功能（实体 SIM 等留后续 Round）。

### 1.3 技术指标
| 指标 | 目标 |
|------|------|
| 编译 | `hardhat compile` 零 error（当前 1 error 中止） |
| 单测 | 新增覆盖验收 A.2（<10 拒绝/=10 通过）、approve/transferFrom、分账正确性、自动发卡权限/时序、锁仓提取、精度；现有 26 it 不回归 |
| 部署 | 31337 + 421614 两处 `deploy.ts` 成功，产出 `deployments/<net>.json`（含 7 proxy + mock USDT + decimals） |
| 精度 | 合约内零硬编码 18/6，`MIN_DEPOSIT=10*10**usdt.decimals()` |
| 资损 | SafeERC20 全覆盖、分账零地址 require、CEI + ReentrancyGuard、整数运算无溢出/截断隐患 |

---

## 二、整体方案

### 2.1 总体顺序（5 大模块，强依赖拓扑）

```
① 修编译（地基，必须最先）
   └─ Oracle: address→接口类型 + 声明 UsageDataSubmitted 事件 + 加 setPayment
   └─ TrafficCardNFT: 定义/删除 DeductionCredit
   └─ Payment: 实现 applyTrafficCardToBill（Oracle 已在调它）
        ↓ 编译通过后才能继续
② ERC20 迁移（Deposit/Payment 资金通道 + SafeERC20 + usdt state + 最小额）
        ↓
③ 分账（Payment.payBill 链上直分 amount→operator.paymentAddress, fee→platform）
   └─ 依赖 ServiceManager 暴露 paymentAddress + 设置入口 + 零地址校验
        ↓
④ 自动发卡（Deposit.issueMonthlyTrafficCards 真实 mint 逻辑 + onlyOracle）
        ↓
⑤ 部署（hardhat.config 加 421614 + mock USDT + deploy.ts 补全 wiring）
        ↓
⑥ 补测（贯穿，TDD 优先）
```

> 顺序铁律：**①修编译必须最先**。当前代码连 `compile` 都过不了，任何 ERC20/分账/发卡逻辑都无法验证。implement 阶段先把 ① 跑到 `compile` 绿，再动 ②③④。

### 2.2 关键选型对比

| 决策点 | 选项 A | 选项 B | **本轮选定** | 理由 |
|--------|--------|--------|-------------|------|
| ERC20 安全转账 | 裸 `IERC20.transfer/transferFrom` | **SafeERC20** | **B** | OZ 5.6.1 自带，处理无返回值/返回 false 的非标准代币；资损敏感必须用 |
| Arbitrum 落地方式 | `upgradeProxy` 升级 16602 旧 proxy | **fresh deploy 全新部署 421614** | **B** | 421614 是新链无旧 proxy；fresh deploy 无 storage layout 约束，最干净（scan §五 已确认 421614 零配置） |
| 16602/localhost 旧 proxy | 升级 | **本轮作废/可选升级** | 作废为主 | 验收只认 31337+421614（R11）；16602 旧产物基于 payable 旧版，留作历史 |
| USDT 精度 | 合约硬编码 6 | **读 `usdt.decimals()`** | **读链** | 锁定决策④；mock 设 6，正式 USDT 也 6，但合约不假设 |
| ReentrancyGuard | 沿用 `ReentrancyGuardTransient`(transient/cancun) | 换 `ReentrancyGuardUpgradeable`(普通 storage) | **见 §六.4 待确认** | 取决于 Arbitrum 对 transient storage(TSTORE) 的支持，arch-review 重点确认项 |
| 自动发卡触发 | 合约内定时（不可能，EVM 无定时器） | **外部 onlyOracle 批量调用** | **B** | 锁定决策③；合约暴露 `issueMonthlyTrafficCards`，后端定时调 |
| mock USDT | 引第三方 mock | **自写 6-decimals ERC20** | 自写 | OZ ERC20 + 重写 `decimals()` 返回 6 + `mint` 便于测试发币 |

### 2.3 模块划分（受影响合约）
- **改写核心**：`Deposit.sol`、`Payment.sol`
- **修编译 + wiring**：`Oracle.sol`、`TrafficCardNFT.sol`
- **加分账入口**：`ServiceManager.sol`（setPaymentAddress + getter）
- **新增**：`mocks/MockUSDT.sol`（仅测试/测试网部署）
- **不改逻辑**：`FeeManager.sol`（满足 R5）、`UserRegistry.sol`（R10 邮箱）
- **接口同步**：`IDeposit.sol`、`IPayment.sol`、`IServiceManager.sol`、（Oracle 内联接口收敛到 interfaces/）

---

## 三、领域模型与状态机

### 3.1 保证金 → 锁仓 → 自动发卡 → NFT 有效期 状态机

```
[未注册] --register(email)--> [已注册]
   │
[已注册] --approve(usdt, deposit, amount)--> [已授权]
   │
[已授权] --deposit(amount≥10USDT)--> [锁仓中]  (lockExpiry = now+30d；复存叠加+30d)
   │                                      │
   │                          withdraw(): require(now≥lockExpiry) ✘ revert "Lock not expired"
   │
[锁仓中] --now≥lockExpiry--> [可发卡 / 可提取]
   │                              │
   │   issueMonthlyTrafficCards(由 Oracle 批量调，校验到期+有存款+无卡)
   │        --> mint NFT --> [持卡中]   (CardInfo.createdAt=now, isDestroyed=false)
   │                              │
   │                         burn(tokenId) --> [服务激活] (ServiceActivated: now+30d)
   │
[可提取] --withdraw()--> usdt.safeTransfer(principal) --> [已提取] (deposits=0, lockExpiry=0)
```

**关键不变量**：
- `issueMonthlyTrafficCards` 复用 `mintTrafficCard` 的三条校验（`now≥lockExpiry` && `deposits>0` && `userCardCount==0`），保证幂等（同一用户重复调不会重复发卡）。
- 锁仓续期：`deposit()` 时若已到期则重置 `now+30d`，未到期则在原到期点 `+30d`（沿用现状 L46-50，不改）。
- NFT 有效期：burn 后服务 30 天（`ServiceActivated`），`DEDUCTION_VALIDITY=30 days` 常量。

### 3.2 账单 创建 → 支付 → 分账 状态机

```
[无账单] --Oracle.monthlySettlement / createBill(onlyOwner|oracle)--> [未支付账单]
            (platformFee = feeManager.calculateFee(amount); 存 bill)
   │
[未支付] --用户 approve(usdt, payment, amount+fee)--> [已授权]
   │
[已授权] --payBill(billId)-->  校验 !isPaid && bill.user==msg.sender && operator.paymentAddress≠0
   │         ├─ usdt.safeTransferFrom(user, operator.paymentAddress, bill.amount)   ← 主体直分
   │         ├─ usdt.safeTransferFrom(user, platformWallet, bill.platformFee)        ← 手续费
   │         └─ isPaid = true
   │         emit BillPaid(billId, user, total, operatorAmount)
   │
[已支付] (终态)

可选：applyTrafficCardToBill(billId)  ← Oracle 在 monthlySettlement 中调，用流量卡抵扣账单
```

> **分账实现要点**：`operator.paymentAddress` 必须非零（部署/设置时 require）。两段 `safeTransferFrom` 各自从 user 拉款（用户 approve 总额 `amount+fee` 给 Payment），避免合约暂存资金（降低重入面）。

---

## 四、接口与事件定义（变更前后对比）

### 4.1 Deposit

| 项 | 变更前（现状） | 变更后（本轮） | 说明 |
|----|----------------|----------------|------|
| `deposit` | `function deposit() external payable` | `function deposit(uint256 amount) external` | 去 payable；前置 `usdt.safeTransferFrom(msg.sender, address(this), amount)`；`require(amount >= MIN_DEPOSIT)` |
| `withdraw` | `.call{value: principal}` 退原生币 | `usdt.safeTransfer(msg.sender, principal)` | CEI：先清零再转账 |
| `issueMonthlyTrafficCards` | `(address[]) external` **空实现** | `(address[]) external` **真实 mint 逻辑** | `require(msg.sender==oracle)`；遍历 users，对满足「到期+有存款+无卡」者 mint，跳过不满足者（不 revert 整批） |
| 新增 state | — | `IERC20 public usdt;` | **追加在 storage 末尾**（UUPS 兼容）；fresh deploy 无约束 |
| 新增 setter | — | `setUSDT(address) onlyOwner` 或 initialize 注入 | 注入 USDT 地址 |
| `IDeposit.deposit` | `() external payable` | `(uint256 amount) external` | 接口同步；**改 function selector**，冲击后端 ABI |
| 事件 | `DepositMade(user,amount)` 不变 | 同前 | 金额单位语义→USDT（值类型不变，监听方需知精度=decimals） |

> ⚠️ **MIN_DEPOSIT 取值**：不可硬编码常量（USDT 6 位 → 10 USDT = 10_000_000）。设计为 `function _minDeposit() internal view returns(uint256){ return 10 * 10**usdt.decimals(); }` 或 initialize 时按 decimals 计算并存 state。验收 A.2 单测 9.999999 USDT 拒绝 / 10.000000 USDT 通过。

### 4.2 Payment

| 项 | 变更前 | 变更后 | 说明 |
|----|--------|--------|------|
| `payBill` | `(uint256) external payable`，fee→platform，**amount 无出口**，多退少补 | `(uint256) external`（去 payable），`safeTransferFrom(user→operator.paymentAddress, amount)` + `safeTransferFrom(user→platformWallet, fee)` | 链上直分；ERC20 精确金额无找零；用户须先 approve `amount+fee` |
| `applyTrafficCardToBill` | **不存在**（Oracle 已调用 → 编译错误源之一） | `(uint256 billId) external` 新增实现 | 用流量卡额度抵扣账单；权限 onlyOracle/owner；最小实现：标记账单部分/全额抵扣（见 §六编译修复） |
| 新增 state | — | `IERC20 public usdt; IServiceManager public serviceManager;` | 查 operator.paymentAddress 需引 ServiceManager |
| 新增 setter | — | `setUSDT/setServiceManager (onlyOwner)` | |
| `IPayment.payBill` | `(uint256) external payable` | `(uint256) external` | 接口同步 |
| `IPayment` 接口 | 无 `applyTrafficCardToBill` | 新增声明 | 收敛 Oracle 内联接口到 IPayment（单一 SSOT） |
| 事件 `BillPaid` | `(billId,user,totalAmount,operatorAmount)` | 不变（语义对齐：operatorAmount=bill.amount） | |

### 4.3 ServiceManager（分账入口）

| 项 | 变更前 | 变更后 | 说明 |
|----|--------|--------|------|
| `updateOperator` | `(id,name,region,requiredDeposit)` **无 paymentAddress** | 增 `setOperatorPaymentAddress(uint256 id, address addr) onlyOwner` | 倾向**新增独立 setter**（减小 ABI 影响面）而非改 updateOperator 签名；`require(addr != address(0))` |
| 11 内置运营商 | `paymentAddress=address(0)` | 部署后逐个 `setOperatorPaymentAddress` 注入真实地址 | 部署脚本补；测试网可用 deployer 派生地址 |
| `requiredDeposit` | `0.0xx ether`（18 位） | 语义复核：切 USDT 后单位错 → **本轮按 USDT 精度重设或标记不参与本轮资金校验** | scan §七雷区；当前 `_operatorRequiredDeposit` 未被使用，arch-review 复核 |

### 4.4 Oracle（修编译 + wiring）

| 项 | 变更前 | 变更后 | 说明 |
|----|--------|--------|------|
| `deposit`/`payment` state | `address public` L10-11 | `IDeposit public deposit; IPayment public payment;` | **编译错误①根因**：address 类型无 `.createBill` 等方法 |
| `UsageDataSubmitted` 事件 | **未声明**（L71 emit 报错） | 声明 `event UsageDataSubmitted(address indexed user, uint256 operatorId, uint256 dataUsage, uint256 callUsage)` | 实测唯一中止点 |
| `setPayment` | **缺失** | 新增 `setPayment(address) onlyOwner` | scan 遗留：Oracle 无法设 payment |
| 内联 `IPayment`/`IDeposit` L105-122 | 文件尾内联 | 收敛：import `./interfaces/IPayment.sol` + `IDeposit.sol`（统一 SSOT） | 避免双份接口漂移 |

### 4.5 TrafficCardNFT（修编译）

| 项 | 变更前 | 变更后 | 说明 |
|----|--------|--------|------|
| `_deductionCredits` | `mapping(address=>DeductionCredit)` L18，**`DeductionCredit` 未定义（编译错误②）** | 二选一：(a) 定义 `struct DeductionCredit{uint256 amount; uint256 expiresAt;}`；(b) 删除该 dead-code mapping + `onlyDeposit` modifier | **倾向 (b) 删 dead code**（最小改动、scan §④ 确认未被任何函数使用）；若 applyTrafficCardToBill 抵扣需要额度则用 (a) |

### 4.6 ERC20 两段式调用序列（approve / transferFrom）

```
存款流程（Deposit）:
  1. user → USDT.approve(depositAddr, amount)          ← 前端第一步（检测 allowance）
  2. user → Deposit.deposit(amount)
       └─ USDT.safeTransferFrom(user, depositAddr, amount)  ← 需 allowance≥amount

支付流程（Payment）:
  1. user → USDT.approve(paymentAddr, amount+fee)       ← approve 总额
  2. user → Payment.payBill(billId)
       ├─ USDT.safeTransferFrom(user, operator.paymentAddress, amount)
       └─ USDT.safeTransferFrom(user, platformWallet, fee)
```

> 前端（子项 3/3）据此做「allowance 不足→approve→deposit/pay」两步态（PRD D.13）。后端（子项 2/3）监听事件不变，但金额按 decimals 解释；注意 `deposit()`→`deposit(uint256)` 改 selector，ABI 需重生成。

---

## 五、关键技术设计（按合约展开）

### 5.1 Deposit — ERC20 改写 + 锁仓 + 最小额
- import `SafeERC20` + `IERC20`（`@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol`），`using SafeERC20 for IERC20;`
- `initialize` 增 `_usdt` 参数（fresh deploy）或新增 `setUSDT`（兼容升级路径）。**fresh deploy 走 initialize 注入更干净**。
- `deposit(uint256 amount)`：`require(userRegistry.isRegistered(msg.sender))` + `require(amount >= 10*10**usdt.decimals(), "Below min deposit")` → `usdt.safeTransferFrom(msg.sender, address(this), amount)` → 累加 + 锁仓续期（沿用 L46-50）→ emit。
- `withdraw()`：CEI——先 `principal=_deposits[...]; _deposits[...]=0; _lockExpiry[...]=0;` 再 `usdt.safeTransfer(msg.sender, principal)`；加 `nonReentrant`。
- `mintTrafficCard`/`dataAmount=_deposits/100000`（L76）：**精度复核**——6 位下 10 USDT=10_000_000，/100000=100；arch 确认 dataAmount 语义（流量配额）。**建议**：发卡 dataAmount 用 `trafficCardQuota`（固定 100MB，L27）而非按存款额除，避免精度耦合。

### 5.2 Payment — 分账 + 手续费走 FeeManager
- 引 `IServiceManager` 查 `getOperator(operatorId).paymentAddress`。
- `payBill`：去 payable + 去找零；`require(operator.paymentAddress != address(0), "Operator payout unset")`；两段 `safeTransferFrom`（amount→operator，fee→platform）；加 `nonReentrant`（沿用现有 guard 或降级，见 §六.4）。
- 手续费仍在 `createBill` 用 `feeManager.calculateFee(amount)` 存入 `bill.platformFee`（满足 R5，不改）。
- `applyTrafficCardToBill(billId)`：Oracle monthlySettlement 调用。最小实现：onlyOracle/owner，校验账单存在未付，用流量卡额度抵扣（全额则标记 isPaid 并 emit，部分则减额）。**arch-review 重点确认抵扣额度来源与精度**。

### 5.3 TrafficCardNFT — mint/销毁/有效期
- 修编译（§4.5）后逻辑基本不动：`mint`（Deposit 持 ownership 调用）、`burn`（emit ServiceActivated now+30d）、`DEDUCTION_VALIDITY=30 days`。
- 自动发卡走 Deposit→NFT.mint 链路（ownership 已在 deploy 转给 Deposit，L88）。

### 5.4 Oracle — wiring + 自动发卡触发
- 修 state 类型为接口（§4.4）、声明事件、加 setPayment、收敛内联接口。
- `monthlySettlement`（L50-89）：onlyOwner（保留），遍历 createBill + `deposit.issueMonthlyTrafficCards(users)` + `applyTrafficCardToBill`。
- **自动发卡触发链**：后端（子项 2/3）定时 → `Oracle.monthlySettlement(...)` 或直接 `Deposit.issueMonthlyTrafficCards(...)`（onlyOracle）。`Deposit.issueMonthlyTrafficCards` 权限 = `require(msg.sender==oracle)`（沿用 L105），实现真实批量 mint（§5.1 三校验幂等）。

### 5.5 mock USDT（6 decimals）
- `contracts/mocks/MockUSDT.sol`：`ERC20`（OZ 标准）+ `decimals() returns(6)` override + `mint(to, amount)`（public，测试网任意发币）+ symbol "USDT"。
- 仅本地/测试网部署；deploy.ts 部署后写 `deployments/<net>.json` 的 `usdt` 字段 + `usdtDecimals:6`。

### 5.6 编译错误清单与逐个修复方案

| # | 文件:行 | 错误 | 根因 | 修复 |
|---|---------|------|------|------|
| ① | `Oracle.sol:71` | `DeclarationError: UsageDataSubmitted` 未声明 | 事件从未定义（实测唯一中止点） | 在 Oracle 声明 `event UsageDataSubmitted(address indexed,uint256,uint256,uint256)` |
| ② | `Oracle.sol:69/78/85`（修①后暴露） | `payment.createBill` / `deposit.issueMonthlyTrafficCards` / `payment.applyTrafficCardToBill` —— `address` 类型无成员 | `deposit`/`payment` 声明为 `address`（L10-11）却当合约调 | 改 state 为 `IDeposit public deposit; IPayment public payment;`，import 接口 |
| ③ | `Payment.sol`（Oracle 调 `applyTrafficCardToBill`） | Payment 无此函数 | Oracle 内联接口声明了但 Payment 未实现 | Payment 实现 `applyTrafficCardToBill(uint256)`（§5.2） |
| ④ | `TrafficCardNFT.sol:18` | `DeclarationError: DeductionCredit` 未定义 | 类型从未定义（dead code，scan §④ 确认未用） | 删除该 mapping + `onlyDeposit`（倾向）或定义 struct |
| ⑤ | `Oracle.sol` 内联接口 vs interfaces/ | 双份 IPayment/IDeposit 漂移风险 | 历史内联 | 收敛为 import interfaces/，IPayment 补 `applyTrafficCardToBill` 声明 |

> implement 顺序：先 ①②④（纯声明/类型修复）→ 跑 compile 看是否过 → 再 ③⑤（需实现逻辑）。每修一类 `npx hardhat compile` 验证。

---

## 六、非功能性设计

### 6.1 资损 checklist（🔴 = arch-review 安全审计重点）

| 项 | 风险 | 设计对策 |
|----|------|----------|
| 🔴 ERC20 重入 | safeTransferFrom 回调（恶意/非标 token） | CEI：状态先改（isPaid/deposits 清零）再转账；payBill/withdraw 加 `nonReentrant` |
| 🔴 授权额度 | approve 总额 vs 实际扣款不一致、无限授权风险 | 合约只 `safeTransferFrom(精确额)`；前端 approve 精确额（不无限授权） |
| 🔴 精度 | 硬编码 18/6 导致金额错算（scan §七雷区） | 全程 `usdt.decimals()`；`MIN_DEPOSIT=10*10**decimals`；`dataAmount` 解耦存款额（用 quota） |
| 🔴 分账零地址 | operator.paymentAddress=address(0) → USDT 转入黑洞 | payBill `require(paymentAddress != address(0))`；setOperatorPaymentAddress 零地址 require；部署脚本补全 |
| 整数溢出 | `amount+fee`、`*rate/10000` | Solidity 0.8 内置 checked；calculateFee 小额截断为 0 可接受（费率本就向下取整） |
| 找零逻辑移除 | ERC20 无 msg.value，旧多退少补删除 | payBill 精确 transferFrom，无残留资金 |
| 部分扣款失败 | 两段 transferFrom 第一段成功第二段失败 | 同一 tx 原子性，任一 revert 全回滚 |

### 6.2 安全（权限矩阵）

| 函数 | 权限 | 校验 |
|------|------|------|
| `Deposit.deposit` | 任意已注册用户 | `isRegistered` + `amount≥MIN` |
| `Deposit.withdraw` | 持仓用户 | `now≥lockExpiry` + `deposits>0` |
| `Deposit.issueMonthlyTrafficCards` | **onlyOracle** | `msg.sender==oracle` + 逐用户三校验 |
| `Deposit.mintTrafficCard` | onlyOwner | 三校验 |
| `Payment.createBill` | onlyOwner（Oracle 经 owner 路径） | amount>0 |
| `Payment.payBill` | 账单本人 | `!isPaid` + `bill.user==sender` + paymentAddress 非零 |
| `Payment.applyTrafficCardToBill` | onlyOracle/owner | 账单存在未付 |
| `Oracle.monthlySettlement` | onlyOwner | 长度一致 |
| `ServiceManager.setOperatorPaymentAddress` | onlyOwner | 非零地址 |
| `setUSDT/setOracle/setPayment/...` | onlyOwner | |

- SafeERC20 全覆盖；ReentrancyGuard 覆盖 payBill/withdraw。

### 6.3 升级兼容（UUPS storage layout）
- **fresh deploy 421614 + 31337**：无旧 storage 约束，新增 `usdt`/`serviceManager` state 随意排布（推荐放显式声明顺序末尾以备未来升级）。
- **若升级 16602 旧 proxy（可选/非验收项）**：新增 state **只能追加到现有 storage 末尾**，不可插入/改序；Oracle 把 `address`→接口类型属同 slot 大小（address 与接口都占 1 slot），storage 兼容，但 OZ upgrades 会校验，需 `unsafeAllow` 或确认 layout。
- 删除 TrafficCardNFT `_deductionCredits` mapping：mapping 头 slot 占位 1，删除对已部署 proxy 是 storage gap 风险——**16602 若升级需保留占位**；fresh deploy 无虑。🔴 arch-review 确认。

### 6.4 Arbitrum 兼容确认项（🔴 implement 前必须验证）
| 项 | 现状 | 风险 | 确认动作 |
|----|------|------|----------|
| 🔴 `evmVersion: cancun` + `viaIR` | hardhat.config L15-16 | Arbitrum Sepolia 对 cancun opcode 支持？ | 部署前在 421614 跑一笔交易验证；Arbitrum Nitro 已支持多数 cancun，但 **transient storage(TSTORE/TLOAD)** 需重点确认 |
| 🔴 `ReentrancyGuardTransient` | Payment L11（依赖 TSTORE） | 若 Arbitrum 不支持 TSTORE → Payment 部署/调用失败 | 二选一：(a) 确认 Arbitrum 支持后保留；(b) 降级 `ReentrancyGuardUpgradeable`（普通 storage，需改 initialize + 删自定义 `_reentrancyGuardInit` assembly L33-38）。**建议 implement 先在 421614 实测 Payment 部署+payBill** |
| 自定义 `_reentrancyGuardInit` assembly | Payment L33-38 | 若换普通 guard 此段须删 | 跟随 (a)/(b) 决策 |
| chainId/RPC | 无 421614 | — | hardhat.config 加网络 |

### 6.5 异常与监控
- 事件不变（DepositMade/BillCreated/BillPaid/CardMinted/ServiceActivated）+ 新增 `UsageDataSubmitted`；后端 event_sync 据此监听（子项 2/3）。
- revert 文案统一英文短串（"Below min deposit"/"Operator payout unset"/"Lock not expired"）。

---

## 七、部署设计

### 7.1 hardhat.config.ts
新增网络（不动现有 31337/16600/16602）：
```ts
arbitrum_sepolia: {
  url: process.env.ARBITRUM_SEPOLIA_RPC || "https://sepolia-rollup.arbitrum.io/rpc",
  chainId: 421614,
  accounts: [DEPLOYER_KEY],
  timeout: 60000,
}
// 可选：显式 localhost(31337) 指向 http://127.0.0.1:8545（deploy:local 用）
```
- `evmVersion`/`viaIR` 维持，待 §6.4 确认；若 transient 不支持则全局或 Payment 单独降级。

### 7.2 deploy.ts 增量
顺序（在现有 7 合约 deployProxy 基础上）：
```
0. 部署 MockUSDT(decimals=6)  ← 测试网/本地；本轮无正式 USDT，全用 mock
... 现有 1-7 proxy 部署 ...
wiring 补全（现状缺失项）：
  + deposit.setUSDT(usdtAddr)          // 或 initialize 注入
  + payment.setUSDT(usdtAddr)
  + payment.setServiceManager(serviceManagerAddr)
  + payment.setOracle(oracleAddr)      // ← 现状缺（scan baseline §二）
  + oracle.setPayment(paymentAddr)     // ← 现状缺（Oracle 无 setPayment，本轮新增）
  + 循环 serviceManager.setOperatorPaymentAddress(id, payoutAddr)  // 11 运营商，非零
  （现有）trafficCardNFT.setDepositContract / deposit.setTrafficCardNFT / deposit.setOracle / oracle.setDeposit / transferOwnership(deposit)
```

### 7.3 deployments 产出
`deployments/<net>.json` 增字段：
```json
{ "chainId": 421614, "proxies": {...7...},
  "implementations": {...7...},
  "usdt": "0x...", "usdtDecimals": 6 }
```
两处部署：`deploy:arbitrum-sepolia`（421614）+ `deploy:local`（31337）。供后端 configs/deployments.json（子项 2/3）+ 前端 contracts.ts（子项 3/3）对齐。

---

## 八、补测策略（验收 A.2 等）

新增/重写测试（建议新建 `test/erc20.ts` 或扩 `linkworld.ts`；TDD 优先）：

| 编号 | 用例 | 验收点 |
|------|------|--------|
| MIN-01 | deposit 9.999999 USDT → revert "Below min deposit" | **A.2 <10 拒绝** |
| MIN-02 | deposit 10.000000 USDT → 成功 | **A.2 =10 通过** |
| ERC-01 | 未 approve → deposit revert（ERC20InsufficientAllowance） | approve 前置 |
| ERC-02 | approve 后 deposit → `_deposits` 增、USDT 余额转移正确 | transferFrom |
| ERC-03 | withdraw 锁仓未到 → revert；到期后 → USDT 退回本金 | 锁仓 + safeTransfer |
| PAY-01 | payBill 未 approve → revert | 支付授权 |
| PAY-02 | payBill 成功 → operator.paymentAddress 收 amount、platformWallet 收 fee、isPaid=true | **分账正确性** |
| PAY-03 | operator.paymentAddress=0 → payBill revert "Operator payout unset" | 零地址校验 |
| PAY-04 | fee = `calculateFee(amount)` 与链上 FeeManager 一致 | R5 |
| ISS-01 | 非 oracle 调 issueMonthlyTrafficCards → revert "Only oracle" | **自动发卡权限** |
| ISS-02 | oracle 调，锁仓到期用户 → mint NFT；未到期用户 → 跳过不 revert | **自动发卡时序/幂等** |
| ISS-03 | 已有卡用户重复 issue → 不重复 mint | 幂等 |
| DEC-01 | mock USDT.decimals()=6；MIN_DEPOSIT 随 decimals 计算 | **精度（不硬编码）** |
| UPG-01（可选） | upgradeProxy 后 storage 不丢（若走升级路径） | UUPS 兼容 |
| 回归 | 现有 FE/UR/SM/TC 26 it（金额改 USDT 精度后）不回归 | 不破坏 |

> 现有测试金额从 `parseEther` 改 USDT 精度辅助：`const usdt=(n)=>BigInt(n)*10n**6n`。FeeManager 的 FE-02/04 用 ETH 计费率与币种无关，可保留或改 USDT 单位。测试需先 mint mock USDT 给 user 并 approve。

---

## 九、落地计划（给 plan 阶段输入）

| 阶段 | 任务 | 顺序 | 风险点 |
|------|------|------|--------|
| P1 修编译 | Oracle 事件+类型+setPayment+收敛接口；TrafficCardNFT 删 DeductionCredit；Payment 实现 applyTrafficCardToBill | **最先，串行** | 改 Oracle state 类型可能触发 OZ upgrades layout 校验（fresh deploy 无虑）；compile 绿是 gate |
| P2 ERC20 迁移 | Deposit/Payment 引 SafeERC20 + usdt state + deposit(amount)/withdraw 改写 + MIN_DEPOSIT | P1 后 | 精度雷区（dataAmount/100000）；用 decimals 不硬编码 |
| P3 分账 | Payment payBill 直分 + ServiceManager setOperatorPaymentAddress + 零地址校验 | P2 后 | operator.paymentAddress 全 0 必须先补；两段 transferFrom 原子性 |
| P4 自动发卡 | Deposit.issueMonthlyTrafficCards 真实 mint + onlyOracle + 幂等 | P3 后 | 时序校验、批量 gas |
| P5 部署 | hardhat.config 421614 + MockUSDT + deploy.ts wiring 补全 | P4 后 | 🔴 Arbitrum cancun/transient 实测（可能回头改 Payment guard） |
| P6 补测 | §八 测试清单 | 贯穿（TDD 可前置到各阶段） | 现有 26 it 不回归 |

**implement 串行铁律**：P1→P6 严格串行（一个完成审查再下一个），尤其 P1 compile gate 不过不进 P2。P5 的 Arbitrum 兼容实测若失败，需回退到「Payment ReentrancyGuard 降级」分支（§6.4）。

---

## 十、🔴 arch-review / security-review 重点审清单

1. **ERC20 重入 + CEI**：payBill/withdraw 状态先改后转、nonReentrant 是否到位。
2. **分账零地址**：operator.paymentAddress 非零校验是否覆盖 payBill + setter + 部署全路径。
3. **精度不硬编码**：全仓 grep 无字面 `* 10**18` / `* 10**6` / `1e18`；MIN_DEPOSIT 与 dataAmount 计算正确。
4. **授权额度**：前端 approve 精确额（非无限），合约只拉精确额。
5. **自动发卡权限/幂等**：onlyOracle + 三校验 + 重复调不重复发。
6. **Arbitrum cancun/transient storage 兼容**：ReentrancyGuardTransient 是否在 421614 可用，否则降级方案。
7. **UUPS storage layout**：fresh deploy 无虑；若升级 16602，Oracle state 类型变更 + TrafficCardNFT mapping 删除的 layout 影响。
8. **接口 SSOT 收敛**：Oracle 内联 IPayment/IDeposit 与 interfaces/ 不漂移；selector 变更（deposit/payBill 去 payable）对后端 ABI 的冲击已通知子项 2/3。
9. **ServiceManager requiredDeposit 单位**：18 位 ether 字面量在 USDT 语境下的语义（当前 `_operatorRequiredDeposit` 未被使用，确认是否引入资金校验）。
