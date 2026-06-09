# LinkWorld 后端(2/3) → Web(3/3) Handoff 清单

> 产出阶段：后端子项目(2/3) design v2（按 arch-review 返工）
> 依据：design.md §4.3 / §6.6 / §7.0 / §12；arch-review.md B1/B2/B3
> 状态：本轮后端对齐引入 **接口契约 breaking change** + **资损预防硬约束**，web 子项目(3/3) 必须据此对齐。
> 移交方向：backend 2/3 → web 3/3。后端 design 已锁定，web 不得反向要求后端置终态。

---

## 1. 接口契约 breaking change（必读，B2/B3）

后端对账模型从「前端自述即记账」改为「链上事件唯一回填终态」。受影响端点：

### 1.1 `POST /api/bills/pay`（语义变更）

| 项 | v1（旧） | v2（新） |
|----|---------|---------|
| 后端动作 | 直接置 `Bill.IsPaid = true` | **只写 pending 意向**（`PayIntentTxHash`），**不动 IsPaid** |
| IsPaid 来源 | 该端点 | **唯一**由后端 event_sync 监听链上 `BillPaid` 事件回填（等 K 块确认） |
| 鉴权 | 无 | **钱包签名（WalletAuth）** |

**web 必做**：
- 用户支付后调本端点仅作「我已发起支付」意向上报；**不能据 200 响应就显示「已付」**。
- 真实「已付」状态：轮询 `GET /api/bills/:wallet` 看 `is_paid`，或等后端确认。中间态显示「支付确认中」。
- 请求需带钱包签名头（见 §3）。

### 1.2 `POST /api/withdraw`（语义变更）

| 项 | v1（旧） | v2（新） |
|----|---------|---------|
| 后端动作 | 凭前端 `tx_hash` 直接写 withdraw 记账 | **不接受 tx_hash 记账**；最多写 pending 意向（不计入余额） |
| 记账来源 | 该端点 | **唯一**由 event_sync 监听 `DepositWithdrawn(user,principal,interest)` 回填 |
| 鉴权 | 错归 AdminAuth | **钱包签名（WalletAuth）**（withdraw 是用户操作） |

**web 必做**：
- 提现仍是 web 侧链上操作（用户钱包调 `Deposit.withdraw()`），流程不变。
- 提现历史/余额以后端**事件确认**为准，不以前端 tx_hash 立即反映。
- 请求需带钱包签名头。

### 1.3 `POST /api/deposit`（鉴权新增）

- 充值仍是 web 侧链上操作（用户钱包调 `Deposit.deposit(amount)`，前置 `usdt.approve`）。
- 后端端点仅记 pending 意向；余额以 `DepositMade` 事件确认为准。
- 请求需带钱包签名头。

---

## 2. exact-amount approve 禁 infinite（资损硬约束，B1）

> owner=deployer 是**平台 root 权限**（能改分账地址 / 授权拓扑 / 凭空造单）。一旦 owner key 被盗用或合约被攻破，用户对合约的 **infinite approve 会被无限抽取**。

**web 必做**：
- `usdt.approve(spender, amount)` **必须按本次精确金额授权**：
  - 充值前：`approve(deposit, amount)`。
  - 付账前：`approve(payment, amount + platformFee)`（精确总额）。
- **禁止 `approve(spender, MaxUint256)` / infinite approve**。
- 每次操作前重新 approve 精确额度；操作后授权应归零或仅余精确额度。
- 精度：amount 用 USDT 6 位最小单位（`deployments.<net>.json.usdtDecimals`，勿硬编码 6）。

---

## 3. WalletAuth 钱包签名鉴权（新增）

用户资金/写端点（bills/pay、withdraw、deposit、service 写）需带钱包签名头：

- web 用用户钱包私钥对「请求体 + nonce/timestamp」签名。
- 后端 `ecrecover` 还原签名地址，必须 == 请求 `wallet` 字段（绑定 msg.sender 语义），否则 401。
- 防重放：带 nonce 或时间窗 timestamp。
- 具体签名消息格式（字段顺序 / EIP-191 or EIP-712）由 implement 阶段双方敲定，本文档锁定「需要钱包签名」这一契约。

---

## 4. deposit/withdraw 状态机（pending → confirmed）

```
[无记录] ──HTTP 意向──▶ [pending] ──event_sync 监听──▶ [seen] ──等 K 块──▶ [confirmed]
                       (不计入余额)                                          (计入余额)
```

**web 展示规则**：
- 余额计算只认 `confirmed` 记录（后端 `GetTotalByUserID` 同口径）。
- pending/seen 记录可显示为「处理中」，不参与可用余额。
- reorg 时后端会回退未确认（< K 块）记录，web 不应缓存未确认态当真。

---

## 5. 金额精度（全链路一致）

- 全链路 USDT **6 位最小单位**（链上 uint256 / 后端 big.Int / DB 字符串）。
- web 展示除 `10^usdtDecimals`；`usdtDecimals` 从 `deployments.<net>.json` 读，**不硬编码 6**。
- 现有 `web/src/services/api/billingApi.ts` 的 `parseFloat` 金额计算需复核 6 位精度语义（避免浮点误差），但金额加减建议用 bigint。

---

## 6. 未变更（web 可沿用）

- `GET /api/bills/:wallet`、`GET /api/deposit/:wallet`、`GET /api/usage/:wallet` 等读端点契约不变（公开，无需鉴权）。
- register、operators、virtual-number 等流程不变。
- 链上合约调用方（deposit/payBill/withdraw 由用户钱包发起）不变——后端不代发这三者。
