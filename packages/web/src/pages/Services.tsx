import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useAccount } from "wagmi";
import { useRegions, useMyNumbers } from "@/hooks/useOperator";
import { EmptyState } from "@/components/shared/EmptyState";

type Tab = "regions" | "numbers";

export default function Services() {
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
      <div className="flex bg-surface-card rounded-xl p-0.5">
        <button
          onClick={() => setTab("regions")}
          className={`flex-1 py-2 rounded-lg text-[13px] font-semibold transition-colors ${
            tab === "regions" ? "bg-brand-blue text-white" : "text-text-muted"
          }`}
        >
          Regions
        </button>
        <button
          onClick={() => setTab("numbers")}
          className={`flex-1 py-2 rounded-lg text-[13px] font-semibold transition-colors ${
            tab === "numbers" ? "bg-brand-blue text-white" : "text-text-muted"
          }`}
        >
          My Numbers {numbers?.length ? `(${numbers.length})` : ""}
        </button>
      </div>

      {tab === "regions" && (
        <>
          <div className="flex items-center gap-2 px-4 py-3 bg-surface-card rounded-xl">
            <span className="text-text-muted">🔍</span>
            <input
              type="text"
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              placeholder="Search country or region..."
              className="flex-1 bg-transparent text-sm text-text-primary outline-none placeholder:text-text-muted"
            />
          </div>
          <div className="space-y-2">
            {filteredRegions?.map((region) => (
              <button
                key={region.code}
                onClick={() => navigate(`/services/${region.code}`)}
                className="w-full flex justify-between items-center p-3.5 bg-surface-card rounded-xl hover:bg-surface-secondary transition-colors"
              >
                <div className="flex items-center gap-2.5">
                  <span className="text-2xl">{region.flag}</span>
                  <div className="text-left">
                    <div className="text-[13px] font-semibold">{region.name}</div>
                    <div className="text-[10px] text-text-muted mt-0.5">{region.operatorCount} operators</div>
                  </div>
                </div>
                <div className="text-right">
                  <div className="text-xs text-brand-blue font-semibold">from ${region.startingPrice}/mo</div>
                  <div className="text-text-muted text-sm mt-0.5">→</div>
                </div>
              </button>
            ))}
          </div>
        </>
      )}

      {tab === "numbers" && (
        <div className="space-y-2">
          {!numbers?.length ? (
            <EmptyState icon="📱" message="No virtual numbers yet. Browse regions to apply." />
          ) : (
            numbers.map((num) => (
              <div key={num.id} className="p-4 bg-surface-card rounded-xl">
                <div className="flex justify-between items-start">
                  <div>
                    <div className="text-[15px] font-bold">{num.number}</div>
                    <div className="text-[11px] text-text-secondary mt-1">{num.operator} · {num.region}</div>
                  </div>
                  <span className={`text-[10px] font-semibold px-2 py-0.5 rounded-md ${
                    num.status === "active" ? "bg-status-success/20 text-status-success" : "bg-surface-secondary text-text-muted"
                  }`}>
                    {num.status}
                  </span>
                </div>
                {num.credentials && (
                  <div className="mt-3 p-2.5 bg-surface-secondary rounded-lg">
                    <div className="text-[10px] text-text-muted uppercase tracking-wider mb-1">
                      {num.credentials.type === "esim" ? "eSIM Config" : "VoIP Account"}
                    </div>
                    <div className="text-[11px] text-text-secondary font-mono break-all">{num.credentials.config}</div>
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
