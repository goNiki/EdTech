# 📋 Функциональная спецификация: Кабинет студента и профиль (Dashboard)

> **Расположение:** `frontend/src/app/dashboard`  
> **Технический контекст:** Защищенная рабочая среда студента (`/dashboard`), список активных программ (`/dashboard/courses`), навигация по структуре курса (`/dashboard/courses/[id]`) и настройки профиля (`/dashboard/settings`).  
> **Дата актуализации:** 2026-10-04  

---

## 🎯 Каталог бизнес-фич модуля

| Фича | Описание возможности | Обеспечивающие функции/компоненты |
|---|---|---|
| **Главный экран обучения студента** | Отображение активных зачислений студента через `/courses/my` с быстрыми ссылками для продолжения обучения и заглушкой при отсутствии курсов. | [`StudentDashboard`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/dashboard/page.tsx#L8-L63) |
| **Реестр курсов с прогресс-барами** | Отображение карточек изучаемых программ с расчетом процента освоения (`progress_percent`), числом уроков и средним баллом за тесты. | [`MyCoursesPage`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/dashboard/courses/page.tsx#L21-L157) |
| **Учебный план и навигатор по курсу** | Иерархическое дерево модулей и уроков курса с аккордеоном, отметками завершения (✓), бейджами типов уроков (теория, тест, практика) и кнопкой быстрого перехода к следующему непройденному уроку. | [`StudentCoursePlayerPage`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/dashboard/courses/%5Bid%5D/page.tsx#L26-L447) |
| **Монитор успеваемости на курсе** | Расчет процента завершения программы и вычисление среднего балла студента за все контрольные тестирования в рамках курса. | [`calculatedAvgScore & progressPercent`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/dashboard/courses/%5Bid%5D/page.tsx#L114-L126) |
| **Редактирование профиля пользователя** | Изменение личных данных (имя, фамилия, биография, ссылка на аватар) через `PATCH /auth/profile` с моментальной синхронизацией в хранилище Zustand. | [`handleSaveProfile()`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/dashboard/settings/page.tsx#L47-L69) |
| **Настройки темы интерфейса** | Переключение между светлой, темной и системной темами оформления с сохранением предпочтений. | [`ProfileAndSettingsPage (preferences)`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/dashboard/settings/page.tsx#L103-L112) |

---

## 🔬 Паспорта функций и компонентов

### ⚡ Компонент: `StudentDashboard`

* **Файл и строки:** [`dashboard/page.tsx#L8-L63`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/dashboard/page.tsx#L8-L63)
* **Бизнес-назначение:** Персональная стартовая панель студента: быстрый доступ ко всем начатым программам.
* **Связанная фича:** Главный экран обучения студента

#### 📥 Входные параметры
* Параметров нет.

#### 🔄 Пошаговый алгоритм работы
1. **Шаг 1:** Выполняет сетевой запрос `GET /courses/my`.
2. **Шаг 2:** Если курсов нет — отображает empty-state с кнопкой перехода в каталог.
3. **Шаг 3:** При наличии курсов — рендерит сетку компонентов `CourseCard`, при клике на которые открывается плеер курса `/dashboard/courses/[id]`.

#### ⚠️ Побочные эффекты (Side Effects)
* **Сетевой запрос:** `GET /courses/my`.

#### 💡 Подсказка для аналитика (Где менять логику?)
* *URL перехода по карточке курса:* [`dashboard/page.tsx#L56`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/dashboard/page.tsx#L56).

---

### ⚡ Компонент: `StudentCoursePlayerPage`

* **Файл и строки:** [`dashboard/courses/[id]/page.tsx#L26-L447`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/dashboard/courses/%5Bid%5D/page.tsx#L26-L447)
* **Бизнес-назначение:** Полнофункциональный интерактивный навигатор студента по материалам курса: отображает дерево уроков, рассчитывает динамический прогресс и подсказывает следующий урок.
* **Связанная фича:** Учебный план и навигатор по курсу

#### 📥 Входные параметры
| Параметр | Тип | Обязателен | Бизнес-смысл и ограничения |
|---|---|:---:|---|
| `params.id` | `string` | Да | Идентификатор курса из динамического сегмента URL. |

#### 🔄 Пошаговый алгоритм работы
1. **Шаг 1 (Загрузка структуры):** `GET /courses/${id}/structure` возвращает модули (sections) и вложенные уроки (lessons).
2. **Шаг 2 (Загрузка прогресса):**
   - `GET /courses/${id}/progress` — сводные данные по прохождению курса.
   - `GET /courses/${id}/progress/lessons` — массив статусов по каждому отдельному уроку (завершен/в процессе, набранный балл).
3. **Шаг 3 (Формирование карты урока):** Создается словарь `lessonProgressMap[lessonId]`, где ключ — ID урока.
4. **Шаг 4 (Определение следующего шага `nextLesson`):** Система обходит все уроки по порядку и находит первый урок со статусом отличным от `'completed'`.
5. **Шаг 5 (Вычисление метрик успеваемости):**
   - `completedCount` — число завершенных уроков.
   - `progressPercent` — процент завершения курса: `Math.round((completedCount / totalLessonsCount) * 100)`.
   - `calculatedAvgScore` — средний балл за все пройденные тесты курса (по шкале 100).
6. **Шаг 6 (Интерактивный аккордеон):** Клик по заголовку модуля переключает его видимость (`toggleSection`). Клик по уроку выполняет переход на экран урока `/lessons/[id]`.

#### ⚠️ Побочные эффекты (Side Effects)
* **Сетевые вызовы:** `GET /courses/{id}/structure`, `GET /courses/{id}/progress`, `GET /courses/{id}/progress/lessons`.
* **Навигация:** Переход к странице конкретного урока `/lessons/{id}`.

#### 💡 Подсказка для аналитика (Где менять логику?)
* *Логика выбора следующего рекомендуемого урока:* [`[id]/page.tsx#L102-L106`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/dashboard/courses/%5Bid%5D/page.tsx#L102-L106).
* *Формула вычисления среднего балла:* [`[id]/page.tsx#L114-L121`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/dashboard/courses/%5Bid%5D/page.tsx#L114-L121).

---

### ⚡ Функция: `ProfileAndSettingsPage.handleSaveProfile(e)`

* **Файл и строки:** [`dashboard/settings/page.tsx#L47-L69`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/dashboard/settings/page.tsx#L47-L69)
* **Бизнес-назначение:** Сохраняет отредактированные персональные данные пользователя и обновляет глобальный стейт приложения.
* **Связанная фича:** Редактирование профиля пользователя

#### 📥 Входные параметры (состояние формы)
| Поле | Тип | Обязателен | Бизнес-смысл |
|---|---|:---:|---|
| `firstName` | `string` | Нет | Имя пользователя. |
| `lastName` | `string` | Нет | Фамилия пользователя. |
| `bio` | `string` | Нет | Текстовое описание о себе / специализация. |
| `avatarUrl` | `string` | Нет | URL ссылки на изображение аватара. |

#### 🔄 Пошаговый алгоритм работы
1. **Шаг 1:** Блокировка кнопки сохранения `setIsSaving(true)`.
2. **Шаг 2:** Отправка HTTP-запроса `PATCH /auth/profile` с объектом обновленных полей.
3. **Шаг 3:** При успешном ответе вызывается `setUser({ ...user, ...updated })` в сторе `useAuth`, обновляя отображение имени и аватара в шапке и боковом меню без перезагрузки страницы.
4. **Шаг 4:** Закрытие режима редактирования `setIsEditing(false)` и показ зеленого индикатора успеха на 3 секунды.

#### ⚠️ Побочные эффекты (Side Effects)
* **Сетевой вызов:** `PATCH /auth/profile`.
* **Zustand State:** Мгновенная мутация `user`.

#### 💡 Подсказка для аналитика (Где менять логику?)
* *Эндпоинт обновления профиля:* [`settings/page.tsx#L51-L56`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/dashboard/settings/page.tsx#L51-L56).
