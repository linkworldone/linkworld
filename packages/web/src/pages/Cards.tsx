import { useEffect, useState } from "react";
import { useAccount } from "wagmi";
import {
  Sparkles,
  Info,
  CreditCard,
  AlertCircle,
  RotateCw,
  Nfc,
  MailPlus,
  CheckCircle2,
  Check,
  Loader2,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Tabs, TabsList, TabsTrigger, TabsContent } from "@/components/ui/tabs";
import { BottomSheet } from "@/components/shared/BottomSheet";
import { EmptyState } from "@/components/shared/EmptyState";
import { AmountDisplay } from "@/components/shared/AmountDisplay";
import {
  useTrafficCards,
  useRedeemForSim,
  type TrafficCardItem,
} from "@/hooks/contracts";
import { useMySims, useClaimSim } from "@/hooks/useSim";
import { WalletAuthRejectedError } from "@/services/api/signedPost";

function formatTimestamp(ts: bigint): string {
  if (ts === 0n) return "—";
  return new Date(Number(ts) * 1000).toLocaleDateString();
}

// 无限流量哨兵值：合约 dataAmount = type(uint256).max。用阈值判定（避免精确比较被极大值绕过）。
const UINT256_MAX = 2n ** 256n - 1n;
// 任何 ≥ 2^200 的值都远超真实流量字节数，按「无限流量」处理。
const INFINITE_THRESHOLD = 2n ** 200n;
function isInfiniteData(dataAmount: bigint): boolean {
  return dataAmount >= INFINITE_THRESHOLD && dataAmount <= UINT256_MAX;
}

// SIM 领取目的地（先以热门旅游地为主：东南亚/东亚数据成本低、出行频次高）。
// comingSoon: 高成本区（美欧澳）数据批发贵，无限卡暂不开放，下拉置灰标「即将开放」。
const SIM_DESTINATIONS: { code: string; name: string; comingSoon?: boolean }[] = [
  { code: "TH", name: "泰国 Thailand" },
  { code: "JP", name: "日本 Japan" },
  { code: "KR", name: "韩国 South Korea" },
  { code: "SG", name: "新加坡 Singapore" },
  { code: "MY", name: "马来西亚 Malaysia" },
  { code: "VN", name: "越南 Vietnam" },
  { code: "ID", name: "印度尼西亚 Indonesia" },
  { code: "PH", name: "菲律宾 Philippines" },
  { code: "HK", name: "中国香港 Hong Kong" },
  { code: "MO", name: "中国澳门 Macau" },
  { code: "TW", name: "中国台湾 Taiwan" },
  { code: "AE", name: "阿联酋 UAE" },
  { code: "MV", name: "马尔代夫 Maldives" },
  { code: "AU", name: "澳大利亚 Australia", comingSoon: true },
  { code: "US", name: "美国 United States", comingSoon: true },
  { code: "GB", name: "英国 United Kingdom", comingSoon: true },
  { code: "FR", name: "法国 France", comingSoon: true },
];

function destinationName(code: string): string {
  return SIM_DESTINATIONS.find((d) => d.code === code)?.name ?? code;
}

