import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { notificationService } from "@/services";

export function useNotifications(address: string | undefined) {
  return useQuery({
    queryKey: ["notifications", address],
    queryFn: () => notificationService.getNotifications(address!),
    enabled: !!address,
  });
}

export function useUnreadCount(address: string | undefined) {
  return useQuery({
    queryKey: ["unreadCount", address],
    queryFn: () => notificationService.getUnreadCount(address!),
    enabled: !!address,
    refetchInterval: 30_000,
  });
}

export function useMarkAsRead() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (notificationId: string) =>
      notificationService.markAsRead(notificationId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["notifications"] });
      queryClient.invalidateQueries({ queryKey: ["unreadCount"] });
    },
  });
}

export function useMarkAllAsRead() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (address: string) =>
      notificationService.markAllAsRead(address),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["notifications"] });
      queryClient.invalidateQueries({ queryKey: ["unreadCount"] });
    },
  });
}
