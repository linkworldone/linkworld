# stage-test-done.md — Link World web(3/3) T12 测试扫尾 + 构建 + guardRule 终检

> 子项目 web(3/3) | 阶段 implement T12（最后一棒） | 状态 **DONE** | 2026-06-10
> 范围：只补测试/验证/产文档，**未改任何业务逻辑**（接链 T1-T10 / 换肤 T6-T11 代码原样）。

---

## 1. Phase 结果（最终全量验证）

| 验证 | 命令 | 结果 |
|------|------|------|
| 类型 | `npx tsc --noEmit` | **0 error** ✅ |
| 构建 | `npm run build`（`tsc -b && vite build`） | **exit 0，built ✓** ✅（仅 chunk>500kB 提示，非错误） |
| 测试 | `npx vitest run` | **23 files / 105 passed**（前序 101 + T12 新增 4，无回归）✅ |
| 构建（指向 31337） | `VITE_CHAIN_ID=31337 npm run build` | **exit 0** ✅（验证 31337 config 加载/链切换无编译·导入错误） |
| coverage | `@vitest/coverage-v8` | **未安装 → 按 spec 跳过**（不擅自加依赖） |

T12 新增测试：`src/config/contracts.test.ts`（CFG-01..04，+4），固化 31337 配置单一出口 +
USDT 6 位精度根来源 + 零地址/未知链兜底（替代被阻塞的 live 31337 冒烟的 web 侧可解析部分）。

---

## 2. 资损红线覆盖审计（对照 design §10 + arch-review B4 + 各 task checkpoint）

| # | 资损/关键红线 | 覆盖测试（文件 · 用例） | 状态 |
|---|---------------|------------------------|------|
| R1 | **精度 6 位**（旧 18 差 10^12） | `utils/format.test.ts` PREC-01：`parseUnits('1.5')===1_500_000n`、`formatAmount(10_000_000n)==='10.00'`、边界 0 / 大额不丢精度 | ✅ 已覆盖 |
| R2 | **MIN_DEPOSIT 精度+值对齐** | `config/constants.test.ts` PREC-02/03：`MIN_DEPOSIT_USDT===10n*10n**6n`(=10_000_000n)、`SUPPORTED_CURRENCIES===['USDT']` | ✅ 已覆盖 |
| R3 | **USDT 精度根来源（从链读不硬编码）** | `config/contracts.test.ts` CFG-02：`getUsdtDecimals(31337)===6`、同源 hardhat.json | ✅ **T12 新补** |
| R4 | **approve exact 两步态**（禁 infinite） | `components/shared/TwoStepAction.test.tsx` TSA-01/04：`approve args===[SPENDER, AMOUNT]`（exact，禁 MaxUint256） | ✅ 已覆盖 |
| R5 | **approve 成功+action 失败→approved-idle，绝不 re-approve** | TwoStepAction TSA-03 + `derivePhase` 纯函数（approvedOnce+error→approved-idle，`approveWrite.not.toHaveBeenCalled()`）；allowance≥amount 跳步 TSA-02 | ✅ 已覆盖 |
| R6 | **付账 approve 额 = amount + calculateFee（读链不自算）** | `hooks/contracts/useFeeManager.test.ts` FEE-01：`calculateFee` 走链（`functionName==='calculateFee'`，args=[amount]，不前端算 amount*rate） | ✅ 已覆盖 |
| R7 | **手续费率读链**（禁写死费率） | useFeeManager FEE-01/02：`getFeeRate` 读链基点 150n→1.5%（/10000）；读链失败→label undefined（兜底 `--` 不写死）；loading→skeleton | ✅ 已覆盖 |
| R8 | **对账不据 200/txHash 置终态** | `services/api/depositApi.test.ts` REC-01/02：充值/提现仅写 pending 意向，不把 tx_hash 当终态依据；余额取链上 `getDepositAmount`（`hooks/useDeposit.balance.test.ts` REC-04，非后端自述） | ✅ 已覆盖 |
| R9 | **pending 不染绿 / reorg vs revert 区分** | `components/shared/TxStatusBadge.test.tsx` BADGE-01：pending『处理中』无 success 类、confirmed 才染绿、failed 区分 reorg『已回退』/revert『失败』、缺省不崩 | ✅ 已覆盖 |
| R10 | **锁仓边界 `>=`（到期即可提）** | `components/shared/LockCountdown.test.tsx` DEP-01a..d：`expiry=0→none`、`now<expiry→locked`、**`now===expiry→unlocked`（边界，不可用 `<`）**、`now>expiry→unlocked` | ✅ 已覆盖 |
| R11 | **WalletAuth EIP-712 会话签名** | `services/api/signedPost.test.ts` AUTH-01/02/03 + 一次性 nonce：GET nonce→signTypedData(域/字段/message 逐项对齐后端 middleware.go)→带签名头 POST；action 各端点对齐；拒签不发 POST 抛 `WalletAuthRejectedError`；每次写取新 nonce 重签（消费式） | ✅ 已覆盖 |
| R12 | **billingApi.toBill bigint（禁 parseFloat 浮点）** | `services/api/billingApi.test.ts` / `billingApi.payIntent.test.ts`：total bigint 加减无浮点误差 | ✅ 已覆盖 |
| R13 | **无二次 parseUnits（双重缩放）** | Billing.tsx / BillDetail.tsx：`totalAmount` 已是 6 位最小单位字符串直接转 bigint，**不再 parseUnits**（代码注释+grep 确认，见 §3） | ✅ 代码层守住（grep 证明） |
| R14 | **31337 地址单一出口 + 未部署兜底** | `config/contracts.test.ts` CFG-01/03/04：7 合约地址全取自 hardhat.json；未知 chainId 抛错；零地址抛错（不把 0x0 当合约调） | ✅ **T12 新补** |

