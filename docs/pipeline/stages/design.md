# Stage: design v2 — web 重构接链 设计分析（子项目 web 3/3）

> **版本**: v2（按 arch-review 返工，闭合 B1-B6 + 关键 ⚠️） | **状态**: reworked，待 re-review | **日期**: 2026-06-10 | **Gate**: 1（设计） | **子项目**: web(3/3) | **角色**: 产品设计师/架构师 | **分支**: web/deep-blue-gold-refactor
> **v2 返工范围**（见文末「arch-review 阻塞闭合对照表」）：B1 WalletAuth 会话级签名铁律 / B2 金色对比度（卡内金额改 navy）/ B3 tsc 绿基线前置 / B4 测试策略 / B5 TwoStepAction 状态机 / B6 confirmed 信号来源 + getLogs 限流；关键 ⚠️ 一并处理。新增节：「测试策略」「TwoStepAction 状态机」「confirmed 信号来源」「arch-review 阻塞闭合对照表」。
> **状态(原 v1)**: completed | **日期**: 2026-06-09 | **Gate**: 1（设计） | **子项目**: web(3/3) | **角色**: 产品设计师 | **分支**: web/deep-blue-gold-refactor
> **输入**：PRD `requirement.md` §三③/§五C·D·E/§七 + 两份 web delta（`docs/design/linkworld-web/{web-alignment-surface,theme-migration}.md`）+ 后端 `handoff-web.md`（对账契约+状态机）+ 合约 `handoff-backend.md`（冻结 ABI/USDT 6 位）+ 搁置深蓝金 `docs/design/linkworld/DESIGN.md`（b18cf37）+ 真实源码逐文件核对。
> **用户已拍板**（2026-06-09）：复用深蓝金视觉系统（DESIGN.md，PRD R12 锁定视觉不变）+ **全力产接链交互稿**；跳过 shotgun 视觉变体探索（且 gstack `$D` 设计器不可用）。
> **注**：本 stage 文件按子项目串行复用——合约(1/3) design 见 git `2b8d4a0`，后端(2/3) 见 `9798789`，本版为 web(3/3)。
> ⚠️ 本文档只产设计分析，不写业务代码。

---

## 0. 本轮 design 形态（与合约/后端子项目的差异）

合约/后端子项目的 design 是「从零设计技术方案」；web 子项目不同：**视觉系统已在上一轮定稿**（DESIGN.md），PRD R12 锁定不变。所以本轮 design 的真正缺口**不是视觉**，而是 DESIGN.md 未覆盖的 **PRD §七 item 8/9 接链交互稿**。本轮 design 工作 = ① 复用 + 轻修订 DESIGN.md；② 新增「接链交互模式」章节（已写入 DESIGN.md）；③ 本文件汇总分析 + 移交。

---

## 1. 已核对真实代码清单（头部强制）

| 核对项 | 结果 | 出处 |
|--------|------|------|
| `index.css :root` 深紫系 oklch（~hue 270）+ body `font-family:"Inter"` + Geist 已自托管 / **Space Grotesk 未安装** | ✓ DESIGN.md 假设全成立 | `packages/web/src/index.css` |
| `tailwind.config.ts` 硬编码 HEX token（surface/brand/status）+ `fontFamily.orbitron` | ✓ theme-migration §1.1 准确 | `packages/web/tailwind.config.ts` |
| dist build 产物残留 `#3b82f6`×3 / `#8b5cf6`×1 / `Inter` / `Orbitron` | ✓ 旧色值确编进产物 | `packages/web/dist/assets/*.css` |
| `Deposit.tsx` 单步 `Confirm`（无 approve 两步态）+ 含 `"ETH"` currency tab + `border-brand-blue` | ✓ | `packages/web/src/pages/Deposit.tsx` |
| `FeeManager.getFeeRate` ABI 存在但**无 hook 读取** | ✓ 需新增 useFeeRate | `src/config/abis/FeeManager.ts:206` |
| 合约 `deposit(uint256 amount)` 非 payable / `withdraw()` 无参 / `getLockExpiry(addr)` 存在 | ✓ | `contracts/Deposit.sol:54,71,125` |
| 锁仓边界：`require(block.timestamp >= _lockExpiry)`（**`>=`**，到期即可提） | ✓ | `Deposit.sol:72` |
| 锁仓累加：首存 `now+30d`，**之后每次充值 `+= 30 days` 顺延**，提现归 0 | ✓ 交互稿据此提示 | `Deposit.sol:61-64,77` |
| `payBill(uint256 billId)` 非 payable nonReentrant / `createBill onlyOracle` | ✓ | `contracts/Payment.sol:102,74` |
| `FeeManager.calculateFee(amount)` 可直读精确费额 | ✓ approve 总额=amount+calculateFee | `contracts/FeeManager.sol:34` |
| `utils/pendingSync.ts` 存在（SIM 领取占位复用） | ✓ | `packages/web/src/utils/pendingSync.ts` |

