# 📦 Модуль: `lessons`
> **Путь:** `frontend/src/app/lessons`  
> **Роль:** Страница интерактивного прохождения урока студентом: полноэкранный плеер материалов Puck, фиксация старта и отправка результатов автопроверки и эссе на бэкенд с защитой от списывания (античит).

---

## 🎯 Назначение и ответственность
Модуль реализует учебную среду студента в режиме полного фокуса (Zen Mode). Он отвечает за:
1. Инициализацию прохождения (`POST /lessons/{id}/start`) и загрузку материалов урока (`GET /lessons/{id}`).
2. Отображение интерактивного плеера `PuckLessonViewer`, поддерживающего текстовые лекции, видео, квизы, упражнения на сопоставление, пропуски и загрузку домашних заданий.
3. Отправку собранной структуры выбранных ответов (`POST /lessons/{id}/complete`) на сервер для античит-проверки и пересчета подтвержденного балла.

---

## ⚠️ Жестко установленные правила (Hard Invariants)
1. **Двухфазный протокол прохождения урока:**
   - **Фаза 1 (Старт):** При открытии страницы вызывается `POST /lessons/${id}/start`, переводящий статус урока в `in_progress` в БД.
   - **Фаза 2 (Завершение):** При клике на кнопку «Завершить урок» или подтверждении в плеере вызывается `POST /lessons/${id}/complete`.
2. **Контракт отправки результатов (LessonCompletionPayload):** В теле запроса к `/lessons/{id}/complete` обязаны передаваться поля `score` (числовой балл 0-100), `answers` (массив `{ block_id, answer }`) и `essays` (массив `{ question_text, answer_text, max_points }`). Для чисто теоретических лекций дефолтный балл составляет `100`.
3. **Фокусный интерфейс без сайдбара:** Страница плеера урока намеренно не использует глобальный `Sidebar`, разворачиваясь на 100% ширины экрана с закрепленным хедером возврата к курсу (`/dashboard/courses/${courseId}`).
4. **Защита доступа:** Маршрут изолирован компонентом `<ProtectedRoute allowedRoles={['student', 'teacher', 'author', 'admin']}>`.
5. **Серверная античит-проверка квизов (Server Anti-Cheat Verification):** Плеер `PuckLessonViewer` отправляет структуру всех выбранных студентом вариантов на сервер. Итоговая оценка фиксируется по подтвержденному ответу бэкенда (`serverData.score`), а результаты валидации блоков подсвечиваются зеленым/красным цветом с бейджами «Сервер: Верно / Неверно».

---

## 📁 Структура файлов модуля
| Файл | Описание роли файла |
|---|---|
| [`[id]/page.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/lessons/%5Bid%5D/page.tsx) | Полноэкранный плеер урока: верхний бар управления, рендер PuckLessonViewer, баннер предыдущих попыток и тосты завершения |

---

## ⚙️ Функции, компоненты, хуки и API
| Компонент / Метод | Файл:Строки | Описание | Маршрут / Сигнатура |
|---|---|---|---|
| [`LessonPlayer`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/lessons/%5Bid%5D/page.tsx) | [`[id]/page.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/lessons/%5Bid%5D/page.tsx) | Главный компонент экрана урока | Динамический роут `/lessons/[id]` |
| `fetchLessonAndProgress` | [`[id]/page.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/lessons/%5Bid%5D/page.tsx) | Запуск урока через `POST /start`, загрузка контента `GET /lessons/{id}` и прогресса `GET /progress` | `() => Promise<void>` |
| `handleComplete` | [`[id]/page.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/lessons/%5Bid%5D/page.tsx) | Серверная валидация квизов и сохранение эссе через `POST /lessons/{id}/complete` | `(payload?: LessonCompletionPayload) => Promise<any>` |

---

## 🔗 Зависимости
- **Входящие (кто вызывает этот модуль):**
  - Экран программы курса [`dashboard/courses/[id]`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/dashboard/courses/%5Bid%5D/page.tsx).
  - Ссылки на уроки из лендинга курса [`courses/[slug]`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/courses/%5Bslug%5D/page.tsx).
- **Исходящие (что импортирует этот модуль):**
  - [`components/player/PuckLessonViewer.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/player/PuckLessonViewer.tsx) — плеер интерактивного контента с серверной проверкой.
  - [`components/ProtectedRoute.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/ProtectedRoute.tsx) — проверка прав доступа.
  - [`lib/api.ts`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/api.ts) — Axios-клиент.

---

## 🧭 Навигатор типовых задач (Where to edit?)
- **Изменить структуру полезной нагрузки при сдаче урока:** редактировать метод `handleComplete` в [`[id]/page.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/lessons/%5Bid%5D/page.tsx).
- **Скорректировать верхнюю панель плеера:** редактировать разметку `<header>` в [`[id]/page.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/lessons/%5Bid%5D/page.tsx).
- **Изменить поведение кнопки возврата к курсу:** изменить обработчики кнопок в [`[id]/page.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/lessons/%5Bid%5D/page.tsx).
