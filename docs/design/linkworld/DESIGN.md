# Design System — Link World Web

> 设计源真理（source of truth）。任何视觉/UI 决策前先读本文件。
> 本轮 Round 1 主题：**The Decentralized Future** — 深蓝渐变底 + 暖米白卡片 + 沉稳金点缀。
> 适用范围：`packages/web`（移动端宽度 web app，max-w 430px）。
> 创建：2026-06-07 · 阶段：design · 已核对真实代码（index.css / tailwind.config.ts / main.tsx + grep）。

---

## Product Context
- **What this is:** 去中心化身份 + 保证金 NFT 流量卡 + 全球号码/运营商对接的 Web3 应用（钱包即身份）。
- **Who it's for:** 数字游民 / Web3 用户，跨境通信 + 链上保证金场景。
- **Space/industry:** Web3 / DePIN 通信，**Arbitrum Sepolia(421614) 测试网**（本轮由 0G 迁入，PRD R2；保证金/支付币种为 **ERC20 USDT**，6 位精度）。
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
| `text.on-dark.gold` | `#D4AF37` | 深底金额 / 金色强调（navy 上 ≈6.5:1，达标） |

### 金色用色铁律（B2 对比度，WCAG，arch-review 返工）
**金 #D4AF37 在暖米白 #F7F3EA 上对比度 ≈2:1，远低于 WCAG AA 正文 4.5:1（大字 3:1 亦不达标）——卡内金色文字/金额一律禁用，会读不清（钱的应用致命）。**

| 场景（底色） | 金额 / 文字色 | 对比 | 结论 |
|------|------|------|------|
| **卡内（暖米白 #F7F3EA）** 金额 / 正文 | **navy `#0C2340`**（`text-text-primary`） | ≈12:1 | ✓ 默认 |
| **深底（navy 画布上，卡外）** 金额 / 强调 | 金 `#D4AF37`（`text-on-dark.gold`） | ≈6.5:1 | ✓ 达标 |
| 金在暖米白上做文字/金额 | — | ≈2:1 | ✗ **禁用** |

**金色在浅底上只允许非文本用途**：CTA 按钮**填充底**（金底 navy 字 #0C2340 ≈ 同上深字达标）、激活态指示、卡片 1px 金线描边、stepper 激活点、图标小面积点缀——这些不承载需读取的文字信息，不受正文对比度约束。

> ⚠️ **B2 修正（覆盖旧「金额改吃金色」决策）：** `AmountDisplay` 默认色**按所在底色分流**，而非全局吃金色：
> - 卡内（绝大多数金额：余额卡、账单、费用明细）→ 默认 `text-text-primary`（navy）；
> - 深底（少数浮在 navy 画布上的金额，如 Landing/hero 数字、深底统计条）→ 传 `colorClass="text-on-dark-gold"`；
> - 组件提供 `tone="auto|gold-on-dark"` prop 或显式 `colorClass`，**默认 navy**。
> warning 让位给橙 `#F08C2E`，金色与 warning 语义仍分离。

### 语义状态色（与金色严格区分）
| token | hex | 说明 |
|-------|-----|------|
| `status.success` | `#22C55E` | 成功 / active 用户状态 |
| `status.warning` | `#F08C2E` | **改为橙**（原 `#F59E0B` 琥珀与金色撞，推向橙以区分） |
| `status.danger` | `#EF4444` | 错误 / 未付 / 角标 / suspended |
| `status.info` | `#2563EB` | 信息（= royal 系） |

> ⚠️ **关键迁移（B2 已修正，详见下方「金色用色铁律」）：** `AmountDisplay` 默认 `colorClass` 从 `text-status-warning` 改走**按底色分流**——卡内金额走 navy `text-text-primary`（对比 ≈12:1），仅深底金额走金 `text-on-dark-gold`（≈6.5:1）；**不再无条件吃金色**（金在米白卡内 ≈2:1 不达标）。warning 让位给橙色。

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

## 接链交互模式（Round 1 · web 3/3 新增 · 接 Arbitrum Sepolia + ERC20 USDT）

> 本节是「接链交互」的设计源真理。契约签名均已核对真实合约源码（`packages/contracts/contracts/{Deposit,Payment,FeeManager}.sol`）。详细状态机/组件复用/移交风险见 `docs/pipeline/stages/design.md`。
> 总原则：**对账已反转为「链上事件唯一回填终态」**（后端 handoff §1/§4）——UI 绝不能据合约 tx 成功或 HTTP 200 就显示「成功/已到账/已付」。

### 通用：交易三态（贯穿充值/提现/付账，最重要）

