> 扫描时间：2026-06-08 | 子项目：contracts(1/3) | 对象：packages/contracts | 已核对真实代码

# 合约子项目结构概览（project-scan）

## 一、目录树

```
packages/contracts/
├── hardhat.config.ts          # 链配置/编译器
├── package.json               # 依赖（OZ 5.6.1 / hardhat 2.22 / toolbox 5）
├── tsconfig.json
├── contracts/
│   ├── Deposit.sol            # 保证金（payable，本轮改写核心①）
│   ├── Payment.sol            # 账单支付（payable，本轮改写核心②）
│   ├── FeeManager.sol         # 手续费（150bps / 10000，已是事实源）
│   ├── TrafficCardNFT.sol     # 流量卡 NFT（mint→burn→30天服务）
│   ├── ServiceManager.sol     # 运营商目录（11 个内置运营商）
│   ├── Oracle.sol             # 计量/月末结算（含一个内联 IPayment/IDeposit）
│   ├── UserRegistry.sol       # 用户注册 + 身份 NFT（邮箱注册）
│   └── interfaces/            # 7 个接口（I*.sol）
├── scripts/
│   ├── deploy.ts             # 正式部署脚本（UUPS proxy 全量）
│   └── smoke*.ts             # 4 个调试探针脚本（非正式链路，见 baseline 文档）
├── deployments/
│   ├── localhost.json        # chainId 31337 地址表
│   └── og_testnet.json       # chainId 16602 地址表
├── .openzeppelin/
│   └── unknown-16602.json    # OZ upgrades proxy manifest（manifestVersion 3.2）
├── artifacts/                # 编译产物（已存在，曾编译成功）
└── cache/
```

## 二、hardhat.config 网络与编译器

**编译器**：solidity `0.8.27`，`optimizer enabled / runs 200`，`evmVersion: cancun`，`viaIR: true`。

> ⚠️ `cancun` EVM 版本 + `viaIR` —— 迁移 Arbitrum Sepolia 时需确认 Arbitrum 对 cancun opcode（尤其 transient storage，Payment 用了 `ReentrancyGuardTransient`）的支持，这是本轮一个待验证点。

**当前网络配置**（共 3 个，**无 hardhat 显式 localhost、无 Arbitrum**）：

| 网络名 | url | chainId | accounts | 用途 |
|--------|-----|---------|----------|------|
| `hardhat` | (内置) | 31337 | 内置 | 本地内存链 |
| `og_mainnet` | https://evmrpc.0g.ai | 16600 | DEPLOYER_KEY | 0G 主网 |
| `og_testnet` | https://evmrpc-testnet.0g.ai | 16602 | DEPLOYER_KEY | 0G 测试网（现部署） |

私钥来自 env `DEPLOYER_PRIVATE_KEY`（缺省回退全 0）。`deploy:local` script 跑的是 `--network localhost`，但 config 里**没有显式 localhost 网络定义**，依赖 hardhat 默认的 `localhost`（指向 `http://127.0.0.1:8545` 外部节点），与内置 `hardhat`(31337) 不同。

## 三、是否 UUPS proxy（可升级）

**全部 7 个合约均为 UUPS 可升级**：
- 均继承 `OwnableUpgradeable + UUPSUpgradeable`，用 `initialize()` 而非 constructor，每个都实现 `_authorizeUpgrade(...) onlyOwner`。
- 部署脚本 `deploy.ts` 用 `upgrades.deployProxy(..., { kind: "uups" })`，配 `@openzeppelin/hardhat-upgrades ^3.9.1`。
- `.openzeppelin/unknown-16602.json`（manifestVersion 3.2）记录了 16602 链上 proxy 列表，全部 `kind: "uups"`。

→ **本轮 ERC20 改写应走 `upgradeProxy` 升级路径**（而非重新部署），需遵守 OZ storage layout 兼容规则（新增 state 变量只能追加，不能插入/改序），且 16602 旧 manifest 与 Arbitrum 新部署是两套独立 manifest。

## 四、现有部署现状（地址表）

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

> 注：两份 json 的 `implementations` 地址完全相同（同一份 impl 字节码部署到两条链得到相同地址，因 CREATE 依赖 deployer+nonce），proxy 地址不同。

## 五、Arbitrum 421614 配置缺口（预期无，已核实确实无）

- `hardhat.config.ts` networks 中**无 arbitrum / 无 chainId 421614**。
- 全仓 grep `arbitrum` / `421614`（.ts/.json/.sol，排除 node_modules/artifacts）→ **零命中**。
- **本轮需新增**：① hardhat.config 加 `arbitrum_sepolia`(421614) 网络 + RPC + accounts；② 新增显式 `localhost`(31337) 网络（或沿用默认）；③ deploy 脚本产出 `deployments/arbitrum_sepolia.json`；④ 自部署 mock USDT(ERC20) 并写入部署产物。

## 六、依赖版本

| 包 | 版本 | 备注 |
|----|------|------|
| `@openzeppelin/contracts` | ^5.6.1 | 含 `ReentrancyGuardTransient`、`SafeERC20`、`IERC20`（本轮 ERC20 改写直接可用，无需升级） |
| `@openzeppelin/contracts-upgradeable` | ^5.6.1 | UUPS / Ownable / ERC721URIStorage upgradeable |
| `@openzeppelin/hardhat-upgrades` | ^3.9.1 | deployProxy/upgradeProxy |
| `@nomicfoundation/hardhat-toolbox` | ^5.0.0 | ethers v6 + chai matchers |
| `hardhat` | ^2.22.0 | |
| `dotenv` | ^17.4.2 | 读 DEPLOYER_PRIVATE_KEY |
| `axios` | ^1.16.0 | (运行期依赖，合约本身不用) |

→ OZ 5.6.1 **已自带 `IERC20` + `SafeERC20`**，ERC20 改写无需新增依赖，直接 `import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol"`。
