import { ethers } from "hardhat";

async function main() {
  const [deployer] = await ethers.getSigners();
  console.log("Deploying with:", deployer.address);

  // 1. UserRegistry
  const UserRegistry = await ethers.getContractFactory("UserRegistry");
  const userRegistry = await UserRegistry.deploy();
  await userRegistry.waitForDeployment();
  const userRegistryAddr = await userRegistry.getAddress();
  console.log("UserRegistry:", userRegistryAddr);

  // 2. FeeManager (250 = 2.5%)
  const FeeManager = await ethers.getContractFactory("FeeManager");
  const feeManager = await FeeManager.deploy(250);
  await feeManager.waitForDeployment();
  const feeManagerAddr = await feeManager.getAddress();
  console.log("FeeManager:", feeManagerAddr);

  // 3. ServiceManager
  const ServiceManager = await ethers.getContractFactory("ServiceManager");
  const serviceManager = await ServiceManager.deploy();
  await serviceManager.waitForDeployment();
  const serviceManagerAddr = await serviceManager.getAddress();
  console.log("ServiceManager:", serviceManagerAddr);

  // 4. Payment
  const Payment = await ethers.getContractFactory("Payment");
  const payment = await Payment.deploy(feeManagerAddr, deployer.address);
  await payment.waitForDeployment();
  const paymentAddr = await payment.getAddress();
  console.log("Payment:", paymentAddr);

  // 5. Deposit
  const Deposit = await ethers.getContractFactory("Deposit");
  const deposit = await Deposit.deploy(userRegistryAddr);
  await deposit.waitForDeployment();
  const depositAddr = await deposit.getAddress();
  console.log("Deposit:", depositAddr);

  // 6. Oracle
  const Oracle = await ethers.getContractFactory("Oracle");
  const oracle = await Oracle.deploy(paymentAddr);
  await oracle.waitForDeployment();
  const oracleAddr = await oracle.getAddress();
  console.log("Oracle:", oracleAddr);

  // 关联合约
  await deposit.setPayment(paymentAddr);
  await deposit.setServiceManager(serviceManagerAddr);
  console.log("Deposit linked to Payment and ServiceManager");

  await payment.setOracle(oracleAddr);
  console.log("Payment linked to Oracle");

  console.log("\n=== Deployment Summary ===");
  console.log("UserRegistry:", userRegistryAddr);
  console.log("FeeManager:", feeManagerAddr);
  console.log("ServiceManager:", serviceManagerAddr);
  console.log("Payment:", paymentAddr);
  console.log("Deposit:", depositAddr);
  console.log("Oracle:", oracleAddr);
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
