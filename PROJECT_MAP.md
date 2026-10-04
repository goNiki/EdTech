# 🗺 Карта проекта (Project Architecture Map)
> Этот документ — центральная точка входа для разработчиков и ИИ-агентов.  
> **Правило:** При выполнении любых задач сначала сверяйтесь с этим файлом, чтобы мгновенно локализовать нужный слой и модуль кодовой базы.

---

## 🏛 Общая архитектура системы

Платформа представляет собой полнофункциональную образовательную экосистему, разделенную на современный клиентский SPA/SSR-фронтенд и строго типизированный бэкенд на языке Go:

```text
============================== FRONTEND CLIENT LAYER ==============================
[ Next.js 16 App Router (React 19, TypeScript, Tailwind CSS 4, Lucide Icons) ]
      │
      ├─ UI & Layouts: TopNavbar, Sidebar, ThemeProvider, CourseCard, Modals
      ├─ Interactive Engines: Puck Visual Builder & PuckLessonViewer Player (Content-as-Data)
      ├─ Client State: Zustand useAuth Store (localStorage JWT synchronization & viewMode)
      └─ Network Client: Axios api instance (Bearer interceptor + auto-refresh on 401)
                   │
                   ▼ (HTTP REST API /api/v1, JSON)
============================== BACKEND GO DDD LAYER ===============================
  1. Entrypoint & Transport (cmd/edtech, chi.Router)
                   │
                   ▼
  2. Presentation Layer (internal/interfaces/handlers, middleware)
     ├─ DTO Contracts & Validation (internal/dto)
     └─ Mappers (internal/interfaces/handlers/converter)
                   │
                   ▼
  3. Business Logic / Use Cases (internal/service)
     ├─ Domain Entities & Invariants (internal/domain)
     └─ Transaction Orchestration (internal/infrastructure/txmanager)
                   │
                   ▼
  4. Data Access Layer (internal/repository)
     ├─ DB Models & Scanners (internal/repository/models)
     └─ QueryExecutor (pgxpool.Pool / pgx.Tx)
                   │
                   ▼
  5. Infrastructure & Database (internal/infrastructure, PostgreSQL 16, goose)
```

---

## ⚡ Сквозная матрица бизнес-доменов (Fast Lookup Matrix)

Используйте эту таблицу для мгновенного сквозного перехода между UI-экранами, сетевым транспортом и бэкенд-модулями конкретной бизнес-фичи:

