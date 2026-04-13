# Stage: design — 设计分析

> **状态**: completed | **日期**: 2026-04-12 | **Gate**: 1

## 产出摘要

基于产品文档和 requirement 阶段确认的需求边界，完成交互流程、状态机、API 映射、组件复用方案和设计规范分析。所有分析结果已整合入设计规格文档。

## 关键决策

| # | 决策 | 理由 |
|---|------|------|
| 1 | 四状态状态机（UNREGISTERED/INACTIVE/ACTIVE/SUSPENDED） | 覆盖产品文档定义的所有用户生命周期 |
| 2 | AppLayout 统一守卫 + GuardCard 引导 | 比路由 loader 更简单，引导卡片即时渲染无跳转 |
| 3 | 底部 Sheet 统一交互模式 | 所有操作确认（存入/提取/支付/申请号码）用同一交互范式 |
| 4 | 5 个 React Query hooks 对应 5 个 Mock Service | 一对一映射，职责清晰，未来替换时只改 service 层 |
| 5 | 深色主题 + 品牌渐变蓝紫 | #0a0a14 背景 + linear-gradient(135deg, #3b82f6, #8b5cf6) 品牌色 |

## 设计分析内容

### 交互流程
- Landing: 连接钱包 → 注册弹窗 → Dashboard
- Dashboard: 状态总览 → Quick Actions 分流
- Deposit: 余额查看 → 存入/提取 Sheet → 历史记录
- Services: 搜索 → 选国家 → 选运营商 → 申请号码 Sheet
- Billing: Unpaid/History 切换 → 账单详情 → 支付 Sheet
- Notifications: New/Earlier 分组 → Mark all read

### 状态机
```
UNREGISTERED → INACTIVE → ACTIVE ⇄ SUSPENDED
                  ↑                      |
                  └──────────────────────┘
```

### API → Entity 映射
- userService → User
- depositService → DepositInfo, DepositRecord
- operatorService → Region, Operator, VirtualNumber
- billingService → Bill
- notificationService → Notification

### 组件复用方案
- AppLayout（Header + TabBar + 守卫）：全部内页共享
- 底部 Sheet：4 种操作共享（Deposit/Withdraw/Pay/Apply）
- StatusBadge：Dashboard + Billing 共享
- AmountDisplay：Dashboard + Deposit + Billing 共享
- GuardCard：INACTIVE/SUSPENDED 受限页面共享

## 产出文件

| 文件 | 内容 |
|------|------|
| docs/superpowers/specs/2026-04-12-linkworld-frontend-design.md | 设计规格文档（含完整设计分析） |
| docs/pipeline/stages/design.md | 本文件 |

## 用户确认的事项

- 全部页面移动端布局 wireframe 通过浏览器 mockup 确认
- 状态守卫方案（AppLayout 层）确认
- 状态管理方案（React Query + wagmi hooks）确认
