import { parseEther } from "viem";
import { apiClient } from "./client";
import type { Region, Operator } from "../../types";

interface ApiOperator {
  id: number;
  name: string;
  region: string;
  country_code: string;
  required_deposit: string;
  is_active: boolean;
}

interface ApiCountry {
  code: string;
  name: string;
  prefix: string;
}

// 国家 flag emoji map
const FLAG_MAP: Record<string, string> = {
  US: "\u{1F1FA}\u{1F1F8}",
  GB: "\u{1F1EC}\u{1F1E7}",
  FR: "\u{1F1EB}\u{1F1F7}",
  RU: "\u{1F1F7}\u{1F1FA}",
  JP: "\u{1F1EF}\u{1F1F5}",
  VN: "\u{1F1FB}\u{1F1F3}",
  LA: "\u{1F1F1}\u{1F1E6}",
  KH: "\u{1F1F0}\u{1F1ED}",
  TH: "\u{1F1F9}\u{1F1ED}",
  MY: "\u{1F1F2}\u{1F1FE}",
  PH: "\u{1F1F5}\u{1F1ED}",
};

function toOperator(api: ApiOperator): Operator {
  return {
    id: String(api.id),
    name: api.name,
    region: api.region,
    requiredDeposit: parseEther(api.required_deposit),
    dataRate: 0.05, // 暂用默认值，后端未返回
    callRate: 0.02,
    isActive: api.is_active,
  };
}

export const operatorApi = {
  async getRegions(): Promise<Region[]> {
    const [operators, countries] = await Promise.all([
      apiClient.get<any, ApiOperator[]>("/api/operators"),
      apiClient.get<any, ApiCountry[]>("/api/countries"),
    ]);

    // groupBy country_code
    const grouped = new Map<string, ApiOperator[]>();
    for (const op of operators) {
      const code = op.country_code;
      if (!grouped.has(code)) grouped.set(code, []);
      grouped.get(code)!.push(op);
    }

    const countryMap = new Map(countries.map((c) => [c.code, c]));

    return Array.from(grouped.entries()).map(([code, ops]) => {
      const country = countryMap.get(code);
      const minDeposit = Math.min(
        ...ops.map((o) => parseFloat(o.required_deposit)),
      );
      return {
        code,
        name: country?.name || ops[0].region,
        flag: FLAG_MAP[code] || "\u{1F30D}",
        operatorCount: ops.length,
        startingPrice: minDeposit,
      };
    });
  },

  async getOperatorsByRegion(regionCode: string): Promise<Operator[]> {
    const operators =
      await apiClient.get<any, ApiOperator[]>("/api/operators");
    return operators
      .filter((op) => op.country_code === regionCode)
      .map(toOperator);
  },
};
