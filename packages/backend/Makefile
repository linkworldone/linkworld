# LinkWorld backend — 开发常用任务
#
# 绑定与 ABI 校验由 T1 引入；其余 target 后续阶段按需补充。

.PHONY: gen-bindings verify-abi build test

# 从 packages/contracts/artifacts 重新生成 8 份合约 abigen 绑定 + 刷新嵌入的 ABI JSON。
# 合约 ABI（selector / 函数 / 事件）变更后跑此命令，再跑 verify-abi 确认指纹同步。
gen-bindings:
	bash scripts/gen-bindings.sh

# 校验后端绑定的 ABI 指纹与 deployments/<net>.json 的 abiHash 一致（formatJson+keccak256）。
verify-abi:
	go test ./internal/blockchain/ -run TestABIHash -v

build:
	go build ./...

test:
	go test ./...
