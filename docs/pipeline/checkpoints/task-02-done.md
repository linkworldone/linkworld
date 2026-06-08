# Task 02 — T1.5 ReentrancyGuard 探针 + 删无效 assembly（P1.5）

> 子项目 contracts(1/3) | 状态 DONE | 2026-06-08

### 产出文件
- packages/contracts/contracts/Payment.sol（删 initialize 第 28 行 `_reentrancyGuardInit()` 调用 + 删 L33-38 整个 `_reentrancyGuardInit()` 函数含 assembly sstore 段；-8 行）

### git commit
b55ef16 fix: 合约 T1.5 删无效 _reentrancyGuardInit + 保留 transient guard

### TDD
删死代码无新行为，不走红绿。验证靠 compile 绿 + 现有 30 测试不回归 + 实证（OZ ReentrancyGuardTransient 用 tload/tstore transient slot，每 tx 自动归零无需 init）。

### 测试结果
hardhat compile：Compiled 1 Solidity file successfully (evm: cancun)，0 error。hardhat test：30 passing / 0 failing，无回归（PM-01/PM-02 通过）。

### code-simplifier
净删 8 行死代码，是简化本身，无需额外处理。

### spec review
按 arch-review §七 A3（删无效 _reentrancyGuardInit）+ design §6.4 执行。transient 兼容判定：保留 ReentrancyGuardTransient（依据 Arbitrum ArbOS32+ 支持 EIP-1153），真·链上 TSTORE 验证按 §6.4 记给 T5 部署阶段。未越界（未动 guard 使用面/ERC20/分账）。

### 设计还原
合约无 UI。以 arch-review A3 落地替代：确认 Payment 继承 ReentrancyGuardTransient 且 guard 无需手动 init → 删 _reentrancyGuardInit（sstore 写持久 slot，guard 从不读 = 死代码）。删后 guard 仍有效。

### 复用检查
复用 OZ ReentrancyGuardTransient（继承，无需 init）；无新建。额外核实：grep 全合约 nonReentrant 0 命中——Payment 当前无函数用 guard，删 init 不改任何行为；payBill/发卡加 nonReentrant 属 T2/T3/T4。

### 设计稿对照
数值对照：删除行数 8 行（initialize 1 行调用 + 函数体 7 行）vs arch-review A3 指定的 L33-38+调用 ✅；compile error 0 个 vs gate 0 ✅；test 30 passing vs 基线 30 ✅ 无回归；transient 兼容判定 = 保留（ArbOS 32 起支持 EIP-1153）；新增网络配置 0（421614 留 T5）✅。

### 新增组件
无新增合约/接口/事件（纯删死代码）。

### 新增色值
无（合约任务）。
