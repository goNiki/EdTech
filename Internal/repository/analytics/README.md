# 📦 Репозиторий: `internal/repository/analytics`
> **Путь:** `internal/repository/analytics`  
> **Роль:** Аналитические SQL-агрегации над курсами, успеваемостью студентов и неручно проверенными заданиями.

---

## 🎯 Назначение и ответственность
Выполняет сложные агрегационные SQL-запросы на базе Common Table Expressions (CTE) для формирования сводок курсов, детальных срезов прохождения уроков отдельными студентами и очередей домашних работ на проверку.

---

## ⚠️ Жестко установленные правила (Hard Invariants)
1. **Единый SQL-раундтрип в `GetCourseAnalyticsSummary`:** Сводка курса (количество студентов, средний % прохождения, средний балл и счетчик долгов по проверке) вычисляется одним запросом через 4 CTE-блока (`student_count`, `course_total_lessons`, `progress_stats`, `pending_hw`), исключая N+1 нагрузку.
2. **Нулевая безопасность (COALESCE):** Все средние значения и агрегаты обернуты в `COALESCE(..., 0.0)` для предотвращения ошибок при отсутствии студентов или уроков на курсе.
3. **Критерий незавершенной проверки:** Непроверенные работы определяются строгим условием `qaa.is_correct IS NULL` в таблице `quiz_attempt_answers`.

---

## 📁 Структура файлов
| Файл | Описание |
|---|---|
| [`repository.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/analytics/repository.go) | Фабрика `NewAnalyticsRepository` |
| [`summary.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/analytics/summary.go) | Однопроходная CTE-агрегация метрик курса (`CourseAnalyticsSummary`) |
| [`drilldown.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/analytics/drilldown.go) | Срез успеваемости студента по каждому уроку курса |
| [`pending.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/analytics/pending.go) | Выборка работ курса с `is_correct IS NULL` с пагинацией |
| [`teacher_pending.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/analytics/teacher_pending.go) | Выборка работ по всем курсам преподавателя со сводкой `courses_summary` |

---

## ⚙️ Функции, методы и API
| Функция / Метод | Файл:Строки | Описание | Сигнатура |
|---|---|---|---|
| `GetCourseAnalyticsSummary` | [`summary.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/analytics/summary.go) | Комплексная сводка по курсу | `GetCourseAnalyticsSummary(ctx context.Context, courseID int64) (domain.CourseAnalyticsSummary, error)` |
| `GetStudentDrilldown` | [`drilldown.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/analytics/drilldown.go) | Детальный список оценок и статусов студента | `GetStudentDrilldown(ctx context.Context, userID, courseID int64) (*domain.StudentDrilldownReport, error)` |
| `ListPendingHomeworks` | [`pending.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/analytics/pending.go) | Список работ курса на ручную проверку | `ListPendingHomeworks(ctx context.Context, courseID, limit, offset int64) ([]domain.PendingHomeworkItem, int64, error)` |
| `ListTeacherPendingHomeworks` | [`teacher_pending.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/analytics/teacher_pending.go) | Кросс-курсовая очередь проверки со сводкой | `ListTeacherPendingHomeworks(ctx context.Context, teacherID, courseID, limit, offset int64) ([]domain.PendingHomeworkItem, int64, []domain.CoursePendingSummaryItem, error)` |

---

## 🔗 Зависимости
- **Входящие:** [`service/analytics`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/analytics).
- **Исходящие:** [`db.QueryExecutor`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/infrastructure/db).
