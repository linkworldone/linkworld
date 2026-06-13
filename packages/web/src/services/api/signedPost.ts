import { signTypedData, getChainId } from "wagmi/actions";
import { apiClient } from "./client";
import { wagmiConfig } from "../../config/wagmi";

// T5 WalletAuth — 写端点钱包身份签名（design §3.7 / arch-review 🔴 N1）。
//
// signedPost = ① GET nonce → ② EIP-712 signTypedData → ③ 带签名头 POST。
// EIP-712 域/字段/action/header **逐字节对齐后端**
// packages/backend/internal/middleware/middleware.go（WalletAuthDigest / NewWalletAuth）：
//   domain  : { name:"LinkWorld", version:"1", chainId }（仅 3 字段，无 verifyingContract）
//   types   : WalletAuth = [ wallet:address, nonce:string, action:string ]（顺序固定）
//   message : { wallet, nonce(string), action }
//   headers : X-Wallet-Address / X-Wallet-Nonce / X-Wallet-Action / X-Wallet-Signature
// chainId 绑定 domain（防跨链重放），必须 == 后端 walletAuthChainID。
//
// 非全局拦截器（design §3.7 红线）：签名经 wagmi **core action** signTypedData(config,…) 完成，
// **不在 axios 拦截器里调 React hook**（拦截器内调 hook 必崩）。core action 仅需 wagmiConfig，
// 可在任意模块（含非组件）调用。

// 后端鉴权头常量（对齐 middleware.go Header* 常量）。
const HEADER = {
  address: "X-Wallet-Address",
  nonce: "X-Wallet-Nonce",
  action: "X-Wallet-Action",
  signature: "X-Wallet-Signature",
} as const;

// EIP-712 domain 常量（对齐 middleware.go eip712Domain*）。
const EIP712_DOMAIN_NAME = "LinkWorld";
const EIP712_DOMAIN_VERSION = "1";

// 写端点 action 取值（对齐 main.go walletAuth("…") 绑定）。
export type WalletAuthAction =
  | "deposit"
  | "withdraw"
  | "bills/pay"
  | "service/activate"
  | "service/deactivate"
  | "sim/claim";

export interface SignedPostOpts {
  /** 发起写操作的钱包地址（签名地址必须 == 该地址，后端 ecrecover 绑 msg.sender 语义）。 */
  wallet: string;
  /** 本端点绑定的语义动作；后端 NewWalletAuth 校验 action == expectedAction。 */
  action: WalletAuthAction;
}

/** 用户在钱包里拒签身份签名时抛出。调用方据此 toast「身份签名被取消」且不进 pending。 */
export class WalletAuthRejectedError extends Error {
  readonly rejected = true as const;
  constructor() {
    // 兜底英文 message；真正展示给用户的地方用 t("errors.authCancelled") 覆盖。
    super("Signature cancelled, action not submitted");
    this.name = "WalletAuthRejectedError";
  }
}

// viem 拒签错误识别：UserRejectedRequestError（name）或 EIP-1193 code 4001。
function isUserRejection(err: unknown): boolean {
  if (!err || typeof err !== "object") return false;
  const e = err as { name?: string; code?: number; cause?: unknown };
  if (e.name === "UserRejectedRequestError" || e.code === 4001) return true;
  // viem 常把底层 provider 错误包在 cause 里。
  if (e.cause && e.cause !== err) return isUserRejection(e.cause);
  return false;
}

interface NonceResponse {
  nonce: string;
}

// 取后端一次性 nonce（GetWalletNonce 路径 /api/auth/nonce/:wallet，main.go）。
async function fetchNonce(wallet: string): Promise<string> {
  const data = await apiClient.get<unknown, NonceResponse>(
    `/api/auth/nonce/${wallet}`,
  );
  return data.nonce;
}

// 会话级签名（design §3.7 B1）说明：
// 后端 nonce 台账是**一次性消费式**（middleware.NewWalletAuth → nonceRepo.Consume，签过即作废）。
// 在该模型下「会话级签名一次」退化为「每次写操作取新 nonce 重签」——不能跨写复用同一 nonce
//（复用必被 Consume 拒为 used）。故 signedPost 每次取新 nonce 现签现用，不缓存已签 nonce。
// clearWalletAuthSession 预留给「换钱包/登出」语义钩子（当前无跨调用缓存，调用即 no-op，
// 但保留显式清理入口，便于后续若后端改为会话式 nonce 时落缓存而不改调用面）。
export function clearWalletAuthSession(): void {
  // 一次性 nonce 模型下无跨调用缓存可清；保留显式入口（换钱包/登出钩子）。
}

/**
 * 带 WalletAuth 身份签名的 POST：取 nonce → EIP-712 签名 → 带签名头调 axios POST。
 * 读端点不用本函数（保持裸 apiClient.get/post）。
 *
 * @throws WalletAuthRejectedError 用户拒签（不发 POST，调用方不应进 pending）。
 */
export async function signedPost<T = unknown>(
  path: string,
  body: unknown,
  { wallet, action }: SignedPostOpts,
): Promise<T> {
  const nonce = await fetchNonce(wallet);
  const chainId = getChainId(wagmiConfig);

  let signature: string;
  try {
    signature = await signTypedData(wagmiConfig, {
      domain: {
        name: EIP712_DOMAIN_NAME,
        version: EIP712_DOMAIN_VERSION,
        chainId,
      },
      types: {
        WalletAuth: [
          { name: "wallet", type: "address" },
          { name: "nonce", type: "string" },
          { name: "action", type: "string" },
        ],
      },
      primaryType: "WalletAuth",
      message: {
        wallet: wallet as `0x${string}`,
        nonce,
        action,
      },
    });
  } catch (err) {
    if (isUserRejection(err)) throw new WalletAuthRejectedError();
    throw err;
  }

  return apiClient.post<unknown, T>(path, body, {
    headers: {
      [HEADER.address]: wallet,
      [HEADER.nonce]: nonce,
      [HEADER.action]: action,
      [HEADER.signature]: signature,
    },
  });
}
