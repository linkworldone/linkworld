# Stage: requirement — 需求探索

> **状态**: completed | **日期**: 2026-04-12 | **Gate**: 1

## 产出摘要

通过交互式需求探索，从产品文档和前端设计规格出发，明确了 Round 1 的完整交付范围、技术决策和 Mock 策略。产出完整设计规格文档。

## 关键决策

| # | 决策 | 理由 |
|---|------|------|
| 1 | Round 1 交付全部 6 个页面 | 用户确认完整交付，非 MVP 裁剪 |
| 2 | shadcn/ui 按需引入 | 只用 Sheet/Button/Badge/Tabs 等基础组件，卡片和布局手写 Tailwind 保持品牌一致性 |
| 3 | React Query + wagmi hooks，无额外 Context | 零冗余，wagmi 自带钱包状态，React Query 管理服务端缓存 |
| 4 | AppLayout 层路由守卫 | 简单直观，引导卡片直接在 Layout 层渲染 |
| 5 | 真实 RainbowKit + Mock 业务逻辑 | 钱包连接真实体验，业务数据走 Mock Service |

## 产出文件

| 文件 | 内容 |
|------|------|
| docs/superpowers/specs/2026-04-12-linkworld-frontend-design.md | 完整前端实现设计规格（10 sections） |
| .superpowers/brainstorm/content/mobile-landing-dashboard.html | Landing + Dashboard 移动端 mockup |
| .superpowers/brainstorm/content/mobile-inner-pages.html | Deposit/Services/Billing/Notifications 移动端 mockup |

## 用户确认的事项

- 全部 6 页面交付范围
- shadcn/ui 按需引入策略
- React Query + wagmi hooks 状态管理方案
- AppLayout 层守卫方案
- 真实 RainbowKit + Mock 注册策略
- 全部 6 个页面的移动端布局 wireframe
- 完整设计规格文档
