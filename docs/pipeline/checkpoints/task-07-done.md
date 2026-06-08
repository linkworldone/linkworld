# Task 07 — T6 补测扫尾 (P6)

> 子项目 contracts(1/3) | 状态 DONE | 2026-06-08 | implement 阶段最后一棒

### 产出文件
- packages/contracts/contracts/mocks/NonStandardUSDT.sol（新建：NoReturnUSDT 无返回值 + FalseReturnUSDT 返回 false，仅测试桩，验 SafeERC20）
- packages/contracts/test/t6-audit.ts（新建：T6 补测 9 用例）
- docs/design/linkworld-contracts/handoff-backend.md（回填 §5.1 GAS-01 压测结论：批量安全上限）

### git commit
be41424 test: 合约 T6 补测扫尾（GAS-01 压测+SafeERC20 非标+续期回归）

### TDD
补测阶段，测试即产物。新增 9 用例全绿，未发现合约 bug（未改任何合约业务逻辑）。

### 测试结果
hardhat compile 0 error。hardhat test：75 passing / 0 failing（原 66 + 新 9：USDT-01/01b/02/02b、REG-01a/01b、ERC-04、GAS-01/01b）。主 Agent 已独立复跑确认 75 passing。

### code-simplifier
仅新增测试 + 测试桩 mock，无业务代码改动。

### spec review
按 design §八补测清单 + arch-review（GAS-01/USDT-01·02/REG）执行：design §八全清单闭合（UPG-01 N-A，本轮 fresh deploy 不升级）；§十复审重点 B6(MS/ATC)/B7(非标 USDT)/批量 gas 上限均落测并回填 handoff。

### 设计还原
合约无 UI。design §八验收测试清单逐项落地（见覆盖审计表）。

### 复用检查
复用现有 test/{linkworld,erc20,payment}.ts 已有用例（不重复，仅补真缺口）；复用 MockUSDT；新建非标 mock 仅测试。

### 设计稿对照
数值对照：最终 75 passing（66+9）vs 无回归 ✅；GAS-01 实测 issueMonthlyTrafficCards per-user≈230k → N≤50、monthlySettlement 全链路 per-user≈432k → N≤25（按 15M 单批预算）已回填 handoff §5.1 ✅；非标 USDT（无返回值/返回 false）SafeERC20 行为正确（入账/revert）vs §6.1 ✅；REG-01 续期不变量（未到期叠加/到期重置）vs §3.1 ✅；compile error 0 ✅。

### 新增组件
新增测试 mock：NonStandardUSDT.sol（NoReturnUSDT/FalseReturnUSDT）。无新增业务合约。

### 新增色值
无（合约任务）。

### implement 阶段总结（T1-T6 全完成）
T1 修编译(e4bd81d) → T1.5 删无效 assembly(b55ef16) → T2 Deposit ERC20(5fcccff) → T3 Payment 分账(c6838e8) → T4 自动发卡+计价(6bd6280) → T5 部署+handoff(ebc4b41) → T6 补测(be41424)。最终 75 passing，无合约 bug。可进 test 阶段。
遗留（非 implement 范畴）：421614 真·上链 + TSTORE 链上实测需配 DEPLOYER_PRIVATE_KEY 后执行（handoff §10）。
