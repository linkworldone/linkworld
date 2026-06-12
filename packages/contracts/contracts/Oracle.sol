// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import "./interfaces/IOracle.sol";
import "./interfaces/IPayment.sol";
import "./interfaces/IDeposit.sol";

/// @title Oracle - 服务验证与计量预言机（v3 简化版）
contract Oracle is IOracle, OwnableUpgradeable, UUPSUpgradeable {
    IDeposit public deposit;
    IPayment public payment;

    /// @notice 月末结算喂价记录事件（v2-A：记录后端链下算好的 USDT 金额，Oracle 不参与计价）
    event UsageDataSubmitted(
        address indexed user,
        uint256 operatorId,
        uint256 amount
    );

    /// @notice Initializer
    function initialize() public initializer {
        __Ownable_init(msg.sender);
    }

    /// @notice 设置 Deposit 合约地址
    function setDeposit(address _deposit) external onlyOwner {
        deposit = IDeposit(_deposit);
    }

    /// @notice 设置 Payment 合约地址
    function setPayment(address _payment) external onlyOwner {
        payment = IPayment(_payment);
    }

    /// @notice 验证用户服务是否活跃（预留接口，未来接入 Chainlink/0G Compute）
    function verifyServiceActive(address user) external view returns (bool) {
        if (address(deposit) == address(0)) return false;
        (bool ok, bytes memory data) = address(deposit).staticcall(
            abi.encodeWithSignature("getLockExpiry(address)", user)
        );
        if (!ok) return false;
        (uint256 lockExpiry) = abi.decode(data, (uint256));
        return lockExpiry > 0 && block.timestamp < lockExpiry;
    }

    /// @notice 月末结算（统一汇总上链，减少Gas消耗）
    /// @dev v2-A/B1：Oracle 不计价。amounts[i] 是后端链下按资费算好的 USDT 金额，
    ///      直接传给 createBill；删除原 totalAmount = dataUsage + callUsage 量纲错误求和。
    /// @param users 用户地址列表
    /// @param operatorIds 运营商ID列表
    /// @param amounts 链下算好的 USDT 金额列表（精度 = usdt.decimals()）
    function monthlySettlement(
        address[] calldata users,
        uint256[] calldata operatorIds,
        uint256[] calldata amounts
    ) external onlyOwner {
        require(users.length == operatorIds.length, "Length mismatch");
        require(users.length == amounts.length, "Length mismatch");

        for (uint256 i = 0; i < users.length; i++) {
            address user = users[i];
            uint256 operatorId = operatorIds[i];
            uint256 amount = amounts[i];

<<<<<<< HEAD
            uint256 dataUsage = dataUsages[i] > 0 ? dataUsages[i] : _monthlyUsage[user][operatorId].dataUsage;
            uint256 callUsage = callUsages[i] > 0 ? callUsages[i] : _monthlyUsage[user][operatorId].callUsage;

            if (dataUsage > 0 || callUsage > 0) {
                uint256 totalAmount = dataUsage + callUsage;
                IPayment(payment).createBill(user, operatorId, totalAmount);

                emit UsageDataSubmitted(user, operatorId, dataUsage, callUsage);
=======
            if (amount > 0) {
                payment.createBill(user, operatorId, amount);
                emit UsageDataSubmitted(user, operatorId, amount);
>>>>>>> 48e2c2a1bafc65ef37ea369d3f85ef9eb004e69d
            }
        }

<<<<<<< HEAD
        if (address(deposit) != address(0)) {
            IDeposit(deposit).issueMonthlyTrafficCards(users);
        }
=======
        // 发卡逻辑已迁移：Deposit 改为"充值即按比例发卡"，不再有月度批量发卡入口。
        // 这里只做账单生成 + 流量卡抵扣（用户在 deposit 时已持卡）。
>>>>>>> 48e2c2a1bafc65ef37ea369d3f85ef9eb004e69d

        if (address(payment) != address(0)) {
            for (uint256 i = 0; i < users.length; i++) {
                IPayment.Bill[] memory unpaidBills = IPayment(payment).getUnpaidBills(users[i]);
                if (unpaidBills.length > 0) {
                    IPayment(payment).applyTrafficCardToBill(unpaidBills[0].id);
                }
            }
        }
    }

    /// @notice 提交使用数据（预留接口，用于 0G Compute/Chainlink 接入）
    /// @dev 预留接口，本轮 v2-A 不参与计价；旁路去留与后端(2/3)对齐（design §4.4 / arch-review ⚠️5）。
    function submitUsage(
        address user,
        uint256,
        uint256,
        uint256
    ) external view {
        require(user != address(0), "Invalid user");
    }

    /// @inheritdoc UUPSUpgradeable
    function _authorizeUpgrade(address newImplementation) internal override onlyOwner {}
}