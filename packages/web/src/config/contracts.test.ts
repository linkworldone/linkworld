import { describe, it, expect } from "vitest";
import {
  CONTRACTS,
  getContractAddress,
  getUsdt,
  getUsdtDecimals,
} from "./contracts";
import hardhatDeployment from "./deployments/hardhat.json";

// T12 31337 配置加载守门（CFG-01..04）。
// 本地 hardhat node 在本分支因 contracts toolchain 未装而跑不起（见 stage-test-done.md），
// 此处以单元测固化 web 侧可解析的那部分：deployments/hardhat.json 单一出口、
// USDT 6 位精度根来源、零地址/未知链兜底抛错。
describe("config/contracts 单一出口（CFG-01..04）", () => {
  it("CFG-01: 31337 地址全部取自 deployments/hardhat.json（无手抄漂移）", () => {
    const names = [
      "UserRegistry",
      "FeeManager",
      "Deposit",
      "ServiceManager",
      "Payment",
      "Oracle",
      "TrafficCardNFT",
    ] as const;
    for (const name of names) {
      expect(getContractAddress(31337, name)).toBe(
        hardhatDeployment.proxies[name]
      );
    }
    // USDT 地址同源。
    expect(getUsdt(31337)).toBe(hardhatDeployment.usdt);
  });

  it("CFG-02: getUsdtDecimals(31337) === 6（精度红线根来源，从 deployments 读不硬编码）", () => {
    expect(getUsdtDecimals(31337)).toBe(6);
    expect(getUsdtDecimals(31337)).toBe(hardhatDeployment.usdtDecimals);
  });

  it("CFG-03: 未知 chainId → 抛错（不静默返回 undefined 地址）", () => {
    expect(() => getContractAddress(999, "Deposit")).toThrow(
      /No contract addresses for chainId 999/
    );
    expect(() => getUsdt(999)).toThrow(/No contract addresses/);
    expect(() => getUsdtDecimals(999)).toThrow(/No contract addresses/);
  });

  it("CFG-04: 零地址（未部署）→ 抛错保护，杜绝把 0x0 当合约调用", () => {
    const ZERO = "0x0000000000000000000000000000000000000000" as const;
    // 注入一条全零地址的临时链，验证抛错分支。
    CONTRACTS[424242] = {
      UserRegistry: ZERO,
      FeeManager: ZERO,
      Deposit: ZERO,
      ServiceManager: ZERO,
      Payment: ZERO,
      Oracle: ZERO,
      TrafficCardNFT: ZERO,
      usdt: ZERO,
      usdtDecimals: 6,
    };
    try {
      expect(() => getContractAddress(424242, "Deposit")).toThrow(
        /Contract Deposit not deployed on chainId 424242/
      );
      expect(() => getUsdt(424242)).toThrow(/USDT not deployed on chainId 424242/);
    } finally {
      delete CONTRACTS[424242];
    }
  });
});
