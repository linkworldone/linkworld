import { ConnectButton as RainbowConnectButton } from "@rainbow-me/rainbowkit";

export function ConnectButton({ label = "Connect Wallet" }: { label?: string }) {
  return (
    <RainbowConnectButton.Custom>
      {({ openConnectModal, account, mounted }) => {
        if (!mounted) return null;
        if (account) return null;
        return (
          <button
            onClick={openConnectModal}
            className="bg-brand-blue text-white px-4 py-2 rounded-lg text-sm font-semibold hover:opacity-90 transition-opacity"
          >
            {label}
          </button>
        );
      }}
    </RainbowConnectButton.Custom>
  );
}
