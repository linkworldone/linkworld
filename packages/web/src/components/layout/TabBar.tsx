import { useLocation, useNavigate } from "react-router-dom";
import { Home, Wallet, History, CreditCard } from "lucide-react";

const tabs = [
  { label: "Home", icon: Home, path: "/dashboard" },
  // Services 暂时隐藏（功能未打通）；恢复时把这项加回，并重新 import Smartphone：{ label: "Services", icon: Smartphone, path: "/services" }
  { label: "Deposit", icon: Wallet, path: "/deposit" },
  { label: "History", icon: History, path: "/billing" },
  { label: "Cards", icon: CreditCard, path: "/cards" },
];

export function TabBar() {
  const location = useLocation();
  const navigate = useNavigate();

  return (
    <nav className="fixed bottom-0 left-0 right-0 max-w-mobile mx-auto bg-surface border-t border-border pb-[env(safe-area-inset-bottom)] z-50">
      <div className="flex justify-around py-2">
        {tabs.map((tab) => {
          const isActive = location.pathname.startsWith(tab.path);
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
              </span>
              <span className={`text-[9px] ${isActive ? "font-semibold" : ""}`}>{tab.label}</span>
            </button>
          );
        })}
      </div>
    </nav>
  );
}
