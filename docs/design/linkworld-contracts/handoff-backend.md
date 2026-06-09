# LinkWorld 合约 → 后端(2/3) Handoff 清单

> 产出阶段：合约子项(1/3) implement T5
> 依据：design §7.3.1
> 状态：本轮 selector 已变更，后端 **必须重新生成 ABI**。本文档为后端 `configs/deployments.json` + 前端 `contracts.ts` 对齐的权威来源。

---

## 1. 冻结 ABI 来源

7 个业务合约 + MockUSDT 的最终 ABI 从编译产物取：

```
packages/contracts/artifacts/contracts/<Name>.sol/<Name>.json  → .abi 字段
packages/contracts/artifacts/contracts/mocks/MockUSDT.sol/MockUSDT.json
```

合约列表：`FeeManager` `UserRegistry` `ServiceManager` `TrafficCardNFT` `Payment` `Deposit` `Oracle` + `MockUSDT`。

ABI 指纹（keccak256(formatJson)）记入 `deployments/<net>.json` 的 `abiHash` 字段，后端可据此比对是否同步。地址来源同样取 `deployments/<net>.json`（本地为 `hardhat.json`，测试网为 `arbitrum_sepolia.json`）。

---

## 2. selector 清单 + 变更标注

| 函数 | 4-byte selector | 变更 | 说明 |
|------|-----------------|------|------|
| `deposit(uint256)` | `0xb6b55f25` | 🔴 **改 selector** | 原 `deposit()` payable → 改为 ERC20，去 payable，新增 `amount` 参数。调用前需先 `usdt.approve(deposit, amount)` |
| `withdraw()` | `0x3ccfd60b` | 不变（语义改） | 改 `safeTransfer` 转 USDT（原转 native），CEI 先清零 |
| `payBill(uint256)` | `0xf0975190` | 🔴 **去 payable** | 原 payable → 链上直分（amount→operator.paymentAddress，fee→platform）。调用前用户须 `usdt.approve(payment, amount+fee)` |
| `createBill(address,uint256,uint256)` | `0xceb323e8` | 🔴 **权限改 onlyOracle** | 仅 Oracle 可调（不再 onlyOwner）。amount 为后端链下算好的 USDT 金额（v2-A）。amount>0 且 operator.paymentAddress≠0 fail-fast |
| `monthlySettlement(address[],uint256[],uint256[])` | `0x01eb00ca` | 🔴 **新签名** | 新增 `amounts[]`（v2-A）。三数组等长。后端分批调用（每批 ≤N，见 §5） |
| `issueMonthlyTrafficCards(address[])` | `0x0948eaad` | 🟢 新增/启用 | onlyOracle；逐用户 `if(!校验) continue` 不 revert 整批；幂等 |
| `applyTrafficCardToBill(uint256)` | `0x0d907c34` | 🟢 新增（受限桩） | onlyOracle；本轮仅校验账单存在 + emit，**不转资金**（见 §6） |
| `setOperatorPaymentAddress(uint256,address)` | `0xadb76801` | 🟢 新增 | onlyOwner；零地址 revert。部署后给 11 内置运营商注入非零分账地址 |
| `Oracle.setPayment(address)` | `0x9ab06fcb` | 🟢 新增 | onlyOwner；B2 授权拓扑：Oracle 调 `payment.createBill` 的前置 |
| `Payment.setOracle(address)` | `0x7adbf973` | 沿用 | B2 授权：`payment.oracle = Oracle` |

---

## 3. 金额精度语义

- **所有金额字段单位 = USDT**，精度 = `usdt.decimals()`，本轮 MockUSDT 为 **6**。事件中的金额（`DepositMade`、`BillPaid`、`UsageDataSubmitted` 等）均按 6 位精度解释。
- `MIN_DEPOSIT = 10 × 10^decimals`（本轮 = 10_000_000，即 10 USDT），合约内动态读 `decimals()`，未硬编码。
- **createBill 的 amount 是后端链下按资费算好的 USDT 金额（v2-A）**：合约不做 usage 求和，不做计价。后端负责计价后传入最终金额。
- **仅支持标准 ERC20，禁 fee-on-transfer**：若 token 实扣 < amount，Deposit 记账虚高 → withdraw 资损。正式 USDT 无此特性；不做实收差额补偿。
- **分账地址要求非黑名单/非暂停**：任一 `safeTransferFrom` revert → 整 tx 原子回滚、账单仍 unpaid，用户可换地址/重试。

---

## 4. 新增 usdt 地址来源

- USDT 地址通过 **initialize 注入**（`Deposit.initialize(userRegistry, usdt)` / `Payment.initialize(feeManager, platformWallet, usdt, serviceManager)`），部署后写入 `deployments/<net>.json` 的 `usdt` + `usdtDecimals` 字段。
- 本轮为 **MockUSDT**（6 位精度，public mint，symbol "USDT"）。正式上线替换为真实 USDT 地址（同样 6 位精度，与前端一致）。
- 后端 / 前端从 `deployments/<net>.json.usdt` 读地址，从 `.usdtDecimals` 读精度，不要硬编码。

---

## 5. 批量调用约定

