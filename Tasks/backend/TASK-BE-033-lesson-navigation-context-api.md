# 🛠 [BE-033] API контекста навигации урока и структуры силлабуса курса

> **Статус:** Completed  
> **Приоритет:** Critical (P0)  
> **Связанные задачи:** FE-035, QA-035  
> **Целевой модуль:** `internal/interfaces/handlers/lesson/`, `internal/service/lesson/`, `internal/repository/lesson/`  
> **Документация модуля:** [internal/service/lesson/README.md](../../internal/service/lesson/README.md)

---

## 🎯 Цель задачи
Обеспечить плеер урока ([`/lessons/[id]`](../../frontend/src/app/lessons/%5Bid%5D/page.tsx)) данными для бесшовного перемещения студента по курсу:
1. Вычисление идентификатора и названия **следующего** (`next_lesson`) и **предыдущего** (`prev_lesson`) урока на основе порядка следования секций и позиций уроков в курсе.
2. Предоставление компактного дерева курса (`course_context`: название курса, список модулей и уроков со статусом завершенности для текущего студента), чтобы плеер мог отображать выпадающее оглавление без выполнения нескольких тяжелых параллельных запросов.

---

## 🔍 Текущее состояние кода
- Сейчас в `GET /api/v1/lessons/{id}` возвращается только объект урока (id, title, content, section_id, course_id).
- Клиент вынужден делать дополнительные тяжелые запросы к `GET /courses/{course_id}/structure` и `GET /courses/{course_id}/progress`, что замедляет инициализацию плеера и часто приводит к неконсистентности данных о следующем уроке.

---

## 📝 Технические требования к реализации

### 1. HTTP Endpoint
`GET /api/v1/lessons/{id}/navigation`
* **Авторизация:** Bearer JWT (студент должен быть зачислен на курс либо иметь роль преподавателя/автора/админа).
* **Параметры пути:** `id` (int, ID текущего урока).

### 2. Response DTO
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "current_lesson": {
      "id": 142,
      "title": "Интерактивный практикум: Орфография",
      "position": 3,
      "section_id": 12
    },
    "course": {
      "id": 13,
      "title": "Русский язык: Полный курс подготовки к ЕГЭ",
      "slug": "russkiy-yazyk-ege"
    },
    "prev_lesson": {
      "id": 141,
      "title": "Теория: Чередующиеся гласные в корне"
    },
    "next_lesson": {
      "id": 143,
      "title": "Контрольный тест по разделу"
    },
    "syllabus": [
      {
        "section_id": 12,
        "section_title": "Модуль 1. Орфография",
        "position": 1,
        "lessons": [
          {
            "id": 140,
            "title": "Введение в блок орфографии",
            "position": 1,
            "is_completed": true,
            "score": 100
          },
          {
            "id": 141,
            "title": "Теория: Чередующиеся гласные в корне",
            "position": 2,
            "is_completed": true,
            "score": 90
          },
          {
            "id": 142,
            "title": "Интерактивный практикум: Орфография",
            "position": 3,
            "is_completed": false,
            "score": 0
          },
          {
            "id": 143,
            "title": "Контрольный тест по разделу",
            "position": 4,
            "is_completed": false,
            "score": 0
          }
        ]
      }
    ]
  }
}
```
* Если текущий урок первый в курсе ➔ `"prev_lesson": null`.
* Если текущий урок последний в курсе ➔ `"next_lesson": null`.

### 3. Бизнес-логика (Service Layer)
- Метод `GetLessonNavigationContext(ctx context.Context, userID, lessonID int) (*dto.LessonNavigationResponse, error)`.
- Алгоритм:
  1. Найти урок и связанный `course_id`.
  2. Проверить доступ пользователя к курсу (`accessService.CanViewCourse`).
  3. Извлечь упорядоченный плоский список всех уроков курса по `ORDER BY sections.position ASC, lessons.position ASC`.
  4. Найти индекс текущего урока `currIdx`.
  5. Определить соседа слева (`currIdx - 1`) и соседа справа (`currIdx + 1`).
  6. Собрать статус прогресса уроков для студента `userID` через `progressRepo.GetLessonsProgressByCourse(ctx, userID, courseID)`.
  7. Сформировать иерархический DTO силлабуса.

---

## ✅ Критерии приёмки (Definition of Done)
- [x] Запрос `GET /api/v1/lessons/{id}/navigation` выполняется быстро (< 50 мс) за счет единого CTE-запроса или двух точечных выборок.
- [x] Граничные уроки (первый и последний) отдают корректный `null` без паники.
- [x] Права доступа валидируются (чужой или незачисленный студент получает 403 Forbidden).
