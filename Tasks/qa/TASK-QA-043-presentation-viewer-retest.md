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
- **Статус задачи:** Ready for QA ⏳
