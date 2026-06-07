# Stage: design — 设计分析

> **状态**: completed | **日期**: 2026-06-07 | **Gate**: 1（设计） | **Round**: 1
> **范围**: `packages/web` 视觉重塑（方案 X：深蓝渐变底 + 暖米白卡片 + 沉稳金点缀）
> **产物**: [DESIGN.md](../../design/linkworld/DESIGN.md)（设计系统真源） + 本文件（设计分析 + 选定方向）
> **方式**: design-consultation（约束：建 DESIGN.md 地基）→ design-shotgun（发散：Dashboard 版式变体，用户挑方向）

---

## 0. 已核对真实代码清单（开聊前 + grep 验证）

| 核对项 | 真实情况（已读源码 / grep） | 对设计的影响 |
|--------|------------------------------|--------------|
| `src/index.css` | 单套 `:root` oklch 深色变量（270 紫蓝色相），body 用 `"Inter"` 与 `--font-sans: Geist` 冲突 | 改写为 navy+gold 单套；删 Inter |
| `tailwind.config.ts` | `theme.extend.colors` 全硬编码 HEX（surface/brand/status/text）；`fontFamily.orbitron` 声明 | token 改 `var(--...)` 引用；删 orbitron |
| `src/main.tsx` | RainbowKit `darkTheme({ accentColor:"#3b82f6" })` + `Toaster theme="dark"` | accentColor 改金 #D4AF37；Toaster 保留 dark |
| 旧色值残留 | **13 文件**：layout 3 / wallet 2 / pages 8 / main.tsx（grep `#3b82f6` `brand-blue` 等） | 色值出口统一方案逐个清洗 |
| emoji 当图标 | 字面量 4 文件（Landing/Notifications/Billing/Services）+ `\u{...}` 转义（TabBar 5 项、AppLayout 2 项 GuardCard） | lucide 映射表全量替换 |
| lucide-react | 已装（package）**零引用** | 直接启用 |
| `AmountDisplay` | 默认 `colorClass = text-status-warning`（#f59e0b 琥珀≈金） | 改默认吃 `text-brand-gold`；warning 推向橙 #F08C2E |
| `EmptyState`/`GuardCard` | `icon: string` 类型 | 改 `icon: LucideIcon`（requirement D5） |
| `ui/*` 原子 | 仅 Button/Badge/Tabs（cva + CSS 变量，规范） | 保留重着色 + 新增 Card/Input |

---

## 1. DESIGN.md 锁定的设计系统（摘要）

完整见 [DESIGN.md](../../design/linkworld/DESIGN.md)。三项关键气质决策已与用户逐项确认：

| 决策 | 结论 | 用户确认 |
|------|------|----------|
| 金色调 | **沉稳金 `#D4AF37`**（古典奢华、金属感） | ✅ |
| 卡片色温 | **暖米白 `#F7F3EA`** + 1px 金线描边（烫金请柬质感） | ✅ |
| 字体 | **Display=Space Grotesk / Body=Geist**；删 Inter、弃 Orbitron | ✅ |

**核心 token：**
- 画布渐变：`--bg-canvas: linear-gradient(180deg,#0A1B33,#0C2340)`；hero：`--gradient-hero: linear-gradient(135deg,#0C2340,#1E40AF)`
- 品牌：navy `#0C2340` / royal `#1E40AF` / gold `#D4AF37`（hover `#E0BE4E` / press `#B8942C`）
- 卡面：card `#F7F3EA` / elevated `#FBF8F1` / input `#EFE9DB` / 金线 `rgba(212,175,55,.55)`
- 浅底文字：primary `#0C2340` / secondary `#5A6478` / muted `#7A8294`
- 深底文字：primary `#F7F3EA` / secondary `#B8C2D9` / muted `#6B7A99`
- 状态：success `#22C55E` / **warning `#F08C2E`（改橙避让金）** / danger `#EF4444` / info `#2563EB`
- 圆角：card 16 / btn 12（`--radius`）/ badge full；间距 base 4px comfortable

**架构级三件事（requirement §七 已交接，DESIGN.md 已给方案）：**
1. **色值唯一真源** = `index.css :root` CSS 变量，`tailwind.config` token 全 `var(--...)` 引用 → 消除双轨 + 13 文件硬编码（满足验收 §六-A4）。
2. **lucide 映射表**：DESIGN.md 已列全量 emoji→lucide（TabBar/Header/GuardCard/EmptyState/各页），`icon` prop 改 `LucideIcon`。
3. **缺失原子**：保留 `ui/{button,badge,tabs}` 重着色 + 新增 `ui/card`（暖米白+金线收口手写 div）、`ui/input`；Dialog 优先复用现有 `BottomSheet`/`vaul`。

