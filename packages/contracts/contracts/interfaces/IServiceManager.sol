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
    }

    struct UserService {
        address user;
        uint256 operatorId;
        string virtualNumber;
        string password;
        uint256 activatedAt;
        bool isActive;
    }

    event OperatorAdded(uint256 indexed operatorId, string name, string region);
    event OperatorUpdated(uint256 indexed operatorId);
    event OperatorDeactivated(uint256 indexed operatorId);
    event UserServiceActivated(address indexed user, uint256 operatorId, string virtualNumber);
    event UserServiceDeactivated(address indexed user, uint256 operatorId);

    function addOperator(string calldata name, string calldata region, string calldata countryCode, uint256 requiredDeposit) external;
    function getOperator(uint256 operatorId) external view returns (Operator memory);
    function getActiveOperators() external view returns (Operator[] memory);
    function getOperatorsByCountry(string calldata countryCode) external view returns (Operator[] memory);
    
    function activateService(uint256 operatorId, string calldata virtualNumber, string calldata password) external;
    function deactivateService() external;
    function getUserService(address user) external view returns (UserService memory);
}
