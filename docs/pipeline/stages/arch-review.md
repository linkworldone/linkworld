# Stage: arch-review — 三角色架构审查 + 合约安全审计（子项目 contracts 1/3）

> **状态**: BLOCKED（本轮 design.md 未通过，含 ❌ 阻塞项，需返工 design 后重跑） | **日期**: 2026-06-08 | **Gate**: 1
> 审查对象：docs/pipeline/stages/design.md（合约技术设计，commit 9584e4a）
> 三角色并行：CEO(plan-ceo-review) / Eng(plan-eng-review) / Security(chain-web3-developer:security-review，替代对合约无意义的 UI design-review)
> 冲突优先级：安全/可行性 > 工程实现 > 视觉

## 一、总裁决：BLOCK — 不可进 plan/implement
三方合计 **7 大 ❌ 阻塞主题**（去重后），其中 B1 计价模型为 🔴 Critical 资损。design 工程纪律 staff 级（逐行核对、实测编译、决策对比表），但存在结构性资损/可行性缺口，必须返工。

## 二、❌ 阻塞项（去重，必须改 design 才能继续）

| # | 阻塞项 | 来源 | 级别 | 修复方向 |
|---|--------|------|------|----------|
| B1 | 账单金额无计价模型：Oracle `monthlySettlement` L68 把 dataUsage(字节)+callUsage(分钟) 直接相加当 USDT 金额传 createBill，ERC20 真实扣款会扣天文数字、量纲错误 | 安全❌-4 | 🔴Critical | 【用户决策待填】后端/Oracle 链下按资费算好金额传入、合约 createBill 只存不算（去掉 usage 求和）；或链上计价表 |
| B2 | createBill 权限链断裂：createBill 是 onlyOwner，Oracle 非 Payment owner（deploy 无 payment.transferOwnership(oracle)），monthlySettlement 调 createBill 必 revert，自动结算跑不通 | 安全❌-3 | High | createBill 改 onlyOracle（或 wiring 转 owner）；§7.2 补授权步骤；design 增"部署后 ownership/授权拓扑图" |
| B3 | 自动发卡权限链+幂等：issueMonthlyTrafficCards(onlyOracle) 若复用 mintTrafficCard(onlyOwner) 内部 mint 会 revert；用 require 会让整批回滚 | 安全❌-2 / Eng | High | 走独立 internal `_mintFor(user)` 不经 onlyOwner；批量循环用 `if(!ok) continue;` 跳过不满足者，不 revert 整批 |
| B4 | dataAmount 精度耦合：mintTrafficCard L76 `_deposits/100000` 在 6 位精度下发废卡；design §5.1 只写"建议"未锁定，§一表口径矛盾 | CEO❌ / 安全❌-1 / Eng⚠️ | High | 【用户决策待填】锁定 dataAmount 来源（固定 quota 或按保证金阶梯重定义换算）；删 _deposits/100000；统一全文口径 |
| B5 | 接口收敛漏声明：§5.6⑤ 要把 Oracle 内联接口收敛到 interfaces/，但 IDeposit.sol 无 issueMonthlyTrafficCards、IPayment.sol 无 applyTrafficCardToBill 声明，收敛后编译不过 | Eng❌-1 | High | IDeposit 补 issueMonthlyTrafficCards、IPayment 补 applyTrafficCardToBill 声明；把"接口签名变更→实现同步"列进 §5.6 连锁编译清单 |
| B6 | 核心新函数零测试：Oracle.monthlySettlement + Payment.applyTrafficCardToBill（本轮新实现的集成路径）design §八测试清单完全未覆盖 | Eng❌-2 | High | 补 monthlySettlement 端到端集成测试 + applyTrafficCardToBill 抵扣测试 |
| B7 | 非标代币边界：未处理 fee-on-transfer（Deposit 记账虚高→withdraw 直接资损）+ 真实 USDT 黑名单/暂停导致 payBill 失败 | 安全❌-5 | High | design §6.1 补：仅支持标准 ERC20、禁 fee-on-transfer；分账地址非黑名单要求 + 支付失败降级策略 |

## 三、⚠️ 需 design 拍板消歧 / 加固项（不单独阻塞但应一并处理）