**结论**：资损红线 R1-R12、R14 均有单元/组件测试固化；R13 由代码（直接 BigInt 转换 + 注释）+ grep 双保险。
审计未发现未覆盖的资损/关键路径缺口；T12 仅就 R3/R14（config 单一出口 + 6 位精度根来源，原无专测）新补 4 测。

---

## 3. guardRule 终检（grep 全 `packages/web/src` 清零，逐项证明）

命令：`rg -n <pattern> packages/web/src`（`--glob '!*.test.*'` 排除测试黑名单镜像）

| 模式 | 结果 |
|------|------|
| `#3b82f6` | **CLEAN** |
| `#0a0a14` | **CLEAN** |
| `#0f0f1a` | **CLEAN** |
| `#8b5cf6` | **CLEAN** |
| `#06b6d4` | **CLEAN** |
| `bg-brand-blue` / `brand-blue`(任意形式，非测试) | **CLEAN** |
| `to-brand-purple` / `brand-purple`(非测试) | **CLEAN** |
| `brand-cyan` / `surface-gradient` / `surface-secondary`(非测试) | **CLEAN** |
| `font-orbitron` / `orbitron`(任意大小写，非测试) | **CLEAN** |
| `"Inter"` font-family | **CLEAN**（已统一走 `--font-sans` Geist） |
| 装饰 emoji（🏠📱💰📄🎟🔔🔍🌐💳⚠🔒🚀✨✅❌… 非测试） | **CLEAN**（全 lucide） |
| 写死费率（交易路径） | **CLEAN**（费率走 `useFeeManager.getFeeRate/calculateFee` 读链） |
| `parseEther` / `10n**18n` / `decimals:18`(非 nativeCurrency) | **CLEAN** |
| 二次 `parseUnits`（双重缩放） | **CLEAN**（Billing/BillDetail 直接 BigInt 转换，注释标红线） |

**合法保留（非违规）**：
- `pages/skin.test.tsx`：`font-orbitron` 等出现在 `DEAD_TOKENS`/`DECO_EMOJI` 黑名单数组里——是 grep 清零的**运行时镜像断言**（`expect(html).not.toContain(...)`），非真实使用。
- 国旗 emoji（`services/mock/data.ts` regionFlags / operatorApi `FLAG_MAP`）：保留为**国家标识数据**，不是装饰图标（spec 明确除外）。
- `config/chains.ts` `nativeCurrency.decimals: 18`：ETH 原生币精度，合法（USDT 走 `usdtDecimals=6`）。
- `services/api/{usageApi,operatorApi}.ts` `0.05/0.02`：mock **展示估算占位**（注释「暂用默认值，后端未返回」），非链上交易费率路径。

**9 页 100% 深蓝金**：T11/task-12 checkpoint 已确认 9 页 + layout/shared/wallet 全覆盖，emoji→lucide 全量，色值单一出口（`:root`）。T12 grep 复跑维持清零。

---

## 4. 本地 31337 冒烟结果

**状态：web 侧 config 加载/链切换 PASS；live node 全链路冒烟 BLOCKED（环境，非 web 缺陷）。**

