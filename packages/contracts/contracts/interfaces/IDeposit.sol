// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

interface IDeposit {
    event DepositMade(address indexed user, uint256 amount);
    event DepositWithdrawn(address indexed user, uint256 principal, uint256 interest);
    event InterestRateUpdated(uint256 newRate);

    function deposit() external payable;
    function withdraw() external;
    function getDepositAmount(address user) external view returns (uint256);
    function getRequiredDeposit(uint256 operatorId) external view returns (uint256);
    function setServiceManager(address _serviceManager) external;
    
    function updateInterestRate(uint256 newRate) external;
    function calculateInterest(address user) external view returns (uint256);
    function getInterestAmount(address user) external view returns (uint256);
    function interestRate() external view returns (uint256);
}