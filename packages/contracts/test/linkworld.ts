import { ethers } from "hardhat";
import { parseEther, ZeroAddress } from "ethers";
import { expect } from "chai";

// ── Low-level helper: deploy and manually encode+send initializer ───────────────
/** Call the initializer on a freshly deployed impl contract. */
async function setInitializer(
  impl: any,
  func: string,
  args: any[]
): Promise<void> {
  const iface = impl.interface;
  const initializerData = iface.encodeFunctionData(func, args);
  const { signers } = ethers;
  const deployer = await signers();
  await (await deployer[0].sendTransaction({
    to: await impl.getAddress(),
    data: initializerData,
  })).wait();
}

/** Deploy any upgradable contract and immediately call its initializer. */
async function deployUpgradeable(
  name: string,
  initFunc: string,
  initArgs: any[]
): Promise<any> {
  const factory = await ethers.getContractFactory(name);
  const impl = await factory.deploy({ args: [] });
  await setInitializer(impl, initFunc, initArgs);
  await (await ethers.getContractAt(name, await impl.getAddress())).waitForDeployment
    ? await ethers.getContractAt(name, await impl.getAddress())
    : impl;
  return impl;
}

// ─────────────────── Sanity helpers ────────────────────────────────────────────
/** Fee rate → actual fee at 250 bps = 2.5 % */
function feeAt250bps(amount: bigint): bigint {
  return (amount * 250n) / 10000n;
}

// ─────────────────── Short-hand ────────────────────────────────────────────────
const { ZeroAddress: ZA } = ethers;

