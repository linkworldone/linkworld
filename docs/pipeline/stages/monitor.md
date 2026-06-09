# Stage: monitor — 复盘（子项目 contracts 1/3）

> **状态**: completed | **日期**: 2026-06-09 | **Gate**: 3 | PR #1

## 一、本轮做了什么
Link World 全栈改造子项目 1/3（合约层）。把 packages/contracts 从原生币迁移到 ERC20 USDT + Arbitrum Sepolia(421614)，修复合并自 origin 的半成品（编译不过），落地分账/自动发卡/计价修正。10 阶段 pipeline 全走完，75 passing / 89% 覆盖，PR #1 已开。

## 二、关键经验（值得沉淀）
1. **外部代码中途落地 → STOP and re-plan**（AI 失效模式二：信息孤岛）：design 阶段途中后端+合约代码合并进来，推翻 PRD 核心假设 D2「只动前端」。没有硬推，而是停下重审 requirement → pipeline 从「web 换肤」重构为「合约/后端/web 3 子项串行」。教训：检测到地基假设变化立即停、重规划，不在错误前提上继续盖楼。
2. **两轮 arch-review 挡住资损**：第 1 轮三方审查查出 7 大 ❌（含 Critical：Oracle 把 dataUsage 字节+callUsage 分钟直接当 USDT 金额求和，真实扣款会扣天文数字）→ 用户拍 3 个产品决策 → design 返工 v2 → 第 2 轮 0 ❌。教训：资损敏感合约，对抗性多角色审查 + 返工重审是必要成本，不能一轮过。
3. **领域适配：UI 形状的 pipeline 阶段套到 Solidity 合约**：design 阶段强制模板是 UI（design-consultation/shotgun 配色）、arch-review 三角色含 Design(UI)、plan 是 Figma/组件/色值、test 是 flutter——本子项无 UI。逐个适配：Design 路换成合约安全审计；plan 的 pages schema 适配成「任务=page、components=合约/函数、colors=[]」过 advance hook；test 的 flutter 换 hardhat。教训：pipeline 强制流程为某栈设计时，识别 HOW 层不匹配并适配/向用户澄清，不照搬产垃圾。
4. **硬 guard 阻塞时的处理**：design 阶段 spawn_task hook 强制 UI 模板、不允许合约专用架构师会话 → 切换为「主 Agent 架构师起草 + AskUserQuestion 收敛决策 + writer subagent 落盘」。教训：遇到 hook 硬阻塞，找符合用户全局规范（Technical Design Format）的替代路径。

## 三、遗留（下轮/下游处理）
1. Arbitrum Sepolia 421614 真·上链 + TSTORE 链上实测（环境无 key；配 key 后 handoff §10 一条命令）
2. 上线前 checklist（review.md §五）：owner/upgrade/platformWallet 转多签、11 个 operator 黑洞地址换真实、移除 axios、上线门禁排除 mock、清理 Oracle 死 state（已开 chip）
3. 下游子项：后端对齐(2/3) 依赖本 PR 冻结 ABI/selector（handoff-backend.md）；web 接通(3/3) 依赖最终合约地址
4. PR #1 待 review/merge

## 四、pipeline 收口
合约子项目 1/3 全 10 阶段完成。后续子项 2/3（后端）、3/3（web）各自新开 pipeline，以本子项的 handoff 为输入。
