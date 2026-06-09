# test 阶段报告 — LinkWorld 合约子项目（1/3）

> 状态：DONE（全绿，可进 review） | 日期：2026-06-09 | 子项目：packages/contracts | Round：1
> 角色：test runner（只跑验证 + 产出文档，不改任何合约 / deploy.ts / hardhat.config 业务内容）
> 资损敏感合约，标准从严。

---

## 〇、结论

**本地 31337/hardhat 全绿，质量门全过，可进 review 阶段。**

- Phase 1 编译门：`hardhat compile` **0 error**，2 个非阻塞 warning（均为已知/无害）。
- Phase 2 测试门：`hardhat test` **75 passing / 0 failing**。
- Phase 3 覆盖率门：`solidity-coverage` 跑通，**全局 89.36% stmts / 88.84% lines / 51.18% branch / 75.34% funcs**，资损核心合约 Deposit/Payment **100% stmts**，远超 ≥60% 目标。
- Phase 4 部署冒烟：`deploy.ts --network hardhat` 全 wiring 断言通过（B2/B3 链就绪）。
- Phase 5 验收映射：requirement §五A/D + design §八补测清单逐条已被测试覆盖（见映射表）。
- Phase 6 grep 自检：硬编码精度字面量 **0**、payable/msg.value 残留 **0**、SafeERC20 全覆盖资金转移点。

**已知限制**：Arbitrum Sepolia 421614 真·上链 + TSTORE 链上实测因本环境无 `DEPLOYER_PRIVATE_KEY`/RPC 无法执行；本地 hardhat（cancun，含 transient storage）全绿、配置就绪（见 §六）。不阻塞本阶段。

---

## 一、Phase 1 — 编译门

命令：`npx hardhat clean && npx hardhat compile`（强制全量重编）

结果：**0 error**，编译 56 个 Solidity 文件成功（evm target: cancun），生成 154 typings。

Warning（2 个，均非阻塞）：
1. `@openzeppelin/contracts/utils/TransientSlot.sol:108` — EIP-1153 transient storage 组合性提示。OZ 库代码，官方明确「用于在调用结束时清除的 reentrancy guard 是安全的」。本项目正是此用法（ReentrancyGuardTransient），无害。
2. `contracts/Oracle.sol:99 submitUsage` — state mutability 可收紧为 `pure`。纯风格提示，不影响逻辑、非资损面。

---

## 二、Phase 2 — 测试门

命令：`npx hardhat test`

结果：**75 passing / 0 failing（4s）**。8 个测试套件：

| 套件 | 文件 | 用例数 | 覆盖点 |
|------|------|--------|--------|
| Deposit ERC20 (T2) | erc20.ts | 9 | DEC-01 精度、MIN-01/02 最小额、ERC-01/02/02b approve+transferFrom、ERC-03 锁仓提取、REG-02 注册校验 |
| LinkWorld Contracts Tests | linkworld.ts | 39 | FeeManager FE-01~09、UserRegistry UR-01~04、ServiceManager SM-01~03、Deposit DP-01~03、Payment PM-01~02、Oracle OR-01~02、TrafficCardNFT TC-01~05、T4 自动发卡 ISS-01~05/DEC/DEC-2、计价 MS-01~04 |
| Payment ERC20 分账 (T3) | payment.ts | 14 | createBill 权限+fail-fast、PAY-01/02/02b/04/05 分账、PAY-03 零地址、0-FEE 跳过、SM-PAY-01~03、ATC-01/02 受限桩 |
| B7 非标 USDT SafeERC20 (T6) | t6-audit.ts | 4 | USDT-01/01b 无返回值入账/退回、USDT-02/02b 返回 false revert |
| REG-01 锁仓续期不变量 (T6) | t6-audit.ts | 2 | REG-01a 未到期叠加 +30d、REG-01b 到期重置 now+30d |
| ERC withdraw 边界 (T6) | t6-audit.ts | 1 | ERC-04 无存款 withdraw revert |
| GAS-01 批量 gas 压测 (T6) | t6-audit.ts | 2 | GAS-01 发卡上限 N≤50、GAS-01b 全链路上限 N≤25 |
| UserRegistry (旧) | UserRegistry.test.ts | （含于 linkworld 统计内的独立旧文件） | 注册基线 |

