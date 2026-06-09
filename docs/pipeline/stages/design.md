# Stage: design — 合约技术设计（子项目 contracts 1/3）

> **状态**: v2（按 arch-review 返工） | **版本**: **v2** | **日期**: 2026-06-08 | **Gate**: 1（设计） | **子项目**: contracts(1/3) | **角色**: 架构师
> **v2 说明**：本版按 arch-review（docs/pipeline/stages/arch-review.md）返工，**逐条闭合 7 大阻塞 B1–B7 + §三 14 项消歧/加固项 + 用户新拍板 3 决策**。变更可追溯见 §十一「arch-review 阻塞闭合对照表」。
> 输入：requirement.md（PRD §三①/§五A·D/§六）+ scan.md（编译不过 + 7 发现）+ 4 份基线（project-scan / contracts-inventory / erc20-migration-surface / test-deploy-baseline）+ arch-review.md（B1–B7）。
> 范围：**只产设计**，不改任何 .sol / 脚本 / 测试 / 配置（实现留 implement 阶段）。
> 注：本文件取代 2026-06-07 旧 design.md（web 视觉版，属 web 子项 3/3，已因全栈反转重新分项）。

### 锁定决策（用户已拍板，不再翻）
**v1 原有 5 项**：
- ① 全补齐不留 stub（**例外见 v2-B**：applyTrafficCardToBill 本轮明确冻结为受限桩）；
- ② 运营商分账=链上直分（amount→operator.paymentAddress，fee→platform）；
- ③ 自动发卡入口由后端/Oracle 定时调用；
- ④ mock USDT 6 位精度，合约一律读 `usdt.decimals()` 不硬编码；
- ⑤ 部署 Arbitrum Sepolia(421614) + 本地 hardhat(31337)，不上主网。

**v2 新增 3 项（arch-review 后用户拍板，最高优先级，全文口径以此为准）**：
- **v2-A｜计价模型（闭合 B1，🔴Critical）**：账单金额由**后端/Oracle 链下按资费算好后**作为已计价的 USDT 金额传入 `createBill`；**合约不做任何 usage 求和/换算**——`Oracle.monthlySettlement` **删除 `totalAmount = dataUsage + callUsage`** 逻辑，改为传入参数 `amounts[]`（已是 USDT 金额）。Oracle 角色收敛为「喂价/触发」，**不做计价**。
- **v2-B｜applyTrafficCardToBill = 受限桩 + 语义冻结**：本轮**只实现为权限受控的桩**——`onlyOracle`，**不转移任何资金**，仅做账单存在性校验 + 可选 `emit` 事件，让 Oracle 调用链与编译走通即可。真实「流量卡抵扣账单」语义（额度来源、抵扣范围、精度）**明确冻结到后续 Round**；本轮补测面/资损面**移除其抵扣逻辑**，只测桩的权限与存在性校验。
- **v2-C｜发卡额度（闭合 B4）**：`dataAmount = 固定 trafficCardQuota`（Deposit L17 已有 state，initialize 设 100MB），**与存款额、与精度彻底解耦**；**删除 `mintTrafficCard` L76 的 `_deposits[user] / 100000`**，全文（§一/§3.1/§4.1/§5.1/§6.1）统一为此口径。

---

## 〇、已核对真实代码清单（含行号引用）

本设计基于逐行核对，引用以下源码（均已 Read）：

| 文件 | 关键行号 / 核对结论 |
|------|---------------------|
| `contracts/Deposit.sol` | `deposit() payable` L41-53；`withdraw() .call{value}` L56-68；`mintTrafficCard onlyOwner` L71-81，`dataAmount=_deposits[user]/100000` L76（**B4 根因·v2-C 改为固定 trafficCardQuota**）；`trafficCardQuota` state L17 + initialize 设 100MB L27；`issueMonthlyTrafficCards(onlyOracle)` L104-107 **空实现 / 注释占位**（B3：复用 onlyOwner `mintTrafficCard` 会 revert）；state L13-21（无 usdt）；`_operatorRequiredDeposit` L21 **从未被读** |
| `contracts/Payment.sol` | `payBill() payable` L76-97；fee→platformWallet L86-89，**amount 无出口**；**`createBill` 是 `onlyOwner` L55**（已定义 `onlyOracle` modifier L20-23 但 createBill 未用 → **B2 根因**）；`ReentrancyGuardTransient` L11 + 自定义 `_reentrancyGuardInit` assembly L33-38（向 `_reentrancyGuardStorageSlot()` `sstore(slot,1)` → **B2 消歧·正确性需厘清**）；**无 `applyTrafficCardToBill`**（Oracle L85 调用 → 编译错误③）；**无 `setUSDT`** |
| `contracts/FeeManager.sol` | `FEE_DENOMINATOR=10000` L10、`MAX_FEE_RATE=1000` L11、`calculateFee=amount*rate/10000` L34-36；满足 R5，本轮不改逻辑 |
| `contracts/TrafficCardNFT.sol` | `mint onlyOwner` L45-62；`mintBatch onlyOwner` L65-80 **L77 `this.mint()` 外部调用会撞 onlyOwner → 自动发卡禁用**；`burn` L83-94 emit `ServiceActivated(now+30days)`；**L18 `DeductionCredit` 类型未定义（编译错误②/④）**；`onlyDeposit` L27-30 / `_deductionCredits` L18 dead code |
| `contracts/Oracle.sol` | `deposit`/`payment` 声明为 `address` L10-11 却调 `.createBill`/`.issueMonthlyTrafficCards`/`.applyTrafficCardToBill`（**编译错误①根因**）；`UsageDataSubmitted` L71 **未声明（实测唯一报错点）**；**L68 `totalAmount = dataUsage + callUsage` → B1 量纲错误根因**；`monthlySettlement onlyOwner` L55；**无 `setPayment`**；文件尾内联 `IPayment`/`IDeposit` L105-122（IPayment 内联含 `applyTrafficCardToBill` L117，但 interfaces/IPayment.sol 无 → B5 漂移） |
| `contracts/ServiceManager.sol` | 11 内置运营商 L21-153，`requiredDeposit` 用 `0.0xx ether`（18 位）；**全部 `paymentAddress=address(0)`**；`addOperator` 已带 paymentAddress 入参 L156-177，`updateOperator` L179-192 **无 paymentAddress 入参** |
| `contracts/UserRegistry.sol` | 邮箱注册 L22-37，无 phone（对齐 R10）；本轮不改 |
| `contracts/interfaces/*.sol` | `IDeposit.deposit() payable` L9，`IDeposit` **无 `issueMonthlyTrafficCards`（B5）**；`IPayment.payBill() payable` L19，`IPayment` **无 `applyTrafficCardToBill`（B5）**；`IServiceManager.Operator.paymentAddress` 已有 L12 |
| `scripts/deploy.ts` | UUPS 全量 L9-89；现有 wiring L80-84（`trafficCardNFT.setDepositContract` / `deposit.setTrafficCardNFT` / `deposit.setOracle` / `payment.setPlatformWallet` / `oracle.setDeposit`）；`trafficCardNFT.transferOwnership(deposit)` L88（NFT owner→Deposit）；**缺 `payment.setOracle`、`oracle.setPayment`、operator.paymentAddress、mock USDT、`payment.setServiceManager`、`*.setUSDT`**；**Payment owner 仍是 deployer，Oracle 非 Payment owner → B2：`monthlySettlement→createBill(onlyOwner)` 必 revert** |
| `hardhat.config.ts` | solidity 0.8.27 / `evmVersion:"cancun"` L15 / `viaIR:true` L16；networks 仅 31337/16600/16602 L19-35，**无 421614** |
| `test/linkworld.ts` | 直接 `Factory.deploy()`+手动 `initialize`（不走 proxy）；全用 `parseEther`（18 位）；**无 withdraw/payBill/mintTrafficCard/最小额/升级 测试** |

**实跑 `npx hardhat compile` 结论**：编译在 `Oracle.sol:71 UsageDataSubmitted` 第一个 DeclarationError 处中止（编译器只报首错）。修掉它后会暴露后续错误（Oracle address-vs-contract、TrafficCardNFT `DeductionCredit`），见 §六编译错误清单。

---

## 一、背景与目标

### 1.1 为什么做这次合约改造
合并自 origin 的合约（commit `7ef9677`）是**编译不过的半成品**，且 PRD Round1 全栈反转后确定把资金通道从原生币切到 ERC20 USDT、目标链切到 Arbitrum Sepolia。本轮合约层（子项 1/3）是后端（2/3）/web（3/3）的地基：合约不编译、ABI 不稳、地址不产出，后两子项无法对齐。

