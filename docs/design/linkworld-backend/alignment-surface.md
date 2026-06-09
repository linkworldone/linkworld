# 后端子项目 alignment-surface 基线（最关键）

> 扫描 2026-06-08 | 子项目 backend(2/3) | 已核对真实代码
>
> 对照 `docs/design/linkworld-contracts/handoff-backend.md`，逐条列后端对齐合约的改动面，每条标「现状 / 需改 / 风险」。

## 0. 全局结论

合约本轮做了 ERC20(USDT) 迁移 + selector 变更 + 计价归后端 + onlyOracle 权限收口。但后端**当前链集成几乎全为 stub**（见 blockchain-integration.md）：不发任何链上写交易、不真正订阅事件、ABI 不全、计价是随机数模拟器、deployments.json 还指向 0G 旧网。

→ 因此「对齐」**不是改几处调用签名，而是从零补齐链上写入 + 事件同步 + 真实计价**。规模评估：**大改/接近重写链侧**，业务 DB 层基本可复用。

---

## 1. 链配置（deployments.json）

**现状**
- `chainId: 16602`、`rpcUrl: https://evm-testnet.0g.ai`、键名 `contracts` 下 7 个 **0G 旧部署地址**。
- Go struct 读 `proxies` 键 → 实际加载为空 map（键名 bug）。
- **无 `usdt` / `usdtDecimals` / `abiHash` 字段**。
- `.env.example` 同步残留 `RPC_URL=https://evm-testnet.0g.ai`、`CHAIN_ID=16602`。

**需改**
- `chainId` → **421614**；`rpcUrl` → Arbitrum Sepolia RPC。
- 7 合约地址 → 取自合约子项 `deployments/arbitrum_sepolia.json.proxies`（待该文件生成；本地联调先用 `hardhat.json` chainId 31337 的 proxies）。
- 新增 `usdt`（MockUSDT 地址）+ `usdtDecimals: 6` + `abiHash`（7 合约指纹）。
- 修 `config.go` 键名：JSON 用 `proxies` 或 struct tag 改回 `contracts`（二选一，建议统一为合约产物的 `proxies`），并补 `Usdt`/`UsdtDecimals`/`AbiHash` 字段。
- 同步 `.env.example` 的 RPC/CHAIN_ID。

**风险**
- 当前 mock USDT 地址：合约 `hardhat.json` 给的是 `0x5FbDB2315678afecb367f032d93F642f64180aa3`（31337 本地）；**Arbitrum Sepolia 上的 MockUSDT 地址尚不存在**（真·上链未执行，handoff §10）。后端不能硬编码，必须等 `arbitrum_sepolia.json` 出来再填，否则联调全失败。
- 键名 bug 不修则 event sync 永远拿空地址表，静默失效。

---

## 2. ABI 重生成 + selector 变更

**现状**
- `abis/` 只有 2 份手写裁剪 ABI（`Deposit.json` 仅事件、`UserRegistry.json` 部分）。缺 `FeeManager/ServiceManager/TrafficCardNFT/Payment/Oracle/MockUSDT`。
- 后端**没有任何地方实际调用** `deposit()` / `payBill()` / `createBill` / `monthlySettlement`（grep 全空）——所以 selector 变更目前「无调用方受影响」，但反过来说**这些链上写功能根本没实现**。
- `signatures.go` 的 `BillCreated` topic 是 5 参数旧签名，与冻结 ABI 的 4 参数不一致（topic hash 错）。

**需改**
- 从 `packages/contracts/artifacts/contracts/<Name>.sol/<Name>.json` 取**全量 `.abi`**，落地 7 业务合约 + MockUSDT 共 8 份，按 `abiHash` 比对。
- 用 abigen 或 `bind.NewBoundContract` 生成绑定，才能编码以下新签名调用：
  - `deposit(uint256)` selector `0xb6b55f25`（原 `deposit()`，调用前需 `usdt.approve(deposit, amount)`）。
  - `payBill(uint256)` `0xf0975190`（去 payable，调用前 `usdt.approve(payment, amount+fee)`）。
  - `createBill(address,uint256,uint256)` `0xceb323e8`（**onlyOracle**，amount 后端算好传入）。
  - `monthlySettlement(address[],uint256[],uint256[])` `0x01eb00ca`（新增 `amounts[]`，三数组等长）。
  - `issueMonthlyTrafficCards(address[])` `0x0948eaad`、`applyTrafficCardToBill(uint256)` `0x0d907c34`、`setOperatorPaymentAddress`、`Oracle.setPayment` 等。
