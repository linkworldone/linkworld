# LinkWorld

### Borderless connectivity, built around people who move across borders.

LinkWorld is an early-stage project exploring a simpler way to connect telecom operators with Web3 users, frequent travelers, and people who need reliable connectivity across countries.

Our long-term goal is a platform where users can discover regional mobile plans, connect with operators, and pay only for the connectivity they use. We believe communication should be easy to access and every payment should bring real value.

> **Project status:** LinkWorld is an early prototype. The demo currently runs on Arbitrum Sepolia and uses test assets. It does not provide real mobile connectivity or commercial eSIM/SIM service yet.

## Demo

[Open the LinkWorld demo](https://186.244.215.115/)

The demo uses the **Arbitrum Sepolia testnet**. Please do not use mainnet assets.

## Try the testnet demo

1. Install and unlock MetaMask, then connect a wallet.
2. Switch MetaMask to **Arbitrum Sepolia** (chain ID `421614`).
3. Make sure the wallet has enough **Arbitrum Sepolia ETH** to pay transaction fees. Sepolia ETH on Ethereum's Sepolia network is a different network; bridge it to Arbitrum Sepolia if needed. Test ETH has no intended real-world value.
4. The deposit flow uses the project's **MockUSDT** test token. To request MockUSDT, email [linkworldone@outlook.com](mailto:linkworldone@outlook.com) and include your **public wallet address**. Never send a seed phrase, private key, password, or recovery code.
5. Choose a deposit tier and review every MetaMask transaction carefully. The demo's intended test flow is to approve the required USDT amount, then deposit.

### If approval succeeds but deposit fails

A successful approval and a successful deposit are separate transactions. If the deposit transaction fails:

- Check that MetaMask is still on Arbitrum Sepolia and that the wallet is registered in the demo.
- Read the error shown in MetaMask before retrying.
- If MetaMask reports that the maximum fee is below the block's base fee, use MetaMask's suggested fee or refresh the estimate. As a troubleshooting example, a maximum fee **above 0.025 Gwei** may help when the current base fee is near that value. Network fees change; this is not a guaranteed or permanent setting.
- Keep the priority fee at or below the maximum fee, and review the token approval amount before confirming. It should match the amount you intend to deposit; do not approve an unexpectedly large amount.
- If the transaction still fails, send the team the public transaction hash and error message. Never share wallet secrets.

## Current test deposit model

For this testnet prototype, the intended deposit tiers are **10, 20, 50, and 100 MockUSDT**. The current product concept is **10 USDT for one unlimited-data eSIM/SIM usage card for one day**, with a choice of supported country or region. Availability in the demo is for product testing only; no real connectivity is delivered yet.

This is an initial experiment, not our final model. We are exploring a more useful platform that connects operators and users, and eventually offers Web3 users and frequent travelers flexible, pay-later connectivity and communication plans—with transparent pricing and less wasted spend.

## What we are building

- A wallet-connected interface for discovering and managing connectivity services
- A way for telecom operators to list and serve users across borders
- On-chain payment and digital service entitlement experiments
- A foundation for flexible regional data and communication plans

## Technology

- **Frontend:** React, TypeScript, Vite
- **Wallet and chain integration:** wagmi, RainbowKit, viem
- **Backend:** Go, Gin, GORM; SQLite and PostgreSQL support
- **Smart contracts:** Solidity, Hardhat, OpenZeppelin
- **Network:** Arbitrum Sepolia for the public test demo

## Development status

LinkWorld is at an early testing stage. The demo is useful for exploring the interface and testnet transaction flow, but it should not be treated as a production service. Contract addresses, test token availability, and network conditions may change as development continues.

## Contact

Questions, feedback, operator partnerships, and tester feedback are welcome:

- **Email:** [linkworldone@outlook.com](mailto:linkworldone@outlook.com)
- **X:** [@Mandyzix](https://x.com/Mandyzix)
- **Telegram:** [@linkworldteam](https://t.me/linkworldteam)

We're just getting started. Thank you for trusting us, trying the demo, and helping us improve. If you have an idea, find a problem, or want to work with us, please reach out.

---

*LinkWorld is an experimental project. Testnet tokens are not real money, and the demo does not guarantee mobile service, coverage, or commercial availability.*
