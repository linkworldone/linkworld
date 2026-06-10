import { describe, it, expect } from "vitest";
import { render } from "@testing-library/react";
import { AmountDisplay } from "./AmountDisplay";

// 金色铁律分流（DESIGN.md §B2）：默认 navy（卡内 ≈12:1）；深底显式金；不再无条件吃金。
describe("AmountDisplay 金色铁律分流（T6）", () => {
  it("默认 tone=auto → 卡内 navy（text-on-light-primary，非金）", () => {
    const { container } = render(<AmountDisplay amount="100.00" currency="USDT" />);
    const span = container.querySelector("span");
    expect(span?.className).toMatch(/text-text-on-light-primary/);
    // 默认不染金（金在米白上 ≈2:1 不达标）。
    expect(span?.className).not.toMatch(/gold/);
    // 金额走等宽数字。
    expect(span?.className).toMatch(/tabular-nums/);
  });

  it("tone=gold-on-dark → 深底金（text-on-dark-gold）", () => {
    const { container } = render(
      <AmountDisplay amount="100.00" tone="gold-on-dark" />
    );
    const span = container.querySelector("span");
    expect(span?.className).toMatch(/text-text-on-dark-gold/);
  });

  it("显式 colorClass 优先级高于 tone", () => {
    const { container } = render(
      <AmountDisplay amount="100.00" tone="gold-on-dark" colorClass="text-status-danger" />
    );
    const span = container.querySelector("span");
    expect(span?.className).toMatch(/text-status-danger/);
    expect(span?.className).not.toMatch(/gold/);
  });
});
