import { useLocation, useNavigate } from "react-router-dom";
import { useAccount } from "wagmi";
import { Home, Smartphone, Wallet, ReceiptText, CreditCard } from "lucide-react";
import { useBills } from "@/hooks/useBilling";

const tabs = [
  { label: "Home", icon: Home, path: "/dashboard" },
  { label: "Services", icon: Smartphone, path: "/services" },
  { label: "Deposit", icon: Wallet, path: "/deposit" },
  { label: "Bills", icon: ReceiptText, path: "/billing", badgeKey: "unpaidBills" as const },
  { label: "Cards", icon: CreditCard, path: "/cards" },
];

export function TabBar() {
  const location = useLocation();
  const navigate = useNavigate();
  const { address } = useAccount();
  const { data: unpaidBills } = useBills(address, "unpaid");

  const badges: Record<string, number | undefined> = {
    unpaidBills: unpaidBills?.length,
  };

  return (
    <nav className="fixed bottom-0 left-0 right-0 max-w-mobile mx-auto bg-surface border-t border-border pb-[env(safe-area-inset-bottom)] z-50">
      <div className="flex justify-around py-2">
        {tabs.map((tab) => {
          const isActive = location.pathname.startsWith(tab.path);
          const badge = tab.badgeKey ? badges[tab.badgeKey] : undefined;
          const Icon = tab.icon;
          return (
            <button
              key={tab.path}
              onClick={() => navigate(tab.path)}
              className={`flex flex-col items-center gap-0.5 min-w-[44px] min-h-[44px] justify-center relative ${
                isActive ? "text-brand-gold" : "text-text-on-dark-muted"
              }`}
            >
              <span className="relative">
                <Icon className="size-5" strokeWidth={isActive ? 2.25 : 1.75} />
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
