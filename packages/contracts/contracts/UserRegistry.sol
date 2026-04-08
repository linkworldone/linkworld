// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "./interfaces/IUserRegistry.sol";

/// @title UserRegistry - 用户注册与 NFT 身份凭证
contract UserRegistry is IUserRegistry, ERC721, Ownable {
    uint256 private _nextTokenId;
    mapping(address => UserInfo) private _users;

    constructor() ERC721("LinkWorld Identity", "LWID") Ownable(msg.sender) {}

    function register(string calldata email) external {
        require(!_users[msg.sender].isActive, "Already registered");

        uint256 tokenId = _nextTokenId++;
        _safeMint(msg.sender, tokenId);

        _users[msg.sender] = UserInfo({
            wallet: msg.sender,
            email: email,
            tokenId: tokenId,
            isActive: true,
            registeredAt: block.timestamp
        });

        emit UserRegistered(msg.sender, email, tokenId);
    }

    function getUserInfo(address wallet) external view returns (UserInfo memory) {
        require(_users[wallet].isActive, "User not found");
        return _users[wallet];
    }

    function isRegistered(address wallet) external view returns (bool) {
        return _users[wallet].isActive;
    }
}
