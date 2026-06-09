# Stage test — 合约子项目(1/3) test 阶段 checkpoint

> 子项目 contracts(1/3) | 状态 **DONE（全绿可进 review）** | 2026-06-09 | test runner（只验证+产文档，不改合约）

### 产出文件
- docs/pipeline/stages/test.md（test 阶段人读报告：结论 + Phase 1-6 证据 + 验收映射 + 已知限制 + 可进 review 判定）
- docs/pipeline/checkpoints/stage-test-done.md（本文件）

> 未新装任何依赖：solidity-coverage 0.8.17 已随 @nomicfoundation/hardhat-toolbox 提供，无 package.json/lock 改动。

### Phase 1-6 结果

| Phase | 门 | 结果 |
|-------|----|----|
| 1 编译 | `hardhat clean && compile` | **0 error**，56 文件编译成功（cancun）。2 warning：OZ TransientSlot EIP-1153 提示（库代码，reentrancy guard 用法安全）、Oracle.sol:99 submitUsage 可 pure（风格，非资损） |
| 2 测试 | `hardhat test` | **75 passing / 0 failing**（4s）。8 套件全绿，无 failing |
| 3 覆盖率 | `hardhat coverage`（跑通，未受 viaIR 限制） | 全局 **89.36% stmts / 88.84% lines / 51.18% branch / 75.34% funcs**；资损核心 Deposit/Payment **100% stmts**、Deposit/FeeManager/UserRegistry 100% lines。远超 ≥60% 目标 |
| 4 部署冒烟 | `deploy.ts --network hardhat` | 7 proxy + MockUSDT(decimals=6) 部署；**All wiring assertions passed (B2/B3 chains ready)**；写 deployments/hardhat.json |
| 5 验收映射 | §五A/D + design §八 | 逐条映射到真实用例（见下） |
| 6 grep 自检 | 精度字面量 / payable / SafeERC20 | 精度字面量 **0**、payable/msg.value **0**、SafeERC20 全覆盖资金转移点 |

### passing 数 / 覆盖率
- **75 passing / 0 failing**（原 66 + T6 补测 9）。本次 test 阶段独立复跑确认，与 T6 checkpoint 一致。
- 覆盖率（solidity-coverage 实跑数字，非验收清单替代）：全局 89.36% stmts / 88.84% lines；核心合约 Deposit 100/100、Payment 100/93.48、FeeManager 100/100、UserRegistry 100/100；interfaces/ 全 100。偏低项 TrafficCardNFT 56% stmts（owner-only admin 薄壳，非资损主路径，核心 mint/getCardInfo 已测）。

### 验收映射摘要

**§五A 合约**：A.1 无 payable（grep=0）✅ / A.2 ≥10 USDT（MIN-01/02）✅ / A.3 部署+deployments（Phase4 hardhat ✅、421614 配置就绪待 key ⏸）/ A.4 mock USDT decimals=6（DEC-01）✅

**§五D 端到端（合约相关）**：D.12 链上 fee（PAY-04）✅ / D.13 approve 两段式（ERC-01/02、PAY-01/02、PAY-05 原子性）✅ / D.14 锁仓提取（ERC-03、REG-01a/b）✅ / D.15 NFT 自动发放（ISS-01~05、DEC、ATC-01/02）✅ / D.17 端到端（MS-03 链路通 + Phase4；421614 上链待 key ⏸）

**design §八补测清单**：MIN/ERC/PAY/ISS/DEC/MS/ATC/USDT/GAS/REG 全 ✅；UPG-01 N-A（fresh deploy 不升级）；旧 26 it USDT 精度重写后不回归 ✅

**design v2 §十复审重点**：MS-01~03 / ATC-01/02 / USDT-01/02 / ISS-04 **全绿** ✅

**arch-review §七加固项**：A1 重入CEI✅ / A3 删无效init✅ / deploy.ts同步✅ / 旁路收敛✅ / Arbitrum transient 实测 ⏸（待 key，本地 cancun 全绿）

### 已知限制（不算失败）
- 421614 真·上链 + TSTORE 链上实测：本环境无 DEPLOYER_PRIVATE_KEY/RPC，无法执行。本地 31337/hardhat 全绿、hardhat.config arbitrum_sepolia 配置就绪（缺 key 优雅降级）。影响 §五D-17 / arch-review §七-3 的「421614 跑通」部分，记 handoff §10 待办，不阻塞本阶段。

### grep 自检结果
- 硬编码精度字面量（`10**18`/`1e18`/`10**6`/`/100000`/`parseEther`）：**0 命中**（MIN_DEPOSIT=`10 * 10 ** decimals()` 动态派生；trafficCardQuota=100MB 是数据量常量非精度）
- payable / msg.value：**0 命中**（已全 ERC20 化）
- SafeERC20：Deposit safeTransferFrom/safeTransfer + Payment 两段 safeTransferFrom 全覆盖，均 `using SafeERC20 for IERC20`

### 三个自检
1. **承诺交付什么**：全量验证（编译/测试/覆盖率/部署冒烟/验收映射/grep 自检）+ test.md + 本 checkpoint。
2. **交付物真实存在吗**：两份文档已写入；75 passing、覆盖率表、部署 wiring 输出均本次实跑取得。
3. **正确连接吗**：映射逐条指向真实用例名/grep/断言；覆盖率来自实跑；deploy wiring 断言全过证明合约 wiring 正确。

### 结论
**DONE — 全绿，可进 review。** 本地质量门全过，资损 grep 全清，核心合约 100% stmts。唯一 ⏸ 为 421614 上链（待 key，handoff 待办），不阻塞 review。
