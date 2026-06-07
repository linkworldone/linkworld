# Design System — Link World Web

> 设计源真理（source of truth）。任何视觉/UI 决策前先读本文件。
> 本轮 Round 1 主题：**The Decentralized Future** — 深蓝渐变底 + 暖米白卡片 + 沉稳金点缀。
> 适用范围：`packages/web`（移动端宽度 web app，max-w 430px）。
> 创建：2026-06-07 · 阶段：design · 已核对真实代码（index.css / tailwind.config.ts / main.tsx + grep）。

---

## Product Context
- **What this is:** 去中心化身份 + 保证金 NFT 流量卡 + 全球号码/运营商对接的 Web3 应用（钱包即身份）。
- **Who it's for:** 数字游民 / Web3 用户，跨境通信 + 链上保证金场景。
- **Space/industry:** Web3 / DePIN 通信，0G Chain 生态。
- **Project type:** 移动优先的 web app（单列、430px 容器、底部 5 项导航）。

## Aesthetic Direction
- **Direction:** Luxury Web3 / Refined-Futuristic（克制的「去中心化未来」）。
- **Decoration level:** intentional —— 主体靠排版 + 金线 + 渐变层次，不堆纹理、不做发光/blob slop。
- **Mood:** 深海军蓝画布上漂浮暖米白卡片，金色像烫金请柬上的金箔——稳重、贵气、可信，而非霓虹赛博。
- **取舍记录（方案 X，requirement D3）:** 放弃全深色，改「深蓝画布 + 浅色卡片」明暗混合体系，卡内文字深色化。这是本轮最大工作量块。

---

## Typography

字体策略：**Display 用 Space Grotesk，正文/UI/数据全用 Geist**。删除现有 `body` 的 `"Inter"` 声明（与 `--font-sans: Geist` 冲突），弃用 Orbitron（crypto 套路字，拉低高级感）。

- **Display / Hero / 品牌标题:** `Space Grotesk` — 几何感 + 未来气息但克制专业，承接「去中心化未来」又不落链圈俗套。仅用于 Landing hero、页面大标题、品牌强调（常配金色）。
- **Body / UI / Labels:** `Geist Variable`（已通过 `@fontsource-variable/geist` 加载）— 干净易读，变量字重。
- **Data / 金额:** `Geist`（启用 `font-variant-numeric: tabular-nums`）— 金额/账单数字等宽对齐，`AmountDisplay` 必用。
- **Code:** 暂无需求；如需 `Geist Mono`。
- **加载策略:** Geist 已自托管（fontsource）；Space Grotesk 用 `@fontsource-variable/space-grotesk` 自托管（与 Geist 同方式，避免 CDN 阻塞）。**清理：** 删 `index.css` body 的 `"Inter"`；删 `tailwind.config.ts` 的 `fontFamily.orbitron`。
- **字号 scale（rem，移动端）:**
  | 级别 | size / line-height | 用途 |
  |------|--------------------|------|
  | display | 2.25rem / 1.1 (Space Grotesk 700) | Landing hero |
  | h1 | 1.5rem / 1.25 (Space Grotesk 600) | 页面主标题 |
  | h2 | 1.25rem / 1.3 (Space Grotesk 600) | 卡片/区块标题 |
  | body | 1rem / 1.5 (Geist 400) | 正文 |
  | sm | 0.875rem / 1.45 (Geist 400) | 次要说明 |
  | xs | 0.75rem / 1.4 (Geist 500) | 标签/角标/tab |
  | amount-lg | 1.75rem / 1.2 (Geist 600 tabular-nums) | 主要金额 |

---

## Color

**Approach:** balanced —— 深蓝为结构主色，金色为灵魂点缀（金额 / 主 CTA / 激活态 / 品牌强调，requirement D4），状态色语义独立。

