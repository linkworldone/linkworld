# Checkpoint: Task 0 — Infrastructure Setup

> **完成时间**: 2026-04-13

### 产出文件
- packages/web/tailwind.config.ts (semantic color tokens)
- packages/web/src/index.css (dark theme CSS variables)
- packages/web/src/config/wagmi.ts (wagmi + RainbowKit config)
- packages/web/src/config/constants.ts (business constants)
- packages/web/src/main.tsx (Provider tree restructure)
- packages/web/src/lib/utils.ts (shadcn cn helper)
- packages/web/components.json (shadcn config)
- packages/web/src/components/ui/button.tsx (shadcn)
- packages/web/src/components/ui/badge.tsx (shadcn)
- packages/web/src/components/ui/tabs.tsx (shadcn)

### git commit
pending — will batch commit with Task 1

### TDD
N/A — infrastructure setup, no business logic

### Figma 还原
N/A — no Figma for this project

### 测试结果
pnpm build: PASS (4515 modules, no errors)

### code-simplifier
N/A — infra setup only

### spec review
Provider 嵌套顺序符合 wagmi 2.x 规范。shadcn v4 与 Tailwind v3 兼容性已验证。

### 复用检查
N/A — all new files

### 设计稿对照
Tailwind tokens 与设计规格 9.1 配色表完全对应。

### 偏差记录
shadcn 安装了 v4 版本（oklch 颜色格式），已适配。
