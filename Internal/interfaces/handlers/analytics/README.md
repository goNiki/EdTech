# 📦 Обработчики: `internal/interfaces/handlers/analytics`
> **Путь:** `internal/interfaces/handlers/analytics`  
> **Роль:** HTTP-контроллеры дашборда аналитики преподавателя, успеваемости студентов и очереди проверки.

---

## 🎯 Назначение и ответственность
Принимает запросы по аналитике курса на маршрутах `/api/v1/courses/{courseid}/*`.

---

## 📁 Структура файлов и Эндпоинты
| Файл | HTTP Метод и Путь | Описание |
|---|---|---|
| [`get_analytics.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/analytics/get_analytics.go) | `GET /api/v1/courses/{courseid}/analytics` | Сводная аналитика курса (студенты, % завершения) |
| [`get_student_drilldown.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/analytics/get_student_drilldown.go) | `GET /api/v1/courses/{courseid}/students/{userid}/drilldown` | Детализированная успеваемость конкретного студента |
| [`list_pending_hw.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/analytics/list_pending_hw.go) | `GET /api/v1/courses/{courseid}/grading/pending` | Список работ курса, ожидающих ручной проверки |
| [`list_teacher_pending.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/analytics/list_teacher_pending.go) | `GET /api/v1/teacher/grading/pending` | Глобальная очередь проверки преподавателя со сводкой по всем его курсам |

---

## 🔗 Зависимости
- **Входящие:** [`internal/app/di.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/app/di.go).
- **Исходящие:** [`service.AnalyticsServices`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/analytics), [`internal/dto`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/dto).
