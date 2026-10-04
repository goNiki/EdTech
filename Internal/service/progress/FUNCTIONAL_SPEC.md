# 📋 Функциональная спецификация: `Успеваемость и прогресс (Progress)`

> **Расположение:** `internal/service/progress`  
> **Технический контекст:** [`README.md`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/progress/README.md)  
> **Дата актуализации:** 2026-10-04  

---

## 🎯 Каталог бизнес-фич модуля
| Фича | Описание возможности | Обеспечивающие функции |
|---|---|---|
| **Старт прохождения урока** | Фиксация начала обучения по уроку и регистрация общего прогресса по курсу | `StartLesson()` |
| **Телеметрия обучения** | Накопительный учет проведенного времени (`time_spent`) и сохранение позиции плеера (`last_position`) | `UpdateLessonProgress()` |
| **Завершение урока и сдача эссе** | Отметка урока как пройденного, фиксация балла, сохранение открытых ответов/эссе и пересчет агрегата курса | `CompleteLesson()` |
| **Мониторинг прогресса студента** | Получение среза прохождения конкретного урока, общего прогресса курса или всех уроков потока | `GetLessonProgress()`, `GetCourseProgress()`, `GetAllLessonProgress()` |

---

## 🔬 Паспорта функций

### ⚡ Функция: `StartLesson(ctx, userID, lessonID)`