### 1.2 本轮范围边界
**做**：① 修编译（4 类错误）；② Deposit/Payment ERC20 USDT 改写（approve+transferFrom、SafeERC20、最小额 10 USDT）；③ Payment 链上直分（amount→operator.paymentAddress，fee→platform）；④ TrafficCardNFT mint/burn/有效期 wiring；⑤ Oracle wiring（setPayment）+ **计价收敛（v2-A：createBill 接收链下算好的 USDT 金额，Oracle 不求和）** + `issueMonthlyTrafficCards` 真实自动发卡逻辑（走独立 `_mintFor`，**v2-C 固定 quota 发卡**）+ `applyTrafficCardToBill` **受限桩（v2-B，不转资金）**；⑥ mock USDT(6 decimals) 自部署；⑦ hardhat.config 加 421614 + deploy.ts 补 wiring（含 B2 ownership/授权拓扑）；⑧ 补测（验收 A.2 + B6 集成测试 + B7 非标 USDT）。

**不做（红线）**：❌ 不上 Arbitrum One 主网；❌ 不加 phone 注册（R10）；❌ 不改业务路由/前端/后端；❌ 不引入新业务功能（实体 SIM 等留后续 Round）；❌ **流量卡抵扣账单真实语义**（v2-B 冻结到后续 Round，本轮只做桩）；❌ **运营商保证金校验**（`_operatorRequiredDeposit` 本轮明确废弃，不引入资金校验，见 §二消歧）；❌ **不升级任何旧 proxy**（16602/localhost 一律 fresh deploy，删除一切 storage layout 升级讨论的无效负担，见 §6.3）。

### 1.3 技术指标
| 指标 | 目标 |
|------|------|
| 编译 | `hardhat compile` 零 error（当前 1 error 中止）。**P1 compile gate 重定义（消歧）：不是「①②④后绿」，而是「§5.6 的 ①–⑤ 全部子项做完后 compile 绿」**——因为 ②Oracle 类型改完后仍缺接口方法（B5）、③/桩未实现前编译不过 |
| 单测 | 新增覆盖验收 A.2（<10 拒绝/=10 通过）、approve/transferFrom、分账正确性、自动发卡权限/时序、锁仓提取、精度；**B6：`monthlySettlement` 端到端集成测试 + `applyTrafficCardToBill` 桩权限测试**；**B7：非标 USDT（transfer 无返回值/返回 false）SafeERC20 分支**；现有 26 it 不回归（金额改 USDT 精度后重写） |
| 部署 | 31337 + 421614 两处 `deploy.ts` 成功，产出 `deployments/<net>.json`（含 7 proxy + mock USDT + decimals） |
| 精度 | 合约内零硬编码 18/6，`MIN_DEPOSIT=10*10**usdt.decimals()` |
| 资损 | SafeERC20 全覆盖、分账零地址 require、CEI + ReentrancyGuard、整数运算无溢出/截断隐患 |

---

## 二、整体方案

### 2.1 总体顺序（5 大模块，强依赖拓扑）

```
① 修编译（地基，必须最先；4 类错误见 §5.6 ①–⑤）
   └─ Oracle: address→接口类型 + 声明 UsageDataSubmitted 事件 + 加 setPayment
   │          + 【v2-A】删 totalAmount=dataUsage+callUsage，monthlySettlement 增 amounts[] 入参
   │          + 接口收敛到 interfaces/（依赖 B5：IDeposit/IPayment 先补声明，否则收敛后编译不过）
   └─ TrafficCardNFT: 删除 DeductionCredit + onlyDeposit dead code
   └─ Payment: createBill 改 onlyOracle（B2）+ 实现 applyTrafficCardToBill 受限桩（v2-B，不转资金）
   └─ interfaces/: IDeposit 补 issueMonthlyTrafficCards、IPayment 补 applyTrafficCardToBill（B5）
        ↓ ①–⑤ 全做完后 compile 应绿（P1 gate，消歧重定义）
② ERC20 迁移（Deposit/Payment 资金通道 + SafeERC20 + usdt state(initialize 注入) + 最小额）
   └─ 【B7】仅支持标准 ERC20，禁 fee-on-transfer；usdt 锁定 initialize 注入（弃 setUSDT 后置）
        ↓
③ 分账（Payment.payBill 链上直分 amount→operator.paymentAddress, fee→platform）
   └─ 依赖 ServiceManager 暴露 paymentAddress + 设置入口 + 零地址校验
   └─ 【B7】分账地址非黑名单要求 + 支付失败降级；【B2】createBill fail-fast 校验 paymentAddress
        ↓
④ 自动发卡（Deposit.issueMonthlyTrafficCards → 独立 internal _mintFor，不经 onlyOwner mintTrafficCard）
   └─ 【B3】批量循环 if(!ok) continue 跳过不 revert 整批；【v2-C】dataAmount=固定 trafficCardQuota
   └─ 【消歧】批量分批上限策略（调用方分批 ≤N 或合约 require maxBatch）+ 压测
        ↓
④.5 【消歧·前移】ReentrancyGuardTransient 在 421614 实测（原 §6.4 从 P5 前移到 P2 之前/与 P1 并行验证）
        ↓
⑤ 部署（hardhat.config 加 421614 + mock USDT + deploy.ts 补全 wiring + B2 ownership/授权拓扑）
        ↓
⑥ 补测（贯穿，TDD 优先；含 B6 集成测试 + B7 非标 USDT + 回归重写）
```

> 顺序铁律：**①修编译必须最先**。当前代码连 `compile` 都过不了，任何 ERC20/分账/发卡逻辑都无法验证。implement 阶段先把 ① 的 5 个子项全做完跑到 `compile` 绿，再动 ②③④。
> **B2 关键**：步骤 ⑤ 必须落地「部署后 ownership/授权拓扑」（见新增 §七.0），否则 `monthlySettlement` 在 deploy 下跑不通。

### 2.2 关键选型对比

| 决策点 | 选项 A | 选项 B | **本轮选定** | 理由 |
|--------|--------|--------|-------------|------|
| ERC20 安全转账 | 裸 `IERC20.transfer/transferFrom` | **SafeERC20** | **B** | OZ 5.6.1 自带，处理无返回值/返回 false 的非标准代币；资损敏感必须用 |
| Arbitrum 落地方式 | `upgradeProxy` 升级 16602 旧 proxy | **fresh deploy 全新部署 421614** | **B** | 421614 是新链无旧 proxy；fresh deploy 无 storage layout 约束，最干净（scan §五 已确认 421614 零配置） |
| 16602/localhost 旧 proxy | 升级 16602 旧 proxy | **一律 fresh deploy，不升级** | **fresh deploy** | **消歧拍板**：验收只认 31337+421614（R11）；16602 旧产物基于 payable 旧版，**本轮不升级、不讨论 storage layout 迁移**（删除该无效负担，§6.3 同步精简） |
| USDT 精度 | 合约硬编码 6 | **读 `usdt.decimals()`** | **读链** | 锁定决策④；mock 设 6，正式 USDT 也 6，但合约不假设 |
| usdt 注入路径 | initialize 注入 | setUSDT 后置 setter | **initialize 注入** | **消歧拍板**：fresh deploy 无升级包袱；避免 usdt 未设时 deposit `safeTransferFrom(address(0))` panic |
| ReentrancyGuard | 沿用 `ReentrancyGuardTransient`(transient/cancun) | 换 `ReentrancyGuardUpgradeable`(普通 storage) | **先实测后定（§6.4）** | 取决于 Arbitrum 对 TSTORE 支持；**实测前移到 P2 之前**（消歧），避免 P5 才发现需返工改 guard |
| createBill 权限 | 沿用 `onlyOwner`（B2 根因，monthlySettlement 必 revert） | **改 `onlyOracle`** | **`onlyOracle`** | **B2 闭合**：Payment 已有 `onlyOracle` modifier L20-23，createBill 改用它即可；Oracle 经 `payment.setOracle(oracle)` wiring 获权（比 transferOwnership 更精确、影响面小） |
| 计价归属 | 合约链上求和 dataUsage+callUsage | **链下算好 USDT 金额传入 createBill** | **链下（v2-A）** | **B1 闭合**：量纲错误（字节+分钟≠USDT）；Oracle 只触发不计价 |
| applyTrafficCardToBill | 本轮实现真实抵扣 | **受限桩，冻结语义** | **受限桩（v2-B）** | **B6 关联**：抵扣额度/范围/精度未定，本轮只让编译+调用链走通，不动资金 |
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
- **【B3】`issueMonthlyTrafficCards` 不复用 `mintTrafficCard`（onlyOwner，会 revert）**，改走独立 internal `_mintFor(user)`：内含与 mintTrafficCard 相同的三条校验（`now≥lockExpiry` && `deposits>0` && `userCardCount==0`）+ 调 `trafficCardNFT.mint(...)`。`mintTrafficCard(onlyOwner)` 重构为 `_mintFor` 的 onlyOwner 外壳，避免逻辑两份。
- **【B3】批量循环幂等 + 不整批回滚**：`issueMonthlyTrafficCards` 遍历 users，每个 user 的三校验失败时 `continue` 跳过（**不 `require` 整批 revert**）；同一用户重复调因 `userCardCount==0` 校验保证不重复发卡。
- **【v2-C】`dataAmount = trafficCardQuota`（固定，Deposit L17 state，initialize 设 100MB）**，**删除 `_deposits[user]/100000`**，与存款额/精度完全解耦。
- 锁仓续期：`deposit()` 时若已到期则重置 `now+30d`，未到期则在原到期点 `+30d`（沿用现状 L46-50，不改）。
- NFT 有效期：burn 后服务 30 天（`ServiceActivated`），`DEDUCTION_VALIDITY=30 days` 常量。

