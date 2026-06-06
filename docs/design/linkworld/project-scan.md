# LinkWorld Web Project Scan Report

## Technology Stack Summary

| Technology | Version | Purpose |
|-----------|---------|---------|
| **React** | 19.0.0 | UI framework |
| **Vite** | 6.0.0 | Build tool & dev server |
| **TypeScript** | 5.6.0 | Type safety |
| **Tailwind CSS** | 3.4.17 | Utility-first styling |
| **shadcn/base-ui** | Button, Tabs from @base-ui/react 1.3.0 | Unstyled accessible components |
| **React Router** | 7.1.0 | Client-side routing |
| **React Query** | 5.62.0 | Server state management (useQuery/useMutation) |
| **wagmi** | 2.14.0 | Ethereum wallet hooks |
| **Viem** | 2.21.0 | Ethereum client library |
| **RainbowKit** | 2.2.0 | Wallet connect UI |
| **Sonner** | 2.0.7 | Toast notifications |
| **Lucide React** | 1.8.0 | Icon library |
| **vaul** | 1.1.2 | Drawer/sheet component |
| **Axios** | 1.15.0 | HTTP client (API requests) |
| **class-variance-authority** | 0.7.1 | Component variant system |
| **clsx** | 2.1.1 | Classname utility |
| **tailwind-merge** | 3.5.0 | Merge Tailwind classes safely |

---

## Directory Structure & Responsibilities

```
packages/web/src/
├── components/
│   ├── ui/              (3 files)       - Shadcn/base-ui primitive components (Button, Badge, Tabs)
│   ├── layout/          (3 files)       - App layout: AppLayout, Header, TabBar
│   ├── shared/          (5 files)       - Reusable business components (GuardCard, StatusBadge, AmountDisplay, BottomSheet, EmptyState)
│   └── wallet/          (2 files)       - Web3 integration: ConnectButton, RegisterSheet
├── hooks/
│   ├── contracts/       (5 files)       - Contract interaction: useUserRegistry, useDepositContract, usePaymentContract, useServiceManager, useTrafficCard
│   ├── useUser.ts                       - User auth & registration flow
│   ├── useNotification.ts               - Notification state (mock)
│   ├── useBilling.ts                    - Bill payment & usage
│   ├── useDeposit.ts                    - Deposit balance & history
│   ├── useOperator.ts                   - Regions, operators, virtual numbers
│   ├── useTransactionFlow.ts            - TX state management & error parsing
│   └── index.ts                         - Exports all
├── services/
│   ├── api/
│   │   ├── client.ts                    - Axios client + response interceptor
│   │   ├── userApi.ts                   - User registration & profile
│   │   ├── operatorApi.ts               - Regions & operators
│   │   ├── depositApi.ts                - Deposit/withdraw history
│   │   ├── billingApi.ts                - Bills & payment records
│   │   ├── usageApi.ts                  - Monthly usage estimate
│   │   └── index.ts                     - Exports all real APIs
│   ├── mock/
│   │   ├── notificationService.ts       - Mock notifications (backend not implemented)
│   │   ├── data.ts                      - Mock data fixtures
│   │   ├── delay.ts                     - Delay helper for mock latency
│   │   ├── userService.ts               - Mock user data
│   │   ├── billingService.ts            - Mock billing data
│   │   ├── depositService.ts            - Mock deposit data
│   │   ├── operatorService.ts           - Mock operator data
│   │   └── notificationService.ts       - Mock notifications
│   └── index.ts                         - Exports APIs (real) + notification mock
├── pages/                   (9 files)   - Routed page components
│   ├── Landing.tsx                      - Public landing page (connect wallet, register)
│   ├── Dashboard.tsx                    - Main dashboard (overview, quick actions)
│   ├── Deposit.tsx                      - Deposit/withdraw ETH
│   ├── Services.tsx                     - Virtual number service catalog
│   ├── RegionDetail.tsx                 - Regional operators & pricing
│   ├── Billing.tsx                      - Bill list & payment
│   ├── BillDetail.tsx                   - Single bill detail
│   ├── Notifications.tsx                - Notification inbox
│   └── Cards.tsx                        - Traffic card NFT management
├── config/
│   ├── abis/                (8 files)   - Contract ABIs (UserRegistry, Deposit, Payment, ServiceManager, TrafficCard, FeeManager, Oracle)
│   ├── api.ts                           - API base URL (from env)
│   ├── constants.ts                     - App constants
│   ├── contracts.ts                     - Contract address resolver (by chainId)
│   ├── wagmi.ts                         - wagmi config + RainbowKit setup
│   └── chains.ts                        - Network configs (hardhatLocal, zgTestnet)
├── lib/
│   └── utils.ts                         - `cn()` utility (clsx + tailwind-merge)
├── utils/
│   ├── format.ts                        - Format functions (shortenAddress, parseUnits, formatAmount, formatDate, timeAgo)
│   └── pendingSync.ts                   - Offline-first sync retry logic
├── types/
│   └── index.ts                         - TypeScript type definitions (User, Bill, DepositInfo, VirtualNumber, etc.)
├── App.tsx                              - Route definitions + page lazy loading
├── main.tsx                             - React entry point + provider setup
└── index.css                            - CSS variables + Tailwind imports
```

