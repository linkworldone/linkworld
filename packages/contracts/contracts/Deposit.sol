// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import "./interfaces/IDeposit.sol";
import "./interfaces/IUserRegistry.sol";
import "./interfaces/IPayment.sol";
import "./interfaces/IServiceManager.sol";

/// @title Deposit - 用户保证金管理
contract Deposit is IDeposit, Ownable, ReentrancyGuard {
    IUserRegistry public userRegistry;
    IPayment public payment;
    IServiceManager public serviceManager;

    // ========== 利息相关状态变量 ==========
    uint256 public interestRate; // 当前年化利率（基点表示，如 500 = 5%）
    uint256 public lastRateUpdate; // 上次利率更新时间
    uint256 public constant SECONDS_PER_YEAR = 365 days;
    uint256 public constant RATE_DENOMINATOR = 10000;

    mapping(address => uint256) private _deposits;
    mapping(uint256 => uint256) private _operatorRequiredDeposit;
    
    // 记录用户存款时间戳
    mapping(address => uint256) private _depositTimestamps;

    constructor(address _userRegistry) Ownable(msg.sender) {
        userRegistry = IUserRegistry(_userRegistry);
        lastRateUpdate = block.timestamp;
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

    function deposit() external payable {
        require(userRegistry.isRegistered(msg.sender), "Not registered");
        require(msg.value > 0, "Zero deposit");

        _deposits[msg.sender] += msg.value;
        
        // 如果是首次存款，记录时间戳
        if (_depositTimestamps[msg.sender] == 0) {
            _depositTimestamps[msg.sender] = block.timestamp;
        }
        
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
        uint256 interest = calculateInterest(msg.sender);
        uint256 total = principal + interest;

        _deposits[msg.sender] = 0;
        _depositTimestamps[msg.sender] = 0; // 重置存款时间戳

        (bool success, ) = msg.sender.call{value: total}("");
        require(success, "Transfer failed");

        emit DepositWithdrawn(msg.sender, principal, interest);
    }

    /// @notice 更新利率（由预言机或管理员调用）
    function updateInterestRate(uint256 newRate) external onlyOwner {
        interestRate = newRate;
        lastRateUpdate = block.timestamp;
        emit InterestRateUpdated(newRate);
    }

    /// @notice 计算用户应得利息
    function calculateInterest(address user) public view returns (uint256) {
        uint256 deposit = _deposits[user];
        if (deposit == 0 || interestRate == 0) return 0;
        
        uint256 timeElapsed = block.timestamp - _depositTimestamps[user];
        if (timeElapsed == 0) return 0;
        
        // 利息 = 本金 × 利率 × 时间（年）
        return (deposit * interestRate * timeElapsed) / (RATE_DENOMINATOR * SECONDS_PER_YEAR);
    }

    /// @notice 查看用户当前可提取的利息
    function getInterestAmount(address user) external view returns (uint256) {
        return calculateInterest(user);
    }

    function getDepositAmount(address user) external view returns (uint256) {
        return _deposits[user];
    }

    function getRequiredDeposit(uint256 operatorId) external view returns (uint256) {
        return _operatorRequiredDeposit[operatorId];
    }
}