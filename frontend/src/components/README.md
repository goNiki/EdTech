# 📦 Модуль: `components`
> **Путь:** `frontend/src/components`  
> **Роль:** Компонентная библиотека платформы: общие лейауты, защищенные маршруты, интерактивный плеер уроков Puck, визуальный редактор контента, виджеты кабинета преподавателя и базовые UI-примитивы Shadcn.

---

## 🎯 Назначение и ответственность
Модуль содержит все переиспользуемые UI-компоненты приложения. Архитектурно он разделен на логические подсистемы:
- **Лейаут и навигация (`layout/`):** Сайдбар со сменой роли (`student` / `teacher`), адаптивная шапка `TopNavbar` с переключением темы.
- **Безопасность (`ProtectedRoute.tsx`):** Клиентский барьер аутентификации и ролевого контроля (RBAC).
- **Интерактивный плеер (`player/PuckLessonViewer.tsx`):** Движок интерактивного прохождения урока студентом со скорингом и сбором эссе в реальном времени.
- **Инструменты визуального редактора (`editor/`):** Интерактивное дерево блоков `CustomOutline` и инлайн-редакторы текста `InlineEditable`.
- **Инструменты преподавателя (`teacher/`):** Модальные окна управления курсом, аналитические карточки, таблицы студентов и очередь ручной проверки ДЗ.
- **Атомарный UI (`ui/`):** Дизайн-система кнопок, инпутов, карточек и аккордеонов на базе Tailwind CSS 4.

---

## ⚠️ Жестко установленные правила (Hard Invariants)
1. **Ролевой контроль в ProtectedRoute:** Компонент обязан блокировать рендеринг дочерних элементов (`children`) до завершения `isChecking` и `isLoading`. Если пользователь не аутентифицирован — выполняется переход на `/login`. Если роль пользователя отсутствует в `allowedRoles` — отображается экран "Доступ запрещен" (403), предотвращая несанкционированный переход в авторский кабинет.
2. **Контракт вычисления баллов в PuckLessonViewer:** Все интерактивные блоки квизов суммируют свои максимальные баллы в `totalMaxPoints` и заработанные очки в `earnedPoints`. Итоговый процент вычисляется по формуле `Math.round((earnedPoints / totalMaxPoints) * 100)`. Текстовые эссе отправляются в массив `essays` для ручной проверки автором.
3. **Безопасный диспатч в InlineEditable:** Хук `usePuckPropUpdater` обязан проверять наличие инстанса Puck через `try { puckApi = usePuck() } catch` во избежание сбоев вне контекста редактора.
4. **Синхронизация состояния модалок:** Все модальные окна создания и редактирования сущностей (`ModalCreateModule`, `ModalCreateLesson`, `ModalEditModule`, `ModalEditLesson`) обязаны принимать пропс `onClose` и колбэк завершения (`onCreated`/`onSaved`), инициирующий рефетч структуры курса.

---

