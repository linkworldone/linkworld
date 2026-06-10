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
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Tabs, TabsList, TabsTrigger, TabsContent } from "@/components/ui/tabs";
import { EmptyState } from "@/components/shared/EmptyState";
import { AmountDisplay } from "@/components/shared/AmountDisplay";
import {
  useTrafficCards,
  useBurnCard,
  type TrafficCardItem,
} from "@/hooks/contracts";
import { savePendingSync } from "@/utils/pendingSync";

function formatTimestamp(ts: bigint): string {
  if (ts === 0n) return "—";
  return new Date(Number(ts) * 1000).toLocaleDateString();
}

// SIM 领取目的地（降级 R9：全球通即将推出，先开放区域 eSIM 领取意向）。
const SIM_DESTINATIONS = [
  { code: "US", name: "美国 United States" },
  { code: "JP", name: "日本 Japan" },
  { code: "GB", name: "英国 United Kingdom" },
  { code: "SG", name: "新加坡 Singapore" },
  { code: "HK", name: "中国香港 Hong Kong" },
];

function TrafficCardsTab() {
  const { address } = useAccount();
  const { cards, isLoading, isError, refetch } = useTrafficCards(address);
  const burn = useBurnCard();

  useEffect(() => {
    if (burn.isSuccess) {
      refetch();
      burn.reset();
    }
  }, [burn.isSuccess]);

  const burnBusy = burn.isPending || burn.isConfirming;

  return (
    <div className="space-y-4">
      {/* 自动发放说明（移除 Admin 发卡按钮：onlyOracle，用户调 revert） */}
      <Card>
        <CardContent className="flex gap-3">
          <Sparkles className="mt-0.5 size-5 shrink-0 text-brand-gold" />
          <div className="space-y-1">
            <div className="flex items-center gap-1.5 text-sm font-semibold text-text-on-light-primary">
              <Info className="size-4 text-brand-royal" />
              流量卡自动发放
            </div>
            <p className="text-xs leading-relaxed text-text-on-light-secondary">
              锁仓满 1 个月后，系统将按你的保证金自动发放流量卡，无需手动领取。
            </p>
            <p className="text-[11px] text-text-on-light-muted">
              流量卡不可转卖；销毁后剩余流量额度有效期为 30 天。
            </p>
          </div>
        </CardContent>
      </Card>

      {/* NFT 列表（读链真实数据） */}
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
              <Button size="sm" variant="outline" onClick={() => refetch()}>
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
              <EmptyState icon={CreditCard} message="暂无流量卡。锁仓满 1 个月后将自动发放。" />
            </CardContent>
          </Card>
        )}

        {!isLoading && !isError && cards.length > 0 && (
          <div className="space-y-2">
            {cards.map((card: TrafficCardItem) => (
              <Card key={card.tokenId.toString()}>
                <CardContent className="flex items-center justify-between">
                  <div>
                    <div className="text-xs text-text-on-light-muted">
                      流量卡 #{card.tokenId.toString()}
                    </div>
                    <div className="mt-0.5">
                      <AmountDisplay amount={card.dataAmount} currency="MB" />
                    </div>
                    <div className="mt-0.5 text-[10px] text-text-on-light-muted">
                      发放于 {formatTimestamp(card.createdAt)} · 不可转卖
                    </div>
                  </div>
                  <Button
                    size="sm"
                    variant="outline"
                    disabled={burnBusy}
                    onClick={() => burn.burnCard(card.tokenId)}
                  >
                    {burnBusy ? "销毁中…" : "销毁兑换额度"}
                  </Button>
                </CardContent>
              </Card>
            ))}
            <p className="px-1 text-[10px] text-text-secondary">
              销毁后剩余流量额度有效期为 30 天。
            </p>
          </div>
        )}

        {burn.error && (
          <p className="mt-2 text-center text-[10px] text-status-danger">
            {burn.error.message.split("\n")[0]}
          </p>
        )}
      </div>
    </div>
  );
}

function SimClaimTab() {
  const { address } = useAccount();
  const [destination, setDestination] = useState(SIM_DESTINATIONS[0].code);
  const [recipient, setRecipient] = useState("");
  const [addressLine, setAddressLine] = useState("");
  const [submitted, setSubmitted] = useState(false);

  const canSubmit = recipient.trim().length > 0 && addressLine.trim().length > 0;

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!canSubmit) return;
    // 降级 R9：后端 SIM 领取端点未就绪 → 先写 pendingSync，待端点可用后补提交。
    savePendingSync("sim_claim", {
      wallet: address ?? null,
      destination,
      recipient: recipient.trim(),
      addressLine: addressLine.trim(),
    });
    setSubmitted(true);
  };

  return (
    <div className="space-y-4">
      <Card>
        <CardContent className="flex gap-3">
          <Nfc className="mt-0.5 size-5 shrink-0 text-brand-gold" />
          <div className="space-y-1">
            <div className="text-sm font-semibold text-text-on-light-primary">
              领取实体 SIM
            </div>
            <p className="text-xs leading-relaxed text-text-on-light-secondary">
              填写收件信息，我们将为你寄送目的地 SIM 卡。
            </p>
            <p className="text-[11px] text-text-on-light-muted">
              全球通 eSIM 即将推出，敬请期待。
            </p>
          </div>
        </CardContent>
      </Card>

      {submitted ? (
        <Card>
          <CardContent className="flex flex-col items-center gap-3 py-8 text-center">
            <CheckCircle2 className="size-8 text-brand-gold" strokeWidth={1.5} />
            <p className="text-sm font-semibold text-text-on-light-primary">
              领取申请已提交
            </p>
            <p className="text-xs text-text-on-light-secondary">
              我们已记录你的收件信息，SIM 卡将尽快寄出。
            </p>
          </CardContent>
        </Card>
      ) : (
        <form onSubmit={handleSubmit} className="space-y-3">
          <div className="space-y-1.5">
            <label className="text-xs font-medium text-text-secondary">目的地</label>
            <select
              value={destination}
              onChange={(e) => setDestination(e.target.value)}
              className="flex h-11 w-full rounded-md border border-surface-card-line bg-surface-input px-3 py-2 text-base text-text-on-light-primary outline-none transition-colors focus-visible:border-brand-royal focus-visible:ring-2 focus-visible:ring-brand-royal/40"
            >
              {SIM_DESTINATIONS.map((d) => (
                <option key={d.code} value={d.code}>
                  {d.name}
                </option>
              ))}
            </select>
          </div>

          <div className="space-y-1.5">
            <label className="text-xs font-medium text-text-secondary">收件人</label>
            <Input
              value={recipient}
              onChange={(e) => setRecipient(e.target.value)}
              placeholder="收件人姓名"
            />
          </div>

          <div className="space-y-1.5">
            <label className="text-xs font-medium text-text-secondary">收件地址</label>
            <Input
              value={addressLine}
              onChange={(e) => setAddressLine(e.target.value)}
              placeholder="详细收件地址"
            />
          </div>

          <Button type="submit" className="w-full py-3" disabled={!canSubmit}>
            <MailPlus className="size-4" />
            提交领取申请
          </Button>
        </form>
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
          <TabsTrigger value="sim">SIM 领取</TabsTrigger>
        </TabsList>
        <TabsContent value="traffic" className="mt-4">
          <TrafficCardsTab />
        </TabsContent>
        <TabsContent value="sim" className="mt-4">
          <SimClaimTab />
        </TabsContent>
      </Tabs>
    </div>
  );
}
