import type { LucideIcon } from "lucide-react";

/*
 * 空态原子（DESIGN.md D5）：icon prop 从 emoji string 改 LucideIcon 组件类型。
 * 调用方传 lucide 图标组件（如 CreditCard）；深蓝金下统一 gold 描金图标 + navy 文案。
 */
export function EmptyState({ icon: Icon, message }: { icon: LucideIcon; message: string }) {
  return (
    <div className="flex flex-col items-center justify-center py-16 text-center">
      <Icon className="mb-3 size-10 text-brand-gold" strokeWidth={1.5} />
      <p className="text-sm text-text-on-light-secondary">{message}</p>
    </div>
  );
}
