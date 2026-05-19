// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

import "@openzeppelin/contracts-upgradeable/token/ERC721/ERC721Upgradeable.sol";
import "@openzeppelin/contracts-upgradeable/token/ERC721/extensions/ERC721URIStorageUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";

/// @title TrafficCardNFT - 流量卡NFT合约（可升级版）
/// @notice 用户持有此NFT可在销毁时获得流量抵扣额度
/// @dev NFT长时间有效，销毁后流量抵扣有效期为30天
contract TrafficCardNFT is ERC721Upgradeable, ERC721URIStorageUpgradeable, OwnableUpgradeable, UUPSUpgradeable {
    uint256 private _nextTokenId;
    
    /// @notice 流量卡信息
    struct CardInfo {
        uint256 dataAmount;      // 流量额度（字节）
        uint256 createdAt;       // 创建时间
        bool isDestroyed;        // 是否已销毁
        uint256 destroyedAt;     // 销毁时间（用于计算30天有效期）
    }
    
    mapping(uint256 => CardInfo) private _cardInfo;   // tokenId -> 卡信息
    mapping(address => uint256) private _userCardCount; // 用户持有数量
    
    /// @notice 销毁后的抵扣额度（有效期30天）
    struct DeductionCredit {
        uint256 amount;         // 可用抵扣额度
        uint256 expiryTime;     // 过期时间
    }
    
    mapping(address => DeductionCredit) private _deductionCredits; // 用户 -> 抵扣额度
    
    uint256 public constant DEDUCTION_VALIDITY = 30 days; // 销毁后抵扣有效期30天
    
    event CardMinted(address indexed user, uint256 tokenId, uint256 dataAmount);
    event CardDestroyed(address indexed user, uint256 tokenId, uint256 creditAmount);
    event CreditUsed(address indexed user, uint256 amount);
    event CreditExpired(address indexed user, uint256 amount);

    /// @inheritdoc ERC721Upgradeable
    function initialize() public initializer {
        __ERC721_init("LinkWorld Traffic Card", "LWTC");
        __ERC721URIStorage_init_unchained();
        __Ownable_init(msg.sender);
    }
    
    /// @notice 铸造流量卡NFT（仅Owner或授权合约调用）
    /// @param to 接收地址
    /// @param dataAmount 流量额度（字节）
    /// @param tokenURI NFT元数据URI
    function mint(
        address to,
        uint256 dataAmount,
        string calldata tokenURI
    ) external onlyOwner {
        require(to != address(0), "Invalid address");
        require(dataAmount > 0, "Zero data amount");
        
        uint256 tokenId = _nextTokenId++;
        _safeMint(to, tokenId);
        _setTokenURI(tokenId, tokenURI);
        
        _cardInfo[tokenId] = CardInfo({
            dataAmount: dataAmount,
            createdAt: block.timestamp,
            isDestroyed: false,
            destroyedAt: 0
        });
        
        _userCardCount[to]++;
        
        emit CardMinted(to, tokenId, dataAmount);
    }
    
    /// @notice 批量铸造流量卡NFT
    function mintBatch(
        address[] calldata to,
        uint256[] calldata dataAmounts,
        string[] calldata tokenURIs
    ) external onlyOwner {
        require(to.length == dataAmounts.length, "Length mismatch");
        require(to.length == tokenURIs.length, "Length mismatch");
        
        for (uint256 i = 0; i < to.length; i++) {
            mint(to[i], dataAmounts[i], tokenURIs[i]);
        }
    }
    
    /// @notice 销毁流量卡NFT并获得抵扣额度
    /// @dev NFT持有者可随时销毁获得抵扣额度
    function burn(uint256 tokenId) external {
        require(_isApprovedOrOwner(msg.sender, tokenId), "Not owner or approved");
        require(!_cardInfo[tokenId].isDestroyed, "Card already destroyed");
        
        CardInfo storage card = _cardInfo[tokenId];
        card.isDestroyed = true;
        card.destroyedAt = block.timestamp;
        
        // 将流量额度转换为抵扣额度，有效期30天
        DeductionCredit storage credit = _deductionCredits[msg.sender];
        credit.amount += card.dataAmount;
        credit.expiryTime = block.timestamp + DEDUCTION_VALIDITY;
        
        _userCardCount[msg.sender]--;
        _burn(tokenId);
        
        emit CardDestroyed(msg.sender, tokenId, card.dataAmount);
    }
    
    /// @notice 使用抵扣额度（由Payment合约调用）
    /// @param user 用户地址
    /// @param amount 使用额度
    function useCredit(address user, uint256 amount) external onlyOwner {
        DeductionCredit storage credit = _deductionCredits[user];
        
        // 检查是否过期
        require(block.timestamp <= credit.expiryTime || credit.expiryTime == 0, "Credit expired");
        require(credit.amount >= amount, "Insufficient credit");
        
        credit.amount -= amount;
        
        emit CreditUsed(user, amount);
    }
    
    /// @notice 查询用户可用抵扣额度
    function getAvailableCredit(address user) external view returns (uint256) {
        DeductionCredit storage credit = _deductionCredits[user];
        
        // 检查是否过期
        if (credit.expiryTime > 0 && block.timestamp > credit.expiryTime) {
            return 0;
        }
        
        return credit.amount;
    }
    
    /// @notice 查询抵扣额度过期时间
    function getCreditExpiry(address user) external view returns (uint256) {
        return _deductionCredits[user].expiryTime;
    }
    
    /// @notice 查询流量卡信息
    function getCardInfo(uint256 tokenId) external view returns (CardInfo memory) {
        require(_exists(tokenId) || _cardInfo[tokenId].createdAt > 0, "Card not found");
        return _cardInfo[tokenId];
    }
    
    /// @notice 查询用户持有流量卡数量
    function getUserCardCount(address user) external view returns (uint256) {
        return _userCardCount[user];
    }
    
    /// @inheritdoc ERC721URIStorageUpgradeable
    function tokenURI(uint256 tokenId) public view override(ERC721Upgradeable, ERC721URIStorageUpgradeable) returns (string memory) {
        return super.tokenURI(tokenId);
    }
    
    /// @inheritdoc ERC721Upgradeable
    function supportsInterface(bytes4 interfaceId) public view override(ERC721Upgradeable, ERC721URIStorageUpgradeable) returns (bool) {
        return super.supportsInterface(interfaceId);
    }
    
    /// @inheritdoc UUPSUpgradeable
    function _authorizeUpgrade(address newImplementation) internal override onlyOwner {}
}