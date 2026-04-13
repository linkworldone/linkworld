import { MOCK_DELAY_MS } from "@/config/constants";

export function delay(ms: number = MOCK_DELAY_MS): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, ms));
}
