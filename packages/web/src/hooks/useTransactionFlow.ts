// src/hooks/useTransactionFlow.ts

export type TxStatus = "idle" | "pending-signature" | "pending-confirmation" | "success" | "error";

export interface TxState {
  status: TxStatus;
  txHash?: string;
  error?: string;
}

// 合约 revert reason → 中文映射
const REVERT_MESSAGES: Record<string, string> = {
  "Not registered": "请先注册账户",
  "Already registered": "该钱包已注册",
  "Zero deposit": "充值金额不能为零",
  "No deposit": "没有可提取的余额",
  "Service still active": "请先停用服务后再提现",
  "Has unpaid bills": "请先支付所有账单后再提现",
  "Not your bill": "该账单不属于您",
  "Already paid": "该账单已支付",
  "Insufficient payment": "支付金额不足",
  "Only oracle": "仅限预言机操作",
};

export function parseContractError(error: unknown): string {
  if (!error) return "未知错误";
  const msg = error instanceof Error ? error.message : String(error);

  // 尝试匹配 revert reason
  for (const [key, value] of Object.entries(REVERT_MESSAGES)) {
    if (msg.includes(key)) return value;
  }

  // 用户拒绝签名
  if (msg.includes("User rejected") || msg.includes("user rejected")) {
    return "交易已取消";
  }

  // 其他错误
  if (msg.includes("insufficient funds")) {
    return "余额不足以支付 Gas 费";
  }

  return "交易失败，请重试";
}

/**
 * 将 wagmi 的 writeContract 状态映射为统一的 TxState。
 *
 * 用法：
 * const { writeContract, data: hash, isPending, error } = useWriteContract();
 * const { isLoading: isConfirming, isSuccess } = useWaitForTransactionReceipt({ hash });
 * const txState = useTxState({ hash, isPending, isConfirming, isSuccess, error });
 */
export function useTxState(params: {
  hash?: `0x${string}`;
  isPending: boolean;
  isConfirming: boolean;
  isSuccess: boolean;
  error: Error | null;
}): TxState {
  const { hash, isPending, isConfirming, isSuccess, error } = params;

  if (error) {
    return { status: "error", error: parseContractError(error), txHash: hash };
  }
  if (isSuccess) {
    return { status: "success", txHash: hash };
  }
  if (isConfirming && hash) {
    return { status: "pending-confirmation", txHash: hash };
  }
  if (isPending) {
    return { status: "pending-signature" };
  }
  return { status: "idle" };
}
