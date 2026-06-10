import { useReadContract, useWriteContract, useWaitForTransactionReceipt, useChainId } from "wagmi";
import { parseUnits } from "viem";
import { DepositABI, MockUSDTABI } from "../../config/abis";
import { getContractAddress, getUsdt, getUsdtDecimals } from "../../config/contracts";

export function useDepositBalance(address: `0x${string}` | undefined) {
  const chainId = useChainId();
  return useReadContract({
    address: getContractAddress(chainId, "Deposit"),
    abi: DepositABI,
    functionName: "getDepositAmount",
    args: address ? [address] : undefined,
    query: { enabled: !!address, staleTime: 30_000 },
  });
}

/**
 * 读锁仓到期时间戳（秒）`getLockExpiry(addr)`（design §3.4 / DESIGN.md §3）。
 * 返回 bigint：0=无锁仓；>0=到期 unix 秒。LockCountdown 据此渲染倒计时/解锁。
 * 锁仓到期对提现可用性敏感，短 staleTime 让充值顺延/提现归 0 后能尽快重读。
 */
export function useLockExpiry(address: `0x${string}` | undefined) {
  const chainId = useChainId();
  return useReadContract({
    address: getContractAddress(chainId, "Deposit"),
    abi: DepositABI,
    functionName: "getLockExpiry",
    args: address ? [address] : undefined,
    query: { enabled: !!address, staleTime: 10_000 },
  });
}

/**
 * 读钱包 USDT(ERC20) 余额 `balanceOf(addr)`（6 位最小单位）。
 * 充值前校验「amount ≤ 钱包 USDT 余额」用（design §3.2）。
 */
export function useUsdtBalance(address: `0x${string}` | undefined) {
  const chainId = useChainId();
  return useReadContract({
    address: getUsdt(chainId),
    abi: MockUSDTABI,
    functionName: "balanceOf",
    args: address ? [address] : undefined,
    query: { enabled: !!address, staleTime: 10_000 },
  });
}

export function useContractDeposit() {
  const chainId = useChainId();
  const { writeContract, data: hash, isPending, error } = useWriteContract();
  const { isPending: rawConfirming, isSuccess } = useWaitForTransactionReceipt({ hash });
  const isConfirming = !!hash && rawConfirming;

  // T1: deposit 改 ERC20 deposit(uint256)（去 payable）。此处仅做保编译最小适配——
  // 按 usdtDecimals(6) 解析金额、去 value。前置 approve 两步态 / allowance 检测留 T3，
  // 完整精度/币种文案统一留 T2。
  const deposit = (amountUsdt: string) => {
    writeContract({
      address: getContractAddress(chainId, "Deposit"),
      abi: DepositABI,
      functionName: "deposit",
      args: [parseUnits(amountUsdt, getUsdtDecimals(chainId))],
    });
  };

  return { deposit, hash, isPending, isConfirming, isSuccess, error };
}

export function useContractWithdraw() {
  const chainId = useChainId();
  const { writeContract, data: hash, isPending, error } = useWriteContract();
  const { isPending: rawConfirming, isSuccess } = useWaitForTransactionReceipt({ hash });
  const isConfirming = !!hash && rawConfirming;

  const withdraw = () => {
    writeContract({
      address: getContractAddress(chainId, "Deposit"),
      abi: DepositABI,
      functionName: "withdraw",
    });
  };

  return { withdraw, hash, isPending, isConfirming, isSuccess, error };
}
