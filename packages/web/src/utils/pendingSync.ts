// src/utils/pendingSync.ts

const STORAGE_PREFIX = "linkworld_pending_";

export interface PendingSyncItem<T = unknown> {
  data: T;
  createdAt: number;
  retryCount: number;
}

export function savePendingSync<T>(key: string, data: T): void {
  const item: PendingSyncItem<T> = {
    data,
    createdAt: Date.now(),
    retryCount: 0,
  };
  try {
    localStorage.setItem(STORAGE_PREFIX + key, JSON.stringify(item));
  } catch {
    // localStorage full or unavailable, silently fail
  }
}

export function getPendingSync<T>(key: string): PendingSyncItem<T> | null {
  try {
    const raw = localStorage.getItem(STORAGE_PREFIX + key);
    if (!raw) return null;
    return JSON.parse(raw) as PendingSyncItem<T>;
  } catch {
    return null;
  }
}

export function clearPendingSync(key: string): void {
  try {
    localStorage.removeItem(STORAGE_PREFIX + key);
  } catch {
    // silently fail
  }
}

export function incrementRetryCount(key: string): number {
  const item = getPendingSync(key);
  if (!item) return 0;
  item.retryCount += 1;
  try {
    localStorage.setItem(STORAGE_PREFIX + key, JSON.stringify(item));
  } catch {
    // silently fail
  }
  return item.retryCount;
}

/**
 * Retry an async function with exponential backoff.
 * Returns true if succeeded, false if all retries exhausted.
 */
export async function retryWithBackoff(
  fn: () => Promise<void>,
  maxRetries = 3,
  baseDelay = 2000
): Promise<boolean> {
  for (let i = 0; i < maxRetries; i++) {
    try {
      await fn();
      return true;
    } catch {
      if (i < maxRetries - 1) {
        await new Promise((r) => setTimeout(r, baseDelay * Math.pow(2, i)));
      }
    }
  }
  return false;
}
