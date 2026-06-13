import { Loader2, CheckCircle2, XCircle } from "lucide-react";
import { useTranslation } from "react-i18next";
import { Badge } from "@/components/ui/badge";
import { cn } from "@/lib/utils";

/**
 * TxStatusBadge —— 通用交易三态徽章（design §3.1 / DESIGN.md 三态表）。
 *
 * 铁律：
 * - `pending`「处理中」：**弱化、绝不染绿**（pending 不计入终态成功，不暗示已到账）。
 * - `confirmed`「已确认」：唯一染绿态（后端 event_sync 等 K 块回填后的终态）。
 * - `failed`「失败/已回退」：区分 `revert`（链上执行失败）vs `reorg`（区块重组回退，
 *   非用户错误，文案安抚为「已回退」）。
 *
 * 用语义 status token 类（T6 换肤上色，本组件不写裸 HEX）。
 */

export type TxBadgeStatus = "pending" | "confirmed" | "failed";
export type TxFailureReason = "revert" | "reorg";

export interface TxStatusBadgeProps {
  status: TxBadgeStatus;
  /** failed 时区分失败原因：revert（执行失败）/ reorg（区块重组回退）。 */
  failureReason?: TxFailureReason;
  className?: string;
}

export function TxStatusBadge({
  status,
  failureReason,
  className,
}: TxStatusBadgeProps) {
  const { t } = useTranslation();
  if (status === "pending") {
    return (
      <Badge
        variant="secondary"
        className={cn("text-text-secondary", className)}
        data-slot="tx-status-badge"
        data-status="pending"
      >
        <Loader2 className="size-3 animate-spin" />
        {t("txStatus.pending")}
      </Badge>
    );
  }

  if (status === "confirmed") {
    return (
      <Badge
        className={cn(
          "bg-status-success/15 text-status-success border-transparent",
          className
        )}
        data-slot="tx-status-badge"
        data-status="confirmed"
      >
        <CheckCircle2 className="size-3" />
        {t("txStatus.confirmed")}
      </Badge>
    );
  }

  // failed
  const isReorg = failureReason === "reorg";
  return (
    <Badge
      variant="destructive"
      className={className}
      data-slot="tx-status-badge"
      data-status="failed"
      data-failure={failureReason ?? "revert"}
    >
      <XCircle className="size-3" />
      {isReorg ? t("txStatus.reverted") : t("txStatus.failed")}
    </Badge>
  );
}
