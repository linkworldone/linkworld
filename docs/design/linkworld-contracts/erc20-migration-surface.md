> 扫描时间：2026-06-08 | 子项目：contracts(1/3) | 对象：packages/contracts | 已核对真实代码

# ERC20 改写面（erc20-migration-surface）

本轮核心：`Deposit`/`Payment` 从原生币（`payable`/`msg.value`）改为 ERC20 USDT（`approve` + `transferFrom`）。本文列出所有触及点，给 plan/implement 阶段当改写清单。

## 一、所有 `payable` / `msg.value` / `.call{value}` 出现位置（grep 实测）

全仓命中 **2 个合约 + 2 个接口**，无其它残留：

```
contracts/Deposit.sol:41    function deposit() external payable {
contracts/Deposit.sol:42        require(msg.value > 0, "Zero deposit");
contracts/Deposit.sol:45        _deposits[msg.sender] += msg.value;
contracts/Deposit.sol:52        emit DepositMade(msg.sender, msg.value);
contracts/Deposit.sol:64        (bool success, ) = msg.sender.call{value: principal}("");   // withdraw 退款
contracts/Payment.sol:76        function payBill(uint256 billId) external payable {
contracts/Payment.sol:82        require(msg.value >= total, "Insufficient payment");
contracts/Payment.sol:87        ... platformWallet.call{value: bill.platformFee}("");        // fee 分账
contracts/Payment.sol:91-92    if (msg.value > total) { msg.sender.call{value: msg.value-total} }  // 找零退款
contracts/interfaces/IDeposit.sol:9    function deposit() external payable;
contracts/interfaces/IPayment.sol:19   function payBill(uint256 billId) external payable;
```

✅ 改写面收敛、规模小（2 合约）。无 ERC20/IERC20/SafeERC20 现有引用（grep 零命中），需全新引入。

## 二、Deposit.sol 现状

**存款** `deposit() payable`：`require(msg.value>0)` + `require(isRegistered)` → 累加 `_deposits[msg.sender] += msg.value` → 锁仓续期（首存 `now+30d`，复存在原到期点叠加 `+30d`）→ emit `DepositMade(sender,msg.value)`。

**提取** `withdraw()`：`require(now>=_lockExpiry)` + `require(_deposits>0)` → 清零本金与锁仓 → `msg.sender.call{value:principal}` 退原生币 → emit `DepositWithdrawn(sender,principal,0)`（interest 恒 0）。

**`_lockExpiry` 锁仓机制**：30 天，叠加续期；前端读 `getLockExpiry(addr)` 做倒计时（R6）。改 ERC20 时**锁仓逻辑可保持不变**，只换资金通道。

**最小额约束现状**：**无**。仅 `msg.value>0`。本轮需加 `require(amount >= MIN_DEPOSIT)`，MIN_DEPOSIT = 10 * 10^(USDT decimals)（验收 A.2 要求 <10 拒绝、=10 通过的单测）。

**改写要点**：
- `deposit()` 去 payable，改 `deposit(uint256 amount)`，前置 `usdt.safeTransferFrom(msg.sender, address(this), amount)`，需用户先 `approve`。
- `withdraw()` 改 `usdt.safeTransfer(msg.sender, principal)`。
- 新增 `IERC20 usdt` state + initialize 注入（注意 UUPS storage layout：只能追加到末尾）。

## 三、Payment.sol 现状

**支付** `payBill(billId) payable`：校验未付/是本人 → `total = amount + platformFee` → `require(msg.value>=total)` → 标记已付 → fee 转 platformWallet → 多余找零退回。

**手续费算法**：`createBill` 阶段调 `feeManager.calculateFee(amount)`（=`amount*feeRate/10000`），存入 `bill.platformFee`。**已用 FeeManager**（满足 R5）。

> ⚠️ **运营商分账缺口**：`payBill` 只转了 platformFee 给 platformWallet，`bill.amount`（运营商应得）留在合约里没有出口。ERC20 改写时需明确：amount 转给谁（运营商 paymentAddress 当前全 0）/ 还是先归集到 platformWallet。这是改写必须定的业务点。

**改写要点**：
- `payBill()` 去 payable，无 msg.value/找零逻辑（ERC20 精确金额）→ `usdt.safeTransferFrom(msg.sender, ...)` 把 fee 转 platformWallet、amount 转运营商（或归集）。
- 需用户先对 Payment 合约 `approve` total。

## 四、FeeManager 费率读取接口现状

