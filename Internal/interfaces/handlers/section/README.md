# 📦 Обработчики: `internal/interfaces/handlers/section`
> **Путь:** `internal/interfaces/handlers/section`  
> **Роль:** HTTP-контроллеры создания, обновления, удаления секций курса и упорядочивания уроков.

---

## 🎯 Назначение и ответственность
Обрабатывает HTTP-запросы по секциям (`/api/v1/sections/*`), принимает запросы на переупорядочивание уроков внутри модуля и изменение статусов публикации.

---

## 📁 Структура файлов и Эндпоинты
| Файл | HTTP Метод и Путь | Описание |
|---|---|---|
| [`handler.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/section/handler.go) | `POST /api/v1/sections`<br>`PATCH /api/v1/sections/{id}`<br>`DELETE /api/v1/sections/{id}` | Базовые CRUD-методы секции |
| [`reorder_lessons.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/section/reorder_lessons.go) | `PUT /api/v1/sections/{id}/reorder-lessons` | Изменение порядка уроков внутри модуля |
| [`update_status.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/section/update_status.go) | `PATCH /api/v1/sections/{id}/status` | Обновление статуса видимости секции |

---

## 🔗 Зависимости
- **Входящие:** [`internal/app/di.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/app/di.go).
- **Исходящие:** [`service.SectionServices`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/section), [`internal/dto`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/dto).
