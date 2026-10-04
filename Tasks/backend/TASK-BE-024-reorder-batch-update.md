# 🛠 [BE-024] Оптимизация переупорядочивания (ReorderLessons & ReorderSections): переход с N+1 на пакетный SQL UPDATE

**Статус:** Completed

> **Приоритет:** High (P1 / Performance & Database Optimization)  
> **Аудит:** [CODE_REVIEW_GOLANG.md](../../CODE_REVIEW_GOLANG.md) (Дефект 2.6)  
> **Целевые модули:**  
> - `internal/repository/lesson/reorderLessons.go`  
> - `internal/repository/section/section.go` (ReorderSections)  

---

## 🎯 Цель задачи
Устранить проблему N+1 сетевых запросов к PostgreSQL при изменении порядка уроков и секций курса.  
Заменить цикл одиночных вызовов `q.Exec(...)` в цикле `for` на единый атомарный пакетный запрос с использованием PostgreSQL `unnest()`, сократив нагрузку на базу и время выполнения операции в десятки раз.

---

## 🔍 Текущее состояние кода
1. В [`internal/repository/lesson/reorderLessons.go#L14-L31`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/lesson/reorderLessons.go#L14-L31):
   ```go
   for idx, id := range lessonIDs {
       pos := int64(idx + 1)
       // ...
       query = `UPDATE lessons SET position = $1, section_id = $2, updated_at = NOW() WHERE id = $3 AND deleted_at IS NULL`
       _, err := q.Exec(ctx, query, args...)
   }
   ```
2. В [`internal/repository/section/section.go#L189-L195`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/section/section.go#L189-L195):
   ```go
   for idx, id := range sectionIDs {
       pos := idx + 1
       _, err := q.Exec(ctx, query, pos, id, courseID)
   }
   ```
- При курсе из 50 уроков или 15 секций выполняется 50 и 15 отдельных сетевых раунд-трипов в рамках одной транзакции.
- Это удерживает блокировки строк, нагружает пул соединений и замедляет ответ API редактора курсов.

---

## 📝 Технические требования к реализации

### 1. Пакетное обновление уроков (`ReorderLessons`)
Сформировать массивы идентификаторов и позиций и выполнить единый запрос:
```go
func (r *repository) ReorderLessons(ctx context.Context, sectionID *int64, lessonIDs []int64) error {
    if len(lessonIDs) == 0 {
        return nil
    }

    positions := make([]int32, len(lessonIDs))
    for i := range lessonIDs {
        positions[i] = int32(i + 1)
    }

    var query string
    var args []any

    if sectionID != nil {
        query = `
            UPDATE lessons AS l
            SET position = v.new_pos, section_id = $1, updated_at = NOW()
            FROM (SELECT unnest($2::bigint[]) AS id, unnest($3::int[]) AS new_pos) AS v
            WHERE l.id = v.id AND l.deleted_at IS NULL
        `
        args = []any{*sectionID, lessonIDs, positions}
    } else {
        query = `
            UPDATE lessons AS l
            SET position = v.new_pos, updated_at = NOW()
            FROM (SELECT unnest($1::bigint[]) AS id, unnest($2::int[]) AS new_pos) AS v
            WHERE l.id = v.id AND l.deleted_at IS NULL
        `
        args = []any{lessonIDs, positions}
    }

    q := txmanager.GetQueryExecutor(ctx, r.pool)
    _, err := q.Exec(ctx, query, args...)
    if err != nil {
        return fmt.Errorf("reorder lessons batch: %w", err)
    }
    return nil
}
```

### 2. Пакетное обновление секций (`ReorderSections`)
Аналогично в `internal/repository/section/`:
```sql
UPDATE sections AS s
SET position = v.new_pos, updated_at = NOW()
FROM (SELECT unnest($1::bigint[]) AS id, unnest($2::int[]) AS new_pos) AS v
WHERE s.id = v.id AND s.course_id = $3 AND s.deleted_at IS NULL
```

---

## ✅ Критерии приёмки (Definition of Done)
- [x] Переупорядочивание N элементов выполняется ровно за 1 SQL-запрос вместо N запросов.
- [x] Порядок позиций (1..N) строго соответствует порядку в переданном массиве ID.
- [x] Написаны unit/интеграционные тесты для проверки сохранения нового порядка элементов.
- [x] Время выполнения операции в тестах сокращается более чем на 80%.
