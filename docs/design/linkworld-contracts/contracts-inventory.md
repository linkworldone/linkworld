> 扫描时间：2026-06-08 | 子项目：contracts(1/3) | 对象：packages/contracts | 已核对真实代码

# 合约清单（contracts-inventory）

7 个合约全部 UUPS 可升级（`OwnableUpgradeable + UUPSUpgradeable`，`initialize()` 初始化，`_authorizeUpgrade onlyOwner`）。

## 总览表

| 合约 | 职责 | 可升级 | 依赖 | 本轮 |
|------|------|--------|------|------|
| **Deposit** | 用户保证金存/取，锁仓 30 天，到期 mint 流量卡 | UUPS | IUserRegistry, ITrafficCardNFT, oracle(addr) | **改写核心** |
| **Payment** | 账单创建 + 支付 + 平台手续费分账 | UUPS | IFeeManager | **改写核心** |
| **FeeManager** | 手续费率（基点）读写与计算 | UUPS | 无 | 事实源，基本不改 |
| **TrafficCardNFT** | 流量卡 ERC721：mint→持有→burn→30天服务 | UUPS | 无（ERC721URIStorage） | 时序相关 |
| ServiceManager | 运营商目录（11 个内置）增删查 | UUPS | 无 | 本轮基本不改 |
| Oracle | 月末结算 + 服务有效性校验 | UUPS | deposit(addr), payment(addr) | 受 Payment/Deposit 改动牵连 |
| UserRegistry | 邮箱注册 + 身份 ERC721 | UUPS | 无 | 本轮不改（R10 邮箱注册保留） |

---

## ① Deposit.sol（改写核心）

**职责**：保证金存款（原生币 payable，锁仓 30 天）、到期提取、到期 mint 流量卡。

**关键 state**：`_deposits[addr]`、`_lockExpiry[addr]`、`trafficCardQuota`(默认 100MB)、`userRegistry`、`trafficCardNFT`、`oracle`(address)、`_operatorRequiredDeposit`(声明但未使用)。

**关键函数签名**：
| 函数 | 当前签名 | 说明 |
|------|----------|------|
| `initialize` | `(address _userRegistry)` | quota=100MB |
| `deposit` | `() external payable` | **payable，本轮改 ERC20** |
| `withdraw` | `() external` | `.call{value:principal}` 退原生币，**本轮改 ERC20 transfer** |
| `mintTrafficCard` | `(address user) onlyOwner returns(uint256)` | 锁仓到期 + 有存款 + 该用户无卡 → mint；`dataAmount = _deposits/100000` |
| `generateTokenURI` | `(address) view returns(string)` | 拼 api.linkworld.io URI |
| `getDepositAmount` | `(address) view returns(uint256)` | |
| `getLockExpiry` | `(address) view returns(uint256)` | 前端倒计时事实源（R6） |
| `issueMonthlyTrafficCards` | `(address[]) external` | 仅 oracle 可调；**函数体为空（仅注释 TODO）** |
| `setOracle/setTrafficCardNFT` | `(address) onlyOwner` | |

**事件**：`DepositMade(user,amount)`、`DepositWithdrawn(user,principal,interest)`、`TrafficCardMinted(user,tokenId,dataAmount)`（继承自 IDeposit）。

**锁仓逻辑**：`deposit()` 时若上次锁仓已过期则 `now+30days`，否则在原到期点 `+30days`（叠加续期）。**当前无最小额约束**（仅 `require(msg.value>0)`），本轮需加 `>= 10 USDT`。

---

## ② Payment.sol（改写核心）

**职责**：创建账单（含手续费）、用户支付（原生币 payable）、平台分账。

**关键 state**：`feeManager`、`platformWallet`、`oracle`、`_nextBillId`、`_bills[id]`、`_userBillIds[addr]`。继承 `ReentrancyGuardTransient`（transient storage，依赖 cancun）。

**关键函数签名**：
| 函数 | 当前签名 | 说明 |
|------|----------|------|
| `initialize` | `(address _feeManager, address _platformWallet)` | 含自定义 `_reentrancyGuardInit()`（assembly 写 slot=1） |
| `createBill` | `(address user, uint256 operatorId, uint256 amount) onlyOwner` | 调 `feeManager.calculateFee(amount)` 算手续费 |
| `payBill` | `(uint256 billId) external payable` | **payable；`msg.value>=amount+fee`；多退少补；fee 转 platformWallet** |
| `getUserBills` | `(address) view returns(Bill[])` | |
| `getUnpaidBills` | `(address) view returns(Bill[])` | |
| `setOracle/setPlatformWallet/setFeeManager` | onlyOwner | |

**Bill 结构**：`{id, user, operatorId, amount, platformFee, createdAt, isPaid}`。

**事件**：`BillCreated(billId,user,totalAmount,platformFee)`、`BillPaid(billId,user,totalAmount,operatorAmount)`。

> ⚠️ **缺口**：当前 `payBill` 仅把 `platformFee` 转给 platformWallet，**`bill.amount`（运营商应得部分）没有转给任何运营商地址**——留在合约里。本轮 ERC20 改写需明确运营商分账去向（ServiceManager.paymentAddress 当前全是 `address(0)`）。
> ⚠️ `IPayment` 接口/Oracle 内联接口里出现 `applyTrafficCardToBill(billId)`，**但 Payment.sol 并未实现该函数**（见 erc20-migration-surface §一致性风险）。

