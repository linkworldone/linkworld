// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

interface ITrafficCardNFT {
    struct CardInfo {
        uint256 dataAmount;
        uint256 createdAt;
        bool isDestroyed;
    }

    event CardMinted(address indexed user, uint256 tokenId, uint256 dataAmount);
    event CardDestroyed(address indexed user, uint256 tokenId, uint256 dataAmount);
    event SimRedeemed(address indexed user, uint256 daysCount, uint256[] tokenIds);
    event ESimRedeemed(address indexed user, uint256 tokenId, string activationCode, string smDpAddress);

    function mint(address to, uint256 dataAmount, string calldata tokenURI) external returns (uint256);
    function mintBatch(
        address[] calldata to,
        uint256[] calldata dataAmounts,
        string[] calldata tokenURIs
    ) external returns (uint256[] memory);
    function burn(uint256 tokenId) external;
    function redeemForSim(uint256[] calldata tokenIds) external returns (uint256 daysCount);
    function getCardInfo(uint256 tokenId) external view returns (CardInfo memory);
    function getUserCardCount(address user) external view returns (uint256);
    function setSmDpAddress(string calldata _addr) external;
    function smDpAddress() external view returns (string memory);
}
