# Stage: scan — 后端基线扫描（子项目 backend 2/3）

> **状态**: completed (DONE_WITH_CONCERNS) | **日期**: 2026-06-08 | **Gate**: 0 | **子项目**: backend(2/3)
> 对象：packages/backend(Go 1.25)；输入：合约 handoff-backend.md

## 产出（4 份基线，commit ba00aec，docs/design/linkworld-backend/）
project-scan.md / blockchain-integration.md / services-api.md / alignment-surface.md（最关键）

## 🔴 关键发现
1. **后端链集成几乎全是 stub**：grep 全 internal/ 无 createBill/monthlySettlement/issueMonthly/approve/abigen 绑定/任何链上写调用；client.go 业务方法返零值；event_sync 主循环只 sleep 30s 空转、process* 永不触发。→ 本子项不是「改签名」而是**从零接入链上写 + 事件同步**，链侧接近重写。
2. **链配置全错**：deployments.json chainId=16602 + rpcUrl=evm-testnet.0g.ai + 7 个 0G 旧地址；缺 usdt/usdtDecimals/abiHash。需换 421614 + Arbitrum RPC + 合约新地址（arbitrum_sepolia.json 尚未生成，依赖合约真·上链）。
3. **config.go 键名 bug**：struct 读 `proxies`，JSON 是 `contracts` → Proxies 永远空 map，event sync 静默失效。
4. **计价缺口（产品级）**：后端「计价」是 OperatorAPISimulator.GetBill 返回 rand.Intn(5000)+500 随机数；FetchAndCreateBills 只写 DB Bill，没碰 monthlySettlement，无 amounts[] 构造/分批/真实资费表。handoff 要求 createBill 金额由后端算好传入——**真实资费规则来源未定义，需产品澄清**。
5. **ABI 严重不全**：abis/ 只 2 份手写裁剪（Deposit 仅事件 + UserRegistry 部分），缺 6 份（FeeManager/ServiceManager/TrafficCardNFT/Payment/Oracle/MockUSDT），无 abigen 绑定。
6. **event_sync 受影响**：signatures.go 的 BillCreated 是旧 5 参签名（冻结 ABI 是 4 参，topic hash 错匹配不到）；缺 TrafficCardApplied(uint256)；金额事件无落库/无 6 位精度解释。
7. **签名非链上可用**：SignData 是 SHA256 拼字符串非 ECDSA/secp256k1；后端不持私钥。handoff §7：monthlySettlement 是 onlyOwner=deployer，后端只需 owner 私钥 + 合约 setter 拓扑保证权限。
8. **RPC 不一致确认**：hardhat 421614 + 旧 16602；后端 evm-testnet.0g.ai。两端历史都指 0G。

## ⚠️ 遗留/风险（移交 design/requirement）
- **真实资费/计价规则未定义**（产品/设计层缺口）——design 阶段需跟用户澄清来源（按用量？固定档？运营商费率表？）。
- arbitrum_sepolia.json 未生成、Arbitrum 上 MockUSDT 地址不存在 → 联调前置依赖合约子项(1/3)真·上链。
- 敏感端点（oracle/monthly-bill、usage/submit）当前无鉴权，接链后成攻击面（review 重点）。

## 移交 design
带着 handoff + 上述发现，确定后端对齐方案：链集成补全（abigen/client/event_sync）、链配置 421614、计价逻辑、签名/owner key、RPC 统一、鉴权。**计价规则是产品决策，design 需先澄清。**
