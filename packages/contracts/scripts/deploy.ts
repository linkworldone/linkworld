import { ethers, ContractTransaction, upgrades, network } from "hardhat";
import * as fs from "fs";
import * as path from "path";

<<<<<<< HEAD
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
    return {
      maxFeePerGas,
      maxPriorityFeePerGas,
    };
  }

  const gasPriceValue = feeData.gasPrice ?? (await ethers.provider.getGasPrice());
  const gasPrice = multiplyFeeValue(gasPriceValue);
  return gasPrice ? { gasPrice } : {};
}
=======
// 11 个内置运营商 ID（ServiceManager.initialize 预置 1..11，全 active）
const BUILTIN_OPERATOR_IDS = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11];
>>>>>>> 48e2c2a1bafc65ef37ea369d3f85ef9eb004e69d

async function main() {
  const [deployer] = await ethers.getSigners();
  console.log("Deploying with:", deployer.address);

  // 0. MockUSDT(decimals=6) —— 本轮无正式 USDT，本地/测试网全用 mock（design §7.2 步骤0）
  const MockUSDT = await ethers.getContractFactory("MockUSDT");
  const usdt = (await MockUSDT.deploy()) as any;
  await usdt.waitForDeployment();
  const usdtAddr = await usdt.getAddress();
  const usdtDecimals = Number(await usdt.decimals());
  console.log("MockUSDT:", usdtAddr, "(decimals=" + usdtDecimals + ")");

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

  // 3. ServiceManager (无依赖) —— 须在 Payment 之前（Payment.initialize 注入 SM 地址）
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

  // 5. Payment —— initialize(feeManager, platformWallet, usdt, serviceManager)（design §7.2）
  //    依赖 FeeManager / ServiceManager / MockUSDT 须已部署
  const Payment = await ethers.getContractFactory("Payment");
  const payment = await upgrades.deployProxy(
    Payment,
    [feeManagerAddr, deployer.address, usdtAddr, serviceManagerAddr],
    { initializer: "initialize", kind: "uups" }
  ) as any;
  await payment.waitForDeployment();
  const paymentAddr = await payment.getAddress();
  console.log("Payment:", paymentAddr);

  // 6. Deposit —— initialize(userRegistry, usdt)（design §7.2）
  const Deposit = await ethers.getContractFactory("Deposit");
  const deposit = await upgrades.deployProxy(
    Deposit,
    [userRegistryAddr, usdtAddr],
    { initializer: "initialize", kind: "uups" }
  ) as any;
  await deposit.waitForDeployment();
  const depositAddr = await deposit.getAddress();
  console.log("Deposit:", depositAddr);

  // 7. Oracle (无外部依赖；payment/deposit 走 setter 后置注入)
  const Oracle = await ethers.getContractFactory("Oracle");
  const oracle = await upgrades.deployProxy(Oracle, [], {
    initializer: "initialize",
    kind: "uups"
  }) as any;
  await oracle.waitForDeployment();
  const oracleAddr = await oracle.getAddress();
  console.log("Oracle:", oracleAddr);

<<<<<<< HEAD
  // 关联合约 (v3 简化)
  const txOverrides = await getTxOverrides();
  await sendAndWait(trafficCardNFT.setDepositContract(depositAddr, txOverrides), "trafficCardNFT.setDepositContract");
  await sendAndWait(deposit.setTrafficCardNFT(trafficCardNFTAddr, txOverrides), "deposit.setTrafficCardNFT");
  await sendAndWait(deposit.setOracle(oracleAddr, txOverrides), "deposit.setOracle");
  await sendAndWait(payment.setPlatformWallet(deployer.address, txOverrides), "payment.setPlatformWallet");
  await sendAndWait(oracle.setDeposit(depositAddr, txOverrides), "oracle.setDeposit");
  console.log("Contracts linked");

  // 将 TrafficCardNFT 的 ownership 转给 Deposit 合约，让其可调用 mint
  await sendAndWait(trafficCardNFT.transferOwnership(depositAddr, txOverrides), "trafficCardNFT.transferOwnership");
=======
  // ===== 授权 wiring（design §7.0 / §7.2，★ = 本轮 B2 必补） =====
  // 现有 wiring（v3）
  await (await trafficCardNFT.setDepositContract(depositAddr)).wait();
  await (await deposit.setTrafficCardNFT(trafficCardNFTAddr)).wait();
  await (await deposit.setOracle(oracleAddr)).wait();           // Deposit.oracle = Oracle（链 B）
  await (await payment.setPlatformWallet(deployer.address)).wait();
  await (await oracle.setDeposit(depositAddr)).wait();          // Oracle.deposit = Deposit（链 B）
  // ★ B2 必补：自动结算账单权限链闭合
  await (await payment.setOracle(oracleAddr)).wait();           // Payment.oracle = Oracle（链 A/C）
  await (await oracle.setPayment(paymentAddr)).wait();          // Oracle.payment = Payment（链 A/C）
  console.log("Contracts linked (incl. payment.setOracle / oracle.setPayment)");

  // ★ B2 必补：11 个内置运营商注入非零分账地址（测试网用 deployer 派生地址）
  //   payBill / createBill fail-fast 要求 operator.paymentAddress != address(0)
  for (const id of BUILTIN_OPERATOR_IDS) {
    // deployer 派生地址：确定性、非零、便于测试网验证
    const payoutAddr = ethers.getAddress(
      "0x" + ethers.keccak256(
        ethers.solidityPacked(["address", "uint256"], [deployer.address, id])
      ).slice(-40)
    );
    await (await serviceManager.setOperatorPaymentAddress(id, payoutAddr)).wait();
  }
  console.log(`OperatorPaymentAddress set for ${BUILTIN_OPERATOR_IDS.length} built-in operators`);

  // 将 TrafficCardNFT 的 ownership 转给 Deposit 合约，让其可调用 mint（链 B）
  await (await trafficCardNFT.transferOwnership(depositAddr)).wait();
>>>>>>> 48e2c2a1bafc65ef37ea369d3f85ef9eb004e69d
  console.log("TrafficCardNFT ownership transferred to Deposit");

  // ===== 部署后断言校验（design §7.0.3 三条权限链前置；漏一项 → 自动结算 revert） =====
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
  // 每个 active operator.paymentAddress != 0
  for (const id of BUILTIN_OPERATOR_IDS) {
    const op = await serviceManager.getOperator(id);
    if (op.paymentAddress === ethers.ZeroAddress) {
      throw new Error(`operator ${id} paymentAddress is zero`);
    }
  }
  console.log("All wiring assertions passed (B2/B3 chains ready)");

  console.log("\n=== Deployment Summary ===");
  console.log("MockUSDT:", usdtAddr, "(decimals=" + usdtDecimals + ")");
  console.log("FeeManager:", feeManagerAddr);
  console.log("UserRegistry:", userRegistryAddr);
  console.log("ServiceManager:", serviceManagerAddr);
  console.log("TrafficCardNFT:", trafficCardNFTAddr);
  console.log("Payment:", paymentAddr);
  console.log("Deposit:", depositAddr);
  console.log("Oracle:", oracleAddr);

  // ===== deployments/<net>.json 产出（design §7.3：+usdt/usdtDecimals/storageLayout/abiHash） =====
  const factories: Record<string, any> = {
    FeeManager, UserRegistry, ServiceManager, TrafficCardNFT, Payment, Deposit, Oracle
  };
  const abiHash: Record<string, string> = {};
  for (const [name, factory] of Object.entries(factories)) {
    const abiJson = factory.interface.formatJson();
    abiHash[name] = ethers.keccak256(ethers.toUtf8Bytes(abiJson));
  }

  const addresses = {
    chainId: network.config.chainId || 31337,
    proxies: {
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
    // USDT 资金通道（本轮为 MockUSDT；正式上线替换为真实 USDT 地址）
    usdt: usdtAddr,
    usdtDecimals,
    // ABI 指纹（供后端比对，本轮 selector 已变后端须重生成）
    abiHash,
    // storage layout 冻结说明（本轮 fresh deploy，无升级包袱；记为后续 Round 升级基线）
    storageLayout: {
      _note: "本轮 fresh deploy，无 unsafeAllow；各 proxy 的 storage layout 以此次部署的实现合约为升级基线。OZ upgrades 插件在 .openzeppelin/ 下记录权威 layout 清单。",
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
