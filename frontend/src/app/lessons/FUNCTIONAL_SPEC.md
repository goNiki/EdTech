# 📋 Функциональная спецификация: Интерактивный плеер урока (Lessons)

> **Расположение:** `frontend/src/app/lessons`  
> **Технический контекст:** Экран полноэкранного прохождения урока (`/lessons/[id]`) с фиксацией начала сессии, рендером Puck-блоков и отправкой итоговых результатов на сервер.  
> **Дата актуализации:** 2026-10-04  

---

## 🎯 Каталог бизнес-фич модуля

| Фича | Описание возможности | Обеспечивающие функции/компоненты |
|---|---|---|
| **Инициализация сессии урока (Start Lesson)** | Автоматическая регистрация старта прохождения урока на бэкенде через `POST /lessons/{id}/start` при загрузке страницы. | [`LessonPlayer.fetchAndStartLesson()`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/lessons/%5Bid%5D/page.tsx#L24-L40) |
| **Загрузка и рендеринг контента Puck** | Получение структурированного JSON-контента урока через `GET /lessons/{id}` и передача его в клиентский плеер `PuckLessonViewer`. | [`LessonPlayer.fetchAndStartLesson()`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/lessons/%5Bid%5D/page.tsx#L29-L33) |
| **Фиксация завершения урока (Complete Lesson)** | Сохранение итогового процента выполнения тестов и отправка развернутых ответов (эссе) на бэкенд через `POST /lessons/{id}/complete`. | [`LessonPlayer.handleComplete()`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/lessons/%5Bid%5D/page.tsx#L42-L54) |
| **Бесшовная навигация к программе курса** | Контекстный возврат в учебный план курса `/dashboard/courses/{course_id}` по завершении урока или клику на стрелку «Назад». | [`LessonPlayer (навигация)`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/lessons/%5Bid%5D/page.tsx#L75-L84) |

---

## 🔬 Паспорта функций и компонентов

### ⚡ Компонент: `LessonPlayer`

* **Файл и строки:** [`lessons/[id]/page.tsx#L10-L147`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/lessons/%5Bid%5D/page.tsx#L10-L147)
* **Бизнес-назначение:** Предоставляет студенту изолированное окружение для фокусированного прохождения урока и сдачи контрольных заданий.
* **Связанная фича:** Загрузка и рендеринг контента Puck

#### 📥 Входные параметры
| Параметр | Тип | Обязателен | Бизнес-смысл и ограничения |
|---|---|:---:|---|
| `params.id` | `string` | Да | Идентификатор урока в системе. |

#### 🔄 Пошаговый алгоритм работы
1. **Шаг 1 (RBAC валидация):** Компонент обернут в `<ProtectedRoute allowedRoles={['student', 'teacher', 'author', 'admin']}>`. Гости перенаправляются на `/login`.
2. **Шаг 2 (Оповещение о начале прохождения):** Сразу после монтирования вызывается `POST /lessons/${id}/start`. Ошибка вызова не блокирует загрузку контента (`catch(() => {})`).
3. **Шаг 3 (Получение контента):** Выполняется `GET /lessons/${id}`. Извлекается объект `data.lesson` с полями `title`, `content` (JSON-строка Puck), `course_id`.
4. **Шаг 4 (Отрисовка плеера):** Передача `lessonData.content` в компонент `<PuckLessonViewer>`.
5. **Шаг 5 (Сдача урока):** Клик по верхней кнопке «Завершить урок» или кнопке в конце теста активирует функцию `handleComplete(payload)`.

#### ⚠️ Побочные эффекты (Side Effects)
* **Сетевые вызовы:** `POST /lessons/{id}/start`, `GET /lessons/{id}`, `POST /lessons/{id}/complete`.
* **State:** Мутация `isCompleted: true`, показ Toast уведомления.

#### 💡 Подсказка для аналитика (Где менять логику?)
* *Смена роута возврата:* [`lessons/[id]/page.tsx#L77-L78`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/lessons/%5Bid%5D/page.tsx#L77-L78).

---

### ⚡ Функция: `handleComplete(payload)`

* **Файл и строки:** [`lessons/[id]/page.tsx#L42-L54`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/lessons/%5Bid%5D/page.tsx#L42-L54)
* **Бизнес-назначение:** Отправляет на бэкенд итоговый расчет баллов студента за данный урок, а также развернутые ответы на эссе и практические задания для ручной проверки учителем.
* **Связанная фича:** Фиксация завершения урока (Complete Lesson)

#### 📥 Входные параметры
| Параметр | Тип | Обязателен | Бизнес-смысл и ограничения |
|---|---|:---:|---|
| `payload` | `LessonCompletionPayload` | Нет | Содержит вычисленный `score` (0–100%) и массив `essays: [{ question_text, answer_text, max_points }]`. |

#### 🔄 Пошаговый алгоритм работы
1. **Шаг 1 (Подготовка пакета данных):** Формируется DTO:
   ```json
   {
     "score": payload?.score ?? 100,
     "essays": payload?.essays ?? []
   }
   ```
2. **Шаг 2 (Сетевой запрос):** Выполняется `POST /lessons/${id}/complete`.
3. **Шаг 3 (Фиксация успешного прохождения):**
   - Устанавливается `isCompleted = true`.
   - Отображается всплывающее Toast-уведомление «Урок успешно завершен! Прогресс и задания сохранены.».
   - Кнопка в верхней панели блокируется и меняет текст на «Урок завершен ✓».

#### ⚠️ Побочные эффекты (Side Effects)
* **БД / Сервер:** Обновление статуса урока студента на `'completed'`, сохранение баллов, добавление эссе в очередь проверки преподавателя.
* **Toast Notification:** Зеленое всплывающее окно в правом нижнем углу.

#### 💡 Подсказка для аналитика (Где менять логику?)
* *Формат отправки результатов урока:* [`lessons/[id]/page.tsx#L44-L47`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/lessons/%5Bid%5D/page.tsx#L44-L47).
