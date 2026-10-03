# 📦 Модуль: `store`
> **Путь:** `frontend/src/store`  
> **Роль:** Глобальное клиентское управление состоянием сессии пользователя, токенов аутентификации и режима интерфейса (студент / преподаватель) на базе Zustand.

---

## 🎯 Назначение и ответственность
Модуль `store` централизует состояние аутентификации пользователя для всего SPA-приложения Next.js. Он синхронизирует состояние сессии между памятью React, сетевыми заголовками Axios и `localStorage` браузера, а также управляет переключением контекста роли (`viewMode`: студент / преподаватель) без необходимости перезагрузки страницы.

---

## ⚠️ Жестко установленные правила (Hard Invariants)
1. **Синхронизация токенов с LocalStorage:** Токены `access_token` и `refresh_token` обязаны сохраняться в `localStorage` синхронно с вызовом `login()`. Метод `logout()` обязан гарантированно очищать оба токена из `localStorage` и сбрасывать `user` в `null`, а `isAuthenticated` в `false`.
2. **SSR/Hydration Safety:** Все чтения и записи в `localStorage` должны быть изолированы проверкой `typeof window !== 'undefined'`. Это предотвращает падение Next.js SSR при компиляции на стороне Node.js.
3. **Автоопределение роли и режима:** Если пользователь обладает ролью `'teacher'`, `'author'` или `'admin'` и в `localStorage` отсутствует сохраненный ключ `view_mode`, хранилище автоматически инициализирует `viewMode = 'teacher'`.
4. **Управление флагом isLoading:** `isLoading` инициализируется значением `true` и переходит в `false` строго после завершения запроса к `/auth/me` (как при успехе, так и при 401/ошибке сети). Это защищает хуки `ProtectedRoute` от ложных перенаправлений на страницу логина до окончания гидратации сессии.

---

## 📁 Структура файлов модуля
| Файл | Описание роли файла |
|---|---|
| [`useAuth.ts`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/store/useAuth.ts#L1-L89) | Zustand-хранилище `useAuth`, интерфейсы профиля `User`, `AuthState` и тип `ViewMode` |

---

## ⚙️ Функции, компоненты, хуки и API
| Функция / Тип | Файл:Строки | Описание | Сигнатура / Вход и Выход |
|---|---|---|---|
| [`User`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/store/useAuth.ts#L4-L17) | [`useAuth.ts#L4-L17`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/store/useAuth.ts#L4-L17) | TypeScript-интерфейс профиля пользователя с бэкенда | `interface User { id: number; email: string; role: string; ... }` |
| [`ViewMode`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/store/useAuth.ts#L19) | [`useAuth.ts#L19`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/store/useAuth.ts#L19) | Тип переключения интерфейса пользователя | `'student' \| 'teacher'` |
| [`AuthState`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/store/useAuth.ts#L21-L31) | [`useAuth.ts#L21-L31`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/store/useAuth.ts#L21-L31) | Интерфейс состояния и методов экшенов стора | `interface AuthState { user, isAuthenticated, login, logout, ... }` |
| [`useAuth`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/store/useAuth.ts#L33-L88) | [`useAuth.ts#L33-L88`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/store/useAuth.ts#L33-L88) | Основной клиентский React-хук состояния аутентификации | `create<AuthState>((set, get) => ({ ... }))` |
| `login` | [`useAuth.ts#L39-L44`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/store/useAuth.ts#L39-L44) | Сохранение JWT-пары в `localStorage` и запуск запроса `fetchUser()` | `(accessToken: string, refreshToken: string) => void` |
| `logout` | [`useAuth.ts#L46-L50`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/store/useAuth.ts#L46-L50) | Очистка токенов из `localStorage` и сброс состояния в дефолт | `() => void` |
| `setViewMode` | [`useAuth.ts#L52-L55`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/store/useAuth.ts#L52-L55) | Переключение роли интерфейса и сохранение ключа `view_mode` | `(mode: ViewMode) => void` |
| `setUser` | [`useAuth.ts#L57-L59`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/store/useAuth.ts#L57-L59) | Локальное обновление объекта пользователя без сетевого запроса | `(user: User) => void` |
| `fetchUser` | [`useAuth.ts#L61-L87`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/store/useAuth.ts#L61-L87) | Запрос профиля `GET /auth/me`, определение роли и установка `isLoading = false` | `() => Promise<void>` |

---

## 🔗 Зависимости
- **Входящие (кто вызывает этот модуль):**
  - [`ProtectedRoute.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/ProtectedRoute.tsx#L3) — проверка авторизации и ролевого доступа.
  - [`Sidebar.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/layout/Sidebar.tsx#L6) — переключатель `viewMode`, кнопка логаута, аватар.
  - [`TopNavbar.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/layout/TopNavbar.tsx#L7) — чтение текущего контекста пользователя.
  - [`login/page.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/%28auth%29/login/page.tsx#L6) и [`register/page.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/%28auth%29/register/page.tsx#L6) — сохранение сессии после ответа бэкенда.
  - [`courses/page.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/courses/page.tsx#L6) и [`courses/[slug]/page.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/courses/[slug]/page.tsx#L5) — проверка зачисления и авторизации.
  - [`settings/page.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/dashboard/settings/page.tsx#L4) — чтение и обновление полей профиля через `setUser`.
- **Исходящие (что импортирует этот модуль):**
  - [`api.ts`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/api.ts#L5) — Axios-клиент для выполнения `GET /auth/me`.
  - `zustand` — библиотека управления состоянием.

---

## 🧭 Навигатор типовых задач (Where to edit?)
- **Добавить новое поле в профиль пользователя:** расширить интерфейс `User` в [`useAuth.ts#L4-L17`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/store/useAuth.ts#L4-L17).
- **Добавить новую роль или режим интерфейса:** обновить объединение типов `ViewMode` в [`useAuth.ts#L19`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/store/useAuth.ts#L19) и логику в `fetchUser` [`useAuth.ts#L72-L75`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/store/useAuth.ts#L72-L75).
- **Сбросить сессию принудительно:** вызвать экшен `useAuth.getState().logout()` из любой точки приложения.