---

## 2. 复用与修订：深蓝金视觉系统

- **复用 DESIGN.md 全部规范**：色值（navy #0C2340 / 沉稳金 #D4AF37 / 暖米白 #F7F3EA）、字体（Display=Space Grotesk，Body/Data=Geist，删 Inter、弃 Orbitron）、lucide 映射表、**CSS 变量单一出口方案**（index.css `:root` 真源，tailwind token 引 `var()`）、缺失原子组件（ui/card / ui/input）。
- **本轮修订**：① Product Context `0G Chain 生态` → `Arbitrum Sepolia(421614) + ERC20 USDT 6 位`；② 香槟金 `#F0C75E`（PRD E18 提及）**明确做**：金线渐变高光端 `#D4AF37→#F0C75E`，停靠点 0%/50%/100%；**不作独立文字/填充色**（与 DESIGN.md 一致）。**主金仍 #D4AF37，不引入新主色**（R12 视觉不变）。
- **换肤 delta** 见 `theme-migration.md`（旧色值/emoji 精确分布 + grep 清零基线），本文件不重复。

---

## 3. 接链交互稿（核心产出）

> 规范级 token / 组件 / 图标 / 文案口径见 DESIGN.md「接链交互模式」章节。本节给**状态机 + 流程 + 契约依据**。

### 3.1 通用交易三态（最重要）
对账模型反转（handoff §1/§4）：链上事件唯一回填终态，HTTP 端点只写 pending 意向。统一三态 `pending → confirmed → failed`：
```
[无记录] ──HTTP 意向──▶ [pending] ──event_sync 监听──▶ [seen] ──等 K 块──▶ [confirmed]
                       (不计入余额)                                          (计入余额)
                                         └────────── reorg 回退 ──────────┘ (不缓存)
```
**铁律**：可用余额只算 confirmed；pending 单列「处理中 +N USDT」弱化不染绿；绝不据 tx 成功 / HTTP 200 显示终态成功。**confirmed 信号来源见下方专节（B6）**——前端不监听链事件做终态，余额读链 + 历史/账单轮询后端 status。pending 超时（>~2min）兜底：Arbiscan 逃生链接 + 「可安全离开，到账后通知」+「约 1-2 分钟」（不暴露 K 块）。

### 3.2 USDT approve 两段式（充值 + 付账）
- 合约依据：`deposit(uint256)` / `payBill(uint256)` 均非 payable，前置 `usdt.approve(spender, exact)`。**exact-amount 禁 infinite**（handoff §2 资损硬约束）。
- approve 额：充值 `amount`；付账 `amount + calculateFee(amount)`（直读合约，不自算）。
- 充值状态机：
  ```
  Idle(输入金额)
   → 校验[amount ≥ 10 USDT(链上 require) & amount ≤ 钱包 USDT 余额]
   → [allowance(user,Deposit) < amount] Approve: 签名→Approving…→Confirming approval…
   → Deposit: 签名→Depositing…→Confirming on chain…
   → 合约成功 → POST /api/deposit(pending 意向)
   → UI=pending「处理中」(余额读链 getDepositAmount / 等 DepositMade 事件 confirmed)
  ```
  allowance 已 ≥ amount → 跳过 Approve 直达 Deposit。
