# 🛠 [BE-003] Обработка коллизий слагов и валидация URL курсов

**Статус:** ✅ Completed (хэш коммита: `b2bb28e`)

> **Приоритет:** High (P0)  
> **Связанные задачи:** FE-003, QA-003  
> **Целевой модуль:** `internal/service/course/`, `internal/repository/course/`  
> **Документация модуля:** [internal/service/course/FUNCTIONAL_SPEC.md](../../internal/service/course/FUNCTIONAL_SPEC.md)

## 🎯 Цель задачи
Обеспечить отказоустойчивость при сохранении курсов с одинаковыми слагами. Сейчас при совпадении слага база возвращает ошибку уникальности `23505`, из-за чего создание курса прерывается с ошибкой 500/400. Задача — добавить автоматическое разрешение коллизий (добавление суффикса `-1`, `-2` или случайного хэша) либо понятную ошибку `ErrCourseSlugAlreadyExists`.

## 🔍 Текущее состояние кода
- В `internal/service/course/create.go` слаг принимается напрямую из запроса.
- При наличии дубликата в `courses (slug)` операция падает на уровне СУБД.
- В DTO `CreateCourseRequest` валидация слага: `validate:"required,min=3,max=100"`.

## 📝 Технические требования к реализации
1. **API / Endpoint:**
   - `POST /api/v1/courses`
   - При совпадении слага:
     - Либо автоматически сгенерировать уникальный слаг (например, `golang-basics-2`), проверив `courseRepo.ExistsBySlug`.
     - Либо вернуть 409 Conflict с деталями `{ "error": "slug_already_exists", "suggested_slug": "golang-basics-1" }`.

2. **Бизнес-логика (Service Layer):**
   - Внедрить нормализацию слага: приведение к нижнему регистру, замена пробелов и подчеркиваний на дефис, удаление спецсимволов.
   - Метод `resolveUniqueSlug(ctx, slug)`:
     ```go
     candidate := slug
     counter := 1
     for {
         exists, err := s.courseRepo.ExistsBySlug(ctx, s.db, candidate)
         if err != nil { return "", err }
         if !exists { return candidate, nil }
         candidate = fmt.Sprintf("%s-%d", slug, counter)
         counter++
     }
     ```

3. **База данных / Хранилище:**
   - Добавить в `internal/repository/course/` метод `ExistsBySlug(ctx context.Context, q db.QueryExecutor, slug string) (bool, error)`.

## ✅ Критерии приёмки (Definition of Done)
- [x] Попытка создания курса с уже существующим слагом не приводит к необработанному сбою БД (500).
- [x] Сервис создает уникальный слаг с суффиксом или возвращает корректный 409 Conflict.
- [x] Слаг гарантированно содержит только `a-z`, `0-9`, `-`.
