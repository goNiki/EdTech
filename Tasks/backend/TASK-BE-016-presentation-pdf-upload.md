# 🛠 [BE-016] Поддержка хранения и валидации PDF-презентаций уроков

> **Приоритет:** Medium (P1)  
> **Связанные задачи:** FE-016, QA-016  
> **Целевой модуль:** `internal/interfaces/handlers/upload/`, `internal/service/upload/`  
> **Документация модуля:** [PROJECT_MAP.md](../../PROJECT_MAP.md)

## 🎯 Цель задачи
Обеспечить поддержку загрузки и быстрой отдачи слайдовых презентаций в формате PDF (экспортированных из PowerPoint, Keynote или Google Slides). Файлы презентаций могут иметь повышенный размер (до 50 МБ) и требуют правильных HTTP-заголовков для встроенного предпросмотра в браузере (inline viewing).

## 🔍 Текущее состояние кода
- В рамках задачи `TASK-BE-002` создается эндпоинт `POST /api/v1/upload`.
- Требуется настроить категорию `presentation` с поддержкой файлов до 50 МБ и гарантировать отдачу заголовка `Content-Disposition: inline` со шлюза статики `/static/*`, чтобы браузер открывал PDF во фрейме/плагине, а не принудительно скачивал его на диск.

## 📝 Технические требования к реализации
1. **API / Endpoint:**
   - `POST /api/v1/upload` (с полем `category = presentation`).
   - Разрешенный MIME-тип: `application/pdf`.
   - Максимальный лимит: 50 МБ для презентаций.
   - Response DTO:
     ```json
     {
       "file_url": "http://localhost:8082/static/uploads/presentation/2026/10/uuid.pdf",
       "file_name": "lecture_01_architecture.pdf",
       "size_bytes": 18450120,
       "mime_type": "application/pdf"
     }
     ```

2. **HTTP Заголовки статического сервера:**
   - Для запросов `GET /static/uploads/presentation/*.pdf` выставлять:
     - `Content-Type: application/pdf`
     - `Content-Disposition: inline; filename="presentation.pdf"`
     - `Accept-Ranges: bytes` (для поддержки потоковой подгрузки страниц без скачивания всего файла целиком).

## ✅ Критерии приёмки (Definition of Done)
- [x] Презентация размером 40 МБ успешно загружается через эндпоинт загрузки.
- [x] Браузер может отображать PDF во встроенном просмотрщике через Range-запросы без принудительного скачивания.
