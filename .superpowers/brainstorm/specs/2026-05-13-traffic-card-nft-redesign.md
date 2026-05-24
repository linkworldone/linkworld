# 流量卡 NFT 化重设计方案

> 创建日期：2026-05-13  
> 阶段：复盘阶段提出，作为 Round 3 候选需求  
> 由：产品 logan 提出意图，claude 整理

## 产品目标

**核心模型**：
- 用户充押金 → 平台拿去稳健理财（DeFi）获利
- 押金锁定满月 → 自动发放一张 ERC-1155 流量卡
- 用户使用流量卡 = 销毁 NFT + 触发链下冲值
- 链上不记账、不收钱、不管实际用量

**用户视角**：押金不亏（可退）+ 流量免费（理财收益换的）。

## 与当前实现的差距

| 维度 | 当前 (rom 版) | 新方案 |
|---|---|---|
| 流量卡载体 | `Deposit._trafficCardBalances` uint256 mapping | ERC-1155 NFT（独立合约） |
| 押金作用 | 资格证明 + 反跑路 | 锁仓 → 产卡的质押凭证 |
| 计费 | Oracle 月末上报用量 → Payment 出账单 → 用户 payBill | **完全删除**，链上无账单 |
| 用户消费 | 抵扣账单（自动） | 主动 burn NFT + 指定 country/number |
| 后端联动 | 月末扣款（未实现） | 订阅 `TrafficCardConsumed` 事件 → 调运营商 API 冲值 |
| 理财 | 押金死躺合约 | 新增 invest/redeem 接口接入 DeFi |
| 手续费 | FeeManager 计算平台费 | **删除**，平台靠理财收益赚钱 |

## 合约改动总览

### 删除
- `Payment.sol` 全文（无账单）
- `Oracle.sol` 全文（无用量上报）
- `FeeManager.sol` 全文（无手续费）
- `Deposit.sol` 内的流量卡相关字段与函数：
  - `trafficCardQuota`、`_trafficCardBalances`、`_monthlyCardIssued`
  - `setTrafficCardQuota`、`issueMonthlyTrafficCards`、`useTrafficCard`、`getTrafficCardBalance`
  - `hasWithdrawnThisMonth`、`getCurrentMonth`、`getMonthFromTimestamp`

### 改造 Deposit.sol
- 保留 `deposit()` / `withdraw()` / `getDepositAmount()`
- **加锁定期**：
  - `_depositTimestamp[user]` 每次充值时刷新
  - `withdraw()` 校验 `block.timestamp - _depositTimestamp >= LOCK_PERIOD`（默认 30 天）
- **加产卡接口**：
  - `claimTrafficCard()`：用户调用
  - 校验：押金 > 0 且距上次 claim ≥ 30 天
  - 计算应发流量 = 押金 × CARD_RATIO（如每 100 USDT → 10 MB）
  - 调 `TrafficCardNFT.mint(user, dataMB)`
- **加理财接口**（onlyOwner / onlyStrategist）：
  - `invest(strategy, amount)`：把合约里的押金调到外部 DeFi（Aave/sDAI/Lido）
  - `redeem(strategy, amount)`：从 DeFi 赎回到合约
  - 用户 `withdraw()` 时合约 ETH 不够 → 自动触发 redeem
  - 累计收益归平台所有，用户得到的回报已经是免费流量卡

### 新增 TrafficCardNFT.sol（ERC-1155）

**tokenId 设计**：
- 推荐方案：**一卡一 id**（每次 mint 递增），便于追踪每张卡的国家/号码绑定记录
- 备选方案：按面额分档（10MB / 100MB / 1GB），同档共享 tokenId — 更省 gas 但失去单卡追踪

**核心接口**：
```solidity
function mint(address user, uint256 dataMB) external onlyDeposit returns (uint256 tokenId);

function useTrafficCard(
    uint256 tokenId,
    string calldata countryCode,
    string calldata phoneNumber
) external {
    require(balanceOf(msg.sender, tokenId) > 0, "Not owner");
    _burn(msg.sender, tokenId, 1);
    emit TrafficCardConsumed(tokenId, msg.sender, countryCode, phoneNumber, _cardData[tokenId], block.timestamp);
}

function getCardInfo(uint256 tokenId) external view returns (uint256 dataMB, uint256 issuedAt, uint256 expiresAt);
```

