# Stage: plan — web 重构接链+换肤 任务拆解（子项目 web 3/3）

> **状态**: completed | **日期**: 2026-06-10 | **Gate**: 1（plan） | **子项目**: web(3/3) | **角色**: 架构师/工程经理 | **分支**: web/deep-blue-gold-refactor
> **输入**：design v2（`docs/pipeline/stages/design.md`，§3 接链交互/§8 TwoStepAction 状态机/§9 confirmed 来源/§10 测试策略/§11 闭合表/§12 落地硬约束）+ DESIGN.md（深蓝金视觉+金色铁律+接链交互模式）+ arch-review PASS（§七 带入 plan 前置）+ 两份 web delta（`web-alignment-surface.md`/`theme-migration.md`）+ 真实源码逐文件核对（codegraph）。
> **本文件只产计划文档，不写代码。** 任务序列 T0-T12 **串行**，遵守 CEO 串行依赖「T0 绿基线 → 接链基建定结构 → 换肤上色」（design §12）。
> **注**：本 stage 文件按子项目串行复用——合约(1/3)/后端(2/3) plan 前版见 git history，本版为 web(3/3)。
> **web DONE 边界**：= 本地 **31337 全链路绿**；Arbitrum 端到端(D17)+对账三态真链行为 = 后置强制验收，阻塞于合约上链，**不计入 web DONE**。

---

## 0. 已核对真实代码（plan 头部强制 · codegraph 实读）

| 核对项 | 结果 | 出处 |
|--------|------|------|
| `RegisterSheet.tsx:22` 解构 `isSuccess` 未用（TS6133，build=`tsc -b && vite build` 现红） | ✓ T0 前置确认 | `packages/web/src/components/wallet/RegisterSheet.tsx:22` |
| `config/contracts.ts` 仅 31337(7地址，已漂移)+16601(全0)；type 无 USDT 字段 | ✓ | `contracts.ts:1-46` |
| `deployments/hardhat.json` 有 `proxies`(7合约)/`implementations`/`usdt`/`usdtDecimals:6`/`abiHash` | ✓ 单一出口源齐备 | `packages/contracts/deployments/hardhat.json` |
| `abis/index.ts` 导出 6 个，**缺 OracleABI**（`Oracle.ts` 文件已在） | ✓ T1 补 | `config/abis/index.ts` |
| `useContractDeposit` 现 `deposit(amountEth){value:parseEther}`（18位+payable+无amount参） | ✓ T1/T3 改 | `useDepositContract.ts:23-30` |
| `useContractPayBill` 现 `payBill(billId,value){args:[billId],value}` | ✓ T1/T3 改 | `usePaymentContract.ts:11-19` |
| `format.ts` `parseUnits/formatAmount` 默认 `decimals=18` | ✓ T2 改 6 位 | `utils/format.ts:5,11` |
| `constants.ts` `PLATFORM_FEE_RATE=0.025`/`MIN_DEPOSIT_USDT=100n*10n**18n`/`SUPPORTED_CURRENCIES=["USDT","ETH"]` | ✓ T2 全改 | `config/constants.ts:1-6` |
| `depositApi` `recordDeposit` 用 `parseEther`(18位)；`recordWithdraw(txHash)` 凭哈希记账 | ✓ T2/T4 改 | `services/api/depositApi.ts:31,46` |
| `billingApi.toBill` `total=(parseFloat+parseFloat-parseFloat).toFixed(2)`（浮点，资损） | ✓ T2 改 bigint | `services/api/billingApi.ts:32-34` |
| `useTxState` 单笔五态（idle/pending-signature/pending-confirmation/success/error） | ✓ T3 编排两实例 | `useTransactionFlow.ts:3,55` |
| `useTrafficCards` `getLogs({fromBlock:0n})` + `catch→setTokenIds([])` silent failure | ✓ T4 修 | `useTrafficCard.ts:91-108` |
| `useIssueMonthlyCards`（Admin 发卡，onlyOracle）存在，Cards 页调用 | ✓ T8 移除按钮 | `useTrafficCard.ts:180-196` |
| `useDeposit` 双写 `recordToBackend`+`currency:"ETH"`；`usePayBill` 双写 `recordPayment` | ✓ T4 去双写 | `useDeposit.ts:15,42` / `useBilling.ts:52` |

