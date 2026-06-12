import { ethers } from "hardhat";
import { expect } from "chai";
import { time } from "@nomicfoundation/hardhat-network-helpers";

// USDT 6 位精度辅助
function usdt(n: string | number): bigint {
  const [int, frac = ""] = String(n).split(".");
  const fracPadded = (frac + "000000").slice(0, 6);
  return BigInt(int) * 1_000_000n + BigInt(fracPadded || "0");
}

const LOCK = 30 * 24 * 60 * 60; // 30 days

// 充值即发卡：把 NFT owner 转给 deposit，让 deposit 内 trafficCardNFT.mint(onlyOwner) 可调
async function deployNFTForDeposit(deposit: any) {
  const TrafficCardNFTFactory = await ethers.getContractFactory("TrafficCardNFT");
  const nft = await TrafficCardNFTFactory.deploy();
  await nft.initialize();
  await deposit.setTrafficCardNFT(await nft.getAddress());
  await nft.setDepositContract(await deposit.getAddress());
  await nft.transferOwnership(await deposit.getAddress());
  return nft;
}

// ============================================================
// B7 非标 USDT：SafeERC20 分支（USDT-01 无返回值 / USDT-02 返回 false）
// ============================================================
describe("B7 非标 USDT SafeERC20 (T6)", function () {
  let userRegistry: any, depositNoRet: any, depositFalse: any;
  let noRetUSDT: any, falseUSDT: any;
  let owner: any, user1: any;

  beforeEach(async function () {
    [owner, user1] = await ethers.getSigners();

    const UserRegistryFactory = await ethers.getContractFactory("UserRegistry");
    userRegistry = await UserRegistryFactory.deploy();
    await userRegistry.initialize();
    await userRegistry.connect(user1).register("u1@linkworld.io");

    const NoRetFactory = await ethers.getContractFactory("NoReturnUSDT");
    noRetUSDT = await NoRetFactory.deploy();

    const FalseFactory = await ethers.getContractFactory("FalseReturnUSDT");
    falseUSDT = await FalseFactory.deploy();

    const DepositFactory = await ethers.getContractFactory("Deposit");
    depositNoRet = await DepositFactory.deploy();
    await depositNoRet.initialize(await userRegistry.getAddress(), await noRetUSDT.getAddress());
    await deployNFTForDeposit(depositNoRet);

    depositFalse = await DepositFactory.deploy();
    await depositFalse.initialize(await userRegistry.getAddress(), await falseUSDT.getAddress());
    await deployNFTForDeposit(depositFalse);
  });

  it("USDT-01  transfer/transferFrom 无返回值的非标 USDT → deposit 仍正确入账（SafeERC20 不 revert）", async function () {
    const amount = usdt("50");
    await noRetUSDT.mint(user1.address, usdt("1000"));
    await noRetUSDT.connect(user1).approve(await depositNoRet.getAddress(), amount);

    await expect(depositNoRet.connect(user1).deposit(amount)).to.not.be.reverted;

    expect(await depositNoRet.getDepositAmount(user1.address)).to.equal(amount);
    expect(await noRetUSDT.balanceOf(await depositNoRet.getAddress())).to.equal(amount);
  });

  it("USDT-01b 无返回值非标 USDT → withdraw 到期后 safeTransfer 正确退回（不 revert）", async function () {
    const amount = usdt("50");
    await noRetUSDT.mint(user1.address, usdt("1000"));
    await noRetUSDT.connect(user1).approve(await depositNoRet.getAddress(), amount);
    await depositNoRet.connect(user1).deposit(amount);

    await time.increase(LOCK + 1);
    const before = await noRetUSDT.balanceOf(user1.address);
    await expect(depositNoRet.connect(user1).withdraw(0)).to.not.be.reverted;
    expect(await noRetUSDT.balanceOf(user1.address)).to.equal(before + amount);
    expect(await depositNoRet.getDepositAmount(user1.address)).to.equal(0n);
  });

  it("USDT-02  transferFrom 返回 false 的恶意 USDT → deposit 整体 revert（SafeERC20 捕获，不静默失败）", async function () {
    const amount = usdt("50");
    await falseUSDT.mint(user1.address, usdt("1000"));
    await falseUSDT.connect(user1).approve(await depositFalse.getAddress(), amount);

    // SafeERC20 检测到 transferFrom 返回 false → revert（SafeERC20FailedOperation）
    await expect(depositFalse.connect(user1).deposit(amount)).to.be.reverted;
    // 资金未被记账（无静默失败）
    expect(await depositFalse.getDepositAmount(user1.address)).to.equal(0n);
  });

  it("USDT-02b payBill 路径：transferFrom 返回 false 的恶意 USDT → payBill 整体 revert（SafeERC20 捕获，账单仍 unpaid）", async function () {
    // Payment 用 FalseReturnUSDT：createBill 不转账（仅记账）可成功，payBill 时 safeTransferFrom 返回 false → revert
    const FeeManagerFactory = await ethers.getContractFactory("FeeManager");
    const feeManager = await FeeManagerFactory.deploy();
    await feeManager.initialize(150);

    const ServiceManagerFactory = await ethers.getContractFactory("ServiceManager");
    const serviceManager = await ServiceManagerFactory.deploy();
    await serviceManager.initialize();
    await serviceManager.setOperatorPaymentAddress(1, owner.address);

    const PaymentFactory = await ethers.getContractFactory("Payment");
    const payment = await PaymentFactory.deploy();
    await payment.initialize(
      await feeManager.getAddress(),
      owner.address,
      await falseUSDT.getAddress(),
      await serviceManager.getAddress()
    );
    await payment.setOracle(owner.address); // owner 充当 oracle 直接调 createBill

    await falseUSDT.mint(user1.address, usdt("1000"));
    await falseUSDT.connect(user1).approve(await payment.getAddress(), usdt("1000"));

    await payment.createBill(user1.address, 1, usdt("100"));
    // payBill 走 safeTransferFrom → FalseReturnUSDT 返回 false → SafeERC20 revert
    await expect(payment.connect(user1).payBill(0)).to.be.reverted;
    // 账单仍未支付（无静默失败）
    const bills = await payment.getUserBills(user1.address);
    expect(bills[0].isPaid).to.equal(false);
  });
});

