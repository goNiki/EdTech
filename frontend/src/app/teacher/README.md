# 📦 Модуль: `teacher`
> **Путь:** `frontend/src/app/teacher`  
> **Роль:** Авторская студия преподавателя: реестр курсов, мастер создания, визуальный конструктор силлабуса (модули и уроки), очередь проверки домашних заданий и WYSIWYG-редактор уроков на Puck.

---

## 🎯 Назначение и ответственность
Модуль представляет собой полноценную CMS-среду для преподавателей, авторов курсов и администраторов:
1. Лейаут (`layout.tsx`): RBAC-защита (`teacher`, `author`, `admin`) и автоматическое переключение в полноэкранный режим без сайдбара для редактора Puck.
2. Реестр и создание курсов (`courses/page.tsx`, `courses/new/page.tsx`): создание черновика с автогенерацией slug и переход в конструктор программы.
3. Конструктор силлабуса (`courses/[id]/curriculum/page.tsx`): управление модулями и уроками, переупорядочивание Drag-and-Drop / стрелками, каскадная публикация курса, аналитика успеваемости студентов и ручное зачисление.
4. Настройки курса (`courses/[id]/settings/page.tsx`): редактирование метаданных, промо-видео и параметров видимости.
5. Глобальная очередь проверки (`grading/page.tsx`): проверка практических работ и выставление баллов за открытые вопросы (эссе).
6. Визуальный редактор уроков (`lessons/[id]/edit/page.tsx`): конструирование интерактивного урока на базе Puck с плагином структуры `CustomOutline`.

---

## ⚠️ Жестко установленные правила (Hard Invariants)
1. **Строгая защита ролей преподавателя:** Доступ ко всем вложенным маршрутам разрешен строго для ролей `['teacher', 'author', 'admin']`.
2. **Адаптивный полноэкранный режим редактора уроков:** В `layout.tsx` реализовано правило: если `pathname.includes('/edit')`, боковой `Sidebar` не рендерится, отдавая всю площадь окна холсту Puck Editor.
3. **Формат сохранения Content-as-Data:** Редактор уроков передает данные в `PATCH /lessons/{id}` строго в виде сериализованного JSON-объекта в поле `content`. Данные уроков не содержат произвольного невалидированного HTML.
4. **Контракт пакетного переупорядочивания (Reordering):** Изменение позиций модулей (`PUT /courses/{id}/reorder-sections`) и уроков (`PUT /sections/{id}/reorder-lessons`) отправляет массив ID в строгом порядке `item_ids: number[]`.
5. **Каскадная публикация:** Публикация курса через `POST /courses/{id}/publish` переводит в статус `published` весь силлабус. Возврат в черновик выполняется через `PATCH /courses/{id}/status` со значением `'draft'`.
6. **Защита прав и обработка 403 Forbidden:** При отсутствии прав редактирования (`permissions.can_edit === false` или пользователь не является автором/админом) силлабус и редактор Puck переходят в режим «Только чтение» (View-Only). Мутирующие операции при получении 403 отображают читаемое уведомление, не затирая локальное состояние форм.
7. **Интеграция загрузки медиа (`POST /upload`):** Для обложек курсов и материалов используется прямая загрузка через `multipart/form-data` (`category=course_cover`, `category=homework`) с валидацией размера/формата на клиенте и мгновенным предпросмотром перед сохранением.

---

