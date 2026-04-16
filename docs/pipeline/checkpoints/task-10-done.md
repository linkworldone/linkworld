### Task 10: wagmi 配置 + 环境变量
- **状态**: DONE
- **产出**: chains.ts 修改、wagmi.ts 修改、.env.local 新建
- **变更**: 新增 hardhatLocal 链定义（31337），wagmi 根据 VITE_CHAIN_ID 切换链，环境变量配置
- **验证**: pnpm dev 启动无错误
