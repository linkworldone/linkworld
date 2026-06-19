import { useEffect, useState } from "react";
import { useAccount } from "wagmi";
import { useTranslation } from "react-i18next";
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
  Smartphone,
  QrCode,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Tabs, TabsList, TabsTrigger, TabsContent } from "@/components/ui/tabs";
import { BottomSheet } from "@/components/shared/BottomSheet";
import { EmptyState } from "@/components/shared/EmptyState";
import { AmountDisplay } from "@/components/shared/AmountDisplay";
import { EsimActivation } from "@/components/shared/EsimActivation";
import {
  useTrafficCards,
  useRedeemForSim,
  type TrafficCardItem,
} from "@/hooks/contracts";
import { useMySims, useClaimSim } from "@/hooks/useSim";
import type { SimDeliveryType, SimRecord } from "@/services/api/simApi";
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
// 展示名走 i18n（t("destinations.<code>")），此处只保留 code 与开放状态。
const SIM_DESTINATIONS: { code: string; comingSoon?: boolean }[] = [
  { code: "TH" },
  { code: "JP" },
  { code: "KR" },
  { code: "SG" },
  { code: "MY" },
  { code: "VN" },
  { code: "ID" },
  { code: "PH" },
  { code: "HK" },
  { code: "MO" },
  { code: "TW" },
  { code: "AE" },
  { code: "MV" },
  { code: "AU", comingSoon: true },
  { code: "US", comingSoon: true },
  { code: "GB", comingSoon: true },
  { code: "FR", comingSoon: true },
];

// 销毁兑换 SIM 的最低选卡张数门槛（产品规则：至少 3 张才能 burn）。
const MIN_BURN = 3;