| Бизнес-домен | Frontend Pages & Components | HTTP Handlers | Use Case Service | Repository | Доменные сущности | Миграции БД |
|---|---|---|---|---|---|---|
| **Auth & Users** | [`(auth)/login`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/%28auth%29/login/page.tsx), [`(auth)/register`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/%28auth%29/register/page.tsx), [`store/useAuth`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/store/useAuth.ts), [`ProtectedRoute`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/ProtectedRoute.tsx) | [`handlers/auth`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/auth/README.md) | [`service/auth`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/auth/README.md) | [`repository/auth`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/auth/README.md), [`refresh`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/refresh/README.md) | [`domain/auth.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/domain/auth.go) | `create_table.sql`, `create_tablejwt.sql` |
| **Курсы (Courses)** | [`courses/page.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/courses/page.tsx), [`courses/[slug]`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/courses/%5Bslug%5D/page.tsx), [`teacher/courses`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/teacher/courses/page.tsx), [`teacher/courses/new`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/teacher/courses/new/page.tsx), [`CourseCard`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/CourseCard.tsx) | [`handlers/courses`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/courses/README.md) | [`service/course`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/course/README.md) | [`repository/course`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/course/README.md) | [`domain/course.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/domain/course.go) | `create_table_course.sql`, `categories.sql` |
| **Категории (Categories)** | [`courses/page.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/courses/page.tsx) | [`handlers/category`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/category/README.md) | [`service/category`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/category/README.md) | [`repository/category`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/category/README.md) | [`domain/category.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/domain/category.go) | `create_categories.sql` |
| **Секции курса** | [`teacher/curriculum`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/teacher/courses/%5Bid%5D/curriculum/page.tsx), [`ModalCreateModule`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/teacher/ModalCreateModule.tsx), [`ModalEditModule`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/teacher/ModalEditModule.tsx) | [`handlers/section`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/section/README.md) | [`service/section`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/section/README.md) | [`repository/section`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/section/README.md) | [`domain/section.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/domain/section.go) | `create_sections.sql` |
| **Уроки и материалы** | [`lessons/[id]`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/lessons/%5Bid%5D/page.tsx), [`PuckLessonViewer`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/player/PuckLessonViewer.tsx), [`teacher/lessons/[id]/edit`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/teacher/lessons/%5Bid%5D/edit/page.tsx), [`lib/puck-config`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/puck-config.tsx) | [`handlers/lesson`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/lesson/README.md) | [`service/lesson`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/lesson/README.md), [`resource`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/resource/README.md) | [`repository/lesson`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/lesson/README.md), [`resource`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/resource/README.md) | [`domain/lesson.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/domain/lesson.go), [`resource.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/domain/resource.go) | `create_table_lesson.sql`, `resource.sql` |
| **Квизы и тесты** | [`lib/puck-config`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/puck-config.tsx), [`PuckLessonViewer`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/player/PuckLessonViewer.tsx), [`teacher/grading`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/teacher/grading/page.tsx), [`ModalGradeHW`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/teacher/ModalGradeHW.tsx) | [`handlers/quiz`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/quiz/README.md) | [`service/quiz`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/quiz/README.md) | [`repository/quiz`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/quiz/README.md) | [`domain/quiz.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/domain/quiz.go) | `create_quizzes.sql`, `attempts.sql` |
| **Прогресс обучения** | [`dashboard`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/dashboard/page.tsx), [`dashboard/courses`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/dashboard/courses/page.tsx), [`dashboard/courses/[id]`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/dashboard/courses/%5Bid%5D/page.tsx), [`CourseAnalyticsCards`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/teacher/CourseAnalyticsCards.tsx) | [`handlers/progress`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/progress/README.md) | [`service/progress`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/progress/README.md) | [`repository/progress`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/progress/README.md) | [`domain/progress.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/domain/progress.go) | `lesson_progress.sql`, `course_progress.sql` |
| **Зачисление (Enroll)** | [`courses/page.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/courses/page.tsx), [`courses/[slug]`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/courses/%5Bslug%5D/page.tsx), [`ModalAddStudent`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/teacher/ModalAddStudent.tsx), [`StudentsTable`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/teacher/StudentsTable.tsx) | [`handlers/enrollment`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/enrollment/README.md) | [`service/enrollment`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/enrollment/README.md) | [`repository/enrollment`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/enrollment/README.md) | [`domain/enroll.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/domain/enroll.go) | `create_table_users_course.sql` |
| **Аналитика и грейдинг** | [`teacher/curriculum`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/teacher/courses/%5Bid%5D/curriculum/page.tsx), [`teacher/grading`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/teacher/grading/page.tsx), [`ModalStudentDrilldown`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/teacher/ModalStudentDrilldown.tsx), [`PendingHomeworksQueue`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/teacher/PendingHomeworksQueue.tsx) | [`handlers/analytics`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/analytics/README.md) | [`service/analytics`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/analytics/README.md) | [`repository/analytics`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/analytics/README.md) | [`domain/analytics.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/domain/analytics.go) | Агрегации над `attempts`, `progress` |
| **Права доступа (RBAC)** | [`ProtectedRoute`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/ProtectedRoute.tsx), [`Sidebar`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/layout/Sidebar.tsx), [`teacher/layout`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/teacher/layout.tsx), [`store/useAuth`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/store/useAuth.ts) | [`middleware/auth`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/middleware/README.md) | [`service/access`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/access/README.md) | [`repository/permission`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/permission/README.md) | [`domain/permission.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/domain/permission.go) | `role_table.sql`, `permissions.sql` |
| **Файлы и хранилище** | [`FileUploadBlock`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/player/FileUploadBlock.tsx) | [`handlers/upload`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/upload/README.md) | [`service/upload`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/upload/README.md) | `storage.LocalStorage` | [`domain/upload.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/domain/upload.go) | Локальное хранилище / S3 |
| **Отзывы и рейтинги** | [`courses/[slug]`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/courses/%5Bslug%5D/page.tsx) | [`handlers/review`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/review/README.md) | [`service/review`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/review/README.md) | [`repository/review`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/review/README.md) | [`domain/review.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/domain/review.go) | `20261004120000_create_course_reviews.sql` |

