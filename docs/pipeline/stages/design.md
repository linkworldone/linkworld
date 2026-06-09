# Stage: design — web 重构接链 设计分析（子项目 web 3/3）

> **状态**: completed | **日期**: 2026-06-09 | **Gate**: 1（设计） | **子项目**: web(3/3) | **角色**: 产品设计师 | **分支**: web/deep-blue-gold-refactor
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
- **本轮修订**：① Product Context `0G Chain 生态` → `Arbitrum Sepolia(421614) + ERC20 USDT 6 位`；② 香槟金 `#F0C75E`（PRD E18 提及）定位为「金线渐变高光端，可选」，**主金仍 #D4AF37，不引入新主色**（R12 视觉不变）。
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
**铁律**：可用余额只算 confirmed；pending 单列「处理中 +N USDT」弱化不染绿；绝不据 tx 成功 / HTTP 200 显示终态成功。

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

### 3.4 锁仓倒计时（Deposit 提现区）
- 合约依据：`getLockExpiry(addr)` 时间戳；提现 `require(now >= expiry)`（边界 `>=`）；**每次充值 expiry += 30d 顺延**；提现归 0。
- 状态：`expiry==0` 无锁仓 / `now<expiry` Withdraw 禁用 + 倒计时「锁仓中 · 剩余 Dd Hh Mm」+ Lock / `now≥expiry` 可点 + Unlock +「锁仓已满，可提取本金+利息」。
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
- 减打扰：建议会话级一次签名（nonce/时间窗），格式 implement 与 `signatures.go` 敲定。

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
| 改 props | `AmountDisplay` | 默认 colorClass `text-status-warning`→`text-brand-gold`；币种 ETH→USDT；6 位精度 |
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
| C10/E18 | 金值 #D4AF37 / 字体 / 渐变停靠点 / 米白色阶 / 色值单一出口 | ✓ DESIGN.md 已锁定，香槟金 #F0C75E 定为高光端可选 |
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

---

## 8. 产物
- `docs/design/linkworld/DESIGN.md` — 新增「接链交互模式」章节 + Product Context 修订 + Decisions Log 两行（设计源真理，长期）。
- `docs/pipeline/stages/design.md` — 本文件（本轮 design 分析 + 核对 + 移交）。
