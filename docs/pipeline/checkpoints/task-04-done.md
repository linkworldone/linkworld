# Task 04 — T3 USDT approve hooks + TwoStepAction 两步态（web 3/3）

> 子项目 web(3/3) | 状态 DONE | 2026-06-10

### 产出文件
- packages/web/src/hooks/contracts/useUsdtContract.ts（新建：useAllowance(owner,spender) 读 MockUSDT.allowance staleTime 2s；useApprove() approve(spender, **exact amount 禁 MaxUint256**)，getUsdt(chainId)+parseUnits 6 位）
- packages/web/src/components/shared/TwoStepAction.tsx（新建：独立编排组件 + derivePhase 纯函数状态机；action 第二笔由 props 注入）
- packages/web/src/hooks/contracts/index.ts（补 export useAllowance/useApprove）
- packages/web/src/vitest.d.ts（新建：jest-dom matcher 类型）
- 两份测试文件（TwoStepAction + useUsdtContract）

### git commit
546f784 feat: web T3 USDT approve hooks + TwoStepAction 两步态组件(exact approve/allowance 跳步/approved-idle 回退)

### TDD
先红后绿：注入 args=[spender,MaxUint256] → TSA-01/04+hook exact 单测变红（证明 exact 断言守资损红线）→ 还原 exact → 全绿。

### 测试结果
npx tsc --noEmit 0 error。npm run build ✓。npm test：6 files / 30 passed（17 前序+13 新，无回归）。用例 TSA-01(allowance<amount 显 Approve+exact)/TSA-02(allowance≥amount 跳过 Approve)/TSA-03(approve 成功+action 失败回 approved-idle 不 re-approve)/TSA-04(exact 非 MaxUint256)+derivePhase 7 纯函数单测+useApprove exact 单测。主 Agent 已独立复跑确认 tsc 0+build ✓+30 测。

### code-simplifier
derivePhase 抽纯函数易测；TwoStepAction 编排两笔交易 + allowance，action 注入式复用（充值/付账）；exact approve 无 MaxUint256。

### spec review
按 design v2 §8 状态机 + §3.2 + arch-review B5(两步态)/exact approve 禁 infinite 执行。phase: idle→[allowance≥跳步]→approve-sign/approving/confirming-approval→action-sign/acting/confirming→done；approve 成功后 action 失败/拒签→approved-idle（approvedOnce 锁定，重试不 re-approve）。未越界（不接页面/不碰对账/WalletAuth/换肤 token，组件用语义类不写死 hex）。

### 设计还原
组件状态机对齐 design §8（含 allowance 跳步 + approved-idle 回退分支）；Stepper 用 lucide + 语义类（text-status-success/brand-blue），换肤 token 由 T6 提供，T3 不写死 hex。

### 复用检查
复用 T1 MockUSDT ABI/getUsdt、T2 parseUnits 6 位、现有 useTxState/ui-button/wagmi useReadContract·useWriteContract·useWaitForTransactionReceipt；TwoStepAction 充值/付账复用。

### 设计稿对照
数值对照：approve args=[spender, exact amount]（非 MaxUint256，TSA-04 断言）✅；allowance≥amount 跳过 Approve（TSA-02）✅；approve 成功+action 失败回 approved-idle（TSA-03）✅；30 测 ✅；tsc 0/build ✓ ✅。

### 新增组件
新增 useAllowance/useApprove hooks + TwoStepAction 组件（+derivePhase）。

### 新增色值
无（组件用语义类，深蓝金 token 定义留 T6 基建）。

### ⚠️ 遗留（带入 T7/T9）
- T7 充值 Deposit.tsx 接 TwoStepAction：spender=Deposit 地址、amount=parseUnits(amt,6)、action.write=deposit、onSuccess 写 POST /api/deposit pending（+T5 signedPost）。
- T9 付账 Billing/BillDetail：spender=Payment、amount=amount+calculateFee(直读合约)、action.write=payBill(billId)、onSuccess 写 POST /api/bills/pay pending。
- 调用方负责：idle 前金额校验(≥10 USDT & ≤钱包余额)、第二笔成功写 pending+WalletAuth 签名(T5)、action 的 useTxState 由各页用 useContractDeposit/useContractPayBill 包装传入。
