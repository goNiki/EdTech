# 🛠 [BE-019] Исправление недопустимой конфигурации CORS и валидация Origin

**Статус:** ✅ Completed

> **Приоритет:** Blocker (P0 / Security & Production Ready)  
> **Аудит:** [CODE_REVIEW_GOLANG.md](../../CODE_REVIEW_GOLANG.md) (Дефект 1.4)  
> **Целевой модуль:** `internal/app/di.go`, `internal/infrastructure/config/`  

---

## 🎯 Цель задачи
Устранить конфликт CORS-заголовков (`Access-Control-Allow-Origin: *` совместно с `Access-Control-Allow-Credentials: true`), который приводит к блокировке авторизованных запросов браузерами (W3C CORS Spec). Внедрить корректную валидацию источника запроса (`Origin`) по белому списку разрешенных доменов из конфигурации.

---

## 🔍 Текущее состояние кода
В [`internal/app/di.go#L442-L458`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/app/di.go#L442-L458):
```go
w.Header().Set("Access-Control-Allow-Origin", "*")
// ...
w.Header().Set("Access-Control-Allow-Credentials", "true")
```
- Согласно спецификации fetch/CORS, браузеры отвергают ответы, где `Allow-Credentials: true` совмещен с wildcard `Allow-Origin: *`.
- Авторизованный фронтенд в production не может совершать запросы к API.

---

## 📝 Технические требования к реализации
1. В конфигурацию сервера (`internal/infrastructure/config/env/server.go`) добавить список разрешенных Origins:
   - Переменная окружения: `CORS_ALLOWED_ORIGINS` (например, `http://localhost:3000,http://localhost:3001,https://edtech.example.com`).
2. Заменить самописный middleware в `di.go` на проверенный `github.com/go-chi/cors`:
   ```go
   r.Use(cors.Handler(cors.Options{
       AllowedOrigins:   cfg.CorsAllowedOrigins(), // []string из конфига
       AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
       AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
       ExposedHeaders:   []string{"Link"},
       AllowCredentials: true,
       MaxAge:           300,
   }))
   ```
   *Примечание:* Библиотека `github.com/go-chi/cors` уже присутствует в `go.mod` (как непрямая зависимость).

---

## ✅ Критерии приёмки (Definition of Done)
- [x] Заголовок `Access-Control-Allow-Origin` возвращает конкретный домен клиента из белого списка, а не wildcard `*`.
- [x] Браузерные preflight-запросы `OPTIONS` с `Authorization` успешно проходят проверку CORS.
- [x] Неразрешенные домены не получают заголовок `Access-Control-Allow-Origin`.
- [x] Список разрешенных Origins конфигурируется через `.env`.
