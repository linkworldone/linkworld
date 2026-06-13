import { useNavigate, useLocation } from "react-router-dom";
import { useAccount } from "wagmi";
import { Bell, ChevronLeft } from "lucide-react";
import { useTranslation } from "react-i18next";
import { shortenAddress } from "@/utils/format";
import { useUnreadCount } from "@/hooks/useNotification";

export function Header() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const location = useLocation();
  const { address } = useAccount();
  const { data: unreadCount } = useUnreadCount(address);
  const isDashboard = location.pathname === "/dashboard";

  if (isDashboard) {
    return (
      <header className="px-4 py-3 flex justify-between items-center">
        <div>
          <div className="text-[11px] text-text-muted">{t("header.welcomeBack")}</div>
          <div className="text-[15px] font-semibold">{address ? shortenAddress(address) : ""}</div>
        </div>
        <div className="flex items-center gap-3">
          <button onClick={() => navigate("/notifications")} className="relative min-w-[44px] min-h-[44px] flex items-center justify-center text-text-on-dark-secondary">
            <Bell className="size-5" strokeWidth={1.75} />
            {(unreadCount ?? 0) > 0 && (
              <span className="absolute top-1.5 right-1.5 bg-status-danger text-white text-[8px] w-3.5 h-3.5 rounded-full flex items-center justify-center">{unreadCount}</span>
            )}
          </button>
          <div className="w-8 h-8 rounded-full bg-gradient-to-br from-brand-royal to-brand-gold" />
        </div>
      </header>
    );
  }

  const pageTitles: Record<string, string> = {
    "/deposit": t("nav.deposit"),
    "/services": t("nav.services"),
    "/billing": t("nav.history"),
    "/notifications": t("nav.notifications"),
  };

  const isSubPage = location.pathname.includes("/services/") || location.pathname.includes("/billing/");
  const title = pageTitles[location.pathname] || "";

  return (
    <header className="px-4 py-3 flex items-center gap-3">
      {isSubPage && (
        <button onClick={() => navigate(-1)} className="min-w-[44px] min-h-[44px] flex items-center justify-center text-text-on-dark-secondary">
          <ChevronLeft className="size-5" />
        </button>
      )}
      <h1 className="text-[17px] font-bold">{title}</h1>
    </header>
  );
}
