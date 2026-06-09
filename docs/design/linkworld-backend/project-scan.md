# 后端子项目 project-scan 基线

> 扫描 2026-06-08 | 子项目 backend(2/3) | 已核对真实代码

## 1. 技术栈与模块树

- 语言：Go 1.25.0（`module linkworld-backend`，`packages/backend/go.mod`）
- Web 框架：Gin v1.12.0 + gin-contrib/cors v1.7.7
- ORM/DB：GORM v1.25.5 + `gorm.io/driver/postgres` v1.5.4 → **PostgreSQL**（非 SQLite）
- 链交互：go-ethereum v1.13.5（仅用 `ethclient` / `common` / `crypto`，**未用 abigen 生成绑定**）
- 测试：stretchr/testify v1.11.1（仅 `services_test.go`）

```
packages/backend/
├── cmd/main.go                       # 入口：DB 初始化 + AutoMigrate + 11 运营商 seed + 依赖装配 + 路由表 + 可选 event sync
├── go.mod / go.sum                   # Go 1.25.0
├── .env.example                      # DB + SERVER_PORT + RPC_URL=0g.ai + CHAIN_ID=16602
├── configs/deployments.json          # 链配置（chainId/rpcUrl/contracts）—— 见下「配置加载」陷阱
├── migrations/001_init.sql           # 唯一一份 SQL migration（实际运行靠 GORM AutoMigrate）
└── internal/
    ├── config/config.go              # LoadDeployments() + InitDB() + getEnv()
    ├── models/models.go              # 6 个 GORM 模型：User/Operator/Deposit/Bill/UserService/UsageData
    ├── repository/repository.go      # 6 个 Repository（CRUD 封装）
    ├── handlers/handlers.go          # Handler 聚合 + 20 个 HTTP 处理函数
    ├── services/
    │   ├── services.go               # UserService/OperatorService/BillingService/OracleService(v1)/NotificationService/DepositService/UserServiceService
    │   ├── oracle.go                  # VirtualNumberGenerator/OperatorAPISimulator/OracleServiceV2/UsageService
    │   └── services_test.go          # 单测
    ├── blockchain/
    │   ├── client.go                 # Client（链连接）—— 业务方法全是 stub，返回 0/false/not implemented
    │   ├── signatures.go             # 事件 topic keccak256 常量
    │   └── abis/{Deposit.json,UserRegistry.json}  # 仅 2 份手写裁剪 ABI（事件为主，非编译产物）
    └── sync/event_sync.go            # EventSync —— syncUserRegistered 是空轮询 stub（TODO）
```

## 2. 分层架构

标准四层，依赖单向：

```
HTTP (Gin 路由) → handlers.Handler → services.* → repository.* → GORM → PostgreSQL
                                          ↘ blockchain.Client / sync.EventSync（链侧，当前全 stub）
```

- handlers 只做参数绑定 + 调 service + 返回 JSON，无业务逻辑。
- services 持有 repository 指针，业务逻辑集中在此层。
- repository 是 GORM 的薄封装。
- blockchain / sync 与主链路**解耦**：仅当 `RPC_URL` 环境变量非空时才在 `main.go` 启动 event sync（go routine 后台跑），且当前为空实现。

## 3. Gin 路由全表（20 端点）

| # | Method | Path | Handler | 职责 |
|---|--------|------|---------|------|
| 1 | POST | `/api/register` | Register | 注册用户（wallet+email+token_id），写 DB |
| 2 | GET | `/api/user/:wallet` | GetUser | 按钱包查用户 |
| 3 | GET | `/api/operators` | GetOperators | 11 运营商列表 |
| 4 | POST | `/api/service/activate` | ActivateService | 开通运营商服务（虚拟号+密码） |
| 5 | POST | `/api/service/deactivate` | DeactivateService | 停用服务 |
| 6 | GET | `/api/service/:wallet` | GetUserService | 查当前激活服务 |
| 7 | GET | `/api/bills/:wallet` | GetBills | 查用户账单 |
| 8 | POST | `/api/bills/pay` | PayBill | 标记账单已付（传 tx_hash，**仅改 DB 状态**） |
| 9 | POST | `/api/deposit` | Deposit | 记录押金（**仅写 DB，不上链**） |
| 10 | GET | `/api/deposit/:wallet` | GetDeposit | 查押金总额（DB 聚合） |
| 11 | GET | `/api/deposit/:wallet/history` | GetDepositHistory | 押金流水 |
| 12 | POST | `/api/withdraw` | Withdraw | 记录提现（链上已完成，后端记账） |
| 13 | POST | `/api/virtual-number/generate` | GenerateVirtualNumber | 按国家码生成虚拟号+密码 |
| 14 | GET | `/api/countries` | GetCountryList | 国家列表 |
| 15 | GET | `/api/usage/:wallet` | GetUsage | 查用量（OperatorAPISimulator 模拟 + 签名 + 落库） |
| 16 | POST | `/api/oracle/monthly-bill` | TriggerMonthlyBill | 触发月度出账（遍历用户，模拟计价后写 Bill） |
| 17 | POST | `/api/usage/submit` | SubmitUsage | 提交用量数据（带签名） |
| 18 | POST | `/api/notification/send` | SendNotification | 发账单邮件（**TODO stub**，仅 log） |

> 实际 `r.POST/GET` 注册行 = 18 条（main.go L103-120）。任务书所称「20 端点」含 CORS OPTIONS 预检 + `r.Run` 监听，按真实业务端点计为 **18 个**。

## 4. 配置加载（含一处真实 bug）

`config.LoadDeployments` 解析的 struct：

```go
type Deployments struct {
    ChainID uint64                    `json:"chainId"`
    RpcURL  string                    `json:"rpcUrl"`
    Proxies map[string]common.Address `json:"proxies"`   // ← 读 "proxies"
}
```

但 `configs/deployments.json` 里键名是 **`"contracts"`**（不是 `"proxies"`）：

```json
{ "chainId": 16602, "rpcUrl": "https://evm-testnet.0g.ai", "contracts": { ... } }
```

→ **`Proxies` 永远反序列化为空 map**。即使配了 `RPC_URL`，传给 `NewClient` / `NewEventSync` 的合约地址表也是空的。这是 scan 阶段需在后续对齐时一并修掉的现存缺陷（键名 + 字段语义双重不一致）。

## 5. 启动流程要点（main.go）

1. `config.InitDB()` 连 PostgreSQL（env 可覆盖，默认 127.0.0.1:5432 / db=linkworld）。
2. `AutoMigrate` 6 张表（GORM 自动建表，`migrations/001_init.sql` 形同备份不参与运行）。
3. 首次启动 seed 11 个内置运营商（T-Mobile/Vodafone/.../Globe），`RequiredDeposit` 是 `"0.01"` 这类**小数字符串**（当前按「native/通用单位」语义，非 USDT 6 位精度）。
4. 装配 repository → service → handler。
5. 可选：`RPC_URL` 非空才启动 `EventSync.Start`（后台 goroutine）。
6. CORS 放行 5173/3000 前端源。
7. `r.Run(":8080")`。
