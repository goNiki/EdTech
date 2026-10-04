# 📦 Репозиторий: `internal/repository/category`
> **Путь:** `internal/repository/category`  
> **Роль:** SQL-запросы к таблице категорий (`categories`).

---

## 🎯 Назначение и ответственность
Обеспечивает CRUD-операции над категориями курсов и агрегационный запрос подсчета опубликованных курсов через `LEFT JOIN courses`.

---

## 📁 Структура файлов
| Файл | Описание |
|---|---|
| [`repository.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/category/repository.go) | Фабрика `NewCategoryRepo` |
| [`category.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/category/category.go) | Реализация методов `ListCategories`, `CreateCategory`, `GetCategoryByID`, `GetCategoryBySlug`, `UpdateCategory`, `DeleteCategory` |

---

## 🔗 Зависимости
- **Входящие:** [`service/category`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/category).
- **Исходящие:** [`db.QueryExecutor`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/infrastructure/db).
