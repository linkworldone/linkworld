import { Injectable } from "@nestjs/common";

@Injectable()
export class AuthService {
  private codes = new Map<string, string>();

  async sendVerificationCode(email: string) {
    const code = Math.floor(100000 + Math.random() * 900000).toString();
    this.codes.set(email, code);
    // TODO: 接入 nodemailer 发送邮件
    return { message: "Code sent" };
  }

  async verifyCode(email: string, code: string) {
    const stored = this.codes.get(email);
    if (!stored || stored !== code) {
      return { success: false, message: "Invalid code" };
    }
    this.codes.delete(email);
    // TODO: 生成 JWT token
    return { success: true, message: "Verified" };
  }
}
