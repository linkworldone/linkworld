# Stage: arch-review — 三角色架构审查

> **状态**: completed | **日期**: 2026-04-13 | **Gate**: 1

## 产出摘要

三个并行审查角色（Product/Tech VP/Senior Eng）对设计规格进行了全面审查。发现 5 个阻塞项，全部已修复。10 个风险项记录到 plan 阶段处理。

## 审查结果

### 阻塞项（已修复）

| # | 来源 | 问题 | 修复 |
|---|------|------|------|
| 1 | Product | 实时用量详情缺失（PRD 条目 5） | Dashboard 增加用量详情区域 |
| 2 | Product | 虚拟号码管理页缺失（PRD 条目 7） | Services 页增加 My Numbers Tab |
| 3 | Eng | Dashboard 数据依赖漏掉 useOperator | 已补充 |
| 4 | Eng | getUserProfile 未注册返回值未定义 | 改为 Promise<User \| null> |
| 5 | Eng | 缺少 sendVerificationCode 接口 | 已补充到 userService |

### 风险项（移交 plan 阶段）

| # | 来源 | 问题 |
|---|------|------|
| 1 | Product | 多虚拟号码场景展示 |
| 2 | Product | 提前结账入口缺失 |
| 3 | Product | 逾期扣款倒计时警告缺失 |
| 4 | CEO | Tailwind 配色需更新为 semantic tokens |
| 5 | CEO | BrowserRouter 位置需从 main.tsx 移除 |
| 6 | CEO | shadcn/ui 初始化需前置 |
| 7 | CEO | 桌面端需 max-w 容器约束 |
| 8 | Eng | Sheet 拖拽手势需 vaul 库 |
| 9 | Eng | Mock Service 需加延迟模拟 |
| 10 | Eng | Bill 金额精度问题 |

## 关键决策

| # | 决策 | 理由 |
|---|------|------|
| 1 | 增加实时用量展示 | PRD 核心体验，不可省略 |
| 2 | Services 页增加 My Numbers 管理 | PRD 明确要求，operatorService 已有对应接口 |
| 3 | getUserProfile 返回 null 表示未注册 | 与状态机 UNREGISTERED 对应，守卫逻辑更清晰 |
| 4 | 拆分 register 和 sendVerificationCode | 注册流程两步操作需要两个独立接口 |

## 产出文件

| 文件 | 内容 |
|------|------|
| docs/superpowers/specs/2026-04-12-linkworld-frontend-design.md | 已更新（5 处修复） |
| docs/pipeline/stages/arch-review.md | 本文件 |

## 用户确认的事项

- 阻塞项修复方案由审查角色提出并直接应用
