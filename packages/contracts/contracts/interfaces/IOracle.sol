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

    event UsageDataSubmitted(address indexed user, uint256 operatorId, uint256 dataUsage, uint256 callUsage);

    function submitUsage(
        address user,
        uint256 operatorId,
        uint256 dataUsage,
        uint256 callUsage
    ) external;

    function getLatestUsage(address user, uint256 operatorId) external view returns (uint256 dataUsage, uint256 callUsage, uint256 timestamp);
}
