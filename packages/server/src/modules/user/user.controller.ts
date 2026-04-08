import { Controller, Get, Param } from "@nestjs/common";
import { UserService } from "./user.service";

@Controller("users")
export class UserController {
  constructor(private readonly userService: UserService) {}

  @Get(":address")
  async getUser(@Param("address") address: string) {
    return this.userService.getUserByAddress(address);
  }
}
