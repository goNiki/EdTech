# 🛠 [BE-021] Защита от DoS при загрузке файлов через http.MaxBytesReader

**Статус:** ✅ Completed

> **Приоритет:** Blocker (P0 / Security & Resource Exhaustion)  
> **Аудит:** [CODE_REVIEW_GOLANG.md](../../CODE_REVIEW_GOLANG.md) (Дефект 1.6)  
> **Целевой модуль:** `internal/interfaces/handlers/upload/handler.go`  

---

## 🎯 Цель задачи
Предотвратить исчерпание дискового пространства сервера (Disk Exhaustion DoS) при загрузке файлов. Ограничить максимальный размер входящего HTTP-тела запроса до его парсинга в multipart-форму.

---

## 🔍 Текущее состояние кода
В [`internal/interfaces/handlers/upload/handler.go#L42`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/upload/handler.go#L42):
```go
if err := r.ParseMultipartForm(32 << 20); err != nil {
    response.HandleError(w, r, h.log, errorsAPP.ErrFileTooLarge, op)
    return
}
```
- Значение `32 << 20` (32 МБ) в `r.ParseMultipartForm` задает только объем данных, удерживаемых в памяти.
- Всё, что превышает этот лимит, Go автоматически сохраняет во временные файлы на диск сервера (`os.TempDir`).
- Отсутствие ограничения входящего потока позволяет злоумышленнику отправлять гигабайты мусорных данных в теле запроса, переполняя системный диск.

---

## 📝 Технические требования к реализации
1. В `UploadHandler.UploadFile` перед вызовом `ParseMultipartForm` обернуть тело запроса в `http.MaxBytesReader`:
   ```go
   // Максимальный размер: лимит файла (25 МБ) + запас на multipart-заголовки и поля формы (5 МБ)
   const maxUploadBodySize = 30 * 1024 * 1024 // 30 MB
   r.Body = http.MaxBytesReader(w, r.Body, maxUploadBodySize)
   ```
2. Обработать ошибку переполнения размера:
   ```go
   if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB в памяти, остальное отсекается MaxBytesReader
       var maxBytesErr *http.MaxBytesError
       if errors.As(err, &maxBytesErr) {
           response.HandleError(w, r, h.log, errorsAPP.ErrFileTooLarge, op)
           return
       }
       response.HandleError(w, r, h.log, fmt.Errorf("%w: %v", errorsAPP.ErrValidationFailed, err), op)
       return
   }
   ```

---

## ✅ Критерии приёмки (Definition of Done)
- [x] Запросы размером более 30 МБ немедленно отклоняются со статусом 400 Bad Request (`FILE_TOO_LARGE`).
- [x] Сервер не создает временные файлы на диске при превышении лимита размера запроса.
- [x] Легитимные файлы размером до 25 МБ успешно загружаются.
- [x] Написан e2e/unit тест на обработчик загрузки с превышением лимита размера.