- 付账状态机：同构，approve(Payment, amount+fee) → payBill(billId) → POST /api/bills/pay(pending) → UI「支付确认中」，is_paid 等 BillPaid 事件。
- 风险：approve 与 deposit/pay 两笔非原子，approve 成功 deposit 失败会留精确额度授权（可接受，下次复用或归零）。

### 3.3 对账 pending 态（三处语义反转）
| 处 | 现状（双写自述） | 改为 |
|----|------------------|------|
| 充值 `useDeposit`/`depositApi` | 合约成功即 toast「Deposit confirmed」+ recordDeposit(txHash) | pending「处理中」，余额等事件 confirmed |
| 提现 `useDeposit`/`depositApi` | 凭 txHash recordWithdraw 记账 | **废弃 txHash 记账**，pending「提现确认中」，等 DepositWithdrawn 事件 |
| 付账 `useBilling`/`billingApi` | recordPayment(txHash) 直接置已付 | **不据 200 置已付**，Bill status 扩展 `unpaid/paying/paid/overdue`，is_paid 等 BillPaid 事件 |
- Bill status 新增 `paying`（确认中）分支：info 蓝 + `Loader2`，禁绿。
- **`billingApi.toBill` 改 bigint（资损升，必改 ⚠️）**：现 `toBill` 用 `parseFloat(operatorFee)+parseFloat(platformFee)-parseFloat(trafficCardDeduction)).toFixed(2)`——对 6 位最小单位字符串做 JS 浮点加减，既单位语义错（最小单位当成元）又有大额超 `MAX_SAFE_INTEGER` 风险。改为 **bigint 加减**（`BigInt(operatorFee)+BigInt(platformFee)-BigInt(trafficCardDeduction)`），展示再经 `formatAmount(total, usdtDecimals)`；`totalAmount` 全程不落 number。

### 3.4 锁仓倒计时（Deposit 提现区）
- 合约依据：`getLockExpiry(addr)` 时间戳；提现 `require(now >= expiry)`（边界 `>=`）；**每次充值 expiry += 30d 顺延**；提现归 0。
- 状态：`expiry==0` 无锁仓 / `now<expiry` Withdraw 禁用 + 倒计时「锁仓中 · 剩余 Dd Hh Mm」+ Lock / `now≥expiry` 可点 + Unlock +「锁仓已满，可提取本金」。
- **关键 UX 提示**：因累加语义，Deposit 页须明示「再次充值将把锁仓期顺延 30 天」。
- 刷新：每分钟（剩余 ≤1 天每秒）；边界与合约 `>=` 对齐。

### 3.5 Cards 双 Tab + 移除 Admin 发卡
- `ui/tabs`（variant=line）：Tab1 流量卡 NFT / Tab2 SIM 领取。
- Tab1：**移除 Admin 发卡按钮**（onlyOracle，用户调 revert）→「锁仓满自动发放」说明卡（Sparkles/Info）；NFT 列表读链真实数据 +「不可转卖」+「销毁后 30 天」文案；空态 EmptyState(CreditCard)。
- Tab2 SIM 领取（降级 R9）：选国家 + 收件表单 → toast 成功 + 写 `pendingSync` +「全球通即将推出」；图标 Nfc/MailPlus。

### 3.6 手续费读链展示（D12）
- 删写死 `PLATFORM_FEE_RATE`；`useFeeRate` 读 `getFeeRate()`（基点 150=1.5%，/10000）；精确费额 `calculateFee(amount)`。
- 展示位：申请号码弹层（RegionDetail）+ 账单页/BillDetail。明细行「平台手续费 (1.5%)：N USDT」。
- loading/读链失败：skeleton 或「--」，不写死兜底。

