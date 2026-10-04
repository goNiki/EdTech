# 📦 Обработчики: `internal/interfaces/handlers/quiz`
> **Путь:** `internal/interfaces/handlers/quiz`  
> **Роль:** HTTP-контроллеры тестов/квизов (создание, старт попытки, сдача теста, ручная оценка преподавателем).

---

## 🎯 Назначение и ответственность
Принимает запросы по квизам на маршрутах `/api/v1/quizzes/*`, `/api/v1/lessons/{lesson_id}/quizzes` и `/api/v1/courses/{course_id}/quizzes/*`.

---

## 📁 Структура файлов и Эндпоинты
| Файл | HTTP Метод и Путь | Описание |
|---|---|---|
| [`createQuiz.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/quiz/createQuiz.go) | `POST /api/v1/lessons/{lesson_id}/quizzes` | Создание теста к уроку |
| [`startAttempt.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/quiz/startAttempt.go) | `POST /api/v1/quizzes/{quiz_id}/attempts/start` | Старт попытки прохождения теста |
| [`submitAttempt.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/quiz/submitAttempt.go) | `POST /api/v1/quizzes/attempts/{attempt_id}/submit` | Отправка ответов на авто- и ручную проверку |
| [`gradeAttemptAnswer.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/quiz/gradeAttemptAnswer.go) | `POST /api/v1/quizzes/attempts/{attempt_id}/answers/{answer_id}/grade` | Оценка открытого ответа преподавателем с валидацией диапазона [0, maxPoints] |
| [`listAttemptsForGrading.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/quiz/listAttemptsForGrading.go) | `GET /api/v1/courses/{course_id}/quizzes/attempts` | Список попыток, требующих проверки |

---

## 🔗 Зависимости
- **Входящие:** [`internal/app/di.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/app/di.go).
- **Исходящие:** [`service.QuizServices`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/quiz), [`internal/dto`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/dto).