GAS-01 实测（本次复跑）：
- `issueMonthlyTrafficCards`：per-user ≈ 237k gas，15M 单批预算 → 安全上限 N≈63（建议 N≤50）。
- `monthlySettlement` 全链路：per-user ≈ 432k gas → 安全上限 N≈34（建议 N≤25）。
- 与 T6 checkpoint / handoff §5.1 结论一致。

---

## 三、Phase 3 — 覆盖率门

命令：`npx hardhat coverage`（solidity-coverage 0.8.17，随 hardhat-toolbox 提供，**无需新装依赖**）

> viaIR + cancun 配置下 coverage 跑通，**未受工具限制**，得到真实数字覆盖率。

| File | % Stmts | % Branch | % Funcs | % Lines |
|------|---------|----------|---------|---------|
| **contracts/ (合计)** | **89.33** | **51.23** | **72.58** | **88.36** |
| Deposit.sol | 100 | 70 | 92.31 | 100 |
| Payment.sol | 100 | 57.89 | 66.67 | 93.48 |
| FeeManager.sol | 100 | 60 | 80 | 100 |
| UserRegistry.sol | 100 | 50 | 80 | 100 |
| Oracle.sol | 88 | 50 | 71.43 | 96 |
| ServiceManager.sol | 90.38 | 38.46 | 66.67 | 87.67 |
| TrafficCardNFT.sol | 56 | 27.27 | 54.55 | 54.84 |
| contracts/interfaces/ (全部) | 100 | 100 | 100 | 100 |
| contracts/mocks/MockUSDT.sol | 100 | 100 | 100 | 100 |
| contracts/mocks/NonStandardUSDT.sol | 87.5 | 50 | 87.5 | 94.12 |
| **All files** | **89.36** | **51.18** | **75.34** | **88.84** |

**评估**：
- 资损核心（Deposit/Payment 资金通道、FeeManager 费率、UserRegistry）statements/lines 全部 ≥93%，Deposit/FeeManager/UserRegistry 达 100%。**满足资损合约关键函数高覆盖要求。**
- 全局 89.36% stmts / 88.84% lines **远超 ≥60% 目标**。
- 偏低项说明（均非资损核心，可接受）：
  - **TrafficCardNFT 56% stmts**：未覆盖部分为 owner-only admin（如 setBaseURI、burn/有效期 view 等非自动发卡主路径的薄壳）。核心 mint/CardMinted/getCardInfo 已测（TC-01~05 + 自动发卡 ISS）。
  - **branch 51%**：多为 require/revert 守卫的「正分支已测、异常分支部分未触发」，以及 OZ 继承分支。关键 revert（最小额、注册、权限、零地址、SafeERC20 false）均有专门用例。
- Istanbul 报告已写入 `packages/contracts/coverage/`（gitignore 范围，未纳入提交）。

---

## 四、Phase 4 — 部署冒烟

命令：`npx hardhat run scripts/deploy.ts --network hardhat`

关键输出：
```
MockUSDT: 0x5FbD...0aa3 (decimals=6)
FeeManager / UserRegistry / ServiceManager / TrafficCardNFT / Payment / Deposit / Oracle  全部部署
Contracts linked (incl. payment.setOracle / oracle.setPayment)
OperatorPaymentAddress set for 11 built-in operators
TrafficCardNFT ownership transferred to Deposit
All wiring assertions passed (B2/B3 chains ready)
Addresses written to deployments/hardhat.json
```

