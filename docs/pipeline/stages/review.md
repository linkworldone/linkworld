# Stage: review — 代码审查 + 安全复审 + 设计审查（子项目 web 3/3）

> **状态**: PASS（代码审查 1 ❌ → fix 1a77bde → 复审 0 ❌；CSO/设计 0 ❌） | **日期**: 2026-06-10 | **Gate**: 3
> 审查对象：web 代码 diff main(05f7d2d)..b7f4873（83 文件 ~5146 增）。三审全员（web 有 UI，Design 路对口）：CSO 安全 / 代码审查(4 维) / 设计审查(视觉)。

## 一、总裁决：PASS（修复后 0 ❌）
代码审查查出 1 个 ❌ 资损（付账授权额/手续费基数二次算费→过额授权）→ 修复 → 复审 0 ❌。CSO + 设计审查首轮即 0 ❌。

## 二、❌ 阻塞项闭合
| # | 阻塞 | 来源 | 修复 | 复审 |
|---|------|------|------|------|
| W1 | 付账 approveAmount=totalAmount+calculateFee(totalAmount) + FeeBreakdown amount=totalAmount：totalAmount 已含 platformFee，在含费合计上再算费→过额授权(违 exact)+展示费偏大 | 代码审查 | fix 1a77bde：approveAmount=totalAmount(exact)；FeeBreakdown 用 bill.platformFee 直传(4 处)不再二次 calculateFee；Landing 2.5%→1.5%；修正固化错误值的 fee 测 | ✅ approveAmount===totalAmount exact、展示=后端 platformFee、RegionDetail 未误动、无 Hooks 违规 |

## 三、✅ 三审认可
- **CSO 安全（0 ❌）**：WalletAuth EIP-712 与后端 middleware.go 逐字节对齐(domain/字段/header/nonce 一次性消费/action 绑定)；资损四红线守住(exact approve 禁 infinite/精度 6 位不绕过/付账额读链不自算/不据 200 置终态/无二次 parseUnits)；无硬编码密钥/.env.local 未提交/无 XSS/无可疑依赖；信任边界正确(链上 source of truth + 200 仅意向)。
- **设计审查（0 ❌，综合 9.7/10）**：金色铁律 10(卡内金额 navy≈12:1/金仅深底 6.5:1，逐处核 8 个 AmountDisplay 无违规)、色值单一出口 10(旧 token 清零，唯一裸 HEX 是 RainbowKit accentColor 设计指定)、emoji→lucide 10(国旗数据豁免)、视觉态 9(三态/两步态/pending loading/拒签/空 vs 失败态完整)、一致性·无 slop 9、9 页+layout/shared/wallet 覆盖 10。
- **代码审查（修复后 0 ❌）**：plan T0-T12 对齐 design v2、105 测、四维通过。

## 四、⚠️ 非阻塞（ship 后/后续 polish）
- 补 FeeBreakdown fee 直传分支单测（修复核心路径目前靠 Billing.test 间接覆盖）。
- Landing 费率已改 1.5%（W1 顺带）；RegisterSheet 用 ui/input 收口未尽（样式已手对齐）；ui/button 原子默认 <44px 靠调用方补；EXPLORER_URL 双源；operatorApi startingPrice parseFloat 展示精度；.env.example 缺失（onboarding）；TwoStepAction stepper 金色措辞偏差。
- 设计审查：EmptyState 金图标在米白上偏淡（图标非文本不受 WCAG 约束）。

## 五、上线前 checklist（继承 + web）
- 合约真·上链(421614) → web deployments/arbitrum_sepolia.json 回填真实地址 + 端到端走查；web DONE 边界=本地 31337 全链路绿(D17 后置)。
- 31337 live-node 冒烟待补 hardhat-toolbox（合约子项收尾清了工具链）。
- WalletAuth chainId 需钱包链==后端 walletAuthChainID。
- 嵌套 packages/web/.git 空壳（无 remote）待清理。

## 六、结论
修复后 0 ❌，资损红线代码级守住 + 付账基数修复双验、深蓝金视觉 9.7/10、WalletAuth 对齐后端。可进 ship。非阻塞 ⚠️ 与上线 checklist 记录在案。

## 七、三审原始摘要
- CSO：0 ❌ + 2 ⚠️(.env.example/startingPrice)。
- 代码审查：1 ❌(W1 付账基数)+2 ⚠️(RegionDetail as any/EXPLORER 双源)→ fix → 复审 0 ❌。
- 设计审查：0 ❌ + 3 ⚠️(Landing 2.5%[已修]/RegisterSheet ui-input/Button 44px)，9.7/10。
