# Task 08 — T7 Deposit 页（两步态+锁仓倒计时+pending+深蓝金）（web 3/3）

> 子项目 web(3/3) | 状态 DONE | 2026-06-10

### 产出文件
- packages/web/src/pages/Deposit.tsx（重写：TwoStepAction 两步充值 + LockCountdown + pending 三态 + 深蓝金换肤）
- packages/web/src/components/shared/LockCountdown.tsx（新建：锁仓倒计时/解锁三态 + 纯函数 lockState/formatRemaining）
- packages/web/src/hooks/contracts/useDepositContract.ts（追加 useLockExpiry/useUsdtBalance）+ index.ts 导出
- packages/web/src/hooks/useDeposit.ts（recordIntent 增强：拒签 WalletAuthRejectedError 上抛、瞬时失败退避+pendingSync）
- packages/web/src/components/shared/LockCountdown.test.tsx + pages/Deposit.test.tsx（新增）

### git commit
cefdb75 feat: web T7 Deposit 页(TwoStepAction 两步充值+LockCountdown 倒计时+pending+深蓝金)

### TDD
先红后绿：DEP 测先写 → 实现后 71 passed。

### 测试结果
npx tsc --noEmit 0 error。npm run build ✓。npx vitest run：16 files / 71 passed（前序 54+新 17，无回归）。用例 DEP-01(LockCountdown 边界 now>=expiry locked/unlocked)/DEP-02(文案无"利息")/DEP-03(<10 USDT 拦截/>钱包余额拦截)+锁仓禁用提现+无 ETH 渲染。grep Deposit.tsx 无 brand-blue/ETH/emoji/sonner。主 Agent 已独立复跑确认 tsc 0+build ✓+71 测。

### code-simplifier
lockState/formatRemaining 抽纯函数易测；TwoStepAction 复用（T3）；ui/card·input·AmountDisplay·TxStatusBadge 复用（T6）。

### spec review
按 design v2 §3.1/§3.2/§3.4 + DESIGN.md 金色铁律 + arch-review B2/pending 兜底执行。锁仓边界 now>=expiry（remaining<=0 非 <）、去"利息"、累加顺延提示、pending 不染绿不据 200、Arbiscan 逃生、拒签不进 pending。未越界（其他页 T8-T11、接链基建 T1-T5 只调用）。

### 设计还原
Deposit 页对齐 design §3.2 两步态/§3.4 锁仓/§3.1 pending + 深蓝金（ui/card 暖米白、AmountDisplay 卡内 navy、lucide 图标 ArrowDownToLine/Lock/ExternalLink、USDT 币种）。

### 复用检查
复用 TwoStepAction(T3)/useDeposit pending(T4)/signedPost(T5)/ui-card·input·AmountDisplay·TxStatusBadge(T6)/getLockExpiry·getDepositAmount 读链；新增 LockCountdown + useLockExpiry/useUsdtBalance。

### 设计稿对照
数值对照：MIN_DEPOSIT 10 USDT 校验 ✅；锁仓边界 now>=expiry（DEP-01）✅；去"利息"（DEP-02）✅；pending「约 1-2 分钟」不染绿 ✅；卡内金额 navy ✅；71 测/tsc 0/build ✓ ✅；grep 无 brand-blue/ETH/emoji ✅。

### 新增组件
新增 LockCountdown（+lockState/formatRemaining）+ useLockExpiry/useUsdtBalance。

### 新增色值
无（用 T6 深蓝金语义类）。

### ⚠️ 遗留（带入 T8-T11）
- T9 付账 useBilling.recordIntent 仍裸 retryWithBackoff（会把拒签当瞬时失败落 pendingSync）→ 对齐 T7 useDeposit 的拒签上抛范式。
- EmptyState icon:string 未改 LucideIcon（DESIGN.md D5）→ T8 Cards/其余页统一处理（T7 历史空态规避未用 EmptyState）。
- DepositInfo.currency 类型仍含 "USDT"|"ETH"，页面已固定 USDT，彻底删 ETH 类型留 T11 全局清理。
- BottomSheet 继续共享复用（BillDetail/RegionDetail/Billing 在用），未动。