---

## 2. 选定视觉方向（design-shotgun 发散结果）

在锁定 token 内对**首页 Dashboard**（最具代表性，一屏承载米白卡/金额金/navy 画布/底部导航）生成 3 个版式方向，用户挑选：

| 方向 | 描述 | 结果 |
|------|------|------|
| A) Vault Ledger 金库账本 | 大余额卡置顶 + 紧凑账本列表，私人银行感、信息密度高 | 未选 |
| B) Passport 护照登机牌 | 卡片做票券/护照存根、金线像烫金齿孔，叙事感强、最独特 | 未选 |
| **C) Hero Spotlight 单焦点** | navy→royal 渐变 banner 托巨大居中金色余额 + 2 列圈角磁贴 | ✅ **选定** |

**选定方向 C 的实现要点：**
- 首屏顶部 1/3 为 `--gradient-hero` banner，居中放保证金余额（金色 50px，Geist tabular-nums），副行带「本月账单预估 $42.50」（金色强调）。
- banner 下方 2×2 圈角磁贴（card 14px 圆角）：充值 / 支付账单 / 我的号码 / 流量卡，图标用 royal 蓝 lucide。
- **取舍记录**：C 首屏信息密度最低（账单/号码详情靠磁贴路由进二级页）。缓解：hero 副行保留本月预估金额，关键数字仍在首屏；磁贴角标可带未付账单数提示。
- 该方向只定义 **Dashboard 骨架**；其余 8 页（Deposit/Services/RegionDetail/Billing/BillDetail/Notifications/Cards/Landing）沿用同一 token + 卡片体系，按各页信息结构套用（列表页用账本式卡，表单页用 input 原子）。

> 工件：手写高保真对比页 `~/.gstack/projects/src/designs/dashboard-20260607/board.html` + `approved.json`（design 二进制需 OpenAI key 未配，改手写 HTML 原型，更贴真实 token，零 API 成本）。

---

## 3. requirement §七 待决细节 — 逐项闭环

| # | requirement 待决项 | design 结论 |
|---|---------------------|-------------|
| 1 | 金色色值 + 应用规则 | #D4AF37；金额/主CTA/激活态/品牌强调，详见 DESIGN.md Color |
| 2 | 深蓝渐变精确参数 | canvas `180deg #0A1B33→#0C2340`；hero `135deg #0C2340→#1E40AF` |
| 3 | 米白卡片 + 卡内文字色阶 | card #F7F3EA；文字 #0C2340 / #5A6478 / #7A8294 三级 |
| 4 | 字体方案 | Space Grotesk(Display) + Geist(Body/数据)；删 Inter、弃 Orbitron |
| 5 | 色值出口统一 | CSS 变量唯一真源 + Tailwind `var()` 引用 |
| 6 | lucide 映射表 | DESIGN.md 全量映射 + `icon` prop 改 LucideIcon |
| 7 | 缺失原子组件 | 保留 ui/* 重着色 + 新增 Card/Input |
| 8 | 9 页高保真 | Dashboard 选定方向 C；其余页同体系套用（实现期落地） |

---

## 4. 移交 arch-review 阶段

- 设计系统已锁（DESIGN.md），视觉主方向已选（C）。架构审查重点：
  - **色值出口统一**的落地方式（CSS 变量改写 `:root` + tailwind `var()` 引用）会触及 `index.css` + `tailwind.config.ts` + 13 个业务文件，是最大工作量块（requirement §五已记），需排期与风险评估。
  - `EmptyState`/`GuardCard` 的 `icon` prop 由 `string` 改 `LucideIcon` 是**类型签名变更**，属 UI 层（甲约束允许），但调用点（AppLayout/Notifications/Services）需同步改。
  - 新增 `ui/card`/`ui/input` 原子的 API 设计。
  - `config/constants.ts` 数值修正（MIN_DEPOSIT 10、PLATFORM_FEE_RATE 0.015）属配置层（甲约束唯一例外），与设计无关但需在任务拆解纳入。

## 完成判据
✅ `docs/pipeline/stages/design.md` + `docs/design/linkworld/DESIGN.md` git tracked + commit → advance 推进 arch-review。
