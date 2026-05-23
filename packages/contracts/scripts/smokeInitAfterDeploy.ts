// smokeInitAfterDeploy.ts — prove that initialize() can be called manually post-deploy
import { ethers } from "hardhat";

function encodeInit(abi: any, func = "initialize", args: any[] = []): string {
  return abi.encodeFunctionData(func, args);
}

async function main() {
  const [deployer] = await ethers.getSigners();

  // ── FeeManager ──
  {
    const F = await ethers.getContractFactory("FeeManager");
    const fm = await F.deploy({ args: [] });         // impl deployed (no-arg ctor)
    const addr = await fm.getAddress();
    const initData = encodeInit(F.interface, "initialize", [250]);
    const receipt = await (await deployer.sendTransaction({
      to: addr, data: initData
    })).wait();
    console.log("FeeManager initialized", addr, "feeRate:",
      (await (await ethers.getContractAt("FeeManager", addr)).getFeeRate()).toString());
  }

  // ── Deposit ──
  {
    const D = await ethers.getContractFactory("Deposit");
    const dep = await D.deploy({ args: [] });
    // Deploy a dummy minimal UserRegistry stub for demo purposes
    // (In real tests this will be the real UserRegistry address)
    console.log("Deposit impl deployed at", await dep.getAddress());
    const initData = encodeInit(D.interface, "initialize", [ethers.ZeroAddress]);
    await (await deployer.sendTransaction({ to: await dep.getAddress(), data: initData })).wait();
    console.log("  Initialized with ZeroAddress");
    console.log(
      "  trafficCardQuota:",
      (await dep.trafficCardQuota()).toString()
    );
  }
}

main().catch(console.error);
