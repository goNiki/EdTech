# 📦 Модуль: `frontend`
> **Путь:** `frontend`  
> **Роль:** Полнофункциональное клиентское веб-приложение образовательной платформы ED.Learn на базе Next.js 16 (App Router), React 19, TypeScript, Tailwind CSS 4, Zustand 5 и визуального конструктора Puck Editor.

---

## 🎯 Назначение и ответственность
Фронтенд платформы ED.Learn предоставляет интерактивный веб-интерфейс для всех ролей пользователей системы:
1. **Гости и студенты:** Публичный лендинг, каталог курсов с фильтрацией, личный кабинет (`dashboard`), отслеживание успеваемости, полноэкранный плеер интерактивных уроков (`lessons/[id]`) с мгновенным автогрейдингом квизов и сдачей эссе.
2. **Преподаватели и авторы (`teacher`):** Авторская студия для создания курсов, визуального управления структурой (модули и уроки), каскадной публикации, ручной проверки домашних заданий (`grading`), детальной аналитики студентов и конструирования уроков в WYSIWYG-редакторе на базе Puck Editor.
3. **Бесшовная интеграция с бэкендом:** Сетевой транспорт на Axios с автоматической ротацией токенов при 401 и глобальным состоянием в Zustand.

---

## ⚠️ Жестко установленные правила (Hard Invariants)
1. **Директива `'use client'` и безопасность SSR:** Все компоненты и страницы, взаимодействующие с хуками React (`useState`, `useEffect`), браузерными API (`localStorage`, `window`) или состоянием Zustand, обязаны содержать директиву `'use client'` первой строкой. Обращения к `localStorage` и `window` изолируются проверками `typeof window !== 'undefined'`.
2. **Архитектура Content-as-Data для уроков:** Содержимое уроков хранится и передается строго в виде типизированного JSON-документа спецификации Puck (12 блоков: текст, видео, тесты с 1 или несколькими ответами, сопоставление пар, пропуски, эссе, файлы). Запрещено сохранять произвольный непроверенный HTML.
3. **Управление токенами и ротация сессии:** Пара JWT-токенов (`access_token`, `refresh_token`) хранится в `localStorage`. При получении `401 Unauthorized` интерцептор Axios выполняет попытку обновления сессии через `POST /auth/refresh`. При сбое рефреша токены очищаются, а пользователь направляется на `/login`. Защита от рекурсивных циклов (`_retry = true`) обязательна.
4. **Контекстное переключение ролей (`viewMode`):** Для пользователей с правами `'teacher'`, `'author'` или `'admin'` интерфейс поддерживает бесшовное переключение вида (`student` <-> `teacher`) через Zustand и `localStorage` без сброса авторизации.
5. **Изоляция лейаутов:** 
   - Общие экраны используют `Sidebar` и `TopNavbar`.
   - Страницы аутентификации `(auth)` изолированы в собственном центрированном `AuthLayout`.
   - Редактор уроков Puck ([`teacher/lessons/[id]/edit`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/teacher/lessons/%5Bid%5D/edit/page.tsx)) и плеер студента ([`lessons/[id]`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/lessons/%5Bid%5D/page.tsx)) отключают сайдбар для полноэкранного режима работы.

---

