> 扫描时间：2026-06-08 | 子项目：contracts(1/3) | 对象：packages/contracts | 已核对真实代码

# 测试与部署基线（test-deploy-baseline）

## 一、现有测试覆盖

### `test/linkworld.ts`（主测试，8.3KB）
用 `ethers v6` + chai。**注意：不通过 proxy 部署**，而是 `Factory.deploy()` 直接部署 impl 再手动调 `initialize()`（绕过 UUPS proxy，测的是实现合约逻辑而非代理行为）。`deployAndWire()` 部署全 7 合约并 wire（setOracle/setTrafficCardNFT/setDeposit）。

| describe | it（编号） | 覆盖点 |
|----------|-----------|--------|
| FeeManager | FE-01~09 | initializer 150bps、calculateFee(1/10/0 ETH)、setFeeRate(500/1000/1001 revert)、非 owner revert、FEE_DENOMINATOR |
| UserRegistry | UR-01~04 | register 成功、重复 revert、未注册 false、name/symbol |
| ServiceManager | SM-01~03 | 初始 active>0、按国家查 US、addOperator(id=12) |
| Deposit | DP-01~03 | deposit({value})+getDepositAmount、未注册 revert、setOracle 存储 |
| Payment | PM-01~02 | createBill emit BillCreated、getUserBills |
| Oracle | OR-01~02 | deposit 地址存储、verifyServiceActive 未注册 false |
| TrafficCardNFT | TC-01~05 | name/symbol、mint ownerOf、emit CardMinted、getUserCardCount、getCardInfo |

### `test/UserRegistry.test.ts`（旧测试，与 linkworld.ts UR 部分重叠）
2 个 it：register 成功 + mint NFT、重复注册 revert。

### 覆盖缺口（给 test 阶段）
- ❌ **withdraw 全无测试**（提取本金、锁仓未到 revert）。
- ❌ **payBill 全无测试**（支付、找零、fee 分账）——本轮 ERC20 改写后必须补。
- ❌ **mintTrafficCard / 锁仓到期 / 发卡时序**无测试。
- ❌ **最小额约束**无测试（验收 A.2 要求新增 <10 拒绝 / =10 通过）。
- ❌ **UUPS 升级路径**（upgradeProxy / storage 兼容）无测试。
- ❌ 测试全用 `parseEther`（18 位），改 USDT(6 位)后需重写金额。
- ❌ Oracle `monthlySettlement` 无测试（且该函数存在编译疑点，见 inventory §⑥）。

**测试链/账户**：默认 hardhat 内存链（chainId 31337），`ethers.getSigners()` 取 owner/user1/user2/user3。

## 二、部署脚本流程（`scripts/deploy.ts`）

UUPS proxy 全量部署，**顺序按依赖拓扑**：

```
1. FeeManager  deployProxy([150], uups)         # 无依赖
2. UserRegistry deployProxy([], uups)            # 无依赖
3. ServiceManager deployProxy([], uups)          # 无依赖
4. TrafficCardNFT deployProxy([], uups)          # 无依赖
5. Payment     deployProxy([feeManager, deployer], uups)   # 依赖 FeeManager；platformWallet=deployer
6. Deposit     deployProxy([userRegistry], uups)           # 依赖 UserRegistry
7. Oracle      deployProxy([], uups)
── wiring ──
trafficCardNFT.setDepositContract(deposit)
deposit.setTrafficCardNFT(trafficCardNFT)
deposit.setOracle(oracle)
payment.setPlatformWallet(deployer)
oracle.setDeposit(deposit)
trafficCardNFT.transferOwnership(deposit)   # 关键：让 Deposit 能调 mint
```

产出写入 `deployments/<networkName>.json`：`{chainId, proxies{7}, implementations{7}}`，impl 地址通过 `upgrades.erc1967.getImplementationAddress` 读取。

