# 📦 Обработчики: `internal/interfaces/handlers/courses`
> **Путь:** `internal/interfaces/handlers/courses`  
> **Роль:** HTTP-контроллеры каталога курсов, создания, управления статусами и структурой курса.

---

## 🎯 Назначение и ответственность
Обрабатывает HTTP-запросы по курсам (`/api/v1/courses/*`). Парсит параметры строки запроса (пагинация, поисковые фильтры), проверяет роли и отдает детальную информацию с рассчитанными правами (`CoursePermissions`).

---

## ⚠️ Жестко установленные правила (Hard Invariants)
1. **Создание доступно только преподавателям:** `CreateCourse` строго проверяет роль в контексте: `userRole != string(domain.RoleTeacher) ➔ 403 Forbidden`.
2. **Опциональный контекст для каталога:** Публичные эндпоинты подключены через `OptionalJWTMiddleware`. Если запрос делает гость, `userID == 0`; если авторизованный пользователь — в ответе `CourseDetailResponse` возвращаются персональные флаги прав (`can_view`, `can_edit`, `can_publish`, `can_enroll`).
3. **Безопасный парсинг Query-параметров:** `query_parser.go` нормализует параметры `page` (дефолт 1) и `pagesize` (дефолт 10, максимум 100).

---

## 📁 Структура файлов и Эндпоинты
| Файл | HTTP Метод и Путь | Описание |
|---|---|---|
| [`listpublic.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/courses/listpublic.go) | `GET /api/v1/courses` | Каталог опубликованных курсов с фильтрами |
| [`get.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/courses/get.go) | `GET /api/v1/courses/{courseid}`<br>`GET /courses/slug/{slug}` | Детальная страница курса с объектом прав |
| [`create.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/courses/create.go) | `POST /api/v1/courses` | Создание нового курса преподавателем |
| [`publish.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/courses/publish.go) | `POST /api/v1/courses/{courseid}/publish` | Публикация курса |
| [`archive.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/courses/archive.go) | `POST /api/v1/courses/{courseid}/archive` | Архивация курса |
| [`delete.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/courses/delete.go) | `DELETE /api/v1/courses/{courseid}` | Мягкое удаление курса создателем |
| [`reorder_sections.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/courses/reorder_sections.go) | `PUT /api/v1/courses/{courseid}/reorder-sections` | Изменение порядка модулей |
| [`list_my.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/courses/list_my.go) | `GET /api/v1/courses/my` | Мои курсы (для студента и автора) |
| [`query_parser.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/courses/query_parser.go) | Утилита | Парсинг пагинации и query параметров URL |

---

## 🔗 Зависимости
- **Входящие:** [`internal/app/di.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/app/di.go).
- **Исходящие:** [`service.CourseServices`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/course), [`service.AccessService`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/access), [`internal/dto`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/dto).
