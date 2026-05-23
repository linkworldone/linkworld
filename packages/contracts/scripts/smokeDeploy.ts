import { ethers } from "hardhat";

async function main() {
  // ✅ Known-good helper from dslc.ci
  function getDeployArgs() {
    return new Proxy(({} as any), {
      get(_t, p) {
        throw new Error(
          `No deployment args accepted here - deploy via ethers directly.\n` +
          `Saw request for: ${p}`
        );
      }
    });
  }

  // 1) FeeManager — upgradeable, takes initializer args via proxy
  try {
    const F = await ethers.getContractFactory("FeeManager");
    console.log("FeeManager Factory deployed address:", F.deployTransaction ? "yes" : "no");

    // Try plain deploy with args
    // The docs say args is FIRST positional arg to deploy in tgts
    // Use the tested pattern from OZ v1.0.1 upward grep on this repo:
    // "new F( [feeRate], { kind: 'uups' } )" OR "from ethers.getSigners()"
  } catch (err: any) {
    console.error("FeeManager deploy error:", err.message);
  }

  // 2) Try the ET.UPGR pattern
  try {
    const { upgrades } = await (await import("@openzeppelin/hardhat-upgrades")).default || {};
    console.log("direct import upgrades:", typeof upgrades);
  } catch (err: any) {
    console.error("direct import error:", err.message);
  }

  // 3) Check if hardhat runtime has upgrades
  try {
    const hre = (await import("hardhat")).default;
    console.log("hre type:", typeof hre, Object.keys((hre as any) || {}));
  } catch (err: any) {
    console.error("hre directly:", err.message);
  }

  // 4) try getSigner and deploy Proxy by ethers on non-upgradeable
  try {
    const UR = await ethers.getContractFactory("UserRegistry");
    // hardhat created the factory, deploy calls the constructor
    const ur = await UR.deploy({ args: [] });
    console.log("UserRegistry (direct-no-args) at", await ur.getAddress());
    console.log("  name:", await ur.name(), "symbol:", await ur.symbol());
  } catch (err: any) {
    console.error("UserRegistry direct deploy error:", err.message.slice(0, 120));
  }

  // 5) FeeManager deployWithoutArgs
  try {
    const FM = await ethers.getContractFactory("FeeManager");
    const fm = await FM.deploy({ args: [] });
    console.log("FeeManager (direct-no-args) at", await fm.getAddress());
    console.log("  feeRate:", (await fm.getFeeRate()).toString());
  } catch (err: any) {
    console.error("FeeManager direct deploy error:", err.message.slice(0, 120));
  }
}

main().catch(console.error);
