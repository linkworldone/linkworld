# Link World Web — 配色映射 (color-mapping)

> 反映 **当前 HEAD** 真实主题系统，仅描述现状，不含改动或迁移建议。当前是 **两套并存** 的色彩体系：
> 1. `index.css` 的 shadcn **oklch** token（`:root` 变量，ui 原语用）
> 2. `tailwind.config.ts` 的 **hex** 语义别名（surface/brand/status/text，业务组件/页面用）
> 第三类是散落在 main.tsx 链上 UI 配置里的硬编码色值。三个来源逐节梳理如下。

## 1. 当前主题体系 A — index.css oklch token（shadcn）

`packages/web/src/index.css` 顶部 import：`tw-animate-css`、`shadcn/tailwind.css`、`@fontsource-variable/geist`，随后 `@tailwind base/components/utilities`。

`.theme` 块定义字体：`--font-heading: var(--font-sans)`、`--font-sans: 'Geist Variable', sans-serif`。

`:root` 全部 token（**oklch 体系**，暗色基调，色相 270≈蓝紫）：

| Token | 值 (oklch) | 含义 |
|-------|-----------|------|
| `--background` | `oklch(0.07 0.02 270)` | 全局背景（极暗蓝紫） |
| `--foreground` | `oklch(0.88 0 0)` | 主前景文字（浅灰） |
| `--card` | `oklch(0.10 0.02 270)` | 卡片背景 |
| `--card-foreground` | `oklch(0.88 0 0)` | 卡片文字 |
| `--popover` | `oklch(0.10 0.02 270)` | 浮层背景 |
| `--popover-foreground` | `oklch(0.88 0 0)` | 浮层文字 |
| `--primary` | `oklch(0.62 0.21 255)` | 主色（亮蓝） |
| `--primary-foreground` | `oklch(1 0 0)` | 主色上文字（白） |
| `--secondary` | `oklch(0.15 0.02 270)` | 次级背景 |
| `--secondary-foreground` | `oklch(0.88 0 0)` | 次级文字 |
| `--muted` | `oklch(0.15 0.02 270)` | 弱化背景 |
| `--muted-foreground` | `oklch(0.53 0 0)` | 弱化文字（中灰） |
| `--accent` | `oklch(0.15 0.02 270)` | 强调背景 |
| `--accent-foreground` | `oklch(0.88 0 0)` | 强调文字 |
| `--destructive` | `oklch(0.577 0.245 27.325)` | 危险/错误（红） |
| `--border` | `oklch(0.18 0.02 270)` | 边框 |
| `--input` | `oklch(0.18 0.02 270)` | 输入框 |
| `--ring` | `oklch(0.62 0.21 255)` | 焦点环（=primary） |
| `--chart-1` | `oklch(0.62 0.21 255)` | 图表蓝 |
| `--chart-2` | `oklch(0.65 0.18 290)` | 图表紫 |
| `--chart-3` | `oklch(0.70 0.15 195)` | 图表青 |
| `--chart-4` | `oklch(0.75 0.18 145)` | 图表绿 |
| `--chart-5` | `oklch(0.70 0.18 55)` | 图表橙 |
| `--radius` | `0.75rem` | 圆角基准 |
| `--sidebar` | `oklch(0.10 0.02 270)` | 侧栏背景 |
| `--sidebar-foreground` | `oklch(0.88 0 0)` | 侧栏文字 |
| `--sidebar-primary` | `oklch(0.62 0.21 255)` | 侧栏主色 |
| `--sidebar-primary-foreground` | `oklch(1 0 0)` | 侧栏主色文字 |
| `--sidebar-accent` | `oklch(0.15 0.02 270)` | 侧栏强调 |
| `--sidebar-accent-foreground` | `oklch(0.88 0 0)` | 侧栏强调文字 |
| `--sidebar-border` | `oklch(0.18 0.02 270)` | 侧栏边框 |
| `--sidebar-ring` | `oklch(0.62 0.21 255)` | 侧栏焦点环 |

`@layer base` 其余规则：
- `* { @apply border-border outline-none }`
- `body { @apply bg-surface text-text-primary; font-family: "Inter", system-ui, sans-serif; ... }` —— **注意冲突点**：body 实际用的是 tailwind.config 的 `bg-surface`/`text-text-primary`（hex 体系），且 font-family 写死 `"Inter"`，与 `.theme` 里声明的 Geist **不一致**。
- `html { @apply font-sans }`

> 这套 oklch token 主要被 **ui 原语**（button/badge/tabs，类似 `bg-primary`/`text-muted-foreground`/`border-border`）消费。

## 2. 当前主题体系 B — tailwind.config.ts hex 语义别名

`packages/web/tailwind.config.ts`：`darkMode: "class"`，`content: ["./index.html","./src/**/*.{js,ts,jsx,tsx}"]`。

`theme.extend`：

**fontFamily**
- `orbitron: ["Orbitron", "sans-serif"]`（仅 Landing logo 用 `font-orbitron`，但 Orbitron 字体未在 package.json/index.css 引入 —— 实际会 fallback sans-serif）

**colors（全 hex，业务组件/页面消费）**

