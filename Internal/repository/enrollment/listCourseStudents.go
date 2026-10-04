package enrollment

import (
	"context"
	"edtech/internal/domain"
	"edtech/internal/infrastructure/txmanager"
	repomodels "edtech/internal/repository/models"
	repoconverter "edtech/internal/repository/models/converter"
	"fmt"
)

func (r *repository) ListCourseStudents(ctx context.Context, courseID int64, limit, offset int64) ([]domain.User, int, error) {
	const op = "repository.enrollment.listcoursestudents"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

	// Count query
	countQuery := `
		SELECT COUNT(*) 
		FROM users_courses uc
		JOIN users u ON uc.user_id = u.id
		WHERE uc.course_id = $1 AND uc.role = 'student' AND u.deleted_at IS NULL
	`
	var total int
	err := q.QueryRow(ctx, countQuery, courseID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}

	// Data query
	query := `
		SELECT 
			u.id, u.email, u.password_hash, u.username, u.first_name, u.last_name, 
			u.avatar_url, u.bio, u.role, u.email_verified, u.is_active, u.is_banned, 
			u.last_login_at, u.created_at, u.updated_at, u.deleted_at 
		FROM users_courses uc
		JOIN users u ON uc.user_id = u.id
		WHERE uc.course_id = $1 AND uc.role = 'student' AND u.deleted_at IS NULL
		ORDER BY uc.enrolled_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := q.Query(ctx, query, courseID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var students []domain.User
	for rows.Next() {
		var user repomodels.User
		err := rows.Scan(
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
		)
		if err != nil {
			return nil, 0, fmt.Errorf("%s: scan: %w", op, err)
		}
		students = append(students, *repoconverter.UserToDomain(&user))
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("%s: rows err: %w", op, err)
	}

	return students, total, nil
}
