# Task 09 — T8 Cards 页（双Tab 去Admin+NFT读链+SIM领取+深蓝金）（web 3/3）

> 子项目 web(3/3) | 状态 DONE | 2026-06-10

### 产出文件
- packages/web/src/pages/Cards.tsx（重写：ui/tabs 双 Tab；Tab1 流量卡去 Admin 发卡按钮+Sparkles/Info 自动发放说明+不可转卖+销毁后 30 天+isError 区分加载失败/真无卡+NFT 读链；Tab2 SIM 领取 Nfc/MailPlus 收件表单→pendingSync+「全球通 eSIM 即将推出」；ui/card+AmountDisplay 卡内 navy+lucide）
- packages/web/src/components/shared/EmptyState.tsx（icon string→LucideIcon，DESIGN.md D5）
- packages/web/src/components/shared/GuardCard.tsx（icon→LucideIcon）
- packages/web/src/components/layout/AppLayout.tsx（GuardCard 调用 emoji→Wallet/AlertTriangle）
- packages/web/src/pages/{Services,Notifications}.tsx（EmptyState 调用 emoji→Smartphone/Bell）
- packages/web/src/pages/Cards.test.tsx（新增 7 用例）

### git commit
e7a438a feat: web T8 Cards 页(双Tab 去Admin发卡+NFT读链 error态+SIM领取 pendingSync+深蓝金)

### TDD
先红后绿：CARD 测先写 → 实现后 78 passed。

### 测试结果
npx tsc --noEmit 0 error。npm run build ✓。npx vitest run：17 files / 78 passed（前序 71+新 7，无回归）。用例 CARD-01(无 Admin 发卡+自动发放说明)/CARD-02(isError 加载失败 vs 真无卡空态)/CARD-03(SIM 表单→pendingSync+成功提示)/CARD-04(双 Tab 切换)。grep Cards.tsx 无旧 token/emoji/Admin 发卡。主 Agent 已独立复跑确认 tsc 0+build ✓+78 测。

### code-simplifier
双 Tab 用 ui/tabs；isError/空态区分清晰；EmptyState/GuardCard icon 改 LucideIcon 收口类型。

### spec review
按 design v2 §3.5 + DESIGN.md + arch-review(getLogs error 态)执行：移除 onlyOracle Admin 发卡按钮、自动发放说明、isError 禁静默空、SIM pendingSync 降级。未越界（其他页 T9-T11、接链基建 T1-T5 只调用）。

### 设计还原
Cards 双 Tab 对齐 design §3.5；深蓝金 ui/card+lucide(CreditCard/Sparkles/Nfc/MailPlus/AlertCircle)+AmountDisplay 卡内 navy。

### 复用检查
复用 ui/tabs·card·input(T6)、useTrafficCards isError(T4)、pendingSync 工具、AmountDisplay；EmptyState/GuardCard icon 类型升级 LucideIcon。

### 设计稿对照
数值对照：双 Tab(NFT/SIM)✅；移除 Admin 发卡(CARD-01)✅；isError 区分(CARD-02)✅；SIM pendingSync(CARD-03)✅；销毁后 30 天文案 ✅；78 测/tsc 0/build ✓ ✅；grep 清 ✅。

### 新增组件
Cards 双 Tab 重写；EmptyState/GuardCard icon→LucideIcon。

### 新增色值
无（用 T6 深蓝金语义类）。

### ⚠️ 遗留（带入 T9-T11）
- EmptyState/GuardCard icon 已升 LucideIcon，4 调用方(AppLayout×2/Services/Notifications)已同步，T11 换这些页无需再改类型只上色。
- Services.tsx:48 搜索框还有 🔍（EmptyState 之外）→ T11 Services 完整换肤处理。
- SIM 领取为降级 R9（仅 localStorage pendingSync），后端 SIM 端点就绪后补真实提交+pending 重试。
- Cards 移除旧 Available Credit 聚合额度卡（旧模型 design §3.5 未要求）。
