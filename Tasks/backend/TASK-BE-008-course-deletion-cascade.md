# 🛠 [BE-008] Аудит мягкого удаления курсов (Soft Delete & Cascade Safety)

> **Приоритет:** High (P1)  
> **Статус:** Completed  
> **Связанные задачи:** FE-008, QA-008  
> **Целевой модуль:** `internal/interfaces/handlers/courses/`, `internal/service/course/`, `internal/repository/course/`  
> **Документация модуля:** [internal/service/course/FUNCTIONAL_SPEC.md](../../internal/service/course/FUNCTIONAL_SPEC.md)

## 🎯 Цель задачи
Гарантировать безопасное мягкое удаление (Soft Delete) курса. Эндпоинт `DELETE /api/v1/courses/{courseid}` должен проверять исключительные права создателя (`CanDeleteCourse`), проставлять `deleted_at = NOW()` и скрывать курс из каталога и списков студентов, исключая физическое разрушение связей в БД при наличии сданных домашних заданий и истории оплат.

## 🔍 Текущее состояние кода
- В `internal/service/course/delete.go` метод `DeleteCourse(ctx, courseID, userID)` проверяет `canDelete, err := s.accessService.CanDeleteCourse(ctx, course, userID)`.
- В `internal/repository/course/delete.go`:
  ```sql
  UPDATE courses SET deleted_at = NOW(), status = 'archived' WHERE id = $1 AND deleted_at IS NULL
  ```
- Требуется аудит: проверить, что все связанные выборки (`ListMyCourses`, `ListPublicCourses`, `GetCourseStructure`) имеют фильтр `deleted_at IS NULL` и что удаленный курс не блокирует имя/слаг для новых курсов.

## 📝 Технические требования к реализации
1. **API / Endpoint:**
   - `DELETE /api/v1/courses/{courseid}`
   - Headers: `Authorization: Bearer <token>`
   - Response DTO (200 OK):
     ```json
     { "message": "course successfully deleted" }
     ```
   - Ошибки: 403 Forbidden (если вызывает не создатель), 404 Not Found.

2. **Бизнес-логика (Service Layer):**
   - Проверка: только пользователь с ролью `creator` (или глобальный `admin`) может удалить курс.
   - Освобождение слага: при удалении курса обновлять его слаг на `deleted_<id>_<slug>`, чтобы освободить оригинальный слаг для создания новых курсов с таким же названием.

## ✅ Критерии приёмки (Definition of Done)
- [x] Только автор курса может его удалить.
- [x] Удаленный курс не появляется в каталоге и в результатах поиска.
- [x] Оригинальный слаг освобождается для повторного использования.
