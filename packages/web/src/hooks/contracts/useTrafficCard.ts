import { useCallback, useEffect, useState } from "react";
import {
  useReadContracts,
  useWriteContract,
  useWaitForTransactionReceipt,
  useChainId,
  usePublicClient,
} from "wagmi";
import { parseAbiItem, type Log } from "viem";
import { TrafficCardNFTABI } from "../../config/abis";
import { getContractAddress } from "../../config/contracts";

const CARD_MINTED_EVENT = parseAbiItem(
  "event CardMinted(address indexed user, uint256 tokenId, uint256 dataAmount)"
);

// getLogs 限窗（design §9 / B6）：现状 fromBlock:0n 全量扫块在 Arbitrum 公共 RPC 必限流/超时。
// 选型：**限定 fromBlock 窗口**（后端 NFT 列表端点未在 handoff 提供契约，故选限窗方案）。
// fromBlock = max(0, latestBlock - LOG_WINDOW_BLOCKS)；本地 31337 块少，窗口足够覆盖；
// Arbitrum Sepolia ~0.25s/块 → 5_000_000 块 ≈ 14 天，覆盖近期发卡且不全量扫历史。
const LOG_WINDOW_BLOCKS = 5_000_000n;

export type TrafficCardItem = {
  tokenId: bigint;
  dataAmount: bigint;
  createdAt: bigint;
  isDestroyed: boolean;
};

export function useTrafficCards(address: `0x${string}` | undefined) {
  const chainId = useChainId();
  const publicClient = usePublicClient();
  const [tokenIds, setTokenIds] = useState<bigint[]>([]);
  const [isLoadingEvents, setIsLoadingEvents] = useState(false);
  // ★ B6：getLogs 失败 → error 态（非静默置空）。区分「加载失败」vs「真无卡」。
  const [logsError, setLogsError] = useState<Error | null>(null);
  const [reloadKey, setReloadKey] = useState(0);

  const contractAddress = address
    ? (() => {
        try {
          return getContractAddress(chainId, "TrafficCardNFT");
        } catch {
          return undefined;
        }
      })()
    : undefined;

  useEffect(() => {
    let cancelled = false;
    async function load() {
      if (!publicClient || !address || !contractAddress) {
        setTokenIds([]);
        setLogsError(null);
        return;
      }
      setIsLoadingEvents(true);
      setLogsError(null);
      try {
        // 限窗：从 latest - LOG_WINDOW_BLOCKS 起扫，避免 fromBlock:0n 全量扫块限流/超时。
        const latest = await publicClient.getBlockNumber();
        const fromBlock =
          latest > LOG_WINDOW_BLOCKS ? latest - LOG_WINDOW_BLOCKS : 0n;
        const logs = (await (publicClient.getLogs as (a: unknown) => Promise<Log[]>)({
          address: contractAddress,
          event: CARD_MINTED_EVENT,
          args: { user: address },
          fromBlock,
          toBlock: "latest",
        })) as Log[];
        const ids = logs
          .map((l) => {
            const args = (l as unknown as { args?: { tokenId?: bigint } }).args;
            return args?.tokenId;
          })
          .filter((v): v is bigint => typeof v === "bigint");
        if (!cancelled) {
          setTokenIds(ids);
          setLogsError(null);
        }
      } catch (err) {
        // ★ 禁静默置空：标记 error 态，让 UI 区分「加载失败，重试」vs「真无卡」。
        if (!cancelled) {
          setLogsError(err instanceof Error ? err : new Error(String(err)));
        }
        // eslint-disable-next-line no-console
        console.warn("useTrafficCards: getLogs failed", err);
      } finally {
        if (!cancelled) setIsLoadingEvents(false);
      }
    }
    load();
    return () => {
      cancelled = true;
    };
  }, [publicClient, address, contractAddress, reloadKey]);

  const infoQ = useReadContracts({
    // 注：此 `as never` 非 T1 占位（占位 = functionName/args 指向不存在函数，已全清），
    // 而是 wagmi 对**动态长度**合约数组的已知类型限制（无法静态推断 tuple 长度）。functionName/args 真实有效。
    contracts: tokenIds.map((id) => ({
      address: contractAddress,
      abi: TrafficCardNFTABI,
      functionName: "getCardInfo",
      args: [id],
    })) as never,
    query: { enabled: tokenIds.length > 0 && !!contractAddress },
  });

  const cards: TrafficCardItem[] = tokenIds
    .map((tokenId, idx) => {
      const r = infoQ.data?.[idx] as
        | { status: string; result?: unknown }
        | undefined;
      if (!r || r.status !== "success") return null;
      // 新 getCardInfo 返回 { dataAmount, createdAt, isDestroyed }（无 destroyedAt）。
      const info = r.result as {
        dataAmount: bigint;
        createdAt: bigint;
        isDestroyed: boolean;
      };
      return {
        tokenId,
        dataAmount: info.dataAmount,
        createdAt: info.createdAt,
        isDestroyed: info.isDestroyed,
      };
    })
    .filter((c): c is TrafficCardItem => c !== null)
    .filter((c) => !c.isDestroyed)
    .sort((a, b) => (a.createdAt > b.createdAt ? -1 : 1));

  const refetch = useCallback(() => {
    setReloadKey((k) => k + 1);
    infoQ.refetch();
  }, [infoQ]);

  return {
    cards,
    isLoading: isLoadingEvents || infoQ.isLoading,
    // 加载失败（getLogs 限流/超时 或 逐卡读链失败）；与 cards.length===0 区分「真无卡」。
    isError: !!logsError || !!infoQ.isError,
    error: logsError ?? (infoQ.error as Error | null) ?? null,
    refetch,
  };
}

/**
 * 可用流量额度（逐卡聚合模型）。
 * 新 TrafficCardNFT 无 getAvailableCredit/getCreditExpiry（旧聚合额度模型已废弃），
 * T4 重写为对未销毁卡的 dataAmount 求和。无链上聚合到期字段 → expiry 恒 0n（按卡逐张计）。
 * 返回形状与旧接口一致，Cards 页（T8）无需改结构。
 */
export function useTrafficCardCredit(address: `0x${string}` | undefined) {
  const { cards, isLoading, isError, error, refetch } = useTrafficCards(address);
  const balance = cards.reduce((sum, c) => sum + c.dataAmount, 0n);
  return {
    balance,
    expiry: 0n,
    isExpired: false,
    isLoading,
    isError,
    error,
    refetch,
  };
}

/**
 * 批量销毁流量卡兑换 SIM 天数（一步式玩法）。
 * 调 redeemForSim(tokenIds[])：每张卡 = 1 天，SIM 天数 = 销毁卡数。
 * 单张 burn 仍在合约内（语义为兑换），但前端多选流程统一走 redeemForSim。
 */
export function useRedeemForSim() {
  const chainId = useChainId();
  const { writeContract, data: hash, isPending, error, reset } = useWriteContract();
  const { isPending: rawConfirming, isSuccess } = useWaitForTransactionReceipt({ hash });
  const isConfirming = !!hash && rawConfirming;

  const redeem = (tokenIds: bigint[]) => {
    writeContract({
      address: getContractAddress(chainId, "TrafficCardNFT"),
      abi: TrafficCardNFTABI,
      functionName: "redeemForSim",
      args: [tokenIds],
    });
  };

  return { redeem, hash, isPending, isConfirming, isSuccess, error, reset };
}
