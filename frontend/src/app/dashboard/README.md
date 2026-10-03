# 📦 Модуль: `dashboard`
> **Путь:** `frontend/src/app/dashboard`  
> **Роль:** Личный кабинет студента: дашборд успеваемости, список записанных курсов, интерактивный навигатор по структуре курса ([id]) и редактирование профиля/настроек.

---

## 🎯 Назначение и ответственность
Модуль представляет собой рабочее пространство учащегося платформы:
1. Лейаут (`layout.tsx`): глобальная изоляция под защитой `ProtectedRoute` с боковой панелью `Sidebar`.
2. Стартовый экран (`page.tsx`) и список курсов (`courses/page.tsx`): отображение активных программ обучения, шкал прогресса (0-100%) и средних баллов за тесты.
3. Просмотр структуры курса (`courses/[id]/page.tsx`): подробный силлабус с индикацией пройденных и заблокированных уроков, кнопкой «Продолжить обучение» и навигацией в плеер урока.
4. Настройки аккаунта (`settings/page.tsx`): редактирование имени, фамилии, био и аватара с сохранением через `PATCH /auth/profile`, а также переключение темы оформления.

---

## ⚠️ Жестко установленные правила (Hard Invariants)
1. **Защита на уровне лейаута:** Весь модуль `dashboard` обернут в `<ProtectedRoute allowedRoles={['student', 'teacher', 'admin', 'author']}>`. Доступ неавторизованных пользователей невозможен ни к одной из вложенных страниц.
2. **Синхронизация прогресса без сбоев (Safe Progress Hydration):** На странице курса `[id]` выполняются три параллельных запроса (`/structure`, `/progress`, `/progress/lessons`). Если записи о прогрессе в БД еще отсутствуют (студент только записался), ошибки мягко перехватываются, а прогресс инициализируется как 0%.
3. **Реактивное обновление состояния пользователя:** При сохранении профиля в `settings/page.tsx` обязателен вызов `setUser({ ...user, ...updated })` для синхронного обновления имени и аватара в сайдбаре без перезагрузки вкладки.
4. **Контекстный роутинг в плеер:** Клик по доступному уроку перенаправляет в отдельный маршрут плеера [`/lessons/${lessonId}`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/lessons/%5Bid%5D/page.tsx).

---

## 📁 Структура файлов модуля
| Файл | Описание роли файла |
|---|---|
| [`layout.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/dashboard/layout.tsx#L1-L18) | Базовый макет с `Sidebar` и ролевой защитой `ProtectedRoute` |
| [`page.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/dashboard/page.tsx#L1-L64) | Главная страница студента со списком активных курсов |
| [`courses/page.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/dashboard/courses/page.tsx#L1-L158) | Список зачисленных курсов («Мои курсы») с прогресс-барами |
| [`courses/[id]/page.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/dashboard/courses/%5Bid%5D/page.tsx#L1-L448) | Экран курса: модули, уроки, статусы выполнения и переход к изучению |
| [`settings/page.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/dashboard/settings/page.tsx#L1-L335) | Личный кабинет: управление профилем, смена темы и параметры аккаунта |

---

## ⚙️ Функции, компоненты, хуки и API
| Компонент / Метод | Файл:Строки | Описание | Маршрут / Сигнатура |
|---|---|---|---|
| [`DashboardLayout`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/dashboard/layout.tsx#L6-L17) | [`layout.tsx#L6-L17`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/dashboard/layout.tsx#L6-L17) | Обертка с проверкой авторизации и навигационным меню | Макет `/dashboard/*` |
| [`StudentDashboard`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/dashboard/page.tsx#L8-L63) | [`page.tsx#L8-L63`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/dashboard/page.tsx#L8-L63) | Обзор активного обучения студента | Маршрут `/dashboard` |
| [`MyCoursesPage`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/dashboard/courses/page.tsx#L21-L157) | [`courses/page.tsx#L21-L157`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/dashboard/courses/page.tsx#L21-L157) | Карточки курсов со статусом прохождения | Маршрут `/dashboard/courses` |
| [`StudentCoursePlayerPage`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/dashboard/courses/%5Bid%5D/page.tsx#L26-L447) | [`courses/[id]/page.tsx#L26-L447`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/dashboard/courses/%5Bid%5D/page.tsx#L26-L447) | Детальная страница курса с аккордеоном модулей | Маршрут `/dashboard/courses/[id]` |
| `toggleSection` | [`courses/[id]/page.tsx#L91-L93`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/dashboard/courses/%5Bid%5D/page.tsx#L91-L93) | Сворачивание и разворачивание модулей курса | `(secId: number) => void` |
| [`ProfileAndSettingsPage`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/dashboard/settings/page.tsx#L23-L334) | [`settings/page.tsx#L23-L334`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/dashboard/settings/page.tsx#L23-L334) | Страница профиля и настроек | Маршрут `/dashboard/settings` |
| `handleSaveProfile` | [`settings/page.tsx#L47-L69`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/dashboard/settings/page.tsx#L47-L69) | Отправка `PATCH /auth/profile` и обновление Zustand | `(e: React.FormEvent) => Promise<void>` |

---

## 🔗 Зависимости
- **Входящие (кто вызывает этот модуль):**
  - Пользователи после логина, по ссылкам из меню `Sidebar` («Мои курсы», «Личный кабинет», «Настройки»).
- **Исходящие (что импортирует этот модуль):**
  - [`store/useAuth.ts`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/store/useAuth.ts#L33) — пользовательские данные и метод `setUser()`.
  - [`lib/api.ts`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/api.ts#L5) — запросы к `/courses/my`, `/courses/{id}/structure`, `/courses/{id}/progress`, `/auth/profile`.
  - [`components/CourseCard.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/CourseCard.tsx#L14), [`components/ProtectedRoute.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/ProtectedRoute.tsx#L7).

---

## 🧭 Навигатор типовых задач (Where to edit?)
- **Изменить расчет доступности следующего урока:** логика поиска следующего незавершенного урока находится в [`courses/[id]/page.tsx#L95-L100`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/dashboard/courses/%5Bid%5D/page.tsx#L95-L100).
- **Добавить новые настройки профиля:** расширить форму и тело запроса `PATCH /auth/profile` в [`settings/page.tsx#L51-L56`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/dashboard/settings/page.tsx#L51-L56).
- **Скорректировать отображение бейджей уроков (лекция, квиз, практика):** редактировать разметку уроков в [`courses/[id]/page.tsx#L320-L400`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/dashboard/courses/%5Bid%5D/page.tsx#L320-L400).