### 3.2 账单 创建 → 支付 → 分账 状态机

```
后端链下按资费算好 amounts[i]（已是 USDT 金额）         ← v2-A：合约不再求和
   │
[无账单] --Oracle.monthlySettlement(users,operatorIds,amounts[])--> createBill(onlyOracle) --> [未支付账单]
            (createBill 内：require(amount>0); platformFee = feeManager.calculateFee(amount); 存 bill)
            (createBill fail-fast：建议同时校验该 operator.paymentAddress≠0，避免生成永远付不了的账单 ← 消歧)
   │
[未支付] --用户 approve(usdt, payment, amount+fee)--> [已授权]
   │
[已授权] --payBill(billId)-->  校验 !isPaid && bill.user==msg.sender && operator.paymentAddress≠0
   │         ├─ isPaid = true                                                     ← CEI：状态先改
   │         ├─ usdt.safeTransferFrom(user, operator.paymentAddress, bill.amount)   ← 主体直分
   │         └─ if(bill.platformFee>0) usdt.safeTransferFrom(user, platformWallet, bill.platformFee)  ← 0-fee 跳过（消歧）
   │         emit BillPaid(billId, user, total, operatorAmount)
   │
[已支付] (终态)

桩：applyTrafficCardToBill(billId)  ← v2-B：本轮受限桩（onlyOracle，不转资金），
                                       仅 require 账单存在 + 可选 emit；Oracle 在 monthlySettlement 中调用以走通调用链。
                                       真实「流量卡抵扣账单」语义冻结到后续 Round。
```

> **分账实现要点**：`operator.paymentAddress` 必须非零（部署/设置时 require + payBill 时 require + createBill fail-fast 校验）。两段 `safeTransferFrom` 各自从 user 拉款（用户 approve 总额 `amount+fee` 给 Payment），避免合约暂存资金（降低重入面）。CEI：先置 `isPaid=true` 再转账。`if(platformFee>0)` 保留（0-fee 跳过，沿用现状 Payment L86 判断）。

---

## 四、接口与事件定义（变更前后对比）

### 4.1 Deposit

| 项 | 变更前（现状） | 变更后（本轮） | 说明 |
|----|----------------|----------------|------|
| `deposit` | `function deposit() external payable` | `function deposit(uint256 amount) external` | 去 payable；前置 `usdt.safeTransferFrom(msg.sender, address(this), amount)`；`require(amount >= MIN_DEPOSIT)` |
| `withdraw` | `.call{value: principal}` 退原生币 | `usdt.safeTransfer(msg.sender, principal)` | CEI：先清零再转账 |
| `issueMonthlyTrafficCards` | `(address[]) external` **空实现** | `(address[]) external` **真实 mint 逻辑** | `require(msg.sender==oracle)`；遍历 users，**走独立 internal `_mintFor(user)`（B3，不经 onlyOwner `mintTrafficCard`）**；对满足「到期+有存款+无卡」者 mint，**不满足者 `continue` 跳过（B3，不 revert 整批）**；`dataAmount=trafficCardQuota`（v2-C） |
| `mintTrafficCard` | `(address) onlyOwner`，`dataAmount=_deposits/100000` L76 | `(address) onlyOwner` 重构为 `_mintFor(user)` 的薄壳，`dataAmount=trafficCardQuota`（v2-C） | **删 `_deposits/100000`**；与批量发卡共用 `_mintFor`，逻辑单一 |
| 新增 internal | — | `_mintFor(address user) internal returns(uint256)` | 三校验 + `trafficCardNFT.mint(user, trafficCardQuota, uri)`；被 `mintTrafficCard`（onlyOwner）与 `issueMonthlyTrafficCards`（onlyOracle）复用 |
| 新增 state | — | `IERC20 public usdt;` | **追加在 storage 末尾**（UUPS 兼容）；fresh deploy 无约束 |
| 注入 USDT | — | **`initialize(_userRegistry, _usdt)` 注入（消歧锁定，弃 setUSDT 后置）** | 避免 usdt 未设时 `safeTransferFrom(address(0))` panic；fresh deploy 无升级包袱 |
| `IDeposit.deposit` | `() external payable` | `(uint256 amount) external` | 接口同步；**改 function selector**，冲击后端 ABI |
| **`IDeposit` 新增声明** | **无 `issueMonthlyTrafficCards`（B5）** | **补 `function issueMonthlyTrafficCards(address[] calldata) external;`** | **B5 闭合**：否则 §5.6⑤ 接口收敛后 Oracle import IDeposit 调用编译不过 |
| 事件 | `DepositMade(user,amount)` 不变 | 同前 | 金额单位语义→USDT（值类型不变，监听方需知精度=decimals） |

> ⚠️ **MIN_DEPOSIT 取值**：不可硬编码常量（USDT 6 位 → 10 USDT = 10_000_000）。设计为 `function _minDeposit() internal view returns(uint256){ return 10 * 10**usdt.decimals(); }` 或 initialize 时按 decimals 计算并存 state。验收 A.2 单测 9.999999 USDT 拒绝 / 10.000000 USDT 通过。
> ⚠️ **【v2-C】dataAmount 全文统一**：发卡额度 = `trafficCardQuota`（固定 100MB），**不再有 `_deposits/100000`**。§一表口径、§3.1、§5.1、§6.1 已同步。

### 4.2 Payment

| 项 | 变更前 | 变更后 | 说明 |
|----|--------|--------|------|
| `createBill` | `(user,operatorId,amount) onlyOwner` L55 | `(user,operatorId,amount) onlyOracle` | **B2 闭合**：改用 Payment 已有 `onlyOracle` modifier（L20-23）；amount 是**链下算好的 USDT 金额**（v2-A，Oracle 不再求和）；**fail-fast：建议 require `serviceManager.getOperator(operatorId).paymentAddress != 0`**（消歧，避免生成永付不了的账单） |
| `payBill` | `(uint256) external payable`，fee→platform，**amount 无出口**，多退少补 | `(uint256) external`（去 payable），CEI 先置 isPaid，再 `safeTransferFrom(user→operator.paymentAddress, amount)` + `if(fee>0) safeTransferFrom(user→platformWallet, fee)` | 链上直分；ERC20 精确金额无找零；用户须先 approve `amount+fee`；0-fee 跳过（消歧） |
| `applyTrafficCardToBill` | **不存在**（Oracle 已调用 → 编译错误③） | **`(uint256 billId) external onlyOracle` 受限桩（v2-B）** | **不转移任何资金**；仅 `require` 账单存在（`billId < _nextBillId` 或对应 bill.user≠0）+ 可选 `emit TrafficCardApplied(billId)`；真实抵扣语义冻结到后续 Round |
| 新增 state | — | `IERC20 public usdt; IServiceManager public serviceManager;` | 查 operator.paymentAddress 需引 ServiceManager |
| 注入路径 | — | **`initialize(_feeManager,_platformWallet,_usdt,_serviceManager)` 注入（消歧锁定）** 或保留 setOracle 后置（Oracle 地址 deploy 时才有） | usdt/serviceManager 走 initialize；oracle 仍走 `setOracle`（部署顺序所限，B2 拓扑详见 §七.0） |
| `IPayment.payBill` | `(uint256) external payable` | `(uint256) external` | 接口同步 |
| **`IPayment` 新增声明** | **无 `applyTrafficCardToBill`（B5）** | **补 `function applyTrafficCardToBill(uint256 billId) external;`** | **B5 闭合**：收敛 Oracle 内联接口到 interfaces/IPayment（单一 SSOT），否则收敛后编译不过 |
| 事件 `BillPaid` | `(billId,user,totalAmount,operatorAmount)` | 不变（语义对齐：operatorAmount=bill.amount） | |
| 新增事件（可选） | — | `event TrafficCardApplied(uint256 indexed billId)` | v2-B 桩的可选 emit；后端据此知 Oracle 已触发（本轮无资金语义） |

