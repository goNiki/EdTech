# 📋 Функциональная спецификация: `Управление уроками (Lessons Management)`

> **Расположение:** `internal/service/lesson`  
> **Технический контекст:** [`README.md`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/lesson/README.md)  
> **Дата актуализации:** 2026-10-04  

---

## 🎯 Каталог бизнес-фич модуля
| Фича | Описание возможности | Обеспечивающие функции |
|---|---|---|
| **Создание урока** | Добавление урока (видео, статья или квиз) с валидацией курса и монотонным инкрементом позиции | `CreateLesson()` |
| **Просмотр контента урока** | Загрузка контента урока (текст, видео-ссылка, длительность) для плеера студента | `GetLesson()` |
| **Редактирование урока** | Обновление обучающих материалов, описания и ссылок | `UpdateLesson()` |
| **Смена статуса урока** | Публикация или снятие урока с показа (`draft`, `published`, `archived`) | `UpdateLessonStatus()` |
| **Удаление урока** | Удаление урока с каскадным удалением вложений | `DeleteLesson()` |

---

## 🔬 Паспорта функций

### ⚡ Функция: `CreateLesson(ctx, lesson)`

* **Файл и строки:** [`createLesson.go#L18-L46`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/lesson/createLesson.go#L18-L46)
* **Бизнес-назначение:** Создание нового учебного шага в программе курса.
* **Связанная фича:** *Создание урока*

#### 📥 Входные параметры
| Параметр | Тип | Обязателен | Бизнес-смысл и ограничения |
|---|---|:---:|---|
| `lesson.CourseID` | `int64` | Да | Идентификатор существующего курса |
| `lesson.Title` | `string` | Да | Заголовок урока (не пустой) |
| `lesson.Description`| `string` | Да | Краткое описание урока |
| `lesson.Type` | `string` | Да | Тип контента: `video`, `text` или `quiz` |
| `lesson.VideoURL` | `*string` | Нет | Ссылка на потоковое видео (для video-уроков) |

#### 🔄 Пошаговый алгоритм работы
1. **Шаг 1 (Валидация полей):** Вызывает `utils.ValidateLesson(courseID, title, description)`.
2. **Шаг 2 (Проверка курса):** Проверяет существование курса в `courserepo.GetCourseByID`. Если курс не найден ➔ `ErrNotFoundCourse`.
3. **Шаг 3 (Расчет позиции):** Вычисляет `position = GetMaxPositionByCourseID + 1`.
4. **Шаг 4 (Сохранение):** Выполняет `INSERT INTO lessons` и возвращает присвоенный ID.

---

### ⚡ Функция: `UpdateLessonStatus(ctx, lessonID, status)`

* **Файл и строки:** [`updateLessonStatus.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/lesson/updateLessonStatus.go)
* **Бизнес-назначение:** Публикация отдельного урока преподавателем.
* **Связанная фича:** *Смена статуса урока*

#### 💡 Подсказка для аналитика (Где менять логику?)
* *Добавить валидацию формата видео-ссылки (YouTube, Vimeo, S3):* добавить проверку регулярного выражения в [`createLesson.go#L24`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/lesson/createLesson.go#L24).
