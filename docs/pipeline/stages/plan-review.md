# Stage: plan — 合约 implement 任务拆解（子项目 contracts 1/3）

> **状态**: completed | **日期**: 2026-06-08 | **Gate**: 1（设计）→ implement | **子项目**: contracts(1/3) | **角色**: 工程经理（任务拆解）
> 输入：design.md（v2，§九落地计划 P1-P6 + §十一闭合表 + §7.0 授权拓扑 + §八测试清单）+ arch-review.md（§七 5 条带入 implement 的 ⚠️ 验收单）+ requirement.md（§五 验收 A/D）。
> 范围：**只产任务计划，不写任何 .sol / 脚本 / 测试 / 配置**（实现留 implement 阶段）。
> 领域适配：本子项是 Solidity 合约，无页面/无色值；下文「任务/Task」= 一个合约实现阶段（P1-P6），「受影响合约/函数」取代 web 形状的「组件」。
>
> **implement 串行铁律**：T1→T6 严格串行，一个任务完成审查后再派下一个。P1 compile gate 不绿不进 T1.5；P1.5 兼容实测失败先切 guard 降级再进 T2。

---

## 1. 任务拆解（串行 T1→T6 / P1-P6）

> 每个任务对应 design §九 一个 P 阶段。gate 列为该任务的「完成判据」，必须达成才进下一任务。