| 态 | 含义 | 视觉 token | 图标 | 文案口径 |
|----|------|-----------|------|---------|
| `pending` | 已提交，等链上确认（HTTP 意向 / 监听中 / 未满 K 块） | `status.info`#2563EB 或弱化 `text.on-*.muted` | `Loader2`(animate-spin) / `Clock` | **统一口径**：「处理中 · 约 1-2 分钟」（不向用户暴露「N 块」，K 块为后端内部口径） |
| `confirmed` | 事件回填且满 K 块 | `status.success`#22C55E | `CheckCircle2` | 「已完成」「已到账」「已支付」 |
| `failed` | 交易 revert（链上拒绝）/ reorg 回退（确认后被推翻） | `status.danger`#EF4444 | `XCircle` | revert→「交易被拒绝：{parseContractError}」；reorg→「网络重组，记录已回退，请重试」 |

**铁律**：「可用余额」只算 `confirmed`；`pending` 金额单列「处理中 +N USDT」并弱化，不计入可用、不染成功绿；reorg 回退的未确认记录不缓存当真。金额一律 USDT 6 位（`usdtDecimals` 从 deployments 读，禁硬编码 6 / 禁 18 位 `parseEther`）。

**confirmed 信号来源（B6 锁定，前端不监听链事件做终态）：**
- **余额**：链上读 `Deposit.getDepositAmount(addr)` = source of truth（已自动只含 confirmed 本金，K 块逻辑在合约/后端）。
- **账单 is_paid / 充值·提现历史终态**：**轮询后端 GET status**（`GET /api/bills/:wallet`、`/api/deposit/:wallet`），pending 期 `refetchInterval ≈ 5s`，转 confirmed 后停轮询（`refetchInterval=false`）。
- 前端**不**用 `useWatchContractEvent` / 不订阅 `BillPaid`/`DepositMade`/`DepositWithdrawn` 事件来置终态；事件回填是后端 event_sync 的职责（handoff §1/§4）。
- **reorg 缓存策略**：pending 相关 react-query 用短 `staleTime`（≈2s）+ 上述短 `refetchInterval`；confirmed 后放长 `staleTime`，停轮询。

**pending 超时兜底（happy path 之外，必做）：** pending 超过约 2 分钟仍未 confirmed → 展示安心文案「确认较慢，可安全离开，到账后会通知」+ **Arbiscan 逃生链接**（`https://sepolia.arbiscan.io/tx/{txHash}` 供用户自查），不报错、不回滚 UI、不染绿；继续轮询直到 confirmed/failed。

**pending 首屏 loading：** 测试网 RPC 慢，pending 态首次进入与读链查询未回前一律用 skeleton（金额位 `--` 占位），不闪「0」或空白。

### 1. USDT approve 两段式（充值 + 付账）
- 合约：`deposit(uint256 amount)` 非 payable；`payBill(uint256 billId)` 非 payable nonReentrant；均需前置 `usdt.approve(spender, exactAmount)`。**exact-amount，禁 infinite approve**（资损硬约束 handoff §2）。
- approve 额：充值 = `amount`；付账 = `amount + FeeManager.calculateFee(amount)`（**直读合约 calculateFee，不自算**）。
- allowance 检测：`allowance(user, spender) ≥ 需求额` → 跳过 Step 1 直达 Step 2。
- Stepper：`Step 1/2 Approve → Step 2/2 Deposit`（激活态金色 #D4AF37）；文案传递安全感「仅授权本次所需 N USDT（不做无限授权）」。
- 按钮文案（idle / 签名中 / 确认中）：`Approve N USDT` / `Waiting for signature…` / `Confirming approval…`；`Deposit N USDT` / `Waiting for signature…` / `Confirming on chain…`。
- 收尾：合约成功 → POST 意向端点 → **UI 进 pending 态**（不显示终态成功）。