---

## 📚 Каталог модулей

### 🌐 Уровень Frontend (Next.js 16 App Router, React 19, Zustand, Puck):
| Модуль | Расположение | Документация | Зона ответственности |
|---|---|---|---|
| **Frontend Master** | `frontend` | [`frontend/README.md`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/README.md) | Архитектура Next.js 16, стек, скрипты сборки, переменные окружения |
| **Auth Store** | `frontend/src/store` | [`src/store/README.md`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/store/README.md) | Zustand `useAuth`, токены, профиль, синхронизация с `localStorage` |
| **Core Lib & Puck** | `frontend/src/lib` | [`src/lib/README.md`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/README.md) | Axios API клиент с 401 авто-рефрешем, 12 блоков Puck Editor, парсер пропусков |
| **UI & Layouts** | `frontend/src/components` | [`src/components/README.md`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/README.md) | Сайдбар, шапка, `ProtectedRoute`, `PuckLessonViewer`, виджеты автора, Shadcn UI |
| **Auth Routing** | `frontend/src/app/(auth)` | [`src/app/(auth)/README.md`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/%28auth%29/README.md) | Экраны входа `/login`, регистрации `/register`, валидация и ролевой редирект |
| **Courses Catalog** | `frontend/src/app/courses` | [`src/app/courses/README.md`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/courses/README.md) | Публичный каталог курсов с фильтрами и детальный лендинг курса `[slug]` |
| **Student Dashboard** | `frontend/src/app/dashboard` | [`src/app/dashboard/README.md`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/dashboard/README.md) | Кабинет студента, список курсов, навигатор по силлабусу и профиль |
| **Interactive Player** | `frontend/src/app/lessons` | [`src/app/lessons/README.md`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/lessons/README.md) | Полноэкранный плеер прохождения урока `lessons/[id]` с автогрейдингом |
| **Teacher Studio** | `frontend/src/app/teacher` | [`src/app/teacher/README.md`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/teacher/README.md) | Студия автора: конструктор силлабуса, проверка заданий и редактор Puck |

### ⚙️ Уровень ядра и архитектурных слоев бэкенда (Go Clean Architecture):
| Модуль | Расположение | Документация | Зона ответственности |
|---|---|---|---|
| **Server Entrypoint** | `cmd/edtech` | [`cmd/edtech/README.md`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/cmd/edtech/README.md) | Запуск веб-сервера, перехват прерываний ОС и graceful shutdown |
| **Migration CLI** | `cmd/migration` | [`cmd/migration/README.md`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/cmd/migration/README.md) | Утилита командной строки и меню для применения/создания SQL-миграций |
| **Application & DI** | `internal/app` | [`internal/app/README.md`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/app/README.md) | Контейнер зависимостей `diContainer`, Chi роутер, HTTP-таймауты |
| **Domain Entities** | `internal/domain` | [`internal/domain/README.md`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/domain/README.md) | Чистые доменные сущности, статусы, роли, инварианты бизнес-логики |
| **DTO Contracts** | `internal/dto` | [`internal/dto/README.md`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/dto/README.md) | JSON-структуры запросов/ответов с тегами валидации |
| **HTTP Handlers** | `internal/interfaces/handlers` | [`handlers/README.md`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/README.md) | Декодирование JSON, валидация параметров, маппинг в домен, HTTP-ответы |
| **HTTP Middlewares** | `internal/interfaces/middleware` | [`middleware/README.md`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/middleware/README.md) | Проверка Bearer JWT, извлечение claims, проверка ролей, slog-логирование |
| **Business Services** | `internal/service` | [`service/README.md`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/README.md) | Реализация всех сценариев (Use Cases), расчет прогресса, скоринга, прав |
| **Repositories** | `internal/repository` | [`repository/README.md`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/README.md) | Выполнение SQL-запросов к PostgreSQL через `QueryExecutor` |
| **Infrastructure** | `internal/infrastructure` | [`infrastructure/README.md`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/infrastructure/README.md) | Пул соединений pgx, JWT-менеджер, Hasher bcrypt/sha256, TxManager |
| **Shared Platform** | `pkg` | [`pkg/README.md`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/pkg/README.md) | Реестр sentinel-ошибок `errorsAPP` и утилиты валидации |
| **SQL Migrations** | `migrators` | [`migrators/README.md`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/migrators/README.md) | 21 версионированная SQL-миграция схемы данных Goose |

