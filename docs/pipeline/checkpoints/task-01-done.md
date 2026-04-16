### Task 1: 前端合约 ABI + 地址配置
- **状态**: DONE
- **产出**: 7 个新文件（5 个 ABI + abis/index.ts + contracts.ts）
- **变更**: 从 artifacts 提取 5 个合约 ABI，配置 localhost 合约地址（从 deployments/localhost.json 读取），含 getContractAddress 工具函数
- **验证**: 文件创建成功，地址与部署输出一致
