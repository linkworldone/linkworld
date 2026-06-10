# Task 06 — T5 WalletAuth signedPost + EIP-712（web 3/3）

> 子项目 web(3/3) | 状态 DONE | 2026-06-10

### 产出文件
- packages/web/src/services/api/signedPost.ts（新建：signedPost(path,body,{wallet,action}) helper + WalletAuthRejectedError + clearWalletAuthSession；GET nonce→signTypedData(EIP-712)→带签名头 POST，不在拦截器调 hook）
- packages/web/src/services/api/signedPost.test.ts（新建 AUTH 用例）
- packages/web/src/services/api/{depositApi,billingApi}.ts（postDepositIntent/postWithdrawIntent/payIntent 改走 signedPost）
- packages/web/src/hooks/useOperator.ts（useApplyNumber /api/service/activate 改走 signedPost）
- 测试同步（depositApi/billingApi/useDeposit.balance）

### git commit
38b9930 feat: web T5 WalletAuth signedPost + EIP-712(对齐后端 middleware) + 会话级签名

### TDD
先红后绿：AUTH 测试先写（signedPost 不存在/EIP-712 字段断言）→ 实现后绿。

### 测试结果
npx tsc --noEmit 0 error。npm run build ✓。npx vitest run：12 files / 49 passed（44 前序+5 新，无回归）。用例 AUTH-01(取 nonce→signTypedData→带头 POST)/AUTH-02(action 绑定对齐后端)/AUTH-03(拒签→不发 POST 抛 WalletAuthRejectedError)/AUTH-04(读端点不带签名头)/一次性 nonce 每次取新。主 Agent 已独立复跑确认 tsc 0+build ✓+49 测。

### code-simplifier
signedPost 单一封装；helper 模块级注入 signer 不在拦截器调 hook；读端点不动。

### spec review
按 design v2 §3.7 + handoff-web + 后端 middleware.go 执行。**EIP-712 逐项对齐后端**：domain(name LinkWorld/version 1/chainId，无 verifyingContract)、primaryType WalletAuth、字段 wallet(address)/nonce(string)/action(string) 同序同类型、header X-Wallet-Address/Nonce/Action/Signature、nonce 路径 /api/auth/nonce/:wallet 一次性消费、action 取值(deposit/withdraw/bills/pay/service/activate·deactivate)。会话级在后端一次性 nonce 下退化为每次取新 nonce 签（按后端实现，注释说明）。未越界（换肤/页面 UI 留 T6+）。

### 设计还原
WalletAuth 签名 UX（拒签态文案「身份签名被取消，操作未提交」就位、与交易签名区分）对齐 design §3.7；视觉/toast 接入待页面改版。

### 复用检查
复用 wagmi signTypedData/getChainId、现有 apiClient；signedPost 供所有写端点复用。

### 设计稿对照
数值对照：EIP-712 domain/字段/header/nonce/action 与后端 middleware.go 8 项逐项一致 ✅；一次性 nonce 每次新签（AUTH 测）✅；读端点不带头(AUTH-04)✅；49 测 ✅；tsc 0/build ✓ ✅。

### 新增组件
新增 signedPost helper + WalletAuthRejectedError + clearWalletAuthSession。

### 新增色值
无（T5 鉴权层，换肤留 T6）。

### ⚠️ 遗留（带入 T7-T10）
- **项目当前无 toast 库**：WalletAuthRejectedError.message 就位，但拒签 UX 需 T7-T10 页面改版时接入实际提示组件（+ShieldCheck/PenLine 图标 + 「请在钱包签名验证身份（不消耗 gas）」前提示，与交易签名视觉区分）；useDeposit/useBilling recordIntent 经 retryWithBackoff 吞错，拒签态需页面侧捕获区分（避免被当网络失败重试存 pendingSync）。
- service/deactivate action 已预留，待 UI 调用点。
- chainId 一致性：signedPost 用 getChainId(wagmiConfig)，端到端真链(D17)需钱包链==后端 walletAuthChainID 否则 ecrecover 失败。
