package blockchain

// 测试用最小桩合约部署字节码（London EVM，单字节 opcode，无 PUSH0/tstore）。
//
// 为何不直接部署真实 Oracle/Deposit 字节码：项目合约用 Solidity 0.8.27 + evmVersion=cancun
// 编译，Payment/Deposit 还用了 transient storage（tstore，Cancun opcode）。而本仓库 pin 的
// go-ethereum v1.13.5 的 accounts/abi/bind/backends.SimulatedBackend 用 ethash faker 共识，
// 硬上限在 London（实测 "ethash does not support shanghai fork"，且 EVM 不认 PUSH0）。因此
// 真实 Cancun 字节码无法在此 simulated.Backend 上执行（升级 geth 或上 hardhat 31337 才行，
// 见 client_test.go 顶部说明 + T8 遗留）。
//
// 折中：部署两个 London 兼容的最小桩，把 *真实的 abigen 绑定*（OracleTransactor /
// DepositCaller / OracleCaller）指向它们，对「编码 calldata → 本地 nonce 取号 → owner 签名
// 发交易 → bind.WaitMined → 据 receipt.Status 判成败 / 解码返回值」这条 *真实链交互* 路径
// 做断言（非 mock）：
//   - stopRuntime（runtime=STOP）：任何无返回值调用（monthlySettlement /
//     issueMonthlyTrafficCards）都成功，receipt.Status=1。
//   - ret32Runtime（runtime=PUSH1 32;PUSH1 0;RETURN）：任何调用返回 32 个零字节，令
//     uint256/bool 读方法解码为 0/false。
//
// 业务合约本身的端到端结算（金额数组语义/事件回填）属 T4/T8，本轮在 hardhat 31337 或升级
// geth 后覆盖；T3 这里只验证 client 的发交易/读机制对真实链与真实绑定成立。
//
// init 字节码生成（codecopy(runtime)+RETURN，见 client_test.go 顶部 python 片段）。

// stopDeployBytecode 部署后 runtime = 0x00 (STOP)。
const stopDeployBytecode = "0x6001600c60003960016000f300"

// ret32DeployBytecode 部署后 runtime = 0x60206000f3 (PUSH1 32, PUSH1 0, RETURN) → 返回 32 零字节。
const ret32DeployBytecode = "0x6005600c60003960056000f360206000f3"
