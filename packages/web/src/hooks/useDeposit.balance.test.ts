import { describe, it, expect, vi } from "vitest";

// REC-04：余额来源 = 链上 getDepositAmount（source of truth, design §9），不取后端自述。
// useDeposit 内部走 useDepositBalance(=useReadContract getDepositAmount)。
vi.mock("./contracts", () => ({
  useDepositBalance: vi.fn(() => ({
    data: 5_000_000n, // 链上读到的 confirmed 本金（5 USDT, 6 位）
    refetch: vi.fn(),
    isLoading: false,
  })),
}));
vi.mock("wagmi", () => ({
  useAccount: () => ({ address: "0xabc" }),
}));
vi.mock("@tanstack/react-query", () => ({
  useQuery: vi.fn(() => ({ data: undefined })),
  useQueryClient: () => ({ invalidateQueries: vi.fn() }),
}));

import { renderHook } from "@testing-library/react";
import { useDeposit } from "./useDeposit";

describe("useDeposit 余额读链（REC-04）", () => {
  it("balance 取自链上 getDepositAmount，currency=USDT（非后端、非 ETH）", () => {
    const { result } = renderHook(() => useDeposit("0xabc"));
    expect(result.current.data?.balance).toBe(5_000_000n);
    expect(result.current.data?.currency).toBe("USDT");
  });
});