**事件**：
```solidity
event TrafficCardIssued(
    uint256 indexed tokenId,
    address indexed user,
    uint256 dataMB,
    uint256 expiresAt
);

event TrafficCardConsumed(
    uint256 indexed tokenId,
    address indexed user,
    string countryCode,
    string phoneNumber,
    uint256 dataMB,
    uint256 timestamp
);
```

## 后端改动

- 删除：账单同步、月末结算 handler、Bill model 相关 API
- 新增：事件监听器
  - 订阅 `TrafficCardConsumed`
  - 收到 → 查 phoneNumber 对应的真实 SIM/eSIM → 调对应运营商 API 冲 dataMB
  - 记录 consumption log（供用户查询历史）
- 简化：Deposit 同步逻辑保留，因为前端要展示锁定时间

## 前端改动

- 删除 Billing / BillDetail 页面相关合约调用（保留页面壳作"流量卡历史"用）
- 改造 Dashboard：展示用户拥有的流量卡列表（按 tokenId 查询）
- 改造 Deposit 页：加锁定期倒计时 + claim 按钮
- 新增 UseCard 页：选卡 + 选国家 + 输号码 → 调 useTrafficCard

## 待定参数（实施前需拍板）

| 参数 | 候选值 | 决策依据 |
|---|---|---|
| LOCK_PERIOD | 30 天 / 60 天 / 90 天 | 平衡理财周期 vs 用户体验 |
| CARD_RATIO | 100 USDT → 10 MB / 50 MB | 取决于运营商批发价 + 理财年化 |
| 卡有效期 | 永久 / 90 天 / 跟随押金 | 防止"屯卡套利" |
| tokenId 策略 | 一卡一 id / 按档分 | 看是否需要单卡追踪 |
| 理财协议 | sDAI / Aave / Lido | 取决于稳定币 vs ETH 押金 |

## 实施路径（建议）

1. 起新分支 `feat/traffic-card-nft`，旧合约不动
2. 先写 `TrafficCardNFT.sol` + 单测，独立可验证
3. 改 Deposit：加锁定 + 加 claimTrafficCard
4. 砍掉 Payment / Oracle / FeeManager（含部署脚本、ABI、前后端引用）
5. 写理财模块（可分多个策略合约）
6. 前后端同步改造
7. 写迁移策略：本地 Hardhat 直接重部署即可；线上有用户的话需要 snapshot 押金导入

## 未解决的问题

1. **押金币种**：用 ETH 还是 USDT？理财协议选型直接相关
2. **理财风险**：DeFi 协议风险归谁？合约亏损时用户押金能否兑付？
3. **运营商批发结算**：链下后端给运营商付钱的资金从哪来？是平台理财收益的一部分？需要明确资金流。
4. **流量卡转赠**：ERC-1155 默认可转，是否允许用户互转/二级市场？涉及 KYC/合规
5. **多 SIM 场景**：一个钱包绑多个 phoneNumber？还是一钱包一号？
6. **首次发卡时机**：用户充押金当天就给一张？还是必须锁满 30 天才有？

## 现存代码的复用价值

- `UserRegistry`：保留（NFT 身份不变）
- `ServiceManager`：保留（运营商目录、用户开号仍需要）
- `Deposit.deposit/withdraw`：保留骨架
- 其他：基本全砍

## 风险与权衡

**优点**：
- 业务模型清晰：押金理财 → 收益换流量
- 链上极简：去掉账单/手续费/Oracle 一大堆复杂度
- 用户体验更好：不用月末再付一次钱
- 卡作为 NFT 可视化、可追踪、可二级市场（如开放）

**风险**：
- 理财收益不稳定时如何承诺流量发放？需要"保底+浮动"机制
- 大规模 mint NFT gas 成本高 → 考虑 batch claim 或按档分 tokenId
- 后端事件监听稳定性 → 漏 event 会导致用户烧卡了但没冲流量

