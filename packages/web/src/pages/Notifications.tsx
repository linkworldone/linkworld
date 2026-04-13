import { useAccount } from "wagmi";
import { useNotifications, useMarkAsRead, useMarkAllAsRead } from "@/hooks/useNotification";
import { timeAgo } from "@/utils/format";
import { EmptyState } from "@/components/shared/EmptyState";
import type { Notification } from "@/types";

const borderColors: Record<Notification["type"], string> = {
  bill_due: "border-l-brand-blue",
  deposit_confirmed: "border-l-status-success",
  payment_confirmed: "border-l-status-success",
  service_suspended: "border-l-status-danger",
  system: "border-l-text-muted",
};

export default function Notifications() {
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
      <div className="flex justify-between items-center">
        <h2 className="text-[17px] font-bold">Notifications</h2>
        {unread.length > 0 && (
          <button
            onClick={() => address && markAllAsRead.mutate(address)}
            className="text-xs text-brand-blue font-semibold"
          >
            Mark all read
          </button>
        )}
      </div>

      {!notifications?.length ? (
        <EmptyState icon="🔔" message="No notifications yet." />
      ) : (
        <>
          {unread.length > 0 && (
            <div>
              <div className="text-[11px] text-text-muted uppercase tracking-wider mb-2">New</div>
              <div className="space-y-2">
                {unread.map((notif) => (
                  <button
                    key={notif.id}
                    onClick={() => handleClick(notif)}
                    className={`w-full text-left p-3.5 bg-brand-blue/[0.08] border-l-[3px] ${borderColors[notif.type]} rounded-xl relative`}
                  >
                    <div className="absolute top-3.5 right-3.5 w-2 h-2 rounded-full bg-brand-blue" />
                    <div className="text-[13px] font-semibold pr-4">{notif.title}</div>
                    <div className="text-[11px] text-text-secondary mt-1 leading-relaxed">{notif.message}</div>
                    <div className="text-[10px] text-text-muted mt-1.5">{timeAgo(notif.createdAt)}</div>
                  </button>
                ))}
              </div>
            </div>
          )}

          {read.length > 0 && (
            <div>
              <div className="text-[11px] text-text-muted uppercase tracking-wider mb-2">Earlier</div>
              <div className="space-y-2">
                {read.map((notif) => (
                  <div
                    key={notif.id}
                    className={`p-3.5 bg-surface-card border-l-[3px] ${borderColors[notif.type]} rounded-xl`}
                  >
                    <div className="text-[13px] font-semibold text-text-secondary">{notif.title}</div>
                    <div className="text-[11px] text-text-muted mt-1 leading-relaxed">{notif.message}</div>
                    <div className="text-[10px] text-text-muted mt-1.5">{timeAgo(notif.createdAt)}</div>
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
