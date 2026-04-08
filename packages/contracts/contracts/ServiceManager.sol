// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

import "@openzeppelin/contracts/access/Ownable.sol";
import "./interfaces/IServiceManager.sol";

/// @title ServiceManager - 运营商信息存储与管理
contract ServiceManager is IServiceManager, Ownable {
    uint256 private _nextOperatorId;
    mapping(uint256 => Operator) private _operators;
    uint256[] private _activeOperatorIds;

    constructor() Ownable(msg.sender) {}

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
}
