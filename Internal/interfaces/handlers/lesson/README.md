# 📦 Обработчики: `internal/interfaces/handlers/lesson`
> **Путь:** `internal/interfaces/handlers/lesson`  
> **Роль:** HTTP-контроллеры уроков (CRUD, получение урока, переключение статуса).

---

## 🎯 Назначение и ответственность
Принимает запросы по управлению уроками на маршрутах `/api/v1/lessons/*`, проводит декодирование и валидацию DTO, вызывает сервис уроков и возвращает форматированный ответ.

---

## 📁 Структура файлов и Эндпоинты
| Файл | HTTP Метод и Путь | Описание |
|---|---|---|
| [`handler.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/lesson/handler.go) | `POST /api/v1/lessons`<br>`PATCH /api/v1/lessons/{id}`<br>`DELETE /api/v1/lessons/{id}` | Создание, редактирование (валидация Puck JSON, лимит 5 МБ) и удаление урока |
| [`handler_test.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/lesson/handler_test.go) | Тесты | Unit-тесты сохранения уроков с 50+ тестами, проверка битого JSON и защиты 5 МБ |
| [`get.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/lesson/get.go) | `GET /api/v1/lessons/{id}` | Получение подробных данных урока по ID |
| [`getNavigation.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/lesson/getNavigation.go) | `GET /api/v1/lessons/{id}/navigation` | Получение контекста навигации (следующий/предыдущий) и силлабуса курса с прогрессом |
| [`update_status.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/lesson/update_status.go) | `PATCH /api/v1/lessons/{id}/status` | Изменение статуса (`draft`, `published`, `archived`) |

---

## 🔗 Зависимости
- **Входящие:** [`internal/app/di.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/app/di.go).
- **Исходящие:** [`service.LessonServices`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/lesson), [`internal/dto`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/dto).
