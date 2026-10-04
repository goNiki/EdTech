# 🧪 [QA-043] Ретест нормализации Canva URL и валидации внешних доменов в PresentationViewer

> **Приоритет:** High (P1)  
> **Связанная задача разработки:** `[FE-043]`  
> **Компонент:** `frontend/src/components/player/PresentationViewer.tsx`  
> **Предыдущий дефект:** `[TASK-QA-016]`  

---

## 🎯 Цель тестирования
Провести повторное тестирование компонента `PresentationViewer` после исправления дефектов парсинга ссылок Canva и добавления валидации поддерживаемых провайдеров.

---

## 📋 Чек-лист проверок (Test Cases)

### Тест-кейс 1: Canva URL с суффиксом `/view`
- **Входной URL:** `https://www.canva.com/design/DAFxxxxxx/view`
- **Ожидаемый результат:** URL в `iframe` должен быть `https://www.canva.com/design/DAFxxxxxx/view?embed` без дублирования `/view/view?embed`.

### Тест-кейс 2: Canva URL без суффикса `/view`
- **Входной URL:** `https://www.canva.com/design/DAFxxxxxx`
- **Ожидаемый результат:** Преобразуется в `https://www.canva.com/design/DAFxxxxxx/view?embed`.

### Тест-кейс 3: Google Slides URL
- **Входной URL:** `https://docs.google.com/presentation/d/123/edit`
- **Ожидаемый результат:** Преобразуется в `https://docs.google.com/presentation/d/123/embed?start=false&loop=false&delayms=3000`.

### Тест-кейс 4: Неподдерживаемый внешний домен (Кейс 2.1 из QA-016)
- **Входной URL:** `https://google.com` или `https://example.com/slide`
- **Ожидаемый результат:** `iframe` не создается. Отображается дружелюбное предупреждение:
  *«Некорректная ссылка на презентацию. Поддерживаются Google Slides, Canva, SpeakerDeck, Office 365»*.

---

## 🏁 Статус
- **Статус задачи:** Passed (Closed) ✅

---

## 📋 Отчёт о результатах тестирования (QA Execution Report)

### 1. Результаты прогона тест-кейсов

| № | Тест-кейс | Входные данные | Ожидаемый результат | Фактический результат | Статус |
|---|---|---|---|---|:---:|
| 1 | Canva URL с суффиксом `/view` | `https://www.canva.com/design/DAFxxxxxx/view` | Преобразование в `https://www.canva.com/design/DAFxxxxxx/view?embed` без дублирования `/view/view?embed` | `https://www.canva.com/design/DAFxxxxxx/view?embed` | ✅ PASSED |
| 2 | Canva URL без суффикса `/view` | `https://www.canva.com/design/DAFxxxxxx` | Корректное добавление `/view?embed` | `https://www.canva.com/design/DAFxxxxxx/view?embed` | ✅ PASSED |
| 3 | Canva URL с замыкающим слешем `/view/` | `https://www.canva.com/design/DAFxxxxxx/view/` | Очистка слеша и суффикса перед подстановкой `?embed` | `https://www.canva.com/design/DAFxxxxxx/view?embed` | ✅ PASSED |
| 4 | Google Slides URL с `/edit` | `https://docs.google.com/presentation/d/123/edit` | Замена на `/embed?start=false&loop=false&delayms=3000` | `https://docs.google.com/presentation/d/123/embed?start=false&loop=false&delayms=3000` | ✅ PASSED |
| 5 | Неподдерживаемый внешний домен | `https://google.com`, `https://example.com/slide` | `isSupportedPresentationUrl` возвращает `false`, `iframe` блокируется, отображается алерт | Возвращает `false`, рендерится `<div className="text-rose-400">...<p>Некорректная ссылка на презентацию</p>...` | ✅ PASSED |

### 2. Сводка дефектов
- Дефект из `TASK-QA-016` (дублирование `view/view?embed` в Canva и отсутствие защиты от произвольных iframe доменов) успешно устранён фронтенд-разработчиком в задаче `TASK-FE-043`.
- Новых дефектов не обнаружено.

### 3. Вердикт
- Ретест пройден на 100%. Компонент `PresentationViewer.tsx` готов к эксплуатации. Задача закрыта.
