# 🛠 [BE-036] Публичный профиль автора курса и расширенные метаданные преподавателя в API

> **Приоритет:** High (P1)  
> **Связанные задачи:** FE-038, QA-038  
> **Целевой модуль:** `internal/interfaces/handlers/courses/`, `internal/service/course/`, `internal/repository/course/`  
> **Документация модуля:** [internal/service/course/README.md](../../internal/service/course/README.md)

---

## 🎯 Цель задачи
Обеспечить передачу расширенной информации о преподавателе / авторе курса в публичных эндпоинтах каталога и лендинга курса:
1. Имя, фамилия, отображаемое имя автора (`author_name`).
2. Аватар преподавателя (`author_avatar_url`).
3. Должность / специализация / регалии (`author_headline`, например: «Ведущий эксперт ЕГЭ, кандидат филологических наук, 12 лет стажа»).
4. Биография преподавателя (`author_bio`).
5. Статистика автора: количество созданных курсов (`courses_count`), общее число обученных студентов (`total_students_count`).

---

## 🔍 Текущее состояние кода
- В схеме БД таблица `courses` содержит `created_by INT REFERENCES users(id)`.
- В `GET /api/v1/courses/{id}` и `GET /api/v1/courses/slug/{slug}` возвращается только числовой `created_by`, без объединения с таблицей `users` и профилем. На странице курса клиент не имеет данных об авторе.

---

## 📝 Технические требования к реализации

### 1. Доработка SQL в `CourseRepository`
В запросы выборки детального курса (`GetCourseByID`, `GetCourseBySlug`) добавить `JOIN users u ON c.created_by = u.id`:
```sql
SELECT 
    c.id, c.title, c.slug, c.description, c.cover_url, c.status, c.created_by,
    u.name AS author_name,
    u.avatar_url AS author_avatar_url,
    COALESCE(u.headline, '') AS author_headline,
    COALESCE(u.bio, '') AS author_bio,
    (SELECT COUNT(*) FROM courses WHERE created_by = u.id AND status = 'published') AS author_courses_count,
    (SELECT COALESCE(SUM(enrolled_count), 0) FROM courses WHERE created_by = u.id) AS author_total_students
FROM courses c
JOIN users u ON c.created_by = u.id
WHERE ...
```

### 2. Response DTO
В объект курса добавить вложенную структуру `author`:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "course": {
      "id": 13,
      "title": "Русский язык: Полный курс подготовки к ЕГЭ",
      "author": {
        "id": 2,
        "name": "Никита Преподаватель",
        "avatar_url": "/static/uploads/avatars/nikita.jpg",
        "headline": "Ведущий методист и преподаватель высшей категории",
        "bio": "Более 8 лет готовлю выпускников к сдаче экзаменов на 90+ баллов...",
        "courses_count": 3,
        "total_students": 1420
      }
    }
  }
}
```

### 3. Обновление профиля автора (`PATCH /api/v1/auth/profile`)
- Добавить поддержку сохранения поля `headline` (до 150 символов) в таблице `users` и эндпоинте обновления своего профиля преподавателем.

---

## ✅ Критерии приёмки (Definition of Done)
- [x] Публичные эндпоинты курса (`/courses/{id}` и `/courses/slug/{slug}`) возвращают заполненный блок `author`.
- [x] Преподаватель может обновлять свой `headline` через `/auth/profile`.
- [x] Неавторизованный гость видит карточку автора без ошибок 401/403.
