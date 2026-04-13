import { useNavigate } from "react-router-dom";
import { useAccount } from "wagmi";
import { useUser } from "@/hooks/useUser";
import { useDeposit } from "@/hooks/useDeposit";
import { useMonthEstimate } from "@/hooks/useBilling";
import { useMyNumbers } from "@/hooks/useOperator";
import { StatusBadge } from "@/components/shared/StatusBadge";
import { AmountDisplay } from "@/components/shared/AmountDisplay";

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

  const quickActions = [
    { icon: "\u{1F4B0}", label: "Top Up", path: "/deposit" },
    { icon: "\u{1F30D}", label: "Switch Region", path: "/services" },
    { icon: "\u{1F4C4}", label: "Bills", path: "/billing" },
    { icon: "\u{1F4B3}", label: "Pay Now", path: "/billing" },
  ];

  return (
    <div className="px-4 space-y-3">
      {/* Status Card */}
      <div className="p-4 rounded-2xl bg-gradient-to-br from-surface-gradient-from to-surface-gradient-to">
        <div className="grid grid-cols-2 gap-3">
          <div>
            <div className="text-[10px] text-text-muted">Account Status</div>
            <div className="mt-1.5">{user && <StatusBadge status={user.status} />}</div>
          </div>
          <div className="text-right">
            <div className="text-[10px] text-text-muted">Deposit Balance</div>
            <div className="mt-1">
              {deposit && <AmountDisplay amount={deposit.balance} currency={deposit.currency} size="md" />}
            </div>
          </div>
          <div>
            <div className="text-[10px] text-text-muted">Virtual Number</div>
            <div className="text-sm font-semibold mt-1">{activeNumber?.number || "\u2014"}</div>
          </div>
          <div className="text-right">
            <div className="text-[10px] text-text-muted">Region</div>
            <div className="text-sm font-semibold mt-1">
              {region ? `${regionFlags[region] || ""} ${region}` : "\u2014"}
            </div>
          </div>
        </div>
      </div>

      {/* Usage Card */}
      <div className="p-4 rounded-xl bg-surface-card flex justify-around">
        <div className="text-center">
          <div className="text-[10px] text-text-muted">Data Used</div>
          <div className="text-[22px] font-extrabold text-brand-blue mt-1">{estimate?.dataGB ?? "\u2014"}</div>
          <div className="text-[10px] text-text-muted">GB</div>
        </div>
        <div className="w-px bg-surface-secondary" />
        <div className="text-center">
          <div className="text-[10px] text-text-muted">Calls</div>
          <div className="text-[22px] font-extrabold text-brand-purple mt-1">{estimate?.callMinutes ?? "\u2014"}</div>
          <div className="text-[10px] text-text-muted">min</div>
        </div>
        <div className="w-px bg-surface-secondary" />
        <div className="text-center">
          <div className="text-[10px] text-text-muted">Est. Bill</div>
          <div className="text-[22px] font-extrabold text-status-warning mt-1">
            {estimate ? `$${estimate.estimatedCost}` : "\u2014"}
          </div>
          <div className="text-[10px] text-text-muted">this month</div>
        </div>
      </div>

      {/* Quick Actions */}
      <div className="grid grid-cols-2 gap-2.5">
        {quickActions.map((action) => (
          <button
            key={action.label}
            onClick={() => navigate(action.path)}
            className="bg-surface-card rounded-xl p-4 text-center hover:bg-surface-secondary transition-colors"
          >
            <div className="text-[22px] mb-1.5">{action.icon}</div>
            <div className="text-xs font-semibold">{action.label}</div>
          </button>
        ))}
      </div>

      {/* Overdue Warning */}
      {user?.status === "suspended" && (
        <div className="p-3 bg-status-danger/15 border border-status-danger/30 rounded-xl flex items-center gap-2">
          <span className="text-base">{"\u26A0\uFE0F"}</span>
          <div className="text-xs text-status-danger">
            Your account is suspended due to unpaid bills. Service will be terminated if not settled within 14 days.
          </div>
        </div>
      )}
    </div>
  );
}
