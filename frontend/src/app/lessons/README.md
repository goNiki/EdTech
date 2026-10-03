# 📦 Модуль: `lessons`
> **Путь:** `frontend/src/app/lessons`  
> **Роль:** Страница интерактивного прохождения урока студентом: полноэкранный плеер материалов Puck, фиксация старта и отправка результатов автопроверки и эссе на бэкенд.

---

## 🎯 Назначение и ответственность
Модуль реализует учебную среду студента в режиме полного фокуса (Zen Mode). Он отвечает за:
1. Инициализацию прохождения (`POST /lessons/{id}/start`) и загрузку материалов урока (`GET /lessons/{id}`).
2. Отображение интерактивного плеера `PuckLessonViewer`, поддерживающего текстовые лекции, видео, квизы, упражнения на сопоставление и пропуски.
3. Отправку итогового скоринга и эссе на сервер (`POST /lessons/{id}/complete`) с пересчетом процента освоения курса.

---

## ⚠️ Жестко установленные правила (Hard Invariants)
1. **Двухфазный протокол прохождения урока:**
   - **Фаза 1 (Старт):** При открытии страницы вызывается `POST /lessons/${id}/start`, переводящий статус урока в `in_progress` в БД.
   - **Фаза 2 (Завершение):** При клике на кнопку «Завершить урок» или автоматическом завершении в плеере вызывается `POST /lessons/${id}/complete`.
2. **Контракт отправки результатов (LessonCompletionPayload):** В теле запроса к `/lessons/{id}/complete` обязаны передаваться поля `score` (числовой балл 0-100) и `essays` (массив `{ question_text, answer_text, max_points }`). Для чисто теоретических лекций дефолтный балл составляет `100`.
3. **Фокусный интерфейс без сайдбара:** Страница плеера урока намеренно не использует глобальный `Sidebar`, разворачиваясь на 100% ширины экрана с закрепленным хедером возврата к курсу (`/dashboard/courses/${courseId}`).
4. **Защита доступа:** Маршрут изолирован компонентом `<ProtectedRoute allowedRoles={['student', 'teacher', 'author', 'admin']}>`.

---

## 📁 Структура файлов модуля
| Файл | Описание роли файла |
|---|---|
| [`[id]/page.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/lessons/%5Bid%5D/page.tsx#L1-L148) | Полноэкранный плеер урока: верхний бар управления, рендер PuckLessonViewer и тосты завершения |

---

## ⚙️ Функции, компоненты, хуки и API
| Компонент / Метод | Файл:Строки | Описание | Маршрут / Сигнатура |
|---|---|---|---|
| [`LessonPlayer`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/lessons/%5Bid%5D/page.tsx#L10-L147) | [`[id]/page.tsx#L10-L147`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/lessons/%5Bid%5D/page.tsx#L10-L147) | Главный компонент экрана урока | Динамический роут `/lessons/[id]` |
| `fetchAndStartLesson` | [`[id]/page.tsx#L24-L40`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/lessons/%5Bid%5D/page.tsx#L24-L40) | Запуск урока через `POST /start` и загрузка контента `GET /lessons/{id}` | `() => Promise<void>` |
| `handleComplete` | [`[id]/page.tsx#L42-L54`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/lessons/%5Bid%5D/page.tsx#L42-L54) | Сохранение результатов и эссе через `POST /lessons/{id}/complete` | `(payload?: LessonCompletionPayload) => Promise<void>` |

---

## 🔗 Зависимости
- **Входящие (кто вызывает этот модуль):**
  - Экран программы курса [`dashboard/courses/[id]`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/dashboard/courses/%5Bid%5D/page.tsx#L410).
  - Ссылки на уроки из лендинга курса [`courses/[slug]`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/courses/%5Bslug%5D/page.tsx).
- **Исходящие (что импортирует этот модуль):**
  - [`components/player/PuckLessonViewer.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/player/PuckLessonViewer.tsx#L39) — плеер интерактивного контента.
  - [`components/ProtectedRoute.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/ProtectedRoute.tsx#L7) — проверка прав доступа.
  - [`lib/api.ts`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/api.ts#L5) — Axios-клиент.

---

## 🧭 Навигатор типовых задач (Where to edit?)
- **Изменить структуру полезной нагрузки при сдаче урока:** редактировать метод `handleComplete` в [`[id]/page.tsx#L42-L54`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/lessons/%5Bid%5D/page.tsx#L42-L54).
- **Скорректировать верхнюю панель плеера:** редактировать разметку `<header>` в [`[id]/page.tsx#L73-L109`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/lessons/%5Bid%5D/page.tsx#L73-L109).
- **Изменить поведение кнопки возврата к курсу:** изменить обработчики кнопок в [`[id]/page.tsx#L76-L81`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/lessons/%5Bid%5D/page.tsx#L76-L81) и [`[id]/page.tsx#L126-L131`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/lessons/%5Bid%5D/page.tsx#L126-L131).
