# 代码审查产出 — LinkWorld 全栈联调 (Round 2)

> 2026-04-15

## 审查结果

### Critical (3) — 已全部修复 ✅
1. **vite proxy 端口冲突** — proxy target 3001→8080，port 3000→删除（用默认 5173）
2. **useRegister 未自动链 backend sync** — 加 useEffect + useRef 存 email，自动调 backendSync
3. **Withdraw handler 写无意义数据** — 移除 depositService.Deposit("0")，改为日志记录

### Important (8) — 修复 2 个，其余记录
4. **operatorApi BigInt 浮点精度** — ✅ 改用 viem parseEther
5. **CORS 端口对齐** — ✅ 后端加 localhost:3000
6. operatorApi 全量拉取 — 接受，加 TODO（后端无筛选端点）
7. DepositInfo.currency 类型 — 接受（"ETH" 在 union 中）
8. import 路径不一致 — 低优，后续统一
9. useMyNumbers 绕过 service 层 — 低优，后续补 serviceApi
10. useBillDetail 返回 null — 后端无单笔接口，保持现状
11. useDepositHistory 返回空 — 后端无 history 端点，保持现状

### Suggestions (5) — 记录
- TxState auto-reset
- pendingSync 启动时 auto-retry（八-A 设计未完整实现）
- rand.Seed deprecated
- Operator.id 类型链（string→bigint 转换可优化）
- txHash 未传递给后端 recordDeposit/recordPayment

## 总结
3 Critical 0 剩余，代码可 ship。
