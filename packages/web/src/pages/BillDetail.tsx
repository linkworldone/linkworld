import { useMemo, useState } from "react";
import { useParams } from "react-router-dom";
import { useAccount } from "wagmi";
import { AlertCircle, Loader2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { BottomSheet } from "@/components/shared/BottomSheet";
import { TwoStepAction, type ActionTx } from "@/components/shared/TwoStepAction";
import { TxStatusBadge } from "@/components/shared/TxStatusBadge";
import { AmountDisplay } from "@/components/shared/AmountDisplay";
import { FeeBreakdown } from "@/components/shared/FeeBreakdown";
import { useBillDetail, usePayBill } from "@/hooks/useBilling";
import { useCalculateFee } from "@/hooks/contracts";
import { getContractAddress } from "@/config/contracts";
import { WalletAuthRejectedError } from "@/services/api/signedPost";
import { formatAmount, formatDate } from "@/utils/format";

export default function BillDetail() {
  const { billId } = useParams<{ billId: string }>();
  const { address, chainId } = useAccount();
  const { data: bill } = useBillDetail(billId);
  const { payBill, txState, recordIntent } = usePayBill();
  const [showPay, setShowPay] = useState(false);
  const [pending, setPending] = useState(false);
  const [rejectMsg, setRejectMsg] = useState<string | null>(null);

  // totalAmount 已是 USDT 6 位最小单位字符串（T2）；直接转 bigint，**不再 parseUnits**（双重缩放资损红线）。
  const totalMinUnit = bill ? BigInt(bill.totalAmount) : 0n;

  // 付账授权额 = 本金 + calculateFee（直读合约，不自算，design §3.6）。读链未回则按本金（approve 不足时第二步 revert，安全）。
  const feeQ = useCalculateFee(totalMinUnit);
  const approveAmount = totalMinUnit + (feeQ.fee ?? 0n);

  // 第二步动作（注入 TwoStepAction）：write=payBill(billId)；state 走 hook 统一 TxState。
  const payAction: ActionTx = useMemo(
    () => ({
      write: () => {
        setRejectMsg(null);
        if (bill) payBill(BigInt(bill.id));
      },
      state: txState,
    }),
    [bill, payBill, txState],
  );

  // 合约成功 → 上报 pending 意向（signedPost）；拒签提示、不进 pending（不据成功置已付，design §3.3）。
  const handlePaySuccess = async () => {
    if (!bill) return;
    try {
      await recordIntent(bill.id);
      setPending(true);
      setShowPay(false);
    } catch (err) {
      if (err instanceof WalletAuthRejectedError) {
        setRejectMsg("身份签名被取消，操作未提交");
      } else {
        setRejectMsg("上报失败，稍后将自动重试");
      }
    }
  };

  if (!bill) return <div className="p-4 text-text-secondary text-sm">加载中…</div>;

  const isPaid = bill.status === "paid";
  const isPaying = bill.status === "paying";
  const paySpender = address ? safeSpender(chainId) : undefined;

  return (
    <div className="px-4 space-y-4">
      <div className="flex justify-between items-center">
        <h2 className="text-lg font-bold text-text-primary">
          {new Date(bill.month + "-01").toLocaleDateString("zh-CN", {
            year: "numeric",
            month: "long",
          })}
        </h2>
        <BillStatusBadge status={bill.status} />
      </div>

      <Card>
        <CardContent className="py-4 space-y-2.5">
          <div className="flex justify-between text-sm">
            <span className="text-text-on-light-secondary">运营商费用</span>
            <span className="font-data tabular-nums text-text-on-light-primary">{formatAmount(BigInt(bill.operatorFee))} USDT</span>
          </div>
          {/* 平台手续费：读链费率 + calculateFee（不写死 2.5%）。 */}
          <FeeBreakdown amount={totalMinUnit} size="sm" />
          {bill.trafficCardDeduction && BigInt(bill.trafficCardDeduction) > 0n && (
            <div className="flex justify-between text-sm">
              <span className="text-text-on-light-secondary">流量卡抵扣</span>
              <span className="text-status-success font-data tabular-nums">
                -{formatAmount(BigInt(bill.trafficCardDeduction))} USDT
              </span>
            </div>
          )}
          <div className="h-px bg-surface-card-line" />
          <div className="flex justify-between items-center">
            <span className="text-base font-bold text-text-on-light-primary">合计</span>
            <AmountDisplay amount={totalMinUnit} currency="USDT" size="md" />
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardContent className="py-4">
          <h3 className="text-[13px] font-semibold mb-3 text-text-on-light-primary">用量详情</h3>
          <div className="flex justify-around">
            <div className="text-center">
              <div className="text-[22px] font-extrabold text-text-on-light-primary font-data tabular-nums">
                {bill.usage.dataGB}
              </div>
              <div className="text-[10px] text-text-on-light-muted">GB 流量</div>
            </div>
            <div className="w-px bg-surface-card-line" />
            <div className="text-center">
              <div className="text-[22px] font-extrabold text-text-on-light-primary font-data tabular-nums">
                {bill.usage.callMinutes}
              </div>
              <div className="text-[10px] text-text-on-light-muted">分钟通话</div>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* pending 提示（处理中 · 约 1-2 分钟）：不据成功置已付，等后端 BillPaid 回填。 */}
      {(pending || isPaying) && (
        <Card data-slot="pending-notice">
          <CardContent className="py-3 flex items-start gap-2">
            <Loader2 className="size-4 mt-0.5 shrink-0 animate-spin text-brand-royal" />
            <div className="text-xs text-text-on-light-secondary space-y-1">
              <p className="font-medium text-text-on-light-primary">付款处理中 · 约 1-2 分钟</p>
              <p>可安全离开本页，到账后账单状态会自动更新。</p>
            </div>
          </CardContent>
        </Card>
      )}

      {!isPaid && !isPaying && (
        <Button onClick={() => setShowPay(true)} className="w-full py-3">
          支付 {formatAmount(totalMinUnit)} USDT
        </Button>
      )}
      {isPaid && bill.paidAt && (
        <div className="text-xs text-text-muted text-center">
          已于 {formatDate(bill.paidAt)} 支付
        </div>
      )}

      <BottomSheet open={showPay} onOpenChange={setShowPay}>
        <h2 className="text-lg font-bold mb-4 text-text-on-light-primary">确认支付</h2>

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
          <FeeBreakdown amount={totalMinUnit} />
          <div className="h-px bg-surface-card-line" />
          <div className="flex justify-between text-xs">
            <span className="text-text-on-light-secondary">需授权总额</span>
            <span className="font-bold font-data tabular-nums text-text-on-light-primary" data-slot="approve-total">
              {feeQ.fee !== undefined ? `${formatAmount(approveAmount)} USDT` : "--"}
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
      </BottomSheet>
    </div>
  );
}

// 账单状态徽章：paying 用 TxStatusBadge(info 蓝 + Loader2，**禁绿**)；paid=secondary、其余=destructive。
function BillStatusBadge({ status }: { status: "unpaid" | "paying" | "paid" | "overdue" }) {
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
