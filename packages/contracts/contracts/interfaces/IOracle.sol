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

    event ServiceVerified(address indexed user, bool isActive);
    event UsageDataSubmitted(address indexed user, uint256 indexed operatorId, uint256 dataUsage, uint256 callUsage);

    function verifyServiceActive(address user) external view returns (bool);
    function submitUsage(
        address user,
        uint256 operatorId,
        uint256 dataUsage,
        uint256 callUsage
    ) external;
}