### 品牌色
| token | hex | 角色 |
|-------|-----|------|
| `brand.navy` | `#0C2340` | 深海军蓝，画布基底 / 卡内主文字 / 结构主色 |
| `brand.navy-deep` | `#0A1B33` | 渐变最深点（画布顶部） |
| `brand.royal` | `#1E40AF` | 皇家蓝，链接 / 次级按钮 / focus ring / info |
| `brand.royal-bright` | `#2952D8` | royal hover |
| `brand.gold` | `#D4AF37` | **沉稳金**：金额 / 主 CTA / 激活态 / 品牌强调 / 卡片金线 |
| `brand.gold-hover` | `#E0BE4E` | gold hover（提亮） |
| `brand.gold-press` | `#B8942C` | gold pressed（压暗） |

### 渐变（requirement §七 item 2 锁定）
- `--gradient-hero`: `linear-gradient(135deg, #0C2340 0%, #1E40AF 100%)` — Landing hero / 品牌时刻（最贴 PRD 参考）。
- `--bg-canvas`: `linear-gradient(180deg, #0A1B33 0%, #0C2340 100%)` — App 全局背景（比 hero 更沉静，让米白卡片更跳）。

### 卡面（暖米白，浮在深蓝上）
| token | hex | 用途 |
|-------|-----|------|
| `surface.card` | `#F7F3EA` | 暖米白卡片主面（D 决策） |
| `surface.card-elevated` | `#FBF8F1` | 更亮的浮起卡/弹层 |
| `surface.input` | `#EFE9DB` | 输入框 / 凹陷区 |
| `surface.card-line` | `#D4AF37` | 卡片 1px 金线描边（可 rgba(212,175,55,.45) 取低调） |

### 文字色阶
**浅底（米白卡内）——深色化三级：**
| token | hex | 对比 / 用途 |
|-------|-----|------|
| `text.on-light.primary` | `#0C2340` | 主文字（navy，超高对比，绑品牌） |
| `text.on-light.secondary` | `#5A6478` | 次文字（slate，≈5:1） |
| `text.on-light.muted` | `#7A8294` | 弱文字/标签（≈3.5:1，仅大字/装饰） |

**深底（navy 画布上，卡外）：**
| token | hex | 用途 |
|-------|-----|------|
| `text.on-dark.primary` | `#F7F3EA` | 主文字（暖白） |
| `text.on-dark.secondary` | `#B8C2D9` | 次文字（浅蓝灰） |
| `text.on-dark.muted` | `#6B7A99` | 弱文字 |

### 语义状态色（与金色严格区分）
| token | hex | 说明 |
|-------|-----|------|
| `status.success` | `#22C55E` | 成功 / active 用户状态 |
| `status.warning` | `#F08C2E` | **改为橙**（原 `#F59E0B` 琥珀与金色撞，推向橙以区分） |
| `status.danger` | `#EF4444` | 错误 / 未付 / 角标 / suspended |
| `status.info` | `#2563EB` | 信息（= royal 系） |

> ⚠️ **关键迁移：** `AmountDisplay` 默认 `colorClass` 从 `text-status-warning` 改为 `text-brand-gold`——金额走金色，warning 让位给橙色，语义彻底分开。

### Dark / Light 说明
本应用是**单一混合主题**（navy 画布 + 米白卡片），非传统 light/dark 双套。`tailwind.config` 仍 `darkMode:"class"` 但实际只有一套体系；`index.css` 只维护一套 `:root`。RainbowKit 用 `darkTheme`（深色弹层贴画布），accentColor 改金色。

---

## 色值出口统一方案（requirement §七 item 5，架构级）

**问题（真实代码核对）：** 当前双轨——业务组件（layout/shared/wallet/pages）吃 `tailwind.config.ts` 的硬编码 HEX token；UI 原子（button/badge/tabs）吃 `index.css` 的 oklch CSS 变量。改主题要两处同步，且 13 个文件散落硬编码 HEX。

