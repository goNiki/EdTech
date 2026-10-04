# 🛠 [BE-038] Валидация диапазона баллов при проверке домашних заданий (Points Validation)

> **Статус:** Completed  
> **Приоритет:** Medium (P2)  
> **Обнаружено в рамках:** [TASK-QA-033-homework-grading-tests.md](../qa/TASK-QA-033-homework-grading-tests.md)  
> **Целевой модуль:** `internal/service/quiz/`, `internal/interfaces/handlers/quiz/`  
> **Документация модуля:** [internal/service/quiz/README.md](../../internal/service/README.md)

---

## 🎯 Описание дефекта (Bug Report)
В эндпоинте выставления оценки за открытый ответ/домашнее задание студента:
`POST /api/v1/quizzes/attempts/{attempt_id}/answers/{answer_id}/grade`
отсутствует валидация входящего значения `points`. Бэкенд принимает любые значения баллов, включая отрицательные (`points: -5`) или многократно превышающие максимальный вес задания (`points: 999`), возвращает `200 OK` и сохраняет некорректное число в таблицу `quiz_attempt_answers`.

---

## 🔬 Шаги воспроизведения (Steps to Reproduce)
1. Авторизоваться под преподавателем (`nikit@mail.ru`).
2. Отправить POST-запрос на оценку попытки с заведомо некорректным баллом:
   ```bash
   curl -X POST http://localhost:8082/api/v1/quizzes/attempts/18/answers/16/grade \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer <TEACHER_ACCESS_TOKEN>" \
     -d '{"points": 999, "feedback": "Over maximum test"}'
   ```
3. Проверить ответ сервера и запись в PostgreSQL:
   ```sql
   SELECT id, attempt_id, points, is_correct, feedback FROM quiz_attempt_answers WHERE id = 16;
   ```

### Фактический результат (Actual Result)
- Сервер возвращает `200 OK`.
- В таблицу `quiz_attempt_answers` записывается `points = 999`.

### Ожидаемый результат (Expected Result)
- Сервер валидирует, что `points >= 0` и `points <= quiz.Points` (максимальный балл задания/квиза).
- При превышении или отрицательном значении сервер возвращает `400 Bad Request` с кодом ошибки валидации: `Балл не может быть меньше 0 или превышать максимальный балл задания (%d)`.

---

## 🛠 Технические требования к реализации

1. В `internal/service/quiz/gradeAttemptAnswer.go` перед вызовом `s.quizRepo.UpdateAttemptAnswer`:
   - Проверить `if points < 0`: вернуть ошибку валидации `400 Bad Request`.
   - Получить максимальный балл для квиза (`quiz.Points` или через `GetQuizTotalPoints`).
   - Если `points > maxPoints`: вернуть ошибку валидации `400 Bad Request`.
2. В DTO / валидаторе хендлера `internal/dto/quiz.go` (если применимо) добавить тег валидации `min=0`.

---

## 📋 Обязательное условие завершения (Hand-off Protocol)
После реализации фикса и покрытия unit/интеграционными тестами разработчик **обязан**:
1. Перевести статус текущей задачи в `Completed`.
2. Создать задачу на верификацию в `Tasks/qa/`:
   `Tasks/qa/TASK-QA-042-homework-grade-points-validation.md`
   с подробным описанием проверочных сценариев для тестировщика.
