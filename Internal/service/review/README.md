# 📦 Сервис: `internal/service/review`
> **Путь:** `internal/service/review`  
> **Роль:** Бизнес-логика отзывов и оценок курсов, валидация прав учащихся и агрегация рейтинга.

---

## 🎯 Назначение и ответственность
Сервис управляет добавлением, обновлением, удалением и просмотром отзывов и оценок (1–5 звезд) для курсов. Он защищает рейтинг курса от накруток, проверяя факт зачисления студента и достаточный прогресс прохождения курса ($\ge 30\%$), а также автоматически пересчитывает средний рейтинг и количество отзывов в таблице `courses`.

---

## ⚠️ Жестко установленные правила (Hard Invariants)
1. **Диапазон оценки:** Оценка (`rating`) должна строго находиться в диапазоне от `1` до `5` включительно (`ErrValidationFailed`).
2. **Проверка зачисления:** Оставлять отзыв может только студент, зачисленный на данный курс (`enrolledRepo.UserExistCourse`). Иначе возвращается `ErrForbidden`.
3. **Порог прохождения курса:** Оставлять отзыв разрешено только при прогрессе по курсу не менее `30%` (`progressRepo.GetCourseProgress` с `Percent >= 30`). Иначе возвращается `ErrForbidden`.
4. **Один отзыв на студента:** У одного пользователя может быть только один отзыв на конкретный курс (Upsert-семантика: повторный вызов обновляет оценку и текст).
5. **Транзакционный пересчет статистики:** При добавлении, редактировании или удалении отзыва статистика курса (`rating` и `reviews_count`) пересчитывается и обновляется в базе данных внутри транзакции.

---

## 📁 Структура файлов
| Файл | Описание |
|---|---|
| [`service.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/review/service.go) | Реализация методов сервиса отзывов (`AddOrUpdateReview`, `DeleteReview`, `ListCourseReviews`, `GetMyReview`) |
| [`service_test.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/review/service_test.go) | Модульные тесты на проверку ограничений (зачисление, прогресс, диапазон оценок) |

---

## ⚙️ Функции, методы и API
| Функция / Метод | Файл:Строки | Описание | Сигнатура |
|---|---|---|---|
| `AddOrUpdateReview` | [`service.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/review/service.go) | Создание или обновление отзыва с проверкой зачисления и прогресса $\ge 30\%$ | `AddOrUpdateReview(ctx context.Context, userID, courseID int64, rating int, comment *string) (*domain.Review, error)` |
| `DeleteReview` | [`service.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/review/service.go) | Удаление отзыва пользователя с пересчетом статистики курса | `DeleteReview(ctx context.Context, userID, courseID int64) error` |
| `ListCourseReviews` | [`service.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/review/service.go) | Пагинированный список отзывов к курсу со сводкой рейтинга | `ListCourseReviews(ctx context.Context, courseID int64, page, pageSize int) ([]*domain.Review, *domain.CourseReviewsSummary, int64, error)` |
| `GetMyReview` | [`service.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/review/service.go) | Получение собственного отзыва текущего пользователя | `GetMyReview(ctx context.Context, userID, courseID int64) (*domain.Review, error)` |

---

## 🔗 Зависимости
- **Входящие:** [`handlers/review`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/review).
- **Исходящие:**
  - [`repository.ReviewRepository`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/review)
  - [`repository.CourseRepository`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/course)
  - [`repository.EnrolledRepository`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/enrollment)
  - [`repository.ProgressRepository`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/progress)
  - [`infrastructure.txmanager.TransactionManager`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/infrastructure/txmanager)
