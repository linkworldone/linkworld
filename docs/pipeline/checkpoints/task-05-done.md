# Task 05 — T4 自动发卡 + 计价修正 (P4)

> 子项目 contracts(1/3) | 状态 DONE | 2026-06-08

### 产出文件
- packages/contracts/contracts/Deposit.sol（引入 ReentrancyGuardTransient；新增 internal _canMint 三校验 + _mintFor 固定 trafficCardQuota 发卡，删 _deposits/100000；mintTrafficCard 重构为 onlyOwner+nonReentrant 薄壳；issueMonthlyTrafficCards 真实实现 onlyOracle+nonReentrant+循环 continue 跳过不整批 revert）
- packages/contracts/contracts/Oracle.sol（monthlySettlement 改签名 (users,operatorIds,amounts[])，删 dataUsage+callUsage 求和，createBill 传 amounts[i]；UsageDataSubmitted 改 (user,operatorId,amount) 仅喂价；length 校验保留）
- packages/contracts/contracts/TrafficCardNFT.sol（mintBatch 加 dead-code 勿用注释）
- packages/contracts/test/linkworld.ts（新增 "T4 自动发卡+计价" describe，11 用例）

### git commit
6bd6280 feat: 合约 T4 自动发卡（_mintFor+onlyOracle+幂等+nonReentrant）+ 固定 quota + monthlySettlement 计价改 amounts[]

### TDD
先红后绿：测试先写，实现前 10 failing（MS 旧4参签名 no matching fragment、ISS/DEC 未发卡 Card not found）→ 实现后 66 passing 全转绿。

### 测试结果
hardhat compile 0 error（仅 2 个 pre-existing warning，非本轮引入）。hardhat test：66 passing / 0 failing（基线 55 + 新 11：ISS-01~05、DEC、DEC-2、MS-01~04）。主 Agent 已独立复跑确认 66 passing（ISS-03 幂等/ISS-04 混合批/ISS-05 固定 quota/MS-01 计价 均绿）。

### code-simplifier
_canMint/_mintFor 抽取消除 mintTrafficCard 与 issueMonthlyTrafficCards 的校验/发卡逻辑重复（DRY）；monthlySettlement 删求和后逻辑更直白。

### spec review
严格按 design §3.1/§4.1/§4.4/v2-A/v2-C + arch-review B1/B3/A1 执行：固定 quota（删 _deposits/100000）、_mintFor internal 不撞 onlyOwner、continue 跳过不整批 revert、幂等 getUserCardCount==0、issueMonthlyTrafficCards+mintTrafficCard 加 nonReentrant、monthlySettlement amounts[] 无求和、禁用 mintBatch。无偏差。

### 设计还原
合约无 UI。design §3.1 状态机落地：锁仓满→issueMonthlyTrafficCards(三校验)→_mintFor→mint NFT(quota 固定)→持卡。§4.4 计价：Oracle 不计价，createBill 收后端算好的 amounts[i]。

### 复用检查
复用 OZ ReentrancyGuardTransient（与 T3 Payment 同栈）、TrafficCardNFT.mint/getUserCardCount、trafficCardQuota state（initialize 已设 100MB）；_mintFor 被 mintTrafficCard 与 issueMonthlyTrafficCards 共用，无重复实现。

### 设计稿对照
数值对照：dataAmount=trafficCardQuota=100MB 固定（ISS-05/DEC 验证存 10 vs 20 USDT 发卡 quota 一致）vs v2-C ✅；发卡三校验 now≥lockExpiry && deposits>0 && cardCount==0 vs §3.1 ✅；monthlySettlement createBill amount==amounts[i]（MS-01，不再 usage 求和）vs v2-A ✅；nonReentrant 覆盖发卡 2 入口 vs A1 ✅；测试 66 passing（55+11）vs 无回归 ✅；compile error 0 ✅。

### 新增组件
无新增合约。新增 Deposit 函数：internal _canMint、_mintFor；issueMonthlyTrafficCards 从空实现→真实实现。

### 新增色值
无（合约任务）。

### 遗留（带入 T5/T6）
- nonReentrant 为发卡路径主防护（NFT _userCardCount++ 在 _safeMint 回调后，依赖 guard 而非 CEI 顺序）——T5 需在 421614 实测 transient guard 真生效。
- monthlySettlement 签名定型 (users,operatorIds,amounts[]) — 给后端子项 2/3 的 handoff。
- Oracle _monthlyUsage/_latestUsage/UsageInfo 成 dead state（新签名不读），未删以保 UUPS layout，后续 Round 评估。
- T5：deploy.ts 补 oracle.setPayment + initialize 参数同步。
