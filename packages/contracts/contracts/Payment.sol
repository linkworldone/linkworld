// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuardTransient.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import "./interfaces/IPayment.sol";
import "./interfaces/IFeeManager.sol";
import "./interfaces/IServiceManager.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";

/// @title Payment - 支付结算（服务商费用 + 平台手续费 + 流量卡抵扣，可升级版，ERC20 USDT 链上直分）
contract Payment is IPayment, OwnableUpgradeable, ReentrancyGuardTransient, UUPSUpgradeable {
    using SafeERC20 for IERC20;

    IFeeManager public feeManager;
    address public platformWallet;
    address public oracle;

    uint256 private _nextBillId;
    mapping(uint256 => Bill) private _bills;
    mapping(address => uint256[]) private _userBillIds;

    /// @notice 资金通道代币（USDT，6 位精度）。仅支持标准 ERC20，禁 fee-on-transfer。
    IERC20 public usdt;
    /// @notice 运营商目录（查 operator.paymentAddress 做链上分账）
    IServiceManager public serviceManager;

    modifier onlyOracle() {
        require(msg.sender == oracle, "Only oracle");
        _;
    }

    /// @notice Initializer
    /// @param _feeManager 手续费合约
    /// @param _platformWallet 平台手续费收款地址
    /// @param _usdt USDT 代币合约（initialize 注入，避免未设时 safeTransferFrom(address(0)) panic）
    /// @param _serviceManager 运营商目录（查 operator.paymentAddress 做分账）
    function initialize(
        address _feeManager,
        address _platformWallet,
        address _usdt,
        address _serviceManager
    ) public initializer {
        __Ownable_init(msg.sender);
        feeManager = IFeeManager(_feeManager);
        platformWallet = _platformWallet;
        usdt = IERC20(_usdt);
        serviceManager = IServiceManager(_serviceManager);
    }

    function setOracle(address _oracle) external onlyOwner {
        oracle = _oracle;
    }

    /// @notice 设置平台钱包
    function setPlatformWallet(address _platformWallet) external onlyOwner {
        platformWallet = _platformWallet;
    }

    /// @notice 设置手续费合约
    function setFeeManager(address _feeManager) external onlyOwner {
        feeManager = IFeeManager(_feeManager);
    }

    /// @notice 设置运营商目录合约
    function setServiceManager(address _serviceManager) external onlyOwner {
        serviceManager = IServiceManager(_serviceManager);
    }

    /// @notice 创建服务购买账单（仅 oracle，B2）
    /// @dev amount 为后端/Oracle 链下按资费算好的 USDT 金额（v2-A，合约不做 usage 求和）
    function createBill(address user, uint256 operatorId, uint256 amount) external onlyOracle {
        require(amount > 0, "Zero amount");
        // fail-fast：避免生成永远付不了的账单（operator 分账地址未设时直接拒绝）
        require(
            serviceManager.getOperator(operatorId).paymentAddress != address(0),
            "Operator payout unset"
        );

        uint256 platformFee = feeManager.calculateFee(amount);
        uint256 billId = _nextBillId++;

        _bills[billId] = Bill({
            id: billId,
            user: user,
            operatorId: operatorId,
            amount: amount,
            platformFee: platformFee,
            createdAt: block.timestamp,
            isPaid: false
        });
        _userBillIds[user].push(billId);

        emit BillCreated(billId, user, amount + platformFee, platformFee);
    }

    /// @notice 用户支付账单（ERC20 USDT 链上直分）
    /// @dev 用户须先 approve(payment, amount+fee)；两段 safeTransferFrom 各自从 user 拉款，
    ///      合约不暂存资金（降低重入面）。CEI：先置 isPaid 再转账。0-fee 跳过第二段。
    function payBill(uint256 billId) external nonReentrant {
        Bill storage bill = _bills[billId];
        require(!bill.isPaid, "Already paid");
        require(bill.user == msg.sender, "Not your bill");

        address operatorPayout = serviceManager.getOperator(bill.operatorId).paymentAddress;
        require(operatorPayout != address(0), "Operator payout unset");

        // CEI：状态先改
        bill.isPaid = true;

        uint256 total = bill.amount + bill.platformFee;

        // 主体分账：user -> operator.paymentAddress
        usdt.safeTransferFrom(msg.sender, operatorPayout, bill.amount);

        // 手续费：user -> platformWallet（0-fee 跳过）
        if (bill.platformFee > 0) {
            usdt.safeTransferFrom(msg.sender, platformWallet, bill.platformFee);
        }

        emit BillPaid(billId, msg.sender, total, bill.amount);
    }

    function getUserBills(address user) external view returns (Bill[] memory) {
        uint256[] memory ids = _userBillIds[user];
        Bill[] memory bills = new Bill[](ids.length);
        for (uint256 i = 0; i < ids.length; i++) {
            bills[i] = _bills[ids[i]];
        }
        return bills;
    }

    function getUnpaidBills(address user) public view returns (Bill[] memory) {
        uint256[] memory ids = _userBillIds[user];
        uint256 count = 0;
        for (uint256 i = 0; i < ids.length; i++) {
            if (!_bills[ids[i]].isPaid) count++;
        }

        Bill[] memory result = new Bill[](count);
        uint256 idx = 0;
        for (uint256 i = 0; i < ids.length; i++) {
            if (!_bills[ids[i]].isPaid) {
                result[idx++] = _bills[ids[i]];
            }
        }
        return result;
    }

    /// @notice 流量卡抵扣账单（受限桩，v2-B）
    /// @dev 本轮仅做权限与账单存在性校验，不转移任何资金；真实抵扣语义冻结到后续 Round
    function applyTrafficCardToBill(uint256 billId) external onlyOracle {
        require(_bills[billId].user != address(0), "Bill not found");
        emit TrafficCardApplied(billId);
    }

    /// @inheritdoc UUPSUpgradeable
    function _authorizeUpgrade(address newImplementation) internal override onlyOwner {}
}