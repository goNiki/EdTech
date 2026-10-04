# 🎨 [FE-045] Расширение полотна урока и адаптивная типографика Word-документов и таблиц

> **Приоритет:** High (P1)  
> **Связанные задачи:** QA-045, FE-044  
> **Компоненты / Страницы:** `frontend/src/app/lessons/[id]/page.tsx`, `frontend/src/components/player/PuckLessonViewer.tsx`, `frontend/src/lib/document-importer.ts`  
> **Документация модуля:** [frontend/src/components/README.md](../../frontend/src/components/README.md)

---

## 🎯 Цель задачи
Устранить эффект «узкого полотна по центру» и решить проблему нечитаемых сжатых таблиц из Word:
1. Расширить базовую ширину контейнера урока с текущих `max-w-4xl` (896px) до **`max-w-6xl` (1152px)** по умолчанию с возможностью переключения в режим **`max-w-7xl` (1280px / 92% экрана)**.
2. Устранить хардкод микроскопического шрифта `text-xs` (12px) в сгенерированных таблицах Word ([`document-importer.ts`](../../frontend/src/lib/document-importer.ts#L70)) — перевести таблицы на масштабируемый размер шрифта читалки.
3. Обернуть таблицы и широкие блоки в контейнеры горизонтальной прокрутки (`overflow-x-auto`) с красивыми мягкими скроллбарами, чтобы даже сложные многоколоночные документы не ломали лейаут.

---

## 🔍 Текущее состояние кода
1. В [`frontend/src/app/lessons/[id]/page.tsx#L420`](../../frontend/src/app/lessons/%5Bid%5D/page.tsx#L420):
   ```tsx
   <main className="flex-1 max-w-4xl w-full mx-auto py-6 sm:py-8 px-4 sm:px-6 space-y-6">
   ```
   Контейнер жестко ограничен 896px. На экранах 1440px и 1920px более половины экрана занято белыми пустыми полями.
2. В [`frontend/src/lib/document-importer.ts#L70`](../../frontend/src/lib/document-importer.ts#L70):
   ```typescript
   '<table class="w-full my-4 border-collapse border border-slate-300 dark:border-slate-700 text-xs ..."'
   ```
   Класс `text-xs` принудительно делает весь текст внутри таблиц 12px, из-за чего документы и заявления невозможно читать.

---

## 📝 Технические требования к реализации

### 1. Расширение контейнера плеера урока (`lessons/[id]/page.tsx`)
- Заменить `max-w-4xl` на динамический класс, зависящий от состояния `contentWidth` (из `ReadingSettingsPopover`):
  ```tsx
  const containerWidthClass = contentWidth === 'wide' 
    ? 'max-w-6xl xl:max-w-7xl' 
    : 'max-w-5xl';
  ```
- Для элементов баннера завершения урока (`QuizPreflightScreen`, `HomeworkFeedbackCard`) сохранять гармоничные отступы и адаптивность.

### 2. Исправление `document-importer.ts`
- В функции `cleanWordHtml`:
  - Заменить `text-xs` в `<table ...>` на класс `text-sm sm:text-base`.
  - Оборачивать каждую таблицу в адаптивный скролл-контейнер:
    ```html
    <div class="my-4 overflow-x-auto rounded-2xl border border-slate-200 dark:border-slate-800 shadow-2xs">
      <table class="w-full border-collapse text-sm sm:text-base ...">
        ...
      </table>
    </div>
    ```
  - Увеличить паддинги ячеек: `p-3.5 sm:p-4` для улучшения визуального чанкинга строк.

### 3. Стилизация `prose` в `PuckLessonViewer.tsx`
- На контейнере RichText и TextBlock расширить `prose`:
  - `prose-slate max-w-none prose-headings:font-black prose-p:leading-relaxed prose-table:my-0`.

---

## ✅ Критерии приёмки (Definition of Done)
- [x] Полотно урока на экране 1440px+ визуально шире, пустое пространство слева и справа гармонично сбалансировано.
- [x] Текст внутри таблиц из Word читается легко (шрифт соответствует общему размеру урока, не мельчит).
- [x] Широкие таблицы прокручиваются горизонтально без разрыва страницы.
- [x] Переключение ширины («Стандарт» / «Широкий экран») плавно меняет контейнер без прыжков верстки.

---

## 🏁 Статус выполнения
- **Статус:** Completed ✅
- **Реализация:**
  1. В [`frontend/src/app/lessons/[id]/page.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/lessons/[id]/page.tsx) базовое полотно плеера расширено с устаревшего `max-w-4xl` до `max-w-5xl` (по умолчанию) и динамического `max-w-6xl xl:max-w-7xl` при выборе широкого экрана в `ReadingSettingsPopover`.
  2. В [`frontend/src/lib/document-importer.ts`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/document-importer.ts) устранен хардкод `text-xs`: таблицы переведены на масштабируемый шрифт `text-sm sm:text-base` с комфортными паддингами `p-3.5 sm:p-4` и автоматически оборачиваются в изолированный контейнер горизонтального скролла (`overflow-x-auto rounded-2xl border border-slate-200 dark:border-slate-800 shadow-2xs`).
  3. В [`frontend/src/components/player/PuckLessonViewer.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/player/PuckLessonViewer.tsx) текстовые блоки и RichTextBlock расширены классами `lesson-content-area prose prose-slate dark:prose-invert max-w-none prose-headings:font-black prose-p:leading-relaxed prose-table:my-0`.
- **Верификация:** Проведено тестирование рендеринга и переключения ширины контейнера в Puppeteer. All checks passed.
