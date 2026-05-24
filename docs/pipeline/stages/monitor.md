# Monitor 阶段产出

> 阶段：复盘  
> 起止：2026-04-15 ~ 2026-05-19  
> Pipeline Round：2  
> 状态：完成

## 本轮目标

LinkWorld 全栈联调 —— 前端（React）从 Mock Service 切换到真实后端 API（Go/Gin）+ 链上合约（Solidity/Hardhat）交互，三端本地环境串通。

## 完成情况

- ✅ 三端环境本地联通：Hardhat (8545) + Backend (8080) + Web (5173)
- ✅ 前端 Mock 全部切换为真实合约 / 后端调用
- ✅ 合约部署到本地链 + 前后端地址对齐
- ✅ 完整冒烟测试通过：注册 → 充押金 → 申请号码 → 月结发卡 → 烧卡抵扣
- ✅ rom 推的 NFT 升级 commit 接入完成（修了他留的 6 个编译 bug）
- ✅ 流量卡 NFT UI 接入完成（Cards 页 + 烧卡 + 月结触发）
- ✅ 下一轮 spec 已成型（traffic-card-nft-redesign，280 行）

## 关键发现 / 业务洞察

### 1. 当前合约设计跟产品意图严重不一致
- 产品意图：用户押金 → 平台理财 → 收益换免费流量
- 当前实现：押金死躺合约 + 月末还要用户付账单
- 中间至少 3 个本质差距，已在 spec 中列明

### 2. rom 的代码质量需要管控
- 第一次 commit (rom TrafficCard)：编译不过
- 第二次 commit (rom upgrade NFT)：还是编译不过
- 模式很稳定：**本地不跑测试就 push**
- 建议下一轮加 pre-push hook 强制跑 `pnpm compile` 拦截

### 3. DeFi 模式 vs 中心化模式的分水岭
本轮讨论清晰确立产品定位：
- **DeFi vault 模式**（不动用户资金分毛）
- **收入靠 Yield Spread**（不是 LP 费、不是利差）
- **合规姿态**：BVI 离岸 + IP 屏蔽 + Soulbound 不可转让

### 4. 数字感的关键校准
- 用户充 100 U / 30 天，理财收益仅 0.5 U
- 减去流量批发成本（0.2 U），平台净利 0.3 U / 月 / 用户
- 平台 net margin 约 3-5%，是"薄利多销 + 规模化"模型
- 真正的目标用户不是"省话费的人"，是"币圈持仓闲置的人"

## 遗留问题（Round 3 候选）

| 优先级 | 项目 | 说明 |
|---|---|---|
| P0 | TrafficCard NFT 重设计 | 按 spec 实施，砍掉账单/计费逻辑 |
| P0 | Pendle PT 等理财策略接入 | sDAI 适配器 + 多策略路由 |
| P0 | 合约层"本金不变量" | totalAssets >= totalPrincipal 硬保证 |
| P0 | 移除合约里明文 password 字段 | 安全隐患 |
| P1 | 后端 Indexer / Subgraph | 替代前端直读 RPC，TVL 起来后必须 |
| P1 | 钱包签名认证 | 后端 API 目前无 auth |
| P1 | Multicall 批量读 | Dashboard 单次刷新 5+ 次 RPC 调用 |
| P1 | rom 代码质量管控 | pre-push hook 强制编译 |
| P2 | History endpoint 增强 | 加分页 + 过滤 |
| P2 | Deposit 单位一致性约束 | 后端强制接收 wei 单位 + 校验 |
| P2 | i18n | 当前中英混杂 |

## 度量

| 项目 | 数字 |
|---|---|
| 修复的合约 bug | 12（rom 的 6 + 本轮新引入的 6） |
| 重抽的 ABI | 6 + 1 新增（TrafficCardNFT） |
| 前端新增 hooks | 4（useTrafficCards / useBurnCard / useIssueMonthlyCards / useTrafficCardCredit） |
| 后端新增 endpoint | 1（GET /api/deposit/:wallet/history） |
| 派出的 subagent | ~15 |
| 提交的 commit | 1（5月19日 NFT 升级修复 + UI） |

## 给下一轮的话

按 spec `2026-05-13-traffic-card-nft-redesign.md` 走，产品方向已经清晰：DeFi vault + ERC-1155 流量卡 + Yield Spread 收入。

技术栈基本不变，但合约层要大改（删 Payment/Oracle/FeeManager，新增 strategy adapter）。
