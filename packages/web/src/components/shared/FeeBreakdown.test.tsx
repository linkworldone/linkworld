import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen } from "@testing-library/react";

// 受控 mock useFeeRate / useCalculateFee（读链结果）。
let rate: { label?: string; isLoading: boolean; isError: boolean } = {
  label: "1.5%",
  isLoading: false,
  isError: false,
};
let fee: { fee?: bigint; isLoading: boolean; isError: boolean } = {
  fee: 1_500_000n,
  isLoading: false,
  isError: false,
};
vi.mock("@/hooks/contracts", () => ({
  useFeeRate: () => rate,
  useCalculateFee: () => fee,
}));

import { FeeBreakdown } from "./FeeBreakdown";

beforeEach(() => {
  rate = { label: "1.5%", isLoading: false, isError: false };
  fee = { fee: 1_500_000n, isLoading: false, isError: false };
});

describe("FeeBreakdown —— 手续费读链展示（design §3.6）", () => {
  it("读链成功 → 「平台手续费 (1.5%)」+ 费额（formatAmount 6 位）", () => {
    render(<FeeBreakdown amount={100_000_000n} />);
    expect(screen.getByText(/Platform fee \(1\.5%\)/)).toBeInTheDocument();
    expect(screen.getByText(/1\.50 USDT/)).toBeInTheDocument();
  });

  it("FEE-02: 读链失败 → 费率/费额均「--」，绝不写死兜底（无 2.5%）", () => {
    rate = { label: undefined, isLoading: false, isError: true };
    fee = { fee: undefined, isLoading: false, isError: true };
    const { container } = render(<FeeBreakdown amount={100_000_000n} />);
    expect(screen.getByText(/Platform fee \(--\)/)).toBeInTheDocument();
    // 费额位也是 --。
    expect(screen.getAllByText("--").length).toBeGreaterThanOrEqual(1);
    // 不出现任何写死费率文案。
    expect(container.textContent).not.toContain("2.5%");
    expect(container.textContent).not.toContain("0.025");
  });

  it("读链 loading → skeleton 占位（不写死兜底数字）", () => {
    rate = { label: undefined, isLoading: true, isError: false };
    fee = { fee: undefined, isLoading: true, isError: false };
    const { container } = render(<FeeBreakdown amount={100_000_000n} />);
    expect(container.querySelectorAll('[data-slot="fee-skeleton"]').length).toBeGreaterThanOrEqual(1);
    expect(container.textContent).not.toContain("2.5%");
  });
});
