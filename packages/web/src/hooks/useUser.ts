import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { userService } from "@/services";
import type { User } from "@/types";

export function useUser(address: string | undefined) {
  return useQuery<User | null>({
    queryKey: ["user", address],
    queryFn: () => userService.getUserProfile(address!),
    enabled: !!address,
  });
}

export function useRegister() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ address, email }: { address: string; email: string }) =>
      userService.register(address, email),
    onSuccess: (user) => {
      queryClient.setQueryData(["user", user.address], user);
    },
  });
}

export function useSendVerificationCode() {
  return useMutation({
    mutationFn: ({ address, email }: { address: string; email: string }) =>
      userService.sendVerificationCode(address, email),
  });
}

export function useVerifyEmail() {
  return useMutation({
    mutationFn: ({ address, code }: { address: string; code: string }) =>
      userService.verifyEmail(address, code),
  });
}
