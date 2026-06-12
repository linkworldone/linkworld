# Stage: plan — 任务拆解（子项目 web 3/3）

> **状态**: completed（用户已审批） | **日期**: 2026-06-10 | **Gate**: 2 | **子项目**: web(3/3)
> 产出：plan-review.md（人读，7 节含表格）+ plan-review.json（hook 读，13 pages 各含 components+colors，已过 advance 校验）

## 13 个串行任务（依赖：T0 绿基线 → 接链基建定结构 → 换肤上色；implement 铁律串行+TDD）
| 任务 | 范围 | gate |
|------|------|------|
| T0 基线 | 清 RegisterSheet.tsx:22 TS6133 + 全量 tsc 绿 + 搭 vitest/@testing-library | npm run build 绿 |
| T1 链配置+ABI | chains 加 421614 + contracts 按 deployments.json 单一出口(修 31337 漂移) + abis 从 artifacts 重生成+补 OracleABI 导出 | tsc 绿 |
| T2 精度+常量 | format 默认 6 位 + MIN_DEPOSIT_USDT=10n*10n**6n + SUPPORTED_CURRENCIES=[USDT] + billingApi.toBill 改 bigint | 单测精度 |
| T3 approve+两步态 | USDT allowance/approve(exact 禁 infinite) hooks + TwoStepAction 组件(allowance 跳步+两笔中间态+approve 成功/action 失败回 approved-idle 不 re-approve) | 组件测 |
| T4 对账重构 | useDeposit/useBilling 去双写 + pending 态 + TxStatusBadge + 余额读链/账单轮询后端 confirmed + getLogs 限流修(后端列表或限窗+补 error 态禁静默) + Bill status 加 paying | 组件测 |
| T5 WalletAuth | signedPost helper(非拦截器调 hook) + EIP-712 签名 + 会话级签名一次(nonce+时间窗)；nonce/字段待与后端 signatures.go 对齐 | 中间件测 |
| T6 深蓝金基建 | index.css :root CSS 变量单一出口 + tailwind token→var() + 字体(删 Inter/弃 Orbitron) + 新增 ui/card·ui/input + 金色铁律(卡内金额 navy) | tsc 绿 |
| T7 Deposit 页 | TwoStepAction 两步态 + LockCountdown(getLockExpiry 倒计时+累加顺延提示) + pending + 深蓝金 | 手动+组件 |
| T8 Cards 页 | 双 Tab(NFT/SIM) 去 Admin 发卡按钮+自动发放说明 + NFT 读链 + SIM pendingSync + 深蓝金 | 手动 |
| T9 Billing/BillDetail | paying 态 + FeeBreakdown 读链(useFeeRate/calculateFee) + 深蓝金 | 手动 |
| T10 RegionDetail | 申请号码弹层手续费明细(读链) + 深蓝金 | 手动 |
| T11 其余页换肤 | Landing/Dashboard/Services/Notifications/layout/TabBar/shared/wallet 深蓝金 + emoji→lucide 全量 + 色值单一出口 | grep 旧色值/emoji 清零 |
| T12 测试+构建 | vitest 全绿 + npm run build 绿 + 本地 31337 全链路冒烟 | 全绿 |

## arch-review 红线/前置（§5）
T0 绿基线前置；WalletAuth nonce/EIP-712 跨端待后端 signatures.go；getLogs 二选一(后端端点 vs 限窗)implement 定；web DONE=本地 31337 全链路绿，Arbitrum 端到端(D17)阻塞合约真·上链后置不计入 web DONE；金色卡内 navy 铁律；对账不据 200/tx 成功置终态。

## 移交 implement
从 T0 严格串行 + TDD（vitest 单元 + 组件 mock wagmi + 31337 冒烟）。每任务完成→主 Agent 审查+写 checkpoint→再派下一个。换肤依赖接链定下的 DOM 结构，顺序不可颠倒。