#### `TwoStepAction` 完整状态机（B5，approve→action 两笔串行，最易写崩处）
```
idle (输入金额，校验 amount≥10 USDT 且 ≤钱包余额)
  │
  ├─[allowance(user,spender) ≥ 需求额]──────────────▶ 跳过 Approve，直达 ②
  │
  └─[allowance < 需求额] ① Approve:
        approve-sign (等钱包签名)
          ├─用户拒签──▶ idle（不进 pending；toast「授权已取消」）
          └─已签 ──▶ approving (Approving…)
                       ├─approve revert/失败──▶ approve-failed（回 idle 分支，可重试 Approve）
                       └─成功 ──▶ confirming-approval (Confirming approval…)
                                    └─满确认 ──▶ ② (allowance 已足)
   ② Deposit/Pay:
        action-sign (等钱包签名)
          ├─用户拒签──▶ ★approved-idle（已授权可重试存入；不 re-approve）
          └─已签 ──▶ acting (Depositing… / Paying…)
                       ├─action revert/失败──▶ ★approved-idle（已授权可重试存入；不 re-approve）
                       └─成功 ──▶ POST 意向 ──▶ pending（通用三态，见上）
```
**关键回退分支（★）：** Approve 已成功后，若 Deposit/Pay 签名被拒或链上失败，**回退到「已授权可重试存入」态（`approved-idle`）而非回到 Approve**——精确额度授权已在链上，重试时检测 `allowance ≥ 需求额` 直接跳 ②，**绝不 re-approve**（避免重复授权 + 多花 gas）。UI 在此态：Step 1 显示「已授权 ✓」、主按钮变「重试存入 N USDT」、副提示「已完成授权，仅需确认存入」。
- **充值/付账复用同一状态机**：仅 `spender`（Deposit / Payment）、需求额（`amount` / `amount + calculateFee`）、文案（Deposit / Pay）不同。
- 现有 `useTxState`（`useTransactionFlow.ts`）是**单笔五态**（idle/pending-signature/pending-confirmation/success/error），撑不起两笔串行 + 跳步 + 回退；`TwoStepAction` 在其上编排两个 `useTxState` 实例（approveTx / actionTx）+ allowance 读值驱动 step 切换。

### 2. 对账 pending 态（充值/提现/付账三处复用通用三态）
- 充值（`Deposit.tsx`）：去掉即时「Deposit confirmed」绿 toast → pending「处理中」，余额读链 `getDepositAmount` 或等 `DepositMade` 事件 confirmed。
- 提现（`Deposit.tsx`）：**废弃凭 txHash 记账** → pending「提现确认中」，等 `DepositWithdrawn` 事件。
- 付账（`Billing.tsx`/`BillDetail.tsx`）：Bill status 由 `unpaid/paid/overdue` **扩展为 `unpaid/paying(确认中)/paid/overdue`**；`paying` 用 info 蓝 + `Loader2`，**禁止显示绿「已付」**；`is_paid` 等 `BillPaid` 事件。
- 历史/列表项：每条带状态点（confirmed 绿 / pending 橙旋转 / failed 红）。

### 3. 锁仓倒计时（Deposit 提现区）
- 合约：`getLockExpiry(addr)` 返回时间戳秒；提现 `require(block.timestamp >= _lockExpiry)`（**边界 `>=`，到期时刻即可提，倒计时归 0 即解锁**）；**每次充值 `lockExpiry += 30 days`（顺延累加）**，首存 = `now + 30d`，提现成功归 0。
- 状态：`expiry==0` → 无锁仓（无押金或已全提）；`now < expiry` → Withdraw **禁用** + 倒计时「锁仓中 · 剩余 Dd Hh Mm」+ `Lock`；`now ≥ expiry` → Withdraw **可点** + `Unlock` +「锁仓已满，可提取本金」（**去「利息」**：`DepositWithdrawn(user,principal,0)` interest 恒 0，旧文案「本金+利息」误导，arch-review ⚠️）。
- **倒计时 `>=` 边界**：`LockCountdown` 解锁判定与合约 `block.timestamp >= _lockExpiry` 严格对齐——`remaining <= 0`（即 `now >= expiry`）即视为解锁可提，不可用 `<`（差 1 秒会卡住到期用户）。
- **关键提示**：因 lockExpiry 累加，Deposit 页必须明示「再次充值将把锁仓期顺延 30 天」，避免用户困惑倒计时变长。
- 刷新：每分钟（剩余 ≤1 天时每秒）；边界判定与合约 `>=` 对齐。

