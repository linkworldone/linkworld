import { useWriteContract, useWaitForTransactionReceipt, useChainId } from "wagmi";
import { PaymentABI } from "../../config/abis";
import { getContractAddress } from "../../config/contracts";

export function useContractPayBill() {
  const chainId = useChainId();
  const { writeContract, data: hash, isPending, error } = useWriteContract();
  const { isLoading: isConfirming, isSuccess } = useWaitForTransactionReceipt({ hash });

  const payBill = (billId: bigint, value: bigint) => {
    writeContract({
      address: getContractAddress(chainId, "Payment"),
      abi: PaymentABI,
      functionName: "payBill",
      args: [billId],
      value,
    });
  };

  return { payBill, hash, isPending, isConfirming, isSuccess, error };
}
