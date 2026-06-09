#!/usr/bin/env bash
#
# gen-bindings.sh — 从 packages/contracts/artifacts 重新生成 8 份合约 abigen 绑定。
#
# 何时跑：合约 ABI 变更（selector / 新增函数 / 事件）后，重新生成并入库 bindings/，
# 让后端 go test 离线可跑。生成后务必跑 `go test ./internal/blockchain/...` 校验 abiHash。
#
# 依赖：node（从 hardhat artifacts 提取 .abi 字段）、go（调 scripts/genbindings）。
#
# 备注：不直接用 `abigen` 命令——go-ethereum v1.13.5 的 cmd/abigen 在 Go 1.21+ 下
# 因 github.com/fjl/memsize 的 runtime.stopTheWorld 符号变更无法链接。改用 bind.Bind
# 库（abigen 内部同款逻辑，产物一致）。详见 scripts/genbindings/main.go。
set -euo pipefail

BACKEND_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
REPO_ROOT="$(cd "$BACKEND_DIR/../.." && pwd)"
ARTIFACTS="$REPO_ROOT/packages/contracts/artifacts/contracts"
BINDINGS_DIR="$BACKEND_DIR/internal/blockchain/bindings"
TMP_ABI="$(mktemp -d)"
trap 'rm -rf "$TMP_ABI"' EXIT

# 合约名:artifact 相对路径（.sol/.json）。用 name:path 平铺列表以兼容 bash 3.2（macOS 自带，无关联数组）。
CONTRACTS="
FeeManager:FeeManager.sol/FeeManager.json
UserRegistry:UserRegistry.sol/UserRegistry.json
ServiceManager:ServiceManager.sol/ServiceManager.json
TrafficCardNFT:TrafficCardNFT.sol/TrafficCardNFT.json
Payment:Payment.sol/Payment.json
Deposit:Deposit.sol/Deposit.json
Oracle:Oracle.sol/Oracle.json
MockUSDT:mocks/MockUSDT.sol/MockUSDT.json
"

mkdir -p "$BINDINGS_DIR"

echo ">> 提取 ABI 并生成绑定 -> $BINDINGS_DIR"
for entry in $CONTRACTS; do
  name="${entry%%:*}"
  relpath="${entry#*:}"
  art="$ARTIFACTS/$relpath"
  abi_out="$TMP_ABI/${name}.abi.json"
  go_out="$BINDINGS_DIR/$(echo "$name" | tr '[:upper:]' '[:lower:]').go"

  # 提取 hardhat artifact 的 .abi 字段（紧凑 JSON 数组）。
  node -e "process.stdout.write(JSON.stringify(require('$art').abi))" > "$abi_out"

  ( cd "$BACKEND_DIR" && go run ./scripts/genbindings -abi "$abi_out" -pkg bindings -type "$name" -out "$go_out" )
done

echo ">> gofmt"
( cd "$BACKEND_DIR" && gofmt -w "$BINDINGS_DIR" )

echo ">> 完成。下一步：cd packages/backend && go build ./... && go test ./internal/blockchain/..."
