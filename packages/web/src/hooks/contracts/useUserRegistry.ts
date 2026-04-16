import { useReadContract, useWriteContract, useWaitForTransactionReceipt } from "wagmi";
import { useChainId } from "wagmi";
import { getContractAddress } from "@/config/contracts";
import { UserRegistryABI } from "@/config/abis";

/* ------------------------------------------------------------------ */
/*  Read: isRegistered                                                */
/* ------------------------------------------------------------------ */

export function useIsRegistered(address: `0x${string}` | undefined) {
  const chainId = useChainId();

  return useReadContract({
    address: address ? getContractAddress(chainId, "UserRegistry") : undefined,
    abi: UserRegistryABI,
    functionName: "isRegistered",
    args: address ? [address] : undefined,
    query: { enabled: !!address },
  });
}

/* ------------------------------------------------------------------ */
/*  Write: register                                                   */
/* ------------------------------------------------------------------ */

export function useContractRegister() {
  const chainId = useChainId();
  const { data: hash, isPending, error, writeContract } = useWriteContract();

  const { isPending: isConfirming, isSuccess } = useWaitForTransactionReceipt({
    hash,
  });

  function register(email: string) {
    writeContract({
      address: getContractAddress(chainId, "UserRegistry"),
      abi: UserRegistryABI,
      functionName: "register",
      args: [email],
    });
  }

  return { register, hash, isPending, isConfirming, isSuccess, error };
}
