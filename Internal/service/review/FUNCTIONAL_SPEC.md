# 📋 Функциональная спецификация: `Отзывы и рейтинги курсов (Course Reviews)`

> **Расположение:** `internal/service/review`  
> **Технический контекст:** [`README.md`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/review/README.md)  
> **Дата актуализации:** 2026-10-04  

---

## 🎯 Каталог бизнес-фич модуля
| Фича | Описание возможности | Обеспечивающие функции |
|---|---|---|
| **Оценка и отзыв о курсе** | Добавление или обновление студентом отзыва (1–5 звезд, текстовый комментарий) с проверкой $\ge 30\%$ прогресса | `AddOrUpdateReview()` |
| **Удаление отзыва** | Удаление студентом своего отзыва с автоматическим пересчетом общего рейтинга курса | `DeleteReview()` |
| **Публичный просмотр отзывов курса** | Пагинированный просмотр всех отзывов курса со сводной статистикой (средний балл, распределение 1–5 звезд) | `ListCourseReviews()` |
| **Просмотр своего отзыва** | Получение авторизованным пользователем ранее оставленного им отзыва на курс | `GetMyReview()` |

---

## 🔬 Паспорта функций

### ⚡ Функция: `AddOrUpdateReview(ctx, userID, courseID, rating, comment)`

* **Файл:** [`service.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/review/service.go)
* **Бизнес-назначение:** Добавление или обновление отзыва учащегося.
* **Связанная фича:** *Оценка и отзыв о курсе*

#### 📥 Входные параметры
| Параметр | Тип | Обязателен | Бизнес-смысл и ограничения |
|---|---|:---:|---|
| `userID` | `int64` | Да | ID авторизованного пользователя |
| `courseID` | `int64` | Да | ID курса |
| `rating` | `int` | Да | Оценка от 1 до 5 звезд включительно |
| `comment` | `*string` | Нет | Текстовый комментарий (может быть nil) |

#### 🔄 Пошаговый алгоритм работы
1. **Шаг 1 (Валидация оценки):** Проверяет, что `rating >= 1 && rating <= 5`. Если нет — `ErrValidationFailed`.
2. **Шаг 2 (Проверка зачисления):** Запрашивает `enrolledRepo.UserExistCourse(ctx, db, userID, courseID)`. Если пользователь не записан на курс — возвращает `ErrForbidden`.
3. **Шаг 3 (Проверка прогресса $\ge 30\%$):** Запрашивает `progressRepo.GetCourseProgress(ctx, db, userID, courseID)`. Если `progress.Percent < 30` — возвращает `ErrForbidden` (недостаточный прогресс для объективной оценки).
4. **Шаг 4 (Транзакционный апсерт и пересчет):**
   - Внутри транзакции сохраняет/обновляет отзыв через `reviewRepo.UpsertReview`.
   - Запрашивает агрегированные данные рейтинга курса через `reviewRepo.GetCourseRatingSummary`.
   - Обновляет агрегированные значения `rating` и `reviews_count` курса в таблице `courses` через `courseRepo.UpdateCourseRatingStats`.
5. **Шаг 5 (Возврат результата):** Возвращает сохраненный доменный объект отзыва.

---

### ⚡ Функция: `DeleteReview(ctx, userID, courseID)`

* **Файл:** [`service.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/review/service.go)
* **Бизнес-назначение:** Удаление своего отзыва.
* **Связанная фича:** *Удаление отзыва*

#### 🔄 Пошаговый алгоритм работы
1. **Шаг 1 (Удаление):** Внутри транзакции удаляет отзыв через `reviewRepo.DeleteReview(ctx, tx, courseID, userID)`. Если отзыв не найден — возвращает `ErrNotFound`.
2. **Шаг 2 (Пересчет статистики):** Пересчитывает агрегаты через `reviewRepo.GetCourseRatingSummary(ctx, tx, courseID)`.
3. **Шаг 3 (Обновление курса):** Сохраняет новые `rating` и `reviews_count` через `courseRepo.UpdateCourseRatingStats(ctx, tx, courseID, summary.AverageRating, summary.ReviewsCount)`.
