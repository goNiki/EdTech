# 📦 Репозиторий: `internal/repository/section`
> **Путь:** `internal/repository/section`  
> **Роль:** SQL-запросы к таблице `sections` (модули курса, сортировка, позиционирование).

---

## 🎯 Назначение и ответственность
Обеспечивает сохранение и выборку разделов курса в PostgreSQL с учетом позиционирования, пакетного обновления порядка через `unnest()` и мягкого удаления (`deleted_at IS NULL`).

---

## ⚠️ Жестко установленные правила (Hard Invariants)
1. **Строгая сортировка по позиции:** `ListSectionsByCourseID` всегда выполняет `ORDER BY position ASC`.
2. **Мягкое удаление (Soft Delete):** Запросы проверяют `deleted_at IS NULL`.
3. **Управление позициями:** Метод `GetMaxPositionByCourseID` через `COALESCE(MAX(position), 0)` гарантирует корректный расчет следующей свободной позиции даже для нового курса.
4. **Пакетное переупорядочивание (`ReorderSections`):** Выполняется за 1 SQL-запрос через PostgreSQL `unnest()` без N+1 раунд-трипов.

---

## 📁 Структура файлов
| Файл | Описание |
|---|---|
| [`repository.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/section/repository.go) | Фабрика `NewSectionRepository` |
| [`section.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/section/section.go) | Реализация методов `CreateSection`, `GetSectionByID`, `ListSectionsByCourseID`, `UpdateSection`, `ReorderSections`, `DeleteSection` |
| [`reorder_test.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/section/reorder_test.go) | Unit-тесты пакетного обновления порядка секций |

---

## ⚙️ Функции, методы и API
| Функция / Метод | Файл:Строки | Описание | Сигнатура |
|---|---|---|---|
| `CreateSection` | [`section.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/section/section.go) | `INSERT INTO sections ... RETURNING` | `CreateSection(ctx context.Context, section *domain.Section) (*domain.Section, error)` |
| `ListSectionsByCourseID` | [`section.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/section/section.go) | Выборка всех секций курса с `ORDER BY position ASC` | `ListSectionsByCourseID(ctx context.Context, courseID int64) ([]domain.Section, error)` |
| `ReorderSections` | [`section.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/section/section.go) | Пакетное обновление порядка секций через `unnest()` | `ReorderSections(ctx context.Context, courseID int64, sectionIDs []int64) error` |
| `GetMaxPositionByCourseID` | [`section.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/section/section.go) | Вычисление текущей максимальной позиции в курсе | `GetMaxPositionByCourseID(ctx context.Context, courseID int64) (int, error)` |

---

## 🔗 Зависимости
- **Входящие:** [`service/section`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/section), [`service/course`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/course).
- **Исходящие:** [`txmanager`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/infrastructure/txmanager), [`models/converter`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/models/converter).
