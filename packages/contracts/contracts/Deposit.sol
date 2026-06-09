// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuardTransient.sol";
import "./interfaces/IDeposit.sol";
import "./interfaces/IUserRegistry.sol";
import "./interfaces/ITrafficCardNFT.sol";
import "@openzeppelin/contracts/utils/Strings.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {IERC20Metadata} from "@openzeppelin/contracts/token/ERC20/extensions/IERC20Metadata.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";

/// @title Deposit - 用户保证金管理（可升级版，ERC20 USDT）
contract Deposit is IDeposit, OwnableUpgradeable, ReentrancyGuardTransient, UUPSUpgradeable {
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

    /// @notice 为用户 mint 流量卡（锁仓到期后调用，仅 owner 手动发卡入口）
    /// @dev onlyOwner 薄壳，复用 internal _mintFor；手动发卡要求三校验全部满足（不满足 revert）。
    function mintTrafficCard(address user) external onlyOwner nonReentrant returns (uint256) {
        require(_canMint(user), "Mint conditions not met");
        return _mintFor(user);
    }

    /// @notice 发卡三校验：锁仓到期 && 有存款 && 当前无卡（幂等）。
    function _canMint(address user) internal view returns (bool) {
        return block.timestamp >= _lockExpiry[user]
            && _deposits[user] > 0
            && trafficCardNFT.getUserCardCount(user) == 0;
    }

    /// @notice 内部发卡（B3）：固定额度 trafficCardQuota（v2-C，与存款额/精度解耦）。
    /// @dev 被 mintTrafficCard(onlyOwner) 与 issueMonthlyTrafficCards(onlyOracle) 复用，逻辑单一。
    ///      调用方负责权限与 nonReentrant；TrafficCardNFT.mint 内 _userCardCount++ 早于 _safeMint 回调（CEI）。
    function _mintFor(address user) internal returns (uint256) {
        uint256 dataAmount = trafficCardQuota;
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

    /// @notice 月度批量自动发卡（仅 oracle 调）。
    /// @dev B3：走独立 internal _mintFor（不经 onlyOwner mintTrafficCard，避免 revert）；
    ///      每个 user 三校验失败则 continue 跳过（不 revert 整批）；幂等靠 getUserCardCount==0。
    ///      A1：nonReentrant（_mintFor→nft.mint→_safeMint 有 onERC721Received 回调入口）。
    ///      消歧：禁用 NFT.mintBatch（其内 this.mint 外部调用会撞 onlyOwner），改逐张 mint。
    function issueMonthlyTrafficCards(address[] calldata users) external nonReentrant {
        require(msg.sender == oracle, "Only oracle");

        for (uint256 i = 0; i < users.length; i++) {
            address user = users[i];
            if (!_canMint(user)) {
                continue;
            }
            _mintFor(user);
        }
    }

    /// @inheritdoc UUPSUpgradeable
    function _authorizeUpgrade(address newImplementation) internal override onlyOwner {}
}