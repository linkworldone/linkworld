// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

import "@openzeppelin/contracts/access/Ownable.sol";
import "./interfaces/IServiceManager.sol";

/// @title ServiceManager - 运营商信息存储管理与用户服务
contract ServiceManager is IServiceManager, Ownable {
    uint256 private _nextOperatorId;
    mapping(uint256 => Operator) private _operators;
    uint256[] private _activeOperatorIds;

    mapping(address => UserService) private _userServices;

    constructor() Ownable(msg.sender) {
        _nextOperatorId = 1;
        _operators[1] = Operator({
            id: 1,
            name: "T-Mobile US",
            region: "United States",
            requiredDeposit: 0.01 ether,
            isActive: true
        });
        _activeOperatorIds.push(1);

        _operators[2] = Operator({
            id: 2,
            name: "Vodafone UK",
            region: "United Kingdom",
            requiredDeposit: 0.008 ether,
            isActive: true
        });
        _activeOperatorIds.push(2);
        _nextOperatorId = 3;
    }

    function addOperator(
        string calldata name,
        string calldata region,
        uint256 requiredDeposit
    ) external onlyOwner {
        uint256 operatorId = _nextOperatorId++;
        _operators[operatorId] = Operator({
            id: operatorId,
            name: name,
            region: region,
            requiredDeposit: requiredDeposit,
            isActive: true
        });
        _activeOperatorIds.push(operatorId);

        emit OperatorAdded(operatorId, name, region);
    }

    function updateOperator(
        uint256 operatorId,
        string calldata name,
        string calldata region,
        uint256 requiredDeposit
    ) external onlyOwner {
        require(_operators[operatorId].isActive, "Operator not found");
        Operator storage op = _operators[operatorId];
        op.name = name;
        op.region = region;
        op.requiredDeposit = requiredDeposit;

        emit OperatorUpdated(operatorId);
    }

    function deactivateOperator(uint256 operatorId) external onlyOwner {
        require(_operators[operatorId].isActive, "Operator not found");
        _operators[operatorId].isActive = false;

        emit OperatorDeactivated(operatorId);
    }

    function getOperator(uint256 operatorId) external view returns (Operator memory) {
        return _operators[operatorId];
    }

    function getActiveOperators() external view returns (Operator[] memory) {
        uint256 count = 0;
        for (uint256 i = 0; i < _activeOperatorIds.length; i++) {
            if (_operators[_activeOperatorIds[i]].isActive) count++;
        }

        Operator[] memory result = new Operator[](count);
        uint256 idx = 0;
        for (uint256 i = 0; i < _activeOperatorIds.length; i++) {
            if (_operators[_activeOperatorIds[i]].isActive) {
                result[idx++] = _operators[_activeOperatorIds[i]];
            }
        }
        return result;
    }

    function activateService(uint256 operatorId, string calldata virtualNumber, string calldata password) external {
        require(_operators[operatorId].isActive, "Operator not found");
        require(!_userServices[msg.sender].isActive, "Service already active");

        _userServices[msg.sender] = UserService({
            user: msg.sender,
            operatorId: operatorId,
            virtualNumber: virtualNumber,
            password: password,
            activatedAt: block.timestamp,
            isActive: true
        });

        emit UserServiceActivated(msg.sender, operatorId, virtualNumber);
    }

    function deactivateService() external {
        require(_userServices[msg.sender].isActive, "No active service");

        uint256 operatorId = _userServices[msg.sender].operatorId;
        _userServices[msg.sender].isActive = false;

        emit UserServiceDeactivated(msg.sender, operatorId);
    }

    function getUserService(address user) external view returns (UserService memory) {
        return _userServices[user];
    }
}
