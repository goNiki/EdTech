# 📦 Модуль: `courses`
> **Путь:** `frontend/src/app/courses`  
> **Роль:** Публичный каталог образовательных курсов с динамической фильтрацией, поиском и детальным лендингом программы ([slug]).

---

## 🎯 Назначение и ответственность
Модуль обеспечивает публичную витрину курсов платформы:
1. Каталог (`courses/page.tsx`): фильтрация по уровню сложности (`beginner`, `intermediate`, `advanced`), языку, поиск по названию/описанию с дебаунсингом, сортировка и модальное окно предпросмотра промо-видео.
2. Детальный лендинг курса (`courses/[slug]/page.tsx`): отображение описания курса, авторов, структуры модулей и уроков, проверка статуса зачисления текущего пользователя и кнопка быстрой записи на курс (`Enroll`).

---

## ⚠️ Жестко установленные правила (Hard Invariants)
1. **Поддержка Slug и Числового ID:** Страница курса `[slug]` обязана поддерживать как строковый URL-слаг, так и числовой первичный ключ ID. Проверка `!isNaN(Number(slug))` определяет, вызывать ли эндпоинт `GET /courses/{id}` или `GET /courses/slug/{slug}`.
2. **Гостевой доступ и защита кнопки Enroll:** Просмотр каталога и программы курса открыт неавторизованным пользователям (гостям). При клике на кнопку «Записаться на курс» неавторизованный пользователь принудительно перенаправляется на `/login`.
3. **Дебаунсинг фильтрации:** Изменение поисковой строки или фильтров сложности дебаунсится на 250 мс через `setTimeout` перед отправкой HTTP-запроса к `/courses`.
4. **Сквозная навигация при зачислении:** После успешного вызова `POST /courses/{id}/enroll` пользователь моментально перенаправляется в интерфейс прохождения курса [`/dashboard/courses/${courseData.id}`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/dashboard/courses/%5Bid%5D/page.tsx).

---

## 📁 Структура файлов модуля
| Файл | Описание роли файла |
|---|---|
| [`page.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/courses/page.tsx#L1-L409) | Страница каталога курсов с фильтрами, сортировкой, сеткой курсов и модалкой промо-видео |
| [`[slug]/page.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/courses/%5Bslug%5D/page.tsx#L1-L223) | Лендинг курса: учебный план (силлабус), метаданные, видео-тизер и действие записи |

---

## ⚙️ Функции, компоненты, хуки и API
| Компонент / Метод | Файл:Строки | Описание | Маршрут / Сигнатура |
|---|---|---|---|
| [`CoursesCatalogPage`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/courses/page.tsx#L40-L408) | [`page.tsx#L40-L408`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/courses/page.tsx#L40-L408) | Витрина каталога с сайдбаром и фильтрами | Публичный маршрут `/courses` |
| `fetchCourses` | [`page.tsx#L67-L96`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/courses/page.tsx#L67-L96) | Запрос к `GET /courses` с параметрами фильтрации и проверка `GET /courses/my` | `() => Promise<void>` |
| `handleEnroll (каталог)` | [`page.tsx#L105-L120`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/courses/page.tsx#L105-L120) | Быстрая запись на курс через `POST /courses/{id}/enroll` | `(courseId: number, e: MouseEvent) => Promise<void>` |
| [`CourseLandingPage`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/courses/%5Bslug%5D/page.tsx#L19-L222) | [`[slug]/page.tsx#L19-L222`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/courses/%5Bslug%5D/page.tsx#L19-L222) | Страница курса с отображением дерева модулей | Динамический маршрут `/courses/[slug]` |
| `handleEnroll (лендинг)` | [`[slug]/page.tsx#L69-L81`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/courses/%5Bslug%5D/page.tsx#L69-L81) | Зачисление студента и переход в кабинет изучения | `() => Promise<void>` |

---

## 🔗 Зависимости
- **Входящие (кто вызывает этот модуль):**
  - Пользователи из главного меню [`Sidebar.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/layout/Sidebar.tsx#L42) («Все курсы»).
  - Ссылки из поисковых систем и главной страницы [`app/page.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/page.tsx).
- **Исходящие (что импортирует этот модуль):**
  - [`store/useAuth.ts`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/store/useAuth.ts#L33) — статус авторизации для записи.
  - [`lib/api.ts`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/api.ts#L5) — Axios-клиент (`/courses`, `/courses/my`, `/courses/{id}/enroll`, `/courses/{id}/structure`).
  - [`layout/Sidebar.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/layout/Sidebar.tsx) и [`layout/TopNavbar.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/layout/TopNavbar.tsx).

---

## 🧭 Навигатор типовых задач (Where to edit?)
- **Добавить новые критерии фильтрации (например, по категории):** добавить параметр в `params.append` в [`page.tsx#L70-L76`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/courses/page.tsx#L70-L76).
- **Скорректировать отображение дерева уроков на лендинге:** изменить секцию отображения `structure.map` в [`[slug]/page.tsx#L140-L200`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/courses/%5Bslug%5D/page.tsx#L140-L200).
- **Изменить модальное окно видео-превью:** отредактировать блок плеера в [`page.tsx#L380-L405`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/courses/page.tsx#L380-L405).
