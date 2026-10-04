# 🛠 [BE-018] Устранение уязвимости Path Traversal в LocalStorage

**Статус:** ✅ Completed

> **Приоритет:** Blocker (P0 / Security)  
> **Аудит:** [CODE_REVIEW_GOLANG.md](../../CODE_REVIEW_GOLANG.md) (Дефект 1.3)  
> **Целевой модуль:** `internal/infrastructure/storage/local.go`  

---

## 🎯 Цель задачи
Устранить уязвимость Path Traversal (CWE-22, OWASP Top 10 A01) в методах локального файлового хранилища `LocalStorage` (`Save`, `Delete`, `Exists`). Исключить возможность выхода за пределы базовой директории `baseDir` через относительные пути с `..`.

---

## 🔍 Текущее состояние кода
В [`internal/infrastructure/storage/local.go#L25, L49, L57`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/infrastructure/storage/local.go#L25):
```go
fullPath := filepath.Join(s.baseDir, filepath.Clean(relativePath))
```
- Функция `filepath.Clean("../../file")` не удаляет префикс `../`.
- `filepath.Join("uploads", "..\..\file")` разрешается в `"..\file"`, выходя из папки `uploads`.
- Злоумышленник может передать вредоносный `relativePath` и перезаписать, удалить или прочитать произвольные файлы сервера.

---

## 📝 Технические требования к реализации
1. Добавить внутренний метод валидации и резолва пути в `LocalStorage`:
   ```go
   func (s *LocalStorage) resolveSafePath(relativePath string) (string, error) {
       cleanBase, err := filepath.Abs(s.baseDir)
       if err != nil {
           return "", fmt.Errorf("resolve base dir: %w", err)
       }

       targetPath := filepath.Clean(filepath.Join(cleanBase, relativePath))
       rel, err := filepath.Rel(cleanBase, targetPath)
       if err != nil || strings.HasPrefix(rel, "..") || rel == "." {
           return "", fmt.Errorf("security violation: path traversal detected for '%s'", relativePath)
       }

       return targetPath, nil
   }
   ```
2. Использовать `resolveSafePath` во всех методах:
   - `Save(ctx context.Context, relativePath string, src io.Reader) error`
   - `Delete(ctx context.Context, relativePath string) error`
   - `Exists(ctx context.Context, relativePath string) (bool, error)`
3. При обнаружении попытки Path Traversal возвращать типизированную ошибку валидации и логировать `slog.Warn`.

---

## ✅ Критерии приёмки (Definition of Done)
- [x] Пути, содержащие `../`, `..\\`, абсолютные пути (`/etc/`, `C:\`) отклоняются с ошибкой безопасности.
- [x] Корректные относительные пути (например, `avatar/2026/04/file.jpg`) успешно сохраняются, проверяются и удаляются.
- [x] Написаны unit-тесты в `internal/infrastructure/storage/local_test.go` с проверкой на Path Traversal payload'ы.
- [x] `go test ./internal/infrastructure/storage/...` проходит успешно.
