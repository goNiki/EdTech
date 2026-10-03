# 📦 Модуль: `courses`
> **Путь:** `frontend/src/app/courses`  
> **Роль:** Публичный каталог образовательных курсов с динамической фильтрацией по категориям, поиском и детальным лендингом программы ([slug]).

---

## 🎯 Назначение и ответственность
Модуль обеспечивает публичную витрину курсов платформы:
1. Каталог (`courses/page.tsx`): фильтрация по категориям (горизонтальный скролл чипсов), уровню сложности (`beginner`, `intermediate`, `advanced`), языку, поиск по названию/описанию с дебаунсингом, сортировка и модальное окно предпросмотра промо-видео.
2. Детальный лендинг курса (`courses/[slug]/page.tsx`): отображение описания курса, авторов, структуры модулей и уроков, проверка статуса зачисления текущего пользователя и кнопка быстрой записи на курс (`Enroll`).

---

## ⚠️ Жестко установленные правила (Hard Invariants)
1. **Поддержка Slug и Числового ID:** Страница курса `[slug]` обязана поддерживать как строковый URL-слаг, так и числовой первичный ключ ID. Проверка `!isNaN(Number(slug))` определяет, вызывать ли эндпоинт `GET /courses/{id}` или `GET /courses/slug/{slug}`.
2. **Гостевой доступ и защита кнопки Enroll:** Просмотр каталога и программы курса открыт неавторизованным пользователям (гостям). При клике на кнопку «Записаться на курс» неавторизованный пользователь принудительно перенаправляется на `/login`.
3. **Дебаунсинг фильтрации:** Изменение поисковой строки или фильтров сложности дебаунсится на 250 мс через `setTimeout` перед отправкой HTTP-запроса к `/courses`.
4. **Сквозная навигация при зачислении:** После успешного вызова `POST /courses/{id}/enroll` пользователь моментально перенаправляется в интерфейс прохождения курса [`/dashboard/courses/${courseData.id}`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/dashboard/courses/%5Bid%5D/page.tsx).
5. **Синхронизация фильтра категорий с URL:** При выборе категории чипс подсвечивается индиго-акцентом, параметр `category_id` заносится в URL без перезагрузки (`window.history.replaceState`), а карточки курсов снабжаются бейджами категории. Загрузка категорий снабжена отказоустойчивым fallback на `DEFAULT_CATEGORIES`.

---

## 📁 Структура файлов модуля
| Файл | Описание роли файла |
|---|---|
| [`page.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/courses/page.tsx) | Страница каталога курсов с горизонтальными чипсами категорий, поиском, сортировкой, сеткой курсов и модалкой промо-видео |
| [`[slug]/page.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/courses/%5Bslug%5D/page.tsx) | Лендинг курса: учебный план (силлабус), метаданные, видео-тизер и действие записи |

---

## ⚙️ Функции, компоненты, хуки и API
| Компонент / Метод | Файл:Строки | Описание | Маршрут / Сигнатура |
|---|---|---|---|
| [`CoursesCatalogPage`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/courses/page.tsx) | [`page.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/courses/page.tsx) | Витрина каталога с сайдбаром, горизонтальными чипсами категорий и фильтрами | Публичный маршрут `/courses` |
| `fetchCourses` | [`page.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/courses/page.tsx) | Запрос к `GET /courses` с `category_id`, `search`, `difficulty` и проверка `GET /courses/my` | `() => Promise<void>` |
| `handleSelectCategory` | [`page.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/courses/page.tsx) | Переключение активной категории и синхронизация с URL query `category_id` | `(categoryId: number \| null) => void` |
| `handleEnroll (каталог)` | [`page.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/courses/page.tsx) | Быстрая запись на курс через `POST /courses/{id}/enroll` | `(courseId: number, e: MouseEvent) => Promise<void>` |
| [`CourseLandingPage`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/courses/%5Bslug%5D/page.tsx) | [`[slug]/page.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/courses/%5Bslug%5D/page.tsx) | Страница курса с отображением дерева модулей | Динамический маршрут `/courses/[slug]` |

---

## 🔗 Зависимости
- **Входящие (кто вызывает этот модуль):**
  - Пользователи из главного меню [`Sidebar.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/layout/Sidebar.tsx) («Все курсы»).
  - Ссылки из поисковых систем и главной страницы [`app/page.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/page.tsx).
- **Исходящие (что импортирует этот модуль):**
  - [`lib/categories.ts`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/categories.ts) — справочник категорий, API `fetchCategories`, бейджи.
  - [`store/useAuth.ts`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/store/useAuth.ts) — статус авторизации для записи.
  - [`lib/api.ts`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/api.ts) — Axios-клиент (`/courses`, `/courses/my`, `/courses/{id}/enroll`, `/courses/{id}/structure`).
  - [`layout/Sidebar.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/layout/Sidebar.tsx) и [`layout/TopNavbar.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/layout/TopNavbar.tsx).

---

## 🧭 Навигатор типовых задач (Where to edit?)
- **Скорректировать стилизацию чипсов категорий:** редактировать секцию чипсов в [`page.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/courses/page.tsx).
- **Скорректировать отображение дерева уроков на лендинге:** изменить секцию отображения `structure.map` в [`[slug]/page.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/courses/%5Bslug%5D/page.tsx).
- **Изменить модальное окно видео-превью:** отредактировать блок плеера в [`page.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/courses/page.tsx).
