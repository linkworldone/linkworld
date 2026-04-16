import { useQuery, useQueryClient } from "@tanstack/react-query";
import { useAccount } from "wagmi";
import { billingApi } from "../services/api/billingApi";
import { usageApi } from "../services/api/usageApi";
import { useContractPayBill } from "./contracts";
import { useTxState } from "./useTransactionFlow";
import { savePendingSync, clearPendingSync, retryWithBackoff } from "../utils/pendingSync";
import type { Bill } from "../types";

// 账单列表：从后端读
export function useBills(address?: string, filter?: "unpaid" | "paid") {
  return useQuery({
    queryKey: ["bills", address, filter],
    queryFn: () => billingApi.getBills(address!, filter),
    enabled: !!address,
    staleTime: 30_000,
  });
}

// 单笔账单：从缓存取
export function useBillDetail(billId?: string) {
  return useQuery({
    queryKey: ["billDetail", billId],
    queryFn: async (): Promise<Bill | null> => {
      // 从 bills 缓存中取，后端没有单笔查询接口
      return null;
    },
    enabled: !!billId,
  });
}

// 当月估算：从后端 usage API
export function useMonthEstimate(address?: string) {
  return useQuery({
    queryKey: ["monthEstimate", address],
    queryFn: () => usageApi.getUsage(address!),
    enabled: !!address,
    staleTime: 60_000,
  });
}

// 支付账单：合约 payBill + 后端记录（双写）
export function usePayBill() {
  const queryClient = useQueryClient();
  const { address } = useAccount();
  const contractPay = useContractPayBill();
  const txState = useTxState(contractPay);

  return {
    payBill: contractPay.payBill,
    txState,
    recordToBackend: async (billId: string) => {
      if (!address) return;
      const ok = await retryWithBackoff(() => billingApi.recordPayment(address, billId));
      if (ok) {
        clearPendingSync(`pay_${billId}`);
        queryClient.invalidateQueries({ queryKey: ["bills"] });
        queryClient.invalidateQueries({ queryKey: ["billDetail"] });
      } else {
        savePendingSync(`pay_${billId}`, { wallet: address, billId });
      }
    },
    isContractPending: contractPay.isPending,
    isConfirming: contractPay.isConfirming,
    isSuccess: contractPay.isSuccess,
  };
}