- usdt 注入路径：锁定为 **initialize 注入**（fresh deploy 无升级包袱），弃 setUSDT 后置（避免 usdt 未设时 deposit panic）。【Eng⚠️-6】
- P1 compile gate 重定义：design §5.6 "①②④后 compile 应绿"不成立（②Oracle 类型改完缺接口方法、③未实现）；应为"P1 全部子项①-⑤做完后 compile 绿"。【Eng⚠️-5】
- 批量 mint / monthlySettlement 无界循环 gas：无上限+无测试+后端静默调用（临界缺口）。design 明确分批策略（调用方分批 ≤N，或合约 require maxBatch）+ 补压测。【Eng⚠️-10】
- applyTrafficCardToBill 抵扣范围：现状只抵第一张未付账单（Oracle L85），单张/全部/按额度语义未定。【Eng⚠️-11 / 关联 B 决策】
- ServiceManager requiredDeposit：18 位 ether 字面量在 USDT 语境语义错，_operatorRequiredDeposit 当前未被使用。design 明确"本轮不引入运营商保证金校验/标记废弃"。【CEO⚠️ / Eng】
- 16602/localhost 旧 proxy：design §2.2/§6.3/§8 升级路径尾巴矛盾。明确"本轮一律 fresh deploy 不升级"，删除 storage layout 升级讨论的无效负担。【Eng⚠️】
- ReentrancyGuardTransient：§6.4 Arbitrum 兼容实测应**前移到 P2 之前**（而非 P5），避免返工；且 Payment L33-38 自定义 _reentrancyGuardInit assembly 本身正确性可疑（向 transient slot 写持久 storage），需厘清而非仅"降级时删"。【Eng⚠️-15 / 安全 6】
- createBill fail-fast：建议 createBill 也校验 operator.paymentAddress != 0，避免生成永远付不了的账单。【安全 4】
- 0-fee 跳过：payBill 第二段 fee=0 时跳过（保留现状 Payment L86 `if(platformFee>0)` 判断）。【Eng⚠️-9 / 安全 9】
- mintBatch 禁用：自动发卡禁止用 NFT.mintBatch（L77 `this.mint()` 外部调用会撞 onlyOwner）。【Eng⚠️-4】
- 交付物补充：合约子项验收应含"冻结的 ABI + selector 清单 + 金额精度语义说明"作为给后端子项(2/3)的正式 handoff（deposit/payBill 改 selector）；本轮 storage layout 冻结并记入 deployments.json。【CEO⚠️】
- 测试补充：非标 USDT(transfer 不返回值/返回 false) 用例验证 SafeERC20 分支；锁仓续期不变量、未注册 deposit revert（旧测试用 {value} 调用，改 deposit(amount) 后失效需重写）。【CEO⚠️ / Eng⚠️-14】

## 四、✅ 三方认可的设计优点
- 修编译设为强依赖拓扑根 + compile P1 gate：排序正确。
- 迁 Arbitrum + ERC20 USDT 对齐跨境/无 KYC 产品定位，币种对齐用户心智货币。
- 3 子项串行、合约作地基先行成立，无反向依赖。
- SafeERC20 选型、CEI 方向、授权精确额、分账零地址 payBill 兜底、fresh deploy storage layout 判断、Arbitrum"先实测再定"姿态——显性资损面处理扎实。
- 补测清单对验收主路径（A.2 最小额/approve/分账/发卡幂等/精度）覆盖完整。

## 五、下一步（硬规则：有 ❌ → 改 design → 重跑 arch-review）
1. 用户已就 B1(计价模型)/B4(发卡额度)/applyTrafficCardToBill 范围 三项产品决策拍板（结论回填本报告 B1/B4 + §三）。
2. 架构师按本报告返工 design.md（闭合 B1-B7 + §三消歧项）。
3. 重跑 arch-review（至少 Eng + Security 复审），0 ❌ 方可标完成、进 plan。

## 六、三方原始意见存档摘要
- CEO：❌1（B4 dataAmount 语义未拍板）+ ⚠️5 + ✅8。结论：方向对、工程扎实，唯一硬伤是"规范真空"——核心发卡功能正确性留给 implement。
- Eng：❌2（B5 接口收敛漏声明 / B6 核心函数零测试）+ 多个二义点 + 临界 gas 缺口。结论：不可直接进 implement，修 2 阻塞+拍板 4 二义点后 P1→P6 可落地。
- Security：BLOCK，❌5（B1计价Critical / B2 createBill权限链 / B3发卡权限链幂等 / B4 dataAmount精度 / B7非标代币）+ ⚠️9 + ✅4。结论：显性资损面扎实，但自动结算两条权限链 deploy 下跑不通 + 计价模型缺失是结构性缺口。
