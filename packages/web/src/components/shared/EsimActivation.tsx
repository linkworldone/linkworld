import { useState } from "react";
import { useTranslation } from "react-i18next";
import { QRCodeSVG } from "qrcode.react";
import { Check, Copy } from "lucide-react";
import { Button } from "@/components/ui/button";

interface EsimActivationProps {
  activationUrl: string;
}

// eSIM 激活展示（纯展示）：二维码 + 激活链接 + 复制按钮。
// 兑换成功态与「我的 SIM 查看」共用，避免重复实现。
export function EsimActivation({ activationUrl }: EsimActivationProps) {
  const { t } = useTranslation();
  const [copied, setCopied] = useState(false);

  async function copyUrl() {
    try {
      await navigator.clipboard.writeText(activationUrl);
      setCopied(true);
      setTimeout(() => setCopied(false), 1500);
    } catch {
      // 复制失败静默（部分环境无 clipboard 权限）；用户仍可手动选取链接文本。
    }
  }

  return (
    <div className="flex flex-col items-center gap-3">
      <div className="rounded-xl bg-white p-3">
        <QRCodeSVG value={activationUrl} size={180} />
      </div>
      <p className="text-xs text-text-on-light-secondary">{t("cards.scanToActivate")}</p>
      <div className="w-full space-y-1.5">
        <div className="text-[11px] font-medium text-text-on-light-muted">
          {t("cards.activationUrl")}
        </div>
        <p className="break-all rounded-md border border-surface-card-line bg-surface-input px-3 py-2 text-[11px] text-text-on-light-primary">
          {activationUrl}
        </p>
        <Button
          type="button"
          variant="outline"
          size="sm"
          className="w-full border-surface-card-line bg-transparent text-text-on-light-primary hover:bg-surface-input"
          onClick={copyUrl}
        >
          {copied ? <Check className="size-3.5" /> : <Copy className="size-3.5" />}
          {copied ? t("cards.copied") : t("cards.copyUrl")}
        </Button>
      </div>
    </div>
  );
}
