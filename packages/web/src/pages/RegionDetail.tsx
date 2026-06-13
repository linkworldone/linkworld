import { useState, useEffect, useRef } from "react";
import { useParams } from "react-router-dom";
import { useAccount } from "wagmi";
import { useTranslation } from "react-i18next";
import { Wifi, Phone, Wallet, ShieldCheck, AlertCircle } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { BottomSheet } from "@/components/shared/BottomSheet";
import { AmountDisplay } from "@/components/shared/AmountDisplay";
import { FeeBreakdown } from "@/components/shared/FeeBreakdown";
import { useOperatorsByRegion, useApplyNumber, useRegions } from "@/hooks/useOperator";
import { WalletAuthRejectedError } from "@/services/api/signedPost";
import { apiClient } from "@/services/api/client";

export default function RegionDetail() {
  const { t } = useTranslation();
  const { regionCode } = useParams<{ regionCode: string }>();
  const { address } = useAccount();
  const { data: regions } = useRegions();
  const { data: operators } = useOperatorsByRegion(regionCode);
  const { applyNumber, txState, invalidate } = useApplyNumber();

  const region = regions?.find((r) => r.code === regionCode);
  const [selectedOp, setSelectedOp] = useState<string | null>(null);
  const [rejectMsg, setRejectMsg] = useState<string | null>(null);
  const pendingServiceRef = useRef<{ operatorId: string; virtualNumber: string; password: string } | null>(null);

  const isApplying = txState.status === "pending-signature" || txState.status === "pending-confirmation";

  // T4：虚拟号码激活已迁后端（useApplyNumber 内部发 /api/service/activate，不再链上 activateService）。
  // success 后仅收尾刷新；后端激活请求由 applyNumber 一次性发出（不再在此重复 POST）。
  useEffect(() => {
    if (txState.status === "success") {
      pendingServiceRef.current = null;
      invalidate();
      setSelectedOp(null);
    }
  }, [txState.status]);

  // 拒签提示对齐（design §3.7）：身份签名被取消，操作未提交，不进 pending。
  useEffect(() => {
    if (txState.status === "error" && txState.error) {
      setRejectMsg(txState.error);
    }
  }, [txState.status, txState.error]);

  const handleApply = async () => {
    if (!address || !selectedOp) return;
    setRejectMsg(null);
    try {
      // 后端期望 country_code，从 regionCode 获取
      const result = await apiClient.post("/api/virtual-number/generate", { country_code: regionCode }) as any;
      const virtualNumber = result.virtual_number;
      const password = result.password;
      pendingServiceRef.current = { operatorId: selectedOp, virtualNumber, password };
      // 后端激活意向（T4：去链上 activateService；T5：signedPost 带身份签名）。
      applyNumber(BigInt(selectedOp), virtualNumber, password);
    } catch (err) {
      // virtual-number 生成失败 / 身份签名取消（WalletAuthRejected）。
      const msg =
        err instanceof WalletAuthRejectedError
          ? t("errors.authCancelled")
          : t("regionDetail.applyFailed");
      setRejectMsg(msg);
      console.error("Failed to generate virtual number:", err);
    }
  };

  const selectedOperator = operators?.find((o) => o.id === selectedOp);

  return (
    <div className="px-4 space-y-3">
      {/* 页头：navy 画布上，深底文字（gold 国名强调） */}
      <div className="text-center mb-2 pt-2">
        <span className="text-3xl" aria-hidden>{region?.flag}</span>
        <h2 className="font-display text-lg font-bold mt-1 text-text-on-dark-gold">
          {region ? t(`destinations.${region.code}`, { defaultValue: region.name }) : ""}
        </h2>
      </div>

      {operators?.map((op) => (
        <Card key={op.id}>
          <CardContent className="space-y-3">
            <div className="flex justify-between items-start">
              <div className="text-sm font-bold text-text-on-light-primary">{op.name}</div>
              <Badge variant={op.isActive ? "default" : "secondary"}>
                {op.isActive ? t("regionDetail.operatorActive") : t("regionDetail.operatorInactive")}
              </Badge>
            </div>
            <div className="grid grid-cols-3 gap-2 text-center">
              <div>
                <div className="flex items-center justify-center gap-1 text-[10px] text-text-on-light-secondary">
                  <Wifi className="size-3" /> {t("regionDetail.dataRate")}
                </div>
                <div className="text-xs font-semibold text-text-on-light-primary">${op.dataRate}/GB</div>
              </div>
              <div>
                <div className="flex items-center justify-center gap-1 text-[10px] text-text-on-light-secondary">
                  <Phone className="size-3" /> {t("regionDetail.callRate")}
                </div>
                <div className="text-xs font-semibold text-text-on-light-primary">${op.callRate}/min</div>
              </div>
              <div>
                <div className="flex items-center justify-center gap-1 text-[10px] text-text-on-light-secondary">
                  <Wallet className="size-3" /> {t("regionDetail.minDeposit")}
                </div>
                {/* 卡内（暖米白底）金额 navy（B2）；6 位精度 USDT */}
                <AmountDisplay amount={op.requiredDeposit} currency="USDT" size="sm" />
              </div>
            </div>
            <Button onClick={() => setSelectedOp(op.id)} className="w-full" disabled={!op.isActive}>
              {t("regionDetail.applyForNumber")}
            </Button>
          </CardContent>
        </Card>
      ))}

      <BottomSheet open={selectedOp !== null} onOpenChange={(o) => !o && setSelectedOp(null)}>
        <h2 className="font-display text-lg font-bold mb-2 text-text-on-light-primary">{t("regionDetail.sheetTitle")}</h2>
        {selectedOperator && (
          <>
            <div className="text-sm text-text-on-light-secondary mb-4">
              <span aria-hidden>{region?.flag}</span> {region ? t(`destinations.${region.code}`, { defaultValue: region.name }) : ""} · {selectedOperator.name}
            </div>

            {/* 费用明细（design §3.6）：押金本金 navy + 平台手续费读链（FeeBreakdown，amount=押金本金）。 */}
            <div className="p-3 bg-surface-input rounded-xl mb-4 space-y-2">
              <div className="flex justify-between items-center text-xs">
                <span className="text-text-on-light-secondary">{t("regionDetail.requiredDeposit")}</span>
                <AmountDisplay amount={selectedOperator.requiredDeposit} currency="USDT" size="sm" />
              </div>
              {/* 平台手续费读链：calculateFee(押金本金)，禁写死（D12 / B 红线） */}
              <FeeBreakdown amount={selectedOperator.requiredDeposit} />
            </div>

            {/* 身份签名提示（design §3.7）：申请走 signedPost 意向，需钱包身份签名（不消耗 gas）。 */}
            <div className="flex items-start gap-2 text-[11px] text-text-on-light-secondary mb-4">
              <ShieldCheck className="size-3.5 mt-0.5 shrink-0 text-brand-royal" />
              <span>{t("regionDetail.signatureHint")}</span>
            </div>

            {rejectMsg && (
              <div className="flex items-start gap-2 text-xs text-status-danger mb-4" data-slot="reject-msg">
                <AlertCircle className="size-3.5 mt-0.5 shrink-0" />
                <span>{rejectMsg}</span>
              </div>
            )}

            <Button onClick={handleApply} disabled={isApplying} className="w-full py-3">
              {isApplying ? t("regionDetail.applying") : t("regionDetail.confirmApplication")}
            </Button>
          </>
        )}
      </BottomSheet>
    </div>
  );
}
