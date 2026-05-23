import { HardhatRuntimeEnvironment } from "hardhat/types";
import { ethers } from "hardhat";

export default async function (hre: HardhatRuntimeEnvironment) {
  console.log("hre keys:", Object.keys(hre));
  console.log("hre.ethers:", typeof hre.ethers);
  console.log("hre.upgrades:", typeof hre.upgrades);
  console.log("hre.run task available:", typeof hre.run);

  if (hre.upgrades) {
    console.log("  UPGRADES AVAILABLE");
    console.log("  deployProxy:", typeof hre.upgrades.deployProxy);
    console.log("  upgradeProxy:", typeof hre.upgrades.upgradeProxy);
  } else {
    console.log("  UPGRADES NOT AVAILABLE IN HRE");
  }
};