function TrafficCardsTab() {
  const { address } = useAccount();
  const { cards, isLoading, isError, refetch } = useTrafficCards(address);
  const redeem = useRedeemForSim();
  const claimSim = useClaimSim(address);

  // 多选：选中的 tokenId 集合（以字符串存，避免 bigint Set 比较坑）。
  const [selected, setSelected] = useState<Set<string>>(new Set());
  const [sheetOpen, setSheetOpen] = useState(false);

  // 收件表单
  const [destination, setDestination] = useState(SIM_DESTINATIONS[0].code);
  const [recipient, setRecipient] = useState("");
  const [addressLine, setAddressLine] = useState("");
  // 流程状态：idle | redeeming(链上) | claiming(后端) | done | error
  const [submitErr, setSubmitErr] = useState<string | null>(null);
  const [done, setDone] = useState(false);

  const selectedCount = selected.size;

  function toggle(tokenId: bigint) {
    const key = tokenId.toString();
    setSelected((prev) => {
      const next = new Set(prev);
      if (next.has(key)) next.delete(key);
      else next.add(key);
      return next;
    });
  }

  function resetForm() {
    setRecipient("");
    setAddressLine("");
    setDestination(SIM_DESTINATIONS[0].code);
    setSubmitErr(null);
    setDone(false);
    redeem.reset();
  }

  function openSheet() {
    if (selectedCount === 0) return;
    resetForm();
    setSheetOpen(true);
  }

  const canSubmit =
    recipient.trim().length > 0 && addressLine.trim().length > 0 && selectedCount > 0;

  // 选中的 tokenId（bigint / number 两种形态备用）
  const selectedTokenIds = cards
    .filter((c) => selected.has(c.tokenId.toString()))
    .map((c) => c.tokenId);

  // 提交：① 链上 redeemForSim 批量销毁 → ② 成功后 simApi.claim。
  function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (!canSubmit || busy) return;
    setSubmitErr(null);
    redeem.redeem(selectedTokenIds);
  }

  // 链上销毁成功 → 调后端 claim。
  useEffect(() => {
    if (!redeem.isSuccess || !redeem.hash) return;
    if (claimSim.isPending || done) return;
    const txHash = redeem.hash;
    const tokenIds = selectedTokenIds.map((id) => Number(id));
    claimSim.mutate(
      {
        destination,
        recipient: recipient.trim(),
        addressLine: addressLine.trim(),
        tokenIds,
        txHash,
      },
      {
        onSuccess: () => {
          setDone(true);
          setSelected(new Set());
          refetch();
        },
        onError: (err) => {
          if (err instanceof WalletAuthRejectedError) {
            setSubmitErr("身份签名被取消，SIM 领取未提交");
          } else {
            setSubmitErr(
              err instanceof Error ? err.message.split("\n")[0] : "SIM 领取提交失败",
            );
          }
        },
      },
    );
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [redeem.isSuccess, redeem.hash]);

  const busy = redeem.isPending || redeem.isConfirming || claimSim.isPending;

  function closeSheet(open: boolean) {
    if (busy) return; // 进行中不允许关闭，避免中途打断
    setSheetOpen(open);
    if (!open && done) {
      // 关闭已完成弹窗：清理表单（选中已清空）
      resetForm();
    }
  }

  return (
    <div className="space-y-4">
      {/* 充值即发说明 */}
      <Card>
        <CardContent className="flex gap-3">
          <Sparkles className="mt-0.5 size-5 shrink-0 text-brand-gold" />
          <div className="space-y-1">
            <div className="flex items-center gap-1.5 text-sm font-semibold text-text-on-light-primary">
              <Info className="size-4 text-brand-royal" />
              充值即得流量卡
            </div>
            <p className="text-xs leading-relaxed text-text-on-light-secondary">
              充值保证金即时发放流量卡：每 10 USDT 得 1 张「1 天无限流量」卡（10/20/50/100 → 1/2/5/10 张），无需等待。
            </p>
            <p className="text-[11px] text-text-on-light-muted">
              勾选多张流量卡，一步销毁即可领取对应天数的 SIM；流量卡不可转卖。
            </p>
          </div>
        </CardContent>
      </Card>

      {/* NFT 列表（读链真实数据，可多选） */}
      <div>
        <h3 className="mb-3 text-[13px] font-semibold text-text-primary">
          我的流量卡{!isLoading && !isError ? `（${cards.length}）` : ""}
        </h3>

        {isLoading && (
          <div className="py-6 text-center text-xs text-text-secondary">加载流量卡…</div>
        )}

        {/* 加载失败 ≠ 真无卡（B6：禁静默空态冒充无卡） */}
        {!isLoading && isError && (
          <Card>
            <CardContent className="flex flex-col items-center gap-3 py-8 text-center">
              <AlertCircle className="size-8 text-status-danger" strokeWidth={1.5} />
              <p className="text-sm text-text-on-light-secondary">
                流量卡加载失败，请重试。
              </p>
              <Button
                size="sm"
                variant="outline"
                className="border-surface-card-line bg-transparent text-text-on-light-primary hover:bg-surface-input"
                onClick={() => refetch()}
              >
                <RotateCw className="size-4" />
                重试
              </Button>
            </CardContent>
          </Card>
        )}

        {/* 真无卡空态 */}
        {!isLoading && !isError && cards.length === 0 && (
          <Card>
            <CardContent className="py-2">
              <EmptyState icon={CreditCard} message="暂无流量卡。充值保证金即可获得。" />
            </CardContent>
          </Card>
        )}

        {!isLoading && !isError && cards.length > 0 && (
          <div className="space-y-2 pb-40">
            {cards.map((card: TrafficCardItem) => {
              const checked = selected.has(card.tokenId.toString());
              return (
                <button
                  type="button"
                  key={card.tokenId.toString()}
                  onClick={() => toggle(card.tokenId)}
                  aria-pressed={checked}
                  className="block w-full text-left"
                >
                  <Card
                    className={
                      checked
                        ? "border-brand-royal ring-2 ring-brand-royal/30"
                        : "border-surface-card-line"
                    }
                  >
                    <CardContent className="flex items-center justify-between gap-3">
                      <div className="flex items-center gap-3">
                        {/* 多选指示 */}
                        <span
                          className={
                            "flex size-5 shrink-0 items-center justify-center rounded-md border " +
                            (checked
                              ? "border-brand-royal bg-brand-royal text-white"
                              : "border-surface-card-line bg-transparent")
                          }
                          aria-hidden
                        >
                          {checked && <Check className="size-3.5" strokeWidth={3} />}
                        </span>
                        <div>
                          <div className="text-xs text-text-on-light-muted">
                            流量卡 #{card.tokenId.toString()}
                          </div>
                          <div className="mt-0.5">
                            {isInfiniteData(card.dataAmount) ? (
                              <div>
                                <span className="text-lg font-bold text-text-on-light-primary">
                                  无限流量
                                </span>
                                <span className="ml-1.5 text-[11px] text-text-on-light-secondary">
                                  1 天
                                </span>
                              </div>
                            ) : (
                              <AmountDisplay amount={card.dataAmount} currency="MB" />
                            )}
                          </div>
                          <div className="mt-0.5 text-[10px] text-text-on-light-muted">
                            发放于 {formatTimestamp(card.createdAt)} · 不可转卖
                          </div>
                        </div>
                      </div>
                    </CardContent>
                  </Card>
                </button>
              );
            })}
            <p className="px-1 text-[10px] text-text-secondary">
              勾选多张流量卡 → 一步销毁领取 SIM（每张 = 1 天）。
            </p>
          </div>
        )}
      </div>

      {/* 底部操作条（选中 > 0 时出现） */}
      {selectedCount > 0 && (
        <div className="fixed inset-x-0 bottom-[calc(64px+env(safe-area-inset-bottom))] z-40 mx-auto max-w-mobile border-t border-surface-card-line bg-surface-card p-4">
          <Button className="w-full py-3" onClick={openSheet}>
            <MailPlus className="size-4" />
            销毁 {selectedCount} 张并领取 SIM（{selectedCount} 天）
          </Button>
        </div>
      )}

      {/* 收件信息弹窗 */}
      <BottomSheet open={sheetOpen} onOpenChange={closeSheet}>
        {done ? (
          <div className="flex flex-col items-center gap-3 py-4 text-center">
            <CheckCircle2 className="size-10 text-brand-gold" strokeWidth={1.5} />
            <p className="text-base font-semibold text-text-on-light-primary">
              SIM 领取已提交
            </p>
            <p className="text-xs text-text-on-light-secondary">
              已销毁流量卡并记录你的收件信息，SIM 卡将尽快寄出。可在「我的 SIM」查看进度。
            </p>
            <Button className="mt-2 w-full py-3" onClick={() => closeSheet(false)}>
              完成
            </Button>
          </div>
        ) : (
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="space-y-1">
              <h2 className="text-base font-bold text-text-on-light-primary">
                领取 SIM（{selectedCount} 天）
              </h2>
              <p className="text-xs text-text-on-light-secondary">
                将销毁 {selectedCount} 张流量卡兑换 {selectedCount} 天 SIM，填写收件信息后提交。
              </p>
            </div>

            <div className="space-y-1.5">
              <label className="text-xs font-medium text-text-on-light-secondary">目的地</label>
              <select
                value={destination}
                onChange={(e) => setDestination(e.target.value)}
                disabled={busy}
                className="flex h-11 w-full rounded-md border border-surface-card-line bg-surface-input px-3 py-2 text-base text-text-on-light-primary outline-none transition-colors focus-visible:border-brand-royal focus-visible:ring-2 focus-visible:ring-brand-royal/40 disabled:opacity-50"
              >
                {SIM_DESTINATIONS.map((d) => (
                  <option key={d.code} value={d.code} disabled={d.comingSoon}>
                    {d.name}{d.comingSoon ? "（即将开放）" : ""}
                  </option>
                ))}
              </select>
            </div>

            <div className="space-y-1.5">
              <label className="text-xs font-medium text-text-on-light-secondary">收件人</label>
              <Input
                value={recipient}
                onChange={(e) => setRecipient(e.target.value)}
                placeholder="收件人姓名"
                disabled={busy}
              />
            </div>

            <div className="space-y-1.5">
              <label className="text-xs font-medium text-text-on-light-secondary">收件地址</label>
              <Input
                value={addressLine}
                onChange={(e) => setAddressLine(e.target.value)}
                placeholder="详细收件地址"
                disabled={busy}
              />
            </div>

            {submitErr && (
              <p className="text-center text-[11px] text-status-danger">{submitErr}</p>
            )}

            <Button type="submit" className="w-full py-3" disabled={!canSubmit || busy}>
              {busy && <Loader2 className="size-4 animate-spin" />}
              {redeem.isPending
                ? "确认钱包签名…"
                : redeem.isConfirming
                  ? "链上销毁中…"
                  : claimSim.isPending
                    ? "提交收件信息…"
                    : `销毁并领取 SIM（${selectedCount} 天）`}
            </Button>
          </form>
        )}
      </BottomSheet>
    </div>
  );
}

