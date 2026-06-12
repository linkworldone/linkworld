# 后端子项目 blockchain-integration 基线

> 扫描 2026-06-08 | 子项目 backend(2/3) | 已核对真实代码

## 总览：链集成当前几乎全为 stub

后端**当前不发起任何链上写交易**，也**未实际订阅任何事件**。链相关代码三件（`client.go` / `signatures.go` / `event_sync.go`）均为骨架。grep 全 `internal/` 无 `createBill` / `monthlySettlement` / `issueMonthly` / `approve` / abigen 绑定 / 任何 `ethClient.Call*` 实际调用。

这意味着后端「计价 → 上链结算 → 落库对账」的链路**尚未实现**，alignment 阶段不是「改调用签名」而是「从零接入链上写 + 事件同步」。

---

## 1. client.go —— 链连接 + 业务方法 stub

```go
type Client struct {
    rpcURL    string
    chainID   uint64
    contracts map[string]common.Address
    ethClient *ethclient.Client
}
```

- `NewClient(rpcURL, chainID, contracts)`：`ethclient.Dial(rpcURL)` 真连，存地址表。**这是唯一真实生效的部分**。
- `EthClient()`：返回底层 `*ethclient.Client`。
- 以下全部是 stub（`ethClient==nil` 直接返回零值，否则也只 log 或返回未实现）：
  - `GetDepositAmount(addr)` → 永远 `big.NewInt(0)`
  - `GetLockExpiry(addr)` → 永远 `0`
  - `VerifyServiceActive(addr)` → 永远 `false`
  - `ListenEvents(ctx)` → `not implemented`
- **没有合约 ABI 绑定**：未用 `bind.NewBoundContract` 或 abigen，无法直接编码/解码合约调用。

## 2. signatures.go —— 事件 topic 常量（keccak256）

运行时用 `crypto.Keccak256Hash([]byte("<sig>"))` 计算事件签名 topic：

| 常量 | 事件签名（当前后端写死） | 与冻结 ABI 是否一致 |
|------|--------------------------|---------------------|
| `UserRegisteredTopic` | `UserRegistered(address,string,uint256)` | ✅ 与 `UserRegistry.json` 一致 |
| `DepositMadeTopic` | `DepositMade(address,uint256)` | ✅ 与 IDeposit 一致 |
| `BillCreatedTopic` | `BillCreated(uint256,address,uint256,uint256,uint256)` | 🔴 **5 参数，错**。冻结 ABI 是 `BillCreated(uint256,address,uint256,uint256)`（4 参数，见 IPayment.sol）→ topic hash 错，永远匹配不到事件 |
| `BillPaidTopic` | `BillPaid(uint256,address,uint256,uint256)` | ✅ 与 IPayment `BillPaid(billId,user,totalAmount,operatorAmount)` 一致 |
| `TrafficCardMintedTopic` | `TrafficCardMinted(address,uint256,uint256)` | ✅ 与 IDeposit 一致 |

> 缺失：`BillCreated` 应改名/对齐为 4 参数；新合约还有 `TrafficCardApplied(uint256)`（applyTrafficCardToBill 桩 emit）后端完全未定义。

## 3. event_sync.go —— 事件同步 stub

```go
func (s *EventSync) Start(ctx)           // 只起 go s.syncUserRegistered(ctx)
func (s *EventSync) syncUserRegistered   // for{} 里 sleep 30s 空转，TODO: 未实现 FilterQuery
func (s *EventSync) processUserRegistered(log) // 有解析骨架：topics[1]→user 地址→Create User
func (s *EventSync) processDepositMade(log)    // 空返回 nil
```

- 只声明同步 `UserRegistered` 一种事件，且**主循环没有真正调 `FilterLogs`/`SubscribeFilterLogs`**，`process*` 永不被触发。
- `processUserRegistered` 落库时 `RegisteredAt: time.Unix(0,0)`（占位），不解析 email/tokenId。
- **金额类事件（DepositMade/BillPaid/...）完全没有落库处理**，无任何「6 位精度解释」逻辑。

## 4. ABI 来源现状

- 后端 ABI 在 `internal/blockchain/abis/`，**只有 2 份**：`Deposit.json`、`UserRegistry.json`。
- 这两份是**手写裁剪版**（只含部分事件 + 个别 function），**不是合约编译产物全量 ABI**，也无 `abiHash` 比对。
- 缺失 5 份：`FeeManager` / `ServiceManager` / `TrafficCardNFT` / `Payment` / `Oracle` / `MockUSDT`（共 6 份业务/代币 ABI 未落地）。
- `Deposit.json` 当前**只有事件**，没有 `deposit(uint256)` / `withdraw()` / `issueMonthlyTrafficCards` 等 function 定义 → 即使要发交易也无法编码。

> handoff §1 要求 ABI 从 `packages/contracts/artifacts/contracts/<Name>.sol/<Name>.json` 的 `.abi` 字段取全量，并按 `deployments/<net>.json.abiHash` 比对。当前后端完全脱节。

## 5. configs/deployments.json 现状

```json
{
  "chainId": 16602,
  "rpcUrl": "https://evm-testnet.0g.ai",
  "contracts": {
    "FeeManager":     "0xF9d4777b760cc3a0F39eE0E11Cc936E34dcfc033",
    "UserRegistry":   "0x0D0E7AeB3437682964d8164835eAE31c86451268",
    "ServiceManager": "0x82CB050c84F3BBEfC01D089d8579805Eb493BA14",
    "TrafficCardNFT": "0x8B29aC425eD0b021CFFb308494707A5f4e6DEd31",
    "Payment":        "0x85Ffe2f47dF883982A6c98f665670e045fd0bfd9",
    "Deposit":        "0x1c73baEceE72d0867b046f939Dd27fbbc714332b",
    "Oracle":         "0x1820f818dF0dE96d29eA3AA7007785eBE46662D1"
  }
}
```

问题清单：
1. **chainId = 16602**（0G testnet），应为 **421614**（Arbitrum Sepolia）。
2. **rpcUrl = `https://evm-testnet.0g.ai`**，应为 Arbitrum Sepolia RPC。
3. **键名 `contracts`** ≠ Go struct 读的 **`proxies`** → 地址表实际加载为空（见 project-scan §4）。
4. **7 个地址全是 0G 上的旧部署**，与合约子项目新部署产物（`packages/contracts/deployments/hardhat.json` 的 proxies，或待生成的 `arbitrum_sepolia.json`）无关。
5. **缺 `usdt` / `usdtDecimals` 字段**：handoff §4 要求从这两个字段读 USDT 地址与精度，后端配置完全没有。
6. **缺 `abiHash` 字段**：无法据 handoff §1 比对 ABI 是否同步。

## 6. 合约子项目最新部署产物（对照源）

`packages/contracts/deployments/hardhat.json`（chainId 31337，本地 fresh deploy，**字段结构是后端应对齐的范式**）：
- `proxies`：7 合约新代理地址（与后端 0G 旧地址完全不同）。
- `implementations`：7 实现地址。
- `usdt`：`0x5FbDB2315678afecb367f032d93F642f64180aa3`（MockUSDT），`usdtDecimals`：6。
- `abiHash`：7 合约 ABI 指纹。
- `storageLayout`：升级基线说明。
- Arbitrum Sepolia(421614) 的 `arbitrum_sepolia.json` **尚未生成**（handoff §10：缺 `DEPLOYER_PRIVATE_KEY`/RPC，真·上链待执行）。
- 另有旧的 `localhost.json` / `og_testnet.json`（旧结构，已过时）。
