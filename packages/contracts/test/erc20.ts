import { ethers } from "hardhat";
import { expect } from "chai";
import { time } from "@nomicfoundation/hardhat-network-helpers";

// USDT 6 位精度辅助：n USDT -> 最小单位
function usdt(n: string | number): bigint {
  const [int, frac = ""] = String(n).split(".");
  const fracPadded = (frac + "000000").slice(0, 6);
  return BigInt(int) * 1_000_000n + BigInt(fracPadded || "0");
}

const LOCK = 30 * 24 * 60 * 60; // 30 days

describe("Deposit ERC20 (T2)", function () {
  let userRegistry: any, deposit: any, mockUSDT: any, trafficCardNFT: any;
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

    // 发卡需 NFT owner 转给 Deposit
    const TrafficCardNFTFactory = await ethers.getContractFactory("TrafficCardNFT");
    trafficCardNFT = await TrafficCardNFTFactory.deploy();
    await trafficCardNFT.initialize();
    await deposit.setTrafficCardNFT(await trafficCardNFT.getAddress());
    await trafficCardNFT.setDepositContract(await deposit.getAddress());
    await trafficCardNFT.transferOwnership(await deposit.getAddress());

    // user1 注册 + 发币
    await userRegistry.connect(user1).register("u1@linkworld.io");
    await mockUSDT.mint(user1.address, usdt(1000));
  });

  describe("MockUSDT 精度", function () {
    it("DEC-01  decimals() == 6", async function () {
      expect(await mockUSDT.decimals()).to.equal(6);
    });
  });

  describe("分档校验（仅 10/20/50/100 USDT）", function () {
    it("TIER-01  9.999999 USDT → Invalid tier", async function () {
      await mockUSDT.connect(user1).approve(await deposit.getAddress(), usdt("9.999999"));
      await expect(
        deposit.connect(user1).deposit(usdt("9.999999"))
      ).to.be.revertedWith("Invalid tier");
    });

    it("TIER-02  非档位金额（30 USDT）→ Invalid tier", async function () {
      await mockUSDT.connect(user1).approve(await deposit.getAddress(), usdt("30"));
      await expect(
        deposit.connect(user1).deposit(usdt("30"))
      ).to.be.revertedWith("Invalid tier");
    });

    it("TIER-03  四个合法档位均通过：getDepositAmount 累加", async function () {
      await mockUSDT.connect(user1).approve(await deposit.getAddress(), usdt("1000"));
      for (const tier of [10, 20, 50, 100]) {
        await deposit.connect(user1).deposit(usdt(tier));
      }
      expect(await deposit.getDepositAmount(user1.address)).to.equal(usdt(180));
      expect(await deposit.getTrancheCount(user1.address)).to.equal(4n);
    });
  });

  describe("充值即按比例发卡（10/20/50/100 → 1/2/5/10）", function () {
    const cases: [number, number][] = [
      [10, 1],
      [20, 2],
      [50, 5],
      [100, 10],
    ];
    for (const [tier, cards] of cases) {
      it(`CARD-${tier}  存 ${tier} USDT → mint ${cards} 张无限流量卡`, async function () {
        await mockUSDT.connect(user1).approve(await deposit.getAddress(), usdt(tier));
        await deposit.connect(user1).deposit(usdt(tier));
        expect(await trafficCardNFT.getUserCardCount(user1.address)).to.equal(BigInt(cards));
      });
    }

    it("CARD-UNL  每张卡 dataAmount == type(uint256).max（无限流量哨兵）", async function () {
      await mockUSDT.connect(user1).approve(await deposit.getAddress(), usdt(20));
      await deposit.connect(user1).deposit(usdt(20));
      const info0 = await trafficCardNFT.getCardInfo(0);
      const info1 = await trafficCardNFT.getCardInfo(1);
      expect(info0.dataAmount).to.equal(ethers.MaxUint256);
      expect(info1.dataAmount).to.equal(ethers.MaxUint256);
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

  describe("逐笔提取（独立锁 + CEI safeTransfer）", function () {
    it("ERC-03  锁仓未到 revert；到期后退回本金、该笔 withdrawn", async function () {
      const amount = usdt("100");
      await mockUSDT.connect(user1).approve(await deposit.getAddress(), amount);
      await deposit.connect(user1).deposit(amount);

      // 锁仓未到 -> revert
      await expect(deposit.connect(user1).withdraw(0)).to.be.revertedWith("Lock not expired");

      // 快进 30 天 + 1 秒
      await time.increase(LOCK + 1);

      const before = await mockUSDT.balanceOf(user1.address);
      await expect(deposit.connect(user1).withdraw(0))
        .to.emit(deposit, "DepositWithdrawn")
        .withArgs(user1.address, amount, 0);

      expect(await mockUSDT.balanceOf(user1.address)).to.equal(before + amount);
      expect(await deposit.getDepositAmount(user1.address)).to.equal(0n);
      expect(await mockUSDT.balanceOf(await deposit.getAddress())).to.equal(0n);
    });

    it("ERC-03b  非法 index → Invalid tranche", async function () {
      await expect(deposit.connect(user1).withdraw(0)).to.be.revertedWith("Invalid tranche");
    });

    it("ERC-03c  重复取回同一笔 → Already withdrawn", async function () {
      await mockUSDT.connect(user1).approve(await deposit.getAddress(), usdt(10));
      await deposit.connect(user1).deposit(usdt(10));
      await time.increase(LOCK + 1);
      await deposit.connect(user1).withdraw(0);
      await expect(deposit.connect(user1).withdraw(0)).to.be.revertedWith("Already withdrawn");
    });

    it("ERC-03d  逐笔互不影响：先取回早笔，晚笔仍受锁", async function () {
      await mockUSDT.connect(user1).approve(await deposit.getAddress(), usdt("1000"));
      await deposit.connect(user1).deposit(usdt(10)); // tranche 0
      await time.increase(LOCK - 100);
      await deposit.connect(user1).deposit(usdt(20)); // tranche 1（晚 ~30d-100s）
      await time.increase(200); // tranche0 到期，tranche1 未到期

      await deposit.connect(user1).withdraw(0);
      expect(await deposit.getDepositAmount(user1.address)).to.equal(usdt(20));
      // tranche1 仍锁
      await expect(deposit.connect(user1).withdraw(1)).to.be.revertedWith("Lock not expired");
    });
  });

  describe("视图函数", function () {
    it("VIEW-01  getTranches 返回每笔 amount/unlockAt/withdrawn", async function () {
      await mockUSDT.connect(user1).approve(await deposit.getAddress(), usdt("1000"));
      await deposit.connect(user1).deposit(usdt(10));
      const tranches = await deposit.getTranches(user1.address);
      expect(tranches.length).to.equal(1);
      expect(tranches[0].amount).to.equal(usdt(10));
      expect(tranches[0].withdrawn).to.equal(false);
      expect(tranches[0].unlockAt).to.be.greaterThan(0n);
    });

    it("VIEW-02  getLockExpiry 返回未取回笔中最早 unlockAt；全取回 → 0", async function () {
      await mockUSDT.connect(user1).approve(await deposit.getAddress(), usdt("1000"));
      await deposit.connect(user1).deposit(usdt(10)); // 早
      await time.increase(100);
      await deposit.connect(user1).deposit(usdt(10)); // 晚
      const tranches = await deposit.getTranches(user1.address);
      expect(await deposit.getLockExpiry(user1.address)).to.equal(tranches[0].unlockAt);

      // 取回早笔后，最早 = 晚笔
      await time.increase(LOCK + 1);
      await deposit.connect(user1).withdraw(0);
      expect(await deposit.getLockExpiry(user1.address)).to.equal(tranches[1].unlockAt);

      // 全取回 → 0
      await deposit.connect(user1).withdraw(1);
      expect(await deposit.getLockExpiry(user1.address)).to.equal(0n);
    });
  });

  describe("注册校验（REG）", function () {
    it("REG-02  未注册用户 deposit revert Not registered", async function () {
      await mockUSDT.mint(user2.address, usdt(100));
      await mockUSDT.connect(user2).approve(await deposit.getAddress(), usdt(10));
      await expect(
        deposit.connect(user2).deposit(usdt(10))
      ).to.be.revertedWith("Not registered");
    });
  });
});
