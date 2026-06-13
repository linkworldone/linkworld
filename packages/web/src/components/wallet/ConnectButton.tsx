import { ConnectButton as RainbowConnectButton } from "@rainbow-me/rainbowkit";
import { useTranslation } from "react-i18next";

export function ConnectButton({ label }: { label?: string }) {
  const { t } = useTranslation();
  const buttonLabel = label ?? t("common.connectWallet");
  return (
    <RainbowConnectButton.Custom>
      {({ openConnectModal, account, mounted }) => {
        if (!mounted) return null;
        if (account) return null;
        return (
          <button
            onClick={openConnectModal}
            className="bg-brand-gold text-brand-navy px-4 py-2 rounded-lg text-sm font-semibold hover:bg-brand-gold-hover transition-colors min-h-[44px]"
          >
            {buttonLabel}
          </button>
        );
      }}
    </RainbowConnectButton.Custom>
  );
}
