import { useState } from "react";
import { useAccount } from "wagmi";
import { Drawer } from "vaul";
import { Button } from "@/components/ui/button";
import { useDeposit, useDepositHistory, useDepositMutation, useWithdrawMutation } from "@/hooks/useDeposit";
import { AmountDisplay } from "@/components/shared/AmountDisplay";
import { formatAmount, formatDate } from "@/utils/format";

type SheetMode = "deposit" | "withdraw" | null;

export default function Deposit() {
  const { address } = useAccount();
  const { data: deposit } = useDeposit(address);
  const { data: history } = useDepositHistory(address);
  const depositMutation = useDepositMutation();
  const withdrawMutation = useWithdrawMutation();

  const [sheetMode, setSheetMode] = useState<SheetMode>(null);
  const [amount, setAmount] = useState("");
  const [currency, setCurrency] = useState<"USDT" | "ETH">("USDT");

  const percentOfMin = deposit
    ? Number((deposit.balance * 100n) / deposit.minimumRequired)
    : 0;

  const handleConfirm = async () => {
    if (!address || !amount) return;
    const wei = BigInt(Math.floor(parseFloat(amount) * 1e18));
    if (sheetMode === "deposit") {
      await depositMutation.mutateAsync({ address, amount: wei, currency });
    } else {
      await withdrawMutation.mutateAsync({ address, amount: wei });
    }
    setSheetMode(null);
    setAmount("");
  };

  const isPending = depositMutation.isPending || withdrawMutation.isPending;

  return (
    <div className="px-4 space-y-4">
      {/* Balance Card */}
      <div className="p-6 rounded-2xl bg-gradient-to-br from-surface-gradient-from to-surface-gradient-to text-center">
        <div className="text-[10px] text-text-muted uppercase tracking-wider">Current Balance</div>
        <div className="mt-2">
          {deposit && <AmountDisplay amount={deposit.balance} currency={deposit.currency} size="lg" />}
        </div>
        <div className="mt-4">
          <div className="flex justify-between text-[10px] text-text-muted mb-1">
            <span>Min Required: {deposit ? formatAmount(deposit.minimumRequired) : "\u2014"} {deposit?.currency}</span>
            <span>{percentOfMin}%</span>
          </div>
          <div className="h-1 bg-surface-secondary rounded-full overflow-hidden">
            <div
              className="h-full bg-gradient-to-r from-status-success to-brand-blue rounded-full transition-all"
              style={{ width: `${Math.min(percentOfMin, 100)}%` }}
            />
          </div>
        </div>
      </div>

      {/* Action Buttons */}
      <div className="flex gap-2.5">
        <Button onClick={() => setSheetMode("deposit")} className="flex-1 bg-status-success hover:bg-status-success/90 py-3">
          Deposit
        </Button>
        <Button onClick={() => setSheetMode("withdraw")} variant="outline" className="flex-1 py-3">
          Withdraw
        </Button>
      </div>

      {/* History */}
      <div>
        <h3 className="text-[13px] font-semibold mb-3">History</h3>
        <div className="space-y-2">
          {history?.map((record) => (
            <div key={record.id} className="flex justify-between items-center p-3 bg-surface-card rounded-xl">
              <div>
                <div className="text-xs font-semibold capitalize">{record.type}</div>
                <div className="text-[10px] text-text-muted mt-0.5">{formatDate(record.timestamp)}</div>
              </div>
              <div className={`text-sm font-bold ${record.type === "deposit" ? "text-status-success" : "text-status-danger"}`}>
                {record.type === "deposit" ? "+" : "-"}{formatAmount(record.amount)}
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* Deposit/Withdraw Sheet */}
      <Drawer.Root open={sheetMode !== null} onOpenChange={(o) => !o && setSheetMode(null)}>
        <Drawer.Portal>
          <Drawer.Overlay className="fixed inset-0 bg-black/60 z-50" />
          <Drawer.Content className="fixed bottom-0 left-0 right-0 max-w-mobile mx-auto bg-surface-card rounded-t-2xl z-50 p-6">
            <div className="w-12 h-1 bg-surface-secondary rounded-full mx-auto mb-6" />
            <h2 className="text-lg font-bold mb-4">{sheetMode === "deposit" ? "Deposit Funds" : "Withdraw Funds"}</h2>

            {sheetMode === "deposit" && (
              <div className="flex gap-2 mb-4">
                {(["USDT", "ETH"] as const).map((c) => (
                  <button
                    key={c}
                    onClick={() => setCurrency(c)}
                    className={`flex-1 py-2 rounded-lg text-sm font-semibold border ${
                      currency === c ? "border-brand-blue bg-brand-blue/10 text-brand-blue" : "border-border text-text-secondary"
                    }`}
                  >
                    {c}
                  </button>
                ))}
              </div>
            )}

            <input
              type="number"
              value={amount}
              onChange={(e) => setAmount(e.target.value)}
              placeholder="Enter amount"
              className="w-full px-4 py-3 bg-surface-secondary rounded-xl text-text-primary text-sm outline-none border border-border focus:border-brand-blue mb-4"
            />

            <Button onClick={handleConfirm} disabled={!amount || isPending} className="w-full py-3">
              {isPending ? "Processing..." : "Confirm"}
            </Button>
          </Drawer.Content>
        </Drawer.Portal>
      </Drawer.Root>
    </div>
  );
}
