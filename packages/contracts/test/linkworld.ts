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

  // 充值即发卡：NFT owner 转给 Deposit，deposit 内 trafficCardNFT.mint(onlyOwner) 才能调。
  // 仅在需要"存款触发发卡"的 suite 调用（TrafficCardNFT 直 mint 测试仍需 owner=测试账户）。
  async function enableDepositMint() {
    await trafficCardNFT.setDepositContract(await deposit.getAddress());
    await trafficCardNFT.transferOwnership(await deposit.getAddress());
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
      await enableDepositMint();
      await userRegistry.connect(user1).register("u1@linkworld.io");
      await mockUSDT.mint(user1.address, usdt(1000));
    });

    it("DP-01  deposit(50) and getDepositAmount", async function () {
      const amount = usdt(50);
      await mockUSDT.connect(user1).approve(await deposit.getAddress(), amount);
      await deposit.connect(user1).deposit(amount);
      expect(await deposit.getDepositAmount(user1.address)).to.equal(amount);
      // 50 USDT → 5 张卡
      expect(await trafficCardNFT.getUserCardCount(user1.address)).to.equal(5n);
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

    // 每卡 = 1 天，销毁兑换 SIM 天数（= 卡数）
    async function mintCards(to, n) {
      for (let i = 0; i < n; i++) {
        await trafficCardNFT.mint(to.address, ethers.MaxUint256, `uri-${i}`);
      }
    }

    it("RDM-01  redeemForSim 批量销毁 3 张：isDestroyed=true、count 归零、emit SimRedeemed(user,3,[0,1,2])", async function () {
      await mintCards(user1, 3);
      expect(await trafficCardNFT.getUserCardCount(user1.address)).to.equal(3n);

      await expect(trafficCardNFT.connect(user1).redeemForSim([0, 1, 2]))
        .to.emit(trafficCardNFT, "SimRedeemed")
        .withArgs(user1.address, 3n, [0n, 1n, 2n]);

      expect(await trafficCardNFT.getUserCardCount(user1.address)).to.equal(0n);
      for (let i = 0; i < 3; i++) {
        expect((await trafficCardNFT.getCardInfo(i)).isDestroyed).to.be.true;
      }
    });

    it("RDM-02  redeemForSim 部分销毁：销毁 2 张，count 由 5 减到 3", async function () {
      await mintCards(user1, 5);
      const tx = await trafficCardNFT.connect(user1).redeemForSim([1, 3]);
      const receipt = await tx.wait();
      expect(await trafficCardNFT.getUserCardCount(user1.address)).to.equal(3n);
      expect((await trafficCardNFT.getCardInfo(1)).isDestroyed).to.be.true;
      expect((await trafficCardNFT.getCardInfo(3)).isDestroyed).to.be.true;
      expect((await trafficCardNFT.getCardInfo(0)).isDestroyed).to.be.false;
    });

    it("RDM-03  redeemForSim 返回 daysCount == 卡数", async function () {
      await mintCards(user1, 4);
      const daysCount = await trafficCardNFT.connect(user1).redeemForSim.staticCall([0, 1, 2, 3]);
      expect(daysCount).to.equal(4n);
    });

    it("RDM-04  非持有者 redeemForSim → revert Not owner or approved", async function () {
      await mintCards(user1, 2);
      await expect(
        trafficCardNFT.connect(user2).redeemForSim([0, 1])
      ).to.be.revertedWith("Not owner or approved");
    });

    it("RDM-05  已销毁卡再次 redeemForSim → revert（token 已 _burn，ownerOf revert ERC721NonexistentToken）", async function () {
      await mintCards(user1, 2);
      await trafficCardNFT.connect(user1).redeemForSim([0]);
      await expect(
        trafficCardNFT.connect(user1).redeemForSim([0])
      ).to.be.revertedWithCustomError(trafficCardNFT, "ERC721NonexistentToken");
    });

    it("RDM-06  空数组 redeemForSim → revert No cards", async function () {
      await expect(
        trafficCardNFT.connect(user1).redeemForSim([])
      ).to.be.revertedWith("No cards");
    });

    it("RDM-07  redeemForSim 含重复 tokenId → 第二次命中已销毁 token，整笔 revert（防双花）", async function () {
      await mintCards(user1, 2);
      await expect(
        trafficCardNFT.connect(user1).redeemForSim([0, 0])
      ).to.be.revertedWithCustomError(trafficCardNFT, "ERC721NonexistentToken");
      // 整笔 revert：第 0 张不会被实际销毁，count 仍为 2
      expect(await trafficCardNFT.getUserCardCount(user1.address)).to.equal(2n);
      expect((await trafficCardNFT.getCardInfo(0)).isDestroyed).to.be.false;
    });

    it("BURN-01  单张 burn 走 SimRedeemed(user,1,[tokenId]) + CardDestroyed，无 ServiceActivated", async function () {
      await mintCards(user1, 1);
      await expect(trafficCardNFT.connect(user1).burn(0))
        .to.emit(trafficCardNFT, "SimRedeemed").withArgs(user1.address, 1n, [0n])
        .and.to.emit(trafficCardNFT, "CardDestroyed");
      expect(await trafficCardNFT.getUserCardCount(user1.address)).to.equal(0n);
      expect((await trafficCardNFT.getCardInfo(0)).isDestroyed).to.be.true;
    });

    it("BURN-02  非持有者 burn → revert", async function () {
      await mintCards(user1, 1);
      await expect(trafficCardNFT.connect(user2).burn(0)).to.be.revertedWith("Not owner or approved");
    });

    it("BURN-03  ServiceActivated 事件已移除（合约 ABI 不含该 fragment）", async function () {
      const hasServiceActivated = trafficCardNFT.interface.fragments.some(
        (f) => f.type === "event" && (f as any).name === "ServiceActivated"
      );
      const hasSimRedeemed = trafficCardNFT.interface.fragments.some(
        (f) => f.type === "event" && (f as any).name === "SimRedeemed"
      );
      expect(hasServiceActivated).to.be.false;
      expect(hasSimRedeemed).to.be.true;
    });
  });

  // ===== 充值即发卡（mint-on-deposit）+ 月末结算（amounts[] 计价，不再月度发卡）=====
  describe("充值即发卡 + 月末结算", function () {
    const LOCK = 30 * 24 * 60 * 60; // 30 days

    // 充值即发卡：NFT owner 转给 Deposit；operator 1 设分账地址供 createBill fail-fast
    async function setupMintOnDeposit() {
      await deployAndWire();
      await enableDepositMint();
      await serviceManager.setOperatorPaymentAddress(1, user2.address);
    }

    // 让 user 注册 + 分档存款（每次 push 一笔 30 天独立锁 + 即发卡）
    async function registerAndDeposit(signer, email, amount) {
      await userRegistry.connect(signer).register(email);
      await mockUSDT.mint(signer.address, amount);
      await mockUSDT.connect(signer).approve(await deposit.getAddress(), amount);
      await deposit.connect(signer).deposit(amount);
    }

    beforeEach(setupMintOnDeposit);

    it("ISS-01  存 50 USDT → 即发 5 张无限流量卡", async function () {
      await registerAndDeposit(user1, "iss01@linkworld.io", usdt(50));
      expect(await trafficCardNFT.getUserCardCount(user1.address)).to.equal(5n);
      const info = await trafficCardNFT.getCardInfo(0);
      expect(info.dataAmount).to.equal(ethers.MaxUint256); // 无限流量哨兵
    });

    it("ISS-02  多笔存款张数累加：10 + 100 → 1 + 10 = 11 张", async function () {
      await userRegistry.connect(user1).register("iss02@linkworld.io");
      await mockUSDT.mint(user1.address, usdt(1000));
      await mockUSDT.connect(user1).approve(await deposit.getAddress(), usdt(1000));
      await deposit.connect(user1).deposit(usdt(10));
      await deposit.connect(user1).deposit(usdt(100));
      expect(await trafficCardNFT.getUserCardCount(user1.address)).to.equal(11n);
      expect(await deposit.getTrancheCount(user1.address)).to.equal(2n);
    });

    it("ISS-03  非档位金额（30 USDT）→ Invalid tier，不发卡", async function () {
      await userRegistry.connect(user1).register("iss03@linkworld.io");
      await mockUSDT.mint(user1.address, usdt(30));
      await mockUSDT.connect(user1).approve(await deposit.getAddress(), usdt(30));
      await expect(deposit.connect(user1).deposit(usdt(30))).to.be.revertedWith("Invalid tier");
      expect(await trafficCardNFT.getUserCardCount(user1.address)).to.equal(0n);
    });

    it("MS-01  monthlySettlement 用 amounts[] → createBill amount==amounts[i]（不再 usage 求和）", async function () {
      await userRegistry.connect(user1).register("ms01@linkworld.io");
      // oracle.monthlySettlement 是 onlyOwner（owner=deployer）；Oracle 经 setPayment 调 createBill
      await oracle.setPayment(await payment.getAddress());
      const amount = usdt(100);
      await oracle.monthlySettlement([user1.address], [1], [amount]);
      const bills = await payment.getUserBills(user1.address);
      expect(bills.length).to.equal(1);
      expect(bills[0].amount).to.equal(amount); // == amounts[i]，绝非 dataUsage+callUsage
    });

    it("MS-02  非 owner 调 monthlySettlement → revert", async function () {
      await oracle.setPayment(await payment.getAddress());
      await expect(
        oracle.connect(user1).monthlySettlement([user1.address], [1], [usdt(100)])
      ).to.be.reverted;
    });

    it("MS-03  端到端：充值即发卡 + monthlySettlement(createBill + applyTrafficCardToBill) 链路通", async function () {
      // 用户存款时已得卡；月末结算只建账单 + 抵扣（不再月度发卡）
      await registerAndDeposit(user1, "ms03@linkworld.io", usdt(50));
      expect(await trafficCardNFT.getUserCardCount(user1.address)).to.equal(5n);

      await oracle.setPayment(await payment.getAddress());
      const amount = usdt(50);
      await expect(
        oracle.monthlySettlement([user1.address], [1], [amount])
      ).to.not.be.reverted;

      // 账单已建
      const bills = await payment.getUserBills(user1.address);
      expect(bills.length).to.equal(1);
      expect(bills[0].amount).to.equal(amount);
    });

    it("MS-04  monthlySettlement 长度不一致 → revert", async function () {
      await oracle.setPayment(await payment.getAddress());
      await expect(
        oracle.monthlySettlement([user1.address], [1, 2], [usdt(100)])
      ).to.be.revertedWith("Length mismatch");
    });
  });
});


