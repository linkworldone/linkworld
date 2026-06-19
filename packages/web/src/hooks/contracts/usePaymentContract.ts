import { useWriteContract, useWaitForTransactionReceipt, useChainId } from "wagmi";
import { PaymentABI } from "../../config/abis";
import { getContractAddress } from "../../config/contracts";
import { localGasOverride } from "./gasOverride";

export function useContractPayBill() {
  const chainId = useChainId();
  const { writeContract, data: hash, isPending, error } = useWriteContract();
  const { isPending: rawConfirming, isSuccess } = useWaitForTransactionReceipt({ hash });
  const isConfirming = !!hash && rawConfirming;

  // payBill 去 payable（链上直分 USDT，不再随 tx 带原生币）。付账前 USDT approve(Payment,
  // amount+calculateFee) 走 TwoStepAction（T9），故此处仅传 billId——不再带 value/总额第二参。
  const payBill = (billId: bigint) => {
    writeContract({
      address: getContractAddress(chainId, "Payment"),
      abi: PaymentABI,
      functionName: "payBill",
      args: [billId],
      ...localGasOverride(chainId),
    });
  };

  return { payBill, hash, isPending, isConfirming, isSuccess, error };
}
