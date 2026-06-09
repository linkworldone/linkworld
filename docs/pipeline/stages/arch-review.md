# Stage: arch-review — 三角色架构审查 + 安全审计（子项目 backend 2/3）

> **状态**: BLOCKED（含 ❌ 阻塞，需返工 design 重跑） | **日期**: 2026-06-09 | **Gate**: 1
> 审查对象：design.md（后端对齐设计）。三角色并行：CEO / Eng / 安全审计（替代 UI Design 路）。冲突优先级：安全/可行性 > 工程。

## 一、总裁决：BLOCK
安全审计 5 ❌ + CEO 1 ❌（与安全 B2 重叠）+ Eng 0 ❌(6 ⚠️)。design 在 owner key 保密、常量时间鉴权、纯整数计价、事件驱动对账方向上判断专业，但**资损「预防控制」系统性缺失**——全篇依赖事后对账告警(检测)而非发交易前硬上限/reorg 确认/单一对账路径(预防)。

## 二、❌ 阻塞项（去重，必须改 design）
| # | 阻塞 | 级别 | 修复方向 |
|---|------|------|----------|
| B1 | owner 单 EOA 持「凭空造单 amounts[] + 扣款」复合权限，无金额护栏（仅 amount>0 挡不住过大；事后告警是检测非预防） | 🔴Critical | 发交易前硬闸：单 user 单 bill 金额上限 MAX_BILL_PER_USER(config 可审计) + 单批总额上限/异常熔断(超历史均值 N 倍人工放行)；owner key 不落 .env 明文长存(启动注入内存)；web 侧 exact-amount approve 禁 infinite（跨 web 3/3，写 handoff-web）；§7.1 措辞升级：owner=deployer=平台 root 权限(能改分账地址/授权拓扑)，非仅结算权限 |
| B2 | `/api/bills/pay`→MarkAsPaid 直接置 IsPaid=true 不验链上 + design §4.3「保留但标注不可信」二选一含保留 = 任意请求伪造已付白嫖 | 🔴Critical | 删除「保留」选项；**IsPaid 唯一由 event_sync 监听 BillPaid 回填**；HTTP 端点删写能力或降级为 pending 意向(绝不置 IsPaid)；§6.6 鉴权清单补 /api/bills/pay；显式声明 web 3/3 接口契约 breaking change 写入 handoff-web |
| B3 | `/api/withdraw`→RecordWithdraw 凭前端任意 txHash 写记账不验链上 + design 把 withdraw 归 AdminAuth（鉴权模型错配：withdraw 是用户操作） | High | withdraw 记账唯一由 DepositWithdrawn 事件回填；HTTP 不接受前端 txHash 记账(最多 pending)；用户写端点鉴权用钱包签名(msg.sender)非 AdminAuth；§6.6 把 withdraw 移出 AdminAuth |
| B4 | usage 无上界 → dataMB×单价 big.Int 不溢出但产天文金额账单；SubmitUsage 接受前端 uint64 无范围校验 | High | PricingService 对 (dataMB,callMin) 设合理上界超界拒绝+告警；amount6 单 bill 硬上限(同 B1 闸)；SubmitUsage gin binding max= 校验 |
| B5 | event_sync 无 reorg/确认数：监听到即置 IsPaid/记押金，块被 reorg 不回退 → 已付/押金虚高资损；(txHash,logIndex) 去重只防重复不防回滚 | High | 资金事件落 confirmed 前等 K 块确认；reorg 检测(last-block 记 blockHash 父哈希不连续→回退重扫)；资金事件 pending→confirmed 两阶段，与 §4.3 状态机统一 |

## 三、⚠️ 需 design/implement 落实（非阻塞）
- **operatorId 固定 ID 映射**（CEO 标最高优先 + Eng#1）：后端 seed 显式写 Operator.ID=链上 operatorId(1..11 合约 initialize 冻结)，**不靠 name 比对**（后端 seed "T-Mobile" vs 链上 "T-Mobile US" 不相等会 miss）；启动期读链 sanity check 不一致 fail/warn。传错=分账打错地址(最隐蔽资损)。
- **event 解码一律走 abigen Filterer**，signatures.go 仅日志用；UsageDataSubmitted 只 user indexed，operatorId/amount 在 data 区(design 多处写法有歧义)；BillCreated 第三参=amount+platformFee 含费总额别当裸 amount。[Eng#2/#5]
- **nonce 本地计数器**(发一笔+1，PendingNonceAt 仅初始化/恢复) + **分批必须等回执 WaitMined** 才判失败批；§6.1 与 §4.1 对"是否等回执"表述不一致需统一。[Eng#3]
- **测试用 go-ethereum simulated.Backend** 部署真实绑定跑分批结算/事件解码，非空 mock RPC。[Eng#4]
- **T1 真前置 = 合约 PR#1 merge(ABI 冻结)**，非"真机上链"；design §7.4 拆开「ABI 冻结→阻塞 T1」vs「真机上链→只阻塞联调」，避免基于未 merge ABI 生成绑定返工。[CEO]
- **验收语言澄清**：本轮=计价引擎正确性，usage 仍 OperatorAPISimulator 模拟(rand)；真实计量是独立后续 Round，§5B 验收不得读成"计费痛点已解决"。[CEO 伪胜利]
- 费率表占位刺眼化(注释 PLACEHOLDER + 启动 warn)+owner+挂真机验收 gate。[CEO]
- usage 单位锁 MB/min 测试固定断言；seed RequiredDeposit 改整数最小单位+存量迁移脚本；string→big.Int 转换校验 ok fail-fast 不静默；chainID 一致校验(transactor chainID==链上)；deployments.json 占位零地址不得用于 event 过滤(零地址拒启动)；结算批次幂等键(month+batchIndex 落 DB 已确认不重发)；弱随机 password 改 crypto/rand；CORS 生产禁 *+credentials；.env gitignore 确认+.env.example 留空。[安全⚠️/Eng]

## 四、✅ 三方认可
- 计价从随机数→可复算整数引擎(纯 big.Int 6 位)、对账从前端自述→链上事件驱动——两个真问题方向解对。
- 后端角色边界收敛漂亮(只代发 monthlySettlement owner，deposit/payBill 归 web 用户侧)，私钥面最小。
- scope 切片有纪律(submitUsage 不做链上签名是聪明 subtraction)；本地 hardhat 解耦让串行不死锁；artifacts 栈式分支可用、abiHash 算法两端一致(可行性验证通过)；go build/vet 现状干净。

## 五、下一步（硬规则：有 ❌ → 改 design → 重跑）
架构师返工 design v2 闭合 B1-B5 + 三 ⚠️ 收紧项 → 重跑 arch-review（Eng+Security 复审）→ 0 ❌ 方可进 plan。

## 六、三方原始摘要
- CEO：1 ❌(B2 对账二选一 one-way door)+ 伪胜利风险 + T1 前置混淆 + owner=root。结论：方向 PASS 但 1 阻塞先消解。
- Eng：0 ❌ DONE_WITH_CONCERNS，6 ⚠️(operatorId/abigen 解码/nonce 回执/simulated.Backend/BillCreated 含费/对账 pending)。现状诊断逐条源码核对属实，可落地。
- 安全：BLOCK 5 ❌(owner 护栏/bills-pay 白嫖/withdraw 伪造/usage 上界/reorg)+9 ⚠️。架构判断专业但资损预防控制系统缺失。
