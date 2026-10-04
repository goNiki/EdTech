# 🛠 [BE-034] API получения результатов проверки и рецензии ДЗ для студента

> **Приоритет:** Critical (P0)  
> **Связанные задачи:** FE-036, QA-036  
> **Целевой модуль:** `internal/interfaces/handlers/quiz/`, `internal/service/quiz/`, `internal/repository/quiz/`  
> **Документация модуля:** [internal/service/quiz/README.md](../../internal/service/quiz/README.md)

---

## 🎯 Цель задачи
Обеспечить передачу студенту подробных результатов ручной проверки его домашнего задания (эссе, загруженных файлов) преподавателем:
1. Оценка в баллах (`points` и `max_points`).
2. Статус проверки (`is_graded: true`, `is_correct: true/false`).
3. Дата и время проверки (`graded_at`).
4. Текстовая обратная связь / рецензия преподавателя (`feedback`).
5. Данные проверившего преподавателя (имя и аватар для отображения в карточке фидбека).

---

## 🔍 Текущее состояние кода
- В БД (`quiz_attempt_answers`) поля `points`, `is_correct`, `feedback` сохраняются в методе `GradeAttemptAnswer`.
- Однако эндпоинты урока `GET /api/v1/lessons/{id}` и прогресса `GET /api/v1/lessons/{id}/progress` не отдают студенту детализацию ответов и комментариев преподавателя. В результате студент видит только сухой процент прогресса, но не знает, почему ему снизили балл и что написал преподаватель.

---

## 📝 Технические требования к реализации

### 1. HTTP Endpoint
`GET /api/v1/lessons/{id}/homework-feedback`
* **Авторизация:** Bearer JWT (авторизованный студент).
* **Параметры пути:** `id` (int, ID урока).

### 2. Response DTO
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "has_submission": true,
    "status": "graded", 
    "attempt_id": 3,
    "submitted_at": "2026-10-04T12:00:00Z",
    "graded_at": "2026-10-04T15:30:00Z",
    "teacher": {
      "id": 2,
      "name": "Никита Преподаватель",
      "avatar_url": "/static/uploads/avatars/teacher.png"
    },
    "answers": [
      {
        "answer_id": 3,
        "question_text": "Напишите развернутое сочинение-рассуждение по предложенному тексту...",
        "student_answer": "В данном тексте автор поднимает важную проблему истинного гуманизма...",
        "points": 23,
        "max_points": 25,
        "is_correct": true,
        "feedback": "Отличная аргументация тезиса. Во втором абзаце обратите внимание на пунктуацию при вводных словах."
      }
    ]
  }
}
```
* Возможные значения `status`:
  - `"not_submitted"` — студент еще не отправлял работу;
  - `"pending"` — работа сдана, ожидает проверки учителем;
  - `"graded"` — работа проверена и оценена.

### 3. Бизнес-логика (Service Layer)
- Метод `GetStudentHomeworkFeedback(ctx context.Context, userID, lessonID int) (*dto.HomeworkFeedbackResponse, error)`.
- Алгоритм:
  1. Найти последнюю попытку студента `quiz_attempts` по `lesson_id` и `user_id`.
  2. Если попытки нет ➔ вернуть `"has_submission": false, "status": "not_submitted"`.
  3. Извлечь ответы на открытые вопросы (`quiz_attempt_answers`) с типом эссе / загрузка файла.
  4. Проверить статус: если есть непроверенные ответы (`is_correct IS NULL`) ➔ `"status": "pending"`.
  5. Если проверено ➔ загрузить данные преподавателя, выставившего оценку, и отдать полный отзыв.

---

## ✅ Критерии приёмки (Definition of Done)
- [x] Endpoint защищен: студент может видеть только свои собственные работы (чужие `user_id` недоступны).
- [x] В ответе передаются реальные тексты эссе, баллы, рецензии преподавателя и дата проверки.
- [x] Документация в `internal/service/quiz/FUNCTIONAL_SPEC.md` актуализирована.
