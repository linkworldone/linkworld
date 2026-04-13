# Stage: plan — 任务拆解

> **状态**: completed | **日期**: 2026-04-13 | **Gate**: 2

## 产出摘要

基于设计规格和 arch-review 审查结果，将前端实现拆解为 11 个 Task（Task 0-10），覆盖从基础设施配置到最终集成测试的完整流程。所有 10 个 arch-review 风险项已分配到具体 Task。

## 关键决策

| # | 决策 | 理由 |
|---|------|------|
| 1 | Task 0 前置基础设施（shadcn/ui + Tailwind tokens + Provider tree） | arch-review 要求先配置再写页面，避免返工 |
| 2 | Bill 金额用 string 不用 number | arch-review 风险项：浮点精度问题 |
| 3 | 使用 vaul 替代 shadcn/ui Sheet | arch-review 风险项：底部弹出 + 拖拽手势需求 |
| 4 | Mock Service 统一加延迟 | arch-review 风险项：无延迟则 loading 状态无法测试 |
| 5 | getUserProfile 返回 User|null | arch-review 阻塞项修复 |
| 6 | VirtualNumber 增加 credentials 字段 | 支持 eSIM/VoIP 凭证展示 |

## Task 列表

| Task | 名称 | 文件数 | 依赖 |
|------|------|--------|------|
| 0 | Infrastructure Setup | 8 | — |
| 1 | Types, Mock Data & Services | 10 | Task 0 |
| 2 | React Query Hooks | 5 | Task 1 |
| 3 | Shared Components + AppLayout | 7 | Task 2 |
| 4 | Wallet + Landing + Router | 5+7 placeholders | Task 3 |
| 5 | Dashboard Page | 1 | Task 4 |
| 6 | Deposit Page | 1 | Task 4 |
| 7 | Services + RegionDetail | 2 | Task 4 |
| 8 | Billing + BillDetail | 2 | Task 4 |
| 9 | Notifications Page | 1 | Task 4 |
| 10 | Final Integration | 0 (verify only) | Task 5-9 |

## 产出文件

| 文件 | 内容 |
|------|------|
| docs/superpowers/plans/2026-04-13-linkworld-frontend.md | 完整实现计划（194 行） |
| docs/pipeline/stages/plan.md | 本文件 |

## 用户确认的事项

- 用户选择 Subagent-Driven 执行模式
