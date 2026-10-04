# 📦 Сервис: `internal/service/progress`
> **Путь:** `internal/service/progress`  
> **Роль:** Отслеживание процесса обучения студента, учет времени просмотра уроков, завершение уроков, управление попытками тестирования и расчет общего прогресса курса.

---

## 🎯 Назначение и ответственность
Сервис фиксирует активность студента: отмечает начало прохождения урока, сохраняет позицию и таймлайн просмотра, обрабатывает завершение урока, валидирует проверочные тесты (anti-cheat guard), управляет попытками тестирования (Pre-flight API, лимиты, старт попыток), сохраняет домашние задания/эссе и автоматически пересчитывает процент прохождения и средний балл курса.

---

## ⚠️ Жестко установленные правила (Hard Invariants)
1. **Транзакционность завершения урока (`CompleteLesson`):** Все шаги завершения (смена статуса урока на `completed`, сохранение эссе, подсчет завершенных уроков курса, пересчет процента и сохранение в `course_progress`) выполняются строго в рамках одной транзакции `txManager.WithTX`.
2. **Идемпотентный Upsert прогресса курса:** Если записи `course_progress` еще нет, она создается при первом старте любого урока. При завершении уроков используется `UpsertCourseProgressWithScore`, предотвращающий гонки и дубликаты.
3. **Формула расчета процента курса:**  
   $$\text{Percent} = \frac{\text{CompletedLessons}}{\text{TotalLessons}} \times 100.0$$
4. **Учет среднего балла:** Средний балл курса рассчитывается по всем завершенным урокам с выставленной оценкой (`score != nil`).
5. **Серверный Scoring Guard и Anti-Cheat:** Для уроков, содержащих тесты (тип `test`/`quiz`, блоки `Quiz...` в `lesson.content` или привязанные `quizzes`), клиентское поле `score` игнорируется. Оценка вычисляется строго сервером на основе правильности ответов или верифицированной попытки `attempt_id`. При досрочном завершении (`is_abandoned = true`) неотвеченные вопросы дают 0 баллов.
6. **Best Score Preservation (Сохранение высшего балла):** В общем прогрессе урока (`lesson_progress.score`) сохраняется максимальный балл среди всех попыток студента (`GREATEST(current, previous)`). Заброшенная или слабая повторная попытка никогда не уменьшает балл и не сбрасывает статус успешно завершенного урока.
7. **Контроль лимита попыток:** При старте новой попытки (`StartLessonAttempt`) проверяется `max_attempts`. Если количество исчерпано, возвращается 403 Forbidden.

---

## 📁 Структура файлов
| Файл | Описание |
|---|---|
| [`service.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/progress/service.go) | Определение сервиса и конструктор `NewProgressService` |
| [`startLesson.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/progress/startLesson.go) | Инициализация прогресса курса и фиксация статуса `in_progress` для урока |
| [`updateLessonProgress.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/progress/updateLessonProgress.go) | Накопление времени просмотра (`additionalTime`) и сохранение последней позиции |
| [`completeLesson.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/progress/completeLesson.go) | Серверный пересчет баллов, защита от читерства, Best Score Preservation и пересчет курса |
| [`attempts.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/progress/attempts.go) | Управление попытками тестирования: старт попытки с проверкой лимита и Pre-flight сводка |
| [`attempts_test.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/progress/attempts_test.go) | Unit-тесты Pre-flight сводки, проверки лимита попыток и сохранения Best Score |
| [`complete_lesson_test.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/progress/complete_lesson_test.go) | Unit-тесты anti-cheat guard, лекций, заброшенных тестов и Best Score Preservation |
| [`getProgress.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/progress/getProgress.go) | Чтение сводного прогресса курса или конкретного урока |
| [`get_progress_test.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/progress/get_progress_test.go) | Unit-тесты чтения прогресса и гидратации сданных заданий |

---

## ⚙️ Функции, методы и API
| Функция / Метод | Файл:Строки | Описание | Сигнатура |
|---|---|---|---|
| `StartLesson` | [`startLesson.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/progress/startLesson.go) | Старт урока, создание записей прогресса | `StartLesson(ctx context.Context, userID, lessonID int64) error` |
| `UpdateLessonProgress` | [`updateLessonProgress.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/progress/updateLessonProgress.go) | Обновление времени и позиции проигрывания | `UpdateLessonProgress(ctx context.Context, userID, lessonID int64, input domain.UpdateProgressInput) error` |
| `CompleteLesson` | [`completeLesson.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/progress/completeLesson.go) | Серверный расчет оценки, завершение урока, Best Score Preservation и пересчет % курса | `CompleteLesson(ctx context.Context, userID, lessonID int64, input domain.CompleteLessonInput) (*domain.LessonCompletionResult, error)` |
| `StartLessonAttempt` | [`attempts.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/progress/attempts.go) | Старт новой попытки тестирования с валидацией лимита | `StartLessonAttempt(ctx context.Context, userID, lessonID int64) (*domain.StartAttemptResult, error)` |
| `GetLessonAttemptsSummary` | [`attempts.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/progress/attempts.go) | Pre-flight сводка: количество попыток, лимит, лучший балл и история | `GetLessonAttemptsSummary(ctx context.Context, userID, lessonID int64) (*domain.LessonAttemptsSummary, error)` |
| `GetCourseProgress` | [`getProgress.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/progress/getProgress.go) | Получение процента и времени по курсу | `GetCourseProgress(ctx context.Context, userID, courseID int64) (*domain.CourseProgress, error)` |

---

## 🔗 Зависимости
- **Входящие:** [`handlers/progress`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/progress), [`service/quiz`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/quiz).
- **Исходящие:**
  - [`repository.ProgressRepository`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/progress)
  - [`repository.LessonRepository`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/lesson)
  - [`repository.QuizRepository`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/quiz)
  - [`infrastructure.txmanager.TransactionManager`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/infrastructure/txmanager)

---

## 🧭 Навигатор типовых задач (Where to edit?)
- **Изменить логику лимитов или Pre-flight сводки:** [`attempts.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/progress/attempts.go).
- **Изменить логику валидации квизов или расчет проходного балла:** [`completeLesson.go#L95-L148`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/progress/completeLesson.go#L95-L148).
- **Изменить формулу расчета общего прогресса:** отредактировать блок расчета в [`completeLesson.go#L210-L245`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/progress/completeLesson.go#L210-L245).
