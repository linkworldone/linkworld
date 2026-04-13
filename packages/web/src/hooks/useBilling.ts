import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { billingService } from "@/services";

export function useBills(address: string | undefined, filter?: "unpaid" | "paid") {
  return useQuery({
    queryKey: ["bills", address, filter],
    queryFn: () => billingService.getBills(address!, filter),
    enabled: !!address,
  });
}

export function useBillDetail(billId: string | undefined) {
  return useQuery({
    queryKey: ["bill", billId],
    queryFn: () => billingService.getBillDetail(billId!),
    enabled: !!billId,
  });
}

export function usePayBill() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (billId: string) => billingService.payBill(billId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["bills"] });
      queryClient.invalidateQueries({ queryKey: ["bill"] });
    },
  });
}

export function useMonthEstimate(address: string | undefined) {
  return useQuery({
    queryKey: ["monthEstimate", address],
    queryFn: () => billingService.getCurrentMonthEstimate(address!),
    enabled: !!address,
  });
}
