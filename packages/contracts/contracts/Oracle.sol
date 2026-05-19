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

    // 添加月末结算完成事件（接口中无，仅在合约内）
    event MonthlySettlementCompleted(uint256 timestamp);

    /// @inheritdoc IOracle
    function initialize(address _payment) public initializer {
        __Ownable_init(msg.sender);
        // __UUPSUpgradeable_init() not needed

        payment = IPayment(_payment);
    }

    function setPayment(address _payment) external onlyOwner {
        payment = IPayment(_payment);
    }

    /// @notice 设置 Deposit 合约地址
    function setDeposit(address _deposit) external onlyOwner {
        deposit = IDeposit(_deposit);
    }

    /// @notice 累计使用数据（不上链，仅记录）
    /// @dev 用于实时记录使用量，月末统一汇总上链
    function recordUsage(
        address user,
        uint256 operatorId,
        uint256 dataUsage,
        uint256 callUsage
    ) external onlyOwner {
        require(user != address(0), "Invalid user");
        require(dataUsage > 0 || callUsage > 0, "Zero usage");

        // 累计到本月使用量（不上链消耗）
        _monthlyUsage[user][operatorId].dataUsage += dataUsage;
        _monthlyUsage[user][operatorId].callUsage += callUsage;
        _monthlyUsage[user][operatorId].timestamp = block.timestamp;

        // 更新最新使用记录
        _latestUsage[user][operatorId] = UsageInfo({
            dataUsage: dataUsage,
            callUsage: callUsage,
            timestamp: block.timestamp
        });

        emit UsageDataSubmitted(user, operatorId, dataUsage, callUsage);
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
            
            // 如果传入数据为0，则使用累计的数据
            uint256 dataUsage = dataUsages[i] > 0 ? dataUsages[i] : _monthlyUsage[user][operatorId].dataUsage;
            uint256 callUsage = callUsages[i] > 0 ? callUsages[i] : _monthlyUsage[user][operatorId].callUsage;
            
            if (dataUsage > 0 || callUsage > 0) {
                uint256 totalAmount = dataUsage + callUsage;
                payment.createBill(user, operatorId, totalAmount);
                
                emit UsageDataSubmitted(user, operatorId, dataUsage, callUsage);
            }
            
            // 重置本月累计数据
            delete _monthlyUsage[user][operatorId];
        }

        // 第二步：为符合条件用户发放流量卡NFT
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

    /// @notice 预言机提交使用数据（保留原有接口，兼容旧系统）
    /// @dev 实际场景中需验证签名，此处简化处理
    function submitUsage(
        address user,
        uint256 operatorId,
        uint256 dataUsage,
        uint256 callUsage
    ) public onlyOwner {
        // 调用新的记录方法（仅记录，不上链）
        recordUsage(user, operatorId, dataUsage, callUsage);
    }

    /// @notice 获取用户最新使用数据
    function getLatestUsage(address user, uint256 operatorId) external view returns (uint256 dataUsage, uint256 callUsage, uint256 timestamp) {
        UsageInfo memory info = _latestUsage[user][operatorId];
        return (info.dataUsage, info.callUsage, info.timestamp);
    }

    /// @inheritdoc UUPSUpgradeable
    function _authorizeUpgrade(address newImplementation) internal override onlyOwner {}
}