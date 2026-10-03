# 📋 Функциональная спецификация: `Зачисление и студенты (Enrollment)`

> **Расположение:** `internal/service/enrollment`  
> **Технический контекст:** [`README.md`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/enrollment/README.md)  
> **Дата актуализации:** 2026-10-04  

---

## 🎯 Каталог бизнес-фич модуля
| Фича | Описание возможности | Обеспечивающие функции |
|---|---|---|
| **Самозапись студента** | Самостоятельное поступление учащегося на открытый курс с проверкой дубликатов | `SelfEnrollCourse()` |
| **Административное зачисление** | Ручное зачисление студента или назначение ассистента преподавателем по Email или ID | `TeacherEnrollCourse()` |
| **Отчисление с курса** | Отчисление студента с курса с уменьшением счетчика учащихся и защитой создателя | `UnenrollUser()` |
| **Смена роли на курсе** | Назначение роли `teacher` или `student` участнику курса | `ChangeUserRole()` |
| **Список учащихся с успеваемостью** | Просмотр студентов курса с прогрессом, средним баллом и индикатором долгов по ДЗ | `ListCourseStudentsWithProgress()` |

---

## 🔬 Паспорта функций

### ⚡ Функция: `SelfEnrollCourse(ctx, req)`

* **Файл и строки:** [`enrollUserToCourse.go#L14-L46`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/enrollment/enrollUserToCourse.go#L14-L46)
* **Бизнес-назначение:** Запись студента на курс в один клик.
* **Связанная фича:** *Самозапись студента*

#### 📥 Входные параметры
| Параметр | Тип | Обязателен | Бизнес-смысл и ограничения |
|---|---|:---:|---|
| `req.UserID` | `int64` | Да | ID зачисляемого студента |
| `req.CourseID` | `int64` | Да | ID курса |

#### 🔄 Пошаговый алгоритм работы
1. **Шаг 1 (Проверка доступности):** Вызывает `course.CanSelfEnroll()`. Курс должен быть `published` и `public`. Иначе `ErrForbidden`.
2. **Шаг 2 (Проверка дубликата):** Вызывает `UserExistCourse`. Если студент уже записан ➔ `ErrUserAlreadyEnrolled`.
3. **Шаг 3 (Транзакционное зачисление):** Внутри `executeEnrollmentTransaction`:
   - Добавляет строку в `users_courses` с ролью `student`.
   - Инкрементирует поле `enrolled_count` в таблице `courses`.

---

### ⚡ Функция: `TeacherEnrollCourse(ctx, req)`

* **Файл и строки:** [`enrollUserToCourse.go#L48-L102`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/enrollment/enrollUserToCourse.go#L48-L102)
* **Бизнес-назначение:** Ручное добавление студентов или ассистентов преподавателем (например, корпоративные группы).
* **Связанная фича:** *Административное зачисление*

#### 🔄 Пошаговый алгоритм работы
1. **Шаг 1 (Авторизация учителя):** Проверяет право преподавателя `accessService.CanManageCourseUsers`.
2. **Шаг 2 (Поиск целевого пользователя):** Если передан `TargetUserID` — ищет по ID; если передан `TargetEmail` — ищет по email. Если пользователь не существует ➔ ошибка валидации.
3. **Шаг 3 (Проверка дубликата):** Проверяет отсутствие в `users_courses`.
4. **Шаг 4 (Зачисление):** Транзакционно добавляет в `users_courses` с переданной ролью (`student` или `teacher`) и инкрементирует `enrolled_count`.

---

### ⚡ Функция: `UnenrollUser(ctx, userID, courseID)`

* **Файл и строки:** [`unenrollUser.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/enrollment/unenrollUser.go)
* **Бизнес-назначение:** Отчисление учащегося с курса.
* **Связанная фича:** *Отчисление с курса*

#### 🔄 Пошаговый алгоритм работы
1. **Шаг 1 (Проверка роли):** Запрашивает текущую роль пользователя на курсе через `enrolledrepo.GetRoleUserInCourse`.
2. **Шаг 2 (Защита создателя):** Если роль `creator` ➔ **Запрещено (`ErrCreatorCannotUnenroll`)**.
3. **Шаг 3 (Транзакционное удаление):** В транзакции:
   - Удаляет запись пользователя из `users_courses`.
   - Безопасно декрементирует счетчик `enrolled_count` с гарантией неотрицательности: `GREATEST(enrolled_count - 1, 0)`.

---

### ⚡ Функция: `TeacherUnenrollUser(ctx, teacherID, targetUserID, courseID)`

* **Файл и строки:** [`unenrollUser.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/enrollment/unenrollUser.go)
* **Бизнес-назначение:** Принудительное отчисление студента преподавателем курса.
* **Связанная фича:** *Отчисление с курса*

#### 🔄 Пошаговый алгоритм работы
1. **Шаг 1 (Проверка прав преподавателя):** Проверяет право преподавателя `accessService.CanManageCourseUsers(ctx, course, teacherID)`. Если прав нет ➔ `ErrForbidden`.
2. **Шаг 2 (Делегирование в UnenrollUser):** Вызывает `UnenrollUser(ctx, targetUserID, courseID)` со всеми встроенными проверками защиты роли `creator` и атомарным декрементом счетчика.

---

#### 💡 Подсказка для аналитика (Где менять логику?)
* *Ограничить максимальное количество студентов на курсе:* добавить проверку `course.EnrolledCount < course.MaxStudents` перед зачислением в [`enrollUserToCourse.go#L33`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/enrollment/enrollUserToCourse.go#L33).