### Документация подпапок сервисов (`internal/service/*`):
* [**`service/auth`**](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/auth/README.md) — Сессии, bcrypt хэши, одиночная активная сессия, бан-правила.
* [**`service/access`**](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/access/README.md) — Матрица разрешений курсов, правила черновиков и прав создателя.
* [**`service/course`**](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/course/README.md) — Каскадная публикация, автозачисление создателя, сборка дерева курса.
* [**`service/section`**](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/section/README.md) — Позиционирование секций и переупорядочивание уроков.
* [**`service/lesson`**](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/lesson/README.md) — Вычисление позиций `max(pos)+1`, CRUD уроков, контекст навигации (prev/next) и силлабус.
* [**`service/enrollment`**](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/enrollment/README.md) — Транзакционный инкремент счетчика студентов, запрет отчисления автора.
* [**`service/progress`**](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/progress/README.md) — Атомарный пересчет курса при завершении урока, прием эссе, таймеры, автосохранение драфтов и активные попытки.
* [**`service/quiz`**](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/quiz/README.md) — Advisory lock от гонок, `FOR UPDATE`, отложенная ручная проверка эссе.
* [**`service/analytics`**](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/analytics/README.md) — Защита доступа преподавателя, аналитика, drilldown и кросс-курсовая очередь проверки.
* [**`service/category`**](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/category/README.md) — Каталог категорий курсов, подсчет опубликованных курсов, CRUD категорий.
* [**`service/resource`**](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/resource/README.md) — Учебные материалы уроков.
* [**`service/review`**](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/review/README.md) — Рейтинги 1–5 звезд, проверка прогресса >= 30%, транзакционный пересчет статистики.
* [**`service/certificate`**](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/certificate/README.md) — Проверка 100% прогресса, генерация кода EDL-YYYY-XXXXXXXX, публичная верификация.
* [**`service/notification`**](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/notification/README.md) — In-app уведомления, получение ленты, пометка прочитанными.
* [**`service/upload`**](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/upload/README.md) — Прием и валидация файлов (до 25 МБ, презентации до 50 МБ), пакетная загрузка Word-изображений.

### Документация подпапок репозиториев (`internal/repository/*`):
* [**`repository/auth`**](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/auth/README.md) — SQL пользователей, перехват уникальности `23505`.
* [**`repository/refresh`**](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/refresh/README.md) — SHA-256 хранение токенов сессий, сброс всех сессий.
* [**`repository/category`**](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/category/README.md) — Выборка категорий с подсчетом курсов через `LEFT JOIN`, CRUD.
* [**`repository/course`**](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/course/README.md) — SQL каталога, атомарные `enrolled_count ± 1`, `UpdateCourseRatingStats`.
* [**`repository/section`**](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/section/README.md) — `ORDER BY position ASC`, `COALESCE(MAX(pos), 0)`.
* [**`repository/lesson`**](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/lesson/README.md) — Пакетное обновление позиций, смена секций.
* [**`repository/enrollment`**](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/enrollment/README.md) — Однопроходный CTE-запрос студентов с прогрессом и долгами.
* [**`repository/progress`**](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/progress/README.md) — Идемпотентный `ON CONFLICT DO UPDATE`, авто-дата 100%.
* [**`repository/quiz`**](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/quiz/README.md) — `pg_advisory_xact_lock`, `SELECT ... FOR UPDATE`, хранение эссе с `is_correct NULL`.
* [**`repository/analytics`**](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/analytics/README.md) — 4-блочный CTE-запрос сводки курса без N+1.
* [**`repository/permission`**](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/permission/README.md) — `SELECT EXISTS(...)` по ролям и действиям.
* [**`repository/resource`**](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/resource/README.md) — Вложения уроков.
* [**`repository/review`**](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/review/README.md) — SQL отзывов, ON CONFLICT upsert, расчет среднего рейтинга и распределения звезд.
* [**`repository/certificate`**](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/certificate/README.md) — Таблица `certificates`, уникальные коды `EDL-YYYY-XXXXXXXX`, выдача и поиск.
* [**`repository/notification`**](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/notification/README.md) — Таблица `notifications`, частичный индекс непрочитанных, выборка ленты.
* [**`repository/models`**](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/models/README.md) — DTO БД таблиц и двусторонние конвертеры в домен.

