# Ship 阶段产出 — LinkWorld 全栈联调 (Round 2)

> 2026-04-15

## 发布信息

- **Commit**: 6fec0d7
- **Branch**: main
- **变更量**: 70 files changed, +4570 / -889

## 变更摘要

### 合约 (packages/contracts)
- deploy.ts 补全 Oracle 部署 + Deposit↔ServiceManager/Oracle↔Payment/Payment↔Oracle linking
- 合约地址输出到 deployments/localhost.json

### 后端 (packages/backend)
- CORS 中间件（允许 localhost:5173/3000）
- Operator model 加 CountryCode，Bill model 加 TxHash
- Register handler 接收 token_id，PayBill/Withdraw handler 补全实现
- 启动时自动 seed 11 个运营商

### 前端 (packages/web)
- 新增 axios client + 5 个 API service（userApi/operatorApi/depositApi/billingApi/usageApi）
- 新增 5 个合约 ABI + 合约地址配置（按 chainId 区分）
- 新增 4 个 wagmi v2 合约 hooks（UserRegistry/Deposit/ServiceManager/Payment）
- 新增 useTransactionFlow 状态机 + revert reason 中文映射
- 新增 pendingSync 双写补偿工具
- 4 个业务 hooks 从 mock → api + contract
- 5 个页面/组件适配新 hook API
- wagmi 配置 hardhatLocal 链（chainId 31337）
- services/index.ts 从 mock 切换到 api 导出

## Review 修复
- vite proxy target 3001→8080
- useRegister 自动链 backend sync（useEffect + useRef）
- Withdraw handler 移除无意义 DB 写入
- operatorApi 用 viem parseEther 替代浮点运算
