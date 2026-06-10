import { useEffect, useState } from "react";
import {
  useReadContract,
  useReadContracts,
  useWriteContract,
  useWaitForTransactionReceipt,
  useChainId,
  usePublicClient,
} from "wagmi";
import { parseAbiItem, type Log } from "viem";
import { TrafficCardNFTABI, DepositABI } from "../../config/abis";
import { getContractAddress } from "../../config/contracts";

const CARD_MINTED_EVENT = parseAbiItem(
  "event CardMinted(address indexed user, uint256 tokenId, uint256 dataAmount)"
);

export function useTrafficCardCredit(address: `0x${string}` | undefined) {
  const chainId = useChainId();
  const enabled = !!address;
  const contractAddress = enabled
    ? getContractAddress(chainId, "TrafficCardNFT")
    : undefined;

  // T4/Cards-rework: 新 TrafficCardNFT 不再暴露 getAvailableCredit/getCreditExpiry（旧聚合额度
  // 模型），改为 getCardInfo/getUserCardCount 逐卡模型。本轮仅 as never 保编译，credit 读链路
  // 重写 + Cards 双 Tab 改造归后续 T。
  const balanceQ = useReadContract({
    address: contractAddress,
    abi: TrafficCardNFTABI,
    functionName: "getAvailableCredit" as never,
    args: (address ? [address] : undefined) as never,
    query: { enabled, staleTime: 15_000 },
  });

  const expiryQ = useReadContract({
    address: contractAddress,
    abi: TrafficCardNFTABI,
    functionName: "getCreditExpiry" as never,
    args: (address ? [address] : undefined) as never,
    query: { enabled, staleTime: 15_000 },
  });

  const balance = (balanceQ.data as bigint | undefined) ?? 0n;
  const expiry = (expiryQ.data as bigint | undefined) ?? 0n;
  const nowSec = BigInt(Math.floor(Date.now() / 1000));
  const isExpired = expiry > 0n && expiry < nowSec;

  return {
    balance,
    expiry,
    isExpired,
    isLoading: balanceQ.isLoading || expiryQ.isLoading,
    refetch: () => {
      balanceQ.refetch();
      expiryQ.refetch();
    },
  };
}

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
        return;
      }
      setIsLoadingEvents(true);
      try {
        const logs = (await (publicClient.getLogs as (a: unknown) => Promise<Log[]>)({
          address: contractAddress,
          event: CARD_MINTED_EVENT,
          args: { user: address },
          fromBlock: 0n,
          toBlock: "latest",
        })) as Log[];
        const ids = logs
          .map((l) => {
            const args = (l as unknown as { args?: { tokenId?: bigint } }).args;
            return args?.tokenId;
          })
          .filter((v): v is bigint => typeof v === "bigint");
        if (!cancelled) setTokenIds(ids);
      } catch (err) {
        if (!cancelled) setTokenIds([]);
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
      const info = r.result as {
        dataAmount: bigint;
        createdAt: bigint;
        isDestroyed: boolean;
        destroyedAt: bigint;
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

  return {
    cards,
    isLoading: isLoadingEvents || infoQ.isLoading,
    refetch: () => {
      setReloadKey((k) => k + 1);
      infoQ.refetch();
    },
  };
}

export function useBurnCard() {
  const chainId = useChainId();
  const { writeContract, data: hash, isPending, error, reset } = useWriteContract();
  const { isPending: rawConfirming, isSuccess } = useWaitForTransactionReceipt({ hash });
  const isConfirming = !!hash && rawConfirming;

  const burnCard = (tokenId: bigint) => {
    writeContract({
      address: getContractAddress(chainId, "TrafficCardNFT"),
      abi: TrafficCardNFTABI,
      functionName: "burn",
      args: [tokenId],
    });
  };

  return { burnCard, hash, isPending, isConfirming, isSuccess, error, reset };
}

export function useIssueMonthlyCards() {
  const chainId = useChainId();
  const { writeContract, data: hash, isPending, error, reset } = useWriteContract();
  const { isPending: rawConfirming, isSuccess } = useWaitForTransactionReceipt({ hash });
  const isConfirming = !!hash && rawConfirming;

  const issue = (users: `0x${string}`[]) => {
    writeContract({
      address: getContractAddress(chainId, "Deposit"),
      abi: DepositABI,
      functionName: "issueMonthlyTrafficCards",
      args: [users],
    });
  };

  return { issue, hash, isPending, isConfirming, isSuccess, error, reset };
}
