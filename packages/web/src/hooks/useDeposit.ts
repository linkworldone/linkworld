import { useQuery, useQueryClient } from "@tanstack/react-query";
import { useAccount } from "wagmi";
import { useDepositBalance, useContractDeposit, useContractWithdraw } from "./contracts";
import { depositApi } from "../services/api/depositApi";
import { useTxState } from "./useTransactionFlow";
import { savePendingSync, clearPendingSync, retryWithBackoff } from "../utils/pendingSync";
import type { DepositInfo, DepositRecord } from "../types";

// 查余额：从合约读（source of truth）
export function useDeposit(address?: string) {
  const { data: balance, ...rest } = useDepositBalance(address as `0x${string}` | undefined);

  const depositInfo: DepositInfo | undefined = balance !== undefined ? {
    balance: balance as bigint,
    minimumRequired: BigInt("100000000000000000000"), // 100 * 1e18
    currency: "ETH",
  } : undefined;

  return { data: depositInfo, ...rest };
}

// 充值历史：从后端读
export function useDepositHistory(address?: string) {
  return useQuery({
    queryKey: ["depositHistory", address],
    queryFn: async (): Promise<DepositRecord[]> => {
      // 后端没有 history 端点，暂时返回空
      // TODO: 后端需要加 GET /api/deposit/:wallet/history
      return [];
    },
    enabled: !!address,
  });
}

// 充值：合约 deposit + 后端记录（双写）
export function useDepositMutation() {
  const queryClient = useQueryClient();
  const { address } = useAccount();
  const contractDeposit = useContractDeposit();
  const txState = useTxState(contractDeposit);

  return {
    deposit: contractDeposit.deposit,
    txState,
    // 后端记录（合约成功后由页面组件调用）
    recordToBackend: async (amount: string) => {
      if (!address) return;
      const ok = await retryWithBackoff(() => depositApi.recordDeposit(address, amount));
      if (ok) {
        clearPendingSync(`deposit_${address}`);
        queryClient.invalidateQueries({ queryKey: ["deposit"] });
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

// 提现：合约 withdraw + 后端记录
export function useWithdrawMutation() {
  const queryClient = useQueryClient();
  const { address } = useAccount();
  const contractWithdraw = useContractWithdraw();
  const txState = useTxState(contractWithdraw);

  return {
    withdraw: contractWithdraw.withdraw,
    txState,
    recordToBackend: async () => {
      if (!address) return;
      const ok = await retryWithBackoff(() => depositApi.recordWithdraw(address));
      if (ok) {
        queryClient.invalidateQueries({ queryKey: ["deposit"] });
        queryClient.invalidateQueries({ queryKey: ["depositHistory"] });
      }
    },
    isContractPending: contractWithdraw.isPending,
    isConfirming: contractWithdraw.isConfirming,
    isSuccess: contractWithdraw.isSuccess,
  };
}
