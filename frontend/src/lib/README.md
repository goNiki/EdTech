# 📦 Модуль: `lib`
> **Путь:** `frontend/src/lib`  
> **Роль:** Ядро инфраструктурных сервисов фронтенда: сетевой HTTP-клиент Axios с циклом авто-ротации токенов, декларативная конфигурация интерактивных блоков Puck Editor и утилиты стилей.

---

## 🎯 Назначение и ответственность
Модуль `lib` объединяет низкоуровневые сервисы клиентского приложения:
1. Сетевой транспорт (`api.ts`): инкапсулирует базовый URL бэкенда, инжектит `Bearer` JWT в каждый запрос, перехватывает ошибки `401 Unauthorized`, прозрачно обновляет токен через `POST /auth/refresh` и возобновляет упавший запрос, либо направляет на экран логина.
2. Движок визуального редактора (`puck-config.tsx`): описывает спецификацию 12 блоков Puck в архитектуре **Content-as-Data**, парсеры шаблонов пропусков (`parseSmartDropdownTemplate`), а также интерактивные формы ввода и редактирования.
3. Утилиты Tailwind (`utils.ts`): хелпер объединения классов `cn` на базе `clsx` и `tailwind-merge`.

---

## ⚠️ Жестко установленные правила (Hard Invariants)
1. **Защита от бесконечного цикла ротации токенов (Infinite Loop Prevention):** В перехватчике ответов `api.interceptors.response` обязательна проверка `!originalRequest._retry` и `originalRequest.url !== '/auth/refresh'`. При получении повторной ошибки 401 запрос не должен зацикливаться — токены удаляются из `localStorage`, а `window.location.href` сбрасывается на `/login`.
2. **Синтаксис шаблонов пропусков Smart Blanks:** Формат текста с выпадающими списками в интерактивных упражнениях строго фиксирован: `{ПравильныйОтвет;Дистрактор1,Дистрактор2}`. Встроенный парсер `parseSmartDropdownTemplate` парсит этот синтаксис регулярным выражением `/\{([^{}]+)\}/g`, автоматически формируя список перемешанных опций.
3. **Строгая типизация Content-as-Data:** Все 12 блоков Puck обязаны иметь зеркальные типы в структуре `PuckProps`. Никаких свободных неструктурированных JSON-полей внутри контента урока — любые интерактивные квизы (Single, Multi, Match, Dropdown, Input, Sequence, Essay, FileUpload) строго типизированы.
4. **Конфигурация API URL:** Клиент по умолчанию использует переменную окружения `NEXT_PUBLIC_API_URL`, а при ее отсутствии откатывается на `http://localhost:8082/api/v1`.

---

