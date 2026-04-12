// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

import "@openzeppelin/contracts/access/Ownable.sol";
import "./interfaces/IOracle.sol";
import "./interfaces/IPayment.sol";

/// @title Oracle - 计量预言机（从运营商获取账单并分发至用户）
contract Oracle is IOracle, Ownable {
    IPayment public payment;
    
    mapping(address => mapping(uint256 => UsageInfo)) private _latestUsage;

    struct UsageInfo {
        uint256 dataUsage;
        uint256 callUsage;
        uint256 timestamp;
    }

    constructor(address _payment) Ownable(msg.sender) {
        payment = IPayment(_payment);
    }

    function setPayment(address _payment) external onlyOwner {
        payment = IPayment(_payment);
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