---

## 1. 页映射与任务拆解

> implement 顺序硬约束（design §12）：**T0 绿基线 → 接链基建定 DOM 结构（T1-T5）→ 换肤上色（T6-T11）→ 测试构建（T12）**。换肤依赖接链结构非正交并行；implement 阶段始终串行（一个完成审查后再派下一个）。`scope` 取值：infra / 接链 / 换肤 / 接链+换肤。

| 任务ID | 代码页 / 模块 | 范围 | 依赖 | gate（验收门槛） |
|--------|--------------|------|------|------------------|
| **T0** | 测试设施 + `RegisterSheet.tsx`（清 TS6133） | infra | — | 全量 `tsc -b && vite build` 绿；`npx vitest run` 可跑（0 用例亦可） |
| **T1** | `config/{chains,contracts,wagmi}.ts` + `config/abis/*`（重生成+补 OracleABI） | 接链 | T0 | `tsc` 绿；31337 地址=deployments；421614 占位；OracleABI 导出 |
| **T2** | `utils/format.ts` + `config/constants.ts` + `services/api/{depositApi,billingApi}.ts`（精度+常量+bigint） | 接链 | T1 | 单元测过（format 6位 / MIN_DEPOSIT / toBill bigint） |
| **T3** | `hooks/contracts/{useDepositContract,usePaymentContract}.ts` + 新增 USDT/ERC20 ABI + `useApprove`/`useAllowance` + `TwoStepAction` | 接链 | T2 | `TwoStepAction` 组件测全分支过 |
| **T4** | `hooks/{useDeposit,useBilling}.ts` + `services/api/{depositApi,billingApi}.ts`（pending）+ `useTrafficCard.ts`（getLogs 修）+ `TxStatusBadge` | 接链 | T3 | 去双写；pending 态；getLogs error 态；徽章三态测过 |
| **T5** | `services/api/client.ts` + 新增 `signedPost`/`useWalletAuth`（EIP-712 会话签名） | 接链 | T4 | signedPost 封装不在 axios 拦截器调 hook；会话级一次 |
| **T6** | `index.css` + `tailwind.config.ts` + 字体 + 新增 `ui/card`·`ui/input` | 换肤 | T1 | CSS 变量单一出口；`ui/card`/`ui/input` 落地；grep 无裸 HEX 定义残留 |
| **T7** | `pages/Deposit.tsx`（两步态+LockCountdown+pending）+ 新增 `LockCountdown`+`useLockExpiry` | 接链+换肤 | T3,T6 | 31337 充值两步态→pending；锁仓顺延提示 |
| **T8** | `pages/Cards.tsx`（双 Tab 去 Admin + NFT 读链 + SIM 领取） | 接链+换肤 | T4,T6 | Tab1 NFT 读链+error 态；Tab2 SIM→pendingSync；无 Admin 按钮 |
| **T9** | `pages/{Billing,BillDetail}.tsx`（paying 态 + `FeeBreakdown` 读链）+ 新增 `FeeBreakdown`/`useFeeRate` | 接链+换肤 | T4,T6 | paying 态 info 蓝禁绿；手续费读链；明细 bigint |
| **T10** | `pages/RegionDetail.tsx`（申请号码弹层手续费） | 接链+换肤 | T9 | 申请弹层展示 `calculateFee` 读链值 |
| **T11** | `pages/{Landing,Dashboard,Services,Notifications}.tsx` + `layout/{AppLayout,Header,TabBar}.tsx` + `shared/{AmountDisplay,EmptyState,GuardCard,StatusBadge}.tsx` + `wallet/{ConnectButton,RegisterSheet}.tsx`（换肤 + emoji→lucide） | 换肤 | T6,T7,T8,T9,T10 | grep 无旧色值（C8）/无 emoji 图标（C9）；9 页 100% 深蓝金（C7） |
| **T12** | 全量测试 + 构建 + 31337 冒烟 | infra | T0-T11 | `vitest run` 绿 + `tsc -b && vite build` 绿 + 31337 全链路冒烟 |

