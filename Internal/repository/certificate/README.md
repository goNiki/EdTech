# 📦 Репозиторий: `internal/repository/certificate`
> **Путь:** `internal/repository/certificate`  
> **Роль:** Работа с таблицей `certificates` в PostgreSQL.

---

## 🎯 Назначение и ответственность
Репозиторий обеспечивает персистентность данных сертификатов: сохранение выданного сертификата (с разрешением коллизий `ON CONFLICT (user_id, course_id)`), чтение по публичному коду сертификата и чтение по связке `(user_id, course_id)`.

---

## 📁 Структура файлов
| Файл | Описание |
|---|---|
| [`repository.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/certificate/repository.go) | Реализация методов интерфейса `repository.CertificateRepository` |

---

## ⚙️ Функции, методы и API
| Метод | Описание |
|---|---|
| `CreateCertificate` | Сохранение сертификата в БД с генерацией ID и фиксацией даты выдачи |
| `GetCertificateByCode` | Чтение сертификата по публичному уникальному коду подтверждения |
| `GetCertificateByUserAndCourse` | Проверка наличия выданного сертификата для студента по курсу |