| 任务ID | 阶段 | 范围 | 受影响合约/函数 | 依赖 | compile/test gate |
|--------|------|------|-----------------|------|-------------------|
| **T1** | P1 修编译 | ①Oracle 声明 `UsageDataSubmitted` 事件 + `deposit`/`payment` state 改接口类型 + 新增 `setPayment` + `monthlySettlement` 改 `(users,operatorIds,amounts[])` 删 L68 求和 + 收敛内联接口到 interfaces/；②interfaces/ 先补声明（IDeposit `issueMonthlyTrafficCards`、IPayment `applyTrafficCardToBill`）；③TrafficCardNFT 删 `DeductionCredit` mapping + `onlyDeposit` dead code；④Payment `createBill` 改 `onlyOracle` + 实现 `applyTrafficCardToBill` **受限桩**（onlyOracle/不转资金/仅存在性校验+可选 emit） | `Oracle.sol`（state/event/setPayment/monthlySettlement/import）、`interfaces/IDeposit.sol`、`interfaces/IPayment.sol`、`TrafficCardNFT.sol`、`Payment.sol`（createBill/applyTrafficCardToBill 桩） | — | **`npx hardhat compile` 零 error**（P1 gate 重定义：①–④ 全做完才绿，无中间绿态） |
| **T1.5** | P1.5 兼容实测（前移） | 421614 实测 Payment 部署 + 一笔 payBill，确认 `ReentrancyGuardTransient`(TSTORE/TLOAD) 兼容；删 Payment L33-38 无效 `_reentrancyGuardInit` assembly（sstore 持久 storage，transient guard 用 tload 读不同空间 → 永不被读=无效）+ initialize L28 对应调用；不兼容则降级 `ReentrancyGuardUpgradeable`（改 initialize + `__ReentrancyGuard_init`） | `Payment.sol`（L11 guard 选型、L33-38 assembly、initialize） | T1 | 421614 实测 payBill 成功（TSTORE 兼容）**或** 降级方案 compile+部署成功；删除无效 assembly 后 compile 仍绿 |
| **T2** | P2 ERC20 迁移 | Deposit/Payment 引 `SafeERC20`+`IERC20` + `using SafeERC20 for IERC20`；`usdt` state 经 **initialize 注入**（Deposit.initialize(userRegistry,usdt) / Payment.initialize(feeManager,platformWallet,usdt,serviceManager)）；`deposit(uint256)` 去 payable + 前置 safeTransferFrom；`withdraw` 改 safeTransfer（CEI）；`MIN_DEPOSIT=10*10**usdt.decimals()`（不硬编码）；禁 fee-on-transfer（注释/handoff） | `Deposit.sol`（deposit/withdraw/initialize/usdt state/MIN_DEPOSIT）、`Payment.sol`（initialize/usdt state）、`IDeposit.deposit` 去 payable、`IPayment.payBill` 去 payable | T1.5 | compile 绿；MIN-01/02、ERC-01/02/03、DEC-01 单测通过 |
| **T3** | P3 分账 | `payBill` 去 payable + CEI 先置 `isPaid=true` + 两段 `safeTransferFrom`（amount→operator.paymentAddress、fee→platformWallet）+ `if(fee>0)` 0-fee 跳过 + `require(paymentAddress!=0)`；`createBill` fail-fast `require(operator.paymentAddress!=0)`；ServiceManager 新增 `setOperatorPaymentAddress(id,addr) onlyOwner` + 零地址 require | `Payment.sol`（payBill/createBill）、`ServiceManager.sol`（setOperatorPaymentAddress） | T2 | compile 绿；PAY-01~04 单测通过 |
| **T4** | P4 自动发卡+计价 | Deposit 抽 internal `_mintFor(user)`（三校验+`trafficCardNFT.mint(user,trafficCardQuota,uri)`）；`mintTrafficCard(onlyOwner)` 重构为 `_mintFor` 薄壳；`issueMonthlyTrafficCards(onlyOracle)` 循环 `_mintFor` + 校验失败 `continue`（不整批 revert）+ `nonReentrant`（A1）；`dataAmount=trafficCardQuota` 固定（删 L76 `_deposits/100000`）；禁用 NFT.mintBatch（发卡链路不接入） | `Deposit.sol`（_mintFor/mintTrafficCard/issueMonthlyTrafficCards）、`TrafficCardNFT.sol`（mintBatch 不接入） | T3 | compile 绿；ISS-01~05 单测通过 |
| **T5** | P5 部署 | hardhat.config 加 `arbitrum_sepolia(421614)`；deploy.ts 步骤 0 部署 `MockUSDT(decimals=6)`；补 wiring（initialize 注入 usdt/serviceManager 即 **deployProxy 参数数组同步**、`payment.setOracle`、`oracle.setPayment`、循环 `setOperatorPaymentAddress`，保留现有 setTrafficCardNFT/setDeposit/setOracle/transferOwnership 等）+ 部署后断言校验；deployments.json 加 `usdt`/`usdtDecimals`/`storageLayout`/`abiHash`；产 handoff（ABI+selector 清单+精度语义）给后端(2/3) | `contracts/mocks/MockUSDT.sol`（新建）、`hardhat.config.ts`、`scripts/deploy.ts`、`deployments/<net>.json`、handoff 文档 | T4 | 31337 + 421614 两处 deploy 成功；断言 `payment.oracle()==oracle`/`oracle.payment()==payment`/`deposit.oracle()==oracle`/`NFT.owner()==deposit`/active operator.paymentAddress≠0 全过 |
| **T6** | P6 补测 | §八 全清单：MIN-01/02、ERC-01~03、PAY-01~04、ISS-01~05、MS-01~03、ATC-01/02、USDT-01/02、DEC-01、GAS-01、REG-01/02；现有 26 it（金额改 USDT 精度）回归不破；旧 `{value}` payable 调用重写为 approve+deposit(amount)/payBill(id) | `test/erc20.ts`（新建或扩 `test/linkworld.ts`）、`contracts/mocks/`（非标 USDT mock：无返回值/返回 false） | T5（TDD 可前置到各阶段） | 全部新增用例 + 26 回归 it 绿 |

> 注：T6 标注「TDD 可前置」——按 design §九，各阶段对应测试（MIN/ERC 入 T2、PAY 入 T3、ISS 入 T4、MS/ATC 入 T5）可在对应任务内 TDD 先写；T6 是兜底全量回归 + 集成/非标/压测的收口。
> ServiceManager `requiredDeposit` / Deposit `_operatorRequiredDeposit`：本轮明确废弃，不读不校验，**不在任何任务范围内**。

---

## 2. 受影响合约与接口变更

