import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import type { DepositRecord } from "@/types";
import type { SimRecord } from "@/services/api/simApi";

// ── 可变 mock 状态 ───────────────────────────────────────────────
let deposits: DepositRecord[] = [];
let sims: SimRecord[] = [];

vi.mock("wagmi", () => ({
  useAccount: () => ({ address: "0xabc" as `0x${string}`, chainId: 31337 }),
  http: () => ({}),
}));

vi.mock("@/hooks/useDeposit", () => ({
  useDepositHistory: () => ({ data: deposits }),
}));

vi.mock("@/hooks/useSim", () => ({
  useMySims: () => ({ data: sims }),
}));

import Billing from "./Billing";

function renderPage() {
  return render(
    <MemoryRouter>
      <Billing />
    </MemoryRouter>,
  );
}

function makeDeposit(over: Partial<DepositRecord> = {}): DepositRecord {
  return {
    id: "d1",
    type: "deposit",
    amount: 50_000_000n, // 50 USDT（6 位最小单位）→ 50/10 = 5 张卡
    currency: "USDT",
    status: "confirmed",
    timestamp: "2026-05-10T00:00:00.000Z",
    txHash: "0xdep",
    ...over,
  };
}

function makeSim(over: Partial<SimRecord> = {}): SimRecord {
  return {
    id: "s1",
    days: 7,
    destination: "日本",
    deliveryType: "physical",
    recipient: "Alice",
    addressLine: "Tokyo",
    activationUrl: "",
    status: "pending",
    createdAt: "2026-05-20T00:00:00.000Z",
    ...over,
  };
}

beforeEach(() => {
  deposits = [];
  sims = [];
});

describe("History 页（流量卡活动时间线）", () => {
  it("空态：无记录时展示空态文案", () => {
    renderPage();
    expect(screen.getByText(/No records yet/)).toBeInTheDocument();
  });

  it("获取记录：50 USDT → 获取 5 张流量卡，副文展示充值金额", () => {
    deposits = [makeDeposit()];
    renderPage();
    expect(screen.getByText("Got 5 traffic cards")).toBeInTheDocument();
    expect(screen.getByText(/Deposited 50\.00 USDT/)).toBeInTheDocument();
  });

  it("销毁记录：展示销毁张数 → 领取 SIM、天数·目的地、状态徽章", () => {
    sims = [makeSim()];
    renderPage();
    expect(screen.getByText("Burned 7 → redeemed SIM")).toBeInTheDocument();
    expect(screen.getByText(/7 days unlimited data · 日本/)).toBeInTheDocument();
    // pending → 处理中
    expect(screen.getByText("Processing")).toBeInTheDocument();
  });

  it("confirmed SIM 展示「已确认」", () => {
    sims = [makeSim({ status: "confirmed" })];
    renderPage();
    expect(screen.getByText("Confirmed")).toBeInTheDocument();
  });

  it("只渲染 type==='deposit' 的获取记录（忽略 withdraw/deduction）", () => {
    deposits = [
      makeDeposit({ id: "d1", amount: 50_000_000n }),
      makeDeposit({ id: "d2", type: "withdraw", amount: 30_000_000n }),
      makeDeposit({ id: "d3", type: "deduction", amount: 20_000_000n }),
    ];
    renderPage();
    expect(screen.getAllByText(/Got \d+ traffic cards/)).toHaveLength(1);
    expect(screen.getByText("Got 5 traffic cards")).toBeInTheDocument();
  });

  it("获取与销毁混排，按日期倒序（较新的排前）", () => {
    deposits = [makeDeposit({ timestamp: "2026-05-10T00:00:00.000Z" })];
    sims = [makeSim({ createdAt: "2026-05-20T00:00:00.000Z" })];
    const { container } = renderPage();
    const titles = Array.from(container.querySelectorAll(".text-sm.font-bold")).map(
      (n) => n.textContent ?? "",
    );
    // SIM(05-20) 在 deposit(05-10) 之前
    const burnIdx = titles.findIndex((t) => t.includes("Burned"));
    const acquireIdx = titles.findIndex((t) => t.includes("Got"));
    expect(burnIdx).toBeGreaterThanOrEqual(0);
    expect(acquireIdx).toBeGreaterThan(burnIdx);
  });
});
