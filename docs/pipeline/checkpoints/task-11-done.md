# Task 11 — T10 RegionDetail（申请弹层手续费+requiredDeposit 6 位+深蓝金）（web 3/3）

> 子项目 web(3/3) | 状态 DONE | 2026-06-10

### 产出文件
- packages/web/src/services/api/operatorApi.ts（requiredDeposit parseEther 18 位→parseUnits 6 位 USDT）
- packages/web/src/pages/RegionDetail.tsx（重写：运营商列表 ui/card+ui/badge+lucide(Wifi/Phone/Wallet)；申请弹层 FeeBreakdown 读链 calculateFee+押金本金+ShieldCheck 身份签名提示+拒签块；AmountDisplay 卡内 navy USDT 6 位；深蓝金）
- packages/web/src/pages/RegionDetail.test.tsx（新增 4 REG 测）

### git commit
8cab315 feat: web T10 RegionDetail(申请弹层手续费读链+requiredDeposit 6 位修+深蓝金)

### TDD
先红后绿：REG 测先写 → 实现后 95 passed。

### 测试结果
npx tsc --noEmit 0 error。npm run build ✓。npx vitest run：21 files / 95 passed（前序 91+新 4，无回归）。用例 REG-01(requiredDeposit 50_000_000→50.00 USDT 不二次缩放/不当 18 位)/REG-02(申请弹层 FeeBreakdown data-amount=押金本金读链 calculateFee+不消耗 gas 提示)/REG-03(ui/card 取代手写卡+无旧 token+lucide)/REG-04(拒签提示不进 pending)。grep RegionDetail 无旧 token/装饰 emoji；operatorApi 无 parseEther。主 Agent 已独立复跑确认 tsc 0+build ✓+95 测。

### code-simplifier
复用 FeeBreakdown/useFeeRate(T9)、ui/card·badge·AmountDisplay；申请流程无臆造链上付费 stepper。

### spec review
按 design v2 §3.6(手续费读链)/§3.7(身份签名)/§3.1 + arch-review B2 + T2 精度遗漏执行。requiredDeposit 6 位修、FeeBreakdown 读链 calculateFee 不自算、申请走 signedPost 意向无链上付费、拒签不进 pending。未越界（其余页 T11、接链基建只调用）。

### 设计还原
RegionDetail 对齐 design §3.6 手续费明细；深蓝金 ui/card+lucide+AmountDisplay 卡内 navy；国名 text-on-dark-gold。

### 复用检查
复用 FeeBreakdown/useFeeRate/useCalculateFee(T9)、ui/card·badge·input(T6)、AmountDisplay、useApplyNumber/signedPost(T4/T5)、parseUnits(6)。

### 设计稿对照
数值对照：requiredDeposit 50_000_000→50.00 USDT（REG-01，不二次缩放）✅；申请弹层 calculateFee 读链不自算（REG-02）✅；ui/card 取代手写卡（REG-03）✅；拒签不进 pending（REG-04）✅；95 测/tsc 0/build ✓ ✅；grep 清 ✅。

### 新增组件
无新增组件（复用 FeeBreakdown 等）；operatorApi 精度修。

### 新增色值
无（用 T6 深蓝金语义类）。

### ⚠️ 遗留（带入 T11/T12）
- T11 其余页换肤：Landing/Dashboard/Services(含 :48 搜索框 🔍)/Notifications/layout(AppLayout/Header/TabBar)/shared/wallet 的旧 token(brand-blue 等残留)+emoji→lucide 全量（权威表 theme-migration.md §2）。
- T11 清理：DepositInfo.currency "USDT"|"ETH" 删 ETH 类型；format.ts formatUSD dead util 删；useBilling 未用导出(isContractPending 等)清。
- 国旗 emoji 保留为国家标识数据（lucide 无逐国国旗）。
- T12 全量回归 + 31337 冒烟。
