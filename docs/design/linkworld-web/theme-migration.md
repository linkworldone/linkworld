# theme-migration.md — Link World web(3/3) 深蓝金换肤 delta

> 扫描 2026-06-09 | 子项目 web(3/3) | 已核对真实代码
>
> 复用基线（不重写）：`docs/design/linkworld/DESIGN.md`（设计源真理：色值/字体/lucide 映射/出口方案）、`docs/design/linkworld/color-mapping.md`（旧→新色值映射）、`docs/design/linkworld/components.md`、`docs/design/linkworld/project-scan.md`、`docs/design/linkworld/utils.md`。
>
> 本文件只产**换肤 delta**：本次 scan 实测的旧色值/emoji 残留清单 + 覆盖范围 + 单一出口方案（与接链改动面正交，接链见 `web-alignment-surface.md`）。

---

## 0. 一句话

DESIGN.md 已锁定「深海军蓝画布 + 暖米白卡片 + 沉稳金点缀」完整规范（色值/字体/图标映射/CSS 变量单一出口）。本次 scan 核对真实代码，确认旧色值与 emoji 残留的**精确分布**，供 implement 阶段按图索骥。换肤与接链可正交并行（不改同文件时），但 `main.tsx`/`wagmi.ts`/`Deposit.tsx`/`Cards.tsx` 接链与换肤都会碰，需串行或协调。

---

## 1. 旧色值残留实测清单（grep 核对 `packages/web/src`）

### 1.1 定义源（tailwind.config.ts）— 改这里是收口关键
`packages/web/tailwind.config.ts` 现硬编码 HEX token（需全部改为 `var(--...)` 引用，DESIGN.md「色值出口统一方案」）：

| token | 现值 | → 新（DESIGN.md） |
|-------|------|-------------------|
| `surface.DEFAULT` | `#0a0a14` | navy 画布 `#0C2340`（实由 body 背景渐变实现） |
| `surface.card` | `#0f0f1a` | 暖米白 `#F7F3EA` |
| `surface.secondary` | `#1a1a2e` | `surface.input #EFE9DB` |
| `surface.gradient.from/to` | `#1a1a3e`/`#0f1a2e` | `--gradient-hero` / `--bg-canvas`（navy 渐变） |
| `brand.blue` | `#3b82f6` | `brand.royal #1E40AF` |
| `brand.purple` | `#8b5cf6` | 删 → 金 `brand.gold #D4AF37` |
| `brand.cyan` | `#06b6d4` | 删 |
| `status.warning` | `#f59e0b` | `#F08C2E`（橙，与金分离） |
| `fontFamily.orbitron` | `["Orbitron"]` | **删**（DESIGN.md 弃 Orbitron） |

新增：`brand.gold/gold-hover/gold-press`、`surface.card-elevated/card-line`、`text.on-light.*` / `text.on-dark.*` 三级文字色阶（DESIGN.md Color 节）。

### 1.2 index.css（shadcn oklch 变量）
单套 `:root` 深色，需整套改为 DESIGN.md「shadcn 变量映射」表的新 oklch 值（`--background`/`--card`/`--card-foreground`/`--primary`/`--accent`(金)/`--ring` 等）。body 删 `font-family:"Inter"`，统一走 `--font-sans`(Geist)；新增 Space Grotesk（Display）。

### 1.3 业务组件硬编码引用（实测分布）
- **`brand-blue`：12 文件**（`bg-brand-blue`/`text-brand-blue`/`border-brand-blue`/`focus:border-brand-blue`）：
  `layout/TabBar.tsx`、`layout/Header.tsx`、`wallet/RegisterSheet.tsx`、`wallet/ConnectButton.tsx`、`pages/{Deposit,Services,Dashboard,Cards,Landing,Notifications,BillDetail,Billing}.tsx`。→ 改 `royal`（链接/次级/focus）或 `gold`（激活态/CTA/金额）按语义区分。
- **`brand-purple`：4 文件**（Header 头像渐变球 `to-brand-purple` 等）→ 删紫，改金或 navy 渐变。
- **`surface-gradient`：3 文件**（`Deposit.tsx`/`Cards.tsx` 的 `bg-gradient-to-br from-surface-gradient-from to-surface-gradient-to` 余额卡）→ 改 navy/金渐变或暖米白卡。
- **`#3b82f6` 裸 HEX：1 文件**（`main.tsx` RainbowKit `accentColor`）→ 金 `#D4AF37`（accentColorForeground `#0C2340`），改 `darkTheme`。
- **`text-status-warning`：3 处**（含 `AmountDisplay` 默认 `colorClass`）→ **金额改吃 `text-brand-gold`**（DESIGN.md ⚠️ 关键迁移），warning 让位橙色。

### 1.4 AmountDisplay 默认色（DESIGN.md 关键迁移）
`components/shared/AmountDisplay.tsx` 现 `colorClass = "text-status-warning"`（琥珀）→ 改 `text-brand-gold`。金额走金色，语义与 warning 彻底分开。

---

## 2. emoji → lucide 残留实测清单（验收 C9）

`lucide-react` 已安装。emoji 有**两种写法**：`\u{...}` 转义 + 字面量，均需替换（DESIGN.md lucide 映射表为权威）。本次 grep 确认出现点：