### 3.7 WalletAuth 签名 UX（新增鉴权）
- 写端点（deposit/withdraw/bills-pay/service 写）带钱包签名头；读端点不加（handoff §6）。
- UX：写操作前多一次签名弹窗，**与交易签名视觉区分**——「请在钱包中签名以验证身份（不消耗 gas）」+ ShieldCheck/PenLine；拒签 toast「身份签名被取消，操作未提交」，不进 pending。
- **会话级签名一次铁律（B1，不得降级到「每次写操作签」）**：一个钱包会话只签一次身份签名（带 `nonce + 时间窗`），内存缓存复用，禁每次 deposit/withdraw/pay 弹签名；过期或换钱包重签。签名方案锁定 **EIP-712**；nonce 来源（后端下发 vs 前端时间窗）跨端待与后端 `signatures.go` 对齐。详见 DESIGN.md「接链交互模式 §6」+ Decisions Log。
- **实现约束**：用 `signedPost(path, body)` helper 封装（内部 `await signTypedData` 取/签缓存 → 带头调 axios），**不在 axios 全局拦截器调 React hook**（拦截器非 React 上下文，调 `useSignTypedData` 必崩）。

---

## 4. 组件复用 / 新增映射

| 类型 | 组件 | 改动 |
|------|------|------|
| 复用重着色 | `ui/button` `ui/badge` `ui/tabs` | 吃 CSS 变量，改 `:root` 即生效；Cards 双 Tab 用 tabs |
| 新增原子（DESIGN.md 已定） | `ui/card`（暖米白+金线）`ui/input` | 收口手写卡/表单 |
| 接链专用新增（本轮建议） | `TwoStepAction` | approve→action 两步态 + stepper + allowance 跳步（充值/付账复用） |
| | `TxStatusBadge` | 通用交易三态徽章 |
| | `LockCountdown` | 读 getLockExpiry → 倒计时/解锁态 |
| | `FeeBreakdown` | 费用明细（读链费率） |
| 改 props | `AmountDisplay` | 默认 colorClass `text-status-warning`→**按底色分流**（卡内 navy `text-text-primary` / 深底金 `text-on-dark-gold`，B2 覆盖旧「改吃金」）；币种 ETH→USDT；6 位精度 |
| | `EmptyState`/`GuardCard` | icon prop `string`→`LucideIcon`（DESIGN.md D5） |
| 业务页改造 | `Deposit.tsx`（两步态+倒计时+pending）`Cards.tsx`（双Tab+去Admin）`Billing/BillDetail.tsx`（paying态+手续费）`RegionDetail.tsx`（申请弹层手续费） | 见 §3 |

---

## 5. 交互稿依赖的真实契约映射

| UI 动作 | 合约 / API | 鉴权 / 事件 |
|---------|-----------|------------|
| 充值授权 | `usdt.approve(Deposit, amount)` | — |
| 充值 | `Deposit.deposit(uint256 amount)` | 余额经 `DepositMade` 事件 confirmed |
| 充值意向 | POST `/api/deposit` | WalletAuth；仅 pending |
| 提现 | `Deposit.withdraw()` | 经 `DepositWithdrawn` 事件 |
| 提现意向 | POST `/api/withdraw` | WalletAuth；仅 pending（废弃 txHash 记账） |
| 锁仓查询 | `Deposit.getLockExpiry(addr)` | 读，无鉴权 |
| 付账授权 | `usdt.approve(Payment, amount+calculateFee(amount))` | — |
| 付账 | `Payment.payBill(uint256 billId)` | 经 `BillPaid` 事件置 is_paid |
| 付账意向 | POST `/api/bills/pay` | WalletAuth；仅 pending |
| 手续费率/费额 | `FeeManager.getFeeRate()` / `calculateFee(amount)` | 读，无鉴权 |
| 流量卡列表 | `TrafficCardNFT`（链上读） | 自动发放由 `Oracle.monthlySettlement`/`Deposit.issueMonthlyTrafficCards` 触发 |
| 账单/押金/用量读 | GET `/api/bills|deposit|usage/:wallet` | **不变、无需鉴权**（handoff §6） |

---

## 6. 验收对照（design 负责锁定项，PRD §五 C/D/E）

