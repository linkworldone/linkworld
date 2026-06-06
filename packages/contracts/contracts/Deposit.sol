// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import "./interfaces/IDeposit.sol";
import "./interfaces/IUserRegistry.sol";
import "./interfaces/ITrafficCardNFT.sol";
import "@openzeppelin/contracts/utils/Strings.sol";

/// @title Deposit - 用户保证金管理（v3: 锁仓→mint流量卡→30天服务）
contract Deposit is IDeposit, OwnableUpgradeable, UUPSUpgradeable {
    IUserRegistry public userRegistry;
    ITrafficCardNFT public trafficCardNFT;

    address public oracle;

    mapping(address => uint256) private _deposits;
    mapping(address => uint256) private _lockExpiry;

    bool private _withdrawLocked;

    modifier onlyOracle() {
        require(msg.sender == oracle, "Only oracle");
        _;
    }

    modifier nonReentrant() {
        require(!_withdrawLocked, "Reentrant");
        _withdrawLocked = true;
        _;
        _withdrawLocked = false;
    }

    /// @notice 用户存入保证金（锁仓30天）
    function deposit() external payable {
        require(msg.value > 0, "Zero deposit");
        require(userRegistry.isRegistered(msg.sender), "Not registered");

        _deposits[msg.sender] += msg.value;
        if (_lockExpiry[msg.sender] < block.timestamp) {
            _lockExpiry[msg.sender] = block.timestamp + 30 days;
        } else {
            _lockExpiry[msg.sender] += 30 days;
        }

        emit DepositMade(msg.sender, msg.value);
    }

    /// @notice 提取保证金（锁仓期满 + 无未销毁卡）
    function withdraw() external nonReentrant {
        require(block.timestamp >= _lockExpiry[msg.sender], "Lock not expired");
        require(_deposits[msg.sender] > 0, "No deposit");
        require(trafficCardNFT.getUserCardCount(msg.sender) == 0, "Card not burned");

        uint256 principal = _deposits[msg.sender];
        _deposits[msg.sender] = 0;
        _lockExpiry[msg.sender] = 0;

        (bool success, ) = msg.sender.call{value: principal}("");
        require(success, "Transfer failed");

        emit DepositWithdrawn(msg.sender, principal, 0);
    }

    /// @notice 为用户 mint 流量卡（锁仓到期后调用，仅 owner）
    function mintTrafficCard(address user) external onlyOwner returns (uint256) {
        require(block.timestamp >= _lockExpiry[user], "Lock not expired");
        require(_deposits[user] > 0, "No deposit");
        require(trafficCardNFT.getUserCardCount(user) == 0, "Card already active");

        uint256 dataAmount = _deposits[user] / 100000;
        uint256 tokenId = trafficCardNFT.mint(user, dataAmount, generateTokenURI(user));

        emit TrafficCardMinted(user, tokenId, dataAmount);
        return tokenId;
    }

    /// @notice 生成 tokenURI
    function generateTokenURI(address user) public view returns (string memory) {
        return string(abi.encodePacked(
            "https://api.linkworld.io/traffic-card/",
            Strings.toHexString(uint256(uint160(user))),
            "-",
            Strings.toString(block.timestamp / 1 days)
        ));
    }

    /// @notice 设置预言机地址
    function setOracle(address _oracle) external onlyOwner {
        oracle = _oracle;
    }

    /// @notice 设置 TrafficCardNFT 合约地址
    function setTrafficCardNFT(address _trafficCardNFT) external onlyOwner {
        trafficCardNFT = ITrafficCardNFT(_trafficCardNFT);
    }

    /// @notice 获取用户保证金余额
    function getDepositAmount(address user) external view returns (uint256) {
        return _deposits[user];
    }

    /// @notice 获取用户锁仓到期时间
    function getLockExpiry(address user) external view returns (uint256) {
        return _lockExpiry[user];
    }

    /// @notice 初始化
        function initialize(address _userRegistry) public initializer {
            __Ownable_init(msg.sender);
            userRegistry = IUserRegistry(_userRegistry);
        }

    /// @inheritdoc UUPSUpgradeable
    function _authorizeUpgrade(address newImplementation) internal override onlyOwner {}
}
