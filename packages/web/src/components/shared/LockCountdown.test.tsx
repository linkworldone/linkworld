import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";

// LockCountdown 读 useLockExpiry（链上 getLockExpiry）→ 渲染倒计时/解锁态。
// mock 链上读值，逐项设置 expiry。
let expiryData: bigint | undefined = undefined;
vi.mock("@/hooks/contracts", () => ({
  useLockExpiry: () => ({ data: expiryData }),
}));

import {
  LockCountdown,
  lockState,
  formatRemaining,
} from "./LockCountdown";

const OWNER = "0xOWNER" as `0x${string}`;
const DAY = 86_400;

describe("lockState 纯函数（边界 >= 对齐合约）", () => {
  it("DEP-01a: expiry=0 → none（无锁仓）", () => {
    const s = lockState(0n, 1000);
    expect(s.status).toBe("none");
  });

  it("DEP-01b: now < expiry → locked + 正剩余", () => {
    const now = 1000;
    const s = lockState(BigInt(now + DAY), now);
    expect(s.status).toBe("locked");
    expect(s.remainingSec).toBe(DAY);
  });

  it("DEP-01c: now >= expiry（边界 ===）→ unlocked（到期即可提，不可用 <）", () => {
    const now = 5000;
    const s = lockState(BigInt(now), now); // now === expiry
    expect(s.status).toBe("unlocked");
    expect(s.remainingSec).toBe(0);
  });

  it("DEP-01d: now > expiry → unlocked", () => {
    const s = lockState(BigInt(4000), 5000);
    expect(s.status).toBe("unlocked");
  });
});

describe("formatRemaining 文案", () => {
  it("Dd Hh Mm 组合", () => {
    expect(formatRemaining(2 * DAY + 3 * 3600 + 5 * 60)).toBe("2d 3h 5m");
  });
  it("不足 1 天显示 0d", () => {
    expect(formatRemaining(3 * 3600 + 5 * 60)).toBe("0d 3h 5m");
  });
});

describe("LockCountdown 渲染", () => {
  it("DEP-01: now < expiry → 锁仓中（含倒计时），Withdraw 由父禁用（暴露 unlocked=false）", () => {
    const now = Math.floor(Date.now() / 1000);
    expiryData = BigInt(now + 2 * DAY);
    const onState = vi.fn();
    render(<LockCountdown address={OWNER} onStateChange={onState} />);
    expect(screen.getByText(/锁仓中/)).toBeInTheDocument();
    // 倒计时含天数
    expect(screen.getByText(/剩余/)).toBeInTheDocument();
    expect(onState).toHaveBeenCalledWith(expect.objectContaining({ status: "locked" }));
  });

  it("DEP-01: now >= expiry → 可提（锁仓已满，可提取本金）", () => {
    const now = Math.floor(Date.now() / 1000);
    expiryData = BigInt(now - 10); // 已过期
    const onState = vi.fn();
    render(<LockCountdown address={OWNER} onStateChange={onState} />);
    expect(screen.getByText(/锁仓已满，可提取本金/)).toBeInTheDocument();
    expect(onState).toHaveBeenCalledWith(expect.objectContaining({ status: "unlocked" }));
  });

  it("DEP-02: 文案不含「利息」（interest 恒 0）", () => {
    const now = Math.floor(Date.now() / 1000);
    expiryData = BigInt(now - 10);
    const { container } = render(<LockCountdown address={OWNER} />);
    expect(container.textContent).not.toContain("利息");
  });

  it("expiry=0 → 无锁仓提示", () => {
    expiryData = 0n;
    render(<LockCountdown address={OWNER} />);
    expect(screen.getByText(/暂无锁仓|无锁仓/)).toBeInTheDocument();
  });
});
