# 📋 Функциональная спецификация: Кабинет преподавателя и конструктор курсов (Teacher)

> **Расположение:** `frontend/src/app/teacher`  
> **Технический контекст:** Пространство управления контентом и курсами для ролей `teacher`, `author` и `admin`, включающее визуальный конструктор Puck, конструктор учебного плана, аналитику и систему грейдинга ДЗ.  
> **Дата актуализации:** 2026-10-04  

---

## 🎯 Каталог бизнес-фич модуля

| Фича | Описание возможности | Обеспечивающие функции/компоненты |
|---|---|---|
| **Реестр курсов автора** | Каталог созданных автором курсов с индикаторами статуса публикации (Черновик / Опубликован), видимости и счетчиками уроков и зачисленных учеников. | [`TeacherCoursesPage`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/teacher/courses/page.tsx#L18-L181) |
| **Мастер создания курса** | Пошаговое создание нового курса с автогенерацией URL-слага из названия и немедленным перенаправлением в конструктор учебного плана. | [`CreateCoursePage.handleSubmit()`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/teacher/courses/new/page.tsx#L36-L72) |
| **Управление учебным планом (Curriculum Builder)** | Создание, редактирование, удаление и сортировка модулей (секций) и уроков с помощью drag-free стрелок реордеринга. | [`handleMoveSection & handleMoveLessonInList`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/teacher/courses/%5Bid%5D/curriculum/page.tsx#L159-L195) |
| **Публикация курса в один клик** | Мгновенный перевод всего курса и дочерних модулей/уроков из статуса «Черновик» в «Опубликован» (`POST /courses/{id}/publish`) и обратно. | [`handleTogglePublishCourse()`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/teacher/courses/%5Bid%5D/curriculum/page.tsx#L136-L157) |
| **Студенческая аналитика курса** | Вычисление метрик курса: число записанных учеников, средний процент прохождения, средний балл тестирования и число неотвеченных ДЗ. | [`TeacherCourseManagementPage.analytics`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/teacher/courses/%5Bid%5D/curriculum/page.tsx#L90-L100) |
| **Управление студентами курса** | Просмотр реестра записанных учеников, ручное зачисление по Email/ID и просмотр детального отчета успеваемости (Drilldown). | [`StudentsTable & ModalStudentDrilldown`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/teacher/courses/%5Bid%5D/curriculum/page.tsx#L102-L117) |
| **Очередь проверки домашних заданий** | Просмотр сданных студентами работ, требующих ручной проверки (эссе, загруженные архивы) с фильтрацией по курсам и формой грейдинга. | [`TeacherGlobalGradingPage`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/teacher/grading/page.tsx#L10-L120) |
| **Визуальный конструктор уроков Puck** | Полноэкранный визуальный редактор урока на базе Puck с поддержкой 12 типов блоков, дерева структуры `CustomOutline` и инлайн-редактирования. | [`LessonEditor.handleSave()`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/teacher/lessons/%5Bid%5D/edit/page.tsx#L49-L60) |
| **Настройки параметров курса** | Изменение метаданных курса, обложки, видео, уровня сложности, языка и категории через `PATCH /courses/{id}`. | [`CourseSettingsPage.handleSubmit()`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/teacher/courses/%5Bid%5D/settings/page.tsx#L67-L106) |

---

## 🔬 Паспорта функций и компонентов

### ⚡ Функция: `CreateCoursePage.handleSubmit(e)`

* **Файл и строки:** [`teacher/courses/new/page.tsx#L36-L72`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/teacher/courses/new/page.tsx#L36-L72)
* **Бизнес-назначение:** Регистрирует новую образовательную программу в системе со статусом «Черновик» и сразу переводит преподавателя к наполнению учебного плана.
* **Связанная фича:** Мастер создания курса

#### 📥 Входные параметры (состояние формы)
| Поле | Тип | Обязателен | Бизнес-смысл и ограничения |
|---|---|:---:|---|
| `title` | `string` | Да | Название курса. Не должно быть пустым. |
| `slug` | `string` | Да | Уникальный буквенно-цифровой идентификатор в URL (автоматически формируется транслитом). |
| `short_description` | `string` | Нет | Краткий анонс курса для карточки витрины. |
| `description` | `string` | Нет | Полная программа и результаты обучения. |
| `cover_url` | `string` | Нет | Ссылка на обложку (фоллбек на плейсхолдер при отсутствии). |
| `intro_video_url` | `string` | Нет | Ссылка на промо-видео (YouTube). |
| `difficulty` | `string` | Да | Уровень сложности: `'beginner'`, `'intermediate'`, `'advanced'`. |
| `language` | `string` | Да | Язык обучения (по умолчанию `'RU'`). |
| `visibility` | `string` | Да | Доступность в каталоге: `'public'` или `'private'`. |

#### 🔄 Пошаговый алгоритм работы
1. **Шаг 1 (Валидация):** Проверка заполненности поля `title.trim()`. Если пусто — вывод предупреждения.
2. **Шаг 2 (Сборка полезной нагрузки):** Формирование объекта курса. Если `slug` пустой, генерируется фоллбек `course-${Date.now()}`.
3. **Шаг 3 (Сетевой запрос):** Вызов `POST /courses`.
4. **Шаг 4 (Редирект в учебный план):** Из ответа извлекается ID созданного курса, и выполняется `router.push('/teacher/courses/' + newId + '/curriculum')`.

#### ⚠️ Побочные эффекты (Side Effects)
* **Сетевой вызов:** `POST /courses`.
* **Навигация:** Переход на страницу конструктора программы курса.

#### 💡 Подсказка для аналитика (Где менять логику?)
* *Правило транслитерации слага:* [`new/page.tsx#L27-L34`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/teacher/courses/new/page.tsx#L27-L34).
* *Дефолтное изображение обложки:* [`new/page.tsx#L50`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/teacher/courses/new/page.tsx#L50).

---

### ⚡ Функция: `TeacherCourseManagementPage.handleTogglePublishCourse()`

* **Файл и строки:** [`teacher/courses/[id]/curriculum/page.tsx#L136-L157`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/teacher/courses/%5Bid%5D/curriculum/page.tsx#L136-L157)
* **Бизнес-назначение:** Позволяет автору публиковать курс для открытия доступа студентам или скрывать его на доработку без удаления данных.
* **Связанная фича:** Публикация курса в один клик

#### 📥 Входные параметры
* Параметры берутся из текущего состояния `courseData.status` и `id`.

#### 🔄 Пошаговый алгоритм работы
1. **Шаг 1 (Проверка текущего статуса):** Определяется текущее состояние публикации (`isCurrentlyPublished`).
2. **Шаг 2 (Сетевой запрос мутации):**
   - Если курс был черновиком: отправляется `POST /courses/${id}/publish`.
   - Если курс был опубликован: отправляется `PATCH /courses/${id}/status` с телом `{ status: 'draft' }`.
3. **Шаг 3 (Уведомление и перезагрузка данных):**
   - Отображается соответствующий Toast: «Курс и все связанные модули и уроки успешно опубликованы!» либо «Курс переведен в статус Черновик».
   - Вызывается `fetchCourseData()` для синхронизации бейджей в шапке.

#### ⚠️ Побочные эффекты (Side Effects)
* **Сетевые вызовы:** `POST /courses/{id}/publish` или `PATCH /courses/{id}/status`.
* **Toast:** Всплывающее системное сообщение.

#### 💡 Подсказка для аналитика (Где менять логику?)
* *Эндпоинт публикации курса:* [`curriculum/page.tsx#L142-L147`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/teacher/courses/%5Bid%5D/curriculum/page.tsx#L142-L147).

---

### ⚡ Функции сортировки: `handleMoveSection` и `handleMoveLessonInList`

* **Файл и строки:** [`curriculum/page.tsx#L159-L195`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/teacher/courses/%5Bid%5D/curriculum/page.tsx#L159-L195)
* **Бизнес-назначение:** Позволяет автору настраивать последовательность изучения тем и уроков студентами.
* **Связанная фича:** Управление учебным планом (Curriculum Builder)

#### 📥 Входные параметры
| Параметр | Тип | Обязателен | Бизнес-смысл |
|---|---|:---:|---|
| `fromIdx` | `number` | Да | Исходный индекс элемента в массиве. |
| `toIdx` | `number` | Да | Целевой индекс смещения (вверх `idx - 1` или вниз `idx + 1`). |
| `item_ids` | `number[]` | Да | Упорядоченный список идентификаторов для сохранения в БД. |

#### 🔄 Пошаговый алгоритм работы
1. **Шаг 1 (Оптимистичное обновление):** Элемент перемещается в локальном массиве `splice(fromIdx, 1)` ➔ `splice(toIdx, 0, moved)`.
2. **Шаг 2 (Сетевая синхронизация порядка):**
   - Для модулей: `PUT /courses/${id}/reorder-sections` с телом `{ item_ids: secIds }`.
   - Для уроков: `PUT /sections/${sectionId}/reorder-lessons` с телом `{ item_ids: lessonIds }`.
3. **Шаг 3 (Оповещение):** Показ Toast-сообщения «Порядок модулей обновлен» / «Порядок уроков сохранен». При сбое состояние откатывается через `fetchCourseData()`.

#### ⚠️ Побочные эффекты (Side Effects)
* **Сетевые запросы:** `PUT /courses/{id}/reorder-sections`, `PUT /sections/{sectionId}/reorder-lessons`.

#### 💡 Подсказка для аналитика (Где менять логику?)
* *Эндпоинты реордеринга:* [`curriculum/page.tsx#L168`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/teacher/courses/%5Bid%5D/curriculum/page.tsx#L168) и [`curriculum/page.tsx#L186`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/teacher/courses/%5Bid%5D/curriculum/page.tsx#L186).

---

### ⚡ Функция: `LessonEditor.handleSave(data)`

* **Файл и строки:** [`teacher/lessons/[id]/edit/page.tsx#L49-L60`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/teacher/lessons/%5Bid%5D/edit/page.tsx#L49-L60)
* **Бизнес-назначение:** Сохраняет визуально скомпонованную структуру урока и все созданные интерактивные тесты в базу данных в формате Content-as-Data.
* **Связанная фича:** Визуальный конструктор уроков Puck

#### 📥 Входные параметры
| Параметр | Тип | Обязателен | Бизнес-смысл |
|---|---|:---:|---|
| `data` | `Data` | Да | Внутреннее дерево состояния редактора Puck (`{ content: [...blocks], root: {...} }`). |

#### 🔄 Пошаговый алгоритм работы
1. **Шаг 1 (Сериализация структуры):** Преобразование объекта блоков в JSON-строку: `contentString = JSON.stringify(data)`.
2. **Шаг 2 (Сетевой запрос):** `PATCH /lessons/${id}` с телом `{ content: contentString }`.
3. **Шаг 3 (Оповещение об успехе):** Отображение анимированного Toast-уведомления: «Контент и интерактивные тесты успешно сохранены!».

#### ⚠️ Побочные эффекты (Side Effects)
* **БД:** Перезапись содержимого поля `content` урока.
* **Toast Notification:** Зеленое уведомление.

#### 💡 Подсказка для аналитика (Где менять логику?)
* *Эндпоинт сохранения контента урока:* [`edit/page.tsx#L52-L54`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/teacher/lessons/%5Bid%5D/edit/page.tsx#L52-L54).
