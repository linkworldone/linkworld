// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import "./interfaces/IDeposit.sol";
import "./interfaces/IUserRegistry.sol";
import "./interfaces/IPayment.sol";

/// @title Deposit - 用户保证金管理
contract Deposit is IDeposit, Ownable, ReentrancyGuard {
    IUserRegistry public userRegistry;
    IPayment public payment;

    mapping(address => uint256) private _deposits;
    mapping(uint256 => uint256) private _operatorRequiredDeposit;

    constructor(address _userRegistry) Ownable(msg.sender) {
        userRegistry = IUserRegistry(_userRegistry);
    }

    function setPayment(address _payment) external onlyOwner {
        payment = IPayment(_payment);
    }

    function setRequiredDeposit(uint256 operatorId, uint256 amount) external onlyOwner {
        _operatorRequiredDeposit[operatorId] = amount;
    }

    function deposit() external payable {
        require(userRegistry.isRegistered(msg.sender), "Not registered");
        require(msg.value > 0, "Zero deposit");

        _deposits[msg.sender] += msg.value;
        emit DepositMade(msg.sender, msg.value);
    }

    function withdraw() external nonReentrant {
        require(_deposits[msg.sender] > 0, "No deposit");

        // 检查是否有未结账单
        if (address(payment) != address(0)) {
            IPayment.Bill[] memory unpaid = payment.getUnpaidBills(msg.sender);
            require(unpaid.length == 0, "Has unpaid bills");
        }

        uint256 amount = _deposits[msg.sender];
        _deposits[msg.sender] = 0;

        (bool success, ) = msg.sender.call{value: amount}("");
        require(success, "Transfer failed");

        emit DepositWithdrawn(msg.sender, amount);
    }

    function getDepositAmount(address user) external view returns (uint256) {
        return _deposits[user];
    }

    function getRequiredDeposit(uint256 operatorId) external view returns (uint256) {
        return _operatorRequiredDeposit[operatorId];
    }
}
