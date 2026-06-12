import * as React from "react"

import { cn } from "@/lib/utils"

/*
 * 暖米白卡片原子（DESIGN.md「缺失原子组件决策」）：
 * bg-surface-card 暖米白 + 1px 金线描边 + radius-lg + 轻投影。
 * 一处定义，全站手写卡（Dashboard/Deposit/Services/Billing/Cards）收口复用。
 * 卡内默认深色化：text-on-light-primary（navy，对比 ≈12:1）。
 * 吃 CSS 变量，换主题只改 :root。
 */
function Card({ className, ...props }: React.ComponentProps<"div">) {
  return (
    <div
      data-slot="card"
      className={cn(
        "rounded-lg border border-surface-card-line bg-surface-card text-text-on-light-primary shadow-card",
        className
      )}
      {...props}
    />
  )
}

function CardHeader({ className, ...props }: React.ComponentProps<"div">) {
  return (
    <div
      data-slot="card-header"
      className={cn("flex flex-col gap-1 px-4 pt-4", className)}
      {...props}
    />
  )
}

function CardTitle({ className, ...props }: React.ComponentProps<"div">) {
  return (
    <div
      data-slot="card-title"
      className={cn("font-display text-base font-semibold text-text-on-light-primary", className)}
      {...props}
    />
  )
}

function CardDescription({ className, ...props }: React.ComponentProps<"div">) {
  return (
    <div
      data-slot="card-description"
      className={cn("text-sm text-text-on-light-secondary", className)}
      {...props}
    />
  )
}

function CardContent({ className, ...props }: React.ComponentProps<"div">) {
  return (
    <div
      data-slot="card-content"
      className={cn("px-4 py-4", className)}
      {...props}
    />
  )
}

function CardFooter({ className, ...props }: React.ComponentProps<"div">) {
  return (
    <div
      data-slot="card-footer"
      className={cn("flex items-center px-4 pb-4", className)}
      {...props}
    />
  )
}

export { Card, CardHeader, CardTitle, CardDescription, CardContent, CardFooter }
