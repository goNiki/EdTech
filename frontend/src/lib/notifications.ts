import { api } from './api';

export interface NotificationItem {
  id: number;
  user_id?: number;
  title: string;
  message: string;
  type: 'homework_graded' | 'new_submission' | 'new_lesson' | 'general' | string;
  link_url?: string;
  is_read: boolean;
  created_at: string;
}

export interface NotificationsResponse {
  notifications: NotificationItem[];
  unread_count: number;
}

const STORAGE_KEY = 'edtech_notifications_store';

const DEFAULT_MOCK_NOTIFICATIONS: NotificationItem[] = [
  {
    id: 101,
    title: 'Домашнее задание проверено',
    message: 'Преподаватель оценил ваш ответ на практическое задание на 95 баллов!',
    type: 'homework_graded',
    link_url: '/dashboard/courses',
    is_read: false,
    created_at: new Date(Date.now() - 1000 * 60 * 25).toISOString(), // 25 минут назад
  },
  {
    id: 102,
    title: 'Опубликован новый урок',
    message: 'В курс «Архитектура микросервисов» добавлен раздел «gRPC и Protocol Buffers».',
    type: 'new_lesson',
    link_url: '/courses',
    is_read: false,
    created_at: new Date(Date.now() - 1000 * 60 * 180).toISOString(), // 3 часа назад
  },
  {
    id: 103,
    title: 'Добро пожаловать в ED.Learn',
    message: 'Исследуйте каталог образовательных программ и начните обучение уже сегодня.',
    type: 'general',
    link_url: '/courses',
    is_read: true,
    created_at: new Date(Date.now() - 1000 * 60 * 60 * 24).toISOString(), // 1 день назад
  },
];

function getStoredNotifications(): NotificationItem[] {
  if (typeof window === 'undefined') return DEFAULT_MOCK_NOTIFICATIONS;
  const raw = localStorage.getItem(STORAGE_KEY);
  if (!raw) {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(DEFAULT_MOCK_NOTIFICATIONS));
    return DEFAULT_MOCK_NOTIFICATIONS;
  }
  try {
    return JSON.parse(raw);
  } catch {
    return DEFAULT_MOCK_NOTIFICATIONS;
  }
}

function saveStoredNotifications(items: NotificationItem[]): void {
  if (typeof window === 'undefined') return;
  localStorage.setItem(STORAGE_KEY, JSON.stringify(items));
}

/**
 * Получить список уведомлений и счетчик непрочитанных
 */
export async function fetchNotifications(): Promise<NotificationsResponse> {
  try {
    const res = await api.get('/notifications');
    const data = res.data?.data || res.data;
    if (data && Array.isArray(data.notifications)) {
      return {
        notifications: data.notifications,
        unread_count: data.unread_count ?? data.notifications.filter((n: NotificationItem) => !n.is_read).length,
      };
    }
  } catch (err) {
    // API endpoint might be pending backend deployment
  }

  const items = getStoredNotifications();
  const unreadCount = items.filter((n) => !n.is_read).length;
  return {
    notifications: items,
    unread_count: unreadCount,
  };
}

/**
 * Пометить уведомление как прочитанное
 */
export async function markNotificationAsRead(id: number): Promise<void> {
  try {
    await api.patch(`/notifications/${id}/read`);
  } catch {
    // API fallback
  }

  const items = getStoredNotifications();
  const updated = items.map((n) => (n.id === id ? { ...n, is_read: true } : n));
  saveStoredNotifications(updated);
}

/**
 * Пометить все уведомления как прочитанные
 */
export async function markAllNotificationsAsRead(): Promise<void> {
  try {
    await api.post('/notifications/read-all');
  } catch {
    // API fallback
  }

  const items = getStoredNotifications();
  const updated = items.map((n) => ({ ...n, is_read: true }));
  saveStoredNotifications(updated);
}
