import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { simApi, type SimClaimPayload, type SimRecord } from "../services/api/simApi";

// 我的 SIM 列表（公开读）。pending 期短轮询促使 pending→confirmed 尽快可见。
export function useMySims(wallet?: string) {
  return useQuery({
    queryKey: ["mySims", wallet],
    queryFn: (): Promise<SimRecord[]> => simApi.getMySims(wallet!),
    enabled: !!wallet,
    staleTime: 2_000,
    refetchInterval: (query) => {
      const data = query.state.data as SimRecord[] | undefined;
      const hasPending = data?.some((r) => r.status === "pending");
      return hasPending ? 5_000 : false;
    },
  });
}

// 领取 SIM（链上 redeemForSim 成功后调用）。成功后刷新「我的 SIM」列表。
export function useClaimSim(wallet?: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: SimClaimPayload) => {
      if (!wallet) throw new Error("钱包未连接");
      return simApi.claim(wallet, payload);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["mySims"] });
    },
  });
}
