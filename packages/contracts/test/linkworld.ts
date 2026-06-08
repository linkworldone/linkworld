import { ethers } from "hardhat";
import { parseEther, ZeroAddress } from "ethers";
import { expect } from "chai";

function feeAt150bps(amount) { return (amount * 150n) / 10000n; }
const ZA = ZeroAddress;

describe("LinkWorld Contracts Tests", function () {
  let userRegistry, feeManager, deposit, serviceManager, payment, oracle, trafficCardNFT, mockUSDT;
  let owner, user1, user2, user3;

  // USDT 6 位精度辅助：n USDT -> 最小单位
  function usdt(n) {
    const [int, frac = ""] = String(n).split(".");
    const fracPadded = (frac + "000000").slice(0, 6);
    return BigInt(int) * 1_000_000n + BigInt(fracPadded || "0");
  }

  async function deployAndWire() {
    [owner, user1, user2, user3] = await ethers.getSigners();

    const MockUSDTFactory = await ethers.getContractFactory("MockUSDT");
    mockUSDT = await MockUSDTFactory.deploy();

    const FeeManagerFactory = await ethers.getContractFactory("FeeManager");
    feeManager = await FeeManagerFactory.deploy();
    await feeManager.initialize(150);

    const UserRegistryFactory = await ethers.getContractFactory("UserRegistry");
    userRegistry = await UserRegistryFactory.deploy();
    await userRegistry.initialize();

    const ServiceManagerFactory = await ethers.getContractFactory("ServiceManager");
    serviceManager = await ServiceManagerFactory.deploy();
    await serviceManager.initialize();

    const PaymentFactory = await ethers.getContractFactory("Payment");
    payment = await PaymentFactory.deploy();
    await payment.initialize(
      await feeManager.getAddress(),
      owner.address,
      await mockUSDT.getAddress(),
      await serviceManager.getAddress()
    );

    const DepositFactory = await ethers.getContractFactory("Deposit");
    deposit = await DepositFactory.deploy();
    await deposit.initialize(await userRegistry.getAddress(), await mockUSDT.getAddress());

    const OracleFactory = await ethers.getContractFactory("Oracle");
    oracle = await OracleFactory.deploy();
    await oracle.initialize();

    const TrafficCardNFTFactory = await ethers.getContractFactory("TrafficCardNFT");
    trafficCardNFT = await TrafficCardNFTFactory.deploy();
    await trafficCardNFT.initialize();

    await deposit.setOracle(await oracle.getAddress());
    await deposit.setTrafficCardNFT(await trafficCardNFT.getAddress());
    await oracle.setDeposit(await deposit.getAddress());

    // Payment 授权 oracle（createBill / applyTrafficCardToBill 的 onlyOracle，B2）
    await payment.setOracle(await oracle.getAddress());
  }

  describe("FeeManager", function () {
    beforeEach(deployAndWire);

    it("FE-01  initializer sets feeRate to 150 bps", async function () {
      expect(await feeManager.getFeeRate()).to.equal(150);
    });

    it("FE-02  1 ETH 1.5% = 0.015 ETH", async function () {
      expect(await feeManager.calculateFee(parseEther("1"))).to.equal(parseEther("0.015"));
    });

    it("FE-03  calculateFee(0) = 0", async function () {
      expect(await feeManager.calculateFee(0)).to.equal(0);
    });

    it("FE-04  10 ETH 1.5% = 0.15 ETH", async function () {
      expect(await feeManager.calculateFee(parseEther("10"))).to.equal(parseEther("0.15"));
    });

    it("FE-05  owner raises fee to 500 bps (5%)", async function () {
      await feeManager.setFeeRate(500);
      expect(await feeManager.getFeeRate()).to.equal(500);
    });

    it("FE-06  MAX_FEE_RATE=1000 accepted", async function () {
      await feeManager.setFeeRate(1000);
      expect(await feeManager.getFeeRate()).to.equal(1000);
    });

    it("FE-07  above 1000 revert", async function () {
      await expect(feeManager.setFeeRate(1001)).to.be.revertedWith("Fee too high");
    });

    it("FE-08  user1.setFeeRate revert", async function () {
      await expect(feeManager.connect(user1).setFeeRate(999)).to.be.reverted;
    });

    it("FE-09  FEE_DENOMINATOR = 10000", async function () {
      expect(Number(await feeManager.FEE_DENOMINATOR())).to.equal(10000);
    });
  });

  describe("UserRegistry", function () {
    beforeEach(deployAndWire);

    it("UR-01  register success", async function () {
      await userRegistry.connect(user1).register("u1@linkworld.io");
      expect(await userRegistry.isRegistered(user1.address)).to.be.true;
      expect(await userRegistry.ownerOf(0)).to.equal(user1.address);
    });

    it("UR-02  duplicate register revert", async function () {
      await userRegistry.connect(user1).register("dup@linkworld.io");
      await expect(userRegistry.connect(user1).register("dup2@linkworld.io"))
        .to.be.revertedWith("Already registered");
    });

    it("UR-03  unregistered user isRegistered=false", async function () {
      expect(await userRegistry.isRegistered(user2.address)).to.be.false;
    });

    it("UR-04  name=LinkWorld Identity, symbol=LWID", async function () {
      expect(await userRegistry.name()).to.equal("LinkWorld Identity");
      expect(await userRegistry.symbol()).to.equal("LWID");
    });
  });

  describe("ServiceManager", function () {
    beforeEach(deployAndWire);

    it("SM-01  initial active operators > 0", async function () {
      const ops = await serviceManager.getActiveOperators();
      expect(ops.length).to.be.greaterThan(0);
    });

it("SM-02  getOperatorsByCountry US non-empty", async function () {
      const us = await serviceManager.getOperatorsByCountry("US");
      expect(us.length).to.be.greaterThan(0);
    });

    it("SM-03  addOperator creates new entry", async function () {
      const tx = await serviceManager.addOperator("TestOp", "TestRegion", "XX", parseEther("0.01"), owner.address);
      const receipt = await tx.wait();
      const event = receipt.logs.find(l => l.fragment?.name === "OperatorAdded");
      expect(event?.args?.operatorId).to.equal(12);
    });
  });

  describe("Deposit", function () {
    beforeEach(async function () {
      await deployAndWire();
      await userRegistry.connect(user1).register("u1@linkworld.io");
      await mockUSDT.mint(user1.address, usdt(1000));
    });

    it("DP-01  deposit and getDepositAmount", async function () {
      const amount = usdt(50);
      await mockUSDT.connect(user1).approve(await deposit.getAddress(), amount);
      await deposit.connect(user1).deposit(amount);
      expect(await deposit.getDepositAmount(user1.address)).to.equal(amount);
    });

    it("DP-02  unregistered user revert", async function () {
      await mockUSDT.mint(user2.address, usdt(100));
      await mockUSDT.connect(user2).approve(await deposit.getAddress(), usdt(50));
      await expect(
        deposit.connect(user2).deposit(usdt(50))
      ).to.be.revertedWith("Not registered");
    });

it("DP-03  setOracle stores address", async function () {
      await deposit.setOracle(await oracle.getAddress());
      expect(await deposit.oracle()).to.equal(await oracle.getAddress());
    });
  });

  describe("Payment", function () {
    beforeEach(async function () {
      await deployAndWire();
      // createBill 现为 onlyOracle（B2）：用 EOA user3 充当 oracle 直接调用；
      // 并为 operator 1 设置分账地址，通过 createBill fail-fast 校验。
      await payment.setOracle(user3.address);
      await serviceManager.setOperatorPaymentAddress(1, user2.address);
    });

    it("PM-01  createBill (onlyOracle) emits BillCreated", async function () {
      const amount = usdt(100);
      const total = amount + feeAt150bps(amount);
      await expect(payment.connect(user3).createBill(user1.address, 1, amount))
        .to.emit(payment, "BillCreated").withArgs(0, user1.address, total, feeAt150bps(amount));
    });

    it("PM-02  getUserBills returns created bills", async function () {
      const amount = usdt(100);
      await payment.connect(user3).createBill(user1.address, 1, amount);
      const bills = await payment.getUserBills(user1.address);
      expect(bills.length).to.equal(1);
      expect(bills[0].amount).to.equal(amount);
    });
  });

  describe("Oracle", function () {
    beforeEach(deployAndWire);

    it("OR-01  deposit address stored", async function () {
      expect(await oracle.deposit()).to.equal(await deposit.getAddress());
    });

    it("OR-02  verifyServiceActive returns false for unregistered", async function () {
      const active = await oracle.verifyServiceActive(user1.address);
      expect(active).to.be.false;
    });
  });

  describe("TrafficCardNFT", function () {
    beforeEach(deployAndWire);

    it("TC-01  name=LinkWorld Traffic Card, symbol=LWTC", async function () {
      expect(await trafficCardNFT.name()).to.equal("LinkWorld Traffic Card");
      expect(await trafficCardNFT.symbol()).to.equal("LWTC");
    });

    it("TC-02  mint ownerOf(0)=user1", async function () {
      await trafficCardNFT.mint(user1.address, 512, "https://example.com/0");
      expect(await trafficCardNFT.ownerOf(0)).to.equal(user1.address);
    });

    it("TC-03  mint emits CardMinted", async function () {
      await expect(trafficCardNFT.mint(user1.address, 512, "uri"))
        .to.emit(trafficCardNFT, "CardMinted").withArgs(user1.address, 0n, 512n);
    });

it("TC-04  getUserCardCount zero initially", async function () {
      expect(await trafficCardNFT.getUserCardCount(user1.address)).to.equal(0);
    });

    it("TC-05  getCardInfo returns correct dataAmount", async function () {
      await trafficCardNFT.mint(user1.address, 512, "https://example.com/0");
      const info = await trafficCardNFT.getCardInfo(0);
      expect(info.dataAmount).to.equal(512);
      expect(info.isDestroyed).to.be.false;
    });
  });
});


