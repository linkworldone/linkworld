import { Outlet, Navigate, useLocation } from "react-router-dom";
import { useAccount } from "wagmi";
import { useTranslation } from "react-i18next";
import { useUser } from "@/hooks/useUser";
import { Header } from "./Header";
import { TabBar } from "./TabBar";
import { GuardCard } from "@/components/shared/GuardCard";
import { Wallet, ShieldAlert } from "lucide-react";

export function AppLayout() {
  const { t } = useTranslation();
  const { address, isConnected } = useAccount();
  const { data: user, isLoading } = useUser(address);
  const location = useLocation();

  if (!isConnected) return <Navigate to="/" replace />;

  if (isLoading) {
    return (
      <div className="max-w-mobile mx-auto min-h-screen bg-surface flex items-center justify-center">
        <div className="text-text-secondary text-sm">{t("common.loading")}</div>
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
          <GuardCard icon={Wallet} title={t("guard.depositRequiredTitle")} message={t("guard.depositRequiredMessage")} actionLabel={t("guard.goToDeposit")} actionPath="/deposit" />
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
          <GuardCard icon={ShieldAlert} title={t("guard.suspendedTitle")} message={t("guard.suspendedMessage")} actionLabel={t("guard.goToBilling")} actionPath="/billing" />
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
