import * as React from "react"

import { cn } from "@/lib/utils"

/*
 * 深蓝金输入框原子（DESIGN.md「缺失原子组件决策」）：
 * bg-surface-input 凹陷米白 + navy 文字 + royal focus ring。
 * 统一 RegisterSheet / SIM 领取表单 / 搜索框（T7-T11 收口）。
 * 吃 CSS 变量，换主题只改 :root。最小触摸区 h-11（44px，移动端）。
 */
function Input({ className, type, ...props }: React.ComponentProps<"input">) {
  return (
    <input
      type={type}
      data-slot="input"
      className={cn(
        "flex h-11 w-full rounded-md border border-surface-card-line bg-surface-input px-3 py-2 text-base text-text-on-light-primary",
        "placeholder:text-text-on-light-muted",
        "transition-colors outline-none",
        "focus-visible:border-brand-royal focus-visible:ring-2 focus-visible:ring-brand-royal/40",
        "disabled:cursor-not-allowed disabled:opacity-50",
        "aria-invalid:border-status-danger aria-invalid:ring-2 aria-invalid:ring-status-danger/30",
        className
      )}
      {...props}
    />
  )
}

export { Input }
