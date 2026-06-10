import { useQuery, useQueryClient } from "@tanstack/react-query";
import { useAccount } from "wagmi";
import { billingApi } from "../services/api/billingApi";
import { usageApi } from "../services/api/usageApi";
import { useContractPayBill } from "./contracts";
import { useTxState } from "./useTransactionFlow";
import { savePendingSync, clearPendingSync, retryWithBackoff } from "../utils/pendingSync";
import { WalletAuthRejectedError } from "../services/api/signedPost";
import type { Bill } from "../types";

// 账单列表：从后端读。is_paid 由后端 BillPaid 事件回填（design §9）。
// paying（确认中）期短轮询（refetchInterval≈5s）促使 paying→paid 尽快可见，
// 短 staleTime（≈2s）不缓存未确认态；无 paying 时停轮询。
export function useBills(address?: string, filter?: Bill["status"]) {
  return useQuery({
    queryKey: ["bills", address, filter],
    queryFn: () => billingApi.getBills(address!, filter),
    enabled: !!address,
    staleTime: 2_000,
    refetchInterval: (query) => {
      const data = query.state.data as Bill[] | undefined;
      const hasPaying = data?.some((b) => b.status === "paying");
      return hasPaying ? 5_000 : false;
    },
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

// 支付账单：合约 payBill + 后端 pending 意向（不据 200 置已付，design §3.3 / handoff §1.1）。
export function usePayBill() {
  const queryClient = useQueryClient();
  const { address } = useAccount();
  const contractPay = useContractPayBill();
  const txState = useTxState(contractPay);

  return {
    payBill: contractPay.payBill,
    txState,
    // 合约成功后由页面调用：仅上报 pending 意向（POST /api/bills/pay）；
    // is_paid 唯一由后端 BillPaid 事件回填——**不据 200 置已付**。轮询 useBills 看 status 转 paid。
    recordIntent: async (billId: string) => {
      if (!address) return;
      // WalletAuth 拒签（T5）≠ 网络瞬时失败：先现签一次探测，拒签直接上抛由页面提示「身份签名被取消」，
      // 不进 pending、不落 pendingSync（对齐 useDeposit，T7 遗留对齐）。
      try {
        await billingApi.payIntent(address, billId);
        clearPendingSync(`pay_${billId}`);
        queryClient.invalidateQueries({ queryKey: ["bills"] });
        queryClient.invalidateQueries({ queryKey: ["billDetail"] });
        return;
      } catch (err) {
        if (err instanceof WalletAuthRejectedError) throw err;
      }
      // 非拒签（网络/后端瞬时失败）→ 退避重试，仍失败落 pendingSync 后台补传。
      const ok = await retryWithBackoff(() => billingApi.payIntent(address, billId));
      if (ok) {
        clearPendingSync(`pay_${billId}`);
        queryClient.invalidateQueries({ queryKey: ["bills"] });
        queryClient.invalidateQueries({ queryKey: ["billDetail"] });
      } else {
        savePendingSync(`pay_${billId}`, { wallet: address, billId });
      }
    },
  };
}
