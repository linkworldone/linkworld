# LinkWorld Utils & Hooks Reference

## Utility Functions

### Format Utilities (`/src/utils/format.ts`)

| Function | Signature | Purpose |
|----------|-----------|---------|
| `shortenAddress` | `(address: string, chars = 4) => string` | Truncate wallet address for display (e.g., `0x1234...5678`) |
| `parseUnits` | `(value: string, decimals = 18) => bigint` | Convert decimal string to bigint with specified decimals |
| `formatAmount` | `(wei: bigint, decimals = 18, displayDecimals = 2) => string` | Format bigint wei to readable decimal string |
| `formatUSD` | `(amount: string) => string` | Format string to USD currency string (e.g., `$1.23`) |
| `formatDate` | `(isoDate: string) => string` | Format ISO date to "MMM D, YYYY" format |
| `timeAgo` | `(isoDate: string) => string` | Convert ISO date to relative time (e.g., "5 minutes ago") |

### Utility Library (`/src/lib/utils.ts`)

| Function | Signature | Purpose |
|----------|-----------|---------|
| `cn` | `(...inputs: ClassValue[]) => string` | Merge clsx + tailwind-merge for conditional Tailwind classes (prevents conflicts) |

### Pending Sync Utils (`/src/utils/pendingSync.ts`)

| Function | Signature | Purpose |
|----------|-----------|---------|
| `savePendingSync` | `(key: string, data: unknown) => void` | Store failed sync requests to localStorage for retry (offline support) |
| `clearPendingSync` | `(key: string) => void` | Remove successful sync from localStorage |
| `retryWithBackoff` | `(fn: () => Promise<boolean>) => Promise<boolean>` | Retry async function with exponential backoff strategy |

---

## Custom Hooks

### User Management (`/src/hooks/useUser.ts`)

| Hook | Return Type | Purpose |
|------|-------------|---------|
| `useUser(address?: string)` | `UseQueryResult<User \| null>` | Fetch user profile from backend by wallet address; staleTime: 30s |
| `useRegister()` | `{ register, txState, backendSync, ... }` | 2-step registration: (1) contract mint NFT, (2) backend API sync with email |
| `useSendVerificationCode()` | `UseMutationResult<{ success }>` | Send verification code to email (currently mock) |
| `useVerifyEmail()` | `UseMutationResult<boolean>` | Verify email with code (currently mock, always passes) |

**Notes**: 
- `useRegister()` uses dual-write pattern: contract tx (source of truth) + backend sync (offline retry via pending sync)
- Email verification is mocked (backend not implemented)

### Notifications (`/src/hooks/useNotification.ts`)

| Hook | Return Type | Purpose |
|------|-------------|---------|
| `useNotifications(address?: string)` | `UseQueryResult<Notification[]>` | Fetch all notifications for user; staleTime: none (default) |
| `useUnreadCount(address?: string)` | `UseQueryResult<number>` | Fetch unread notification count; refetchInterval: 30s |
| `useMarkAsRead()` | `UseMutationResult<void>` | Mark single notification as read; invalidates both queries |
| `useMarkAllAsRead()` | `UseMutationResult<void>` | Mark all notifications as read; invalidates both queries |

**Service**: Uses `/src/services/mock/notificationService` (notification API not yet on backend)

### Billing (`/src/hooks/useBilling.ts`)

| Hook | Return Type | Purpose |
|------|-------------|---------|
| `useBills(address?, filter?: "unpaid" \| "paid")` | `UseQueryResult<Bill[]>` | Fetch bills from backend with optional filter; staleTime: 30s |
| `useBillDetail(billId?: string)` | `UseQueryResult<Bill \| null>` | Fetch single bill detail (currently returns null - no backend endpoint) |
| `useMonthEstimate(address?)` | `UseQueryResult<MonthEstimate>` | Fetch current month usage estimate & projected cost from usage API; staleTime: 60s |
| `usePayBill()` | `{ payBill, txState, recordToBackend, ... }` | Pay single bill: (1) contract payBill(), (2) backend sync with offline retry |

**Pattern**: Dual-write (contract + backend), with pending sync for offline scenarios

### Deposit (`/src/hooks/useDeposit.ts`)

| Hook | Return Type | Purpose |
|------|-------------|---------|
| `useDeposit(address?)` | `{ data: DepositInfo, refetch, ... }` | Read deposit balance from contract (source of truth); wrapper around `useDepositBalance()` |
| `useDepositHistory(address?)` | `UseQueryResult<DepositRecord[]>` | Fetch deposit/withdraw/deduction transaction history from backend; staleTime: 10s |
| `useDepositMutation()` | `{ deposit, txState, recordToBackend, ... }` | Deposit ETH: (1) contract deposit(), (2) backend record with offline retry |
| `useWithdrawMutation()` | `{ withdraw, txState, recordToBackend, ... }` | Withdraw ETH: (1) contract withdraw(), (2) backend record with offline retry |

