# Stage: review — 代码审查 + 安全复审（子项目 contracts 1/3）

> **状态**: PASS（三方 0 ❌ 阻塞，可进 ship） | **日期**: 2026-06-09 | **Gate**: 3
> 审查对象：**实现代码** diff 7ef9677..1a516c9（vs arch-review 审设计）。资损敏感，标准从严。
> 三 reviewer 并行：合约安全审计(chain-web3-developer:security-review) / 代码审查(code-reviewer) / CSO 基础设施安全(cso)

## 一、总裁决：PASS，0 ❌，可进 ship
三方一致 0 阻塞。7 大阻塞 B1-B7 + v2-A/B/C + arch-review 5 项 ⚠️ 加固**全部在代码里逐条落地**，实现紧贴 design v2，无遗漏/无超范围。

## 二、合约安全审计（资损，核实"代码真做了"）— 0 ❌
| 审计点 | 结论 | 代码位置 |
|--------|------|----------|
| CEI / 重入 | ✅ | Deposit.withdraw 先清零再 safeTransfer(L75-79)；Payment.payBill 先 isPaid=true(L111) 再两段转账 + nonReentrant(L102)；mintTrafficCard(L86)/issueMonthlyTrafficCards(L134) 带 nonReentrant（A1 闭合）|
| SafeERC20 | ✅ | 全部资金点 safe*，裸 transfer/call{value} 命中 0 |
| 精度 | ✅ | MIN_DEPOSIT=10*10**usdt.decimals()(Deposit:56)；dataAmount 固定 quota，_deposits/100000 已删 |
| 分账零地址 | ✅ | createBill fail-fast(L77) + payBill(L108) + setOperatorPaymentAddress(SM:204) 三处 require |
| 权限链 | ✅ | createBill/issueMonthlyTrafficCards/applyTrafficCardToBill onlyOracle；deploy wiring 后三链可跑通（B2/B3）|
| 自动发卡幂等 | ✅ | 循环 if(!_canMint)continue 不 revert 整批；userCardCount==0 幂等 |
| applyTrafficCardToBill 受限桩 | ✅ | onlyOracle，仅校验+emit，不转资金 |
| A3 _reentrancyGuardInit | ✅ | 已确认删除 |
- ⚠️ Low SEC-001：MockUSDT public mint 无权限（测试制品，上线门禁须排除 mock）。
- 报告：docs/security-review/security-cr-20260609-contracts-erc20-usdt-1a516c9.md

## 三、代码审查（4 维度）— 0 ❌
- Plan alignment ✅：B1-B7 + v2-A/B/C + 5 加固代码级落地，紧贴 design v2，无遗漏/超范围。
- Code quality ✅：_mintFor 复用消重、错误信息清晰、Solidity 风格一致。
- Architecture ✅：接口 SSOT 收敛、initialize 注入、UUPS 追加规则、事件设计合理。
- Testing ✅：75 passing，断言真验行为（非走过场）。
- ⚠️ 非阻塞：① Oracle `_monthlyUsage`/`_latestUsage` 经 v2-A 后成死代码（已开 spawn chip 供单独清理）；② 测试走 Factory.deploy() 非 proxy 路径，但 deploy.ts proxy 冒烟兜底。

## 四、CSO 基础设施安全 — 0 ❌
- 密钥 ✅：无硬编码私钥/RPC token；git 历史无 .env 误提交；.gitignore 覆盖 .env*；.env.example 纯占位；新网络 accounts 无 key 时空数组（不 fallback 默认 key）。
- 依赖 ✅：本 diff 无新增依赖；OZ ^5.6.1 + lockfile 兜底。
- handoff 文档 ✅：无敏感信息泄漏。
- ⚠️ 非阻塞观察：deployments/ 与 .openzeppelin/ 被 git 追踪（公开信息，OZ 推荐入库，非泄漏）。

## 五、⚠️ 上线前 checklist（本测试网轮次不阻塞，下游/后续 Round 必做）
1. 所有合约 owner + UUPS upgrade 权限 + platformWallet 从 deployer EOA **转多签**。
2. 11 个 operator paymentAddress（当前 deployer 派生的黑洞地址，不可提取）**换真实运营商地址**。
3. 移除未使用的 `axios` 依赖（dead dep）。
4. hardhat.config og_* 网络的全 0 私钥 fallback 统一为空数组模式。
5. 上线门禁排除 mock 制品（MockUSDT/NonStandardUSDT）。
6. 清理 Oracle 死 state（_monthlyUsage 等，已开 chip）。

## 六、已知限制（不阻塞）
- Arbitrum Sepolia 421614 真·上链 + TSTORE 链上实测：本环境无 DEPLOYER_PRIVATE_KEY/RPC，未执行。本地 31337/hardhat(cancun, transient) 全绿、配置就绪。记 handoff §10，ship 阶段若配 key 可补做。

## 七、结论
三方 0 ❌，合约实现代码资损面扎实、紧贴设计、测试充分。**可进 ship。** 上线前 checklist（§五）为正式网/下游事项，本测试网子项不阻塞。
