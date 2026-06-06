# Stage: scan — 项目基线扫描

> **状态**: completed | **日期**: 2026-06-05 | **Gate**: 0

## 产出摘要
扫描 packages/web 现有代码，产出 4 份基线文档到 docs/design/linkworld/。现有 web 端完成度约 80%：9 个页面（Landing→Dashboard→Deposit/Services→Billing/Notifications/Cards），核心用户流程（注册→存保证金→购买服务→支付账单）已打通。业务逻辑层（hooks/services/utils/types/config）100% 可复用；UI 层约 70% 待按新深蓝主题重设计。

## 关键决策
| # | 决策 | 理由 |
|---|------|------|
| 1 | 保留全部业务逻辑层 | hooks/services/utils/types/config 与 UI 解耦，重构只动表现层 |
| 2 | 重设计 pages/* 与 layout/*，替换 index.css 色值 | 主题色换深蓝渐变 #0C2340→#1E40AF，当前为 #3b82f6 亮蓝 |
| 3 | UI 基础组件(ui/*)按新设计系统决定保留或替换 | 由后续 design 阶段定夺 |

## 产出文件
| 文件 | 内容 |
|------|------|
| docs/design/linkworld/components.md | 13 个组件清单(3 UI+3 layout+5 shared+2 wallet)，含 props/变体 |
| docs/design/linkworld/color-mapping.md | 28 个 CSS 变量 + 15 个 Tailwind 扩展 + 12 处硬编码色值 + 深蓝替换计划 |
| docs/design/linkworld/utils.md | 26+ 个 hooks/工具函数 + 5 真实 API + 1 mock(notifications) |
| docs/design/linkworld/project-scan.md | 技术栈 + 目录树 + 9 页面路由 + 服务层 + 7 合约 Web3 集成 + 复用度矩阵 |

## 关键发现
- 当前主题色：primary #3b82f6 亮蓝，背景 #0a0a14 近黑；目标深蓝渐变 #0C2340→#1E40AF
- 5 个真实后端 API（User/Operator/Deposit/Billing/Usage）+ 1 个 mock（Notifications，后端未就绪）
- 7 个合约已集成（UserRegistry/Deposit/Payment/ServiceManager/TrafficCard/FeeManager/Oracle），wagmi+RainbowKit，0G Chain testnet
- 待补功能：账单详情接口、真实邮件验证、真实通知、流量卡 NFT 页面
- 离线优先：pendingSync 用 localStorage 重试失败的后端同步（合约+后端双写兜底）

## 用户确认的事项
- 仅重构 web 端（packages/web）
- 去掉 figma 与 visual-test 阶段
- mode = strict
- 主题色换深蓝渐变 linear-gradient(135deg, #0C2340 0%, #1E40AF 50%)