| 合约/接口 | 变更 | 来源(design §) | 风险 |
|-----------|------|----------------|------|
| `Oracle.sol` | 声明 `UsageDataSubmitted` 事件；`deposit`/`payment` `address`→`IDeposit`/`IPayment`；新增 `setPayment(address) onlyOwner`；`monthlySettlement` 改 `(users,operatorIds,amounts[])` 删 L68 求和；收敛内联接口到 interfaces/ | §4.4、§5.4、§5.6①②⑤ | 编译错误①根因；删求和=B1 闭合，量纲修复 |
| `interfaces/IDeposit.sol` | `deposit()`去 payable→`deposit(uint256)`；补 `issueMonthlyTrafficCards(address[])` | §4.1、§5.6⑤、B5 | selector 变更冲击后端 ABI；漏补→收敛后编译不过 |
| `interfaces/IPayment.sol` | `payBill()`去 payable→`payBill(uint256)`；补 `applyTrafficCardToBill(uint256)` | §4.2、§5.6⑤、B5 | 同上 |
| `Payment.sol` | `createBill` onlyOwner→**onlyOracle**(B2)+fail-fast paymentAddress≠0；`payBill` 去 payable+CEI+两段 safeTransferFrom+0-fee 跳过；`applyTrafficCardToBill` 受限桩(v2-B)；`usdt`/`serviceManager` state initialize 注入；guard 选型(T1.5) | §4.2、§5.2、§7.0、B2/B7/v2-B | 权限链 deploy 跑通依赖 §7.0 wiring；ERC20 重入 |
| `Deposit.sol` | `deposit(uint256)`+`withdraw` safeTransfer；`usdt` state initialize 注入；`MIN_DEPOSIT=10*10**decimals`；internal `_mintFor`+`mintTrafficCard` 薄壳+`issueMonthlyTrafficCards` onlyOracle/continue/nonReentrant；`dataAmount=trafficCardQuota`(删 _deposits/100000) | §4.1、§5.1、B3/B4/v2-C/A1 | 精度雷区；发卡幂等+不整批回滚；A1 重入 |
| `TrafficCardNFT.sol` | 删 `DeductionCredit` mapping + `onlyDeposit` dead code；mintBatch 不接入发卡链路 | §4.5、§5.3、编译错误②④ | 删 dead code（scan 确认未用） |
| `ServiceManager.sol` | 新增 `setOperatorPaymentAddress(id,addr) onlyOwner` + 零地址 require；11 运营商 paymentAddress 部署后注入 | §4.3、§5.2 | 全 address(0) 必先补，否则分账转黑洞 |
| `contracts/mocks/MockUSDT.sol` | 新建：OZ ERC20 + `decimals()=6` + public `mint` + symbol "USDT" | §5.5 | 仅测试网/本地 |
| `hardhat.config.ts` | 加 `arbitrum_sepolia(421614)`；evmVersion/viaIR 维持（待 T1.5 确认） | §7.1 | TSTORE 兼容（T1.5 实测） |
| `scripts/deploy.ts` | 步骤0 MockUSDT；initialize 参数数组同步；补 setOracle/setPayment/setOperatorPaymentAddress 循环；断言校验；handoff | §7.2、§7.3、arch ⚠️4 | 授权拓扑漏一项→自动结算链 revert |
| `deployments/<net>.json` | 加 `usdt`/`usdtDecimals`/`storageLayout`/`abiHash` | §7.3 | storage layout 冻结作升级基线 |

---

## 3. 验收标准映射

> 把 requirement §五 A（合约）与 arch-review §七 5 条带入 implement 的 ⚠️ 加固项，映射到任务与验收方式。

