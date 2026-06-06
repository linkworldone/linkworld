import { useEffect } from "react";
import { useAccount } from "wagmi";
import { Button } from "@/components/ui/button";
import {
  useTrafficCardCredit,
  useTrafficCards,
  useBurnCard,
  useIssueMonthlyCards,
  type TrafficCardItem,
} from "@/hooks/contracts";

function formatBytes(bytes: bigint): string {
  if (bytes === 0n) return "0 B";
  const mb = bytes / (1024n * 1024n);
  if (mb > 0n) return `${mb.toString()} MB`;
  const kb = bytes / 1024n;
  if (kb > 0n) return `${kb.toString()} KB`;
  return `${bytes.toString()} B`;
}

function formatTimestamp(ts: bigint): string {
  if (ts === 0n) return "—";
  return new Date(Number(ts) * 1000).toLocaleDateString();
}

export default function Cards() {
  const { address } = useAccount();
  const credit = useTrafficCardCredit(address);
  const { cards, isLoading, refetch: refetchCards } = useTrafficCards(address);
  const burn = useBurnCard();
  const issue = useIssueMonthlyCards();

  // Refresh after burn/issue confirmed
  useEffect(() => {
    if (burn.isSuccess) {
      refetchCards();
      credit.refetch();
      burn.reset();
    }
  }, [burn.isSuccess]);

  useEffect(() => {
    if (issue.isSuccess) {
      refetchCards();
      issue.reset();
    }
  }, [issue.isSuccess]);

  const burnBusy = burn.isPending || burn.isConfirming;
  const issueBusy = issue.isPending || issue.isConfirming;

  return (
    <div className="px-4 space-y-4">
      {/* Available Credit */}
      <div className="p-6 rounded-2xl bg-gradient-to-br from-surface-gradient-from to-surface-gradient-to">
        <div className="text-[10px] text-text-muted uppercase tracking-wider">
          Available Credit
        </div>
        <div className="mt-2 flex items-baseline justify-between">
          <div className="text-3xl font-extrabold text-brand-blue">
            {formatBytes(credit.balance)}
          </div>
          <div className="text-[11px] text-text-muted text-right">
            {credit.expiry > 0n ? (
              <>
                <div>Expires</div>
                <div
                  className={
                    credit.isExpired ? "text-status-danger font-semibold" : "font-semibold"
                  }
                >
                  {formatTimestamp(credit.expiry)}
                </div>
              </>
            ) : (
              <span>No active credit</span>
            )}
          </div>
        </div>
      </div>

      {/* Cards list */}
      <div>
        <h3 className="text-[13px] font-semibold mb-3">
          Cards in wallet ({cards.length})
        </h3>

        {isLoading && (
          <div className="text-center text-xs text-text-muted py-6">Loading cards...</div>
        )}

        {!isLoading && cards.length === 0 && (
          <div className="text-center text-xs text-text-muted py-6 bg-surface-card rounded-xl">
            No active traffic cards. Issue one below.
          </div>
        )}

        <div className="space-y-2">
          {cards.map((card: TrafficCardItem) => (
            <div
              key={card.tokenId.toString()}
              className="p-4 bg-surface-card rounded-xl flex items-center justify-between"
            >
              <div>
                <div className="text-xs text-text-muted">
                  Card #{card.tokenId.toString()}
                </div>
                <div className="text-lg font-bold mt-0.5">
                  {formatBytes(card.dataAmount)}
                </div>
                <div className="text-[10px] text-text-muted mt-0.5">
                  Minted: {formatTimestamp(card.createdAt)}
                </div>
              </div>
              <Button
                size="sm"
                variant="outline"
                disabled={burnBusy}
                onClick={() => burn.burnCard(card.tokenId)}
              >
                {burnBusy ? "Burning..." : "Burn for Credit"}
              </Button>
            </div>
          ))}
        </div>
      </div>

      {/* Issue monthly card (admin/self) */}
      <div className="pt-2">
        <Button
          className="w-full py-3"
          disabled={!address || issueBusy}
          onClick={() => address && issue.issue([address])}
        >
          {issueBusy ? "Issuing..." : "Issue Monthly Card (Admin)"}
        </Button>
        <p className="text-[10px] text-text-muted mt-2 text-center">
          Mints a card if you have deposit &gt; 0 and haven&apos;t claimed this month.
        </p>
        {issue.error && (
          <p className="text-[10px] text-status-danger mt-2 text-center">
            {issue.error.message.split("\n")[0]}
          </p>
        )}
        {burn.error && (
          <p className="text-[10px] text-status-danger mt-2 text-center">
            {burn.error.message.split("\n")[0]}
          </p>
        )}
      </div>
    </div>
  );
}
