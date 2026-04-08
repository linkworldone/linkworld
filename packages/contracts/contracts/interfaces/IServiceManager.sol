// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

interface IServiceManager {
    struct Operator {
        uint256 id;
        string name;
        string region;
        uint256 requiredDeposit;
        bool isActive;
    }

    event OperatorAdded(uint256 indexed operatorId, string name, string region);
    event OperatorUpdated(uint256 indexed operatorId);
    event OperatorDeactivated(uint256 indexed operatorId);

    function addOperator(string calldata name, string calldata region, uint256 requiredDeposit) external;
    function getOperator(uint256 operatorId) external view returns (Operator memory);
    function getActiveOperators() external view returns (Operator[] memory);
}
