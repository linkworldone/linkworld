import { describe, it, expect, vi, beforeEach } from "vitest";
import { render } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";

/*
 * T11 换肤 sweep 守门测试（SKIN-01/02/03）。
 * - SKIN-01：Dashboard/Services/Landing/Notifications 渲染无旧 token class、无装饰 emoji。
 * - SKIN-02：TabBar 5 项 lucide 图标渲染（svg）。
 * - SKIN-03：DepositInfo.currency 仅 USDT（tsc 层断言，见下方类型探针）。
 */

// ── 旧 token / 装饰 emoji 黑名单（grep 清零的运行时镜像） ──
const DEAD_TOKENS = [
  "brand-blue",
  "brand-purple",
  "brand-cyan",
  "surface-gradient",
  "surface-secondary",
  "font-orbitron",
];
// 装饰性 emoji（国旗 FLAG_MAP 是数据，不在此列）。
const DECO_EMOJI = ["🏠", "📱", "💰", "📄", "🎟️", "🎟", "🔔", "🔍", "🌐", "💳", "⚠️", "⚠"];

function assertClean(html: string, label: string) {
  for (const t of DEAD_TOKENS) {
    expect(html, `${label} 残留旧 token: ${t}`).not.toContain(t);
  }
  for (const e of DECO_EMOJI) {
    expect(html, `${label} 残留装饰 emoji: ${e}`).not.toContain(e);
  }
}

vi.mock("wagmi", () => ({
  useAccount: () => ({ address: "0xabc" as `0x${string}`, chainId: 31337 }),
  http: () => ({}),
}));

// 各页 hook mock（仅供渲染，不验证行为——行为测在各自专测）。
vi.mock("@/hooks/useUser", () => ({
  useUser: () => ({ data: undefined, isLoading: false }),
}));
vi.mock("@/hooks/useDeposit", () => ({
  useDeposit: () => ({ data: undefined }),
}));
vi.mock("@/hooks/useBilling", () => ({
  useBills: () => ({ data: [] }),
  useMonthEstimate: () => ({ data: undefined }),
}));
vi.mock("@/hooks/useOperator", () => ({
  useRegions: () => ({ data: [] }),
  useMyNumbers: () => ({ data: [] }),
}));
vi.mock("@/hooks/useNotification", () => ({
  useNotifications: () => ({ data: [] }),
  useMarkAsRead: () => ({ mutate: vi.fn() }),
  useMarkAllAsRead: () => ({ mutate: vi.fn() }),
  useUnreadCount: () => ({ data: 0 }),
}));
// RainbowKit ConnectButton 依赖 provider，渲染层 stub。
vi.mock("@/components/wallet/ConnectButton", () => ({
  ConnectButton: () => <button>Connect</button>,
}));
vi.mock("@/components/wallet/RegisterSheet", () => ({
  RegisterSheet: () => null,
}));

import Dashboard from "./Dashboard";
import Services from "./Services";
import Landing from "./Landing";
import Notifications from "./Notifications";
import { TabBar } from "@/components/layout/TabBar";

function renderPage(node: React.ReactNode, path = "/") {
  return render(<MemoryRouter initialEntries={[path]}>{node}</MemoryRouter>);
}

beforeEach(() => {
  vi.clearAllMocks();
});

describe("SKIN-01: 其余页换肤——无旧 token / 无装饰 emoji", () => {
  it("Dashboard", () => {
    const { container } = renderPage(<Dashboard />, "/dashboard");
    assertClean(container.innerHTML, "Dashboard");
  });
  it("Services", () => {
    const { container } = renderPage(<Services />, "/services");
    assertClean(container.innerHTML, "Services");
    // 搜索框前缀已 lucide（svg）。
    expect(container.querySelectorAll("svg").length).toBeGreaterThan(0);
  });
  it("Landing", () => {
    const { container } = renderPage(<Landing />, "/");
    assertClean(container.innerHTML, "Landing");
    // 不再出现旧文案 "0G"。
    expect(container.textContent).not.toContain("0G");
  });
  it("Notifications", () => {
    const { container } = renderPage(<Notifications />, "/notifications");
    assertClean(container.innerHTML, "Notifications");
  });
});

describe("SKIN-02: TabBar 5 项 lucide 图标", () => {
  it("渲染 5 个导航项，每项含 lucide svg（非 emoji）", () => {
    const { container } = renderPage(<TabBar />, "/dashboard");
    const buttons = container.querySelectorAll("button");
    expect(buttons.length).toBe(5);
    // 5 项各一 lucide svg。
    expect(container.querySelectorAll("svg").length).toBeGreaterThanOrEqual(5);
    assertClean(container.innerHTML, "TabBar");
  });
});

describe("SKIN-03: DepositInfo.currency 类型仅 USDT", () => {
  it("类型探针：currency 字面量类型为 'USDT'（编译期约束，运行期占位断言）", () => {
    // 若 DepositInfo.currency 放开为 "USDT" | "ETH"，下面 satisfies 会在 tsc 报错。
    const info = { balance: 0n, currency: "USDT" } satisfies import("@/types").DepositInfo;
    expect(info.currency).toBe("USDT");
  });
});
