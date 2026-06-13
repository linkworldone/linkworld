import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useAccount } from "wagmi";
import { useTranslation } from "react-i18next";
import { useRegions, useMyNumbers } from "@/hooks/useOperator";
import { EmptyState } from "@/components/shared/EmptyState";
import { Smartphone, Search, ChevronRight } from "lucide-react";

type Tab = "regions" | "numbers";

export default function Services() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const { address } = useAccount();
  const { data: regions } = useRegions();
  const { data: numbers } = useMyNumbers(address);
  const [tab, setTab] = useState<Tab>("regions");
  const [search, setSearch] = useState("");

  const filteredRegions = regions?.filter(
    (r) =>
      r.name.toLowerCase().includes(search.toLowerCase()) ||
      r.code.toLowerCase().includes(search.toLowerCase())
  );

  return (
    <div className="px-4 space-y-3">
      <div className="flex bg-surface-card border border-surface-card-line rounded-xl p-0.5">
        <button
          onClick={() => setTab("regions")}
          className={`flex-1 py-2 rounded-lg text-[13px] font-semibold transition-colors min-h-[44px] ${
            tab === "regions" ? "bg-brand-gold text-brand-navy" : "text-text-on-light-secondary"
          }`}
        >
          {t("services.tabRegions")}
        </button>
        <button
          onClick={() => setTab("numbers")}
          className={`flex-1 py-2 rounded-lg text-[13px] font-semibold transition-colors min-h-[44px] ${
            tab === "numbers" ? "bg-brand-gold text-brand-navy" : "text-text-on-light-secondary"
          }`}
        >
          {t("services.tabNumbers")} {numbers?.length ? `(${numbers.length})` : ""}
        </button>
      </div>

      {tab === "regions" && (
        <>
          <div className="flex items-center gap-2 px-4 py-3 bg-surface-card border border-surface-card-line rounded-xl">
            <Search className="size-4 shrink-0 text-text-on-light-muted" />
            <input
              type="text"
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              placeholder={t("services.searchPlaceholder")}
              className="flex-1 bg-transparent text-sm text-text-on-light-primary outline-none placeholder:text-text-on-light-muted"
            />
          </div>
          <div className="space-y-2">
            {filteredRegions?.map((region) => (
              <button
                key={region.code}
                onClick={() => navigate(`/services/${region.code}`)}
                className="w-full flex justify-between items-center p-3.5 bg-surface-card border border-surface-card-line shadow-card rounded-xl hover:bg-surface-card-elevated transition-colors"
              >
                <div className="flex items-center gap-2.5">
                  <span className="text-2xl">{region.flag}</span>
                  <div className="text-left">
                    <div className="text-[13px] font-semibold text-text-on-light-primary">{region.name}</div>
                    <div className="text-[10px] text-text-on-light-muted mt-0.5">{t("services.operators", { count: region.operatorCount })}</div>
                  </div>
                </div>
                <div className="flex items-center gap-1.5">
                  <span className="text-xs text-brand-royal font-semibold">{t("services.fromPrice", { price: region.startingPrice })}</span>
                  <ChevronRight className="size-4 text-text-on-light-muted" />
                </div>
              </button>
            ))}
          </div>
        </>
      )}

      {tab === "numbers" && (
        <div className="space-y-2">
          {!numbers?.length ? (
            <EmptyState icon={Smartphone} message={t("services.emptyNumbers")} />
          ) : (
            numbers.map((num) => (
              <div key={num.id} className="p-4 bg-surface-card border border-surface-card-line shadow-card rounded-xl">
                <div className="flex justify-between items-start">
                  <div>
                    <div className="text-[15px] font-bold text-text-on-light-primary">{num.number}</div>
                    <div className="text-[11px] text-text-on-light-secondary mt-1">{num.operator} · {t(`destinations.${num.region}`, { defaultValue: num.region })}</div>
                  </div>
                  <span className={`text-[10px] font-semibold px-2 py-0.5 rounded-md ${
                    num.status === "active" ? "bg-status-success/20 text-status-success" : "bg-surface-input text-text-on-light-muted"
                  }`}>
                    {num.status === "active" ? t("services.statusActive") : t("services.statusInactive")}
                  </span>
                </div>
                {num.credentials && (
                  <div className="mt-3 p-2.5 bg-surface-input rounded-lg">
                    <div className="text-[10px] text-text-on-light-muted uppercase tracking-wider mb-1">
                      {num.credentials.type === "esim" ? t("services.esimConfig") : t("services.voipAccount")}
                    </div>
                    <div className="text-[11px] text-text-on-light-secondary font-mono break-all">{num.credentials.config}</div>
                  </div>
                )}
              </div>
            ))
          )}
        </div>
      )}
    </div>
  );
}
