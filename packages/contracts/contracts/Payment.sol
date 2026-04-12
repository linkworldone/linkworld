// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import "./interfaces/IPayment.sol";
import "./interfaces/IFeeManager.sol";

/// @title Payment - 支付结算（服务商费用 + 平台手续费）
contract Payment is IPayment, Ownable, ReentrancyGuard {
    IFeeManager public feeManager;

    address public platformWallet;
    address public oracle; // 预言机地址，有权创建账单

    uint256 private _nextBillId;
    mapping(uint256 => Bill) private _bills;
    mapping(address => uint256[]) private _userBillIds;

    modifier onlyOracle() {
        require(msg.sender == oracle, "Only oracle");
        _;
    }

    constructor(
        address _feeManager,
        address _platformWallet
    ) Ownable(msg.sender) {
        feeManager = IFeeManager(_feeManager);
        platformWallet = _platformWallet;
    }

    function setOracle(address _oracle) external onlyOwner {
        oracle = _oracle;
    }

    function setPlatformWallet(address _platformWallet) external onlyOwner {
        platformWallet = _platformWallet;
    }

    /// @notice 预言机创建账单
    function createBill(
        address user,
        uint256 operatorId,
        uint256 amount
    ) external onlyOracle {
        uint256 fee = feeManager.calculateFee(amount);
        uint256 billId = _nextBillId++;

        _bills[billId] = Bill({
            id: billId,
            user: user,
            operatorId: operatorId,
            amount: amount,
            platformFee: fee,
            createdAt: block.timestamp,
            isPaid: false
        });
        _userBillIds[user].push(billId);

        emit BillCreated(billId, user, amount, fee);
    }

    /// @notice 用户支付账单
    function payBill(uint256 billId) external payable nonReentrant {
        Bill storage bill = _bills[billId];
        require(bill.user == msg.sender, "Not your bill");
        require(!bill.isPaid, "Already paid");

        uint256 total = bill.amount + bill.platformFee;
        require(msg.value >= total, "Insufficient payment");

        bill.isPaid = true;

        // 手续费转平台
        (bool feeOk, ) = platformWallet.call{value: bill.platformFee}("");
        require(feeOk, "Fee transfer failed");

        // 运营商费用留在合约，由管理员分发（或直接转运营商钱包）

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
                    isPaid: false
                });
                _userBillIds[users[i]].push(billId);

                emit BillCreated(billId, users[i], amounts[i], fee);
            }
        }
    }
}
