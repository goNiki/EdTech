# 📦 Репозиторий: `internal/repository/quiz`
> **Путь:** `internal/repository/quiz`  
> **Роль:** SQL-запросы к таблицам тестов (`quizzes`, `quiz_questions`, `quiz_answers`, `quiz_attempts`, `quiz_attempt_answers`).

---

## 🎯 Назначение и ответственность
Обеспечивает хранение тестов, вопросов, вариантов ответов, пользовательских попыток и отдельных ответов студентов. Реализует транзакционные advisory-блокировки PostgreSQL для исключения гонок при старте тестов и пессимистические блокировки попыток при проверке.

---

## ⚠️ Жестко установленные правила (Hard Invariants)
1. **Advisory блокировки PostgreSQL (`pg_advisory_xact_lock`):** Метод `AcquireAdvisoryLock(userID, quizID)` захватывает транзакционную блокировку, которая снимается автоматически только при `COMMIT`/`ROLLBACK`.
2. **Пессимистический Lock попытки (`FOR UPDATE`):** Метод `GetAttemptForUpdate` выполняет `SELECT ... FOR UPDATE`, гарантируя, что ни один параллельный поток не изменит попытку во время расчета скоринга или проверки.
3. **Открытые ответы и эссе:** Для вопросов с ручной проверкой и эссе поле `is_correct` сохраняется как `NULL` в `quiz_attempt_answers`, что служит триггером статуса `NeedsGrading`.
4. **Агрегация баллов в базе:** Методы `SumAttemptPoints` и `GetQuizTotalPoints` выполняют суммирование баллов на стороне PostgreSQL через `SUM(points)`.

---

## 📁 Структура файлов
| Файл | Описание |
|---|---|
| [`repository.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/quiz/repository.go) | Фабрика `NewQuizRepository` |
| [`quiz.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/quiz/quiz.go) | `CreateQuiz`, `GetQuizByID`, подсчет максимальных баллов теста |
| [`attempt.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/quiz/attempt.go) | Создание, чтение (`FOR UPDATE`), подсчет попыток пользователя и списки на грейдинг |
| [`answer.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/quiz/answer.go) | Пакетная вставка ответов (`CreateBatchAnswers`), обновление оценки учителем (`UpdateAttemptAnswer`) |
| [`lock.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/quiz/lock.go) | Вызов `SELECT pg_advisory_xact_lock($1, $2)` |
| [`essay.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/quiz/essay.go) | Сохранение эссе/домашних работ как попыток со статусом ручной проверки |

---

## ⚙️ Функции, методы и API
| Функция / Метод | Файл:Строки | Описание | Сигнатура |
|---|---|---|---|
| `AcquireAdvisoryLock` | [`lock.go#L8-L12`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/quiz/lock.go#L8-L12) | Транзакционный лок по связке (пользователь, тест) | `AcquireAdvisoryLock(ctx, q db.QueryExecutor, userID, quizID int64) error` |
| `GetAttemptForUpdate` | [`attempt.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/quiz/attempt.go) | Выборка попытки с блокировкой строки `FOR UPDATE` | `GetAttemptForUpdate(ctx, q db.QueryExecutor, attemptID int64) (*domain.QuizAttempt, error)` |
| `CreateBatchAnswers` | [`answer.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/quiz/answer.go) | Пакетная вставка ответов студента в `quiz_attempt_answers` | `CreateBatchAnswers(ctx, q db.QueryExecutor, answers []domain.QuizAttemptAnswer) error` |
| `UpdateAttemptAnswer` | [`answer.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/quiz/answer.go) | Запись оценки преподавателя, отзыва и статуса `is_correct` | `UpdateAttemptAnswer(ctx, q, answerID, points, feedback, isCorrect) error` |
| `SaveEssaySubmission` | [`essay.go#L12-L72`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/quiz/essay.go#L12-L72) | Сохранение развернутого текстового ответа на проверку | `SaveEssaySubmission(ctx, q, userID, courseID, lessonID, essay) error` |
| `GetLessonSubmissions` | [`essay.go#L74-L122`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/quiz/essay.go#L74-L122) | Получение списка сданных работ студента с оценками и фидбеком | `GetLessonSubmissions(ctx, q, userID, lessonID int64) ([]domain.LessonSubmissionDetail, error)` |

---

## 🔗 Зависимости
- **Входящие:** [`service/quiz`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/quiz), [`service/progress`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/progress).
- **Исходящие:** [`db.QueryExecutor`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/infrastructure/db), [`models/converter`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/models/converter).
