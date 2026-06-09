// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/token/ERC721/extensions/ERC721URIStorageUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import "./interfaces/ITrafficCardNFT.sol";

/// @title TrafficCardNFT - 流量卡凭证（v3: mint→持有→burn→30天服务）
contract TrafficCardNFT is ITrafficCardNFT, OwnableUpgradeable, ERC721URIStorageUpgradeable, UUPSUpgradeable {
    uint256 private _nextTokenId;

    mapping(uint256 => CardInfo) private _cardInfo;
    mapping(address => uint256) private _userCardCount;

    address public depositContract;

    uint256 public constant DEDUCTION_VALIDITY = 30 days;

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

    /// @notice 用户销毁 NFT 激活服务（仅 owner 可调）
    function burn(uint256 tokenId) external {
        require(_isAuthorized(msg.sender, address(0), tokenId), "Not owner or approved");
        require(!_cardInfo[tokenId].isDestroyed, "Card already destroyed");

        _cardInfo[tokenId].isDestroyed = true;
        _userCardCount[msg.sender]--;

        _burn(tokenId);

        emit CardDestroyed(msg.sender, tokenId, _cardInfo[tokenId].dataAmount);
        emit ServiceActivated(msg.sender, block.timestamp + 30 days);
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