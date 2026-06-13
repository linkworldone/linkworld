import { ethers, ContractTransaction, upgrades, network } from "hardhat";
import * as fs from "fs";
import * as path from "path";

async function sendAndWait(txPromise: Promise<ContractTransaction>, description: string) {
  const tx = await txPromise;
  console.log(`${description} tx hash: ${tx.hash}`);
  await tx.wait();
}

function multiplyFeeValue(value: unknown | undefined, factor = 2n) {
  if (value === undefined || value === null) return undefined;
  if (typeof value === "bigint") return value * factor;
  if (typeof value === "number") return BigInt(value) * factor;
  if (typeof value === "string") return BigInt(value) * factor;
  if (typeof value === "object" && value !== null && "toBigInt" in value) {
    return (value as any).toBigInt() * factor;
  }
  return undefined;
}

async function getTxOverrides() {
  const feeData = await ethers.provider.getFeeData();
  const maxFeePerGas = multiplyFeeValue(feeData.maxFeePerGas);
  const maxPriorityFeePerGas = multiplyFeeValue(feeData.maxPriorityFeePerGas);
  if (maxFeePerGas && maxPriorityFeePerGas) {
    return { maxFeePerGas, maxPriorityFeePerGas };
  }
  const gasPriceValue = feeData.gasPrice ?? (await ethers.provider.getGasPrice());
  const gasPrice = multiplyFeeValue(gasPriceValue);
  return gasPrice ? { gasPrice } : {};
}

const BUILTIN_OPERATOR_IDS = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11];

