import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { depositService } from "@/services";

export function useDeposit(address: string | undefined) {
  return useQuery({
    queryKey: ["deposit", address],
    queryFn: () => depositService.getDeposit(address!),
    enabled: !!address,
  });
}

export function useDepositHistory(address: string | undefined) {
  return useQuery({
    queryKey: ["depositHistory", address],
    queryFn: () => depositService.getDepositHistory(address!),
    enabled: !!address,
  });
}

export function useDepositMutation() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ address, amount, currency }: { address: string; amount: bigint; currency: string }) =>
      depositService.deposit(address, amount, currency),
    onSuccess: (_, { address }) => {
      queryClient.invalidateQueries({ queryKey: ["deposit", address] });
      queryClient.invalidateQueries({ queryKey: ["depositHistory", address] });
    },
  });
}

export function useWithdrawMutation() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ address, amount }: { address: string; amount: bigint }) =>
      depositService.withdraw(address, amount),
    onSuccess: (_, { address }) => {
      queryClient.invalidateQueries({ queryKey: ["deposit", address] });
      queryClient.invalidateQueries({ queryKey: ["depositHistory", address] });
    },
  });
}
