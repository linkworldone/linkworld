// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import "./interfaces/IDeposit.sol";
import "./interfaces/IUserRegistry.sol";
import "./interfaces/ITrafficCardNFT.sol";
import "@openzeppelin/contracts/utils/Strings.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {IERC20Metadata} from "@openzeppelin/contracts/token/ERC20/extensions/IERC20Metadata.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";

/// @title Deposit - 用户保证金管理（可升级版，ERC20 USDT）
contract Deposit is IDeposit, OwnableUpgradeable, UUPSUpgradeable {
    using SafeERC20 for IERC20;

    IUserRegistry public userRegistry;
    ITrafficCardNFT public trafficCardNFT;

    address public oracle;
    uint256 public trafficCardQuota;

    mapping(address => uint256) private _deposits;
    mapping(address => uint256) private _lockExpiry;
    mapping(uint256 => uint256) private _operatorRequiredDeposit;

    /// @notice 资金通道代币（USDT，6 位精度）。仅支持标准 ERC20，禁 fee-on-transfer。
    IERC20 public usdt;

    /// @notice Initializer
    /// @param _userRegistry 用户注册合约
    /// @param _usdt USDT 代币合约（initialize 注入，避免未设时 safeTransferFrom(address(0)) panic）
    function initialize(address _userRegistry, address _usdt) public initializer {
        __Ownable_init(msg.sender);
        userRegistry = IUserRegistry(_userRegistry);
        usdt = IERC20(_usdt);
        trafficCardQuota = 100 * 1024 * 1024; // 默认 100M
    }

    /// @notice 设置预言机地址
    function setOracle(address _oracle) external onlyOwner {
        oracle = _oracle;
    }

    /// @notice 设置 TrafficCardNFT 合约地址
    function setTrafficCardNFT(address _trafficCardNFT) external onlyOwner {
        trafficCardNFT = ITrafficCardNFT(_trafficCardNFT);
    }

    /// @notice 用户存入保证金（USDT，锁仓30天）。需先 approve(deposit, amount)。
    /// @param amount 存入的 USDT 数量（最小单位，需 ≥ 10 USDT）
    function deposit(uint256 amount) external {
        require(userRegistry.isRegistered(msg.sender), "Not registered");
        require(amount >= 10 * 10 ** IERC20Metadata(address(usdt)).decimals(), "Below min deposit");

        usdt.safeTransferFrom(msg.sender, address(this), amount);

        _deposits[msg.sender] += amount;
        if (_lockExpiry[msg.sender] < block.timestamp) {
            _lockExpiry[msg.sender] = block.timestamp + 30 days;
        } else {
            _lockExpiry[msg.sender] += 30 days;
        }

        emit DepositMade(msg.sender, amount);
    }

    /// @notice 提取保证金（USDT）。CEI：先清零状态再 safeTransfer。
    function withdraw() external {
        require(block.timestamp >= _lockExpiry[msg.sender], "Lock not expired");
        require(_deposits[msg.sender] > 0, "No deposit");

        uint256 principal = _deposits[msg.sender];
        _deposits[msg.sender] = 0;
        _lockExpiry[msg.sender] = 0;

        usdt.safeTransfer(msg.sender, principal);

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

    /// @notice 获取用户保证金余额
    function getDepositAmount(address user) external view returns (uint256) {
        return _deposits[user];
    }

    /// @notice 获取用户锁仓到期时间
    function getLockExpiry(address user) external view returns (uint256) {
        return _lockExpiry[user];
    }

    /// @notice 月度发放流量卡
    function issueMonthlyTrafficCards(address[] calldata users) external {
        require(msg.sender == oracle, "Only oracle");
        // Implementation for monthly traffic card issuance
    }

    /// @inheritdoc UUPSUpgradeable
    function _authorizeUpgrade(address newImplementation) internal override onlyOwner {}
}