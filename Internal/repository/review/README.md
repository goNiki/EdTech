# 📦 Репозиторий: `internal/repository/review`
> **Путь:** `internal/repository/review`  
> **Роль:** Работа с таблицей `course_reviews` в PostgreSQL.

---

## 🎯 Назначение и ответственность
Репозиторий обеспечивает персистентность данных отзывов студентов: сохранение/обновление записей (Upsert по уникальному ключу `course_id, user_id`), удаление, чтение отдельного отзыва, выборку пагинированного списка отзывов и вычисление сводной статистики (средний балл и распределение оценок).

---

## 📁 Структура файлов
| Файл | Описание |
|---|---|
| [`repository.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/review/repository.go) | Реализация методов интерфейса `repository.ReviewRepository` |

---

## ⚙️ Функции, методы и API
| Метод | Описание |
|---|---|
| `UpsertReview` | Сохранение или обновление отзыва с `ON CONFLICT (course_id, user_id)` |
| `DeleteReview` | Удаление отзыва пользователя по `course_id` и `user_id` |
| `GetReviewByUserAndCourse` | Получение отзыва студента на указанный курс |
| `ListReviewsByCourse` | Пагинированный список отзывов с сортировкой по дате создания (новые первыми) |
| `GetCourseRatingSummary` | Агрегация `average_rating`, `reviews_count` и распределения звезд (1–5) |
