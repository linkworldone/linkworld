import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { useAccount } from "wagmi";
import { useTranslation } from "react-i18next";
import { Globe } from "lucide-react";
import { ConnectButton } from "@/components/wallet/ConnectButton";
import { RegisterSheet } from "@/components/wallet/RegisterSheet";
import { LanguageToggle } from "@/components/shared/LanguageToggle";
import { useUser } from "@/hooks/useUser";

export default function Landing() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const { address, isConnected } = useAccount();
  const { data: user, isLoading } = useUser(address);
  const [showRegister, setShowRegister] = useState(false);

  useEffect(() => {
    if (!isConnected || isLoading) return;
    if (user) {
      navigate("/dashboard", { replace: true });
    } else if (isConnected && !user) {
      setShowRegister(true);
    }
  }, [isConnected, user, isLoading, navigate]);

  return (
    <div className="min-h-screen bg-surface text-text-primary flex flex-col max-w-mobile mx-auto">
      <div className="px-5 py-4 flex justify-between items-center">
        <span className="font-display text-lg font-extrabold text-brand-gold">
          LinkWorld
        </span>
        <div className="flex items-center gap-3">
          <LanguageToggle />
          <ConnectButton />
        </div>
      </div>

      <div className="flex-1 flex flex-col items-center justify-center px-6 text-center">
        <div className="w-[72px] h-[72px] rounded-full bg-gradient-hero flex items-center justify-center mb-7 ring-1 ring-brand-gold/40">
          <Globe className="size-9 text-brand-gold" strokeWidth={1.5} />
        </div>
        <h1 className="text-[32px] font-extrabold leading-tight mb-3">
          {t("landing.title")}<br />{t("landing.titleLine2")}
        </h1>
        <p className="text-sm text-text-secondary leading-relaxed mb-9">
          {t("landing.subtitle")}
        </p>
        <ConnectButton label={t("common.getStarted")} />
        <button className="mt-3 text-sm text-text-secondary border border-border rounded-xl px-12 py-3 w-full max-w-[280px]">
          {t("common.learnMore")}
        </button>
      </div>

      <div className="flex justify-around py-6 border-t border-border">
        <div className="text-center">
          <div className="text-xl font-extrabold text-brand-gold">{t("common.countriesValue")}</div>
          <div className="text-[11px] text-text-muted mt-0.5">{t("common.countries")}</div>
        </div>
        <div className="text-center">
          <div className="text-xl font-extrabold text-status-success">{t("common.platformFeeValue")}</div>
          <div className="text-[11px] text-text-muted mt-0.5">{t("common.platformFee")}</div>
        </div>
        <div className="text-center">
          <div className="text-xl font-extrabold text-brand-gold">{t("common.kycRequiredValue")}</div>
          <div className="text-[11px] text-text-muted mt-0.5">{t("common.kycRequired")}</div>
        </div>
      </div>

      <div className="text-center py-3 text-[10px] text-text-muted">
        {t("landing.poweredBy")}
      </div>

      {address && (
        <RegisterSheet
          address={address}
          open={showRegister}
          onClose={() => setShowRegister(false)}
          onSuccess={() => navigate("/dashboard", { replace: true })}
        />
      )}
    </div>
  );
}
