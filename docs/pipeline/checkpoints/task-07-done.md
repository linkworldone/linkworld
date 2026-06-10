# Task 07 — T6 深蓝金视觉基建（web 3/3）

> 子项目 web(3/3) | 状态 DONE | 2026-06-10

### 产出文件
- packages/web/src/index.css（重写 :root 深蓝金语义变量单一出口 navy/gold/royal/surface-card/text-on-light·on-dark 三级/status/渐变；shadcn 变量同源；body navy 渐变背景；删 Inter；Space Grotesk import）
- packages/web/tailwind.config.ts（colors 全改 var(--...)；新增 fontFamily.{sans,display,data}+backgroundImage 渐变+shadow-card；删 fontFamily.orbitron）
- packages/web/index.html（删 Orbitron CDN）
- packages/web/src/main.tsx（RainbowKit accentColor #D4AF37 / foreground #0C2340）
- packages/web/src/components/ui/{card,input}.tsx（新建：暖米白+金线 card / 深蓝金 input 44px 触摸高）
- packages/web/src/components/shared/AmountDisplay.tsx（金色铁律分流 tone auto→卡内 navy / gold-on-dark→深底金；font-data tabular-nums）
- package.json/pnpm-lock.yaml（@fontsource-variable/space-grotesk）
- 3 个新测试

### git commit
61b4fe4 feat: web T6 深蓝金基建(CSS 变量单一出口/tailwind token→var/字体/ui-card·input/金色铁律)

### TDD
基建以 token/组件为主；新增 card/input/AmountDisplay 分流 5 测证明组件可渲染 + 金额底色分流正确。

### 测试结果
npx tsc --noEmit 0 error。npm run build ✓。npx vitest run：14 files / 54 passed（前序 49+新 5，无回归）。grep：tailwind.config 裸 HEX=0（全 var()）；旧值 #3b82f6/#0a0a14/#8b5cf6/Inter/Orbitron 在 config/css/html 清零（index.css :root 29 处 hex 是唯一真源）。主 Agent 已独立复跑确认 tsc 0+build ✓+54 测。

### code-simplifier
CSS 变量单一出口消除双轨色值（业务 token 与 shadcn 同喂一套真源）；ui/card·input 收口手写卡/表单。

### spec review
按 DESIGN.md（深蓝金色值/金色铁律/字体/CSS 变量单一出口/ui-card·input）+ arch-review B2（金色对比度卡内 navy）执行。AmountDisplay 默认卡内 navy、深底金、不写死 hex。未越界（页面级 emoji→lucide/布局/精修留 T7-T11；接链逻辑 T1-T5 未碰）。

### 设计还原
深蓝金 token 系统对齐 DESIGN.md（navy #0C2340/金 #D4AF37/暖米白 #F7F3EA）；金色铁律（卡内金额 navy）落 AmountDisplay；字体 Space Grotesk+Geist。

### 复用检查
复用 shadcn CSS 变量机制、viem、现有 ui 原子；新增 ui/card·input 供 T7-T11 收口；token var() 让现有语义类自动改色。

### 设计稿对照
数值对照：navy #0C2340/gold #D4AF37/cream #F7F3EA（DESIGN.md）✅；tailwind 裸 HEX 0（全 var()）✅；字体 Space Grotesk(display)+Geist(sans/data) 删 Inter/Orbitron ✅；AmountDisplay 卡内 navy/深底金分流（3 测）✅；RainbowKit accentColor 金 ✅；54 测/tsc 0/build ✓ ✅。

### 新增组件
新增 ui/card·ui/input；AmountDisplay 加 tone 分流。

### 新增色值
深蓝金全套语义变量（index.css :root 单一出口）：brand-navy #0C2340/gold #D4AF37/royal + surface-card 暖米白 #F7F3EA/card-line 金线 + text-on-light·on-dark 三级 + 渐变 hero/canvas/gold-line。

### ⚠️ 遗留（带入 T7-T11）
- 旧 token 现为 no-op（已从 config 移除，页面渲染无效类，build 绿但视觉不全，待替换）：brand-blue 残留 13 文件、brand-purple/cyan/surface-gradient/surface-secondary/font-orbitron 残留 11 文件 → 各 re-skin 任务替换为 royal/gold/暖米白卡。
- 基建组件供下游：ui/card 收口手写卡、ui/input 收口表单、AmountDisplay 深底传 tone="gold-on-dark"。
- emoji→lucide 清单（T11，权威表 theme-migration.md §2）：layout/{AppLayout,Header,TabBar}、pages/{Billing,Dashboard,Landing,Notifications,Services}、operatorApi.ts。
