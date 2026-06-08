# Stage: plan — 任务拆解与增量比对（子项目 contracts 1/3）

> **状态**: completed（用户已审批） | **日期**: 2026-06-08 | **Gate**: 2 | **子项目**: contracts(1/3)
> 产出：docs/pipeline/stages/plan-review.md（人读）+ plan-review.json（hook 读，已过 advance 校验：7 pages，每页含 components+colors）
> 领域适配：plan 模板为 web/UI 形状（Figma/页面/色值），本子项为 Solidity 合约——已适配为「任务=page、components=受影响合约/函数、colors=[]」。

## 7 个串行任务（implement 铁律：串行，一个完成审查再下一个；P1 compile gate 不绿不进 P2）
| 任务 | 阶段 | 范围 | gate |
|------|------|------|------|
| T1 | P1 | 修编译：Oracle 补 UsageDataSubmitted 事件/state 改接口类型/setPayment/收敛接口；IDeposit 补 issueMonthlyTrafficCards、IPayment 补 applyTrafficCardToBill；删 TrafficCardNFT DeductionCredit dead code；Payment createBill 改 onlyOracle + applyTrafficCardToBill 受限桩(onlyOracle/不转资金) | hardhat compile 绿 |
| T1.5 | P1.5 | ReentrancyGuard 探针（前移）：421614 实测 TSTORE 兼容；删 Payment L33-38 无效 _reentrancyGuardInit assembly；不兼容降级 ReentrancyGuardUpgradeable | 实测/降级定案 |
| T2 | P2 | ERC20 迁移：SafeERC20 + usdt initialize 注入 + deposit(amount)/withdraw safeTransfer + MIN_DEPOSIT=10*10**usdt.decimals() | 单测 MIN/ERC |
| T3 | P3 | 分账：payBill 两段 safeTransferFrom(amount→operator,fee→platform) + CEI + 0-fee 跳过 + createBill fail-fast；ServiceManager setOperatorPaymentAddress + 零地址 require | 单测 PAY |
| T4 | P4 | 自动发卡+计价：internal _mintFor（不撞 onlyOwner）+ issueMonthlyTrafficCards onlyOracle + continue 跳过 + nonReentrant(A1)；dataAmount=固定 trafficCardQuota（删 _deposits/100000）；monthlySettlement 改签名 (users,operatorIds,amounts[]) 删 usage 求和；禁用 mintBatch | 单测 ISS/MS |
| T5 | P5 | 部署：hardhat.config 加 421614；MockUSDT(decimals=6) 步骤0；deploy.ts 补 wiring(setUSDT×2/setServiceManager/payment.setOracle/oracle.setPayment/循环 setOperatorPaymentAddress)+initializer 参数数组同步+部署后断言；deployments.json 加 usdt/usdtDecimals+storageLayout 冻结；产后端 handoff(ABI+selector+精度语义) | 31337+421614 部署成功 |
| T6 | P6 | 补测（TDD 贯穿）：MIN-01/02、ERC-01/02、PAY-01~05、ISS-01~05、MS-01~03、ATC-01/02、USDT-01/02、DEC-01、REG-01/02、GAS-01；现有 26 it 回归不破 | 全绿 |

## arch-review 5 项 ⚠️ 映射
A1 发卡重入→T4 nonReentrant；A3 删无效 assembly→T1.5；TSTORE 实测→T1.5；deploy initializer 参数同步→T5；monthlySettlement 旧喂价旁路对齐→T1/T4（amounts[] 改签名时处置 _monthlyUsage/submitUsage 旁路）。

## 范围确认
- ServiceManager requiredDeposit/_operatorRequiredDeposit 本轮废弃不校验。
- 不上主网；仅 31337 + Arbitrum Sepolia 421614。

## 移交 implement
从 T1 开始严格串行 + TDD。每个任务完成 → 主 Agent 审查 + 写 checkpoint → 再派下一个（implement 阶段 checkpoint 必写，不可批量补）。
