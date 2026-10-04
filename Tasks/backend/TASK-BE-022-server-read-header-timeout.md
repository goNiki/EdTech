# 🛠 [BE-022] Защита от Slowloris DoS через установку ReadHeaderTimeout в HTTP-сервере

**Статус:** ✅ Completed

> **Приоритет:** Blocker (P0 / Security & Reliability)  
> **Аудит:** [CODE_REVIEW_GOLANG.md](../../CODE_REVIEW_GOLANG.md) (Дефект 1.7)  
> **Целевой модуль:** `internal/app/app.go`  

---

## 🎯 Цель задачи
Устранить уязвимость к атаке **Slowloris** (CWE-400), настроив обязательный таймаут чтения заголовков (`ReadHeaderTimeout`) для `http.Server`.

---

## 🔍 Текущее состояние кода
В [`internal/app/app.go#L69-L75`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/app/app.go#L69-L75):
```go
a.httpServer = &http.Server{
    Addr:         ":" + cfg.Port(),
    Handler:      a.di.Router(),
    ReadTimeout:  cfg.TimeOut(),
    WriteTimeout: 15 * time.Second,
    IdleTimeout:  cfg.Idletimeout(),
}
```
- Параметр `ReadHeaderTimeout` не установлен.
- Согласно официальной документации `net/http` и рекомендациям безопасности Go, `ReadTimeout` не защищает от медленной передачи HTTP-заголовков. Злоумышленник может слать по 1 байту заголовка раз в несколько секунд, удерживая соединение открытым и блокируя горутины сервера до исчерпания пула сокетов ОС.

---

## 📝 Технические требования к реализации
1. В метод `a.initHttpServer()` структуры `App` добавить явную установку `ReadHeaderTimeout`:
   ```go
   a.httpServer = &http.Server{
       Addr:              ":" + cfg.Port(),
       Handler:           a.di.Router(),
       ReadHeaderTimeout: 5 * time.Second, // Защита от медленных клиентов и Slowloris
       ReadTimeout:       cfg.TimeOut(),
       WriteTimeout:      15 * time.Second,
       IdleTimeout:       cfg.Idletimeout(),
   }
   ```
2. По желанию вынести дефолтное значение (5s) в конфигурацию сервера (`internal/infrastructure/config/env/server.go`).

---

## ✅ Критерии приёмки (Definition of Done)
- [x] Сервер имеет явно заданный `ReadHeaderTimeout <= 5s`.
- [x] Клиентские соединения, не успевшие передать заголовки за 5 секунд, разрываются сервером.
- [x] Сервис корректно компилируется и стартует.
