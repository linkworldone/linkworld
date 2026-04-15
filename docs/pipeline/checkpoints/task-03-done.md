### Task 3: 后端改造
- **状态**: DONE
- **产出**: 6 个文件修改（go.mod, main.go, models.go, handlers.go, repository.go, services.go）
- **变更**: CORS 中间件 + Operator 加 CountryCode + Bill 加 TxHash + Register 加 TokenID + PayBill 补全 + Withdraw 补全 + 11 个运营商种子数据
- **验证**: go build ./... 通过
