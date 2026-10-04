# 📦 Обработчики: `internal/interfaces/handlers/auth`
> **Путь:** `internal/interfaces/handlers/auth`  
> **Роль:** HTTP-контроллеры регистрации, авторизации, профиля, смены пароля и администрирования пользователей.

---

## 🎯 Назначение и ответственность
Принимает HTTP-запросы на эндпоинты `/api/v1/auth/*` и `/api/v1/admin/*`, валидирует входящие JSON-тела через `validator.Validate`, извлекает идентификаторы пользователей из JWT-контекста и передает управление в `service.AuthService`.

---

## ⚠️ Жестко установленные правила (Hard Invariants)
1. **Проверка авторизации (`GetUserID`):** Защищенные маршруты извлекают `userID` из контекста. Если `userID == 0`, немедленно возвращается `401 Unauthorized`.
2. **Административные права:** Эндпоинты `/api/v1/admin/*` строго проверяют, что `userRole == string(domain.RoleAdmin)`.
3. **Единая обработка ошибок:** Все ошибки передаются в `response.HandleError(w, r, log, err, op)`, гарантируя сокрытие внутренних деталей БД и корректный маппинг в HTTP статус-коды.

---

## 📁 Структура файлов и Эндпоинты
| Файл | HTTP Метод и Путь | Описание |
|---|---|---|
| [`register.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/auth/register.go) | `POST /api/v1/auth/register` | Регистрация нового аккаунта |
| [`login.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/auth/login.go) | `POST /api/v1/auth/login` | Вход и выдача токенов |
| [`refresh.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/auth/refresh.go) | `POST /api/v1/auth/refresh` | Обновление пары JWT токенов |
| [`logout.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/auth/logout.go) | `POST /api/v1/auth/logout` | Инвалидация refresh-токена |
| [`me.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/auth/me.go) | `GET /api/v1/auth/me` | Данные текущего авторизованного пользователя |
| [`profile.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/auth/profile.go) | `PATCH /api/v1/auth/profile` | Частичное обновление данных профиля |
| [`password.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/auth/password.go) | `POST /api/v1/auth/change-password` | Смена пароля учетной записи |
| [`verify_email.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/auth/verify_email.go) | `POST /api/v1/auth/verify-email` | Подтверждение email |
| [`admin.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/auth/admin.go) | `GET /api/v1/admin/users`, `PATCH /api/v1/admin/users/{id}/*` | Получение списка пользователей с фильтрами/пагинацией, изменение роли и бан |

---

## 🔗 Зависимости
- **Входящие:** [`internal/app/di.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/app/di.go) (регистрация в маршрутизаторе).
- **Исходящие:** [`service.AuthService`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/auth), [`internal/dto`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/dto), [`internal/interfaces/response`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/response).
