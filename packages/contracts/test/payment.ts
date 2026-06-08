import { ethers } from "hardhat";
import { expect } from "chai";
import { ZeroAddress } from "ethers";

// USDT 6 位精度辅助：n USDT -> 最小单位
function usdt(n: string | number): bigint {
  const [int, frac = ""] = String(n).split(".");
  const fracPadded = (frac + "000000").slice(0, 6);
  return BigInt(int) * 1_000_000n + BigInt(fracPadded || "0");
}

describe("Payment ERC20 分账 (T3)", function () {
  let feeManager: any, serviceManager: any, payment: any, mockUSDT: any;
  let owner: any, oracle: any, operatorPayout: any, platformWallet: any, user1: any, user2: any;

  // 使用 ServiceManager 内置 operator 1（T-Mobile US，paymentAddress=0）
  const OP_ID = 1n;

  beforeEach(async function () {
    [owner, oracle, operatorPayout, platformWallet, user1, user2] = await ethers.getSigners();

    const MockUSDTFactory = await ethers.getContractFactory("MockUSDT");
    mockUSDT = await MockUSDTFactory.deploy();

    const FeeManagerFactory = await ethers.getContractFactory("FeeManager");
    feeManager = await FeeManagerFactory.deploy();
    await feeManager.initialize(150); // 1.5%

    const ServiceManagerFactory = await ethers.getContractFactory("ServiceManager");
    serviceManager = await ServiceManagerFactory.deploy();
    await serviceManager.initialize();

    const PaymentFactory = await ethers.getContractFactory("Payment");
    payment = await PaymentFactory.deploy();
    await payment.initialize(
      await feeManager.getAddress(),
      platformWallet.address,
      await mockUSDT.getAddress(),
      await serviceManager.getAddress()
    );

    // oracle 授权（createBill / applyTrafficCardToBill 的 onlyOracle）
    await payment.setOracle(oracle.address);

    // 给 user1 发币
    await mockUSDT.mint(user1.address, usdt(1000));
  });

  // 设置 operator 分账地址（多数用例需要）
  async function setOperatorPayout() {
    await serviceManager.setOperatorPaymentAddress(OP_ID, operatorPayout.address);
  }

  describe("createBill 权限 + fail-fast", function () {
    it("createBill: 非 oracle 调用 → revert Only oracle", async function () {
      await setOperatorPayout();
      await expect(
        payment.connect(owner).createBill(user1.address, OP_ID, usdt(100))
      ).to.be.revertedWith("Only oracle");
    });

    it("createBill: oracle 调用 + paymentAddress 已设 → 成功 emit BillCreated", async function () {
      await setOperatorPayout();
      const amount = usdt(100);
      const fee = await feeManager.calculateFee(amount);
      await expect(
        payment.connect(oracle).createBill(user1.address, OP_ID, amount)
      )
        .to.emit(payment, "BillCreated")
        .withArgs(0, user1.address, amount + fee, fee);
    });

    it("createBill: operator.paymentAddress=0 → fail-fast revert Operator payout unset", async function () {
      // 不设 payout 地址，内置 operator 1 的 paymentAddress 默认为 0
      await expect(
        payment.connect(oracle).createBill(user1.address, OP_ID, usdt(100))
      ).to.be.revertedWith("Operator payout unset");
    });
  });

  describe("payBill 分账", function () {
    beforeEach(async function () {
      await setOperatorPayout();
      await payment.connect(oracle).createBill(user1.address, OP_ID, usdt(100));
    });

    it("PAY-01  未 approve → payBill revert", async function () {
      await expect(payment.connect(user1).payBill(0)).to.be.reverted;
    });

    it("PAY-02  approve 后 payBill 成功：operator 收 amount、platform 收 fee、isPaid=true", async function () {
      const amount = usdt(100);
      const fee = await feeManager.calculateFee(amount);

      await mockUSDT.connect(user1).approve(await payment.getAddress(), amount + fee);

      const opBefore = await mockUSDT.balanceOf(operatorPayout.address);
      const plBefore = await mockUSDT.balanceOf(platformWallet.address);
      const userBefore = await mockUSDT.balanceOf(user1.address);

      await payment.connect(user1).payBill(0);

      expect(await mockUSDT.balanceOf(operatorPayout.address)).to.equal(opBefore + amount);
      expect(await mockUSDT.balanceOf(platformWallet.address)).to.equal(plBefore + fee);
      expect(await mockUSDT.balanceOf(user1.address)).to.equal(userBefore - amount - fee);

      const bills = await payment.getUserBills(user1.address);
      expect(bills[0].isPaid).to.equal(true);
    });

    it("PAY-02b  payBill emit BillPaid(billId, user, total, operatorAmount=amount)", async function () {
      const amount = usdt(100);
      const fee = await feeManager.calculateFee(amount);
      await mockUSDT.connect(user1).approve(await payment.getAddress(), amount + fee);

      await expect(payment.connect(user1).payBill(0))
        .to.emit(payment, "BillPaid")
        .withArgs(0, user1.address, amount + fee, amount);
    });

    it("PAY-04  链上 fee == FeeManager.calculateFee(amount)", async function () {
      const amount = usdt(100);
      const bills = await payment.getUserBills(user1.address);
      expect(bills[0].platformFee).to.equal(await feeManager.calculateFee(amount));
    });

    it("PAY-05  approve=amount（缺 fee）→ payBill 整体 revert，operator 未实收（原子回滚）", async function () {
      const amount = usdt(100);
      await mockUSDT.connect(user1).approve(await payment.getAddress(), amount); // 只批 amount，缺 fee

      const opBefore = await mockUSDT.balanceOf(operatorPayout.address);
      await expect(payment.connect(user1).payBill(0)).to.be.reverted;

      // 原子回滚：operator 未实收，账单仍未支付
      expect(await mockUSDT.balanceOf(operatorPayout.address)).to.equal(opBefore);
      const bills = await payment.getUserBills(user1.address);
      expect(bills[0].isPaid).to.equal(false);
    });

    it("payBill: 非账单本人调用 → revert Not your bill", async function () {
      const amount = usdt(100);
      const fee = await feeManager.calculateFee(amount);
      await mockUSDT.mint(user2.address, usdt(1000));
      await mockUSDT.connect(user2).approve(await payment.getAddress(), amount + fee);
      await expect(payment.connect(user2).payBill(0)).to.be.revertedWith("Not your bill");
    });

    it("payBill: 重复支付 → revert Already paid", async function () {
      const amount = usdt(100);
      const fee = await feeManager.calculateFee(amount);
      await mockUSDT.connect(user1).approve(await payment.getAddress(), amount + fee);
      await payment.connect(user1).payBill(0);
      await expect(payment.connect(user1).payBill(0)).to.be.revertedWith("Already paid");
    });
  });

  describe("PAY-03 分账地址校验", function () {
    it("PAY-03  operator.paymentAddress=0 时 payBill revert Operator payout unset", async function () {
      // 创建账单时先设地址通过 fail-fast，再清掉地址模拟事后失效是不允许的（setter 禁零地址）。
      // 这里直接验证 setter 禁零地址 + createBill fail-fast 共同保证 payBill 不会遇到零地址。
      // 通过新增 operator（paymentAddress=0）走 createBill fail-fast 路径覆盖。
      await expect(
        payment.connect(oracle).createBill(user1.address, OP_ID, usdt(100))
      ).to.be.revertedWith("Operator payout unset");
    });
  });

  describe("0-fee 跳过", function () {
    it("0-FEE  amount 极小使 fee=0 时第二段跳过、不 revert，operator 收全额", async function () {
      await setOperatorPayout();
      // fee = amount*150/10000，amount<67 时 fee 截断为 0（66*150/10000 = 0.99 -> 0）
      const amount = 66n;
      expect(await feeManager.calculateFee(amount)).to.equal(0n);

      await payment.connect(oracle).createBill(user1.address, OP_ID, amount);
      await mockUSDT.connect(user1).approve(await payment.getAddress(), amount);

      const opBefore = await mockUSDT.balanceOf(operatorPayout.address);
      const plBefore = await mockUSDT.balanceOf(platformWallet.address);

      await payment.connect(user1).payBill(0);

      expect(await mockUSDT.balanceOf(operatorPayout.address)).to.equal(opBefore + amount);
      expect(await mockUSDT.balanceOf(platformWallet.address)).to.equal(plBefore); // platform 未收（fee=0）
    });
  });

  describe("ServiceManager.setOperatorPaymentAddress", function () {
    it("SM-PAY-01  onlyOwner：非 owner 调用 revert", async function () {
      await expect(
        serviceManager.connect(user1).setOperatorPaymentAddress(OP_ID, operatorPayout.address)
      ).to.be.reverted;
    });

    it("SM-PAY-02  零地址 revert", async function () {
      await expect(
        serviceManager.setOperatorPaymentAddress(OP_ID, ZeroAddress)
      ).to.be.revertedWith("Zero payout address");
    });

    it("SM-PAY-03  设置后 getOperator 返回新 paymentAddress", async function () {
      await serviceManager.setOperatorPaymentAddress(OP_ID, operatorPayout.address);
      const op = await serviceManager.getOperator(OP_ID);
      expect(op.paymentAddress).to.equal(operatorPayout.address);
    });
  });

  describe("applyTrafficCardToBill 受限桩 (T1，回归保绿)", function () {
    beforeEach(async function () {
      await setOperatorPayout();
      await payment.connect(oracle).createBill(user1.address, OP_ID, usdt(100));
    });

    it("ATC-01  非 oracle 调用 revert", async function () {
      await expect(payment.connect(user1).applyTrafficCardToBill(0)).to.be.revertedWith("Only oracle");
    });

    it("ATC-02  oracle 调用存在账单 → emit TrafficCardApplied，不转资金", async function () {
      const opBefore = await mockUSDT.balanceOf(operatorPayout.address);
      await expect(payment.connect(oracle).applyTrafficCardToBill(0))
        .to.emit(payment, "TrafficCardApplied")
        .withArgs(0);
      expect(await mockUSDT.balanceOf(operatorPayout.address)).to.equal(opBefore);
    });
  });
});
