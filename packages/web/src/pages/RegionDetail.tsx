import { useState, useEffect, useRef } from "react";
import { useParams } from "react-router-dom";
import { useAccount } from "wagmi";
import { Button } from "@/components/ui/button";
import { BottomSheet } from "@/components/shared/BottomSheet";
import { useOperatorsByRegion, useApplyNumber, useRegions } from "@/hooks/useOperator";
import { formatAmount } from "@/utils/format";
import { apiClient } from "@/services/api/client";

export default function RegionDetail() {
  const { regionCode } = useParams<{ regionCode: string }>();
  const { address } = useAccount();
  const { data: regions } = useRegions();
  const { data: operators } = useOperatorsByRegion(regionCode);
  const { applyNumber, txState, invalidate } = useApplyNumber();

  const region = regions?.find((r) => r.code === regionCode);
  const [selectedOp, setSelectedOp] = useState<string | null>(null);
  const pendingServiceRef = useRef<{ operatorId: string; virtualNumber: string; password: string } | null>(null);

  const isApplying = txState.status === "pending-signature" || txState.status === "pending-confirmation";

  // 合约成功后同步后端并刷新号码列表
  useEffect(() => {
    if (txState.status === "success" && address && pendingServiceRef.current) {
      const { operatorId, virtualNumber, password } = pendingServiceRef.current;
      apiClient.post("/api/service/activate", {
        wallet: address,
        operator_id: Number(operatorId),
        virtual_number: virtualNumber,
        password: password,
      }).catch((err: any) => console.error("Backend sync failed:", err));
      pendingServiceRef.current = null;
      invalidate();
      setSelectedOp(null);
    }
  }, [txState.status]);

  const handleApply = async () => {
    if (!address || !selectedOp) return;
    try {
      // 后端期望 country_code，从 regionCode 获取
      const result = await apiClient.post("/api/virtual-number/generate", { country_code: regionCode }) as any;
      const virtualNumber = result.virtual_number;
      const password = result.password;
      // 存下来供合约成功后同步后端
      pendingServiceRef.current = { operatorId: selectedOp, virtualNumber, password };
      // 调合约激活
      applyNumber(BigInt(selectedOp), virtualNumber, password);
    } catch (err) {
      console.error("Failed to generate virtual number:", err);
    }
  };

  const selectedOperator = operators?.find((o) => o.id === selectedOp);

  return (
    <div className="px-4 space-y-3">
      <div className="text-center mb-2">
        <span className="text-3xl">{region?.flag}</span>
        <h2 className="text-lg font-bold mt-1">{region?.name}</h2>
      </div>

      {operators?.map((op) => (
        <div key={op.id} className="p-4 bg-surface-card rounded-xl">
          <div className="flex justify-between items-start mb-3">
            <div className="text-sm font-bold">{op.name}</div>
            <span className={`text-[10px] font-semibold px-2 py-0.5 rounded-md ${
              op.isActive ? "bg-status-success/20 text-status-success" : "bg-surface-secondary text-text-muted"
            }`}>
              {op.isActive ? "Active" : "Inactive"}
            </span>
          </div>
          <div className="grid grid-cols-3 gap-2 text-center mb-3">
            <div>
              <div className="text-[10px] text-text-muted">Data Rate</div>
              <div className="text-xs font-semibold">${op.dataRate}/GB</div>
            </div>
            <div>
              <div className="text-[10px] text-text-muted">Call Rate</div>
              <div className="text-xs font-semibold">${op.callRate}/min</div>
            </div>
            <div>
              <div className="text-[10px] text-text-muted">Min Deposit</div>
              <div className="text-xs font-semibold">{formatAmount(op.requiredDeposit)}</div>
            </div>
          </div>
          <Button onClick={() => setSelectedOp(op.id)} className="w-full" disabled={!op.isActive}>
            Apply for Number
          </Button>
        </div>
      ))}

      <BottomSheet open={selectedOp !== null} onOpenChange={(o) => !o && setSelectedOp(null)}>
        <h2 className="text-lg font-bold mb-2">Apply for Virtual Number</h2>
        {selectedOperator && (
          <>
            <div className="text-sm text-text-secondary mb-4">
              {region?.flag} {region?.name} · {selectedOperator.name}
            </div>
            <div className="p-3 bg-surface-secondary rounded-xl mb-4">
              <div className="flex justify-between text-xs">
                <span className="text-text-muted">Required Deposit</span>
                <span className="font-semibold">{formatAmount(selectedOperator.requiredDeposit)} USDT</span>
              </div>
            </div>
            <Button onClick={handleApply} disabled={isApplying} className="w-full py-3">
              {isApplying ? "Applying..." : "Confirm Application"}
            </Button>
          </>
        )}
      </BottomSheet>
    </div>
  );
}
