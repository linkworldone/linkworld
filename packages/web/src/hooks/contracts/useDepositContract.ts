import { useReadContract, useWriteContract, useWaitForTransactionReceipt, useChainId } from "wagmi";
import { parseUnits } from "viem";
import { DepositABI } from "../../config/abis";
import { getContractAddress, getUsdtDecimals } from "../../config/contracts";

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
