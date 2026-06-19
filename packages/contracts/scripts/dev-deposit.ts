import { ethers } from "hardhat";
import * as fs from "fs";
import * as path from "path";

// 绕过钱包，直接在本地链上替指定钱包充值（impersonate）。
// 用法：WALLET=0x.. AMOUNT=50 npx hardhat run scripts/dev-deposit.ts --network localhost
// AMOUNT 必须是档位 10/20/50/100。
async function main() {
  const wallet = process.env.WALLET || "0x040fB0390FC18ac705e0a3aC9825bb49ac5BCCB2";
  const amount = process.env.AMOUNT || "50";
  const dep = JSON.parse(fs.readFileSync(path.join(__dirname, "../deployments/localhost.json"), "utf8"));

  const usdt = await ethers.getContractAt("MockUSDT", dep.usdt);
  const deposit = await ethers.getContractAt("Deposit", dep.proxies.Deposit);
  const signer = await ethers.getImpersonatedSigner(wallet);
  const amt = ethers.parseUnits(amount, 6);

  // approve（不足才补）
  if ((await usdt.allowance(wallet, dep.proxies.Deposit)) < amt) {
    await (await usdt.connect(signer).approve(dep.proxies.Deposit, amt)).wait();
    console.log(`approved ${amount} USDT`);
  }
  // deposit
  const tx = await deposit.connect(signer).deposit(amt);
  await tx.wait();
  console.log(`✅ deposit ${amount} USDT  tx=${tx.hash}`);
  console.log("   链上押金余额:", ethers.formatUnits(await deposit.getDepositAmount(wallet), 6), "USDT");
}

main().catch((e) => { console.error(e); process.exit(1); });
