// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import "./interfaces/IOracle.sol";

/// @title Oracle - 服务验证与计量预言机（v3 简化版）
contract Oracle is IOracle, OwnableUpgradeable, UUPSUpgradeable {
    address public deposit;

    /// @notice Initializer
    function initialize() public initializer {
        __Ownable_init(msg.sender);
    }

    /// @notice 设置 Deposit 合约地址
    function setDeposit(address _deposit) external onlyOwner {
        deposit = _deposit;
    }

    /// @notice 验证用户服务是否活跃（预留接口，未来接入 Chainlink/0G Compute）
    function verifyServiceActive(address user) external view returns (bool) {
        if (address(deposit) == address(0)) return false;
        (bool ok, bytes memory data) = address(deposit).staticcall(
            abi.encodeWithSignature("getLockExpiry(address)", user)
        );
        if (!ok) return false;
        (uint256 lockExpiry) = abi.decode(data, (uint256));
        return lockExpiry > 0 && block.timestamp < lockExpiry;
    }

    /// @notice 提交使用数据（预留接口，用于 0G Compute/Chainlink 接入）
    function submitUsage(
        address user,
        uint256,
        uint256,
        uint256
    ) external view {
        require(user != address(0), "Invalid user");
    }

    /// @inheritdoc UUPSUpgradeable
    function _authorizeUpgrade(address newImplementation) internal override onlyOwner {}
}
