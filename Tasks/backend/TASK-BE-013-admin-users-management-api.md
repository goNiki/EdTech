# 🛠 [BE-013] Расширение эндпоинтов администратора: получение списка пользователей (Admin Users API)

> **Статус:** Completed  
> **Приоритет:** Low (P2)  
> **Связанные задачи:** FE-013, QA-013  
> **Целевой модуль:** `internal/interfaces/handlers/auth/`, `internal/service/auth/`, `internal/repository/auth/`  
> **Документация модуля:** [internal/service/auth/FUNCTIONAL_SPEC.md](../../internal/service/auth/FUNCTIONAL_SPEC.md)

## 🎯 Цель задачи
В кодовой базе реализованы методы смены роли пользователя (`PATCH /api/v1/admin/users/{id}/role`) и блокировки (`PATCH /api/v1/admin/users/{id}/ban`), но отсутствует базовый эндпоинт получения списка пользователей для панели администратора. Требуется создать `GET /api/v1/admin/users` с поиском, пагинацией и фильтрами по роли и статусу блокировки.

## 🔍 Текущее состояние кода
- В `internal/app/di.go` ([L528-L533](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/app/di.go#L528-L533)) зарегистрированы только мутирующие админские методы:
  - `PATCH /api/v1/admin/users/{id}/role`
  - `PATCH /api/v1/admin/users/{id}/ban`
- Нет эндпоинта выборки пользователей.

## 📝 Технические требования к реализации
1. **API / Endpoint:**
   - `GET /api/v1/admin/users`
   - Headers: `Authorization: Bearer <admin_token>`
   - Query Params:
     - `search` (string, поиск по email/username/имени)
     - `role` (string: `student`, `teacher`, `author`, `admin`)
     - `is_banned` (bool, optional)
     - `page`, `page_size`
   - Response DTO (200 OK):
     ```json
     {
       "users": [
         {
           "id": 1,
           "email": "teacher@edtech.io",
           "username": "teacher",
           "first_name": "Анна",
           "last_name": "Смирнова",
           "role": "teacher",
           "is_banned": false,
           "created_at": "2026-09-01T12:00:00Z"
         }
       ],
       "total": 45,
       "page": 1,
       "page_size": 20
     }
     ```
   - Ошибки: 401 (не авторизован), 403 (пользователь не является `admin`).

2. **Бизнес-логика (Service Layer):**
   - Проверка прав: `admin.IsAdmin()`, иначе `errorsAPP.ErrForbidden`.
   - Репозиторий: постраничный запрос `SELECT ... FROM users WHERE ... ORDER BY created_at DESC`.

## ✅ Критерии приёмки (Definition of Done)
- [x] Администратор может получить список пользователей с пагинацией и фильтрами.
- [x] Пользователи без роли `admin` получают 403 Forbidden.
- [x] Ролевые изменения и бан можно отслеживать в реальном времени.
