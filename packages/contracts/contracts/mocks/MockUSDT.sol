// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";

/// @title MockUSDT - 测试/测试网用 6 位精度 ERC20（标准实现，无 fee-on-transfer）
/// @dev 仅供测试与测试网部署。public mint 便于任意发币。
contract MockUSDT is ERC20 {
    constructor() ERC20("Mock USDT", "USDT") {}

    /// @notice USDT 标准精度为 6 位
    function decimals() public pure override returns (uint8) {
        return 6;
    }

    /// @notice 公开 mint，测试任意发币
    function mint(address to, uint256 amount) external {
        _mint(to, amount);
    }
}
