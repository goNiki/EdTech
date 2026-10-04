# 📦 Сервис: `internal/service/analytics`
> **Путь:** `internal/service/analytics`  
> **Роль:** Расчет и сбор образовательной аналитики, детальной успеваемости студентов и очереди работ на ручную проверку.

---

## 🎯 Назначение и ответственность
Сервис предоставляет аналитические данные для авторов курсов и преподавателей. Собирает метрики прохождения курса (общий процент завершения, средний балл за квизы, количество активных и отчисленных студентов), формирует детальный срез по конкретному студенту (`drilldown`) и выдает список неручно проверенных домашних заданий/эссе.

---

## ⚠️ Жестко установленные правила (Hard Invariants)
1. **Обязательная авторизация преподавателя (`checkTeacherAccess`):** Любой вызов аналитики курса или успеваемости студента валидирует право `accessService.CanEditCourse`. Студенты и посторонние пользователи получают `403 Forbidden`.
2. **Изоляция данных курса:** Аналитика и очередь работ строго фильтруются по `courseID`. Преподаватель видит работы только в рамках курсов, к которым он прикреплен.
3. **Агрегация только по завершенным урокам/попыткам:** Расчет средней успеваемости и баллов учитывает только сданные и проверенные попытки (`CompletedAt != nil`).

---

## 📁 Структура файлов
| Файл | Описание |
|---|---|
| [`service.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/analytics/service.go) | Определение сервиса и внутренняя проверка доступа `checkTeacherAccess` |
| [`getCourseAnalytics.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/analytics/getCourseAnalytics.go) | Получение сводной аналитики по курсу (`CourseAnalyticsSummary`) |
| [`getStudentDrilldown.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/analytics/getStudentDrilldown.go) | Детальный отчет по урокам и тестам конкретного студента |
| [`listPendingHomeworks.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/analytics/listPendingHomeworks.go) | Пагинированный список заданий конкретного курса, ожидающих проверки |
| [`list_teacher_pending.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/analytics/list_teacher_pending.go) | Глобальная очередь непроверенных работ преподавателя по всем его курсам с агрегированной сводкой |

---

## ⚙️ Функции, методы и API
| Функция / Метод | Файл:Строки | Описание | Сигнатура |
|---|---|---|---|
| `GetCourseAnalytics` | [`getCourseAnalytics.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/analytics/getCourseAnalytics.go) | Сводный отчет по курсу (% завершения, средний балл) | `(s *analyticsService) GetCourseAnalytics(ctx, teacherID, courseID int64) (domain.CourseAnalyticsSummary, error)` |
| `GetStudentDrilldown` | [`getStudentDrilldown.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/analytics/getStudentDrilldown.go) | Детализированный отчет по студенту на курсе | `(s *analyticsService) GetStudentDrilldown(ctx, teacherID, courseID, studentID int64) (*domain.StudentDrilldownReport, error)` |
| `ListPendingHomeworks` | [`listPendingHomeworks.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/analytics/listPendingHomeworks.go) | Список работ курса со статусом `NeedsGrading` для ручной оценки | `(s *analyticsService) ListPendingHomeworks(ctx, teacherID, courseID, page, pageSize int64) ([]domain.PendingHomeworkItem, int64, error)` |
| `ListTeacherPendingHomeworks` | [`list_teacher_pending.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/analytics/list_teacher_pending.go) | Кросс-курсовая очередь домашних заданий преподавателя со сводкой | `(s *analyticsService) ListTeacherPendingHomeworks(ctx, teacherID, courseID, page, pageSize int64) (*domain.TeacherPendingHomeworksResult, error)` |

---

## 🔗 Зависимости
- **Входящие:** [`handlers/analytics`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/analytics).
- **Исходящие:** [`repository.AnalyticsRepository`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/analytics), [`repository.CourseRepository`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/course), [`service.AccessService`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/access).

---

## 🧭 Навигатор типовых задач (Where to edit?)
- **Добавить новую метрику в дашборд курса (например, среднее время прохождения):** расширить структуру `domain.CourseAnalyticsSummary` и SQL-запрос в [`repository/analytics/summary.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/analytics/summary.go).
- **Изменить критерии попадания работ в очередь проверки:** [`listPendingHomeworks.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/analytics/listPendingHomeworks.go) и [`repository/analytics/pending.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/analytics/pending.go).
