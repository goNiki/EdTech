# 📦 Обработчики: `internal/interfaces/handlers/certificate`
> **Путь:** `internal/interfaces/handlers/certificate`  
> **Роль:** HTTP-контроллеры выдачи и публичной верификации цифровых сертификатов об окончании курсов.

---

## 🎯 Назначение и ответственность
Принимает HTTP-запросы студентов на получение/генерацию сертификата об окончании курса (`/api/v1/courses/{courseid}/certificate`) и публичные запросы на проверку подлинности сертификата по коду (`/api/v1/certificates/verify/{code}`).

---

## 📁 Структура файлов и Эндпоинты
| Файл | HTTP Метод и Путь | Авторизация | Описание |
|---|---|:---:|---|
| [`handler.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/certificate/handler.go) | `GET /api/v1/courses/{courseid}/certificate` | JWT (Student) | Проверка прохождения курса и выдача/возврат сертификата |
| [`handler.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/certificate/handler.go) | `GET /api/v1/certificates/verify/{code}` | Публичный (Без токена) | Валидация подлинности сертификата по уникальному коду |

---

## 🔗 Зависимости
- **Входящие:** [`internal/app/di.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/app/di.go)
- **Исходящие:** [`service.CertificateServices`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/certificate), [`internal/dto`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/dto), [`internal/interfaces/response`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/response)