// ============================================================
// REG-01 逐笔独立锁不变量（每次 deposit 独立计时 30 天，互不影响）
// ============================================================
describe("REG-01 逐笔独立锁不变量 (T6)", function () {
  let userRegistry: any, deposit: any, mockUSDT: any;
  let owner: any, user1: any;

  beforeEach(async function () {
    [owner, user1] = await ethers.getSigners();

    const MockUSDTFactory = await ethers.getContractFactory("MockUSDT");
    mockUSDT = await MockUSDTFactory.deploy();

    const UserRegistryFactory = await ethers.getContractFactory("UserRegistry");
    userRegistry = await UserRegistryFactory.deploy();
    await userRegistry.initialize();
    await userRegistry.connect(user1).register("reg01@linkworld.io");

    const DepositFactory = await ethers.getContractFactory("Deposit");
    deposit = await DepositFactory.deploy();
    await deposit.initialize(await userRegistry.getAddress(), await mockUSDT.getAddress());
    await deployNFTForDeposit(deposit);

    await mockUSDT.mint(user1.address, usdt("1000"));
  });

  it("REG-01a  每笔独立计时：晚笔 unlockAt 比早笔晚（不叠加到同一到期点）", async function () {
    const amt = usdt("10");
    await mockUSDT.connect(user1).approve(await deposit.getAddress(), usdt("1000"));

    const tx1 = await deposit.connect(user1).deposit(amt);
    const b1 = await ethers.provider.getBlock(tx1.blockNumber!);

    await time.increase(1 * 24 * 3600); // +1 天后复存
    const tx2 = await deposit.connect(user1).deposit(amt);
    const b2 = await ethers.provider.getBlock(tx2.blockNumber!);

    const tranches = await deposit.getTranches(user1.address);
    expect(tranches.length).to.equal(2);
    // 每笔 = 本笔存款时间 + 30d（互相独立，非叠加）
    expect(tranches[0].unlockAt).to.equal(BigInt(b1!.timestamp) + BigInt(LOCK));
    expect(tranches[1].unlockAt).to.equal(BigInt(b2!.timestamp) + BigInt(LOCK));
    // 锁仓总额累加
    expect(await deposit.getDepositAmount(user1.address)).to.equal(amt * 2n);
  });

  it("REG-01b  早笔到期可单独取回，晚笔仍受锁；取回不影响晚笔", async function () {
    const amt = usdt("10");
    await mockUSDT.connect(user1).approve(await deposit.getAddress(), usdt("1000"));

    await deposit.connect(user1).deposit(amt); // tranche 0
    await time.increase(LOCK - 100);
    await deposit.connect(user1).deposit(amt); // tranche 1（晚 ~30d-100s）
    await time.increase(200); // tranche0 到期，tranche1 未到期

    await deposit.connect(user1).withdraw(0);
    expect(await deposit.getDepositAmount(user1.address)).to.equal(amt);
    await expect(deposit.connect(user1).withdraw(1)).to.be.revertedWith("Lock not expired");
  });
});

