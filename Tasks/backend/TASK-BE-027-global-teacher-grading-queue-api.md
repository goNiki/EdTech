# 🛠 [BE-027] Глобальный эндпоинт очереди проверки домашних заданий преподавателя (Cross-Course Aggregation)

**Статус:** Completed

> **Приоритет:** High (P1)  
> **Связанные задачи:** FE-026, QA-026  
> **Целевой модуль:** `internal/interfaces/handlers/analytics/`, `internal/service/analytics/`, `internal/repository/analytics/`  
> **Документация модуля:** [internal/service/analytics/README.md](../../internal/service/README.md)

---

## 🎯 Цель задачи
Реализовать единый серверный эндпоинт агрегированной очереди проверки заданий для преподавателя.
Сейчас в системе существует только эндпоинт конкретного курса: `GET /api/v1/courses/{courseid}/grading/pending`.
Из-за этого глобальная страница проверки заданий (`/teacher/grading`) вынуждена наугад выбирать один курс из списка, что приводит к ложному статусу «Все домашние задания проверены!», если у первого выбранного курса нет сданных работ, даже если на других курсах автора ожидают проверки десятки студенческих ответов.

Новый API позволит получить:
1. Все ожидающие проверки работы (`qaa.is_correct IS NULL`) по **всем курсам**, где текущий пользователь является автором или назначенным преподавателем.
2. Опциональную фильтрацию по конкретному `course_id`.
3. Сводку счетчиков по каждому курсу (`courses_summary`), чтобы интерфейс мог отобразить бейджи количества не проверенных заданий напротив каждого курса в выпадающем списке.

---

## 🔍 Текущее состояние кода
- В `internal/interfaces/handlers/analytics/list_pending_hw.go` метод `ListPendingHomeworks` жестко считывает параметр пути `courseid` из Chi router (`chi.URLParam(r, "courseid")`).
- В `internal/repository/analytics/pending.go` SQL-запрос `WHERE qz.course_id = $1 AND qaa.is_correct IS NULL` работает только для одного `course_id`.
- В `internal/app/di.go` зарегистрирован только маршрут `/api/v1/courses/{courseid}/grading/pending`.

---

## 📝 Технические требования к реализации

### 1. API / Endpoints
Зарегистрировать маршрут в `internal/app/di.go` (под мидлварью авторизации):
- `GET /api/v1/teacher/grading/pending`
  - **Query-параметры:**
    - `course_id` (опционально, `int64`): если передан и `> 0`, фильтрует по конкретному курсу. Если не передан или `0` — возвращает работы по всем курсам преподавателя.
    - `page` (`int64`, default `1`).
    - `page_size` (`int64`, default `20`, max `100`).
  - **Response DTO (`dto.PaginatedTeacherPendingHomeworksResponse`):**
    ```json
    {
      "items": [
        {
          "attempt_id": 14,
          "answer_id": 42,
          "student_id": 5,
          "student_name": "Иван Смирнов",
          "student_email": "student@example.com",
          "student_username": "smirnov",
          "course_id": 2,
          "course_title": "Русский язык",
          "lesson_id": 10,
          "lesson_title": "Синтаксис и пунктуация",
          "question_id": 7,
          "question_text": "Напишите развернутое сочинение-рассуждение",
          "student_answer": "Текст работы студента...",
          "attachment_url": "https://storage.../essay.pdf",
          "max_points": 25,
          "rubric": "Оценивается структура и логика",
          "submitted_at": "2026-10-04T10:15:00Z"
        }
      ],
      "total": 4,
      "page": 1,
      "page_size": 20,
      "courses_summary": [
        {
          "course_id": 2,
          "course_title": "Русский язык",
          "pending_count": 4
        },
        {
          "course_id": 5,
          "course_title": "Основы Go",
          "pending_count": 0
        }
      ]
    }
    ```

### 2. Бизнес-логика (Service Layer)
- В `internal/service/analytics/`:
  - Создать метод `ListTeacherPendingHomeworks(ctx, teacherID, courseID, page, pageSize)`:
    1. Проверить права: пользователь должен иметь роль `teacher` или `admin`.
    2. Если `courseID > 0`, проверить, что преподаватель имеет доступ к этому курсу (`checkTeacherAccess`).
    3. Вызвать репозиторий для получения агрегированного списка и сводки по курсам.

### 3. Репозиторий (Repository Layer)
- В `internal/repository/analytics/`:
  - Реализовать запрос выборки через JOIN с таблицей курсов:
    ```sql
    SELECT 
        qa.id AS attempt_id,
        qaa.id AS answer_id,
        u.id AS student_id,
        COALESCE(NULLIF(TRIM(u.first_name || ' ' || u.last_name), ''), u.username) AS student_name,
        u.email AS student_email,
        u.username AS student_username,
        c.id AS course_id,
        c.title AS course_title,
        l.id AS lesson_id,
        l.title AS lesson_title,
        qq.id AS question_id,
        qq.question_text AS question_text,
        COALESCE(qaa.user_answer, '') AS student_answer,
        qaa.attachment_url,
        COALESCE(qz.points, 10) AS max_points,
        qa.created_at AS submitted_at
    FROM quiz_attempt_answers qaa
    JOIN quiz_attempts qa ON qaa.attempt_id = qa.id
    JOIN quizzes qz ON qa.quiz_id = qz.id
    JOIN courses c ON qz.course_id = c.id
    JOIN lessons l ON qz.lesson_id = l.id
    JOIN quiz_questions qq ON qaa.question_id = qq.id
    JOIN users u ON qa.user_id = u.id
    WHERE qaa.is_correct IS NULL
      AND (c.creator_id = $1 OR EXISTS (
          SELECT 1 FROM course_instructors ci WHERE ci.course_id = c.id AND ci.user_id = $1
      ))
      AND ($2 = 0 OR c.id = $2)
    ORDER BY qa.created_at ASC
    LIMIT $3 OFFSET $4
    ```
  - Реализовать запрос `courses_summary`:
    ```sql
    SELECT c.id, c.title, COUNT(qaa.id) FILTER (WHERE qaa.is_correct IS NULL) AS pending_count
    FROM courses c
    LEFT JOIN quizzes qz ON qz.course_id = c.id
    LEFT JOIN quiz_attempts qa ON qa.quiz_id = qz.id
    LEFT JOIN quiz_attempt_answers qaa ON qaa.attempt_id = qa.id AND qaa.is_correct IS NULL
    WHERE (c.creator_id = $1 OR EXISTS (
        SELECT 1 FROM course_instructors ci WHERE ci.course_id = c.id AND ci.user_id = $1
    ))
    GROUP BY c.id, c.title
    ORDER BY pending_count DESC, c.title ASC
    ```

---

## ✅ Критерии приёмки (Definition of Done)
- [x] Запрос `GET /api/v1/teacher/grading/pending` без параметров возвращает все непроверенные работы преподавателя по всем его курсам, отсортированные от старых к новым (FIFO).
- [x] Запрос `GET /api/v1/teacher/grading/pending?course_id=2` фильтрует работы только указанного курса.
- [x] Поле `courses_summary` возвращает точный расклад счетчиков непроверенных работ по каждому курсу преподавателя.
- [x] Доступ защищен авторизацией (401 для гостей, 403 при попытке запросить чужой `course_id`).
- [x] Ошибки БД и пустые результаты обрабатываются без паник, возвращая пустые слайсы `[]`.
