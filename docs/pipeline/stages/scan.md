# Stage: scan — web 基线扫描 + 接链改动面（子项目 web 3/3）

> **状态**: completed (DONE_WITH_CONCERNS) | **日期**: 2026-06-09 | **Gate**: 0 | **子项目**: web(3/3)
> 对象：packages/web(React)。复用现有基线(docs/design/linkworld/*)，聚焦接链/对账/换肤 delta。

## 产出（commit 7d334de）
- docs/design/linkworld-web/web-alignment-surface.md（接链/对账改动面，最关键）
- docs/design/linkworld-web/theme-migration.md（深蓝金换肤 delta）
- 复用：docs/design/linkworld/{DESIGN.md,color-mapping,components,project-scan,utils}.md（UI 结构未变）

## 🔴 关键发现（接链改动规模大）
1. **链配置全失配**：chains.ts 仅 0G(16600/16601)+hardhat 缺 421614；contracts.ts 仅 31337+16601(全0)，31337 地址也漂移(与 deployments/hardhat.json 对不上)；wagmi.ts 非本地落 0G 需改 Arbitrum。
2. **ABI 旧 0G 快照、selector 已变(breaking)**：abis/Deposit.ts deposit 仍无参 payable、缺 getLockExpiry；abis/Payment.ts payBill 仍 payable。真合约已 deposit(uint256)/payBill 去 payable/createBill onlyOracle。必须从 artifacts 全量重生成 + 同步 hooks/contracts/*。
3. **abis/index.ts 漏导出 OracleABI**（PRD 验收 D11）。
4. **USDT 接入空白**：现全程原生币(value+parseEther)，无 ERC20/allowance/approve。需 exact-amount approve 两段式(禁 infinite)+Deposit 两步态。
5. **精度资损红线**：format.ts 默认 18 位、MIN_DEPOSIT_USDT=100*10^18(位数+起充双错)、PLATFORM_FEE_RATE=0.025 写死、currency 写死 ETH。USDT 6 位，漏改任一处金额错 10^12 倍。
6. **对账三处降级(breaking)**：useDeposit/useBilling 双写(recordDeposit/recordWithdraw 凭 txHash/recordPayment 标记已付) → 改 pending 意向+事件驱动确认，UI 不能据 200 显示已付/已到账。
7. **WalletAuth 缺失**：api/client.ts 裸 axios 无签名头，受保护写端点需加 EIP-712/191 钱包签名拦截器（与后端对齐）。
8. **换肤覆盖面**：brand-blue×12 文件、brand-purple×4、surface-gradient×3、emoji 7+ 处(TabBar/AppLayout/Dashboard/Services/Notifications/Landing/Billing)、AmountDisplay 默认色；tailwind.config 硬编码 HEX+orbitron 需收口 CSS 变量单一出口。

## ⚠️ 遗留/风险
- **合约真·上链未做**：deployments/arbitrum_sepolia.json 不存在，421614 无真实地址 → web 421614 段只能全 0 占位，**Arbitrum 端到端验收(PRD D17)阻塞**于合约先上链；本地 31337 可先验全链路。
- 31337 地址漂移 → implement 做「deployments json → contracts.ts」单一出口脚本，杜绝手抄。
- WalletAuth 签名格式未定(EIP-712 vs 191/nonce 来源)，implement 期与后端 signatures.go 敲定。
- 换肤⇄接链同文件(main.tsx/wagmi.ts/Deposit.tsx/Cards.tsx)两线都碰，需串行/协调。

## 移交 design
带着两份 delta + 搁置深蓝金 DESIGN.md + 两份 handoff，design 阶段敲定：深蓝金视觉规范(复用 b18cf37 DESIGN.md)+接链交互稿(USDT approve 两步态/锁仓倒计时/Cards 双Tab/手续费读链)+对账降级 UI 形态。web 有 UI，design 走 UI 路径(design-consultation/shotgun)。