### Документация подпапок обработчиков (`internal/interfaces/handlers/*`):
* [**`handlers/auth`**](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/auth/README.md) — Регистрация, вход, профиль, админ-панель.
* [**`handlers/category`**](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/category/README.md) — Публичный каталог категорий и создание новых.
* [**`handlers/courses`**](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/courses/README.md) — Каталог, создание, управление курсом.
* [**`handlers/section`**](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/section/README.md) — CRUD секций и сортировка уроков.
* [**`handlers/lesson`**](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/lesson/README.md) — CRUD уроков, статусы, контекст навигации и силлабус.
* [**`handlers/enrollment`**](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/enrollment/README.md) — Самозапись, отчисление, списки учащихся.
* [**`handlers/progress`**](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/progress/README.md) — Старт, таймлайн, завершение уроков, автосохранение драфтов и активные попытки.
* [**`handlers/quiz`**](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/quiz/README.md) — Прохождение тестов и ручной грейдинг.
* [**`handlers/analytics`**](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/analytics/README.md) — Дашборд успеваемости, срез drilldown и глобальная очередь проверки заданий преподавателя.
* [**`handlers/review`**](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/review/README.md) — REST API отзывов курсов и оценок.
* [**`handlers/certificate`**](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/certificate/README.md) — Получение сертификата студентом и публичная верификация.
* [**`handlers/notification`**](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/notification/README.md) — REST API ленты уведомлений и отметок о прочтении.
* [**`handlers/upload`**](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/upload/README.md) — Multipart-загрузка файлов и StaticFileServer с поддержкой Range-запросов и inline PDF.
* [**`handlers/converter`**](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/converter/README.md) — Мапперы Domain <-> HTTP DTO.

---

## 🔄 Сквозные сценарии (End-to-End Data Flows)

### 1. Клиентская аутентификация и авто-ротация токенов (Login & Token Refresh)
```text
Пользователь (UI Form: /login)
  ➔ [src/app/(auth)/login/page.tsx]     : Вызывает api.post('/auth/login', { email, password })
  ➔ [src/lib/api.ts]                    : Передает запрос на бэкенд
  ➔ [interfaces/handlers/auth/login.go] : Декодирует dto.LoginRequest
  ➔ [service/auth/login.go]             : Находит пользователя, сверяет bcrypt хэш пароля
  ➔ [infrastructure/jwt]               : Генерирует JWT access (15м) и refresh (7д)
  ➔ [src/store/useAuth.ts:login]        : Сохраняет токены в localStorage, isAuthenticated = true
  ➔ [src/store/useAuth.ts:fetchUser]    : Выполняет GET /auth/me, определяет роль (student/teacher)
...
Срок действия access-токена истек (401 Unauthorized):
  ➔ [src/lib/api.ts:interceptor]        : Перехватывает 401, видит !_retry, читает refresh_token
  ➔ [POST /auth/refresh]                : Бэкенд валидирует SHA-256 хэш сессии, выдает новую пару
  ➔ [src/lib/api.ts]                    : Обновляет localStorage, повторяет упавший запрос
  ➔ (Если refresh невалиден)            : Удаляет токены из localStorage, redirect на /login
```

