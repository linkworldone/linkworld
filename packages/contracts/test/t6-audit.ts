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

    depositFalse = await DepositFactory.deploy();
    await depositFalse.initialize(await userRegistry.getAddress(), await falseUSDT.getAddress());
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
    await expect(depositNoRet.connect(user1).withdraw()).to.not.be.reverted;
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
// REG-01 锁仓续期不变量（design §八 回归块）
// ============================================================
describe("REG-01 锁仓续期不变量 (T6)", function () {
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

    await mockUSDT.mint(user1.address, usdt("1000"));
  });

  it("REG-01a  未到期复存 → lockExpiry 在原到期点叠加 +30d（不重置到 now+30d）", async function () {
    const amt = usdt("10");
    await mockUSDT.connect(user1).approve(await deposit.getAddress(), usdt("1000"));

    await deposit.connect(user1).deposit(amt);
    const expiry1 = await deposit.getLockExpiry(user1.address);

    // 未到期：仅快进 1 天后复存
    await time.increase(1 * 24 * 3600);
    await deposit.connect(user1).deposit(amt);
    const expiry2 = await deposit.getLockExpiry(user1.address);

    // 叠加：expiry2 == expiry1 + 30d（在原到期点上加，而非 now+30d）
    expect(expiry2).to.equal(expiry1 + BigInt(30 * 24 * 3600));
    // 存款累加
    expect(await deposit.getDepositAmount(user1.address)).to.equal(amt * 2n);
  });

  it("REG-01b  已到期复存 → lockExpiry 重置为 now+30d", async function () {
    const amt = usdt("10");
    await mockUSDT.connect(user1).approve(await deposit.getAddress(), usdt("1000"));

    await deposit.connect(user1).deposit(amt);

    // 快进超过锁仓期，使其到期
    await time.increase(LOCK + 100);
    const tx = await deposit.connect(user1).deposit(amt);
    const block = await ethers.provider.getBlock(tx.blockNumber!);
    const now = BigInt(block!.timestamp);

    const expiry = await deposit.getLockExpiry(user1.address);
    expect(expiry).to.equal(now + BigInt(30 * 24 * 3600));
  });
});

// ============================================================
// 边界回归：withdraw 无存款 revert "No deposit"（合约 L73，design §3.1）
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

  it("ERC-04  无存款用户 withdraw（锁仓已过：lockExpiry=0）→ revert No deposit", async function () {
    // 从未存款者 lockExpiry=0 ≤ now，先过 require(lockExpiry) 再撞 require(_deposits>0)
    await expect(deposit.connect(user1).withdraw()).to.be.revertedWith("No deposit");
  });
});

