# Task 13 — T12 测试扫尾 + 构建 + guardRule 终检（web 3/3，implement 末棒）

> 子项目 web(3/3) | 状态 DONE | 2026-06-10

### 产出文件
- packages/web/src/config/contracts.test.ts（新增 4 测：CFG-01..04 单一出口加载/usdtDecimals=6/未知链兜底/零地址保护）
- docs/pipeline/checkpoints/stage-test-done.md（test 阶段验证结论：覆盖审计+grep 证明+冒烟阻塞+三自检）

### git commit
f5404da test: web T12 测试扫尾 + 构建 + guardRule 终检

### TDD
补 config 单一出口+6 位精度根来源专测（原无）；其余红线 T0-T11 已覆盖，T12 审计确认无缺口。

### 测试结果
npx tsc --noEmit 0 error。npm run build ✓（含 VITE_CHAIN_ID=31337 绿）。npx vitest run：23 files / 105 passed（前序 101+新 4，无回归）。主 Agent 已独立确认。coverage 插件未装按 spec 跳过。

### code-simplifier
仅补测+审计文档，无业务改动。

### spec review
按 design §10 + arch-review §五红线执行。14 条资损红线覆盖审计全绿（精度/approve exact 两步态/对账不据 200/锁仓边界 >=/手续费读链/WalletAuth EIP-712/billingApi bigint/无二次 parseUnits/单一出口）。未改业务逻辑、未发现 bug。

### 设计还原
guardRule 终检 grep 全 src 清零达标（旧色值/emoji/写死费率/二次 parseUnits 全 CLEAN，国旗数据/chains ETH decimals 合法保留）；9 页 100% 深蓝金。

### 复用检查
复用全部前序测试设施 + contracts 单一出口；无新依赖。

### 设计稿对照
数值对照：14 红线覆盖（2 新补 CFG-02/CFG-01·03·04）✅；105 测/tsc 0/build ✓（含 31337）✅；grep 旧色值/emoji 0 残留 ✅；usdtDecimals=6 单一出口 ✅。

### 新增组件
无（补测 + 审计文档）。

### 新增色值
无。

### ⚠️ 遗留 / 已知限制
- **31337 live-node 冒烟阻塞**：npx hardhat node 报 Cannot find module @nomicfoundation/hardhat-toolbox（合约子项收尾 commit c7616e5 清工具链所致，非 web 缺陷）；补装属合约范围。web 侧 build(31337)+config 加载断言已绿，完整复跑步骤记 stage-test-done.md §4。
- **Arbitrum 421614 端到端**：后置阻塞合约真·上链（web DONE 边界=本地 31337 绿）。
- **嵌套 packages/web/.git 空壳（无 remote）**：待主 Agent 与用户确认清理。
- web implement T0-T12 全完成（13 任务）；进 test 阶段（基础验证 T12 已做，正式 test 走 vitest 全量+构建+visual 抽查）。