- `monthlySettlement` / `issueMonthlyTrafficCards` 由后端**分批调用**，每批 `users` 数量 ≤ N。
- 合约内**不强制 maxBatch**（保留调用方分批灵活度）。安全上限 N 由 GAS-01 压测量出（T6 产出），后端遵守该上限。

### 5.1 GAS-01 压测结论（T6 回填，权威）

压测环境：hardhat（cancun，optimizer runs=200，viaIR），多梯度线性外推；单批安全预算取区块 gas 上限 30,000,000 的 50% = **15,000,000 gas**（留波动 / L1 calldata 成本 / 并发余量）。per-user gas 随批量增大略降（基础 overhead 摊薄），线性无二次方膨胀。

| 入口 | per-user gas（实测） | 15M 预算下理论上限 | **后端遵守的安全上限 N** |
|------|----------------------|--------------------|--------------------------|
| `issueMonthlyTrafficCards(address[])` | ≈ 230,000（含 NFT mint + safeMint 回调） | ≈ 65 | **N ≤ 50** |
| `Oracle.monthlySettlement(address[],uint256[],uint256[])`（全链路：createBill + 发卡 + getUnpaidBills 二次循环 + applyTrafficCardToBill 桩） | ≈ 432,000 | ≈ 34 | **N ≤ 25** |

实测样本：
- `issueMonthlyTrafficCards`：N=10→2.35M、N=25→5.77M、N=50→11.51M gas（均 < 区块上限，可上链）。
- `monthlySettlement`：N=10→4.37M、N=20→8.68M、N=30→12.97M gas。

**后端落地建议**：
- 调 `Oracle.monthlySettlement` 走全链路时，每批 `users` **≤ 25**（最紧约束，含发卡+账单+桩三段）。
- 仅调 `Deposit.issueMonthlyTrafficCards` 单独发卡时，每批 **≤ 50**。
- 真机（Arbitrum Sepolia 421614）上链后建议复测一次校准（L2 calldata 计价与本地略有差异），但本地实测的 per-user 量级与梯度线性度可直接作为分批上限依据。

---

## 6. applyTrafficCardToBill 语义（v2-B 冻结）

- 本轮为**受限桩**：onlyOracle，仅校验账单存在并 emit `TrafficCardApplied`，**不抵扣资金、不转账**。
- 后端**不要据此做资金抵扣对账**。流量卡抵扣的真实资金逻辑留待后续 Round。

---

## 7. 授权拓扑（部署后终态，供后端理解调用链）

- owner：`FeeManager/UserRegistry/ServiceManager/Payment/Oracle/Deposit` = deployer(EOA)；`TrafficCardNFT.owner = Deposit(合约)`。
- 授权字段（setter 注入）：`Payment.oracle = Oracle`、`Deposit.oracle = Oracle`、`Oracle.payment = Payment`、`Oracle.deposit = Deposit`、`Payment.usdt/Deposit.usdt = MockUSDT`、`Payment.serviceManager = SM`。
- 自动结算链：后端 → `Oracle.monthlySettlement`(onlyOwner=deployer) → `Payment.createBill`(onlyOracle，放行 Oracle) + `Deposit.issueMonthlyTrafficCards`(onlyOracle) + `Payment.applyTrafficCardToBill`(onlyOracle 桩)。

---

## 8. storage layout 冻结

- 本轮 fresh deploy，无 `unsafeAllow`。各 proxy 的 storage layout 以此次部署的实现合约为后续 Round 升级基线。
- OZ upgrades 插件在 `packages/contracts/.openzeppelin/` 下记录权威 layout 清单（升级时据此校验兼容性）。
- `deployments/<net>.json.storageLayout` 字段记说明 + manifest 路径。

---

## 9. 部署网络

| 网络 | chainId | deployments 文件 | 状态 |
|------|---------|-----------------|------|
| 本地 hardhat | 31337 | `deployments/hardhat.json`（in-process）/ `localhost.json` | ✅ 已部署，断言全过 |
| Arbitrum Sepolia | 421614 | `deployments/arbitrum_sepolia.json` | ⏳ 待配 `DEPLOYER_PRIVATE_KEY` + RPC 后执行（见 §10） |

---

## 10. Arbitrum Sepolia(421614) 真·上链待执行（遗留）

本轮**未执行真·上链**，原因：环境无 `DEPLOYER_PRIVATE_KEY` / `ARBITRUM_SEPOLIA_RPC`，无资金账户。配 key 后执行：

```bash
cd packages/contracts
cp .env.example .env   # 填入 DEPLOYER_PRIVATE_KEY（需含测试网 ETH）+ 可选 ARBITRUM_SEPOLIA_RPC
npx hardhat run scripts/deploy.ts --network arbitrum_sepolia
```

**同时需做（design §6.4 / T1.5 链上侧验证）**：在 421614 实测一笔 `payBill`，确认 `ReentrancyGuardTransient`（TSTORE/TLOAD）在 Arbitrum 链上真实可用。本地 hardhat(cancun) 已验证编译+运行通过；Arbitrum 真机确认留待上链时一并完成。若链上不支持 TSTORE，需降级 `ReentrancyGuardUpgradeable`（改 `Payment.initialize` + `__ReentrancyGuard_init`）。