---

## ③ FeeManager.sol（手续费事实源，R5）

**常量**：`FEE_DENOMINATOR = 10000`、`MAX_FEE_RATE = 1000`(10%)。
**state**：`_feeRate`（基点，部署初始化为 `150` = 1.5%）。
**函数**：`getFeeRate() view→uint256`、`setFeeRate(uint256) onlyOwner`(≤1000)、`calculateFee(uint256 amount) view→ amount*_feeRate/10000`。
**事件**：`FeeRateUpdated(oldRate,newRate)`。

→ 已满足 R5「动态基点 150 / 分母 10000」，前端直接读 `getFeeRate()`/`calculateFee()`。本轮基本不改，但需注意 USDT 6 位精度下小额账单 `calculateFee` 的取整（截断为 0）。

---

## ④ TrafficCardNFT.sol（流量卡，时序相关）

**ERC721URIStorageUpgradeable**，name `LinkWorld Traffic Card` / symbol `LWTC`。

**关键 state**：`_nextTokenId`、`_cardInfo[tokenId]{dataAmount,createdAt,isDestroyed}`、`_userCardCount[addr]`、`depositContract`、`_deductionCredits[addr]`(声明)、常量 `DEDUCTION_VALIDITY = 30 days`。

**关键函数**：
| 函数 | 签名 | 权限 | 说明 |
|------|------|------|------|
| `mint` | `(address to, uint256 dataAmount, string tokenURI_) returns(uint256)` | **onlyOwner** | 部署后 ownership 转给 Deposit，故实际由 Deposit 调 |
| `mintBatch` | `(address[], uint256[], string[]) returns(uint256[])` | onlyOwner | 循环调 `this.mint` |
| `burn` | `(uint256 tokenId)` | `_isAuthorized` 持有人/授权 | 标记 isDestroyed + `_burn` + emit `ServiceActivated(now+30days)` |
| `getCardInfo` | `(uint256) view returns(CardInfo)` | | |
| `getUserCardCount` | `(address) view returns(uint256)` | | |

**事件**：`CardMinted`、`CardDestroyed`、`CreditUsed`、`CreditExpired`（合约内额外声明）、`ServiceActivated`（来自接口）。

> 注：`onlyDeposit` modifier 与 `_deductionCredits` 已声明但**未被任何函数使用**（dead code / 预留）。`mint` 用 `onlyOwner`，靠 deploy 时 `transferOwnership(Deposit)` 授权 Deposit。

**Deposit↔发卡时序**：`Deposit.mintTrafficCard(user)`（onlyOwner，锁仓到期手动调）→ `TrafficCardNFT.mint`。另有 `Deposit.issueMonthlyTrafficCards`（仅 oracle，**空实现**）。

---

## ⑤ ServiceManager.sol（运营商目录，本轮基本不改）

`initialize()` 内置 11 个运营商（id 1~11，US/GB/FR/RU/JP/VN/LA/KH/TH/MY/PH），`requiredDeposit` 用 `0.0xx ether`（原生币计价 ⚠️ 切 USDT 后该字段单位语义需复核）。
函数：`addOperator/updateOperator/deactivateOperator`(onlyOwner)、`getOperator/getActiveOperators/getOperatorsByCountry`(view)。
**所有 `paymentAddress` 当前为 `address(0)`** —— 影响 Payment 运营商分账。
事件：`OperatorAdded/OperatorUpdated/OperatorDeactivated`。

---

## ⑥ Oracle.sol（计量/结算）

**state**：`deposit`(address)、`payment`(address)、`_monthlyUsage`、`_latestUsage`、内联 `UsageInfo` struct。
**函数**：`initialize()`、`setDeposit(address)`、`verifyServiceActive(address) view`（staticcall `getLockExpiry` 判活）、`monthlySettlement(users[],operatorIds[],dataUsages[],callUsages[]) onlyOwner`、`submitUsage(...) view`（预留空壳）。
**文件尾部内联了 `IPayment`/`IDeposit` 两个接口**（与 interfaces/ 下的版本并存，含 `applyTrafficCardToBill`）。

> ⚠️ **一致性风险（高优先级）**：`monthlySettlement` 里 `payment` / `deposit` 声明为 `address` 类型，却直接 `payment.createBill(...)`、`deposit.issueMonthlyTrafficCards(...)`、`payment.applyTrafficCardToBill(...)` 调用——`address` 类型无这些方法，按当前源码**应无法编译**（疑似与 artifacts 不同步，artifacts 可能是旧版本）。且 `setPayment` 缺失（无法设置 payment 地址）。本轮 arch-review 需确认该文件真实可编译状态。

---

## ⑦ UserRegistry.sol（注册，本轮不改）

`ERC721URIStorageUpgradeable`，name `LinkWorld Identity` / symbol `LWID`。
`register(string email)` → mint 身份 NFT + 存 `UserInfo{wallet,email,tokenId,isActive,registeredAt}`；`getUserInfo`、`isRegistered`。
事件 `UserRegistered(user,email,tokenId)`。**无 phone 字段**（对齐 PRD R10：本轮保留邮箱注册）。
