# 📦 Сервис: `internal/service/lesson`
> **Путь:** `internal/service/lesson`  
> **Роль:** Управление уроками курса, типами контента (видео, текст, квиз), авто-позиционированием и статусами.

---

## 🎯 Назначение и ответственность
Обеспечивает CRUD-операции над уроками курса и секций. Вычисляет порядковый номер урока в курсе, валидирует обязательные поля и управляет видимостью урока (`draft`, `published`, `archived`).

---

## ⚠️ Жестко установленные правила (Hard Invariants)
1. **Проверка существования родительского курса:** Урок не может быть создан без привязки к существующему курсу (`ErrNotFoundCourse`).
2. **Автоматическое вычисление позиции (`Position`):** При создании урока позиция рассчитывается как `max(position) + 1` по курсу (`GetMaxPositionByCourseID`), гарантируя монотонный порядок.
3. **Обязательная валидация:** `ValidateLesson` строго требует наличия непустого `title`, `description` и валидного `courseID`.
4. **Изоляция уроков:** Удаление урока проверяет существование и каскадно очищает связанные материалы в БД.
5. **Валидация структуры контента (Puck JSON):** При сохранении контента (`Content`) валидируется корректность JSON, обязательное наличие `type`, уникальность `props.id` (дубликаты блоков исключены) и ограничение 5 МБ на размер тела запроса.

---

## 📁 Структура файлов
| Файл | Описание |
|---|---|
| [`service.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/lesson/service.go) | Определение сервиса и конструктор `NewLessonService` |
| [`createLesson.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/lesson/createLesson.go) | Создание урока с автовычислением `position + 1` |
| [`getLesson.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/lesson/getLesson.go) | Получение данных урока по ID |
| [`updateLesson.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/lesson/updateLesson.go) | Обновление названия, описания, видео-ссылки или типа урока |
| [`updateLessonStatus.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/lesson/updateLessonStatus.go) | Смена статуса видимости урока (`draft`/`published`) |
| [`deleteLesson.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/lesson/deleteLesson.go) | Удаление урока по ID |
| [`listLessonsByCourseID.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/lesson/listLessonsByCourseID.go) | Получение всех уроков курса |
| [`getLessonNavigation.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/lesson/getLessonNavigation.go) | Контекст навигации урока (следующий/предыдущий) и дерево силлабуса курса с прогрессом |

---

## ⚙️ Функции, методы и API
| Функция / Метод | Файл:Строки | Описание | Сигнатура |
|---|---|---|---|
| `CreateLesson` | [`createLesson.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/lesson/createLesson.go) | Проверка курса, расчет позиции и создание записи | `CreateLesson(ctx context.Context, lesson *domain.Lesson) (int64, error)` |
| `GetLesson` | [`getLesson.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/lesson/getLesson.go) | Чтение урока по ID | `GetLesson(ctx context.Context, lessonID int64) (*domain.Lesson, error)` |
| `UpdateLesson` | [`updateLesson.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/lesson/updateLesson.go) | Редактирование параметров урока | `UpdateLesson(ctx context.Context, lesson *domain.Lesson) error` |
| `UpdateLessonStatus` | [`updateLessonStatus.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/lesson/updateLessonStatus.go) | Переключение статуса урока | `UpdateLessonStatus(ctx context.Context, lessonID int64, status string) error` |
| `DeleteLesson` | [`deleteLesson.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/lesson/deleteLesson.go) | Удаление урока из системы | `DeleteLesson(ctx context.Context, lessonID int64) error` |
| `GetLessonNavigationContext` | [`getLessonNavigation.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/lesson/getLessonNavigation.go) | Сборка контекста навигации (prev/next) и силлабуса с успеваемостью | `GetLessonNavigationContext(ctx context.Context, userID int64, lessonID int64) (*domain.LessonNavigationContext, error)` |

---

## 🔗 Зависимости
- **Входящие:** [`handlers/lesson`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/lesson), [`handlers/courses`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/courses).
- **Исходящие:** [`repository.LessonRepository`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/lesson), [`repository.CourseRepository`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/course), [`pkg/utils`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/pkg/utils/validate.go).

---

## 🧭 Навигатор типовых задач (Where to edit?)
- **Добавить новый тип урока (например, `webinar` или `code_challenge`):** обновить доменные константы в `domain/lesson.go` и валидацию в [`createLesson.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/lesson/createLesson.go).
- **Привязать урок к секции обязательно:** добавить проверку `lesson.SectionID != nil` в [`createLesson.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/lesson/createLesson.go).
