import { useState, useEffect } from "react";
import { useParams } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { BottomSheet } from "@/components/shared/BottomSheet";
import { Badge } from "@/components/ui/badge";
import { useBillDetail, usePayBill } from "@/hooks/useBilling";
import { formatUSD, formatDate, parseUnits } from "@/utils/format";

export default function BillDetail() {
  const { billId } = useParams<{ billId: string }>();
  const { data: bill } = useBillDetail(billId);
  const { payBill, isContractPending, isSuccess, recordToBackend } = usePayBill();
  const [showPay, setShowPay] = useState(false);

  // 合约支付成功后同步后端
  useEffect(() => {
    if (isSuccess && bill) {
      recordToBackend(bill.id);
      setShowPay(false);
    }
  }, [isSuccess]);

  if (!bill) return <div className="p-4 text-text-secondary text-sm">Loading...</div>;

  const handlePay = () => {
    const totalAmountWei = parseUnits(bill.totalAmount.toString());
    payBill(BigInt(bill.id), totalAmountWei);
  };

  return (
    <div className="px-4 space-y-4">
      <div className="flex justify-between items-center">
        <h2 className="text-lg font-bold">
          {new Date(bill.month + "-01").toLocaleDateString("en-US", { month: "long", year: "numeric" })}
        </h2>
        <Badge variant={bill.status === "paid" ? "secondary" : "destructive"} className="text-[10px]">
          {bill.status}
        </Badge>
      </div>

      <div className="p-4 bg-surface-card rounded-xl space-y-2">
        <div className="flex justify-between text-sm">
          <span className="text-text-secondary">Operator Fee</span>
          <span>{formatUSD(bill.operatorFee)}</span>
        </div>
        <div className="flex justify-between text-sm">
          <span className="text-text-secondary">Platform Fee (2.5%)</span>
          <span>{formatUSD(bill.platformFee)}</span>
        </div>
        <div className="h-px bg-surface-secondary" />
        <div className="flex justify-between text-base font-bold">
          <span>Total</span>
          <span>{formatUSD(bill.totalAmount)}</span>
        </div>
      </div>

      <div className="p-4 bg-surface-card rounded-xl">
        <h3 className="text-[13px] font-semibold mb-3">Usage Details</h3>
        <div className="flex justify-around">
          <div className="text-center">
            <div className="text-[22px] font-extrabold text-brand-blue">{bill.usage.dataGB}</div>
            <div className="text-[10px] text-text-muted">GB Data</div>
          </div>
          <div className="w-px bg-surface-secondary" />
          <div className="text-center">
            <div className="text-[22px] font-extrabold text-brand-purple">{bill.usage.callMinutes}</div>
            <div className="text-[10px] text-text-muted">Min Calls</div>
          </div>
        </div>
      </div>

      {bill.status !== "paid" && (
        <Button onClick={() => setShowPay(true)} className="w-full py-3">
          Pay {formatUSD(bill.totalAmount)}
        </Button>
      )}
      {bill.paidAt && <div className="text-xs text-text-muted text-center">Paid on {formatDate(bill.paidAt)}</div>}

      <BottomSheet open={showPay} onOpenChange={setShowPay}>
        <h2 className="text-lg font-bold mb-4">Confirm Payment</h2>
        <div className="p-3 bg-surface-secondary rounded-xl space-y-2 mb-4">
          <div className="flex justify-between text-xs">
            <span className="text-text-muted">Amount</span>
            <span className="font-bold">{formatUSD(bill.totalAmount)}</span>
          </div>
        </div>
        <Button onClick={handlePay} disabled={isContractPending} className="w-full py-3">
          {isContractPending ? "Processing..." : "Confirm Payment"}
        </Button>
      </BottomSheet>
    </div>
  );
}
