# Hook 注入日志

> 生成时间: 06-08 14:02
> 阶段: design
> 本文件每次 session 启动时自动覆盖

[CLAWPIPE — 全 pipeline 通用规则]

## ⚡ CodeGraph 优先（如可用）

如果项目根目录有 `.codegraph/`（codegraph CLI 已索引），**优先用 mcp__codegraph__* MCP tool** 而非 grep / Read / Glob 扫文件：

| 想做的事 | 用 |
|---|---|
| 找符号位置 | `codegraph_search` |
| 理解某区域 / "X 怎么工作" | `codegraph_explore` |
| 看调用链 | `codegraph_callers` / `codegraph_callees` |
| 改前看影响范围 | `codegraph_impact` |
| 单符号完整源（含 overload） | `codegraph_node` |
| 列文件 / 项目结构 | `codegraph_files` |
| 看索引状态 / pending sync | `codegraph_status` |

⛔ 调 codegraph_* 工具的返回结果**视为已读**，**不要再 grep / Read 验证**（除非看到 ⚠️ pending sync banner，那时直接 Read 实时文件）。

若无 `.codegraph/` 目录 → 此段跳过，照旧用 grep / Read。
若 codegraph CLI 未装 → 强烈建议装：
  `curl -fsSL https://raw.githubusercontent.com/colbymchenry/codegraph/main/install.sh | sh && codegraph install`
  装完在 scan 阶段会自动 `codegraph init -i` 建索引。
⚠️ 以下规则在 design 阶段强制执行。

[PIPELINE ENFORCE — design 阶段]

## 模式：按 figma 状态走双路径

## 路径判断

```bash
node -e 'try{const p=JSON.parse(require("fs").readFileSync(".claude/pipeline.json","utf8"));console.log(p.figmaFileKey?"HAS_FIGMA":"NO_FIGMA")}catch(e){console.log("NO_FIGMA")}'
```

- `HAS_FIGMA` → 路径 A（主 Agent 自己整理 figma 产出）
- `NO_FIGMA` → 路径 B（派 spawn_task 跟用户对话出 DESIGN.md + design.md）

---

## 路径 A: 有 figma（主 Agent 自己整理）

figma 阶段已定视觉方向，design 阶段只是**整理 + 加实现层细节**。不要派 spawn_task。

### 必做
1. Read figma 产出（`docs/pipeline/context/*.design.md` + `screenshots/`）
2. Read 项目基线（`docs/design/<项目>/components.md` 等）
3. Read 关键源码（至少 3 个，如 `src/index.css` / `tailwind.config` / 关键 hook / ABI），grep build 产物验证假设
4. 写 `docs/pipeline/stages/design.md`（含状态机 / API 映射 / 组件复用 / 头部"已核对真实代码"清单）
5. git commit + advance 推进 arch-review

---

## 路径 B: 无 figma（派 spawn_task）

### 流程

1. 调 `mcp__ccd_session__spawn_task`，prompt 原样复制下方 ~~~ 之间的模板
2. 派完停手，告诉用户「点 chip 进入新 session 对话，完事回主 session 说继续」
3. 用户说"继续"时回流：Read 产物 + advance 推进

### ⛔ 红线

- 不要 Read requirement.md / scan / 真实源码（让 spawned 自己读）
- 不要提炼"设计种子"塞 prompt
- 模板原样复制，不二次创作

### spawn_task PROMPT 模板（原样复制 ~~~ 之间）

~~~
你是 design 阶段执行者。读 requirement + 真实代码 + 跟用户对话出设计系统 + 设计分析。

【开聊前必读】
- docs/pipeline/stages/requirement.md（上一阶段 PRD）
- docs/pipeline/stages/scan.md
- docs/design/<项目>/{components,color-mapping,utils}.md
- 至少 3 个关键源码（src/index.css / tailwind.config.ts / 关键 hook / ABI），grep build 产物验证假设

