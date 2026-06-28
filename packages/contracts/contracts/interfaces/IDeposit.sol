// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

interface IDeposit {
    /// @notice 逐笔独立锁仓单元（每次 deposit 产生一笔，互不影响）
    struct Tranche {
        uint256 amount;
        uint256 unlockAt;
        bool withdrawn;
    }

    event DepositMade(address indexed user, uint256 amount);
    event DepositWithdrawn(address indexed user, uint256 amount, uint256 trancheIndex);
    event TrafficCardMinted(address indexed user, uint256 tokenId, uint256 dataAmount);

    function deposit(uint256 amount) external;
    function withdraw(uint256 trancheIndex) external;

    function getTranches(address user) external view returns (Tranche[] memory);
    function getTrancheCount(address user) external view returns (uint256);
    function getDepositAmount(address user) external view returns (uint256);
    function getLockExpiry(address user) external view returns (uint256);

    function setOracle(address _oracle) external;
    function setTrafficCardNFT(address _trafficCardNFT) external;
    function setSmDpAddress(string calldata _addr) external;
    function upgradeTrafficCardNFT(address _newImpl) external;
}
