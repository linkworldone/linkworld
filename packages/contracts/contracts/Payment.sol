// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import "./interfaces/IPayment.sol";
import "./interfaces/IFeeManager.sol";

/// @title Payment - 购买运营商服务（v3: 1.5% 手续费，实时结算）
contract Payment is IPayment, OwnableUpgradeable, UUPSUpgradeable {
    IFeeManager public feeManager;
    address public platformWallet;

    uint256 private _nextBillId;
    mapping(uint256 => Bill) private _bills;
    mapping(address => uint256[]) private _userBillIds;

    /// @notice Initializer
    function initialize(address _feeManager, address _platformWallet) public initializer {
        __Ownable_init(msg.sender);
        feeManager = IFeeManager(_feeManager);
        platformWallet = _platformWallet;
    }

    /// @notice 设置平台钱包
    function setPlatformWallet(address _platformWallet) external onlyOwner {
        platformWallet = _platformWallet;
    }

    /// @notice 设置手续费合约
    function setFeeManager(address _feeManager) external onlyOwner {
        feeManager = IFeeManager(_feeManager);
    }

    /// @notice 创建服务购买账单（仅 oracle/owner）
    function createBill(address user, uint256 operatorId, uint256 amount) external onlyOwner {
        require(amount > 0, "Zero amount");

        uint256 platformFee = feeManager.calculateFee(amount);
        uint256 billId = _nextBillId++;

        _bills[billId] = Bill({
            id: billId,
            user: user,
            operatorId: operatorId,
            amount: amount,
            platformFee: platformFee,
            createdAt: block.timestamp,
            isPaid: false
        });
        _userBillIds[user].push(billId);

        emit BillCreated(billId, user, amount + platformFee, platformFee);
    }

    /// @notice 用户支付账单
    function payBill(uint256 billId) external payable {
        Bill storage bill = _bills[billId];
        require(!bill.isPaid, "Already paid");
        require(bill.user == msg.sender, "Not your bill");

        uint256 total = bill.amount + bill.platformFee;
        require(msg.value >= total, "Insufficient payment");

        bill.isPaid = true;

        if (bill.platformFee > 0) {
            (bool feeOk, ) = platformWallet.call{value: bill.platformFee}("");
            require(feeOk, "Fee transfer failed");
        }

        if (msg.value > total) {
            (bool refundOk, ) = msg.sender.call{value: msg.value - total}("");
            require(refundOk, "Refund failed");
        }

        emit BillPaid(billId, msg.sender, total, bill.amount);
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

    /// @inheritdoc UUPSUpgradeable
    function _authorizeUpgrade(address newImplementation) internal override onlyOwner {}
}
