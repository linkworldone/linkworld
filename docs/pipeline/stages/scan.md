# Stage: scan — 合约基线扫描（子项目 contracts 1/3）

> **状态**: completed (DONE_WITH_CONCERNS) | **日期**: 2026-06-08 | **Gate**: 0 | **子项目**: contracts(1/3)
> 对象：packages/contracts（全栈反转后，合约从 origin 合并落地，需重新摸底）

## 产出文件（4 份基线，commit 4249a92）
| 文件 | 内容 |
|------|------|
| docs/design/linkworld-contracts/project-scan.md | 结构/网络/编译器/UUPS proxy/部署现状/依赖 |
| docs/design/linkworld-contracts/contracts-inventory.md | 7 合约逐个清单（职责/函数/事件/proxy）|
| docs/design/linkworld-contracts/erc20-migration-surface.md | ERC20 改写面（payable/msg.value 位置 + Deposit/Payment/FeeManager/TrafficCardNFT 现状）|
| docs/design/linkworld-contracts/test-deploy-baseline.md | 测试覆盖 + 部署脚本/产物基线 |

## 🔴 已验证的关键事实：当前合约编译不通过
主 Agent 实跑 `npx hardhat compile` 确认失败：
```
DeclarationError: Undeclared identifier.
  --> contracts/Oracle.sol:71:22  emit UsageDataSubmitted(...)   ← 事件未声明
Error HH600: Compilation failed
```
合并自 origin 的合约代码（commit「rom upgrade test commit 0606」）是**编译不过的半成品**。ERC20 迁移前，design/plan 必须先把合约修到可编译（这是地基前置）。

## 关键发现（给 design 阶段）
1. **ERC20 改写面收敛**：payable/msg.value/.call{value} 仅在 Deposit.sol(5处)+Payment.sol(4处)+IDeposit/IPayment 接口各1处。全仓无 ERC20/SafeERC20 引用；OZ ^5.6.1 自带 IERC20+SafeERC20，无需加依赖。
2. **7 合约全部 UUPS 可升级**：已有 .openzeppelin/unknown-16602.json manifest。改写应走 upgradeProxy + storage layout 只追加规则（新增 usdt state 加末尾）；Arbitrum 是全新部署无此约束。
3. **Arbitrum 421614 配置完全缺失**：hardhat.config.ts 只有 31337/16600/16602，需新增网络+RPC+部署脚本+mock USDT。
4. **无 10 USDT 下限**：Deposit 当前仅 require(msg.value>0)；锁仓 30 天+续期逻辑可保留，只换资金通道。
5. **USDT 精度雷区**：现状全按 18 位。Deposit.mintTrafficCard 的 _deposits/100000、ServiceManager 的 ether requiredDeposit、测试 parseEther、待加 MIN_DEPOSIT 都受 6 位精度冲击。建议读 usdt.decimals() 不硬编码。
6. **Payment 运营商分账缺口**：payBill 只转 platformFee 到 platformWallet，bill.amount 无出口；ServiceManager paymentAddress 全 address(0)。ERC20 改写须定 amount 去向。
7. **测试覆盖严重不足**：withdraw/payBill/mintTrafficCard/锁仓时序/最小额/UUPS 升级全无测试，且全用 18 位 parseEther。验收 A.2 的 <10 拒绝/=10 通过单测目前不存在。

## ⚠️ 需 design/arch 复核的遗留问题
- Oracle.sol 把 payment/deposit 声明为 address 却调其方法；applyTrafficCardToBill 在 Payment.sol 未实现；Oracle 缺 setPayment。当前编译已失败（见上），artifacts/ 与源码不同步。
- Deposit.issueMonthlyTrafficCards（PRD R7 自动发卡入口）当前是空实现，真正自动发放链路不存在，arch 需补设计。
- 部署脚本 wiring 缺 payment.setOracle(...)。
- 编译器用 evmVersion: cancun + viaIR，Payment 依赖 ReentrancyGuardTransient(transient storage)——迁 Arbitrum Sepolia 前需确认 cancun opcode 支持。

## 移交 design
带着 PRD §三①/§五A/§六 + 上述发现，先确定"修编译→ERC20 改写→Arbitrum 部署+mock USDT→补测"的合约设计方案与状态机/接口签名。
