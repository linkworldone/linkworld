import { useTranslation } from "react-i18next";
import type { UserStatus } from "@/types";

const statusConfig: Record<UserStatus, { labelKey: string; color: string; dot: string }> = {
  active: { labelKey: "accountStatus.active", color: "text-status-success", dot: "bg-status-success" },
  inactive: { labelKey: "accountStatus.inactive", color: "text-status-warning", dot: "bg-status-warning" },
  suspended: { labelKey: "accountStatus.suspended", color: "text-status-danger", dot: "bg-status-danger" },
};

export function StatusBadge({ status }: { status: UserStatus }) {
  const { t } = useTranslation();
  const config = statusConfig[status];
  return (
    <div className="flex items-center gap-1.5">
      <span className={`w-2 h-2 rounded-full ${config.dot}`} />
      <span className={`text-sm font-semibold ${config.color}`}>{t(config.labelKey)}</span>
    </div>
  );
}