### 4. Cards 双 Tab + 移除 Admin 发卡
- `ui/tabs`（TabsList variant=line）：Tab1「流量卡 NFT」/ Tab2「SIM 领取」。
- Tab1：**移除「Issue Monthly Card (Admin)」按钮**（`onlyOracle`，用户调必 revert）→ 换「自动发放」说明卡（`Sparkles`/`Info` +「保证金锁仓满 30 天后每月自动发放，无需手动领取」）；NFT 列表读链真实数据，每卡含「不可转卖」标记 +「销毁后 30 天」规则文案；空态 `EmptyState`(icon=`CreditCard`)。
- **NFT 列表取数修正（B6，getLogs 限流 + 禁静默置空）：** 现 `useTrafficCards` 用 `getLogs({fromBlock:0n, toBlock:'latest'})` 全量扫块，在 Arbitrum 公共 RPC 必限流/超时，且 `catch` 后 `setTokenIds([])` 是 **silent failure**（用户误以为「没有卡」）。修复二选一（implement 定）：① **后端给列表端点**（首选，彻底避开扫块）；② 退而 **限定 fromBlock 窗口**（从 deployments 记录的合约部署块号开始，而非 0n）。无论哪种，**必须补 error 态**：取数失败渲染「加载失败，点击重试」+ `refetch`，**禁止 catch 后静默置空**（要能区分「真无卡」与「加载失败」）。
- Tab2 SIM 领取（前端降级 R9）：选国家 + 收件表单 → 提交 toast 成功 + 写 `localStorage pendingSync`（`utils/pendingSync.ts` 已存在）+「全球通即将推出」；图标 `Nfc`/`MailPlus`。

### 5. 手续费读链展示（D12）
- 删写死 `PLATFORM_FEE_RATE`；新增 `useFeeRate` 读 `FeeManager.getFeeRate()`（基点，150=1.5%，分母 10000）；精确费额直读 `calculateFee(amount)`。
- 展示位：① 申请号码弹层（RegionDetail）② 账单页 / BillDetail 费用明细。
- 形态：明细行「平台手续费 (1.5%)：N USDT」，1.5% = `getFeeRate/100`，金额 = `calculateFee(amount)`。
- loading / 读链失败：显示 skeleton 或「--」，**不写死兜底**（避免误导）。

### 6. WalletAuth 签名 UX（新增鉴权）
- 受保护写端点（deposit/withdraw/bills-pay/service 写）请求带钱包签名头；**读端点（GET）不加**（handoff §6）。
- UX：写操作前多一次钱包签名弹窗，**必须与交易签名视觉区分**——文案「请在钱包中签名以验证身份（不消耗 gas）」，图标 `ShieldCheck`/`PenLine`。
- 拒签/失败：toast「身份签名被取消，操作未提交」，不进 pending。
- **B1 铁律（会话级签名一次，硬约束，不得降级到「每次写操作签」）：** 一个钱包会话只签一次身份签名（携带 `nonce + timestamp 时间窗`），签名结果缓存在内存（非 localStorage 持久化，避免泄露重放），会话内后续写操作复用，**禁止每次 deposit/withdraw/pay 都弹签名**（频繁弹窗劝退 + 与交易签名混淆）。会话过期（时间窗到期）或切换钱包 → 重新签一次。
- **已签名会话态视觉**：会话已持有有效身份签名时，写操作直接进交易流程不再弹身份签名；可在 Header/钱包区用 `ShieldCheck` + 「已验证」浅金角标提示当前会话已鉴权，过期则降级提示需重签。
- **nonce 来源（跨端待对齐）：** nonce 由后端下发还是前端生成时间窗，需与后端 `signatures.go` 对齐；签名方案锁定 **EIP-712**（结构化、钱包可读、优于 EIP-191）。字段顺序/domain 在 implement 阶段与后端敲定，本设计锁定「会话级一次 + EIP-712」两条铁律。
- **实现约束（架构坑规避）：** 用 `signedPost(path, body)` helper 封装（内部 `await signTypedData()` 拿缓存或新签的签名 → 带签名头调 axios），**不在 axios 全局拦截器里调 React hook**（拦截器非 React 上下文，调 `useSignTypedData` 必崩）。

### 7. 交互态 / 响应式补全（arch-review ⚠️「视觉态空洞」闭合）
- **香槟金 #F0C75E 定稿（不再「可选」）：明确做**——仅用于卡片 1px 金线描边的**渐变高光端**：`linear-gradient(135deg, #D4AF37 0%, #F0C75E 50%, #D4AF37 100%)`（停靠点锁死 0%/50%/100%，中点 #F0C75E 为高光），制造烫金请柬的金属反光。**主金仍 #D4AF37**，#F0C75E 不作独立文字/填充色、不进 token 表语义色（R12 视觉不变，不引入新主色）。
- **approve 成功 / deposit(pay) 失败中间态 UI**：即上方 `approved-idle` 态——Step 1「已授权 ✓」（金 check）、主按钮「重试存入/付款」、副文「已完成授权，仅需确认」；不报红错（授权未损失），失败原因 toast 走 `parseContractError`。
- **failed 区分 reorg vs revert**：revert（链上拒绝）→ 红 + `parseContractError` 具体原因 + 「重试」；reorg（确认后回退）→ 红 + 「网络重组，记录已回退」+ 「重试」，**两者文案不同**（见三态表）。
- **320 / 375 窄屏**：max-w 430px 容器内布局须在 320px 不溢出（金额大字 + stepper + 双按钮换行而非压缩）；**所有可点元素最小触摸区 44×44px**（按钮、tab、倒计时旁的图标、关闭/返回），padding 撑够不靠字号。
- **pending 文案三处统一**：充值/提现/付账 pending 均用「处理中 · 约 1-2 分钟」口径（见三态表），不各写各的。

