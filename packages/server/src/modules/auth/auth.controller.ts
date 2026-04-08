import { Controller, Post, Body } from "@nestjs/common";
import { AuthService } from "./auth.service";

@Controller("auth")
export class AuthController {
  constructor(private readonly authService: AuthService) {}

  @Post("send-code")
  async sendCode(@Body("email") email: string) {
    return this.authService.sendVerificationCode(email);
  }

  @Post("verify")
  async verify(
    @Body("email") email: string,
    @Body("code") code: string
  ) {
    return this.authService.verifyCode(email, code);
  }
}
