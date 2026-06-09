# requirement 阶段 · PRD — Link World 全栈重构（Round 1 · 修订版）

> 状态：completed | 日期：2026-06-08 | Gate：1（设计） | Round：1
> 产出方式：brainstorming（发散）+ office-hours 收敛纪律协作产出，决策已拷问到可验收条件。
> **本版取代 2026-06-07 旧 PRD**（旧版前提 D1/D2「只动前端、后端合约不配合」已作废，详见 §〇）。

---

## 〇、为什么重审（前提反转）

旧 PRD（2026-06-07）核心约束是 **D1「换肤+对齐、只动前端」+ D2「后端/合约不配合，依赖项只能前端降级」**。

2026-06-08 后端 + 合约代码已合并（commit `7ef9677`）并部署到 og_testnet（chainId 16602），逐条核对真实代码后：

- 旧 PRD 多数「前端降级项」其实链上/后端已有真实链路（锁仓到期、NFT 发放、提取/历史/通知）。
- 用户在重审中决定**进一步把合约切到 Arbitrum + ERC20 USDT**（见 §二决策表）。

→ **D1/D2 彻底反转**。本轮从「web 端换肤」升级为**跨 3 子系统的全栈改造**（合约 + 后端 + web）。

---

## 一、用户故事（主流程不变，落地链路换真实）

数字游民 / Web3 用户：连接钱包（Arbitrum 网络）→ 邮箱注册 → **approve USDT → 存保证金（≥10 USDT，锁仓 30 天）→ 锁仓满自动获 NFT 流量卡 → 选目的地运营商申请号码 → 用量产生账单 → USDT 支付（含链上动态手续费）**。

与旧版差异：充值/支付从「原生币 `msg.value`」改为「ERC20 USDT（approve + transferFrom）两段式」；所有数值取链上真实值。

---

## 二、关键决策（本轮重审逐项确认）

| # | 决策点 | 结论 | 备注 |
|---|--------|------|------|
| **R1** | 本轮范围 | **换肤 + 全量接通后端合约**（全栈） | 推翻旧 D1/D2 |
| **R2** | 合约目标链 | **迁移到 Arbitrum Sepolia（chainId 421614，测试网）** | 验收用测试网，可水龙头免费走真交易 |
| **R3** | 保证金/支付币种 | **ERC20 USDT**（对齐文档「10u 起」） | `Deposit.sol`/`Payment.sol` 从 `payable`→`transferFrom`+`approve` 改写；min 10 USDT 可链上 `require` 强约束 |
| **R4** | 测试 USDT 来源 | **自部署 mock ERC20**（Arbitrum Sepolia 无官方 USDT） | 合约阶段确定 mock 地址 + 精度（USDT 通常 6 位小数，需核对 mock） |
| **R5** | 手续费事实源 | **读链上 `FeeManager.getFeeRate()`（动态基点，默认 150=1.5%，分母 10000）** | 前端**不再写死** `PLATFORM_FEE_RATE` |
| **R6** | 锁仓校验 | **动态**：读 `Deposit.getLockExpiry(addr)` 做 Withdraw 禁用 + 倒计时 | 推翻旧版「只能静态提示」 |
| **R7** | NFT 自动发放 | 接真实链路（`Deposit.issueMonthlyTrafficCards`/`mintTrafficCard` + `TrafficCardNFT.mint`） | 移除 Admin 发卡按钮，改「锁仓满自动发放」说明 + 真实状态展示 |
| **R8** | 提取 / 历史 / 通知 | 接真实 API：`/api/withdraw`、`/api/deposit/:wallet/history`、`/api/notification/send`（原 mock 转真实） | |
| **R9** | 实体 SIM 领取（文档4.6） | **保留前端降级占位**（后端/合约无任何端点） | Cards 页「SIM 领取」Tab：选国家+收件表单 → 提交 toast + `pendingSync`(localStorage) + 「全球通即将推出」 |
| **R10** | 本地号码注册（文档3.3） | **仍邮箱注册**（后端 `User` 无 phone 字段） | 留后续 Round |
| **R11** | 开发链路 | **本地 hardhat(31337) 先行开发验证 → Arbitrum Sepolia(421614) 部署/验收** | web `chains.ts`/`contracts.ts` 同配 31337+421614，`VITE_CHAIN_ID` 切换（沿用现有 wagmi.ts 模式） |
| **R12** | 视觉基调 | **深蓝渐变底 + 米白/浅色卡片 + 金色点缀**（沿用旧 D3/D4，不变） | 金色仅金额/主CTA/激活态/品牌强调；emoji→lucide 全量替换 |
| **R13** | pipeline 形态 | **建议拆 3 子项目串行**（合约迁移 → 后端对齐 → web 重构接通），主 Agent 重审 pipeline | 见 §六 |

