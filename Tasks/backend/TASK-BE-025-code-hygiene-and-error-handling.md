# 🛠 [BE-025] Комплексная чистка кодовой базы: Error Handling, пакетная гигиена и оптимизация аллокаций

**Статус:** Completed

> **Приоритет:** High (P1 / Code Quality & Reliability)  
> **Аудит:** [CODE_REVIEW_GOLANG.md](../../CODE_REVIEW_GOLANG.md) (Дефекты 2.4, 2.5, 2.7, 2.8, 2.9)  
> **Целевые модули:**  
> - `pkg/errors/`, `pkg/utils/`  
> - `internal/interfaces/handlers/courses/get.go`  
> - `internal/domain/course.go`  
> - `internal/service/auth/`, `internal/service/quiz/`  
> - `cmd/migration/migration.go`  

---

## 🎯 Цель задачи
Устранить накопленный технический долг, антипаттерны обработки ошибок и нарушения границ слоев:
1. Заменить хрупкие строковые проверки ошибок (`err.Error() != "..."`) на `errors.Is`.
2. Вынести компиляцию регулярного выражения слага из метода валидации в переменную пакета.
3. Разорвать недопустимую зависимость `pkg/utils` -> `internal/dto`.
4. Привести именование пакета `pkg/errors` к стандартам Go и удалить дублирующиеся алиасы ошибок.
5. Устранить опасное подавление ошибок через `_ = ` в критических путях аутентификации и прохождения квизов.

---

## 🔍 Текущее состояние кода и требования

### 1. Строковые проверки ошибок (Дефект 2.4)
- **Где:** [`internal/interfaces/handlers/courses/get.go#L53, L104`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/courses/get.go#L53).
- **Было:**
  ```go
  if err.Error() != "user not found" && err.Error() != "not found" && err.Error() != "no rows in result set"
  ```
- **Стало:**
  Сервис `GetRoleUserInCource` уже возвращает `"", nil`, если пользователь не записан на курс. В хэндлере убрать проверку строк. Любая вернувшаяся ошибка является реальным сбоем и должна обрабатываться через `response.HandleError`.

### 2. Аллокация регулярки в методе валидации (Дефект 2.5)
- **Где:** [`internal/domain/course.go#L139`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/domain/course.go#L139).
- **Было:** `slugRegex := regexp.MustCompile("^[a-z0-9-]+$")` внутри `c.Validate()`.
- **Стало:** Вынести на уровень пакета `var slugRegex = regexp.MustCompile("^[a-z0-9-]+$")`.

### 3. Нарушение границ слоев в `pkg/utils` (Дефект 2.7)
- **Где:** [`pkg/utils/validate.go#L4`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/pkg/utils/validate.go#L4).
- **Проблема:** Пакет `pkg/utils` импортирует `edtech/internal/dto` и содержит устаревшие функции валидации с `// TODO`.
- **Решение:** Перенести валидацию `ValidateEnrolle(dto.EnrollRequest)` в слой хэндлера / DTO валидатора (`go-playground/validator`). Удалить зависимость `pkg/` от `internal/`. Разделить пакет `utils` на специализированные пакеты: `pkg/slug` и `pkg/uuid`.

### 4. Именование пакета ошибок и дубли констант (Дефект 2.8)
- **Где:** [`pkg/errors/errorsAPP.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/pkg/errors/errorsAPP.go), [`pkg/errors/auth.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/pkg/errors/auth.go).
- **Было:** `package errorsAPP` (не совпадает с папкой `errors`, camelCase). Дубли `ErrNotFoundCourse = ErrCourseNotFound`, `ErrNotFoundLesson = ErrLessonNotFound`.
- **Стало:** Переименовать пакет в `package apperrors` (или `package errors`). Удалить дубли, унифицировать использование единого имени по всему проекту (например, `ErrCourseNotFound`).

### 5. Подавление ошибок через `_ = ` (Дефект 2.9)
- **Где:**
  - [`internal/service/auth/login.go#L41`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/auth/login.go#L41): `_ = s.repo.UpdateLastLogin(...)` -> логировать ошибку через `slog.Warn`, если обновление последнего входа не удалось, а не подавлять без следа.
  - [`internal/service/auth/refresh.go#L31`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/auth/refresh.go#L31): `_ = s.refreshRepo.DeleteRefreshToken(...)` -> возвращать ошибку или логировать.
  - [`internal/service/quiz/submitAttempt.go#L125-L126`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/quiz/submitAttempt.go#L125-L126): `_ = s.progressService.CompleteLesson(...)` -> проверять ошибку и прерывать выполнение/откатывать транзакцию при сбое!
  - [`cmd/migration/migration.go#L23-L25`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/cmd/migration/migration.go#L23-L25): пустой `if err != nil { }` -> логировать ошибку загрузки конфига и завершать процесс `log.Error(...) os.Exit(1)`.

---

## ✅ Критерии приёмки (Definition of Done)
- [x] В проекте нет ни одного сравнения ошибок через `err.Error() == "..."` или `!= "..."`.
- [x] Регулярное выражение слага компилируется ровно 1 раз при старте пакета.
- [x] Пакеты в `pkg/` не имеют импортов из `internal/`.
- [x] Имя пакета ошибок приведено к единому стилю (`apperrors`), дублирующие алиасы удалены.
- [x] В критических путях аутентификации и сабмита квизов все возвращаемые ошибки проверяются.
- [x] `go vet ./...` и `go test ./...` проходят без замечаний.
