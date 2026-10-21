import {
  useReadContract,
  useWriteContract,
  useWaitForTransactionReceipt,
  useChainId,
} from "wagmi";
import { parseUnits } from "viem";
import { DepositABI, MockUSDTABI } from "../../config/abis";
import {
  getUsdtDecimals,
  getOptionalContractAddress,
  getOptionalUsdtAddress,
} from "../../config/contracts";
import { localGasOverride } from "./gasOverride";

export function useDepositBalance(address: `0x${string}` | undefined) {
  const chainId = useChainId();
  const contractAddress = getOptionalContractAddress(chainId, "Deposit");

  return useReadContract({
    address: contractAddress,
    abi: DepositABI,
    functionName: "getDepositAmount",
    args: address ? [address] : undefined,
    query: {
      enabled: Boolean(address && contractAddress),
      staleTime: 30_000,
    },
  });
}

export function useLockExpiry(address: `0x${string}` | undefined) {
  const chainId = useChainId();
  const contractAddress = getOptionalContractAddress(chainId, "Deposit");

  return useReadContract({
    address: contractAddress,
    abi: DepositABI,
    functionName: "getLockExpiry",
    args: address ? [address] : undefined,
    query: {
      enabled: Boolean(address && contractAddress),
      staleTime: 10_000,
    },
  });
}

export function useUsdtBalance(address: `0x${string}` | undefined) {
  const chainId = useChainId();
  const usdtAddress = getOptionalUsdtAddress(chainId);

  return useReadContract({
    address: usdtAddress,
    abi: MockUSDTABI,
    functionName: "balanceOf",
    args: address ? [address] : undefined,
    query: {
      enabled: Boolean(address && usdtAddress),
      staleTime: 10_000,
    },
  });
}

export function useContractDeposit() {
  const chainId = useChainId();
  const {
    writeContract,
    data: hash,
    isPending,
    error,
  } = useWriteContract();

  const {
    isPending: rawConfirming,
    isSuccess,
  } = useWaitForTransactionReceipt({ hash });

  const isConfirming = Boolean(hash) && rawConfirming;

  const deposit = (amountUsdt: string) => {
    const contractAddress = getOptionalContractAddress(chainId, "Deposit");
    const usdtAddress = getOptionalUsdtAddress(chainId);

    if (!contractAddress || !usdtAddress) {
      throw new Error(
        `Deposit or USDT contract is not configured for chainId ${chainId}`
      );
    }

    writeContract({
      address: contractAddress,
      abi: DepositABI,
      functionName: "deposit",
      args: [parseUnits(amountUsdt, getUsdtDecimals(chainId))],
      ...localGasOverride(chainId),
    });
  };

  return {
    deposit,
    hash,
    isPending,
    isConfirming,
    isSuccess,
    error,
  };
}

export interface Tranche {
  amount: bigint;
  unlockAt: bigint;
  withdrawn: boolean;
}

export function useTranches(address: `0x${string}` | undefined) {
  const chainId = useChainId();
  const contractAddress = getOptionalContractAddress(chainId, "Deposit");

  const query = useReadContract({
    address: contractAddress,
    abi: DepositABI,
    functionName: "getTranches",
    args: address ? [address] : undefined,
    query: {
      enabled: Boolean(address && contractAddress),
      staleTime: 10_000,
    },
  });

  const tranches =
    (query.data as readonly Tranche[] | undefined) ?? [];

  return {
    ...query,
    tranches,
  };
}

export function useContractWithdraw() {
  const chainId = useChainId();
  const {
    writeContract,
    data: hash,
    isPending,
    error,
  } = useWriteContract();

  const {
    isPending: rawConfirming,
    isSuccess,
  } = useWaitForTransactionReceipt({ hash });

  const isConfirming = Boolean(hash) && rawConfirming;

  const withdraw = (trancheIndex: number | bigint) => {
    const contractAddress = getOptionalContractAddress(chainId, "Deposit");

    if (!contractAddress) {
      throw new Error(
        `Deposit contract is not configured for chainId ${chainId}`
      );
    }

    writeContract({
      address: contractAddress,
      abi: DepositABI,
      functionName: "withdraw",
      args: [BigInt(trancheIndex)],
      ...localGasOverride(chainId),
    });
  };

  return {
    withdraw,
    hash,
    isPending,
    isConfirming,
    isSuccess,
    error,
  };
}
