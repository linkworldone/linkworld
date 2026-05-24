// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import "./interfaces/IOracle.sol";
import "./interfaces/IPayment.sol";
import "./interfaces/IDeposit.sol";

/// @title Oracle - 计量预言机（从运营商获取账单并分发至用户，可升级版）
contract Oracle is IOracle, OwnableUpgradeable, UUPSUpgradeable {
    IPayment public payment;
    IDeposit public deposit; // 关联 Deposit 合约

    // 累计使用量（本月）- 不上链，仅月末汇总
    mapping(address => mapping(uint256 => UsageInfo)) private _monthlyUsage;
    // 最新提交的使用数据（用于查询）
    mapping(address => mapping(uint256 => UsageInfo)) private _latestUsage;

    struct UsageInfo {
        uint256 dataUsage;
        uint256 callUsage;
        uint256 timestamp;
    }

    /// @notice Initializer
    function initialize(address _payment) public initializer {
        __Ownable_init(msg.sender);

        payment = IPayment(_payment);
    }

    function setPayment(address _payment) external onlyOwner {
        payment = IPayment(_payment);
    }

    /// @notice 设置 Deposit 合约地址
    function setDeposit(address _deposit) external onlyOwner {
        deposit = IDeposit(_deposit);
    }

    /// @notice 月末结算（统一汇总上链，减少Gas消耗）
    /// @param users 用户地址列表
    /// @param operatorIds 运营商ID列表
    /// @param dataUsages 流量使用量列表（字节）- 如果为0则使用累计数据
    /// @param callUsages 通话时长列表（分钟）- 如果为0则使用累计数据
    function monthlySettlement(
        address[] calldata users,
        uint256[] calldata operatorIds,
        uint256[] calldata dataUsages,
        uint256[] calldata callUsages
    ) external onlyOwner {
        require(users.length == operatorIds.length, "Length mismatch");
        require(users.length == dataUsages.length, "Length mismatch");
        require(users.length == callUsages.length, "Length mismatch");

        // 第一步：统一汇总上链，生成账单
        for (uint256 i = 0; i < users.length; i++) {
            address user = users[i];
            uint256 operatorId = operatorIds[i];

            uint256 dataUsage = dataUsages[i] > 0 ? dataUsages[i] : _monthlyUsage[user][operatorId].dataUsage;
            uint256 callUsage = callUsages[i] > 0 ? callUsages[i] : _monthlyUsage[user][operatorId].callUsage;

            if (dataUsage > 0 || callUsage > 0) {
                uint256 totalAmount = dataUsage + callUsage;
                payment.createBill(user, operatorId, totalAmount);

                emit UsageDataSubmitted(user, operatorId, dataUsage, callUsage);
            }

            delete _monthlyUsage[user][operatorId];
        }

        // 第二步：为符合条件用户发放流量卡NFT
        if (address(deposit) != address(0)) {
            deposit.issueMonthlyTrafficCards(users);
        }

        // 第三步：为每个用户的未支付账单应用流量卡抵扣
        if (address(payment) != address(0)) {
            for (uint256 i = 0; i < users.length; i++) {
                // viaIR: cast payment to IPayment so call opens a proper LSP interface path
                IPayment.Bill[] memory unpaidBills = IPayment(payment).getUnpaidBills(users[i]);
                if (unpaidBills.length > 0) {
                    payment.applyTrafficCardToBill(unpaidBills[0].id);
                }
            }
        }
    }

    // submitUsage 内联了 recordUsage 逻辑（viaIR 下跨外部函数调用可见性异常）
    function submitUsage(
        address user,
        uint256 operatorId,
        uint256 dataUsage,
        uint256 callUsage
    ) public onlyOwner {
        require(user != address(0), "Invalid user");
        require(dataUsage > 0 || callUsage > 0, "Zero usage");

        _monthlyUsage[user][operatorId].dataUsage += dataUsage;
        _monthlyUsage[user][operatorId].callUsage += callUsage;
        _monthlyUsage[user][operatorId].timestamp = block.timestamp;

        _latestUsage[user][operatorId] = UsageInfo({
            dataUsage: dataUsage,
            callUsage: callUsage,
            timestamp: block.timestamp
        });

        emit UsageDataSubmitted(user, operatorId, dataUsage, callUsage);
    }

    function getLatestUsage(address user, uint256 operatorId)
        external
        view
        returns (uint256 dataUsage, uint256 callUsage, uint256 timestamp)
    {
        UsageInfo memory info = _latestUsage[user][operatorId];
        return (info.dataUsage, info.callUsage, info.timestamp);
    }

    /// @inheritdoc UUPSUpgradeable
    function _authorizeUpgrade(address newImplementation) internal override onlyOwner {}
}
