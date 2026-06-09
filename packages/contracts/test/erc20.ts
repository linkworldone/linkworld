import { ethers } from "hardhat";
import { expect } from "chai";
import { time } from "@nomicfoundation/hardhat-network-helpers";

// USDT 6 位精度辅助：n USDT -> 最小单位
function usdt(n: string | number): bigint {
  const [int, frac = ""] = String(n).split(".");
  const fracPadded = (frac + "000000").slice(0, 6);
  return BigInt(int) * 1_000_000n + BigInt(fracPadded || "0");
}

describe("Deposit ERC20 (T2)", function () {
  let userRegistry: any, deposit: any, mockUSDT: any;
  let owner: any, user1: any, user2: any;

  beforeEach(async function () {
    [owner, user1, user2] = await ethers.getSigners();

    const MockUSDTFactory = await ethers.getContractFactory("MockUSDT");
    mockUSDT = await MockUSDTFactory.deploy();

    const UserRegistryFactory = await ethers.getContractFactory("UserRegistry");
    userRegistry = await UserRegistryFactory.deploy();
    await userRegistry.initialize();

    const DepositFactory = await ethers.getContractFactory("Deposit");
    deposit = await DepositFactory.deploy();
    await deposit.initialize(await userRegistry.getAddress(), await mockUSDT.getAddress());

    // user1 注册 + 发币
    await userRegistry.connect(user1).register("u1@linkworld.io");
    await mockUSDT.mint(user1.address, usdt(1000));
  });

  describe("MockUSDT 精度", function () {
    it("DEC-01  decimals() == 6", async function () {
      expect(await mockUSDT.decimals()).to.equal(6);
    });
  });

  describe("最小存款额（验收 A.2）", function () {
    it("MIN-01  9.999999 USDT 拒绝", async function () {
      await mockUSDT.connect(user1).approve(await deposit.getAddress(), usdt("9.999999"));
      await expect(
        deposit.connect(user1).deposit(usdt("9.999999"))
      ).to.be.revertedWith("Below min deposit");
    });

    it("MIN-02  10.000000 USDT 通过：_deposits 增、合约 USDT 余额增", async function () {
      const amount = usdt("10");
      await mockUSDT.connect(user1).approve(await deposit.getAddress(), amount);
      await deposit.connect(user1).deposit(amount);

      expect(await deposit.getDepositAmount(user1.address)).to.equal(amount);
      expect(await mockUSDT.balanceOf(await deposit.getAddress())).to.equal(amount);
    });
  });

  describe("approve / transferFrom", function () {
    it("ERC-01  未 approve → deposit revert", async function () {
      await expect(
        deposit.connect(user1).deposit(usdt("10"))
      ).to.be.reverted;
    });

    it("ERC-02  approve 后 deposit 成功，余额正确转移", async function () {
      const amount = usdt("50");
      const before = await mockUSDT.balanceOf(user1.address);
      await mockUSDT.connect(user1).approve(await deposit.getAddress(), amount);
      await deposit.connect(user1).deposit(amount);

      expect(await mockUSDT.balanceOf(user1.address)).to.equal(before - amount);
      expect(await mockUSDT.balanceOf(await deposit.getAddress())).to.equal(amount);
      expect(await deposit.getDepositAmount(user1.address)).to.equal(amount);
    });

    it("ERC-02b  emit DepositMade(user, amount)", async function () {
      const amount = usdt("10");
      await mockUSDT.connect(user1).approve(await deposit.getAddress(), amount);
      await expect(deposit.connect(user1).deposit(amount))
        .to.emit(deposit, "DepositMade")
        .withArgs(user1.address, amount);
    });
  });

  describe("提取（锁仓 + CEI safeTransfer）", function () {
    it("ERC-03  锁仓未到 revert；到期后退回本金、_deposits 清零", async function () {
      const amount = usdt("100");
      await mockUSDT.connect(user1).approve(await deposit.getAddress(), amount);
      await deposit.connect(user1).deposit(amount);

      // 锁仓未到 -> revert
      await expect(deposit.connect(user1).withdraw()).to.be.revertedWith("Lock not expired");

      // 快进 30 天 + 1 秒
      await time.increase(30 * 24 * 3600 + 1);

      const before = await mockUSDT.balanceOf(user1.address);
      await deposit.connect(user1).withdraw();

      expect(await mockUSDT.balanceOf(user1.address)).to.equal(before + amount);
      expect(await deposit.getDepositAmount(user1.address)).to.equal(0n);
      expect(await deposit.getLockExpiry(user1.address)).to.equal(0n);
      expect(await mockUSDT.balanceOf(await deposit.getAddress())).to.equal(0n);
    });
  });

  describe("注册校验（REG，重写自旧 DP 用例）", function () {
    it("REG-02  未注册用户 deposit revert Not registered", async function () {
      await mockUSDT.mint(user2.address, usdt(100));
      await mockUSDT.connect(user2).approve(await deposit.getAddress(), usdt(10));
      await expect(
        deposit.connect(user2).deposit(usdt(10))
      ).to.be.revertedWith("Not registered");
    });
  });
});
