# Task 05 — T4 对账重构 + service/NFT 读链路重写（web 3/3）

> 子项目 web(3/3) | 状态 DONE | 2026-06-10

### 产出文件
- packages/web/src/services/api/{depositApi,billingApi}.ts（recordDeposit/Withdraw→postDepositIntent/postWithdrawIntent 6 位去 tx_hash 终态；recordPayment→payIntent；toBill 加 paying 态）
- packages/web/src/hooks/{useDeposit,useBilling}.ts（recordToBackend→recordIntent；pending 感知 refetchInterval 5s/staleTime 2s 轮询）
- packages/web/src/hooks/contracts/{useServiceManager,useTrafficCard}.ts（运营商模型 getOperator/getActiveOperators + 逐卡 getCardInfo 重写，清 as never 占位）+ contracts/index.ts + useOperator.ts(useApplyNumber→后端意向)
- packages/web/src/components/shared/TxStatusBadge.tsx（新建：三态 pending 不染绿/confirmed/failed 区分 reorg vs revert）
- packages/web/src/types/index.ts（DepositRecord 加 status；Bill status 加 paying）
- 页面对账逻辑衔接（非 UI 改版）：Deposit/Billing/BillDetail/RegionDetail
- 多份测试文件

### git commit
e2a0b42 feat: web T4 对账重构(去双写/pending/事件驱动 confirmed) + getLogs 修 + service/NFT 读链路重写

### TDD
先红后绿：首轮 7 failed（postDepositIntent 不存在/status 未推 paying/TxStatusBadge import 失败）→ 实现后全绿。

### 测试结果
npx tsc --noEmit 0 error。npm run build ✓。npx vitest run：11 files / 44 passed（前序 30+新 14，无回归）。用例 REC-01/02(pending 意向去 tx_hash)/REC-03(不据 200 置 paying)/REC-04(余额读链 getDepositAmount)/LOG-01(getLogs 失败 isError 非静默空)/BADGE-01(pending 不染绿)。主 Agent 已独立复跑确认 tsc 0+build ✓+44 测。

### code-simplifier
recordIntent 统一去双写；TxStatusBadge 复用三态；service/NFT 按新 ABI 重写消除 as never 占位。

### spec review
按 design v2 §3.1/§3.3/§9 + handoff-web 对账契约 + arch-review B6 执行：三处去双写改 pending、IsPaid 等后端事件回填、余额读链 source of truth、账单轮询后端 confirmed 不前端监听、getLogs 限窗+error 态。未越界（WalletAuth 留 T5、换肤 token 留 T6、页面 UI 改版留 T7/T8/T9）。

### 设计还原
对账三态(pending 不染绿)+ TxStatusBadge 对齐 design §3.1；getLogs 修区分加载失败 vs 真无卡。

### 复用检查
复用 react-query refetchInterval、getDepositAmount 读链、新 ABI(getOperator/getCardInfo)；TxStatusBadge 供 T9 paying 渲染复用。

### 设计稿对照
数值对照：对账去双写 3 处(充值/提现/付账)✅；余额读链 getDepositAmount(REC-04)✅；getLogs 限窗 latest-5M 块+isError(LOG-01)✅；as never 清零(grep 验证，剩 2 处为 wagmi 动态数组类型限制非占位)✅；44 测 ✅；tsc 0/build ✓ ✅。

### 新增组件
新增 TxStatusBadge；新增 postDepositIntent/postWithdrawIntent/payIntent；service/NFT hooks 按新模型重写。

### 新增色值
无（组件用语义类，换肤 token 留 T6）。

### ⚠️ 遗留（带入 T5/T7/T8/T9/T10 + 跨端）
- T5：WalletAuth/signedPost 未做；pending 意向 POST 需带签名头。
- T7 Deposit：toast 已改 pending 文案；两步态/LockCountdown/pending UI 待 T7；Billing/BillDetail 的 payBill(_value) 调用行仍传 parseUnits 18 位（payBill 已忽略 _value）待 T9 清理。
- T8 Cards：useTrafficCards 暴露 isError/error，Cards.tsx 需用其区分「加载失败重试」vs「真无卡」；Admin 发卡移除归 T8。
- T9 Billing：Bill status 加 paying，UI 用 TxStatusBadge info 蓝+Loader2 禁绿。
- T10 RegionDetail：useApplyNumber 改后端意向；手续费展示待 T10；operatorApi.requiredDeposit parseEther(18 位)仍 T2 遗漏，改它同步 RegionDetail 展示。
- **跨端待对齐**：后端 ApiBill 需返回 pay_intent_tx_hash、deposit history 需返回 status 字段（paying/pending 推导依赖）。
