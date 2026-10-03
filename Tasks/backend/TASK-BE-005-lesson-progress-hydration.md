# 🛠 [BE-005] Оптимизация эндпоинта получения прогресса урока и отдачи истории прохождения

> **Статус:** Completed  
> **Приоритет:** Medium (P1)  
> **Связанные задачи:** FE-005, QA-005  
> **Целевой модуль:** `internal/interfaces/handlers/progress/`, `internal/service/progress/`  
> **Документация модуля:** [internal/service/progress/FUNCTIONAL_SPEC.md](../../internal/service/progress/FUNCTIONAL_SPEC.md)

## 🎯 Цель задачи
Обеспечить фронтенд полной информацией о состоянии прохождения конкретного урока студентом. Эндпоинт `GET /api/v1/lessons/{lesson_id}/progress` должен отдавать статус прохождения (`not_started`, `in_progress`, `completed`), набранный балл, количество потраченного времени и историю сданных ответов на эссе/задания, чтобы клиент мог гидратировать состояние интерфейса.

## 🔍 Текущее состояние кода
- В `internal/interfaces/handlers/progress/getLessonProgress.go` метод возвращает DTO `LessonProgress`.
- В `internal/service/progress/getProgress.go` возвращается базовая модель `LessonProgress`.
- Требуется расширить выдачу: прикреплять детали о проверке сданных эссе (если работа уже проверена преподавателем: выставленный балл и текстовый фидбек преподавателя).

## 📝 Технические требования к реализации
1. **API / Endpoint:**
   - `GET /api/v1/lessons/{lesson_id}/progress`
   - Headers: `Authorization: Bearer <token>`
   - Response DTO (200 OK):
     ```json
     {
       "lesson_id": 42,
       "status": "completed",
       "score": 95,
       "time_spent": 1200,
       "last_position": 0,
       "started_at": "2026-10-01T10:00:00Z",
       "completed_at": "2026-10-01T10:20:00Z",
       "submissions": [
         {
           "question_text": "Практическое задание №1",
           "student_answer": "...",
           "points_awarded": 25,
           "max_points": 25,
           "teacher_feedback": "Отличная работа!",
           "is_graded": true
         }
       ]
     }
     ```
   - Если прогресса еще нет — отдавать 200 OK с дефолтным объектом `{ "status": "not_started", "score": null }` вместо 404, чтобы фронтенду было удобнее инициализировать плеер.

2. **Бизнес-логика (Service Layer):**
   - Получить `LessonProgress` из `progressRepo`.
   - Если есть сданные попытки квизов/эссе для этого урока и пользователя — подтянуть их статус и оценку через `quizRepo`.

## ✅ Критерии приёмки (Definition of Done)
- [x] Запрос возвращает статус урока, балл и комментарии преподавателя.
- [x] Для нового урока возвращается объект со статусом `not_started` без падения в ошибку.
- [x] Документация в `FUNCTIONAL_SPEC.md` обновлена.