**Summary**:
- **Components**: 13 total (3 UI + 3 layout + 5 shared + 2 wallet)
- **Pages**: 9 routed pages
- **Hooks**: 12 custom + 5 contract-specific
- **Services**: 5 real APIs + 1 mock (notifications)
- **Config**: Contract ABIs, chains, wagmi setup

---

## Pages & Routes (from App.tsx)

| Route | Page Component | Status | Purpose |
|-------|---|--------|---------|
| `/` | Landing | ✓ Complete | Public landing, wallet connect, registration modal |
| `/dashboard` | Dashboard | ✓ Complete | Overview: account status, deposit balance, usage, quick actions |
| `/deposit` | Deposit | ✓ Complete | Deposit/withdraw ETH with TX state flow |
| `/services` | Services | ✓ Complete | Browse regions + operators |
| `/services/:regionCode` | RegionDetail | ✓ Complete | Region detail + operator list + pricing |
| `/billing` | Billing | ✓ Complete | Unpaid/paid bill list + payment flow |
| `/billing/:billId` | BillDetail | ✓ Complete | Single bill detail (payment modal from parent) |
| `/notifications` | Notifications | ✓ Complete | Notification inbox (mock) |
| `/cards` | Cards | ✓ Complete | Traffic card NFT management (to be implemented) |

**Route Protection**: AppLayout enforces auth:
- Must be connected (wagmi)
- Must be registered (useUser returns non-null)
- Inactive users blocked from: Services, Billing, Notifications
- Suspended users blocked from: Services, Deposit

---

## Services Layer Status

### Real Backend APIs (5 services)

| API | Endpoints | Status | Notes |
|-----|-----------|--------|-------|
| **User API** | POST `/api/register`, GET `/api/user/{wallet}` | ✓ Real | User registration & profile fetch |
| **Operator API** | GET `/api/operators`, GET `/api/countries` | ✓ Real | Regions & operator catalog |
| **Deposit API** | GET `/api/deposit/{wallet}`, POST `/api/deposit`, POST `/api/withdraw` | ✓ Real | Deposit history & transaction records |
| **Billing API** | GET `/api/bills?wallet=...&filter=...`, POST `/api/bill/{id}/pay` | ✓ Real | Bill list & payment records |
| **Usage API** | GET `/api/usage/{wallet}` | ✓ Real | Monthly usage estimate |

### Mock Services (1 service)

| Service | Status | Reason |
|---------|--------|--------|
| **Notification Service** | 🟡 Mock | Backend API not yet implemented; localStorage + delay-based mock |

### Dual-Write Pattern (Contract + Backend)

Operations that sync blockchain + backend:
1. **Registration**: Contract mint UserRegistry NFT → Backend store profile
2. **Deposit**: Contract deposit() → Backend record transaction
3. **Withdrawal**: Contract withdraw() → Backend record transaction
4. **Bill Payment**: Contract payBill() → Backend mark as paid
5. **Service Activation**: Contract activateService() → Backend archive service

**Offline Resilience**: Uses `pendingSync` (localStorage + retry-with-backoff) for failed backend syncs.

---

## Web3 Integration

### Wallet & Chain Configuration

- **Wallet Connectors**: RainbowKit (supports MetaMask, WalletConnect, Coinbase, etc.)
- **Primary Chains**: 
  - `zgTestnet` (production target, 0G Chain testnet)
  - `hardhatLocal` (local dev, chainId: 31337)
- **Config File**: `/src/config/wagmi.ts` + `/src/config/chains.ts`

### Contract Interactions

| Contract | Module | Write Functions | Read Functions |
|----------|--------|-----------------|-----------------|
| **UserRegistry** | `useUserRegistry.ts` | `register(email)` | `isRegistered(address)` |
| **Deposit** | `useDepositContract.ts` | `deposit()`, `withdraw(amount)` | `balanceOf(address)` |
| **Payment** | `usePaymentContract.ts` | `payBill(billId, amount)` | - |
| **ServiceManager** | `useServiceManager.ts` | `activateService(operatorId, serviceId)` | - |
| **TrafficCard** | `useTrafficCard.ts` | - | - (TBD) |
| **FeeManager** | config/abis | - | - (oracle-only) |
| **Oracle** | config/abis | - | - (oracle-only) |

### ABIs Location

All contract ABIs stored in `/src/config/abis/`:
- `UserRegistry.ts` (register, isRegistered)
- `Deposit.ts` (deposit, withdraw, balanceOf)
- `Payment.ts` (payBill, payWithCard)
- `ServiceManager.ts` (activateService, deactivateService)
- `TrafficCardNFT.ts` (mint, burn, balanceOf)
- `FeeManager.ts` (oracle only)
- `Oracle.ts` (price feed)