**接链基建任务**（先定结构，T1-T5）：链配置/ABI 重生成（T1）、精度常量（T2）、approve hooks+TwoStepAction（T3）、对账重构+getLogs 修（T4）、WalletAuth（T5）。**9 页**（Landing/Dashboard/Services/RegionDetail/Deposit/Billing/BillDetail/Cards/Notifications）由 T7-T11 覆盖；**换肤任务**（后上色，T6-T11）。

---

## 2. 组件复用审查

| 组件 | 现状（复用/重着色/新建/改 props） | 来源 / 处理 |
|------|--------------------------------|-------------|
| `ui/button` | 复用重着色 | cva + CSS 变量，吃 `:root` 即生效（T6 改变量）；ConnectButton 建议改用 `<Button>` |
| `ui/badge` | 复用重着色 | 同上；交易/状态徽章基底 |
| `ui/tabs` | 复用重着色 | 同上；Cards 双 Tab 用 `TabsList variant=line`（T8） |
| `ui/card` | **新建** | DESIGN.md「缺失原子组件」：暖米白 `bg-surface-card` + 1px 金线 + radius-lg + 轻阴影；收口 Dashboard/Deposit/Services/Billing/Cards 手写卡（T6） |
| `ui/input` | **新建** | `bg-surface-input` + navy 文字 + royal focus ring；统一 RegisterSheet/SIM 表单/搜索框（T6） |
| `TwoStepAction` | **新建** | approve→action 两步态 + stepper + allowance 跳步 + ★approved-idle 回退（不 re-approve）；编排两个 `useTxState`；充值/付账复用（T3，design §8） |
| `TxStatusBadge` | **新建** | 通用交易三态徽章 pending/confirmed/failed（含 reorg vs revert 文案区分）（T4，DESIGN.md 三态表） |
| `LockCountdown` | **新建** | 读 `getLockExpiry` → 倒计时/解锁态；解锁判定 `now>=expiry`（`>=` 边界）（T7，design §3.4） |
| `FeeBreakdown` | **新建** | 费用明细（小计+手续费+合计），读链 `getFeeRate`/`calculateFee`；total bigint（T9，design §3.6） |
| `useFeeRate` | **新建 hook** | 读 `FeeManager.getFeeRate()`（基点/10000）；loading/失败 skeleton 不写死（T9） |
| `useApprove`/`useAllowance` | **新建 hook** | ERC20 `approve`(exact，禁 infinite)/`allowance`/`balanceOf`/`decimals`（T3） |
| `useLockExpiry` | **新建 hook** | 读 `Deposit.getLockExpiry(addr)`（T7） |
| `useWalletAuth`/`signedPost` | **新建** | EIP-712 会话级签名一次 + 内存缓存；`signedPost(path,body)` helper（**不在 axios 全局拦截器调 hook**）（T5，design §3.7） |
| `AmountDisplay` | **改 props** | 默认 `colorClass` `text-status-warning` → **按底色分流**（卡内 navy `text-text-primary` / 深底金 `text-on-dark-gold`，B2 覆盖旧「改吃金」）；增 `tone="auto\|gold-on-dark"`；币种 ETH→USDT；6 位精度（T2 改精度/T11 改色） |
| `EmptyState` | **改 props** | `icon` prop `string`→`LucideIcon`（DESIGN.md D5）（T11） |
| `GuardCard` | **改 props** | `icon` prop `string`→`LucideIcon`（DESIGN.md D5）（T11） |
| `StatusBadge` | 复用重着色 | 吃 status token（T11 改色） |
| `BottomSheet` | 复用 | vaul Drawer，SIM 领取/申请弹层优先复用（T8/T10） |

---

## 3. 色值映射表

> **CSS 变量单一出口**（DESIGN.md「色值出口统一方案」）：`index.css :root` 为唯一真源，`tailwind.config.ts` token 全改 `var(--...)`；业务组件用语义 Tailwind 类，禁裸 HEX。**金色铁律（B2）**：金 #D4AF37 在暖米白 #F7F3EA 上 ≈2:1 不达标——卡内金额/文字一律 navy；金色仅深底文字/CTA 填充/激活态/卡片金线（非文本用途豁免）。

