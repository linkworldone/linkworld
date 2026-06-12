# web-alignment-surface.md — Link World web(3/3) 接链/对账改动面

> 扫描 2026-06-09 | 子项目 web(3/3) | 已核对真实代码
>
> 来源：PRD `docs/pipeline/stages/requirement.md`（§三③ / §五 C/D）；合约 handoff `docs/design/linkworld-contracts/handoff-backend.md`；后端 handoff `docs/design/linkworld-backend/handoff-web.md`；真实 web 源码 `packages/web/src/**`；合约源码 `packages/contracts/contracts/{Deposit,Payment,Oracle}.sol`；部署文件 `packages/contracts/deployments/hardhat.json`。
>
> 本文件只产「接链/对账改动面」delta，逐条「现状 / 需改 / 风险」。换肤 delta 见同目录 `theme-migration.md`。web 设计基线（components/utils/project-scan/color-mapping）复用 `docs/design/linkworld/`，不重写。

---

## 0. 全局结论（改动规模一句话）

web 当前是 **0G 原生币（msg.value）模型** 的旧基线，与已 merge 入 main 的合约/后端**全面失配**。本轮 web 接链改动面集中在 6 块：① 链配置（加 421614）；② **ABI 全量重生成（已变 selector）**；③ USDT ERC20 approve 两段式；④ 计价/精度（18→6 位）；⑤ 对账模型（双写自述 → 事件驱动 pending）；⑥ WalletAuth 钱包签名。其中 ②④⑤ 是 breaking，影响面最大。

---

## 1. 链配置（chains.ts / contracts.ts / wagmi.ts）

### 1.1 chains.ts — 缺 Arbitrum Sepolia 421614
- **现状**：`packages/web/src/config/chains.ts` 仅定义 `zgMainnet(16600)` / `zgTestnet(16601)` / `hardhatLocal(31337)`。无 Arbitrum Sepolia。注意 PRD §〇 提及合约旧部署在 og_testnet **16602**，而 chains.ts 里是 16601，本就不一致。
- **需改**：新增 `arbitrumSepolia`（`id: 421614`，`nativeCurrency` ETH 18 位，`rpcUrls` = Arbitrum Sepolia 公共 RPC 或 `VITE_*` 注入，`blockExplorers` = `https://sepolia.arbiscan.io`）。0G 链定义本轮可保留不用或删（PRD R2 锁定迁 Arbitrum；建议保留 hardhatLocal + arbitrumSepolia 两条即可）。
- **风险**：RPC 选择——Arbitrum Sepolia 官方 RPC 限流，event 监听/读多调用建议走 Alchemy/Infura key（与后端 `configs/deployments.json` 的 RPC 对齐，避免本地/线上不一致，PRD §三② 已点名 RPC 不一致问题）。

### 1.2 contracts.ts — 地址全失配 + 421614 缺位
- **现状**：`CONTRACTS` 只有 `31337`（7 个合约地址）+ `16601`（全 0 占位）。**31337 的地址也已过期**：现 `contracts.ts` 里 31337 的 `UserRegistry=0xe7f1725E...`，但 `deployments/hardhat.json` 里 `UserRegistry` proxy 是 `0xDc64a140...`（`0xe7f1725E...` 现在是 implementations.FeeManager）——本地地址都对不上了，必须按 `deployments/hardhat.json.proxies` 重填。`16601` 整段应删，换 `421614`。
- **需改**：
  - 31337 段：用 `deployments/hardhat.json.proxies` 的 7 个地址覆盖。
  - 新增 421614 段：**合约真·上链未执行**（见 §7 风险），`deployments/arbitrum_sepolia.json` **当前不存在**。先全 0 占位 + `getContractAddress` 抛错保护（现有逻辑已对全 0 抛错，沿用）。
  - `ContractAddresses` type 需新增 USDT 地址项（现 type 无 `USDT` 字段，见 §3）。建议改为从 `deployments/<net>.json` 同步（含 `usdt` + `usdtDecimals`），单一出口。
