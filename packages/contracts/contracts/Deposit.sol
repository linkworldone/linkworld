// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import "./interfaces/IDeposit.sol";
import "./interfaces/IUserRegistry.sol";
import "./interfaces/IPayment.sol";
import "./interfaces/IServiceManager.sol";

/// @title Deposit - 用户保证金管理（可升级版）
contract Deposit is IDeposit, OwnableUpgradeable, ReentrancyGuard, UUPSUpgradeable {
    IUserRegistry public userRegistry;
    IPayment public payment;
    IServiceManager public serviceManager;

    address public oracle; // 预言机地址，有权触发流量卡发放

    // ========== 流量卡相关状态变量 ==========
    uint256 public trafficCardQuota; // 流量卡额度（字节），默认 100M
    mapping(address => uint256) private _trafficCardBalances; // 用户流量卡余额
    mapping(address => uint256) private _lastWithdrawTimestamp; // 用户最后提款时间戳
    mapping(address => mapping(uint256 => bool)) private _monthlyCardIssued; // 用户×月份 → 是否已发放

    // ========== 保证金相关状态变量 ==========
    mapping(address => uint256) private _deposits;
    mapping(uint256 => uint256) private _operatorRequiredDeposit;

    modifier onlyOracle() {
        require(msg.sender == oracle, "Only oracle");
        _;
    }

    /// @inheritdoc IDeposit
    function initialize(address _userRegistry) public initializer {
        __Ownable_init(msg.sender);
        // 手动初始化 ReentrancyGuard 存储槽
        _reentrancyGuardInit();

        userRegistry = IUserRegistry(_userRegistry);
        trafficCardQuota = 100 * 1024 * 1024; // 默认 100M
    }

    // 内部 ReentrancyGuard 初始化（因为构造函数不会在代理模式中调用）
    function _reentrancyGuardInit() internal {
        _reentrancyGuardStorageSlot().getUint256Slot().value = 1; // NOT_ENTERED = 1
    }

    function setPayment(address _payment) external onlyOwner {
        payment = IPayment(_payment);
    }

    function setServiceManager(address _serviceManager) external onlyOwner {
        serviceManager = IServiceManager(_serviceManager);
    }

    function setRequiredDeposit(uint256 operatorId, uint256 amount) external onlyOwner {
        _operatorRequiredDeposit[operatorId] = amount;
    }

    function setOracle(address _oracle) external onlyOwner {
        oracle = _oracle;
    }

    function deposit() external payable {
        require(userRegistry.isRegistered(msg.sender), "Not registered");
        require(msg.value > 0, "Zero deposit");

        _deposits[msg.sender] += msg.value;
        emit DepositMade(msg.sender, msg.value);
    }

    function withdraw() external nonReentrant {
        require(_deposits[msg.sender] > 0, "No deposit");

        if (address(serviceManager) != address(0)) {
            IServiceManager.UserService memory service = serviceManager.getUserService(msg.sender);
            require(!service.isActive, "Service still active");
        }

        if (address(payment) != address(0)) {
            IPayment.Bill[] memory unpaid = payment.getUnpaidBills(msg.sender);
            require(unpaid.length == 0, "Has unpaid bills");
        }

        uint256 principal = _deposits[msg.sender];
        _deposits[msg.sender] = 0;
        _lastWithdrawTimestamp[msg.sender] = block.timestamp; // 记录提款时间

        (bool success, ) = msg.sender.call{value: principal}("");
        require(success, "Transfer failed");

        emit DepositWithdrawn(msg.sender, principal, 0);
    }

    /// @notice 管理员设置流量卡额度
    function setTrafficCardQuota(uint256 quota) external onlyOwner {
        trafficCardQuota = quota;
        emit TrafficCardQuotaUpdated(quota);
    }

    /// @notice 月末为符合条件的用户发放流量卡（由Oracle或Owner调用）
    function issueMonthlyTrafficCards(address[] calldata users) external {
        require(msg.sender == owner() || msg.sender == oracle, "Only owner or oracle");
        uint256 currentMonth = getCurrentMonth();

        for (uint256 i = 0; i < users.length; i++) {
            address user = users[i];
            // 检查是否已发放
            if (_monthlyCardIssued[user][currentMonth]) continue;
            
            // 检查本月是否提过款
            if (hasWithdrawnThisMonth(user, currentMonth)) continue;
            
            // 检查用户是否有存款
            if (_deposits[user] == 0) continue;

            // 发放流量卡
            _trafficCardBalances[user] += trafficCardQuota;
            _monthlyCardIssued[user][currentMonth] = true;
            
            emit TrafficCardIssued(user, trafficCardQuota, currentMonth);
        }
    }

    /// @notice 用户使用流量卡抵扣（由Payment合约调用）
    function useTrafficCard(address user, uint256 amount) external {
        require(msg.sender == address(payment), "Only payment contract");
        require(_trafficCardBalances[user] >= amount, "Insufficient traffic card balance");
        
        _trafficCardBalances[user] -= amount;
        emit TrafficCardUsed(user, amount);
    }

    /// @notice 查询用户流量卡余额
    function getTrafficCardBalance(address user) external view returns (uint256) {
        return _trafficCardBalances[user];
    }

    /// @notice 查询用户本月是否提过款
    function hasWithdrawnThisMonth(address user, uint256 month) public view returns (bool) {
        uint256 withdrawTime = _lastWithdrawTimestamp[user];
        if (withdrawTime == 0) return false;
        
        uint256 withdrawMonth = getMonthFromTimestamp(withdrawTime);
        return withdrawMonth == month;
    }

    /// @notice 获取当前月份（从epoch开始计算的月数）
    function getCurrentMonth() public view returns (uint256) {
        return getMonthFromTimestamp(block.timestamp);
    }

    /// @notice 从时间戳获取月份
    function getMonthFromTimestamp(uint256 timestamp) internal pure returns (uint256) {
        uint256 secondsPerMonth = 30 days;
        return timestamp / secondsPerMonth;
    }

    function getDepositAmount(address user) external view returns (uint256) {
        return _deposits[user];
    }

    function getRequiredDeposit(uint256 operatorId) external view returns (uint256) {
        return _operatorRequiredDeposit[operatorId];
    }

    /// @inheritdoc UUPSUpgradeable
    function _authorizeUpgrade(address newImplementation) internal override onlyOwner {}
}