【加载顺序：先约束 → 后发散】
1. 加载 design-consultation skill（角色：约束 — 建立 DESIGN.md 作为后续探索的设计系统地基；写到 docs/design/<项目>/DESIGN.md）
2. design-consultation 完成后，加载 design-shotgun skill（角色：发散 — 基于 DESIGN.md 视觉变体探索，用户挑方向）

汇总：把 DESIGN.md + 选定方向 + 真实代码核对结果写到 docs/pipeline/stages/design.md。

⚠️ 任务信封 vs HOW 分层（v1.6.4 含冲突裁决）：

**省略轴**（模板没提）：模板没提某步骤 → **按 skill checklist 走**，不要因为本模板没列就跳过。

**冲突轴**（直接打架）：skill 自带指令跟模板信封正面冲突时，**模板赢**：

- skill 说"终点跳 X skill" vs 模板说"接 Y skill" → 以**模板**为准
- skill 说"产物写到路径 A" vs 模板说"写到路径 B" → 以**模板**为准

具体例子（design 阶段必踩冲突）：
- design-consultation 默认查/写**项目根** `DESIGN.md`（`ls DESIGN.md design-system.md`）→ 本场景**覆盖**为「写 `docs/design/<项目>/DESIGN.md`」
- design-shotgun 默认写 `~/.gstack/` 全局目录 → 本场景 mockup HTML 工件**可留在 `~/.gstack/`**（临时工件 OK），但最终设计决策必须汇总到 `docs/pipeline/stages/design.md`

裁决原则总结：
- **HOW**（每步怎么做 / 视觉变体生成方式 / 与用户对话节奏）→ **skill 权威**
- **顺序 / 产物路径 / 下一个接哪个 skill** → **模板权威**

【产出】
- docs/pipeline/stages/design.md
- docs/design/<项目>/DESIGN.md

git add docs/pipeline/stages/design.md docs/design/<项目>/DESIGN.md
git commit -m "feat: design 阶段产出（DESIGN.md + 设计分析）"

【限制】
- 不写业务代码（只产设计分析）
- **跟用户对话自己来**，不要把对话委派给 subagent（避免用户感觉被踢皮球）
- 但 `$D generate` / mockup 视觉生成等**机械并行工作**按 skill 默认机制走 — design-shotgun 用 Task 派 N 个 fresh Agent 并行生成 variant 是核心多样性来源，**允许**
- commit 后停下，告诉用户回主 session 说"继续"
~~~

## 完成判据
`docs/pipeline/stages/design.md` git tracked + commit → advance 自动推进到 arch-review。

[阶段完成协议] 完成 design 时: ①写 stages/design.md ②更新 pipeline.json ③汇报

---
> [pipeline-init-check] 06-08 14:03

📋 当前项目无架构地图（.claude/CLAUDE.md），建议让 Claude 扫描项目生成一份，后续对话可快速了解项目脉络。
📍 Pipeline: [项目基线扫描 ⬜] [需求探索与明确 ⬜] [→设计分析 🔄] [三角色架构审查 ⬜] [任务拆解与增量比对 ⬜] [编码实现 ⬜] [测试 ⬜] [代码审查 ⬜] [发布 PR ⬜] [复盘 ⬜]
📝 需求: Link World web 端重构 — 应用新深蓝金配色方案，重构钱包即身份/保证金NFT流量卡/自动转网/去中心化身份/运营商对接(1.5%手续费)/实体SIM卡领取 等核心功能模块。目标 packages/web。
Gate 1（设计）| 当前: 设计分析 | Skills: 无

---
> [pipeline-init-check] 06-11 17:33

📋 当前项目无架构地图（.claude/CLAUDE.md），建议让 Claude 扫描项目生成一份，后续对话可快速了解项目脉络。

---
> [pipeline-init-check] 06-11 20:31

📋 当前项目无架构地图（.claude/CLAUDE.md），建议让 Claude 扫描项目生成一份，后续对话可快速了解项目脉络。

---
> [pipeline-init-check] 06-12 09:40

📋 当前项目无架构地图（.claude/CLAUDE.md），建议让 Claude 扫描项目生成一份，后续对话可快速了解项目脉络。

