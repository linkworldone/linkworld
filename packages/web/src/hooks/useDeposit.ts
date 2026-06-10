import { useQuery, useQueryClient } from "@tanstack/react-query";
import { useAccount } from "wagmi";
import { useDepositBalance, useContractDeposit, useContractWithdraw } from "./contracts";
import { depositApi } from "../services/api/depositApi";
import { useTxState } from "./useTransactionFlow";
import { savePendingSync, clearPendingSync, retryWithBackoff } from "../utils/pendingSync";
import { WalletAuthRejectedError } from "../services/api/signedPost";
import type { DepositInfo, DepositRecord } from "../types";

// 查余额：从链上读 getDepositAmount（source of truth, design §9）。
// 链上余额已只含 confirmed 本金（K 块逻辑在合约/后端），前端只读 confirmed 字段。
export function useDeposit(address?: string) {
  const { data: balance, refetch, ...rest } = useDepositBalance(address as `0x${string}` | undefined);

  const depositInfo: DepositInfo | undefined = balance !== undefined ? {
    balance: balance as bigint,
    currency: "USDT",
  } : undefined;

  return { data: depositInfo, refetch, ...rest };
}

// 充值/提现历史：从后端读，终态由 event_sync 回填（design §9）。
// pending 期短轮询（refetchInterval≈5s）促使 pending→confirmed 尽快可见；
// reorg 期短 staleTime（≈2s）不缓存未确认态当真。终态后由调用方/页面停轮询。
export function useDepositHistory(address?: string) {
  return useQuery({
    queryKey: ["depositHistory", address],
    queryFn: (): Promise<DepositRecord[]> => depositApi.getHistory(address!),
    enabled: !!address,
    staleTime: 2_000,
    refetchInterval: (query) => {
      const data = query.state.data as DepositRecord[] | undefined;
      // 还有 pending 记录 → 5s 轮询；全部 confirmed → 停轮询。
      const hasPending = data?.some((r) => r.status === "pending");
      return hasPending ? 5_000 : false;
    },
  });
}

// 充值：合约 deposit + 后端 pending 意向（不据成功置终态，design §3.3）。
export function useDepositMutation() {
  const queryClient = useQueryClient();
  const { address } = useAccount();
  const contractDeposit = useContractDeposit();
  const txState = useTxState(contractDeposit);

  return {
    deposit: contractDeposit.deposit,
    txState,
    // 合约成功后由页面调用：仅上报 pending 意向（POST /api/deposit），**不置终态**。
    // 真实到账以链上 getDepositAmount + 后端 event_sync 回填为准。
    recordIntent: async (amount: string) => {
      if (!address) return;
      // WalletAuth 拒签（T5）≠ 网络瞬时失败：先现签一次探测，拒签直接上抛由页面提示，不入 pendingSync。
      try {
        await depositApi.postDepositIntent(address, amount);
        clearPendingSync(`deposit_${address}`);
        queryClient.invalidateQueries({ queryKey: ["depositHistory"] });
        return;
      } catch (err) {
        if (err instanceof WalletAuthRejectedError) throw err;
      }
      // 非拒签（网络/后端瞬时失败）→ 退避重试，仍失败落 pendingSync 后台补传。
      const ok = await retryWithBackoff(() => depositApi.postDepositIntent(address, amount));
      if (ok) {
        clearPendingSync(`deposit_${address}`);
        queryClient.invalidateQueries({ queryKey: ["depositHistory"] });
      } else {
        savePendingSync(`deposit_${address}`, { wallet: address, amount });
      }
    },
    isContractPending: contractDeposit.isPending,
    isConfirming: contractDeposit.isConfirming,
    isSuccess: contractDeposit.isSuccess,
  };
}

// 提现：合约 withdraw + 后端 pending 意向（废弃凭 tx_hash 记账，design §3.3 / handoff §1.2）。
export function useWithdrawMutation() {
  const queryClient = useQueryClient();
  const { address } = useAccount();
  const contractWithdraw = useContractWithdraw();
  const txState = useTxState(contractWithdraw);

  return {
    withdraw: contractWithdraw.withdraw,
    txState,
    // 仅上报 pending 意向（POST /api/withdraw，无 tx_hash）；记账由后端 DepositWithdrawn 事件回填。
    recordIntent: async () => {
      if (!address) return;
      // 拒签（T5）上抛由页面提示「身份签名被取消」不进 pending；网络瞬时失败退避重试。
      try {
        await depositApi.postWithdrawIntent(address);
        queryClient.invalidateQueries({ queryKey: ["depositHistory"] });
        return;
      } catch (err) {
        if (err instanceof WalletAuthRejectedError) throw err;
      }
      await retryWithBackoff(() => depositApi.postWithdrawIntent(address));
      queryClient.invalidateQueries({ queryKey: ["depositHistory"] });
    },
    isContractPending: contractWithdraw.isPending,
    isConfirming: contractWithdraw.isConfirming,
    isSuccess: contractWithdraw.isSuccess,
  };
}
