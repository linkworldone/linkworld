# Stage: scan — 项目基线扫描

> **状态**: completed | **日期**: 2026-06-07 | **Gate**: 0 | **Round**: 1

## 产出摘要
扫描 `packages/web` 现有代码，产出 4 份基线文档到 `docs/design/linkworld/`。本轮 scan 复用上一轮成果——CodeGraph 索引报告 "Already up to date"（代码自上次索引后无变更），4 份基线文档（2026-06-06 扫描）经主 Agent 逐份复核，内容完整且与当前代码一致，无需重新扫描。现有 web 端完成度约 80%：9 个页面，核心用户流程（注册→存保证金→购买服务→支付账单）已打通。业务逻辑层（hooks/services/utils/types/config）100% 可复用；UI 层约 70% 待按新深蓝主题重设计。

## 关键决策
| # | 决策 | 理由 |
|---|------|------|
| 1 | 保留全部业务逻辑层 | hooks/services/utils/types/config 与 UI 解耦，重构只动表现层 |
| 2 | 重设计 pages/* 与 layout/*，替换 index.css 色值 | 主题色换深蓝渐变 #0C2340→#1E40AF，当前为 #3b82f6 亮蓝 |
| 3 | UI 基础组件(ui/*)按新设计系统决定保留或替换 | 由后续 design 阶段定夺 |
| 4 | 本轮 scan 直接复用现有基线文档，不重扫 | CodeGraph 索引 up-to-date + 文档经复核完整且最新，重扫只会生成相同内容 |

## 产出文件
| 文件 | 内容 |
|------|------|
| docs/design/linkworld/components.md | 13 个组件清单(3 UI+3 layout+5 shared+2 wallet)，含 props/变体 |
| docs/design/linkworld/color-mapping.md | shadcn CSS 变量 + Tailwind 扩展色 token + 硬编码色值 + 深蓝替换计划 |
| docs/design/linkworld/utils.md | 26+ 个 hooks/工具函数 + 5 真实 API + 1 mock(notifications) |
| docs/design/linkworld/project-scan.md | 技术栈 + 目录树 + 9 页面路由 + 服务层 + 7 合约 Web3 集成 + 复用度矩阵 |

## 关键发现
- 当前主题色：primary #3b82f6 亮蓝，背景 #0a0a14 近黑；目标深蓝渐变 #0C2340→#1E40AF + 金色点缀
- 色值双轨：业务组件吃 tailwind.config.ts 的 HEX token，UI 原子吃 index.css 的 oklch CSS 变量——改主题需两处同步，重构应统一出口
- 只有 3 个 UI 原子组件(Button/Badge/Tabs)，缺 Card/Input/Dialog/Sheet 封装，页面手写较多
- 大量 emoji 当图标，已装 lucide-react 却基本没用——新主题建议统一换 lucide
- 5 个真实后端 API（User/Operator/Deposit/Billing/Usage）+ 1 个 mock（Notifications，后端未就绪）
- 7 个合约已集成（UserRegistry/Deposit/Payment/ServiceManager/TrafficCard/FeeManager/Oracle），wagmi+RainbowKit，0G Chain
- 注意：config/abis/index.ts 漏导出 OracleABI（文件存在）；main.tsx RainbowKit accentColor 硬编码 #3b82f6
- body 字体 "Inter" 与 --font-sans: Geist 冲突；Orbitron 字体已声明未用

## 用户确认的事项
- 仅重构 web 端（packages/web）
- 去掉 figma 与 visual-test 阶段（无 Figma 源 + React 非 Flutter）
- mode = strict，10 阶段完整 pipeline
- 主题色换深蓝渐变 linear-gradient(135deg, #0C2340 0%, #1E40AF 50%) + 金色点缀
