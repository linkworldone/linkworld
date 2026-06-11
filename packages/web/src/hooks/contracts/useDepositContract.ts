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

/**
 * 单笔保证金（合约 `Tranche`）：amount=本金（6 位最小单位），unlockAt=解锁 unix 秒，
 * withdrawn=是否已取回。每笔独立锁 30 天，各自到期单独取（withdraw(trancheIndex)）。
 */
export interface Tranche {
  amount: bigint;
  unlockAt: bigint;
  withdrawn: boolean;
}

/**
 * 读全部保证金笔次 `getTranches(addr)`，供提现页按笔列出。
 * 数组下标即合约 trancheIndex（withdraw 逐笔取回的入参）。
 * 提现对到期/已取回敏感 → 短 staleTime 让取回后能尽快重读刷新状态。
 */
export function useTranches(address: `0x${string}` | undefined) {
  const chainId = useChainId();
  const query = useReadContract({
    address: getContractAddress(chainId, "Deposit"),
    abi: DepositABI,
    functionName: "getTranches",
    args: address ? [address] : undefined,
    query: { enabled: !!address, staleTime: 10_000 },
  });
  const tranches = (query.data as readonly Tranche[] | undefined) ?? [];
  return { ...query, tranches };
}

/**
 * 逐笔取回保证金 `withdraw(trancheIndex)`。每笔独立锁 30 天，到期单独取。
 * trancheIndex = getTranches 返回数组的下标。
 */
export function useContractWithdraw() {
  const chainId = useChainId();
  const { writeContract, data: hash, isPending, error } = useWriteContract();
  const { isPending: rawConfirming, isSuccess } = useWaitForTransactionReceipt({ hash });
  const isConfirming = !!hash && rawConfirming;

  const withdraw = (trancheIndex: number | bigint) => {
    writeContract({
      address: getContractAddress(chainId, "Deposit"),
      abi: DepositABI,
      functionName: "withdraw",
      args: [BigInt(trancheIndex)],
    });
  };

  return { withdraw, hash, isPending, isConfirming, isSuccess, error };
}
