import { signedPost } from "./signedPost";

// 绑定邮箱（两步：发码 + 验证）。两端点均走 WalletAuth 身份签名，
// 与 depositApi 写端点同模式：signedPost 取 nonce → EIP-712 签名 → 带签名头 POST。
// 后端契约：
//   POST /api/email/bind   action=email/bind   body { email }          200 {message} / 429 {error} / 400 {error}
//   POST /api/email/verify action=email/verify body { email, code(6位) } 200 {message} / 400 {error}
// 错误经 client 拦截器抛 Error（带 .status，调用方据 429 区分限流提示）。
export const emailApi = {
  /** 发送验证码到邮箱（同钱包 60s 限流 → 429）。拒签抛 WalletAuthRejectedError。 */
  async bindEmail(wallet: string, email: string): Promise<void> {
    await signedPost(
      "/api/email/bind",
      { email },
      { wallet, action: "email/bind" },
    );
  },

  /** 校验 6 位验证码完成绑定。验证码错误/过期 → 400 {error}。 */
  async verifyEmail(wallet: string, email: string, code: string): Promise<void> {
    await signedPost(
      "/api/email/verify",
      { email, code },
      { wallet, action: "email/verify" },
    );
  },
};
