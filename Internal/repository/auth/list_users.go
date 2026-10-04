package auth

import (
	"context"
	"fmt"
	"strings"

	"edtech/internal/domain"
	"edtech/internal/infrastructure/txmanager"
	repomodels "edtech/internal/repository/models"
	repoconverter "edtech/internal/repository/models/converter"
	errorsAPP "edtech/pkg/errors"
)

func (r *repository) ListUsers(ctx context.Context, filter domain.UserFilter, pagination domain.Pagination) ([]domain.User, int64, error) {
	const op = "repository.auth.ListUsers"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

	pagination.Sanitize()

	var whereClauses []string
	var args []any
	argIdx := 1

	whereClauses = append(whereClauses, "deleted_at IS NULL")

	if filter.Search != nil && strings.TrimSpace(*filter.Search) != "" {
		searchTerm := "%" + strings.TrimSpace(*filter.Search) + "%"
		whereClauses = append(whereClauses, fmt.Sprintf("(email ILIKE $%d OR username ILIKE $%d OR COALESCE(first_name, '') ILIKE $%d OR COALESCE(last_name, '') ILIKE $%d)", argIdx, argIdx, argIdx, argIdx))
		args = append(args, searchTerm)
		argIdx++
	}

	if filter.Role != nil && *filter.Role != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("role = $%d", argIdx))
		args = append(args, string(*filter.Role))
		argIdx++
	}

	if filter.IsBanned != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("is_banned = $%d", argIdx))
		args = append(args, *filter.IsBanned)
		argIdx++
	}

	whereSQL := strings.Join(whereClauses, " AND ")

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM users WHERE %s", whereSQL)
	var total int64
	if err := q.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("%s: count error: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	if total == 0 {
		return []domain.User{}, 0, nil
	}

	limitIdx := argIdx
	offsetIdx := argIdx + 1
	dataArgs := make([]interface{}, len(args), len(args)+2)
	copy(dataArgs, args)
	dataArgs = append(dataArgs, pagination.PageSize, pagination.Offset())

	dataQuery := fmt.Sprintf(`
		SELECT 
			id, 
			email, 
			password_hash, 
			username, 
			first_name, 
			last_name, 
			avatar_url, 
			bio, 
			role, 
			email_verified, 
			is_active, 
			is_banned, 
			last_login_at, 
			created_at, 
			updated_at, 
			deleted_at
		FROM users
		WHERE %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereSQL, limitIdx, offsetIdx)

	rows, err := q.Query(ctx, dataQuery, dataArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: query error: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}
	defer rows.Close()

	users := make([]domain.User, 0, pagination.PageSize)
	for rows.Next() {
		var user repomodels.User
		if err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.PasswordHash,
			&user.Username,
			&user.FirstName,
			&user.LastName,
			&user.AvatarURL,
			&user.Bio,
			&user.Role,
			&user.EmailVerified,
			&user.IsActive,
			&user.IsBanned,
			&user.LastLoginAt,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.DeletedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("%s: scan error: %w", op, err)
		}

		domainUser := repoconverter.UserToDomain(&user)
		if domainUser != nil {
			users = append(users, *domainUser)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("%s: rows error: %w: %w", op, errorsAPP.ErrInternalDB, err)
	}

	return users, total, nil
}
