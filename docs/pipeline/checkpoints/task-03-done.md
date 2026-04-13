# Checkpoint: Task 3 — Shared Components + AppLayout

> **完成时间**: 2026-04-13

### 产出文件
- src/components/shared/StatusBadge.tsx
- src/components/shared/AmountDisplay.tsx
- src/components/shared/GuardCard.tsx
- src/components/shared/EmptyState.tsx
- src/components/layout/TabBar.tsx
- src/components/layout/Header.tsx
- src/components/layout/AppLayout.tsx

### git commit
635c2bc

### TDD
N/A — UI components, verified via tsc

### Figma 还原
N/A — no Figma, follows design spec mockups

### 测试结果
npx tsc --noEmit: PASS (0 errors)

### code-simplifier
Components are focused and minimal, each under 80 lines

### spec review
- AppLayout: max-w-mobile (430px) addresses arch-review risk #7
- TabBar: 5 tabs with badges matching spec 6.2
- Header: Dashboard variant + generic variant matching spec 6.3
- GuardCard: INACTIVE/SUSPENDED guards matching spec 4.3

### 复用检查
N/A — all new files

### 设计稿对照
- AppLayout: max-w-mobile = 430px (spec: 430px), min-h-screen
- TabBar: fixed bottom-0, py-2, 5 tabs, min-w-[44px] min-h-[44px] touch targets (spec: 44px), badge size w-3.5 h-3.5 text-[8px]
- Header dashboard: px-4 py-3, address text-[15px] font-semibold, muted text-[11px], avatar w-8 h-8 rounded-full
- Header generic: px-4 py-3, title text-[17px] font-bold
- TabBar label: text-[9px], icon text-xl, active color brand-blue
- Main content: pb-[80px] for TabBar clearance

### 偏差记录
无偏差
