// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

interface ITrafficCardNFT {
    struct CardInfo {
        uint256 dataAmount;
        uint256 createdAt;
        bool isDestroyed;
        uint256 destroyedAt;
    }
    
    event CardMinted(address indexed user, uint256 tokenId, uint256 dataAmount);
    event CardDestroyed(address indexed user, uint256 tokenId, uint256 creditAmount);
    event CreditUsed(address indexed user, uint256 amount);
    event CreditExpired(address indexed user, uint256 amount);
    
    function mint(address to, uint256 dataAmount, string calldata tokenURI) external;
    function mintBatch(
        address[] calldata to,
        uint256[] calldata dataAmounts,
        string[] calldata tokenURIs
    ) external;
    function burn(uint256 tokenId) external;
    function useCredit(address user, uint256 amount) external;
    function getAvailableCredit(address user) external view returns (uint256);
    function getCreditExpiry(address user) external view returns (uint256);
    function getCardInfo(uint256 tokenId) external view returns (CardInfo memory);
    function getUserCardCount(address user) external view returns (uint256);
}