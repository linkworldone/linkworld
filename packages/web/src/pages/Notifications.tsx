import { useAccount } from "wagmi";
import { useTranslation } from "react-i18next";
import { useNotifications, useMarkAsRead, useMarkAllAsRead } from "@/hooks/useNotification";
import { timeAgo } from "@/utils/format";
import { EmptyState } from "@/components/shared/EmptyState";
import { Bell } from "lucide-react";
import type { Notification } from "@/types";

const borderColors: Record<Notification["type"], string> = {
  bill_due: "border-l-brand-gold",
  deposit_confirmed: "border-l-status-success",
  payment_confirmed: "border-l-status-success",
  service_suspended: "border-l-status-danger",
  system: "border-l-text-on-light-muted",
};

export default function Notifications() {
  const { t } = useTranslation();
  const { address } = useAccount();
  const { data: notifications } = useNotifications(address);
  const markAsRead = useMarkAsRead();
  const markAllAsRead = useMarkAllAsRead();

  const unread = notifications?.filter((n) => !n.read) || [];
  const read = notifications?.filter((n) => n.read) || [];

  const handleClick = (notif: Notification) => {
    if (!notif.read) {
      markAsRead.mutate(notif.id);
    }
  };

  return (
    <div className="px-4 space-y-4">
      {unread.length > 0 && (
        <div className="flex justify-end">
          <button
            onClick={() => address && markAllAsRead.mutate(address)}
            className="text-xs text-brand-royal font-semibold"
          >
            {t("notifications.markAllRead")}
          </button>
        </div>
      )}

      {!notifications?.length ? (
        <EmptyState icon={Bell} message={t("notifications.empty")} />
      ) : (
        <>
          {unread.length > 0 && (
            <div>
              <div className="text-[11px] text-text-muted uppercase tracking-wider mb-2">{t("notifications.sectionNew")}</div>
              <div className="space-y-2">
                {unread.map((notif) => (
                  <button
                    key={notif.id}
                    onClick={() => handleClick(notif)}
                    className={`w-full text-left p-3.5 bg-surface-card-elevated border border-surface-card-line shadow-card border-l-[3px] ${borderColors[notif.type]} rounded-xl relative`}
                  >
                    <div className="absolute top-3.5 right-3.5 w-2 h-2 rounded-full bg-brand-gold" />
                    <div className="text-[13px] font-semibold pr-4 text-text-on-light-primary">{notif.title}</div>
                    <div className="text-[11px] text-text-on-light-secondary mt-1 leading-relaxed">{notif.message}</div>
                    <div className="text-[10px] text-text-on-light-muted mt-1.5">{timeAgo(notif.createdAt)}</div>
                  </button>
                ))}
              </div>
            </div>
          )}

          {read.length > 0 && (
            <div>
              <div className="text-[11px] text-text-muted uppercase tracking-wider mb-2">{t("notifications.sectionEarlier")}</div>
              <div className="space-y-2">
                {read.map((notif) => (
                  <div
                    key={notif.id}
                    className={`p-3.5 bg-surface-card border border-surface-card-line border-l-[3px] ${borderColors[notif.type]} rounded-xl opacity-80`}
                  >
                    <div className="text-[13px] font-semibold text-text-on-light-secondary">{notif.title}</div>
                    <div className="text-[11px] text-text-on-light-muted mt-1 leading-relaxed">{notif.message}</div>
                    <div className="text-[10px] text-text-on-light-muted mt-1.5">{timeAgo(notif.createdAt)}</div>
                  </div>
                ))}
              </div>
            </div>
          )}
        </>
      )}
    </div>
  );
}
