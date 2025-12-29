'use client';

import { useState } from 'react';
import { useNotificationStore } from '@/lib/stores/notification-store';
import { Bell, CheckCircle2, X, AlertTriangle, ExternalLink } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { toast } from 'sonner';

export function NotificationCenter() {
  const { notifications, unreadCount, markAsRead, markAllAsRead, clearAll, removeNotification } = useNotificationStore();
  const [isOpen, setIsOpen] = useState(false);

  const handleClearAll = () => {
    clearAll();
    toast.success('All notifications cleared');
  };

  const handleMarkAllRead = () => {
    markAllAsRead();
    toast.success('All notifications marked as read');
  };

  return (
    <div className="relative">
      <Button
        variant="ghost"
        size="icon"
        className="relative"
        onClick={() => setIsOpen(!isOpen)}
      >
        <Bell className="h-5 w-5" />
        {unreadCount > 0 && (
          <span className="absolute -top-1 -right-1 flex h-4 w-4 items-center justify-center rounded-full bg-red-500 text-xs text-white">
            {unreadCount}
          </span>
        )}
      </Button>

      {isOpen && (
        <div className="absolute right-0 mt-2 w-80 bg-white dark:bg-neutral-900 border border-neutral-200 dark:border-neutral-800 rounded-lg shadow-lg z-50">
          <div className="p-4 border-b flex justify-between items-center">
            <h3 className="font-semibold">Price Drop Notifications</h3>
            {notifications.length > 0 && (
              <div className="flex gap-2">
                <Button
                  variant="ghost"
                  size="icon"
                  onClick={handleMarkAllRead}
                  className="text-xs h-6 px-2"
                >
                  Mark All Read
                </Button>
                <Button
                  variant="ghost"
                  size="icon"
                  onClick={handleClearAll}
                  className="text-xs h-6 px-2"
                >
                  Clear All
                </Button>
              </div>
            )}
          </div>

          {notifications.length === 0 ? (
            <div className="p-4 text-center text-neutral-500">
              <AlertTriangle className="h-6 w-6 mx-auto mb-2" />
              <p>No price drops detected</p>
              <p className="text-xs mt-1">
                We'll notify you when prices drop to your target!
              </p>
            </div>
          ) : (
            <>
              <div className="max-h-96 overflow-y-auto">
                {notifications.map((notification) => (
                  <div
                    key={notification.id}
                    className={`p-3 border-b last:border-b-0 hover:bg-neutral-50 dark:hover:bg-neutral-800 ${
                      !notification.isRead ? 'bg-neutral-50 dark:bg-neutral-800' : ''
                    }`}
                  >
                    <div className="flex items-start gap-3">
                      <CheckCircle2 className="h-5 w-5 text-green-500 mt-0.5 flex-shrink-0" />
                      <div className="flex-1 min-w-0">
                        <div className="flex justify-between items-start">
                          <p className="font-medium truncate">{notification.itemName}</p>
                          {!notification.isRead && (
                            <span className="text-xs bg-blue-500 text-white px-1 rounded">new</span>
                          )}
                        </div>
                        <p className="text-sm text-neutral-600 dark:text-neutral-400">
                          Price: <span className="font-bold text-green-600">${notification.currentPrice.toFixed(2)}</span>
                          {notification.targetPrice && (
                            <> (Target: ${notification.targetPrice.toFixed(2)})</>
                          )}
                        </p>
                        <p className="text-xs text-neutral-500 mt-1">
                          {new Date(notification.createdAt).toLocaleString()}
                        </p>
                        <div className="mt-2 flex gap-2">
                          <Button
                            variant="ghost"
                            size="icon"
                            className="h-6 px-2 text-xs"
                            onClick={() => {
                              window.open(notification.url, '_blank');
                              markAsRead(notification.id);
                            }}
                          >
                            <ExternalLink className="h-3 w-3 mr-1" />
                            View Product
                          </Button>
                          <Button
                            variant="ghost"
                            size="icon"
                            className="h-6 px-2 text-xs"
                            onClick={() => removeNotification(notification.id)}
                          >
                            <X className="h-3 w-3 mr-1" />
                            Dismiss
                          </Button>
                        </div>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
              <div className="p-2 border-t text-xs text-neutral-500 flex justify-between">
                <span>{notifications.length} notification{notifications.length !== 1 ? 's' : ''}</span>
                <span>{unreadCount} unread</span>
              </div>
            </>
          )}
        </div>
      )}
    </div>
  );
}