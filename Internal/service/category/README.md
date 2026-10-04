# ⚙️ Сервис: `internal/service/category`
> **Путь:** `internal/service/category`  
> **Роль:** Бизнес-логика каталога и управления категориями курсов.  
> **Спецификация:** [`FUNCTIONAL_SPEC.md`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/category/FUNCTIONAL_SPEC.md)

---

## 🎯 Назначение и ответственность
Сервис обеспечивает получение дерева и списка категорий платформы с подсчетом активных курсов, а также создание и редактирование категорий.

---

## 📁 Структура файлов
| Файл | Описание |
|---|---|
| [`service.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/category/service.go) | Реализация методов интерфейса `service.CategoryServices` |
| [`category_test.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/category/category_test.go) | Юнит-тесты сервиса категорий |

---

## 🔗 Зависимости
- **Входящие:** [`interfaces/handlers/category`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/category).
- **Исходящие:** [`repository.CategoryRepository`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository), [`db.QueryExecutor`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/infrastructure/db).