### 4.3 ServiceManager（分账入口）

| 项 | 变更前 | 变更后 | 说明 |
|----|--------|--------|------|
| `updateOperator` | `(id,name,region,requiredDeposit)` **无 paymentAddress** | 增 `setOperatorPaymentAddress(uint256 id, address addr) onlyOwner` | 倾向**新增独立 setter**（减小 ABI 影响面）而非改 updateOperator 签名；`require(addr != address(0))` |
| 11 内置运营商 | `paymentAddress=address(0)` | 部署后逐个 `setOperatorPaymentAddress` 注入真实地址 | 部署脚本补；测试网可用 deployer 派生地址 |
| `requiredDeposit` | `0.0xx ether`（18 位，initialize L26/38/50…） | **本轮明确废弃：不引入任何运营商保证金校验/标记，字段保留但不读不校验** | **消歧拍板**：scan §七雷区；`_operatorRequiredDeposit`（Deposit L21）从未被读，本轮不动它、不引入资金校验，避免 USDT 语境下 18 位字面量语义错引发误用 |

### 4.4 Oracle（修编译 + wiring）

| 项 | 变更前 | 变更后 | 说明 |
|----|--------|--------|------|
| `deposit`/`payment` state | `address public` L10-11 | `IDeposit public deposit; IPayment public payment;` | **编译错误①根因**：address 类型无 `.createBill` 等方法 |
| `UsageDataSubmitted` 事件 | **未声明**（L71 emit 报错） | 声明 `event UsageDataSubmitted(address indexed user, uint256 operatorId, uint256 dataUsage, uint256 callUsage)` | 实测唯一中止点 |
| **`monthlySettlement` 签名** | `(users,operatorIds,dataUsages,callUsages)`，L68 `totalAmount=dataUsage+callUsage` | **`(users,operatorIds,amounts[])`**，直接 `payment.createBill(user,operatorId,amounts[i])` | **B1/v2-A 闭合**：删除 L68 求和；`amounts[i]` 是后端链下算好的 USDT 金额；Oracle 不计价。`dataUsages/callUsages` 可保留为 `UsageDataSubmitted` 事件的喂价记录（仅 emit，不参与金额） |
| `setPayment` | **缺失** | 新增 `setPayment(address) onlyOwner` | scan 遗留：Oracle 无法设 payment（B2 拓扑依赖此 setter） |
| 内联 `IPayment`/`IDeposit` L105-122 | 文件尾内联（含 `applyTrafficCardToBill`/`issueMonthlyTrafficCards`） | 收敛：import `./interfaces/IPayment.sol` + `IDeposit.sol`（统一 SSOT） | **依赖 B5**：interfaces/ 两接口先补对应声明，否则收敛后编译不过 |

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
- **`initialize(_userRegistry, _usdt)` 注入 usdt（消歧锁定，弃 setUSDT 后置）**——避免 usdt 未设时 `safeTransferFrom(address(0))` panic；fresh deploy 无升级包袱。
- **【B7】仅支持标准 ERC20，禁 fee-on-transfer**：若 token 转账实扣 < amount，Deposit 记账虚高 → withdraw 资损。本轮 mock 为标准 ERC20；正式 USDT 无 fee-on-transfer。设计层禁用 fee-on-transfer token（README/handoff 注明），不做「实收差额」补偿逻辑。
- `deposit(uint256 amount)`：`require(userRegistry.isRegistered(msg.sender))` + `require(amount >= 10*10**usdt.decimals(), "Below min deposit")` → `usdt.safeTransferFrom(msg.sender, address(this), amount)` → 累加 + 锁仓续期（沿用 L46-50）→ emit。
- `withdraw()`：CEI——先 `principal=_deposits[...]; _deposits[...]=0; _lockExpiry[...]=0;` 再 `usdt.safeTransfer(msg.sender, principal)`；加 `nonReentrant`。
- **【B3 / v2-C】发卡逻辑收敛到 internal `_mintFor(user)`**：三校验（`now≥lockExpiry` && `deposits>0` && `userCardCount==0`）+ `trafficCardNFT.mint(user, trafficCardQuota, generateTokenURI(user))`。`mintTrafficCard(onlyOwner)` 调 `_mintFor` 并 require 成功；`issueMonthlyTrafficCards(onlyOracle)` 循环调 `_mintFor`，校验失败 `continue` 跳过（不 revert 整批）。
- **【v2-C】`dataAmount = trafficCardQuota`（固定 100MB，Deposit L17 state / initialize L27），删除 L76 `_deposits[user]/100000`**——与存款额、与 USDT 精度彻底解耦，杜绝 6 位精度下发废卡（旧逻辑 10 USDT=1e7 /100000=100 单位含义不明）。全文口径统一为此。

### 5.2 Payment — 分账 + 手续费走 FeeManager
- 引 `IServiceManager` 查 `getOperator(operatorId).paymentAddress`。
- **【B2】`createBill` 改 `onlyOracle`**（用 Payment 已有 modifier L20-23）；**【消歧】createBill fail-fast：require `serviceManager.getOperator(operatorId).paymentAddress != 0`**，避免生成永远付不了的账单。amount 为链下算好金额（v2-A）。
- `payBill`：去 payable + 去找零；**CEI 先置 `bill.isPaid=true`**；`require(operator.paymentAddress != address(0), "Operator payout unset")`；`safeTransferFrom(user→operator, amount)` + `if(bill.platformFee>0) safeTransferFrom(user→platform, fee)`（**0-fee 跳过，消歧**）；加 `nonReentrant`（沿用现有 guard 或降级，见 §六.4）。
- **【B7】支付失败降级**：分账地址要求非黑名单/非暂停（部署时选用非黑名单地址；handoff 注明）；任一 `safeTransferFrom` revert → 整 tx 原子回滚，账单仍 unpaid，用户可换地址/重试。
- 手续费仍在 `createBill` 用 `feeManager.calculateFee(amount)` 存入 `bill.platformFee`（满足 R5，不改）。
- **【v2-B】`applyTrafficCardToBill(uint256 billId) external onlyOracle` = 受限桩**：**不转移任何资金**；仅 `require` 账单存在（如 `require(_bills[billId].user != address(0), "Bill not found")`）+ 可选 `emit TrafficCardApplied(billId)`。Oracle 在 monthlySettlement 中调用以走通编译/调用链。**真实「流量卡抵扣账单」语义（额度来源/抵扣范围/精度）明确冻结到后续 Round**，本轮不实现、补测面只覆盖权限与存在性校验（见 §八 ATC-01/02）。

### 5.3 TrafficCardNFT — mint/销毁/有效期
- 修编译（§4.5）后逻辑基本不动：`mint`（Deposit 持 ownership 调用）、`burn`（emit ServiceActivated now+30d）、`DEDUCTION_VALIDITY=30 days`。
- 自动发卡走 Deposit→NFT.mint 链路（ownership 已在 deploy 转给 Deposit，L88）。
- **【消歧】自动发卡禁用 `NFT.mintBatch`**：`mintBatch` L77 内 `this.mint(...)` 是**外部调用**，会以 NFT 合约自身为 `msg.sender` 撞 `mint` 的 `onlyOwner`（owner 已是 Deposit）→ revert。批量发卡一律由 Deposit 侧 `issueMonthlyTrafficCards` 循环单张 `trafficCardNFT.mint(...)` 完成。mintBatch 本轮不用（保留或留注释，不接入发卡链路）。