- 修 `signatures.go`：`BillCreated` 改 4 参数 `BillCreated(uint256,address,uint256,uint256)`；新增 `TrafficCardApplied(uint256)`。

**风险**
- ABI 不全/绑定缺失 → 链上写功能无从实现，这是阻塞链路的根因。
- selector 写错或 ABI 与链上实现不符 → 交易 revert 难排查。
- approve 前置漏掉 → ERC20 transferFrom 失败，整 tx 回滚。

---

## 3. USDT 6 位精度

**现状**
- 金额全程用 **string 原样存取**（`Deposit.Amount` / `Bill.Amount` / `Bill.PlatformFee` 都是 `string`）。
- seed 运营商 `RequiredDeposit` 是 `"0.01"` 这类小数；`OperatorAPISimulator.GetBill` 返回随机整数字符串。
- event_sync **不解析任何金额事件**，无落库，无单位换算。
- 没有任何地方读 `usdtDecimals` 或做 10^6 缩放。

**需改**
- 统一金额精度语义为 **USDT 6 位**（最小单位 = amount × 10^6）。建议链上金额用 `*big.Int`（6 位最小单位），DB 存最小单位字符串/整数，展示层再除 10^6。
- event_sync 落库 `DepositMade`/`BillPaid`/`BillCreated` 金额时按 6 位精度解释。
- 校验 `MIN_DEPOSIT = 10 × 10^6`（10 USDT），与合约动态 `decimals()` 对齐（从 `usdtDecimals` 读，不硬编码）。
- 重新定义 seed `RequiredDeposit` 的单位（当前 "0.01" 语义不明）。

**风险**
- 精度混用（有的地方当 USDT 整数、有的当 0.01）→ 押金/账单金额错误，直接资损或显示错乱。
- fee-on-transfer 禁用（handoff §3）：后端记账须信任 amount = 实收，正式 USDT 无此问题，但 mock 阶段要确认。

---

## 4. 计价逻辑（最关键缺口）

**现状**
- 后端**有「计价」但是假的**：`OperatorAPISimulator.GetBill` 返回 `rand.Intn(5000)+500` 随机金额 + 2.5% fee。
- `OracleServiceV2.FetchAndCreateBills` 拿到随机金额后**只写 DB Bill，不上链**。
- **完全没有 `monthlySettlement` 调用**，自然也没有 `amounts[]` 数组构造逻辑。

**需改**（handoff §3：createBill 的 amount 由后端链下按资费算好传入，合约不计价）
- 实现**真实计价**：按 usage（data/call）× 运营商资费表算出每用户 USDT 金额（替换随机数模拟器，或保留模拟器但产出确定金额）。
- 实现 `Oracle.monthlySettlement(users[], operatorIds[], amounts[])` 调用：构造三个等长数组，amounts 为算好的 USDT 6 位金额，**分批**调用（见 §5）。
- 后端需持有 **oracle 调用账户私钥**（onlyOracle，见 §7）来签发交易。
- 当前 `SignData` 是 SHA256 字符串，**不是链上可用签名**，需替换为真实交易签名流程（go-ethereum `bind.NewKeyedTransactorWithChainID`）。

**风险**
- 计价错误 → 账单金额错，链上资金错配。
- 没有真实计价表（资费规则未在后端定义）→ 这是产品/设计层缺口，scan 阶段需上抛澄清「资费规则从哪来」。

---

## 5. 批量上限（GAS-01 压测结论）

**现状**
- `FetchAndCreateBills` 单循环遍历**全部用户**写 DB，**无分批**。
- 无任何 `monthlySettlement` / `issueMonthlyTrafficCards` 链上批量调用。

**需改**（handoff §5.1 权威上限）
- `Oracle.monthlySettlement`（全链路）每批 `users` **≤ 25**。
- 仅 `Deposit.issueMonthlyTrafficCards` 单独发卡每批 **≤ 50**。
- 后端实现按上限切片 + 逐批发交易 + 逐批确认/重试。
- 上链后建议在 421614 真机复测校准（L2 calldata 计价差异）。

**风险**
- 不分批或批过大 → 单 tx 超 15M gas 预算甚至超区块上限，交易失败。
- 分批中途失败需幂等：`issueMonthlyTrafficCards` 合约侧幂等 + 逐用户 continue 不整批 revert；后端要据此做断点续跑，避免重复结算。

---

