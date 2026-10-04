# 🛠 [BE-007] Реализация API категорий курсов (Categories CRUD & Listing)

> **Статус:** Completed  
> **Приоритет:** Medium (P1)  
> **Связанные задачи:** FE-007, QA-007  
> **Целевой модуль:** `internal/interfaces/handlers/category/`, `internal/service/category/`, `internal/repository/category/`  
> **Документация модуля:** [PROJECT_MAP.md](../../PROJECT_MAP.md)

## 🎯 Цель задачи
В базе данных существует таблица категорий (`migrators/20250912230010_create_categories.sql`), а в таблице `courses` есть внешний ключ `category_id`. Однако в бэкенде отсутствуют репозиторий, сервис и HTTP-эндпоинты для категорий. Из-за этого клиенты не могут получить список категорий для фильтрации витрины курсов и для привязки категории при создании курса.

## 🔍 Текущее состояние кода
- Таблица `categories` (поля: `id`, `name`, `slug`, `description`, `icon_url`, `created_at`).
- В `internal/app/di.go` нет роутов и зависимостей для категорий.
- В `internal/repository/course/listPublicCourses.go` уже реализована фильтрация `WHERE category_id = $x`, но передавать туда реальные ID из UI невозможно.

## 📝 Технические требования к реализации
1. **API / Endpoints:**
   - `GET /api/v1/categories` (публичный эндпоинт, получение активных категорий).
   - Response DTO (200 OK):
     ```json
     {
       "categories": [
         {
           "id": 1,
           "name": "Программирование",
           "slug": "programming",
           "description": "Курсы по разработке ПО",
           "icon_url": "https://...",
           "courses_count": 12
         }
       ]
     }
     ```
   - Опционально для админа: `POST /api/v1/categories` (создание новой категории).

2. **Архитектурные слои:**
   - **Repository:** `internal/repository/category/` с методом `ListCategories(ctx, q) ([]domain.Category, error)` с подсчетом опубликованных курсов через `LEFT JOIN courses`.
   - **Service:** `internal/service/category/` с методом `GetAllCategories(ctx) ([]domain.Category, error)`.
   - **Handler:** `internal/interfaces/handlers/category/` с методом `ListCategories`.
   - Регистрация в Chi Router (`di.go`): `r.Get("/api/v1/categories", d.CategoryHdl().ListCategories)`.

## ✅ Критерии приёмки (Definition of Done)
- [x] Запрос `GET /api/v1/categories` возвращает список доступных категорий.
- [x] В выдаче присутствует количество привязанных курсов.
- [x] Модуль зарегистрирован в DI контейнере и покрыт документацией в `PROJECT_MAP.md`.
