import { useReadContract, useChainId } from "wagmi";
import { FeeManagerABI } from "../../config/abis";
import { getContractAddress } from "../../config/contracts";

// FeeManager 手续费读链（design §3.6 / D12）：删写死 PLATFORM_FEE_RATE，费率/费额一律实读链。
// 基点制：getFeeRate() 返回基点（150 = 1.5%），分母 10000（链上 FEE_DENOMINATOR）。
const FEE_DENOMINATOR = 10_000n;

function safeFeeManager(chainId: number): `0x${string}` | undefined {
  try {
    return getContractAddress(chainId, "FeeManager");
  } catch {
    // 未部署链（零地址占位）→ 不读，交由调用方按 isError/「--」兜底。
    return undefined;
  }
}

export interface FeeRateResult {
  /** 原始基点（如 150）。读链未回 / 失败时为 undefined。 */
  bps?: bigint;
  /** 百分比数值（如 1.5）。读链未回 / 失败时为 undefined。 */
  percent?: number;
  /** 展示用百分比字符串（如 "1.5%"）。读链未回 / 失败时为 undefined（调用方兜底 "--"）。 */
  label?: string;
  isLoading: boolean;
  isError: boolean;
}

/**
 * 读链平台手续费率 getFeeRate()（基点，150=1.5%，/10000）。
 * **铁律**：禁写死兜底——loading→skeleton、失败→「--」由调用方处理（design §3.6）。
 */
export function useFeeRate(): FeeRateResult {
  const chainId = useChainId();
  const address = safeFeeManager(chainId);
  const enabled = !!address;
  const q = useReadContract({
    address,
    abi: FeeManagerABI,
    functionName: "getFeeRate",
    query: { enabled, staleTime: 60_000 },
  });

  const bps = q.data as bigint | undefined;
  const percent =
    bps !== undefined ? Number((bps * 100n * 1000n) / FEE_DENOMINATOR) / 1000 : undefined;
  const label = percent !== undefined ? `${formatPercent(percent)}%` : undefined;

  return {
    bps,
    percent,
    label,
    // 未部署链视为读链失败（isError），不当 loading 一直转圈。
    isLoading: enabled && q.isLoading,
    isError: !enabled || q.isError,
  };
}

/**
 * 读链精确费额 calculateFee(amount)（直读合约，**不自算**，design §3.6 / B 红线）。
 * amount<=0 或未部署链 → 不读。
 */
export function useCalculateFee(amount?: bigint): {
  fee?: bigint;
  isLoading: boolean;
  isError: boolean;
} {
  const chainId = useChainId();
  const address = safeFeeManager(chainId);
  const enabled = !!address && amount !== undefined && amount > 0n;
  const q = useReadContract({
    address,
    abi: FeeManagerABI,
    functionName: "calculateFee",
    args: enabled ? [amount] : undefined,
    query: { enabled, staleTime: 60_000 },
  });

  return {
    fee: q.data as bigint | undefined,
    isLoading: enabled && q.isLoading,
    isError: enabled && q.isError,
  };
}

// 1.5 → "1.5"，2 → "2"（去尾零，最多 3 位小数）。
function formatPercent(p: number): string {
  return p
    .toFixed(3)
    .replace(/\.?0+$/, "");
}
