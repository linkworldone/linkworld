# Task 07 — T7 鉴权 + 端点加固 + main 装配（后端 2/3）

> 子项目 backend(2/3) | 状态 DONE | 2026-06-09

### 产出文件
- packages/backend/internal/middleware/middleware.go（新：NewAdminAuth ConstantTimeCompare 缺 key error；WalletAuthDigest EIP-712 TypedData 绑 name/version/chainId；NewWalletAuth ecrecover 绑 wallet + 消费 nonce）
- packages/backend/internal/middleware/*_test.go（ADMIN-01/WALLET-01~04）
- packages/backend/internal/models/models.go（新增 WalletNonce 一次性台账：Wallet+Nonce 唯一+Used+ExpiresAt）
- packages/backend/internal/repository/repository.go（WalletNonceRepository：Issue crypto/rand 签发 + Consume 条件 UPDATE 原子消费防并发重放 + Migrate；BillRepository.SetPayIntent）
- packages/backend/internal/services/services.go（RecordPayIntent B2；Deposit Status=pending；RecordPendingWithdraw B3 废弃凭 txHash 记账；MarkAsPaid 标注仅 event_sync 内部）
- packages/backend/internal/handlers/handlers.go（PayBill→RecordPayIntent、Withdraw→RecordPendingWithdraw、SubmitUsage binding max=、TriggerMonthlyBill→SettleMonthlyOnChain 返 SettlementSummary、新增 GetWalletNonce 公开端点）+ handlers_test.go
- packages/backend/cmd/main.go（nonceRepo Migrate；AdminAuth 缺 key fail-fast；WalletAuth 绑 chainID；结算编排器条件装配；CORS 收口放行鉴权头无 *；按 §6.6 挂中间件）

### git commit
ef8513e feat: 后端 T7 鉴权(AdminAuth+WalletAuth nonce台账)+端点加固+bills/pay降级+main.go装配+CORS收口

### TDD
先红后绿：①middleware 测试先写→undefined(WalletNonce/NewAdminAuth/WalletAuthDigest) 编译红→实现绿；②handler 测试先写→PAY-INTENT-01 IsPaid=true 断言红 + USAGE-01 超 max got 200 红→改 service/handler 绿。

### 测试结果
go build ./... 0 error。go test ./... 全绿（config/handlers/middleware/services/sync 全 ok，T1-T6 无回归）。用例：AdminAuth(缺key启动fail/正确过/错误·缺header 401)、WalletAuth(ValidSignature过/**ReplayNonce 拒绝 🔴/**WrongChainID/WalletMismatch/UnknownNonce 拒绝)、PayBill_WritesIntent_NotPaid、Withdraw_PendingIntent_NotCountedInBalance、SubmitUsage_OverMax_400/WithinMax_OK。主 Agent 已独立复跑确认 build exit 0 + 全包 test ok。

### code-simplifier
中间件单一职责；nonce 消费用条件 UPDATE 原子化（无锁）；复用 crypto/subtle/go-ethereum crypto。

### spec review
按 design v2 §6.6/§4.3/§4.4 + arch-review N1 红线/B2/B3/CORS 执行。**N1：服务端一次性 nonce 台账(消费式 Used 标记，签过即作废，非纯 timestamp)+EIP-712 绑 chainId/domain**——WALLET-02 重放测拒绝验证。未越界（T8 测试扫尾留 T8；client/event_sync/settlement 核心只调用）。

### 设计还原
后端无 UI。design §6.6 鉴权 + §4.3 对账降级逐项落地：端点鉴权清单挂载、bills/pay 写 PayIntentTxHash 不置 IsPaid、withdraw 不凭 txHash 记账。

### 复用检查
复用 go-ethereum crypto(Ecrecover/SigToPub)、crypto/subtle ConstantTimeCompare、T4 的 IsPaid 唯一事件回填(MarkPaidByOnChainID)、T6 的 SettleMonthlyOnChain；MarkAsPaid 不导出 HTTP 仅 event_sync 备用。

### 设计稿对照
数值对照：WalletAuth 一次性 nonce 重放拒绝(WALLET-02)vs N1 红线 ✅；EIP-712 绑 chainId(WALLET-03 错 chainId 拒)vs N1 ✅；bills/pay IsPaid 仍 false(PAY-INTENT-01)vs B2 ✅；withdraw pending 不计入余额(WD-01)vs B3 ✅；usage 超 max=1000000/100000→400(USAGE-01)vs B4 ✅；CORS 无 * 固定白名单 ✅；端点鉴权 18 端点全覆盖 ✅；go build 0 ✅。

### 新增组件
新增 middleware(AdminAuth/WalletAuth/WalletAuthDigest)、WalletNonce 模型、WalletNonceRepository、GetWalletNonce 端点、RecordPayIntent/RecordPendingWithdraw/SetPayIntent。无新增业务合约。

### 新增色值
无（后端任务）。

### ⚠️ 遗留（带入 T8）
- T8 端到端可用 middleware.WalletAuthDigest 构造签名头驱动受保护端点。
- MarkAsPaid 现无 HTTP 调用方（保留供 event_sync 备用），T8 可决定是否清理。
- WalletNonce 仅 Used 标记+10min 兜底过期，未加过期清理 job（不影响正确性，台账会增长）；后续可加。
- SettlementSummary=null 降级分支序列化 T8 注意。
