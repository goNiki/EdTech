# 🛠 [BE-029] История попыток тестирования, стратегия сохранения высшего балла (Best Score) и Pre-flight API

> **Статус:** Completed (хэш коммита: ad240e3)  
> **Приоритет:** High (P1)  
> **Связанные задачи:** FE-030, QA-030  
> **Целевой модуль:** `internal/interfaces/handlers/progress/`, `internal/service/progress/`, `internal/repository/progress/`  
> **Документация модуля:** [internal/service/progress/README.md](../../internal/service/README.md)

---

## 🎯 Цель задачи
Реализовать серверную поддержку стратегии **Best Score Preservation (Сохранение лучшего результата)** и историю попыток студента по уроку-тесту.

Сейчас, когда студент повторно заходит в тест, его предыдущий результат теряется, а новая попытка стирает статус `completed`.
Необходимо:
1. Хранить полную историю каждой попытки прохождения теста (`quiz_attempts`), включая набранный балл, процент правильных ответов, время начала и завершения.
2. В общем прогрессе урока (`lesson_progress.score`) **всегда сохранять наивысший результат (High Score)**: если попытка 1 была на 90%, а попытка 2 — на 40% или была заброшена, итоговый балл студента в курсе остается **90%**!
3. Предоставить эндпоинт `GET /api/v1/lessons/{id}/attempts/summary` для Pre-flight экрана, сообщающий количество совершенных попыток, оставшийся лимит и лучший балл.

---

## 🔍 Текущее состояние кода
- Таблица `quiz_attempts` уже существует в БД (`id`, `user_id`, `quiz_id`, `score`, `max_score`, `created_at`, `completed_at`).
- Однако таблица `lesson_progress` при вызове `CompleteLesson` просто перезаписывает поле `score` последним пришедшим значением:
  `UPDATE lesson_progress SET score = $1, status = $2...`
- Отсутствует сводный API для получения истории попыток конкретного урока для авторизованного студента.

---

## 📝 Технические требования к реализации

### 1. API / Endpoints
- `GET /api/v1/lessons/{id}/attempts/summary`
  - **Авторизация:** JWT токен студента.
  - **Response DTO (`dto.LessonAttemptsSummaryResponse`):**
    ```json
    {
      "lesson_id": 10,
      "total_attempts_made": 2,
      "max_attempts_allowed": 3,
      "can_start_new_attempt": true,
      "best_score": 90,
      "best_score_percentage": 90,
      "is_passed": true,
      "passing_threshold": 70,
      "last_attempt": {
        "attempt_id": 15,
        "score": 85,
        "submitted_at": "2026-10-04T10:30:00Z"
      },
      "attempts_history": [
        {
          "attempt_id": 12,
          "score": 90,
          "submitted_at": "2026-10-03T15:20:00Z"
        },
        {
          "attempt_id": 15,
          "score": 85,
          "submitted_at": "2026-10-04T10:30:00Z"
        }
      ]
    }
    ```

- `POST /api/v1/lessons/{id}/attempts/start`
  - Создает новую запись попытки в `quiz_attempts` со статусом `in_progress` и таймстемпом `started_at = NOW()`.
  - Проверяет лимит попыток: если `total_attempts_made >= max_attempts_allowed`, возвращает `403 Forbidden` с сообщением «Лимит попыток исчерпан».
  - Возвращает `{ "attempt_id": 16, "started_at": "..." }`.

### 2. Бизнес-логика (Best Score Strategy в `CompleteLesson`)
При завершении попытки `POST /api/v1/lessons/{id}/complete`:
1. Зафиксировать результат текущей попытки в `quiz_attempts` (`score`, `completed_at = NOW()`).
2. Запросить из БД максимальный балл среди ВСЕХ завершенных попыток данного студента по этому уроку:
   `SELECT COALESCE(MAX(score), 0) FROM quiz_attempts WHERE user_id = $1 AND lesson_id = $2 AND completed_at IS NOT NULL`.
3. Записать в `lesson_progress.score = GREATEST(currentScore, bestPreviousScore)`.
4. Статус `lesson_progress.status`:
   - Если `bestScore >= passingThreshold` ➔ `'completed'` (урок навсегда остается зачтенным).
   - Если `bestScore < passingThreshold` ➔ `'failed'`.

---

## ✅ Критерии приёмки (Definition of Done)
- [x] При прохождении теста со второй попытки на более низкий балл итоговая оценка в прогрессе курса не уменьшается.
- [x] Заброшенная или прерванная попытка не сбрасывает статус ранее успешно сданного урока.
- [x] Эндпоинт `GET /lessons/{id}/attempts/summary` отдает корректное количество попыток, историю и лучший результат.
- [x] При достижении лимита попыток сервер блокирует создание новой попытки со статусом 403.
