// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import "./interfaces/IOracle.sol";

/// @title Oracle - 服务验证与计量预言机（v3 简化版）
contract Oracle is IOracle, OwnableUpgradeable, UUPSUpgradeable {
    address public deposit;
    address public payment;

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
    function initialize() public initializer {
        __Ownable_init(msg.sender);
    }

    /// @notice 设置 Deposit 合约地址
    function setDeposit(address _deposit) external onlyOwner {
        deposit = _deposit;
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

        if (address(deposit) != address(0)) {
            deposit.issueMonthlyTrafficCards(users);
        }

        if (address(payment) != address(0)) {
            for (uint256 i = 0; i < users.length; i++) {
                IPayment.Bill[] memory unpaidBills = IPayment(payment).getUnpaidBills(users[i]);
                if (unpaidBills.length > 0) {
                    payment.applyTrafficCardToBill(unpaidBills[0].id);
                }
            }
        }
    }

    /// @notice 提交使用数据（预留接口，用于 0G Compute/Chainlink 接入）
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

interface IPayment {
    struct Bill {
        uint256 id;
        address user;
        uint256 operatorId;
        uint256 amount;
        uint256 platformFee;
        uint256 createdAt;
        bool isPaid;
    }
    function createBill(address user, uint256 operatorId, uint256 amount) external;
    function getUnpaidBills(address user) external view returns (Bill[] memory);
    function applyTrafficCardToBill(uint256 billId) external;
}

interface IDeposit {
    function issueMonthlyTrafficCards(address[] calldata users) external;
}