# color-mapping.md — Link World Web 基线·色值映射

> Web/Tailwind 项目（**非 Flutter**）。色值分两套并存：
> 1. `tailwind.config.ts` 的 `theme.extend.colors` 自定义 token（业务组件直接用，HEX 硬编码）
> 2. `src/index.css` 的 shadcn CSS 变量（oklch，UI 原子组件用）
> 扫描时间：2026-06-06

## 1. Tailwind 自定义色 token（`tailwind.config.ts` → theme.extend.colors）

业务组件（layout/shared/wallet/pages）主要消费这套，**值为硬编码 HEX**。`darkMode: "class"`。

| token | 值 | 用途 / 出现位置 |
|-------|-----|----------------|
| `surface.DEFAULT` | `#0a0a14` | 页面主背景（`bg-surface`，AppLayout/TabBar 等） |
| `surface.card` | `#0f0f1a` | 卡片/抽屉背景（BottomSheet、RegisterSheet `bg-surface-card`） |
| `surface.secondary` | `#1a1a2e` | 次级背景/输入框/拖拽条（`bg-surface-secondary`） |
| `surface.gradient.from` | `#1a1a3e` | 渐变起点（备用，当前页面少用） |
| `surface.gradient.to` | `#0f1a2e` | 渐变终点 |
| `brand.blue` | `#3b82f6` | 主品牌色 / 激活态 / 链接（TabBar 激活、ConnectButton、聚焦边框）；亦是 RainbowKit accentColor |
| `brand.purple` | `#8b5cf6` | 辅助品牌色（Header 头像渐变球 `to-brand-purple`） |
| `brand.cyan` | `#06b6d4` | 点缀色（当前页面少用） |
| `status.success` | `#22c55e` | 成功/激活状态（StatusBadge active） |
| `status.warning` | `#f59e0b` | 警告/金额默认色（AmountDisplay 默认、StatusBadge inactive） |
| `status.danger` | `#ef4444` | 错误/未付/角标（角标背景、校验错误、StatusBadge suspended） |
| `text.primary` | `#e0e0e0` | 主文本（`text-text-primary`，body 默认） |
| `text.secondary` | `#888888` | 次文本（说明/币种后缀 `text-text-secondary`） |
| `text.muted` | `#555555` | 弱文本（标签/未激活 tab `text-text-muted`） |
| `border.DEFAULT` | `#1a1a2e` | 边框（`border-border`，与 surface.secondary 同值） |

其它 extend：`fontFamily.orbitron: ["Orbitron"]`（科技风字体，当前未见使用）；`maxWidth.mobile: "430px"`（移动端容器宽度，全局布局用）。

## 2. shadcn CSS 变量（`src/index.css` → `:root`，oklch 色彩空间）

UI 原子组件（button/badge/tabs）消费这套。**只有一套 `:root`（深色），无独立 light 主题块** —— 即应用强制深色，`darkMode:"class"` 的 `.dark` 覆盖未定义。

字体变量：`--font-sans: 'Geist Variable'`，`--font-heading: var(--font-sans)`（在 `.theme` 下）。body 实际 `font-family: "Inter", system-ui`（与 Geist 变量不一致，存在冲突）。

| 变量 | 值（oklch） | 近似含义 |
|------|-------------|----------|
| `--background` | `oklch(0.07 0.02 270)` | 近黑带紫蓝调主背景 |
| `--foreground` | `oklch(0.88 0 0)` | 浅灰主文本 |
| `--card` / `--popover` | `oklch(0.10 0.02 270)` | 卡片/浮层背景 |
| `--card-foreground` / `--popover-foreground` | `oklch(0.88 0 0)` | 卡片文本 |
| `--primary` | `oklch(0.62 0.21 255)` | 蓝色主色（≈ #3b82f6 同色系） |
| `--primary-foreground` | `oklch(1 0 0)` | 纯白（主色上文字） |
| `--secondary` / `--muted` / `--accent` | `oklch(0.15 0.02 270)` | 次级/静音/强调背景（深紫蓝） |
| `--secondary-foreground` / `--accent-foreground` | `oklch(0.88 0 0)` | 对应文本 |
| `--muted-foreground` | `oklch(0.53 0 0)` | 中灰弱文本 |
| `--destructive` | `oklch(0.577 0.245 27.325)` | 红色（危险） |
| `--border` / `--input` | `oklch(0.18 0.02 270)` | 边框/输入边框 |
| `--ring` | `oklch(0.62 0.21 255)` | 聚焦环（= primary） |
| `--radius` | `0.75rem` | 圆角基准 |
| `--chart-1..5` | 255/290/195/145/55 色相 | 图表配色（蓝/紫/青/绿/橙） |
| `--sidebar*` | 同 card/primary 系列 | 侧边栏（当前移动端布局未用 sidebar） |

## 3. 当前主题风格总结 & 与新主题对比

**当前风格**：深色「赛博/科技夜」—— 主背景近黑偏紫蓝（`#0a0a14` / oklch 270 色相），主色亮蓝 `#3b82f6`，辅以紫 `#8b5cf6`、青 `#06b6d4` 点缀，状态色用标准红/黄/绿。整体偏冷、霓虹感，强制深色无浅色主题。

**两套色值不统一**是最大问题：业务组件吃 Tailwind HEX token，UI 原子吃 oklch CSS 变量，二者色相/明度大体对齐（蓝主色基本一致）但**各自为政，改主题需两处同步改**。

**新主题目标**：深蓝渐变 `#0C2340 → #1E40AF` + 金色点缀。

| 维度 | 当前 | 新主题（目标） | 迁移动作 |
|------|------|----------------|----------|
| 主背景 | `#0a0a14`（近黑紫蓝） | `#0C2340`（深海军蓝）渐变到 `#1E40AF` | 改 `surface.*` + `--background`/`--card`；启用 `surface.gradient` |
| 主色 | `#3b82f6`（亮蓝） | `#1E40AF`（皇家蓝）系 | 改 `brand.blue` + `--primary`/`--ring`；同步 main.tsx RainbowKit accentColor |
| 点缀色 | purple/cyan | **金色**（新增，如 #D4AF37/#F0C75E） | `brand` 下新增 `gold` token + 对应 CSS 变量；replace purple/cyan 点缀 |
| 色相 | 270（紫蓝） | ~255-265（更纯蓝/海军蓝） | oklch 色相微调 |
| 状态色 | 标准红黄绿 | 可沿用 | 基本不动 |

**关键迁移点**：① 统一色值出口（建议业务组件改吃 CSS 变量或语义化 Tailwind token，消除 HEX 硬编码）；② `main.tsx` 的 `accentColor: "#3b82f6"` 与多处 `text-white`/`bg-brand-blue` 硬编码需一并改；③ body 字体 `"Inter"` 与 `--font-sans: Geist` 冲突，重构时顺手对齐；④ Orbitron 字体已声明未用，金色科技标题可考虑启用。
