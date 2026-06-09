#!/usr/bin/env bash
# sync-deployments.sh — 把合约侧部署清单同步到 backend/configs/deployments.json。
#
# 合约 1/3 真·上链后，packages/contracts/deployments/<net>.json 会生成真实地址 + usdt + abiHash。
# 本脚本把指定网络的清单拷贝/合并进后端 config，一键回填占位零地址。
#
# 用法：
#   bash scripts/sync-deployments.sh arbitrum_sepolia   # 默认目标
#   bash scripts/sync-deployments.sh hardhat            # 本地联调
#
# 键名范式两端一致（proxies/usdt/usdtDecimals/abiHash），design §6.7。
set -euo pipefail

NET="${1:-arbitrum_sepolia}"

# 脚本位于 packages/backend/scripts/，定位到 backend 根与 contracts 根。
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKEND_DIR="$(dirname "$SCRIPT_DIR")"
CONTRACTS_DEPLOY_DIR="$(cd "$BACKEND_DIR/../contracts/deployments" && pwd)"

SRC="$CONTRACTS_DEPLOY_DIR/$NET.json"
DST="$BACKEND_DIR/configs/deployments.json"

if [[ ! -f "$SRC" ]]; then
  echo "ERROR: 源清单不存在：$SRC" >&2
  echo "       合约 1/3 真·上链后才会生成（见 handoff-backend.md §10）。" >&2
  exit 1
fi

# 用 jq 抽取 backend 需要的字段（chainId/rpcUrl/proxies/usdt/usdtDecimals/abiHash），
# 缺失的字段保留 DST 原值（rpcUrl 通常在 backend 侧配置，不一定在合约清单里）。
if ! command -v jq >/dev/null 2>&1; then
  echo "ERROR: 需要 jq。请先安装（brew install jq / apt-get install jq）。" >&2
  exit 1
fi

EXISTING_RPC="$(jq -r '.rpcUrl // ""' "$DST" 2>/dev/null || echo "")"

jq --arg rpc "$EXISTING_RPC" '{
  chainId: .chainId,
  rpcUrl: (.rpcUrl // $rpc),
  proxies: .proxies,
  usdt: .usdt,
  usdtDecimals: (.usdtDecimals // 6),
  abiHash: (.abiHash // {})
}' "$SRC" > "$DST.tmp"

mv "$DST.tmp" "$DST"
echo "OK: 已从 $SRC 同步到 $DST"
echo "    chainId=$(jq -r '.chainId' "$DST")  proxies=$(jq -r '.proxies | length' "$DST") 个"
