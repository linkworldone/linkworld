// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

import "@openzeppelin/contracts-upgradeable/token/ERC721/extensions/ERC721URIStorageUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import "./interfaces/IUserRegistry.sol";

/// @title UserRegistry - 用户注册与 NFT 身份凭证（可升级版）
contract UserRegistry is IUserRegistry, ERC721URIStorageUpgradeable, OwnableUpgradeable, UUPSUpgradeable {
    uint256 private _nextTokenId;
    mapping(address => UserInfo) private _users;

    /// @notice Initializer
    function initialize() public initializer {
        __ERC721_init("LinkWorld Identity", "LWID");
        __ERC721URIStorage_init_unchained();
        __Ownable_init(msg.sender);
        // __UUPSUpgradeable_init() not needed
    }

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

    /// @inheritdoc UUPSUpgradeable
    function _authorizeUpgrade(address newImplementation) internal override onlyOwner {}
}