| 类名前缀 | 子键 | hex | 用途 |
|----------|------|-----|------|
| `surface` | DEFAULT | `#0a0a14` | 页面背景（`bg-surface`） |
| `surface` | card | `#0f0f1a` | 卡片（`bg-surface-card`） |
| `surface` | secondary | `#1a1a2e` | 次级面（`bg-surface-secondary`） |
| `surface` | gradient.from | `#1a1a3e` | 渐变起（`from-surface-gradient-from`） |
| `surface` | gradient.to | `#0f1a2e` | 渐变止（`to-surface-gradient-to`） |
| `brand` | blue | `#3b82f6` | 品牌蓝（active/链接/CTA） |
| `brand` | purple | `#8b5cf6` | 品牌紫（渐变/头像） |
| `brand` | cyan | `#06b6d4` | 品牌青 |
| `status` | success | `#22c55e` | 成功 |
| `status` | warning | `#f59e0b` | 警告（金额默认色） |
| `status` | danger | `#ef4444` | 危险/红点 |
| `text` | primary | `#e0e0e0` | 主文字 |
| `text` | secondary | `#888888` | 次文字 |
| `text` | muted | `#555555` | 弱文字 |
| `border` | DEFAULT | `#1a1a2e` | 边框 |

**maxWidth**
- `mobile: "430px"`（移动端容器 `max-w-mobile`，全站外壳/抽屉/导航用）

`plugins: []`。

> 两套体系数值并不完全对齐：oklch `--background` 与 hex `surface.DEFAULT #0a0a14` 是各自独立定义的，色相/明度接近但不等价。

## 3. 散落硬编码颜色热点（文件:行）

对 `src/pages` + `src/components` grep `#hex / linear-gradient / radial-gradient / bg-gradient / rgba`：

| 文件:行 | 内容 | 类型 |
|---------|------|------|
| `components/layout/Header.tsx:27` | `bg-gradient-to-br from-brand-blue to-brand-purple`（头像圆点） | tailwind 渐变 |
| `pages/Landing.tsx:26` | `bg-gradient-to-r from-brand-blue to-brand-purple bg-clip-text`（logo 文字渐变） | tailwind 渐变 |
| `pages/Landing.tsx:33` | `bg-gradient-to-br from-brand-blue to-brand-purple`（圆形 icon） | tailwind 渐变 |
| `pages/Dashboard.tsx:35` | `bg-gradient-to-br from-surface-gradient-from to-surface-gradient-to`（余额卡） | tailwind 渐变 |
| `pages/Deposit.tsx:65` | `bg-gradient-to-br from-surface-gradient-from to-surface-gradient-to`（余额卡） | tailwind 渐变 |
| `pages/Cards.tsx:55` | `bg-gradient-to-br from-surface-gradient-from to-surface-gradient-to`（卡面） | tailwind 渐变 |

热点统计：
- **裸 hex / rgba() / CSS linear-gradient 写在组件源码：0 处**（所有 hex 都集中在 tailwind.config.ts）。
- **tailwind `bg-gradient` 工具类：6 处**（2 类：品牌渐变 brand-blue→brand-purple 3 处；面板渐变 surface-gradient 3 处）。
- `bg-black/60` 半透明遮罩 2 处（BottomSheet、RegisterSheet），未计入 hex。
- main.tsx 有 2 处链上 UI 写死蓝色：`darkTheme({ accentColor: "#3b82f6" })`、`ConnectButton.tsx` 的 `bg-brand-blue`。

**结论**：颜色高度集中——hex 全部落在 tailwind.config.ts、oklch 全部落在 index.css；组件/页面源码内无任何裸 hex。仅有的散落点是 6 处 tailwind `bg-gradient` 工具类、2 处 `bg-black/60` 遮罩，以及 main.tsx 的 RainbowKit `accentColor` 与 ConnectButton 的 `bg-brand-blue`。

## 4. 第三来源 — 链上 UI 硬编码色值（main.tsx / ConnectButton）

这两处不经 tailwind.config 也不经 index.css，是独立写死的色值来源：

| 文件:行 | 写法 | 色值 | 说明 |
|---------|------|------|------|
| `src/main.tsx:26` | `darkTheme({ accentColor: "#3b82f6" })` | `#3b82f6` | RainbowKit 钱包弹窗主题强调色（裸 hex，与 brand.blue 同值但独立定义） |
| `src/components/wallet/ConnectButton.tsx:12` | `bg-brand-blue text-white` | `#3b82f6` / `#ffffff` | 自定义连接按钮，引用 tailwind brand.blue + `text-white`（未走 ui/Button 原语） |

## 5. 三来源现状小结

| 来源 | 文件 | 色值形态 | 消费方 |
|------|------|----------|--------|
| 体系 A | `src/index.css` `:root` | oklch（蓝紫 270 色相，暗色） | shadcn ui 原语（button/badge/tabs，`bg-primary`/`border-border` 等语义类） |
| 体系 B | `tailwind.config.ts` `theme.extend.colors` | hex（surface/brand/status/text/border） | 业务组件 + 页面（`bg-surface`/`text-brand-blue`/`bg-status-danger` 等别名类） |
| 第三类 | `main.tsx` + `ConnectButton.tsx` | 裸 hex / `text-white` | RainbowKit 弹窗 + 自定义连接按钮 |

> 现状事实：体系 A 与体系 B 各自独立定义、数值不完全对齐（如 oklch `--background` vs hex `surface.DEFAULT #0a0a14`，色相/明度接近但不等价）；body 用体系 B 的 `bg-surface`/`text-text-primary` 且 font-family 写死 `"Inter"`，与 `.theme` 声明的 Geist 不一致；`font-orbitron` 在 Landing 引用但 Orbitron 字体未在 package.json/index.css 安装，实际 fallback sans-serif。
</content>
