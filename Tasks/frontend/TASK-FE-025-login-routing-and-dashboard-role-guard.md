# 🎨 [FE-025] Исправление перенаправления после входа (Login Routing) и защита дашборда для ролей преподавателя

> **Приоритет:** Critical (P0)  
> **Связанные задачи:** QA-025  
> **Компоненты / Страницы:**  
> - `frontend/src/app/(auth)/login/page.tsx`  
> - `frontend/src/app/dashboard/page.tsx`  
> - `frontend/src/store/useAuth.ts`

---

## 🎯 Цель задачи
Устранить критический баг пользовательского опыта, когда преподаватель/автор/админ после успешного ввода логина и пароля перенаправляется на студенческую страницу `/dashboard` («Мое обучение») и видит пустой экран с сообщением «У вас пока нет курсов».

В рамках задачи необходимо:
1. Обеспечить надежное и детерминированное перенаправление в Авторскую студию (`/teacher/courses`) для пользователей с ролями `teacher`, `author`, `admin` или с активным режимом автора (`viewMode === 'teacher'`).
2. Исправить парсинг ответа API в `frontend/src/app/dashboard/page.tsx`, из-за которого список курсов студента всегда оставался пустым (`data.data?.courses` вместо `data.courses`).
3. Добавить на `/dashboard` ролевой шлюз (Role Guard) и баннер быстрого переключения в Авторскую студию, если преподаватель зашел на студенческую страницу по прямой ссылке.

---

## 🔍 Текущее состояние кода и коренные причины (Root Causes)

1. **В `frontend/src/app/(auth)/login/page.tsx` (строки 34–43):**
   ```tsx
   const { data } = await api.post('/auth/login', { email, password });
   login(data.access_token, data.refresh_token);
   await useAuth.getState().fetchUser();
   
   const currentUser = useAuth.getState().user;
   if (currentUser?.role === 'teacher') {
     router.push('/teacher/courses');
   } else {
     router.push('/dashboard');
   }
   ```
   - Жесткая проверка `currentUser?.role === 'teacher'` не учитывает роль `admin`, роль `author`, а также состояние `viewMode`. Если пользователь является администратором или у него сохранен режим автора, он принудительно выталкивается на `/dashboard`.
   - В методе `useAuth.login()` вызов `get().fetchUser()` запускается асинхронно без возврата `Promise`, создавая состояние гонки между сохранением токенов в `localStorage` и последующим чтением профиля.

2. **В `frontend/src/app/dashboard/page.tsx` (строка 15–16):**
   ```tsx
   const { data } = await api.get('/courses/my');
   setCourses(data.data?.courses || []);
   ```
   - Бэкенд в `internal/interfaces/handlers/courses/list_my.go` возвращает структуру `dto.PaginatedCourses` напрямую через `response.OK(w, r, respData)` без обертки в поле `data`.
   - В результате `axios` возвращает `response.data = { courses: [...], page: 1, total: ... }`. Поле `data.data` равно `undefined`, а выражение `data.data?.courses || []` ВСЕГДА возвращает пустой массив `[]`! Студенты с реальными записями на курсы также видели пустой экран.

---

## 📝 Технические требования к реализации

### 1. Доработка авторизации в `frontend/src/store/useAuth.ts`
- Сделать метод `login` асинхронным либо возвращающим промис:
  ```ts
  login: async (accessToken: string, refreshToken: string) => {
    localStorage.setItem('access_token', accessToken);
    localStorage.setItem('refresh_token', refreshToken);
    set({ isAuthenticated: true });
    await get().fetchUser();
  }
  ```
- Гарантировать, что при логине пользователя с ролью `teacher`, `author` или `admin` в состояние `viewMode` устанавливается значение `'teacher'`.

### 2. Маршрутизация в `frontend/src/app/(auth)/login/page.tsx`
- После получения токенов и загрузки пользователя определять целевой маршрут:
  - Если в URL есть query param `?redirect=/path` — перенаправлять по указанному пути.
  - Если роль пользователя входит в список `['teacher', 'author', 'admin']` ИЛИ `viewMode === 'teacher'` — перенаправлять на `/teacher/courses`.
  - В остальных случаях (студент) — перенаправлять на `/dashboard`.
- Использовать `router.replace(...)` вместо `router.push(...)`, чтобы страница логина не оставалась в истории переходов браузера (кнопка «Назад»).

### 3. Исправление и защита `frontend/src/app/dashboard/page.tsx`
- **Корректный парсинг данных курсов:**
  ```tsx
  const res = await api.get('/courses/my');
  const payload = res.data?.data || res.data;
  const courseList = payload?.courses || (Array.isArray(payload) ? payload : []);
  setCourses(courseList);
  ```
- **Role Guard / Авто-редирект преподавателей:**
  - При монтировании страницы проверять `user` из `useAuth`:
  - Если `user && ['teacher', 'author', 'admin'].includes(user.role) && viewMode === 'teacher'`:
    - Выполнять мягкий автоматический переход на `/teacher/courses`.
    - Во время перехода показывать приветственный спиннер: «Перенаправление в Авторскую студию...».
  - Если преподаватель намеренно переключил `viewMode` на `'student'` для просмотра сайта глазами ученика — разрешать просмотр `/dashboard`, но показывать верхний информационный баннер:
    - *«Вы просматриваете кабинет как студент. Переключиться в режим управления курсами -> [В Авторскую студию]»*.

---

## ✅ Критерии приёмки (Definition of Done)
- [x] При авторизации с аккаунтом преподавателя (роль `teacher`, `author` или `admin`) пользователь сразу перенаправляется на `/teacher/courses`.
- [x] Если в строке логина был параметр `?redirect=...`, после входа открывается запрошенная страница.
- [x] На студенческом `/dashboard` курсы корректно считываются из ответа бэкенда (`payload.courses`) и отображаются карточками.
- [x] Если авторизованный преподаватель вручную вводит адрес `/dashboard`, он либо автоматически перенаправляется на `/teacher/courses`, либо видит переключатель режима с возможностью в один клик вернуться к своим курсам.
- [x] Исключены белые экраны и циклические перезагрузки при обновлении страницы по F5.

## 📁 Затронутые файлы
- `frontend/src/store/useAuth.ts`
- `frontend/src/app/(auth)/login/page.tsx`
- `frontend/src/app/dashboard/page.tsx`

## ✅ Результаты валидации
1. Проверка типов `npx tsc --noEmit` успешно пройдена (0 ошибок).
2. Браузерные тесты в MCP Puppeteer подтвердили:
   - Вход под `nikit@mail.ru` детерминированно перенаправляет на `/teacher/courses` с `viewMode: "teacher"`.
   - При переходе на `/dashboard` срабатывает ролевой шлюз и возвращает преподавателя в Авторскую студию.
   - В режиме студента на `/dashboard` отображается плашка предварительного просмотра с кнопкой возврата в авторский режим.

