import { describe, it, expect } from "vitest";
import { cn } from "./utils";

// Smoke test: 验证 vitest 测试设施可跑（T0 基线）
describe("cn", () => {
  it("merges class names", () => {
    expect(cn("a", "b")).toBe("a b");
  });

  it("dedupes conflicting tailwind classes (twMerge)", () => {
    expect(cn("px-2", "px-4")).toBe("px-4");
  });

  it("ignores falsy values", () => {
    expect(cn("a", false, null, undefined, "b")).toBe("a b");
  });
});
