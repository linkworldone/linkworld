import { useWriteContract, useWaitForTransactionReceipt, useChainId } from "wagmi";
import { PaymentABI } from "../../config/abis";
import { getContractAddress } from "../../config/contracts";

export function useContractPayBill() {
  const chainId = useChainId();
  const { writeContract, data: hash, isPending, error } = useWriteContract();
  const { isPending: rawConfirming, isSuccess } = useWaitForTransactionReceipt({ hash });
  const isConfirming = !!hash && rawConfirming;

  // T1: payBill 去 payable（链上直分 USDT，不再随 tx 带原生币）。仅做保编译最小适配——
  // 去 value。付账前 approve(Payment, amount+fee) 两步态留 T3，第二参 _value 暂保留兼容
  // 旧调用点（Billing/BillDetail），T3 接入 approve 后移除。
  const payBill = (billId: bigint, _value?: bigint) => {
    void _value;
    writeContract({
      address: getContractAddress(chainId, "Payment"),
      abi: PaymentABI,
      functionName: "payBill",
      args: [billId],
    });
  };

  return { payBill, hash, isPending, isConfirming, isSuccess, error };
}
