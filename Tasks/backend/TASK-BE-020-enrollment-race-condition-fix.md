# 🛠 [BE-020] Атомарная запись на курс и устранение Race Condition (TOCTOU)

**Статус:** ✅ Completed

> **Приоритет:** Blocker (P0 / Data Consistency)  
> **Аудит:** [CODE_REVIEW_GOLANG.md](../../CODE_REVIEW_GOLANG.md) (Дефект 1.5)  
> **Целевой модуль:** `internal/service/enrollment/`, `internal/repository/enrollment/`, `internal/repository/course/`  

---

## 🎯 Цель задачи
Устранить состояние гонки (Race Condition / TOCTOU) при записи пользователей на курс. Обеспечить строгую атомарность добавления записи в `users_courses` и инкремента счетчика `courses.enrolled_count`, предотвратив повторную запись и рассинхрон аналитики.

---

## 🔍 Текущее состояние кода
В [`internal/service/enrollment/enrollUserToCourse.go#L26-L40`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/enrollment/enrollUserToCourse.go#L26-L40):
```go
isEnrolled, err := s.enrolledrepo.UserExistCourse(ctx, s.db, req.UserID, req.CourseID)
// ...
if isEnrolled {
    return fmt.Errorf("%s: %w", op, errorsAPP.ErrUserAlreadyEnrolled)
}
// ...
err = s.executeEnrollmentTransaction(ctx, enroll)
```
- Проверка `UserExistCourse` выполняется **вне транзакции** через общий пул соединений `s.db`.
- Два конкурентных клика/запроса пользователя одновременно читают `isEnrolled == false`.
- Оба запроса переходят к выполнению транзакции `executeEnrollmentTransaction`, что приводит либо к двойному инкременту `enrolled_count`, либо к ошибкам 500 из-за падения уникального ограничения `unique_user_course`.

---

## 📝 Технические требования к реализации
1. **Атомарная вставка с `ON CONFLICT`:**  
   В `internal/repository/enrollment/`:
   ```sql
   INSERT INTO users_courses (user_id, course_id, role, enrolled_at)
   VALUES ($1, $2, $3, NOW())
   ON CONFLICT (user_id, course_id) DO NOTHING;
   ```
   Метод репозитория возвращает булев флаг (была ли создана новая запись) на основе `tag.RowsAffected() == 1`.
2. **Транзакционная бизнес-логика:**  
   В методе `SelfEnrollCourse` и `TeacherEnrollCourse`:
   - Вся логика выполняется внутри транзакции `s.txManager.WithTX(ctx, ...)`.
   - Если запись уже существует (`inserted == false`), возвращать доменную ошибку `errorsAPP.ErrUserAlreadyEnrolled` без инкремента счетчика курса.
   - Инкремент `s.courserepo.IncrementEnrolledCount` вызывается строго в той же транзакции только если `inserted == true`.

---

## ✅ Критерии приёмки (Definition of Done)
- [x] Конкурентные запросы на запись одного и того же пользователя не приводят к двойному инкременту `enrolled_count`.
- [x] При повторной записи пользователь получает ошибку 409 Conflict (`ErrUserAlreadyEnrolled`), а не 500 Internal Server Error.
- [x] Написан конкурентный интеграционный/unit тест (`t.Parallel()`, несколько горутин с параллельной записью).
- [x] Тесты `go test ./internal/service/enrollment/...` проходят без ошибок и гонок (`go test -race`).