- **风险**：地址硬编码在 TS 里 → 与 `deployments/<net>.json` 双源易漂移（31337 已发生）。建议 implement 阶段做一个「从 deployments json 生成 contracts.ts」的脚本或构建期导入，杜绝手抄漂移。

### 1.3 wagmi.ts — 链选择
- **现状**：`const chain = isLocalChain ? hardhatLocal : zgTestnet;`，非本地一律落 0G testnet。
- **需改**：`zgTestnet` → `arbitrumSepolia`。`VITE_CHAIN_ID==="31337"` 走 hardhatLocal，否则 421614。`getDefaultConfig` 的 `chains`/`transports` 同步。
- **风险**：`main.tsx` 的 RainbowKit accentColor 同时要换金色（见 theme-migration），属换肤但同文件，implement 时一并改。

---

## 2. ABI 重生成（abis/ → hooks/contracts/*）

> **最关键**：web 的 ABI 是旧 0G 模型快照，selector 已变。必须从 `packages/contracts/artifacts/contracts/<Name>.sol/<Name>.json` 重新导出（合约 handoff §1 冻结来源），并据 `deployments/<net>.json.abiHash` 比对同步。

### 2.1 abis/Deposit.ts — `deposit` selector 变了
- **现状**：`DepositABI` 里 `deposit` 是 `{"inputs":[],"stateMutability":"payable"}`（无参 payable）；**无 `getLockExpiry`**；`initialize(address _userRegistry)` 单参。
- **真合约**（`Deposit.sol`）：`deposit(uint256 amount) external`（去 payable，加 amount，selector `0xb6b55f25`）；新增 `getLockExpiry(address) view returns(uint256)`；`initialize(address _userRegistry, address _usdt)` 双参。
- **需改**：重生成 `DepositABI`。受影响调用：`hooks/contracts/useDepositContract.ts` 的 `useContractDeposit`（现 `functionName:"deposit", value: parseEther(amountEth)` → 改 `args:[amount], 去 value`）；`useDepositBalance`（`getDepositAmount` 不变，但精度按 6 位，见 §4）；新增 `useLockExpiry` 读 `getLockExpiry`（PRD D14 倒计时）。
- **风险**：`deposit` 改 selector 后，旧 ABI 调用会 revert（函数不存在）；`value` 传 ETH 会被新合约拒（非 payable）。必须 ABI + hook 同步改，不能只改一处。

### 2.2 abis/Payment.ts — `payBill` 去 payable + 多了 `autoSettle`/缺对齐
- **现状**：`PaymentABI` 里 `payBill(uint256)` 是 `payable`；有 `autoSettle(address[],uint256[],uint256[])`、`createBill(address,uint256,uint256)`、`getUnpaidBills`/`getUserBills`（返回 Bill tuple 含 `platformFee`/`trafficCardDeduction`）。
- **真合约**（`Payment.sol`）：`payBill(uint256) external nonReentrant`（**去 payable**，selector `0xf0975190`，调用前用户须 `usdt.approve(payment, amount+fee)`）；`createBill(...) onlyOracle`（selector `0xceb323e8`，权限改）；结算入口在 `Oracle.monthlySettlement(address[],uint256[],uint256[])`（selector `0x01eb00ca`，合约 handoff §2）——web 端 Cards 页现调的发卡/结算需核对（见 §6 Cards）。
- **需改**：重生成 `PaymentABI`；`usePaymentContract.ts` 的 `useContractPayBill` 现 `args:[billId], value` → 改 `args:[billId], 去 value`（付款金额走 USDT transferFrom，不再随 tx 带原生币），并前置 approve（§3）。
- **风险**：web 现 `payBill(billId, value)` 第二参 `value` 是 msg.value，新合约不收；不改会 revert。