### 5.4 Oracle — wiring + 自动发卡触发
- 修 state 类型为接口（§4.4）、声明事件、加 setPayment、收敛内联接口（依赖 B5 接口先补声明）。
- **【B1/v2-A】`monthlySettlement` 改签名 `(users, operatorIds, amounts[])`**：删除 L68 `totalAmount = dataUsage + callUsage`；循环内直接 `payment.createBill(user, operatorId, amounts[i])`（amounts 为后端链下算好的 USDT 金额）。`dataUsages/callUsages` 若保留仅用于 `emit UsageDataSubmitted`（喂价记录，不参与金额）。onlyOwner 保留。
- **【B2】调用权限闭合**：`monthlySettlement → payment.createBill` 现可成功，因 createBill 改 `onlyOracle` 且 deploy 执行 `payment.setOracle(oracleAddr)`（见 §七.0 拓扑），Oracle 即为 Payment 授权的 oracle。`deposit.issueMonthlyTrafficCards` 同理依赖 `deposit.setOracle(oracleAddr)`（已在 deploy.ts L82）。
- **【v2-B】applyTrafficCardToBill 桩调用**：`monthlySettlement` 末段对每个 user 的首张未付账单调 `payment.applyTrafficCardToBill(billId)`（桩，不转资金，仅走通调用链）。
- **【消歧·批量 gas】无界循环上限策略**：`monthlySettlement` / `issueMonthlyTrafficCards` 由后端**分批调用（每批 users ≤ N，建议 N≤100）**；本轮**不在合约内强制 maxBatch**（保留调用方分批），但补**压测用例**（§八 GAS-01）量出单批安全上限，写入 handoff 供后端遵守。
- **自动发卡触发链**：后端（子项 2/3）定时 → `Oracle.monthlySettlement(...)` 或直接 `Deposit.issueMonthlyTrafficCards(...)`（onlyOracle）。`Deposit.issueMonthlyTrafficCards` 权限 = `require(msg.sender==oracle)`（沿用 L105），实现真实批量 mint（§5.1 `_mintFor` + B3 跳过不 revert）。

### 5.5 mock USDT（6 decimals）
- `contracts/mocks/MockUSDT.sol`：`ERC20`（OZ 标准）+ `decimals() returns(6)` override + `mint(to, amount)`（public，测试网任意发币）+ symbol "USDT"。
- 仅本地/测试网部署；deploy.ts 部署后写 `deployments/<net>.json` 的 `usdt` 字段 + `usdtDecimals:6`。

### 5.6 编译错误清单与逐个修复方案

| # | 文件:行 | 错误 | 根因 | 修复 |
|---|---------|------|------|------|
| ① | `Oracle.sol:71` | `DeclarationError: UsageDataSubmitted` 未声明 | 事件从未定义（实测唯一中止点） | 在 Oracle 声明 `event UsageDataSubmitted(address indexed,uint256,uint256,uint256)` |
| ② | `Oracle.sol:69/78/85`（修①后暴露） | `payment.createBill` / `deposit.issueMonthlyTrafficCards` / `payment.applyTrafficCardToBill` —— `address` 类型无成员 | `deposit`/`payment` 声明为 `address`（L10-11）却当合约调 | 改 state 为 `IDeposit public deposit; IPayment public payment;`，import 接口 |
| ③ | `Payment.sol`（Oracle 调 `applyTrafficCardToBill`） | Payment 无此函数 | Oracle 内联接口声明了但 Payment 未实现 | **【v2-B】Payment 实现 `applyTrafficCardToBill(uint256) onlyOracle` 受限桩**（不转资金，§5.2）；同时 `createBill` 改 `onlyOracle`（B2） |
| ④ | `TrafficCardNFT.sol:18` | `DeclarationError: DeductionCredit` 未定义 | 类型从未定义（dead code，scan §④ 确认未用） | 删除该 mapping + `onlyDeposit` modifier（dead code，倾向删） |
| ⑤ | `Oracle.sol` 内联接口 vs interfaces/ | 双份 IPayment/IDeposit 漂移风险 | 历史内联 | 收敛为 import interfaces/。**【B5 连锁编译清单】先在 interfaces/IDeposit.sol 补 `issueMonthlyTrafficCards(address[])`、interfaces/IPayment.sol 补 `applyTrafficCardToBill(uint256)`**，否则 Oracle import 后调用编译不过；**接口签名变更 → 实现同步**：IDeposit.deposit 去 payable → Deposit 实现同步；IPayment.payBill 去 payable → Payment 实现同步 |

> **【B5 连锁编译清单】接口签名变更 → 实现/调用方必须同步**（任一漏改即编译失败）：
> 1. `IDeposit.deposit() payable` → `deposit(uint256)`：同步 Deposit.sol 实现 + 调用方。
> 2. `IDeposit` 补 `issueMonthlyTrafficCards(address[])`：Oracle 收敛后据此调用。
> 3. `IPayment.payBill() payable` → `payBill(uint256)`：同步 Payment.sol。
> 4. `IPayment` 补 `applyTrafficCardToBill(uint256)`：Oracle 收敛 + Payment 桩实现据此。
> 5. `IPayment.createBill` 权限收紧（onlyOracle）属实现内修饰符，不改接口签名但改 deploy 授权拓扑（§七.0）。
>
> **implement 顺序 / P1 compile gate（消歧重定义）**：①（事件/类型/monthlySettlement 改签名）+ ④（删 dead code）+ ⑤（先补 interfaces/ 声明）+ ②（state 类型）+ ③（createBill onlyOracle + applyTrafficCardToBill 桩）—— **必须 ①–⑤ 全部做完，compile 才会绿**（②类型改完后缺接口方法、③桩未实现前都编译不过，不存在「②④后即绿」的中间态）。每完成一类 `npx hardhat compile` 观察剩余报错收敛。

---

## 六、非功能性设计

### 6.1 资损 checklist（🔴 = arch-review 安全审计重点）

| 项 | 风险 | 设计对策 |
|----|------|----------|
| 🔴 ERC20 重入 | safeTransferFrom 回调（恶意/非标 token） | CEI：状态先改（isPaid/deposits 清零）再转账；payBill/withdraw 加 `nonReentrant` |
| 🔴 授权额度 | approve 总额 vs 实际扣款不一致、无限授权风险 | 合约只 `safeTransferFrom(精确额)`；前端 approve 精确额（不无限授权） |
| 🔴 精度 | 硬编码 18/6 导致金额错算（scan §七雷区） | 全程 `usdt.decimals()`；`MIN_DEPOSIT=10*10**decimals`；**【v2-C】`dataAmount=trafficCardQuota` 固定，删除 `_deposits/100000`，与精度彻底解耦** |
| 🔴 计价量纲 | **【B1】Oracle L68 `dataUsage(字节)+callUsage(分钟)` 当 USDT 金额传 createBill → ERC20 扣天文数字** | **【v2-A】合约不计价**：amount 由后端链下按资费算好的 USDT 金额传入 createBill；删除 L68 求和；Oracle 仅喂价/触发 |
| 🔴 分账零地址 | operator.paymentAddress=address(0) → USDT 转入黑洞 | payBill `require(paymentAddress != address(0))`；setOperatorPaymentAddress 零地址 require；**createBill fail-fast 校验（消歧）**；部署脚本补全 |
| 🔴 非标代币 fee-on-transfer | **【B7】token 实扣 < amount → Deposit 记账虚高 → withdraw 资损** | **仅支持标准 ERC20，禁 fee-on-transfer**（handoff/README 注明）；mock 为标准 ERC20，正式 USDT 无此特性；不做实收差额补偿 |
| 🔴 真实 USDT 黑名单/暂停 | **【B7】operator/platform 地址被 USDT 黑名单或合约暂停 → payBill 永久失败** | 分账地址要求**非黑名单**（部署选用，handoff 注明）；**支付失败降级**：任一 transferFrom revert → 整 tx 回滚、账单仍 unpaid、用户可换地址重试 |
| 整数溢出 | `amount+fee`、`*rate/10000` | Solidity 0.8 内置 checked；calculateFee 小额截断为 0 可接受（费率本就向下取整） |
| 找零逻辑移除 | ERC20 无 msg.value，旧多退少补删除 | payBill 精确 transferFrom，无残留资金 |
| 部分扣款失败 | 两段 transferFrom 第一段成功第二段失败 | 同一 tx 原子性，任一 revert 全回滚 |
| applyTrafficCardToBill 资损面 | v1 设想真实抵扣 → 额度/精度未定有资损风险 | **【v2-B】本轮受限桩不转资金 → 零资损面**；真实抵扣语义连同其资损评估冻结到后续 Round |

### 6.2 安全（权限矩阵）

| 函数 | 权限 | 校验 |
|------|------|------|
| `Deposit.deposit` | 任意已注册用户 | `isRegistered` + `amount≥MIN` |
| `Deposit.withdraw` | 持仓用户 | `now≥lockExpiry` + `deposits>0` |
| `Deposit.issueMonthlyTrafficCards` | **onlyOracle** | `msg.sender==oracle` + 逐用户三校验（失败 continue，不 revert 整批 / B3）；走 `_mintFor` |
| `Deposit.mintTrafficCard` | onlyOwner | 三校验（薄壳调 `_mintFor`） |
| `Payment.createBill` | **onlyOracle（B2）** | amount>0 + fail-fast paymentAddress≠0（消歧）；Oracle 经 `payment.setOracle` 获权 |
| `Payment.payBill` | 账单本人 | `!isPaid` + `bill.user==sender` + paymentAddress 非零；CEI 先置 isPaid |
| `Payment.applyTrafficCardToBill` | **onlyOracle（桩 / v2-B）** | 仅账单存在；**不转资金** |
| `Oracle.monthlySettlement` | onlyOwner | 长度一致；amounts 链下算好（v2-A，无求和） |
| `ServiceManager.setOperatorPaymentAddress` | onlyOwner | 非零地址 |
| `setUSDT/setOracle/setPayment/...` | onlyOwner | |

