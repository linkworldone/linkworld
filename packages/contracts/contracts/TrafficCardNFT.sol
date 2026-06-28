// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/token/ERC721/extensions/ERC721URIStorageUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import "./interfaces/ITrafficCardNFT.sol";
import "@openzeppelin/contracts/utils/Strings.sol";
using Strings for uint256;

contract TrafficCardNFT is ITrafficCardNFT, OwnableUpgradeable, ERC721URIStorageUpgradeable, UUPSUpgradeable {
    uint256 private _nextTokenId;
    mapping(uint256 => CardInfo) private _cardInfo;
    mapping(address => uint256) private _userCardCount;
    address public depositContract;

    struct DeductionCredit {
        uint256 amount;
        uint256 expiresAt;
    }
    mapping(address => DeductionCredit) private _deductionCredits;
    uint256 public constant DEDUCTION_VALIDITY = 30 days;

    // eSIM redemption - new variables appended at end for UUPS upgrade compatibility
    mapping(uint256 => string) private _activationCodes;
    string private _smDpAddress;

    function setSmDpAddress(string calldata _addr) external onlyOwner {
        _smDpAddress = _addr;
    }

    function smDpAddress() external view returns (string memory) {
        return _smDpAddress;
    }

    function _generateActivationCode(uint256 tokenId) private view returns (string memory) {
        return Strings.toString(uint256(keccak256(abi.encodePacked(tokenId, block.timestamp))) % 1000000);
    }

    modifier onlyDeposit() {
        require(msg.sender == depositContract, "Not deposit");
        _;
    }

    function initialize() public initializer {
        __Ownable_init(msg.sender);
        __ERC721_init("LinkWorld Traffic Card", "LWTC");
        __ERC721URIStorage_init();
    }

    function setDepositContract(address _deposit) external onlyOwner {
        depositContract = _deposit;
    }

    function mint(address to, uint256 dataAmount, string calldata tokenURI_) external onlyOwner returns (uint256) {
        require(to != address(0), "Invalid address");
        require(dataAmount > 0, "Zero data amount");
        uint256 tokenId = _nextTokenId++;
        _safeMint(to, tokenId);
        _setTokenURI(tokenId, tokenURI_);
        _cardInfo[tokenId] = CardInfo({ dataAmount: dataAmount, createdAt: block.timestamp, isDestroyed: false });
        _userCardCount[to]++;
        emit CardMinted(to, tokenId, dataAmount);
        return tokenId;
    }

    function mintBatch(address[] calldata to, uint256[] calldata dataAmounts, string[] calldata tokenURIs_)
        external onlyOwner returns (uint256[] memory)
    {
        require(to.length == dataAmounts.length && to.length == tokenURIs_.length, "Length mismatch");
        uint256[] memory tokenIds = new uint256[](to.length);
        for (uint256 i = 0; i < to.length; i++) {
            tokenIds[i] = this.mint(to[i], dataAmounts[i], tokenURIs_[i]);
        }
        return tokenIds;
    }

    function burn(uint256 tokenId) external {
        string memory activationCode = _generateActivationCode(tokenId);
        _activationCodes[tokenId] = activationCode;
        _redeem(tokenId);
        uint256[] memory ids = new uint256[](1);
        ids[0] = tokenId;
        emit SimRedeemed(msg.sender, 1, ids);
        emit ESimRedeemed(msg.sender, tokenId, activationCode, _smDpAddress);
    }

    function redeemForSim(uint256[] calldata tokenIds) external returns (uint256 daysCount) {
        require(tokenIds.length > 0, "No cards");
        for (uint256 i = 0; i < tokenIds.length; i++) {
            string memory activationCode = _generateActivationCode(tokenIds[i]);
            _activationCodes[tokenIds[i]] = activationCode;
            _redeem(tokenIds[i]);
            emit ESimRedeemed(msg.sender, tokenIds[i], activationCode, _smDpAddress);
        }
        daysCount = tokenIds.length;
        emit SimRedeemed(msg.sender, daysCount, tokenIds);
        return daysCount;
    }

    function getActivationCode(uint256 tokenId) external view returns (string memory) {
        return _activationCodes[tokenId];
    }

    function _redeem(uint256 tokenId) internal {
        address tokenOwner = ownerOf(tokenId);
        require(_isAuthorized(tokenOwner, msg.sender, tokenId), "Not owner or approved");
        require(!_cardInfo[tokenId].isDestroyed, "Card already destroyed");
        _cardInfo[tokenId].isDestroyed = true;
        _userCardCount[tokenOwner]--;
        _burn(tokenId);
        emit CardDestroyed(tokenOwner, tokenId, _cardInfo[tokenId].dataAmount);
    }

    function getCardInfo(uint256 tokenId) external view returns (CardInfo memory) {
        require(_cardInfo[tokenId].createdAt > 0, "Card not found");
        return _cardInfo[tokenId];
    }

    function getUserCardCount(address user) external view returns (uint256) {
        return _userCardCount[user];
    }

    function tokenURI(uint256 tokenId) public view override(ERC721URIStorageUpgradeable) returns (string memory) {
        return super.tokenURI(tokenId);
    }

    function supportsInterface(bytes4 interfaceId)
        public view override(ERC721URIStorageUpgradeable) returns (bool)
    {
        return super.supportsInterface(interfaceId);
    }

    function _setTokenURI(uint256 tokenId, string memory tokenURI_) internal override {
        super._setTokenURI(tokenId, tokenURI_);
    }

    function _authorizeUpgrade(address newImplementation) internal override onlyOwner {}
}