| 接口 | 签名 | 值 |
|------|------|-----|
| `getFeeRate()` | `view returns(uint256)` | 基点，部署初始化 150（1.5%） |
| `calculateFee(amount)` | `view returns(uint256)` | `amount * feeRate / FEE_DENOMINATOR` |
| `FEE_DENOMINATOR` | `constant` | 10000 |
| `MAX_FEE_RATE` | `constant` | 1000（10%）|
| `setFeeRate(newRate)` | `onlyOwner`，≤1000 | emit FeeRateUpdated |

→ 满足 R5，前端直接读。FeeManager 本身**无需改 ERC20**（纯计算）。

## 五、TrafficCardNFT 发卡/销毁/有效期 + Deposit 时序

- **mint**：`mint(to,dataAmount,uri) onlyOwner`（ownership 转给 Deposit 后由 Deposit 调）；记 `CardInfo{dataAmount,createdAt,isDestroyed}` + `_userCardCount++`。
- **burn**：持有人/授权可 burn → `isDestroyed=true` + `_burn` + emit `ServiceActivated(now+30days)`（销毁后 30 天服务，对齐前端文案）。
- **有效期**：常量 `DEDUCTION_VALIDITY = 30 days`（`_deductionCredits` 预留未用）。
- **Deposit↔发卡时序**：
  1. 用户 `Deposit.deposit()` 锁仓 30 天。
  2. 锁仓到期，owner 调 `Deposit.mintTrafficCard(user)`（校验到期 + 有存款 + 无卡）→ `dataAmount = _deposits/100000` → `TrafficCardNFT.mint`。
  3. `Deposit.issueMonthlyTrafficCards`（仅 oracle）当前**空实现**——PRD R7「锁仓满自动发放」目前**靠手动 mintTrafficCard，无真正自动链路**，arch 需补设计。
- **ERC20 影响**：发卡本身不涉及资金转移，**不直接受 ERC20 改写影响**；但 `dataAmount = _deposits/100000` 这个除数在 USDT 6 位精度下含义会变（见 §七），需复核。

## 六、ERC20 改动触及的 interface 签名 + 事件

| 文件 | 当前 | 改写后 |
|------|------|--------|
| `IDeposit.deposit` | `() external payable` | 去 payable，建议 `deposit(uint256 amount)` |
| `IDeposit` 事件 | `DepositMade(user,amount)` 等 | 金额语义改为 USDT 单位（值类型不变，监听方需知精度）|
| `IPayment.payBill` | `(uint256) external payable` | 去 payable（金额由 transferFrom 决定）|
| `IPayment` 事件 | `BillCreated/BillPaid` | 同上，金额单位变 USDT |
| Deposit/Payment 合约 | — | 新增 `IERC20 usdt` + initialize 注入；用 `SafeERC20` |

> ⚠️ 接口签名变更会**冲击后端**（子系统②）：`event_sync.go`/`signatures.go` 的事件签名、ABI、金额单位都要对齐（PRD §三②已标注）。后端按事件 topic hash 监听，`deposit()`→`deposit(uint256)` 会改函数 selector（但事件签名不变除非改事件参数）。

> ⚠️ **UUPS 升级兼容**：Deposit/Payment 已部署为 proxy。新增 `usdt` state 变量必须**追加在现有 storage 末尾**，不能插入。建议升级而非重部署（保留地址/manifest）。Arbitrum 是全新部署，无此约束。

## 七、USDT 精度问题（6 vs 18）相关现状代码

ETH/原生币是 18 位小数，USDT 通常 **6 位**。现状代码中所有金额都按 18 位语义写：

| 位置 | 现状 | 6 位精度下的风险 |
|------|------|------------------|
| `Deposit.mintTrafficCard` | `dataAmount = _deposits[user] / 100000` | 18 位下 0.1 ETH→大数；6 位下 10 USDT = 10_000_000，/100000 = 100。除数语义需重定 |
| `FeeManager.calculateFee` | `amount * feeRate / 10000` | 小额账单（如 1 USDT=1_000_000）* 150 /10000 = 15000，OK；但更小额可能截断为 0 |
| `ServiceManager` requiredDeposit | `0.003~0.012 ether`（18 位硬编码） | 切 USDT 后这些值语义完全错，需全部按 USDT 6 位重设或弃用 |
| 测试 `linkworld.ts` | 全用 `parseEther`(18 位) | 改 ERC20 后测试金额需改 USDT 精度 |
| MIN_DEPOSIT（待加） | — | 必须 = `10 * 10**usdt.decimals()`，**从 mock 实读 decimals 而非写死 6**（R4 要求核对 mock 精度）|

→ **强烈建议**：合约内不硬编码 6，用 `usdt.decimals()` 或部署参数注入；mock USDT 部署时显式设 decimals=6 并在 deployments json 记录，前端/后端统一读。