function MySimsTab() {
  const { address } = useAccount();
  const { data: sims, isLoading, isError, refetch } = useMySims(address);

  return (
    <div className="space-y-4">
      <Card>
        <CardContent className="flex gap-3">
          <Nfc className="mt-0.5 size-5 shrink-0 text-brand-gold" />
          <div className="space-y-1">
            <div className="text-sm font-semibold text-text-on-light-primary">
              我的 SIM
            </div>
            <p className="text-xs leading-relaxed text-text-on-light-secondary">
              销毁流量卡领取的 SIM 在此查看，含天数余额与配送状态。
            </p>
            <p className="text-[11px] text-text-on-light-muted">
              全球通 eSIM 即将推出，敬请期待。
            </p>
          </div>
        </CardContent>
      </Card>

      {isLoading && (
        <div className="py-6 text-center text-xs text-text-secondary">加载我的 SIM…</div>
      )}

      {!isLoading && isError && (
        <Card>
          <CardContent className="flex flex-col items-center gap-3 py-8 text-center">
            <AlertCircle className="size-8 text-status-danger" strokeWidth={1.5} />
            <p className="text-sm text-text-on-light-secondary">SIM 加载失败，请重试。</p>
            <Button
              size="sm"
              variant="outline"
              className="border-surface-card-line bg-transparent text-text-on-light-primary hover:bg-surface-input"
              onClick={() => refetch()}
            >
              <RotateCw className="size-4" />
              重试
            </Button>
          </CardContent>
        </Card>
      )}

      {!isLoading && !isError && (!sims || sims.length === 0) && (
        <Card>
          <CardContent className="py-2">
            <EmptyState icon={Nfc} message="暂无 SIM。销毁流量卡即可领取 SIM。" />
          </CardContent>
        </Card>
      )}

      {!isLoading && !isError && sims && sims.length > 0 && (
        <div className="space-y-2">
          {sims.map((sim) => (
            <Card key={sim.id}>
              <CardContent className="space-y-2">
                <div className="flex items-center justify-between">
                  <div>
                    <span className="text-lg font-bold text-text-on-light-primary">
                      {sim.days} 天
                    </span>
                    <span className="ml-2 text-xs text-text-on-light-secondary">
                      {destinationName(sim.destination)}
                    </span>
                  </div>
                  <span
                    className={
                      "rounded-full px-2 py-0.5 text-[11px] font-medium " +
                      (sim.status === "confirmed"
                        ? "bg-status-success/15 text-status-success"
                        : "bg-brand-gold/15 text-brand-gold")
                    }
                  >
                    {sim.status === "confirmed" ? "已确认" : "处理中"}
                  </span>
                </div>
                <div className="text-[11px] text-text-on-light-muted">
                  收件人：{sim.recipient || "—"} · {sim.addressLine || "—"}
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}

export default function Cards() {
  return (
    <div className="px-4">
      <Tabs defaultValue="traffic">
        <TabsList variant="line" className="w-full">
          <TabsTrigger value="traffic">流量卡</TabsTrigger>
          <TabsTrigger value="sim">我的 SIM</TabsTrigger>
        </TabsList>
        <TabsContent value="traffic" className="mt-4">
          <TrafficCardsTab />
        </TabsContent>
        <TabsContent value="sim" className="mt-4">
          <MySimsTab />
        </TabsContent>
      </Tabs>
    </div>
  );
}
