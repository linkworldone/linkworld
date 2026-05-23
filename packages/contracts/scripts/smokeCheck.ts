import { ethers } from "hardhat";
import { HardhatUpgradesPlugin } from "@openzeppelin/hardhat-upgrades";

async function main() {
  try {
    // HardhatUpgradesPlugin lives at the MODULE level
    const defMod = await import("@openzeppelin/hardhat-upgrades");
    console.log("default export:", typeof defMod.default);
    console.log("named exports:", Object.keys(defMod));

    // Try if default exports the plugin factory
    if (defMod.default) {
      const hreModule = (defMod.default as any)();
      console.log("hreModule keys:", Object.keys(hreModule));
    }
  } catch (err: any) {
    console.error("import-plugin error:", err.message);
  }

  // Now try accessing hardhat runtime environment upgrades directly
  // We need to use it inside a task context
  const [signer] = await ethers.getSigners();

  // Common approach in test: type-cast using `ethers.run`
  // ethers.run("compile") is different; let's use transactions directly
  const F = await ethers.getContractFactory("FeeManager");
  // We need to call _initialize bytes data to proxy
  const iface = F.interface;
  const initData = iface.encodeFunctionData("initialize", [250]);
  console.log("initialize(250) calldata:", initData.substring(0, 40));
}

main().catch(console.error);
