// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

interface IUserRegistry {
    struct UserInfo {
        address wallet;
        string email;
        uint256 tokenId; // NFT token ID
        bool isActive;
        uint256 registeredAt;
    }

    event UserRegistered(address indexed wallet, string email, uint256 tokenId);
    event UserDeactivated(address indexed wallet);

    function register(string calldata email) external;
    function getUserInfo(address wallet) external view returns (UserInfo memory);
    function isRegistered(address wallet) external view returns (bool);
}
