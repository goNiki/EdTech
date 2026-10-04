# TASK-FE-022: Унификация панелей инструментов Canvas и Word редакторов, поддержка списков, цитат, заголовков H1-H3 и Tab-отступов

## 📋 Описание проблемы
1. **Списки и цитаты в Canvas-блоке:** При нажатии на кнопки маркированного/нумерованного списка и цитаты разметка вставлялась, но визуально маркеры и стили цитат отсутствовали из-за глобального сброса стилей Tailwind CSS (`list-style: none`, сброс `blockquote`).
2. **Различие панелей управления:** Текстовый блок (Canvas) и Лекция/Word имели разные наборы кнопок и разный пользовательский опыт, хотя логически решают смежные задачи оформления контента.
3. **Отсутствие визуальной иерархии заголовков:** По умолчанию заголовки H1, H2, H3 и обычный текст P не имели четкой ступенчатой визуальной градации по размеру шрифта и отступам.
4. **Клавиша Tab:** Нажатие Tab должно создавать абзацный отступ или вложенный список без потери фокуса и конфликтов с Puck.

## 🎯 Цели и задачи
- [x] Добавить в `globals.css` явные селекторы для `.edtech-inline-editable` и `[contenteditable]`:
  - `h1`: 2rem / font-extrabold / line-height 1.25,
  - `h2`: 1.5rem / font-bold / line-height 1.3,
  - `h3`: 1.25rem / font-semibold / line-height 1.35,
  - `p`: 1rem / leading-relaxed,
  - `ul`: `list-style-type: disc !important`, padding-left 1.75rem,
  - `ol`: `list-style-type: decimal !important`, padding-left 1.75rem,
  - `li`: `display: list-item !important`,
  - `blockquote`: левая индиго-граница 4px, фон, курсив, скругление.
- [x] Привести панели инструментов `RichTextCanvasEditor` и `RichTextWordEditor` к единому богатому интерфейсу с одинаковым расположением управляющих элементов.
- [x] Доработать обработчик клавиши Tab:
  - в списках — вызывать `exec('indent')` / `exec('outdent')` (при Shift+Tab),
  - в тексте — вставлять абзацный отступ `&emsp;&emsp;`.
- [x] Протестировать работу списков, цитат, H1-H3 и Tab в браузере через MCP Puppeteer.

## 📁 Затронутые файлы
- `frontend/src/app/globals.css`
- `frontend/src/components/editor/InlineEditable.tsx`

## ✅ Результаты валидации
1. `npx tsc --noEmit` пройден без ошибок (0 errors).
2. Браузерные тесты Puppeteer подтвердили:
   - `h1`: 32px (font-extrabold 800)
   - `h2`: 24px (font-bold 700)
   - `h3`: 20px (font-semibold 600)
   - `p`: 16px
   - `ul / li`: маркеры disc активны, padding-left 28px
   - `ol / li`: нумерация decimal активна, padding-left 28px
   - `blockquote`: граница 4px rgb(99, 102, 241), курсив, скругленные углы.
   - Tab корректно делает отступ абзаца и управляет вложенностью списков.