### 2.3 abis/index.ts — 未导出 OracleABI（PRD 验收 D11 明列）
- **现状**：`abis/index.ts` 导出 6 个（UserRegistry/Deposit/ServiceManager/Payment/FeeManager/TrafficCardNFT），**漏 OracleABI**（文件 `abis/Oracle.ts` 存在但未 re-export）。
- **需改**：补 `export { OracleABI } from "./Oracle";`，并核对 `Oracle.ts` 内容是否含 `monthlySettlement` 新签名（合约已加，web 旧 ABI 可能缺）。
- **风险**：低，但 PRD 验收 D11 明确要查，必补。

### 2.4 其余合约 ABI 同步
- FeeManager（`getFeeRate`，见 §4 手续费）、TrafficCardNFT、ServiceManager、UserRegistry：按 `abiHash` 逐个比对重生成。Deposit 新增事件 `TrafficCardIssued`、Payment `BillCreated/BillPaid/TrafficCardApplied` 等若 web 要监听，需对齐事件签名（对账走事件，§5）。

---

## 3. USDT 接入（原生币 → ERC20 两段式 approve）

- **现状**：充值 `useContractDeposit` 用 `value: parseEther(amountEth)`；支付 `useContractPayBill` 用 `value`。**全程当原生币**，无 ERC20、无 approve、无 allowance 检测。`ContractAddresses` type 无 USDT 地址，无 `MockUSDT`/`ERC20` ABI。
- **需改**：
  1. 新增 USDT 地址来源：从 `deployments/<net>.json.usdt` 读（本地 hardhat = `0x5FbDB231...`，MockUSDT 6 位）；新增 ERC20/MockUSDT ABI（至少 `approve` / `allowance` / `balanceOf` / `decimals`）。
  2. 充值前：检测 `allowance(user, Deposit)` → 不足则 `usdt.approve(Deposit, amount)` → 等确认 → `deposit(amount)`。
  3. 付账前：`usdt.approve(Payment, amount + platformFee)`（精确总额）→ `payBill(billId)`。
  4. **exact-amount approve，禁 infinite**（后端 handoff §2 资损硬约束）：每次按精确金额 approve，**不得** `approve(MaxUint256)`。
  5. UI 两步态：PRD D13「allowance 不足时先 approve，再 deposit/pay，UI 有两步态」——现 `Deposit.tsx` 只有单步 `Confirm` + `Waiting for signature/Confirming`，需扩成 `Approve → Deposit` 两阶段。
- **风险**：
  - 两笔交易（approve + deposit）非原子，approve 成功 deposit 失败会留下精确额度授权（可接受，下次复用或归零）。
  - approve 金额必须 = 链上精确所需（含手续费），算错会 transferFrom revert。
  - 精度：amount 用 `10^usdtDecimals`（6 位），不可沿用 `parseEther`（18 位），见 §4。

---

## 4. 计价 / 手续费 / 精度（constants.ts + format.ts + FeeManager）

### 4.1 constants.ts 全部失配
- **现状**：`PLATFORM_FEE_RATE = 0.025`（写死，且值 2.5% 与合约默认 1.5% 都不符）；`MIN_DEPOSIT_USDT = 100n * 10n ** 18n`（**18 位 + 100 起，双错**：应 6 位 + 10 USDT）；`SUPPORTED_CURRENCIES = ["USDT","ETH"]`。
- **需改**：
  - **删 `PLATFORM_FEE_RATE`**，改读链 `FeeManager.getFeeRate()`（基点，默认 150 = 1.5%，分母 10000）。PRD R5 / D12：手续费 UI 必须 = 链上实读值。
  - `MIN_DEPOSIT_USDT` → `10n * 10n ** BigInt(usdtDecimals)`（= 10_000_000，10 USDT 6 位），`usdtDecimals` 从 deployments 读，不硬编码。
  - `SUPPORTED_CURRENCIES` → 仅 `["USDT"]`（PRD R3，去 ETH/A0GI）。
