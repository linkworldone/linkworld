import { describe, it, expect, vi, beforeEach } from "vitest";

// T5 WalletAuth：signedPost = GET nonce → EIP-712 signTypedData → 带签名头 POST。
// EIP-712 域/字段/action/header 必须逐字节对齐后端 packages/backend/internal/middleware/middleware.go
//（WalletAuthDigest / NewWalletAuth），否则后端 ecrecover 验签失败。
//
// mock 边界：
//  - ./client          → axios 实例（get 取 nonce / post 带头写）
//  - wagmi/actions      → signTypedData(config,…)（**core action，非 React hook**）+ getChainId
//  - ../../config/wagmi → wagmiConfig（core action 需要的 config 入参）

vi.mock("./client", () => ({
  apiClient: { get: vi.fn(), post: vi.fn() },
}));

vi.mock("wagmi/actions", () => ({
  signTypedData: vi.fn(),
  getChainId: vi.fn(() => 31337),
}));

vi.mock("../../config/wagmi", () => ({
  wagmiConfig: { __mock: "config" },
}));

import { apiClient } from "./client";
import { signTypedData, getChainId } from "wagmi/actions";
import { signedPost, clearWalletAuthSession, type WalletAuthAction } from "./signedPost";

const mockGet = apiClient.get as unknown as ReturnType<typeof vi.fn>;
const mockPost = apiClient.post as unknown as ReturnType<typeof vi.fn>;
const mockSign = signTypedData as unknown as ReturnType<typeof vi.fn>;
const mockChainId = getChainId as unknown as ReturnType<typeof vi.fn>;

const WALLET = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266";

beforeEach(() => {
  mockGet.mockReset();
  mockPost.mockReset();
  mockSign.mockReset();
  mockChainId.mockReset();
  clearWalletAuthSession();

  // 后端一次性 nonce 台账（Consume 消费式）：每次取新 nonce。
  let n = 0;
  mockGet.mockImplementation(() => Promise.resolve({ nonce: `nonce-${++n}` }));
  mockSign.mockResolvedValue("0xsignature");
  mockPost.mockResolvedValue({ ok: true });
  mockChainId.mockReturnValue(31337);
});

describe("signedPost — WalletAuth EIP-712（AUTH-01/02/03）", () => {
  it("AUTH-01: 先 GET nonce → signTypedData(EIP-712 域/字段对齐后端) → 带签名头 POST", async () => {
    await signedPost("/api/deposit", { wallet: WALLET, amount: "1500000" }, {
      wallet: WALLET,
      action: "deposit",
    });

    // ① 先取 nonce（后端 GetWalletNonce 路径 /api/auth/nonce/:wallet）
    expect(mockGet).toHaveBeenCalledTimes(1);
    expect(mockGet.mock.calls[0][0]).toBe(`/api/auth/nonce/${WALLET}`);

    // ② EIP-712 typedData 与后端 WalletAuthDigest 逐项对齐
    expect(mockSign).toHaveBeenCalledTimes(1);
    const typed = mockSign.mock.calls[0][1];
    // domain：name=LinkWorld / version=1 / chainId=链上当前 chainId（无 verifyingContract）
    expect(typed.domain).toEqual({ name: "LinkWorld", version: "1", chainId: 31337 });
    // 主类型 + 字段名/类型/顺序：wallet(address) nonce(string) action(string)
    expect(typed.primaryType).toBe("WalletAuth");
    expect(typed.types.WalletAuth).toEqual([
      { name: "wallet", type: "address" },
      { name: "nonce", type: "string" },
      { name: "action", type: "string" },
    ]);
    // message：wallet/nonce(string)/action
    expect(typed.message).toEqual({
      wallet: WALLET,
      nonce: "nonce-1",
      action: "deposit",
    });
    // signTypedData 第一个入参是 wagmiConfig（core action，非 hook）
    expect(mockSign.mock.calls[0][0]).toEqual({ __mock: "config" });

    // ③ 带签名头 POST：header 名对齐后端 middleware 常量
    expect(mockPost).toHaveBeenCalledTimes(1);
    const [path, body, cfg] = mockPost.mock.calls[0];
    expect(path).toBe("/api/deposit");
    expect(body).toEqual({ wallet: WALLET, amount: "1500000" });
    expect(cfg.headers).toEqual({
      "X-Wallet-Address": WALLET,
      "X-Wallet-Nonce": "nonce-1",
      "X-Wallet-Action": "deposit",
      "X-Wallet-Signature": "0xsignature",
    });
  });

  it("AUTH-02: action 绑定正确（withdraw/bills/pay/service 各端点 action 对齐后端）", async () => {
    const cases: Array<[string, WalletAuthAction]> = [
      ["/api/withdraw", "withdraw"],
      ["/api/bills/pay", "bills/pay"],
      ["/api/service/activate", "service/activate"],
      ["/api/service/deactivate", "service/deactivate"],
    ];
    for (const [path, action] of cases) {
      mockSign.mockClear();
      mockPost.mockClear();
      await signedPost(path, { wallet: WALLET }, { wallet: WALLET, action });
      expect(mockSign.mock.calls[0][1].message.action).toBe(action);
      const cfg = mockPost.mock.calls[0][2];
      expect(cfg.headers["X-Wallet-Action"]).toBe(action);
    }
  });

  it("AUTH-03: 拒签 → 不发 POST，抛 WalletAuthRejectedError（不进 pending）", async () => {
    // viem UserRejectedRequestError 携带 name=UserRejectedRequestError / code 4001
    const rejected = Object.assign(new Error("User rejected the request."), {
      name: "UserRejectedRequestError",
      code: 4001,
    });
    mockSign.mockRejectedValueOnce(rejected);

    await expect(
      signedPost("/api/deposit", { wallet: WALLET }, { wallet: WALLET, action: "deposit" }),
    ).rejects.toMatchObject({ rejected: true });

    expect(mockPost).not.toHaveBeenCalled();
  });

  it("一次性 nonce：同会话每次写都取新 nonce 重签（后端 Consume 消费式，不复用同 nonce）", async () => {
    await signedPost("/api/deposit", { wallet: WALLET }, { wallet: WALLET, action: "deposit" });
    await signedPost("/api/withdraw", { wallet: WALLET }, { wallet: WALLET, action: "withdraw" });

    expect(mockGet).toHaveBeenCalledTimes(2);
    expect(mockSign.mock.calls[0][1].message.nonce).toBe("nonce-1");
    expect(mockSign.mock.calls[1][1].message.nonce).toBe("nonce-2");
    // header 里 nonce 不复用
    expect(mockPost.mock.calls[0][2].headers["X-Wallet-Nonce"]).toBe("nonce-1");
    expect(mockPost.mock.calls[1][2].headers["X-Wallet-Nonce"]).toBe("nonce-2");
  });
});
