# 🎯 [TASK-FE-043] Фикс нормализации Canva URL и валидация внешних ссылок в PresentationViewer

## 📌 Описание проблемы
В ходе выполнения QA-тестирования `[TASK-QA-016]` в компоненте `PresentationViewer.tsx` обнаружены дефекты обработки ссылок на презентации:

1. **Дублирование `/view` в URL Canva:**
   Если автор вставляет стандартную ссылку общего доступа Canva вида:
   `https://www.canva.com/design/DAFxxxxxx/view`
   функция `normalizePresentationUrl` делит строку по `?` (`base = trimmed.split('?')[0]`) и конкатенирует `${base}/view?embed`, в результате чего формируется некорректный URL:
   `https://www.canva.com/design/DAFxxxxxx/view/view?embed`, что приводит к ошибке 404 в iframe.

2. **Отсутствие валидации невалидных доменов (Кейс 2.1):**
   При вводе ссылок сторонних сервисов (например, `https://google.com`), плеер пытается вставить их в `iframe`, что приводит к CSP / X-Frame-Options ошибке браузера без информативного сообщения пользователю.

## 🛠️ Требования к реализации
1. В `frontend/src/components/player/PresentationViewer.tsx`:
   - Скорректировать нормализацию Canva: удалять суффикс `/view` перед добавлением `/view?embed`:
     ```typescript
     if (trimmed.includes('canva.com/design') && !trimmed.includes('view?embed')) {
       let base = trimmed.split('?')[0];
       base = base.replace(/\/view\/?$/, '');
       return `${base}/view?embed`;
     }
     ```
   - Добавить валидацию ссылки на известные поддерживаемые провайдеры (Google Slides, Canva, SpeakerDeck, OneDrive/Office 365).
   - Если ссылка не совпадает с поддерживаемыми сервисами или не является валидным embed URL, выводить в плеере предупреждение:
     *«Некорректная ссылка на презентацию. Поддерживаются Google Slides, Canva, SpeakerDeck, Office 365»*.

## 🧪 Требование к тестированию
После реализации изменений разработчик обязан создать встречную задачу в папку `Tasks/qa/` (например, `TASK-QA-043-presentation-viewer-retest.md`) с описанием исправлений для повторного ретеста.

---

## 🏁 Статус выполнения
- **Статус:** Completed ✅
- **Реализация:** В `frontend/src/components/player/PresentationViewer.tsx` исправлена нормализация Canva URL (удаление завершающего `/view` перед формированием `${base}/view?embed`), добавлена поддержка Office 365 / OneDrive, создана экспортируемая функция валидации `isSupportedPresentationUrl(url)`. Для невалидных доменов (например, `google.com`) отображается информативный блок с `AlertCircle` вместо сломанного iframe.
- **Встречная QA задача:** Создана в [`Tasks/qa/TASK-QA-043-presentation-viewer-retest.md`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/Tasks/qa/TASK-QA-043-presentation-viewer-retest.md).