* **Файл и строки:** [`startLesson.go#L14-L68`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/progress/startLesson.go#L14-L68)
* **Бизнес-назначение:** Первичная инициализация прохождения урока учащимся.
* **Связанная фича:** *Старт прохождения урока*

#### 📥 Входные параметры
| Параметр | Тип | Обязателен | Бизнес-смысл и ограничения |
|---|---|:---:|---|
| `userID` | `int64` | Да | Идентификатор студента |
| `lessonID` | `int64` | Да | Идентификатор запускаемого урока |

#### 🔄 Пошаговый алгоритм работы
1. **Шаг 1 (Получение контекста):** Запрашивает урок по ID для выяснения привязки к курсу (`lesson.CourseID`).
2. **Шаг 2 (Определение базы курса):** Запрашивает все уроки курса для вычисления общего знаменателя `totalLessons` (если уроков нет, по умолчанию 1).
3. **Шаг 3 (Транзакционная фиксация):** В транзакции `WithTX`:
   - Создает запись `course_progress` (если еще не создана): `TotalLessons`, `CompletedLess: 0`, `Percent: 0`, `TotalWatchTime: 0`, `StartedAt: now`, `LastAccessedAt: now`.
   - Создает запись `lesson_progress`: `Status = in_progress`, `StartedAt: now`.

#### ⚠️ Побочные эффекты (Side Effects)
* **База данных:** Создание строк в `course_progress` и `lesson_progress`.

---

### ⚡ Функция: `UpdateLessonProgress(ctx, userID, lessonID, input)`

* **Файл и строки:** [`updateLessonProgress.go#L13-L35`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/progress/updateLessonProgress.go#L13-L35)
* **Бизнес-назначение:** Фоновый сбор телеметрии просмотра видеоматериалов/лекций с фронтенда.
* **Связанная фича:** *Телеметрия обучения*

#### 📥 Входные параметры
| Параметр | Тип | Обязателен | Бизнес-смысл и ограничения |
|---|---|:---:|---|
| `userID` | `int64` | Да | Идентификатор студента |
| `lessonID` | `int64` | Да | Идентификатор урока |
| `input.TimeSpent` | `int` | Нет | Дельта времени в секундах (накапливается к общему времени урока) |
| `input.LastPosition` | `int` | Нет | Текущий таймкод/позиция в секундах |
| `input.Status` | `ProgressStatus` | Нет | Опциональное обновление статуса (`not_started`, `in_progress`, `completed`) |

#### 🔄 Пошаговый алгоритм работы
1. В транзакции при наличии `TimeSpent > 0` или `LastPosition > 0` выполняет инкремент времени и обновление позиции.
2. При наличии непустого `Status` обновляет статус прохождения урока.

---

### ⚡ Функция: `CompleteLesson(ctx, userID, lessonID, input)`

* **Файл и строки:** [`completeLesson.go#L24-L81`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/progress/completeLesson.go#L24-L81)
* **Бизнес-назначение:** Защищенное завершение урока с серверной верификацией ответов квизов, предотвращением читерства (авто-100 баллов), сохранением наивысшего балла (Best Score Preservation) и пересчетом прогресса курса.
* **Связанная фича:** *Завершение урока и сдача эссе / Anti-Cheat Quiz Guard*

#### 📥 Входные параметры
| Параметр | Тип | Обязателен | Бизнес-смысл и ограничения |
|---|---|:---:|---|
| `userID` | `int64` | Да | Идентификатор студента |
| `lessonID` | `int64` | Да | Идентификатор завершаемого урока |
| `input.Score` | `*int` | Нет | Оценка от клиента (учитывается ТОЛЬКО для лекций; для тестов игнорируется) |
| `input.Answers` | `[]LessonAnswerSubmission` | Нет | Ответы на Puck-блоки тестов |
| `input.Essays` | `[]EssaySubmission` | Нет | Массив развернутых ответов/эссе на проверку преподавателю |
| `input.AttemptID` | `*int64` | Нет | ID связанной попытки реляционного теста (`quiz_attempts`) |
| `input.IsAbandoned` | `bool` | Нет | Флаг досрочного выхода из тестирования (неотвеченные = 0 баллов) |

#### 🔄 Пошаговый алгоритм работы
1. **Шаг 1 (Определение характера урока и настроек тестирования):**
   - Запрашивает настройки урока `lesson.GetQuizSettings()` (`TimeLimitMinutes`, `QuestionTimeLimitSeconds`, `MaxAttempts`, `PassingScorePercent`, `FeedbackMode`, `ShuffleQuestions`).
   - Проверяет наличие Puck-блоков тестирования через `quiz.ValidateQuizSubmission`.
   - Проверяет тип урока (`lesson.Type == 'test' || 'quiz'`) и наличие привязанных тестов в `quizzes`.
2. **Шаг 2 (Серверная валидация времени и Grace Period):**
   - Если `quiz_settings.TimeLimitMinutes > 0`, сервер рассчитывает `timeLimit = timeLimitMinutes + 15s (Grace Period)`.
   - При наличии `attempt_id`: если `time.Since(attempt.StartedAt) > timeLimit`, фиксируется статус `timed_out`, попытка отклоняется с 0 баллов, `is_passed = false`.
   - При отсутствии `attempt_id` и `TimeSpent > timeLimit`: статус фиксируется как `timed_out`, `is_passed = false`.
3. **Шаг 3 (Серверный расчет оценки и Anti-Cheat Guard):**
   - **Лекция без тестов:** устанавливается балл 100 (или переданный `input.Score`), `is_passed = true`.
   - **Урок с тестами:** переданный клиентом `input.Score` **игнорируется**.
     - Если передан `attempt_id`: проверяется принадлежность студенту и уроку, рассчитывается результат попытки. При `is_abandoned = true` фиксируется фактический балл за решенные задания. Сдача засчитывается при `score >= passing_score_percent` (дефолт 70%).
     - Если присутствуют Puck-блоки: балл вычисляется строго из `validationResult.Score` (при отсутствии ответов балл = 0). Сдача засчитывается при `score >= passing_score_percent`.
     - Если тест без блоков и без `attempt_id`: балл = 0, `is_passed = false`.
4. **Шаг 4 (Best Score Preservation):**
   - Если у студента в `lesson_progress` уже был зафиксирован более высокий результат по этому уроку, в прогрессе сохраняется максимум: `progressScore = max(current, previous)`. Ранее сданный статус `completed` сохраняется.
5. **Шаг 5 (Транзакционная фиксация):** В `WithTX`:
   - Обновляет или создает запись `lesson_progress` (статус `targetStatus`, балл `progressScore`).
   - Сохраняет эссе/домашние задания через `quizRepo.SaveEssaySubmission` (если не было таймаута).
   - Пересчитывает агрегатный прогресс курса (`completedCount`, `percent`, `avgScore` с учетом нулевых оценок) и обновляет `course_progress`.
6. **Шаг 6 (Возврат детализации и режим `exam_blind`):**
   - В режиме `feedback_mode == "exam_blind"` из ответа полностью удаляются правильные ответы (`correct_answer`), флаги верности и пояснения (`results = nil`).
   - Возвращает `LessonCompletionResult` с полями `lesson_id`, `status` (`completed`, `in_progress` или `timed_out`), `score`, `earned_points`, `total_max_points`, `is_passed` и картой верификации блоков.

#### ⚠️ Побочные эффекты (Side Effects)
* **База данных:**
  - Обновление/создание `lesson_progress` (статус `completed`, дата завершения, оценка `progressScore`).
  - Финализация `quiz_attempts` (при наличии `attempt_id`).
  - Создание записей в таблице ответов/эссе (при наличии домашних заданий).
  - Пересчет и обновление строки в `course_progress` (`completed_lessons`, `percent`, `avg_score`, `completed_at`).

---

### ⚡ Функция: `GetLessonProgress(ctx, userID, lessonID)`

* **Файл и строки:** [`getProgress.go#L14-L46`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/progress/getProgress.go#L14-L46)
* **Бизнес-назначение:** Предоставление актуального состояния прохождения урока студентом, включая статус, балл, затраченное время и детальную историю сданных домашних заданий/эссе с комментариями преподавателя.
* **Связанная фича:** *Мониторинг прогресса студента*

#### 📥 Входные параметры
| Параметр | Тип | Обязателен | Бизнес-смысл и ограничения |
|---|---|:---:|---|
| `userID` | `int64` | Да | Идентификатор студента |
| `lessonID` | `int64` | Да | Идентификатор запрашиваемого урока |

#### 🔄 Пошаговый алгоритм работы
1. **Шаг 1 (Запрос прогресса):** Обращается к `progressRepo.GetLessonProgress`.
2. **Шаг 2 (Обработка отсутствия прогресса):** Если запись прогресса не найдена (`ErrLessonProgressNotFound` или `pgx.ErrNoRows`), не падает с ошибкой, а формирует дефолтный объект со статусом `not_started`, `Score = nil`, `TimeSpent = 0`, `LastPos = 0` и пустым списком сданных заданий `submissions: []`.
3. **Шаг 3 (Гидратация истории сданных заданий):** Через `quizRepo.GetLessonSubmissions` запрашивает все сданные ответы студента по уроку (формулировка вопроса, ответ студента, начисленные баллы, максимальные баллы, комментарий преподавателя `teacher_feedback`, статус проверки `is_graded`).
4. **Шаг 4 (Возврат данных):** Прикрепляет список `Submissions` к прогрессу урока и возвращает клиенту (200 OK).

---

### ⚡ Функции чтения: `GetCourseProgress()`, `GetAllLessonProgress()`

* **Файл и строки:** [`getProgress.go#L48-L68`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/progress/getProgress.go#L48-L68)
* **Бизнес-назначение:** Предоставление агрегированных срезов прогресса курса или всех уроков потока.
* **Связанная фича:** *Мониторинг прогресса студента*

---

### ⚡ Функция: `GetLessonAttemptsSummary(ctx, userID, lessonID)`

* **Файл и строки:** [`attempts.go#L69-L151`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/progress/attempts.go#L69-L151)
* **Бизнес-назначение:** Предоставление данных для Pre-flight экрана перед началом тестирования: количество совершенных попыток, разрешенный лимит, возможность начать новую попытку, наивысший балл и хронологическая история попыток.
* **Связанная фича:** *Pre-flight экран и история попыток (Best Score)*

#### 📥 Входные параметры
| Параметр | Тип | Обязателен | Бизнес-смысл и ограничения |
|---|---|:---:|---|
| `userID` | `int64` | Да | Идентификатор студента |
| `lessonID` | `int64` | Да | Идентификатор урока с тестом |

#### 🔄 Пошаговый алгоритм работы
1. **Шаг 1:** Проверяет существование урока в базе данных.
2. **Шаг 2:** Извлекает настройки квиза (`quizzes`), включая проходной порог `passing_score` (дефолт 70) и лимит попыток `max_attempts` (0 = безлимитно).
3. **Шаг 3:** Извлекает все попытки студента по уроку (`ListUserAttemptsByLessonID`) и текущий сохраненный балл из `lesson_progress`.
4. **Шаг 4:** Агрегирует статистику:
   - Подсчитывает завершенные попытки и формирует массив хронологической истории `AttemptsHistory`.
   - Вычисляет наивысший балл `BestScore = max(attempts.score, progress.score)`.
   - Проверяет статус `IsPassed = BestScore >= PassingThreshold` (или если статус урока уже `completed`).
   - Проверяет флаг `CanStartNewAttempt = (MaxAttemptsAllowed <= 0 || TotalAttemptsMade < MaxAttemptsAllowed)`.
5. **Шаг 5:** Возвращает структуру `LessonAttemptsSummary`.

---

### ⚡ Функция: `StartLessonAttempt(ctx, userID, lessonID)`

* **Файл и строки:** [`attempts.go#L13-L67`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/progress/attempts.go#L13-L67)
* **Бизнес-назначение:** Инициализация новой попытки прохождения теста с валидацией лимита попыток.
* **Связанная фича:** *Старт попытки тестирования*

#### 📥 Входные параметры
| Параметр | Тип | Обязателен | Бизнес-смысл и ограничения |
|---|---|:---:|---|
| `userID` | `int64` | Да | Идентификатор студента |
| `lessonID` | `int64` | Да | Идентификатор запускаемого теста |

#### 🔄 Пошаговый алгоритм работы
1. **Шаг 1:** Проверяет существование урока в БД.
2. **Шаг 2:** Находит привязанный тест в таблице `quizzes` (при отсутствии — регистрирует дефолтный квиз урока).
3. **Шаг 3 (Контроль лимита):** Подсчитывает число существующих попыток студента (`CountUserAttempts`). Если `max_attempts > 0` и `count >= max_attempts`, возвращает ошибку `ErrForbidden` (HTTP 403 Forbidden: «Лимит попыток исчерпан»).
4. **Шаг 4:** Создает новую строку в `quiz_attempts` со статусом `started_at = NOW()`, `score = 0`, `passed = false`.
5. **Шаг 5:** Возвращает `StartAttemptResult` с идентификатором попытки `AttemptID` и временем старта.

---

## 💡 Подсказка для аналитика (Где менять логику?)
* *Изменить лимиты и логику Pre-flight сводки:* [`attempts.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/progress/attempts.go).
* *Изменить дефолтный балл за завершение лекции без тестирования:* [`completeLesson.go#L95-L105`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/progress/completeLesson.go#L95-L105).
* *Формула вычисления процента завершения курса:* [`completeLesson.go#L210-L245`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/progress/completeLesson.go#L210-L245).
* *Политика дефолтной длины курса при отсутствии уроков:* [`startLesson.go#L28-L30`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/progress/startLesson.go#L28-L30).

