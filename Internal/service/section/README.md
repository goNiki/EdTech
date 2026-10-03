# 📦 Сервис: `internal/service/section`
> **Путь:** `internal/service/section`  
> **Роль:** Управление тематическими модулями/секциями курса и сортировкой уроков внутри модуля.

---

## 🎯 Назначение и ответственность
Сервис управляет разделами (секциями) курсов: создание, редактирование заголовков и описаний, изменение статуса видимости, удаление секций и переупорядочивание уроков внутри секции (`ReorderLessons`).

---

## ⚠️ Жестко установленные правила (Hard Invariants)
1. **Автоматическое назначение позиции:** Если позиция секции не задана (`section.Position == 0`), она вычисляется как `max(position) + 1` по данному курсу через `GetMaxPositionByCourseID`.
2. **Статус по умолчанию:** Новая секция всегда создается в статусе `draft`, если статус не передан явно.
3. **Сортировка уроков в секции:** Метод `ReorderLessons` принимает упорядоченный срез идентификаторов уроков `lessonIDs` и гарантирует атомарное обновление их позиций.

---

## 📁 Структура файлов
| Файл | Описание |
|---|---|
| [`service.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/section/service.go) | Определение сервиса и конструктор `NewSectionService` |
| [`createSection.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/section/createSection.go) | Создание секции с автовычислением позиции |
| [`updateSection.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/section/updateSection.go) | Редактирование заголовка и описания секции |
| [`updateStatus.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/section/updateStatus.go) | Смена статуса (`draft`/`published`/`archived`) |
| [`reorderLessons.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/section/reorderLessons.go) | Переупорядочивание уроков внутри секции |
| [`deleteSection.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/section/deleteSection.go) | Удаление секции по ID |

---

## ⚙️ Функции, методы и API
| Функция / Метод | Файл:Строки | Описание | Сигнатура |
|---|---|---|---|
| `CreateSection` | [`createSection.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/section/createSection.go) | Создание модуля курса с вычислением позиции | `CreateSection(ctx context.Context, section *domain.Section) (*domain.Section, error)` |
| `UpdateSection` | [`updateSection.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/section/updateSection.go) | Обновление информации о секции | `UpdateSection(ctx context.Context, section *domain.Section) error` |
| `ReorderLessons` | [`reorderLessons.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/section/reorderLessons.go) | Пакетное обновление порядка уроков в секции | `ReorderLessons(ctx context.Context, sectionID int64, lessonIDs []int64) error` |
| `DeleteSection` | [`deleteSection.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/section/deleteSection.go) | Удаление секции | `DeleteSection(ctx context.Context, sectionID int64) error` |

---

## 🔗 Зависимости
- **Входящие:** [`handlers/section`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/section), [`handlers/courses`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/courses).
- **Исходящие:** [`repository.SectionRepository`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/section), [`repository.LessonRepository`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/lesson), [`infrastructure.txmanager.TransactionManager`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/infrastructure/txmanager).

---

## 🧭 Навигатор типовых задач (Where to edit?)
- **Запретить удаление секций, содержащих уроки:** добавить проверку количества уроков в [`deleteSection.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/section/deleteSection.go).
- **Изменить поведение drag-and-drop сортировки уроков:** [`reorderLessons.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/section/reorderLessons.go).
