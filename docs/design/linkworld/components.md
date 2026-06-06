# LinkWorld Web Components Inventory

## UI Components (shadcn/base-ui System)

| Component | File Path | Purpose | Key Props/Variants |
|-----------|-----------|---------|-------------------|
| Button | `/src/components/ui/button.tsx` | Base interactive button component | `variant: "default" \| "outline" \| "secondary" \| "ghost" \| "destructive" \| "link"` <br> `size: "default" \| "xs" \| "sm" \| "lg" \| "icon" \| "icon-xs" \| "icon-sm" \| "icon-lg"` |
| Badge | `/src/components/ui/badge.tsx` | Status/tag badge component | `variant: "default" \| "secondary" \| "destructive" \| "outline" \| "ghost" \| "link"` |
| Tabs | `/src/components/ui/tabs.tsx` | Tabbed navigation component | `Tabs.Root`, `Tabs.List` (variant: "default" \| "line"), `Tabs.Trigger`, `Tabs.Content` |

**Base-UI Integration**: Components use `@base-ui/react` primitives (Button, Tabs) with class-variance-authority styling. Core utilities: `/src/lib/utils.ts` exports `cn()` for merging Tailwind classes.

---

## Layout Components

| Component | File Path | Purpose | Key Props/Dependencies |
|-----------|-----------|---------|----------------------|
| AppLayout | `/src/components/layout/AppLayout.tsx` | Main app wrapper with auth guard | Uses `useAccount()` (wagmi), `useUser()` hook; protects routes based on `user.status` (active/inactive/suspended) |
| Header | `/src/components/layout/Header.tsx` | Top navigation bar | Conditional render: dashboard welcome vs. page title; shows notification bell with unread count |
| TabBar | `/src/components/layout/TabBar.tsx` | Bottom fixed navigation | 5 tabs: Home, Services, Deposit, Bills (with badge), Cards; uses emoji icons; active state = `text-brand-blue` |

---

## Business Components (Shared)

| Component | File Path | Purpose | Key Props |
|-----------|-----------|---------|-----------|
| GuardCard | `/src/components/shared/GuardCard.tsx` | Access restriction screen | `icon: string`, `title: string`, `message: string`, `actionLabel: string`, `actionPath: string` |
| StatusBadge | `/src/components/shared/StatusBadge.tsx` | User status indicator | `status: "active" \| "inactive" \| "suspended"` → renders colored dot + label |
| AmountDisplay | `/src/components/shared/AmountDisplay.tsx` | Formatted amount display | `amount: bigint \| string`, `currency?: string`, `size: "sm" \| "md" \| "lg"`, `colorClass?: string` |
| BottomSheet | `/src/components/shared/BottomSheet.tsx` | Drawer modal (vaul) | `open: boolean`, `onOpenChange: (open: boolean) => void`, `children: ReactNode` |
| EmptyState | `/src/components/shared/EmptyState.tsx` | No data state placeholder | TBD - not yet read |

---

## Wallet Components

| Component | File Path | Purpose | Key Props |
|-----------|-----------|---------|-----------|
| ConnectButton | `/src/components/wallet/ConnectButton.tsx` | RainbowKit wallet connect button | `label?: string` (default: "Connect Wallet"); returns null if already connected |
| RegisterSheet | `/src/components/wallet/RegisterSheet.tsx` | Account registration flow (2-step) | `address: string`, `open: boolean`, `onClose: () => void`, `onSuccess: () => void` <br> Steps: email → verification code → contract mint + backend sync |

---

## Summary
- **Total UI Components**: 3 (button, badge, tabs) - all shadcn/base-ui wraps
- **Total Layout**: 3 (AppLayout, Header, TabBar)
- **Total Shared**: 5 (GuardCard, StatusBadge, AmountDisplay, BottomSheet, EmptyState)
- **Total Wallet**: 2 (ConnectButton, RegisterSheet)
- **Total Across Project**: 13 components

**Styling**: All components use Tailwind 3.4 + custom colors (brand-blue, brand-purple, status-success/warning/danger, text-primary/secondary/muted)