> ⚠️ **wiring 缺口（与本轮相关）**：
> - **`payment.setOracle(oracle)` 没调** → Oracle 无法以 onlyOracle 身份操作 Payment（且 Oracle 侧也无 `setPayment`，见 inventory §⑥）。
> - `oracle.setDeposit` 调了，但 Oracle 的 `payment` 地址全程没被设置。
> - ServiceManager 的运营商 `paymentAddress` 全 0，未在部署中补真实地址。
>
> 本轮 ERC20 + Arbitrum 改造需顺带补全 wiring，并新增 **mock USDT 部署 + 注入 Deposit/Payment**。

**部署命令**（package.json scripts）：
- `deploy:local` → `--network localhost`（外部节点，config 无显式 localhost 定义）
- `deploy:0g-testnet` → `--network og_testnet`(16602)
- `deploy:0g` → `--network og_mainnet`(16600)
- **缺 Arbitrum Sepolia(421614) 部署 script**，本轮需新增。

## 三、smoke 脚本作用（`scripts/smoke*.ts`）

这 4 个是**早期调试/探针脚本**，非正式部署链路（探测 hardhat-upgrades 插件可用性、proxy 部署写法），**不应在生产流程使用**：

| 脚本 | 作用 |
|------|------|
| `smokeCheck.ts` | 探测 `upgrades` 插件导出方式、试 `F.deploy({args})`（多为试错日志） |
| `smokeDeploy.ts` | 探测 deployProxy 调用模式（含废弃 Proxy 报错 helper） |
| `smokeInitAfterDeploy.ts` | 验证 impl 部署后可手动 `sendTransaction` 调 initialize（FeeManager/Deposit）|
| `smokeTask.ts` | 打印 hre.upgrades 是否可用 |

→ 建议本轮可清理或忽略，正式链路只认 `deploy.ts`。（本扫描阶段不动，仅记录。）

## 四、现有部署产物地址表

### localhost.json（chainId 31337）
| 合约 | proxy |
|------|-------|
| FeeManager | 0x2db016b079a38e6bBDF67230afa64F6a130ea148 |
| UserRegistry | 0x7d298bEcF48370e256d74372De3A8344d0a34CDD |
| ServiceManager | 0x3755Ce7d04f3F74f9B20298fC03a79e5ae890A7D |
| TrafficCardNFT | 0x33D88b2B4096646F49227722d936D094081Ce96E |
| Payment | 0x1781bA57D052e4524Faa549f68408EB41133A3e5 |
| Deposit | 0xffcFc38663Ff7BA540490F403beaDE4F29835611 |
| Oracle | 0x45d2E95790AC00d9c47E4feDbeF8633836C717Da |

### og_testnet.json（chainId 16602）
| 合约 | proxy |
|------|-------|
| FeeManager | 0xF9d4777b760cc3a0F39eE0E11Cc936E34dcfc033 |
| UserRegistry | 0x0D0E7AeB3437682964d8164835eAE31c86451268 |
| ServiceManager | 0x82CB050c84F3BBEfC01D089d8579805Eb493BA14 |
| TrafficCardNFT | 0x8B29aC425eD0b021CFFb308494707A5f4e6DEd31 |
| Payment | 0x85Ffe2f47dF883982A6c98f665670e045fd0bfd9 |
| Deposit | 0x1c73baEceE72d0867b046f939Dd27fbbc714332b |
| Oracle | 0x1820f818dF0dE96d29eA3AA7007785eBE46662D1 |

OZ upgrades manifest：`.openzeppelin/unknown-16602.json`（manifestVersion 3.2，记录 16602 上各 proxy 的 address/txHash/kind=uups）。**无 localhost(31337) 对应 manifest 文件**（hardhat 网络的 manifest 通常不持久化）。**无 Arbitrum manifest**（本轮新增）。

> 现有地址全部基于**原生币（payable）旧版合约**。本轮 ERC20 改写后，og_testnet/localhost 旧产物作废或需升级；Arbitrum Sepolia(421614) 为全新部署，产出 `deployments/arbitrum_sepolia.json` + mock USDT 地址。
