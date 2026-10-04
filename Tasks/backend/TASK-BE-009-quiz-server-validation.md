# 🛠 [BE-009] Серверная валидация квизов и фильтрация правильных ответов (Античит)

> **Приоритет:** Medium (P1)  
> **Статус:** Completed  
> **Связанные задачи:** FE-009, QA-009  
> **Целевой модуль:** `internal/interfaces/handlers/lesson/`, `internal/service/progress/`, `internal/service/lesson/`  
> **Документация модуля:** [internal/service/quiz/FUNCTIONAL_SPEC.md](../../internal/service/quiz/FUNCTIONAL_SPEC.md)

## 🎯 Цель задачи
Устранить уязвимость к списыванию: сейчас при запросе урока `GET /api/v1/lessons/{id}` поле `content` содержит сырой JSON Puck Editor со всеми правильными ответами и дистракторами (`isCorrect: true`, `correctValue`, соответствия). Студент может открыть консоль браузера или сетевую вкладку и увидеть все ответы. Кроме того, студент шлет на сервер сырое число `score`, которое принимается на веру. Задача — санитизировать JSON урока перед отдачей студенту или реализовать серверный подсчет баллов.

## 🔍 Текущее состояние кода
- В `internal/interfaces/handlers/progress/completeLesson.go` принимается `req.Score *int` напрямую от клиента.
- В `internal/service/lesson/getLesson.go` отдается исходный `lesson.Content` как строка JSON.

## 📝 Технические требования к реализации
1. **API / Handlers:**
   - В `GET /api/v1/lessons/{id}`:
     - Если запрашивающий пользователь — `student` (не автор курса и не админ), пропускать `lesson.Content` через санитизатор `SanitizeLessonContentForStudent(contentJSON)`.
     - Удалять ключи `isCorrect`, `correct`, `correctPair`, заменять `{Правильный; Дистракторы}` на обычные пропуски без указания, какой именно вариант правильный.
   - В `POST /api/v1/lessons/{id}/complete`:
     - Принимать в запросе массив выбранных студентом ответов:
       ```json
       {
         "answers": [
           { "block_id": "test-1", "selected_option": 2 },
           { "block_id": "match-2", "pairs": { "1": "B", "2": "A" } }
         ],
         "essays": [...]
       }
       ```
     - На стороне сервера сопоставлять ответы с эталонным `lesson.Content`, вычислять процент выполнения и сохранять подтвержденный балл.

2. **Бизнес-логика (Service Layer):**
   - Пакет санитизации и проверки JSON квизов `internal/service/quiz/sanitizer.go`.
   - Защита от подделки: максимальный балл за урок рассчитывается строго на бэкенде.

## ✅ Критерии приёмки (Definition of Done)
- [x] Студент не может найти правильные ответы в ответе `GET /lessons/{id}` через вкладку Network.
- [x] Оценка рассчитывается бэкендом на основе переданных выборов, а не принимается слепо числом из браузера.
