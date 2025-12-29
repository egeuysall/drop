import { create } from 'zustand';

interface PriceDropNotification {
  id: string;
  itemId: string;
  itemName: string;
  currentPrice: number;
  targetPrice?: number;
  url: string;
  createdAt: string;
  isRead: boolean;
}

interface NotificationState {
  notifications: PriceDropNotification[];
  unreadCount: number;
  addNotification: (notification: Omit<PriceDropNotification, 'id' | 'isRead'>) => void;
  markAsRead: (id: string) => void;
  markAllAsRead: () => void;
  clearAll: () => void;
  removeNotification: (id: string) => void;
}

export const useNotificationStore = create<NotificationState>((set) => ({
  notifications: [],
  unreadCount: 0,
  
  addNotification: (notification) => {
    const newNotification = {
      ...notification,
      id: Date.now().toString(),
      isRead: false,
      createdAt: new Date().toISOString()
    };
    
    set((state) => {
      const updatedNotifications = [newNotification, ...state.notifications];
      return {
        notifications: updatedNotifications,
        unreadCount: updatedNotifications.filter(n => !n.isRead).length
      };
    });
  },
  
  markAsRead: (id) => {
    set((state) => {
      const updatedNotifications = state.notifications.map(n =>
        n.id === id ? { ...n, isRead: true } : n
      );
      return {
        notifications: updatedNotifications,
        unreadCount: updatedNotifications.filter(n => !n.isRead).length
      };
    });
  },
  
  markAllAsRead: () => {
    set((state) => {
      const updatedNotifications = state.notifications.map(n => ({ ...n, isRead: true }));
      return {
        notifications: updatedNotifications,
        unreadCount: 0
      };
    });
  },
  
  clearAll: () => {
    set({ notifications: [], unreadCount: 0 });
  },
  
  removeNotification: (id) => {
    set((state) => {
      const updatedNotifications = state.notifications.filter(n => n.id !== id);
      return {
        notifications: updatedNotifications,
        unreadCount: updatedNotifications.filter(n => !n.isRead).length
      };
    });
  }
}));