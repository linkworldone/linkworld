import { Outlet, Navigate, useLocation } from "react-router-dom";
import { useAccount } from "wagmi";
import { useUser } from "@/hooks/useUser";
import { Header } from "./Header";
import { TabBar } from "./TabBar";
import { GuardCard } from "@/components/shared/GuardCard";

export function AppLayout() {
  const { address, isConnected } = useAccount();
  const { data: user, isLoading } = useUser(address);
  const location = useLocation();

  if (!isConnected) return <Navigate to="/" replace />;

  if (isLoading) {
    return (
      <div className="max-w-mobile mx-auto min-h-screen bg-surface flex items-center justify-center">
        <div className="text-text-secondary text-sm">Loading...</div>
      </div>
    );
  }

  if (!user) return <Navigate to="/" replace />;

  const path = location.pathname;
  const status = user.status;

  if (status === "inactive") {
    const restricted = ["/services", "/billing", "/notifications"];
    if (restricted.some((r) => path.startsWith(r))) {
      return (
        <div className="max-w-mobile mx-auto min-h-screen bg-surface">
          <Header />
          <GuardCard icon={"\u{1F4B0}"} title="Deposit Required" message="You need to deposit funds before accessing this feature." actionLabel="Go to Deposit" actionPath="/deposit" />
          <TabBar />
        </div>
      );
    }
  }

  if (status === "suspended") {
    const restricted = ["/services", "/deposit"];
    if (restricted.some((r) => path.startsWith(r))) {
      return (
        <div className="max-w-mobile mx-auto min-h-screen bg-surface">
          <Header />
          <GuardCard icon={"\u26A0\uFE0F"} title="Account Suspended" message="Please settle your outstanding bills to restore access." actionLabel="Go to Billing" actionPath="/billing" />
          <TabBar />
        </div>
      );
    }
  }

  return (
    <div className="max-w-mobile mx-auto min-h-screen bg-surface flex flex-col">
      <Header />
      <main className="flex-1 overflow-y-auto pb-[80px]">
        <Outlet />
      </main>
      <TabBar />
    </div>
  );
}
