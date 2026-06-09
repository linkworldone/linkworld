# Stage: arch-review — 三角色审查（CEO/Design/Eng）（子项目 web 3/3）

> **状态**: BLOCKED（6 ❌ 阻塞，需返工 design 重跑） | **日期**: 2026-06-09 | **Gate**: 1
> 审查对象：design.md + DESIGN.md（深蓝金视觉+接链交互稿）。三角色全员（web 真有 UI，Design 路首次对口）。冲突优先级：安全/可行性 > 工程 > 视觉。

## 一、总裁决：BLOCK（6 ❌）
方向对、契约理解准确（Eng 逐行核对真实 ABI/合约/hook 零虚报）、视觉系统成熟（零 AI-slop），但工程闭环不全 + 1 个无障碍硬伤 + 1 个决策留虚。

## 二、❌ 阻塞项（去重）
| # | 阻塞 | 来源 | 修复方向 |
|---|------|------|----------|
| B1 | WalletAuth 签名频次只写"建议会话级一次"，降级甩 implement（唯一能压门槛的杠杆）| CEO | design 锁死铁律：会话级签名一次(nonce+时间窗)，禁每次写操作签；写进 DESIGN.md Decisions Log，作为 signatures.go 对齐硬约束 |
| B2 | 金色 #D4AF37 在暖米白 #F7F3EA 卡内对比度≈2.0:1(WCAG AA 4.5:1 不达标)，而金额默认改吃金色且浮在米白卡=钱的应用读不清致命 | Design | 卡内金额改 navy #0C2340 字色，金色仅用于深底(navy 上金额/CTA/激活态 6.5:1 达标)；DESIGN.md text 色阶+§3 金额改金处同步修 |
| B3 | tsc -b main 基线已红(RegisterSheet.tsx:22 'isSuccess' TS6133)，build=tsc -b&&vite build 现在就过不了，接链重写无绿基线分不清新旧错 | Eng | implement T0 先清此 pre-existing error；design 测试策略注明绿基线前置 |
| B4 | 项目零测试设施(devDeps 无 vitest、src 零 *.test.*)，资损敏感(6 位精度)+对账反转+两笔非原子授权重写无自动化验证网，design 无测试策略 | Eng | design 新增「测试策略」章节：装 vitest + 单元(format 精度/constants/billingApi bigint/approve 额/LockCountdown 边界)必测 + 组件(TwoStepAction/TxStatusBadge mock wagmi)+本地 31337 冒烟 |
| B5 | TwoStepAction 两步态状态机未画：useTxState 是单笔五态，撑不起 approve→deposit 两笔串行+allowance 跳步+approve 成功/deposit 失败回退(接链最易写崩处) | Eng | design 画出 TwoStepAction 状态机(allowance≥amount 跳步、两笔中间态、失败回退 UI 分支) |
| B6 | 「前端如何确认 confirmed」三选一未定 + useTrafficCards getLogs(fromBlock:0n) 全量扫块在 Arbitrum 公共 RPC 必限流/超时且 catch 后静默置空(silent failure) | Eng | design 锁定：余额读链 getDepositAmount + 账单/历史轮询后端 status(refetchInterval)，不前端监听链事件做终态；NFT 列表迁后端/限定 fromBlock 窗口 + 补 error 态(不静默吞错) |

## 三、⚠️ 同批建议修（非阻塞，强烈建议随返工处理）
- **pending 超时兜底缺失**(CEO/Design)：状态机只画 happy+reorg，缺「事件不达/超时」分支 → 补「确认中超时→Arbiscan 逃生链接+安心文案『可安全离开，到账后通知』」；「约 N 块」翻译成「约 1-2 分钟」。
- **billingApi.toBill parseFloat 资损升必改**(Eng)：对 6 位最小单位字符串做 JS 浮点加减(.toFixed(2) 单位语义错+大额超 MAX_SAFE_INTEGER)→ 改 bigint。
- **WalletAuth 拦截器不能调 React hook**(Eng)：axios 拦截器调 useSignTypedData 架构坑 → 封装 signedPost(path,body) helper(内部 await signTypedData 再带头调 axios)，非全局拦截器；EIP-712 优于 191。
- **DepositWithdrawn interest 恒 0**(Eng)：Deposit.sol:80 emit(...,principal,0)，design §3.4 文案「可提取本金+利息」误导 → 去掉「利息」。
- **换肤依赖接链新结构的隐性串行**(CEO)：文档叙述为正交并行，实际换肤依赖接链定下的 DOM 结构(Deposit 两步态/余额卡)→ implement 顺序「接链先定结构→换肤后上色」，否则返工。
- **web DONE 验收边界重画**(CEO)：显式定义 web DONE=本地 31337 全链路绿；Arbitrum 端到端(D17)+对账三态真链行为=后置强制验收，阻塞于合约上链，不计入 web DONE 也别成孤儿。
- **视觉态空洞**(Design)：香槟金 #F0C75E 可选→明确做(金线渐变停靠点)或 NOT in scope；渐变中间停靠点锁死；approve 成功/deposit 失败中间态 UI；WalletAuth「已签名」会话态视觉；pending 首屏 loading/skeleton(测试网慢 RPC)；320/375 窄屏+44px 触摸区；pending 文案三处口径统一；failed 态区分 reorg vs revert。
- **reorg react-query 缓存**(Eng)：pending 相关 query 用短 staleTime/refetchInterval，confirmed 后放长；K 块逻辑全留后端，前端只轮询 confirmed 字段。

## 四、✅ 三方认可
- CEO：对账反转方向正确(消灭前端自述假数据唯一正解)、跳过 shotgun 复用深蓝金对、计费/精度伪胜利防死、诚实标阻塞。
- Design：色值单一出口架构级正解、金色克制有 taste、烫金请柬隐喻一致、换肤覆盖 grep 实测可验收、对账 pending 不染绿铁律对、新增组件抽象合理。总体 7.4/10。
- Eng：契约理解准确零虚报(selector/payable/6 位/事件/onlyOracle/锁仓累加全对得上)、余额读链 source-of-truth、deployments 单一出口方向对、exact approve 额算法对、6 位精度漏改面圈全。

## 五、下一步（硬规则：有 ❌ → 改 design → 重跑）
架构师返工 design v2 + DESIGN.md 闭合 B1-B6 + 关键 ⚠️ → 重跑 arch-review(至少 Design+Eng 复审) → 0 ❌ 进 plan。B3(tsc 红)记为 implement T0 前置。

## 六、三方原始摘要
- CEO：1 ❌(B1 WalletAuth 留虚)+3 ⚠️+4 ✅。结论：方向对不返工重做，WalletAuth 锁铁律+pending 兜底+验收边界重画。
- Design：1 ❌(B2 金色对比度)+6 ⚠️，总体 7.4/10。结论：视觉优秀但金额金色 2:1 必修，交互态系统性空洞待补。
- Eng：4 ❌(B3 tsc 红/B4 无测试/B5 TwoStepAction 状态机/B6 confirmed 来源+getLogs 限流)+多 ⚠️。结论：交互/契约设计扎实但非可直接 implement 的工程方案，补 4 块闭环。
