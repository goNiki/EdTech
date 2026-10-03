# 📋 Функциональная спецификация: Каталог и лендинг курса (Courses)

> **Расположение:** `frontend/src/app/courses`  
> **Технический контекст:** Публичный и авторизованный каталог образовательных программ (`/courses`) и промо-страница курса с программой обучения (`/courses/[slug]`).  
> **Дата актуализации:** 2026-10-04  

---

## 🎯 Каталог бизнес-фич модуля

| Фича | Описание возможности | Обеспечивающие функции/компоненты |
|---|---|---|
| **Поиск и фильтрация каталога курсов** | Полнотекстовый поиск с дебаунсом 250 мс, фильтрация по сложности (`difficulty`: начальный, средний, продвинутый), языку (`language`) и сортировка по дате/рейтингу. | [`CoursesCatalogPage.fetchCourses()`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/courses/page.tsx#L67-L96) |
| **Быстрая запись на курс из каталога** | Зачисление авторизованного студента в один клик через `POST /courses/{id}/enroll` без перехода на детальную страницу. Перенаправление гостя на `/login`. | [`CoursesCatalogPage.handleEnroll()`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/courses/page.tsx#L105-L120) |
| **Предпросмотр промо-видео курса** | Всплывающее модальное окно с плеером трейлера курса без покидания каталога. | [`CoursesCatalogPage (previewVideoUrl)`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/courses/page.tsx#L55-L58) |
| **Универсальная страница курса (Slug/ID)** | Загрузка метаданных курса как по человекопонятному URL (`slug`), так и по числовому ID через единый роут `[slug]`. | [`CourseLandingPage.fetchCourse()`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/courses/%5Bslug%5D/page.tsx#L29-L67) |
| **Отображение учебного плана и программы** | Выгрузка полной иерархии модулей (секций) и уроков через `/courses/{id}/structure` с подсветкой бесплатных уроков. | [`CourseLandingPage.structure`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/courses/%5Bslug%5D/page.tsx#L43-L46) |
| **Адаптивный Call-to-Action (Запись / Продолжить)** | Контекстная кнопка: если студент не записан — «Записаться на курс» (с последующим переходом в кабинет); если уже записан — «Продолжить обучение». | [`CourseLandingPage.handleEnroll()`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/courses/%5Bslug%5D/page.tsx#L69-L81) |

---

## 🔬 Паспорта функций и компонентов

### ⚡ Функция: `CoursesCatalogPage.fetchCourses()`

* **Файл и строки:** [`courses/page.tsx#L67-L96`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/courses/page.tsx#L67-L96)
* **Бизнес-назначение:** Загружает витрину доступных курсов с учетом выбранных пользователем фильтров и определяет, на какие из них текущий студент уже зачислен.
* **Связанная фича:** Поиск и фильтрация каталога курсов

#### 📥 Входные параметры (состояние фильтров)
| Параметр | Тип | Обязателен | Бизнес-смысл и ограничения |
|---|---|:---:|---|
| `search` | `string` | Нет | Поисковая подстрока (по названию или описанию курса). |
| `difficulty` | `string` | Нет | Уровень сложности: `'beginner'`, `'intermediate'`, `'advanced'`. |
| `language` | `string` | Нет | Язык курса (например, `'ru'`, `'en'`). |
| `sortBy` | `string` | Да | Поле сортировки (по умолчанию `'created_at'`). |
| `sortOrder` | `'asc' \| 'desc'` | Да | Направление сортировки (по умолчанию `'desc'`). |

#### 🔄 Пошаговый алгоритм работы
1. **Шаг 1 (Сбор Query-параметров):** Формируется объект `URLSearchParams` из заполненных полей фильтра.
2. **Шаг 2 (Сетевой запрос каталога):** `GET /courses?${params.toString()}` с получением массива курсов.
3. **Шаг 3 (Проверка личных записей студента):** Если `isAuthenticated === true`, параллельно выполняется запрос `GET /courses/my`. Из ответа извлекаются ID курсов, на которые записан пользователь, и сохраняются в стейт `enrolledCourseIds`.
4. **Шаг 4 (Отображение карточек):** Карточки курсов помечаются бейджами «Вы записаны» при совпадении ID.

#### ⚠️ Побочные эффекты (Side Effects)
* **Сетевые вызовы:** `GET /courses`, `GET /courses/my`.
* **State:** Обновление массивов `courses`, `enrolledCourseIds`, снятие флага `isLoading`.

#### 💡 Подсказка для аналитика (Где менять логику?)
* *Таймаут дебаунса поисковой строки:* [`courses/page.tsx#L99-L102`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/courses/page.tsx#L99-L102).
* *Дефолтное поле и порядок сортировки:* [`courses/page.tsx#L52-L53`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/courses/page.tsx#L52-L53).

---

### ⚡ Функция: `CoursesCatalogPage.handleEnroll(courseId, e)`

* **Файл и строки:** [`courses/page.tsx#L105-L120`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/courses/page.tsx#L105-L120)
* **Бизнес-назначение:** Позволяет студенту записаться на курс сразу со страницы каталога.
* **Связанная фича:** Быстрая запись на курс из каталога

#### 📥 Входные параметры
| Параметр | Тип | Обязателен | Бизнес-смысл и ограничения |
|---|---|:---:|---|
| `courseId` | `number` | Да | Идентификатор выбранного курса. |
| `e` | `React.MouseEvent` | Да | Событие клика для вызова `e.stopPropagation()`. |

#### 🔄 Пошаговый алгоритм работы
1. **Шаг 1 (Изоляция клика):** `e.stopPropagation()` предотвращает переход на страницу описания курса при клике по кнопке записи внутри карточки.
2. **Шаг 2 (Проверка авторизации):** Если `!isAuthenticated`, пользователя перенаправляет на `/login`.
3. **Шаг 3 (Сетевой запрос записи):** `POST /courses/${courseId}/enroll`.
4. **Шаг 4 (Мгновенное обновление UI):** `courseId` добавляется в локальный массив `enrolledCourseIds`, карточка переключается в состояние «Вы записаны», отображается Toast-уведомление «Вы успешно записались на курс!».
5. **Шаг 5 (Обработка ошибок):** Показ серверного сообщения об ошибке через Toast.

#### ⚠️ Побочные эффекты (Side Effects)
* **Сетевой вызов:** `POST /courses/{courseId}/enroll`.
* **Toast:** Всплывающее уведомление в нижнем правом углу.

#### 💡 Подсказка для аналитика (Где менять логику?)
* *Поведение для гостей:* [`courses/page.tsx#L107-L110`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/courses/page.tsx#L107-L110).

---

### ⚡ Функция: `CourseLandingPage.handleEnroll()`

* **Файл и строки:** [`courses/[slug]/page.tsx#L69-L81`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/courses/%5Bslug%5D/page.tsx#L69-L81)
* **Бизнес-назначение:** Запись на курс со страницы лендинга и моментальный перевод студента в рабочий кабинет курса для начала обучения.
* **Связанная фича:** Адаптивный Call-to-Action (Запись / Продолжить)

#### 📥 Входные параметры
* Параметры берутся из состояния страницы (`courseData.id`, `isAuthenticated`).

#### 🔄 Пошаговый алгоритм работы
1. **Шаг 1 (Проверка гостя):** Если пользователь не авторизован, редирект на `/login`.
2. **Шаг 2 (Сетевой запрос):** `POST /courses/${courseData.id}/enroll`.
3. **Шаг 3 (Переход в учебную среду):** Устанавливается `isEnrolled = true`, и вызывается `router.push('/dashboard/courses/' + courseData.id)`.

#### ⚠️ Побочные эффекты (Side Effects)
* **Навигация:** Редирект в учебный интерфейс курса `/dashboard/courses/[id]`.
* **Сетевой запрос:** Создание записи enrollment на бэкенде.

#### 💡 Подсказка для аналитика (Где менять логику?)
* *Маршрут перехода после успешной записи:* [`courses/[slug]/page.tsx#L77`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/courses/%5Bslug%5D/page.tsx#L77).
