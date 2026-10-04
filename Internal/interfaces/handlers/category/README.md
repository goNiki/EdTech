# 🌐 HTTP Обработчик: `internal/interfaces/handlers/category`
> **Путь:** `internal/interfaces/handlers/category`  
> **Роль:** HTTP-эндпоинты для категорий курсов (`/api/v1/categories`).

---

## 🎯 Назначение и ответственность
Принимает HTTP-запросы на получение списка категорий и создание новых категорий, валидирует входные данные и отдает JSON-ответы клиенту.

---

## 📁 Структура файлов
| Файл | Описание |
|---|---|
| [`category.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/category/category.go) | `ListCategories` (GET /api/v1/categories), `CreateCategory` (POST /api/v1/categories) |

---

## 🔗 Зависимости
- **Входящие:** `chi.Router` в [`internal/app/di.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/app/di.go).
- **Исходящие:** [`service.CategoryServices`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service).
