# Checkpoint: Task 6 — Deposit Page

> **完成时间**: 2026-04-13

### 产出文件
- src/pages/Deposit.tsx

### git commit
pending

### TDD
N/A — UI page, verified via tsc

### Figma 还原
N/A — follows mobile mockup wireframe

### 测试结果
npx tsc --noEmit: PASS

### code-simplifier
Single page, ~120 lines with inline Drawer

### spec review
- Balance card with progress bar matches spec 5.3
- Deposit(green)/Withdraw(outline) dual buttons matches spec
- History list with colored +/- matches spec
- Currency selector (USDT/ETH) in deposit sheet
- Uses vaul Drawer (arch-review risk #8)

### 复用检查
Uses AmountDisplay from Task 3, Button from shadcn

### 设计稿对照
- Balance card: p-6 rounded-2xl, label text-[10px] uppercase, balance size="lg" (text-4xl), progress bar h-1
- Action buttons: flex gap-2.5, py-3, deposit bg-status-success, withdraw variant="outline"
- History items: p-3 bg-surface-card rounded-xl, type text-xs, date text-[10px], amount text-sm font-bold
- Drawer: p-6, handle w-12 h-1, title text-lg font-bold, input px-4 py-3 rounded-xl text-sm
- Currency selector: py-2 rounded-lg text-sm font-semibold, active border-brand-blue

### 偏差记录
无偏差