结论：7 个 UUPS proxy + MockUSDT(decimals=6) 全部部署成功；§7.0 授权拓扑（payment.setOracle / oracle.setPayment / setOperatorPaymentAddress×11 / NFT ownership→Deposit）全部 wiring 断言通过。`deployments/hardhat.json` 含 7 proxy + impl + usdt + usdtDecimals=6 + abiHash。

---

## 五、Phase 5 — 验收映射

### requirement §五A 合约

| 验收 | 内容 | 测试用例 / 证据 | 状态 |
|------|------|------------------|------|
| A.1 | 无 payable/msg.value 残留 | Phase 6 grep = 0；全资金路径 safeTransferFrom/safeTransfer | ✅ |
| A.2 | amount ≥ 10 USDT 强约束（<10 拒 / =10 通过） | MIN-01（9.999999 revert）、MIN-02（10.0 通过） | ✅ |
| A.3 | 31337 + 421614 两处可部署，写入 deployments/ | Phase 4 hardhat 部署成功 + deployments/hardhat.json；421614 配置就绪（hardhat.config arbitrum_sepolia） | ✅ 本地 / ⏸ 421614 待 key |
| A.4 | mock USDT 已部署，精度/符号与前端一致 | DEC-01（decimals()==6）、deploy 输出 decimals=6、deployments usdtDecimals=6 | ✅ |

### requirement §五D 端到端（合约相关项）

| 验收 | 内容 | 测试用例 / 证据 | 状态 |
|------|------|------------------|------|
| D.12 | 手续费 = FeeManager 链上实读 | PAY-04（链上 fee==calculateFee）、FE-01~09 | ✅ |
| D.13 | approve 两段式（allowance 不足先 approve 再 deposit/pay） | ERC-01/02（deposit）、PAY-01/02（payBill）；两段 safeTransferFrom 原子性 PAY-05 | ✅（合约侧） |
| D.14 | 锁仓未到禁提 + 到期可提（getLockExpiry） | ERC-03（未到 revert / 到期退本金）、REG-01a/b 续期不变量 | ✅（合约侧） |
| D.15 | NFT 自动发放（移除 Admin 发卡，锁仓满自动） | ISS-01~05、DEC/DEC-2（quota 固定与存款解耦）、ATC-01/02 受限桩 | ✅（合约侧） |
| D.17 | 主流程端到端在 31337 + 421614 跑通 | MS-03 端到端链路通（createBill+发卡+applyTrafficCardToBill）、Phase 4 部署；**421614 上链待 key** | ✅ 本地 / ⏸ 421614 |

> D.7~11/16 为 web 视觉/前端项，归 web 子项目（3/3），本子项目 N-A。

### design §八 补测清单逐项

| 编号 | 用例 | 状态 |
|------|------|------|
| MIN-01/02 | 最小额拒绝/通过 | ✅ |
| ERC-01/02/03 | approve 前置 / transferFrom / 锁仓 safeTransfer | ✅ |
| PAY-01/02/03/04 | 支付授权 / 分账正确 / 零地址 / fee 一致 | ✅ |
| ISS-01~05 | 自动发卡权限 / 时序 / 幂等 / B3 混合批 / v2-C 固定 quota | ✅ |
| DEC-01 | mock decimals=6、MIN_DEPOSIT 随 decimals 派生 | ✅ |
| MS-01/02/03 | v2-A amounts[] 不求和 / createBill onlyOracle / B6 集成主路径 | ✅ |
| ATC-01/02 | v2-B 桩权限 + 不转资金 + 存在性校验 | ✅ |
| USDT-01/02 | B7 SafeERC20 无返回值入账 / 返回 false revert | ✅ |
| GAS-01 | 批量 gas 上限（已写入 handoff §5.1） | ✅ |
| REG-01/02 | 续期不变量 / 注册重写 | ✅ |
| 回归 26 it | 旧用例 USDT 精度重写后不回归 | ✅（75 全绿含全部旧用例重写） |
| UPG-01 | upgradeProxy storage（v1 已删） | N-A（本轮 fresh deploy 不升级） |

