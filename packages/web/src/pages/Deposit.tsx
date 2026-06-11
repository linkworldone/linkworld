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
import { BottomSheet } from "@/components/shared/BottomSheet";
import { TwoStepAction, type ActionTx } from "@/components/shared/TwoStepAction";
import { LockCountdown } from "@/components/shared/LockCountdown";
import { AmountDisplay } from "@/components/shared/AmountDisplay";
import { TxStatusBadge } from "@/components/shared/TxStatusBadge";
import {
  useDeposit,
  useDepositHistory,
  useDepositMutation,
  useWithdrawMutation,
} from "@/hooks/useDeposit";
import { useUsdtBalance, useTranches, type Tranche } from "@/hooks/contracts";
import { getContractAddress } from "@/config/contracts";
import { arbitrumSepolia, hardhatLocal } from "@/config/chains";
import { WalletAuthRejectedError } from "@/services/api/signedPost";
import { formatAmount, formatDate } from "@/utils/format";
import type { DepositRecord } from "@/types";

type SheetMode = "deposit" | "withdraw" | null;

// 充值档位（合约 deposit 只接受 10/20/50/100 USDT，非档位 revert "Invalid tier"）。
// minUnit = 6 位最小单位；cards = 立即 mint 的无限流量卡张数（amount / 10 USDT）。
const DEPOSIT_TIERS = [
  { usdt: 10, minUnit: 10_000_000n, cards: 1 },
  { usdt: 20, minUnit: 20_000_000n, cards: 2 },
  { usdt: 50, minUnit: 50_000_000n, cards: 5 },
  { usdt: 100, minUnit: 100_000_000n, cards: 10 },
] as const;

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
  const { tranches, refetch: refetchTranches } = useTranches(address);

  const depositMutation = useDepositMutation();
  const withdrawMutation = useWithdrawMutation();

  const [sheetMode, setSheetMode] = useState<SheetMode>(null);
  // 选中的充值档位（DEPOSIT_TIERS 下标），null=未选。
  const [tierIdx, setTierIdx] = useState<number | null>(null);
  const [pending, setPending] = useState<PendingNotice | null>(null);
  const [rejectMsg, setRejectMsg] = useState<string | null>(null);
  // 正在取回的笔次下标（按笔禁用对应按钮）。
  const [withdrawingIdx, setWithdrawingIdx] = useState<number | null>(null);

  const usdtBalance = (usdtBalanceRaw as bigint | undefined) ?? undefined;
  const selectedTier = tierIdx !== null ? DEPOSIT_TIERS[tierIdx] : null;

  // ── Deposit 第二步动作（注入 TwoStepAction）：write 用选中档位最小单位；state 走 hook 统一 TxState ──
  const depositAction: ActionTx = useMemo(
    () => ({
      write: () => {
        setRejectMsg(null);
        if (selectedTier) depositMutation.deposit(String(selectedTier.usdt));
      },
      state: depositMutation.txState,
    }),
    [selectedTier, depositMutation.txState],
  );

  // 充值合约成功 → 上报 pending 意向（signedPost）；拒签提示，不进 pending。
  const handleDepositSuccess = async () => {
    if (!selectedTier) return;
    try {
      await depositMutation.recordIntent(String(selectedTier.usdt));
      refetchBalance();
      setPending({ kind: "deposit" });
      setSheetMode(null);
      setTierIdx(null);
    } catch (err) {
      if (err instanceof WalletAuthRejectedError) {
        setRejectMsg("身份签名被取消，操作未提交");
      } else {
        setRejectMsg("上报失败，稍后将自动重试");
      }
    }
  };

  // 提现（逐笔）合约成功 → 上报 pending 意向（无 tx_hash 记账，design §3.3）。
  useEffect(() => {
    if (!withdrawMutation.isSuccess) return;
    (async () => {
      try {
        await withdrawMutation.recordIntent();
        refetchBalance();
        refetchTranches();
        setPending({ kind: "withdraw" });
        setWithdrawingIdx(null);
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

  // 充值校验：选了档位 且 钱包 USDT 余额够该档位。
  const overBalance =
    selectedTier !== undefined &&
    selectedTier !== null &&
    usdtBalance !== undefined &&
    selectedTier.minUnit > usdtBalance;
  const depositValid = selectedTier !== null && !overBalance;

  const withdrawTxState = withdrawMutation.txState;
  const withdrawBusy =
    withdrawTxState.status === "pending-signature" ||
    withdrawTxState.status === "pending-confirmation";
  // 提现入口：有任意笔次即可打开弹层（按笔列出，逐笔取回/历史）；无记录则禁用。
  const hasTranches = tranches.length > 0;

  const handleWithdraw = (index: number) => {
    setRejectMsg(null);
    setWithdrawingIdx(index);
    withdrawMutation.withdraw(index);
  };

  const openSheet = (mode: SheetMode) => {
    setRejectMsg(null);
    setTierIdx(null);
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

      {/* 最早一笔锁仓倒计时（逐笔取回详情见提现弹层） */}
      <Card>
        <CardContent className="py-4 space-y-2">
          <LockCountdown address={address} />
          <p className="text-[11px] text-text-on-light-muted">
            每笔充值独立锁仓 30 天，各自到期后单独取回
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
          disabled={!hasTranches}
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
        <h2 className="text-lg font-bold mb-4 text-text-on-light-primary">
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
            <p className="mb-3 text-xs text-text-on-light-secondary">
              选择充值档位，充值后立即获得对应数量的无限流量卡。
            </p>
            <div className="mb-2 grid grid-cols-2 gap-2.5" data-slot="tier-grid">
              {DEPOSIT_TIERS.map((tier, idx) => {
                const active = tierIdx === idx;
                return (
                  <button
                    key={tier.usdt}
                    type="button"
                    onClick={() => setTierIdx(idx)}
                    aria-pressed={active}
                    data-slot="tier-option"
                    className={
                      "flex flex-col items-start rounded-lg border px-3 py-2.5 text-left transition-colors " +
                      (active
                        ? "border-brand-royal bg-brand-royal/5 ring-2 ring-brand-royal/40"
                        : "border-surface-card-line bg-surface-input hover:border-brand-royal/50")
                    }
                  >
                    <span className="text-base font-bold text-text-on-light-primary">
                      {tier.usdt} USDT
                    </span>
                    <span className="mt-0.5 text-[11px] text-text-on-light-secondary">
                      → 得 {tier.cards} 张无限流量卡
                    </span>
                  </button>
                );
              })}
            </div>
            {overBalance && (
              <p className="mb-3 text-xs text-status-danger" data-slot="amount-error">
                超出钱包 USDT 余额
              </p>
            )}

            <div className="mt-2">
              <TwoStepAction
                owner={address}
                spender={depositSpender}
                amount={depositValid && selectedTier ? selectedTier.minUnit : 0n}
                action={depositAction}
                actionLabel="存入"
                onSuccess={handleDepositSuccess}
                disabled={!depositValid}
              />
            </div>
          </>
        )}

        {sheetMode === "withdraw" && (
          <div className="space-y-3">
            <p className="text-sm text-text-on-light-secondary">
              每笔保证金独立锁仓 30 天，到期后可单独取回本金。
            </p>
            <div className="space-y-2" data-slot="tranche-list">
              {tranches.map((tranche: Tranche, idx: number) => (
                <TrancheRow
                  key={idx}
                  index={idx}
                  tranche={tranche}
                  onWithdraw={handleWithdraw}
                  busy={withdrawBusy && withdrawingIdx === idx}
                  disabledAll={withdrawBusy}
                />
              ))}
              {tranches.length === 0 && (
                <p className="text-xs text-text-on-light-muted">暂无保证金记录。</p>
              )}
            </div>
            {withdrawTxState.status === "error" && (
              <p className="text-xs text-status-danger">{withdrawTxState.error}</p>
            )}
          </div>
        )}
      </BottomSheet>
    </div>
  );
}

// 单笔保证金行：金额 + 解锁时间 + 状态（锁定中/可取回/已取回）。
// 到期且未取回 → 给「取回」按钮调 withdraw(index)；锁定中禁用并显示解锁日期。
function TrancheRow({
  index,
  tranche,
  onWithdraw,
  busy,
  disabledAll,
}: {
  index: number;
  tranche: Tranche;
  onWithdraw: (index: number) => void;
  busy: boolean;
  disabledAll: boolean;
}) {
  const nowSec = Math.floor(Date.now() / 1000);
  const unlocked = nowSec >= Number(tranche.unlockAt);
  const unlockDate = new Date(Number(tranche.unlockAt) * 1000).toLocaleDateString();

  let statusLabel: string;
  let statusClass: string;
  if (tranche.withdrawn) {
    statusLabel = "已取回";
    statusClass = "text-text-on-light-muted";
  } else if (unlocked) {
    statusLabel = "可取回";
    statusClass = "text-status-success";
  } else {
    statusLabel = `锁定中 · ${unlockDate} 解锁`;
    statusClass = "text-text-on-light-secondary";
  }

  return (
    <div
      className="flex items-center justify-between rounded-md bg-surface-input px-3 py-2.5"
      data-slot="tranche-row"
      data-tranche-index={index}
      data-withdrawn={tranche.withdrawn ? "true" : "false"}
    >
      <div className="min-w-0">
        <div className="text-sm font-semibold text-text-on-light-primary">
          {formatAmount(tranche.amount)} USDT
        </div>
        <div className={"mt-0.5 text-[11px] " + statusClass}>{statusLabel}</div>
      </div>
      {!tranche.withdrawn && (
        <Button
          size="sm"
          variant="outline"
          disabled={!unlocked || disabledAll}
          onClick={() => onWithdraw(index)}
          data-action="withdraw-tranche"
        >
          {busy && <Loader2 className="size-4 animate-spin" />}
          {busy ? "取回中…" : "取回"}
        </Button>
      )}
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