### 已跑通（web 侧可解析部分）
1. `VITE_CHAIN_ID=31337 npm run build` → **exit 0**：证明 wagmi.ts `isLocalChain` 切换 + `config/contracts.ts` 静态 import `deployments/hardhat.json`（7 合约代理地址 + usdt + usdtDecimals=6）编译·打包无误。
2. `config/contracts.test.ts`（CFG-01..04）：运行时断言 31337 地址全取自 hardhat.json、`getUsdtDecimals(31337)===6`、未知链/零地址抛错兜底。
3. `packages/web/src/config/deployments/hardhat.json` 已存在（确定性 Hardhat 地址 + `usdtDecimals:6`），web config 链路可达。

### 阻塞步骤 + 原因
尝试 `cd packages/contracts && npx hardhat node`（后台）→ **启动失败**：
```
Error: Cannot find module '@nomicfoundation/hardhat-toolbox'
  at hardhat.config.ts:2:1
```
- **根因**：本分支 `web/deep-blue-gold-refactor` 在 commit c7616e5「web 分支清除后端子项」时清掉了 contracts/backend 子项工具链，`packages/contracts` 的 `@nomicfoundation/hardhat-toolbox` 等 devDeps **未在本分支安装**。
- **为何不强跑/不修**：补装 contracts 依赖会触碰**禁区**（合约子项目），超出 web 测试 subagent 范围；且 spec 明确「环境跑不起 node → 记录步骤+阻塞原因，不强跑」。
- **完整冒烟步骤（待 contracts toolchain 就绪后可执行）**：
  1. `cd packages/contracts && pnpm install`（补 hardhat-toolbox 等）
  2. `npx hardhat node`（后台，:8545，chainId 31337）
  3. `npm run deploy:local` → 部署 7 合约 + MockUSDT(6 位)
  4. 部署地址同步到 `packages/web/src/config/deployments/hardhat.json`（确定性地址，当前已是该组）
  5. `cd packages/web && VITE_CHAIN_ID=31337 npm run dev` → 浏览器连 MetaMask(31337) 跑充值两步态→pending→读链 confirmed / 付账 / 提现 / 锁仓边界 / 手续费读链 / WalletAuth 会话签名

### Arbitrum 421614 端到端
**明确后置阻塞**于合约真·上链（`deployments/arbitrum_sepolia.json` 当前为零地址占位）。按 design §11/§12「web DONE 边界」：**web DONE = 本地 31337 全链路绿**；421614 端到端(D17) + 对账三态真链行为 = 后置强制验收，不计入 web DONE 也不成孤儿。

---

## 5. 发现的 bug

**无。** 全程未发现需修的业务逻辑 bug；tsc/build/vitest 全绿，grep 清零达标，资损红线测试全覆盖。

---

## 6. 已知限制

1. **live 31337 全链路冒烟阻塞于 contracts toolchain（本分支未装），非 web 缺陷**（详见 §4）。web 侧 config 加载/链切换以 build + CFG-01..04 单元测固化。
2. **421614 真链端到端后置**，阻塞合约上链（design 既定边界）。
3. **coverage 未跑**：`@vitest/coverage-v8` 未安装，按 spec 不擅自加依赖（如需，`pnpm add -D @vitest/coverage-v8` 后 `npx vitest run --coverage`）。
4. **遗留**：嵌套 `packages/web/.git` 空壳仓库（无 remote）仍在，T12 未碰；T0-T11 改动均已归并入外层 linkworld 仓库，待主 Agent 与用户确认后清理（与 task-12 checkpoint §⚠️ 一致）。

---

## 7. 三自检

1. **承诺**：T12 只补测试/验证/产文档，未改任何业务逻辑（接链 T1-T10 / 换肤 T6-T11 源码 0 改动）；未碰合约/后端/pipeline.json/嵌套 web .git；git add 仅限实际改动的 web 测试文件 + 本 checkpoint。✅
2. **交付物存在**：`packages/web/src/config/contracts.test.ts`（新增，4 测通过）+ 本文件 `docs/pipeline/checkpoints/stage-test-done.md`。✅
3. **正确连接**：新测试 import 自 `./contracts` + `./deployments/hardhat.json`（真实生产路径，非 mock 复制），断言对齐 hardhat.json 实际内容；全量 23 files/105 passed 含新测；覆盖审计逐项指向真实测试文件:用例号。✅

---

## 8. 结论：可进 test 阶段

implement T0-T12 全部 DONE：tsc 0 / build exit 0 / **23 files 105 passed** / grep guardRule 全清零 / 9 页 100% 深蓝金 / 资损红线 R1-R14 测试覆盖。web DONE 边界（本地 31337 全链路绿）中 **web 侧可解析部分已绿**，live node 全链路冒烟阻塞于 contracts toolchain（环境，已记录完整复跑步骤），421614 端到端后置阻塞合约上链。**可进 test 阶段。**
