// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import "./interfaces/IPayment.sol";
import "./interfaces/IFeeManager.sol";
import "./interfaces/IDeposit.sol";

/// @title Payment - 支付结算（服务商费用 + 平台手续费 + 流量卡抵扣，可升级版）
contract Payment is IPayment, OwnableUpgradeable, ReentrancyGuard, UUPSUpgradeable {
    IFeeManager public feeManager;
    IDeposit public deposit;

    address public platformWallet;
    address public oracle; // 预言机地址，有权创建账单

    uint256 private _nextBillId;
    mapping(uint256 => Bill) private _bills;
    mapping(address => uint256[]) private _userBillIds;

    modifier onlyOracle() {
        require(msg.sender == oracle, "Only oracle");
        _;
    }

    /// @inheritdoc IPayment
    function initialize(
        address _feeManager,
        address _platformWallet
    ) public initializer {
        __Ownable_init(msg.sender);
        _reentrancyGuardInit();

        feeManager = IFeeManager(_feeManager);
        platformWallet = _platformWallet;
    }

    // 内部 ReentrancyGuard 初始化
    function _reentrancyGuardInit() internal {
        _reentrancyGuardStorageSlot().getUint256Slot().value = 1; // NOT_ENTERED = 1
    }

    function setOracle(address _oracle) external onlyOwner {
        oracle = _oracle;
    }

    function setPlatformWallet(address _platformWallet) external onlyOwner {
        platformWallet = _platformWallet;
    }

    function setDeposit(address _deposit) external onlyOwner {
        deposit = IDeposit(_deposit);
    }

    function setDeposit(address _deposit) external onlyOwner {
        deposit = IDeposit(_deposit);
    }

    function createBill(
        address user,
        uint256 operatorId,
        uint256 amount
    ) external onlyOracle {
        require(amount > 0, "Zero amount");

        uint256 fee = feeManager.calculateFee(amount);
        uint256 billId = _nextBillId++;

        _bills[billId] = Bill({
            id: billId,
            user: user,
            operatorId: operatorId,
            amount: amount,
            platformFee: fee,
            createdAt: block.timestamp,
            isPaid: false,
            trafficCardDeduction: 0
        });
        _userBillIds[user].push(billId);

        emit BillCreated(billId, user, amount, fee);
    }

    /// @notice 为账单应用流量卡抵扣（由Oracle在月末结算时调用）
    function applyTrafficCardToBill(uint256 billId) public onlyOracle {
        Bill storage bill = _bills[billId];
        require(!bill.isPaid, "Bill already paid");

        address user = bill.user;
        uint256 availableCredit = deposit.getTrafficCardBalance(user);
        
        // 抵扣金额不超过账单金额+手续费，且不超过可用余额
        uint256 totalBill = bill.amount + bill.platformFee;
        uint256 deduction = availableCredit > totalBill ? totalBill : availableCredit;
        
        if (deduction > 0) {
            bill.trafficCardDeduction = deduction;
            deposit.useTrafficCard(user, deduction);
            
            emit TrafficCardApplied(billId, user, deduction);
        }
    }

    /// @notice 用户支付账单
    function payBill(uint256 billId) external payable nonReentrant {
        Bill storage bill = _bills[billId];
        require(bill.user == msg.sender, "Not your bill");
        require(!bill.isPaid, "Already paid");

        // 计算实际应付金额：账单金额 + 平台手续费 - 流量卡抵扣
        uint256 total = (bill.amount + bill.platformFee) - bill.trafficCardDeduction;
        require(msg.value >= total, "Insufficient payment");

        bill.isPaid = true;

        // 手续费转平台
        if (bill.platformFee > 0) {
            (bool feeOk, ) = platformWallet.call{value: bill.platformFee}("");
            require(feeOk, "Fee transfer failed");
        }

        // 退还多付
        if (msg.value > total) {
            (bool refundOk, ) = msg.sender.call{value: msg.value - total}("");
            require(refundOk, "Refund failed");
        }

        emit BillPaid(billId, msg.sender, total);
    }

    function getUserBills(address user) external view returns (Bill[] memory) {
        uint256[] memory ids = _userBillIds[user];
        Bill[] memory bills = new Bill[](ids.length);
        for (uint256 i = 0; i < ids.length; i++) {
            bills[i] = _bills[ids[i]];
        }
        return bills;
    }

    function getUnpaidBills(address user) external view returns (Bill[] memory) {
        uint256[] memory ids = _userBillIds[user];
        uint256 count = 0;
        for (uint256 i = 0; i < ids.length; i++) {
            if (!_bills[ids[i]].isPaid) count++;
        }

        Bill[] memory result = new Bill[](count);
        uint256 idx = 0;
        for (uint256 i = 0; i < ids.length; i++) {
            if (!_bills[ids[i]].isPaid) {
                result[idx++] = _bills[ids[i]];
            }
        }
        return result;
    }

    /// @notice 自动结算（月底触发，将预言机数据生成账单发送给用户）
    function autoSettle(
        address[] calldata users,
        uint256[] calldata operatorIds,
        uint256[] calldata amounts
    ) external onlyOracle {
        require(users.length == operatorIds.length, "Length mismatch");
        require(users.length == amounts.length, "Length mismatch");

        for (uint256 i = 0; i < users.length; i++) {
            if (amounts[i] > 0) {
                uint256 fee = feeManager.calculateFee(amounts[i]);
                uint256 billId = _nextBillId++;

                _bills[billId] = Bill({
                    id: billId,
                    user: users[i],
                    operatorId: operatorIds[i],
                    amount: amounts[i],
                    platformFee: fee,
                    createdAt: block.timestamp,
                    isPaid: false,
                    trafficCardDeduction: 0
                });
                _userBillIds[users[i]].push(billId);

                emit BillCreated(billId, users[i], amounts[i], fee);
            }
        }

        // 自动应用流量卡抵扣（如果Deposit合约已设置）
        if (address(deposit) != address(0)) {
            deposit.issueMonthlyTrafficCards(users);
            for (uint256 i = 0; i < users.length; i++) {
                IPayment.Bill[] memory unpaidBills = getUnpaidBills(users[i]);
                if (unpaidBills.length > 0) {
                    applyTrafficCardToBill(unpaidBills[0].id);
                }
            }
        }
    }

    /// @inheritdoc UUPSUpgradeable
    function _authorizeUpgrade(address newImplementation) internal override onlyOwner {}
}