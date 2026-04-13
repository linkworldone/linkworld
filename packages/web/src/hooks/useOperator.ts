import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { operatorService } from "@/services";

export function useRegions() {
  return useQuery({
    queryKey: ["regions"],
    queryFn: () => operatorService.getRegions(),
  });
}

export function useOperatorsByRegion(regionCode: string | undefined) {
  return useQuery({
    queryKey: ["operators", regionCode],
    queryFn: () => operatorService.getOperatorsByRegion(regionCode!),
    enabled: !!regionCode,
  });
}

export function useMyNumbers(address: string | undefined) {
  return useQuery({
    queryKey: ["myNumbers", address],
    queryFn: () => operatorService.getMyNumbers(address!),
    enabled: !!address,
  });
}

export function useApplyNumber() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ address, operatorId }: { address: string; operatorId: string }) =>
      operatorService.applyNumber(address, operatorId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["myNumbers"] });
    },
  });
}