## 📁 Структура файлов модуля
| Файл | Описание роли файла |
|---|---|
| [`layout.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/teacher/layout.tsx#L1-L26) | Лейаут авторского кабинета с отключением сайдбара в режиме редактирования Puck |
| [`courses/page.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/teacher/courses/page.tsx#L1-L182) | Дашборд созданных курсов автора со статистикой уроков и студентов |
| [`courses/new/page.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/teacher/courses/new/page.tsx#L1-L247) | Мастер создания нового курса с транслитерацией слага и переходом в силлабус |
| [`courses/[id]/page.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/teacher/courses/%5Bid%5D/page.tsx#L1-L7) | Серверный редирект на маршрут конструктора программы `curriculum` |
| [`courses/[id]/curriculum/page.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/teacher/courses/%5Bid%5D/curriculum/page.tsx#L1-L643) | Центральный конструктор курса: управление модулями/уроками, аналитика, студенты и ДЗ |
| [`courses/[id]/settings/page.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/teacher/courses/%5Bid%5D/settings/page.tsx#L1-L333) | Редактирование метаданных, промо-материалов и статуса курса |
| [`grading/page.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/teacher/grading/page.tsx#L1-L121) | Единый реестр проверки заданий студентов с фильтрацией по курсам |
| [`lessons/[id]/edit/page.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/teacher/lessons/%5Bid%5D/edit/page.tsx#L1-L128) | Полноэкранный визуальный редактор Puck с плагинами блоков и кастомным аутлайном |

---

## ⚙️ Функции, компоненты, хуки и API
| Компонент / Хэндлер | Файл:Строки | Описание | Маршрут / Сигнатура |
|---|---|---|---|
| [`TeacherLayout`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/teacher/layout.tsx#L7-L25) | [`layout.tsx#L7-L25`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/teacher/layout.tsx#L7-L25) | Защищенный каркас кабинета автора с проверкой `/edit` | Макет `/teacher/*` |
| [`TeacherCoursesPage`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/teacher/courses/page.tsx#L18-L181) | [`courses/page.tsx#L18-L181`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/teacher/courses/page.tsx#L18-L181) | Список курсов автора | Маршрут `/teacher/courses` |
| [`CreateCoursePage`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/teacher/courses/new/page.tsx#L9-L246) | [`courses/new/page.tsx#L9-L246`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/teacher/courses/new/page.tsx#L9-L246) | Форма создания курса | Маршрут `/teacher/courses/new` |
| [`TeacherCourseManagementPage`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/teacher/courses/%5Bid%5D/curriculum/page.tsx#L36-L642) | [`curriculum/page.tsx#L36-L642`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/teacher/courses/%5Bid%5D/curriculum/page.tsx#L36-L642) | Конструктор структуры и аналитика курса | Маршрут `/teacher/courses/[id]/curriculum` |
| `handleTogglePublishCourse` | [`curriculum/page.tsx#L136-L157`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/teacher/courses/%5Bid%5D/curriculum/page.tsx#L136-L157) | Публикация или перевод в черновик | `() => Promise<void>` |
| `handleMoveSection` | [`curriculum/page.tsx#L159-L176`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/teacher/courses/%5Bid%5D/curriculum/page.tsx#L159-L176) | Переупорядочивание модулей курса | `(fromIdx: number, toIdx: number) => Promise<void>` |
| [`CourseSettingsPage`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/teacher/courses/%5Bid%5D/settings/page.tsx#L9-L332) | [`settings/page.tsx#L9-L332`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/teacher/courses/%5Bid%5D/settings/page.tsx#L9-L332) | Управление параметрами курса | Маршрут `/teacher/courses/[id]/settings` |
| [`TeacherGlobalGradingPage`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/teacher/grading/page.tsx#L10-L120) | [`grading/page.tsx#L10-L120`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/teacher/grading/page.tsx#L10-L120) | Очередь непроверенных заданий | Маршрут `/teacher/grading` |
| [`LessonEditor`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/teacher/lessons/%5Bid%5D/edit/page.tsx#L12-L127) | [`lessons/[id]/edit/page.tsx#L12-L127`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/teacher/lessons/%5Bid%5D/edit/page.tsx#L12-L127) | Интерактивный редактор урока Puck | Маршрут `/teacher/lessons/[id]/edit` |
| `handleSave (Puck)` | [`lessons/[id]/edit/page.tsx#L49-L60`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/teacher/lessons/%5Bid%5D/edit/page.tsx#L49-L60) | Сериализация и сохранение контента урока через `PATCH /lessons/{id}` | `(data: any) => Promise<void>` |

---

## 🔗 Зависимости
- **Входящие (кто вызывает этот модуль):**
  - Пользователи с ролью `teacher`, `author` или `admin` через переключатель режима в [`Sidebar.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/layout/Sidebar.tsx#L31).
  - Ссылки из страницы входа [`login/page.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/%28auth%29/login/page.tsx#L40).
- **Исходящие (что импортирует этот модуль):**
  - [`lib/api.ts`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/api.ts#L5) — REST API клиент.
  - [`lib/puck-config.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/puck-config.tsx#L211) — схема блоков редактора.
  - [`components/editor/CustomOutline.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/editor/CustomOutline.tsx#L295) — кастомное дерево блоков.
  - Компоненты преподавателя: [`CourseAnalyticsCards`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/teacher/CourseAnalyticsCards.tsx), [`StudentsTable`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/teacher/StudentsTable.tsx), [`PendingHomeworksQueue`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/teacher/PendingHomeworksQueue.tsx), модальные окна `Modal*`.

---

## 🧭 Навигатор типовых задач (Where to edit?)
- **Добавить новый параметр в карточку курса:** обновить форму в [`courses/new/page.tsx#L45-L56`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/teacher/courses/new/page.tsx#L45-L56) и настройки в [`courses/[id]/settings/page.tsx#L77-L88`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/teacher/courses/%5Bid%5D/settings/page.tsx#L77-L88).
- **Скорректировать процесс публикации курса:** редактировать хэндлер `handleTogglePublishCourse` в [`curriculum/page.tsx#L136-L157`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/teacher/courses/%5Bid%5D/curriculum/page.tsx#L136-L157).
- **Изменить поведение панели управления в редакторе Puck:** редактировать `<header>` и пропсы `<Puck>` в [`lessons/[id]/edit/page.tsx#L76-L116`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/teacher/lessons/%5Bid%5D/edit/page.tsx#L76-L116).
