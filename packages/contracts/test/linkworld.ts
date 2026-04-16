import { ethers } from "hardhat";
import { expect } from "chai";

describe("LinkWorld Contracts", function () {
  let userRegistry: any;
  let feeManager: any;
  let deposit: any;
  let serviceManager: any;
  let payment: any;
  let oracle: any;
  let owner: any;
  let user1: any;
  let user2: any;

  beforeEach(async function () {
    [owner, user1, user2] = await ethers.getSigners();

    const UserRegistry = await ethers.getContractFactory("UserRegistry");
    userRegistry = await UserRegistry.deploy();
    await userRegistry.waitForDeployment();

    const FeeManager = await ethers.getContractFactory("FeeManager");
    feeManager = await FeeManager.deploy(250);
    await feeManager.waitForDeployment();

    const ServiceManager = await ethers.getContractFactory("ServiceManager");
    serviceManager = await ServiceManager.deploy();
    await serviceManager.waitForDeployment();

    const Payment = await ethers.getContractFactory("Payment");
    payment = await Payment.deploy(await feeManager.getAddress(), owner.address);
    await payment.waitForDeployment();

    const Oracle = await ethers.getContractFactory("Oracle");
    oracle = await Oracle.deploy(await payment.getAddress());
    await oracle.waitForDeployment();

    await payment.setOracle(await oracle.getAddress());

    const Deposit = await ethers.getContractFactory("Deposit");
    deposit = await Deposit.deploy(await userRegistry.getAddress());
    await deposit.waitForDeployment();
  });

  describe("UserRegistry", function () {
    it("should register a new user", async function () {
      await userRegistry.register("test@example.com");
      const isActive = await userRegistry.isRegistered(user1.address);
      expect(isActive).to.equal(true);
    });

    it("should get user info", async function () {
      await userRegistry.register("test@example.com");
      const userInfo = await userRegistry.getUserInfo(user1.address);
      expect(userInfo.wallet).to.equal(user1.address);
      expect(userInfo.email).to.equal("test@example.com");
    });

    it("should prevent duplicate registration", async function () {
      await userRegistry.register("test@example.com");
      await expect(userRegistry.register("test2@example.com")).to.be.revertedWith("Already registered");
    });
  });

  describe("FeeManager", function () {
    it("should deploy with 2.5% fee", async function () {
      expect(await feeManager.getFeeRate()).to.equal(250);
    });

    it("should calculate fee correctly", async function () {
      const fee = await feeManager.calculateFee(ethers.parseEther("1"));
      expect(fee).to.equal(ethers.parseEther("0.025"));
    });

    it("should allow owner to update fee rate", async function () {
      await feeManager.setFeeRate(300);
      expect(await feeManager.getFeeRate()).to.equal(300);
    });
  });

  describe("ServiceManager", function () {
    it("should deploy with 11 operators", async function () {
      const operators = await serviceManager.getActiveOperators();
      expect(operators.length).to.equal(11);
    });

    it("should get operators by country", async function () {
      const usOperators = await serviceManager.getOperatorsByCountry("US");
      expect(usOperators.length).to.greaterThan(0);
    });

    it("should activate user service", async function () {
      await serviceManager.connect(user1).activateService(1, "+1234567890", "password123");
      const service = await serviceManager.getUserService(user1.address);
      expect(service.isActive).to.equal(true);
    });

    it("should not allow duplicate service", async function () {
      await serviceManager.connect(user1).activateService(1, "+1234567890", "password123");
      await expect(
        serviceManager.connect(user1).activateService(2, "+9876543210", "password456")
      ).to.be.revertedWith("Service already active");
    });

    it("should deactivate service", async function () {
      await serviceManager.connect(user1).activateService(1, "+1234567890", "password123");
      await serviceManager.connect(user1).deactivateService();
      const service = await serviceManager.getUserService(user1.address);
      expect(service.isActive).to.equal(false);
    });
  });

  describe("Payment + Oracle", function () {
    it("should create bills via oracle", async function () {
      await oracle.submitUsage(user1.address, 1, 1000, 100);
      const bills = await payment.getUserBills(user1.address);
      expect(bills.length).to.equal(1);
    });

    it("should get unpaid bills", async function () {
      await oracle.submitUsage(user1.address, 1, 1000, 100);
      const unpaid = await payment.getUnpaidBills(user1.address);
      expect(unpaid.length).to.equal(1);
    });

    it("should auto settle bills", async function () {
      await oracle.submitUsage(user1.address, 1, 1000, 100);
      await payment.autoSettle([user1.address], [1], [ethers.parseEther("0.2")]);
      const bills = await payment.getUserBills(user1.address);
      expect(bills.length).to.equal(2);
    });
  });

  describe("Deposit", function () {
    it("should allow deposits from registered user", async function () {
      await userRegistry.connect(user1).register("test@example.com");
      await deposit.connect(user1).deposit({ value: ethers.parseEther("0.1") });
      const amount = await deposit.getDepositAmount(user1.address);
      expect(amount).to.equal(ethers.parseEther("0.1"));
    });

    it("should require registration before deposit", async function () {
      await expect(
        deposit.connect(user2).deposit({ value: ethers.parseEther("0.1") })
      ).to.be.revertedWith("Not registered");
    });
  });
});