| 旧色值 | 新深蓝金 token | 涉及文件 / 页 |
|--------|---------------|--------------|
| `surface.DEFAULT #0a0a14` | navy 画布 `#0C2340`（body 渐变 `--bg-canvas`） | `tailwind.config.ts`（定义源）；全局背景 |
| `surface.card #0f0f1a` | 暖米白 `surface.card #F7F3EA` | `tailwind.config.ts`；Dashboard/Deposit/Services/Billing/Cards 卡 |
| `surface.secondary #1a1a2e` | `surface.input #EFE9DB` | `tailwind.config.ts`；输入框/凹陷区 |
| `surface.gradient.from/to #1a1a3e/#0f1a2e` | `--gradient-hero`/`--bg-canvas`（navy 渐变） | `Deposit.tsx`/`Cards.tsx` 余额卡（×3 文件） |
| `brand.blue #3b82f6` | `brand.royal #1E40AF`（链接/次级/focus）或 `gold #D4AF37`（激活/CTA/金额，按语义） | brand-blue ×12 文件：TabBar/Header/RegisterSheet/ConnectButton/Deposit/Services/Dashboard/Cards/Landing/Notifications/BillDetail/Billing |
| `brand.purple #8b5cf6` | 删 → 金 `#D4AF37` 或 navy 渐变 | brand-purple ×4 文件（Header 头像渐变球等） |
| `brand.cyan #06b6d4` | 删 | `tailwind.config.ts` |
| `status.warning #f59e0b` | 橙 `#F08C2E`（与金分离） | `tailwind.config.ts`；warning 语义处 |
| `fontFamily.orbitron` | **删**（弃 Orbitron） | `tailwind.config.ts` |
| `#3b82f6` 裸 HEX（RainbowKit accentColor） | 金 `#D4AF37`（accentColorForeground `#0C2340`） | `main.tsx` `darkTheme` |
| `AmountDisplay colorClass="text-status-warning"` | 按底色分流：卡内 navy `text-text-primary` / 深底金 `text-on-dark-gold`（**不无条件吃金**，B2） | `shared/AmountDisplay.tsx` |
| index.css `:root` 单套深紫 oklch（hue~270）+ body `"Inter"` | DESIGN.md shadcn 变量映射表新 oklch（`--background`/`--card`/`--accent`金/`--ring` 等）；删 `"Inter"` 走 Geist；加 Space Grotesk | `index.css` |
| 香槟金 #F0C75E | **明确做**：仅金线渐变高光端 `linear-gradient(135deg,#D4AF37 0%,#F0C75E 50%,#D4AF37 100%)`（停靠点锁 0/50/100%）；主金仍 #D4AF37，不作独立文字/填充色 | `ui/card` 金线（T6） |

---

## 4. 接链 / 对账 / 前置改动面

