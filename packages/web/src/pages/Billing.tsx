import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { useAccount } from "wagmi";
import { Button } from "@/components/ui/button";
import { BottomSheet } from "@/components/shared/BottomSheet";
import { useBills, usePayBill } from "@/hooks/useBilling";
import { formatDate, formatUSD, parseUnits } from "@/utils/format";
import { Badge } from "@/components/ui/badge";

type Filter = "unpaid" | "paid";

export default function Billing() {
  const navigate = useNavigate();
  const { address } = useAccount();
  const [filter, setFilter] = useState<Filter>("unpaid");
  const { data: bills } = useBills(address, filter);
  const { payBill, isContractPending, isSuccess, recordToBackend } = usePayBill();
  const [payingBillId, setPayingBillId] = useState<string | null>(null);

  const unpaidBills = bills?.filter((b) => b.status !== "paid") || [];
  const payingBill = bills?.find((b) => b.id === payingBillId);

  // 合约支付成功后同步后端
  useEffect(() => {
    if (isSuccess && payingBillId) {
      recordToBackend(payingBillId);
      setPayingBillId(null);
    }
  }, [isSuccess]);

  const handlePay = () => {
    if (!payingBillId || !payingBill) return;
    const totalAmountWei = parseUnits(payingBill.totalAmount.toString());
    payBill(BigInt(payingBillId), totalAmountWei);
  };

  return (
    <div className="px-4 space-y-3">
      <div className="flex bg-surface-card rounded-xl p-0.5">
        {(["unpaid", "paid"] as const).map((f) => (
          <button
            key={f}
            onClick={() => setFilter(f)}
            className={`flex-1 py-2 rounded-lg text-[13px] font-semibold capitalize transition-colors ${
              filter === f ? "bg-brand-blue text-white" : "text-text-muted"
            }`}
          >
            {f === "unpaid" ? "Unpaid" : "History"}
          </button>
        ))}
      </div>

      {filter === "unpaid" && unpaidBills.length > 0 && (
        <div className="p-3 bg-status-danger/15 border border-status-danger/30 rounded-xl flex items-center gap-2">
          <span className="text-base">⚠️</span>
          <div className="text-xs text-status-danger">
            {unpaidBills.length} unpaid bill{unpaidBills.length > 1 ? "s" : ""} · Due by{" "}
            {formatDate(unpaidBills[0].dueDate)}
          </div>
        </div>
      )}

      <div className="space-y-2.5">
        {bills?.map((bill) => (
          <div key={bill.id} className={`p-4 bg-surface-card rounded-xl ${bill.status === "paid" ? "opacity-70" : ""}`}>
            <div className="flex justify-between items-center mb-3">
              <button onClick={() => navigate(`/billing/${bill.id}`)} className="text-sm font-bold hover:text-brand-blue">
                {new Date(bill.month + "-01").toLocaleDateString("en-US", { month: "long", year: "numeric" })}
              </button>
              <Badge variant={bill.status === "paid" ? "secondary" : "destructive"} className="text-[10px]">
                {bill.status === "paid" ? "Paid" : "Unpaid"}
              </Badge>
            </div>
            <div className="space-y-1.5 mb-3">
              <div className="flex justify-between text-xs">
                <span className="text-text-secondary">Operator Fee</span>
                <span>{formatUSD(bill.operatorFee)}</span>
              </div>
              <div className="flex justify-between text-xs">
                <span className="text-text-secondary">Platform Fee (2.5%)</span>
                <span>{formatUSD(bill.platformFee)}</span>
              </div>
              <div className="h-px bg-surface-secondary my-1" />
              <div className="flex justify-between text-sm font-bold">
                <span>Total</span>
                <span>{formatUSD(bill.totalAmount)}</span>
              </div>
            </div>
            {bill.status !== "paid" && (
              <Button onClick={() => setPayingBillId(bill.id)} className="w-full">
                Pay Now
              </Button>
            )}
            {bill.status === "paid" && bill.paidAt && (
              <div className="text-[10px] text-text-muted">Paid on {formatDate(bill.paidAt)}</div>
            )}
          </div>
        ))}
      </div>

      <BottomSheet open={payingBillId !== null} onOpenChange={(o) => !o && setPayingBillId(null)}>
        <h2 className="text-lg font-bold mb-4">Confirm Payment</h2>
        {payingBill && (
          <>
            <div className="p-3 bg-surface-secondary rounded-xl space-y-2 mb-4">
              <div className="flex justify-between text-xs">
                <span className="text-text-muted">Amount</span>
                <span className="font-bold">{formatUSD(payingBill.totalAmount)}</span>
              </div>
              <div className="flex justify-between text-xs">
                <span className="text-text-muted">Source</span>
                <span>Deposit Balance</span>
              </div>
            </div>
            <Button onClick={handlePay} disabled={isContractPending} className="w-full py-3">
              {isContractPending ? "Processing..." : "Confirm Payment"}
            </Button>
          </>
        )}
      </BottomSheet>
    </div>
  );
}