- SafeERC20 全覆盖；ReentrancyGuard 覆盖 payBill/withdraw。

### 6.3 升级兼容（UUPS storage layout）
> **消歧拍板：本轮 16602/localhost 一律 fresh deploy，不升级任何旧 proxy。** 因此 v1 关于「升级 16602 旧 proxy 的 storage layout 约束 / mapping 删除占位风险」的讨论**全部删除（无效负担）**。下列只保留 fresh deploy 口径。
- **fresh deploy 421614 + 31337**：无旧 storage 约束。新增 `usdt`/`serviceManager` state 与删除 `_deductionCredits` mapping 均无 layout 风险，按显式声明顺序排布即可（推荐新增 state 放各合约 storage 末尾，仅为未来 Round 升级预留习惯，非本轮约束）。
- **storage layout 冻结（交付物）**：本轮部署后，将各 proxy 的 storage layout 与 ABI **冻结记入 `deployments/<net>.json`**（见 §七.3），作为后续 Round 升级的基线。本轮不做升级，故无需 `unsafeAllow`。

### 6.4 Arbitrum 兼容确认项（🔴 **消歧：实测前移到 P2 之前 / 与 P1 并行**，不留到 P5）
> **消歧拍板**：原 v1 把 transient storage 兼容实测排在 P5（部署阶段），风险是 P5 才发现不支持 → 已写好的 Payment guard 逻辑返工。**本版前移：P1 compile 绿后、动 ②ERC20 之前**，先在 421614 实测 Payment 部署 + 一笔 payBill，确认 TSTORE 可用再继续；不可用则**先**切 guard 方案再写后续。
| 项 | 现状 | 风险 | 确认动作（前移） |
|----|------|------|----------|
| 🔴 `evmVersion: cancun` + `viaIR` | hardhat.config L15-16 | Arbitrum Sepolia 对 cancun opcode 支持？ | **P2 之前**在 421614 跑一笔交易验证；Arbitrum Nitro 已支持多数 cancun，但 **transient storage(TSTORE/TLOAD)** 需重点确认 |
| 🔴 `ReentrancyGuardTransient` | Payment L11（依赖 TSTORE） | 若 Arbitrum 不支持 TSTORE → Payment 部署/调用失败 | 二选一：(a) 确认支持后保留；(b) 降级 `ReentrancyGuardUpgradeable`（普通 storage，需改 initialize + 删自定义 `_reentrancyGuardInit` assembly L33-38）。**实测前移到 P2 之前** |
| 🔴 自定义 `_reentrancyGuardInit` assembly 正确性 | Payment L33-38：`sstore(_reentrancyGuardStorageSlot(), 1)` | **消歧·正确性可疑**：`ReentrancyGuardTransient` 用的是 **transient slot（tstore/tload）**，而此处用 `sstore`（持久 storage）写 NOT_ENTERED，**slot 类型与 guard 读取方式可能不一致**——可能既无效又浪费 storage | implement 阶段**厘清** OZ `ReentrancyGuardTransient` 实际读哪个 slot/用哪个 opcode：若 guard 用 tload 读，则此 sstore 初始化无意义应删；若保留 transient guard 通常**无需手动 init**（transient 每 tx 自动归零）。**不只是「降级时删」，而是先确认其当前是否正确**。若走 (b) 降级则随 `ReentrancyGuardUpgradeable` 的 `__ReentrancyGuard_init()` 一并删除此段 |
| chainId/RPC | 无 421614 | — | hardhat.config 加网络 |

### 6.5 异常与监控
- 事件不变（DepositMade/BillCreated/BillPaid/CardMinted/ServiceActivated）+ 新增 `UsageDataSubmitted`（Oracle 喂价记录）+ 可选 `TrafficCardApplied(billId)`（v2-B 桩）；后端 event_sync 据此监听（子项 2/3）。
- revert 文案统一英文短串（"Below min deposit"/"Operator payout unset"/"Lock not expired"/"Only oracle"/"Bill not found"）。

---

## 七、部署设计

### 7.0 部署后 ownership / 授权拓扑图（🔴 B2/B3 核心交付）

> 这是 B2/B3 的核心闭合：自动结算两条权限链（createBill / issueMonthlyTrafficCards / applyTrafficCardToBill）在 deploy 下必须能跑通。本节明确**谁是谁的 owner、哪些 setter 在 deploy 调用、每条权限链如何闭合**。

#### 7.0.1 owner 关系（部署后终态）
```
deployer(EOA)
  ├── owner of: FeeManager, UserRegistry, ServiceManager, Payment, Oracle  ← 保留为 deployer
  ├── owner of: Deposit                                                     ← 保留为 deployer
  └── NFT ownership: TrafficCardNFT.owner = Deposit(合约)                    ← deploy.ts L88 transferOwnership(deposit)

注：Payment / Oracle 的 owner 仍是 deployer（不再 transferOwnership 给谁）。
   Oracle 调用 Payment/Deposit 的权限不通过 ownership，而通过各自的 oracle 授权字段（见下）。
```

#### 7.0.2 调用授权字段（非 ownership，通过 setter 注入）
| 授权关系 | 注入 setter（deploy 调用） | 用途 | 现状 |
|----------|---------------------------|------|------|
| `Payment.oracle = Oracle` | **`payment.setOracle(oracleAddr)`** | createBill / applyTrafficCardToBill 的 `onlyOracle` 放行 Oracle | **deploy.ts 现缺 → B2 必补** |
| `Deposit.oracle = Oracle` | `deposit.setOracle(oracleAddr)` | issueMonthlyTrafficCards 的 `msg.sender==oracle` 放行 | deploy.ts L82 已有 ✅ |
| `Oracle.payment = Payment` | **`oracle.setPayment(paymentAddr)`** | Oracle 调 `payment.createBill/applyTrafficCardToBill` | **Oracle 现无 setPayment + deploy 缺 → B2 必补** |
| `Oracle.deposit = Deposit` | `oracle.setDeposit(depositAddr)` | Oracle 调 `deposit.issueMonthlyTrafficCards` | deploy.ts L84 已有 ✅ |
| `Deposit.trafficCardNFT = NFT` | `deposit.setTrafficCardNFT(nftAddr)` | Deposit `_mintFor` 调 `nft.mint` | deploy.ts L81 已有 ✅ |
| `TrafficCardNFT.owner = Deposit` | `trafficCardNFT.transferOwnership(depositAddr)` | NFT.mint 的 `onlyOwner` 放行 Deposit | deploy.ts L88 已有 ✅ |
| `Payment.serviceManager = SM` | **`payment.setServiceManager(smAddr)`**（或 initialize 注入） | payBill / createBill 查 operator.paymentAddress | **现缺 → 必补** |
| `Payment.usdt`, `Deposit.usdt` | **initialize 注入 MockUSDT**（消歧锁定） | ERC20 资金通道 | **现缺 → 必补** |
| `ServiceManager` 11 运营商 paymentAddress | **循环 `serviceManager.setOperatorPaymentAddress(id, addr)`** | 分账非零地址 | **现全 address(0) → 必补** |

#### 7.0.3 三条权限链闭合验证（B2/B3）
```
链 A：自动结算账单（B2）
  后端 → Oracle.monthlySettlement(onlyOwner=deployer 调)
       → Payment.createBill(onlyOracle)
         放行条件：Payment.oracle == Oracle   ← 由 payment.setOracle(oracle) 注入（B2 必补）
       ✅ 闭合后不再 revert

链 B：自动发卡（B3）
  Oracle.monthlySettlement → Deposit.issueMonthlyTrafficCards(onlyOracle)
         放行条件：Deposit.oracle == Oracle   ← deploy.ts L82 已注入
       → 内部走 _mintFor（不经 onlyOwner mintTrafficCard）
       → TrafficCardNFT.mint(onlyOwner)
         放行条件：NFT.owner == Deposit        ← deploy.ts L88 已注入
       逐用户 if(!三校验) continue（不 revert 整批）
       ✅ 闭合，幂等

链 C：applyTrafficCardToBill 桩（v2-B）
  Oracle.monthlySettlement → Payment.applyTrafficCardToBill(onlyOracle)
         放行条件：Payment.oracle == Oracle   ← 同链 A
       → 桩：require 账单存在 + 可选 emit，不转资金
       ✅ 走通调用链，零资损面
```