**方案：CSS 变量为唯一真源，Tailwind token 全部 `var(--...)` 引用。**
1. `index.css` `:root` 定义全部语义变量（hex 或 oklch 皆可，建议沿用 oklch 与 shadcn 一致）：`--brand-navy`、`--brand-gold`、`--surface-card`、`--text-primary` …
2. `tailwind.config.ts` 的 `theme.extend.colors` 改为引用：`gold: 'var(--brand-gold)'`、`surface: { card: 'var(--surface-card)' }` …
3. 业务组件继续用语义 Tailwind 类（`bg-surface-card` / `text-text-primary` / `bg-brand-gold`），**禁止再写裸 HEX**。
4. UI 原子继续用 shadcn 语义类（`bg-primary` 等），其变量也指向同一 `:root`。

**收益：** 单一出口，换主题只改 `:root`；满足验收 §六-A4「业务组件不再吃硬编码 HEX」。

**shadcn 变量映射（index.css 改写）：**
| 变量 | 新值（oklch 近似） | 对应 |
|------|------|------|
| `--background` | navy 渐变由 body 背景图实现；`--background` 设 `oklch(0.20 0.06 260)`≈#0C2340 | 画布基色 |
| `--foreground` | `oklch(0.95 0.01 80)`≈#F7F3EA | 深底主文字 |
| `--card` | `oklch(0.96 0.015 85)`≈#F7F3EA | 暖米白卡 |
| `--card-foreground` | `oklch(0.25 0.05 260)`≈#0C2340 | 卡内主文字 |
| `--primary` | `oklch(0.42 0.16 265)`≈#1E40AF | 皇家蓝 |
| `--primary-foreground` | `oklch(0.98 0 0)` | 蓝上白字 |
| `--accent` | `oklch(0.76 0.13 85)`≈#D4AF37 | 金 |
| `--accent-foreground` | `oklch(0.25 0.05 260)`≈#0C2340 | 金上深字 |
| `--ring` | = primary（royal） | focus 环 |
| `--destructive` | `oklch(0.62 0.23 27)`≈#EF4444 | 危险 |
| `--border` | `oklch(0.30 0.04 260)` 深底 / 米白卡内用 `--card-line` | 边框 |
| `--radius` | `0.75rem`（保留） | 圆角基准 |

---

## lucide 图标映射表（requirement §七 item 6，验收 §六-A3）

全量替换 emoji（字面量 + `\u{...}` 转义两种写法）。`lucide-react` 已装零用。
**`EmptyState` / `GuardCard` 的 `icon` prop 类型从 `string` 改为 `LucideIcon`（组件）**（requirement D5）。

| 位置 | 现状 | → lucide | 备注 |
|------|------|----------|------|
| TabBar · Home | `\u{1F3E0}` 🏠 | `Home` | |
| TabBar · Services | `\u{1F4F1}` 📱 | `Smartphone` | 地区/号码 |
| TabBar · Deposit | `\u{1F4B0}` 💰 | `Wallet` | |
| TabBar · Bills | `\u{1F4C4}` 📄 | `ReceiptText` | 带未付角标 |
| TabBar · Cards | `\u{1F39F}` 🎟️ | `CreditCard` | 流量卡 NFT |
| Header · 通知铃铛 | 🔔 | `Bell` | |
| AppLayout GuardCard · 押金 | `\u{1F4B0}` 💰 | `Wallet` | Deposit Required |
| AppLayout GuardCard · 封禁 | `⚠️` ⚠️ | `ShieldAlert` | Account Suspended |
| Landing · 品牌 | 🌐 | `Globe` | hero 图标 |
| Services · 搜索 | 🔍 | `Search` | 输入框前缀 |
| Services · 空号码 | 📱 | `SmartphoneNfc` | EmptyState |
| Notifications · 空 | 🔔 | `BellOff` | EmptyState |
| Billing · 警示 | ⚠️ | `AlertTriangle` | 逾期提示 |
| 通用 · 返回 | （Header 箭头） | `ChevronLeft` | 非 emoji，统一 lucide |
| SIM 领取 Tab（新增） | — | `Nfc` 或 `MailPlus` | Cards 双 Tab |

