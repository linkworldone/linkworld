# 测试阶段产出 — LinkWorld 全栈联调 (Round 2)

> 2026-04-15

## 编译验证

| 端 | 工具 | 结果 |
|----|------|------|
| 前端 | `npx tsc --noEmit` | ✅ 零错误 |
| 前端 | `pnpm dev` (Vite) | ✅ 启动成功 (localhost:3000) |
| 后端 | `go build ./...` | ✅ 编译通过 |
| 合约 | `npx hardhat test` | ✅ 2/2 测试通过 |
| 合约 | `npx hardhat compile` | ✅ 8 个 Solidity 文件编译通过 |

## 页面适配修复

TypeScript 编译发现 5 个页面/组件使用旧 hook API（`.mutateAsync` / `.isPending`），已全部修复：
- RegisterSheet.tsx → useRegister 新 API + useEffect 同步后端
- Deposit.tsx → useDepositMutation/useWithdrawMutation 新 API
- Billing.tsx → usePayBill 新 API + parseEther 计算
- BillDetail.tsx → 同 Billing 模式
- RegionDetail.tsx → useApplyNumber + apiClient 生成号码

## 已知限制（非阻塞）

1. 通知功能仍使用 mock（后端无通知 API）
2. 充值历史暂返回空数组（后端无 history 端点）
3. 单笔账单详情从列表缓存取（后端无单笔查询端点）
4. 运营商 dataRate/callRate 使用默认值（后端未返回）

## 结论

三端代码全部编译通过，无阻塞错误。基础设施层（ABI、API client、合约 hooks、双写补偿）已搭建完成，页面层已适配新 hook API。