### 2. Создание курса и учебного плана преподавателем (Curriculum Builder & Puck Editor)
```text
Преподаватель (UI: /teacher/courses/new)
  ➔ [src/app/teacher/courses/new]       : Валидирует поля, транслитерирует slug, POST /courses
  ➔ [handlers/courses/create.go]        : Валидирует dto.CreateCourseRequest, статус 'draft'
  ➔ [service/course/create.go]          : Сохраняет курс, делает автора 'creator' в БД
  ➔ [src/app/teacher/.../curriculum]    : Открывает конструктор учебного плана
     ├─ ModalCreateModule               : POST /sections (добавление модулей в курс)
     ├─ ModalCreateLesson               : POST /lessons (добавление лекций/квизов)
     ├─ Drag/Buttons reorder            : PUT /courses/{id}/reorder-sections ({ item_ids })
     └─ Кнопка «Открыть редактор»       : Переход на /teacher/lessons/{lessonId}/edit
  ➔ [src/app/teacher/lessons/[id]/edit] : Загружает Puck Editor с CustomOutline
  ➔ Редактирование блоков               : Header, RichText Word, QuizSingle/Multi, Blanks, Essay
  ➔ Кнопка «Опубликовать» (Puck)        : PATCH /lessons/{id} с контентом JSON в формате Content-as-Data
  ➔ Кнопка «Опубликовать курс»          : POST /courses/{id}/publish -> статус 'published'
```

### 3. Интерактивное прохождение урока студентом с автогрейдингом
```text
Студент (UI: /dashboard/courses/{id})
  ➔ Клик по доступному уроку            : Переход на /lessons/{id} (Zen Study Player)
  ➔ [src/app/lessons/[id]/page.tsx]     : Вызывает POST /lessons/{id}/start (статус in_progress)
  ➔ [GET /lessons/{id}]                 : Загружает JSON контента урока
  ➔ [PuckLessonViewer.tsx]              : Динамически рендерит 12 интерактивных блоков
     ├─ QuizSingle / QuizMulti          : Локальный учет выбранных ответов студента
     ├─ QuizMatch / Sequence            : Сопоставление пар и переупорядочивание
     ├─ QuizDropdown / InputBlank       : Парсер parseSmartDropdownTemplate и проверка пропусков
     └─ QuizEssay                       : Ввод развернутого текстового ответа (эссе)
  ➔ Кнопка «Завершить урок»             : PuckLessonViewer или Stepper передает answers, essays, attempt_id, is_abandoned
  ➔ [src/app/lessons/[id]/page.tsx]     : POST /lessons/{id}/complete
  ➔ [handlers/progress/completeLesson]  : Валидирует параметры, передает CompleteLessonInput в сервис
  ➔ [service/progress/completeLesson]   : Anti-Cheat Guard (клиентский балл для тестов игнорируется), серверный пересчет, Best Score Preservation и пересчет курса
  ➔ [Если есть эссе]                    : Сохраняет в quiz_answers (is_correct=NULL, NeedsGrading=true)
  ➔ Преподаватель (/teacher/grading)    : Видит работу в очереди, ставит балл через ModalGradeHW
```

### 4. Тестирование: сдача попытки на бэкенде (`POST /api/v1/quizzes/attempts/{id}/submit`)
```text
Student POST /api/v1/quizzes/attempts/{id}/submit
  ➔ [service/quiz/submitAttempt]       : Захватывает advisory lock для предотвращения гонок
  ➔ [repository/quiz]                  : Загружает вопросы и эталонные ответы
  ➔ [domain/quiz.go]                   : Автопроверка single/multi choice, расчет первичных баллов
  ➔ [service/quiz]                     : Если есть open_text вопросы -> ставит NeedsGrading = true
  ➔ [repository/quiz/attempt.go]       : Сохраняет попытку с баллами и статусом Passed
  ➔ [service/progress]                 : Если сдан успешно -> обновляет прогресс урока (Best Score)
```

### 5. Pre-flight сводка и старт попытки тестирования (`GET /lessons/{id}/attempts/summary`, `POST /attempts/start`)
```text
Student GET /api/v1/lessons/{id}/attempts/summary
  ➔ [handlers/progress/getAttemptsSummary] : Запрос сводки перед тестом
  ➔ [service/progress/attempts]            : Агрегирует попытки, лимит, лучший балл и историю
Student POST /api/v1/lessons/{id}/attempts/start
  ➔ [handlers/progress/startAttempt]       : Валидирует доступ и параметры
  ➔ [service/progress/attempts]            : Проверяет лимит max_attempts (403 при исчерпании) и создает запись attempt
```

