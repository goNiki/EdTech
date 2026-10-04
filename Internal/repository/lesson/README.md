# 📦 Репозиторий: `internal/repository/lesson`
> **Путь:** `internal/repository/lesson`  
> **Роль:** SQL-запросы к таблице `lessons` (уроки, видеоматериалы, позиции, статусы).

---

## 🎯 Назначение и ответственность
Отвечает за сохранение уроков в PostgreSQL, выборку по курсу с сортировкой по позиции, пакетное перемещение уроков между секциями и обновление статусов.

---

## ⚠️ Жестко установленные правила (Hard Invariants)
1. **Перемещение между секциями (`ReorderLessons`):** Метод поддерживает пакетный перенос уроков между модулями: если передан `sectionID`, обновляются одновременно и позиция `idx + 1`, и ссылка на секцию `section_id` за 1 SQL-запрос через PostgreSQL `unnest()`.
2. **Упорядочивание уроков:** Запросы списка уроков курса всегда отсортированы по `ORDER BY position ASC`.
3. **Безопасная обработка NULL в секции:** Урок может существовать как внутри секции, так и без нее (`section_id NULL`).

---

## 📁 Структура файлов
| Файл | Описание |
|---|---|
| [`repository.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/lesson/repository.go) | Фабрика `NewLessonRepo` |
| [`createLesson.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/lesson/createLesson.go) | `INSERT INTO lessons ... RETURNING` |
| [`getLessonByID.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/lesson/getLessonByID.go) | Чтение урока по первичному ключу |
| [`getLessonsByCourseID.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/lesson/getLessonsByCourseID.go) | Получение всех уроков курса |
| [`getMaxPositionByCourseID.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/lesson/getMaxPositionByCourseID.go) | Запрос `COALESCE(MAX(position), 0)` по курсу |
| [`reorderLessons.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/lesson/reorderLessons.go) | Пакетное обновление позиций уроков через `unnest()` |
| [`reorder_test.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/lesson/reorder_test.go) | Unit-тесты пакетного обновления порядка уроков |
| [`updateLesson.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/lesson/updateLesson.go) | Обновление текстовых полей, типа и ссылки на видео |
| [`updateStatus.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/lesson/updateStatus.go) | Обновление статусов уроков |
| [`deleteLessonByID.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/lesson/deleteLessonByID.go) | Удаление урока |

---

## ⚙️ Функции, методы и API
| Функция / Метод | Файл:Строки | Описание | Сигнатура |
|---|---|---|---|
| `CreateLesson` | [`createLesson.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/lesson/createLesson.go) | Вставка урока в БД | `CreateLesson(ctx context.Context, lesson *domain.Lesson) error` |
| `GetLessonsByCourseID` | [`getLessonsByCourseID.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/lesson/getLessonsByCourseID.go) | Выборка всех уроков курса (`ORDER BY position ASC`) | `GetLessonsByCourseID(ctx context.Context, courseID int64) ([]domain.Lesson, error)` |
| `ReorderLessons` | [`reorderLessons.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/lesson/reorderLessons.go) | Пакетное сохранение нового порядка уроков через `unnest()` | `ReorderLessons(ctx context.Context, sectionID *int64, lessonIDs []int64) error` |

---

## 🔗 Зависимости
- **Входящие:** [`service/lesson`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/lesson), [`service/course`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/course), [`service/progress`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/progress).
- **Исходящие:** [`txmanager`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/infrastructure/txmanager), [`models/converter`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/models/converter).