---

> 决策完待定参数后，本 spec 即可进入 Round 3 的 brainstorm / plan-eng-review 流程。

## 决策记录（2026-05-13 讨论补充）

### 1. 锁仓模式：硬锁 + 即时发卡（A 方案）

权衡 A/B 后选定 **A：用户充值即锁仓，锁仓期内不可退，但充值时立即发放可用流量卡**。

理由：
- 用户体验：充值即用，留存率高（B 方案要等 30 天才能用，流失大）
- 经济闭环：用户预支福利 = 平台用锁仓期理财收益覆盖发卡成本，30 天到期用户全额退本金
- 理财协议：可用期限型（如 Pendle PT），到期日匹配，年化更高
- 风险归属：用户本金锁在合约里跑不掉，平台无对手方风险

发卡额度公式：
```
expectedYield = principal × APR × lockDays / 365
maxTrafficMB = (expectedYield × (1 - profitMargin)) / wholesalePricePerMB
```

### 2. 产品定位：DeFi 模式（不动用户资金分毛）

明确产品形态为 **DeFi vault**，不是中心化金融产品。核心承诺：

**合约层硬不变量**：合约总资产 ≥ 全部用户本金之和（含理财协议中的资产），永远成立。

```solidity
contract Deposit {
    mapping(address => uint256) public principal;
    uint256 public totalPrincipal;
    
    // 平台只能收割 yield，不能动本金
    function harvest() external {
        uint256 totalValue = IStrategy(strategy).totalAssets();
        require(totalValue > totalPrincipal, "No yield");
        uint256 yield = totalValue - totalPrincipal;
        IStrategy(strategy).withdraw(yield);
        usdc.transfer(platformTreasury, yield);
    }
}
```

合约里**根本没有"转走用户押金"的函数**，从代码层杜绝吸储嫌疑。

### 3. 收入模型：Yield Spread（不是 LP 费）

LinkWorld 收入来源 = **理财收益 - 流量批发成本**，叫 Yield Spread（收益差）。

具体测算（参考）：
- 用户押 100 USDC，锁 30 天
- 平台投 sDAI，年化 6% → 30 天收益 ≈ 0.49 USDC
- 给 10MB 流量卡，批发成本 ≈ 0.20 USDC
- **平台净赚 ≈ 0.29 USDC / 用户 / 月**

平台**不靠克扣用户本金赚钱**，跟 Compound 等存款型协议同模式（区别：Compound 是借贷利差，LinkWorld 是 yield - 商品成本差）。

#### 收入模式 vs 其他 DeFi 模式对照

| 模式 | 代表 | LinkWorld 是否采用 |
|---|---|---|
| LP 费（交易手续费分成） | Uniswap/Curve | ❌ 不是 DEX |
| 借贷利差 | Aave/Compound | ❌ 不做借贷 |
| Performance Fee（yield 抽成） | Yearn/Convex | ✅ 部分采用 |
| **Yield Spread**（收益差） | sDAI/Pendle | ✅ **核心模式** |
| Protocol Fee | 1inch/GMX | ❌ 不适用 |

### 4. Owner 权限矩阵（最小化原则）

| 操作 | Owner 能做 | 限制 |
|---|---|---|
| 直接转用户押金 | ❌ | 合约里无此函数 |
| 收割 yield | ✅ | 仅可转到 platformTreasury，金额硬限 = totalAssets - totalPrincipal |
| 切换理财策略 | ✅ | 仅白名单中、且 24-72h timelock |
| 加白名单策略 | ✅ | 7 天 timelock |
| 暂停新存款 | ✅ | **不能阻止 withdraw** |
| 暂停 withdraw | ❌ 或可选 | 若加须 + 24h 自动恢复 |
| 升级合约 | ✅ | UUPS + 7 天 timelock |

Owner key 用 **Gnosis Safe 3/5 多签**，避免单点。

### 5. 理财策略选型（按风险分档）