> **B2 关键修正**：v2 选定 **createBill 改 `onlyOracle` + `payment.setOracle(oracle)` 授权**，而非 `payment.transferOwnership(oracle)`。理由：Payment 的 owner 还要保留给 setPlatformWallet/setFeeManager/setServiceManager 等管理操作；用 oracle 字段精确授权调用权，影响面最小、语义最清。

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

### 7.2 deploy.ts 增量（授权拓扑见 §7.0）
顺序（在现有 7 合约 deployProxy 基础上）。**usdt/serviceManager 走 initialize 注入（消歧锁定），oracle 因部署顺序仍走 setter**：
```
0. 部署 MockUSDT(decimals=6)  ← 测试网/本地；本轮无正式 USDT，全用 mock
1-7. proxy 部署（initialize 注入 usdt/serviceManager）：
   - Deposit.initialize(userRegistry, usdt)                       // ← usdt initialize 注入
   - Payment.initialize(feeManager, platformWallet, usdt, serviceManager)  // ← usdt/SM initialize 注入
   （其余合约 initialize 不变）

§7.0 授权 wiring 补全（★ = 现状缺，B2 必补）：
  ★ + payment.setOracle(oracleAddr)        // B2：createBill/applyTrafficCardToBill 的 onlyOracle 放行
  ★ + oracle.setPayment(paymentAddr)       // B2：Oracle 无 setPayment，本轮新增
  ★ + 循环 serviceManager.setOperatorPaymentAddress(id, payoutAddr)  // 11 运营商，非零（测试网用 deployer 派生地址）
    （现有）trafficCardNFT.setDepositContract / deposit.setTrafficCardNFT / deposit.setOracle / oracle.setDeposit / payment.setPlatformWallet / transferOwnership(deposit, NFT)
```
> 校验：deploy 完成后断言 `payment.oracle()==oracle`、`oracle.payment()==payment`、`deposit.oracle()==oracle`、`NFT.owner()==deposit`、每个 active operator.paymentAddress≠0（§7.0.3 三条链闭合的前置）。

### 7.3 deployments 产出 + 给后端(2/3)的 handoff 交付物（消歧·CEO⚠️）
`deployments/<net>.json` 增字段：
```json
{ "chainId": 421614, "proxies": {...7...},
  "implementations": {...7...},
  "usdt": "0x...", "usdtDecimals": 6,
  "storageLayout": { "Deposit": "...", "Payment": "...", ... },   // ← storage layout 冻结记入（消歧）
  "abiHash": { "Deposit": "0x...", "Payment": "0x..." }            // ← ABI 指纹，供后端比对
}
```
两处部署：`deploy:arbitrum-sepolia`（421614）+ `deploy:local`（31337）。供后端 configs/deployments.json（子项 2/3）+ 前端 contracts.ts（子项 3/3）对齐。

#### 7.3.1 合约子项 → 后端子项(2/3) 正式 handoff 清单（验收交付物）
合约子项验收**必须产出**以下 handoff（消歧·CEO⚠️），交给后端对齐：
1. **冻结 ABI**：7 合约 + MockUSDT 的最终 ABI（本轮 selector 已变，后端必须重生成）。
2. **selector 清单 + 变更标注**：重点标 `deposit()`→`deposit(uint256)`、`payBill()`→`payBill(uint256)`（去 payable 改 selector）、`monthlySettlement` 新签名（含 amounts[]，v2-A）、新增 `issueMonthlyTrafficCards`/`applyTrafficCardToBill`/`setOperatorPaymentAddress`/`oracle.setPayment`。
3. **金额精度语义说明**：所有金额字段单位 = USDT，精度 = `usdt.decimals()`（mock 6）；**createBill 的 amount 是后端链下算好的 USDT 金额（v2-A），后端负责计价**；MIN_DEPOSIT=10×10^decimals。
4. **storage layout 冻结**：记入 deployments.json（见上），作为后续 Round 升级基线。
5. **批量调用约定**：monthlySettlement/issueMonthlyTrafficCards 后端分批 ≤N（压测 GAS-01 量出的安全上限）。
6. **applyTrafficCardToBill 语义说明**：本轮为受限桩（不抵扣资金），后端不要据此做抵扣对账（v2-B 冻结）。

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
| ISS-02 | oracle 调，锁仓到期用户 → mint NFT；未到期用户 → 跳过不 revert | **自动发卡时序 / B3 不整批回滚** |
| ISS-03 | 已有卡用户重复 issue → 不重复 mint | 幂等 |
| ISS-04 | 一批含「合格+不合格」混合用户 → 合格者 mint、不合格者 continue 跳过、整批不 revert | **B3 核心** |
| ISS-05 | issueMonthlyTrafficCards 发卡 dataAmount == trafficCardQuota（固定，与存款额无关） | **v2-C** |
| DEC-01 | mock USDT.decimals()=6；MIN_DEPOSIT 随 decimals 计算 | **精度（不硬编码）** |
| **B6 集成测试** | | |
| MS-01 | deploy 全 wiring 后 owner 调 `monthlySettlement(users,operatorIds,amounts[])` → 对应 user 生成账单 amount==amounts[i]（**v2-A 不求和**）、platformFee 正确 | **B1/B6：计价 + createBill onlyOracle 链通** |
| MS-02 | createBill 的调用方非授权 oracle → revert "Only oracle" | **B2：createBill onlyOracle** |
| MS-03 | monthlySettlement 端到端：createBill + issueMonthlyTrafficCards + applyTrafficCardToBill 桩 全程不 revert（验证 §7.0.3 三条链闭合） | **B6 集成主路径** |
| ATC-01 | 非 oracle 调 applyTrafficCardToBill → revert "Only oracle" | **B6/v2-B：桩权限** |
| ATC-02 | oracle 调 applyTrafficCardToBill(存在账单) → 不 revert、不发生任何 USDT 转移（余额不变）；调不存在账单 → revert "Bill not found" | **v2-B：桩不转资金 + 存在性校验** |
| **B7 非标 USDT** | | |
| USDT-01 | transfer/transferFrom 不返回值的非标 USDT → SafeERC20 仍正确入账（不 revert） | **B7：SafeERC20 分支** |
| USDT-02 | transferFrom 返回 false 的恶意 token → deposit/payBill revert（SafeERC20 捕获） | **B7：SafeERC20 分支** |
| **压测** | | |
| GAS-01 | issueMonthlyTrafficCards / monthlySettlement 批量 N 用户，量 gas，确定单批安全上限 N | **消歧：批量 gas 上限 → 写入 handoff** |
| **回归（重写）** | | |
| REG-01 | 锁仓续期不变量：未到期复存 → lockExpiry 原点+30d；到期复存 → now+30d | **消歧：续期不变量** |
| REG-02 | 未注册用户 deposit → revert "Not registered"（旧测试用 `{value}` payable 调用，改 `deposit(amount)` 后失效，**重写为 approve+deposit(amount)**） | **消歧：回归重写** |
| 回归 | 现有 FE/UR/SM/TC 26 it（金额改 USDT 精度后）不回归 | 不破坏 |

> **已删除**：v1 的 `UPG-01（upgradeProxy storage 不丢）` —— 本轮一律 fresh deploy 不升级（消歧），无升级路径可测。
> 现有测试金额从 `parseEther` 改 USDT 精度辅助：`const usdt=(n)=>BigInt(n)*10n**6n`。FeeManager 的 FE-02/04 用 ETH 计费率与币种无关，可保留或改 USDT 单位。测试需先 mint mock USDT 给 user 并 approve。**旧测试直接 `Factory.deploy()`+手动 initialize 且用 `{value}` 调 deposit/payBill 的用例，改 ERC20 后全部失效，需按 approve+deposit(amount)/payBill(id) 重写。**

---

## 九、落地计划（给 plan 阶段输入）

