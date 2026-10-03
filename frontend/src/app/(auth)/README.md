# 📦 Модуль: `(auth)`
> **Путь:** `frontend/src/app/%28auth%29`  
> **Роль:** Модуль аутентификации и регистрации: изолированный лейаут входа, валидация учетных данных, сохранение JWT-сессии и ролевой редирект.

---

## 🎯 Назначение и ответственность
Модуль обслуживает публичные маршруты входа (`/login`) и регистрации (`/register`). Он принимает учетные данные пользователя, взаимодействует с эндпоинтами бэкенда `/auth/login` и `/auth/register`, передает полученные токены в глобальное хранилище `useAuth` и направляет пользователя в соответствующий личный кабинет в зависимости от назначенной роли.

---

## ⚠️ Жестко установленные правила (Hard Invariants)
1. **Сквозная регистрация с авто-логином (Seamless Onboarding):** При успешном создании учетной записи через `POST /auth/register` форма обязана автоматически выполнить вызов `POST /auth/login`, сохранить полученные токены через `login(tokens)` и загрузить профиль `fetchUser()`. Пользователю не требуется вводить пароль повторно.
2. **Ролевая маршрутизация (Role-Based Landing):** После успешной аутентификации выполняется проверка `currentUser?.role`:
   - Если роль `teacher` (преподаватель) — редирект на [`/teacher/courses`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/teacher/courses/page.tsx).
   - В остальных случаях (студент) — редирект на [`/dashboard`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/dashboard/page.tsx).
3. **Graceful Network Error Handling:** При ошибках `ERR_NETWORK` или `Network Error` (когда бэкенд выключен или недоступен) формы обязаны выводить понятное сообщение: *"Нет связи с сервером. Убедитесь, что бэкенд запущен."*, исключая необработанные падения приложения.
4. **Изоляция лейаута:** Роут-группа `(auth)` исключена из общей структуры страниц с `Sidebar` и рендерится в центрированном контейнере с промо-панелью.

---

## 📁 Структура файлов модуля
| Файл | Описание роли файла |
|---|---|
| [`layout.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/%28auth%29/layout.tsx#L1-L45) | Изолированный каркас страниц аутентификации с логотипом ED.Learn и ссылкой на главную |
| [`login/page.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/%28auth%29/login/page.tsx#L1-L240) | Страница входа в личный кабинет с валидацией, переключателем видимости пароля и чеклистом преимуществ |
| [`register/page.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/%28auth%29/register/page.tsx#L1-L275) | Страница создания аккаунта со сквозным входом и согласием с условиями платформы |

---

## ⚙️ Функции, компоненты, хуки и API
| Компонент / Хэндлер | Файл:Строки | Описание | Маршрут / Сигнатура |
|---|---|---|---|
| [`AuthLayout`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/%28auth%29/layout.tsx#L5-L44) | [`layout.tsx#L5-L44`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/%28auth%29/layout.tsx#L5-L44) | Базовый макет с брендовой шапкой и футером | `{ children: React.ReactNode }` |
| [`LoginPage`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/%28auth%29/login/page.tsx#L17-L239) | [`login/page.tsx#L17-L239`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/%28auth%29/login/page.tsx#L17-L239) | Клиентский компонент формы авторизации | Маршрут `/login` |
| `handleSubmit (login)` | [`login/page.tsx#L28-L53`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/%28auth%29/login/page.tsx#L28-L53) | Отправка `POST /auth/login`, запись токенов и ролевой редирект | `(e: React.FormEvent) => Promise<void>` |
| [`RegisterPage`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/%28auth%29/register/page.tsx#L18-L274) | [`register/page.tsx#L18-L274`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/%28auth%29/register/page.tsx#L18-L274) | Клиентский компонент формы регистрации | Маршрут `/register` |
| `handleSubmit (register)` | [`register/page.tsx#L30-L72`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/%28auth%29/register/page.tsx#L30-L72) | Отправка `POST /auth/register` + бесшовный вход через `POST /auth/login` | `(e: React.FormEvent) => Promise<void>` |

---

## 🔗 Зависимости
- **Входящие (кто вызывает этот модуль):**
  - Пользователи при переходе по прямым ссылкам `/login`, `/register`.
  - Автоматические редиректы из [`ProtectedRoute.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/ProtectedRoute.tsx#L30) и интерцептора [`api.ts`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/api.ts#L45) при истечении сессии.
- **Исходящие (что импортирует этот модуль):**
  - [`store/useAuth.ts`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/store/useAuth.ts#L33) — методы `login()` и `fetchUser()`.
  - [`lib/api.ts`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/api.ts#L5) — Axios-клиент для обращений к `/auth/login` и `/auth/register`.
  - `next/navigation` (`useRouter`).

---

## 🧭 Навигатор типовых задач (Where to edit?)
- **Изменить адрес редиректа после логина:** отредактировать блок проверки роли в [`login/page.tsx#L38-L43`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/%28auth%29/login/page.tsx#L38-L43).
- **Добавить валидацию сложности пароля при регистрации:** расширить метод `handleSubmit` в [`register/page.tsx#L30-L40`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/%28auth%29/register/page.tsx#L30-L40).
- **Изменить брендинг или текст промо-блока:** отредактировать левую панель в [`login/page.tsx#L58-L115`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/%28auth%29/login/page.tsx#L58-L115).