- **风险**：`PLATFORM_FEE_RATE` 被哪些组件消费需 grep（申请号码弹层 + 账单页都要展示手续费，PRD D12），删后改读 FeeManager hook。

### 4.2 format.ts — formatAmount 默认 18 位（USDT 是 6 位）
- **现状**：`utils/format.ts` 的 `formatAmount(wei, decimals=18, ...)` 与 `parseUnits(value, decimals=18)` 默认 18 位。`AmountDisplay`、`Deposit.tsx`、`depositApi`（`parseEther`）全链路按 18 位。
- **需改**：金额展示/解析统一改 6 位（USDT），`decimals` 从 `usdtDecimals` 取。`depositApi.recordDeposit` 里 `parseEther(amount)`（18 位）→ 改 6 位 `parseUnits(amount, usdtDecimals)`。
- **风险**：**资损/显示错位**——18 位当 6 位算会差 10^12 倍。这是全链路一致性硬要求（后端 handoff §5、合约 handoff §3），漏改任一处即金额错。`billingApi.toBill` 用 `parseFloat` 加减（浮点）也需复核，后端建议用 bigint。

### 4.3 currency 字段写死 "ETH"
- **现状**：`useDeposit` 里 `currency: "ETH"`；`depositApi.getHistory` map `currency: "ETH"`；`Deposit.tsx` currency state 默认含 ETH tab。
- **需改**：全改 "USDT"。AmountDisplay 后缀、Deposit tab、历史记录展示统一 USDT。

---

## 5. 对账对齐（双写自述 → 事件驱动 pending）【breaking】

> 后端 handoff §1：对账模型从「前端自述即记账」改为「链上事件唯一回填终态」。web 现有 **双写**模式全部要降级为 pending 意向上报。

### 5.1 充值（useDeposit.ts / depositApi.ts）
- **现状**：`useDepositMutation.recordToBackend(amount)` 在合约成功后调 `depositApi.recordDeposit(wallet, amount, txHash)` POST `/api/deposit` 写记录；`Deposit.tsx` 合约成功即 `toast.success("Deposit confirmed")`。
- **需改**：`POST /api/deposit` 现仅记 pending 意向（后端不再据此置余额，余额由 `DepositMade` 事件确认）。web 不能据 200 显示「已充值/已到账」，中间态显示「处理中」，余额读链上 `getDepositAmount`（source of truth，已在做）或轮询后端 confirmed。`recordDeposit` 传 `txHash` 后端已**不据此记账**（最多关联意向），web 不应依赖。
- **风险**：用户体验——充值后余额需等事件确认才更新，UI 要有 pending 态，不能立即显示到账造成误导。

### 5.2 提现（useDeposit.ts / depositApi.ts）
- **现状**：`useWithdrawMutation.recordToBackend()` 合约成功后调 `depositApi.recordWithdraw(wallet, txHash)` POST `/api/withdraw`，**凭 txHash 记账**。
- **需改**：后端 handoff §1.2——`/api/withdraw` **不再接受 txHash 记账**，最多写 pending（不计入余额），记账唯一由 `DepositWithdrawn` 事件回填。web 提现历史/余额以事件确认为准。**废弃 `recordWithdraw` 凭 txHash 标记**。
- **风险**：提现历史展示需区分 pending/confirmed；reorg 时后端回退未确认记录，web 不缓存未确认态当真（handoff §4 状态机）。

