// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

interface IOracle {
    struct UsageData {
        address user;
        uint256 operatorId;
        uint256 dataUsage;
        uint256 callUsage;
        uint256 timestamp;
        bytes signature;
    }

    event UsageDataSubmitted(address indexed user, uint256 operatorId, uint256 dataUsage, uint256 callUsage);
    event MonthlySettlementCompleted(uint256 timestamp);

    function submitUsage(
        address user,
        uint256 operatorId,
        uint256 dataUsage,
        uint256 callUsage
    ) external;

    function getLatestUsage(address user, uint256 operatorId) external view returns (uint256 dataUsage, uint256 callUsage, uint256 timestamp);
    
    function setDeposit(address _deposit) external;
    
    // 月末结算
    function monthlySettlement(
        address[] calldata users,
        uint256[] calldata operatorIds,
        uint256[] calldata dataUsages,
        uint256[] calldata callUsages
    ) external;
}