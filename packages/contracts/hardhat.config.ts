import { HardhatUserConfig } from "hardhat/config";
import "@nomicfoundation/hardhat-toolbox";
import "@openzeppelin/hardhat-upgrades";
import * as dotenv from "dotenv";

dotenv.config();

const DEPLOYER_KEY = process.env.DEPLOYER_PRIVATE_KEY || "0x" + "0".repeat(64);

const config: HardhatUserConfig = {
  solidity: {
    version: "0.8.27",
    settings: {
      optimizer: { enabled: true, runs: 200 },
      evmVersion: "cancun",
      viaIR: true,
    },
  },
  networks: {
    hardhat: {
      chainId: 31337,
      // 持续出块：交易即时打包(auto) + 每 4 秒出一个空块(interval)，让后端 event_sync 的 K=5 确认块能凑齐、
      // deposit/兑换等资金事件从 pending 自动转 confirmed。
      // 注:本地用 MetaMask 时记得开 VPN——它的交易模拟/安全/gas 等远程服务在国外,不开 VPN 会超时转圈。
      mining: {
        auto: true,
        interval: 4000,
      },
    },
    og_mainnet: {
      url: "https://evmrpc.0g.ai",
      chainId: 16600,
      accounts: [DEPLOYER_KEY],
      timeout: 60000,
    },
    og_testnet: {
      url: "https://evmrpc-testnet.0g.ai",
      chainId: 16602,
      accounts: [DEPLOYER_KEY],
      timeout: 60000,
    },
    sepolia: {
      url: "https://ethereum-sepolia-rpc.publicnode.com",
      chainId: 11155111,
      accounts: [DEPLOYER_KEY],
    },
    arbitrum_sepolia: {
      url: process.env.ARBITRUM_SEPOLIA_RPC || "https://sepolia-rollup.arbitrum.io/rpc",
      chainId: 421614,
      accounts: process.env.DEPLOYER_PRIVATE_KEY ? [process.env.DEPLOYER_PRIVATE_KEY] : [],
      timeout: 60000,
    },
  },
  paths: {
    sources: "./contracts",
    tests: "./test",
    cache: "./cache",
    artifacts: "./artifacts",
  },
};

export default config;