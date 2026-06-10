import { useNavigate } from "react-router-dom";
import { useAccount } from "wagmi";
import { Wallet, Globe, ReceiptText, CreditCard, AlertTriangle, type LucideIcon } from "lucide-react";
import { useUser } from "@/hooks/useUser";
import { useDeposit } from "@/hooks/useDeposit";
import { useMonthEstimate } from "@/hooks/useBilling";
import { useMyNumbers } from "@/hooks/useOperator";
import { StatusBadge } from "@/components/shared/StatusBadge";
import { AmountDisplay } from "@/components/shared/AmountDisplay";
import { Card, CardContent } from "@/components/ui/card";

export default function Dashboard() {
  const navigate = useNavigate();
  const { address } = useAccount();
  const { data: user } = useUser(address);
  const { data: deposit } = useDeposit(address);
  const { data: estimate } = useMonthEstimate(address);
  const { data: numbers } = useMyNumbers(address);

  const activeNumber = numbers?.find((n) => n.status === "active");
  const region = activeNumber?.region;

  const regionFlags: Record<string, string> = {
    JP: "\u{1F1EF}\u{1F1F5}", US: "\u{1F1FA}\u{1F1F8}", GB: "\u{1F1EC}\u{1F1E7}", SG: "\u{1F1F8}\u{1F1EC}", KR: "\u{1F1F0}\u{1F1F7}", DE: "\u{1F1E9}\u{1F1EA}", AU: "\u{1F1E6}\u{1F1FA}", TH: "\u{1F1F9}\u{1F1ED}",
  };

  const quickActions: { icon: LucideIcon; label: string; path: string }[] = [
    { icon: Wallet, label: "Top Up", path: "/deposit" },
    { icon: Globe, label: "Switch Region", path: "/services" },
    { icon: ReceiptText, label: "Bills", path: "/billing" },
    { icon: CreditCard, label: "Pay Now", path: "/billing" },
  ];

  return (
    <div className="px-4 space-y-3">
      {/* Status Card\uff08navy \u6e10\u53d8\u6df1\u5e95\uff0c\u6587\u5b57 on-dark\uff0c\u91d1\u989d gold-on-dark\uff09 */}
      <div className="p-4 rounded-2xl bg-gradient-hero ring-1 ring-brand-gold/25">
        <div className="grid grid-cols-2 gap-3">
          <div>
            <div className="text-[10px] text-text-on-dark-muted">Account Status</div>
            <div className="mt-1.5">{user && <StatusBadge status={user.status} />}</div>
          </div>
          <div className="text-right">
            <div className="text-[10px] text-text-on-dark-muted">Deposit Balance</div>
            <div className="mt-1">
              {deposit && (
                <AmountDisplay amount={deposit.balance} currency={deposit.currency} size="md" tone="gold-on-dark" />
              )}
            </div>
          </div>
          <div>
            <div className="text-[10px] text-text-on-dark-muted">Virtual Number</div>
            <div className="text-sm font-semibold mt-1 text-text-on-dark-primary">{activeNumber?.number || "\u2014"}</div>
          </div>
          <div className="text-right">
            <div className="text-[10px] text-text-on-dark-muted">Region</div>
            <div className="text-sm font-semibold mt-1 text-text-on-dark-primary">
              {region ? `${regionFlags[region] || ""} ${region}` : "\u2014"}
            </div>
          </div>
        </div>
      </div>

      {/* Usage Card\uff08\u6696\u7c73\u767d ui/card\uff0c\u6587\u5b57 on-light\uff0c\u6570\u503c navy/\u6a59\uff09 */}
      <Card>
        <CardContent className="flex justify-around">
          <div className="text-center">
            <div className="text-[10px] text-text-on-light-muted">Data Used</div>
            <div className="text-[22px] font-extrabold text-text-on-light-primary mt-1 font-data tabular-nums">{estimate?.dataGB ?? "\u2014"}</div>
            <div className="text-[10px] text-text-on-light-muted">GB</div>
          </div>
          <div className="w-px bg-surface-card-line" />
          <div className="text-center">
            <div className="text-[10px] text-text-on-light-muted">Calls</div>
            <div className="text-[22px] font-extrabold text-text-on-light-primary mt-1 font-data tabular-nums">{estimate?.callMinutes ?? "\u2014"}</div>
            <div className="text-[10px] text-text-on-light-muted">min</div>
          </div>
          <div className="w-px bg-surface-card-line" />
          <div className="text-center">
            <div className="text-[10px] text-text-on-light-muted">Est. Bill</div>
            <div className="text-[22px] font-extrabold text-status-warning mt-1 font-data tabular-nums">
              {estimate ? `$${estimate.estimatedCost}` : "\u2014"}
            </div>
            <div className="text-[10px] text-text-on-light-muted">this month</div>
          </div>
        </CardContent>
      </Card>

      {/* Quick Actions（暖米白卡 + 金描线图标 + navy 文字） */}
      <div className="grid grid-cols-2 gap-2.5">
        {quickActions.map((action) => {
          const Icon = action.icon;
          return (
            <button
              key={action.label}
              onClick={() => navigate(action.path)}
              className="bg-surface-card border border-surface-card-line shadow-card rounded-xl p-4 text-center hover:bg-surface-card-elevated transition-colors min-h-[44px]"
            >
              <Icon className="size-6 mb-1.5 mx-auto text-brand-gold" strokeWidth={1.75} />
              <div className="text-xs font-semibold text-text-on-light-primary">{action.label}</div>
            </button>
          );
        })}
      </div>

      {/* Overdue Warning */}
      {user?.status === "suspended" && (
        <div className="p-3 bg-status-danger/15 border border-status-danger/30 rounded-xl flex items-center gap-2">
          <AlertTriangle className="size-4 shrink-0 text-status-danger" />
          <div className="text-xs text-status-danger">
            Your account is suspended due to unpaid bills. Service will be terminated if not settled within 14 days.
          </div>
        </div>
      )}
    </div>
  );
}
