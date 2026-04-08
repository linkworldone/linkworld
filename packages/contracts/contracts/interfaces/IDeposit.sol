// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

interface IDeposit {
    event DepositMade(address indexed user, uint256 amount);
    event DepositWithdrawn(address indexed user, uint256 amount);

    function deposit() external payable;
    function withdraw() external;
    function getDepositAmount(address user) external view returns (uint256);
    function getRequiredDeposit(uint256 operatorId) external view returns (uint256);
}
