import { useEffect, useMemo, useState } from "react";
import { useNavigate } from "react-router-dom";
import { useAccount } from "wagmi";
import { AlertCircle, AlertTriangle } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { BottomSheet } from "@/components/shared/BottomSheet";
import { TwoStepAction, type ActionTx } from "@/components/shared/TwoStepAction";
import { TxStatusBadge } from "@/components/shared/TxStatusBadge";
import { AmountDisplay } from "@/components/shared/AmountDisplay";
import { FeeBreakdown } from "@/components/shared/FeeBreakdown";
import { useBills, usePayBill } from "@/hooks/useBilling";
import { getContractAddress } from "@/config/contracts";
import { WalletAuthRejectedError } from "@/services/api/signedPost";
import { formatAmount, formatDate } from "@/utils/format";
import type { Bill } from "@/types";

type Filter = "unpaid" | "paid";

export default function Billing() {
  const navigate = useNavigate();
  const { address, chainId } = useAccount();
  const [filter, setFilter] = useState<Filter>("unpaid");
  const { data: bills } = useBills(address);
  const { payBill, txState, recordIntent } = usePayBill();
  const [payingBillId, setPayingBillId] = useState<string | null>(null);
  const [rejectMsg, setRejectMsg] = useState<string | null>(null);

  // 列表过滤：paid → status==="paid"；unpaid → 非 paid（含 unpaid/paying/overdue）。
  const shown = (bills ?? []).filter((b) =>
    filter === "paid" ? b.status === "paid" : b.status !== "paid",
  );
  const unpaidBills = (bills ?? []).filter((b) => b.status === "unpaid" || b.status === "overdue");
  const payingBill = bills?.find((b) => b.id === payingBillId);

  // totalAmount 已是 USDT 6 位最小单位字符串（T2）= amount + platformFee（链上实收 = 两段转账之和）；
  // 直接转 bigint，**不再 parseUnits**（双重缩放资损红线）。
  const totalMinUnit = payingBill ? BigInt(payingBill.totalAmount) : 0n;
  // 授权额 = totalAmount（exact = 链上实收）；**不再在已含费的合计上叠 calculateFee**（过额授权资损红线）。
  const approveAmount = totalMinUnit;

  const payAction: ActionTx = useMemo(
    () => ({
      write: () => {
        setRejectMsg(null);
        if (payingBill) payBill(BigInt(payingBill.id));
      },
      state: txState,
    }),
    [payingBill, payBill, txState],
  );

  // 合约成功 → 上报 pending 意向；拒签提示、不进 pending（不据成功置已付，design §3.3）。
  const handlePaySuccess = async () => {
    if (!payingBill) return;
    try {
      await recordIntent(payingBill.id);
      setPayingBillId(null);
    } catch (err) {
      if (err instanceof WalletAuthRejectedError) {
        setRejectMsg("身份签名被取消，操作未提交");
      } else {
        setRejectMsg("上报失败，稍后将自动重试");
      }
    }
  };

  // 弹层关闭时清提示。
  useEffect(() => {
    if (payingBillId === null) setRejectMsg(null);
  }, [payingBillId]);

  const paySpender = address ? safeSpender(chainId) : undefined;

  return (
    <div className="px-4 space-y-3">
      <div className="flex bg-surface-card border border-surface-card-line rounded-xl p-0.5">
        {(["unpaid", "paid"] as const).map((f) => (
          <button
            key={f}
            onClick={() => setFilter(f)}
            className={`flex-1 py-2 rounded-lg text-[13px] font-semibold transition-colors min-h-[44px] ${
              filter === f ? "bg-brand-gold text-brand-navy" : "text-text-on-light-secondary"
            }`}
          >
            {f === "unpaid" ? "待支付" : "历史"}
          </button>
        ))}
      </div>

      {filter === "unpaid" && unpaidBills.length > 0 && (
        <div className="p-3 bg-status-danger/15 border border-status-danger/30 rounded-xl flex items-center gap-2">
          <AlertTriangle className="size-4 shrink-0 text-status-danger" />
          <div className="text-xs text-status-danger">
            {unpaidBills.length} 笔待支付账单 · 截止 {formatDate(unpaidBills[0].dueDate)}
          </div>
        </div>
      )}

      <div className="space-y-2.5">
        {shown.map((bill) => {
          const isPaid = bill.status === "paid";
          const billTotal = BigInt(bill.totalAmount);
          const billPlatformFee = BigInt(bill.platformFee);
          return (
            <Card key={bill.id} className={isPaid ? "opacity-70" : ""}>
              <CardContent className="py-4">
                <div className="flex justify-between items-center mb-3">
                  <button
                    onClick={() => navigate(`/billing/${bill.id}`)}
                    className="text-sm font-bold text-text-on-light-primary hover:text-brand-royal"
                  >
                    {new Date(bill.month + "-01").toLocaleDateString("zh-CN", {
                      year: "numeric",
                      month: "long",
                    })}
                  </button>
                  <BillStatusBadge status={bill.status} />
                </div>
                <div className="space-y-1.5 mb-3">
                  <div className="flex justify-between text-xs">
                    <span className="text-text-on-light-secondary">运营商费用</span>
                    <span className="font-data tabular-nums text-text-on-light-primary">
                      {formatAmount(BigInt(bill.operatorFee))} USDT
                    </span>
                  </div>
                  <FeeBreakdown fee={billPlatformFee} />
                  {bill.trafficCardDeduction && BigInt(bill.trafficCardDeduction) > 0n && (
                    <div className="flex justify-between text-xs">
                      <span className="text-text-on-light-secondary">流量卡抵扣</span>
                      <span className="text-status-success font-data tabular-nums">
                        -{formatAmount(BigInt(bill.trafficCardDeduction))} USDT
                      </span>
                    </div>
                  )}
                  <div className="h-px bg-surface-card-line my-1" />
                  <div className="flex justify-between items-center">
                    <span className="text-sm font-bold text-text-on-light-primary">合计</span>
                    <AmountDisplay amount={billTotal} currency="USDT" size="sm" />
                  </div>
                </div>
                {bill.status === "unpaid" || bill.status === "overdue" ? (
                  <Button onClick={() => setPayingBillId(bill.id)} className="w-full">
                    立即支付
                  </Button>
                ) : null}
                {isPaid && bill.paidAt && (
                  <div className="text-[10px] text-text-muted">已于 {formatDate(bill.paidAt)} 支付</div>
                )}
              </CardContent>
            </Card>
          );
        })}
      </div>

      <BottomSheet
        open={payingBillId !== null}
        onOpenChange={(o) => !o && setPayingBillId(null)}
      >
        <h2 className="text-lg font-bold mb-4 text-text-on-light-primary">确认支付</h2>
        {payingBill && (
          <>
            {rejectMsg && (
              <div
                className="mb-4 flex items-center gap-2 rounded-md bg-status-danger/10 px-3 py-2 text-xs text-status-danger"
                data-slot="reject-notice"
              >
                <AlertCircle className="size-4 shrink-0" />
                {rejectMsg}
              </div>
            )}
            <div className="p-3 bg-surface-input rounded-xl space-y-2 mb-4">
              <div className="flex justify-between text-xs">
                <span className="text-text-on-light-secondary">账单金额</span>
                <span className="font-bold font-data tabular-nums text-text-on-light-primary">
                  {formatAmount(totalMinUnit)} USDT
                </span>
              </div>
              <FeeBreakdown fee={BigInt(payingBill.platformFee)} />
              <div className="h-px bg-surface-card-line" />
              <div className="flex justify-between text-xs">
                <span className="text-text-on-light-secondary">需授权总额</span>
                <span className="font-bold font-data tabular-nums text-text-on-light-primary" data-slot="approve-total">
                  {`${formatAmount(approveAmount)} USDT`}
                </span>
              </div>
            </div>
            <TwoStepAction
              owner={address}
              spender={paySpender}
              amount={approveAmount}
              action={payAction}
              actionLabel="支付"
              onSuccess={handlePaySuccess}
            />
          </>
        )}
      </BottomSheet>
    </div>
  );
}

// 账单状态徽章：paying 用 TxStatusBadge(info 蓝 + Loader2，**禁绿**)；paid=secondary、其余=destructive。
function BillStatusBadge({ status }: { status: Bill["status"] }) {
  if (status === "paying") return <TxStatusBadge status="pending" />;
  const label =
    status === "paid" ? "已支付" : status === "overdue" ? "已逾期" : "待支付";
  return (
    <Badge variant={status === "paid" ? "secondary" : "destructive"} className="text-[10px]">
      {label}
    </Badge>
  );
}

// 未部署链 getContractAddress 抛错；弹层未打开/钱包未连不应崩页，故包一层。
function safeSpender(chainId: number | undefined): `0x${string}` | undefined {
  if (chainId === undefined) return undefined;
  try {
    return getContractAddress(chainId, "Payment");
  } catch {
    return undefined;
  }
}
