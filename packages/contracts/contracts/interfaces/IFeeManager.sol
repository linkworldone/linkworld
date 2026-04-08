// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

interface IFeeManager {
    event FeeRateUpdated(uint256 oldRate, uint256 newRate);

    function getFeeRate() external view returns (uint256);
    function setFeeRate(uint256 newRate) external;
    function calculateFee(uint256 amount) external view returns (uint256);
}
