# Модуль репозитория уведомлений (`internal/repository/notification`)

## Назначение
Реализация уровня доступа к данным (PostgreSQL) для таблицы `notifications`.

## Архитектура
Реализует интерфейс `repository.NotificationRepository`:

```go
type NotificationRepository interface {
    CreateNotification(ctx context.Context, n *domain.Notification) (*domain.Notification, error)
    GetUserNotifications(ctx context.Context, userID int64, limit int) ([]domain.Notification, error)
    CountUnreadNotifications(ctx context.Context, userID int64) (int64, error)
    MarkAsRead(ctx context.Context, userID, notificationID int64) error
    MarkAllAsRead(ctx context.Context, userID int64) error
}
```

## Таблица БД
* `notifications`:
  * `id` (`BIGSERIAL PRIMARY KEY`)
  * `user_id` (`BIGINT REFERENCES users(id) ON DELETE CASCADE`)
  * `title` (`TEXT NOT NULL`)
  * `message` (`TEXT NOT NULL`)
  * `type` (`VARCHAR(32) NOT NULL`)
  * `link_url` (`TEXT`)
  * `is_read` (`BOOLEAN NOT NULL DEFAULT FALSE`)
  * `created_at` (`TIMESTAMPTZ NOT NULL DEFAULT NOW()`)

## Индексы
* `idx_notifications_user_unread ON notifications(user_id) WHERE is_read = FALSE` — частичный индекс для быстрого подсчета непрочитанных уведомлений.
* `idx_notifications_user_created ON notifications(user_id, created_at DESC)` — составной индекс для эффективной выборки последних событий пользователя.
