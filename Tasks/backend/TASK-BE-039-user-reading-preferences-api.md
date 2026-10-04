# 🛠 [BE-039] API пользовательских настроек отображения и чтения (Preferences & Accessibility Sync)

> **Приоритет:** Medium (P1)  
> **Связанные задачи:** FE-044, QA-044  
> **Целевой модуль:** `internal/interfaces/handlers/auth/`, `internal/service/auth/`, `internal/repository/auth/`  
> **Документация модуля:** [internal/service/auth/README.md](../../internal/service/auth/README.md)

---

## 🎯 Цель задачи
Обеспечить сохранение и синхронизацию индивидуальных настроек доступности и комфортного чтения студента (размер шрифта обучающего контента, предпочтительная ширина полотна урока, тема оформления ридера) между его устройствами (ноутбук, смартфон, планшет).

---

## 🔍 Текущее состояние кода
- В профиле пользователя (`users`) хранятся базовые поля: имя, фамилия, аватар, био, headline.
- Настройки интерфейса (например, размер шрифта) на клиенте хранятся только в `localStorage` конкретного браузера и сбрасываются при смене устройства или авторизации в режиме инкогнито.

---

## 📝 Технические требования к реализации

### 1. Схема данных (PostgreSQL)
Добавить в таблицу `users` колонку `preferences`:
```sql
ALTER TABLE users ADD COLUMN IF NOT EXISTS preferences JSONB DEFAULT '{"font_scale": "medium", "content_width": "wide", "line_height": "normal", "reading_theme": "system"}';
```

### 2. HTTP Endpoint
`PATCH /api/v1/auth/preferences`
* **Авторизация:** Bearer JWT (любой авторизованный пользователь).
* **Request DTO:**
  ```json
  {
    "font_scale": "large",
    "content_width": "wide",
    "line_height": "relaxed",
    "reading_theme": "sepia"
  }
  ```
  * `font_scale`: `"compact"` (14px), `"medium"` (16px), `"large"` (18px), `"xlarge"` (20px).
  * `content_width`: `"standard"` (1024px), `"wide"` (1280px), `"full"` (100%).
  * `line_height`: `"normal"`, `"relaxed"`.
  * `reading_theme`: `"system"`, `"light"`, `"dark"`, `"sepia"`.

### 3. Response DTO
```json
{
  "code": 200,
  "message": "preferences updated successfully",
  "data": {
    "preferences": {
      "font_scale": "large",
      "content_width": "wide",
      "line_height": "relaxed",
      "reading_theme": "sepia"
    }
  }
}
```

### 4. Включение в `GET /api/v1/auth/me`
Поле `preferences` должно возвращаться в стандартном объекте пользователя при авторизации и запросе текущего профиля (`UserResponse`), чтобы фронтенд сразу восстанавливал предпочтения при логине.

---

## ✅ Критерии приёмки (Definition of Done)
- [x] Поле `preferences` сохраняется в БД без потери предыдущих ключей (JSONB shallow merge).
- [x] Невалидные значения валидируются и отклоняются со статусом 400 Bad Request.
- [x] При вызове `GET /auth/me` актуальные настройки возвращаются клиенту.
