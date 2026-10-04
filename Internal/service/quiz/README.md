# 📦 Сервис: `internal/service/quiz`
> **Путь:** `internal/service/quiz`  
> **Роль:** Создание тестов/квизов, контроль попыток, тайм-лимитов, автопроверка и ручная оценка преподавателем.

---

## 🎯 Назначение и ответственность
Сервис реализует систему тестирования: создание вопросов разных типов (`single_choice`, `multiple_choice`, `open_text`), запуск попытки с контролем лимитов, автоматическую проверку закрытых тестов, изоляцию открытых вопросов для ручной проверки преподавателем и фиксацию сдачи теста в прогрессе урока.

---

## ⚠️ Жестко установленные правила (Hard Invariants)
1. **Защита от состояния гонки (Advisory Locks):** В методе `StartAttempt` перед проверкой лимита попыток (`MaxAttempts`) берется блокировка `AcquireAdvisoryLock(userID, quizID)` в рамках транзакции, исключая параллельный запуск лишних попыток.
2. **Пессимистическая блокировка попытки (`FOR UPDATE`):** Методы `SubmitAttempt` и `GradeAttemptAnswer` блокируют строку попытки через `GetAttemptForUpdate`, исключая повторную сдачу или параллельное выставление оценок.
3. **Контроль лимита времени (`TimeLimit`):** Если у теста установлен тайм-лимит, сдача попытки за пределами отведенного времени завершается ошибкой `ErrTimeLimitExceeded`.
4. **Гибридный скоринг (Авто + Ручной):**
   - Вопросы с вариантами ответов проверяются автоматически.
   - При наличии хотя бы одного открытого вопроса (`open_text`) попытка помечается флагом `NeedsGrading = true`, а статус сдачи откладывается до проверки преподавателем.
5. **Автоматическое завершение урока при успехе:** Когда все вопросы проверены (`CountUngradedAnswers == 0`) и балл превышает проходной (`attempt.Passed == true`), автоматически вызывается `progressService.CompleteLesson`, засчитывая урок студенту.
6. **Валидация диапазона выставляемых баллов:** При оценке работы студента (`GradeAttemptAnswer`) выставляемый балл обязан находиться в диапазоне `0 <= points <= maxPoints`. При попытке выставить отрицательный балл или значение выше веса задания возвращается `ErrInvalidGradePoints` (HTTP 400 Bad Request).

---

## 📁 Структура файлов
| Файл | Описание |
|---|---|
| [`service.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/quiz/service.go) | Определение сервиса и конструктор `NewQuizService` |
| [`createQuiz.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/quiz/createQuiz.go) | Создание структуры теста с вопросами и вариантами ответов |
| [`startAttempt.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/quiz/startAttempt.go) | Старт попытки с advisory lock и проверкой `MaxAttempts` |
| [`submitAttempt.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/quiz/submitAttempt.go) | Прием ответов, автоскоринг и перевод в ручную проверку |
| [`gradeAttemptAnswer.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/quiz/gradeAttemptAnswer.go) | Выставление оценки преподавателем с валидацией диапазона [0, maxPoints], отзыв и финальный пересчет сдачи |
| [`listAttemptsForGrading.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/quiz/listAttemptsForGrading.go) | Выборка попыток курса, ожидающих оценки |
| [`homework_feedback.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/quiz/homework_feedback.go) | Получение студентом результатов проверки и рецензии домашнего задания |
| [`sanitizer.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/quiz/sanitizer.go) | Античит санитизация контента урока Puck и серверная валидация квизов |

---

## ⚙️ Функции, методы и API
| Функция / Метод | Файл:Строки | Описание | Сигнатура |
|---|---|---|---|
| `CreateQuiz` | [`createQuiz.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/quiz/createQuiz.go) | Валидация и сохранение квиза с вопросами | `CreateQuiz(ctx context.Context, userID int64, quiz *domain.Quiz) (*domain.Quiz, error)` |
| `StartAttempt` | [`startAttempt.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/quiz/startAttempt.go) | Запуск попытки с транзакционной блокировкой | `StartAttempt(ctx context.Context, userID, quizID int64) (*domain.QuizAttempt, error)` |
| `SubmitAttempt` | [`submitAttempt.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/quiz/submitAttempt.go) | Сдача попытки, автопроверка тестов | `SubmitAttempt(ctx context.Context, userID, attemptID int64, answers []domain.QuizAttemptAnswer) (*domain.QuizAttempt, error)` |
| `GradeAttemptAnswer` | [`gradeAttemptAnswer.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/quiz/gradeAttemptAnswer.go) | Проверка ответа учителем с валидацией баллов [0, maxPoints], завершение урока при успехе | `GradeAttemptAnswer(ctx context.Context, teacherID, attemptID, answerID int64, points int, feedback *string) (*domain.QuizAttempt, error)` |
| `GetStudentHomeworkFeedback` | [`homework_feedback.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/quiz/homework_feedback.go) | Получение результатов проверки и рецензии ДЗ студентом | `GetStudentHomeworkFeedback(ctx context.Context, userID, lessonID int64) (*domain.StudentHomeworkFeedback, error)` |
| `SanitizeLessonContentForStudent` | [`sanitizer.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/quiz/sanitizer.go) | Очистка правильных ответов Puck для студентов | `SanitizeLessonContentForStudent(contentJSON string) (string, error)` |
| `ValidateQuizSubmission` | [`sanitizer.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/quiz/sanitizer.go) | Серверный расчет баллов по ответам на квизы | `ValidateQuizSubmission(contentJSON string, answers []domain.LessonAnswerSubmission) (*domain.LessonCompletionResult, error)` |

---

## 🔗 Зависимости
- **Входящие:** [`handlers/quiz`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/quiz).
- **Исходящие:**
  - [`repository.QuizRepository`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/quiz)
  - [`service.ProgressServices`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/progress)
  - [`service.AccessService`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/access)
  - [`infrastructure.txmanager.TransactionManager`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/infrastructure/txmanager)

---

## 🧭 Навигатор типовых задач (Where to edit?)
- **Изменить логику вычисления проходного балла:** отредактировать `domain/quiz.go` и метод `finalizeGrading` в [`gradeAttemptAnswer.go#L87-L110`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/quiz/gradeAttemptAnswer.go#L87-L110).
- **Добавить штраф по времени за опоздание со сдачей:** метод `checkTimeLimit` в [`submitAttempt.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/quiz/submitAttempt.go).
