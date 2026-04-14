// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

import "@openzeppelin/contracts/access/Ownable.sol";
import "./interfaces/IServiceManager.sol";

/// @title ServiceManager - 运营商信息存储管理与用户服务
contract ServiceManager is IServiceManager, Ownable {
    uint256 private _nextOperatorId;
    mapping(uint256 => Operator) private _operators;
    uint256[] private _activeOperatorIds;
    mapping(string => uint256[]) private _countryOperatorIds;

    mapping(address => UserService) private _userServices;

    constructor() Ownable(msg.sender) {
        _nextOperatorId = 1;
        
        _operators[1] = Operator({
            id: 1,
            name: "T-Mobile US",
            region: "United States",
            countryCode: "US",
            requiredDeposit: 0.01 ether,
            isActive: true
        });
        _activeOperatorIds.push(1);
        _countryOperatorIds["US"].push(1);

        _operators[2] = Operator({
            id: 2,
            name: "Vodafone UK",
            region: "United Kingdom",
            countryCode: "GB",
            requiredDeposit: 0.008 ether,
            isActive: true
        });
        _activeOperatorIds.push(2);
        _countryOperatorIds["GB"].push(2);

        _operators[3] = Operator({
            id: 3,
            name: "Orange France",
            region: "France",
            countryCode: "FR",
            requiredDeposit: 0.008 ether,
            isActive: true
        });
        _activeOperatorIds.push(3);
        _countryOperatorIds["FR"].push(3);

        _operators[4] = Operator({
            id: 4,
            name: "MTS Russia",
            region: "Russia",
            countryCode: "RU",
            requiredDeposit: 0.005 ether,
            isActive: true
        });
        _activeOperatorIds.push(4);
        _countryOperatorIds["RU"].push(4);

        _operators[5] = Operator({
            id: 5,
            name: "SoftBank Japan",
            region: "Japan",
            countryCode: "JP",
            requiredDeposit: 0.012 ether,
            isActive: true
        });
        _activeOperatorIds.push(5);
        _countryOperatorIds["JP"].push(5);

        _operators[6] = Operator({
            id: 6,
            name: "Viettel Vietnam",
            region: "Vietnam",
            countryCode: "VN",
            requiredDeposit: 0.003 ether,
            isActive: true
        });
        _activeOperatorIds.push(6);
        _countryOperatorIds["VN"].push(6);

        _operators[7] = Operator({
            id: 7,
            name: "Unitel Laos",
            region: "Laos",
            countryCode: "LA",
            requiredDeposit: 0.003 ether,
            isActive: true
        });
        _activeOperatorIds.push(7);
        _countryOperatorIds["LA"].push(7);

        _operators[8] = Operator({
            id: 8,
            name: "Smart Cambodia",
            region: "Cambodia",
            countryCode: "KH",
            requiredDeposit: 0.003 ether,
            isActive: true
        });
        _activeOperatorIds.push(8);
        _countryOperatorIds["KH"].push(8);

        _operators[9] = Operator({
            id: 9,
            name: "AIS Thailand",
            region: "Thailand",
            countryCode: "TH",
            requiredDeposit: 0.004 ether,
            isActive: true
        });
        _activeOperatorIds.push(9);
        _countryOperatorIds["TH"].push(9);

        _operators[10] = Operator({
            id: 10,
            name: "Maxis Malaysia",
            region: "Malaysia",
            countryCode: "MY",
            requiredDeposit: 0.004 ether,
            isActive: true
        });
        _activeOperatorIds.push(10);
        _countryOperatorIds["MY"].push(10);

        _operators[11] = Operator({
            id: 11,
            name: "Globe Philippines",
            region: "Philippines",
            countryCode: "PH",
            requiredDeposit: 0.003 ether,
            isActive: true
        });
        _activeOperatorIds.push(11);
        _countryOperatorIds["PH"].push(11);

        _nextOperatorId = 12;
    }

    function addOperator(
        string calldata name,
        string calldata region,
        string calldata countryCode,
        uint256 requiredDeposit
    ) external onlyOwner {
        uint256 operatorId = _nextOperatorId++;
        _operators[operatorId] = Operator({
            id: operatorId,
            name: name,
            region: region,
            countryCode: countryCode,
            requiredDeposit: requiredDeposit,
            isActive: true
        });
        _activeOperatorIds.push(operatorId);
        _countryOperatorIds[countryCode].push(operatorId);

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

    function getOperatorsByCountry(string calldata countryCode) external view returns (Operator[] memory) {
        uint256[] storage opIds = _countryOperatorIds[countryCode];
        uint256 count = 0;
        for (uint256 i = 0; i < opIds.length; i++) {
            if (_operators[opIds[i]].isActive) count++;
        }

        Operator[] memory result = new Operator[](count);
        uint256 idx = 0;
        for (uint256 i = 0; i < opIds.length; i++) {
            if (_operators[opIds[i]].isActive) {
                result[idx++] = _operators[opIds[i]];
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
