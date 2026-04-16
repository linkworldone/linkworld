import { useReadContract, useWriteContract, useWaitForTransactionReceipt, useChainId } from "wagmi";
import { parseEther } from "viem";
import { DepositABI } from "../../config/abis";
import { getContractAddress } from "../../config/contracts";

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
  const { isPending: isConfirming, isSuccess } = useWaitForTransactionReceipt({ hash });

  const deposit = (amountEth: string) => {
    writeContract({
      address: getContractAddress(chainId, "Deposit"),
      abi: DepositABI,
      functionName: "deposit",
      value: parseEther(amountEth),
    });
  };

  return { deposit, hash, isPending, isConfirming, isSuccess, error };
}

export function useContractWithdraw() {
  const chainId = useChainId();
  const { writeContract, data: hash, isPending, error } = useWriteContract();
  const { isPending: isConfirming, isSuccess } = useWaitForTransactionReceipt({ hash });

  const withdraw = () => {
    writeContract({
      address: getContractAddress(chainId, "Deposit"),
      abi: DepositABI,
      functionName: "withdraw",
    });
  };

  return { withdraw, hash, isPending, isConfirming, isSuccess, error };
}