---
> [pipeline-init-check] 06-12 09:41

📋 当前项目无架构地图（.claude/CLAUDE.md），建议让 Claude 扫描项目生成一份，后续对话可快速了解项目脉络。

---
> [pipeline-init-check] 06-12 09:42

📋 当前项目无架构地图（.claude/CLAUDE.md），建议让 Claude 扫描项目生成一份，后续对话可快速了解项目脉络。

---
> [pipeline-init-check] 06-12 09:42

📋 当前项目无架构地图（.claude/CLAUDE.md），建议让 Claude 扫描项目生成一份，后续对话可快速了解项目脉络。

---
> [pipeline-init-check] 06-12 09:52

📋 当前项目无架构地图（.claude/CLAUDE.md），建议让 Claude 扫描项目生成一份，后续对话可快速了解项目脉络。

---
> [pipeline-init-check] 06-12 10:46

📋 当前项目无架构地图（.claude/CLAUDE.md），建议让 Claude 扫描项目生成一份，后续对话可快速了解项目脉络。

---
> [pipeline-init-check] 06-12 16:04

📋 当前项目无架构地图（.claude/CLAUDE.md），建议让 Claude 扫描项目生成一份，后续对话可快速了解项目脉络。

---
> [pipeline-init-check] 06-12 16:16

📋 当前项目无架构地图（.claude/CLAUDE.md），建议让 Claude 扫描项目生成一份，后续对话可快速了解项目脉络。

---
> [pipeline-init-check] 06-12 17:42

📋 当前项目无架构地图（.claude/CLAUDE.md），建议让 Claude 扫描项目生成一份，后续对话可快速了解项目脉络。

---
> [pipeline-init-check] 06-12 17:43

📋 当前项目无架构地图（.claude/CLAUDE.md），建议让 Claude 扫描项目生成一份，后续对话可快速了解项目脉络。

---
> [pipeline-init-check] 06-12 20:37

📋 当前项目无架构地图（.claude/CLAUDE.md），建议让 Claude 扫描项目生成一份，后续对话可快速了解项目脉络。

---
> [pipeline-init-check] 06-12 20:41

📋 当前项目无架构地图（.claude/CLAUDE.md），建议让 Claude 扫描项目生成一份，后续对话可快速了解项目脉络。

---
> [pipeline-init-check] 06-13 08:11

📋 当前项目无架构地图（.claude/CLAUDE.md），建议让 Claude 扫描项目生成一份，后续对话可快速了解项目脉络。

---
> [pipeline-init-check] 06-13 14:02

📋 当前项目无架构地图（.claude/CLAUDE.md），建议让 Claude 扫描项目生成一份，后续对话可快速了解项目脉络。

---
> [pipeline-init-check] 06-13 14:05

📋 当前项目无架构地图（.claude/CLAUDE.md），建议让 Claude 扫描项目生成一份，后续对话可快速了解项目脉络。

---
> [pipeline-init-check] 06-13 14:16

📋 当前项目无架构地图（.claude/CLAUDE.md），建议让 Claude 扫描项目生成一份，后续对话可快速了解项目脉络。

---
> [pipeline-init-check] 06-13 14:36

📋 当前项目无架构地图（.claude/CLAUDE.md），建议让 Claude 扫描项目生成一份，后续对话可快速了解项目脉络。

---
> [pipeline-init-check] 06-13 14:36

📋 当前项目无架构地图（.claude/CLAUDE.md），建议让 Claude 扫描项目生成一份，后续对话可快速了解项目脉络。

---
> [pipeline-init-check] 06-13 14:37

📋 当前项目无架构地图（.claude/CLAUDE.md），建议让 Claude 扫描项目生成一份，后续对话可快速了解项目脉络。

---
> [pipeline-init-check] 06-13 14:48

📋 当前项目无架构地图（.claude/CLAUDE.md），建议让 Claude 扫描项目生成一份，后续对话可快速了解项目脉络。

---
> [pipeline-init-check] 06-13 15:15

