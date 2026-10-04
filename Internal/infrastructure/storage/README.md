# 📦 Модуль: `internal/infrastructure/storage`
> **Путь:** `internal/infrastructure/storage`  
> **Роль:** Файловое хранилище (LocalStorage) с защитой от атак Path Traversal.

---

## 🎯 Назначение и ответственность
Предоставляет интерфейс `Storage` для сохранения, удаления и проверки наличия файлов на локальном диске (`baseDir`).
- **Защита от Path Traversal (CWE-22):** Метод `resolveSafePath` проверяет относительные пути, отклоняет попытки использования `..`, абсолютные пути, дисковые буквы Windows и пути, выходящие за пределы `baseDir`.
- **Логирование угроз:** Попытки Path Traversal логируются через `slog.Warn`.

---

## 📁 Структура файлов модуля
| Файл | Описание |
|---|---|
| [`storage.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/infrastructure/storage/storage.go) | Определение интерфейса `Storage` |
| [`local.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/infrastructure/storage/local.go) | Реализация `LocalStorage` и метода `resolveSafePath` |
| [`local_test.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/infrastructure/storage/local_test.go) | Unit-тесты на корректную работу и устойчивость к Path Traversal payload'ам |