### 5.3 付账（useBilling.ts / billingApi.ts）
- **现状**：`usePayBill.recordToBackend(billId)` 合约成功后调 `billingApi.recordPayment(wallet, billId, txHash)` POST `/api/bills/pay`，**直接触发后端置已付**。
- **需改**：后端 handoff §1.1——`/api/bills/pay` v2 **只写 pending 意向（PayIntentTxHash），不动 IsPaid**；`is_paid` 唯一由 `BillPaid` 事件回填（等 K 块）。web 支付后调本端点仅作「已发起支付」意向，**不能据 200 显示「已付」**，中间态显示「支付确认中」，真实状态轮询 `GET /api/bills/:wallet` 看 `is_paid`。**废弃 `recordPayment` 凭前端标记已付**。
- **风险**：`billingApi.toBill` 按 `is_paid` 算 status，需保证读的是事件确认后的值；中间态「确认中」需新增 status 分支（现只有 unpaid/paid/overdue）。

### 5.4 状态机统一（handoff §4）
- pending → seen → confirmed（等 K 块）。web 余额/历史只认 confirmed；pending/seen 显示「处理中」不计入可用余额；reorg 回退不缓存。建议抽一个 pending 态展示规范，三处（充值/提现/付账）复用。

---

## 6. WalletAuth 钱包签名（新增，后端 handoff §3）

- **现状**：`services/api/client.ts` 是裸 axios，无任何签名头；写端点（deposit/withdraw/bills-pay/service 写）直接 POST。
- **需改**：受保护写端点需带钱包签名头——web 用用户钱包对「请求体 + nonce/timestamp」签名（EIP-712 / EIP-191，具体格式 implement 阶段双方敲定，handoff 锁定「需要钱包签名」契约）。后端 `ecrecover` 还原地址须 == 请求 `wallet` 字段否则 401，带 nonce 或时间窗防重放。建议在 `apiClient` 加请求拦截器，对写端点用 wagmi `signTypedData`/`signMessage` 注入签名头。
- **风险**：
  - 每次写操作多一次钱包签名弹窗，UX 要顺；nonce 需后端下发或本地时间窗，防重放语义要和后端 `signatures.go` 对齐。
  - 读端点（`GET /api/bills/:wallet` 等）handoff §6 明确**不变、无需鉴权**，不要误加签名。

---

## 7. UI 功能对齐 PRD（Deposit 倒计时 / Cards 双 Tab / 手续费展示）

### 7.1 Deposit 锁仓倒计时（PRD R6 / D14）
- **现状**：`Deposit.tsx` Withdraw 按钮无条件可点；无 `getLockExpiry` 读取、无倒计时、无禁用。
- **需改**：读 `Deposit.getLockExpiry(addr)`（合约已有），锁仓未到 → Withdraw 禁用 + 显示剩余倒计时；到期 → 可提取。需新增 `useLockExpiry` hook（§2.1）。
- **风险**：倒计时需轮询/定时刷新；到期边界（==now）的判定与合约 `require` 对齐。

### 7.2 Cards 双 Tab + 移除 Admin 发卡（PRD R7 / D15 / R9 / D16）
- **现状**：`Cards.tsx` 单视图；有 **「Issue Monthly Card (Admin)」按钮**直接调 `issue.issue([address])`（`useIssueMonthlyCards`）——PRD R7 明确**移除**该按钮，改「锁仓满自动发放」说明。无 SIM 领取 Tab。
- **需改**：
  - 移除 Admin 发卡按钮 → 替换为「锁仓满自动发放」说明 + NFT 自动发放状态（来自链上真实数据，发卡由后端 `Oracle.monthlySettlement`/`Deposit.issueMonthlyTrafficCards` 触发，web 不主动发）。
  - 新增**双 Tab：流量卡 NFT / SIM 领取**。SIM 领取为前端降级占位（R9/D16）：选国家 + 收件表单 → 提交 toast 成功 + 写 `localStorage pendingSync`（项目已有 `utils/pendingSync`）+ 「全球通即将推出」。
  - 不可转卖 / 销毁后 30 天规则文案到位。
- **风险**：现 Cards 的 `useIssueMonthlyCards` 是 onlyOracle 受限调用，普通用户调会 revert（合约权限已收紧）——移除按钮正好规避；保留会暴露 revert 报错。`useTrafficCards` 读链路需核对新 TrafficCardNFT ABI。

