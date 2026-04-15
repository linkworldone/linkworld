import { ethers } from "hardhat";
import * as fs from "fs";
import * as path from "path";

async function main() {
  const [deployer] = await ethers.getSigners();
  console.log("Deploying with:", deployer.address);

  // 1. UserRegistry
  const UserRegistry = await ethers.getContractFactory("UserRegistry");
  const userRegistry = await UserRegistry.deploy();
  await userRegistry.waitForDeployment();
  console.log("UserRegistry:", await userRegistry.getAddress());

  // 2. FeeManager (250 = 2.5%)
  const FeeManager = await ethers.getContractFactory("FeeManager");
  const feeManager = await FeeManager.deploy(250);
  await feeManager.waitForDeployment();
  console.log("FeeManager:", await feeManager.getAddress());

  // 3. Deposit
  const Deposit = await ethers.getContractFactory("Deposit");
  const deposit = await Deposit.deploy(await userRegistry.getAddress());
  await deposit.waitForDeployment();
  console.log("Deposit:", await deposit.getAddress());

  // 4. ServiceManager
  const ServiceManager = await ethers.getContractFactory("ServiceManager");
  const serviceManager = await ServiceManager.deploy();
  await serviceManager.waitForDeployment();
  console.log("ServiceManager:", await serviceManager.getAddress());

  // 5. Payment
  const Payment = await ethers.getContractFactory("Payment");
  const payment = await Payment.deploy(
    await feeManager.getAddress(),
    deployer.address // platformWallet
  );
  await payment.waitForDeployment();
  console.log("Payment:", await payment.getAddress());

  // 6. Oracle
  const Oracle = await ethers.getContractFactory("Oracle");
  const oracle = await Oracle.deploy(await payment.getAddress());
  await oracle.waitForDeployment();
  console.log("Oracle:", await oracle.getAddress());

  // 关联合约
  await deposit.setPayment(await payment.getAddress());
  console.log("Deposit linked to Payment");

  await deposit.setServiceManager(await serviceManager.getAddress());
  console.log("Deposit linked to ServiceManager");

  await oracle.setPayment(await payment.getAddress());
  console.log("Oracle linked to Payment");

  await payment.setOracle(await oracle.getAddress());
  console.log("Payment linked to Oracle");

  // 输出合约地址 JSON
  const addresses = {
    chainId: 31337,
    UserRegistry: await userRegistry.getAddress(),
    FeeManager: await feeManager.getAddress(),
    Deposit: await deposit.getAddress(),
    ServiceManager: await serviceManager.getAddress(),
    Payment: await payment.getAddress(),
    Oracle: await oracle.getAddress(),
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