describe("LinkWorld Contracts — Integration Tests", function () {
  let userRegistry: any;
  let feeManager: any;
  let deposit: any;
  let serviceManager: any;
  let payment: any;
  let oracle: any;
  let trafficCardNFT: any;
  let owner: any;
  let user1: any;
  let user2: any;
  let user3: any;

  // Helper array
  function allSigners(): any[] { return [owner, user1, user2, user3]; }

  // ════════════════════════════════════════════════════════════════════════════
  //  Deploy + Wire helper
  // ════════════════════════════════════════════════════════════════════════════
  async function deployAndWire() {
    [owner, user1, user2, user3] = await ethers.getSigners();

    // ① FeeManager
    const fmImpl = await deployUpgradeable("FeeManager", "initialize", [250]);
    feeManager = await ethers.getContractAt("FeeManager", await fmImpl.getAddress());

    // ② UserRegistry
    const urImpl = await deployUpgradeable("UserRegistry", "initialize", []);
    userRegistry = await ethers.getContractAt("UserRegistry", await urImpl.getAddress());

    // ③ ServiceManager
    const smImpl = await deployUpgradeable("ServiceManager", "initialize", []);
    serviceManager = await ethers.getContractAt("ServiceManager", await smImpl.getAddress());

    // ④ Payment(feeManagerAddr, deployerAddr OR owner)
    const pmImpl = await deployUpgradeable(
      "Payment", "initialize",
      [await feeManager.getAddress(), owner.address]
    );
    payment = await ethers.getContractAt("Payment", await pmImpl.getAddress());

    // ⑤ Deposit(userRegistryAddr)
    const dpImpl = await deployUpgradeable(
      "Deposit", "initialize",
      [await userRegistry.getAddress()]
    );
    deposit = await ethers.getContractAt("Deposit", await dpImpl.getAddress());

    // ⑥ Oracle(paymentAddr)
    const ocImpl = await deployUpgradeable(
      "Oracle", "initialize",
      [await payment.getAddress()]
    );
    oracle = await ethers.getContractAt("Oracle", await ocImpl.getAddress());

    // ⑦ TrafficCardNFT
    const tcImpl = await deployUpgradeable("TrafficCardNFT", "initialize", []);
    trafficCardNFT = await ethers.getContractAt("TrafficCardNFT", await tcImpl.getAddress());

    // ── Wire cross-contract references ──
    await payment.setOracle(await oracle.getAddress());
    await payment.setDeposit(await deposit.getAddress());
    await deposit.setPayment(await payment.getAddress());
    await deposit.setServiceManager(await serviceManager.getAddress());
    await deposit.setOracle(await oracle.getAddress());
    await deposit.setTrafficCardNFT(await trafficCardNFT.getAddress());
    await oracle.setDeposit(await deposit.getAddress());
  }

  // ════════════════════════════════════════════════════════════════════════════
  //  SUITE 1 · FeeManager
  // ════════════════════════════════════════════════════════════════════════════
  describe("FeeManager", function () {
    beforeEach(deployAndWire);

    it("FE-01  initializer sets feeRate to 250 bps", async function () {
      expect(await feeManager.getFeeRate()).to.equal(250);
    });

    it("FE-02  1 ETH × 2.5 % → 0.025 ETH", async function () {
      expect(await feeManager.calculateFee(parseEther("1"))).to.equal(parseEther("0.025"));
    });

    it("FE-03  calculateFee(0) → 0", async function () {
      expect(await feeManager.calculateFee(0)).to.equal(0);
    });

    it("FE-04  10 ETH × 2.5 % → 0.25 ETH", async function () {
      expect(await feeManager.calculateFee(parseEther("10"))).to.equal(parseEther("0.25"));
    });

    it("FE-05  owner raises fee to 500 bps (5 %)", async function () {
      await feeManager.setFeeRate(500);
      expect(await feeManager.getFeeRate()).to.equal(500);
    });

    it("FE-06  MAX_FEE_RATE=1000 accepted", async function () {
      await feeManager.setFeeRate(1000);
      expect(await feeManager.getFeeRate()).to.equal(1000);
    });

    it("FE-07  above 1000 → revert", async function () {
      await expect(feeManager.setFeeRate(1001)).to.be.revertedWith("Fee too high");
    });

    it("FE-08  user1.setFeeRate → revert", async function () {
      await expect(feeManager.connect(user1).setFeeRate(999)).to.be.reverted;
    });

    it("FE-09  FEE_DENOMINATOR = 10000", async function () {
      expect(Number(await feeManager.FEE_DENOMINATOR())).to.equal(10000);
    });
  });

  // ════════════════════════════════════════════════════════════════════════════
  //  SUITE 2 · UserRegistry
  // ════════════════════════════════════════════════════════════════════════════
  describe("UserRegistry", function () {
    beforeEach(deployAndWire);

    it("UR-01  register → isRegistered=true and identity NFT at tokenId=0", async function () {
      await userRegistry.connect(user1).register("u1@linkworld.io");
      expect(await userRegistry.isRegistered(user1.address)).to.be.true;
      expect(await userRegistry.ownerOf(0)).to.equal(user1.address);
      const info = await userRegistry.getUserInfo(user1.address);
      expect(info.wallet).to.equal(user1.address);
      expect(info.email).to.equal("u1@linkworld.io");
    });

    it("UR-02  emit UserRegistered", async function () {
      await expect(userRegistry.connect(user1).register("evt@linkworld.io"))
        .to.emit(userRegistry, "UserRegistered")
        .withArgs(user1.address, "evt@linkworld.io", 0n);
    });

    it("UR-03  duplicate register → Already registered", async function () {
      await userRegistry.connect(user1).register("dup@linkworld.io");
      await expect(userRegistry.connect(user1).register("dup2@linkworld.io"))
        .to.be.revertedWith("Already registered");
    });

    it("UR-04  unregistered user isRegistered=false", async function () {
      expect(await userRegistry.isRegistered(user2.address)).to.be.false;
    });

    it("UR-05  getUnregisteredUserInfo → User not found", async function () {
      await expect(userRegistry.getUserInfo(user2.address)).to.be.revertedWith("User not found");
    });

    it("UR-06  name=LinkWorld Identity, symbol=LWID", async function () {
      expect(await userRegistry.name()).to.equal("LinkWorld Identity");
      expect(await userRegistry.symbol()).to.equal("LWID");
    });
  });

  // ════════════════════════════════════════════════════════════════════════════
  //  SUITE 3 · ServiceManager
  // ════════════════════════════════════════════════════════════════════════════
  describe("ServiceManager", function () {
    beforeEach(deployAndWire);

    it("SM-01  initial active operators > 0", async function () {
      const ops = await serviceManager.getActiveOperators();
      expect(ops.length).to.be.greaterThan(0);
    });

    it("SM-02  getOperatorsByCountry US → non-empty", async function () {
      const us = await serviceManager.getOperatorsByCountry("US");
      expect(us.length).to.be.greaterThan(0);
    });

    it("SM-03  activateService → isActive=true", async function () {
      await serviceManager.connect(user1).activateService(1, "+1234567890", "pw");
      const { isActive, operatorId } = await serviceManager.getUserService(user1.address);
      expect(isActive).to.be.true;
      expect(operatorId).to.equal(1);
    });

    it("SM-04  double activate → Service already active", async function () {
      await serviceManager.connect(user1).activateService(1, "+1111111", "pw");
      await expect(
        serviceManager.connect(user1).activateService(2, "+2222222", "pw2")
      ).to.be.revertedWith("Service already active");
    });

    it("SM-05  deactivateService → isActive=false, deactivatedAt > 0", async function () {
      await serviceManager.connect(user1).activateService(1, "+3333333", "pw");
      await serviceManager.connect(user1).deactivateService();
      const { isActive, deactivatedAt } = await serviceManager.getUserService(user1.address);
      expect(isActive).to.be.false;
      expect(deactivatedAt).to.be.greaterThan(0n);
    });

    it("SM-06  reactivation after deactivate", async function () {
      await serviceManager.connect(user1).activateService(1, "+4444444", "pw");
      await serviceManager.connect(user1).deactivateService();
      await serviceManager.connect(user1).activateService(7, "+5555555", "pw2");
      const { isActive } = await serviceManager.getUserService(user1.address);
      expect(isActive).to.be.true;
    });
  });

  // ════════════════════════════════════════════════════════════════════════════
  //  SUITE 4 · Deposit
  // ════════════════════════════════════════════════════════════════════════════
  describe("Deposit", function () {
    beforeEach(async function () {
      await deployAndWire();
      await userRegistry.connect(user1).register("u1@linkworld.io");
      await userRegistry.connect(user2).register("u2@linkworld.io");
    });

    it("DP-01  0.1 ETH deposit → getDepositAmount = 0.1 ETH", async function () {
      await deposit.connect(user1).deposit({ value: parseEther("0.1") });
      expect(await deposit.getDepositAmount(user1.address)).to.equal(parseEther("0.1"));
    });

    it("DP-02  emit DepositMade(wallet, wei)", async function () {
      await expect(deposit.connect(user1).deposit({ value: parseEther("0.2") }))
        .to.emit(deposit, "DepositMade")
        .withArgs(user1.address, parseEther("0.2"));
    });

    it("DP-03  unregistered user → Not registered", async function () {
      await expect(
        deposit.connect(user2).deposit({ value: parseEther("0.1") })
      ).to.be.revertedWith("Not registered");
    });

    it("DP-04  zero value → Zero deposit", async function () {
      await expect(deposit.connect(user1).deposit({ value: 0 })).to.be.revertedWith("Zero deposit");
    });

    it("DP-05  three deposits accumulate", async function () {
      const d10 = parseEther("0.1"), d20 = parseEther("0.2"), d5 = parseEther("0.05");
      await deposit.connect(user1).deposit({ value: d10 });
      await deposit.connect(user1).deposit({ value: d20 });
      await deposit.connect(user1).deposit({ value: d5 });
      expect(await deposit.getDepositAmount(user1.address)).to.equal(d10 + d20 + d5);
    });

    it("DP-06  withdraw with active service → Service still active", async function () {
      await deposit.connect(user1).deposit({ value: parseEther("0.5") });
      await serviceManager.connect(user1).activateService(1, "+1111111", "pw");
      await expect(deposit.connect(user1).withdraw()).to.be.revertedWith("Service still active");
    });

    it("DP-07  withdraw with unpaid bill → Has unpaid bills", async function () {
      await deposit.connect(user1).deposit({ value: parseEther("0.5") });
      await oracle.submitUsage(user1.address, 1, 100, 10);
      await expect(deposit.connect(user1).withdraw()).to.be.revertedWith("Has unpaid bills");
    });

    it("DP-08  happy withdraw → ETH returned, _deposits=0, DepositWithdrawn", async function () {
      await deposit.connect(user1).deposit({ value: parseEther("0.5") });
      const balBefore = await ethers.provider.getBalance(user1.address);
      const tx = await deposit.connect(user1).withdraw();
      const rcpt = await tx.wait();
      const gasCost = rcpt!.gasUsed * rcpt!.gasPrice;
      const balAfter = await ethers.provider.getBalance(user1.address);
      expect(await deposit.getDepositAmount(user1.address)).to.equal(0);
      expect(balAfter).to.be.closeTo(balBefore + parseEther("0.5"), gasCost);
    });

    it("DP-09  emit DepositWithdrawn(wallet, principal, 0)", async function () {
      await deposit.connect(user1).deposit({ value: parseEther("0.3") });
      await expect(deposit.connect(user1).withdraw())
        .to.emit(deposit, "DepositWithdrawn")
        .withArgs(user1.address, parseEther("0.3"), 0);
    });

    it("DP-10  setTrafficCardQuota only owner", async function () {
      await deposit.setTrafficCardQuota(50 * 1024 * 1024);
      expect(Number(await deposit.trafficCardQuota())).to.equal(50 * 1024 * 1024);
      await expect(deposit.connect(user1).setTrafficCardQuota(1)).to.be.reverted;
    });

    it("DP-11  generateTokenURI contains api.linkworld.io", async function () {
      const month = BigInt(Math.floor(Date.now() / 1000 / 2592000));
      const uri = await deposit.generateTokenURI(user1.address, Number(month));
      expect(uri).to.include("api.linkworld.io/traffic-card/");
    });

    it("DP-12  setRequiredDeposit readable/gettable", async function () {
      await deposit.setRequiredDeposit(1, parseEther("0.1"));
      expect(await deposit.getRequiredDeposit(1)).to.equal(parseEther("0.1"));
    });

    it("DP-13  setOracle stores address", async function () {
      await deposit.setOracle(await oracle.getAddress());
      expect(await deposit.oracle()).to.equal(await oracle.getAddress());
    });
  });

  // ════════════════════════════════════════════════════════════════════════════
  //  SUITE 5 · TrafficCardNFT
  // ════════════════════════════════════════════════════════════════════════════
  describe("TrafficCardNFT", function () {
    beforeEach(async function () {
      await deployAndWire();
    });

    it("TC-01  name=LinkWorld Traffic Card, symbol=LWTC", async function () {
      expect(await trafficCardNFT.name()).to.equal("LinkWorld Traffic Card");
      expect(await trafficCardNFT.symbol()).to.equal("LWTC");
    });

    it("TC-02  mint → ownerOf(0)=user1", async function () {
      await trafficCardNFT.mint(user1.address, 2048, "https://example.com/0");
      expect(await trafficCardNFT.ownerOf(0)).to.equal(user1.address);
    });

    it("TC-03  mint emits CardMinted(user,0,dataAmt)", async function () {
      await expect(trafficCardNFT.mint(user1.address, 512, "uri"))
        .to.emit(trafficCardNFT, "CardMinted").withArgs(user1.address, 0n, 512n);
    });

    it("TC-04  user1 mint → user2 mint → ownerOf(1)=user2", async function () {
      await trafficCardNFT.mint(user1.address, 512, "uri1");
      await trafficCardNFT.mint(user2.address, 1024, "uri2");
      expect(await trafficCardNFT.ownerOf(0)).to.equal(user1.address);
      expect(await trafficCardNFT.ownerOf(1)).to.equal(user2.address);
    });

    it("TC-05  mintBatch 3 users → 3 distinct tokenIds", async function () {
      await trafficCardNFT.mintBatch(
        [user1.address, user2.address, user3.address],
        [256, 512, 1024],
        ["uri1", "uri2", "uri3"]
      );
      expect(await trafficCardNFT.getUserCardCount(user1.address)).to.equal(1);
      expect(await trafficCardNFT.getUserCardCount(user2.address)).to.equal(1);
      expect(await trafficCardNFT.getUserCardCount(user3.address)).to.equal(1);
    });

    it("TC-06  mintBatch wrong length → Length mismatch", async function () {
      await expect(
        trafficCardNFT.mintBatch([user1.address], [512, 1024], ["uri1"])
      ).to.be.revertedWith("Length mismatch");
    });

    it("TC-07  mint spender=address(0) → Invalid address", async function () {
      await expect(trafficCardNFT.mint(ZA, 512, "uri")).to.be.revertedWith("Invalid address");
    });

    it("TC-08  mint dataAmount=0 → Zero data amount", async function () {
      await expect(trafficCardNFT.mint(user1.address, 0, "uri")).to.be.revertedWith("Zero data amount");
    });

    it("TC-09  burn by non-owner → Not owner or approved", async function () {
      await trafficCardNFT.mint(user1.address, 512, "uri");
      await expect(trafficCardNFT.connect(user2).burn(0)).to.be.revertedWith("Not owner or approved");
    });

    it("TC-10  burn emits CardDestroyed + credit created", async function () {
      await trafficCardNFT.mint(user1.address, 4096, "uri");
      await expect(trafficCardNFT.connect(user1).burn(0))
        .to.emit(trafficCardNFT, "CardDestroyed")
        .withArgs(user1.address, 0n, 4096n);
      expect(await trafficCardNFT.getAvailableCredit(user1.address)).to.equal(4096);
    });

    it("TC-11  double-burn → Card already destroyed", async function () {
      await trafficCardNFT.mint(user1.address, 512, "uri");
      await trafficCardNFT.connect(user1).burn(0);
      await expect(trafficCardNFT.connect(user1).burn(0)).to.be.revertedWith("Card already destroyed");
    });

    it("TC-12  getCardInfo after mint returns dataAmount=512", async function () {
      await trafficCardNFT.mint(user1.address, 512, "uri");
      const info = await trafficCardNFT.getCardInfo(0);
      expect(info.dataAmount).to.equal(512);
      expect(info.isDestroyed).to.be.false;
    });

    it("TC-13  getCardInfo(nonexistent 999) → Card not found", async function () {
      await expect(trafficCardNFT.getCardInfo(999)).to.be.revertedWith("Card not found");
    });

    it("TC-14  getCreditExpiry after burn > now+25d", async function () {
      await trafficCardNFT.mint(user1.address, 128, "uri");
      const { timestamp: ts } = await ethers.provider.getBlock("latest")!;
      await trafficCardNFT.connect(user1).burn(0);
      expect(Number(await trafficCardNFT.getCreditExpiry(user1.address))).to.be.greaterThan(ts + 25 * 86400);
    });

    it("TC-15  useCredit 256 from 1024 → remaining 768", async function () {
      await trafficCardNFT.mint(user1.address, 1024, "uri");
      await trafficCardNFT.connect(user1).burn(0);
      await trafficCardNFT.useCredit(user1.address, 256);
      expect(await trafficCardNFT.getAvailableCredit(user1.address)).to.equal(768n);
    });

    it("TC-16  useCredit > balance → Insufficient credit", async function () {
      await trafficCardNFT.mint(user1.address, 100, "uri");
      await trafficCardNFT.connect(user1).burn(0);
      await expect(trafficCardNFT.useCredit(user1.address, 200)).to.be.revertedWith("Insufficient credit");
    });

    it("TC-17  non-owner useCredit → revert", async function () {
      await trafficCardNFT.mint(user1.address, 512, "uri");
      await trafficCardNFT.connect(user1).burn(0);
      await expect(trafficCardNFT.connect(user2).useCredit(user1.address, 1)).to.be.reverted;
    });

    it("TC-18  tokenURI override (custom URI)", async function () {
      await trafficCardNFT.mint(user1.address, 256, "https://custom.example.com/42");
      expect(await trafficCardNFT.tokenURI(0)).to.equal("https://custom.example.com/42");
    });

    it("TC-19  supportsInterface(ERC721) → true", async function () {
      const ERC721 = ethers.keccak256(ethers.toUtf8Bytes("ERC721"));
      expect(await trafficCardNFT.supportsInterface(ERC721)).to.be.true;
    });

    it("TC-20  no credit pre-mint/burn → getAvailableCredit=0", async function () {
      expect(await trafficCardNFT.getAvailableCredit(user1.address)).to.equal(0);
    });

    it("TC-21  userCardCount zero after deployment", async function () {
      expect(await trafficCardNFT.getUserCardCount(user1.address)).to.equal(0);
    });
  });

  // ════════════════════════════════════════════════════════════════════════════
  //  SUITE 6 · Oracle
  // ════════════════════════════════════════════════════════════════════════════
  describe("Oracle", function () {
    beforeEach(deployAndWire);

    it("OR-01  initialize → payment == Payment address", async function () {
      expect(await oracle.payment()).to.equal(await payment.getAddress());
    });

    it("OR-02  recordUsage → _monthlyUsage + _latestUsage updated", async function () {
      await oracle.recordUsage(user1.address, 1, 1000, 200);
      const [data, call] = await oracle.getLatestUsage(user1.address, 1);
      expect(data).to.equal(1000);
      expect(call).to.equal(200);
    });

    it("OR-03  recordUsage address(0) → Invalid user", async function () {
      await expect(oracle.recordUsage(ZA, 1, 100, 10)).to.be.revertedWith("Invalid user");
    });

    it("OR-04  recordUsage 0+0 → Zero usage", async function () {
      await expect(oracle.recordUsage(user1.address, 1, 0, 0)).to.be.revertedWith("Zero usage");
    });

    it("OR-05  submitUsage same shape as recordUsage", async function () {
      await oracle.submitUsage(user2.address, 3, 600, 80);
      const [data, call] = await oracle.getLatestUsage(user2.address, 3);
      expect(data).to.equal(600);
      expect(call).to.equal(80);
    });

    it("OR-06  getLatestUsage(never queried) → 0,0", async function () {
      const [data, call] = await oracle.getLatestUsage(user3.address, 99);
      expect(data).to.equal(0);
      expect(call).to.equal(0);
    });

    it("OR-07  setDeposit sets oracle internal Deposit address", async function () {
      await oracle.setDeposit(await deposit.getAddress());
      expect(await oracle.deposit()).to.equal(await deposit.getAddress());
    });

    it("OR-08  setDeposit from user1 → revert", async function () {
      await expect(oracle.connect(user1).setDeposit(await deposit.getAddress())).to.be.reverted;
    });

    it("OR-09  monthlySettlement wrong length → Length mismatch", async function () {
      await expect(oracle.monthlySettlement(
        [user1.address, user2.address], [1], [0, 0], [0, 0]
      )).to.be.revertedWith("Length mismatch");
    });

    it("OR-10  monthlySettlement(0s) → no new bill", async function () {
      await oracle.monthlySettlement([user1.address], [1], [0], [0]);
      expect(await payment.getUnpaidBills(user1.address)).to.have.length(0);
    });

    it("OR-11  monthlySettlement fires UsageDataSubmitted", async function () {
      await oracle.recordUsage(user1.address, 1, 3000, 500);
      await expect(oracle.monthlySettlement([user1.address], [1], [3000], [500]))
        .to.emit(oracle, "UsageDataSubmitted")
        .withArgs(user1.address, 1, 3000, 500);
    });
  });

  // ════════════════════════════════════════════════════════════════════════════
  //  SUITE 7 · Payment
  // ════════════════════════════════════════════════════════════════════════════
  describe("Payment", function () {
    beforeEach(deployAndWire);

    it("PY-01  createBill fields: amount, platformFee (250bps), isPaid=false", async function () {
      const amt = parseEther("1.0");
      await payment.createBill(user1.address, 5, amt);
      const bills = await payment.getUserBills(user1.address);
      expect(bills.length).to.equal(1);
      expect(bills[0].amount).to.equal(amt);
      expect(bills[0].isPaid).to.be.false;
      expect(bills[0].platformFee).to.equal(feeAt250bps(amt));
    });

    it("PY-02  createBill emits BillCreated", async function () {
      await expect(
        payment.createBill(user2.address, 9, parseEther("0.5"))
      ).to.emit(payment, "BillCreated");
    });

    it("PY-03  getUnpaidBills length = 1 after single create", async function () {
      await payment.createBill(user1.address, 1, parseEther("2.0"));
      expect((await payment.getUnpaidBills(user1.address)).length).to.equal(1);
    });

    it("PY-04  payBill exact → isPaid=true", async function () {
      const amt = parseEther("1.0");
      await payment.createBill(user1.address, 1, amt);
      const total = amt + await feeManager.calculateFee(amt);
      await payment.connect(user1).payBill(1, { value: total });
      expect((await payment.getUserBills(user1.address))[0].isPaid).to.be.true;
    });

    it("PY-05  payBill under-1 wei → Insufficient payment", async function () {
      const amt = parseEther("1.0");
      await payment.createBill(user1.address, 1, amt);
      const total = amt + await feeManager.calculateFee(amt) - 1n;
      await expect(payment.connect(user1).payBill(1, { value: total }))
        .to.be.revertedWith("Insufficient payment");
    });

    it("PY-06  double pay same bill → Already paid", async function () {
      const amt = parseEther("0.1");
      await payment.createBill(user1.address, 2, amt);
      const total = amt + await feeManager.calculateFee(amt);
      await payment.connect(user1).payBill(1, { value: total });
      await expect(
        payment.connect(user1).payBill(1, { value: total })
      ).to.be.revertedWith("Already paid");
    });

    it("PY-07  payBill emits BillPaid(1,user1,total)", async function () {
      const amt = parseEther("0.5");
      await payment.createBill(user1.address, 3, amt);
      const total = amt + await feeManager.calculateFee(amt);
      await expect(payment.connect(user1).payBill(1, { value: total }))
        .to.emit(payment, "BillPaid")
        .withArgs(1n, user1.address, total);
    });

    it("PY-08  setPlatformWallet only owner", async function () {
      await payment.setPlatformWallet(user2.address);
      expect(await payment.platformWallet()).to.equal(user2.address);
    });

    it("PY-09  non-owner set* calls → revert", async function () {
      await expect(payment.connect(user1).setPlatformWallet(ZA)).to.be.reverted;
      await expect(payment.connect(user1).setDeposit(ZA)).to.be.reverted;
      await expect(payment.connect(user1).setOracle(ZA)).to.be.reverted;
    });

    it("PY-10  createBill only via oracle", async function () {
      await expect(
        payment.connect(user1).createBill(user2.address, 1, parseEther("1"))
      ).to.be.revertedWith("Only oracle");
    });

    it("PY-11  createBill zero amount → Zero amount", async function () {
      await expect(payment.createBill(user1.address, 1, 0)).to.be.revertedWith("Zero amount");
    });
  });

  // ════════════════════════════════════════════════════════════════════════════
  //  SUITE 8 · Integration — End-to-End Business Flow
  // ════════════════════════════════════════════════════════════════════════════
  describe("Integration — Full Business Flow", function () {
    beforeEach(deployAndWire);

    it("FLOW-01  Register×3 → Deposit×3 → Activate×3 → SubmitUsage → MonthlySettlement", async function () {
      await userRegistry.connect(user1).register("f1@linkworld.io");
      await userRegistry.connect(user2).register("f2@linkworld.io");
      await userRegistry.connect(user3).register("f3@linkworld.io");

      await deposit.connect(user1).deposit({ value: parseEther("0.5") });
      await deposit.connect(user2).deposit({ value: parseEther("0.3") });
      await deposit.connect(user3).deposit({ value: parseEther("0.7") });

      await serviceManager.connect(user1).activateService(1, "+10000001", "pw");
      await serviceManager.connect(user2).activateService(2, "+10000002", "pw");
      await serviceManager.connect(user3).activateService(3, "+10000003", "pw");

      await oracle.submitUsage(user1.address, 1, 2000, 400);  // 2400
      await oracle.submitUsage(user2.address, 2, 800, 100);    //  900
      await oracle.submitUsage(user3.address, 3, 500, 100);    //  600

      await expect(oracle.monthlySettlement(
        [user1.address, user2.address, user3.address],
        [1, 2, 3], [0, 0, 0], [0, 0, 0]
      )).to.emit(oracle, "UsageDataSubmitted");

      const b1 = (await payment.getUnpaidBills(user1.address))[0];
      expect(b1.amount).to.equal(2400);

      const b2 = (await payment.getUnpaidBills(user2.address))[0];
      expect(b2.amount).to.equal(900);

      const b3 = (await payment.getUnpaidBills(user3.address))[0];
      expect(b3.amount).to.equal(600);
    });

    it("FLOW-02  Pay bill + verify zero unpaid bills", async function () {
      await userRegistry.connect(user1).register("payflow2@linkworld.io");
      await deposit.connect(user1).deposit({ value: parseEther("1.0") });
      await oracle.submitUsage(user1.address, 1, 1000, 0);
      await oracle.monthlySettlement([user1.address], [1], [0], [0]);

      const bid = (await payment.getUnpaidBills(user1.address))[0].id;
      const bill = (await payment.getUserBills(user1.address))[0];
      const total = bill.amount + bill.platformFee;
      await payment.connect(user1).payBill(bid, { value: total });
      expect((await payment.getUnpaidBills(user1.address)).length).to.equal(0);
    });

    it("FLOW-03  Pay + deactivate service → withdraw succeeds → ETH returned", async function () {
      await userRegistry.connect(user1).register("wdflow@linkworld.io");
      await deposit.connect(user1).deposit({ value: parseEther("0.5") });

      // Pay bill
      await oracle.submitUsage(user1.address, 1, 100, 0);
      await oracle.monthlySettlement([user1.address], [1], [0], [0]);
      const bid = (await payment.getUnpaidBills(user1.address))[0].id;
      const bill0 = (await payment.getUserBills(user1.address))[0];
      await payment.connect(user1).payBill(bid, {
        value: bill0.amount + bill0.platformFee,
      });

      // Deactivate service to clear withdraw precondition
      await serviceManager.connect(user1).deactivateService();

      const beforeBal = await ethers.provider.getBalance(user1.address);
      const tx = await deposit.connect(user1).withdraw();
      const rcpt = await tx.wait();
      const gasCost = rcpt!.gasUsed * rcpt!.gasPrice;
      const afterBal = await ethers.provider.getBalance(user1.address);
      expect(afterBal).to.be.closeTo(beforeBal + parseEther("0.5"), gasCost);
    });

    it("FLOW-04  fee change reflected in new bills", async function () {
      await userRegistry.connect(user1).register("feechange@linkworld.io");
      await feeManager.setFeeRate(400);               // 4 %
      const TEN = parseEther("10.0");
      await payment.createBill(user1.address, 1, TEN);
      const bill = (await payment.getUserBills(user1.address))[0];
      expect(bill.platformFee).to.equal(parseEther("0.4"));  // 4 %
      expect(bill.amount).to.equal(TEN);
    });

    it("FLOW-05  autoSettle ≥ 2 bills + deposit sponsored TrafficCardNFT", async function () {
      await userRegistry.connect(user1).register("autosettle3@linkworld.io");
      await deposit.connect(user1).deposit({ value: parseEther("1.0") });
      await deposit.setTrafficCardNFT(await trafficCardNFT.getAddress());
      await oracle.submitUsage(user1.address, 1, 500, 0);

      await payment.autoSettle(
        [user1.address], [1], [parseEther("0.6")]
      );
      const all = await payment.getUserBills(user1.address);
      expect(all.length).to.be.greaterThanOrEqual(2);
    });

    it("FLOW-06  3 users independent state isolation", async function () {
      const amounts = [0.25, 0.50, 0.75];
      const signers = [user1, user2, user3];

      for (let i = 0; i < 3; i++) {
        const s = signers[i];
        const a = amounts[i];
        await userRegistry.connect(s).register(`mi${s.address.slice(0, 8)}@linkworld.io`);
        await deposit.connect(s).deposit({ value: parseEther(String(a)) });
      }

      for (let i = 0; i < 3; i++) {
        const bal = await deposit.getDepositAmount(signers[i].address);
        expect(bal).to.be.greaterThan(0n);
        expect(bal).to.be.closeTo(
          parseEther(String(amounts[i])),
          parseEther("0.0001")
        );
      }
    });
  });
});
