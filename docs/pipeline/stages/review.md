# Stage: review -- Code Review Report

> **Status**: completed | **Date**: 2026-04-12 | **Gate**: PASS

## Summary

The LinkWorld frontend implementation is well-structured and closely follows the design spec. The three-layer architecture (Pages -> Hooks -> Services) is consistently applied. Code is clean with zero type safety violations, no dead code, and consistent patterns across all files. Four non-blocking issues were identified.

## Issues Found

| # | Severity | File | Issue | Recommendation |
|---|----------|------|-------|---------------|
| 1 | Important | `services/mock/data.ts` | `mockUser` is exported by value (`let` export). The re-export in `data.ts` line 73 exports the initial `null` binding; consumers importing `mockUser` directly get a stale reference after `setMockUser()` mutates it. Only `userService.getUserProfile()` works correctly because it reads the module-level variable at call time. | Replace direct export of `mockUser` with a getter function `getMockUser()`, or always access via service. Currently not a runtime bug because all consumers go through services, but it is a latent trap for future developers. |
| 2 | Important | `pages/Deposit.tsx` (line 28) | `BigInt(Math.floor(parseFloat(amount) * 1e18))` has floating-point precision loss for amounts with many decimals (e.g. entering "0.123456789012345678" will silently truncate). | Use a string-based decimal-to-wei conversion (split on ".", pad fractional part to 18 chars, concatenate, parse as BigInt). Acceptable for mock phase but should be fixed before real contract integration. |
| 3 | Important | `components/layout/AppLayout.tsx` | When user is `null` (UNREGISTERED but wallet is connected and query finished loading), no redirect or guard is triggered. The Outlet renders with undefined user. Pages handle this gracefully because they use optional chaining, but the spec says UNREGISTERED should only access Landing. | Add a guard: if `!isLoading && !user`, redirect to `/` (Landing) so the registration flow triggers. |
| 4 | Important | `services/mock/billingService.ts` | `payBill` does not update user status. When all unpaid bills are paid, the user should transition from SUSPENDED back to ACTIVE per the state machine. Mock service should call `setMockUser({...mockUser, status: 'active'})` after payment clears all unpaid bills. | Add status transition logic to `payBill` or document this as a known mock limitation. |
| 5 | Suggestion | `pages/Dashboard.tsx` | Missing the real-time usage detail section specified in spec section 5.2 ("Data Usage / Call Log with mock timed increments"). Current implementation shows the summary estimate but not the detailed log. | Add a simple usage log list with mock data. Low priority for Round 1. |
| 6 | Suggestion | `components/wallet/ConnectButton.tsx` | When wallet is connected (`account` exists), the component returns `null`. This means on Landing page, after connecting, the Connect Wallet button disappears but there is no disconnect option or address display in its place until redirect happens. | Consider showing a minimal connected state (abbreviated address) or use RainbowKit's built-in connected display for the Landing header. |
| 7 | Suggestion | `pages/Landing.tsx` | The "Learn More" button has no handler -- it is purely decorative. | Add a scroll-to or link behavior, or add `disabled` styling to signal it is not yet functional. |
| 8 | Suggestion | `hooks/useNotification.ts` | `refetchInterval: 30_000` on `useUnreadCount` is good, but `useNotifications` has no refetch interval. If a user stays on the Notifications page, new notifications will not appear until manual refresh. | Add `refetchInterval` to `useNotifications` as well, or rely on invalidation from other mutations. |
| 9 | Suggestion | Multiple Drawer usages | All Drawer implementations repeat the same boilerplate (Overlay + Content + drag handle). | Extract a `BottomSheet` wrapper component to reduce duplication. Not blocking. |

## Spec Coverage

