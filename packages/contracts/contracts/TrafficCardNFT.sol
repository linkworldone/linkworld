// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/token/ERC721/extensions/ERC721URIStorageUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import "./interfaces/ITrafficCardNFT.sol";

/// @title TrafficCardNFT - 流量卡凭证（v4: mint→持有→burn/redeemForSim→兑换 SIM 天数，每卡 1 天）
contract TrafficCardNFT is ITrafficCardNFT, OwnableUpgradeable, ERC721URIStorageUpgradeable, UUPSUpgradeable {
    uint256 private _nextTokenId;

    mapping(uint256 => CardInfo) private _cardInfo;
    mapping(address => uint256) private _userCardCount;

    address public depositContract;

    /// @notice ERC721 initializer
    function initialize() public initializer {
        __Ownable_init(msg.sender);
        __ERC721_init("LinkWorld Traffic Card", "LWTC");
        __ERC721URIStorage_init();
    }

    /// @notice 设置 Deposit 合约地址（允许其 mint）
    function setDepositContract(address _deposit) external onlyOwner {
        depositContract = _deposit;
    }

    /// @notice Owner 直接 mint
    function mint(address to, uint256 dataAmount, string calldata tokenURI_) external onlyOwner returns (uint256) {
        require(to != address(0), "Invalid address");
        require(dataAmount > 0, "Zero data amount");

        uint256 tokenId = _nextTokenId++;
        _safeMint(to, tokenId);
        _setTokenURI(tokenId, tokenURI_);

        _cardInfo[tokenId] = CardInfo({
            dataAmount: dataAmount,
            createdAt: block.timestamp,
            isDestroyed: false
        });

        _userCardCount[to]++;
        emit CardMinted(to, tokenId, dataAmount);
        return tokenId;
    }

    /// @notice 批量 mint
    /// @dev ⚠️ 勿用于自动发卡：内部 this.mint(...) 是外部调用，以本合约为 msg.sender 会撞 mint 的
    ///      onlyOwner（owner 已转给 Deposit）→ revert。自动发卡走 Deposit.issueMonthlyTrafficCards
    ///      逐张 trafficCardNFT.mint(...)。本函数本轮未接入发卡链路（dead code，保留待后续 Round）。
    function mintBatch(
        address[] calldata to,
        uint256[] calldata dataAmounts,
        string[] calldata tokenURIs_
    )
        external
        onlyOwner
        returns (uint256[] memory)
    {
        require(to.length == dataAmounts.length && to.length == tokenURIs_.length, "Length mismatch");
        uint256[] memory tokenIds = new uint256[](to.length);
        for (uint256 i = 0; i < to.length; i++) {
            tokenIds[i] = this.mint(to[i], dataAmounts[i], tokenURIs_[i]);
        }
        return tokenIds;
    }

    /// @notice 用户销毁单张流量卡兑换 1 天 SIM（与 redeemForSim 同语义的便捷入口）
    /// @dev 每张卡 = 1 天，emit SimRedeemed(user, 1, [tokenId])。链上仅销毁 + emit 天数，SIM 走链下。
    function burn(uint256 tokenId) external {
        _redeem(tokenId);
        uint256[] memory ids = new uint256[](1);
        ids[0] = tokenId;
        emit SimRedeemed(msg.sender, 1, ids);
    }

    /// @notice 批量销毁流量卡兑换 SIM 天数（销毁 N 张 = N 天，天数累加）
    /// @dev 链上仅负责销毁卡 + emit 天数（= 卡数）；SIM 本身走链下后端。
    function redeemForSim(uint256[] calldata tokenIds) external returns (uint256 daysCount) {
        require(tokenIds.length > 0, "No cards");
        for (uint256 i = 0; i < tokenIds.length; i++) {
            _redeem(tokenIds[i]);
        }
        daysCount = tokenIds.length; // 每张 = 1 天
        emit SimRedeemed(msg.sender, daysCount, tokenIds);
        return daysCount;
    }

    /// @notice 单张销毁的内部逻辑（CEI + 计数顺序）
    /// @dev _isAuthorized(owner, spender, id)：owner 必须是实际持有者，spender 为 msg.sender。
    ///      计数挂在实际持有者名下，而非 msg.sender（支持被授权方代销毁）。
    function _redeem(uint256 tokenId) internal {
        address tokenOwner = ownerOf(tokenId); // tokenId 不存在时 revert（ERC721NonexistentToken）
        require(_isAuthorized(tokenOwner, msg.sender, tokenId), "Not owner or approved");
        require(!_cardInfo[tokenId].isDestroyed, "Card already destroyed");

        _cardInfo[tokenId].isDestroyed = true;
        _userCardCount[tokenOwner]--;

        _burn(tokenId);

        emit CardDestroyed(tokenOwner, tokenId, _cardInfo[tokenId].dataAmount);
    }

    /// @notice 查询卡片信息
    function getCardInfo(uint256 tokenId) external view returns (CardInfo memory) {
        require(tokenId < _nextTokenId, "Card not found");
        return _cardInfo[tokenId];
    }

    /// @notice 查询用户持有的卡片数量
    function getUserCardCount(address user) external view returns (uint256) {
        return _userCardCount[user];
    }

    function tokenURI(uint256 tokenId) public view override(ERC721URIStorageUpgradeable) returns (string memory) {
        return super.tokenURI(tokenId);
    }

    function supportsInterface(bytes4 interfaceId)
        public
        view
        override(ERC721URIStorageUpgradeable)
        returns (bool)
    {
        return super.supportsInterface(interfaceId);
    }

    /// @notice 内部设置 tokenURI
    function _setTokenURI(uint256 tokenId, string memory tokenURI_) internal override {
        super._setTokenURI(tokenId, tokenURI_);
    }

    /// @inheritdoc UUPSUpgradeable
    function _authorizeUpgrade(address newImplementation) internal override onlyOwner {}
}