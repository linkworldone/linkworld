import { useState, useEffect } from "react";
import { useAccount } from "wagmi";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { BottomSheet } from "@/components/shared/BottomSheet";
import { useDeposit, useDepositHistory, useDepositMutation, useWithdrawMutation } from "@/hooks/useDeposit";
import { AmountDisplay } from "@/components/shared/AmountDisplay";
import { formatAmount, formatDate } from "@/utils/format";

type SheetMode = "deposit" | "withdraw" | null;

export default function Deposit() {
  const { address } = useAccount();
  const { data: deposit, refetch: refetchBalance } = useDeposit(address);
  const { data: history } = useDepositHistory(address);
  const depositMutation = useDepositMutation();
  const withdrawMutation = useWithdrawMutation();

  const [sheetMode, setSheetMode] = useState<SheetMode>(null);
  const [amount, setAmount] = useState("");
  const [currency, setCurrency] = useState<"USDT" | "ETH">("USDT");

  // 合约成功后仅上报 pending 意向（不据成功置终态，design §3.3）+ 刷新读链余额。
  // 文案为「处理中」而非「confirmed」：真实到账以链上 getDepositAmount + 后端事件回填为准。
  useEffect(() => {
    if (depositMutation.isSuccess && amount) {
      depositMutation.recordIntent(amount);
      refetchBalance();
      toast("充值处理中 · 约 1-2 分钟到账");
      setSheetMode(null);
      setAmount("");
    }
  }, [depositMutation.isSuccess]);

  useEffect(() => {
    if (withdrawMutation.isSuccess) {
      withdrawMutation.recordIntent();
      refetchBalance();
      toast("提现处理中 · 约 1-2 分钟到账");
      setSheetMode(null);
      setAmount("");
    }
  }, [withdrawMutation.isSuccess]);

  const handleConfirm = () => {
    if (!address || !amount) return;
    if (sheetMode === "deposit") {
      depositMutation.deposit(amount);
    } else {
      withdrawMutation.withdraw();
    }
  };

  const isWaitingSignature = depositMutation.isContractPending || withdrawMutation.isContractPending;
  const isConfirming = depositMutation.isConfirming || withdrawMutation.isConfirming;
  const isPending = isWaitingSignature || isConfirming;
  const buttonText = isWaitingSignature
    ? "Waiting for signature..."
    : isConfirming
      ? "Confirming on chain..."
      : "Confirm";

  return (
    <div className="px-4 space-y-4">
      {/* Balance Card */}
      <div className="p-6 rounded-2xl bg-gradient-to-br from-surface-gradient-from to-surface-gradient-to text-center">
        <div className="text-[10px] text-text-muted uppercase tracking-wider">Current Balance</div>
        <div className="mt-2">
          {deposit && <AmountDisplay amount={deposit.balance} currency={deposit.currency} size="lg" />}
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
      <BottomSheet open={sheetMode !== null} onOpenChange={(o) => !o && setSheetMode(null)}>
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
          {buttonText}
        </Button>
      </BottomSheet>
    </div>
  );
}
