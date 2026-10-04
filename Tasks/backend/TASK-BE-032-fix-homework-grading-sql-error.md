# 🛠 [BE-032] Устранение 500 ошибки при утверждении оценки ДЗ (SQL Column Points & Transaction Abort Fix)

> **Статус:** Completed  
> **Приоритет:** Critical (P0)  
> **Связанные задачи:** QA-033  
> **Целевой модуль:** `internal/repository/quiz/`, `internal/service/quiz/`  
> **Документация модуля:** [internal/service/quiz/README.md](../../internal/service/README.md)

---

## 🎯 Цель задачи
Устранить ошибку HTTP 500 Internal Server Error, возникающую при попытке преподавателя выставить оценку и сохранить домашнее задание студента через модальное окно проверки (`ModalGradeHW`).

---

## 🔍 Коренная причина ошибки (Root Cause Analysis)

1. **Несуществующая колонка в таблице `quiz_questions`:**
   В методе `GetQuizTotalPoints` ([`internal/repository/quiz/quiz.go`](../../internal/repository/quiz/quiz.go)) выполнялся SQL-запрос:
   ```sql
   SELECT COALESCE(SUM(points), 1) FROM quiz_questions WHERE quiz_id = $1
   ```
   В схеме БД таблица `quiz_questions` (`migrators/20251001000010_create_quiz_questions.sql`) содержит только поля `id, quiz_id, question_text, explanation, position, created_at, updated_at`. Поле `points` в ней **отсутствует** (оно определено в родительской таблице `quizzes.points`).
   В результате запрос в PostgreSQL падал с фатальной ошибкой:
   `ERROR: column "points" does not exist (SQLSTATE 42703)`.

2. **Состояние Aborted Transaction (`SQLSTATE 25P02`):**
   В методе `finalizeGrading` ([`internal/service/quiz/gradeAttemptAnswer.go`](../../internal/service/quiz/gradeAttemptAnswer.go)) ошибка вызова `GetQuizTotalPoints` игнорировалась конструкцией:
   ```go
   maxPoints, err := s.quizRepo.GetQuizTotalPoints(ctx, attempt.QuizID)
   if err != nil { maxPoints = 1 }
   ```
   Хотя код на Go проигнорировал ошибку, транзакция PostgreSQL уже была переведена в состояние ошибки. Следующий запрос транзакции (`UpdateAttempt`) немедленно падал с:
   `ERROR: current transaction is aborted, commands ignored until end of transaction block (SQLSTATE 25P02)`.

3. **Некорректный подсчет баллов попытки в `SumAttemptPoints`:**
   В [`internal/repository/quiz/answer.go`](../../internal/repository/quiz/answer.go) метод `SumAttemptPoints` делал `COUNT(*)` вместо `SUM(points)`, из-за чего оценка за эссе (например, 20 баллов) приводилась к 1 баллу.

4. **Ошибка в `GetAnswerPointsAndCorrectness`:**
   Выборка `qq.points` падала по аналогичной причине отсутствия поля в `quiz_questions`.

---

## 📝 Реализованные изменения

1. **В `internal/repository/quiz/quiz.go`:**
   Заменен SQL-запрос в `GetQuizTotalPoints` на выборку из таблицы `quizzes`:
   ```go
   err := q.QueryRow(ctx, "SELECT COALESCE(points, 25) FROM quizzes WHERE id = $1", quizID).Scan(&maxPoints)
   ```

2. **В `internal/repository/quiz/answer.go`:**
   - В `SumAttemptPoints`:
     ```sql
     SELECT COALESCE(SUM(points), 0) FROM quiz_attempt_answers WHERE attempt_id = $1 AND is_correct = TRUE
     ```
   - В `GetAnswerPointsAndCorrectness`:
     Добавлен JOIN с `quizzes`:
     ```sql
     SELECT qa.is_correct, COALESCE(qz.points, 1) 
     FROM quiz_answers qa 
     JOIN quiz_questions qq ON qa.question_id = qq.id 
     JOIN quizzes qz ON qq.quiz_id = qz.id 
     WHERE qa.id = $1
     ```

3. **В `internal/service/quiz/gradeAttemptAnswer.go`:**
   В `finalizeGrading` добавлена строгая проверка ошибки:
   ```go
   maxPoints, err := s.quizRepo.GetQuizTotalPoints(ctx, attempt.QuizID)
   if err != nil {
       return fmt.Errorf("get quiz total points: %w", err)
   }
   ```

---

## ✅ Критерии приёмки (Definition of Done)
- [x] Эндпоинт `POST /api/v1/quizzes/attempts/{attempt_id}/answers/{answer_id}/grade` отрабатывает со статусом 200 OK.
- [x] Оценка и отзыв преподавателя сохраняются в `quiz_attempt_answers`.
- [x] Транзакция фиксируется без ошибок `SQLSTATE 25P02`.
- [x] При оценке всех вопросов попытка финализируется, рассчитывается процент `score` и обновляется прогресс студента.
