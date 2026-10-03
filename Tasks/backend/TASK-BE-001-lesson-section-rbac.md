# 🛠 [BE-001] Добавление RBAC и проверки прав владения курсом в CRUD секций и уроков

**Статус:** ✅ Completed (коммит: `867ed88`)

> **Приоритет:** Critical (P0)  
> **Связанные задачи:** FE-001, QA-001  
> **Целевой модуль:** `internal/service/lesson/`, `internal/service/section/`, `internal/interfaces/handlers/lesson/`, `internal/interfaces/handlers/section/`  
> **Документация модуля:** [internal/service/lesson/FUNCTIONAL_SPEC.md](../../internal/service/lesson/FUNCTIONAL_SPEC.md), [internal/service/section/FUNCTIONAL_SPEC.md](../../internal/service/section/FUNCTIONAL_SPEC.md)

## 🎯 Цель задачи
Устранить критическую уязвимость BOLA / IDOR (Broken Object Level Authorization). В текущей реализации любой авторизованный пользователь (включая студента) может создавать, редактировать и удалять модули и уроки в любых чужих курсах, поскольку методы сервисов и хендлеров не принимают `userID` и не проверяют права через `AccessService.CanEditCourse`.

## 🔍 Текущее состояние кода
- В `internal/service/lesson/updateLesson.go` и `deleteLesson.go` методы `UpdateLesson(ctx, lesson)` и `DeleteLesson(ctx, id)` не имеют параметра `userID` и не вызывают `accessService.CanEditCourse`.
- В `internal/service/lesson/createLesson.go` метод `CreateLesson(ctx, lesson)` не проверяет права создателя.
- В `internal/service/section/` методы `CreateSection`, `UpdateSection`, `DeleteSection`, `ReorderLessons` также не проверяют права пользователя.
- Хендлеры `LessonHandler` и `SectionHandler` не передают `userID` из JWT в сервисы.

## 📝 Технические требования к реализации
1. **API / Handlers:**
   - Извлечь `userID := h.authMiddleware.GetUserID(r.Context())` во всех методах:
     - `POST /api/v1/sections`
     - `PATCH /api/v1/sections/{id}`
     - `PATCH /api/v1/sections/{id}/status`
     - `PUT /api/v1/sections/{id}/reorder-lessons`
     - `DELETE /api/v1/sections/{id}`
     - `POST /api/v1/lessons`
     - `PATCH /api/v1/lessons/{id}`
     - `PATCH /api/v1/lessons/{id}/status`
     - `DELETE /api/v1/lessons/{id}`
   - Если `userID == 0` — отдавать `errorsAPP.ErrUnauthorized` (401).

2. **Бизнес-логика (Service Layer):**
   - Расширить интерфейсы `LessonServices` и `SectionServices`: добавить аргумент `userID int64` в каждый мутирующий метод.
   - Внедрить зависимость `accessService service.AccessService` в `lessonService` и `sectionService`.
   - В каждом методе:
     - Получить целевой `courseID` (напрямую или загрузив существующий урок/секцию из репозитория).
     - Загрузить курс через `courseRepo.GetCourseByID`.
     - Вызвать `canEdit, err := s.accessService.CanEditCourse(ctx, course, userID)`.
     - Если `!canEdit` — возвращать `errorsAPP.ErrForbidden` (403).

3. **База данных / Хранилище:**
   - Изменений схемы БД не требуется (используются существующие таблицы `courses`, `users_courses`, `role_permissions`).

4. **Побочные эффекты:**
   - Логирование попыток несанкционированного доступа (Warn уровень) с указанием `userID`, `targetCourseID`, `action`.

## ✅ Критерии приёмки (Definition of Done)
- [x] Попытка модификации или удаления урока/секции пользователем, не являющимся автором или преподавателем курса, возвращает HTTP 403 Forbidden.
- [x] Преподаватель и создатель курса успешно сохраняют свои изменения.
- [x] Ошибки авторизации логируются с контекстом запроса.
- [x] Модульная документация в `internal/service/lesson/FUNCTIONAL_SPEC.md` и `section/FUNCTIONAL_SPEC.md` обновлена.
