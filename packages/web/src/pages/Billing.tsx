import { useMemo } from "react";
import { useAccount } from "wagmi";
import { useTranslation } from "react-i18next";
import { ArrowDownToLine, Flame, History as HistoryIcon } from "lucide-react";
import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { EmptyState } from "@/components/shared/EmptyState";
import { useDepositHistory } from "@/hooks/useDeposit";
import { useMySims } from "@/hooks/useSim";
import { formatAmount, formatDate } from "@/utils/format";
import type { DepositRecord } from "@/types";
import type { SimRecord } from "@/services/api/simApi";

// 流量卡 NFT 活动时间线（History）：获取（充值发卡）+ 销毁（兑换 SIM）两类记录混排倒序。
// 合约规则：每 10 USDT 发 1 张卡。DepositRecord.amount 为 USDT 6 位最小单位 bigint，
// 故卡数 N = amount / (10 USDT) = amount / 10_000_000n。
const CARDS_PER_DEPOSIT_UNIT = 10_000_000n; // 10 USDT × 10^6（6 位最小单位）

type TimelineItem =
  | { kind: "acquire"; key: string; date: string; cards: bigint; usdtMinUnit: bigint }
  | { kind: "burn"; key: string; date: string; days: number; destination: string; status: SimRecord["status"] };

export default function Billing() {
  const { t } = useTranslation();
  const { address } = useAccount();
  const { data: deposits } = useDepositHistory(address);
  const { data: sims } = useMySims(address);

  const items = useMemo<TimelineItem[]>(() => {
    const acquire: TimelineItem[] = (deposits ?? [])
      .filter((r: DepositRecord) => r.type === "deposit")
      .map((r) => ({
        kind: "acquire" as const,
        key: `deposit-${r.id}`,
        date: r.timestamp,
        cards: r.amount / CARDS_PER_DEPOSIT_UNIT,
        usdtMinUnit: r.amount,
      }));

    const burn: TimelineItem[] = (sims ?? []).map((s: SimRecord) => ({
      kind: "burn" as const,
      key: `sim-${s.id}`,
      date: s.createdAt,
      days: s.days,
      destination: s.destination,
      status: s.status,
    }));

    return [...acquire, ...burn].sort(
      (a, b) => new Date(b.date).getTime() - new Date(a.date).getTime(),
    );
  }, [deposits, sims]);

  if (items.length === 0) {
    return (
      <div className="px-4">
        <EmptyState
          icon={HistoryIcon}
          message={t("billing.empty")}
        />
      </div>
    );
  }

  return (
    <div className="px-4 space-y-2.5">
      {items.map((item) =>
        item.kind === "acquire" ? (
          <Card key={item.key}>
            <CardContent className="py-4 flex items-start gap-3">
              <div className="mt-0.5 flex size-9 shrink-0 items-center justify-center rounded-full bg-brand-gold/15">
                <ArrowDownToLine className="size-5 text-brand-gold" strokeWidth={1.75} />
              </div>
              <div className="flex-1 min-w-0">
                <div className="text-sm font-bold text-text-on-light-primary">
                  {t("billing.acquireCards", { count: Number(item.cards) })}
                </div>
                <div className="text-xs text-text-on-light-secondary mt-0.5">
                  {t("billing.depositedAmount", { amount: formatAmount(item.usdtMinUnit) })}
                </div>
                <div className="text-[10px] text-text-muted mt-1.5">{formatDate(item.date)}</div>
              </div>
            </CardContent>
          </Card>
        ) : (
          <Card key={item.key}>
            <CardContent className="py-4 flex items-start gap-3">
              <div className="mt-0.5 flex size-9 shrink-0 items-center justify-center rounded-full bg-brand-gold/15">
                <Flame className="size-5 text-brand-gold" strokeWidth={1.75} />
              </div>
              <div className="flex-1 min-w-0">
                <div className="flex items-center gap-2">
                  <span className="text-sm font-bold text-text-on-light-primary">
                    {t("billing.burnRedeem", { count: item.days })}
                  </span>
                  <Badge
                    variant={item.status === "confirmed" ? "secondary" : "outline"}
                    className="text-[10px]"
                  >
                    {item.status === "confirmed" ? t("billing.statusConfirmed") : t("billing.statusPending")}
                  </Badge>
                </div>
                <div className="text-xs text-text-on-light-secondary mt-0.5">
                  {t("billing.burnUsage", {
                    days: item.days,
                    destination: t(`destinations.${item.destination}`, { defaultValue: item.destination }),
                  })}
                </div>
                <div className="text-[10px] text-text-muted mt-1.5">{formatDate(item.date)}</div>
              </div>
            </CardContent>
          </Card>
        ),
      )}
    </div>
  );
}
