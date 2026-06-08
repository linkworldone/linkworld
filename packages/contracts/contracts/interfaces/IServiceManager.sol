// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

interface IServiceManager {
    struct Operator {
        uint256 id;
        string name;
        string region;
        string countryCode;
        uint256 requiredDeposit;
        bool isActive;
        address paymentAddress;
    }

    event OperatorAdded(uint256 indexed operatorId, string name, string region);
    event OperatorUpdated(uint256 indexed operatorId);
    event OperatorDeactivated(uint256 indexed operatorId);

    function addOperator(
        string calldata name,
        string calldata region,
        string calldata countryCode,
        uint256 requiredDeposit,
        address paymentAddress
    ) external;
    function updateOperator(
        uint256 operatorId,
        string calldata name,
        string calldata region,
        uint256 requiredDeposit
    ) external;
    function deactivateOperator(uint256 operatorId) external;
    function getOperator(uint256 operatorId) external view returns (Operator memory);
    function getActiveOperators() external view returns (Operator[] memory);
    function getOperatorsByCountry(string calldata countryCode) external view returns (Operator[] memory);
}