**Contract Addresses**: Resolved via `/src/config/contracts.ts` by `chainId`

---

## Current Web Completion Status

### Fully Implemented Features ✓

1. **Authentication**
   - Wallet connect (RainbowKit)
   - User registration (email → verification code → contract mint)
   - Session persistence (localStorage + useUser hook)

2. **Core Pages**
   - Landing (public, wallet connect, register modal)
   - Dashboard (overview: account status, deposit, usage, quick actions)
   - Deposit (deposit/withdraw ETH with TX state)
   - Services (region browser)
   - Region Detail (operators + pricing)
   - Billing (bill list, payment modal, status filter)
   - Notifications (mock inbox)

3. **State Management**
   - React Query (useQuery/useMutation)
   - wagmi (wallet + contract interaction)
   - Offline-first sync (pending sync for dual-write)

4. **UI System**
   - 13 components (ui, layout, shared, wallet)
   - Tailwind 3.4 + custom colors (dark theme: #0a0a14 bg, #3b82f6 primary)
   - shadcn/base-ui integration
   - Responsive mobile-first (max-width: 430px)

### Partially Implemented Features 🟡

1. **Cards/Traffic Card**
   - Page exists but functionality TBD
   - ABIs defined but hooks not fully implemented

2. **Email Verification**
   - Currently mock (always passes)
   - Backend not sending real emails yet

3. **Notifications**
   - Using mock service (localStorage)
   - Backend notification API not yet live

### Not Yet Implemented ❌

1. **Bill Detail Page** - Detail fetch endpoint not on backend (returns null)
2. **Service Suspension Flow** - Guard messages show but no appeal/reinstatement
3. **Error Recovery UI** - Generic error pages for API failures
4. **Pagination** - All lists load full results

---

## Reusability Assessment for Reconstruction

### High-Value Reusable Components (Easy to Reuse in Refactor)

| Component | Reuse Level | Notes |
|-----------|-----------|-------|
| Button, Badge, Tabs | ⭐⭐⭐⭐⭐ | Pure UI, no dependencies, can copy-paste to new design |
| StatusBadge | ⭐⭐⭐⭐ | Self-contained, just update colors/styling |
| AmountDisplay | ⭐⭐⭐⭐ | Formatting logic reusable, just reskin |
| BottomSheet | ⭐⭐⭐⭐ | vaul-based, no app-specific logic |
| GuardCard | ⭐⭐⭐ | Layout reusable, update colors |

### Medium-Value Components (Partial Reuse)

| Component | Reuse Level | Notes |
|-----------|-----------|-------|
| Header, TabBar | ⭐⭐⭐ | Navigation logic reusable, styling will change |
| ConnectButton | ⭐⭐ | Wrapper around RainbowKit, light refactor |

### Low Reuse (Specific Business Logic)

| Component | Reuse Level | Notes |
|-----------|-----------|-------|
| RegisterSheet | ⭐⭐ | 2-step flow tied to current auth, may need redesign |
| AppLayout | ⭐ | Route protection logic reusable, overall structure may change |
| Page Components | ⭐ | Layout/UX will likely be redesigned |

### Hooks & Utils (High Reuse Potential)

| Hook/Util | Reuse Level | Notes |
|-----------|-----------|-------|
| useUser, useBilling, useDeposit, useOperator | ⭐⭐⭐⭐⭐ | Business logic independent of UI, reuse as-is |
| useTransactionFlow, parseContractError | ⭐⭐⭐⭐⭐ | Pure abstraction, essential for Web3 flow |
| Format utils (shortenAddress, formatAmount, etc.) | ⭐⭐⭐⭐⭐ | Standalone, no UI dependency |
| pendingSync utilities | ⭐⭐⭐⭐ | Offline-first pattern, reusable |

---

## Key Files for Reconstruction Reference

**Must Keep/Reuse**:
- `/src/hooks/*` (all business logic)
- `/src/utils/*` (formatting + offline sync)
- `/src/services/*` (API layer)
- `/src/config/*` (contracts, chains, wagmi)
- `/src/types/*` (TypeScript definitions)

**Can Redesign**:
- `/src/pages/*` (layout/UX)
- `/src/components/layout/*` (navigation structure)
- `/src/index.css` (colors need update to deep blue gradient)

**May Refactor**:
- `/src/components/ui/*` (can stay if shadcn/base-ui works; else use new component lib)
- Auth flow in RegisterSheet (depends on new design)

---

## Summary Statement

**Current Status**: MVP feature-complete for core user flows (register → deposit → buy service → pay bill). 9 pages, 13 components, 12+ hooks, dual-write blockchain+backend pattern with offline sync support.

**Completion %**: ~80% (missing: final bill detail, real email, real notifications, traffic card NFT).

**Reconstruction Approach**: 
- Reuse 100% of hooks, services, utils, types, config
- Redesign pages & layout (70% new UI)
- Keep or refactor UI components based on new design system choice
- Target deep blue gradient (`#0C2340` → `#1E40AF`) for primary brand colors
- No breaking changes to Web3 or backend integration patterns
