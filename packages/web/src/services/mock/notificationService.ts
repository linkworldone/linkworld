import type { Notification } from "@/types";
import { delay } from "./delay";
import { mockNotifications } from "./data";

export const notificationService = {
  async getNotifications(_address: string): Promise<Notification[]> {
    await delay();
    return [...mockNotifications];
  },
  async getUnreadCount(_address: string): Promise<number> {
    await delay(200);
    return mockNotifications.filter((n) => !n.read).length;
  },
  async markAsRead(notificationId: string): Promise<void> {
    await delay(300);
    const notif = mockNotifications.find((n) => n.id === notificationId);
    if (notif) notif.read = true;
  },
  async markAllAsRead(_address: string): Promise<void> {
    await delay(300);
    mockNotifications.forEach((n) => (n.read = true));
  },
};
