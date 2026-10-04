# 📦 Сервис: `internal/service/certificate`
> **Путь:** `internal/service/certificate`  
> **Роль:** Бизнес-логика выдачи, генерации и публичной верификации цифровых сертификатов об окончании курса.

---

## 🎯 Назначение и ответственность
Сервис отвечает за выдачу официальных электронных сертификатов об окончании курсов. Он гарантирует, что сертификат выдается только при достижении 100% прогресса студента по курсу, генерирует уникальный криптографический код подтверждения формата `EDL-YYYY-XXXXXXXX` и предоставляет публичный метод валидации подлинности сертификата.

---

## ⚠️ Жестко установленные правила (Hard Invariants)
1. **100% завершение курса:** Выдача сертификата возможна только если прогресс студента по курсу составляет ровно `100%` и пройдены все уроки (`CompletedLess == TotalLessons`). При невыполнении возвращается `ErrCourseNotCompleted`.
2. **Идемпотентность выдачи:** Повторный запрос сертификата тем же пользователем для того же курса не создает дубликатов, а возвращает ранее выпущенный сертификат.
3. **Формат кода подтверждения:** Уникальный код генерируется криптографически стойким генератором `crypto/rand` в верхнем регистре вида `EDL-<ГОД>-<HEX8>` (например, `EDL-2026-F9A12B8C`).
4. **Фиксация имени и названия:** В момент выдачи сертификата в базу фиксируются актуальные имя/фамилия студента (`StudentName`) и название курса (`CourseTitle`) для неизменности сертификата даже при последующем переименовании курса или профиля.
5. **Публичная доступность верификации:** Метод `VerifyCertificate` не требует авторизации и возвращает метаданные сертификата для проверки работодателями или внешними сервисами.

---

## 📁 Структура файлов
| Файл | Описание |
|---|---|
| [`service.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/certificate/service.go) | Реализация методов сервиса (`GetOrIssueCertificate`, `VerifyCertificate`) и вспомогательных функций форматирования |
| [`service_test.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/certificate/service_test.go) | Модульные тесты на успешную выдачу, отказ при неполном прохождении и валидацию по коду |

---

## ⚙️ Функции, методы и API
| Функция / Метод | Файл:Строки | Описание | Сигнатура |
|---|---|---|---|
| `GetOrIssueCertificate` | [`service.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/certificate/service.go) | Проверяет прогресс и выдает новый или возвращает существующий сертификат | `GetOrIssueCertificate(ctx context.Context, userID, courseID int64) (*domain.Certificate, error)` |
| `VerifyCertificate` | [`service.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/certificate/service.go) | Валидирует сертификат по его уникальному коду | `VerifyCertificate(ctx context.Context, code string) (*domain.Certificate, error)` |

---

## 🔗 Зависимости
- **Входящие:** [`handlers/certificate`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/certificate)
- **Исходящие:**
  - [`repository.CertificateRepository`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/certificate)
  - [`repository.CourseRepository`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/course)
  - [`repository.UserRepository`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/auth)
  - [`repository.ProgressRepository`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/progress)
