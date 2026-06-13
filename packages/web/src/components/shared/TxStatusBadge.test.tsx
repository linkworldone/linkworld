import { describe, it, expect } from "vitest";
import { render, screen } from "@testing-library/react";
import { TxStatusBadge } from "./TxStatusBadge";

// BADGE-01：通用交易三态徽章渲染（design §3.1 / DESIGN.md 三态表）。
// 铁律：pending「处理中」弱化、绝不染绿；confirmed 才染绿；failed 区分 reorg vs revert。
describe("TxStatusBadge 三态渲染（BADGE-01）", () => {
  it("pending → 文案『处理中』，且不染绿（无 status-success 类）", () => {
    const { container } = render(<TxStatusBadge status="pending" />);
    expect(screen.getByText("Processing")).toBeInTheDocument();
    const badge = container.querySelector("[data-slot='tx-status-badge']");
    expect(badge).toHaveAttribute("data-status", "pending");
    // 不染绿（铁律）：class 不含任何 success 语义类。
    expect(badge?.className).not.toMatch(/success/);
  });

  it("confirmed → 文案『已确认』，染绿（success 语义）", () => {
    const { container } = render(<TxStatusBadge status="confirmed" />);
    expect(screen.getByText("Confirmed")).toBeInTheDocument();
    const badge = container.querySelector("[data-slot='tx-status-badge']");
    expect(badge).toHaveAttribute("data-status", "confirmed");
    expect(badge?.className).toMatch(/success/);
  });

  it("failed(revert) → 文案『失败』", () => {
    const { container } = render(<TxStatusBadge status="failed" failureReason="revert" />);
    expect(screen.getByText("Failed")).toBeInTheDocument();
    const badge = container.querySelector("[data-slot='tx-status-badge']");
    expect(badge).toHaveAttribute("data-status", "failed");
    expect(badge).toHaveAttribute("data-failure", "revert");
  });

  it("failed(reorg) → 文案区分『已回退』（与 revert 文案不同）", () => {
    const { container } = render(<TxStatusBadge status="failed" failureReason="reorg" />);
    expect(screen.getByText("Reverted")).toBeInTheDocument();
    const badge = container.querySelector("[data-slot='tx-status-badge']");
    expect(badge).toHaveAttribute("data-failure", "reorg");
  });

  it("failed 缺省 failureReason → 不崩，默认『失败』", () => {
    render(<TxStatusBadge status="failed" />);
    expect(screen.getByText("Failed")).toBeInTheDocument();
  });
});
