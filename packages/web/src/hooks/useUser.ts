import { useEffect, useRef } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useAccount } from "wagmi";
import { userApi } from "../services/api/userApi";
import { useContractRegister } from "./contracts";
import { useTxState } from "./useTransactionFlow";
import { savePendingSync, clearPendingSync } from "../utils/pendingSync";

// 查用户：从后端 API 获取
export function useUser(address?: string) {
  return useQuery({
    queryKey: ["user", address],
    queryFn: () => userApi.getUser(address!),
    enabled: !!address,
    staleTime: 30_000,
  });
}

// 注册：合约铸 NFT → 后端存 profile（双写）
export function useRegister() {
  const queryClient = useQueryClient();
  const { address } = useAccount();
  const contractRegister = useContractRegister();
  const txState = useTxState(contractRegister);
  const emailRef = useRef<string>("");

  // 当合约 tx 成功后，自动调后端 API
  const backendSync = useMutation({
    mutationFn: async ({ wallet, email, tokenId }: { wallet: string; email: string; tokenId?: number }) => {
      await userApi.register(wallet, email, tokenId);
    },
    onSuccess: () => {
      if (address) {
        clearPendingSync(`register_${address}`);
        queryClient.invalidateQueries({ queryKey: ["user"] });
      }
    },
    onError: () => {
      // 后端失败，存 pending sync（保留注册时的 email）
      if (address) {
        savePendingSync(`register_${address}`, { wallet: address, email: emailRef.current });
      }
    },
  });

  // 合约交易成功后自动触发后端同步
  useEffect(() => {
    if (contractRegister.isSuccess && address && emailRef.current) {
      backendSync.mutate({ wallet: address, email: emailRef.current });
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [contractRegister.isSuccess, address]);

  const register = (email: string) => {
    emailRef.current = email;
    contractRegister.register(email);
  };

  return {
    register,
    txState,
    backendSync,
    isContractPending: contractRegister.isPending,
    isConfirming: contractRegister.isConfirming,
    isSuccess: contractRegister.isSuccess,
  };
}

// 邮箱验证保持 mock（后端未实现真邮件）
export function useSendVerificationCode() {
  return useMutation({
    mutationFn: async (_params: { address: string; email: string }) => {
      // Mock: always succeed
      return { success: true };
    },
  });
}

export function useVerifyEmail() {
  return useMutation({
    mutationFn: async (_params: { address: string; code: string }) => {
      // Mock: always pass
      return true;
    },
  });
}
