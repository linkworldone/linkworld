// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

import "@openzeppelin/contracts/access/Ownable.sol";
import "./interfaces/IFeeManager.sol";

/// @title FeeManager - 平台手续费管理（默认 2.5%，管理员可调）
contract FeeManager is IFeeManager, Ownable {
    uint256 public constant FEE_DENOMINATOR = 10000;
    uint256 public constant MAX_FEE_RATE = 1000; // 最高 10%

    uint256 private _feeRate; // 以基点表示，250 = 2.5%

    constructor(uint256 initialFeeRate) Ownable(msg.sender) {
        require(initialFeeRate <= MAX_FEE_RATE, "Fee too high");
        _feeRate = initialFeeRate;
    }

    function getFeeRate() external view returns (uint256) {
        return _feeRate;
    }

    function setFeeRate(uint256 newRate) external onlyOwner {
        require(newRate <= MAX_FEE_RATE, "Fee too high");
        uint256 oldRate = _feeRate;
        _feeRate = newRate;
        emit FeeRateUpdated(oldRate, newRate);
    }

    function calculateFee(uint256 amount) external view returns (uint256) {
        return (amount * _feeRate) / FEE_DENOMINATOR;
    }
}
