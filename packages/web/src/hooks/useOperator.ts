import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { useAccount } from "wagmi";
import { operatorApi } from "../services/api/operatorApi";
import { parseContractError, type TxState } from "./useTransactionFlow";
import type { VirtualNumber } from "../types";
import { apiClient } from "../services/api/client";
import { signedPost, WalletAuthRejectedError } from "../services/api/signedPost";

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
          // 查 operator 信息
          let operatorName = String(service.operator_id);
          let region = "";
          try {
            const operators = await apiClient.get<any, any[]>("/api/operators");
            const op = operators?.find((o: any) => o.id === service.operator_id);
            if (op) {
              operatorName = op.name;
              region = op.country_code || op.region || "";
            }
          } catch {}

          return [{
            id: String(service.id),
            number: service.virtual_number,
            region,
            operator: operatorName,
            status: service.is_active ? "active" as const : "inactive" as const,
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

// 申请号码（T4 service 写链路重写）：
// 新 ServiceManager 是运营商注册中心模型，**不再有链上 activateService 用户态**（旧 0G 模型）。
// 虚拟号码激活改走后端 `/api/service/activate`（handoff §6「virtual-number 等流程不变」=后端流程）。
// applyNumber 直接发后端激活意向，并用本地 TxState 反映 pending→success，保持调用面（txState/invalidate）
// 与旧接口一致，T10 申请弹层无需改结构。完整 T10 申请弹层（手续费展示）改造后续。
export function useApplyNumber() {
  const queryClient = useQueryClient();
  const { address } = useAccount();
  const { t } = useTranslation();
  const [txState, setTxState] = useState<TxState>({ status: "idle" });

  const applyNumber = async (
    operatorId: bigint,
    virtualNumber: string,
    password: string,
  ) => {
    if (!address) return;
    setTxState({ status: "pending-confirmation" });
    try {
      // T5：service 写端点经 signedPost 带 WalletAuth 身份签名（action="service/activate"）。
      await signedPost(
        "/api/service/activate",
        {
          wallet: address,
          operator_id: Number(operatorId),
          virtual_number: virtualNumber,
          password,
        },
        { wallet: address, action: "service/activate" },
      );
      setTxState({ status: "success" });
      queryClient.invalidateQueries({ queryKey: ["myNumbers"] });
    } catch (err) {
      // 拒签：身份签名被取消（与交易签名/合约 revert 文案区分，design §3.7）。
      const error =
        err instanceof WalletAuthRejectedError
          ? t("errors.authCancelled")
          : parseContractError(err, t);
      setTxState({ status: "error", error });
    }
  };

  return {
    applyNumber,
    txState,
    invalidate: () => queryClient.invalidateQueries({ queryKey: ["myNumbers"] }),
  };
}