### lucide 图标映射补充（接 DESIGN.md 上方映射表）
| 场景 | lucide |
|------|--------|
| 交易确认中 / loading | `Loader2`(animate-spin) / `Clock` |
| 已完成 | `CheckCircle2` |
| 失败 | `XCircle` |
| 锁仓中 / 可提现 | `Lock` / `Unlock`(或 `LockOpen`) |
| 自动发放说明 | `Sparkles` / `Info` |
| SIM 领取 Tab | `Nfc` / `MailPlus` |
| 手续费 | `Percent` / `ReceiptText` |
| 身份签名(WalletAuth) | `ShieldCheck` / `PenLine` |
| 两步 approve | `KeyRound`(授权) → `ArrowRight`(下一步) |

### 接链专用新增组件建议（接 DESIGN.md「缺失原子组件决策」）
- `TwoStepAction` — 封装 approve→action 两步态（stepper + 按钮文案 + allowance 跳步），充值/付账复用。
- `TxStatusBadge` — 通用交易三态徽章（pending/confirmed/failed）。
- `LockCountdown` — 读 `getLockExpiry` → 渲染倒计时 / 解锁态。
- `FeeBreakdown` — 费用明细（小计 + 手续费 + 合计），读链费率。

---

## Decisions Log
| Date | Decision | Rationale |
|------|----------|-----------|
| 2026-06-07 | 金色 = 沉稳金 #D4AF37 | 用户选；古典奢华、金属感，最贴「Decentralized Future」 |
| 2026-06-07 | 卡片 = 暖米白 #F7F3EA + 1px 金线 | 用户选；navy+cream+gold 经典奢华三件套，烫金请柬质感 |
| 2026-06-07 | Display=Space Grotesk / Body=Geist，删 Inter、弃 Orbitron | 用户选；现代锐利有科技品味，避开 crypto 套路字 |
| 2026-06-07 | warning #F59E0B→#F08C2E，金额改吃 gold | 解决金色与 warning 撞色，语义分离 |
| 2026-06-07 | 色值唯一真源 = CSS 变量，Tailwind token 引用 var() | 消除双轨与 13 文件硬编码 HEX，满足验收 A4 |
| 2026-06-07 | 保留 ui/* 重着色 + 新增 Card/Input | 集中暖米白卡+金线处理，收口手写 div |
| 2026-06-09 | 接链交互模式：交易三态(pending/confirmed/failed)、approve 两段式、锁仓倒计时、Cards 双Tab、手续费读链、WalletAuth 签名 | web 3/3 接 Arbitrum+USDT；对账反转为事件驱动，UI 禁止据 tx/200 显示终态 |
| 2026-06-09 | 香槟金 #F0C75E 仅作金线渐变高光端(可选)，主金仍 #D4AF37 | 回应 PRD E18 两金值；R12 锁定视觉不变，不引入新主色 |
| 2026-06-10 | **WalletAuth 会话级签名一次铁律**（nonce+时间窗，EIP-712），禁每次写操作签 | arch-review B1：唯一能压门槛杠杆不得降级 implement；signedPost helper 封装、不在 axios 拦截器调 hook；nonce 来源跨端与 signatures.go 对齐 |
| 2026-06-10 | **金色对比度修正（覆盖「金额吃金」）**：卡内金额/文字用 navy #0C2340，金色仅用于深底金额/CTA填充/激活态/卡片金线 | arch-review B2：金在米白 ≈2:1 WCAG 不达标（钱的应用致命）；金在 navy ≈6.5:1 达标；AmountDisplay 默认按底色分流 |
| 2026-06-10 | confirmed 终态来源=读链 getDepositAmount + 轮询后端 GET status（pending refetch≈5s），前端不监听链事件做终态 | arch-review B6：getLogs(fromBlock:0n) 公共 RPC 限流；K 块逻辑全留后端 |
| 2026-06-10 | 香槟金 #F0C75E 由「可选」改「明确做」：仅金线渐变高光端，停靠点 0/50/100% | arch-review ⚠️ 视觉态空洞需定稿 |
