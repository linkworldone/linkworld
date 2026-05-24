import { ethers, upgrades } from "hardhat";
import * as fs from "fs";
import * as path from "path";

async function main() {
  const [deployer] = await ethers.getSigners();
  console.log("Deploying with:", deployer.address);

  // 1. UserRegistry (无依赖)
  const UserRegistry = await ethers.getContractFactory("UserRegistry");
  const userRegistry = await upgrades.deployProxy(UserRegistry, [], {
    initializer: "initialize",
    kind: "uups"
  }) as any;
  await userRegistry.waitForDeployment();
  const userRegistryAddr = await userRegistry.getAddress();
  console.log("UserRegistry:", userRegistryAddr);

  // 2. FeeManager (无依赖，initializer: initialize(uint256))
  const FeeManager = await ethers.getContractFactory("FeeManager");
  const feeManager = await upgrades.deployProxy(FeeManager, [250], {
    initializer: "initialize",
    kind: "uups"
  }) as any;
  await feeManager.waitForDeployment();
  const feeManagerAddr = await feeManager.getAddress();
  console.log("FeeManager:", feeManagerAddr);

  // 3. ServiceManager (无依赖)
  const ServiceManager = await ethers.getContractFactory("ServiceManager");
  const serviceManager = await upgrades.deployProxy(ServiceManager, [], {
    initializer: "initialize",
    kind: "uups"
  }) as any;
  await serviceManager.waitForDeployment();
  const serviceManagerAddr = await serviceManager.getAddress();
  console.log("ServiceManager:", serviceManagerAddr);

  // 4. Payment (依赖 FeeManager)
  const Payment = await ethers.getContractFactory("Payment");
  const payment = await upgrades.deployProxy(Payment, [feeManagerAddr, deployer.address], {
    initializer: "initialize",
    kind: "uups"
  }) as any;
  await payment.waitForDeployment();
  const paymentAddr = await payment.getAddress();
  console.log("Payment:", paymentAddr);

  // 5. Deposit (依赖 UserRegistry)
  const Deposit = await ethers.getContractFactory("Deposit");
  const deposit = await upgrades.deployProxy(Deposit, [userRegistryAddr], {
    initializer: "initialize",
    kind: "uups"
  }) as any;
  await deposit.waitForDeployment();
  const depositAddr = await deposit.getAddress();
  console.log("Deposit:", depositAddr);

  // 6. Oracle (依赖 Payment)
  const Oracle = await ethers.getContractFactory("Oracle");
  const oracle = await upgrades.deployProxy(Oracle, [paymentAddr], {
    initializer: "initialize",
    kind: "uups"
  }) as any;
  await oracle.waitForDeployment();
  const oracleAddr = await oracle.getAddress();
  console.log("Oracle:", oracleAddr);

  // 7. TrafficCardNFT (无依赖)
  const TrafficCardNFT = await ethers.getContractFactory("TrafficCardNFT");
  const trafficCardNFT = await upgrades.deployProxy(TrafficCardNFT, [], {
    initializer: "initialize",
    kind: "uups"
  }) as any;
  await trafficCardNFT.waitForDeployment();
  const trafficCardNFTAddr = await trafficCardNFT.getAddress();
  console.log("TrafficCardNFT:", trafficCardNFTAddr);

  // 关联合约
  await deposit.setPayment(paymentAddr);
  await deposit.setServiceManager(serviceManagerAddr);
  await deposit.setOracle(oracleAddr);
  await deposit.setTrafficCardNFT(trafficCardNFTAddr);
  console.log("Deposit linked to Payment, ServiceManager, Oracle, TrafficCardNFT");

  await payment.setOracle(oracleAddr);
  await payment.setDeposit(depositAddr);
  console.log("Payment linked to Oracle and Deposit");

  await oracle.setDeposit(depositAddr);
  console.log("Oracle linked to Deposit");

  // 将 TrafficCardNFT 的 ownership 转给 Deposit 合约，让其可调用 mint
  await trafficCardNFT.transferOwnership(depositAddr);
  console.log("TrafficCardNFT ownership transferred to Deposit");

  console.log("\n=== Deployment Summary ===");
  console.log("UserRegistry:", userRegistryAddr);
  console.log("FeeManager:", feeManagerAddr);
  console.log("ServiceManager:", serviceManagerAddr);
  console.log("Payment:", paymentAddr);
  console.log("Deposit:", depositAddr);
  console.log("Oracle:", oracleAddr);
  console.log("TrafficCardNFT:", trafficCardNFTAddr);

  // 输出合约地址 JSON (同时保存实现地址用于后续升级)
  const addresses = {
    chainId: 31337,
    proxies: {
      UserRegistry: userRegistryAddr,
      FeeManager: feeManagerAddr,
      ServiceManager: serviceManagerAddr,
      Payment: paymentAddr,
      Deposit: depositAddr,
      Oracle: oracleAddr,
      TrafficCardNFT: trafficCardNFTAddr,
    },
    implementations: {
      UserRegistry: await upgrades.erc1967.getImplementationAddress(userRegistryAddr),
      FeeManager: await upgrades.erc1967.getImplementationAddress(feeManagerAddr),
      ServiceManager: await upgrades.erc1967.getImplementationAddress(serviceManagerAddr),
      Payment: await upgrades.erc1967.getImplementationAddress(paymentAddr),
      Deposit: await upgrades.erc1967.getImplementationAddress(depositAddr),
      Oracle: await upgrades.erc1967.getImplementationAddress(oracleAddr),
      TrafficCardNFT: await upgrades.erc1967.getImplementationAddress(trafficCardNFTAddr),
    }
  };

  const deploymentsDir = path.resolve(__dirname, "../deployments");
  fs.mkdirSync(deploymentsDir, { recursive: true });
  fs.writeFileSync(
    path.join(deploymentsDir, "localhost.json"),
    JSON.stringify(addresses, null, 2)
  );
  console.log("Addresses written to deployments/localhost.json");
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});