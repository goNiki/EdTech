# TASK-FE-020: Фикс инлайн-форматирования выделенного текста в Canvas и Word редакторах (Iframe Selection Fix)

## 📋 Описание проблемы
В визуальном конструкторе уроков Puck весь холст с блоками (`RichTextCanvasEditor` и `RichTextWordEditor`) монтируется внутри изолированного вложенного фрейма (`iframe#preview-frame`).
При вызове браузерных API работы с выделением и командами редактирования (`window.getSelection()`, `document.execCommand()`, `document.addEventListener('selectionchange')`, `document.createElement('span')`) компоненты обращались к верхнеуровневому объекту `window`/`document` родительской страницы, а не к `ownerDocument` и `ownerDocument.defaultView` самого iframe.
В результате:
1. Выделенный фрагмент текста (слово, предложение) не определялся (`window.getSelection()` возвращал пустую строку).
2. Нажатие на кнопки панели инструментов (Жирный, Курсив, Подчёркивание, Шрифт, Размер, Палитра цветов) не применяло стили к выделенному тексту либо вызывало ошибку кросс-документного добавления узлов.
3. Работало только выравнивание всего контейнера блока через прямое изменение стилей родителя.

## 🎯 Цели и задачи
- [x] Перевести `RichTextCanvasEditor` на динамическое получение `ownerDocument` и `defaultView` от `editorRef.current`.
- [x] Исправить методы `saveSelection` и `restoreSelection` для сохранения и восстановления Range в контексте правильного окна iframe.
- [x] Переписать метод `exec` для выполнения `doc.execCommand` в контексте `ownerDocument`.
- [x] Переписать метод `applyInlineStyle` (для цветов, шрифтов и размеров) с созданием элементов через `doc.createElement('span')` и `doc.createRange()`.
- [x] Добавить подписку на `selectionchange` как к `doc`, так и к родительскому `document`.
- [x] Реализовать аналогичные исправления для `RichTextWordEditor`.
- [x] Добавить быстрые кнопки добавления блоков в шапку редактора (`+ Текст`, `+ Статья`).
- [x] Провести E2E тестирование в Puppeteer: выделить слово, применить жирный шрифт, курсив, размер шрифта (A+) и проверить результирующий DOM.

## 📁 Затронутые файлы
- `frontend/src/components/editor/InlineEditable.tsx`
- `frontend/src/app/teacher/lessons/[id]/edit/page.tsx`

## Статус: Completed ✅
