# 📦 Сервис: `internal/service/course`
> **Путь:** `internal/service/course`  
> **Роль:** Управление жизненным циклом курсов, каталог, иерархическая структура и сортировка разделов.

---

## 🎯 Назначение и ответственность
Оркестрирует создание, редактирование, публикацию, архивацию и удаление курсов. Обеспечивает сборку иерархического дерева курса (`CourseStructure`), фильтрацию каталога публичных курсов и переупорядочивание секций.

---

## ⚠️ Жестко установленные правила (Hard Invariants)
1. **Транзакционное создание с привязкой создателя:** При создании нового курса он всегда инициализируется в статусе `draft`. В рамках этой же транзакции создатель курса автоматически зачисляется с ролью `creator` в таблицу `users_courses`.
2. **Каскадная публикация:** Публикация курса (`PublishCourse`) через транзакцию каскадно переводит все дочерние секции и уроки этого курса в статус `published`.
3. **Защита доступа через `AccessService`:** Все операции модификации, публикации, архивации и удаления предварительно проверяют права вызывающего пользователя через методы `accessService.CanEditCourse`, `CanPublishCourse`, `CanDeleteCourse`.
4. **Уникальность Slug:** Идентификатор URL курса (`slug`) обязан быть уникальным. При попытке создать дубликат возвращается `ErrSlugAlreadyExists`.
5. **Фильтрация черновиков в каталоге:** Метод `ListPublicCourses` возвращает только курсы в статусе `published` с видимостью `public`.

---

## 📁 Структура файлов
| Файл | Описание |
|---|---|
| [`service.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/course/service.go) | Определение сервиса и конструктор `NewCourseService` |
| [`create.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/course/create.go) | Создание курса в транзакции с автозачислением создателя (`creator`) |
| [`publish.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/course/publish.go) | Публикация курса с каскадным обновлением секций и уроков |
| [`archive.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/course/archive.go) | Перевод курса в статус `archived` |
| [`get.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/course/get.go) | Получение информации о курсе по ID или Slug с публичным профилем автора и статистикой |
| [`getCourseStructure.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/course/getCourseStructure.go) | Сборка полного иерархического дерева (курс -> секции -> уроки) |
| [`list.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/course/list.go) | Листинг публичных курсов и списка "Мои курсы" |
| [`update.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/course/update.go) | Частичное обновление данных курса (название, описание, обложка) |
| [`update_status.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/course/update_status.go) | Ручная смена статуса курса (`draft`, `published`, `archived`) |
| [`reorderSections.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/course/reorderSections.go) | Пакетное сохранение нового порядка секций внутри курса |
| [`delete.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/course/delete.go) | Удаление курса с проверкой прав `CanDeleteCourse` |

---

## ⚙️ Функции, методы и API
| Функция / Метод | Файл:Строки | Описание | Сигнатура |
|---|---|---|---|
| `CreateCourse` | [`create.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/course/create.go) | Создание черновика курса и зачисление создателя | `CreateCourse(ctx context.Context, course *domain.Course) (*domain.Course, error)` |
| `PublishCourse` | [`publish.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/course/publish.go) | Каскадный перевод курса, секций и уроков в `published` | `PublishCourse(ctx context.Context, userID, courseID int64) error` |
| `GetCourseStructure` | [`getCourseStructure.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/course/getCourseStructure.go) | Сборка полного древовидного представления курса | `GetCourseStructure(ctx context.Context, courseID int64) (domain.CourseStructure, error)` |
| `ReorderSections` | [`reorderSections.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/course/reorderSections.go) | Обновление позиций секций на основе переданного массива ID | `ReorderSections(ctx context.Context, userID, courseID int64, sectionIDs []int64) error` |
| `ListPublicCourses` | [`list.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/course/list.go) | Выборка опубликованных курсов с фильтрацией и пагинацией | `ListPublicCourses(ctx context.Context, p domain.Pagination, f domain.CourseFilter) (domain.PaginatedCourses, error)` |

---

## 🔗 Зависимости
- **Входящие:** [`handlers/courses`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/courses).
- **Исходящие:**
  - [`repository.CourseRepository`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/course)
  - [`repository.SectionRepository`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/section)
  - [`repository.LessonRepository`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/lesson)
  - [`repository.EnrolledRepository`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/enrollment)
  - [`service.AccessService`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/access)
  - [`infrastructure.txmanager.TransactionManager`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/infrastructure/txmanager)

---

## 🧭 Навигатор типовых задач (Where to edit?)
- **Запретить публикацию курсов без уроков или описания:** добавить бизнес-проверку в [`publish.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/course/publish.go).
- **Изменить поля, доступные для редактирования:** обновить валидацию и маппинг в [`update.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/course/update.go).
- **Изменить правила сортировки каталога:** отредактировать передаваемые фильтры в [`list.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/course/list.go).