---

## 三、范围（按子系统）

### 子系统 ①：合约（`packages/contracts`）
- `Deposit.sol`：`payable`/`msg.value` → ERC20 USDT（`approve` 前置 + `transferFrom`）；保留 `_lockExpiry` 30 天锁仓；`require(amount >= 10 USDT)` 强约束最小额。
- `Payment.sol`：账单支付改 USDT `transferFrom`；手续费走 `FeeManager.calculateFee`。
- 自部署 mock USDT（ERC20）测试币（R4）。
- 部署脚本 + `deployments/` 产出 Arbitrum Sepolia(421614) 新地址/ABI；本地 31337 同步可部署。
- ⚠️ ERC20 改写涉及 `IDeposit`/`IPayment` 接口签名变更，需同步 interfaces + 事件。

### 子系统 ②：后端（`packages/backend`）
- `configs/deployments.json`：chainId/rpcUrl 换 Arbitrum Sepolia + 7 合约新地址 + mock USDT 地址。
- `blockchain/event_sync.go`、`signatures.go`：跟 ERC20 改动（事件签名、金额单位）对齐。
- ⚠️ 现存 RPC 不一致需顺带核对（hardhat `evmrpc-testnet.0g.ai` vs backend `evm-testnet.0g.ai`），统一为 Arbitrum Sepolia RPC。

### 子系统 ③：web（`packages/web`）
- 深蓝金换肤（9 页 + layout/shared/wallet 组件，emoji→lucide）。
- `config/chains.ts`：新增 Arbitrum Sepolia(421614) 链定义；`contracts.ts`：填 421614 + 31337 真实地址；`abis/index.ts` 补漏导出 `OracleABI`。
- `config/constants.ts`：删除写死的 `PLATFORM_FEE_RATE`（改读链）；`SUPPORTED_CURRENCIES` 改为 USDT（去掉 A0GI/ETH 框架）；min 10 USDT 与链上一致。
- 充值/支付加 **USDT approve 两段式流程**（检测 allowance → approve → deposit/pay）。
- Deposit 页 Withdraw 区：读 `getLockExpiry` 做动态禁用 + 倒计时。
- 费用明细（1.5% 手续费）在「申请号码弹层」+「支付账单页」按链上费率展示。
- Cards 页：双 Tab（流量卡 / SIM 领取占位 R9）；移除 Admin 发卡按钮 → 自动发放说明；不可转卖 / 销毁后 30 天规则文案。

---

## 四、范围红线（绝不碰 / 留后续）

- ❌ 本地号码注册（后端无 phone 字段，R10）→ 留后续 Round，本轮保持邮箱。
- ❌ 实体 SIM **真实**领取流程（无端点，R9）→ 本轮仅前端表单 + pendingSync 占位。
- ❌ 不上 Arbitrum One 主网（R2 锁定 Sepolia 测试网，不动真实资金）。
- ❌ 不改主流程的业务路由结构、不新增底部导航第 6 项。

---

## 五、验收标准（Acceptance Criteria）

### A. 合约（子系统①）
1. `Deposit`/`Payment` 已改 ERC20 USDT 模型：存款/支付走 `approve`+`transferFrom`，无 `payable`/`msg.value` 残留（grep 核查）。
2. `Deposit` 强约束 `amount >= 10 USDT`（含单元测试覆盖 <10 拒绝、=10 通过）。
3. 部署脚本可在 **本地 31337 + Arbitrum Sepolia 421614** 两处成功部署，产出地址写入 `deployments/`。
4. mock USDT(ERC20) 已部署，精度/符号与前端展示一致。

### B. 后端（子系统②）
5. `configs/deployments.json` chainId=421614、RPC=Arbitrum Sepolia、7 合约地址 + USDT 地址全部真实非零。
6. event_sync 能监听 Arbitrum Sepolia 上 ERC20 改写后的真实事件（启动无报错、能落库）。

