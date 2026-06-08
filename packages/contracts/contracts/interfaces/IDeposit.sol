// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

interface IDeposit {
    event DepositMade(address indexed user, uint256 amount);
    event DepositWithdrawn(address indexed user, uint256 principal, uint256 interest);
    event TrafficCardMinted(address indexed user, uint256 tokenId, uint256 dataAmount);

    function deposit() external payable;
    function withdraw() external;
    function getDepositAmount(address user) external view returns (uint256);
    function getLockExpiry(address user) external view returns (uint256);
    function mintTrafficCard(address user) external returns (uint256);
    function issueMonthlyTrafficCards(address[] calldata users) external;
    function setOracle(address _oracle) external;
    function setTrafficCardNFT(address _trafficCardNFT) external;
}
