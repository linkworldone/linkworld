import { defineChain } from "viem";

// Arbitrum Sepolia 测试网（本轮主网络，替代旧 0G testnet）
// 真·上链待合约部署后，地址回填 deployments/arbitrum_sepolia.json
export const arbitrumSepolia = defineChain({
  id: 421614,
  name: "Arbitrum Sepolia",
  nativeCurrency: { name: "Ether", symbol: "ETH", decimals: 18 },
  rpcUrls: {
    default: { http: ["https://sepolia-rollup.arbitrum.io/rpc"] },
  },
  blockExplorers: {
    default: { name: "Arbiscan", url: "https://sepolia.arbiscan.io" },
  },
  testnet: true,
});

// 本地链 RPC 可通过 VITE_RPC_URL 覆盖（公网部署时指向 Caddy 反代的 /rpc），缺省回退到本机 hardhat
const LOCAL_RPC = import.meta.env.VITE_RPC_URL || "http://127.0.0.1:8545";

export const hardhatLocal = defineChain({
  id: 31337,
  name: "Hardhat Local",
  nativeCurrency: { name: "Ether", symbol: "ETH", decimals: 18 },
  rpcUrls: {
    default: { http: [LOCAL_RPC] },
  },
});
