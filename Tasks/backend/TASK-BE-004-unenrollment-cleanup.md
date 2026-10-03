# 🛠 [BE-004] Аудит и обеспечение консистентности при отчислении студентов из курса

**Статус:** ✅ Completed (хэш коммита: `38129ab`)

> **Приоритет:** High (P1)  
> **Связанные задачи:** FE-004, QA-004  
> **Целевой модуль:** `internal/service/enrollment/`, `internal/repository/enrollment/`  
> **Документация модуля:** [internal/service/enrollment/FUNCTIONAL_SPEC.md](../../internal/service/enrollment/FUNCTIONAL_SPEC.md)

## 🎯 Цель задачи
Обеспечить полную целостность данных при самоотчислении студента (`DELETE /courses/{courseid}/enroll`) и принудительном отчислении преподавателем (`DELETE /courses/{courseid}/students/{userid}`). Проверить, что счетчик `enrolled_count` декрементируется атомарно (не уходит в минус), права создателя защищены от отчисления, а статус в `users_courses` архивируется или удаляется корректно.

## 🔍 Текущее состояние кода
- Эндпоинты `SelfEnroll` и `TeacherUnenroll` зарегистрированы в `internal/app/di.go` ([L473, L476](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/app/di.go#L473)).
- В `internal/service/enrollment/unenrollUser.go`:
  - Запрет отчисления автора: `if role == "creator" { return errorsAPP.ErrCreatorCannotUnenroll }`.
  - Выполняются вызовы `enrolledrepo.UnenrollUser` и `courserepo.DecrementEnrolledCount`.
- Требуется аудит: убедиться, что `DecrementEnrolledCount` имеет защиту `GREATEST(enrolled_count - 1, 0)` и что прогресс студента по курсу сохраняется или помечается статусом `archived`.

## 📝 Технические требования к реализации
1. **API / Endpoints:**
   - `DELETE /api/v1/courses/{courseid}/enroll` (самоотчисление текущего пользователя)
   - `DELETE /api/v1/courses/{courseid}/students/{userid}` (отчисление преподавателем)
   - При успешном выполнении возвращать:
     ```json
     { "message": "successfully unenrolled from course" }
     ```

2. **Бизнес-логика (Service Layer):**
   - Проверка в `TeacherUnenroll`: преподаватель может отчислять только из своего курса (`accessService.CanEditCourse`).
   - Запрет отчисления пользователя с ролью `creator` или преподавателя с активными правами автора.
   - Защита счетчика в `DecrementEnrolledCount`:
     ```sql
     UPDATE courses SET enrolled_count = GREATEST(enrolled_count - 1, 0) WHERE id = $1
     ```

## ✅ Критерии приёмки (Definition of Done)
- [x] Студент может отчислиться сам, счетчик студентов уменьшается на 1.
- [x] Преподаватель может удалить студента из списка курса.
- [x] Попытка отчислить создателя курса отклоняется с ошибкой 400 (`ErrCreatorCannotUnenroll`).
- [x] Счетчик `enrolled_count` никогда не становится меньше нуля.
