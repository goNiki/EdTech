# 📦 Обработчики: `internal/interfaces/handlers/progress`
> **Путь:** `internal/interfaces/handlers/progress`  
> **Роль:** HTTP-контроллеры процесса обучения: старт урока, фиксация времени просмотра, завершение урока, Pre-flight сводка и старт попыток тестирования, чтение прогресса курса.

---

## 🎯 Назначение и ответственность
Принимает запросы прогресса студента как на маршрутах уроков (`/api/v1/lessons/{lesson_id}/*`), так и курсов (`/api/v1/courses/{course_id}/progress*`).

---

## 📁 Структура файлов и Эндпоинты
| Файл | HTTP Метод и Путь | Описание |
|---|---|---|
| [`startLesson.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/progress/startLesson.go) | `POST /api/v1/lessons/{lesson_id}/start` | Фиксация начала прохождения урока |
| [`updateLessonProgress.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/progress/updateLessonProgress.go) | `PATCH /api/v1/lessons/{lesson_id}/progress` | Обновление времени просмотра и позиции |
| [`completeLesson.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/progress/completeLesson.go) | `POST /api/v1/lessons/{lesson_id}/complete` | Завершение урока, серверная верификация квизов (anti-cheat), досрочный выход и эссе |
| [`getAttemptsSummary.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/progress/getAttemptsSummary.go) | `GET /api/v1/lessons/{lesson_id}/attempts/summary` | Pre-flight экран: количество попыток, лимиты, лучший балл и история |
| [`startAttempt.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/progress/startAttempt.go) | `POST /api/v1/lessons/{lesson_id}/attempts/start` | Старт новой попытки тестирования с серверным контролем лимитов |
| [`getLessonProgress.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/progress/getLessonProgress.go) | `GET /api/v1/lessons/{lesson_id}/progress` | Текущий прогресс по уроку |
| [`getCourseProgress.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/progress/getCourseProgress.go) | `GET /api/v1/courses/{course_id}/progress` | Сводный прогресс по курсу |
| [`getAllLessonProgress.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/progress/getAllLessonProgress.go) | `GET /api/v1/courses/{course_id}/progress/lessons` | Прогресс по каждому уроку курса |

---

## 🔗 Зависимости
- **Входящие:** [`internal/app/di.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/app/di.go).
- **Исходящие:** [`service.ProgressServices`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/progress), [`internal/dto`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/dto).
