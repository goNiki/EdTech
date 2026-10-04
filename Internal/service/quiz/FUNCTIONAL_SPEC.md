# 📋 Функциональная спецификация: `Тестирование и проверка знаний (Quiz)`

> **Расположение:** `internal/service/quiz`  
> **Технический контекст:** [`README.md`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/quiz/README.md)  
> **Дата актуализации:** 2026-10-04  

---

## 🎯 Каталог бизнес-фич модуля
| Фича | Описание возможности | Обеспечивающие функции |
|---|---|---|
| **Конструирование квизов** | Создание тестов с вариантами выбора и открытыми вопросами преподавателем | `CreateQuiz()` |
| **Старт попытки тестирования** | Контролируемый запуск попытки с защитой от гонок и проверкой лимита попыток | `StartAttempt()` |
| **Сдача теста и автопроверка** | Прием ответов, валидация таймлимита (+30 сек грейс-период), автопроверка тестов и перевод эссе на ручную проверку | `SubmitAttempt()` |
| **Ручная проверка заданий преподавателем** | Оценивание развернутых ответов/эссе с комментарием и автоматический зачет урока при сдаче | `GradeAttemptAnswer()` |
| **Очередь работ на проверку** | Пагинированный реестр попыток с открытыми вопросами для преподавателей курса | `ListAttemptsForGrading()` |
| **Рецензия и результаты проверки ДЗ** | Получение студентом баллов, рецензий преподавателя и карточки проверяющего по сданному заданию | `GetStudentHomeworkFeedback()` |
| **Античит санитизация урока** | Фильтрация правильных ответов, подсказок и ключей из Puck JSON для студентов | `SanitizeLessonContentForStudent()` |
| **Серверная проверка квизов** | Независимый расчет баллов и проверка ответов на бэкенде по эталонному контенту урока | `ValidateQuizSubmission()` |

---

## 🔬 Паспорта функций

### ⚡ Функция: `CreateQuiz(ctx, userID, quiz)`

