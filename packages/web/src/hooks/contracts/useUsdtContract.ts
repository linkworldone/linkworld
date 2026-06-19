import {
  useReadContract,
  useWriteContract,
  useWaitForTransactionReceipt,
  useChainId,
} from "wagmi";
import { MockUSDTABI } from "../../config/abis";
import { getUsdt } from "../../config/contracts";
import { localGasOverride } from "./gasOverride";

/**
 * 读 USDT(ERC20) 对某 spender 的授权额度。
 * 用于 TwoStepAction 判断「allowance ≥ 需求额 → 跳过 Approve」。
 */
export function useAllowance(
  owner: `0x${string}` | undefined,
  spender: `0x${string}` | undefined
) {
  const chainId = useChainId();
  const enabled = !!owner && !!spender;
  return useReadContract({
    address: enabled ? getUsdt(chainId) : undefined,
    abi: MockUSDTABI,
    functionName: "allowance",
    args: enabled ? [owner, spender] : undefined,
    // 授权额度对资损敏感且会被自身 approve 改写，短 staleTime 让 Approve 成功后能尽快重读到新值。
    query: { enabled, staleTime: 2_000 },
  });
}

/**
 * USDT(ERC20) approve(spender, exact)。
 * **铁律：exact-amount，禁 infinite/MaxUint256**（handoff §2 资损硬约束）。
 * 调用方传入精确需求额（充值=amount；付账=amount+calculateFee(amount)）。
 */
export function useApprove() {
  const chainId = useChainId();
  const { writeContract, data: hash, isPending, error } = useWriteContract();
  const { isPending: rawConfirming, isSuccess } = useWaitForTransactionReceipt({
    hash,
  });
  const isConfirming = !!hash && rawConfirming;

  const approve = (spender: `0x${string}`, amount: bigint) => {
    writeContract({
      address: getUsdt(chainId),
      abi: MockUSDTABI,
      functionName: "approve",
      // exact amount —— 绝不传 MaxUint256/infinite。
      args: [spender, amount],
      ...localGasOverride(chainId),
    });
  };

  return { approve, hash, isPending, isConfirming, isSuccess, error };
}
