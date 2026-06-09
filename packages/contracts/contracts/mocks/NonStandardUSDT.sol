// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

/// @title 非标 ERC20 测试桩（仅供 SafeERC20 分支测试 USDT-01/02）
/// @dev 不继承 OZ ERC20，手写最小账本，以便精确控制 transfer/transferFrom 的返回值/无返回值行为。
///      用于验证 SafeERC20 对真实 USDT 这类非标代币（transfer 无 bool 返回值 / 返回 false）的处理。

/// @notice USDT-01：transfer/transferFrom 不返回任何值（如真实 Tether USDT）。
///         验证 SafeERC20 仍能正确入账（不因缺返回值而 revert）。
contract NoReturnUSDT {
    string public constant name = "NoReturn USDT";
    string public constant symbol = "USDT";
    uint8 public constant decimals = 6;

    mapping(address => uint256) public balanceOf;
    mapping(address => mapping(address => uint256)) public allowance;

    function mint(address to, uint256 amount) external {
        balanceOf[to] += amount;
    }

    function approve(address spender, uint256 amount) external {
        allowance[msg.sender][spender] = amount;
        // 无返回值（非标）
    }

    /// @notice 无 bool 返回值（非标），但真实执行转账并在余额不足时 revert
    function transfer(address to, uint256 amount) external {
        require(balanceOf[msg.sender] >= amount, "insufficient balance");
        balanceOf[msg.sender] -= amount;
        balanceOf[to] += amount;
        // 无返回值
    }

    /// @notice 无 bool 返回值（非标）
    function transferFrom(address from, address to, uint256 amount) external {
        require(balanceOf[from] >= amount, "insufficient balance");
        uint256 allowed = allowance[from][msg.sender];
        require(allowed >= amount, "insufficient allowance");
        if (allowed != type(uint256).max) {
            allowance[from][msg.sender] = allowed - amount;
        }
        balanceOf[from] -= amount;
        balanceOf[to] += amount;
        // 无返回值
    }
}

/// @notice USDT-02：transfer/transferFrom 返回 false（恶意/异常代币），不 revert。
///         验证 SafeERC20 捕获 false 并使 deposit/payBill 整体 revert（不静默失败）。
contract FalseReturnUSDT {
    string public constant name = "FalseReturn USDT";
    string public constant symbol = "USDT";
    uint8 public constant decimals = 6;

    mapping(address => uint256) public balanceOf;
    mapping(address => mapping(address => uint256)) public allowance;

    function mint(address to, uint256 amount) external {
        balanceOf[to] += amount;
    }

    function approve(address spender, uint256 amount) external returns (bool) {
        allowance[msg.sender][spender] = amount;
        return true;
    }

    /// @notice 总是返回 false，且不实际转账（静默失败的典型恶意行为）
    function transfer(address, uint256) external pure returns (bool) {
        return false;
    }

    /// @notice 总是返回 false
    function transferFrom(address, address, uint256) external pure returns (bool) {
        return false;
    }
}
