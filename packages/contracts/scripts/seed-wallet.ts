import { ethers } from "hardhat";
import * as fs from "fs";
import * as path from "path";

// 给指定钱包充值（ETH + USDT）并在链上注册。
// 钱包地址：环境变量 WALLET 覆盖，默认本地常用测试钱包。
// 合约地址从 deployments/localhost.json 读取（不硬编码，重新部署也不用改）。
async function main() {
  const wallet = process.env.WALLET || "0x040fB0390FC18ac705e0a3aC9825bb49ac5BCCB2";
  const dep = JSON.parse(
    fs.readFileSync(path.join(__dirname, "../deployments/localhost.json"), "utf8"),
  );
  const usdtAddr: string = dep.usdt;
  const urAddr: string = dep.proxies.UserRegistry;

  // 1) 原生币 ETH（付 gas）—— 100 ETH
  await ethers.provider.send("hardhat_setBalance", [
    wallet,
    "0x" + (100n * 10n ** 18n).toString(16),
  ]);

  // 2) USDT —— mint 500
  const usdt = await ethers.getContractAt("MockUSDT", usdtAddr);
  await (await usdt.mint(wallet, ethers.parseUnits("500", 6))).wait();

  // 3) 链上注册（未注册才注册；deposit 合约要求 isRegistered）
  const ur = await ethers.getContractAt("UserRegistry", urAddr);
  if (!(await ur.isRegistered(wallet))) {
    const signer = await ethers.getImpersonatedSigner(wallet);
    await (await ur.connect(signer).register("user@linkworld.io")).wait();
  }

  const ethBal = await ethers.provider.getBalance(wallet);
  const usdtBal = await usdt.balanceOf(wallet);
  console.log(`✅ ${wallet}`);
  console.log(`   ETH=${ethers.formatUnits(ethBal, 18)}  USDT=${ethers.formatUnits(usdtBal, 6)}  registered=${await ur.isRegistered(wallet)}`);
}

main().catch((e) => { console.error(e); process.exit(1); });
