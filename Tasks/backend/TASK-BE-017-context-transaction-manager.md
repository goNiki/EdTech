# 🛠 [BE-017] Рефакторинг Transaction Manager на паттерн Context-Based TxManager

**Статус:** ✅ Completed

> **Приоритет:** Blocker (P0)  
> **Аудит:** [CODE_REVIEW_GOLANG.md](../../CODE_REVIEW_GOLANG.md) (Дефекты 1.1, 1.2, 2.3)  
> **Целевые модули:**  
> - `internal/infrastructure/txmanager/`  
> - `internal/repository/` (все пакеты репозиториев)  
> - `internal/service/` (все сервисы с транзакциями)  

---

## 🎯 Цель задачи
Устранить критические дефекты взаимоблокировок (Deadlocks), исчерпания пула соединений (Pool Starvation) и утечек соединений при `panic()` за счет перехода на контекстно-зависимый Transaction Manager (`Context-Based TxManager`).  
Полностью очистить сигнатуры методов репозиториев и сервисов от сквозного прокидывания инфраструктурного параметра `q db.QueryExecutor`.

---

## 🔍 Текущее состояние кода
1. **Вложенные транзакции (Nested TX):**  
   В [`internal/service/quiz/submitAttempt.go#L125`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/quiz/submitAttempt.go#L125) и [`internal/service/quiz/gradeAttemptAnswer.go#L105`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/quiz/gradeAttemptAnswer.go#L105) транзакция открывается через `WithTX`. Внутри вызывается `CompleteLesson`, который в [`internal/service/progress/completeLesson.go#L29`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/progress/completeLesson.go#L29) открывает **вторую транзакцию на втором соединении из пула**. Это приводит к исчерпанию пула соединений и дедлокам PostgreSQL (`deadlock detected`).
2. **Отсутствие безопасного отката (Panic / Canceled Context):**  
   В [`internal/infrastructure/txmanager/txmanager.go#L28-L47`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/infrastructure/txmanager/txmanager.go#L28-L47) нет `defer tx.Rollback` и нет защиты от паники. При панике соединение зависает в пуле. При отмене клиентского контекста `tx.Rollback(ctx)` возвращает `context.Canceled`, не завершая откат.
3. **Загрязнение интерфейсов:**  
   Все методы в [`internal/repository/repository.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/repository.go) требуют `q db.QueryExecutor`, вынуждая сервисы хранить `db` и нарушая чистую архитектуру.

---

## 📝 Технические требования к реализации

### 1. Новая реализация `internal/infrastructure/txmanager/txmanager.go`
1. Определить приватный тип ключа контекста:
   ```go
   type txKey struct{}
   ```
2. Интерфейс `TransactionManager`:
   ```go
   type TransactionManager interface {
       WithTX(ctx context.Context, opts pgx.TxOptions, fn func(ctx context.Context) error) error
   }
   ```
3. Логика метода `WithTX`:
   - **Проверка существующей транзакции:** Если в `ctx` уже присутствует активная транзакция (`ctx.Value(txKey{}).(pgx.Tx)`), новая транзакция НЕ создается — функция `fn(ctx)` выполняется в рамках текущей транзакции (Flat nesting).
   - **Корневая транзакция:** Создается через `tm.pool.BeginTx(ctx, opts)`.
   - **Безопасный откат:** Использовать `defer` с `recover()` и откатом через `context.Background()`:
     ```go
     defer func() {
         if p := recover(); p != nil {
             _ = tx.Rollback(context.Background())
             panic(p)
         } else if err != nil {
             _ = tx.Rollback(context.Background())
         }
     }()
     ```
   - Завершение через `tx.Commit(ctx)`.

### 2. Хелпер `GetQueryExecutor`
В пакете `txmanager` создать хелпер:
```go
func GetQueryExecutor(ctx context.Context, defaultPool QueryExecutor) QueryExecutor {
    if tx, ok := ctx.Value(txKey{}).(pgx.Tx); ok {
        return tx
    }
    return defaultPool
}
```

### 3. Очистка репозиториев
1. Удалить аргумент `q db.QueryExecutor` из сигнатур всех методов в `internal/repository/` и интерфейса `repository.go`.
2. В каждом репозитории извлекать исполнитель через `q := txmanager.GetQueryExecutor(ctx, r.pool)`.
3. Сигнатура методов становится идиоматичной:
   ```go
   CreateCourse(ctx context.Context, course *domain.Course) (*domain.Course, error)
   ```

### 4. Очистка сервисов
1. Удалить поле `db db.QueryExecutor` из структур сервисов (`service.db`).
2. В вызовах `WithTX` коллбэк принимает только `ctx`:
   ```go
   err := s.txManager.WithTX(ctx, pgx.TxOptions{}, func(ctx context.Context) error {
       return s.courserepo.CreateCourse(ctx, course)
   })
   ```

---

## ✅ Критерии приёмки (Definition of Done)
- [x] `WithTX` переиспользует существующую транзакцию при вложенных вызовах без запроса нового соединения из пула.
- [x] При возникновении `panic()` внутри `fn` транзакция гарантированно откатывается, а соединение возвращается в пул в чистом виде.
- [x] `Rollback` корректно отрабатывает при отмененном клиентском контексте (`context.Canceled`).
- [x] Ни один метод репозитория не принимает `q db.QueryExecutor`.
- [x] Ни один сервис не хранит поле `db db.QueryExecutor`.
- [x] Написаны unit-тесты на `TxManager` (проверка вложенности, проверки отката при ошибке и панике).
- [x] Все существующие тесты (`go test ./...`) успешно проходят.
