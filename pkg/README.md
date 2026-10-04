# 📦 Модуль: `pkg`
> **Путь:** `pkg`  
> **Роль:** Общие платформенные пакеты: глобальный каталог ошибок (`apperrors`), утилиты для работы со слагами (`slug`) и генерации UUID (`uuid`).

---

## 🎯 Назначение и ответственность
Предоставляет типизированные константы и sentinel-ошибки предметной области и инфраструктуры (`pkg/errors`), используемые всеми слоями приложения, а также независимые утилиты (`pkg/slug`, `pkg/uuid`). Пакет `pkg` строго не зависит от `internal/`.

---

## 📁 Структура файлов модуля
| Файл / Пакет | Описание роли |
|---|---|
| [`errors/errorsAPP.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/pkg/errors/errorsAPP.go) | Единый реестр sentinel-ошибок `apperrors` (доменные, конфликты, валидация, права, инфраструктура БД и JWT) |
| [`errors/auth.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/pkg/errors/auth.go) | Специфические ошибки аутентификации (дубликаты email/username, деактивация, блокировка) |
| [`slug/slug.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/pkg/slug/slug.go) | Генерация и нормализация безопасных URL-слагов |
| [`uuid/uuid.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/pkg/uuid/uuid.go) | Генерация криптографически стойких RFC 4122 v4 UUID идентификаторов |
| [`utils/`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/pkg/utils) | Фасад обратной совместимости для `slug` и `uuid` |

---

## ⚙️ Функции, методы и API
| Функция / Символ | Файл | Описание | Сигнатура |
|---|---|---|---|
| `Err*` (Sentinel Errors) | [`errors/errorsAPP.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/pkg/errors/errorsAPP.go) | Константы ошибок (`ErrUserNotFound`, `ErrForbidden`, `ErrCourseNotFound`, `ErrSlugAlreadyExists`...) | `var Err* error` |
| `NormalizeSlug` | [`slug/slug.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/pkg/slug/slug.go) | Преобразование строки в безопасный слаг `[a-z0-9-]` | `func NormalizeSlug(input string) string` |
| `GenerateUUID` | [`uuid/uuid.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/pkg/uuid/uuid.go) | Генерация случайного UUID v4 | `func GenerateUUID() (string, error)` |

---

## 🔗 Зависимости
- **Входящие (кто вызывает этот модуль):**
  - Все слои системы (`internal/domain`, `internal/service`, `internal/repository`, `internal/interfaces/handlers`, `internal/interfaces/response`).
- **Исходящие (что импортирует этот модуль):**
  - `errors`, `crypto/rand`, `fmt`, `regexp`, `strings` (строго только стандартная библиотека, никаких `internal/`).

---

## 🧭 Навигатор типовых задач (Where to edit?)
- **Добавить новую бизнес-ошибку:** объявить новую переменную ошибки в [`errors/errorsAPP.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/pkg/errors/errorsAPP.go).
- **Связать ошибку с правильным HTTP-кодом (400, 403, 404, 409):** зарегистрировать её в `switch` функции [`internal/interfaces/response/handler.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/response/handler.go).
