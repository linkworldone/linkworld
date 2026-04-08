import { Injectable } from "@nestjs/common";

@Injectable()
export class UserService {
  async getUserByAddress(address: string) {
    // TODO: 从链上读取用户信息
    return { address, status: "pending" };
  }
}
