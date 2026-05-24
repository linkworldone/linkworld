# LinkWorld 经验教训

## Round 2 教训（2026-05-19）

### 1. wagmi `useWaitForTransactionReceipt` 的 isPending 陷阱
**症状**：hash 为 undefined 时（用户还没点交易），它会返回 `isPending: true`，导致 UI 误显示"Confirming on chain"。

**正确写法**：必须加 hash 守卫
```typescript
const { isPending: rawConfirming, isSuccess } = useWaitForTransactionReceipt({ hash });
const isConfirming = !!hash && rawConfirming;
```

**本轮踩了 2 次**：useTransactionFlow 里修过一次，底层合约 hook 又复发一次。所有用 useWaitForTransactionReceipt 的地方都要加 guard，**没有例外**。

### 2. 前后端数据单位不一致是 bug 高发区
**症状**：用户输 "1" (ETH)，前端合约调用走 parseEther 正确，但 recordToBackend 直接发字符串 "1"，后端存 "1"，读回当 wei 处理 → 显示 0.00。

**规则**：跨边界的金额传输**永远用 wei 单位字符串**，不要混用人类可读的小数字符串。后端建议加输入校验，拒绝小于 10^12 的金额（一定是单位错了）。

### 3. MetaMask gas 21M 是 estimateGas 失败的信号
**症状**：MetaMask 显示 "网络费 不可用"，gas limit 默认 21000000，提交时 Hardhat 报 "exceeds cap 16777216"。

**根因**：合约调用要 revert，estimateGas 失败，MetaMask fallback 到硬编码 21M。
**正确反应**：不是改 gas，是**查为什么会 revert**（已注册？参数错？合约地址错？）。

### 4. rom 模式：本地不编译就 push
**已观察 3 次**（rom TrafficCard / rom upgrade NFT / 还有更早的）：
- 第一次：5 处编译错
- 第二次：6 处编译错（含上次没修的旧 bug）

**对策**：
- 加 git pre-push hook：`pnpm compile && pnpm test`
- 或 GitHub Actions PR 拦截
- 不要相信外部 commit "应该能跑"

### 5. Pull 前 stash 比 Merge 干净
本轮在 rom push 后处理冲突：先 `git stash push -u -m` 保护本地修改，再 pull，然后丢弃 stash 重新基于新代码修复。比 merge 三方解决更线性。

### 6. 后端 mock 数据要么真要么 0，不要随机
`OperatorAPISimulator.GetUsage` 每次返回 `rand.Intn(10000)`，造成前端每刷新数字都变，用户误以为系统在乱算账。**所有 mock 应该可预测**：要么读 DB（真），要么返回 0（明显的空状态）。

### 7. ReentrancyGuard 在 UUPS 代理下要用 Transient 版
**症状**：rom 用了非 upgradeable 的 `ReentrancyGuard`，又调用了不存在的 `_reentrancyGuardStorageSlot()`。

**正确**：OZ 5.x 在代理模式下用 `ReentrancyGuardTransient`（EIP-1153 瞬态存储，无 constructor，proxy 安全）。

### 8. Maturity transformation + 硬锁仓 = 安全 DeFi 模式
**洞察**：传统 DeFi vault 怕挤兑，因为用户可秒退。但加上"硬锁仓"机制后，挤兑变成可预测的现金流——平台知道未来每天有多少钱要出去，可以提前规划 PT 到期日。

**结论**：池化 + 硬锁是最优组合，比"每用户独立锁档"APR 提升 50%+。

### 9. 三端独立部署的接口治理
本轮的链上 / 后端 / 前端三向同步成本高：合约地址变 → 前端 contracts.ts 改 + ABI 重抽 + 后端 model 跟着改。

**建议**：
- 部署脚本自动写入 `deployments/<chain>.json`
- 前端 / 后端 build 时从这个 JSON 读地址，不写死
- ABI 通过 codegen 从 artifacts 自动生成，不手抽