* **Файл и строки:** [`createQuiz.go#L11-L48`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/quiz/createQuiz.go#L11-L48)
* **Бизнес-назначение:** Добавление проверочного тестирования к уроку преподавателем.
* **Связанная фича:** *Конструирование квизов*

#### 📥 Входные параметры
| Параметр | Тип | Обязателен | Бизнес-смысл и ограничения |
|---|---|:---:|---|
| `userID` | `int64` | Да | Идентификатор создателя |
| `quiz.LessonID` | `int64` | Да | Привязка к конкретному уроку |
| `quiz.PassingScor` | `int` | Да | Проходной балл в процентах (0-100) |
| `quiz.MaxAttempts` | `*int` | Нет | Ограничение количества попыток (nil = безлимитно) |
| `quiz.TimeLimit` | `*int` | Нет | Ограничение по времени в секундах |

#### 🔄 Пошаговый алгоритм работы
1. Валидация структуры квиза и вопросов (`quiz.Validate()`).
2. Проверка прав преподавателя через `accessService.CanEditCourse` (403 при отсутствии прав).
3. Сохранение квиза, вопросов и вариантов ответов в БД.

---

### ⚡ Функция: `StartAttempt(ctx, userID, quizID)`

* **Файл и строки:** [`startAttempt.go#L15-L62`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/quiz/startAttempt.go#L15-L62)
* **Бизнес-назначение:** Безопасный запуск прохождения теста студентом.
* **Связанная фича:** *Старт попытки тестирования*

#### 📥 Входные параметры
| Параметр | Тип | Обязателен | Бизнес-смысл и ограничения |
|---|---|:---:|---|
| `userID` | `int64` | Да | Идентификатор студента |
| `quizID` | `int64` | Да | Идентификатор теста |

#### 🔄 Пошаговый алгоритм работы
1. Запрашивает метаданные квиза (`quizRepo.GetQuizByID`).
2. Открывает транзакцию `WithTX`:
   - **Блокировка от гонок:** берет `AcquireAdvisoryLock(userID, quizID)` на уровне PostgreSQL транзакции. Исключает создание лишних параллельных попыток в обход лимита.
   - Подсчитывает число предыдущих попыток `CountUserAttemptsForUpdate`.
   - Если `quiz.MaxAttempts != nil` и `count >= *quiz.MaxAttempts` ➔ возвращает ошибку `ErrMaxAttemptsReached`.
   - Создает запись `QuizAttempt` (`Score = 0`, `Passed = false`, `StartedAt = now`).
3. Возвращает созданную попытку со стартовым временем.

---

### ⚡ Функция: `SubmitAttempt(ctx, userID, attemptID, answers)`

* **Файл и строки:** [`submitAttempt.go#L15-L73`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/quiz/submitAttempt.go#L15-L73)
* **Бизнес-назначение:** Завершение тестирования студентом, расчет набранных баллов или отправка на ручную проверку.
* **Связанная фича:** *Сдача теста и автопроверка*

#### 📥 Входные параметры
| Параметр | Тип | Обязателен | Бизнес-смысл и ограничения |
|---|---|:---:|---|
| `userID` | `int64` | Да | Идентификатор сдающего студента |
| `attemptID` | `int64` | Да | Идентификатор активной попытки |
| `answers` | `[]QuizAttemptAnswer` | Да | Список ответов на вопросы квиза |

#### 🔄 Пошаговый алгоритм работы
1. **Шаг 1 (Пессимистическая блокировка):** Выбирает попытку `GetAttemptForUpdate(attemptID)`.
2. **Шаг 2 (Валидация владельца и статуса):**
   - Проверяет `attempt.UserID == userID` (иначе `ErrForbidden`).
   - Проверяет `attempt.CompletedAt == nil` (иначе `ErrAttemptAlreadyCompleted`).
3. **Шаг 3 (Проверка лимита времени):**
   - Если `quiz.TimeLimit > 0`, допускается превышение не более 30 секунд (сетевой люфт):
     $$\text{AllowedDuration} = \text{TimeLimit} + 30\text{ сек}$$
   - При превышении ➔ `ErrTimeLimitExceeded`.
4. **Шаг 4 (Автоматический скоринг / Маркировка эссе):**
   - Для вариантов с выбором: проверяет правильность и начисляет баллы.
   - Для открытых вопросов (`AnswerID == nil`): помечает `hasOpenText = true`, `IsCorrect = nil`, `Points = 0`.
   - Сохраняет пакет ответов в `quiz_attempt_answers`.
   - Проставляет `attempt.CompletedAt = now`.
5. **Шаг 5 (Ветвление логики финализации):**
   - **Ветка А (Есть открытые вопросы):**
     - `attempt.NeedsGrading = true`
     - `attempt.Passed = false`
     - `attempt.Score = 0` (ожидает проверки учителем).
   - **Ветка Б (Только тесты с автопроверкой):**
     - Вычисляет процент: $\text{Score} = (\text{CorrectPoints} / \text{MaxPoints}) \times 100$.
     - Проверяет условие сдачи: $\text{Passed} = (\text{Score} \ge \text{PassingScore})$.
     - **Побочный эффект:** Если `Passed == true`, автоматически вызывает `progressService.CompleteLesson(...)` и обновляет прогресс по уроку.

---

### ⚡ Функция: `GradeAttemptAnswer(ctx, teacherID, attemptID, answerID, points, feedback)`

* **Файл и строки:** [`gradeAttemptAnswer.go#L14-L63`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/quiz/gradeAttemptAnswer.go#L14-L63)
* **Бизнес-назначение:** Проверка преподавателем конкретного открытого ответа студента, выставление баллов и комментария.
* **Связанная фича:** *Ручная проверка заданий преподавателем*

#### 📥 Входные параметры
| Параметр | Тип | Обязателен | Бизнес-смысл и ограничения |
|---|---|:---:|---|
| `teacherID` | `int64` | Да | Идентификатор проверяющего преподавателя |
| `attemptID` | `int64` | Да | Идентификатор проверяемой попытки |
| `answerID` | `int64` | Да | Идентификатор ответа студента |
| `points` | `int` | Да | Присужденное количество баллов ($0 \le \text{points} \le \text{maxPoints}$) |
| `feedback` | `*string` | Нет | Текстовый комментарий / ревью преподавателя |

#### 🔄 Пошаговый алгоритм работы
1. **Шаг 1 (Первичная валидация):** Проверяет `points >= 0`. При `points < 0` ➔ `ErrInvalidGradePoints` (400 Bad Request).
2. **Шаг 2 (Авторизация учителя):** Проверяет право преподавателя `accessService.CanEditCourse` для курса, которому принадлежит урок.
3. **Шаг 3 (Валидация максимального веса):** Запрашивает `maxPoints` задания/квиза через `quizRepo.GetQuizTotalPoints`. Если `points > maxPoints` ➔ `ErrInvalidGradePoints` (400 Bad Request: «Балл не может быть меньше 0 или превышать максимальный балл задания»).
4. **Шаг 4 (Оценка ответа):** Обновляет строку ответа: `points`, `feedback`, `is_correct = (points > 0)`, `graded_by = teacherID`.
5. **Шаг 5 (Проверка оставшихся неотвеченных эссе):**
   - Запрашивает `CountUngradedAnswers(attemptID)`.
   - Если `ungradedCount == 0` (все эссе попытки проверены) ➔ запускает `finalizeGrading`:
     - Суммирует все баллы попытки (`SumAttemptPoints`).
     - Рассчитывает финальный балл `Score` и флаг `Passed`.
     - Если `Passed == true`:
       - Вызывает `progressService.CompleteLesson` с итоговым баллом.
       - Вызывает `quizRepo.UpdateLessonProgressAfterQuiz`.

---

### ⚡ Функция: `ListAttemptsForGrading(ctx, userID, courseID, quizID, page, pageSize)`

* **Файл и строки:** [`listAttemptsForGrading.go#L11-L41`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/quiz/listAttemptsForGrading.go#L11-L41)
* **Бизнес-назначение:** Получение преподавателем списка сданных работ, ожидающих ручной проверки (`needs_grading = true`).
* **Связанная фича:** *Очередь работ на проверку*

---

### ⚡ Функция: `SanitizeLessonContentForStudent(contentJSON)`

* **Файл и строки:** [`sanitizer.go#L20-L55`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/quiz/sanitizer.go#L20-L55)
* **Бизнес-назначение:** Предотвращение списывания студентами путем очистки правильных ответов и подсказок из контента урока Puck перед отправкой в браузер.
* **Связанная фича:** *Античит санитизация урока*

#### 📥 Входные параметры
| Параметр | Тип | Обязателен | Бизнес-смысл и ограничения |
|---|---|:---:|---|
| `contentJSON` | `string` | Да | Исходный JSON разметки Puck Editor из базы данных |

#### 🔄 Пошаговый алгоритм работы
1. Парсит JSON структуру контента урока.
2. Для блоков `QuizSingleBlock` и `QuizMultiBlock`: удаляет поля `isCorrect`, `is_correct`, `correct`, `explain` из всех вариантов.
3. Для блоков `QuizMatchBlock`: извлекает правые части соответствий, сортирует их по алфавиту и перемешивает соответствия, удаляя `correctPair`.
4. Для блоков `QuizDropdownBlankBlock`: сортирует варианты в шаблоне `{Правильный; Дистракторы}` по алфавиту, исключая выдачу первого элемента как правильного ответа, удаляет `correctAnswer` и `correctIndex`.
5. Для блоков `QuizInputBlankBlock`: заменяет скрытый правильный ответ на нейтральный заполнитель `{blank_N}`, удаляет поле `correctAnswer`.
6. Для блоков `QuizSequenceBlock`: инвертирует/перемешивает начальный порядок шагов.
7. Рекурсивно очищает все чувствительные ключи ответов.
8. Возвращает санитизированную строку JSON.

---

### ⚡ Функция: `ValidateQuizSubmission(contentJSON, answers)`

* **Файл и строки:** [`sanitizer.go#L240-L365`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/quiz/sanitizer.go#L240-L365)
* **Бизнес-назначение:** Серверный расчет набранных баллов и проверка правильности ответов на основе эталонного контента урока без доверия к браузерному счетчику.
* **Связанная фича:** *Серверная проверка квизов*

#### 📥 Входные параметры
| Параметр | Тип | Обязателен | Бизнес-смысл и ограничения |
|---|---|:---:|---|
| `contentJSON` | `string` | Да | Эталонный (несанитизированный) JSON урока из БД |
| `answers` | `[]LessonAnswerSubmission` | Да | Выбранные студентом ответы по ID блоков |

#### 🔄 Пошаговый алгоритм работы
1. Сопоставляет переданные ответы студента с эталонными блоками урока.
2. Для каждого проверочного блока сверяет выбранный ответ с правильным (`QuizSingleBlock`, `QuizMultiBlock`, `QuizMatchBlock`, `QuizDropdownBlankBlock`, `QuizInputBlankBlock`, `QuizSequenceBlock`).
3. Начисляет баллы пропорционально правильности решения.
4. Вычисляет итоговый процент:
   $$\text{Score} = \left(\frac{\text{EarnedPoints}}{\text{TotalMaxPoints}}\right) \times 100$$
5. Формирует структуру `LessonCompletionResult` с подтвержденным баллом и результатами по каждому блоку для подсветки в плеере.

---

### ⚡ Функция: `GetStudentHomeworkFeedback(ctx, userID, lessonID)`

* **Файл и строки:** [`homework_feedback.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/quiz/homework_feedback.go)
* **Бизнес-назначение:** Предоставление студенту детализированных результатов проверки его домашнего задания (эссе, открытых ответов) преподавателем.
* **Связанная фича:** *Рецензия и результаты проверки ДЗ*

#### 📥 Входные параметры
| Параметр | Тип | Обязателен | Бизнес-смысл и ограничения |
|---|---|:---:|---|
| `userID` | `int64` | Да | Идентификатор авторизованного студента (из JWT) |
| `lessonID` | `int64` | Да | Идентификатор урока с домашним заданием |

#### 🔄 Пошаговый алгоритм работы
1. Валидация входных идентификаторов (`userID > 0`, `lessonID > 0`).
2. Проверка существования урока и связанного курса (`lessonRepo.GetLessonByID`, `courseRepo.GetCourseByID`).
3. Проверка прав доступа студента к курсу через `accessService.CanViewCourse` (403 при отсутствии доступа).
4. Поиск последней сдачи студента в `quizRepo.GetStudentHomeworkFeedback`:
   - Если сдачи нет ➔ `has_submission: false`, `status: "not_submitted"`.
   - Если есть непроверенные ответы (`is_correct IS NULL`) ➔ `status: "pending"`, дата проверки и данные преподавателя скрыты.
   - Если все ответы проверены ➔ `status: "graded"`, возвращаются баллы, максимальные баллы, рецензия (`feedback`), дата проверки (`graded_at`) и карточка преподавателя (`teacher: { id, name, avatar_url }`).

---

## 💡 Подсказка для аналитика (Где менять логику?)
* *Сетевой допуск на задержку сдачи теста (сейчас +30 секунд):* [`submitAttempt.go#L77`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/quiz/submitAttempt.go#L77).
* *Правило выставления флага `IsCorrect` при ручной проверке (сейчас `points > 0`):* [`gradeAttemptAnswer.go#L35`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/quiz/gradeAttemptAnswer.go#L35).
* *Формула расчета процента сдачи теста:* метод `CalculateScore` сущности `domain.QuizAttempt`.
* *Определение статуса проверки ДЗ (`not_submitted` / `pending` / `graded`):* [`homework_feedback.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/quiz/homework_feedback.go).


