# 📦 Модуль: `internal/service/upload`
> **Путь:** `internal/service/upload`  
> **Роль:** Слой бизнес-логики загрузки и валидации бинарных файлов (изображения, документы, архивы).

---

## 🎯 Назначение и ответственность
Модуль отвечает за валидацию загружаемых файлов (размер до 25 МБ, MIME-типы, блокировка исполняемых файлов), генерацию уникальных безопасных имен (UUID v4) и сохранение в постоянное хранилище (локальный диск / S3) с возвратом публичного постоянного URL.

---

## 📁 Структура файлов модуля
| Файл | Описание роли файла |
|---|---|
| [`service.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/upload/service.go) | Реализация `UploadServices`: проверка magic bytes, ограничений размера, расширений и сохранение |
| [`upload_test.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/upload/upload_test.go) | Unit-тесты сервиса загрузки: успешная загрузка, блокировка .exe, проверка лимитов и пустых файлов |
| [`FUNCTIONAL_SPEC.md`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/upload/FUNCTIONAL_SPEC.md) | Функциональная спецификация и паспорта функций для бизнес-аналитиков |

---

## ⚙️ Функции, методы и API
| Метод | Сигнатура | Описание |
|---|---|---|
| `UploadFile` | `(ctx context.Context, file io.Reader, filename string, size int64, category string) (*domain.FileUploadResult, error)` | Валидация файла, проверка magic bytes, сохранение и генерация URL |
| `UploadImagesBatch` | `(ctx context.Context, files []domain.BatchFileItem, category string) ([]domain.BatchUploadResultItem, error)` | Пакетная валидация картинок Word (до 50 шт, до 50 МБ) и параллельное сохранение через errgroup |

---

## 🔗 Зависимости
- **Входящие (кто вызывает этот модуль):**
  - [`internal/interfaces/handlers/upload`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/upload/README.md)
- **Исходящие (что импортирует этот модуль):**
  - [`internal/infrastructure/storage`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/infrastructure/storage)
  - [`internal/domain`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/domain)
  - [`pkg/errors`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/pkg/errors)
  - [`pkg/utils`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/pkg/utils)