## 📁 Структура файлов модуля
| Файл | Описание роли файла |
|---|---|
| [`api.ts`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/api.ts#L1-L52) | Экземпляр Axios с интерцепторами авторизации и прозрачным refresh-циклом |
| [`puck-config.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/puck-config.tsx#L1-L1370) | Спецификация 12 интерактивных блоков Puck, типы данных и парсеры шаблонов |
| [`utils.ts`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/utils.ts#L1-L7) | Утилита объединения CSS-классов `cn` (`clsx` + `twMerge`) |

---

## ⚙️ Функции, компоненты, хуки и API

### Сетевой транспорт (`api.ts`)
| Экспорт / Перехватчик | Файл:Строки | Описание | Сигнатура / Поведение |
|---|---|---|---|
| `api` | [`api.ts#L5-L10`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/api.ts#L5-L10) | Базовый настроенный клиент Axios с заголовком `application/json` | `AxiosInstance` |
| `Request Interceptor` | [`api.ts#L12-L20`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/api.ts#L12-L20) | Инжекция заголовка `Authorization: Bearer <access_token>` из `localStorage` | `(config: InternalAxiosRequestConfig) => InternalAxiosRequestConfig` |
| `Response Interceptor` | [`api.ts#L22-L51`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/api.ts#L22-L51) | Автоматический перехват 401, запрос к `/auth/refresh`, сохранение новой пары токенов и ретрай | `(response) => response, async (error) => Promise<any>` |

### Движок Puck (`puck-config.tsx`)
| Функция / Конфигурация | Файл:Строки | Описание | Сигнатура / Вход и Выход |
|---|---|---|---|
| [`PuckProps`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/puck-config.tsx#L18-L90) | [`puck-config.tsx#L18-L90`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/puck-config.tsx#L18-L90) | Полная типовая карта всех 12 блоков урока | `type PuckProps = { HeaderBlock: ...; QuizSingleBlock: ...; ... }` |
| [`parseSmartDropdownTemplate`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/puck-config.tsx#L105-L176) | [`puck-config.tsx#L105-L176`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/puck-config.tsx#L105-L176) | Парсер шаблона с пропусками вида `{правильный;дистрактор1,дистрактор2}` | `(templateText: string, existingBlanks?) => { tokens, blanks }` |
| [`parseTemplateParts`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/puck-config.tsx#L178-L209) | [`puck-config.tsx#L178-L209`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/puck-config.tsx#L178-L209) | Разбиение строки на текстовые узлы и плейсхолдеры `{ключ}` | `(text: string) => Array<{ type, value, key }>` |
| [`config`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/puck-config.tsx#L211-L1370) | [`puck-config.tsx#L211-L1370`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/puck-config.tsx#L211-L1370) | Экземпляр конфигурации Puck Editor с полями и рендерами 12 блоков | `Config<PuckProps>` |

### Каталог 12 блоков Puck Editor:
1. `HeaderBlock` ([#L213-L268](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/puck-config.tsx#L213-L268)) — Заголовок уровня H1/H2/H3 с подзаголовком.
2. `TextBlock` ([#L270-L300](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/puck-config.tsx#L270-L300)) — Абзац текста с поддержкой Markdown.
3. `RichTextBlock` ([#L302-L344](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/puck-config.tsx#L302-L344)) — Полноценная статья лекции с форматированием (Word-like).
4. `VideoBlock` ([#L346-L383](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/puck-config.tsx#L346-L383)) — Видео-плеер по URL (YouTube, Vimeo, MP4) с подписью.
5. `QuizSingleBlock` ([#L385-L545](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/puck-config.tsx#L385-L545)) — Вопрос с одиночным выбором ответа и пояснением.
6. `QuizMultiBlock` ([#L547-L706](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/puck-config.tsx#L547-L706)) — Вопрос с несколькими правильными ответами.
7. `QuizMatchBlock` ([#L708-L834](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/puck-config.tsx#L708-L834)) — Задание на сопоставление пар терминов.
8. `QuizDropdownBlankBlock` ([#L836-L997](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/puck-config.tsx#L836-L997)) — Текст с пропусками и выпадающими списками вариантов.
9. `QuizInputBlankBlock` ([#L999-L1121](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/puck-config.tsx#L999-L1121)) — Текст с пропусками для ручного ввода правильного слова.
10. `QuizSequenceBlock` ([#L1123-L1240](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/puck-config.tsx#L1123-L1240)) — Упражнение на выстраивание элементов в верном хронологическом/логическом порядке.
11. `QuizEssayBlock` ([#L1242-L1298](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/puck-config.tsx#L1242-L1298)) — Открытый вопрос (эссе) с критериями оценивания для преподавателя.
12. `FileUploadBlock` ([#L1300-L1369](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/puck-config.tsx#L1300-L1369)) — Задание с загрузкой практического файла-решения студентом.

### Утилиты стилей (`utils.ts`)
| Функция | Файл:Строки | Описание | Сигнатура |
|---|---|---|---|
| [`cn`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/utils.ts#L4-L6) | [`utils.ts#L4-L6`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/utils.ts#L4-L6) | Объединение и разрешение конфликтов CSS-классов Tailwind | `(...inputs: ClassValue[]) => string` |

---

## 🔗 Зависимости
- **Входящие (кто вызывает этот модуль):**
  - Все страницы каталога, уроков, кабинетов студента и преподавателя импортируют `api` для сетевых запросов.
  - Редактор уроков [`edit/page.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/teacher/lessons/[id]/edit/page.tsx#L5) импортирует `config`.
  - Плеер студента [`PuckLessonViewer.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/player/PuckLessonViewer.tsx#L4) импортирует `parseSmartDropdownTemplate`.
  - Все компоненты Shadcn UI используют функцию `cn`.
- **Исходящие (что импортирует этот модуль):**
  - `axios` — HTTP-библиотека.
  - `@puckeditor/core` — ядро конструктора Puck.
  - `InlineEditable.tsx` — компоненты встроенного редактирования текста.
  - `lucide-react` — иконки интерфейса.

---

## 🧭 Навигатор типовых задач (Where to edit?)
- **Добавить новый интерактивный блок в конструктор уроков:**
  1. Описать его пропсы в `PuckProps` в [`puck-config.tsx#L18-L90`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/puck-config.tsx#L18-L90).
  2. Добавить конфигурацию полей и рендер в объект `config.components` в [`puck-config.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/puck-config.tsx#L211).
  3. Добавить обработку блока в интерактивный плеер [`PuckLessonViewer.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/player/PuckLessonViewer.tsx#L125).
  4. Добавить иконку и превью в сайдбар аутлайна [`CustomOutline.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/editor/CustomOutline.tsx#L41).
- **Изменить поведение refresh-токенов или таймауты запросов:** редактировать интерцепторы в [`api.ts#L12-L51`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/api.ts#L12-L51).
- **Сменить базовый адрес API бэкенда:** передать `NEXT_PUBLIC_API_URL` в `.env.local` или изменить fallback в [`api.ts#L3`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/api.ts#L3).