| Spec Section | Covered? | Notes |
|-------------|----------|-------|
| 5.1 Landing Page | Yes | All elements present: logo, hero, stats, Powered by, ConnectButton, RegisterSheet |
| 5.2 Dashboard | Mostly | Status card (2x2), usage card (3-col), quick actions (2x2), suspended warning all present. Missing: detailed usage log with timed increments |
| 5.3 Deposit | Yes | Balance card with progress bar, deposit/withdraw buttons, history list, currency selector sheet |
| 5.4 Services + RegionDetail | Yes | Regions/My Numbers tabs, search, country list, operator cards, apply number sheet, eSIM credentials display |
| 5.5 Billing + BillDetail | Yes | Unpaid/History tabs, alert banner, fee breakdown, pay sheet, bill detail with usage |
| 5.6 Notifications | Yes | New/Earlier groups, colored left borders, unread blue bg + dot, mark all read |
| 6.1 AppLayout | Yes | max-w-mobile, Header + Outlet + TabBar, route guards for INACTIVE and SUSPENDED |
| 6.2 TabBar | Yes | 5 tabs, active highlight, badges for unpaid bills and unread notifications |
| 6.3 Bottom Sheet | Yes | Using vaul Drawer with overlay, drag handle, full-width confirm button |
| 6.4 GuardCard | Yes | Icon + title + message + action button |
| 7.x Types | Yes | All types match spec. Bill uses string amounts (arch-review risk #10 addressed). VirtualNumber has credentials field. |
| 8.x Services | Yes | All method signatures match spec exactly. |
| 9.x Design tokens | Yes | Semantic color tokens in tailwind.config.ts match spec colors. max-w-mobile = 430px. |

## Arch-Review Risk Items

| # | Risk | Addressed? | How |
|---|------|-----------|-----|
| 1 | Multi-virtual-number scenario | Yes | Services page has My Numbers tab showing all numbers with credentials |
| 2 | Early billing entry missing | Yes | Billing accessible from Dashboard quick actions and TabBar |
| 3 | Overdue deduction countdown warning | Yes | Dashboard shows suspended warning with 14-day message |
| 4 | Tailwind semantic color tokens | Yes | tailwind.config.ts defines surface/brand/status/text/border token groups |
| 5 | BrowserRouter position in main.tsx | Yes | BrowserRouter is inside RainbowKitProvider, matching spec Provider order |
| 6 | shadcn/ui initialization | Yes | components.json exists, Button/Badge/Tabs installed, lib/utils.ts present |
| 7 | Desktop max-w container | Yes | AppLayout uses max-w-mobile mx-auto; TabBar and Drawers also constrained |
| 8 | Sheet drag gesture (vaul) | Yes | All sheets use vaul Drawer with proper overlay and drag handle |
| 9 | Mock service delay simulation | Yes | delay.ts helper used in all services with configurable timing (600ms default, 1200ms for mutations) |
| 10 | Bill amount precision | Yes | Bill type uses string for operatorFee/platformFee/totalAmount instead of number |

## Positive Observations

- **Excellent type safety**: Zero `any` types, zero ts-ignore directives, zero console.log statements across the entire codebase. All hooks properly type their query/mutation parameters.
- **Consistent architecture**: Every page follows the same pattern -- wagmi `useAccount()` for address, custom hooks for data, local state for UI. No exceptions or shortcuts.
- **Smart mock design**: Services return shallow copies (`{...mockDeposit}`, `[...mockBills]`) preventing accidental mutation of shared state. Mutation functions properly use setter helpers.
- **Proper React Query usage**: Mutations invalidate related queries correctly. `enabled` flag prevents queries from firing without required parameters. `staleTime` and `refetchInterval` are configured appropriately.
- **Mobile-first design tokens**: Semantic color system in Tailwind config perfectly mirrors the spec's design table. max-w-mobile constraint consistently applied.
- **Lazy loading**: All pages use React.lazy with a clean Suspense fallback, good for bundle splitting.
- **Clean service abstraction**: `services/index.ts` re-exports all services, making future mock-to-real swap straightforward as planned.
- **Guard implementation**: AppLayout correctly implements the status-based route guard logic from the spec with clear path matching.

## Recommendation

**Ship it.** The 4 important issues are real but non-blocking for the mock/demo phase. Issues #3 (UNREGISTERED guard) and #4 (payBill status transition) should be addressed before any user testing. Issues #1 and #2 are latent traps that should be fixed before real contract integration.