// ============================================================
// 边界回归：withdraw 非法 index → "Invalid tranche"
// ============================================================
describe("ERC withdraw 边界 (T6)", function () {
  let userRegistry: any, deposit: any, mockUSDT: any;
  let owner: any, user1: any;

  beforeEach(async function () {
    [owner, user1] = await ethers.getSigners();
    const MockUSDTFactory = await ethers.getContractFactory("MockUSDT");
    mockUSDT = await MockUSDTFactory.deploy();
    const UserRegistryFactory = await ethers.getContractFactory("UserRegistry");
    userRegistry = await UserRegistryFactory.deploy();
    await userRegistry.initialize();
    const DepositFactory = await ethers.getContractFactory("Deposit");
    deposit = await DepositFactory.deploy();
    await deposit.initialize(await userRegistry.getAddress(), await mockUSDT.getAddress());
  });

  it("ERC-04  无任何笔的用户 withdraw(0) → revert Invalid tranche", async function () {
    await expect(deposit.connect(user1).withdraw(0)).to.be.revertedWith("Invalid tranche");
  });
});

// ============================================================
// GAS-01 批量压测：deposit 充值即发卡（量单笔 deposit 的 gas，含逐张 mint）
// ============================================================
describe("GAS-01 充值即发卡 gas 压测 (T6)", function () {
  this.timeout(600000);

  const BLOCK_GAS_LIMIT = 30_000_000n; // hardhat 默认 / 主流链区块上限参考

  let userRegistry: any, deposit: any, mockUSDT: any, trafficCardNFT: any;
  let owner: any;

  beforeEach(async function () {
    [owner] = await ethers.getSigners();

    const MockUSDTFactory = await ethers.getContractFactory("MockUSDT");
    mockUSDT = await MockUSDTFactory.deploy();

    const UserRegistryFactory = await ethers.getContractFactory("UserRegistry");
    userRegistry = await UserRegistryFactory.deploy();
    await userRegistry.initialize();

    const DepositFactory = await ethers.getContractFactory("Deposit");
    deposit = await DepositFactory.deploy();
    await deposit.initialize(await userRegistry.getAddress(), await mockUSDT.getAddress());

    trafficCardNFT = await deployNFTForDeposit(deposit);
  });

  it("GAS-01  各档位单笔 deposit gas（含逐张 mint）均 < 区块上限，每卡成本趋稳", async function () {
    const tiers: [number, number][] = [
      [10, 1],
      [20, 2],
      [50, 5],
      [100, 10],
    ];
    const samples: { tier: number; cards: number; gas: bigint; perCard: bigint }[] = [];

    let i = 0;
    for (const [tier, cards] of tiers) {
      const w = ethers.Wallet.createRandom().connect(ethers.provider);
      await owner.sendTransaction({ to: w.address, value: ethers.parseEther("1") });
      await userRegistry.connect(w).register(`gas-${Date.now()}-${i++}@linkworld.io`);
      await mockUSDT.mint(w.address, usdt(tier));
      await mockUSDT.connect(w).approve(await deposit.getAddress(), usdt(tier));

      const tx = await deposit.connect(w).deposit(usdt(tier));
      const receipt = await tx.wait();
      const gas = receipt!.gasUsed;
      samples.push({ tier, cards, gas, perCard: gas / BigInt(cards) });

      expect(await trafficCardNFT.getUserCardCount(w.address)).to.equal(BigInt(cards));
    }

    console.log("\n  [GAS-01] deposit 充值即发卡压测:");
    for (const s of samples) {
      console.log(`    tier=${s.tier} USDT → ${s.cards} 卡: gasUsed=${s.gas}  per-card≈${s.perCard}`);
    }

    // 最大档（100 USDT → 10 张）单笔 < 区块上限（可上链）
    const top = samples[samples.length - 1];
    expect(top.gas).to.be.lessThan(BLOCK_GAS_LIMIT);
    // per-card 趋稳：最大档每卡成本 ≤ 最小档（基础 overhead 摊薄）
    expect(top.perCard).to.be.lessThanOrEqual(samples[0].perCard);
  });
});
