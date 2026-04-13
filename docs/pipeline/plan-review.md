# LinkWorld 前端实现计划 — Plan Review

> **日期**: 2026-04-13 | **Round**: 1

## 概述

将设计规格拆解为 11 个 Task（Task 0-10），从基础设施配置到最终集成测试。

## 页面→Task 映射

| 代码页 | Task | 组件依赖 | 色值依赖 |
|--------|------|---------|---------|
| Infrastructure | Task 0 | shadcn/ui (Button, Badge, Tabs) | Tailwind semantic tokens |
| Types + Services | Task 1 | — | — |
| Hooks | Task 2 | — | — |
| Layout + Shared | Task 3 | StatusBadge, AmountDisplay, GuardCard, EmptyState, Header, TabBar, AppLayout | surface, brand, status, text |
| Landing | Task 4 | ConnectButton, RegisterSheet | brand-blue, brand-purple |
| Dashboard | Task 5 | StatusBadge, AmountDisplay | surface-gradient, brand-blue, brand-purple, status-success/warning/danger |
| Deposit | Task 6 | AmountDisplay, Button | status-success, surface-gradient |
| Services + RegionDetail | Task 7 | EmptyState, Button | brand-blue, surface-card |
| Billing + BillDetail | Task 8 | Badge, Button | status-danger, brand-blue |
| Notifications | Task 9 | EmptyState | brand-blue, status-success, status-danger, text-muted |
| Integration | Task 10 | — | — |

## 组件增量比对

| 组件 | 状态 | 来源 |
|------|------|------|
| Button | 复用 | shadcn/ui (Task 0 安装) |
| Badge | 复用 | shadcn/ui (Task 0 安装) |
| Tabs | 复用 | shadcn/ui (Task 0 安装) |
| Drawer | 复用 | vaul (Task 0 安装) |
| StatusBadge | 新建 | Task 3 |
| AmountDisplay | 新建 | Task 3 |
| GuardCard | 新建 | Task 3 |
| EmptyState | 新建 | Task 3 |
| Header | 新建 | Task 3 |
| TabBar | 新建 | Task 3 |
| AppLayout | 新建 | Task 3 |
| ConnectButton | 新建 | Task 4 |
| RegisterSheet | 新建 | Task 4 |

## 色值增量比对

| 色值 | Tailwind Token | 来源 |
|------|---------------|------|
| #0a0a14 | surface-DEFAULT | Task 0 tailwind.config.ts |
| #0f0f1a | surface-card | Task 0 |
| #1a1a2e | surface-secondary / border | Task 0 |
| #3b82f6 | brand-blue | Task 0 |
| #8b5cf6 | brand-purple | Task 0 |
| #06b6d4 | brand-cyan | Task 0 |
| #22c55e | status-success | Task 0 |
| #f59e0b | status-warning | Task 0 |
| #ef4444 | status-danger | Task 0 |
| #e0e0e0 | text-primary | Task 0 |
| #888888 | text-secondary | Task 0 |
| #555555 | text-muted | Task 0 |

## Arch-Review 风险项处理

| # | 风险 | 处理 Task |
|---|------|----------|
| 1 | Tailwind 配色 semantic tokens | Task 0 ✅ |
| 2 | BrowserRouter 位置 | Task 0 ✅ |
| 3 | shadcn/ui 初始化前置 | Task 0 ✅ |
| 4 | 桌面端 max-w 约束 | Task 3 ✅ |
| 5 | vaul 替代 Sheet | Task 0 ✅ |
| 6 | Mock 延迟模拟 | Task 1 ✅ |
| 7 | Bill 金额精度 | Task 1 ✅ |
| 8 | 多虚拟号码 | Task 7 |
| 9 | 提前结账入口 | Task 8 |
| 10 | 逾期扣款警告 | Task 5 ✅ |
