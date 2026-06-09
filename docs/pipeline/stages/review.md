# Stage: review — 代码审查 + 安全复审（子项目 backend 2/3）

> **状态**: PASS（第1轮 BLOCKED 1 越权 → fix 5d7aea7 → 第2轮双安全 reviewer 0 ❌ 通过） | **日期**: 2026-06-09 | **Gate**: 3
> 审查对象：后端实现代码 diff 68e8382..HEAD。三 reviewer 并行（合约安全审计/代码审查/CSO）+ 修复后双安全复审。

## 一、总裁决：PASS（修复后 0 ❌）
第1轮两个独立安全 reviewer 同时确认 1 个 ❌（WalletAuth 横向越权）→ fixer 修复 → 第2轮双安全 reviewer 逐条核实代码+测试闭合，0 ❌。代码审查第1轮即 0 ❌。

## 二、❌ 阻塞项闭合（R1）
| # | 阻塞 | 修复 | 复审 |
|---|------|------|------|
| R1 | WalletAuth 横向越权(OWASP A01)：中间件验签 signer 并 c.Set(CtxWallet)，但 5 个 handler 从 body 取 wallet 操作他人资源(DeactivateService 即时停用他人服务) | fix 5d7aea7：5 个 handler 改用 c.GetString(CtxWallet) 作操作主体忽略 body wallet；PayBill SetPayIntent owner 基于 authed userID(WHERE id=? AND user_id=?)；FindByWallet 大小写不敏感 | ✅ 合约安全+CSO 双确认闭合 |
| Medium | WalletAuth action 未绑端点(跨端点 nonce 复用) | NewWalletAuth 增 expectedAction，在 ecrecover/消费 nonce 之前校验；main.go 5 路由各绑 action | ✅ 跨 action nonce 拒绝(不消费)闭合 |

新增用例(全绿)：WALLET-AUTHZ-01/02、WALLET-ACTION-01、EmptyAction。最终 79 passing。

## 三、✅ 第1轮三方认可（资损红线逐条核实在代码里）
- owner key 仅内存/不日志/缺失降级/chainID 校验；WalletAuth 一次性 nonce 台账 Consume 原子条件 UPDATE 防并发重放、EIP-712 绑 chainId 非纯 timestamp。
- 三层金额闸 L1/L2/L3 复用单一常量不漂移；冷启动回退无样本不除零仅绝对闸。
- 对账单一路径：IsPaid 唯一 event_sync BillPaid 回填；bills/pay 仅写 PayIntentTxHash；withdraw 仅 pending；GetTotalByUserID 仅计 confirmed。
- event_sync：K 块确认两阶段、blockHash 父哈希断回退、(txHash,logIndex)去重、占位零地址跳过、事件来源合约地址校验。
- 精度 6 位全链路 big.Int、string→big.Int fail-fast；AdminAuth 常量时间比较+缺 key fail；CORS 无 *+credentials；go.mod 固定+verify 通过+无 replace；密钥无硬编码/无误提交/.env gitignore/example 纯占位。
- 代码质量 staff-level、abigen 纯 Go abiHash 复刻、测试断言验行为、simulated.Backend/Cancun 遗留合理。

## 四、⚠️ 非阻塞清理项（后续/随手）
- ActivateServiceRequest.Wallet body 死字段（handler 已忽略）→ 清理。
- 死代码 MarkAsPaid/CreateConfirmed（无调用方）。
- BillCreated 对账校验收窄（processBillCreated 丢弃 totalAmount，design §4.3 要求校验告警）。
- oracle.go SignData(SHA256)残留+Generate 死分支；自建表多处 AutoMigrate 竞态技术债；生产 chainID=0 单链假设；Withdraw 日志 tx_hash 未 %q 转义。
- 补 Activate/Deactivate/Deposit 端点级 e2e 越权回归测试。

## 五、上线前 checklist（继承+后端新增，正式网必做）
- owner key 走 secret manager/KMS；占位常量(MAX_BILL/MAX_BATCH/N/K)拍真值；占位合约地址待合约 1/3 真·上链回填(sync-deployments.sh)。
- 真机端到端：合约上链后 hardhat 31337 或 geth v1.16+ 跑业务合约结算（Cancun 限制，design §7.4 不阻塞本轮）。

## 六、结论
修复后 0 ❌，资损红线代码级闭合 + 越权修复双安全确认，可进 ship。非阻塞清理项与上线 checklist 记录在案。