| 阶段 | 任务 | 顺序 | 风险点 |
|------|------|------|--------|
| P1 修编译 | ①Oracle 事件+类型+setPayment+**monthlySettlement 改 amounts[] 签名(v2-A)**；②state 类型；③Payment **createBill 改 onlyOracle(B2)** + **applyTrafficCardToBill 受限桩(v2-B)**；④TrafficCardNFT 删 DeductionCredit/onlyDeposit；⑤**interfaces/ 先补 issueMonthlyTrafficCards/applyTrafficCardToBill 声明(B5)** 再收敛 Oracle 内联接口 | **最先，串行** | compile 绿是 gate；**P1 gate 重定义：①–⑤ 全做完才绿（消歧），无中间绿态** |
| **P1.5 兼容实测（前移）** | 在 421614 实测 Payment 部署 + 一笔 payBill，确认 TSTORE/ReentrancyGuardTransient 可用；厘清 `_reentrancyGuardInit` assembly 正确性 | **P1 后、P2 前（消歧前移）** | 🔴 不支持则先切 `ReentrancyGuardUpgradeable` 再继续，避免 P5 返工 |
| P2 ERC20 迁移 | Deposit/Payment 引 SafeERC20 + **usdt 走 initialize 注入** + deposit(amount)/withdraw 改写 + MIN_DEPOSIT；**禁 fee-on-transfer(B7)** | P1.5 后 | 精度雷区；用 decimals 不硬编码 |
| P3 分账 | Payment payBill CEI 直分 + ServiceManager setOperatorPaymentAddress + 零地址校验 + **createBill fail-fast** + **0-fee 跳过** | P2 后 | operator.paymentAddress 全 0 必须先补；两段 transferFrom 原子性；黑名单地址降级(B7) |
| P4 自动发卡 | Deposit.issueMonthlyTrafficCards 走 **独立 `_mintFor`(B3)** + onlyOracle + **continue 跳过不整批回滚(B3)** + **dataAmount=trafficCardQuota(v2-C)**；**禁用 NFT.mintBatch** | P3 后 | 时序校验、批量 gas（GAS-01 压测定上限） |
| P5 部署 | hardhat.config 421614 + MockUSDT + **deploy.ts 补 B2 授权拓扑(§7.0：payment.setOracle/oracle.setPayment/setOperatorPaymentAddress)** + 校验断言 + **handoff 交付物(§7.3.1)** | P4 后 | 授权拓扑漏一项则自动结算链 revert（§7.0.3 校验兜底） |
| P6 补测 | §八 全清单（含 B6 MS/ATC、B7 USDT、GAS-01、回归重写 REG） | 贯穿（TDD 可前置到各阶段） | 现有 26 it 不回归、需重写 payable 调用 |

**implement 串行铁律**：P1→P6 严格串行（一个完成审查再下一个），尤其 P1 compile gate 不过不进 P1.5。**P1.5 Arbitrum 兼容实测**（消歧前移，原 P5）若失败，先回退到「Payment ReentrancyGuard 降级」分支（§6.4）再进 P2。

---

## 十、🔴 重跑 arch-review 重点复审清单（v2）

1. **计价归属（B1/v2-A）**：Oracle 无任何 usage 求和；createBill amount = 链下 USDT 金额；量纲正确无天文扣款。
2. **两条权限链 deploy 可跑通（B2）**：§7.0.3 链 A/B 闭合；deploy 确含 `payment.setOracle` + `oracle.setPayment`；createBill 为 onlyOracle。
3. **批量发卡幂等 + 不整批回滚（B3）**：走 `_mintFor` 不撞 onlyOwner；`continue` 跳过；混合批 ISS-04 通过。
4. **发卡额度解耦（B4/v2-C）**：dataAmount=trafficCardQuota 固定；全仓无 `_deposits/100000`；口径全文统一。
5. **接口 SSOT 收敛 + 连锁编译（B5）**：interfaces/ 已补两声明；签名变更→实现同步清单（§5.6）；selector 变更已进 handoff。
6. **核心新函数测试覆盖（B6）**：MS-01~03 集成 + ATC-01/02 桩权限存在性。
7. **非标代币边界（B7）**：禁 fee-on-transfer；黑名单降级；USDT-01/02 SafeERC20 分支。
8. **applyTrafficCardToBill 受限桩（v2-B）**：onlyOracle、不转资金、仅存在性校验；真实抵扣语义已冻结。
9. **ERC20 重入 + CEI**：payBill/withdraw 状态先改后转、nonReentrant 到位。
10. **精度不硬编码**：全仓 grep 无字面 `* 10**18` / `* 10**6` / `1e18`；MIN_DEPOSIT 随 decimals。
11. **Arbitrum transient 兼容（消歧前移）**：P1.5 实测 ReentrancyGuardTransient + 厘清 `_reentrancyGuardInit` assembly 正确性。
12. **fresh deploy 口径**：无 16602 升级讨论残留；storage layout 冻结记入 deployments.json。
13. **ServiceManager requiredDeposit 废弃**：本轮不读不校验，无误用。

---

## 十一、arch-review 阻塞闭合对照表（B1–B7 + 用户 3 决策）

| 阻塞 | 原问题 | v2 解法 | 落点章节 |
|------|--------|---------|----------|
| **B1**🔴 | Oracle L68 `dataUsage+callUsage` 当 USDT 金额 → 天文扣款/量纲错 | **v2-A**：删 L68 求和；`monthlySettlement(...,amounts[])` 传入链下算好的 USDT 金额；Oracle 只喂价/触发不计价 | 头部 v2-A、2.2、3.2、4.4、5.4、6.1、§八 MS-01 |
| **B2** | createBill 是 onlyOwner，Oracle 非 owner 且 deploy 无授权 → monthlySettlement revert | createBill 改 **onlyOracle** + deploy 补 **`payment.setOracle(oracle)` / `oracle.setPayment`**；新增 **§7.0 ownership/授权拓扑图** + §7.2 授权步骤 | 2.2、4.2、4.4、5.2、6.2、**§7.0**、7.2、§八 MS-02 |
| **B3** | issueMonthlyTrafficCards 复用 onlyOwner mintTrafficCard 会 revert；require 整批回滚 | 走独立 internal **`_mintFor`**（不经 onlyOwner）；批量循环 **`if(!ok) continue`** 跳过不 revert 整批 | 3.1、4.1、5.1、5.4、6.2、§7.0.3 链B、§八 ISS-04 |
| **B4** | `_deposits/100000` 在 6 位精度下发废卡；口径矛盾 | **v2-C**：`dataAmount=固定 trafficCardQuota`，删除 `_deposits/100000`，全文统一 | 头部 v2-C、3.1、4.1、5.1、6.1、§八 ISS-05 |
| **B5** | IDeposit 无 issueMonthlyTrafficCards、IPayment 无 applyTrafficCardToBill → 接口收敛后编译不过 | interfaces/ 先补两声明；**§5.6 连锁编译清单**（接口签名变更→实现/调用方同步） | 〇、4.1、4.2、4.4、**5.6** |
| **B6** | monthlySettlement + applyTrafficCardToBill 零测试 | 补 **MS-01~03 集成测试** + **ATC-01/02 桩权限/存在性测试** | 1.3、§八（B6 集成测试块） |
| **B7** | 未处理 fee-on-transfer（withdraw 资损）+ USDT 黑名单/暂停 payBill 失败 | 6.1 补：**仅标准 ERC20、禁 fee-on-transfer**；分账地址非黑名单 + **支付失败降级**；测 USDT-01/02 | 2.1、5.1、5.2、**6.1**、§八（B7 块） |

**§三消歧项闭合**（一并处理）：

| 消歧项 | v2 处理 | 落点 |
|--------|---------|------|
| usdt 注入路径 | 锁定 **initialize 注入**（弃 setUSDT 后置） | 2.2、4.1、4.2、5.1、7.2 |
| P1 compile gate 重定义 | **①–⑤ 全做完才绿**（无中间绿态） | 1.3、2.1、5.6、§九 |
| 批量 mint 分批上限/压测 | 调用方分批 ≤N + **GAS-01 压测**定上限，写入 handoff | 5.4、§八 GAS-01、7.3.1 |
| ServiceManager requiredDeposit | **本轮明确废弃，不读不校验** | §一红线、4.3、§十.13 |
| 16602 一律 fresh deploy 不升级 | 删除所有 storage layout 升级讨论；fresh deploy 口径 | 2.2、§一红线、**6.3** |
| ReentrancyGuardTransient 兼容实测前移 + assembly 正确性 | 实测**前移到 P1.5（P2 前）**；厘清 `_reentrancyGuardInit` sstore 正确性 | **6.4**、§九 P1.5 |
| createBill fail-fast 校验 paymentAddress | createBill require operator.paymentAddress≠0 | 3.2、5.2、6.1 |
| 0-fee 跳过 | 保留 `if(platformFee>0)` | 3.2、4.2、5.2 |
| 禁用 NFT.mintBatch | L77 `this.mint()` 撞 onlyOwner，发卡链路禁用 | 〇、5.3 |
| handoff 交付物（ABI+selector+精度语义）+ storage layout 冻结 | **§7.3.1 handoff 清单** + deployments.json 记 storageLayout/abiHash | 6.3、**7.3 / 7.3.1** |
| 非标 USDT 测试 + 锁仓续期/未注册 deposit 回归重写 | USDT-01/02 + REG-01/02 重写 | §八（B7 块、回归块） |