### 6. Завершение тестирования с валидацией таймера и режимом обратной связи (`POST /lessons/{id}/complete`)
```text
Student POST /api/v1/lessons/{id}/complete
  ➔ [handlers/progress/completeLesson]     : Прием ответов и параметров завершения
  ➔ [service/progress/completeLesson]      : Валидация дедлайна time_limit_minutes (Grace Period 15s) -> timed_out
                                             Расчет баллов, учет passing_score_percent и Best Score Preservation
                                             В режиме exam_blind: скрытие correct_answer и feedback (results: nil)
  ➔ [repository/progress]                  : Сохранение прогресса в транзакции WithTX
```


---

## 🤖 Руководство по навигации для ИИ-агентов (AI Agent Protocols)

При получении задачи от пользователя следуйте строгому алгоритму локализации:

### Фронтенд задачи (Frontend Layer):
1. **Если нужно изменить UI страницы или роутинг:**
   - Найдите маршрут в [`frontend/src/app/`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/README.md#L45-L60) (`(auth)`, `courses`, `dashboard`, `lessons`, `teacher`).
   - Если требуется изменить общий лейаут или навигацию — смотрите [`frontend/src/components/layout/`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/layout/).
2. **Если нужно добавить или изменить интерактивный блок урока:**
   - Спецификация и типы блока: [`frontend/src/lib/puck-config.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/puck-config.tsx).
   - Интерактивный плеер студента: [`frontend/src/components/player/PuckLessonViewer.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/player/PuckLessonViewer.tsx).
   - Дерево блоков аутлайна редактора: [`frontend/src/components/editor/CustomOutline.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/editor/CustomOutline.tsx).
3. **Если нужно изменить состояние сессии, пользователя или роли:**
   - Обновите интерфейс и методы в [`frontend/src/store/useAuth.ts`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/store/useAuth.ts).
4. **Если нужно изменить логику сетевых запросов или перехватчиков:**
   - Настройте Axios-интерцепторы в [`frontend/src/lib/api.ts`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/lib/api.ts).
5. **Если нужно добавить новый модальный диалог в кабинет автора:**
   - Добавьте компонент в [`frontend/src/components/teacher/`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/components/teacher/).
   - Подключите модалку в [`frontend/src/app/teacher/courses/[id]/curriculum/page.tsx`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/frontend/src/app/teacher/courses/%5Bid%5D/curriculum/page.tsx).

### Бэкенд задачи (Backend Go Layer):
1. **Если нужно изменить API контракт / ответ:**
   - Измените DTO в [`internal/dto/`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/dto/README.md).
   - Обновите конвертер в [`internal/interfaces/handlers/converter/`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/converter/README.md).
2. **Если нужно добавить новый эндпоинт:**
   - Добавьте метод хэндлера в соответствующую папку [`internal/interfaces/handlers/`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/handlers/README.md).
   - Зарегистрируйте URL-путь в [`internal/app/di.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/app/di.go#L405-L538).
3. **Если нужно изменить бизнес-логику / расчеты / доступы:**
   - Найдите сервис в [`internal/service/`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/README.md).
   - Для проверок прав курса смотрите [`internal/service/access/`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/service/access/README.md).
4. **Если нужно изменить SQL-запрос или структуру таблицы:**
   - Создайте новую миграцию в [`migrators/`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/migrators/README.md).
   - Измените SQL и `Scan` в соответствующем репозитории [`internal/repository/`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/README.md).
   - Обновите структуру в [`internal/repository/models/`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/repository/models/README.md).
5. **Если возникла новая ошибка:**
   - Зарегистрируйте sentinel-ошибку в [`pkg/errors/errorsAPP.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/pkg/errors/errorsAPP.go).
   - Добавьте маппинг на HTTP статус код в [`internal/interfaces/response/handler.go`](file:///c:/Users/gogol/OneDrive/Desktop/EdTech/internal/interfaces/response/handler.go#L24-L118).
