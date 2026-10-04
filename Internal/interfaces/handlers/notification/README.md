# Модуль HTTP-обработчиков уведомлений (`internal/interfaces/handlers/notification`)

## Назначение
Предоставление REST API для клиентского центра уведомлений пользователя.

## Эндпоинты
| Метод | Путь | Описание | Авторизация |
| :--- | :--- | :--- | :--- |
| `GET` | `/api/v1/notifications` | Получить ленту последних уведомлений со счетчиком `unread_count`. Параметр `limit` (опционально, default 20, max 100). | Bearer JWT |
| `PATCH` | `/api/v1/notifications/{id}/read` | Пометить конкретное уведомление как прочитанное. | Bearer JWT |
| `POST` | `/api/v1/notifications/read-all` | Пометить все уведомления пользователя как прочитанные. | Bearer JWT |

## Формат ответов
* `GET /api/v1/notifications`:
```json
{
  "notifications": [
    {
      "id": 1,
      "user_id": 10,
      "title": "Домашнее задание проверено",
      "message": "Преподаватель оценил ваш ответ на 100 баллов",
      "type": "homework_graded",
      "link_url": "/lessons/42",
      "is_read": false,
      "created_at": "2026-10-04T10:00:00Z"
    }
  ],
  "unread_count": 1
}
```
* `PATCH /api/v1/notifications/{id}/read`:
```json
{
  "message": "notification marked as read"
}
```
* `POST /api/v1/notifications/read-all`:
```json
{
  "message": "all notifications marked as read"
}
```
