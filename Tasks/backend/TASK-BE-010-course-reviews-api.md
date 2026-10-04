# 🛠 [BE-010] Разработка модуля отзывов и рейтингов курсов (Reviews & Ratings)

> **Приоритет:** Medium (P2)  
> **Статус:** Completed  
> **Связанные задачи:** FE-010, QA-010  
> **Целевой модуль:** `internal/interfaces/handlers/review/`, `internal/service/review/`, `internal/repository/review/`  
> **Документация модуля:** [PROJECT_MAP.md](../../PROJECT_MAP.md)

## 🎯 Цель задачи
Реализовать механизм сбора обратной связи и оценки курсов по шкале 1–5 звезд. Отзывы могут оставлять только зачисленные студенты, завершившие не менее 30% программы курса. Рейтинг должен агрегироваться и выводиться в карточке и лендинге курса.

## 🔍 Текущее состояние кода
- В схеме базы данных нет таблицы `course_reviews`.
- В структуре `courses` нет колонок `rating` и `reviews_count`.

## 📝 Технические требования к реализации
1. **База данных / Миграция Goose:**
   - Создать миграцию `migrators/20261004120000_create_course_reviews.sql`:
     ```sql
     CREATE TABLE course_reviews (
         id BIGSERIAL PRIMARY KEY,
         course_id BIGINT NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
         user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
         rating INT NOT NULL CHECK (rating >= 1 AND rating <= 5),
         comment TEXT,
         created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
         updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
         CONSTRAINT uq_course_user_review UNIQUE (course_id, user_id)
     );
     ALTER TABLE courses ADD COLUMN rating NUMERIC(3, 2) DEFAULT 0.00;
     ALTER TABLE courses ADD COLUMN reviews_count INT DEFAULT 0;
     ```

2. **API / Endpoints:**
   - `GET /api/v1/courses/{courseid}/reviews` — список отзывов с пагинацией и средним баллом.
   - `POST /api/v1/courses/{courseid}/reviews` — создание/обновление отзыва (Auth required).
     - Body: `{ "rating": 5, "comment": "Отличный курс!" }`
   - `DELETE /api/v1/courses/{courseid}/reviews` — удаление своего отзыва.

3. **Бизнес-логика (Service Layer):**
   - Проверка права оставить отзыв:
     - Студент зачислен в курс (`enrolledRepo.IsUserEnrolled`).
     - Прогресс студента >= 30% (`progressRepo.GetCourseProgressPercent >= 30`).
   - Автоматический пересчет `rating` и `reviews_count` курса в триггере БД или в транзакции сервиса.

## ✅ Критерии приёмки (Definition of Done)
- [x] Студент с прогрессом >= 30% может поставить оценку и написать отзыв.
- [x] Повторный отзыв того же студента обновляет существующую запись, а не дублирует.
- [x] Незачисленный студент или студент с нулевым прогрессом получает 403 Forbidden.
- [x] Средний рейтинг курса корректно вычисляется и отображается в `CourseDTO`.
