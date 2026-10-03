# 🛠 [BE-002] Разработка сервиса загрузки и хранения файлов (Object Storage / Local Uploads)

**Статус:** ✅ Completed (хэш коммита: `be2a5c5`)

> **Приоритет:** Critical (P0)  
> **Связанные задачи:** FE-002, QA-002  
> **Целевой модуль:** `internal/interfaces/handlers/upload/`, `internal/service/upload/`, `internal/infrastructure/storage/`  
> **Документация модуля:** [PROJECT_MAP.md](../../PROJECT_MAP.md)

## 🎯 Цель задачи
В проекте полностью отсутствует возможность загрузки бинарных файлов (изображений обложек, аватарок, практических решений студентов в формате PDF/ZIP). Блок `FileUploadBlock` в плеере сейчас является фикцией. Задача — реализовать серверный эндпоинт загрузки файлов с валидацией MIME-типов, ограничением размера, сохранением в локальное хранилище (или S3/MinIO) и генерацией публичного постоянного URL.

## 🔍 Текущее состояние кода
- В кодовой базе нет ни одного обработчика `multipart/form-data`.
- Все ссылки на изображения и файлы в БД и фронтенде ожидают готовый URL (`cover_url`, `avatar_url`, `attachment_url`).
- В БД есть миграция `migrators/20250912230121_create_table_resource.sql`, но она не привязана к API.

## 📝 Технические требования к реализации
1. **API / Endpoint:**
   - `POST /api/v1/upload`
   - Content-Type: `multipart/form-data`
   - Поля формы:
     - `file` (Binary) — файл для загрузки
     - `category` (string, optional: `avatar`, `course_cover`, `homework`, `general`)
   - Response DTO (201 Created):
     ```json
     {
       "file_url": "http://localhost:8082/static/uploads/homework/2026/10/uuid.pdf",
       "file_name": "solution_v1.pdf",
       "size_bytes": 1048576,
       "mime_type": "application/pdf"
     }
     ```
   - Ошибочные коды: 400 (превышен размер или запрещенный тип), 401 (не авторизован).

2. **Бизнес-логика (Service Layer & Infrastructure):**
   - Максимальный размер файла: 25 МБ (настраивается в `.env`).
   - Разрешенные MIME-типы:
     - Изображения: `image/jpeg`, `image/png`, `image/webp`, `image/gif`.
     - Документы и архивы: `application/pdf`, `application/zip`, `application/x-zip-compressed`.
   - Проверка «магических байт» (magic bytes) файла через `http.DetectContentType` (не доверять заголовку клиента).
   - Генерация уникального безопасного имени: UUID v4 + оригинальное расширение.
   - Сохранение файлов на диск в директорию `uploads/{category}/{year}/{month}/`.

3. **HTTP Статический сервер:**
   - В Chi Router (`di.go`) зарегистрировать раздачу статики:  
     `r.Handle("/static/*", http.StripPrefix("/static", http.FileServer(http.Dir("./uploads"))))`

4. **Побочные эффекты:**
   - Логирование успешных загрузок с размером и типом файла.

## ✅ Критерии приёмки (Definition of Done)
- [x] Эндпоинт `POST /api/v1/upload` принимает multipart форму и возвращает доступный по HTTP URL файла.
- [x] Загрузка исполняемых файлов (`.exe`, `.sh`, `.php`, `.js`) строго блокируется с ошибкой 400.
- [x] Файлы более 25 МБ отклоняются.
- [x] Загруженные файлы доступны по GET-запросу по возвращенному `file_url`.
