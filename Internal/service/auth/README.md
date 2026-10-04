# 📦 Сервис: `internal/service/auth`
> **Путь:** `internal/service/auth`  
> **Роль:** Аутентификация, управление учетными записями, выпуск и ротация JWT-токенов, администрирование пользователей.

---

## 🎯 Назначение и ответственность
Обеспечивает безопасный жизненный цикл пользователей: регистрацию с валидацией дубликатов, аутентификацию с проверкой bcrypt-хэшей паролей, генерацию Access/Refresh токенов, ротацию сессий, смену паролей, обновление профилей и права системного администратора (блокировка, смена ролей).

---

## ⚠️ Жестко установленные правила (Hard Invariants)
1. **Безопасность паролей:** Пароли в открытом виде никогда не сохраняются и не логируются. Хэширование выполняется строго через `bcrypt` (`DefaultCost`).
2. **Хэширование Refresh-токенов в БД:** Refresh-токены хранятся в PostgreSQL исключительно в виде SHA-256 хэша (`hasherManager.HashRefreshToken`).
3. **Правило единой сессии (Session Rotation):** При каждом успешном входе (`Login`) или ротации (`RefreshToken`) вызывается `DeleteAllByUserID`, отзывая все предыдущие сессии пользователя.
4. **Инвариант `CanLogin`:** Забаненные (`IsBanned == true`) и деактивированные (`IsActive == false`) пользователи не могут авторизоваться (`ErrUserBanned`, `ErrUserDeactivated`).
5. **Защита администратора от самоблокировки:** Администратор не может забанить сам себя (`adminID == targetUserID`) или понизить свою роль (`ErrCannotModifySelf`).
6. **Мгновенный отзыв токенов при бане:** При блокировке пользователя транзакционно выставляется флаг `is_banned = true` и удаляются все связанные refresh-токены.

---

## 📁 Структура файлов
| Файл | Описание |
|---|---|
| [`service.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/auth/service.go) | Определение структуры сервиса и конструктор `NewAuthService` |
| [`register.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/auth/register.go) | Регистрация пользователя с проверкой уникальности email и username |
| [`login.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/auth/login.go) | Проверка пароля, проверка банов, генерация JWT токенов, фиксация `LastLoginAt` |
| [`refresh.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/auth/refresh.go) | Ротация refresh-токена с проверкой срока жизни |
| [`logout.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/auth/logout.go) | Инвалидация refresh-токена в базе данных |
| [`get.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/auth/get.go) | Получение данных текущего пользователя по ID |
| [`update.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/auth/update.go) | Обновление полей профиля и смена пароля с проверкой старого |
| [`verify_email.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/auth/verify_email.go) | Подтверждение адреса электронной почты |
| [`admin.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/auth/admin.go) | Административные действия: список пользователей с фильтрацией, изменение глобальной роли и бан пользователей |

---

## ⚙️ Функции, методы и API
| Функция / Метод | Файл:Строки | Описание | Сигнатура |
|---|---|---|---|
| `Register` | [`register.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/auth/register.go) | Проверка дубликатов email/username, хэширование пароля и сохранение | `Register(ctx context.Context, input domain.RegisterInput) (*domain.User, error)` |
| `Login` | [`login.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/auth/login.go) | Проверка bcrypt хэша, генерация JWT пары и запись токена в БД | `Login(ctx context.Context, email, password string) (domain.AuthTokens, error)` |
| `RefreshToken` | [`refresh.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/auth/refresh.go) | Проверка срока жизни токена и перевыпуск новой пары | `RefreshToken(ctx context.Context, refreshToken string) (domain.AuthTokens, error)` |
| `Logout` | [`logout.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/auth/logout.go) | Удаление SHA-256 хэша refresh-токена из базы | `Logout(ctx context.Context, refreshToken string) error` |
| `ChangePassword` | [`update.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/auth/update.go) | Сверка старого пароля и обновление на новый bcrypt хэш | `ChangePassword(ctx context.Context, userID int64, oldPassword, newPassword string) error` |
| `ChangeUserRole` | [`admin.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/auth/admin.go) | Изменение роли пользователя администратором | `ChangeUserRole(ctx context.Context, adminID, targetUserID int64, newRole domain.Role) error` |
| `SetUserBanned` | [`admin.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/auth/admin.go) | Блокировка пользователя и транзакционный сброс активных токенов | `SetUserBanned(ctx context.Context, adminID, targetUserID int64, isBanned bool) error` |
| `ListUsersForAdmin` | [`admin.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/auth/admin.go) | Получение списка пользователей с фильтрацией по имени/email/роли/бану и пагинацией | `ListUsersForAdmin(ctx context.Context, adminID int64, filter domain.UserFilter, pagination domain.Pagination) ([]domain.User, int64, error)` |

---

## 🔗 Зависимости
- **Входящие:** [`handlers/auth`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/auth).
- **Исходящие:**
  - [`repository.UserRepository`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/auth)
  - [`repository.RefreshRepository`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/refresh)
  - [`infrastructure.jwt.TokenManager`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/infrastructure/jwt)
  - [`infrastructure.hasher.HasherManager`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/infrastructure/hasher)
  - [`infrastructure.txmanager.TransactionManager`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/infrastructure/txmanager)

---

## 🧭 Навигатор типовых задач (Where to edit?)
- **Изменить требования к паролю:** валидаторы в DTO и метод `Register` в [`register.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/auth/register.go).
- **Разрешить множественные сессии с разных устройств:** убрать вызов `DeleteAllByUserID` и сохранять каждый активный токен отдельно в [`login.go#L62-L68`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/auth/login.go#L62-L68).
- **Добавить отправку письма с кодом подтверждения:** расширить вызов в [`verify_email.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/auth/verify_email.go).
