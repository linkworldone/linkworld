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

/// @title Deposit - 用户充值管理（可升级版，ERC20 USDT）
/// @notice 分档充值（10/20/50/100 USDT）→ 逐笔独立锁仓 30 天 + 充值即按比例发流量卡。
contract Deposit is IDeposit, OwnableUpgradeable, ReentrancyGuardTransient, UUPSUpgradeable {
    using SafeERC20 for IERC20;

    IUserRegistry public userRegistry;
    ITrafficCardNFT public trafficCardNFT;

    address public oracle;

    /// @notice 每个用户的逐笔锁仓记录（每次 deposit push 一笔）
    mapping(address => Tranche[]) private _tranches;

    /// @notice 资金通道代币（USDT，6 位精度）。仅支持标准 ERC20，禁 fee-on-transfer。
    IERC20 public usdt;

    /// @notice 锁仓周期（每笔独立计时）
    uint256 public constant LOCK_PERIOD = 30 days;

    /// @notice 发卡基准档位（10 USDT → 1 张），张数 = amount / TIER_UNIT
    uint256 private constant TIER_MULTIPLIER = 10;

    /// @notice 无限流量哨兵值（每张流量卡的 dataAmount）
    uint256 private constant UNLIMITED_DATA = type(uint256).max;

    /// @notice Initializer
    /// @param _userRegistry 用户注册合约
    /// @param _usdt USDT 代币合约（initialize 注入，避免未设时 safeTransferFrom(address(0)) panic）
    function initialize(address _userRegistry, address _usdt) public initializer {
        __Ownable_init(msg.sender);
        userRegistry = IUserRegistry(_userRegistry);
        usdt = IERC20(_usdt);
    }

    /// @notice 设置预言机地址
    function setOracle(address _oracle) external onlyOwner {
        oracle = _oracle;
    }

    /// @notice 设置 TrafficCardNFT 合约地址
    function setTrafficCardNFT(address _trafficCardNFT) external onlyOwner {
        trafficCardNFT = ITrafficCardNFT(_trafficCardNFT);
    }

    /// @notice 设置 SM-DP 地址（转发给 TrafficCardNFT）
    function setSmDpAddress(string calldata _addr) external onlyOwner {
        trafficCardNFT.setSmDpAddress(_addr);
    }

    /// @notice 用户分档充值（仅 10/20/50/100 USDT），充值即按比例发卡，并新增一笔 30 天独立锁。
    /// @dev 需先 approve(deposit, amount)。CEI：先收款 + 记账，再 mint（mint 有 onERC721Received 回调，靠 nonReentrant 兜底）。
    /// @param amount 存入的 USDT 数量（最小单位）。按 IERC20Metadata.decimals() 折算后必须等于 10/20/50/100 档。
    function deposit(uint256 amount) external nonReentrant {
        require(userRegistry.isRegistered(msg.sender), "Not registered");

        uint256 unit = 10 ** IERC20Metadata(address(usdt)).decimals();
        require(
            amount == 10 * unit || amount == 20 * unit || amount == 50 * unit || amount == 100 * unit,
            "Invalid tier"
        );

        // 收款
        usdt.safeTransferFrom(msg.sender, address(this), amount);

        // 记一笔独立锁
        _tranches[msg.sender].push(Tranche({
            amount: amount,
            unlockAt: block.timestamp + LOCK_PERIOD,
            withdrawn: false
        }));

        emit DepositMade(msg.sender, amount);

        // 充值即发卡：张数 = amount / (10 * unit)（10→1 / 20→2 / 50→5 / 100→10）
        uint256 cardCount = amount / (TIER_MULTIPLIER * unit);
        for (uint256 i = 0; i < cardCount; i++) {
            uint256 tokenId = trafficCardNFT.mint(
                msg.sender,
                UNLIMITED_DATA,
                generateTokenURI(msg.sender, i)
            );
            emit TrafficCardMinted(msg.sender, tokenId, UNLIMITED_DATA);
        }
    }

    /// @notice 逐笔提取：仅退该笔本金，各笔互不影响。CEI：先置 withdrawn=true 再 safeTransfer。
    /// @param trancheIndex _tranches[msg.sender] 的下标
    function withdraw(uint256 trancheIndex) external nonReentrant {
        Tranche[] storage list = _tranches[msg.sender];
        require(trancheIndex < list.length, "Invalid tranche");

        Tranche storage t = list[trancheIndex];
        require(!t.withdrawn, "Already withdrawn");
        require(block.timestamp >= t.unlockAt, "Lock not expired");

        uint256 principal = t.amount;
        t.withdrawn = true;

        usdt.safeTransfer(msg.sender, principal);

        emit DepositWithdrawn(msg.sender, principal, trancheIndex);
    }

    /// @notice 生成 tokenURI（同 tx 多张带索引避免完全重复）
    function generateTokenURI(address user, uint256 index) public view returns (string memory) {
        return string(abi.encodePacked(
            "https://api.linkworld.io/traffic-card/",
            Strings.toHexString(uint256(uint160(user))),
            "-",
            Strings.toString(block.timestamp / 1 days),
            "-",
            Strings.toString(index)
        ));
    }

    /// @notice 获取用户全部锁仓笔（前端按笔展示解锁时间 + 取回按钮）
    function getTranches(address user) external view returns (Tranche[] memory) {
        return _tranches[user];
    }

    /// @notice 获取用户锁仓笔数
    function getTrancheCount(address user) external view returns (uint256) {
        return _tranches[user].length;
    }

    /// @notice 获取用户当前锁仓总额（未取回笔 amount 之和）。保持函数名给前端/后端兼容。
    function getDepositAmount(address user) external view returns (uint256) {
        Tranche[] storage list = _tranches[user];
        uint256 total;
        for (uint256 i = 0; i < list.length; i++) {
            if (!list[i].withdrawn) {
                total += list[i].amount;
            }
        }
        return total;
    }

    /// @notice 兼容旧接口：返回未取回笔中最早的解锁时间（无未取回笔则 0）。
    function getLockExpiry(address user) external view returns (uint256) {
        Tranche[] storage list = _tranches[user];
        uint256 earliest;
        for (uint256 i = 0; i < list.length; i++) {
            if (!list[i].withdrawn) {
                uint256 u = list[i].unlockAt;
                if (earliest == 0 || u < earliest) {
                    earliest = u;
                }
            }
        }
        return earliest;
    }

    /// @notice 升级 TrafficCardNFT 实现合约（仅部署者可调用）
    function upgradeTrafficCardNFT(address _newImpl) external onlyOwner {
        (bool success, ) = address(trafficCardNFT).call(
            abi.encodeWithSignature("upgradeTo(address)", _newImpl)
        );
        require(success, "Upgrade failed");
    }

    /// @inheritdoc UUPSUpgradeable
    function _authorizeUpgrade(address newImplementation) internal override onlyOwner {}
}
