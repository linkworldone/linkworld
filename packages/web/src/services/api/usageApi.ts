import { apiClient } from "./client";
import type { MonthEstimate } from "../../types";

interface ApiUsage {
  data_usage: number;
  call_usage: number;
  signature: string;
}

export const usageApi = {
  async getUsage(wallet: string): Promise<MonthEstimate> {
    const data = await apiClient.get<any, ApiUsage>(
      `/api/usage/${wallet}`,
    );
    // data_usage 单位是 MB，转 GB
    const dataGB = data.data_usage / 1000;
    const callMinutes = data.call_usage;
    // 估算费用（简化）
    const estimatedCost = (dataGB * 0.05 + callMinutes * 0.02).toFixed(2);
    return {
      dataGB: parseFloat(dataGB.toFixed(2)),
      callMinutes,
      estimatedCost,
    };
  },
};
