# 🛠 [BE-006] Аудит безопасности и сброс сессий при смене пароля

> **Статус:** Completed  
> **Приоритет:** Medium (P1)  
> **Связанные задачи:** FE-006, QA-006  
> **Целевой модуль:** `internal/interfaces/handlers/auth/`, `internal/service/auth/`  
> **Документация модуля:** [internal/service/auth/FUNCTIONAL_SPEC.md](../../internal/service/auth/FUNCTIONAL_SPEC.md)

## 🎯 Цель задачи
Обеспечить максимальную надежность процедуры смены пароля: валидация сложности нового пароля, обязательная проверка совпадения старого пароля через bcrypt и отзыв всех остальных активных refresh-сессий пользователя в базе при успешной смене, чтобы защитить аккаунт от компрометации.

## 🔍 Текущее состояние кода
- В `internal/interfaces/handlers/auth/password.go` эндпоинт `POST /api/v1/auth/change-password` вызывает `h.authService.ChangePassword(ctx, userID, req.OldPassword, req.NewPassword)`.
- В `internal/service/auth/update.go` или `password.go`: пароль хешируется и обновляется.
- Требуется аудит: убедиться, что при смене пароля удаляются все refresh-токены пользователя из `refresh_tokens`, кроме текущей сессии (или полная инвалидация всех сессий с требованием повторной авторизации).

## 📝 Технические требования к реализации
1. **API / Endpoint:**
   - `POST /api/v1/auth/change-password`
   - Headers: `Authorization: Bearer <token>`
   - Request DTO:
     ```json
     {
       "old_password": "CurrentPassword123!",
       "new_password": "NewStrongPassword456!"
     }
     ```
   - Валидация: `new_password` — минимум 8 символов, не равен `old_password`.
   - Response DTO (200 OK):
     ```json
     { "message": "password successfully changed" }
     ```

2. **Бизнес-логика (Service Layer):**
   - Получить пользователя по `userID`.
   - Сверить `hasher.Compare(user.PasswordHash, oldPassword)`. Если не совпадает — `ErrInvalidCredentials` (401/400).
   - Захешировать `newPassword` через bcrypt (cost 10+).
   - В транзакции:
     1. Обновить `password_hash` и `updated_at`.
     2. Отозвать все сессии: `refreshRepo.DeleteAllByUserID(ctx, tx, userID)`.

## ✅ Критерии приёмки (Definition of Done)
- [x] Смена пароля без правильного старого пароля отклоняется.
- [x] Новый пароль не может совпадать со старым.
- [x] Все старые сессии (refresh-токены) аннулируются, предотвращая доступ злоумышленника.
