# 🛠 [BE-026] Актуализация библиотек по Best Practices: JWT v5, pgx v5, Goose v3

**Статус:** Completed

> **Приоритет:** Medium (P2 / Libraries Best Practices & Security)  
> **Аудит:** [CODE_REVIEW_GOLANG.md](../../CODE_REVIEW_GOLANG.md) (Раздел 3)  
> **Целевые модули:**  
> - `internal/infrastructure/jwt/jwt.go`  
> - `internal/infrastructure/migrator/migrator.go`  

---

## 🎯 Цель задачи
Привести использование сторонних библиотек к актуальным стандартам их авторов (согласно рекомендациям из Context7):
1. **`golang-jwt/jwt/v5`**: Использовать `ParserOption` (`jwt.WithValidMethods`, `jwt.WithIssuer`, `jwt.WithExpirationRequired`) вместо устаревшей ручной проверки типа алгоритма.
2. **`pressly/goose/v3`**: Отказаться от мутации глобального состояния пакета (`goose.SetDialect`, `goose.SetTableName`) в пользу потокобезопасного `goose.NewProvider()`.

---

## 🔍 Текущее состояние кода и требования

### 1. Безопасный парсер токенов в `internal/infrastructure/jwt/jwt.go`
- **Проблема:**  
  В [`internal/infrastructure/jwt/jwt.go#L93-L108`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/infrastructure/jwt/jwt.go#L93-L108) метод подписи проверяется вручную через каст `t.Method.(*jwt.SigningMethodHMAC)`. Кроме того, по умолчанию в `jwt/v5` claim `exp` является опциональным, если явно не затребован парсером.
- **Требование:**  
  Использовать опции парсера v5:
  ```go
  token, err := jwt.ParseWithClaims(tokenStr, &ClaimsAccessToken{}, func(t *jwt.Token) (any, error) {
      return []byte(j.Secret), nil
  },
      jwt.WithValidMethods([]string{"HS256"}),
      jwt.WithIssuer("auth"),
      jwt.WithExpirationRequired(),
  )
  ```

### 2. Потокобезопасный Goose Provider в `internal/infrastructure/migrator/migrator.go`
- **Проблема:**  
  В [`internal/infrastructure/migrator/migrator.go#L25-L36`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/infrastructure/migrator/migrator.go#L25-L36) вызываются глобальные функции `goose.SetDialect` и `goose.SetTableName`.
- **Требование:**  
  В `goose/v3` рекомендуется создавать изолированный провайдер:
  ```go
  provider, err := goose.NewProvider(
      goose.DialectPostgres,
      sqlDB,
      os.DirFS(migDir),
      goose.WithTableName(tableName),
  )
  ```
  И вызывать методы `provider.Up(ctx)`, `provider.Down(ctx)`, `provider.Status(ctx)`.

---

## ✅ Критерии приёмки (Definition of Done)
- [x] Парсинг токенов JWT отклоняет токены без claim `exp` и токены с несовпадающим `iss`.
- [x] Отсутствуют ручные проверки алгоритмов подписи в `keyFunc`.
- [x] Мигратор не модифицирует глобальное состояние пакета `goose`.
- [x] Все миграции и тесты авторизации успешно работают.
