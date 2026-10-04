# 🛠 [BE-023] Декомпозиция монолитных интерфейсов сервисов по потребителям

**Статус:** ✅ Completed

> **Приоритет:** Medium (P2 / Architecture & Clean Code)  
> **Аудит:** [CODE_REVIEW_GOLANG.md](../../CODE_REVIEW_GOLANG.md) (Дефект 2.2)  
> **Целевые модули:** `internal/service/`, `internal/interfaces/handlers/`  

---

## 🎯 Цель задачи
Устранить избыточное зацепление (Interface Pollution) монолитного файла `internal/service/service.go`.  
Разбить гигантские интерфейсы сервисов на целевые интерфейсы, соответствующие реальным сценариям использования потребителями (обработчиками и смежными сервисами), без искусственного дробления до 1 метода.

---

## 🔍 Текущее состояние кода
В [`internal/service/service.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/service.go):
- Все 10 интерфейсов системы объявлены в одном файле.
- Например, `CourseServices` содержит 13 методов: от публичного каталога (`ListPublicCourses`, `GetCourseBySlug`) до авторского редактирования (`CreateCourse`, `UpdateCourse`, `ReorderSections`).
- Хэндлеры и смежные сервисы зависят от полного монструозного интерфейса, что усложняет написание mock-объектов и модульное тестирование.

---

## 📝 Технические требования к реализации
1. Вынести интерфейсы из единого файла `service/service.go` в целевые пакеты сервисов либо объявить их на стороне потребителей (в хэндлерах):
   - **Публичное чтение курсов (Reader / Catalog):**
     - `GetCourseByID`, `GetCourseBySlug`, `GetCourseStructure`, `ListPublicCourses`
   - **Управление курсами (Manager / Author):**
     - `CreateCourse`, `UpdateCourse`, `UpdateCourseStatus`, `PublishCourse`, `ArchiveCourse`, `DeleteCourse`, `ReorderSections`
   - **Студенческие курсы (MyCourses):**
     - `ListMyCourses`
2. Аналогично сгруппировать интерфейсы для `AuthService` (Auth / Profile / Admin), `EnrolledServices` (SelfEnroll / TeacherManagement) и `ProgressServices`.
3. Обновить внедрение зависимостей в хэндлерах (`CourseHandler`, `ProgressHandler` и т.д.): они должны требовать только тот интерфейс, чьи методы фактически вызываются.

---

## ✅ Критерии приёмки (Definition of Done)
- [x] Монолитный файл `service/service.go` очищен или разделен на логические интерфейсные контракты.
- [x] Хэндлеры зависят от сфокусированных интерфейсов под свои сценарии.
- [x] Облегчено создание mock-ов в unit-тестах.
- [x] Весь проект собирается без ошибок (`go build ./...`).
