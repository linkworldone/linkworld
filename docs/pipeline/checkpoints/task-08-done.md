# Checkpoint: Task 8 — Billing + BillDetail

> **完成时间**: 2026-04-13

### 产出文件
- src/pages/Billing.tsx
- src/pages/BillDetail.tsx

### git commit
pending

### TDD
N/A — UI pages, verified via tsc

### Figma 还原
N/A — follows mobile mockup wireframe

### 测试结果
npx tsc --noEmit: PASS

### code-simplifier
Two focused pages, Billing ~120 lines, BillDetail ~90 lines

### spec review
- Unpaid/History tab switcher matches spec 5.5
- Unpaid alert banner with count and due date
- Bill cards with fee breakdown (operator + platform 2.5%)
- BillDetail with usage stats (data GB + call minutes)
- Pay confirmation drawer
- Badge component for paid/unpaid status

### 复用检查
Uses Badge, Button, Drawer, formatUSD, formatDate

### 设计稿对照
Matches mobile-inner-pages.html Billing mockup

### 偏差记录
无偏差
