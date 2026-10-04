# 📦 Модуль: `internal/domain`
> **Путь:** `internal/domain`  
> **Роль:** Чистые доменные сущности (Entities), структуры данных предметной области и инварианты бизнес-логики.

---

## 🎯 Назначение и ответственность
Содержит ядро предметной области EdTech платформы без внешних зависимостей от БД, HTTP или фреймворков. Определяет структуры сущностей, их статусы, правила валидации, бизнес-методы (например, расчет баллов квиза, проверка прав публикации, архивация, логин).

---

## 📁 Структура файлов модуля
| Файл | Описание роли файла |
|---|---|
| [`auth.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/domain/auth.go) | Сущности `User`, `Role` (student/teacher/admin), токены авторизации `AuthTokens`, `RefreshTokenData` |
| [`course.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/domain/course.go) | Сущность `Course`, статусы (draft, published, archived), видимость, генерация slug, методы `Publish`, `Archive` |
| [`section.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/domain/section.go) | Сущность `Section` (модуль курса), статусы видимости, позиции в курсе |
| [`lesson.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/domain/lesson.go) | Сущность `Lesson`, правила тестирования `QuizSettings` (таймеры, лимиты, режимы обратной связи, проходной порог), валидация контента Puck |
| [`quiz.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/domain/quiz.go) | Сущности `Quiz`, `QuizQuestion`, `QuizAnswer`, `QuizAttempt`, типы вопросов (single, multi, open_text), расчет скоринга |
| [`progress.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/domain/progress.go) | Сущности `LessonProgress`, `CourseProgress`, результаты тестирования `LessonCompletionResult`, сводка попыток `LessonAttemptsSummary` |
| [`enroll.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/domain/enroll.go) | Запись зачисления студента/преподавателя на курс (`UserCourse`, `Enrollment`) |
| [`analytics.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/domain/analytics.go) | Аналитические срезы `CourseAnalytics`, `StudentDrilldown`, очереди на ручную проверку `PendingHomework` |
| [`permission.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/domain/permission.go) | Системные права доступа `Permission`, `RolePermission` |
| [`category.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/domain/category.go) | Категории курсов `Category` |
| [`resource.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/domain/resource.go) | Вложения и учебные материалы уроков `Resource` |
| [`pagination.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/domain/pagination.go) | Базовая структура пагинации списка сущностей |

---

## ⚙️ Функции, методы и API
| Функция / Класс | Файл:Строки | Описание | Сигнатура / Вход и Выход |
|---|---|---|---|
| [`User.CanLogin`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/domain/auth.go#L63-L72) | [`auth.go#L63-L72`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/domain/auth.go#L63-L72) | Проверка, не забанен ли пользователь и активен ли аккаунт | `(u *User) CanLogin() error` |
| [`User.UpdateProfile`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/domain/auth.go#L111-L124) | [`auth.go#L111-L124`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/domain/auth.go#L111-L124) | Мутация полей профиля пользователя | `(u *User) UpdateProfile(input UpdateProfileInput)` |
| [`Lesson.GetQuizSettings`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/domain/lesson.go#L39-L59) | [`lesson.go#L39-L59`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/domain/lesson.go#L39-L59) | Получение эффективных настроек тестирования урока с дефолтными значениями | `(l *Lesson) GetQuizSettings() QuizSettings` |
| [`Course.Publish`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/domain/course.go#L73-L81) | [`course.go#L73-L81`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/domain/course.go#L73-L81) | Перевод курса в статус published с фиксацией времени | `(c *Course) Publish(time time.Time)` |
| [`Course.Archive`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/domain/course.go#L60-L64) | [`course.go#L60-L64`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/domain/course.go#L60-L64) | Перевод курса в статус archived | `(c *Course) Archive(time time.Time)` |
| [`Course.CanSelfEnroll`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/domain/course.go#L83-L88) | [`course.go#L83-L88`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/domain/course.go#L83-L88) | Проверка доступности курса для самостоятельной записи | `(c *Course) CanSelfEnroll() error` |
| [`Quiz.Validate`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/domain/quiz.go#L29-L46) | [`quiz.go#L29-L46`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/domain/quiz.go#L29-L46) | Валидация проходного балла, лимитов попыток и времени теста | `(q *Quiz) Validate() error` |
| [`QuizAttempt.CalculateScore`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/domain/quiz.go#L76-L84) | [`quiz.go#L76-L84`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/domain/quiz.go#L76-L84) | Автоматический подсчет процента и выставление флага сдачи `Passed` | `(a *QuizAttempt) CalculateScore(correct, total, pass int)` |
| [`LessonProgress.IsCompleted`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/domain/progress.go) | [`progress.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/domain/progress.go) | Проверка завершенности урока студентом | `(lp *LessonProgress) IsCompleted() bool` |

---

## 🔗 Зависимости
- **Входящие (кто вызывает этот модуль):**
  - [`internal/service/*`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/README.md) (бизнес-сценарии)
  - [`internal/repository/*`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/README.md) (маппинг из базы в домен)
  - [`internal/interfaces/handlers/*`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/README.md) (через DTO и конвертеры)
- **Исходящие (что импортирует этот модуль):**
  - [`pkg/errors`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/pkg/README.md)
  - Стандартная библиотека Go (`time`, `regexp`).

---

## 🧭 Навигатор типовых задач (Where to edit?)
- **Добавить новое поле в сущность пользователя (например, телефон или город):** отредактировать структуры `User`, `UpdateProfileInput` в [`auth.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/domain/auth.go).
- **Добавить новый статус или тип курса/урока:** обновить константы в [`course.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/domain/course.go) или [`lesson.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/domain/lesson.go).
- **Изменить формулу скоринга квизов:** модифицировать метод `CalculateScore` в [`quiz.go#L76-L84`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/domain/quiz.go#L76-L84).
- **Добавить новые правила бизнес-валидации перед переходом статуса:** добавлять методы `Can*` непосредственно на структурах в `internal/domain`.