function TrafficCardsTab() {
  const { t } = useTranslation();
  const { address } = useAccount();
  const { cards, isLoading, isError, refetch } = useTrafficCards(address);
  const redeem = useRedeemForSim();
  const claimSim = useClaimSim(address);

  // 多选：选中的 tokenId 集合（以字符串存，避免 bigint Set 比较坑）。
  const [selected, setSelected] = useState<Set<string>>(new Set());
  const [sheetOpen, setSheetOpen] = useState(false);

  // 收件表单
  const [destination, setDestination] = useState(SIM_DESTINATIONS[0].code);
  // 交付方式：默认 eSIM（扫码激活，无需邮寄地址）。
  const [deliveryType, setDeliveryType] = useState<SimDeliveryType>("esim");
  const [recipient, setRecipient] = useState("");
  const [addressLine, setAddressLine] = useState("");
  // 流程状态：idle | redeeming(链上) | claiming(后端) | done | error
  const [submitErr, setSubmitErr] = useState<string | null>(null);
  const [done, setDone] = useState(false);
  // 兑换成功后拿到的 SIM 记录（esim 含 activationUrl，用于渲染二维码）。
  const [claimedSim, setClaimedSim] = useState<SimRecord | null>(null);

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
    setDeliveryType("esim");
    setSubmitErr(null);
    setDone(false);
    setClaimedSim(null);
    redeem.reset();
  }

  function openSheet() {
    if (selectedCount < MIN_BURN) return;
    resetForm();
    setSheetOpen(true);
  }

  // esim 只需选够卡数；physical 还需收件人 + 地址。
  const canSubmit =
    selectedCount >= MIN_BURN &&
    (deliveryType === "esim" ||
      (recipient.trim().length > 0 && addressLine.trim().length > 0));

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
    const isEsim = deliveryType === "esim";
    claimSim.mutate(
      {
        destination,
        deliveryType,
        recipient: isEsim ? "" : recipient.trim(),
        addressLine: isEsim ? "" : addressLine.trim(),
        tokenIds,
        txHash,
      },
      {
        onSuccess: (sim) => {
          setClaimedSim(sim);
          setDone(true);
          setSelected(new Set());
          refetch();
        },
        onError: (err) => {
          if (err instanceof WalletAuthRejectedError) {
            setSubmitErr(t("cards.authCancelled"));
          } else {
            setSubmitErr(
              err instanceof Error ? err.message.split("\n")[0] : t("cards.claimFailed"),
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
              {t("cards.instantTitle")}
            </div>
            <p className="text-xs leading-relaxed text-text-on-light-secondary">
              {t("cards.instantDesc")}
            </p>
            <p className="text-[11px] text-text-on-light-muted">
              {t("cards.instantNote")}
            </p>
          </div>
        </CardContent>
      </Card>

      {/* NFT 列表（读链真实数据，可多选） */}
      <div>
        <h3 className="mb-3 text-[13px] font-semibold text-text-primary">
          {!isLoading && !isError
            ? t("cards.myCardsCount", { count: cards.length })
            : t("cards.myCards")}
        </h3>

        {isLoading && (
          <div className="py-6 text-center text-xs text-text-secondary">{t("cards.loadingCards")}</div>
        )}

        {/* 加载失败 ≠ 真无卡（B6：禁静默空态冒充无卡） */}
        {!isLoading && isError && (
          <Card>
            <CardContent className="flex flex-col items-center gap-3 py-8 text-center">
              <AlertCircle className="size-8 text-status-danger" strokeWidth={1.5} />
              <p className="text-sm text-text-on-light-secondary">
                {t("cards.loadCardsError")}
              </p>
              <Button
                size="sm"
                variant="outline"
                className="border-surface-card-line bg-transparent text-text-on-light-primary hover:bg-surface-input"
                onClick={() => refetch()}
              >
                <RotateCw className="size-4" />
                {t("common.retry")}
              </Button>
            </CardContent>
          </Card>
        )}

        {/* 真无卡空态 */}
        {!isLoading && !isError && cards.length === 0 && (
          <Card>
            <CardContent className="py-2">
              <EmptyState icon={CreditCard} message={t("cards.emptyCards")} />
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
                            {t("cards.cardLabel", { id: card.tokenId.toString() })}
                          </div>
                          <div className="mt-0.5">
                            {isInfiniteData(card.dataAmount) ? (
                              <div>
                                <span className="text-lg font-bold text-text-on-light-primary">
                                  {t("cards.unlimitedData")}
                                </span>
                                <span className="ml-1.5 text-[11px] text-text-on-light-secondary">
                                  {t("cards.oneDay")}
                                </span>
                              </div>
                            ) : (
                              <AmountDisplay amount={card.dataAmount} currency="MB" />
                            )}
                          </div>
                          <div className="mt-0.5 text-[10px] text-text-on-light-muted">
                            {t("cards.issuedOn", { date: formatTimestamp(card.createdAt) })}
                          </div>
                        </div>
                      </div>
                    </CardContent>
                  </Card>
                </button>
              );
            })}
            <p className="px-1 text-[10px] text-text-secondary">
              {t("cards.selectHint")}
            </p>
          </div>
        )}
      </div>

      {/* 底部操作条（选中 > 0 时出现） */}
      {selectedCount > 0 && (
        <div className="fixed inset-x-0 bottom-[calc(64px+env(safe-area-inset-bottom))] z-40 mx-auto max-w-mobile border-t border-surface-card-line bg-surface-card p-4">
          <Button
            className="w-full py-3"
            onClick={openSheet}
            disabled={selectedCount < MIN_BURN}
          >
            <MailPlus className="size-4" />
            {selectedCount < MIN_BURN
              ? t("cards.minBurnHint", { min: MIN_BURN })
              : t("cards.burnAndRedeem", { count: selectedCount })}
          </Button>
        </div>
      )}

      {/* 收件信息弹窗 */}
      <BottomSheet open={sheetOpen} onOpenChange={closeSheet}>
        {done ? (
          claimedSim?.deliveryType === "esim" && claimedSim.activationUrl ? (
            <div className="flex flex-col items-center gap-4 py-2 text-center">
              <CheckCircle2 className="size-10 text-brand-gold" strokeWidth={1.5} />
              <p className="text-base font-semibold text-text-on-light-primary">
                {t("cards.esimReadyTitle")}
              </p>
              <EsimActivation activationUrl={claimedSim.activationUrl} />
              <Button className="mt-2 w-full py-3" onClick={() => closeSheet(false)}>
                {t("common.done")}
              </Button>
            </div>
          ) : (
            <div className="flex flex-col items-center gap-3 py-4 text-center">
              <CheckCircle2 className="size-10 text-brand-gold" strokeWidth={1.5} />
              <p className="text-base font-semibold text-text-on-light-primary">
                {t("cards.redeemSubmittedTitle")}
              </p>
              <p className="text-xs text-text-on-light-secondary">
                {t("cards.redeemSubmittedDesc")}
              </p>
              <Button className="mt-2 w-full py-3" onClick={() => closeSheet(false)}>
                {t("common.done")}
              </Button>
            </div>
          )
        ) : (
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="space-y-1">
              <h2 className="text-base font-bold text-text-on-light-primary">
                {t("cards.redeemTitle", { count: selectedCount })}
              </h2>
              <p className="text-xs text-text-on-light-secondary">
                {t("cards.redeemDesc", { count: selectedCount })}
              </p>
            </div>

            <div className="space-y-1.5">
              <label className="text-xs font-medium text-text-on-light-secondary">{t("cards.destination")}</label>
              <select
                value={destination}
                onChange={(e) => setDestination(e.target.value)}
                disabled={busy}
                className="flex h-11 w-full rounded-md border border-surface-card-line bg-surface-input px-3 py-2 text-base text-text-on-light-primary outline-none transition-colors focus-visible:border-brand-royal focus-visible:ring-2 focus-visible:ring-brand-royal/40 disabled:opacity-50"
              >
                {SIM_DESTINATIONS.map((d) => (
                  <option key={d.code} value={d.code} disabled={d.comingSoon}>
                    {t(`destinations.${d.code}`)}{d.comingSoon ? t("common.comingSoon") : ""}
                  </option>
                ))}
              </select>
            </div>

            {/* 交付方式：eSIM（扫码激活）/ 实体邮寄 */}
            <div className="space-y-1.5">
              <label className="text-xs font-medium text-text-on-light-secondary">{t("cards.deliveryMethod")}</label>
              <div className="grid grid-cols-2 gap-2">
                {(
                  [
                    { value: "esim" as const, label: t("cards.esim"), icon: Smartphone },
                    { value: "physical" as const, label: t("cards.physicalMail"), icon: MailPlus },
                  ]
                ).map(({ value, label, icon: Icon }) => {
                  const active = deliveryType === value;
                  return (
                    <button
                      key={value}
                      type="button"
                      onClick={() => !busy && setDeliveryType(value)}
                      disabled={busy}
                      aria-pressed={active}
                      className={
                        "flex h-11 items-center justify-center gap-1.5 rounded-md border text-sm font-medium transition-colors disabled:opacity-50 " +
                        (active
                          ? "border-brand-royal bg-brand-royal/10 text-brand-royal"
                          : "border-surface-card-line bg-surface-input text-text-on-light-secondary")
                      }
                    >
                      <Icon className="size-4" />
                      {label}
                    </button>
                  );
                })}
              </div>
            </div>

            {deliveryType === "esim" ? (
              <p className="text-[11px] text-text-on-light-muted">{t("cards.esimNoAddressNote")}</p>
            ) : (
              <>
                <div className="space-y-1.5">
                  <label className="text-xs font-medium text-text-on-light-secondary">{t("cards.recipient")}</label>
                  <Input
                    value={recipient}
                    onChange={(e) => setRecipient(e.target.value)}
                    placeholder={t("cards.recipientPlaceholder")}
                    disabled={busy}
                  />
                </div>

                <div className="space-y-1.5">
                  <label className="text-xs font-medium text-text-on-light-secondary">{t("cards.addressLine")}</label>
                  <Input
                    value={addressLine}
                    onChange={(e) => setAddressLine(e.target.value)}
                    placeholder={t("cards.addressPlaceholder")}
                    disabled={busy}
                  />
                </div>
              </>
            )}

            {submitErr && (
              <p className="text-center text-[11px] text-status-danger">{submitErr}</p>
            )}

            <Button type="submit" className="w-full py-3" disabled={!canSubmit || busy}>
              {busy && <Loader2 className="size-4 animate-spin" />}
              {redeem.isPending
                ? t("cards.btnConfirmWallet")
                : redeem.isConfirming
                  ? t("cards.btnBurning")
                  : claimSim.isPending
                    ? t("cards.btnSubmittingInfo")
                    : t("cards.btnBurnAndRedeem", { count: selectedCount })}
            </Button>
          </form>
        )}
      </BottomSheet>
    </div>
  );
}