### C. web 视觉（子系统③ · 机器+人工核查）
7. 9 页 + layout/shared/wallet 组件 100% 应用深蓝金主题，无遗漏页。
8. grep 查无旧色值残留（`#3b82f6`/`#0a0a14`/`#0f0f1a`/`bg-brand-blue`/`to-brand-purple`/`#8b5cf6`/`#06b6d4`）。
9. grep 查无 emoji 当图标残留（全部 lucide）。
10. 卡片米白/浅底、背景深蓝渐变、金色仅金额/CTA/激活态/品牌强调；色值单一出口（语义 token/CSS 变量）。

### D. web 功能接通（子系统③ · 端到端）
11. `chains.ts` 含 421614、`contracts.ts` 421614+31337 地址非零、`abis/index.ts` 已导出 `OracleABI`。
12. 手续费 UI 显示值 = 链上 `FeeManager.getFeeRate()` 实读值（不是写死的 0.015）；申请号码弹层 + 账单页均展示。
13. 充值/支付走 USDT approve 两段式：allowance 不足时先 approve，再 deposit/pay，UI 有两步态。
14. Deposit 页 Withdraw 区：锁仓未到 → 按钮禁用 + 显示剩余倒计时（读 `getLockExpiry`）；到期 → 可提取。
15. Cards 页：双 Tab；Admin 发卡按钮已移除并替换为自动发放说明；NFT 自动发放状态来自链上真实数据；不可转卖/30天规则文案到位。
16. SIM 领取表单可填可提交，提交后 toast 成功 + 写 `pendingSync`。
17. **主流程端到端走通**（连钱包→注册→approve→存 USDT→锁仓→自动发卡→选运营商→账单→USDT 支付）在 **本地 31337 + Arbitrum Sepolia** 均跑通 —— 由 test 阶段端到端走查验收。

### E. 设计细节
18. 金色具体色值（沉稳金 `#D4AF37` / 香槟金 `#F0C75E`）、字体方案（修 Inter/Geist 冲突、是否启用 Orbitron）、深蓝渐变精确停靠点、米白卡片色阶 —— **以 design 阶段产出为准**；本验收只要求「已应用 design 阶段确定的规范」。

---

## 六、移交说明（给主 Agent）

⚠️ **本轮已超出原 pipeline「只动 web、10 阶段」的建立前提**，强烈建议主 Agent 重审 pipeline 配置：

1. **建议拆 3 个子项目串行**：①合约迁移（Arbitrum+USDT 改写+部署）→ ②后端对齐（配置+event_sync）→ ③web 重构接通（换肤+接链+approve 流程）。各自独立 spec→plan→实现，否则一个设计/任务拆解阶段扛不住，且资损/合约审计风险压不住。
2. **新增风险面**需在 arch-review 重点审：ERC20 approve 重入/授权额度、USDT 精度（6 位小数 vs 18）、锁仓与自动发卡的时序、Arbitrum 与 0G 的 gas/确认差异、合约升级（UUPS proxy）兼容。
3. 合约改写属**资损敏感**，arch-review 阶段务必含 security-review（合约审计 checklist）。

## 七、移交 design 阶段（待决细节）

design 阶段需产出并锁定：
1. 金色色值（沉稳金 vs 香槟金）+ 应用规则（明度/对比度）。
2. 深蓝渐变精确参数（`linear-gradient(135deg, #0C2340 0%, #1E40AF 50%, ...)` 完整停靠点）。
3. 米白/浅色卡片色值 + 卡内文字色阶（主/次/弱文本深色映射）。
4. 字体方案（修 Inter/Geist 冲突；Orbitron 取舍）。
5. 色值出口统一（Tailwind 语义 token 与 oklch CSS 变量收敛为单一来源）。
6. lucide 图标映射表（每处 emoji → 对应图标）。
7. 缺失原子组件（Card/Input/Dialog）补齐 vs 替换 `ui/*` 决策。
8. **USDT approve 两段式流程的交互稿**（allowance 检测态、approve 等待态、deposit/pay 态）。
9. **锁仓倒计时 + 自动发卡状态**的展示形态。
10. 9 页 + 组件新主题高保真。

---

## 八、留待后续 Round（本轮不做，仅记录）

- 「提交本地号码」注册（需后端加 phone 字段）。
- 实体 SIM 卡**真实**领取流程（需后端接口；本轮仅前端表单 + pendingSync 占位）。
- Arbitrum One 主网部署（本轮锁 Sepolia 测试网）。
</content>
</invoke>
