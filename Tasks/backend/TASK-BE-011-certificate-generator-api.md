# 🛠 [BE-011] Генерация и верификация сертификатов об окончании курса (Certificates Engine)

> **Статус:** Completed  
> **Приоритет:** Low (P2)  
> **Связанные задачи:** FE-011, QA-011  
> **Целевой модуль:** `internal/interfaces/handlers/certificate/`, `internal/service/certificate/`, `internal/repository/certificate/`  
> **Документация модуля:** [PROJECT_MAP.md](../../PROJECT_MAP.md)

## 🎯 Цель задачи
Реализовать выдачу цифровых сертификатов об успешном окончании курса при достижении 100% прогресса. Сертификат должен содержать уникальный публичный идентификатор (UUID/код), имя студента, название курса, средний балл, дату выпуска и ссылку для онлайн-верификации подлинности.

## 🔍 Текущее состояние кода
- В кодовой базе нет сущности сертификатов.
- В таблице `course_progress` фиксируется `completion_percentage` и `completed_at`.

## 📝 Технические требования к реализации
1. **База данных / Миграция Goose:**
   - Создать таблицу `certificates`:
     ```sql
     CREATE TABLE certificates (
         id BIGSERIAL PRIMARY KEY,
         certificate_code VARCHAR(64) UNIQUE NOT NULL,
         user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
         course_id BIGINT NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
         student_name TEXT NOT NULL,
         course_title TEXT NOT NULL,
         final_score NUMERIC(5, 2) NOT NULL,
         issued_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
         CONSTRAINT uq_user_course_certificate UNIQUE (user_id, course_id)
     );
     ```

2. **API / Endpoints:**
   - `GET /api/v1/courses/{courseid}/certificate` — получение/генерация сертификата для текущего студента.
     - Доступно строго при `course_progress.completion_percentage == 100`.
   - `GET /api/v1/certificates/verify/{code}` — публичная проверка подлинности сертификата без авторизации.

3. **Бизнес-логика (Service Layer):**
   - Проверить прогресс студента по курсу. Если < 100% — возвращать 400 Bad Request (`course not completed`).
   - Сгенерировать криптостойкий код сертификата: `EDL-YYYY-XXXXXX`.
   - Сохранить запись в `certificates` (идемпотентно: повторный вызов отдает уже выданный ранее сертификат).

## ✅ Критерии приёмки (Definition of Done)
- [x] Сертификат генерируется только при 100% завершении курса.
- [x] Публичный эндпоинт `/verify/{code}` отдает валидные данные сертификата для работодателей.
- [x] Повторные запросы возвращают тот же код сертификата.
- [x] Написаны модульные тесты (`internal/service/certificate/service_test.go`), все тесты и линтеры пройдены.
- [x] Синхронизирована документация (README.md, FUNCTIONAL_SPEC.md, PROJECT_MAP.md, FEATURE_CATALOG.md).