| 验收点 | 对应任务 | 验收方式 |
|--------|----------|----------|
| **req A.1** Deposit/Payment 改 ERC20，无 payable/msg.value 残留 | T2/T3 | grep 核查无 `payable`/`msg.value` + ERC-02/PAY-02 单测 |
| **req A.2** Deposit 强约束 amount≥10 USDT（<10 拒绝、=10 通过） | T2 | MIN-01（9.999999 revert）/ MIN-02（10.000000 成功）单测 |
| **req A.3** 31337 + 421614 两处部署成功，地址写 deployments/ | T5 | `deploy:local` + `deploy:arbitrum-sepolia` 实跑，产出 json |
| **req A.4** mock USDT 已部署，精度/符号与前端一致 | T5 | MockUSDT 部署 + DEC-01（decimals=6）单测 |
| **arch ⚠️1 / A1** 发卡路径重入：`_mintFor`/`mintTrafficCard`/`issueMonthlyTrafficCards` 加 nonReentrant 或 CEI | T4 | 代码核查 nonReentrant + ISS-04 混合批不整批 revert |
| **arch ⚠️2 / A3** 删 Payment L33-38 无效 `_reentrancyGuardInit` + initialize 调用 | T1.5 | grep 无 `_reentrancyGuardInit`；删后 compile 绿 |
| **arch ⚠️3** Arbitrum TSTORE 实测 + ReentrancyGuardUpgradeable 降级 fallback | T1.5 | 421614 实测 payBill 成功 或 降级方案部署成功 |
| **arch ⚠️4** deploy.ts initializer 参数个数同步 deployProxy 参数数组 + MockUSDT 步骤0 | T5 | deploy 成功 + 断言 usdt 注入正确（safeTransferFrom 不 panic） |
| **arch ⚠️5** monthlySettlement 旧喂价旁路去留与后端(2/3)对齐 | T1（标注）/handoff | design §4.4 dataUsages/callUsages 仅 emit UsageDataSubmitted；handoff 注明 |
| **B1/v2-A** Oracle 不计价，amount=链下 USDT 金额 | T1 | MS-01：账单 amount==amounts[i]（不求和）单测 |
| **B2** createBill onlyOracle + 权限链 deploy 跑通 | T1+T5 | MS-02（非授权 revert）+ MS-03（三链不 revert）单测 + deploy 断言 |
| **B3** 发卡走 _mintFor 不撞 onlyOwner + continue 跳过幂等 | T4 | ISS-02/03/04 单测 |
| **B4/v2-C** dataAmount=trafficCardQuota 固定，无 _deposits/100000 | T4 | grep 无 `/100000` + ISS-05 单测 |
| **B5** interfaces/ 补两声明 + 连锁编译同步 | T1 | compile 绿（收敛后无 undefined 方法） |
| **B6** monthlySettlement + applyTrafficCardToBill 测试覆盖 | T5/T6 | MS-01~03 + ATC-01/02 单测 |
| **B7** 仅标准 ERC20/禁 fee-on-transfer/黑名单降级 | T2/T3+T6 | USDT-01（无返回值入账）/ USDT-02（返回 false revert）单测 |
| **精度不硬编码** 全仓无 `1e18`/`* 10**18`/`* 10**6` 字面 | T2/T4 | grep 核查 + DEC-01/ISS-05 单测 |
| **回归** 现有 26 it 不破 + payable 旧用例重写 | T6 | REG-01（续期不变量）/REG-02（未注册 revert，重写）+ 26 it 绿 |
| **批量 gas 上限** 量出单批安全 N 写入 handoff | T6 | GAS-01 压测 |

---

## 4. 风险与 implement 铁律

| 铁律 | 内容 |
|------|------|
| **串行** | T1→T6 严格串行，一个完成审查后再派下一个（global CLAUDE.md：implement 阶段始终串行）。 |
| **P1 compile gate** | T1 的 ①–④ 子项**全做完** `hardhat compile` 才会绿（②类型改完缺接口方法、桩未实现前都不过），无「②④后即绿」中间态。不绿不进 T1.5。 |
| **P1.5 前移** | TSTORE 兼容实测在 T2 之前完成；失败先切 `ReentrancyGuardUpgradeable` 降级再进 T2，避免 P5 才返工。 |
| **wiring 兜底** | T5 deploy 后必须跑断言（§7.0.3 三条链前置）；漏 `payment.setOracle`/`oracle.setPayment`/`setOperatorPaymentAddress` 任一项 → 自动结算 revert。 |
| **零资损面铁律** | applyTrafficCardToBill 本轮**受限桩不转资金**（v2-B）；createBill fail-fast + payBill 零地址 require + CEI + nonReentrant + SafeERC20 全覆盖。 |
| **handoff 强制** | T5 必产给后端(2/3) 的 handoff：冻结 ABI + selector 变更清单（deposit/payBill 改 selector）+ 金额精度语义 + storage layout 冻结 + 批量分批约定 + applyTrafficCardToBill 桩语义说明。 |
| **每任务 checkpoint** | implement 每个 Task 完成后立即写 checkpoint（global MEMORY 规则），不批量补写。 |