### arch-review §七 带入 implement 的 ⚠️ 加固项（test 阶段对照）

| # | 加固项 | 测试对照 | 状态 |
|---|--------|----------|------|
| 1 | A1 发卡路径重入 / CEI | ISS 系列 + _userCardCount CEI；发卡走独立 _mintFor | ✅ 落地（编译+测试绿） |
| 2 | A3 删无效 _reentrancyGuardInit | T1.5 已删（commit b55ef16），编译绿、guard 有效 | ✅ |
| 3 | Arbitrum transient 实测 | 本地 cancun 编译+测试全绿；**421614 实测待 key** | ⏸ 待 key |
| 4 | deploy.ts 参数同步 | Phase 4 部署成功（initializer 参数 + MockUSDT 步骤0） | ✅ |
| 5 | monthlySettlement 旧喂价旁路 | MS-01 用 amounts[] 不求和；旁路去留交后端(2/3) | ✅（合约侧收敛） |

> design v2 §十复审重点（MS-01~03 / ATC-01/02 / USDT-01/02 / ISS-04）**全部绿**。

---

## 六、Phase 6 — 编码规范 / 资损 grep 自检

在 `packages/contracts/contracts/` 下：

| 自检项 | 命令 | 结果 |
|--------|------|------|
| 硬编码精度字面量（`10**18`/`1e18`/`10**6`/`/100000`/`parseEther`） | rg --glob '*.sol' | **0 命中**（T2/T4 已清除；MIN_DEPOSIT=`10 * 10 ** decimals()` 动态派生） |
| payable / msg.value 残留 | rg --glob '*.sol' | **0 命中**（已全 ERC20 化） |
| SafeERC20 覆盖资金转移点 | rg Deposit.sol/Payment.sol | **全覆盖**：Deposit `safeTransferFrom`(入)/`safeTransfer`(出)；Payment 两段 `safeTransferFrom`(operator+platform)；均 `using SafeERC20 for IERC20` |

> 注：Deposit.sol `trafficCardQuota = 100*1024*1024`（100MB 流量额度）是数据量常量，**非代币精度字面量**，与精度解耦（v2-C 设计要求），不算违规。

---

## 七、已知限制（如实记录，不算失败）

- **Arbitrum Sepolia 421614 真·上链 + TSTORE 链上实测**：本环境无 `DEPLOYER_PRIVATE_KEY` / RPC，无法执行。本地 31337/hardhat（cancun，含 EIP-1153 transient storage）编译+测试+部署全绿，`hardhat.config.ts` 已配 `arbitrum_sepolia`(421614) 网络且 key 缺失时优雅降级（空账户数组，不报错）。
- 影响验收 §五D-17 / arch-review §七-3 的「421614 跑通 / TSTORE 链上实测」部分 → 记为 handoff §10 待办（配 key 后执行），**不阻塞本阶段**（本地全绿 + 配置就绪）。

---

## 八、三个自检

1. **承诺交付什么**：test 阶段全量验证（编译/测试/覆盖率/部署冒烟/验收映射/grep 自检）+ 两份产出文档（本报告 + checkpoint）。
2. **交付物真实存在吗**：本报告 `docs/pipeline/stages/test.md` + checkpoint `docs/pipeline/checkpoints/stage-test-done.md` 真实写入；测试 75 passing、覆盖率表、部署输出均来自本次实跑（非引用旧 checkpoint）。
3. **正确连接吗**：验收映射逐条指向真实用例名 / grep 结果 / 部署断言；覆盖率数字来自 solidity-coverage 实跑；deploy.ts wiring 断言全过证明合约间 wiring 正确连接。

---

## 九、是否可进 review

**可进 review。** 本地质量门全过，资损 grep 自检全清，覆盖率达标（核心合约 100% stmts）。唯一 ⏸ 项为 421614 上链（待 key），属 handoff 待办、不阻塞 review。
