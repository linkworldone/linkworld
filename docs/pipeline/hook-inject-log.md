# Hook 注入日志

> 生成时间: 2026-04-12T23:49:33.458Z
> 阶段: implement
> 本文件每次 session 启动时自动覆盖

⚠️ 以下规则在写代码时强制执行。

[PIPELINE ENFORCE — implement 阶段]

## 编码规范
- 遵循项目 .claude/CLAUDE.md 中的编码规范（颜色常量、导航方式、字符串国际化、按钮参数等项目特定规则）
- 颜色: 使用项目色值常量，禁止硬编码 hex
- icon: 从 Figma 导出 SVG/PNG，禁止框架内置 icon，禁止手绘 SVG
- 图片: PNG > 500KB 必须压缩

## 组件复用规则
- 先查项目已有组件目录（见 CLAUDE.md），有就直接用
- 没有就封装为独立 Widget/Component 到组件目录，禁止写进 page 文件
- 基础 UI 组件库作为底层积木用，不需要专门去找

## UI 还原硬标准
- 间距: 与设计稿 ±2px
- 颜色: hex 值完全匹配（查色值映射表）
- 字号/字重: 与设计稿完全匹配
- 圆角/阴影: 与设计稿一致

## UI Task 流程（主 Agent 执行）
1. 读 plan-review.json 当前页的组件/色值/资源（hook 已注入到上方）
2. 读 docs/design/*/context/{codePage}.design.md（React+Tailwind 结构）+ {codePage}.json（精确数值）
3. 将设计结构 + 截图 + 编码规范注入 subagent prompt
4. subagent 对照设计结构 + 精确数值写实际代码
5. subagent 返回后，主 Agent git commit
6. 对比设计结构中的 fontSize/padding/color 与实际代码值
7. 差异写入 checkpoint 的"设计稿对照"字段
8. 有 blocker 差异 → 派 subagent 修复
9. 有模拟器/浏览器 → 截图存证

## TDD（硬编码，不靠 AI 读 TDD.md）
- Entity/Controller 层: 强制先写测试再写实现（红→绿→重构）
- Widget 层: 推荐写 widget test
- code-simplifier: 新增 >50 行必跑，>20% 改动需人工确认

## Figma 规则
- 设计结构已在 figma 阶段下载，存在 context/{codePage}.design.md + {codePage}.json
- 直接读文件，subagent 按 techStack 转对应框架代码
- icon 导出: use_figma exportAsync 按完整组件节点
- SVG 导出后检查 viewBox（icon 通常 24x24，异常说明导了零件）
- 资源命名: ic_xxx.svg（snake_case），禁止 Figma 自动名

注意: checkpoint 文件已由 advance hook 批量创建。你只需要填写字段。
hook 会硬阻塞未填完的 checkpoint — 填完才能继续。

[阶段完成协议] 完成 implement 时: ①写 stages/implement.md ②更新 pipeline.json ③汇报

---
> [pipeline-init-check] 2026-04-12T23:49:33.485Z

📋 当前项目无架构地图（.claude/CLAUDE.md），建议让 Claude 扫描项目生成一份，后续对话可快速了解项目脉络。
📍 Pipeline: [需求探索 ✅] [设计分析 ✅] [三角色架构审查 ✅] [任务拆解 ✅] [→编码实现 ⬜] [测试 ⬜] [代码审查 ⬜]
📝 需求: LinkWorld 前端 Web 应用实现 —— 基于 0G Chain 的去中心化全球通信平台，React 19 + Vite 6 + TailwindCSS + shadcn/ui + wagmi/RainbowKit，Mock 数据驱动，移动端优先
Gate gate-2（）| 当前: 编码实现 | Skills: subagent-driven-development, test-driven-development, code-simplifier

---
> [orchestrator-guard] 2026-04-12T23:49:52.120Z

⛔ [ORCHESTRATOR] 主 Agent 不能直接写 /Users/bm/linkworld/packages/web/tailwind.config.ts。
请派 subagent 执行此操作。主 Agent 只能写 pipeline.json。

---
> [orchestrator-guard] 2026-04-12T23:50:01.331Z

⛔ [ORCHESTRATOR] 主 Agent 不能直接写 /Users/bm/linkworld/packages/web/src/index.css。
请派 subagent 执行此操作。主 Agent 只能写 pipeline.json。

---
> [orchestrator-guard] 2026-04-12T23:50:02.877Z

⛔ [ORCHESTRATOR] 主 Agent 不能直接写 /Users/bm/linkworld/packages/web/src/config/wagmi.ts。
请派 subagent 执行此操作。主 Agent 只能写 pipeline.json。

---
> [orchestrator-guard] 2026-04-12T23:50:05.223Z

⛔ [ORCHESTRATOR] 主 Agent 不能直接写 /Users/bm/linkworld/packages/web/src/config/constants.ts。
请派 subagent 执行此操作。主 Agent 只能写 pipeline.json。

---
> [orchestrator-guard] 2026-04-12T23:50:08.404Z

⛔ [ORCHESTRATOR] 主 Agent 不能直接写 /Users/bm/linkworld/packages/web/src/main.tsx。
请派 subagent 执行此操作。主 Agent 只能写 pipeline.json。

---
> [orchestrator-guard] 2026-04-12T23:50:33.241Z

⛔ [ORCHESTRATOR] 主 Agent 不能直接写 /Users/bm/linkworld/packages/web/tailwind.config.ts。
请派 subagent 执行此操作。主 Agent 只能写 pipeline.json。

---
> [subagent-inject] 2026-04-12T23:51:07.169Z | stage=implement | task=-

注入 skills: subagent-driven-development, test-driven-development
未找到: code-simplifier
总字符: 4169/15000

---
> [pipeline-context] 2026-04-12T23:51:22.512Z

📍 [需求探索 ✅] [设计分析 ✅] [三角色架构审查 ✅] [任务拆解 ✅] [→编码实现] [测试 ⬜] [代码审查 ⬜]
当前: 编码实现 | Skills: subagent-driven-development, test-driven-development, code-simplifier
需求: LinkWorld 前端 Web 应用实现 —— 基于 0G Chain 的去中心化全球通信平台，React 19 + Vite 6 + TailwindCSS + shadcn/ui + wagmi/RainbowKit，Mock 数据驱动，移动端优先

---
> [subagent-inject] 2026-04-12T23:54:16.610Z | stage=implement | task=-

注入 skills: subagent-driven-development, test-driven-development
未找到: code-simplifier
总字符: 4169/15000

---
> [subagent-inject] 2026-04-12T23:56:21.616Z | stage=implement | task=-

注入 skills: subagent-driven-development, test-driven-development
未找到: code-simplifier
总字符: 4169/15000

---
> [pipeline-context] 2026-04-12T23:56:38.241Z

📍 [需求探索 ✅] [设计分析 ✅] [三角色架构审查 ✅] [任务拆解 ✅] [→编码实现] [测试 ⬜] [代码审查 ⬜]
当前: 编码实现 | Skills: subagent-driven-development, test-driven-development, code-simplifier
需求: LinkWorld 前端 Web 应用实现 —— 基于 0G Chain 的去中心化全球通信平台，React 19 + Vite 6 + TailwindCSS + shadcn/ui + wagmi/RainbowKit，Mock 数据驱动，移动端优先

---
> [subagent-inject] 2026-04-12T23:57:04.321Z | stage=implement | task=-

注入 skills: subagent-driven-development, test-driven-development
未找到: code-simplifier
总字符: 4169/15000

---
> [checkpoint-guard] 2026-04-12T23:57:18.824Z

⛔ [BLOCKED] task-00-done.md 有未填字段: 产出文件, git commit, TDD, Figma 还原, 测试结果, code-simplifier, spec review, 复用检查, 设计稿对照
请先填完此 checkpoint 的以下字段再继续:
  - ### 产出文件
  - ### git commit
  - ### TDD
  - ### Figma 还原
  - ### 测试结果
  - ### code-simplifier
  - ### spec review
  - ### 复用检查
  - ### 设计稿对照

---
> [checkpoint-guard] 2026-04-12T23:57:21.355Z

⛔ [BLOCKED] task-00-done.md 有未填字段: 产出文件, git commit, TDD, Figma 还原, 测试结果, code-simplifier, spec review, 复用检查, 设计稿对照
请先填完此 checkpoint 的以下字段再继续:
  - ### 产出文件
  - ### git commit
  - ### TDD
  - ### Figma 还原
  - ### 测试结果
  - ### code-simplifier
  - ### spec review
  - ### 复用检查
  - ### 设计稿对照

---
> [orchestrator-guard] 2026-04-12T23:58:35.157Z

⛔ [ORCHESTRATOR] 主 Agent 不能直接写 EOF。
请派 subagent 执行此操作。主 Agent 只能写 pipeline.json。

---
> [subagent-inject] 2026-04-12T23:58:44.369Z | stage=implement | task=-

注入 skills: subagent-driven-development, test-driven-development
未找到: code-simplifier
总字符: 4169/15000
