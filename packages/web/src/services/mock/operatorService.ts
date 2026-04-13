import type { Region, Operator, VirtualNumber } from "@/types";
import { delay } from "./delay";
import { mockRegions, mockOperators, mockNumbers, addMockNumber } from "./data";

export const operatorService = {
  async getRegions(): Promise<Region[]> {
    await delay();
    return [...mockRegions];
  },
  async getOperatorsByRegion(regionCode: string): Promise<Operator[]> {
    await delay();
    return mockOperators[regionCode] || [];
  },
  async getMyNumbers(_address: string): Promise<VirtualNumber[]> {
    await delay();
    return [...mockNumbers];
  },
  async applyNumber(_address: string, operatorId: string): Promise<VirtualNumber> {
    await delay(1200);
    const allOps = Object.values(mockOperators).flat();
    const op = allOps.find((o) => o.id === operatorId);
    if (!op) throw new Error(`Operator ${operatorId} not found`);
    const region = mockRegions.find((r) => r.code === op.region)!;
    const num: VirtualNumber = {
      id: `vn-${Date.now()}`,
      number: `+${Math.floor(Math.random() * 90 + 10)} ${Math.floor(Math.random() * 9000 + 1000)}-${Math.floor(Math.random() * 9000 + 1000)}`,
      region: region.code, operator: op.name, status: "active", activatedAt: new Date().toISOString(),
      credentials: { type: "esim", config: `LPA:1$rsp.${op.name.toLowerCase().replace(/\s/g, "")}.com$${Math.random().toString(36).slice(2, 18).toUpperCase()}` },
    };
    addMockNumber(num);
    return num;
  },
};
