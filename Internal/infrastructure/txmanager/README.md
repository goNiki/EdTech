# 📦 Модуль: `internal/infrastructure/txmanager`
> **Путь:** `internal/infrastructure/txmanager`  
> **Роль:** Управление транзакциями базы данных PostgreSQL с поддержкой контекстного распространения (Context-Based TxManager) и плоской вложенности (Flat Nesting).

---

## 🎯 Назначение и ответственность
Модуль предоставляет интерфейс `TransactionManager` и контекстно-зависимый хелпер `GetQueryExecutor`:
- **Flat Nesting (Плоская вложенность):** При вызове `WithTX` проверяется наличие уже открытой транзакции в переданном `ctx`. Если транзакция существует, функция выполняется в ее рамках без запроса нового соединения из пула, предотвращая исчерпание пула (`Pool Starvation`) и дедлоки (`Deadlocks`).
- **Безопасный откат:** Гарантированный откат через `context.Background()` в `defer` при панике (`recover()`) или возврате ошибки. После отката при панике паника пробрасывается дальше (`panic(p)`).
- **Чистая архитектура:** Устраняет необходимость прокидывания `db.QueryExecutor` через методы репозиториев и хранение `db` в структурах сервисов.

---

## 📁 Структура модуля
| Файл | Описание |
|---|---|
| [`txmanager.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/infrastructure/txmanager/txmanager.go) | Реализация `txManager`, метод `WithTX`, хелпер `GetQueryExecutor` |
| [`txmanager_test.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/infrastructure/txmanager/txmanager_test.go) | Unit-тесты: коммит, откат при ошибке, откат при панике, плоская вложенность |

---

## ⚙️ API и использование

### Вызов транзакции в сервисе:
```go
err := s.txManager.WithTX(ctx, pgx.TxOptions{}, func(ctx context.Context) error {
    if err := s.repoA.OperationA(ctx, data); err != nil {
        return err
    }
    return s.repoB.OperationB(ctx, data)
})
```

### Извлечение исполнителя в репозитории:
```go
func (r *Repository) SomeQuery(ctx context.Context, id int64) error {
    q := txmanager.GetQueryExecutor(ctx, r.Pool)
    _, err := q.Exec(ctx, "UPDATE ...", id)
    return err
}
```