async function main() {
  const [deployer] = await ethers.getSigners();
  console.log("Deploying with:", deployer.address);

  // 0. MockUSDT (used on testnets; replace with real USDT for mainnet)
  const MockUSDT = await ethers.getContractFactory("MockUSDT");
  const usdt = (await MockUSDT.deploy()) as any;
  await usdt.waitForDeployment();
  const usdtAddr = await usdt.getAddress();
  const usdtDecimals = Number(await usdt.decimals());
  console.log("MockUSDT:", usdtAddr, "(decimals=" + usdtDecimals + ")");

  // 1. FeeManager (1.5% = 150 bps)
  const FeeManager = await ethers.getContractFactory("FeeManager");
  const feeManager = await upgrades.deployProxy(FeeManager, [150], {
    initializer: "initialize", kind: "uups"
  }) as any;
  await feeManager.waitForDeployment();
  const feeManagerAddr = await feeManager.getAddress();
  console.log("FeeManager:", feeManagerAddr);

  // 2. UserRegistry
  const UserRegistry = await ethers.getContractFactory("UserRegistry");
  const userRegistry = await upgrades.deployProxy(UserRegistry, [], {
    initializer: "initialize", kind: "uups"
  }) as any;
  await userRegistry.waitForDeployment();
  const userRegistryAddr = await userRegistry.getAddress();
  console.log("UserRegistry:", userRegistryAddr);

  // 3. ServiceManager
  const ServiceManager = await ethers.getContractFactory("ServiceManager");
  const serviceManager = await upgrades.deployProxy(ServiceManager, [], {
    initializer: "initialize", kind: "uups"
  }) as any;
  await serviceManager.waitForDeployment();
  const serviceManagerAddr = await serviceManager.getAddress();
  console.log("ServiceManager:", serviceManagerAddr);

  // 4. TrafficCardNFT
  const TrafficCardNFT = await ethers.getContractFactory("TrafficCardNFT");
  const trafficCardNFT = await upgrades.deployProxy(TrafficCardNFT, [], {
    initializer: "initialize", kind: "uups"
  }) as any;
  await trafficCardNFT.waitForDeployment();
  const trafficCardNFTAddr = await trafficCardNFT.getAddress();
  console.log("TrafficCardNFT:", trafficCardNFTAddr);

  // 5. Payment (depends on FeeManager, MockUSDT, ServiceManager)
  const Payment = await ethers.getContractFactory("Payment");
  const payment = await upgrades.deployProxy(
    Payment,
    [feeManagerAddr, deployer.address, usdtAddr, serviceManagerAddr],
    { initializer: "initialize", kind: "uups" }
  ) as any;
  await payment.waitForDeployment();
  const paymentAddr = await payment.getAddress();
  console.log("Payment:", paymentAddr);

  // 6. Deposit (depends on UserRegistry, MockUSDT)
  const Deposit = await ethers.getContractFactory("Deposit");
  const deposit = await upgrades.deployProxy(
    Deposit,
    [userRegistryAddr, usdtAddr],
    { initializer: "initialize", kind: "uups" }
  ) as any;
  await deposit.waitForDeployment();
  const depositAddr = await deposit.getAddress();
  console.log("Deposit:", depositAddr);

  // 7. Oracle (no external deps; payment/deposit wired after via setters)
  const Oracle = await ethers.getContractFactory("Oracle");
  const oracle = await upgrades.deployProxy(Oracle, [], {
    initializer: "initialize", kind: "uups"
  }) as any;
  await oracle.waitForDeployment();
  const oracleAddr = await oracle.getAddress();
  console.log("Oracle:", oracleAddr);

  // ===== Wiring =====
  const txOverrides = await getTxOverrides();
  await sendAndWait(trafficCardNFT.setDepositContract(depositAddr), "trafficCardNFT.setDepositContract");
  await sendAndWait(deposit.setTrafficCardNFT(trafficCardNFTAddr), "deposit.setTrafficCardNFT");
  await sendAndWait(deposit.setOracle(oracleAddr), "deposit.setOracle");
  await sendAndWait(payment.setPlatformWallet(deployer.address), "payment.setPlatformWallet");
  await sendAndWait(oracle.setDeposit(depositAddr), "oracle.setDeposit");
  await sendAndWait(payment.setOracle(oracleAddr), "payment.setOracle");
  await sendAndWait(oracle.setPayment(paymentAddr), "oracle.setPayment");
  console.log("Contracts linked");

  // Set operator payment addresses for built-in operators
  for (const id of BUILTIN_OPERATOR_IDS) {
    const payoutAddr = ethers.getAddress(
      "0x" + ethers.keccak256(
        ethers.solidityPacked(["address", "uint256"], [deployer.address, id])
      ).slice(-40)
    );
    await (await serviceManager.setOperatorPaymentAddress(id, payoutAddr)).wait();
  }
  console.log(`OperatorPaymentAddress set for ${BUILTIN_OPERATOR_IDS.length} built-in operators`);

  // Transfer TrafficCardNFT ownership to Deposit
  await sendAndWait(trafficCardNFT.transferOwnership(depositAddr, txOverrides), "trafficCardNFT.transferOwnership");
  console.log("TrafficCardNFT ownership transferred to Deposit");

  // ===== Post-deploy assertions =====
  const assertEq = (actual: string, expected: string, label: string) => {
    if (actual.toLowerCase() !== expected.toLowerCase()) {
      throw new Error(`Wiring assertion failed: ${label} = ${actual}, expected ${expected}`);
    }
  };
  assertEq(await payment.oracle(), oracleAddr, "payment.oracle()");
  assertEq(await oracle.payment(), paymentAddr, "oracle.payment()");
  assertEq(await oracle.deposit(), depositAddr, "oracle.deposit()");
  assertEq(await deposit.oracle(), oracleAddr, "deposit.oracle()");
  assertEq(await deposit.trafficCardNFT(), trafficCardNFTAddr, "deposit.trafficCardNFT()");
  assertEq(await trafficCardNFT.owner(), depositAddr, "trafficCardNFT.owner()");
  assertEq(await payment.usdt(), usdtAddr, "payment.usdt()");
  assertEq(await deposit.usdt(), usdtAddr, "deposit.usdt()");
  assertEq(await payment.serviceManager(), serviceManagerAddr, "payment.serviceManager()");
  for (const id of BUILTIN_OPERATOR_IDS) {
    const op = await serviceManager.getOperator(id);
    if (op.paymentAddress === ethers.ZeroAddress) {
      throw new Error(`operator ${id} paymentAddress is zero`);
    }
  }
  console.log("All wiring assertions passed");

  console.log("\n=== Deployment Summary ===");
  console.log("MockUSDT:", usdtAddr, "(decimals=" + usdtDecimals + ")");
  console.log("FeeManager:", feeManagerAddr);
  console.log("UserRegistry:", userRegistryAddr);
  console.log("ServiceManager:", serviceManagerAddr);
  console.log("TrafficCardNFT:", trafficCardNFTAddr);
  console.log("Payment:", paymentAddr);
  console.log("Deposit:", depositAddr);
  console.log("Oracle:", oracleAddr);

  const factories: Record<string, any> = {
    MockUSDT, FeeManager, UserRegistry, ServiceManager, TrafficCardNFT, Payment, Deposit, Oracle
  };
  const abiHash: Record<string, string> = {};
  for (const [name, factory] of Object.entries(factories)) {
    const abiJson = (factory as any).interface.formatJson();
    abiHash[name] = ethers.keccak256(ethers.toUtf8Bytes(abiJson));
  }

  const addresses = {
    chainId: network.config.chainId || 31337,
    proxies: {
      MockUSDT: usdtAddr,
      FeeManager: feeManagerAddr,
      UserRegistry: userRegistryAddr,
      ServiceManager: serviceManagerAddr,
      TrafficCardNFT: trafficCardNFTAddr,
      Payment: paymentAddr,
      Deposit: depositAddr,
      Oracle: oracleAddr
    },
    implementations: {
      FeeManager: await upgrades.erc1967.getImplementationAddress(feeManagerAddr),
      UserRegistry: await upgrades.erc1967.getImplementationAddress(userRegistryAddr),
      ServiceManager: await upgrades.erc1967.getImplementationAddress(serviceManagerAddr),
      TrafficCardNFT: await upgrades.erc1967.getImplementationAddress(trafficCardNFTAddr),
      Payment: await upgrades.erc1967.getImplementationAddress(paymentAddr),
      Deposit: await upgrades.erc1967.getImplementationAddress(depositAddr),
      Oracle: await upgrades.erc1967.getImplementationAddress(oracleAddr)
    },
    usdt: usdtAddr,
    usdtDecimals,
    abiHash,
    storageLayout: {
      _note: "Fresh deploy; storage layout baseline for future upgrades.",
      manifest: ".openzeppelin/"
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
