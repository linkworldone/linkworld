import { useEffect, useMemo, useState } from "react";
import { useAccount } from "wagmi";
import {
  ArrowDownToLine,
  ArrowUpFromLine,
  ExternalLink,
  Loader2,
  AlertCircle,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { BottomSheet } from "@/components/shared/BottomSheet";
import { TwoStepAction, type ActionTx } from "@/components/shared/TwoStepAction";
import { LockCountdown, type LockStateResult } from "@/components/shared/LockCountdown";
import { AmountDisplay } from "@/components/shared/AmountDisplay";
import { TxStatusBadge } from "@/components/shared/TxStatusBadge";
import {
  useDeposit,
  useDepositHistory,
  useDepositMutation,
  useWithdrawMutation,
} from "@/hooks/useDeposit";
import { useUsdtBalance } from "@/hooks/contracts";
import { getContractAddress } from "@/config/contracts";
import { arbitrumSepolia, hardhatLocal } from "@/config/chains";
import { WalletAuthRejectedError } from "@/services/api/signedPost";
import { MIN_DEPOSIT_USDT } from "@/config/constants";
import { parseUnits, formatAmount, formatDate } from "@/utils/format";
import type { DepositRecord } from "@/types";

type SheetMode = "deposit" | "withdraw" | null;

// pending 提示态（design §3.1）：成功后展示「处理中·约 1-2 分钟」+ Arbiscan 逃生链接 + 可安全离开。
// 不据 200 / tx 成功置终态；终态由后端 event_sync 回填，余额读链。
interface PendingNotice {
  kind: "deposit" | "withdraw";
}

const EXPLORER_URL =
  (import.meta.env.VITE_CHAIN_ID === "31337" ? hardhatLocal : arbitrumSepolia)
    .blockExplorers?.default.url;

export default function Deposit() {
  const { address, chainId } = useAccount();
  const { data: deposit, refetch: refetchBalance } = useDeposit(address);
  const { data: history } = useDepositHistory(address);
  const { data: usdtBalanceRaw } = useUsdtBalance(address);

  const depositMutation = useDepositMutation();
  const withdrawMutation = useWithdrawMutation();

  const [sheetMode, setSheetMode] = useState<SheetMode>(null);
  const [amount, setAmount] = useState("");
  const [pending, setPending] = useState<PendingNotice | null>(null);
  const [rejectMsg, setRejectMsg] = useState<string | null>(null);
  const [lockUnlocked, setLockUnlocked] = useState(false);

  const usdtBalance = (usdtBalanceRaw as bigint | undefined) ?? undefined;

  // ── Deposit 第二步动作（注入 TwoStepAction）：write 用当前 amount 闭包；state 走 hook 统一 TxState ──
  const depositAction: ActionTx = useMemo(
    () => ({
      write: () => {
        setRejectMsg(null);
        depositMutation.deposit(amount);
      },
      state: depositMutation.txState,
    }),
    [amount, depositMutation.txState],
  );

  // 充值合约成功 → 上报 pending 意向（signedPost）；拒签提示，不进 pending。
  const handleDepositSuccess = async () => {
    try {
      await depositMutation.recordIntent(amount);
      refetchBalance();
      setPending({ kind: "deposit" });
      setSheetMode(null);
      setAmount("");
    } catch (err) {
      if (err instanceof WalletAuthRejectedError) {
        setRejectMsg("身份签名被取消，操作未提交");
      } else {
        setRejectMsg("上报失败，稍后将自动重试");
      }
    }
  };

  // 提现合约成功 → 上报 pending 意向（无 tx_hash 记账，design §3.3）。
  useEffect(() => {
    if (!withdrawMutation.isSuccess) return;
    (async () => {
      try {
        await withdrawMutation.recordIntent();
        refetchBalance();
        setPending({ kind: "withdraw" });
        setSheetMode(null);
        setAmount("");
      } catch (err) {
        if (err instanceof WalletAuthRejectedError) {
          setRejectMsg("身份签名被取消，操作未提交");
        } else {
          setRejectMsg("上报失败，稍后将自动重试");
        }
      }
    })();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [withdrawMutation.isSuccess]);

  // 充值金额校验（design §3.2）：≥ 10 USDT（链上 require）& ≤ 钱包 USDT 余额。
  const amountMinUnit = amount ? parseUnits(amount, 6) : 0n;
  const tooSmall = amountMinUnit > 0n && amountMinUnit < MIN_DEPOSIT_USDT;
  const overBalance =
    usdtBalance !== undefined && amountMinUnit > 0n && amountMinUnit > usdtBalance;
  const depositAmountValid = amountMinUnit > 0n && !tooSmall && !overBalance;

  const withdrawTxState = withdrawMutation.txState;
  const withdrawBusy =
    withdrawTxState.status === "pending-signature" ||
    withdrawTxState.status === "pending-confirmation";

  const openSheet = (mode: SheetMode) => {
    setRejectMsg(null);
    setAmount("");
    setSheetMode(mode);
  };

  const depositSpender = address ? safeSpender(chainId) : undefined;

  return (
    <div className="px-4 space-y-4">
      {/* 余额卡（暖米白卡内金额 navy，B2 铁律） */}
      <Card>
        <CardContent className="text-center py-6">
          <div className="text-[11px] uppercase tracking-wider text-text-on-light-muted">
            保证金余额
          </div>
          <div className="mt-2">
            {deposit ? (
              <AmountDisplay amount={deposit.balance} currency="USDT" size="lg" />
            ) : (
              <span className="text-text-on-light-muted">—</span>
            )}
          </div>
          {usdtBalance !== undefined && (
            <div className="mt-2 text-xs text-text-on-light-secondary">
              钱包可用 {formatAmount(usdtBalance)} USDT
            </div>
          )}
        </CardContent>
      </Card>

      {/* 锁仓倒计时（联动 Withdraw 禁用） */}
      <Card>
        <CardContent className="py-4 space-y-2">
          <LockCountdown
            address={address}
            onStateChange={(s: LockStateResult) =>
              setLockUnlocked(s.status === "unlocked")
            }
          />
          <p className="text-[11px] text-text-on-light-muted">
            再次充值将把锁仓期顺延 30 天
          </p>
        </CardContent>
      </Card>

      {/* pending 提示（处理中 · 约 1-2 分钟 + 逃生链接） */}
      {pending && (
        <Card data-slot="pending-notice">
          <CardContent className="py-3 flex items-start gap-2">
            <Loader2 className="size-4 mt-0.5 shrink-0 animate-spin text-brand-royal" />
            <div className="text-xs text-text-on-light-secondary space-y-1">
              <p className="font-medium text-text-on-light-primary">
                {pending.kind === "deposit" ? "充值处理中" : "提现处理中"} · 约 1-2 分钟
              </p>
              <p>可安全离开本页，到账后余额会自动更新。</p>
              {EXPLORER_URL && (
                <a
                  href={EXPLORER_URL}
                  target="_blank"
                  rel="noreferrer"
                  className="inline-flex items-center gap-1 text-brand-royal hover:underline"
                >
                  在区块浏览器查看 <ExternalLink className="size-3" />
                </a>
              )}
            </div>
          </CardContent>
        </Card>
      )}

      {/* 操作按钮 */}
      <div className="flex gap-2.5">
        <Button onClick={() => openSheet("deposit")} className="flex-1 py-3">
          <ArrowDownToLine className="size-4" />
          充值
        </Button>
        <Button
          onClick={() => openSheet("withdraw")}
          variant="outline"
          className="flex-1 py-3"
          disabled={!lockUnlocked || (deposit ? deposit.balance === 0n : true)}
        >
          <ArrowUpFromLine className="size-4" />
          提现
        </Button>
      </div>

      {/* 历史 */}
      <div>
        <h3 className="text-[13px] font-semibold mb-3 text-text-primary">交易记录</h3>
        <div className="space-y-2">
          {history?.map((record: DepositRecord) => (
            <Card key={record.id}>
              <CardContent className="py-3 flex justify-between items-center">
                <div>
                  <div className="text-xs font-semibold text-text-on-light-primary">
                    {record.type === "deposit"
                      ? "充值"
                      : record.type === "withdraw"
                        ? "提现"
                        : "扣费"}
                  </div>
                  <div className="text-[10px] text-text-on-light-muted mt-0.5">
                    {formatDate(record.timestamp)}
                  </div>
                </div>
                <div className="flex items-center gap-2">
                  <TxStatusBadge
                    status={record.status === "confirmed" ? "confirmed" : "pending"}
                  />
                  <span className="text-sm font-bold font-data tabular-nums text-text-on-light-primary">
                    {record.type === "deposit" ? "+" : "-"}
                    {formatAmount(record.amount)}
                  </span>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      </div>

      {/* Deposit / Withdraw 弹层 */}
      <BottomSheet open={sheetMode !== null} onOpenChange={(o) => !o && setSheetMode(null)}>
        <h2 className="text-lg font-bold mb-4 text-text-primary">
          {sheetMode === "deposit" ? "充值保证金" : "提取保证金"}
        </h2>

        {rejectMsg && (
          <div
            className="mb-4 flex items-center gap-2 rounded-md bg-status-danger/10 px-3 py-2 text-xs text-status-danger"
            data-slot="reject-notice"
          >
            <AlertCircle className="size-4 shrink-0" />
            {rejectMsg}
          </div>
        )}

        {sheetMode === "deposit" && (
          <>
            <Input
              type="number"
              inputMode="decimal"
              value={amount}
              onChange={(e) => setAmount(e.target.value)}
              placeholder="输入充值金额（USDT）"
              aria-invalid={tooSmall || overBalance}
              className="mb-2"
            />
            {tooSmall && (
              <p className="mb-3 text-xs text-status-danger" data-slot="amount-error">
                单次充值不少于 {formatAmount(MIN_DEPOSIT_USDT)} USDT
              </p>
            )}
            {overBalance && (
              <p className="mb-3 text-xs text-status-danger" data-slot="amount-error">
                超出钱包 USDT 余额
              </p>
            )}

            <TwoStepAction
              owner={address}
              spender={depositSpender}
              amount={depositAmountValid ? amountMinUnit : 0n}
              action={depositAction}
              actionLabel="存入"
              onSuccess={handleDepositSuccess}
              disabled={!depositAmountValid}
            />
          </>
        )}

        {sheetMode === "withdraw" && (
          <div className="space-y-3">
            <p className="text-sm text-text-secondary">
              提取全部可用保证金（本金）。锁仓满后方可提现。
            </p>
            <Button
              onClick={() => {
                setRejectMsg(null);
                withdrawMutation.withdraw();
              }}
              disabled={!lockUnlocked || withdrawBusy}
              className="w-full"
            >
              {withdrawBusy && <Loader2 className="size-4 animate-spin" />}
              {withdrawTxState.status === "pending-signature"
                ? "请在钱包中确认…"
                : withdrawTxState.status === "pending-confirmation"
                  ? "交易确认中…"
                  : "提取本金"}
            </Button>
            {withdrawTxState.status === "error" && (
              <p className="text-xs text-status-danger">{withdrawTxState.error}</p>
            )}
          </div>
        )}
      </BottomSheet>
    </div>
  );
}

// 未部署链 getContractAddress 抛错；弹层未打开/钱包未连不应崩页，故包一层。
function safeSpender(chainId: number | undefined): `0x${string}` | undefined {
  if (chainId === undefined) return undefined;
  try {
    return getContractAddress(chainId, "Deposit");
  } catch {
    return undefined;
  }
}