## 6. RPC 不一致（scan 1/3 已发现）

**现状**
- 后端 `deployments.json` + `.env.example` 用 `https://evm-testnet.0g.ai`（chainId 16602）。
- scan 阶段 1/3 发现 hardhat 配置侧用 `evmrpc-testnet.0g.ai`（host 前缀不同）。
- 两端 host 不一致，且**都指向 0G 而非 Arbitrum**。

**需改**
- 统一三端（hardhat / 后端 / 前端）到 **Arbitrum Sepolia 421614** 的同一 RPC endpoint。
- 后端 RPC 走 env (`RPC_URL`) + `deployments.json.rpcUrl` 双来源，需确认优先级一致（main.go 当前用 `os.Getenv("RPC_URL")` 起 sync，但合约地址走 deployments.json）。

**风险**
- RPC 指向错网 → 读到错链数据 / 交易发错链。
- env 与 json 两处 RPC 不一致 → 读写分裂（连 A 链发 B 链）。

---

## 7. oracle 调用权限

**现状**
- 后端**不持有任何私钥**，不作为 onlyOracle 角色调任何合约。
- `oracle/monthly-bill`、`usage/submit` 端点**无任何鉴权**，任何人可触发（仅写 DB，暂无链上后果，但接链后会成为高危）。

**需改**（handoff §7 授权拓扑：后端 → `Oracle.monthlySettlement`(onlyOwner=deployer) → `Payment.createBill`(onlyOracle 放行 Oracle) + `Deposit.issueMonthlyTrafficCards`(onlyOracle) + `Payment.applyTrafficCardToBill`(onlyOracle 桩)）
- 注意权限拆分：
  - `Oracle.monthlySettlement` 是 **onlyOwner=deployer**（后端调它要用 deployer/owner 账户）。
  - `Oracle` 合约本身被授权为 `Payment.oracle` / `Deposit.oracle`，由 Oracle 合约内部去调 onlyOracle 函数。
  - 即后端**对外只需一个能调 `Oracle.monthlySettlement` 的 owner 私钥**，链上权限传导由合约 setter 拓扑保证（`Oracle.setPayment` / `Payment.setOracle` / `Deposit.setOracle` 已配）。
- 后端需安全管理该私钥（env / KMS），并给 `oracle/monthly-bill` 等敏感端点加鉴权。

**风险**
- 私钥泄露 = 可任意发起月度结算，高危资损面。
- 用错账户（非 owner）调 `monthlySettlement` → revert。
- 敏感端点裸奔（无鉴权）一旦接链立即成攻击面。

---

## 8. 其他对齐项

- **`applyTrafficCardToBill` 是受限桩**（handoff §6）：仅 emit `TrafficCardApplied(uint256)`，不转资金。后端**不要据此做资金抵扣对账**；`Bill.TrafficCardDeduction` 字段本轮应保持 0 / 不参与真实抵扣。
- **event_sync 受影响点汇总**：①`BillCreated` topic 改 4 参数；②新增 `TrafficCardApplied` 订阅；③`DepositMade/BillPaid/...` 需真正 `FilterLogs` 落库（当前主循环空转）；④金额按 6 位解释；⑤合约地址表先修键名 bug 才有地址可订阅。
- **`processUserRegistered`** 落库 `RegisteredAt=Unix(0,0)`、不解析 email/tokenId，需从事件 data 正确解码。

---

## 9. 改动规模速评

| 模块 | 改动量 | 说明 |
|------|--------|------|
| `configs/deployments.json` + `config.go` | 中 | 换网/换址/补字段 + 修键名 bug |
| `abis/` | 大 | 补 6 份全量 ABI + abigen 绑定（当前 0 绑定） |
| `blockchain/client.go` | 大 | stub → 真实链上读 + 写交易（approve/createBill/monthlySettlement/issue...） |
| `signatures.go` | 小 | 修 BillCreated + 加 TrafficCardApplied |
| `sync/event_sync.go` | 大 | 空转 → 真实 FilterLogs 多事件落库 + 6 位精度 |
| `services/oracle.go` 计价 | 大 | 随机模拟器 → 真实资费计价 + 分批 monthlySettlement + 真实签名上链 |
| `services.go` 押金/账单 | 中 | 接链上写 + 事件回填闭环（当前纯 DB） |
| 鉴权/私钥管理 | 中 | 新增 oracle 账户 + 敏感端点鉴权 |
| 业务 DB / handlers / models | 小 | 基本复用，金额精度语义微调 |
