# 🛠 [BE-035] Архитектурное исправление вложенных транзакций и дедлоков (TxContext & Concurrency Safety)

> **Приоритет:** Critical (P0)  
> **Связанные задачи:** QA-037  
> **Целевой модуль:** `internal/infrastructure/txmanager/`, `internal/service/quiz/`, `internal/service/progress/`  
> **Документация модуля:** [internal/infrastructure/README.md](../../internal/infrastructure/README.md), [CODE_REVIEW_GOLANG.md](../../CODE_REVIEW_GOLANG.md)

---

## 🎯 Цель задачи
Устранить критический архитектурный дефект, выявленный в аудите кодовой базы ([`CODE_REVIEW_GOLANG.md`](../../CODE_REVIEW_GOLANG.md#11-deadlock-и-поломка-атомарности-через-вложенные-транзакции-nested-transactions-in-txmanager)):
В методах `SubmitAttempt` ([`internal/service/quiz/submitAttempt.go#L125`](../../internal/service/quiz/submitAttempt.go#L125)) и `GradeAttemptAnswer` ([`internal/service/quiz/gradeAttemptAnswer.go#L105`](../../internal/service/quiz/gradeAttemptAnswer.go#L105)) запускается транзакция через `s.txManager.WithTX(...)`. 
Внутри этой транзакции вызывается `s.progressService.CompleteLesson(...)`, который внутри себя открывает **вторую отдельную транзакцию** на новом соединении из `pgxpool.Pool`.

**Последствия текущей проблемы:**
1. **Pool Starvation & Deadlock:** При 10–20 одновременных запросах пулу не хватает соединений, и воркеры зависают в дедлоке ожидания соединения.
2. **Postgres Row Deadlock:** Внешняя транзакция держит лок на строке `quiz_attempts`, а внутренняя пытается обновить `lesson_progress` того же студента на другом соединении.
3. **Нарушение атомарности (ACID):** Если внешняя транзакция падает после вызова `CompleteLesson`, прогресс урока уже зафиксирован в БД, а сабмит квиза откатился.

---

## 🔍 Текущее состояние кода
- [`internal/infrastructure/txmanager/txmanager.go`](../../internal/infrastructure/txmanager/txmanager.go): `TxManager.WithTX` не сохраняет открытую транзакцию в `context.Context` и не поддерживает реентерабельность (вложенные вызовы). При панике отсутствует `defer tx.Rollback(context.Background())`.

---

## 📝 Технические требования к реализации

### 1. Поддержка контекстной транзакции в `TxManager`
1. Реализовать функции работы с контекстом:
   ```go
   type txKey struct{}

   func WithTxContext(ctx context.Context, tx pgx.Tx) context.Context {
       return context.WithValue(ctx, txKey{}, tx)
   }

   func GetTxFromContext(ctx context.Context) (pgx.Tx, bool) {
       tx, ok := ctx.Value(txKey{}).(pgx.Tx)
       return tx, ok
   }
   ```
2. Модифицировать `TxManager.WithTX`:
   - Если в `ctx` уже присутствует активная транзакция (`GetTxFromContext(ctx)`), **не открывать новую транзакцию из пула**, а выполнить функцию `fn` в рамках существующей транзакции (reentrant no-op / savepoint).
   - Если транзакции нет — открыть `pool.BeginTx`, поместить `tx` в контекст через `WithTxContext`, запустить `fn`.
   - Добавить безопасный откат при панике:
     ```go
     defer func() {
         if p := recover(); p != nil {
             _ = tx.Rollback(context.Background())
             panic(p) // re-throw panic after rollback
         }
     }()
     ```

### 2. Рефакторинг сервисов `QuizService` и `ProgressService`
- В `CompleteLesson` использовать переданный `ctx`: если транзакция уже открыта во внешнем методе (`SubmitAttempt` или `GradeAttemptAnswer`), `CompleteLesson` не захватывает второе соединение из пула, а выполняется в рамках единой неделимой транзакции.
- При ошибке в `SubmitAttempt` откатываются **все** изменения (и попытка квиза, и прогресс урока).

---

## ✅ Критерии приёмки (Definition of Done)
- [x] Сервис `SubmitAttempt` и `CompleteLesson` гарантируют атомарность: либо сохраняются оба изменения, либо ни одно.
- [x] Нагрузочный тест на 20 одновременных сабмитов в пуле из 5 соединений проходит без зависаний и без ошибки `deadlock detected`.
- [x] При панике внутри транзакции соединение корректно откатывается и не зависает в Postgres.
