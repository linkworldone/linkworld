#!/usr/bin/env bash
# LinkWorld 本地全栈一键启动 + 钱包充值/注册。
#
# 在你自己的独立终端跑（不要在别的进程托管环境里跑，否则进程可能被回收）：
#   bash tasks/dev-up.sh
#   WALLET=0x你的地址 bash tasks/dev-up.sh   # 给别的钱包充值
#
# 做的事：①清旧进程 ②起 hardhat 链 ③部署+同步后端配置 ④起 Go 后端 ⑤起前端 ⑥给钱包充 100ETH+500USDT+链上注册
# 注意：每次跑都是全新链（fresh deploy，地址确定不变）。链一重启 MetaMask 要清一次活动数据/重连网络。
set -uo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"
LOG=/tmp/linkworld; mkdir -p "$LOG"
WALLET="${WALLET:-0x040fB0390FC18ac705e0a3aC9825bb49ac5BCCB2}"

wait_port() { # $1=端口 $2=超时秒(默认40)
  for _ in $(seq 1 "${2:-40}"); do
    lsof -nP -iTCP:"$1" -sTCP:LISTEN >/dev/null 2>&1 && return 0
    sleep 1
  done
  echo "  ⚠ 端口 $1 等待超时，看 $LOG 日志" >&2; return 1
}

echo "▶ 0/6 清理旧进程…"
pkill -f "hardhat node" 2>/dev/null || true
pkill -f "go run ./cmd/main.go" 2>/dev/null || true
pkill -f "exe/main" 2>/dev/null || true
pkill -f "@linkworld/web" 2>/dev/null || true
sleep 1

echo "▶ 1/6 hardhat 链 (8545)…"
nohup pnpm -F @linkworld/contracts exec hardhat node > "$LOG/hardhat.log" 2>&1 & disown
wait_port 8545 30 || exit 1

echo "▶ 2/6 部署合约…"
pnpm -F @linkworld/contracts deploy:local > "$LOG/deploy.log" 2>&1 || { echo "  部署失败，看 $LOG/deploy.log" >&2; exit 1; }

echo "▶ 3/6 同步部署清单到后端配置…"
bash packages/backend/scripts/sync-deployments.sh localhost > /dev/null
# 新链=块从 0：重置后端同步游标 + 清旧链的事件/业务记录。
# 否则 event_sync 把旧游标(几百块)当 reorg 卡死，新 deposit/兑换永远 pending。
psql -d linkworld -c "UPDATE sync_states SET last_block=0; DELETE FROM chain_events; DELETE FROM deposits; DELETE FROM sims;" >/dev/null 2>&1 || true

echo "▶ 4/6 Go 后端 (8080)…"
( cd packages/backend && set -a && . ./.env && set +a && nohup go run ./cmd/main.go > "$LOG/backend.log" 2>&1 & disown )
wait_port 8080 60 || exit 1

echo "▶ 5/6 前端 (5173)…"
nohup pnpm dev:web > "$LOG/web.log" 2>&1 & disown
wait_port 5173 30 || exit 1

echo "▶ 6/6 充值 + 链上注册: $WALLET"
WALLET="$WALLET" pnpm -F @linkworld/contracts exec hardhat run scripts/seed-wallet.ts --network localhost

cat <<EOF

✅ 全部就绪
   链   http://localhost:8545  (chainId 31337)
   后端  http://localhost:8080
   前端  http://localhost:5173
   日志  $LOG/{hardhat,deploy,backend,web}.log
   钱包  $WALLET  → 100 ETH + 500 USDT + 已注册

下一步：MetaMask 清一次「活动和 nonce 数据」（链重启了），再用前端。
停止全部：pkill -f 'hardhat node'; pkill -f 'go run'; pkill -f '@linkworld/web'
EOF
