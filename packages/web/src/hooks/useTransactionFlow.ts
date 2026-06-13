// src/hooks/useTransactionFlow.ts
import { useTranslation } from "react-i18next";

export type TxStatus = "idle" | "pending-signature" | "pending-confirmation" | "success" | "error";

export interface TxState {
  status: TxStatus;
  txHash?: string;
  error?: string;
}

/** 极简翻译函数签名（兼容 react-i18next 的 t，测试可传 (k)=>k 或 i18n.t）。 */
type TFunc = (key: string) => string;

// 合约 revert reason → i18n key 映射（文案在 locales errors.tx.*）。
const REVERT_MESSAGES: Record<string, string> = {
  "Not registered": "errors.tx.notRegistered",
  "Already registered": "errors.tx.alreadyRegistered",
  "Zero deposit": "errors.tx.zeroDeposit",
  "No deposit": "errors.tx.noDeposit",
  "Service still active": "errors.tx.serviceStillActive",
  "Has unpaid bills": "errors.tx.hasUnpaidBills",
  "Not your bill": "errors.tx.notYourBill",
  "Already paid": "errors.tx.alreadyPaid",
  "Insufficient payment": "errors.tx.insufficientPayment",
  "Only oracle": "errors.tx.onlyOracle",
};

/**
 * 解析合约/交易错误为已翻译文案。
 * @param error 任意错误对象。
 * @param t i18n 翻译函数（调用方在 hook/组件里 useTranslation 拿到）。
 */
export function parseContractError(error: unknown, t: TFunc): string {
  if (!error) return t("errors.tx.unknown");
  const msg = error instanceof Error ? error.message : String(error);

  // 尝试匹配 revert reason
  for (const [key, value] of Object.entries(REVERT_MESSAGES)) {
    if (msg.includes(key)) return t(value);
  }

  // 用户拒绝签名
  if (msg.includes("User rejected") || msg.includes("user rejected")) {
    return t("errors.tx.cancelled");
  }

  // 其他错误
  if (msg.includes("insufficient funds")) {
    return t("errors.tx.insufficientGas");
  }

  return t("errors.tx.failed");
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
  const { t } = useTranslation();

  if (error) {
    return { status: "error", error: parseContractError(error, t), txHash: hash };
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
