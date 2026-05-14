// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

import "@openzeppelin/contracts/access/Ownable.sol";
import "./interfaces/IOracle.sol";
import "./interfaces/IPayment.sol";
import "./interfaces/IDeposit.sol";

/// @title Oracle - 计量预言机（从运营商获取账单并分发至用户）
contract Oracle is IOracle, Ownable {
    IPayment public payment;
    IDeposit public deposit; // 关联 Deposit 合约
    
    mapping(address => mapping(uint256 => UsageInfo)) private _latestUsage;

    struct UsageInfo {
        uint256 dataUsage;
        uint256 callUsage;
        uint256 timestamp;
    }

    // 添加月末结算完成事件（接口中无，仅在合约内）
    event MonthlySettlementCompleted(uint256 timestamp);

    constructor(address _payment) Ownable(msg.sender) {
        payment = IPayment(_payment);
    }

    function setPayment(address _payment) external onlyOwner {
        payment = IPayment(_payment);
    }

    /// @notice 设置 Deposit 合约地址
    function setDeposit(address _deposit) external onlyOwner {
        deposit = IDeposit(_deposit);
    }

    /// @notice 月末结算（提交流量使用数据并发放流量卡）
    /// @param users 用户地址列表
    /// @param operatorIds 运营商ID列表
    /// @param dataUsages 流量使用量列表（字节）
    /// @param callUsages 通话时长列表（分钟）
    function monthlySettlement(
        address[] calldata users,
        uint256[] calldata operatorIds,
        uint256[] calldata dataUsages,
        uint256[] calldata callUsages
    ) external onlyOwner {
        require(users.length == operatorIds.length, "Length mismatch");
        require(users.length == dataUsages.length, "Length mismatch");
        require(users.length == callUsages.length, "Length mismatch");

        // 第一步：提交使用数据，生成账单
        for (uint256 i = 0; i < users.length; i++) {
            submitUsage(users[i], operatorIds[i], dataUsages[i], callUsages[i]);
        }

        // 第二步：为符合条件用户发放流量卡
        if (address(deposit) != address(0)) {
            deposit.issueMonthlyTrafficCards(users);
        }

        // 第三步：为每个用户的未支付账单应用流量卡抵扣
        if (address(payment) != address(0)) {
            for (uint256 i = 0; i < users.length; i++) {
                // 获取该用户最新的未支付账单
                IPayment.Bill[] memory unpaidBills = payment.getUnpaidBills(users[i]);
                if (unpaidBills.length > 0) {
                    // 应用流量卡抵扣到最新账单
                    payment.applyTrafficCardToBill(unpaidBills[0].id);
                }
            }
        }

        emit MonthlySettlementCompleted(block.timestamp);
    }

    /// @notice 预言机提交使用数据（由预言机角色调用）
    /// @dev 实际场景中需验证签名，此处简化处理
    function submitUsage(
        address user,
        uint256 operatorId,
        uint256 dataUsage,
        uint256 callUsage
    ) external onlyOwner {
        require(user != address(0), "Invalid user");
        require(dataUsage > 0 || callUsage > 0, "Zero usage");

        _latestUsage[user][operatorId] = UsageInfo({
            dataUsage: dataUsage,
            callUsage: callUsage,
            timestamp: block.timestamp
        });

        uint256 totalAmount = dataUsage + callUsage;
        
        payment.createBill(user, operatorId, totalAmount);

        emit UsageDataSubmitted(user, operatorId, dataUsage, callUsage);
    }

    /// @notice 获取用户最新使用数据
    function getLatestUsage(address user, uint256 operatorId) external view returns (uint256 dataUsage, uint256 callUsage, uint256 timestamp) {
        UsageInfo memory info = _latestUsage[user][operatorId];
        return (info.dataUsage, info.callUsage, info.timestamp);
    }
}

// 补充事件定义（在合约外部，实际应放在合约内）