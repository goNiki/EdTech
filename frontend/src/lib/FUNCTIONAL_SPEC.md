# 📋 Функциональная спецификация: Библиотеки ядра и конфигурация контента (Lib)

> **Расположение:** `frontend/src/lib`  
> **Технический контекст:** Транспортный слой API-клиента (Axios с ротацией JWT) и ядро декларативной конфигурации визуального редактора контента Puck Editor.  
> **Дата актуализации:** 2026-10-04  

---

## 🎯 Каталог бизнес-фич модуля

| Фича | Описание возможности | Обеспечивающие функции |
|---|---|---|
| **Авто-инъекция JWT в HTTP-запросы** | Автоматическое добавление заголовка авторизации `Authorization: Bearer <token>` во все исходящие сетевые запросы при наличии активной пользовательской сессии. | [`api.interceptors.request`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/api.ts#L12-L20) |
| **Прозрачная ротация токенов при 401** | Бесшовное фоновое обновление просроченного `access_token` через эндпоинт `/auth/refresh` без разлогинивания пользователя и повторение исходного запроса. | [`api.interceptors.response`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/api.ts#L22-L51) |
| **Аварийный логаут и редирект на логин** | Принудительная очистка невалидных токенов из `localStorage` и перенаправление пользователя на `/login` при сбое обновления сессии или отсутствии refresh-токена. | [`api.interceptors.response (catch)`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/api.ts#L42-L47) |
| **Парсер шаблонов интерактивных пропусков** | Извлечение интерактивных зон из разметки `{Правильный; Дистрактор1, Дистрактор2}` и `{Слово}`, разбиение на токены, разделение правильных ответов и дистракторов. | [`parseSmartDropdownTemplate()`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/puck-config.tsx#L105-L176) |
| **Токенизация текстовых шаблонов** | Базовое разбиение произвольного текста с заполнителями `{placeholder}` на последовательность текстовых и пустых токенов. | [`parseTemplateParts()`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/puck-config.tsx#L178-L209) |
| **Декларативный реестр блоков Puck** | Спецификация 12 типов учебных блоков (теория, видео, квизы 8 видов, загрузка ДЗ) с поддержкой встроенного инлайн-редактирования (InlineEditable). | [`config: Config<PuckProps>`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/puck-config.tsx#L211-L1369) |
| **Утилита слияния стилей интерфейса** | Конкатенация и разрешение конфликтов Tailwind CSS классов. | [`cn()`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/utils.ts#L4-L6) |

---

## 🔬 Паспорта функций и компонентов

### ⚡ Функция: `api.interceptors.request`

* **Файл и строки:** [`api.ts#L12-L20`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/api.ts#L12-L20)
* **Бизнес-назначение:** Обеспечивает сквозную авторизацию всех обращений к бэкенду. Пользователю не требуется повторно передавать учетные данные при каждом действии.
* **Связанная фича:** Авто-инъекция JWT в HTTP-запросы

#### 📥 Входные параметры
| Параметр | Тип | Обязателен | Бизнес-смысл и ограничения |
|---|---|:---:|---|
| `config` | `InternalAxiosRequestConfig` | Да | Конфигурация исходящего HTTP-запроса (URL, метод, заголовки, тело). |

#### 🔄 Пошаговый алгоритм работы
1. **Шаг 1 (Проверка среды исполнения):** Проверяется условие `typeof window !== 'undefined'`, исключающее попытки обращения к браузерному `localStorage` на этапе SSR (Server-Side Rendering).
2. **Шаг 2 (Извлечение токена):** Считывается ключ `access_token` из `localStorage`.
3. **Шаг 3 (Инъекция заголовка):** При наличии токена в `config.headers.Authorization` устанавливается строка формата `Bearer ${token}`.
4. **Шаг 4 (Возврат конфигурации):** Возвращается модифицированный объект конфигурации запроса.

#### ⚠️ Побочные эффекты (Side Effects)
* **Заголовки запроса:** Модификация HTTP-заголовка `Authorization`.

#### 📤 Результат и ошибки
* **Успешный результат:** `config` с проставленным заголовком авторизации.

#### 💡 Подсказка для аналитика (Где менять логику?)
* *Смена префикса или формата токена:* [`api.ts#L16-L17`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/api.ts#L16-L17).

---

### ⚡ Функция: `api.interceptors.response`

* **Файл и строки:** [`api.ts#L22-L51`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/api.ts#L22-L51)
* **Бизнес-назначение:** Предотвращает внезапное прерывание пользовательской работы при истечении срока действия access-токена (например, в середине прохождения теста или редактирования урока).
* **Связанная фича:** Прозрачная ротация токенов при 401 и Аварийный логаут

#### 📥 Входные параметры
| Параметр | Тип | Обязателен | Бизнес-смысл и ограничения |
|---|---|:---:|---|
| `error` | `AxiosError` | Да | Объект сетевой или серверной ошибки, содержащий `status`, `config` и `response`. |

#### 🔄 Пошаговый алгоритм работы
1. **Шаг 1 (Проверка критериев обновления):** Запрос перехватывается, если:
   - Статус ответа строго равен `401 Unauthorized`.
   - Запрос еще не повторялся ранее (`!originalRequest._retry`).
   - Запрос не являлся запросом на обновление токена (`originalRequest.url !== '/auth/refresh'`).
2. **Шаг 2 (Защита от зацикливания):** Выставляется флаг повтора `originalRequest._retry = true`.
3. **Шаг 3 (Чтение refresh-токена):** Считывается `refresh_token` из `localStorage`. Если он отсутствует, генерируется исключение.
4. **Шаг 4 (Запрос нового токена):** Через немодифицированный экземпляр `axios.post` отправляется запрос на `${API_URL}/auth/refresh` с телом `{ refresh_token }`.
5. **Шаг 5 (Сохранение новых ключей):**
   - Новый `data.access_token` сохраняется в `localStorage`.
   - Если сервер вернул новый `data.refresh_token`, он также обновляется в хранилище.
6. **Шаг 6 (Повтор исходной операции):** В `originalRequest.headers.Authorization` подставляется новый токен, после чего запрос повторно отправляется через `api(originalRequest)`.
7. **Шаг 7 (Аварийный сброс при неудаче):** Если ротация токена провалилась (просрочен refresh-токен или невалиден):
   - Очищаются ключи `access_token` и `refresh_token` из `localStorage`.
   - Выполняется принудительный переход `window.location.href = '/login'`.
   - Промис реджектится с ошибкой `refreshError`.

#### ⚠️ Побочные эффекты (Side Effects)
* **Web Storage:** Запись новых токенов либо их полное удаление.
* **Навигация:** Редирект браузера на страницу входа `/login`.
* **Сетевой запрос:** Отправка POST `/auth/refresh` и повтор оригинального запроса.

#### 📤 Результат и ошибки
* **Успешный результат:** Прозрачный возврат результата повторенного запроса.
* **Возможные ошибки:** `refreshError` при невалидном refresh-токене (завершается редиректом).

#### 💡 Подсказка для аналитика (Где менять логику?)
* *URL эндпоинта обновления токенов:* [`api.ts#L33`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/api.ts#L33).
* *Поведение при истечении сессии (например, открытие модального окна входа вместо жесткого редиректа):* [`api.ts#L43-L46`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/api.ts#L43-L46).

---

### ⚡ Функция: `parseSmartDropdownTemplate(templateText, existingBlanks)`

* **Файл и строки:** [`puck-config.tsx#L105-L176`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/puck-config.tsx#L105-L176)
* **Бизнес-назначение:** Позволяет авторам курсов быстро создавать интерактивные задания с пропусками в формате текста (Inline DSL), автоматически формируя выпадающие списки с правильным ответом и дистракторами (ложными вариантами).
* **Связанная фича:** Парсер шаблонов интерактивных пропусков

#### 📥 Входные параметры
| Параметр | Тип | Обязателен | Бизнес-смысл и ограничения |
|---|---|:---:|---|
| `templateText` | `string` | Да | Текст с синтаксисом `{Правильный; Дистрактор1, Дистрактор2}` или `{Ответ}`. |
| `existingBlanks` | `Array<DropdownBlankItem \| InputBlankItem>` | Нет | Ранее сохраненные пропуски для обратной совместимости. |

#### 🔄 Пошаговый алгоритм работы
1. **Шаг 1 (Проверка входных данных):** Если строка пустая, возвращается `{ tokens: [], blanks: [] }`.
2. **Шаг 2 (Поиск вхождений):** Текст сканируется регулярным выражением `/\{([^{}]+)\}/g`.
3. **Шаг 3 (Формирование текстовых токенов):** Символы до найденной фигурной скобки добавляются как токен `{ type: 'text', value: ... }`.
4. **Шаг 4 (Парсинг структуры пропуска):**
   - Если внутри скобок присутствует символ `;`, строка разделяется на две части: до точки с запятой — `correctAnswer`, после точки с запятой — список дистракторов через запятую `,`.
   - Если точки с запятой нет, вся строка внутри скобок признается `correctAnswer`.
5. **Шаг 5 (Сборка уникальных опций):**
   - Генерируется уникальный ключ поля `blank_1`, `blank_2`, ...
   - Формируется дедуплицированный массив опций: `[correctAnswer, ...distractors]`.
6. **Шаг 6 (Формирование токена пропуска):** Создается объект `ParsedDropdownBlank` и пушится в массивы `tokens` (тип `'blank'`) и `blanks`.
7. **Шаг 7 (Хвостовой текст):** Текст после последнего совпадения добавляется в `tokens` с типом `'text'`.

#### ⚠️ Побочные эффекты (Side Effects)
* Чистая функция без побочных эффектов.

#### 📤 Результат и ошибки
* **Успешный результат:** Объект `{ tokens: ParsedTemplateToken[], blanks: ParsedDropdownBlank[] }`.
* **Пример результата:**
  ```json
  {
    "tokens": [
      { "type": "text", "value": "Протокол " },
      { "type": "blank", "value": "{gRPC; REST, SOAP}", "blank": { "key": "blank_1", "correctAnswer": "gRPC", "options": ["gRPC", "REST", "SOAP"] } }
    ],
    "blanks": [...]
  }
  ```

#### 💡 Подсказка для аналитика (Где менять логику?)
* *Синтаксис разделителя правильного ответа и дистракторов (сейчас `;` и `,`):* [`puck-config.tsx#L135-L144`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/puck-config.tsx#L135-L144).

---

### ⚡ Объект конфигурации: `config` (Puck Editor Component Registry)

* **Файл и строки:** [`puck-config.tsx#L211-L1369`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/puck-config.tsx#L211-L1369)
* **Бизнес-назначение:** Центральный реестр компонентов визуального редактора уроков. Задает схему полей формы (Sidebar inspector) и способ интерактивного рендеринга на канвасе в режиме редактирования.
* **Связанная фича:** Декларативный реестр блоков Puck

#### Состав зарегистрированных блоков:
1. `HeaderBlock` ([L213-L268](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/puck-config.tsx#L213-L268)): Заголовки разделов H1, H2, H3 с подзаголовком.
2. `TextBlock` ([L270-L300](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/puck-config.tsx#L270-L300)): Базовый текст лекции с поддержкой Markdown.
3. `RichTextBlock` ([L302-L344](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/puck-config.tsx#L302-L344)): Полнофункциональный Word-подобный WYSIWYG редактор статьи (форматирование, списки, врезки советов).
4. `VideoBlock` ([L346-L383](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/puck-config.tsx#L346-L383)): Встроенный видеоплеер YouTube / MP4 с подписью.
5. `QuizSingleBlock` ([L385-L545](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/puck-config.tsx#L385-L545)): Тест с 1 правильным ответом, динамическим добавлением/удалением вариантов прямо на холсте.
6. `QuizMultiBlock` ([L547-L706](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/puck-config.tsx#L547-L706)): Тест с множественным выбором правильных ответов.
7. `QuizMatchBlock` ([L708-L834](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/puck-config.tsx#L708-L834)): Задание на сопоставление терминов и определений из двух колонок.
8. `QuizDropdownBlankBlock` ([L836-L997](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/puck-config.tsx#L836-L997)): Пропуски в предложении с выбором из выпадающего меню.
9. `QuizInputBlankBlock` ([L999-L1121](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/puck-config.tsx#L999-L1121)): Пропуски с клавиатурным вводом точного термина.
10. `QuizSequenceBlock` ([L1123-L1240](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/puck-config.tsx#L1123-L1240)): Задание на расстановку этапов алгоритма в правильном хронологическом порядке.
11. `QuizEssayBlock` ([L1242-L1298](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/puck-config.tsx#L1242-L1298)): Открытый вопрос с рубрикой критериев для последующей ручной проверки преподавателем.
12. `FileUploadBlock` ([L1300-L1367](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/puck-config.tsx#L1300-L1367)): Практическое домашнее задание с прикреплением файла (zip, pdf, код).

#### 💡 Подсказка для аналитика (Где менять логику?)
* *Дефолтные значения баллов и вопросов:* поля `defaultProps` у каждого блока в [`puck-config.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/puck-config.tsx).
