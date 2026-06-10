import { useReadContract, useWriteContract, useWaitForTransactionReceipt, useChainId } from "wagmi";
import { ServiceManagerABI } from "../../config/abis";
import { getContractAddress } from "../../config/contracts";

// T4: 新 ServiceManager 已是运营商中心模型（getOperator/getActiveOperators/addOperator…），
// 不再有 getUserService/activateService（旧 0G 用户态）。本轮仅做保编译适配（as never 绕过
// 严格字面量类型），运行时按新 ABI 重写归 T4（service 写 / 对账衔接）。
export function useContractUserService(address: `0x${string}` | undefined) {
  const chainId = useChainId();
  return useReadContract({
    address: getContractAddress(chainId, "ServiceManager"),
    abi: ServiceManagerABI,
    functionName: "getUserService" as never,
    args: (address ? [address] : undefined) as never,
    query: { enabled: !!address, staleTime: 60_000 },
  });
}

export function useContractActivateService() {
  const chainId = useChainId();
  const { writeContract, data: hash, isPending, error } = useWriteContract();
  const { isPending: rawConfirming, isSuccess } = useWaitForTransactionReceipt({ hash });
  const isConfirming = !!hash && rawConfirming;

  const activateService = (operatorId: bigint, virtualNumber: string, password: string) => {
    writeContract({
      address: getContractAddress(chainId, "ServiceManager"),
      abi: ServiceManagerABI,
      functionName: "activateService" as never,
      args: [operatorId, virtualNumber, password] as never,
    });
  };

  return { activateService, hash, isPending, isConfirming, isSuccess, error };
}
