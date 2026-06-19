// 本地 hardhat 链(31337)写死 gas limit，绕过 estimateGas（本地链上 estimateGas 慢/易失败，
// 表现为交易要等几分钟才弹窗、"费用不可用"）；真实网络(如 Arbitrum Sepolia 421614)返回空对象，
// 保持钱包/viem 自动估算不变。只覆盖 gas(gas limit)，不碰 maxFeePerGas 等 fee 字段。
//
// 2_000_000 足以覆盖 deposit 充值即 mint 多张 NFT 的最大开销；hardhat 本地链不实际收费，偏大无害。
export function localGasOverride(chainId?: number): { gas?: bigint } {
  return chainId === 31337 ? { gas: 2_000_000n } : {};
}