function MySimsTab() {
  const { t } = useTranslation();
  const { address } = useAccount();
  const { data: sims, isLoading, isError, refetch } = useMySims(address);
  // 查看 eSIM 二维码弹窗（点击 eSIM 记录的「查看」打开）。
  const [viewSim, setViewSim] = useState<SimRecord | null>(null);

  return (
    <div className="space-y-4">
      <Card>
        <CardContent className="flex gap-3">
          <Nfc className="mt-0.5 size-5 shrink-0 text-brand-gold" />
          <div className="space-y-1">
            <div className="text-sm font-semibold text-text-on-light-primary">
              {t("cards.mySimTitle")}
            </div>
            <p className="text-xs leading-relaxed text-text-on-light-secondary">
              {t("cards.mySimDesc")}
            </p>
            <p className="text-[11px] text-text-on-light-muted">
              {t("cards.mySimNote")}
            </p>
          </div>
        </CardContent>
      </Card>

      {isLoading && (
        <div className="py-6 text-center text-xs text-text-secondary">{t("cards.loadingSims")}</div>
      )}

      {!isLoading && isError && (
        <Card>
          <CardContent className="flex flex-col items-center gap-3 py-8 text-center">
            <AlertCircle className="size-8 text-status-danger" strokeWidth={1.5} />
            <p className="text-sm text-text-on-light-secondary">{t("cards.loadSimsError")}</p>
            <Button
              size="sm"
              variant="outline"
              className="border-surface-card-line bg-transparent text-text-on-light-primary hover:bg-surface-input"
              onClick={() => refetch()}
            >
              <RotateCw className="size-4" />
              {t("common.retry")}
            </Button>
          </CardContent>
        </Card>
      )}

      {!isLoading && !isError && (!sims || sims.length === 0) && (
        <Card>
          <CardContent className="py-2">
            <EmptyState icon={Nfc} message={t("cards.emptySims")} />
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
                      {t("cards.simDays", { days: sim.days })}
                    </span>
                    <span className="ml-2 text-xs text-text-on-light-secondary">
                      {t(`destinations.${sim.destination}`, { defaultValue: sim.destination })}
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
                    {sim.status === "confirmed" ? t("cards.simConfirmed") : t("cards.simPending")}
                  </span>
                </div>
                {sim.deliveryType === "esim" ? (
                  sim.activationUrl ? (
                    <Button
                      size="sm"
                      className="w-full"
                      onClick={() => setViewSim(sim)}
                    >
                      <QrCode className="size-3.5" />
                      {t("cards.viewEsim")}
                    </Button>
                  ) : null
                ) : (
                  <div className="text-[11px] text-text-on-light-muted">
                    {t("cards.simRecipient", {
                      recipient: sim.recipient || "—",
                      address: sim.addressLine || "—",
                    })}
                  </div>
                )}
              </CardContent>
            </Card>
          ))}
        </div>
      )}

      {/* 查看 eSIM 二维码弹窗 */}
      <BottomSheet open={!!viewSim} onOpenChange={(open) => !open && setViewSim(null)}>
        {viewSim && (
          <div className="flex flex-col items-center gap-4 py-2 text-center">
            <p className="text-base font-semibold text-text-on-light-primary">
              {t("cards.esimReadyTitle")}
            </p>
            <EsimActivation activationUrl={viewSim.activationUrl} />
            <Button className="mt-2 w-full py-3" onClick={() => setViewSim(null)}>
              {t("common.done")}
            </Button>
          </div>
        )}
      </BottomSheet>
    </div>
  );
}

export default function Cards() {
  const { t } = useTranslation();
  return (
    <div className="px-4">
      <Tabs defaultValue="traffic">
        <TabsList variant="line" className="w-full">
          <TabsTrigger value="traffic">{t("cards.tabTraffic")}</TabsTrigger>
          <TabsTrigger value="sim">{t("cards.tabSim")}</TabsTrigger>
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
