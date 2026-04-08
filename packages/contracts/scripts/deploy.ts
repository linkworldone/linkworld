import { ethers } from "hardhat";

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

  // 关联 Deposit <-> Payment
  await deposit.setPayment(await payment.getAddress());
  console.log("Deposit linked to Payment");
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