| # | 验收点 | design 锁定状态 |
|---|--------|----------------|
| C10/E18 | 金值 #D4AF37 / 字体 / 渐变停靠点 / 米白色阶 / 色值单一出口 | ✓ DESIGN.md 已锁定；**B2 金色对比度修正**（卡内金额改 navy）；香槟金 #F0C75E **明确做**=金线渐变高光端(停靠点 0/50/100%) |
| D12 | 手续费 UI = 链上 getFeeRate 实读 | ✓ §3.6（useFeeRate + calculateFee，禁写死） |
| D13 | approve 两段式 + 两步态 | ✓ §3.2（exact approve，stepper，禁 infinite） |
| D14 | 锁仓倒计时 + 到期可提 | ✓ §3.4（边界 >=，累加顺延提示） |
| D15 | Cards 双 Tab + 移除 Admin 发卡 + 自动发放说明 | ✓ §3.5 |
| D16 | SIM 领取表单 + pendingSync | ✓ §3.5 Tab2 |
| §七8/9 | approve 交互稿 + 锁仓/自动发卡展示 | ✓ §3.2/3.4/3.5 |
| 对账反转 | 充值/提现/付账 pending 态（不据 200 显示终态） | ✓ §3.1/3.3 |

---

## 7. 移交 arch-review（风险点 / 待 implement 敲定）

1. **exact approve 非原子**：approve 成功 / deposit(pay) 失败会留精确额度授权——arch-review 评估是否需失败后归零提示。
2. **精度 6 位全链路**：format.ts/constants.ts/depositApi 任一漏改差 10^12 倍（资损红线）；design 已锁「金额一律 6 位、usdtDecimals 从 deployments 读」。
3. **锁仓累加语义易误解**：充值顺延 30 天，UI 必须明示（§3.4），否则用户投诉「倒计时变长」。
4. **pending 态 reorg 回退**：UI 不缓存未确认态当真；arch-review 确认 K 块确认数（后端占位 5）与前端轮询/事件订阅策略。
5. **WalletAuth 签名格式未定**：EIP-712 vs EIP-191 / 字段顺序 / nonce 来源——implement 与后端 `signatures.go` 敲定；design 已给 UX 规范（与交易签名区分、拒签态）。
6. **421614 真上链未做**：`deployments/arbitrum_sepolia.json` 不存在，端到端 Arbitrum 验收（D17）阻塞于合约先上链；本地 31337 可先验全链路。
7. **31337 地址漂移**：contracts.ts 与 deployments/hardhat.json 不一致；建议「deployments json → contracts.ts」单一出口脚本（实现层，arch-review 关注）。
8. **Oracle ABI 未导出（接链改动面）**：`config/abis/index.ts` 当前导出 6 个（UserRegistry/Deposit/ServiceManager/Payment/FeeManager/TrafficCardNFT）但**缺 OracleABI**（`Oracle.ts` 文件已在，FeeManager 已导出做对照）；implement 需补 `export { OracleABI }`（D11 验收 + 若读结算状态需要）。

---

## 8. TwoStepAction 状态机（B5 · approve→action 两笔串行）

`useTxState`（`useTransactionFlow.ts`）是**单笔五态**（idle/pending-signature/pending-confirmation/success/error），撑不起两笔串行 + allowance 跳步 + approve 成功后 action 失败的回退。`TwoStepAction` 在其上编排两个 `useTxState`（approveTx/actionTx）+ allowance 读值。

```
idle (输入金额，校验 amount≥10 USDT & ≤钱包USDT余额)
  ├─[allowance(user,spender) ≥ 需求额]──────────────▶ 跳过 Approve，直达 ②
  └─[allowance < 需求额] ① Approve:
        approve-sign → (拒签→idle, toast「授权已取消」)
                     → approving → (失败→approve-failed，回 idle 可重试 Approve)
                                  → confirming-approval → ②
   ② Deposit/Pay:
        action-sign → (拒签→★approved-idle)
                    → acting → (失败→★approved-idle)
                             → 成功 → POST 意向 → pending（通用三态）
```
- **★ 关键回退分支**：Approve 已成功后，Deposit/Pay 签名被拒或链上失败 → 回退到 **`approved-idle`（已授权可重试存入）**，**不回 Approve、不 re-approve**；重试时 `allowance ≥ 需求额` 直接跳 ②。UI：Step1「已授权 ✓」、主按钮「重试存入/付款 N USDT」。
- 需求额：充值 `amount`；付账 `amount + calculateFee(amount)`（直读合约）。spender：Deposit / Payment。**充值/付账复用同一状态机**，仅 spender/需求额/文案不同。
- 完整画法 + 视觉态见 DESIGN.md「接链交互模式 §1 + §7」。

