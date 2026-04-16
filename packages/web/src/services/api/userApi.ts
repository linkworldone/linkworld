import { apiClient } from "./client";
import type { User } from "../../types";

// 后端 response 类型
interface ApiUser {
  id: number;
  wallet_addr: string;
  email: string;
  token_id: number;
  is_active: boolean;
  registered_at: string;
  created_at: string;
  updated_at: string;
}

function toUser(api: ApiUser): User {
  return {
    address: api.wallet_addr,
    email: api.email,
    status: api.is_active ? "active" : "inactive",
    registeredAt: api.registered_at,
    nftTokenId: api.token_id,
  };
}

export const userApi = {
  async getUser(wallet: string): Promise<User | null> {
    try {
      const data = await apiClient.get<any, ApiUser>(`/api/user/${wallet}`);
      return toUser(data);
    } catch {
      return null;
    }
  },

  async register(
    wallet: string,
    email: string,
    tokenId?: number,
  ): Promise<void> {
    await apiClient.post("/api/register", {
      wallet,
      email,
      token_id: tokenId || 0,
    });
  },
};
