# Task 04 — T3 Payment ERC20 + 链上分账 (P3)

> 子项目 contracts(1/3) | 状态 DONE | 2026-06-08

### 产出文件
- packages/contracts/contracts/Payment.sol（import SafeERC20/IERC20/IServiceManager + state usdt/serviceManager + initialize 4 参注入 + setServiceManager + createBill 改 onlyOracle+fail-fast + payBill 去 payable→nonReentrant+CEI+双段 safeTransferFrom+零地址校验+0-fee 跳过）
- packages/contracts/contracts/ServiceManager.sol（setOperatorPaymentAddress(id,addr) onlyOwner + 零地址 require + emit OperatorPaymentAddressSet）
- packages/contracts/contracts/interfaces/IPayment.sol（payBill 去 payable）
- packages/contracts/contracts/interfaces/IServiceManager.sol（补 setOperatorPaymentAddress + 事件声明）
- packages/contracts/test/payment.ts（新建：17 用例）
- packages/contracts/test/linkworld.ts（REG：deployAndWire 改 4 参 initialize+payment.setOracle；PM-01/02 改 oracle 调 createBill+预设 paymentAddress+USDT 精度金额）

### git commit
c6838e8 feat: 合约 T3 Payment ERC20 分账（双段 safeTransferFrom + createBill onlyOracle + 零地址校验）

### TDD
先红后绿：先写 test/payment.ts → RED（Payment.initialize 4 参不存在 fragment 报错）→ 实现 Payment/ServiceManager/接口 → GREEN 17 用例全过。

### 测试结果
hardhat clean+compile：55 Solidity files 0 error。hardhat test：55 passing / 0 failing（基线 38 + 新 17：createBill 权限+fail-fast、PAY-01~05、PAY-02b、0-FEE、SM-PAY-01~03、ATC-01/02、Not your bill、Already paid）。主 Agent 已独立复跑确认 55 passing（PAY-02 分账/PAY-03 零地址/PAY-05 原子回滚 均绿）。

### code-simplifier
using SafeERC20 简化转账；分账两段各自从 user 拉款（合约不暂存资金，降重入面）；改动聚焦 Payment+ServiceManager。

### spec review
严格按 design §4.2/§5.2/§3.2/§4.3/§6.1/§6.2 + arch-review B2 执行：createBill onlyOracle+fail-fast、payBill 双段 safeTransferFrom（amount→operator、fee→platform）、CEI 先置 isPaid、nonReentrant、零地址三重校验（setter+createBill+payBill）、0-fee 跳过、ServiceManager 独立 setter 方案。无偏差。

### 设计还原
合约无 UI。design §3.2 账单分账状态机逐项落地：未付→approve→payBill（校验 !isPaid + bill.user==sender + paymentAddress≠0）→双段转账→isPaid=true 终态。

### 复用检查
复用 OZ SafeERC20、FeeManager.calculateFee（未改，PAY-04 验证一致）、Payment 已有 onlyOracle modifier（T1）、ServiceManager Operator.paymentAddress 字段（IServiceManager 已有）；新增独立 setter 而非改 updateOperator 签名（减小 ABI 影响）。

### 设计稿对照
数值对照：payBill 转账段数 2 段（amount→operator + fee→platform）vs design §3.2 ✅；零地址校验点 3 处（setter+createBill+payBill）vs §6.1 ✅；fee=calculateFee(amount)=amount*150/10000（PAY-04 链上一致）✅；测试 55 passing（38+17）vs 无回归 ✅；Payment initialize 参数 4 个（feeManager/platformWallet/usdt/serviceManager）vs design §4.2 ✅；compile error 0 ✅。

### 新增组件
无新增合约。新增函数：Payment.setServiceManager、ServiceManager.setOperatorPaymentAddress；新增事件：ServiceManager.OperatorPaymentAddressSet。

### 新增色值
无（合约任务）。
