import { useNavigate, useLocation } from "react-router-dom";
import { useAccount } from "wagmi";
import { shortenAddress } from "@/utils/format";
import { useUnreadCount } from "@/hooks/useNotification";

export function Header() {
  const navigate = useNavigate();
  const location = useLocation();
  const { address } = useAccount();
  const { data: unreadCount } = useUnreadCount(address);
  const isDashboard = location.pathname === "/dashboard";

  if (isDashboard) {
    return (
      <header className="px-4 py-3 flex justify-between items-center">
        <div>
          <div className="text-[11px] text-text-muted">Welcome back</div>
          <div className="text-[15px] font-semibold">{address ? shortenAddress(address) : ""}</div>
        </div>
        <div className="flex items-center gap-3">
          <button onClick={() => navigate("/notifications")} className="relative text-xl min-w-[44px] min-h-[44px] flex items-center justify-center">
            {"\u{1F514}"}
            {unreadCount && unreadCount > 0 && (
              <span className="absolute top-1 right-1 bg-status-danger text-white text-[8px] w-3.5 h-3.5 rounded-full flex items-center justify-center">{unreadCount}</span>
            )}
          </button>
          <div className="w-8 h-8 rounded-full bg-gradient-to-br from-brand-blue to-brand-purple" />
        </div>
      </header>
    );
  }

  const pageTitles: Record<string, string> = {
    "/deposit": "Deposit",
    "/services": "Services",
    "/billing": "Billing",
    "/notifications": "Notifications",
  };

  const isSubPage = location.pathname.includes("/services/") || location.pathname.includes("/billing/");
  const title = pageTitles[location.pathname] || "";

  return (
    <header className="px-4 py-3 flex items-center gap-3">
      {isSubPage && (
        <button onClick={() => navigate(-1)} className="text-lg min-w-[44px] min-h-[44px] flex items-center justify-center">{"\u2190"}</button>
      )}
      <h1 className="text-[17px] font-bold">{title}</h1>
    </header>
  );
}
