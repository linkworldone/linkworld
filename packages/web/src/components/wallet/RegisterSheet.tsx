import { useEffect } from "react";

import { Drawer } from "vaul";

import { useTranslation } from "react-i18next";

import { Button } from "@/components/ui/button";

import { useRegister } from "@/hooks/useUser";

interface RegisterSheetProps {
  address: string;
  open: boolean;
  onClose: () => void;
  onSuccess: () => void;
}

export function RegisterSheet({
  address,
  open,
  onClose,
  onSuccess,
}: RegisterSheetProps) {
  const { t } = useTranslation();

  const { register, backendSync, isContractPending } = useRegister();

  useEffect(() => {
    if (backendSync.isSuccess) {
      onSuccess();
    }
  }, [backendSync.isSuccess, onSuccess]);

  const handleRegister = () => {
    // Temporary compatibility:
    // the legacy contract/backend still have an email field,
    // but LinkWorld no longer asks users for one.
    register("");
  };

  return (
    <Drawer.Root
      open={open}
      onOpenChange={(o) => !o && onClose()}
    >
      <Drawer.Portal>
        <Drawer.Overlay className="fixed inset-0 bg-black/60 z-50" />

        <Drawer.Content className="fixed bottom-0 left-0 right-0 max-w-mobile mx-auto bg-surface-card rounded-t-2xl z-50 p-6">
          <div className="w-12 h-1 bg-surface-input rounded-full mx-auto mb-6" />

          <h2 className="text-lg font-bold mb-1 text-text-on-light-primary">
            Connect your wallet
          </h2>

          <p className="text-sm text-text-on-light-secondary mb-6">
            Your wallet is your LinkWorld identity. No email or password is required.
          </p>

          <div className="mb-4 p-4 rounded-xl bg-surface-input">
            <div className="text-xs text-text-on-light-muted mb-1">
              Wallet
            </div>

            <div className="text-sm font-mono text-text-on-light-primary break-all">
              {address}
            </div>
          </div>

          <Button
            onClick={handleRegister}
            disabled={isContractPending || backendSync.isPending}
            className="w-full py-3"
          >
            {isContractPending || backendSync.isPending
              ? "Creating your identity..."
              : "Continue with Wallet"}
          </Button>
        </Drawer.Content>
      </Drawer.Portal>
    </Drawer.Root>
  );
}
