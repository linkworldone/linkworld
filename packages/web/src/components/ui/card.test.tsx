import { describe, it, expect } from "vitest";
import { render, screen } from "@testing-library/react";
import { Card, CardHeader, CardTitle, CardContent } from "./card";
import { Input } from "./input";

// T6 基建原子组件最简渲染测：吃 CSS 变量语义类，不写死 hex。
describe("ui/card 暖米白卡（T6）", () => {
  it("渲染 + 暖米白底 + 金线 + 卡内 navy 文字语义类", () => {
    const { container } = render(
      <Card>
        <CardHeader>
          <CardTitle>余额</CardTitle>
        </CardHeader>
        <CardContent>内容</CardContent>
      </Card>
    );
    const card = container.querySelector("[data-slot='card']");
    expect(card).toBeInTheDocument();
    // 暖米白底 + 金线描边 + 卡内深色（金色铁律：卡内文字 navy 非金）。
    expect(card?.className).toMatch(/bg-surface-card/);
    expect(card?.className).toMatch(/border-surface-card-line/);
    expect(card?.className).toMatch(/text-text-on-light-primary/);
    expect(screen.getByText("余额")).toBeInTheDocument();
    expect(screen.getByText("内容")).toBeInTheDocument();
  });
});

describe("ui/input 深蓝金输入框（T6）", () => {
  it("渲染 + 凹陷米白底 + navy 文字 + royal focus ring", () => {
    const { container } = render(<Input placeholder="搜索" />);
    const input = container.querySelector("[data-slot='input']");
    expect(input).toBeInTheDocument();
    expect(input?.className).toMatch(/bg-surface-input/);
    expect(input?.className).toMatch(/text-text-on-light-primary/);
    expect(input?.className).toMatch(/focus-visible:border-brand-royal/);
  });
});
