# 🎨 [FE-002] Интеграция реальной загрузки файлов в плеер уроков и настройки профиля

> **Приоритет:** Critical (P0)  
> **Связанные задачи:** BE-002, QA-002  
> **Компоненты / Страницы:** `frontend/src/components/player/PuckLessonViewer.tsx`, `frontend/src/app/dashboard/settings/page.tsx`, `frontend/src/app/teacher/courses/new/page.tsx`

## 🎯 Цель задачи
Оживить декоративную заглушку в `FileUploadBlock` в интерактивном плеере уроков: подключить скрытый `<input type="file">`, отправлять файл на бэкенд через `multipart/form-data`, отображать статус и прогресс загрузки, а полученный URL передавать в пейлоад завершения урока (`completeLesson`). Также добавить загрузку файлов для аватара и обложки курса вместо ручного ввода ссылок.

## 🖥 Пользовательский интерфейс (UI / UX)
- **Блок `FileUploadBlock` в `PuckLessonViewer`:**
  - Клик по дропзоне открывает системный файловый диалог (атрибут `accept=".zip,.pdf"`).
  - Поддержка Drag-and-Drop: подсветка рамки при перетаскивании файла.
  - Состояние загрузки: спиннер / прогресс-бар («Загрузка файла 45%...»).
  - Успешное состояние: плашка с иконкой файла, именем файла, размером и кнопкой «Заменить файл / Удалить».
  - Состояние ошибки: красный тост с текстом ошибки (например, «Файл превышает 25 МБ»).
- **Аватарка (`dashboard/settings`):**
  - Кнопка «Загрузить фото» поверх аватара, отправляющая выбранный файл на `POST /upload?category=avatar` и подставляющая полученный URL в инпут `avatarUrl`.
- **Обложка курса (`teacher/courses/new` и `[id]/settings`):**
  - Кнопка выбора локального изображения с мгновенным предпросмотром.

## 🔌 Интеграция с API (Контракт с Бэкендом)
- **Сетевой клиент:**
  ```typescript
  const formData = new FormData();
  formData.append('file', file);
  formData.append('category', 'homework');
  const res = await api.post('/upload', formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
  });
  const fileUrl = res.data.file_url;
  ```
- **Пейлоад завершения урока (`handleFinish`):**
  - В массив `essays` вместо хардкода отправлять реальный URL:
    ```typescript
    submittedEssays.push({
      question_text: props.title || 'Загрузка практической работы',
      answer_text: `Файл решения: ${uploadedFile.name} (скачать: ${fileUrl})`,
      max_points: Number(props.points) || 50,
    });
    ```

## ✅ Критерии приёмки (Definition of Done)
- [x] Студент может выбрать и загрузить файл `.pdf` или `.zip` в уроке.
- [x] Файл физически уходит на сервер, а не имитируется текстом.
- [x] Преподаватель в модальном окне проверки (`ModalGradeHW`) видит кликабельную ссылку на реальный файл студента.
- [x] В профиле пользователя работает загрузка фото с диска.

---

## 🏁 Статус выполнения
- **Статус:** ✅ Completed
- **Реализовано:**
  1. `frontend/src/components/player/PuckLessonViewer.tsx`: добавлен реальный аплоад файлов (`POST /upload?category=homework`), прогресс загрузки, drag-and-drop, бейдж с именем/размером файла, замена/удаление, передача реальной ссылки на скачивание в payload завершения урока.
  2. `frontend/src/app/dashboard/settings/page.tsx`: интерактивный аплоад аватара через hover-оверлей и кнопку «Загрузить с диска» (`POST /upload?category=avatar`), резолвинг относительных путей, удаление.
  3. `frontend/src/app/teacher/courses/new/page.tsx`: выбор файла обложки (`POST /upload?category=course_cover`), превью и очистка.
  4. `frontend/src/app/teacher/courses/[id]/settings/page.tsx`: выбор файла обложки, превью, сохранение в параметры курса.
  5. Сборка `next build` проверена и успешно завершается.