## 🛠 Технологический стек
- **Фреймворк:** [Next.js 16.3.3](https://nextjs.org) (App Router, Turbopack)
- **Библиотека компонентов:** [React 19.2.8](https://react.dev)
- **Язык программирования:** TypeScript 5
- **Управление состоянием:** [Zustand 5.0.15](https://zustand-demo.pmnd.rs)
- **Стилизация:** Tailwind CSS 4, `@tailwindcss/postcss`, `tw-animate-css`
- **Визуальный конструктор контента:** [@puckeditor/core 0.23.0](https://puckeditor.com)
- **UI-компоненты:** Radix UI / Base UI, Lucide Icons, Shadcn UI
- **Сетевой клиент:** Axios 1.20.0
- **Темизация:** `next-themes` (светлая, темная, системная)

---

## 📁 Структура папок и файлов
```text
frontend/
├── package.json              # Зависимости и скрипты запуска
├── tsconfig.json             # Конфигурация TypeScript (алиас @/* -> ./src/*)
├── next.config.ts            # Настройки Next.js (allowedDevOrigins)
├── postcss.config.mjs        # Настройки PostCSS для Tailwind 4
├── src/
│   ├── app/                  # Маршрутизация Next.js App Router
│   │   ├── (auth)/           # Маршруты аутентификации (/login, /register)
│   │   ├── courses/          # Каталог курсов и лендинг курса ([slug])
│   │   ├── dashboard/        # Личный кабинет студента (/courses, /courses/[id], /settings)
│   │   ├── lessons/          # Плеер интерактивного урока студента ([id])
│   │   ├── teacher/          # Авторский кабинет (/courses, /curriculum, /grading, /edit)
│   │   ├── globals.css       # Глобальные стили Tailwind 4 и переменные тем
│   │   ├── layout.tsx        # Корневой лейаут приложения с ThemeProvider
│   │   └── page.tsx          # Главная страница (промо-лендинг платформы)
│   ├── components/           # Библиотека React-компонентов
│   │   ├── layout/           # Сайдбар и шапка (Sidebar, TopNavbar)
│   │   ├── player/           # Интерактивный плеер урока (PuckLessonViewer)
│   │   ├── editor/           # Дерево блоков CustomOutline, InlineEditable
│   │   ├── teacher/          # Модалки и виджеты учителя (Curriculum, Grading)
│   │   ├── ui/               # Базовые примитивы Shadcn (Button, Card, Input...)
│   │   ├── CourseCard.tsx    # Карточка курса
│   │   ├── ProtectedRoute.tsx# HOC защиты маршрутов по ролям
│   │   └── ThemeProvider.tsx # Провайдер тем
│   ├── lib/                  # Инфраструктурные модули
│   │   ├── api.ts            # Сетевой клиент Axios с перехватчиками
│   │   ├── puck-config.tsx   # Конфигурация 12 блоков Puck и парсеры шаблонов
│   │   └── utils.ts          # Хелпер стилей cn
│   └── store/                # Клиентские хранилища
│       └── useAuth.ts        # Zustand-хранилище сессии и режима роли
```

---

## 🌐 Переменные окружения (.env)
Создайте файл `.env.local` в корне каталога `frontend/`:

```env
# Базовый адрес API бэкенда платформы
NEXT_PUBLIC_API_URL=http://localhost:8082/api/v1
```

> **Примечание:** Если переменная `NEXT_PUBLIC_API_URL` не задана, сетевой клиент [`api.ts`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/api.ts#L3) автоматически использует значение по умолчанию `http://localhost:8082/api/v1`.

---

## 🚀 Команды запуска и сборки
Все команды выполняются из каталога `frontend`:

```bash
# 1. Установка зависимостей
npm install

# 2. Запуск локального сервера разработки (порт 3000)
npm run dev

# 3. Проверка типов и линтинг ESLint
npm run lint

# 4. Производственная сборка (Production Build)
npm run build

# 5. Запуск собранного production-приложения
npm run start
```

---

## 📚 Каталог фронтенд-модулей
Подробная техническая документация каждого слоя и раздела фронтенда:
- [**`src/store`**](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/store/README.md) — Zustand-хранилище `useAuth`, токены, профиль и синхронизация с `localStorage`.
- [**`src/lib`**](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/README.md) — Сетевой клиент Axios (`api.ts`), спецификация 12 блоков Puck (`puck-config.tsx`), утилиты `cn`.
- [**`src/components`**](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/README.md) — Компоненты интерфейса, плеер уроков `PuckLessonViewer`, виджеты учителя и примитивы Shadcn UI.
- [**`src/app/%28auth%29`**](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/%28auth%29/README.md) — Авторизация и регистрация пользователя (`/login`, `/register`).
- [**`src/app/courses`**](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/courses/README.md) — Публичный каталог курсов и детальный лендинг курса (`[slug]`).
- [**`src/app/dashboard`**](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/dashboard/README.md) — Кабинет студента, список курсов, навигатор по силлабусу и профиль.
- [**`src/app/lessons`**](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/lessons/README.md) — Страница прохождения интерактивного урока студентом (`lessons/[id]`).
- [**`src/app/teacher`**](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/teacher/README.md) — Кабинет преподавателя: конструктор учебного плана, очередь проверки ДЗ и редактор Puck.

---

## 🧭 Навигатор типовых задач фронтенда (Where to edit?)
- **Добавить новый интерактивный блок в уроки:**
  - Конфигурация полей и типов: [`src/lib/puck-config.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/puck-config.tsx).
  - Интерактивный плеер для студента: [`src/components/player/PuckLessonViewer.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/player/PuckLessonViewer.tsx).
  - Отображение в боковом аутлайне редактора: [`src/components/editor/CustomOutline.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/editor/CustomOutline.tsx).
- **Изменить интерфейс или логику работы с API бэкенда:**
  - Сетевые перехватчики и ретраи: [`src/lib/api.ts`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/api.ts).
  - Состояние сессии и данных пользователя: [`src/store/useAuth.ts`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/store/useAuth.ts).
- **Скорректировать стили или тему оформления:**
  - Цветовая палитра и переменные CSS: [`src/app/globals.css`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/globals.css).
  - Провайдер темы: [`src/components/ThemeProvider.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/ThemeProvider.tsx).
- **Добавить новую страницу или маршрут:**
  - Создать папку маршрута в [`src/app/`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/) с файлом `page.tsx`.
  - При необходимости ролевой защиты обернуть страницу в [`<ProtectedRoute>`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/ProtectedRoute.tsx).
