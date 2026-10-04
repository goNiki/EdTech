# 🛠 [BE-031] Серверное автосохранение черновика попытки (Quiz Draft Autosave) и восстановление активной сессии

> **Приоритет:** High (P1)  
> **Связанные задачи:** FE-032, QA-032  
> **Целевой модуль:** `internal/interfaces/handlers/quiz/`, `internal/service/quiz/`  
> **Документация модуля:** [internal/service/README.md](../../internal/service/README.md)

---

## 🎯 Цель задачи
Обеспечить сохранность ответов студента при сбоях сети, случайном закрытии вкладки или перезагрузке устройства.
Реализовать серверное хранение промежуточного состояния незавершенной попытки:
1. Эндпоинт периодического фонового сохранения драфта ответов `PATCH /api/v1/lessons/{id}/attempts/{attempt_id}/draft`.
2. Эндпоинт проверки наличия активной попытки `GET /api/v1/lessons/{id}/attempts/active`. Если у студента есть незаконченная попытка с неистекшим таймером, сервер возвращает ее идентификатор, время старта и сохраненные ответы.

---

## 🔍 Текущее состояние кода
- Таблица `quiz_attempt_answers` сохраняет ответы только постфактум при финализации теста.
- При перезагрузке страницы студент теряет все введенные данные, если они не зафиксированы на сервере.

---

## 📝 Технические требования к реализации

### 1. Схема данных в БД
В таблицу `quiz_attempts` добавить поле для JSON-снапшота черновика:
```sql
ALTER TABLE quiz_attempts ADD COLUMN IF NOT EXISTS draft_answers JSONB DEFAULT '{}'::jsonb;
ALTER TABLE quiz_attempts ADD COLUMN IF NOT EXISTS current_step INT DEFAULT 1;
```

### 2. API / Endpoints
- `PATCH /api/v1/lessons/{id}/attempts/{attempt_id}/draft`
  - **Request Body:**
    ```json
    {
      "current_step": 3,
      "answers": {
        "QuizSingleBlock-1": 2,
        "QuizMultiBlock-2": [0, 3],
        "QuizMatchBlock-3": { "0": "AST", "1": "Puck" }
      }
    }
    ```
  - **Response:** `200 OK`, `{ "saved_at": "2026-10-04T12:00:00Z" }`.

- `GET /api/v1/lessons/{id}/attempts/active`
  - Возвращает активную попытку студента:
    ```json
    {
      "has_active_attempt": true,
      "attempt": {
        "id": 18,
        "started_at": "2026-10-04T11:55:00Z",
        "time_limit_minutes": 20,
        "remaining_seconds": 900,
        "current_step": 3,
        "draft_answers": { ... }
      }
    }
    ```
  - Если активных попыток нет (или время истекло) ➔ `{ "has_active_attempt": false }`.

---

## ✅ Критерии приёмки (Definition of Done)
- [x] Сервер сохраняет черновик ответов в `quiz_attempts.draft_answers`.
- [x] Запрос активной попытки корректно рассчитывает `remaining_seconds` от серверного времени `NOW() - started_at`.
- [x] Попытка с истекшим временем автоматически помечается как `timed_out` и не возвращается как активная.