## 9. confirmed 信号来源（B6 · 前端不监听链事件做终态）

| 信号 | 来源 | 策略 |
|------|------|------|
| **余额** | 链上读 `Deposit.getDepositAmount(addr)` = source of truth | 已只含 confirmed 本金（K 块逻辑在合约/后端）；staleTime 15s |
| **账单 is_paid** | 轮询后端 `GET /api/bills/:wallet` | pending(paying) 期 `refetchInterval≈5s`，转 paid 停轮询 |
| **充值/提现历史终态** | 轮询后端 `GET /api/deposit/:wallet` | 同上，pending→confirmed 停轮询 |

- 前端 **不** `useWatchContractEvent` / 不订阅 `BillPaid`/`DepositMade`/`DepositWithdrawn` 置终态；事件回填是后端 event_sync 职责（handoff §1/§4）。**K 块逻辑全留后端，前端只读 confirmed 字段。**
- **getLogs 限流修复**：`useTrafficCards` 现 `getLogs({fromBlock:0n})` 全量扫块在 Arbitrum 公共 RPC 必限流/超时，且 `catch→setTokenIds([])` 是 silent failure。修复二选一：① 后端给 NFT 列表端点（首选）；② 限定 `fromBlock` 为合约部署块号窗口。**必补 error 态**（「加载失败，重试」+ refetch），**禁 catch 后静默置空**（区分「真无卡」vs「加载失败」）。
- **reorg 缓存**：pending 相关 react-query 短 `staleTime`(~2s)+短 `refetchInterval`(~5s)；confirmed 后放长 staleTime、停轮询。

## 10. 测试策略（B4 · 项目零测试设施，资损敏感必补）

> **B3 绿基线前置**：implement **T0** 先清 `RegisterSheet.tsx:22` 'isSuccess' 未用 TS6133（现 `tsc -b && vite build` 过不了），跑到 **全量 tsc 绿** 才有回归基线，否则接链重写新旧错混淆。

**设施**：devDeps 装 `vitest` + `@testing-library/react` + `@testing-library/jest-dom`（jsdom 环境）；接链重写期同步补测。

**单元（必测，资损红线）：**
- `utils/format.ts`：`parseUnits("1.5", 6)===1500000n`、`formatAmount(1500000n, 6, 2)==="1.50"`；**6 位精度**（旧默认 18 → 接链必传 6，差 10^12 倍）；边界 0 / 大额。
- `config/constants.ts`：`MIN_DEPOSIT_USDT` 现 `100n*10n**18n`，含**两个独立改动**：① **精度修复**——18 位精度 bug，应为 6 位最小单位（接链 USDT 6 位，18→6 差 10^12 倍）；② **值对齐**——值 100→10，对齐链上 `require amount≥10 USDT` 下限。断言目标 `MIN_DEPOSIT_USDT === 10n * 10n**6n`。
- `billingApi.toBill`：改 **bigint** 后断言 total 用 `BigInt` 加减、无浮点误差、大额不溢出（见 §3.3）。
- approve 额算法：充值 `=amount`、付账 `=amount + calculateFee(amount)`（exact，禁 infinite）。
- `LockCountdown` 解锁判定：`now >= expiry`（**`>=` 边界**，到期即解锁；`now===expiry` 必须可提）。
- `parseContractError`：revert reason 映射、User rejected、insufficient funds 各命中。

