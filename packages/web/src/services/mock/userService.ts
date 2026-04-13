import type { User } from "@/types";
import { delay } from "./delay";
import { mockUser, setMockUser } from "./data";

export const userService = {
  async getUserProfile(address: string): Promise<User | null> {
    await delay();
    return mockUser && mockUser.address === address ? { ...mockUser } : null;
  },
  async register(address: string, email: string): Promise<User> {
    await delay(800);
    const user: User = { address, email, status: "inactive", registeredAt: new Date().toISOString(), nftTokenId: Math.floor(Math.random() * 10000) };
    setMockUser(user);
    return { ...user };
  },
  async sendVerificationCode(_address: string, _email: string): Promise<{ success: boolean }> {
    await delay(500);
    return { success: true };
  },
  async verifyEmail(_address: string, _code: string): Promise<boolean> {
    await delay(500);
    return true;
  },
};