**Notes**:
- Contract is source of truth for balance (not backend)
- Backend records transaction history
- Uses pending sync for offline-first deposit/withdraw records

### Operator/Services (`/src/hooks/useOperator.ts`)

| Hook | Return Type | Purpose |
|------|-------------|---------|
| `useRegions()` | `UseQueryResult<Region[]>` | Fetch all regions from operator API; staleTime: 5min (cached aggregation) |
| `useOperatorsByRegion(regionCode?)` | `UseQueryResult<Operator[]>` | Fetch operators in specific region; staleTime: 5min |
| `useMyNumbers(address?)` | `UseQueryResult<VirtualNumber[]>` | Fetch user's active virtual numbers from backend + enrich with operator info |
| `useApplyNumber()` | `{ applyNumber, txState, invalidate }` | Request virtual number: contract activate + backend archive + cache invalidation |

### Transaction State (`/src/hooks/useTransactionFlow.ts`)

| Hook/Function | Signature | Purpose |
|---|---|---|
| `useTxState(params)` | `(params: { hash?, isPending, isConfirming, isSuccess, error }) => TxState` | Map wagmi write/wait states to unified `TxState` (idle/pending-signature/pending-confirmation/success/error) |
| `parseContractError()` | `(error: unknown) => string` | Parse contract revert reason and return Chinese error message or fallback |

**TxStatus States**:
```
"idle" → no transaction active
"pending-signature" → waiting for user to sign in wallet
"pending-confirmation" → signed, waiting for block confirmation
"success" → confirmed on-chain
"error" → revert or user rejection
```

---

## Contract Interaction Hooks

### User Registry (`/src/hooks/contracts/useUserRegistry.ts`)

| Hook | Purpose |
|------|---------|
| `useIsRegistered(address?)` | Read: Check if wallet is registered on UserRegistry contract |
| `useContractRegister()` | Write: Register wallet with email on contract; returns `{ register, hash, isPending, isConfirming, isSuccess, error }` |

### Deposit Contract (`/src/hooks/contracts/useDepositContract.ts`)

| Hook | Purpose |
|------|---------|
| `useDepositBalance(address?)` | Read: Get current deposit balance from contract |
| `useContractDeposit()` | Write: Deposit ETH into contract |
| `useContractWithdraw()` | Write: Withdraw ETH from contract |

### Payment Contract (`/src/hooks/contracts/usePaymentContract.ts`)

| Hook | Purpose |
|------|---------|
| `useContractPayBill()` | Write: Pay bill by billId + amount on Payment contract |

### Service Manager (`/src/hooks/contracts/useServiceManager.ts`)

| Hook | Purpose |
|------|---------|
| `useContractActivateService()` | Write: Activate virtual number service on contract |

### Traffic Card / Other (`/src/hooks/contracts/useTrafficCard.ts`)

| Hook | Purpose |
|------|---------|
| TBD | NFT traffic card interaction hooks (to be read) |

---

## Service Layer

### API Client (`/src/services/api/client.ts`)

```typescript
const apiClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
  headers: { "Content-Type": "application/json" },
});
```

**Response Interceptor**: Auto-extracts `.data`, maps errors to `Promise.reject(Error(message))`

### API Services

| Service | Module | Key Methods |
|---------|--------|-------------|
| **User API** | `/src/services/api/userApi.ts` | `getUser(wallet)`, `register(wallet, email, tokenId)` |
| **Operator API** | `/src/services/api/operatorApi.ts` | `getRegions()`, `getOperatorsByRegion(regionCode)` |
| **Deposit API** | `/src/services/api/depositApi.ts` | `getHistory(wallet)`, `recordDeposit(wallet, amount)`, `recordWithdraw(wallet)` |
| **Billing API** | `/src/services/api/billingApi.ts` | `getBills(wallet, filter?)`, `recordPayment(wallet, billId)` |
| **Usage API** | `/src/services/api/usageApi.ts` | `getUsage(wallet)` (returns `MonthEstimate`) |

### Mock Services

| Service | Module | Status |
|---------|--------|--------|
| **Notifications** | `/src/services/mock/notificationService.ts` | Mock (backend API not implemented) |

---

## Summary Statistics

- **Total Utility Functions**: 9 (format.ts + lib + pendingSync)
- **Total Custom Hooks**: 17+ (user, notification, billing, deposit, operator, txFlow, contract-specific)
- **Total API Services**: 5 (all real backend)
- **Total Mock Services**: 1 (notifications)
- **Primary Pattern**: React Query (useQuery/useMutation) + wagmi (useWriteContract/useReadContract/useWaitForTransactionReceipt)
- **Offline Support**: Pending sync for dual-write operations (contract + backend)
