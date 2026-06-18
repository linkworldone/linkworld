import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useAccount, useDisconnect } from "wagmi";
import { useQueryClient } from "@tanstack/react-query";
import { useTranslation } from "react-i18next";
import { Headset, LogOut, Mail, Send, ChevronRight } from "lucide-react";
import { LanguageToggle } from "@/components/shared/LanguageToggle";
import { BottomSheet } from "@/components/shared/BottomSheet";
import { Button } from "@/components/ui/button";
import { useUser } from "@/hooks/useUser";
import { emailApi } from "@/services/api";
import { WalletAuthRejectedError } from "@/services/api/signedPost";

type BindStep = "email" | "verify";

export default function Settings() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const { address } = useAccount();
  const { disconnect } = useDisconnect();
  const queryClient = useQueryClient();
  const { data: user } = useUser(address);

  const [supportOpen, setSupportOpen] = useState(false);

  // 绑定邮箱弹层两步流程状态
  const [bindOpen, setBindOpen] = useState(false);
  const [step, setStep] = useState<BindStep>("email");
  const [email, setEmail] = useState("");
  const [code, setCode] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const isValidEmail = /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email);

  const handleLogout = () => {
    disconnect();
    navigate("/");
  };

  const openBind = () => {
    setStep("email");
    setEmail(user?.email ?? "");
    setCode("");
    setError("");
    setBindOpen(true);
  };

  const handleSendCode = async () => {
    if (!address || !isValidEmail || loading) return;
    setLoading(true);
    setError("");
    try {
      await emailApi.bindEmail(address, email);
      setCode("");
      setStep("verify");
    } catch (e) {
      if (e instanceof WalletAuthRejectedError) {
        setError(t("errors.authCancelled"));
      } else if ((e as { status?: number }).status === 429) {
        setError(t("settings.rateLimited"));
      } else {
        setError((e as Error).message || t("settings.invalidEmail"));
      }
    } finally {
      setLoading(false);
    }
  };

  const handleVerify = async () => {
    if (!address || code.length !== 6 || loading) return;
    setLoading(true);
    setError("");
    try {
      await emailApi.verifyEmail(address, email, code);
      setBindOpen(false);
      // 刷新用户信息，使入口处邮箱更新
      queryClient.invalidateQueries({ queryKey: ["user"] });
    } catch (e) {
      if (e instanceof WalletAuthRejectedError) {
        setError(t("errors.authCancelled"));
      } else {
        setError((e as Error).message || t("settings.invalidCode"));
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="px-4 space-y-3">
      {/* 语言切换 */}
      <div className="flex items-center justify-between p-3.5 bg-surface-card border border-surface-card-line shadow-card rounded-xl">
        <span className="text-[14px] font-semibold text-text-on-light-primary">{t("settings.language")}</span>
        <LanguageToggle />
      </div>

      {/* 绑定邮箱 */}
      <button
        onClick={openBind}
        className="w-full flex items-center justify-between p-3.5 bg-surface-card border border-surface-card-line shadow-card rounded-xl text-left"
      >
        <span className="flex items-center gap-2.5">
          <Mail className="size-5 text-brand-gold" strokeWidth={1.75} />
          <span className="text-[14px] font-semibold text-text-on-light-primary">{t("settings.bindEmail")}</span>
        </span>
        <span className="flex items-center gap-1.5">
          <span className="text-[13px] text-text-on-light-secondary max-w-[150px] truncate">
            {user?.email || t("settings.notBound")}
          </span>
          <ChevronRight className="size-4 text-text-on-light-muted" />
        </span>
      </button>

      {/* 客服 */}
      <button
        onClick={() => setSupportOpen(true)}
        className="w-full flex items-center justify-between p-3.5 bg-surface-card border border-surface-card-line shadow-card rounded-xl text-left"
      >
        <span className="flex items-center gap-2.5">
          <Headset className="size-5 text-brand-gold" strokeWidth={1.75} />
          <span className="text-[14px] font-semibold text-text-on-light-primary">{t("settings.customerService")}</span>
        </span>
        <ChevronRight className="size-4 text-text-on-light-muted" />
      </button>

      {/* 登出 */}
      <button
        onClick={handleLogout}
        className="w-full flex items-center gap-2.5 p-3.5 bg-surface-card border border-surface-card-line shadow-card rounded-xl text-left"
      >
        <LogOut className="size-5 text-status-danger" strokeWidth={1.75} />
        <span className="text-[14px] font-semibold text-status-danger">{t("settings.logout")}</span>
      </button>

      {/* 绑定邮箱弹层 */}
      <BottomSheet open={bindOpen} onOpenChange={setBindOpen}>
        <div className="text-[15px] font-semibold text-text-on-light-primary mb-4">{t("settings.bindEmail")}</div>

        {step === "email" ? (
          <>
            <input
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder={t("settings.emailPlaceholder")}
              className="w-full px-4 py-3 bg-surface-input rounded-xl text-text-on-light-primary text-sm outline-none border border-surface-card-line focus:border-brand-royal mb-4"
            />
            {email && !isValidEmail && (
              <p className="text-xs text-status-danger mb-2">{t("settings.invalidEmail")}</p>
            )}
            {error && <p className="text-xs text-status-danger mb-2">{error}</p>}
            <Button onClick={handleSendCode} disabled={!isValidEmail || loading} className="w-full py-3">
              {loading ? t("settings.sending") : t("settings.sendCode")}
            </Button>
          </>
        ) : (
          <>
            <p className="text-sm text-text-on-light-secondary mb-3">{t("settings.codeSent")}</p>
            <input
              type="text"
              inputMode="numeric"
              value={code}
              onChange={(e) => setCode(e.target.value.replace(/\D/g, "").slice(0, 6))}
              placeholder={t("settings.codePlaceholder")}
              maxLength={6}
              className="w-full px-4 py-3 bg-surface-input rounded-xl text-text-on-light-primary text-sm outline-none border border-surface-card-line focus:border-brand-royal mb-4 text-center tracking-widest text-lg"
            />
            {error && <p className="text-xs text-status-danger mb-2">{error}</p>}
            <Button onClick={handleVerify} disabled={code.length !== 6 || loading} className="w-full py-3">
              {loading ? t("settings.binding") : t("settings.confirmBind")}
            </Button>
            <button
              onClick={() => {
                setStep("email");
                setCode("");
                setError("");
              }}
              disabled={loading}
              className="w-full mt-3 text-[13px] text-text-on-light-secondary disabled:opacity-50"
            >
              {t("settings.resendCode")}
            </button>
          </>
        )}
      </BottomSheet>

      <BottomSheet open={supportOpen} onOpenChange={setSupportOpen}>
        <div className="text-[15px] font-semibold text-text-on-light-primary mb-4">{t("settings.supportTitle")}</div>
        <div className="space-y-2">
          <a
            href="mailto:linkworldone@outlook.com"
            className="flex items-center gap-3 p-3 bg-surface-card-elevated border border-surface-card-line rounded-xl"
          >
            <Mail className="size-5 text-brand-gold shrink-0" strokeWidth={1.75} />
            <span className="flex flex-col">
              <span className="text-[13px] font-semibold text-text-on-light-primary">{t("settings.contactEmail")}</span>
              <span className="text-[12px] text-text-on-light-secondary">linkworldone@outlook.com</span>
            </span>
          </a>
          <a
            href="https://t.me/linkworldteam"
            target="_blank"
            rel="noopener noreferrer"
            className="flex items-center gap-3 p-3 bg-surface-card-elevated border border-surface-card-line rounded-xl"
          >
            <Send className="size-5 text-brand-gold shrink-0" strokeWidth={1.75} />
            <span className="flex flex-col">
              <span className="text-[13px] font-semibold text-text-on-light-primary">{t("settings.contactTelegram")}</span>
              <span className="text-[12px] text-text-on-light-secondary">@linkworldteam</span>
            </span>
          </a>
        </div>
      </BottomSheet>
    </div>
  );
}