### 7.3 手续费展示（PRD D12）
- **现状**：手续费靠写死 `PLATFORM_FEE_RATE=0.025`（§4.1），且未见在「申请号码弹层 + 账单页」按 1.5% 展示链上值。`billingApi` 的 `platform_fee` 来自后端但 UI 是否展示需核对各页。
- **需改**：申请号码弹层 + 账单页/BillDetail 均展示手续费 = 链上 `FeeManager.getFeeRate()` 实读（基点/10000 = 1.5%），不是写死 0.015/0.025。
- **风险**：FeeManager.getFeeRate 是基点（150），展示需 /10000 转百分比，别和小数 0.015 混。

---

## 8. 改动面汇总表（影响文件清单）

| 块 | 改动文件 | breaking? |
|----|----------|-----------|
| 链配置 | `config/chains.ts`、`config/contracts.ts`、`config/wagmi.ts` | 是（地址/链全换） |
| ABI 重生成 | `config/abis/{Deposit,Payment,Oracle,FeeManager,...}.ts`、`config/abis/index.ts`（补 OracleABI） | 是（selector 变） |
| 合约 hook | `hooks/contracts/{useDepositContract,usePaymentContract,useTrafficCard,...}.ts`（去 value、加 amount、加 approve、加 getLockExpiry） | 是 |
| USDT approve | 新增 ERC20/USDT ABI + allowance/approve hook；`Deposit.tsx` 两步态 | 是（新流程） |
| 计价/精度 | `config/constants.ts`（删 FEE_RATE、改 MIN/CURRENCIES）、`utils/format.ts`（18→6 位） | 是（资损敏感） |
| 对账 | `hooks/useDeposit.ts`、`hooks/useBilling.ts`、`services/api/{depositApi,billingApi}.ts`、`Deposit.tsx`/`Billing.tsx` pending 态 | 是（语义反转） |
| WalletAuth | `services/api/client.ts`（签名拦截器） | 是（新增鉴权） |
| UI 功能 | `pages/Deposit.tsx`（倒计时）、`pages/Cards.tsx`（双 Tab + 去 Admin）、申请号码弹层/`pages/BillDetail.tsx`（手续费） | 否（功能补齐） |
| 换肤 | 见 `theme-migration.md` | 否（视觉） |

---

## 9. 遗留 / 阻塞（implement 必读）

1. **合约真·上链未执行**（合约 handoff §10）：`deployments/arbitrum_sepolia.json` **不存在**，421614 上无真实地址。web 421614 段只能**全 0 占位**，端到端 Arbitrum 验收（PRD D17）**阻塞**于合约先上链。本地 31337 可全链路验证。
2. **31337 地址已漂移**：`contracts.ts` 的 31337 地址与 `deployments/hardhat.json.proxies` 不一致（实测对不上），即便本地也跑不通，必须按 deployments 重填。建议建「deployments json → contracts.ts」单一出口避免再漂。
3. **WalletAuth 签名格式未定**：handoff §3 只锁「需要钱包签名」，EIP-712 vs EIP-191 / 字段顺序 / nonce 来源由 implement 双方敲定，需与后端 `signatures.go` 同步。
4. **对账 breaking 影响 UX**：充值/提现/付账三处都从「立即显示成功」改「处理中等事件确认」，需统一 pending 态设计，否则用户误判到账/已付。
5. **精度 18→6 是资损红线**：format.ts/constants.ts/depositApi 任一漏改即金额错 10^12 倍，implement 时全链路 grep `parseEther`/`decimals=18`/`10n ** 18n` 清零。
6. **og_testnet(16602) 旧链路**：PRD R2 已锁定迁 Arbitrum，0G 链定义（chains.ts 16600/16601）本轮可清理，避免误连。
