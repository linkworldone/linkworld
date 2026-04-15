### Task 0: 合约部署脚本补全
- **状态**: DONE
- **产出**: packages/contracts/scripts/deploy.ts, packages/contracts/deployments/localhost.json
- **变更**: 新增 Oracle 部署 + 4 个 linking（Deposit↔ServiceManager, Oracle↔Payment, Payment↔Oracle）+ 地址输出 JSON
- **验证**: 部署成功，6 个合约地址输出到 deployments/localhost.json
