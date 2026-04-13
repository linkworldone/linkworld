import { useLocation, useNavigate } from "react-router-dom";
import { useAccount } from "wagmi";
import { useUnreadCount } from "@/hooks/useNotification";
import { useBills } from "@/hooks/useBilling";

const tabs = [
  { label: "Home", icon: "\u{1F3E0}", path: "/dashboard" },
  { label: "Services", icon: "\u{1F4F1}", path: "/services" },
  { label: "Deposit", icon: "\u{1F4B0}", path: "/deposit" },
  { label: "Bills", icon: "\u{1F4C4}", path: "/billing", badgeKey: "unpaidBills" as const },
  { label: "Alerts", icon: "\u{1F514}", path: "/notifications", badgeKey: "unread" as const },
];

export function TabBar() {
  const location = useLocation();
  const navigate = useNavigate();
  const { address } = useAccount();
  const { data: unreadCount } = useUnreadCount(address);
  const { data: unpaidBills } = useBills(address, "unpaid");

  const badges: Record<string, number | undefined> = {
    unread: unreadCount,
    unpaidBills: unpaidBills?.length,
  };

  return (
    <nav className="fixed bottom-0 left-0 right-0 max-w-mobile mx-auto bg-surface border-t border-border pb-[env(safe-area-inset-bottom)] z-50">
      <div className="flex justify-around py-2">
        {tabs.map((tab) => {
          const isActive = location.pathname.startsWith(tab.path);
          const badge = tab.badgeKey ? badges[tab.badgeKey] : undefined;
          return (
            <button
              key={tab.path}
              onClick={() => navigate(tab.path)}
              className={`flex flex-col items-center gap-0.5 min-w-[44px] min-h-[44px] justify-center relative ${
                isActive ? "text-brand-blue" : "text-text-muted"
              }`}
            >
              <span className="text-xl relative">
                {tab.icon}
                {badge && badge > 0 && (
                  <span className="absolute -top-1 -right-2 bg-status-danger text-white text-[8px] w-3.5 h-3.5 rounded-full flex items-center justify-center font-semibold">
                    {badge > 9 ? "9+" : badge}
                  </span>
                )}
              </span>
              <span className={`text-[9px] ${isActive ? "font-semibold" : ""}`}>{tab.label}</span>
            </button>
          );
        })}
      </div>
    </nav>
  );
}
