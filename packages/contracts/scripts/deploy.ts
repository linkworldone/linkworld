import { ethers, upgrades, network } from "hardhat";
import * as fs from "fs";
import * as path from "path";

async function main() {
  const [deployer] = await ethers.getSigners();
  console.log("Deploying with:", deployer.address);

  // 1. FeeManager (无依赖, 1.5% = 150 bps)
  const FeeManager = await ethers.getContractFactory("FeeManager");
  const feeManager = await upgrades.deployProxy(FeeManager, [150], {
    initializer: "initialize",
    kind: "uups"
  }) as any;
  await feeManager.waitForDeployment();
  const feeManagerAddr = await feeManager.getAddress();
  console.log("FeeManager:", feeManagerAddr);

  // 2. UserRegistry (无依赖)
  const UserRegistry = await ethers.getContractFactory("UserRegistry");
  const userRegistry = await upgrades.deployProxy(UserRegistry, [], {
    initializer: "initialize",
    kind: "uups"
  }) as any;
  await userRegistry.waitForDeployment();
  const userRegistryAddr = await userRegistry.getAddress();
  console.log("UserRegistry:", userRegistryAddr);

  // 3. ServiceManager (无依赖)
  const ServiceManager = await ethers.getContractFactory("ServiceManager");
  const serviceManager = await upgrades.deployProxy(ServiceManager, [], {
    initializer: "initialize",
    kind: "uups"
  }) as any;
  await serviceManager.waitForDeployment();
  const serviceManagerAddr = await serviceManager.getAddress();
  console.log("ServiceManager:", serviceManagerAddr);

  // 4. TrafficCardNFT (无依赖)
  const TrafficCardNFT = await ethers.getContractFactory("TrafficCardNFT");
  const trafficCardNFT = await upgrades.deployProxy(TrafficCardNFT, [], {
    initializer: "initialize",
    kind: "uups"
  }) as any;
  await trafficCardNFT.waitForDeployment();
  const trafficCardNFTAddr = await trafficCardNFT.getAddress();
  console.log("TrafficCardNFT:", trafficCardNFTAddr);

  // 5. Payment (依赖 FeeManager)
  const Payment = await ethers.getContractFactory("Payment");
  const payment = await upgrades.deployProxy(Payment, [feeManagerAddr, deployer.address], {
    initializer: "initialize",
    kind: "uups"
  }) as any;
  await payment.waitForDeployment();
  const paymentAddr = await payment.getAddress();
  console.log("Payment:", paymentAddr);

  // 6. Deposit (依赖 UserRegistry)
  const Deposit = await ethers.getContractFactory("Deposit");
  const deposit = await upgrades.deployProxy(Deposit, [userRegistryAddr], {
    initializer: "initialize",
    kind: "uups"
  }) as any;
  await deposit.waitForDeployment();
  const depositAddr = await deposit.getAddress();
  console.log("Deposit:", depositAddr);

  // 7. Oracle (无外部依赖)
  const Oracle = await ethers.getContractFactory("Oracle");
  const oracle = await upgrades.deployProxy(Oracle, [], {
    initializer: "initialize",
    kind: "uups"
  }) as any;
  await oracle.waitForDeployment();
  const oracleAddr = await oracle.getAddress();
  console.log("Oracle:", oracleAddr);

  // 关联合约 (v3 简化)
  await trafficCardNFT.setDepositContract(depositAddr);
  await deposit.setTrafficCardNFT(trafficCardNFTAddr);
  await deposit.setOracle(oracleAddr);
  await payment.setPlatformWallet(deployer.address);
  await oracle.setDeposit(depositAddr);
  console.log("Contracts linked");

  console.log("\n=== Deployment Summary ===");
  console.log("FeeManager:", feeManagerAddr);
  console.log("UserRegistry:", userRegistryAddr);
  console.log("ServiceManager:", serviceManagerAddr);
  console.log("TrafficCardNFT:", trafficCardNFTAddr);
  console.log("Payment:", paymentAddr);
  console.log("Deposit:", depositAddr);
  console.log("Oracle:", oracleAddr);

  const addresses = {
    chainId: network.config.chainId || 31337,
    proxies: {
      FeeManager: feeManagerAddr,
      UserRegistry: userRegistryAddr,
      ServiceManager: serviceManagerAddr,
      TrafficCardNFT: trafficCardNFTAddr,
      Payment: paymentAddr,
      Deposit: depositAddr,
      Oracle: oracleAddr,
    },
    implementations: {
      FeeManager: await upgrades.erc1967.getImplementationAddress(feeManagerAddr),
      UserRegistry: await upgrades.erc1967.getImplementationAddress(userRegistryAddr),
      ServiceManager: await upgrades.erc1967.getImplementationAddress(serviceManagerAddr),
      TrafficCardNFT: await upgrades.erc1967.getImplementationAddress(trafficCardNFTAddr),
      Payment: await upgrades.erc1967.getImplementationAddress(paymentAddr),
      Deposit: await upgrades.erc1967.getImplementationAddress(depositAddr),
      Oracle: await upgrades.erc1967.getImplementationAddress(oracleAddr),
    }
  };

  const deploymentsDir = path.resolve(__dirname, "../deployments");
  fs.mkdirSync(deploymentsDir, { recursive: true });
  const networkName = network.name;
  fs.writeFileSync(
    path.join(deploymentsDir, `${networkName}.json`),
    JSON.stringify(addresses, null, 2)
  );
  console.log(`Addresses written to deployments/${networkName}.json`);
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
