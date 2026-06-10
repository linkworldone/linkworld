# Task 03 — T2 精度 6 位 + 常量对齐（web 3/3）

> 子项目 web(3/3) | 状态 DONE | 2026-06-10

### 产出文件
- packages/web/src/utils/format.ts（parseUnits/formatAmount 默认精度 18→6，提取 USDT_DECIMALS=6；修 formatAmount 尾点 bug）
- packages/web/src/config/constants.ts（MIN_DEPOSIT_USDT 100n*10n**18n→10n*10n**6n；SUPPORTED_CURRENCIES→["USDT"]；PLATFORM_FEE_RATE 标 @deprecated）
- packages/web/src/services/api/billingApi.ts（toBill.total 浮点 parseFloat→BigInt 加减，totalAmount 全程 6 位最小单位字符串）
- packages/web/src/{utils/format.test.ts, config/constants.test.ts, services/api/billingApi.test.ts}（新建 PREC/BILL 用例）

### git commit
bb6aabe feat: web T2 精度 6 位 + 常量对齐(MIN_DEPOSIT 10 USDT/SUPPORTED USDT/billingApi bigint)

### TDD
先红后绿：先写 3 测 → 10 failed（parseUnits 返 1e19/MIN 100e18/SUPPORTED 含 ETH/toBill 浮点丢精度）→ 改源 → 17 passed；红还暴露并修了 formatAmount(x,6,0) 尾点 bug。

### 测试结果
npx tsc --noEmit 0 error。npm run build ✓。npx vitest run：4 files / 17 passed（含 T0 smoke 不回归）。用例 PREC-01(parseUnits/formatAmount 6 位 8 例)/PREC-02(MIN_DEPOSIT===10n*10n**6n)/PREC-03(SUPPORTED===["USDT"])/BILL-01(toBill bigint 含超 MAX_SAFE_INTEGER 大额)。主 Agent 已独立复跑确认 tsc 0+build ✓+17 测。

### code-simplifier
USDT_DECIMALS=6 提常量；toBill 改 bigint 消浮点资损；调用方仍可显式传精度。

### spec review
按 design v2 §3.3/§7.2/§10 + arch-review B4(精度资损红线/billingApi bigint) 执行。MIN_DEPOSIT 两处独立修正(精度 18→6 + 值 100→10 对齐链上 require)。PLATFORM_FEE_RATE 标 @deprecated(删留 T9，避免破坏现有引用)。未越界 approve/对账/换肤。

### 设计还原
T2 精度层无 UI 还原。design §3.3/§7.2 6 位全链路 + billingApi bigint 落地。

### 复用检查
复用 viem parseUnits、contracts getUsdtDecimals；format 默认 6 位供 T3 approve 复用。

### 设计稿对照
数值对照：parseUnits("10",6)===10_000_000n ✅；MIN_DEPOSIT_USDT===10n*10n**6n ✅；SUPPORTED===["USDT"] ✅；toBill bigint 大额不丢精度 ✅；tsc 0/build ✓/17 测 ✅。

### 新增组件
无新增业务组件。新增测试 3 文件 + USDT_DECIMALS 常量。

### 新增色值
无（T2 精度层，换肤留 T6+）。

### ⚠️ 遗留（带入 T3/T9，重要）
- T3：approve exact 额用 parseUnits(amount, getUsdtDecimals) 即可（精度已统一 6）。
- **T9 重要**：toBill.totalAmount 语义已变（美元字符串→6 位最小单位字符串）。Billing.tsx/BillDetail.tsx 仍 parseUnits(bill.totalAmount)（会二次缩放 10^6=错误支付额）+ formatUSD(operatorFee)（当美元）——T9 必须改：展示走 formatAmount(total, usdtDecimals)、支付额直接用最小单位 bigint 去掉二次 parseUnits。否则支付额错 10^6 倍（资损）。
- T9：PLATFORM_FEE_RATE @deprecated → useFeeRate 读链(getFeeRate 基点 150/10000)+calculateFee+FeeBreakdown 后删 + 改 BillDetail 写死 "2.5%" 文案。