| 阶段 | TVL | 策略 | 预期年化 | 归零概率 |
|---|---|---|---|---|
| 种子期 | < $100K | 100% sDAI | 5-7% | < 0.1% / 年 |
| 早期 | $100K-$1M | 80% sDAI + 20% Aave USDC | 5-8% | < 0.5% / 年 |
| 成长期 | > $1M | 60% sDAI + 30% Pendle PT + 10% buffer | 7-10% | 1-2% / 年 |

**绝对不碰**：Ethena 类合成稳定币、算稳、未审计协议、TVL < $100M 的协议。

集中度限制：单协议不超过押金总量 40%。
保留 5-10% 链上 buffer 应对短期赎回。

### 6. 合规姿态（DeFi 化的关键好处）

| 监管问题 | DeFi 模式的回答 |
|---|---|
| 是否吸储？ | 否。用户押金在合约里，平台无权动用，链上可查 |
| 是否承诺收益？ | 否。流量卡是用户押金锁仓产生的派生品，跟平台无关 |
| 用户的钱在哪？ | 在用户授权（通过合约调用）的白名单 DeFi 协议中 |
| 出事谁担责？ | 用户使用 DeFi 协议的智能合约风险，跟用户用 Aave 一样 |
| 平台收入来源？ | 收益差（yield spread），仅取 yield 部分，不取本金 |
| 用户能绕过前端吗？ | 能。合约公开，前端只是 UI，可用 IPFS 部署冗余 |

**主体注册建议**：BVI / 塞舌尔 / 开曼。前端屏蔽美国 + 中国 IP。NFT 流量卡**强烈建议 soulbound（不可转让）**，避免被认定为代币化证券。

### 7. 待补的设计点

- [ ] LOCK_PERIOD 具体值（30/60/90 天）
- [ ] CARD_RATIO 具体值（押金 → MB 兑换比）
- [ ] 多档锁定期 + 多档面值的具体公式
- [ ] Strategy 白名单的初始名单（仅 sDAI？还是多个？）
- [ ] Yield buffer 比例（5% / 10%）
- [ ] Emergency pause 条件（策略亏损 > X% 时触发）
- [ ] Soulbound 还是允许转让的最终决策

### 8. 最终决策（2026-05-19 拍板）

所有 7 个待补设计点已决策完毕：

| # | 项目 | 决策值 | 备注 |
|---|---|---|---|
| 1 | LOCK_PERIOD | **4 档：30 / 90 / 180 / 365 天** | 与 sDAI（30d 无锁）+ Pendle PT（90/180/365d）的标准到期日对齐 |
| 2 | CARD_RATIO | **公式驱动** | 不写死，按公式计算 |
| 3 | 多档公式 | `MB = principal × APR × lockDays / 365 × (1 − platformMargin) / wholesalePerMB` | platformMargin=20%, wholesalePerMB=$0.01 |
| 4 | Strategy 白名单 | **MVP 单策略 sDAI**，不做白名单架构 | 后续切策略需重新部署或合约升级 |
| 5 | Buffer 比例 | **5%** | 5% 留作即时赎回缓冲，95% 进 sDAI |
| 6 | Emergency pause | **仅 pause deposit，不做 pause withdraw** | 避免给中心化把柄，符合 DeFi "不动用户资金"承诺 |
| 7 | Soulbound | **不可转让** | 避免代币化证券嫌疑 |

#### 4 档预设参考表（公式输出）

| 档位 | 锁定期 | 协议 | APR | 100 U 换 | 1000 U 换 |
|---|---|---|---|---|---|
| 试用 | 30 天 | sDAI | 6% | 40 MB | 400 MB |
| 标准 | 90 天 | Pendle PT 90d | 8% | 160 MB | 1.6 GB |
| 推荐 | 180 天 | Pendle PT 180d | 10% | 400 MB | 4 GB |
| 长期 | 365 天 | Pendle PT 1y | 12% | 1 GB | 10 GB |

#### Emergency pause 实施细则
- 仅 Owner 多签可触发 pause / unpause
- 不设自动触发条件（避免被攻击者利用）
- 公告 24 小时后才能生效（用户有时间撤离）
- pause 状态下：deposit 函数 revert；withdraw、burn、harvest 全部正常
