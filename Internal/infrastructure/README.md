# 📦 Модуль: `internal/infrastructure`
> **Путь:** `internal/infrastructure`  
> **Роль:** Техническая инфраструктура: пул соединений PostgreSQL (pgxpool), менеджер транзакций, токены JWT v5, хэширование bcrypt/sha256, логирование slog, валидатор и goose-мигратор на базе `goose.Provider`.

---

## 🎯 Назначение и ответственность
Модуль инкапсулирует низкоуровневые технические детали взаимодействия с базой данных, криптографией, чтением конфигурации среды и системными библиотеками, защищая доменный слой и сервисы от деталей реализации.

---

## 📁 Структура файлов модуля
| Файл | Описание роли файла |
|---|---|
| [`db/postgres.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/infrastructure/db/postgres.go) | Пул соединений PostgreSQL на базе `github.com/jackc/pgx/v5/pgxpool` с пингом и хелсчеками |
| [`db/query_executor.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/infrastructure/db/query_executor.go) | Интерфейс `QueryExecutor` (`Exec`, `Query`, `QueryRow`), объединяющий `pgxpool.Pool` и `pgx.Tx` |
| [`txmanager/txmanager.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/infrastructure/txmanager/txmanager.go) | Управление транзакциями `TxManager.WithTX` с авто-rollback при ошибках и commit |
| [`jwt/jwt.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/infrastructure/jwt/jwt.go) | Генерация и валидация HMAC-SHA256 JWT Access & Refresh токенов (`JwtManager`) с `ParserOption` (`WithValidMethods`, `WithIssuer`, `WithExpirationRequired`) |
| [`jwt/jwt_test.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/infrastructure/jwt/jwt_test.go) | Unit-тесты для парсера JWT (валидация exp, iss, alg) |
| [`hasher/hasher.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/infrastructure/hasher/hasher.go) | Хэширование паролей (`bcrypt`) и генерация хэшей refresh-токенов (`SHA-256`) |
| [`logger/logger.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/infrastructure/logger/logger.go) | Контекстный логгер `slog.JSONHandler` с обогащением операцией `op` и `request_id` |
| [`logger/sl/sl.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/infrastructure/logger/sl/sl.go) | Хелпер `sl.Error(err)` для форматирования ошибки в атрибут `slog.Attr` |
| [`validator/validator.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/infrastructure/validator/validator.go) | Обертка над валидатором структур `go-playground/validator/v10` |
| [`migrator/migrator.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/infrastructure/migrator/migrator.go) | Потокобезопасная обертка над изолированным `goose.NewProvider` (`Up`, `Down`, `Status`, `Create`) |
| [`config/config.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/infrastructure/config/config.go) | Загрузка конфигурации из `.env` файлов и переменных окружения |
| [`config/env/*`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/infrastructure/config/env) | Провайдеры конфигураций: `postgres.go`, `jwt.go`, `server.go`, `logger.go` |

---

## ⚙️ Функции, методы и API
| Функция / Класс | Файл:Строки | Описание | Сигнатура / Вход и Выход |
|---|---|---|---|
| [`db.New`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/infrastructure/db/postgres.go) | [`db/postgres.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/infrastructure/db/postgres.go) | Создание пула pgxpool с таймаутом и проверкой соединения | `func New(cfg config.Postgres) (*Postgres, error)` |
| [`TxManager.WithTX`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/infrastructure/txmanager/txmanager.go) | [`txmanager/txmanager.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/infrastructure/txmanager/txmanager.go) | Обертка транзакционного выполнения блока с автоматическим откатом | `WithTX(ctx context.Context, opts pgx.TxOptions, fn func(...) error) error` |
| [`JwtManager.GenerateAccessToken`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/infrastructure/jwt/jwt.go) | [`jwt/jwt.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/infrastructure/jwt/jwt.go) | Создание JWT токена доступа с claims (`ID`, `Username`, `Role`, `iss: auth`) | `(j *JwtManager) GenerateAccessToken(u *domain.User, now time.Time) (string, error)` |
| [`JwtManager.ParseToken`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/infrastructure/jwt/jwt.go) | [`jwt/jwt.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/infrastructure/jwt/jwt.go) | Безопасный парсинг и верификация токена через опции v5 (`WithValidMethods`, `WithIssuer`, `WithExpirationRequired`) | `(j *JwtManager) ParseToken(tokenStr string) (int64, string, error)` |
| [`migrator.NewMigrator`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/infrastructure/migrator/migrator.go) | [`migrator/migrator.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/infrastructure/migrator/migrator.go) | Инициализация изолированного экземпляра `goose.Provider` без мутации глобального состояния | `func NewMigrator(pool *pgxpool.Pool, migDir string) (*Migrator, error)` |

---

## 🔗 Зависимости
- **Входящие (кто вызывает этот модуль):**
  - [`internal/app/di.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/app/di.go) (собирает инфраструктуру в контейнер)
  - [`internal/service/*`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/README.md) (использует JWT, Hasher, TxManager, Logger)
  - [`internal/repository/*`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/README.md) (использует `QueryExecutor` и `db.Postgres`)
- **Исходящие (что импортирует этот модуль):**
  - Внешние библиотеки: `pgx/v5`, `golang-jwt/jwt/v5`, `pressly/goose/v3`, `golang.org/x/crypto/bcrypt`.
  - [`pkg/errors`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/pkg/README.md).
