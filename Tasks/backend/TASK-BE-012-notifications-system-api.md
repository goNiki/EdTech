# 🛠 [BE-012] Серверная система уведомлений (In-App Notification Center)

> **Статус:** Completed  
> **Приоритет:** Low (P2)  
> **Связанные задачи:** FE-012, QA-012  
> **Целевой модуль:** `internal/interfaces/handlers/notification/`, `internal/service/notification/`, `internal/repository/notification/`  
> **Документация модуля:** [PROJECT_MAP.md](../../PROJECT_MAP.md)

## 🎯 Цель задачи
Реализовать механизм внутрисистемных уведомлений пользователей. События:
1. Студенту: домашнее задание проверено преподавателем с оценкой и фидбеком.
2. Преподавателю: студент сдал новую работу на проверку.
3. Студенту: автор опубликовал новый урок в курсе, на который студент записан.

## 🔍 Текущее состояние кода
- В кодовой базе нет таблицы `notifications` и сервиса отправки уведомлений.
- Проверка ДЗ (`GradeAttemptAnswer`) завершается без уведомления студента.

## 📝 Технические требования к реализации
1. **База данных / Миграция Goose:**
   - Таблица `notifications`:
     ```sql
     CREATE TABLE notifications (
         id BIGSERIAL PRIMARY KEY,
         user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
         title TEXT NOT NULL,
         message TEXT NOT NULL,
         type VARCHAR(32) NOT NULL, -- homework_graded, new_submission, new_lesson
         link_url TEXT,
         is_read BOOLEAN NOT NULL DEFAULT FALSE,
         created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
     );
     CREATE INDEX idx_notifications_user_unread ON notifications(user_id) WHERE is_read = FALSE;
     ```

2. **API / Endpoints:**
   - `GET /api/v1/notifications` — список последних 20 уведомлений со счетчиком непрочитанных `unread_count`.
   - `PATCH /api/v1/notifications/{id}/read` — пометить уведомление прочитанным.
   - `POST /api/v1/notifications/read-all` — пометить все уведомления прочитанными.

3. **Интеграция событий в бизнес-слой:**
   - В `service.quiz.GradeAttemptAnswer`: после успешного выставления оценки создавать уведомление студенту:
     - Заголовок: *«Домашнее задание проверено»*
     - Сообщение: *«Преподаватель оценил ваш ответ на {points} баллов»*.
     - Ссылка: `/lessons/{lessonId}`.

## ✅ Критерии приёмки (Definition of Done)
- [x] Пользователь получает уведомление в реальном времени или при следующем опросе.
- [x] Поддерживается отметка одного или всех уведомлений как прочитанных.
- [x] Непрочитанные уведомления эффективно выбираются по частичному индексу.
- [x] Написаны модульные тесты (`internal/service/notification/service_test.go`), все тесты (`go test ./...`) и линтер (`golangci-lint run`) пройдены без ошибок.
- [x] Актуализирована техническая и аналитическая документация (README.md, FUNCTIONAL_SPEC.md, PROJECT_MAP.md, FEATURE_CATALOG.md).