**组件（mock wagmi）：**
- `TwoStepAction`：状态机全分支——allowance 跳步、approve 拒签回 idle、**approve 成功+action 失败→approved-idle（不 re-approve）**、成功进 pending。
- `TxStatusBadge`：pending/confirmed/failed 三态渲染（含 reorg vs revert 文案区分）。
- pending 渲染：首屏 skeleton（mock 慢 RPC 读链未回）、pending 不染绿。

**集成/冒烟**：本地 **31337 全链路冒烟**（充值两步态 → pending → 读链 confirmed；付账；提现；锁仓边界）。

**后置（不计入 web DONE）**：Arbitrum Sepolia(421614) 真链端到端（D17）**阻塞于合约上链**，合约上链后补。

## 11. arch-review 阻塞闭合对照表

| # | 阻塞 | 闭合处 | 状态 |
|---|------|--------|------|
| B1 | WalletAuth 签名频次留虚 | DESIGN.md §6 + Decisions Log「会话级一次铁律(nonce+时间窗,EIP-712)，禁每次签」；design §3.7 | ✓ 锁死铁律，nonce 来源标跨端对齐 signatures.go |
| B2 | 金色卡内 ≈2:1 不达标 | DESIGN.md 文字色阶 + 「金色用色铁律」表 + AmountDisplay 默认按底色分流（卡内 navy / 深底金）；WCAG 标注 | ✓ 卡内金额改 navy #0C2340 |
| B3 | tsc 基线红 | design §10 测试策略「T0 前置」+ 落地计划注明全量 tsc 绿前置 | ✓ 记 implement T0 |
| B4 | 零测试设施 | design §10「测试策略」新增（vitest + 单元/组件/31337 冒烟，Arbitrum 后置） | ✓ |
| B5 | TwoStepAction 状态机未画 | DESIGN.md §1 状态机 + design §8（含 ★approved-idle 回退不 re-approve、allowance 跳步） | ✓ |
| B6 | confirmed 来源未定 + getLogs 限流 | DESIGN.md 三态表「confirmed 信号来源」+ Cards 取数修正；design §9 | ✓ 读链 getDepositAmount + 轮询后端 status；getLogs 迁后端/限窗 + error 态 |

**关键 ⚠️ 闭合**：pending 超时兜底(Arbiscan+「约1-2分钟」+「可安全离开」)✓ / billingApi bigint(§3.3)✓ / signedPost 非全局拦截器(§3.7)✓ / interest 恒 0 去「利息」(§3.4 DESIGN.md)✓ / 换肤依赖接链结构→implement「接链定结构→换肤上色」✓(见 §12) / web DONE=本地 31337 全链路绿,Arbitrum D17 后置✓ / 香槟金明确做(金线渐变停靠点 0/50/100%)✓ / approve 失败中间态+WalletAuth 已签名态+pending 首屏 loading+320 窄屏 44px+pending 文案统一+failed 区分 reorg/revert+reorg 短 staleTime ✓(DESIGN.md §7)。

## 12. 落地计划补充（implement 顺序硬约束）

- **绿基线前置（T0）**：清 `RegisterSheet.tsx:22` TS6133 → 全量 `tsc -b && vite build` 绿。
- **隐性串行（换肤依赖接链结构）**：先「接链定 DOM 结构」（Deposit 两步态/余额卡/Cards 双 Tab），后「换肤上色」（吃 `:root` 变量）；二者非正交并行，否则换肤返工。
- **implement 阶段始终串行**（一个 Task 完成审查后再派下一个）。
- **web DONE 验收边界**：= 本地 **31337 全链路绿**（充值/提现/付账三态 + 锁仓 + 手续费读链 + WalletAuth 会话签名）；Arbitrum 端到端(D17)+对账三态真链行为 = **后置强制验收**，阻塞于合约上链，不计入 web DONE 也不成孤儿。

## 13. 产物
- `docs/design/linkworld/DESIGN.md` — 「接链交互模式」章节（v2 补：金色用色铁律/TwoStepAction 状态机/confirmed 来源/pending 超时兜底/交互态补全）+ Decisions Log 四行（设计源真理，长期）。
- `docs/pipeline/stages/design.md` — 本文件 v2（按 arch-review 返工，新增 §8-§12）。
