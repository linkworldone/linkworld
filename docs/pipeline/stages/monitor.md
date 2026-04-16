# 复盘 — LinkWorld 全栈联调 (Round 2)

> 2026-04-15

## Round 2 总结

### 目标达成
将前端从 Mock Service 切换到真实后端 API + 链上合约交互，三端本地环境串通。

### 数据
- **阶段**: 11 个（scan → monitor 全量）
- **Task**: 11 个，全部完成
- **代码量**: 70 files, +4570 / -889
- **耗时**: 跨 2 天（中间撞 subagent 限额暂停 1 天）
- **Review 修复**: 3 Critical + 2 Important

### 做得好
1. **三端数据模型对齐彻底** — API service 层的 snake_case→camelCase 转换覆盖了所有字段
2. **混合调用策略清晰** — 资金操作走合约，查询走后端，注册/激活双写
3. **双写补偿设计** — pendingSync + retryWithBackoff + 链上状态兜底，arch-review 阶段就识别了风险
4. **页面零改动原则** — 只换底层数据源，8 个页面 UI 完全不动，最后只改了 5 个页面的 hook 调用方式
5. **合约部署脚本完善** — Oracle 部署 + 全部 linking + 地址输出 JSON，一键可复现

### 做得不好
1. **subagent 限额中断** — 批次 3 三个 subagent 同时撞限额，导致暂停 1 天。应该控制并发数或错开时间
2. **页面适配遗漏** — implement 阶段没有同步更新页面组件，导致 test 阶段发现 16 个 TS 错误。应该在 Task 8（hooks 改造）时就同步改页面
3. **useDepositHistory 和 useBillDetail 返空** — 后端缺失端点，前端 hook 只能返回空数组/null，留下了功能缺口
4. **pendingSync 启动时 auto-retry 未实现** — 设计文档八-A 写了启动检测，实现时跳过了。技术债

### Round 3 建议

| 优先级 | 项目 | 理由 |
|--------|------|------|
| P0 | 后端 wallet 签名认证 | 当前 API 裸奔，部署测试网前必须加 |
| P0 | 合约移除明文密码存储 | ServiceManager 存明文 password，上公链前必须改 |
| P1 | 事件索引器替代前端双写补偿 | pendingSync 是临时方案，需要后端监听合约 event |
| P1 | 后端补 deposit history + bill detail 端点 | 前端两个 hook 返回空 |
| P2 | pendingSync 启动时 auto-retry | 完成设计文档八-A 的完整实现 |
| P2 | 通知 API（后端 + 前端去 mock） | 当前通知仍 mock |
| P3 | ABI 自动生成（wagmi CLI 或 deploy 脚本输出） | 手动提取易过期 |

### 流程改进
1. **implement 阶段对 hook API 变更做 breaking change 检查** — hooks 返回值变了但没同步改调用方，应加个检查步骤
2. **subagent 并发控制** — 最多同时 2 个，避免撞限额
3. **checkpoint 写入时机** — Round 1 的教训已吸取，Round 2 每个 Task 完成后立即写 checkpoint
