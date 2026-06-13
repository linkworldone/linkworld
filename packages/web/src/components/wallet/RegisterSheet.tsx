import { useState, useEffect } from "react";
import { Drawer } from "vaul";
import { useTranslation } from "react-i18next";
import { Button } from "@/components/ui/button";
import { useRegister, useSendVerificationCode, useVerifyEmail } from "@/hooks/useUser";

interface RegisterSheetProps {
  address: string;
  open: boolean;
  onClose: () => void;
  onSuccess: () => void;
}

export function RegisterSheet({ address, open, onClose, onSuccess }: RegisterSheetProps) {
  const { t } = useTranslation();
  const [step, setStep] = useState<"email" | "verify">("email");
  const [email, setEmail] = useState("");
  const [code, setCode] = useState("");

  const isValidEmail = /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email);

  const sendCode = useSendVerificationCode();
  const verifyEmail = useVerifyEmail();
  const { register, backendSync, isContractPending } = useRegister();

  // 合约注册成功 → useRegister 内部 useEffect 自动调 backendSync
  // 这里只处理页面跳转
  useEffect(() => {
    if (backendSync.isSuccess) {
      onSuccess();
    }
  }, [backendSync.isSuccess]);

  const handleSendCode = async () => {
    await sendCode.mutateAsync({ address, email });
    setStep("verify");
  };

  const handleVerify = async () => {
    const verified = await verifyEmail.mutateAsync({ address, code });
    if (verified) {
      register(email);
    }
  };

  return (
    <Drawer.Root open={open} onOpenChange={(o) => !o && onClose()}>
      <Drawer.Portal>
        <Drawer.Overlay className="fixed inset-0 bg-black/60 z-50" />
        <Drawer.Content className="fixed bottom-0 left-0 right-0 max-w-mobile mx-auto bg-surface-card rounded-t-2xl z-50 p-6">
          <div className="w-12 h-1 bg-surface-input rounded-full mx-auto mb-6" />
          <h2 className="text-lg font-bold mb-1 text-text-on-light-primary">{t("register.title")}</h2>
          <p className="text-sm text-text-on-light-secondary mb-6">
            {step === "email" ? t("register.emailHint") : t("register.codeHint")}
          </p>

          {step === "email" ? (
            <>
              <input
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder={t("register.emailPlaceholder")}
                className="w-full px-4 py-3 bg-surface-input rounded-xl text-text-on-light-primary text-sm outline-none border border-surface-card-line focus:border-brand-royal mb-4"
              />
              {email && !isValidEmail && (
                <p className="text-xs text-status-danger mb-2">{t("register.invalidEmail")}</p>
              )}
              <Button onClick={handleSendCode} disabled={!isValidEmail || sendCode.isPending} className="w-full py-3">
                {sendCode.isPending ? t("register.sending") : t("register.sendCode")}
              </Button>
            </>
          ) : (
            <>
              <input
                type="text"
                value={code}
                onChange={(e) => setCode(e.target.value)}
                placeholder={t("register.codePlaceholder")}
                maxLength={6}
                className="w-full px-4 py-3 bg-surface-input rounded-xl text-text-on-light-primary text-sm outline-none border border-surface-card-line focus:border-brand-royal mb-4 text-center tracking-widest text-lg"
              />
              <Button onClick={handleVerify} disabled={code.length < 6 || verifyEmail.isPending || isContractPending} className="w-full py-3">
                {verifyEmail.isPending || isContractPending ? t("register.verifying") : t("register.verifyRegister")}
              </Button>
            </>
          )}
        </Drawer.Content>
      </Drawer.Portal>
    </Drawer.Root>
  );
}