| 文件:行 | 现状 | → lucide（DESIGN.md 映射） |
|---------|------|----------------------------|
| `layout/TabBar.tsx:6` | `\u{1F3E0}` 🏠 | `Home` |
| `layout/TabBar.tsx:7` | `\u{1F4F1}` 📱 | `Smartphone` |
| `layout/TabBar.tsx:8` | `\u{1F4B0}` 💰 | `Wallet` |
| `layout/TabBar.tsx`（Bills/Cards 行） | 📄 / 🎟️ | `ReceiptText` / `CreditCard` |
| `layout/AppLayout.tsx:34` | `\u{1F4B0}` 💰 | `Wallet`（GuardCard 押金）；封禁态 `ShieldAlert` |
| `pages/Dashboard.tsx:26` | `\u{1F4B0}` 💰 | `Wallet` |
| `pages/Services.tsx:47` | 🔍 | `Search` |
| `pages/Services.tsx:83` | 📱 | `SmartphoneNfc`（空号码 EmptyState） |
| `pages/Notifications.tsx:45` | 🔔 | `BellOff`（空态）/ Header 铃铛 `Bell` |
| `pages/Landing.tsx:34` | 🌐 | `Globe`（hero 品牌） |
| `pages/Billing.tsx:55` | ⚠️ | `AlertTriangle`（逾期警示） |
| Header 通知铃铛 | 🔔 | `Bell` |
| Cards SIM 领取 Tab（新增） | — | `Nfc` / `MailPlus` |

> DESIGN.md D5：`EmptyState`/`GuardCard` 的 `icon` prop 类型从 `string` 改为 `LucideIcon`（组件）。implement 期按页面实查补全（本表已覆盖 grep 确认的全部出现点）。

---

## 3. 9 页 + 组件覆盖清单（验收 C7：100% 应用，无遗漏页）

**9 页**（`pages/`）：`Landing`、`Dashboard`、`Services`、`RegionDetail`、`Deposit`、`Billing`、`BillDetail`、`Cards`、`Notifications`。
**组件**：
- `layout/`：`AppLayout`、`TabBar`、`Header`（emoji + brand-blue 重灾区）。
- `shared/`：`AmountDisplay`（金额改金）、`BottomSheet`、`StatusBadge`、`EmptyState`、`GuardCard` 等。
- `wallet/`：`ConnectButton`、`RegisterSheet`（brand-blue）。
- `ui/`：`button`/`badge`/`tabs` 保留重着色（吃 CSS 变量，改 `:root` 即生效）；新增 `ui/card`（暖米白+金线，收口手写卡）、`ui/input`（DESIGN.md「缺失原子组件决策」）。

逐页换肤 = 替 emoji + 替 brand-blue/purple/surface-gradient 类 + 余额卡/手写卡改暖米白 + 金额改金。

---

## 4. 色值单一出口方案（DESIGN.md「色值出口统一方案」复用）

**问题**：现双轨——业务组件吃 `tailwind.config.ts` 硬编码 HEX token，UI 原子吃 `index.css` oklch 变量，改主题要两处同步。
**方案**：`index.css :root` 为唯一真源（语义变量 `--brand-navy/--brand-gold/--surface-card/--text-*`），`tailwind.config.ts` 的 token 全改 `var(--...)` 引用；业务组件继续用语义 Tailwind 类（`bg-surface-card`/`text-brand-gold`），**禁止再写裸 HEX**；UI 原子 shadcn 类指向同一 `:root`。
**收益**：换主题只改 `:root`，满足验收 C10「色值单一出口」+ C8「grep 查无旧色值残留」。

---

## 5. 换肤验收对照（PRD §五 C）

| # | 验收点 | 本 scan 核对状态 |
|---|--------|------------------|
| C7 | 9 页 + layout/shared/wallet 100% 深蓝金 | 待 implement；覆盖清单见 §3 |
| C8 | grep 无旧色值（`#3b82f6`/`#0a0a14`/`#0f0f1a`/`bg-brand-blue`/`to-brand-purple`/`#8b5cf6`/`#06b6d4`） | 现存：brand-blue×12、brand-purple×4、#3b82f6×1(main.tsx)、tailwind.config 定义源；清零靠改定义源 + 业务类 |
| C9 | grep 无 emoji 当图标（全 lucide） | 现存 emoji 见 §2（TabBar/AppLayout/Dashboard/Services/Notifications/Landing/Billing） |
| C10 | 卡片米白/深底渐变/金色仅金额·CTA·激活·品牌；色值单一出口 | 方案见 §4；金额改金见 §1.4 |
| E18 | 金色 #D4AF37 / 香槟金 #F0C75E、字体、渐变停靠点、米白色阶以 design 为准 | DESIGN.md 已锁定，按其执行 |

---

## 6. 换肤遗留 / 协调点

1. **换肤 ⇄ 接链同文件协调**：`main.tsx`（accentColor 金 + wagmi 链）、`wagmi.ts`、`Deposit.tsx`（余额卡 + 两步态）、`Cards.tsx`（余额卡 + 双 Tab）四处两条线都碰，需串行或主 Agent 协调，避免冲突。
2. **字体新增自托管**：DESIGN.md 要 `@fontsource-variable/space-grotesk`（Display）；Geist 已自托管。删 body `"Inter"` + tailwind `orbitron`。
3. **新增原子组件** `ui/card`/`ui/input` 属换肤范畴但要收口全站手写卡（Dashboard/Deposit/Services/Billing/Cards），与接链页改动重叠，注意顺序。
4. **grep 清零基线**：implement 完成后按 C8/C9 的 grep 模式跑一遍，本文件 §1.3/§2 的精确清单即归零核对表。