> 实现期按页面实查补全；本表覆盖已 grep 确认的全部 emoji 出现点（TabBar 5 + AppLayout 2 + Landing/Services/Notifications/Billing 字面量）。

---

## 缺失原子组件决策（requirement §七 item 7）

**保留 + 重新着色** `ui/button` `ui/badge` `ui/tabs`（cva + CSS 变量，最规范；只需变量指向新 `:root`）。

**新增 3 个原子**（沉淀手写 div，集中暖米白卡 + 金线处理）：
- `ui/card.tsx` — 暖米白卡：`bg-surface-card` + 1px 金线 + radius-lg + 轻阴影。**一处定义，全站卡片复用**（Dashboard/Deposit/Services/Billing/Cards 等手写卡全部收口）。
- `ui/input.tsx` — `bg-surface-input` + navy 文字 + royal focus ring，统一表单（RegisterSheet、SIM 领取表单、搜索框）。
- `ui/dialog.tsx`（可选）— 与现有 `vaul` Drawer / `BottomSheet` 并存或替换，按实现期评估；优先复用现有 `BottomSheet`。

---

## Spacing
- **Base unit:** 4px；**Density:** comfortable（移动端可点）。
- **Scale:** 2xs(2) xs(4) sm(8) md(12) lg(16) xl(24) 2xl(32) 3xl(48) 4xl(64)。
- **容器:** `max-w-mobile = 430px`（保留）；页面内边距 16px。

## Layout
- **Approach:** grid-disciplined（移动单列，卡片堆叠）。
- **结构:** navy 画布全屏背景 → 内容区 max-w 430px 居中 → 暖米白卡片承载内容 → 固定底部 5 项 TabBar。
- **Border radius:** sm 8px / md 12px(`--radius`) / lg 16px(卡片) / full 9999px(badge/头像/pill)。
- **卡片阴影:** 轻投影 `0 4px 16px rgba(10,27,51,.25)`（深底上浮起感）+ 1px 金线。

## Motion
- **Approach:** intentional（克制，不做 slop 发光/弹跳）。
- **Easing:** enter `ease-out` / exit `ease-in` / move `ease-in-out`。
- **Duration:** micro 80ms（hover/press）/ short 200ms（卡片/tab 切换）/ medium 320ms（BottomSheet/页面过渡）。
- **金色交互:** CTA hover 用 `gold-hover` 提亮 + micro 过渡；不加 shimmer/呼吸动画（避免廉价感）。

---

## main.tsx / 全局收尾（真实代码核对）
- RainbowKit：`darkTheme({ accentColor: "#D4AF37", accentColorForeground: "#0C2340" })`（原 `#3b82f6`）。
- `Toaster`：保持 `theme="dark"` + `richColors`（贴 navy 画布）。
- body：删 `font-family: "Inter"`，统一走 `--font-sans`(Geist)。

## Decisions Log
| Date | Decision | Rationale |
|------|----------|-----------|
| 2026-06-07 | 金色 = 沉稳金 #D4AF37 | 用户选；古典奢华、金属感，最贴「Decentralized Future」 |
| 2026-06-07 | 卡片 = 暖米白 #F7F3EA + 1px 金线 | 用户选；navy+cream+gold 经典奢华三件套，烫金请柬质感 |
| 2026-06-07 | Display=Space Grotesk / Body=Geist，删 Inter、弃 Orbitron | 用户选；现代锐利有科技品味，避开 crypto 套路字 |
| 2026-06-07 | warning #F59E0B→#F08C2E，金额改吃 gold | 解决金色与 warning 撞色，语义分离 |
| 2026-06-07 | 色值唯一真源 = CSS 变量，Tailwind token 引用 var() | 消除双轨与 13 文件硬编码 HEX，满足验收 A4 |
| 2026-06-07 | 保留 ui/* 重着色 + 新增 Card/Input | 集中暖米白卡+金线处理，收口手写 div |