| 项 | 现状 | 改为 | 风险 |
|----|------|------|------|
| 链配置 chains.ts | 仅 zg(16600/16601)+hardhat(31337)，无 421614 | 新增 `arbitrumSepolia(421614)`（ETH 18位 nativeCurrency / 公共 RPC 或 `VITE_*` / `sepolia.arbiscan.io`）；保留 hardhatLocal | Arbitrum 公共 RPC 限流，event/读多调用建议 Alchemy/Infura key，与后端 deployments RPC 对齐 |
| contracts.ts 地址 | 31337 地址**已漂移**（与 deployments/hardhat.json 不一致）；16601 全 0；type 无 USDT | 31337 用 `deployments/hardhat.json.proxies` 7 地址覆盖；删 16601 换 421614 全 0 占位（`getContractAddress` 全 0 抛错沿用）；type 增 USDT；**建议「deployments json → contracts.ts」单一出口脚本/构建期导入** | 硬编码 TS 与 deployments 双源易漂移（31337 已发生） |
| ABI 重生成 | 旧 0G 模型快照，selector 已变（`deposit`/`payBill` 去 payable 加 amount/精度变） | 从 `artifacts/contracts/<Name>.sol/<Name>.json` 重生成，据 `deployments.abiHash` 比对；**补 `export { OracleABI }`**（缺，D11）；新增 ERC20/MockUSDT ABI（approve/allowance/balanceOf/decimals） | 旧 ABI 调用 revert；`value` 传 ETH 被非 payable 拒；ABI+hook 必须同步改 |
| USDT approve 两步 | 充值/付账全程 `value`（原生币），无 ERC20/approve/allowance | `useApprove`/`useAllowance`；充值前 `approve(Deposit,amount)`；付账前 `approve(Payment,amount+calculateFee)`；**exact-amount 禁 infinite**（资损硬约束）；`useContractDeposit` 去 value 加 `args:[amount]`；`useContractPayBill` 去 value | approve+deposit/pay 两笔非原子（approve 成功 action 失败留精确授权，可接受）；approve 额算错 transferFrom revert |
| 精度 6 位 | `format.ts` 默认 18 位；`MIN_DEPOSIT_USDT=100n*10n**18n`；`depositApi` `parseEther` | 金额一律 6 位，`usdtDecimals` 从 deployments 读；`MIN_DEPOSIT_USDT=10n*10n**6n`（精度修+值 100→10 对齐链上 `≥10`）；`SUPPORTED_CURRENCIES=["USDT"]`；`depositApi` 改 `parseUnits(amount,6)` | **资损红线**：18 当 6 差 10^12 倍；format/constants/depositApi 任一漏改即金额错；implement grep `parseEther`/`decimals=18`/`10n**18n` 清零 |
| billingApi bigint | `toBill.total=(parseFloat+...-...).toFixed(2)`（6 位最小单位当元做浮点加减） | `BigInt(operatorFee)+BigInt(platformFee)-BigInt(trafficCardDeduction)`，展示经 `formatAmount(total,usdtDecimals)`；`totalAmount` 不落 number | 单位语义错 + 大额超 `MAX_SAFE_INTEGER`（资损升） |
| 对账 pending | 充值/付账双写自述即记账（`recordDeposit`/`recordPayment(txHash)`）；提现凭 txHash `recordWithdraw` | 充值/付账/提现仅 POST pending 意向，**不据 200/txHash 置终态**；Bill status 扩 `unpaid/paying/paid/overdue`，`paying`=info 蓝+Loader2 禁绿；余额读链 `getDepositAmount`；is_paid/历史轮询后端 status（pending `refetchInterval≈5s`，confirmed 停） | 不据成功显示终态铁律；reorg 短 `staleTime≈2s` 不缓存未确认；前端**不**监听链事件做终态（后端 event_sync 职责） |
| getLogs 修 | `useTrafficCards` `getLogs({fromBlock:0n})` 全量扫块 + `catch→setTokenIds([])` silent failure | **二选一（implement 定）**：① 后端 NFT 列表端点（首选）；② 限定 `fromBlock` 为合约部署块号窗口。**必补 error 态**（「加载失败，重试」+refetch），禁 catch 静默置空（区分真无卡 vs 加载失败） | 公共 RPC 限流/超时；silent failure 误导「没有卡」 |
| WalletAuth | `client.ts` 裸 axios 无签名头 | 写端点带钱包签名头（读端点不加）；**EIP-712 会话级签名一次**（nonce+时间窗，内存缓存复用，禁每次签）；`signedPost(path,body)` helper（内部 `await signTypedData` 取/签缓存 → 带头调 axios），**不在 axios 全局拦截器调 React hook**（必崩） | nonce 来源（后端下发 vs 前端时间窗）+ EIP-712 字段/domain **跨端待与后端 `signatures.go` 对齐**，web 单方不可闭环 |
| 手续费读链 | 写死 `PLATFORM_FEE_RATE=0.025` | 删写死；`useFeeRate` 读 `getFeeRate()`（150=1.5%，/10000）；精确费额 `calculateFee(amount)`；展示位 RegionDetail 申请弹层 + Billing/BillDetail；loading/失败 skeleton 或「--」不写死兜底 | 基点 150 展示需 /10000，别和小数 0.015 混 |

---

## 5. 实现红线与前置