📋 当前项目无架构地图（.claude/CLAUDE.md），建议让 Claude 扫描项目生成一份，后续对话可快速了解项目脉络。

---
> [pipeline-init-check] 06-13 16:21

📋 当前项目无架构地图（.claude/CLAUDE.md），建议让 Claude 扫描项目生成一份，后续对话可快速了解项目脉络。

---
> [pipeline-init-check] 06-13 16:26

📋 当前项目无架构地图（.claude/CLAUDE.md），建议让 Claude 扫描项目生成一份，后续对话可快速了解项目脉络。

---
> [pipeline-init-check] 06-13 16:26

📋 当前项目无架构地图（.claude/CLAUDE.md），建议让 Claude 扫描项目生成一份，后续对话可快速了解项目脉络。

---
> [pipeline-init-check] 06-13 16:28

📋 当前项目无架构地图（.claude/CLAUDE.md），建议让 Claude 扫描项目生成一份，后续对话可快速了解项目脉络。

---
> [pipeline-init-check] 06-13 16:29

📋 当前项目无架构地图（.claude/CLAUDE.md），建议让 Claude 扫描项目生成一份，后续对话可快速了解项目脉络。

---
> [pipeline-init-check] 06-13 16:31

📋 当前项目无架构地图（.claude/CLAUDE.md），建议让 Claude 扫描项目生成一份，后续对话可快速了解项目脉络。

---
> [pipeline-init-check] 06-13 16:33

📋 当前项目无架构地图（.claude/CLAUDE.md），建议让 Claude 扫描项目生成一份，后续对话可快速了解项目脉络。

---
> [pipeline-init-check] 06-13 16:33

📋 当前项目无架构地图（.claude/CLAUDE.md），建议让 Claude 扫描项目生成一份，后续对话可快速了解项目脉络。

---
> [pipeline-init-check] 06-13 16:34

📋 当前项目无架构地图（.claude/CLAUDE.md），建议让 Claude 扫描项目生成一份，后续对话可快速了解项目脉络。

---
> [pipeline-init-check] 06-13 16:34

📋 当前项目无架构地图（.claude/CLAUDE.md），建议让 Claude 扫描项目生成一份，后续对话可快速了解项目脉络。

---
> [pipeline-init-check] 06-13 16:35

📋 当前项目无架构地图（.claude/CLAUDE.md），建议让 Claude 扫描项目生成一份，后续对话可快速了解项目脉络。

---
> [pipeline-init-check] 06-13 16:40

📋 当前项目无架构地图（.claude/CLAUDE.md），建议让 Claude 扫描项目生成一份，后续对话可快速了解项目脉络。

---
> [pipeline-init-check] 06-13 16:44

📋 当前项目无架构地图（.claude/CLAUDE.md），建议让 Claude 扫描项目生成一份，后续对话可快速了解项目脉络。

---
> [pipeline-init-check] 06-13 16:48

📋 当前项目无架构地图（.claude/CLAUDE.md），建议让 Claude 扫描项目生成一份，后续对话可快速了解项目脉络。

---
> [pipeline-init-check] 06-13 16:48

📋 当前项目无架构地图（.claude/CLAUDE.md），建议让 Claude 扫描项目生成一份，后续对话可快速了解项目脉络。

---
> [pipeline-init-check] 06-13 17:36

📋 当前项目无架构地图（.claude/CLAUDE.md），建议让 Claude 扫描项目生成一份，后续对话可快速了解项目脉络。

---
> [pipeline-init-check] 06-13 18:06

📋 当前项目无架构地图（.claude/CLAUDE.md），建议让 Claude 扫描项目生成一份，后续对话可快速了解项目脉络。

---
> [pipeline-init-check] 06-13 19:53

📋 当前项目无架构地图（.claude/CLAUDE.md），建议让 Claude 扫描项目生成一份，后续对话可快速了解项目脉络。

---
> [pipeline-init-check] 06-13 20:33

📋 当前项目无架构地图（.claude/CLAUDE.md），建议让 Claude 扫描项目生成一份，后续对话可快速了解项目脉络。
