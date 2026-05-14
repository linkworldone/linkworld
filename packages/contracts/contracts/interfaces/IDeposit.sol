// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

interface IDeposit {
    event DepositMade(address indexed user, uint256 amount);
    event DepositWithdrawn(address indexed user, uint256 principal, uint256 interest); // interest 字段已弃用，始终为0
    event TrafficCardQuotaUpdated(uint256 newQuota);
    event TrafficCardIssued(address indexed user, uint256 amount, uint256 month);
    event TrafficCardUsed(address indexed user, uint256 amount);

    function deposit() external payable;
    function withdraw() external;
    function getDepositAmount(address user) external view returns (uint256);
    function getRequiredDeposit(uint256 operatorId) external view returns (uint256);
    function setServiceManager(address _serviceManager) external;
    
    // 流量卡管理
    function setTrafficCardQuota(uint256 quota) external;
    function issueMonthlyTrafficCards(address[] calldata users) external;
    function useTrafficCard(address user, uint256 amount) external;
    function getTrafficCardBalance(address user) external view returns (uint256);
    function hasWithdrawnThisMonth(address user, uint256 month) external view returns (bool);
    function getCurrentMonth() external view returns (uint256);
    
    // 查询
    function trafficCardQuota() external view returns (uint256);
}