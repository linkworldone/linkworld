import { useNavigate } from "react-router-dom";
import { useAccount } from "wagmi";
import { Wallet, CreditCard, History, AlertTriangle, type LucideIcon } from "lucide-react";
import { useUser } from "@/hooks/useUser";
import { useDeposit } from "@/hooks/useDeposit";
import { useMyNumbers } from "@/hooks/useOperator";
import { useTrafficCards } from "@/hooks/contracts";
import { useMySims } from "@/hooks/useSim";
import { StatusBadge } from "@/components/shared/StatusBadge";
import { AmountDisplay } from "@/components/shared/AmountDisplay";
import { Card, CardContent } from "@/components/ui/card";

export default function Dashboard() {
  const navigate = useNavigate();
  const { address } = useAccount();
  const { data: user } = useUser(address);
  const { data: deposit } = useDeposit(address);
  const { data: numbers } = useMyNumbers(address);
  const { cards } = useTrafficCards(address);
  const { data: sims } = useMySims(address);

  const activeNumber = numbers?.find((n) => n.status === "active");
  const region = activeNumber?.region;

  const heldCards = cards?.length ?? 0;
  const simDays = (sims ?? []).reduce((sum, s) => sum + (s.days ?? 0), 0);

  const regionFlags: Record<string, string> = {
    JP: "\u{1F1EF}\u{1F1F5}", US: "\u{1F1FA}\u{1F1F8}", GB: "\u{1F1EC}\u{1F1E7}", SG: "\u{1F1F8}\u{1F1EC}", KR: "\u{1F1F0}\u{1F1F7}", DE: "\u{1F1E9}\u{1F1EA}", AU: "\u{1F1E6}\u{1F1FA}", TH: "\u{1F1F9}\u{1F1ED}",
  };

  const quickActions: { icon: LucideIcon; label: string; path: string }[] = [
    { icon: Wallet, label: "充值", path: "/deposit" },
    { icon: CreditCard, label: "流量卡", path: "/cards" },
    { icon: History, label: "历史", path: "/billing" },
  ];

  return (
    <div className="px-4 space-y-3">
      {/* 账户状态卡（navy 渐变深底，文字 on-dark，金额 gold-on-dark） */}
      <div className="p-4 rounded-2xl bg-gradient-hero ring-1 ring-brand-gold/25">
        <div className="grid grid-cols-2 gap-3">
          <div>
            <div className="text-[10px] text-text-on-dark-muted">账户状态</div>
            <div className="mt-1.5">{user && <StatusBadge status={user.status} />}</div>
          </div>
          <div className="text-right">
            <div className="text-[10px] text-text-on-dark-muted">保证金余额</div>
            <div className="mt-1">
              {deposit && (
                <AmountDisplay amount={deposit.balance} currency={deposit.currency} size="md" tone="gold-on-dark" />
              )}
            </div>
          </div>
          <div>
            <div className="text-[10px] text-text-on-dark-muted">虚拟号码</div>
            <div className="text-sm font-semibold mt-1 text-text-on-dark-primary">{activeNumber?.number || "—"}</div>
          </div>
          <div className="text-right">
            <div className="text-[10px] text-text-on-dark-muted">地区</div>
            <div className="text-sm font-semibold mt-1 text-text-on-dark-primary">
              {region ? `${regionFlags[region] || ""} ${region}` : "—"}
            </div>
          </div>
        </div>
      </div>

      {/* 资产卡（持有流量卡 + SIM 天数）暖米白 ui/card，文字 on-light */}
      <Card>
        <CardContent className="flex justify-around">
          <div className="text-center">
            <div className="text-[10px] text-text-on-light-muted">持有流量卡</div>
            <div className="text-[22px] font-extrabold text-text-on-light-primary mt-1 font-data tabular-nums">{heldCards}</div>
            <div className="text-[10px] text-text-on-light-muted">张 · 无限流量</div>
          </div>
          <div className="w-px bg-surface-card-line" />
          <div className="text-center">
            <div className="text-[10px] text-text-on-light-muted">SIM 天数</div>
            <div className="text-[22px] font-extrabold text-brand-royal mt-1 font-data tabular-nums">{simDays}</div>
            <div className="text-[10px] text-text-on-light-muted">天 · 已兑换</div>
          </div>
        </CardContent>
      </Card>

      {/* 快捷入口（暖米白卡 + 金描线图标 + navy 文字） */}
      <div className="grid grid-cols-3 gap-2.5">
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

      {/* 账户暂停提示 */}
      {user?.status === "suspended" && (
        <div className="p-3 bg-status-danger/15 border border-status-danger/30 rounded-xl flex items-center gap-2">
          <AlertTriangle className="size-4 shrink-0 text-status-danger" />
          <div className="text-xs text-status-danger">
            账户已暂停，请联系客服恢复服务。
          </div>
        </div>
      )}
    </div>
  );
}
