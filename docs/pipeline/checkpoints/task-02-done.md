# Checkpoint: Task 2 — React Query Hooks

> **完成时间**: 2026-04-13

### 产出文件
- src/hooks/useUser.ts
- src/hooks/useDeposit.ts
- src/hooks/useBilling.ts
- src/hooks/useOperator.ts
- src/hooks/useNotification.ts

### git commit
35079e5

### TDD
N/A — hooks are thin wrappers over services, type-checked via tsc

### Figma 还原
N/A

### 测试结果
npx tsc --noEmit: PASS (0 errors)

### code-simplifier
N/A — hooks follow uniform pattern

### spec review
All hooks match spec Section 8 service interfaces. useUnreadCount has 30s refetchInterval. All mutations invalidate related queries.

### 复用检查
N/A — all new files

### 设计稿对照
N/A — 纯数据层，无视觉输出。5 个 hooks 对应 spec Section 8 的 5 个 service 接口，refetchInterval 30000ms 匹配 spec 要求。

### 偏差记录
无偏差
