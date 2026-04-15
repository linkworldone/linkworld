import { useQuery, useQueryClient } from "@tanstack/react-query";
import { operatorApi } from "../services/api/operatorApi";
import { useContractActivateService } from "./contracts";
import { useTxState } from "./useTransactionFlow";
import type { VirtualNumber } from "../types";
import { apiClient } from "../services/api/client";

// 地区列表：从后端 API 聚合
export function useRegions() {
  return useQuery({
    queryKey: ["regions"],
    queryFn: () => operatorApi.getRegions(),
    staleTime: 5 * 60_000, // 5 分钟
  });
}

// 地区运营商：从后端 API 筛选
export function useOperatorsByRegion(regionCode?: string) {
  return useQuery({
    queryKey: ["operators", regionCode],
    queryFn: () => operatorApi.getOperatorsByRegion(regionCode!),
    enabled: !!regionCode,
    staleTime: 5 * 60_000,
  });
}

// 我的号码：从后端 API
export function useMyNumbers(address?: string) {
  return useQuery({
    queryKey: ["myNumbers", address],
    queryFn: async () => {
      if (!address) return [];
      try {
        const service = await apiClient.get<any, any>(`/api/service/${address}`);
        if (service && service.virtual_number) {
          return [{
            id: String(service.id),
            number: service.virtual_number,
            region: "",
            operator: String(service.operator_id),
            status: service.is_active ? "active" : "inactive",
            activatedAt: service.activated_at,
            credentials: { type: "voip" as const, config: `SIP: ${service.virtual_number}` },
          }] as VirtualNumber[];
        }
        return [];
      } catch {
        return [];
      }
    },
    enabled: !!address,
  });
}

// 申请号码：三步（后端生成号码 → 合约激活 → 后端存档）
export function useApplyNumber() {
  const queryClient = useQueryClient();
  const contractActivate = useContractActivateService();
  const txState = useTxState(contractActivate);

  return {
    applyNumber: contractActivate.activateService,
    txState,
    invalidate: () => queryClient.invalidateQueries({ queryKey: ["myNumbers"] }),
  };
}
