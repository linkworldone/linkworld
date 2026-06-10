import { useReadContract, useReadContracts, useChainId } from "wagmi";
import { ServiceManagerABI } from "../../config/abis";
import { getContractAddress } from "../../config/contracts";

/**
 * 新 ServiceManager 是**运营商注册中心模型**（getOperator/getActiveOperators/
 * getOperatorsByCountry/addOperator…），不再有旧 0G 的 getUserService/activateService
 * 用户服务态。本轮（T4）重写 service 读链路，去掉占位 `as never`。
 *
 * 链上 Operator 结构（getOperator 返回 tuple）：
 *   { id, name, region, countryCode, requiredDeposit, isActive, paymentAddress }
 */

export interface OnChainOperator {
  id: bigint;
  name: string;
  region: string;
  countryCode: string;
  requiredDeposit: bigint;
  isActive: boolean;
  paymentAddress: `0x${string}`;
}

/** 读单个运营商详情（getOperator(operatorId)）。 */
export function useContractOperator(operatorId: bigint | undefined) {
  const chainId = useChainId();
  const enabled = operatorId !== undefined;
  const q = useReadContract({
    address: getContractAddress(chainId, "ServiceManager"),
    abi: ServiceManagerABI,
    functionName: "getOperator",
    args: enabled ? [operatorId] : undefined,
    query: { enabled, staleTime: 60_000 },
  });
  return {
    ...q,
    operator: q.data as unknown as OnChainOperator | undefined,
  };
}

/** 读活跃运营商 id 列表（getActiveOperators() → uint256[]）。 */
export function useContractActiveOperators() {
  const chainId = useChainId();
  const q = useReadContract({
    address: getContractAddress(chainId, "ServiceManager"),
    abi: ServiceManagerABI,
    functionName: "getActiveOperators",
    query: { staleTime: 60_000 },
  });
  return {
    ...q,
    operatorIds: (q.data as readonly bigint[] | undefined) ?? [],
  };
}

/** 批量读多个运营商详情（按 id 列表逐个 getOperator）。 */
export function useContractOperators(operatorIds: readonly bigint[]) {
  const chainId = useChainId();
  const address = (() => {
    try {
      return getContractAddress(chainId, "ServiceManager");
    } catch {
      return undefined;
    }
  })();
  const q = useReadContracts({
    // 注：`as never` 是 wagmi 对动态长度合约数组的已知类型限制（非占位）；functionName/args 真实有效。
    contracts: operatorIds.map((id) => ({
      address,
      abi: ServiceManagerABI,
      functionName: "getOperator",
      args: [id],
    })) as never,
    query: { enabled: operatorIds.length > 0 && !!address },
  });
  const results = (q.data ?? []) as ReadonlyArray<{
    status: string;
    result?: unknown;
  }>;
  const operators: OnChainOperator[] = results
    .map((r) =>
      r && r.status === "success"
        ? (r.result as unknown as OnChainOperator)
        : null
    )
    .filter((o): o is OnChainOperator => o !== null);
  return { ...q, operators };
}
