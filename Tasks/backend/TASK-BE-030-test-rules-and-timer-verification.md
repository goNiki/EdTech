# 🛠 [BE-030] Настройки правил тестирования: таймеры (общий и блиц), лимиты попыток, скрытие ответов и серверная валидация времени

> **Приоритет:** High (P1)  
> **Связанные задачи:** FE-031, QA-031  
> **Целевой модуль:** `internal/domain/`, `internal/interfaces/handlers/lesson/`, `internal/service/progress/`  
> **Документация модуля:** [internal/domain/README.md](../../internal/domain/README.md)

---

## 🎯 Цель задачи
Реализовать серверную поддержку гибких правил проведения тестирования в уроке:
1. **Таймеры времени:**
   - Общий таймер на тест (`time_limit_minutes`, например 15 или 30 минут).
   - Повопросный блиц-таймер (`question_time_limit_seconds`, например 30 или 45 секунд на вопрос).
2. **Лимит попыток (`max_attempts`):**
   - 0 / null — неограниченно (тренажер).
   - 1 — экзамен / срез знаний.
   - N — контрольная работа с фиксированным числом попыток.
3. **Режим обратной связи / Скрытие ответов (`feedback_mode`):**
   - `'immediate'` — обучающий (сразу показывает правильный ответ и пояснение).
   - `'exam_blind'` — экзаменационный режим (правильные ответы и конкретные ошибки НЕ показываются студенту, чтобы исключить слив ответов сокурсникам).
4. **Проходной балл (`passing_score_percent`, default 70%):**
   - Минимальный порог в процентах, необходимый для зачета урока.
5. **Серверная защита от манипуляции временем:**
   - Расчет истечения времени производится строго на сервере от `attempt.started_at`, клиентские часы не могут заморозить таймер.

---

## 🔍 Текущее состояние кода
- В таблице `lessons` и домене `domain.Lesson` параметры тестирования не выделены в структурированные поля (хранятся только `type`, `title`, `content`).
- Таблица `quiz_attempts` имеет `created_at`, но сервер не проверяет дедлайн сдачи относительно `time_limit_minutes`.

---

## 📝 Технические требования к реализации

### 1. Схема данных и миграция БД
Добавить в таблицу `lessons` поле настроек тестирования:
```sql
ALTER TABLE lessons ADD COLUMN IF NOT EXISTS quiz_settings JSONB DEFAULT '{
  "time_limit_minutes": 0,
  "question_time_limit_seconds": 0,
  "max_attempts": 0,
  "passing_score_percent": 70,
  "feedback_mode": "immediate",
  "shuffle_questions": false
}'::jsonb;
```

В структуру `domain.Lesson` и `dto.Lesson`:
```go
type QuizSettings struct {
    TimeLimitMinutes         int    `json:"time_limit_minutes"`          // 0 = без таймера
    QuestionTimeLimitSeconds int    `json:"question_time_limit_seconds"` // 0 = без блица
    MaxAttempts              int    `json:"max_attempts"`               // 0 = неограниченно
    PassingScorePercent      int    `json:"passing_score_percent"`       // напр. 70
    FeedbackMode             string `json:"feedback_mode"`              // "immediate" | "exam_blind"
    ShuffleQuestions         bool   `json:"shuffle_questions"`
}
```

### 2. Серверная валидация времени при завершении попытки
В `POST /api/v1/lessons/{id}/complete`:
1. Если `quiz_settings.time_limit_minutes > 0`:
   - Вычислить реальное прошедшее время:  
     `elapsed := time.Since(attempt.StartedAt)`.
   - Допустить допустимую сетевую погрешность (Grace Period) в 15 секунд.
   - Если `elapsed > (timeLimit + gracePeriod)`:
     - Зафиксировать статус попытки как `timed_out`.
     - Засчитать только те ответы, которые были отправлены до истечения лимита.
2. В режиме `feedback_mode == "exam_blind"`:
   - В теле ответа `CompleteLessonResponse` НЕ возвращать поля `correct_answer`, `options.is_correct` и `explain`.
   - Возвращать только суммарный набранный балл (`score`, `total_points`, `is_passed`).

---

## ✅ Критерии приёмки (Definition of Done)
- [x] Настройки `quiz_settings` сохраняются и считываются в CRUD урока.
- [x] Если время теста вышло, сервер отклоняет запоздалые ответы и завершает попытку по таймауту.
- [x] В режиме `exam_blind` из JSON-ответов сервера удаляются все правильные ответы и пояснения.
- [x] Урок помечается как `completed` только при преодолении `passing_score_percent`.