## 📁 Структура файлов модуля
| Каталог / Файл | Описание роли файла |
|---|---|
| [`CourseCard.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/CourseCard.tsx#L1-L49) | Карточка курса с обложкой, бейджем сложности и переходом в каталог или плеер |
| [`ProtectedRoute.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/ProtectedRoute.tsx#L1-L61) | HOC-компонент защиты маршрутов по авторизации и списку ролей |
| [`ThemeProvider.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/ThemeProvider.tsx#L1-L12) | Провайдер темной/светлой темы на базе `next-themes` |
| **`layout/`** | |
| [`layout/Sidebar.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/layout/Sidebar.tsx#L1-L217) | Боковая панель навигации, переключатель режима Студент/Преподаватель, диалог выхода |
| [`layout/TopNavbar.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/layout/TopNavbar.tsx#L1-L55) | Верхняя панель со стрелкой возврата, заголовком страницы и кнопкой темы |
| **`player/`** | |
| [`player/PuckLessonViewer.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/player/PuckLessonViewer.tsx#L1-L1008) | Интерактивный плеер контента урока: тесты, сопоставление, пропуски, эссе, автогрейдинг |
| **`editor/`** | |
| [`editor/CustomOutline.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/editor/CustomOutline.tsx#L1-L296) | Дерево структуры урока в Puck с переупорядочиванием, дублированием и превью блоков |
| [`editor/InlineEditable.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/editor/InlineEditable.tsx#L1-L1234) | Компоненты прямого редактирования текста, панель форматирования Word-like и хук `usePuckPropUpdater` |
| **`teacher/`** | |
| [`teacher/CourseAnalyticsCards.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/teacher/CourseAnalyticsCards.tsx#L1-L78) | Сводные KPI-карточки курса: студенты, средний прогресс, средний балл, очередь ДЗ |
| [`teacher/ModalAddStudent.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/teacher/ModalAddStudent.tsx#L1-L111) | Модальное окно ручного зачисления студента по Email или ID пользователя |
| [`teacher/ModalCreateLesson.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/teacher/ModalCreateLesson.tsx#L1-L217) | Модальное окно создания урока в секции с выбором типа и длительности |
| [`teacher/ModalCreateModule.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/teacher/ModalCreateModule.tsx#L1-L153) | Модальное окно создания нового модуля (секции) курса |
| [`teacher/ModalEditLesson.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/teacher/ModalEditLesson.tsx#L1-L172) | Модальное окно редактирования метаданных и статуса урока |
| [`teacher/ModalEditModule.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/teacher/ModalEditModule.tsx#L1-L193) | Модальное окно редактирования модуля и изменения порядка уроков |
| [`teacher/ModalGradeHW.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/teacher/ModalGradeHW.tsx#L1-L146) | Модальное окно выставления баллов и обратной связи по эссе |
| [`teacher/ModalStudentDrilldown.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/teacher/ModalStudentDrilldown.tsx#L1-L198) | Модальное окно детальной аналитики успеваемости конкретного студента |
| [`teacher/PendingHomeworksQueue.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/teacher/PendingHomeworksQueue.tsx#L1-L101) | Список сданных домашних заданий, ожидающих проверки преподавателем |
| [`teacher/StudentsTable.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/teacher/StudentsTable.tsx#L1-L140) | Таблица зачисленных студентов со шкалами прогресса и кнопками перехода к отчету |
| **`ui/`** | |
| [`ui/button.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/ui/button.tsx), [`card.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/ui/card.tsx), [`input.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/ui/input.tsx), ... | Атомарные UI-примитивы Shadcn / Radix UI |

---

## ⚙️ Функции, компоненты, хуки и API
| Компонент / Хук | Файл:Строки | Описание | Пропсы / Сигнатура |
|---|---|---|---|
| [`CourseCard`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/CourseCard.tsx#L14-L48) | [`CourseCard.tsx#L14-L48`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/CourseCard.tsx#L14-L48) | Карточка отображения курса | `{ course: Course, href?: string }` |
| [`ProtectedRoute`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/ProtectedRoute.tsx#L7-L60) | [`ProtectedRoute.tsx#L7-L60`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/ProtectedRoute.tsx#L7-L60) | Барьер авторизации и ролевого доступа | `{ children: ReactNode, allowedRoles?: string[] }` |
| [`Sidebar`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/layout/Sidebar.tsx#L21-L216) | [`Sidebar.tsx#L21-L216`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/layout/Sidebar.tsx#L21-L216) | Главный навигационный сайдбар со сворачиванием | `() => JSX.Element` |
| [`TopNavbar`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/layout/TopNavbar.tsx#L15-L54) | [`TopNavbar.tsx#L15-L54`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/layout/TopNavbar.tsx#L15-L54) | Шапка страницы с историей навигации | `{ title?: string, subtitle?: string, showBack?: boolean }` |
| [`PuckLessonViewer`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/player/PuckLessonViewer.tsx#L39-L1007) | [`PuckLessonViewer.tsx#L39-L1007`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/player/PuckLessonViewer.tsx#L39-L1007) | Интерактивный движок решения урока | `{ contentJson: string \| object, onComplete?, onNavigateBack?, initialProgress? }` |
| [`usePuckPropUpdater`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/editor/InlineEditable.tsx#L42-L100) | [`InlineEditable.tsx#L42-L100`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/editor/InlineEditable.tsx#L42-L100) | Хук реактивного обновления свойств блока Puck | `(blockId?: string) => { selectThisBlock, updateProp }` |
| [`CustomOutline`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/editor/CustomOutline.tsx#L160-L294) | [`CustomOutline.tsx#L160-L294`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/editor/CustomOutline.tsx#L160-L294) | Сайдбар дерева компонентов редактора Puck | Кастомный оверлей для Puck plugin |
| [`CourseAnalyticsCards`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/teacher/CourseAnalyticsCards.tsx#L13-L77) | [`CourseAnalyticsCards.tsx#L13-L77`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/teacher/CourseAnalyticsCards.tsx#L13-L77) | Блок 4 ключевых метрик успеваемости | `{ totalStudents, avgProgress, avgScore, pendingCount }` |
| [`StudentsTable`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/teacher/StudentsTable.tsx#L23-L139) | [`StudentsTable.tsx#L23-L139`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/teacher/StudentsTable.tsx#L23-L139) | Таблица студентов курса с действиями | `{ students, onOpenDrilldown, onAddStudent, onRemoveStudent? }` |
| [`PendingHomeworksQueue`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/teacher/PendingHomeworksQueue.tsx#L26-L100) | [`PendingHomeworksQueue.tsx#L26-L100`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/teacher/PendingHomeworksQueue.tsx#L26-L100) | Список непроверенных работ студентов | `{ items, onOpenGradeModal }` |
| [`ModalGradeHW`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/teacher/ModalGradeHW.tsx#L14-L145) | [`ModalGradeHW.tsx#L14-L145`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/teacher/ModalGradeHW.tsx#L14-L145) | Диалог проверки практической работы | `{ hw, onClose, onGraded }` |
| [`ModalStudentDrilldown`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/teacher/ModalStudentDrilldown.tsx#L13-L197) | [`ModalStudentDrilldown.tsx#L13-L197`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/teacher/ModalStudentDrilldown.tsx#L13-L197) | Диалог детального прогресса учащегося | `{ courseId, studentId, onClose }` |

---

## 🔗 Зависимости
- **Входящие (кто вызывает этот модуль):**
  - Все страницы приложения (`app/courses`, `app/dashboard`, `app/lessons`, `app/teacher`).
- **Исходящие (что импортирует этот модуль):**
  - [`store/useAuth.ts`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/store/useAuth.ts#L33) — глобальное состояние пользователя.
  - [`lib/api.ts`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/api.ts#L5) — Axios-клиент для отправки оценок, зачисления и сохранения уроков.
  - [`lib/puck-config.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/puck-config.tsx#L105) — парсер шаблонов пропусков.
  - `@puckeditor/core`, `lucide-react`, `next-themes`.

---

## 🧭 Навигатор типовых задач (Where to edit?)
- **Изменить логику защиты маршрутов:** скорректировать проверки в [`ProtectedRoute.tsx#L18-L59`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/ProtectedRoute.tsx#L18-L59).
- **Добавить пункт навигации в левое меню:** обновить массивы `studentNavItems` или `teacherNavItems` в [`Sidebar.tsx#L41-L55`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/layout/Sidebar.tsx#L41-L55).
- **Изменить интерфейс или логику скоринга урока:** редактировать [`PuckLessonViewer.tsx#L112-L150`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/player/PuckLessonViewer.tsx#L112-L150).
- **Скорректировать поля создания урока или модуля:** редактировать [`ModalCreateLesson.tsx#L42-L65`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/teacher/ModalCreateLesson.tsx#L42-L65) или [`ModalCreateModule.tsx#L35-L51`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/teacher/ModalCreateModule.tsx#L35-L51).
