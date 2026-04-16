### Task 8: 业务 hooks 改造
- **状态**: DONE
- **产出**: 4 个文件改造（useUser, useOperator, useDeposit, useBilling）
- **变更**: mock service → api + contract 调用，双写流程含 pendingSync 补偿，useTransactionFlow 状态集成
- **验证**: hooks 文件 TypeScript 无错误，页面层有预期内的类型不匹配（待页面适配）
