import type { UserStatus } from "@/types";

const statusConfig: Record<UserStatus, { label: string; color: string; dot: string }> = {
  active: { label: "Active", color: "text-status-success", dot: "bg-status-success" },
  inactive: { label: "Inactive", color: "text-status-warning", dot: "bg-status-warning" },
  suspended: { label: "Suspended", color: "text-status-danger", dot: "bg-status-danger" },
};

export function StatusBadge({ status }: { status: UserStatus }) {
  const config = statusConfig[status];
  return (
    <div className="flex items-center gap-1.5">
      <span className={`w-2 h-2 rounded-full ${config.dot}`} />
      <span className={`text-sm font-semibold ${config.color}`}>{config.label}</span>
    </div>
  );
}
