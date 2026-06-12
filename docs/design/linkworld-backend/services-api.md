# 后端子项目 services-api 基线

> 扫描 2026-06-08 | 子项目 backend(2/3) | 已核对真实代码

## 1. services.go —— 核心业务服务（7 个）

| 服务 | 职责 | 关键方法 | 链交互 |
|------|------|----------|--------|
| `UserService` | 用户注册/查询 | `Register(wallet,email,tokenID)`、`GetUser` | 无（纯 DB） |
| `OperatorService` | 运营商列表 | `GetAll`、`GetByID` | 无 |
| `BillingService` | 账单查询/状态 | `GetBills`、`MarkAsPaid(billID,txHash)`、`GetUnpaidBills` | 无，`MarkAsPaid` 仅改 DB `IsPaid` |
| `OracleService`(v1) | 提交用量 | `SubmitUsage(wallet,opID,data,call,sig)` | 无，仅落 `UsageData` |
| `NotificationService` | 账单邮件 | `SendBillEmail` | **TODO stub**，仅 `log.Printf`，未接 SMTP/SendGrid |
| `DepositService` | 押金记账 | `Deposit(wallet,amount)`、`GetDepositAmount`、`RecordWithdraw(wallet,txHash)`、`GetHistory` | 无，押金/提现**仅写 DB 字符串金额**，不上链不校验 |
| `UserServiceService` | 开通/停用服务 | `Activate`、`Deactivate`、`GetUserService` | 无 |

要点：
- `DepositService.Deposit` 把 `amount string` 原样存库（`Type="deposit"`），**不做精度换算、不校验 MIN_DEPOSIT**。
- `RecordWithdraw` 取当前总额写一条 `Type="withdraw"` 记录，**不做链上提现，也不清零**——纯流水追加，逻辑偏弱。
- `BillingService.MarkAsPaid` 信任前端传来的 `txHash`，无链上确认。

## 2. oracle.go —— 计价/用量/虚拟号（4 个组件）

### 2.1 VirtualNumberGenerator
- 按 11 国 `CountryCode` 生成 `+区号+前缀+随机数字` 虚拟号 + 8 位随机密码。
- `GetCountryList()` 返回 code/name/prefix。
- 与链/合约无关，纯工具。

### 2.2 OperatorAPISimulator（**计价模拟器**）
- `GetUsage(userID,opID)` → 随机 data/call 用量。
- `GetBill(userID,opID,month)` → **随机基础金额 + fee**：
  ```go
  baseAmount := rand.Intn(5000)+500       // 随机
  fee := baseAmount * 25 / 1000           // 2.5% 平台费
  return fmt.Sprintf("%d", baseAmount), fmt.Sprintf("%d", fee)
  ```
- 这是后端**目前唯一的「计价」来源**：金额是随机整数字符串，**无精度单位语义**（不是 USDT 6 位、也不是 wei），纯演示数据。

### 2.3 OracleServiceV2（月度出账核心）
- `FetchAndCreateBills()`：遍历所有用户 → 取激活服务 → 调 `OperatorAPISimulator.GetBill` 拿到 (amount,fee) → 写 `Bill{Amount,PlatformFee,IsPaid:false}` 到 DB。
  - **只写 DB，不调任何合约 `createBill` / `monthlySettlement`**。
  - 无分批逻辑，无 oracle 私钥签名上链。
- `FetchUsage(wallet)`：取用户激活服务的模拟用量。
- `SignData(userID,opID,data,call)`：**SHA256(拼接字符串)** 生成「签名」——**不是 ECDSA / EIP-191 / secp256k1**，不是合约可验证的 oracle 签名，仅作占位标识。

### 2.4 UsageService
- `QueryUsage(wallet)`：拉模拟用量 → `SignData` → 落 `UsageData`（带 signature 字符串）。
- `TriggerMonthlyBill()`：转调 `OracleServiceV2.FetchAndCreateBills`，是 `/api/oracle/monthly-bill` 的入口。

## 3. handlers.go 现状

- 单一 `Handler` struct 聚合全部 10 个 service 指针，`NewHandler` 一次性注入。
- 每个 handler 模式统一：`ShouldBindJSON` 校验 → 调 service → `c.JSON`。
- 请求体 binding 用 tag（如 `binding:"required,email"`）。
- 无中间件鉴权（除 CORS）；**oracle/monthly-bill、usage/submit 等敏感端点无任何权限校验**——任何人都能触发月度出账、提交用量。

## 4. 18 个 API 端点职责映射（见 project-scan §3 全表）

按业务域归类：
- **用户域**：register、user/:wallet
- **运营商/服务域**：operators、service/activate、service/deactivate、service/:wallet、virtual-number/generate、countries
- **押金域**：deposit、deposit/:wallet、deposit/:wallet/history、withdraw
- **账单域**：bills/:wallet、bills/pay
- **用量/Oracle 域**：usage/:wallet、usage/submit、oracle/monthly-bill
- **通知域**：notification/send

## 5. 与合约职责的当前落差（详见 alignment-surface.md）

- 「计价」= 随机数模拟器，**无真实资费表、无 USDT 精度**。
- 「月度结算」= 只写 DB，**完全没碰 `Oracle.monthlySettlement` / `Payment.createBill`**。
- 「oracle 签名」= SHA256 字符串，**非链上可验证签名，无私钥来源**。
- 押金/提现/付账 = DB 记账，**无链上写 + 无事件回填闭环**。
