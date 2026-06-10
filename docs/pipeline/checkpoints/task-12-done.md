# Task 12 — T11 其余页换肤 sweep + emoji→lucide 全量 + 旧 token 清零（web 3/3）

> 子项目 web(3/3) | 状态 DONE | 2026-06-10

### 产出文件
- packages/web/src/components/layout/{TabBar,Header,AppLayout}.tsx（5 项导航 emoji→lucide+44px、Bell/ChevronLeft、封禁 ShieldAlert、头像 navy/royal/gold）
- packages/web/src/pages/{Landing,Dashboard,Services,Notifications,Billing,BillDetail}.tsx（旧 token 清零 + emoji→lucide + 深蓝金 ui/card + AmountDisplay 卡内 navy/深底金；Landing 0G→Arbitrum；Services 🔍→Search）
- packages/web/src/components/{shared/{FeeBreakdown,BottomSheet,TwoStepAction},wallet/{ConnectButton,RegisterSheet}}.tsx（残留 brand-blue/surface-secondary→深蓝金语义类）
- packages/web/src/types/index.ts（DepositInfo.currency 仅 USDT 删 ETH）
- packages/web/src/utils/format.ts（删 dead formatUSD）
- packages/web/src/hooks/useBilling.ts（删 usePayBill 未用导出）
- packages/web/src/pages/skin.test.tsx（新增 SKIN-01/02/03）

### git commit
bab9658 feat: web T11 其余页换肤 sweep + emoji→lucide 全量 + 旧 token 清零 + 死代码清理

### TDD
先红后绿：SKIN 测先写 → 实现后 101 passed。

### 测试结果
npx tsc --noEmit 0 error。npm run build ✓。npx vitest run：22 files / 101 passed（前序 95+新 6，无回归）。**全量 grep 清零**：#3b82f6/#0a0a14/#8b5cf6/brand-blue/surface-gradient/font-orbitron/装饰 emoji/formatUSD/Landing "0G" 全 CLEAN（国旗 FLAG_MAP emoji 保留为数据、chains.ts nativeCurrency ETH 合法）。主 Agent 已独立复跑确认 tsc 0+build ✓+101 测。

### code-simplifier
全量旧 token→深蓝金语义类单一出口；删 formatUSD/usePayBill 死代码；emoji→lucide 统一。

### spec review
按 DESIGN.md 深蓝金/金色铁律/lucide 全量 + theme-migration §2 权威表 + arch-review B2/44px 触摸执行。卡内金额 navy、深底金、金仅激活/CTA/品牌/卡金线、44px 触摸、pending 文案统一。未越界（接链逻辑 T1-T10 仅换 class/图标/文案）。

### 设计还原
9 页 + layout/shared/wallet 100% 深蓝金覆盖；emoji→lucide 全量；色值单一出口；grep 旧色值/emoji 清零达 guardRule。

### 复用检查
复用 T6 深蓝金 token/ui-card·input、AmountDisplay 分流、lucide；无新组件，纯换肤+清理。

### 设计稿对照
数值对照：旧色值 grep 0 残留（guardRule 达标）✅；emoji→lucide 全量（国旗数据除外）✅；TabBar 5 项 lucide+44px ✅；DepositInfo.currency 仅 USDT ✅；101 测/tsc 0/build ✓ ✅；Landing 0G→Arbitrum ✅。

### 新增组件
无新增（纯换肤）；删 formatUSD/usePayBill 未用导出。

### 新增色值
无（全用 T6 深蓝金语义类，本任务是把旧 token 替换为它们）。

### ⚠️ 遗留 / 重要结构事项
- **嵌套 packages/web/.git 空壳仓库（无 remote）**：T11 subagent 误提交进它，已由主 Agent 把 T11 改动归并入外层 linkworld 仓库（bab9658，与 T0-T10 一致）。该嵌套 .git 是遗留空壳，待主 Agent 与用户确认后清理（rm packages/web/.git）以防后续混淆。
- 国旗 emoji（operatorApi FLAG_MAP / Dashboard regionFlags）保留为国家标识数据。
- T12：全量测试 + 构建 + 本地 31337 链路冒烟（Arbitrum 421614 端到端阻塞合约上链，web DONE=本地 31337 绿）。