| 红线/前置 | 内容 | 关联任务 | 性质 |
|-----------|------|----------|------|
| **T0 绿基线前置（B3）** | 先清 `RegisterSheet.tsx:22` `isSuccess` TS6133（现 `tsc -b && vite build` 过不了），跑到**全量 tsc 绿**才有回归基线；同步搭 vitest（`vitest`+`@testing-library/react`+`@testing-library/jest-dom`，jsdom）。否则接链重写新旧错混淆。 | T0 | 硬约束 |
| **WalletAuth 跨端待对齐（B1）** | 会话级一次 + EIP-712 两条铁律 design 锁死，但 **nonce 来源 + 字段/domain 跨端待与后端 `signatures.go` 对齐**，web 单方不可闭环。**禁降级到「每次写操作签」**（唯一能压门槛的杠杆）。 | T5 | 跨端待对齐 |
| **getLogs 二选一（B6）** | 后端 NFT 列表端点 vs 限定 fromBlock 窗口，由 implement 定；无论哪种**必补 error 态**，禁 catch 静默置空。 | T4/T8 | implement 定 |
| **web DONE 边界** | = 本地 **31337 全链路绿**（充值/提现/付账三态 + 锁仓 + 手续费读链 + WalletAuth 会话签名）；Arbitrum 端到端(D17) + 对账三态真链行为 = **后置强制验收**，阻塞于合约上链（`deployments/arbitrum_sepolia.json` 不存在），不计入 web DONE 也不成孤儿。 | D17 后置 | 验收边界 |
| **金色铁律（B2）** | 卡内金额/文字一律 navy `#0C2340`（≈12:1）；金色仅深底文字（≈6.5:1）/CTA 填充/激活态/卡片金线（非文本用途）。`AmountDisplay` 默认按底色分流，**禁无条件吃金**。 | T11 | 视觉红线 |
| **对账不据 200/txHash 置终态** | 充值/提现/付账三处统一 pending「处理中 · 约 1-2 分钟」（不暴露 K 块）；pending 不染绿、不计入可用余额；超时（>~2min）兜底 Arbiscan 逃生链接 +「可安全离开，到账后通知」；reorg 回退不缓存当真。 | T4/T7/T9 | 对账红线 |
| **精度 6 位资损红线** | format/constants/depositApi 任一漏改差 10^12 倍；implement 必 grep `parseEther`/`decimals=18`/`10n**18n` 清零并以单元测固化。 | T2/T12 | 资损红线 |
| **TwoStepAction ★approved-idle 回退** | approve 成功后 action 失败/拒签 → 回 `approved-idle`（已授权可重试），**绝不 re-approve**（allowance≥需求额直跳 ②）。 | T3 | 状态机红线 |
| **隐性串行（换肤依赖接链结构）** | 先接链定 DOM 结构（Deposit 两步态/余额卡/Cards 双 Tab），后换肤上色；`main.tsx`/`wagmi.ts`/`Deposit.tsx`/`Cards.tsx` 接链与换肤都碰，串行处理避免返工。implement 阶段始终串行。 | T1-T11 | 执行顺序 |

---

## 6. 自检（writing-plans Self-Review）

| 自检项 | 结论 |
|--------|------|
| **spec 覆盖** | design §3.1-3.7 / §8 / §9 / §10 / §11 / §12 + DESIGN.md 接链交互 §1-§7 + 两份 delta 全部映射到 T0-T12；arch-review §七 4 项前置 → T0(1)/T5(2)/T4·T8(3)/T12·DONE 边界(4)。 |
| **金色铁律 B2** | 第 2/3/5 节均落（AmountDisplay 分流 + 卡内 navy）。 |
| **WalletAuth B1** | T5 + 第 5 节 WalletAuth 红线（会话级一次+EIP-712+signedPost 非拦截器）。 |
| **精度/bigint** | T2 + 第 4/5 节精度红线（format/constants/depositApi/billingApi.toBill）。 |
| **类型一致** | `TwoStepAction`/`TxStatusBadge`/`LockCountdown`/`FeeBreakdown` 在 §1/§2 命名一致；`signedPost`/`useWalletAuth`/`useFeeRate`/`useLockExpiry`/`useApprove`/`useAllowance` 一致。 |
| **每页一任务/依赖顺序/可验收/无超范围** | T0-T12 各对应单一模块面，依赖链 T0→T1-T5→T6-T11→T12 串行无环；每任务 §1 表带 gate 验收门槛；范围限 packages/web 接链+换肤，无超范围扩张。 |

> **下一步**：进 implement，按 T0→T12 **串行**派 subagent（一个完成审查后再派下一个），implement T0 先清 tsc + 搭 vitest。
