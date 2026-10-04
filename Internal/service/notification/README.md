# Модуль сервиса уведомлений (`internal/service/notification`)

## Назначение
Сервисный слой управления внутрисистемными уведомлениями пользователей (In-App Notification Center). Предоставляет бизнес-логику создания уведомлений о событиях обучения (проверка домашних заданий, сдача работ, публикация уроков), получение пользовательской ленты с подсчетом непрочитанных и пометку прочитанными.

## Архитектура
Пакет реализует интерфейс `services.NotificationServices`:

```go
type NotificationServices interface {
    CreateNotification(ctx context.Context, userID int64, title, message string, nType domain.NotificationType, linkURL *string) (*domain.Notification, error)
    GetFeed(ctx context.Context, userID int64, limit int) (*domain.NotificationFeed, error)
    MarkAsRead(ctx context.Context, userID, notificationID int64) error
    MarkAllAsRead(ctx context.Context, userID int64) error
}
```

## Файловая структура
| Файл | Описание |
| :--- | :--- |
| `service.go` | Реализация сервиса уведомлений, валидация и взаимодействие с репозиторием. |
| `service_test.go` | Модульные тесты для сценариев отправки, чтения ленты и обработки ошибок. |
| `README.md` | Техническая документация модуля. |
| `FUNCTIONAL_SPEC.md` | Аналитическая спецификация бизнес-логики и сценариев уведомлений. |

## Зависимости
* `edtech/internal/domain` — модели уведомлений и типов.
* `edtech/internal/repository` — интерфейс `NotificationRepository`.
* `edtech/pkg/errors` — доменные ошибки валидации и 404.
