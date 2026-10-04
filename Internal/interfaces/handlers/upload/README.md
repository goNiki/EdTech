# 📦 Модуль: `internal/interfaces/handlers/upload`
> **Путь:** `internal/interfaces/handlers/upload`  
> **Роль:** HTTP-хендлер загрузки файлов (multipart/form-data).

---

## 🎯 Назначение и ответственность
Принимает входящие HTTP POST запросы на загрузку файлов в формате `multipart/form-data`, проверяет наличие авторизации пользователя, извлекает файл и категорию, вызывает `UploadServices` и возвращает ответ со статусом `201 Created`.

---

## 📁 Структура файлов модуля
| Файл | Описание роли файла |
|---|---|
| [`handler.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/upload/handler.go) | Реализация `UploadHandler` с эндпоинтами `UploadFile` и `UploadBatch` |

---

## ⚙️ Эндпоинты API
| Метод | URL | Описание | Auth | DTO Ответа |
|---|---|---|---|---|
| `POST` | `/api/v1/upload` | Загрузка файла (multipart/form-data: `file`, `category`) | Bearer JWT | [`dto.FileUploadResponse`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/dto/upload.go) |
| `POST` | `/api/v1/upload/batch` | Пакетная загрузка изображений Word (multipart/form-data: `files[]`, `category`) | Bearer JWT | [`dto.BatchImageUploadResponse`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/dto/upload.go) |
