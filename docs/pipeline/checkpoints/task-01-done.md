# Task 01 — T0 绿基线 + vitest 测试设施（web 3/3）

> 子项目 web(3/3) | 状态 DONE | 2026-06-10

### 产出文件
- packages/web/src/components/wallet/RegisterSheet.tsx（删第 22 行 useRegister() 未用解构 isSuccess，清 TS6133）
- packages/web/package.json（devDeps: vitest/@testing-library/react/@testing-library/jest-dom/jsdom；scripts: test=vitest run、test:watch）
- packages/web/vitest.config.ts（新建：jsdom+globals+setupFiles+@ alias，不动 vite.config.ts 保护 build）
- packages/web/vitest.setup.ts（新建：jest-dom/vitest）
- packages/web/src/lib/utils.test.ts（新建：cn() 3 smoke 断言）
- pnpm-lock.yaml（仅 vitest 相关依赖新增）

### git commit
3ee2397 chore: web T0 清 tsc 基线 + 搭 vitest 测试设施

### TDD
T0 搭设施非业务行为；smoke 测试(cn() 3 断言)证明设施可跑。业务单测/组件测 T2-T12 各自 TDD 写。

### 测试结果
npx tsc --noEmit 0 error（全量仅 RegisterSheet:22 一处，已清，无其他 pre-existing 错）。npm test：1 file / 3 tests passed。npm run build：✓ built（仅 chunk>500kB 信息级 warning 非 error）。主 Agent 已确认 build 绿。

### code-simplifier
删未用解构最干净（无 eslint 抑制）；测试配置独立 vitest.config.ts 与 build 解耦。

### spec review
按 design v2 §10 + arch-review B3/B4 执行：清 RegisterSheet TS6133 拿绿基线(B3)、搭 vitest 设施(B4)。未越界（链配置/精度/换肤/接链全留 T1+）。关键决策：项目是 pnpm workspace（根 pnpm-workspace.yaml+pnpm-lock.yaml），web 的 package-lock.json 是陈旧遗留，用 pnpm --filter 装依赖（用 npm 会破坏 workspace）。

### 设计还原
T0 为基线设施，无 UI 还原。B3 绿基线 + B4 测试设施就绪，供 T1-T12 TDD。

### 复用检查
复用现有 vite/tsconfig/@ alias；vitest 配置独立不污染 build；smoke 测复用现有 cn()。

### 设计稿对照
数值对照：tsc error 1→0 ✅；测试 0→3 passed ✅；build 绿 ✅；新增 devDeps 4 个（vitest/testing-library×2/jsdom）✅；新建配置文件 2 + smoke 测 1。

### 新增组件
无新增业务组件。新增测试设施（vitest.config/setup + smoke 测）。

### 新增色值
无（T0 基线，换肤留 T6+）。

### ⚠️ 遗留
- 陈旧 packages/web/package-lock.json 建议后续顺手删（本轮严禁碰，未动）。
- 测试设施就绪：T2-T12 直接 import vitest + @testing-library/react；jsdom 不渲染真实样式，视觉回归走 build/browser 不在单测。