// ============================================================
// GAS-01 批量压测：issueMonthlyTrafficCards / monthlySettlement
// 量出单批安全上限 N（写入 handoff）
// ============================================================
describe("GAS-01 批量 gas 压测 (T6)", function () {
  this.timeout(600000);

  // 单批安全预算：取常见区块 gas limit 30,000,000 的 ~50% 作单笔 tx 安全余量
  // （留波动/L1 calldata 成本/并发空间）。据此 + 实测 per-user gas 外推单批安全上限 N。
  const BLOCK_GAS_LIMIT = 30_000_000n;       // hardhat 默认 / 主流链区块上限参考
  const SAFE_GAS_BUDGET = 15_000_000n;       // 单笔 tx 安全预算（区块上限 50%）

  let userRegistry: any, deposit: any, mockUSDT: any, trafficCardNFT: any;
  let owner: any;
  let oracleSigner: any;

  // 批量构造 N 个满足「注册+存款+锁仓到期」的随机钱包用户
  async function buildEligibleUsers(n: number): Promise<string[]> {
    const provider = ethers.provider;
    const addrs: string[] = [];
    const depAddr = await deposit.getAddress();
    const minDep = usdt("10");

    for (let i = 0; i < n; i++) {
      const w = ethers.Wallet.createRandom().connect(provider);
      // fund ETH for gas
      await owner.sendTransaction({ to: w.address, value: ethers.parseEther("1") });
      await userRegistry.connect(w).register(`gas-${Date.now()}-${i}@linkworld.io`);
      await mockUSDT.mint(w.address, minDep);
      await mockUSDT.connect(w).approve(depAddr, minDep);
      await deposit.connect(w).deposit(minDep);
      addrs.push(w.address);
    }
    return addrs;
  }

  beforeEach(async function () {
    [owner] = await ethers.getSigners();
    oracleSigner = owner; // owner 充当 oracle（EOA 直接调）

    const MockUSDTFactory = await ethers.getContractFactory("MockUSDT");
    mockUSDT = await MockUSDTFactory.deploy();

    const UserRegistryFactory = await ethers.getContractFactory("UserRegistry");
    userRegistry = await UserRegistryFactory.deploy();
    await userRegistry.initialize();

    const DepositFactory = await ethers.getContractFactory("Deposit");
    deposit = await DepositFactory.deploy();
    await deposit.initialize(await userRegistry.getAddress(), await mockUSDT.getAddress());

    const TrafficCardNFTFactory = await ethers.getContractFactory("TrafficCardNFT");
    trafficCardNFT = await TrafficCardNFTFactory.deploy();
    await trafficCardNFT.initialize();

    await deposit.setTrafficCardNFT(await trafficCardNFT.getAddress());
    await trafficCardNFT.setDepositContract(await deposit.getAddress());
    await trafficCardNFT.transferOwnership(await deposit.getAddress());
    await deposit.setOracle(oracleSigner.address);
  });

  it("GAS-01  issueMonthlyTrafficCards 多梯度 gas 压测 → 量单批安全上限 N", async function () {
    const gradients = [10, 25, 50];
    const samples: { n: number; gas: bigint; perUser: bigint }[] = [];

    for (const n of gradients) {
      const users = await buildEligibleUsers(n);
      await time.increase(LOCK + 1);

      const tx = await deposit.connect(oracleSigner).issueMonthlyTrafficCards(users);
      const receipt = await tx.wait();
      const gas = receipt!.gasUsed;
      samples.push({ n, gas, perUser: gas / BigInt(n) });

      // 校验全部发卡成功
      for (const a of users) {
        expect(await trafficCardNFT.getUserCardCount(a)).to.equal(1);
      }
    }

    // 用最大梯度 per-user（最保守，含基础 overhead 摊薄后趋稳）外推安全上限
    const top = samples[samples.length - 1];
    const perUser = top.perUser;
    const safeN = SAFE_GAS_BUDGET / perUser;

    console.log("\n  [GAS-01] issueMonthlyTrafficCards 压测:");
    for (const s of samples) {
      console.log(`    N=${s.n}: gasUsed=${s.gas}  per-user≈${s.perUser}`);
    }
    console.log(`    安全预算 ${SAFE_GAS_BUDGET} gas / per-user≈${perUser} → 单批安全上限 N≈${safeN}（建议取整 N≤50）`);

    // 断言：实测最大梯度 N=50 单批 < 区块 gas 上限（可上链）
    expect(top.gas).to.be.lessThan(BLOCK_GAS_LIMIT);
    // per-user 趋稳（线性，无意外二次方膨胀）：最大梯度 per-user ≤ 最小梯度 per-user
    expect(top.perUser).to.be.lessThanOrEqual(samples[0].perUser);
    // 安全上限至少容纳一个有意义的批量
    expect(safeN).to.be.greaterThanOrEqual(30n);
  });

  it("GAS-01b monthlySettlement 全链路（createBill+发卡+applyTrafficCardToBill 桩）多梯度 gas 压测", async function () {
    // monthlySettlement 需 Payment + Oracle wiring + operator paymentAddress
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
      await mockUSDT.getAddress(),
      await serviceManager.getAddress()
    );

    const OracleFactory = await ethers.getContractFactory("Oracle");
    const oracle = await OracleFactory.deploy();
    await oracle.initialize();

    // wiring：Oracle 为 monthlySettlement 触发者（onlyOwner=owner）
    await oracle.setDeposit(await deposit.getAddress());
    await oracle.setPayment(await payment.getAddress());
    await payment.setOracle(await oracle.getAddress());
    // Deposit 的 oracle 必须是 Oracle 合约（monthlySettlement 内部调 issueMonthlyTrafficCards）
    await deposit.setOracle(await oracle.getAddress());

    // monthlySettlement 单 user 含 createBill + 发卡(mint) + getUnpaidBills 二次循环 + applyTrafficCardToBill 桩，
    // per-user 远高于纯发卡，故用更小梯度（10/20/30）避免单 tx 超区块上限。
    const gradients = [10, 20, 30];
    const samples: { n: number; gas: bigint; perUser: bigint }[] = [];

    for (const n of gradients) {
      const users = await buildEligibleUsers(n);
      await time.increase(LOCK + 1);

      const operatorIds = users.map(() => 1n);
      const amounts = users.map(() => usdt("10"));

      const tx = await oracle.monthlySettlement(users, operatorIds, amounts);
      const receipt = await tx.wait();
      const gas = receipt!.gasUsed;
      samples.push({ n, gas, perUser: gas / BigInt(n) });

      // 全链路生效校验：账单已建 + 卡已发
      const bills = await payment.getUserBills(users[0]);
      expect(bills.length).to.be.greaterThan(0);
      expect(await trafficCardNFT.getUserCardCount(users[0])).to.equal(1);
    }

    const top = samples[samples.length - 1];
    const safeN = SAFE_GAS_BUDGET / top.perUser;

    console.log("\n  [GAS-01b] monthlySettlement 全链路压测（createBill+发卡+applyTrafficCardToBill 桩）:");
    for (const s of samples) {
      console.log(`    N=${s.n}: gasUsed=${s.gas}  per-user≈${s.perUser}`);
    }
    console.log(`    安全预算 ${SAFE_GAS_BUDGET} gas / per-user≈${top.perUser} → 单批安全上限 N≈${safeN}（全链路上限更紧，建议 N≤25）`);

    // 实测最大梯度单批 < 区块上限（可上链）
    expect(top.gas).to.be.lessThan(BLOCK_GAS_LIMIT);
    expect(safeN).to.be.greaterThanOrEqual(20n);
  });
});
