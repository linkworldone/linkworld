import { useReadContract, useWriteContract, useWaitForTransactionReceipt, useChainId } from "wagmi";
import { ServiceManagerABI } from "../../config/abis";
import { getContractAddress } from "../../config/contracts";

export function useContractUserService(address: `0x${string}` | undefined) {
  const chainId = useChainId();
  return useReadContract({
    address: getContractAddress(chainId, "ServiceManager"),
    abi: ServiceManagerABI,
    functionName: "getUserService",
    args: address ? [address] : undefined,
    query: { enabled: !!address, staleTime: 60_000 },
  });
}

export function useContractActivateService() {
  const chainId = useChainId();
  const { writeContract, data: hash, isPending, error } = useWriteContract();
  const { isLoading: isConfirming, isSuccess } = useWaitForTransactionReceipt({ hash });

  const activateService = (operatorId: bigint, virtualNumber: string, password: string) => {
    writeContract({
      address: getContractAddress(chainId, "ServiceManager"),
      abi: ServiceManagerABI,
      functionName: "activateService",
      args: [operatorId, virtualNumber, password],
    });
  };

  return { activateService, hash, isPending, isConfirming, isSuccess, error };
}
